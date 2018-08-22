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
		{
			Type: newUnion("string"),
			Name: "app_id",
		},
		{
			Type: newUnion("string"),
			Name: "app_version",
		},
	},
}

type dep struct {
	record Record
	deps   []string
}

// Generate converts structs from parsed files to avro protocol.
// Every struct ending with `V\d` will be parsed as Separate top level event and will be generated to separate protocol.
// e.g. MetricsV1 will be generated as separate protocol and Metrics will be not.
func Generate(sources map[string]astparser.ParsedFile, namespace string) map[string]Protocol {

	r := regexp.MustCompile(".*V\\d+$")
	deps := map[string]dep{}
	versions := map[string]string{}

	for _, parsedFile := range sources {
		// build dependencies map
		for _, s := range parsedFile.Structs {
			// skip events
			if r.Match([]byte(s.Name)) {
				continue
			}
			deps[s.Name] = parseDep(s)
		}

		// build minor version map
		for _, c := range parsedFile.Constants {
			if strings.HasPrefix(c.Name, "minorVersion") {
				ss := strings.Split(c.Name, "minorVersion")
				if len(ss) == 2 {
					versions[ss[1]] = c.Value
				}
			}
		}
	}

	result := map[string]Protocol{}
	for _, parsedFile := range sources {
		for _, s := range parsedFile.Structs {
			// pass only events ends on
			if !r.Match([]byte(s.Name)) {
				continue
			}
			p := avroProtocol(s, deps, namespace, versions[s.Name])
			result[s.Name] = p
		}

	}

	return result
}

func avroProtocol(s astparser.StructDef, deps map[string]dep, namespace, minorVersion string) Protocol {
	base := avroBaseV1Type(s.Name, minorVersion)
	base.Fields = append(base.Fields, Field{
		Name: "payload",
		Type: payloadName(s.Name),
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
	rs.Name = payloadName(rs.Name)

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

func avroBaseV1Type(name, minorVersion string) Record {
	return Record{
		Name: name,
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
				Doc:  fmt.Sprintf("minorVersion=%s", minorVersion),
			},
			{
				Name: "auth",
				Type: newUnion("Auth"),
			},
		},
	}
}

func avroDepName(tpe interface{}) string {
	switch t := tpe.(type) {
	case string:
		return t
	case Map:
		return t.Values
	case Array:
		return avroDepName(t.Items)
	case Union:
		return avroDepName(t[1])
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
			field.Type = newUnion(field.Type)
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
			field.Type = newUnion(field.Type)
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
		return newUnion(avroType(v.InnerType))

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

func newUnion(t interface{}) Union {
	switch v := t.(type) {
	case Union:
		return v
	case Array:
		return Union{"null", t}
	case Map:
		return Union{"null", t}
	default:
		return Union{"null", fmt.Sprint(t)}
	}
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

func payloadName(name string) string {
	return fmt.Sprintf("Payload%s", name)
}
