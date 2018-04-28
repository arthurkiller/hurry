package hurry

import (
	"sync"
	"sync/atomic"
)

// map + link_list for both random get and put performence.
type LRU struct {
	sync.RWMutex
	lkl      *LinkList
	m        map[string]*Node
	nodePool NodePool
}

func NewLRU(n int) *LRU {
	return &LRU{
		lkl:      NewLinkList(n),
		m:        make(map[string]*Node, n),
		nodePool: NewNodePool(n),
	}
}

func (l *LRU) Get(key string) interface{} {
	if n, ok := l.m[key]; ok {
		return n
	}
	return nil
}

func (l *LRU) Put(key string, obj interface{}) {
	l.RLock()
	if n, ok := l.m[key]; ok {
		l.RUnlock()
		// TODO FIXME
		n.Lock()
		n.value = obj
		n.Unlock()
		l.lkl.Up(n)
		return
	}

	l.RUnlock()
	n := l.nodePool.Get()
	n.key = key
	n.value = obj

	l.Lock()
	l.m[key] = n
	l.lkl.Push(n)
	l.Unlock()

	l.Lock()
	for atomic.LoadInt64(&l.lkl.Len) > int64(l.lkl.maxLen) {
		if nr, ok := l.lkl.Pop(); ok {
			delete(l.m, nr.key)
			l.nodePool.Put(nr)
		}
	}
	l.Unlock()
	return
}

func (l *LRU) Delete(key ...string) {
	for _, k := range key {
		if n, ok := l.m[k]; ok {
			l.Lock()
			delete(l.m, k)
			l.lkl.UnLink(n)
			l.nodePool.Put(n)
			l.Unlock()
		}
	}
}

func (l *LRU) Exist(key string) bool {
	l.RLock()
	defer l.RUnlock()
	if _, ok := l.m[key]; ok {
		return ok
	}
	return false
}

func (l *LRU) Len() int64 {
	// TODO fixme
	return atomic.LoadInt64(&l.lkl.Len)
}

func (l *LRU) GetFirstN(n int) []interface{} {
	return nil
}
func (l *LRU) GetLastN(n int) []interface{} {
	return nil
}
func (l *LRU) GetAll() []interface{} {
	return nil
}
