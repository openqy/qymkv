package database

func (d *DB) Set(key string, val interface{}) bool {
	return d.data.Put(key, val)
}
