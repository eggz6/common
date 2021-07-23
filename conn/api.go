package conn

type Protoc uint8

const (
	Unknown Protoc = iota
	JSONProtocol
	SimpleString
)

type Session interface {
	Get(key string) interface{}
	Set(key string, val interface{})
}

type Auth interface {
	Pass(id, key string) (bool, error)
}
