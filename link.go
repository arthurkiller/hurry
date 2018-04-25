package hurry

import "sync"

var MaxLen int // TODO
var _NodePool chan *Node = make(chan *Node, MaxLen)

type Node struct {
	pre   *Node
	next  *Node
	value interface{}
	sync.RWMutex
}

type linkNode struct {
	Head *Node
	Tail *Node
}
