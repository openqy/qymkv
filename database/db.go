package database

import (
	"qymkv/datastruct/dict"
	"qymkv/datastruct/lock"
)

const (
	dataSize  = 1 << 16
	lockCount = 1024
)

type DB struct {
	data dict.Dict

	locker *lock.Locks
}

func MakeDB() *DB {
	db := &DB{
		dict.MakeConcurrent(dataSize),
		lock.MakeLocks(lockCount),
	}
	return db
}
