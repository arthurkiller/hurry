package hurry

import "testing"

// test racing in leakbuf
func TestLeakybuf(t *testing.T) {
	b := NewLeakyBuf(100, 100)
	go func() {
		bf := b.Get()
		b.Put(bf)
	}()
}

// benchmark leakbuf
func BenchmarkGetPut(b *testing.B) {
	bf := NewLeakyBuf(100, 100)
	for i := 0; i < b.N; i++ {
		k := bf.Get()
		bf.Put(k)
	}
}

// benchmark leakbuf parallel
func BenchmarkGetPutParallel(b *testing.B) {
	bf := NewLeakyBuf(100, 100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			k := bf.Get()
			bf.Put(k)
		}
	})
}
