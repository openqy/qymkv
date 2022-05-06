package network

type StringOperation interface {
	Set(key, val string) *Reply
}

type Reply struct {
	Success bool
	Msg     string
}
