package astparser

import "go/ast"

type ParsedFile struct {
	Structs   []StructDef
	Constants []ConstantDef
}

// Type represent parsed type.
type Type interface{}

// ConstantDef describes defined constants
type ConstantDef struct {
	Name  string
	Value string
}

// StructDef describes parsed go struct.
type StructDef struct {
	Name     string
	Fields   []FieldDef
	Comments []string
}

// Tag contains parsed field tags.
type Tag struct {
	JsonName string

	Omitempty bool
}

// FieldDef described parsed go struct field.
type FieldDef struct {
	FieldName string
	FieldType Type
	JsonName  string
	Omitempty bool
	Comments  []string
}

// TypeSimple indicates that type is a primitive golang type like int or string.
type TypeSimple struct {
	Name string
}

// TypeArray indicates that type is golang array or slice.
// Inner type could be any type golang supports.
type TypeArray struct {
	InnerType Type
}

// TypeMap indicates that type is golang map.
// Both keys and values could be any type golang supports.
type TypeMap struct {
	KeyType   Type
	ValueType Type
}

// TypeCustomer indicates that type is a defined struct or type alias.
type TypeCustom struct {
	Name string
	Expr ast.Expr
}

// TypePointer indicates that type is a point with underlying any golang type
type TypePointer struct {
	InnerType Type
}
