package network

import "qymkv/database"

func NewStringOp(db *database.DB) StringOperation {
	stringOp := &StringOp{
		db: db,
	}
	return stringOp
}
