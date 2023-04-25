package abstract

type QueueInterface interface {
	SyncProducer
	Consumer
	ConsumerWithHandler
}
