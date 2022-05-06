package network

import (
	"errors"
	"qymkv/database"
)

const (
	PutReply    = "put k/v success"
	UpdateReply = "update k/v success"
)

var (
	NULLDBERROR = errors.New("operate a nil database")
)

type StringOp struct {
	db *database.DB
}

func (s *StringOp) Set(key string, val interface{}) (*Reply, error) {
	if s.db == nil {
		return nil, NULLDBERROR
	}
	put := s.db.Set(key, val)
	reply := &Reply{Success: true}
	if put {
		reply.Msg = PutReply
	} else {
		reply.Msg = UpdateReply
	}
	return reply, nil
}

func (s *StringOp) Get(key string) (*Reply, error) {
	if s.db == nil {
		return nil, NULLDBERROR
	}
	val, exist := s.db.Get(key)
	reply := &Reply{Success: false}
	if exist {
		reply.Success = true
		reply.Msg = val
	}
	return reply, nil
}
