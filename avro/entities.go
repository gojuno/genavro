package avro

type Protocol struct {
	Namespace string   `json:"namespace"`
	Protocol  string   `json:"protocol"`
	Types     []Record `json:"types"`
	Doc       string   `json:"doc,omitempty"`
}

type Record struct {
	Type      string  `json:"type"`
	Name      string  `json:"name"`
	Namespace string  `json:"namespace,omitempty"`
	Doc       string  `json:"doc,omitempty"`
	Fields    []Field `json:"fields"`
}

type Field struct {
	Name    string      `json:"name"`
	Doc     string      `json:"doc,omitempty"`
	Type    interface{} `json:"type"`
	Default interface{} `json:"default,omitempty"`
}

type Array struct {
	Type  string      `json:"type"`
	Items interface{} `json:"items"`
}

type Map struct {
	Type   string `json:"type"`
	Values string `json:"values"`
}

type Union [2]string
