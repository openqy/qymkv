package dict

type Dict interface {
	Len() int
	Get(key string) (interface{}, bool)
	Put(key string, val interface{}) bool
	Remove(key string) bool
}
