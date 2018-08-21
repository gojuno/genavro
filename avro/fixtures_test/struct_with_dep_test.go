package fixtures_test

type Dep struct {
	Int int `json:"int"`
}

type Optional struct {
	Int int `json:"int"`
}

const minorVersionStructV1 = "1"

type StructV1 struct {
	Dep      Dep       `json:"dep"`
	Optional *Optional `json:"optional"`
}
