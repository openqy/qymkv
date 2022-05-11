package test

import (
	"qymkv/database"
	"qymkv/network"
	"strconv"
	"sync"
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
		t.Errorf("expected true, get %v, msg=%v", get.Success, get.Msg)
	}
	if get.Msg.(string) != "1" {
		t.Errorf("expected 1, get %v", get.Msg)
	}
}

func TestMultiSet(t *testing.T) {
	db := newStringToTest()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			set, err := db.Set(strconv.Itoa(i), i)
			if err != nil {
				t.Error(err)
			}
			if !set.Success {
				t.Errorf("expected true, get %v, msg=%v", set.Success, set.Msg)
			}
		}(i / 500)
	}
	wg.Wait()

	get, err := db.Get("0")
	if err != nil {
		t.Error(err)
	}
	if !get.Success {
		t.Errorf("expected true, get %v, msg=%v", get.Success, get.Msg)
	}
	if get.Msg.(string) != "0" {
		t.Errorf("expected 0, get %v", get.Msg)
	}

	get, err = db.Get("1")
	if err != nil {
		t.Error(err)
	}
	if !get.Success {
		t.Errorf("expected true, get %v, msg=%v", get.Success, get.Msg)
	}
	if get.Msg.(string) != "1" {
		t.Errorf("expected 1, get %v", get.Msg)
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

func BenchmarkSet(b *testing.B) {
	db := newStringToTest()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			set, err := db.Set(strconv.Itoa(i), i)
			if err != nil {
				b.Error(err)
			}
			if !set.Success {
				b.Errorf("expected true, get %v, msg=%v", set.Success, set.Msg)
			}
		}(i)
	}
	wg.Wait()
}
