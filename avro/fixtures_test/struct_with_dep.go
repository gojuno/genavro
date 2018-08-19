package fixtures_test

type Dep struct {
	Int int `json:"int"`
}

type StructV1 struct {
	Dep Dep `json:"dep"`
}
