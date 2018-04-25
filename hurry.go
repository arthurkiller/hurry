package hurry

// Hurry is a parallel safe LRU cacue focused on performance
type Hurry interface {
	Get(key string) interface{}
	Put(key string, obj interface{})
	Delete(key ...string)

	Exist(key string) bool
	Len() int
	GetFirstN(n int) []interface{}
	GetLastN(n int) []interface{}
	GetAll() []interface{}
}
