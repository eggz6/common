package poker

type EventEntry interface {
	Name() string
	Data() interface{}
	Target() interface{}
}
