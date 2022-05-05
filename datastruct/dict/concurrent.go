package dict

import (
	"math"
	"sync"
	"sync/atomic"
)

/*
	godis作者说由于各种原因没有实现渐进式rehash，因此使用了分段锁策略：
	其思想是将key分散到固定数量的shard中，shard是有锁保护的map，当shard进行rehash时会进行rehash，
	当shard进行rehash时会阻塞shard内的读写，但不会影响其他的shard,
	ConcurrentMap 可以保证对单个 key 操作的并发安全性.
*/

type ConcurrentDict struct {
	table []*shard
	count uint32
}

type shard struct {
	mutex sync.RWMutex
	m     map[string]interface{}
}

func computeCapacity(param int) (size int) {
	// 小于16的都返回16
	if param <= 16 {
		return 16
	}
	// 这个算法的大致意思就是返回 最接近 param的 2的n次方那个数
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	} else {
		return n + 1
	}
}

func MakeConcurrent(shardCount int) *ConcurrentDict {
	shardCount = computeCapacity(shardCount)
	table := make([]*shard, shardCount)
	for i := range table {
		table[i] = &shard{
			m: make(map[string]interface{}),
		}
	}
	return &ConcurrentDict{
		table: table,
		count: 0,
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

/*
	spread 用来获取要存储在哪个shard，这个位与运算意思是：
	根据hashcode来获取应该存在哪个 shard
	举个例子：比如 tableSize = 16 ， hashcode = 30；
	则 tableSize - 1 = 1111 （需要减一因为索引从0开始）， hashCode = 11110，
	& 运算取得低（tableSize - 1）位，得到 index = 1110，即索引 14
*/
func (c *ConcurrentDict) shard(hashCode uint32) uint32 {
	if c == nil {
		panic("shard a nil ConcurrentDict")
	}
	return uint32(len(c.table)-1) & hashCode
}

func (c *ConcurrentDict) getShard(index uint32) *shard {
	if c == nil {
		panic("get shard from a nil ConcurrentDict")
	}
	return c.table[index]
}

// Len 返回 ConcurrentDict 中元素的数量
func (c *ConcurrentDict) Len() int {
	if c.table == nil {
		panic("get len from a nil ConcurrentDict")
	}
	// https://blog.huoding.com/2021/10/08/958
	// 为什么需要 LoadInt32？ 加载 32位不是肯定原子吗？ 因为要统一函数接口,用于不同的os实现和更老于32位的系统
	// 怎么说呢，大概就是为了实现多平台的正确性吧
	return int(atomic.LoadUint32(&c.count))
}

// Get 获取对应的key的value，返回value和是否存在
func (c *ConcurrentDict) Get(key string) (interface{}, bool) {
	if c.table == nil {
		panic("get value from a nil ConcurrentDict")
	}
	hashCode := fnv32(key)
	index := c.shard(hashCode)
	shard := c.getShard(index)
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()
	val, ok := shard.m[key]
	return val, ok
}

// Put 放置kv并返回是否是新插入的
func (c *ConcurrentDict) Put(key string, val interface{}) bool {
	if c.table == nil {
		panic("put value to a nil ConcurrentDict")
	}
	hashCode := fnv32(key)
	index := c.shard(hashCode)
	shard := c.getShard(index)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; ok {
		shard.m[key] = val
		return false
	}
	shard.m[key] = val
	c.addCountL()
	return true
}

// Remove 删除 key 并返回是否有删除元素（可能key不存在
func (c *ConcurrentDict) Remove(key string) bool {
	if c.table == nil {
		panic("remove value from a nil ConcurrentDict")
	}
	hashCode := fnv32(key)
	index := c.shard(hashCode)
	shard := c.getShard(index)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; ok {
		delete(shard.m, key)
		c.decCountL()
		return true
	}
	return false
}

func (c *ConcurrentDict) addCountL() {
	atomic.AddUint32(&c.count, 1)
}

func (c *ConcurrentDict) decCountL() {
	atomic.StoreUint32(&c.count, c.count-1)
}
