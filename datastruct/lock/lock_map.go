package lock

import (
	"sort"
	"sync"
)

type Locks struct {
	locks []*sync.RWMutex
}

func MakeLocks(locksSize int) *Locks {
	locks := make([]*sync.RWMutex, locksSize)
	for i := range locks {
		locks[i] = &sync.RWMutex{}
	}
	return &Locks{
		locks: locks,
	}
}

// 哈希算法选择FNV算法:
const prime32 = uint32(16777619)

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (l *Locks) shard(hashCode uint32) uint32 {
	if l == nil {
		panic("shard a nil Locks")
	}
	return uint32(len(l.locks)-1) & hashCode
}

func (l *Locks) Lock(key string) {
	index := l.shard(fnv32(key))
	l.locks[index].Lock()
}

func (l *Locks) UnLock(key string) {
	index := l.shard(fnv32(key))
	l.locks[index].Unlock()
}

/*
在锁定多个key时需要注意，若协程A持有键a的锁试图获得键b的锁，此时协程B持有键b的锁试图获得键a的锁则会形成死锁。

解决方法是所有协程都按照相同顺序加锁，若两个协程都想获得键a和键b的锁，那么必须先获取键a的锁后获取键b的锁，
这样就可以避免循环等待。
*/
func (l *Locks) toLockIndices(keys []string) []uint32 {
	// 去重
	m := make(map[uint32]bool)
	for _, val := range keys {
		index := l.shard(fnv32(val))
		m[index] = true
	}
	// 收集所有需要上锁的区域，并排序（也就是按照一定的顺序去上锁
	indices := make([]uint32, 0, len(m))
	for i := range m {
		indices = append(indices, i)
	}

	sort.Slice(indices, func(i, j int) bool {
		return indices[i] > indices[j]
	})

	return indices
}

// RWLocks 允许 read 和 write 中有一样的key
func (l *Locks) RWLocks(writeKeys, readKeys []string) {
	keys := append(writeKeys, readKeys...)
	indices := l.toLockIndices(keys)
	writeIndices := l.toLockIndices(writeKeys)
	write := make(map[uint32]bool)
	for i := range writeIndices {
		write[writeIndices[i]] = true
	}
	// 读写可能有相同的key，这种情况需要上写锁
	for i := range indices {
		_, w := write[indices[i]]
		mu := l.locks[indices[i]]
		if w {
			mu.Lock()
		} else {
			mu.RLock()
		}
	}
}

func (l *Locks) RWUnLocks(writeKeys, readKeys []string) {
	keys := append(writeKeys, readKeys...)
	indices := l.toLockIndices(keys)
	writeIndices := l.toLockIndices(writeKeys)
	write := make(map[uint32]bool)
	for i := range writeIndices {
		write[writeIndices[i]] = true
	}
	for i := range indices {
		_, w := write[indices[i]]
		mu := l.locks[indices[i]]
		if w {
			mu.Unlock()
		} else {
			mu.RUnlock()
		}
	}
}
