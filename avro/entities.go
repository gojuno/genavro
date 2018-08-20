package avro

// Protocol reflects limited to types avro protocol schema.
type Protocol struct {
	Namespace string   `json:"namespace"`
	Protocol  string   `json:"protocol"`
	Types     []Record `json:"types"`
	Doc       string   `json:"doc,omitempty"`
}

// Record reflects avro record type schema.
type Record struct {
	Type      string  `json:"type"`
	Name      string  `json:"name"`
	Namespace string  `json:"namespace,omitempty"`
	Doc       string  `json:"doc,omitempty"`
	Fields    []Field `json:"fields"`
}

// Field reflects field in avro record type.
type Field struct {
	Name    string      `json:"name"`
	Doc     string      `json:"doc,omitempty"`
	Type    interface{} `json:"type"`
	Default interface{} `json:"default,omitempty"`
}

// Array is a array type of the field.
type Array struct {
	Type  string      `json:"type"`
	Items interface{} `json:"items"`
}

// Map is a map type of the field.
type Map struct {
	Type   string `json:"type"`
	Values string `json:"values"`
}

// Union is a union type of the field. used for nullable fields.
type Union [2]string
