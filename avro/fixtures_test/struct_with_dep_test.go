package fixtures_test

type Dep1 struct {
	Str string `json:"str"`
}
type Dep2 struct {
	Str string `json:"str"`
}
type Dep3 struct {
	Str string `json:"str"`
}
type Dep4 struct {
	Str string `json:"str"`
}

type Dep struct {
	Int       int             `json:"int"`
	Dep1      Dep1            `json:"dep1"`
	Dep2Opt   *Dep2           `json:"dep2_opt,omitempty"`
	Dep3Array []Dep3          `json:"dep3_array,omitempty"`
	Dep4Map   map[string]Dep4 `json:"dep4_map,omitempty"`
}

type Optional struct {
	Int int `json:"int"`
}

const minorVersionStructV1 = "1"

type StructV1 struct {
	Dep      Dep       `json:"dep"`
	Optional *Optional `json:"optional"`
}
