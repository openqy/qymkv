package database

import "fmt"

func (d *DB) Set(key string, val interface{}) bool {
	switch val.(type) {
	case string:
		return d.data.Put(key, val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		// 把上面的数字类型和布尔类型都转换为 string
		valStr := fmt.Sprintf("%v", val)
		return d.data.Put(key, valStr)
	default:
		return false
	}
}

func (d *DB) Get(key string) (interface{}, bool) {
	val, exist := d.data.Get(key)
	return val, exist
}
