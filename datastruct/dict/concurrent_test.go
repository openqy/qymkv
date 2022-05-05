package dict

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNormalFunc(t *testing.T) {
	fmt.Println("----------------computeCapacity--------------")
	cap := computeCapacity(1 << 16)
	fmt.Println(cap)
	//cap = computeCapacity(15)
	//fmt.Println(cap)
	//cap = computeCapacity(16)
	//fmt.Println(cap)
	//cap = computeCapacity(20)
	//fmt.Println(cap)
	//cap = computeCapacity(32)
	//fmt.Println(cap)
	//cap = computeCapacity(33)
	//fmt.Println(cap)
	//fmt.Println(math.MaxInt32)
	//cap = computeCapacity(math.MaxInt32+2)
	//fmt.Println(cap)
	//cap = computeCapacity(-1)
	//fmt.Println(cap)

	m := MakeConcurrent(10)
	fmt.Println(len(m.table) == computeCapacity(10))
	fmt.Println(m.count == 0)
	fmt.Println("----------------shard--------------")
	index := m.shard(30)
	fmt.Println(index == 14)
	if index != 14 {
		t.Errorf("index:%v unequal, we need %v", index, 14)
	}
}

func TestCommands(t *testing.T) {
	m := MakeConcurrent(0)
	put := m.Put("lll", "mmm")
	if !put {
		t.Error("Put error : ", put)
	}
	val, ok := m.Get("lll")
	if !ok || val != "mmm" {
		t.Error("Get error, expect: mmm , get ", val)
	}
	if m.count != 1 {
		t.Error("Get error, expect: 1 , get ", m.count)
	}

	remove := m.Remove("lll")
	if !remove {
		t.Errorf("Expect remove , but not")
	}
	if m.count != 0 {
		t.Error("Expect 0, get :", m.count)
	}
}

func TestConcurrentPut(t *testing.T) {
	d := MakeConcurrent(0)
	count := 1000000
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			// insert
			key := "k" + strconv.Itoa(i)
			ret := d.Put(key, i)
			if ret != true { // insert 1
				t.Error("put test failed: expected result true, actual: " + strconv.FormatBool(ret) + ", key: " + key)
			}
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				if intVal != i {
					t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal) + ", key: " + key)
				}
			} else {
				_, ok := d.Get(key)
				t.Error("put test failed: expected true, actual: false, key: " + key + ", retry: " + strconv.FormatBool(ok))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(d.count)
}

