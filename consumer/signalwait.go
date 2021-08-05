package consumer

type signalWait struct {
	count   int // number of potential 'alive' signals
	signals chan struct{}
}

type signaler func()

func newSignalWait() *signalWait {
	return &signalWait{
		signals: make(chan struct{}, 1),
	}
}

func (s *signalWait) Wait() {
	if s.count == 0 {
		return
	}

	<-s.signals
	s.count--
}

func (s *signalWait) Add(fn signaler) {
	s.count++
	go func() {
		fn()
		var v struct{}
		s.signals <- v
	}()
}

func (s *signalWait) AddChan(c <-chan struct{}) {
	s.Add(waitChannel(c))
}

func waitChannel(c <-chan struct{}) signaler {
	return func() { <-c }
}
