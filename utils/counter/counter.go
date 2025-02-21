package counter

import (
	"fmt"
	"github.com/kercylan98/minotaur/utils/generic"
	"github.com/kercylan98/minotaur/utils/hash"
	"sync"
	"time"
)

// NewCounter 创建一个计数器
func NewCounter[K comparable, V generic.Number]() *Counter[K, V] {
	c := &Counter[K, V]{
		c:  make(map[K]any),
		dr: make(map[K]int64),
	}
	return c
}

// newSubCounter 创建一个子计数器
func newSubCounter[K comparable, V generic.Number]() *Counter[K, V] {
	c := &Counter[K, V]{
		c:   make(map[K]any),
		dr:  make(map[K]int64),
		sub: true,
	}
	return c
}

// Counter 计数器
type Counter[K comparable, V generic.Number] struct {
	lock sync.RWMutex
	sub  bool
	dr   map[K]int64
	drm  map[string]int64
	c    map[K]any
}

// AddWithMark 添加计数，根据 mark 判断是否重复计数
func (slf *Counter[K, V]) AddWithMark(mark, key K, value V, expired time.Duration) {
	slf.lock.Lock()
	if expired > 0 {
		now := time.Now().Unix()
		mk := fmt.Sprintf("%v_%v", mark, key)
		expiredTime, exist := slf.drm[mk]
		if exist {
			if expiredTime > now {
				slf.lock.Unlock()
				return
			}
		}
		slf.drm[mk] = now + int64(expired/time.Second)
	}

	v, exist := slf.c[key]
	if !exist {
		slf.c[key] = value
		slf.lock.Unlock()
		return
	}
	if v, ok := v.(V); !ok {
		slf.lock.Unlock()
		panic("counter value is sub counter")
	} else {
		slf.c[key] = v + value
	}
	slf.lock.Unlock()
}

// Add 添加计数
//   - 当设置了 expired 时，在 expired 时间内，不会重复计数
//   - 需要注意的是，当第一次设置了 expired，第二次未设置时，第二次的计数将生效
func (slf *Counter[K, V]) Add(key K, value V, expired ...time.Duration) {
	slf.lock.Lock()
	if len(expired) > 0 && expired[0] > 0 {
		now := time.Now().Unix()
		expiredTime, exist := slf.dr[key]
		if exist {
			if expiredTime > now {
				slf.lock.Unlock()
				return
			}
		}
		slf.dr[key] = now + int64(expired[0]/time.Second)
	}

	v, exist := slf.c[key]
	if !exist {
		slf.c[key] = value
		slf.lock.Unlock()
		return
	}
	if v, ok := v.(V); !ok {
		slf.lock.Unlock()
		panic("counter value is sub counter")
	} else {
		slf.c[key] = v + value
	}
	slf.lock.Unlock()
}

// Get 获取计数
func (slf *Counter[K, V]) Get(key K) V {
	slf.lock.RLock()
	v, exist := slf.c[key]
	if !exist {
		slf.lock.RUnlock()
		return 0
	}
	if v, ok := v.(V); !ok {
		slf.lock.RUnlock()
		panic("counter value is sub counter")
	} else {
		slf.lock.RUnlock()
		return v
	}
}

// Reset 重置特定 key 的计数
//   - 当 key 为一个子计数器时，将会重置该子计数器
func (slf *Counter[K, V]) Reset(key K) {
	slf.lock.Lock()
	delete(slf.c, key)
	delete(slf.dr, key)
	slf.lock.Unlock()
}

// ResetMark 重置特定 mark 的 key 的计数
func (slf *Counter[K, V]) ResetMark(mark, key K) {
	slf.lock.Lock()
	mk := fmt.Sprintf("%v_%v", mark, key)
	delete(slf.c, key)
	delete(slf.drm, mk)
	slf.lock.Unlock()
}

// ResetExpired 重置特定 key 的过期时间
func (slf *Counter[K, V]) ResetExpired(key K) {
	slf.lock.Lock()
	delete(slf.dr, key)
	slf.lock.Unlock()
}

// ResetExpiredMark 重置特定 mark 的 key 的过期时间
func (slf *Counter[K, V]) ResetExpiredMark(mark, key K) {
	slf.lock.Lock()
	mk := fmt.Sprintf("%v_%v", mark, key)
	delete(slf.drm, mk)
	slf.lock.Unlock()
}

// ResetAll 重置所有计数
func (slf *Counter[K, V]) ResetAll() {
	slf.lock.Lock()
	hash.Clear(slf.c)
	hash.Clear(slf.dr)
	hash.Clear(slf.drm)
	slf.lock.Unlock()
}

// SubCounter 获取子计数器
func (slf *Counter[K, V]) SubCounter(key K) *Counter[K, V] {
	slf.lock.Lock()
	v, exist := slf.c[key]
	if !exist {
		counter := newSubCounter[K, V]()
		slf.c[key] = counter
		slf.lock.Unlock()
		return counter
	}
	if v, ok := v.(*Counter[K, V]); !ok {
		slf.lock.Unlock()
		panic("counter value is count value")
	} else {
		slf.lock.Unlock()
		return v
	}
}

// GetCounts 获取计数器的所有计数
func (slf *Counter[K, V]) GetCounts() map[K]V {
	counts := make(map[K]V)
	slf.lock.RLock()
	for k, v := range slf.c {
		if v, ok := v.(V); !ok {
			continue
		} else {
			counts[k] = v
		}
	}
	slf.lock.RUnlock()
	return counts
}

// GetSubCounters 获取计数器的所有子计数器
func (slf *Counter[K, V]) GetSubCounters() map[K]*Counter[K, V] {
	counters := make(map[K]*Counter[K, V])
	slf.lock.RLock()
	for k, v := range slf.c {
		if v, ok := v.(*Counter[K, V]); !ok {
			continue
		} else {
			counters[k] = v
		}
	}
	slf.lock.RUnlock()
	return counters
}

// Shadow 获取该计数器的影子计数器，影子计数器的任何操作都不会影响到原计数器
//   - 适用于持久化等场景
func (slf *Counter[K, V]) Shadow() *Shadow[K, V] {
	return newShadow(slf)
}
