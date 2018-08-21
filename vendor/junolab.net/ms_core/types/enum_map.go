package types

type enumMap struct {
	ab map[interface{}]interface{}
	ba map[interface{}]interface{}
}

func NewEnumMap() *enumMap {
	return &enumMap{make(map[interface{}]interface{}), make(map[interface{}]interface{})}
}

func (m *enumMap) Put(name, value interface{}) *enumMap {
	m.ab[name] = value
	m.ba[value] = name
	return m
}

func (m *enumMap) Value(name interface{}) (value interface{}, exists bool) {
	value, exists = m.ab[name]
	return
}

func (m *enumMap) Name(value interface{}) (name interface{}, exists bool) {
	name, exists = m.ba[value]
	return
}
