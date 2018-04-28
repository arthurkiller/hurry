package hurry

import "testing"

// test racing
func TestLeakybuf(t *testing.T) {
	b := NewNodePool(100)
	go func() {
		bf := b.Get()
		b.Put(bf)
	}()
}

// benchmark
func BenchmarkGetPut(b *testing.B) {
	bf := NewNodePool(100)
	for i := 0; i < b.N; i++ {
		k := bf.Get()
		bf.Put(k)
	}
}

// benchmark in parallel
func BenchmarkGetPutParallel(b *testing.B) {
	bf := NewNodePool(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := bf.Get()
			bf.Put(k)
		}
	})
}
