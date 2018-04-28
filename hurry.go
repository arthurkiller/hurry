package hurry

type Mode int

const (
	Default Mode = iota
	LRUK
	TwoQueue
	KQueue
)

// Hurry is a parallel safe LRU cacue focused on performance
type Hurry interface {
	Get(key string) interface{}
	Put(key string, obj interface{})
	Delete(key ...string)

	Exist(key string) bool
	Len() int64

	//GetFirstN(n int) []interface{}
	//GetLastN(n int) []interface{}
	//GetAll() []interface{}
}

// NewHurry will gen a Hurry implement LRU cache
func NewHurry(n int, m Mode) Hurry {
	// TODO: Switch different LRU implement
	return NewLRU(n)
}
