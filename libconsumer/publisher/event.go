package publisher

type Batch interface {
	// signals
	ACK()
	Drop()
	Retry()
	Canceled()
}

type Event struct {
}
