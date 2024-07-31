package queue

type Producer interface {
	Produce(topic string, message []byte) error
}
