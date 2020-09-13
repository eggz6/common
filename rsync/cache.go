package rsync

type Action func(key string, val interface{}) interface{}
type EggMap interface {
	Get(key string) (interface{}, bool)
	Put(key string, val interface{})
	GetAndDo(key string, ac Action) (interface{}, bool)
}
