package engine

import (
	"hash/fnv"
	"math"
	"sync"
)

// BloomFilter 布隆过滤器实现
type BloomFilter struct {
	bitset    []bool
	size      uint
	hashFuncs uint
	mutex     sync.RWMutex
}

// NewBloomFilter 创建一个新的布隆过滤器
// expectedItems: 预期元素数量
// falsePositiveRate: 可接受的误判率 (0.0 - 1.0)
func NewBloomFilter(expectedItems uint, falsePositiveRate float64) *BloomFilter {
	// 计算最佳大小和哈希函数数量
	size := calculateOptimalSize(expectedItems, falsePositiveRate)
	hashFuncs := calculateOptimalHashFunctions(size, expectedItems)

	return &BloomFilter{
		bitset:    make([]bool, size),
		size:      size,
		hashFuncs: hashFuncs,
		mutex:     sync.RWMutex{},
	}
}

// Add 添加元素到布隆过滤器
func (bf *BloomFilter) Add(item string) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	// 计算多个哈希值并设置对应位
	for i := uint(0); i < bf.hashFuncs; i++ {
		position := bf.hash(item, i) % bf.size
		bf.bitset[position] = true
	}
}

// Contains 检查元素是否可能存在于布隆过滤器中
func (bf *BloomFilter) Contains(item string) bool {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	// 检查所有哈希位置是否都被设置
	for i := uint(0); i < bf.hashFuncs; i++ {
		position := bf.hash(item, i) % bf.size
		if !bf.bitset[position] {
			// 如果有任何一个位未设置，则元素肯定不存在
			return false
		}
	}

	// 所有位都被设置，元素可能存在（有误判可能性）
	return true
}

// hash 计算哈希值
func (bf *BloomFilter) hash(item string, seed uint) uint {
	h := fnv.New64a()
	h.Write([]byte(item))
	h.Write([]byte{byte(seed), byte(seed >> 8)})
	return uint(h.Sum64())
}

// calculateOptimalSize 计算布隆过滤器的最佳大小
func calculateOptimalSize(n uint, p float64) uint {
	return uint(math.Ceil(-float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
}

// calculateOptimalHashFunctions 计算最佳哈希函数数量
func calculateOptimalHashFunctions(m, n uint) uint {
	return uint(math.Ceil(float64(m) / float64(n) * math.Log(2)))
}