func TestConcurrentRemove(t *testing.T) {
	d := MakeConcurrent(0)
	totalCount := 100
	// remove head node
	for i := 0; i < totalCount; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	if d.Len() != totalCount {
		t.Error("put test failed: expected len is 100, actual: " + strconv.Itoa(d.Len()))
	}
	for i := 0; i < totalCount; i++ {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != true {
			t.Error("remove test failed: expected result 1, actual: " + strconv.FormatBool(ret) + ", key:" + key)
		}
		if d.Len() != totalCount-i-1 {
			t.Error("put test failed: expected len is 99, actual: " + strconv.Itoa(d.Len()))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != false {
			t.Error("remove test failed: expected result 0 actual: " + strconv.FormatBool(ret))
		}
		if d.Len() != totalCount-i-1 {
			t.Error("put test failed: expected len is 99, actual: " + strconv.Itoa(d.Len()))
		}
	}

	// remove tail node
	d = MakeConcurrent(0)
	for i := 0; i < 100; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	for i := 9; i >= 0; i-- {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != true {
			t.Error("remove test failed: expected result 1, actual: " + strconv.FormatBool(ret))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != false {
			t.Error("remove test failed: expected result 0 actual: " + strconv.FormatBool(ret))
		}
	}

	// remove middle node
	d = MakeConcurrent(0)
	d.Put("head", 0)
	for i := 0; i < 10; i++ {
		// insert
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	d.Put("tail", 0)
	for i := 9; i >= 0; i-- {
		key := "k" + strconv.Itoa(i)

		val, ok := d.Get(key)
		if ok {
			intVal, _ := val.(int)
			if intVal != i {
				t.Error("put test failed: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
			}
		} else {
			t.Error("put test failed: expected true, actual: false")
		}

		ret := d.Remove(key)
		if ret != true {
			t.Error("remove test failed: expected result 1, actual: " + strconv.FormatBool(ret))
		}
		_, ok = d.Get(key)
		if ok {
			t.Error("remove test failed: expected true, actual false")
		}
		ret = d.Remove(key)
		if ret != false {
			t.Error("remove test failed: expected result 0 actual: " + strconv.FormatBool(ret))
		}
	}
}

func BenchmarkHash(b *testing.B) {
	m := MakeConcurrent(0)
	for i := 0; i < b.N; i++ {
		m.Put(strconv.Itoa(i), i)
	}
	fmt.Println(m.count)
}

func TestPutTime(t *testing.T) {
	m := MakeConcurrent(0)
	var wg sync.WaitGroup
	wg.Add(10000000)
	ti := time.Now()
	for i := 0; i < 10000000; i++ {
		go func(i int) {
			m.Put(strconv.Itoa(i), i)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("my:", time.Now().Sub(ti))

	var m2 sync.Map
	var wg2 sync.WaitGroup
	wg2.Add(10000000)
	ti2 := time.Now()
	for i := 0; i < 10000000; i++ {
		go func(i int) {
			m2.Store(strconv.Itoa(i), i)
			wg2.Done()
		}(i)
	}
	wg2.Wait()
	fmt.Println(time.Now().Sub(ti2))
}

func TestGetTime(t *testing.T) {
	m := MakeConcurrent(0)
	for i := 0; i < 10000000; i++ {
		m.Put(strconv.Itoa(i), i)
	}

	var m2 sync.Map
	for i := 0; i < 10000000; i++ {
		m2.Store(strconv.Itoa(i), i)
	}
	//
	for i := 0; i < 5; i++ {
		ti := time.Now()
		for i := 0; i < 10000000; i++ {
			val, _ := m.Get(strconv.Itoa(i))
			if val.(int) != i {
				t.Error("concurrent map error")
			}
		}
		fmt.Println("my:", time.Now().Sub(ti))

		ti2 := time.Now()
		for i := 0; i < 10000000; i++ {
			val, _ := m2.Load(strconv.Itoa(i))
			if val.(int) != i {
				t.Error("sync map error")
			}
		}
		fmt.Println("golang:", time.Now().Sub(ti2))
	}
}

func TestReadWriteTime(t *testing.T) {
	m := MakeConcurrent(0)
	ti := time.Now()
	for i := 0; i < 1000000; i++ {
		m.Put(strconv.Itoa(i), i)
		m.Get(strconv.Itoa(i))
	}
	fmt.Println(time.Now().Sub(ti))

	var m2 sync.Map
	ti2 := time.Now()
	for i := 0; i < 1000000; i++ {
		m2.Store(strconv.Itoa(i), i)
		m2.Load(strconv.Itoa(i))
	}
	fmt.Println(time.Now().Sub(ti2))
}

func TestTestReadWriteTime(t *testing.T) {
	m := MakeConcurrent(0)
	var m2 sync.Map
	ti := time.Now()
	for i := 0; i < 100; i++ {
		m2.Store(strconv.Itoa(i), i)
		for j := 0; j < 1000; j++ {
			m2.Load(strconv.Itoa(i))
		}
	}
	fmt.Println(time.Now().Sub(ti))

	ti2 := time.Now()
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), i)
		for j := 0; j < 1000; j++ {
			val, _ := m.Get(strconv.Itoa(i))
			if val != i {
				fmt.Println("err")
			}
		}
	}
	fmt.Println(time.Now().Sub(ti2))
}
