package database

func (d *DB) Set(key string, val interface{}) bool {
	return d.data.Put(key, val)
}

func (d *DB) Get(key string) (interface{}, bool) {
	val, exist := d.data.Get(key)
	return val, exist
}
