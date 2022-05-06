package network

type StringOperation interface {
	Set(key string, val interface{}) (*Reply, error)
}

type DictOperation interface {
}

type ListOperation interface {
}

type SetOperation interface {
}

type KeysOperation interface {
}

type Reply struct {
	Success bool
	Msg     string
}
