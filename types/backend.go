package types

type Backend interface {
	Publish(interface{}) error
	Consume() interface{}
}
