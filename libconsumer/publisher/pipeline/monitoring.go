package pipeline

type observer interface {
	pipelineObserver
	clientObserver
	queueObserver
	outputObserver

	cleanup()
}

type pipelineObserver interface {
	clientConnected()
	clientClosing()
	clientClosed()
}

type clientObserver interface {
	newEvent()
	filteredEvent()
	publishedEvent()
	failedPublishEvent()
}

type queueObserver interface {
	queueACKed(n int)
	queueMaxEvents(n int)
}

type outputObserver interface {
	updateOutputGroup()
	eventsFailed(int)
	eventsDropped(int)
	eventsRetry(int)
	outBatchSend(int)
	outBatchACKed(int)
}

type emptyObserver struct{}

var nilObserver observer = (*emptyObserver)(nil)

func (*emptyObserver) cleanup()            {}
func (*emptyObserver) clientConnected()    {}
func (*emptyObserver) clientClosing()      {}
func (*emptyObserver) clientClosed()       {}
func (*emptyObserver) newEvent()           {}
func (*emptyObserver) filteredEvent()      {}
func (*emptyObserver) publishedEvent()     {}
func (*emptyObserver) failedPublishEvent() {}
func (*emptyObserver) queueACKed(n int)    {}
func (*emptyObserver) queueMaxEvents(int)  {}
func (*emptyObserver) updateOutputGroup()  {}
func (*emptyObserver) eventsFailed(int)    {}
func (*emptyObserver) eventsDropped(int)   {}
func (*emptyObserver) eventsRetry(int)     {}
func (*emptyObserver) outBatchSend(int)    {}
func (*emptyObserver) outBatchACKed(int)   {}
