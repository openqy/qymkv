package lock

import (
	"qymkv/datastruct/dict"
	"sync"
	"testing"
)

func TestMakeLocks(t *testing.T) {
	locks := MakeLocks(10)
	for i := range locks.locks {
		locks.locks[i].Lock()
		locks.locks[i].Unlock()
	}
}

func TestLocks(t *testing.T) {
	locks := MakeLocks(10)
	m := dict.MakeConcurrent(16)
	key := "ll"
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			// get 和 set 之间不是原子的，所以需要这样
			locks.Lock(key)
			val, has := m.Get(key)
			if has {
				m.Put(key, val.(int)+1)
			} else {
				m.Put(key, 1)
			}
			locks.UnLock(key)
			wg.Done()
		}()
	}
	wg.Wait()
	val, b := m.Get(key)
	if !b {
		t.Error("key not exist")
	}
	if val != 20 {
		t.Error("expect 20, but got : ", val)
	}
}

func TestRWLocks(t *testing.T) {
	locks := MakeLocks(10)
	m := dict.MakeConcurrent(16)
	keys := []string{"ll", "mm"}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			locks.RWLocks(keys, keys)
			for _, key := range keys {
				val, has := m.Get(key)
				if has {
					m.Put(key, val.(int)+1)
				} else {
					m.Put(key, 1)
				}
			}
			locks.RWUnLocks(keys, keys)
			wg.Done()
		}()
	}
	wg.Wait()
	for _, key := range keys {
		val, b := m.Get(key)
		if !b {
			t.Error("key not exist")
		}
		if val != 20 {
			t.Error("expect 20, but got : ", val)
		}
	}
}
