package test

import (
	"qymkv/database"
	"qymkv/network"
	"testing"
)

func newStringToTest() network.StringOperation {
	db := database.MakeDB()
	op := network.NewStringOp(db)
	return op
}

func TestSet(t *testing.T) {
	db := newStringToTest()
	set, err := db.Set("test", 1)
	if err != nil {
		t.Error(err)
	}
	if !set.Success {
		t.Errorf("expected true, get %v, msg=%v", set.Success, set.Msg)
	}
	get, err := db.Get("test")
	if err != nil {
		t.Error(err)
	}
	if !get.Success {
		t.Errorf("expected true, get %v, msg=%v", set.Success, set.Msg)
	}
	if get.Msg.(int) != 1 {
		t.Errorf("expected 1, get %v", set.Msg)
	}
}

func TestGet(t *testing.T) {
	db := newStringToTest()
	get, err := db.Get("test")
	if err != nil {
		t.Error(err)
	}
	if get.Success {
		t.Errorf("expected false, get %v, msg=%v", get.Success, get.Msg)
	}
}
