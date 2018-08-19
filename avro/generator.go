package avro

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/mkorolyov/astparser"
)

var avroAuthType = Record{
	Name: "Auth",
	Type: "record",
	Fields: []Field{
		{
			Type: newUnion("string"),
			Name: "session_id",
		},
		{
			Type: newUnion("string"),
			Name: "user_id",
		},
	},
}

var avroBaseV1Type = Record{
	Name: "BaseV1",
	Type: "record",
	Fields: []Field{
		{
			Name: "event_id",
			Type: "string",
		},
		{
			Name: "request_id",
			Type: "string",
		},
		{
			Name: "event_ts",
			Type: "long",
		},
		{
			Name: "type",
			Type: "string",
		},
		{
			Name: "minor_version",
			Type: "string",
		},
		{
			Name: "auth",
			Type: newUnion("Auth"),
		},
	},
}

type dep struct {
	record Record
	deps   []string
}

func Generate(sources map[string][]astparser.StructDef, namespace string) map[string]Protocol {

	r := regexp.MustCompile(".*V\\d+$")
	deps := map[string]dep{}
	for _, structs := range sources {
		for _, s := range structs {
			// skip events
			if r.Match([]byte(s.Name)) {
				continue
			}
			deps[s.Name] = parseDep(s)
		}

	}

	result := map[string]Protocol{}
	for _, structs := range sources {
		for _, s := range structs {
			// pass only events ends on
			if !r.Match([]byte(s.Name)) {
				continue
			}
			result[s.Name] = avroProtocol(s, deps, namespace)
		}

	}

	return result
}

func avroProtocol(s astparser.StructDef, deps map[string]dep, namespace string) Protocol {
	base := avroBaseV1Type
	base.Fields = append(base.Fields, Field{
		Name: "payload",
		Type: s.Name,
	})

	protocol := Protocol{
		Namespace: namespace,
		Protocol:  s.Name,
		Types:     []Record{avroAuthType, base},
	}

	notUniqueDeps := map[int]dep{}
	depIndex := 0
	rs := avroRecord(s, func(tpe interface{}) {
		if d := avroDep(deps, tpe); d != nil {
			notUniqueDeps[depIndex] = *d
			depIndex++
		}
	})

	uniqueDepsIndex := map[string]int{}
	for i, d := range notUniqueDeps {
		uniqueDepsIndex[d.record.Name] = i
	}

	uniqueDeps := make([]Record, 0, len(uniqueDepsIndex))
	for _, index := range uniqueDepsIndex {
		uniqueDeps = append(uniqueDeps, notUniqueDeps[index].record)
	}

	sort.Slice(uniqueDeps, func(i, j int) bool {
		return uniqueDepsIndex[uniqueDeps[i].Name] < uniqueDepsIndex[uniqueDeps[j].Name]
	})

	protocol.Types = append(append(uniqueDeps, rs), protocol.Types...)

	return protocol
}

func avroDepName(tpe interface{}) string {
	switch t := tpe.(type) {
	case string:
		return t
	case Map:
		return t.Values
	case Array:
		return avroDepName(t.Items)
	default:
		return ""
	}
}

func avroDep(deps map[string]dep, tpe interface{}) *dep {
	if dep, ok := deps[avroDepName(tpe)]; ok {
		return &dep
	}
	return nil
}

func avroRecord(s astparser.StructDef, collectDeps func(tpe interface{})) Record {
	fields := make([]Field, 0, len(s.Fields))
	for _, f := range s.Fields {
		field := Field{
			Name: f.JsonName,
			Doc:  strings.Join(f.Comments, ", ")}

		field.Type = avroType(f.FieldType)
		if f.Omitempty {
			field.Default = getDefaultValue(field.Type)
		}

		fields = append(fields, field)

		if collectDeps != nil {
			collectDeps(field.Type)
		}
	}

	return Record{
		Name:   s.Name,
		Type:   "record",
		Doc:    strings.Join(s.Comments, ", "),
		Fields: fields,
	}
}

func parseDep(s astparser.StructDef) dep {
	var deps []string
	fields := make([]Field, 0, len(s.Fields))
	for _, f := range s.Fields {
		field := Field{
			Name: f.JsonName,
			Doc:  strings.Join(f.Comments, ", ")}

		field.Type = avroType(f.FieldType)
		if f.Omitempty {
			field.Default = getDefaultValue(field.Type)
		}

		if !avroIsSimpleType(field.Type) {
			deps = append(deps, fmt.Sprint(field.Type))
		}

		fields = append(fields, field)
	}

	return dep{record: Record{
		Name:   s.Name,
		Type:   "record",
		Doc:    strings.Join(s.Comments, ", "),
		Fields: fields,
	}, deps: deps}
}

func avroType(t astparser.Type) interface{} {
	switch v := t.(type) {
	case astparser.TypeSimple:
		return avroSimpleType(v.Name)
	case astparser.TypePointer:
		return newUnion(fmt.Sprint(avroType(v.InnerType)))

	case astparser.TypeArray:
		return Array{Type: "array", Items: fmt.Sprint(avroType(v.InnerType))}

	case astparser.TypeMap:
		return Map{Type: "map", Values: fmt.Sprint(avroType(v.ValueType))}

	case astparser.TypeCustom:
		switch v.Name {
		// core.ID
		case "ID":
			return "string"
		// timeapi.Time, //timeapi.Duration
		case "Time", "Duration":
			return "long"
		default:
			return v.Name
		}

	default:
		log.Fatalf("unexpected go type %+[1]v: %[1]T", t)
		return nil
	}
}

func getDefaultValue(tpe interface{}) interface{} {
	switch tpe {
	case "int", "long", "float", "double":
		return 0
	case "string", "bytes":
		return ""
	case "boolean":
		return false
	default:
		//log.Fatalf("unsupported avro type %s", tpe)
		return nil
	}
}

func newUnion(t string) Union {
	return Union{"null", t}
}

func avroSimpleType(gotype string) string {
	switch gotype {
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
		return "int"
		// TODO add support timeapi.Time -> timestamp-millis transition
	case "int64", "uint64", "uintptr":
		return "long"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "bool":
		return "boolean"
	case "[]byte":
		return "bytes"
	case "string":
		return "string"
	default:
		log.Fatalf("unsupported go type %s", gotype)
		return ""
	}
}

func avroIsSimpleType(avroType interface{}) bool {
	switch v := avroType.(type) {
	case string:
		switch v {
		case "int", "long", "float", "double", "boolean", "bytes", "string":
			return true
		}
	}

	return false
}
