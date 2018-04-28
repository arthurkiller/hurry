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

func (n *Node) free() {
	n.pre.next = n.next
	n.next.pre = n.pre
	n.next, n.pre = nil, nil
}

// TODO not implement yet
func (n *Node) linsert(s *Node) {}
func (n *Node) rinsert(s *Node) {}

// LinkList is a double linked list implementation
type LinkList struct {
	Head     *Node
	Tail     *Node
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
	atomic.AddInt64(&nl.Len, -1)
	n := nl.Tail.pre
	n.Lock()
	n.pre.Lock()
	n.free()
	n.pre.Unlock()
	n.next.Unlock()
	n.Unlock()
	return n, true
}

// Push will push a Node to head
func (nl *LinkList) Push(n *Node) {
	nl.Head.Lock()
	nl.Head.next.Lock()
	n.next = nl.Head.next
	n.pre = nl.Head
	n.pre.next = n
	n.next.pre = n
	atomic.AddInt64(&nl.Len, 1)
	n.pre.Unlock()
	n.next.Unlock()
	for atomic.LoadInt64(&nl.Len) > int64(nl.maxLen) {
		nl.Pop()
	}
}

// Up make the given Node push to top
func (nl *LinkList) Up(n *Node) {
	if n == nil || n.pre == nil || n.next == nil {
		panic("invalid link node")
	}
	nl.Push(nl.UnLink(n))
}

// Unlink free a node from LinkList
func (nl *LinkList) UnLink(n *Node) *Node {
	if n == nil || n.pre == nil || n.next == nil {
		panic("invalid link node")
	}
	n.Lock()
	atomic.AddInt64(&nl.Len, -1)

	n.pre.Lock()
	n.next.Lock()
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

	var n *Node = nl.Head.next
	n.RLock()
	for n.key != key {
		n.next.RLock()
		n = n.next
		n.pre.RUnlock()
	}
	n.RUnlock()
	return nl.UnLink(n), true
}
