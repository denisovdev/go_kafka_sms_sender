package queue

type Queue interface {
	Produce(topic string, message []byte) error
}
