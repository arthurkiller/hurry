package hurry

import (
	"sync"
	"sync/atomic"
)

// Node is the LinkList elemenet
type Node struct {
	pre   *Node
	next  *Node
	key   string
	value interface{}
	sync.RWMutex
}

// unsafe free node
func (n *Node) free() {
	n.pre.next = n.next
	n.next.pre = n.pre
	n.next, n.pre = nil, nil
}

// unsafe ljoin
func (n *Node) ljoin(s *Node) {
	s.pre = n.pre
	s.next = n
	n.pre.next = s
	n.pre = s
}

// unsafe rjoin
func (n *Node) rjoin(s *Node) {
	s.pre = n
	s.next = n.next
	n.next.pre = s
	n.next = s
}

// LinkList is a double linked list implementation
type LinkList struct {
	Head *Node
	Tail *Node
	// TODO should not be modified outside
	Len      int64
	nodePool NodePool
	maxLen   int
}

// NewLinkList retuen a
func NewLinkList(n int) *LinkList {
	no := LinkList{
		Head:     &Node{},
		Tail:     &Node{},
		nodePool: NewNodePool(n),
		maxLen:   n,
	}
	no.Head.next = no.Tail
	no.Tail.pre = no.Head

	return &no
}

// Pop will return the last Node in list
func (nl *LinkList) Pop() (*Node, bool) {
	if atomic.LoadInt64(&nl.Len) == 0 {
		return nil, false
	}

	nl.Tail.Lock()
	n := nl.Tail.pre
	n.Lock()
	n.pre.Lock()
	atomic.AddInt64(&nl.Len, -1)
	n.free()
	n.pre.Unlock()
	n.next.Unlock()
	n.Unlock()
	return n, true
}

// Push will push a Node to head
func (nl *LinkList) Push(n *Node) {
	n.Lock()
	nl.Head.Lock()
	nl.Head.next.Lock()
	atomic.AddInt64(&nl.Len, 1)
	nl.Head.rjoin(n)
	n.pre.Unlock()
	n.next.Unlock()
	n.Unlock()
}

// Up make the given Node push to top
// FIXME maybe race
func (nl *LinkList) Up(n *Node) {
	if n == nil || n.pre == nil || n.next == nil {
		panic("invalid link node")
	}
	n.Lock()
	n.pre.Lock()
	n.next.Lock()
	n.free()
	n.pre.Unlock()
	n.next.Unlock()
	nl.Head.Lock()
	nl.Head.next.Lock()
	nl.Head.rjoin(n)
	n.pre.Unlock()
	n.next.Unlock()
	n.Unlock()
}

// Unlink free a node from LinkList
func (nl *LinkList) UnLink(n *Node) *Node {
	if n == nil || n.pre == nil || n.next == nil {
		panic("invalid link node")
	}
	n.Lock()
	n.pre.Lock()
	n.next.Lock()
	atomic.AddInt64(&nl.Len, -1)
	n.free()
	n.pre.Unlock()
	n.next.Unlock()
	n.Unlock()
	return n
}

// Find return the Node pointed to the key
func (nl *LinkList) Find(key string) (*Node, bool) {
	if atomic.LoadInt64(&nl.Len) == 0 {
		return nil, false
	}

	nl.Head.next.RLock()
	var n *Node = nl.Head.next
	for n.key != key && n != nl.Tail {
		n.next.RLock()
		n = n.next
		n.pre.RUnlock()
	}
	n.RUnlock()
	if n == nl.Tail {
		return nil, false
	}
	return nl.UnLink(n), true
}
