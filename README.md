# hurry
hurry is a LRU cache in go

## Features
* LRU
* LRUK // TODO
* TwoQueue // TODO
* KQueue // TODO

## Powerfull Interface

``` golang
type Hurry interface {
	Get(key string) interface{}
	Put(key string, obj interface{})
	Delete(key ...string)
	Exist(key string) bool
	Len() int64
}
```

## QuickStart
>> go get arthurkiller/hurry

```golang
lru := NewHurry(100, hurry.LRU)
lru.Get("foo")
```
