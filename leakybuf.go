// leakybuf provides leaky buffer, based on the example in Effective Go.
package hurry

type NodePool chan *Node

// NewNodePool creates a leaky buffer which can hold at most n buffer, each
// with bufSize bytes.
func NewNodePool(n int) NodePool {
	return make(chan *Node, n)
}

// Get returns a buffer from the leaky buffer or create a new buffer.
func (lb *NodePool) Get() (n *Node) {
	select {
	case n = <-*lb:
	default:
		n = &Node{}
	}
	return
}

// Put add the buffer into the free buffer pool for reuse. Panic if the buffer
// size is not the same with the leaky buffer's. This is intended to expose
// error usage of leaky buffer.
func (lb *NodePool) Put(n *Node) {
	select {
	case *lb <- n:
	default:
	}
	return
}
