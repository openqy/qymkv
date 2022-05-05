package dict

type Map struct {
	data map[string]interface{}
}

func MakeMap() *Map {
	return &Map{
		data: make(map[string]interface{}),
	}
}

func (m *Map) Len() int {
	return len(m.data)
}

func (m *Map) Put(key string, val interface{}) bool {
	if _, ok := m.data[key]; ok {
		m.data[key] = val
		return true
	}
	m.data[key] = val
	return false
}

func (m *Map) Get(key string) (interface{}, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *Map) Remove(key string) bool {
	if _, ok := m.data[key]; ok {
		delete(m.data, key)
		return true
	}
	return false
}
