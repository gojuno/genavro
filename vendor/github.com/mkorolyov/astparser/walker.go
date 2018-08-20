package astparser

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"github.com/pkg/errors"
)

// TODO parse type comments
// Walker implements go/ast.Visitor to walk through golang
// structs and constants to parse them.
type Walker struct {
	Structs   []StructDef
	Constants []ConstantDef
}

// A Walkers's Visit method is invoked for each node encountered by go/ast.Walk.
// If the result visitor w is not nil, go/ast.Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
func (w *Walker) Visit(node ast.Node) ast.Visitor {
	switch spec := node.(type) {
	case *ast.TypeSpec:
		w.visitStruct(spec)
		return nil
	case *ast.ValueSpec:
		w.visitConstant(spec)
	}

	return w
}

func (w *Walker) visitConstant(astValueSpec *ast.ValueSpec) {
	if len(astValueSpec.Names) < 1 || len(astValueSpec.Values) < 1 {
		return
	}
	name := astValueSpec.Names[0].Name
	var value string
	switch v := astValueSpec.Values[0].(type) {
	case *ast.BasicLit:
		value = removeQuotes(v.Value)
	default:
		return
	}

	w.Constants = append(w.Constants, ConstantDef{
		Name: name, Value: value,
	})
}

func (w *Walker) visitStruct(astTypeSpec *ast.TypeSpec) {
	structName := astTypeSpec.Name.Name

	switch v := astTypeSpec.Type.(type) {
	case *ast.StructType:
		astFields := v.Fields.List

		s := StructDef{
			Name:     structName,
			Comments: parseComments(astTypeSpec.Doc)}

		for _, astField := range astFields {
			field, err := parseField(astField)
			if err != nil {
				log.Fatalf("failed to parse struct %s: %v", structName, err)
			}
			s.Fields = append(s.Fields, field)
		}

		w.Structs = append(w.Structs, s)

	default:
		log.Fatalf("unexpected type for typeSpec: %s, %+v: %T", structName, astTypeSpec, astTypeSpec.Type)
	}

}

func parseField(astField *ast.Field) (FieldDef, error) {
	fieldName, err := parseFieldName(astField.Names)
	if err != nil {
		return FieldDef{}, errors.Wrapf(err, "failed to parse field %+v name", *astField)
	}
	fieldType, err := parseFieldType(astField.Type)
	if err != nil {
		return FieldDef{}, errors.Wrapf(err, "failed to parse field %s type", fieldName)
	}
	tag, err := parseJSONTag(astField.Tag)
	if err != nil {
		return FieldDef{}, errors.Wrapf(err, "failed to parse field %s tags", fieldName)
	}
	field := FieldDef{
		FieldName: fieldName,
		FieldType: fieldType,
		Omitempty: tag.Omitempty,
		JsonName:  tag.JsonName,
		Comments:  parseComments(astField.Doc),
	}
	return field, nil
}

func parseComments(group *ast.CommentGroup) []string {
	if group == nil || len(group.List) == 0 {
		return nil
	}

	var comments []string
	for _, c := range group.List {
		if c == nil {
			continue
		}
		comments = append(comments, strings.TrimSpace(strings.Replace(c.Text, "//", "", -1)))
	}
	return comments
}

func removeQuotes(s string) string {
	if len(s) < 2 {
		log.Fatalf("bad input for removing quotes: %s", s)
	}

	return s[1 : len(s)-1]
}

func parseJSONTag(astTag *ast.BasicLit) (Tag, error) {
	if astTag == nil {
		return Tag{}, nil
	}
	tagString := astTag.Value
	if tagString == "" {
		return Tag{}, nil
	}
	tagString = removeQuotes(tagString) //clean from `json:"place_type,omitempty"` to  json:"place_type,omitempty"
	splittedTags := strings.Split(tagString, " ")
	t := Tag{}
	for _, tagWithName := range splittedTags {
		if tagWithName == "" {
			continue
		}
		v := strings.SplitN(tagWithName, ":", 2)
		if len(v) != 2 {
			return Tag{}, fmt.Errorf("invalid tag %s", tagWithName)
		}
		tagName := strings.Trim(v[0], " ")
		if tagName != "json" {
			continue
		}
		tagValues := strings.Split(strings.TrimSpace(v[1]), ",")
		if len(tagValues) == 1 { // e.g. json:"field_name"
			t.JsonName = removeQuotes(tagValues[0])
		} else { // e.g. handle json:"field_name,omitempty" where
			// tagValues[0] = '"field_name', so we need to cut first char "
			t.JsonName = tagValues[0][1:len(tagValues[0])]
			t.Omitempty = true
		}
		break
	}
	return t, nil
}

func parseFieldType(t ast.Expr) (Type, error) {
	switch v := t.(type) {
	case *ast.Ident:
		if st := simpleType(v.Name); st != nil {
			return st, nil
		}
		return TypeCustom{Name: v.Name}, nil
	case *ast.SelectorExpr:
		return TypeCustom{Name: v.Sel.Name, Expr: t}, nil
	case *ast.ArrayType:
		t, err := parseFieldType(v.Elt)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse array nested type %+v", t)
		}
		return TypeArray{InnerType: t}, nil
	case *ast.StarExpr:
		t, err := parseFieldType(v.X)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse star expr type %+v", t)
		}
		return TypePointer{InnerType: t}, nil
	case *ast.MapType:
		kt, err := parseFieldType(v.Key)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse map key type %+v", v.Key)
		}
		vt, er := parseFieldType(v.Value)
		if er != nil {
			return nil, errors.Wrapf(er, "failed to parse map value type %+v", v.Value)
		}
		return TypeMap{KeyType: kt, ValueType: vt}, nil
	default:
		return nil, fmt.Errorf("unexpected %+[1]v with type %[1]T", t)
	}
}

func parseFieldName(fieldNames []*ast.Ident) (string, error) {
	if len(fieldNames) == 0 {
		return "", fmt.Errorf("anonimuous fields are not supported")
	}
	return fieldNames[0].Name, nil
}

func simpleType(fieldType string) Type {
	switch fieldType {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "bool", "byte":
		return TypeSimple{Name: fieldType}
	default:
		return nil
	}
}
