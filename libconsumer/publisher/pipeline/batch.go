package pipeline

import (
	"sync"

	"github.com/elastic/beats/libbeat/publisher"
	"github.com/elastic/beats/libbeat/publisher/queue"
)

type Batch interface {
	publisher.Batch

	reduceTTL() bool
}

type batch struct {
	original queue.Batch
	ctx      *batchContext
	ttl      int
	events   []publisher.Event
}

type batchContext struct {
	observer outputObserver
	retryer  *retryer
}

var batchPool = sync.Pool{
	New: func() interface{} {
		return &batch{}
	},
}

func newBatch(ctx *batchContext, original queue.Batch, ttl int) *batch {
	if original == nil {
		panic("empty batch")
	}

	b := batchPool.Get().(*batch)
	*b = batch{
		original: original,
		ctx:      ctx,
		ttl:      ttl,
		events:   original.Events(),
	}
	return b
}

func releaseBatch(b *batch) {
	*b = batch{} // clear batch
	batchPool.Put(b)
}

func (b *batch) Events() []publisher.Event {
	return b.events
}

func (b *batch) ACK() {
	if b.ctx != nil {
		b.ctx.observer.outBatchACKed(len(b.events))
	}
	b.original.ACK()
	releaseBatch(b)
}

func (b *batch) Drop() {
	b.original.ACK()
	releaseBatch((b))
}

func (b *batch) Retry() {
	b.ctx.retryer.retry(b)
}

func (b *batch) Cancelled() {
	b.ctx.retryer.cancelled(b)
}

func (b *batch) RetryEvents(events []publisher.Event) {
	b.updEvents(events)
	b.Retry()
}

func (b *batch) CancelledEvents(events []publisher.Event) {
	b.updEvents(events)
	b.Cancelled()
}

func (b *batch) updEvents(events []publisher.Event) {
	l1 := len(b.events)
	l2 := len(events)
	if l1 > l2 {
		b.ctx.observer.outBatchACKed(l1 - l2)
	}

	b.events = events
}

// reduceTTL reduces the time to live for all events that have no 'guaranteed'
// sending requirements. reduceTTL returns true if the batch is still alive.
func (b *batch) reduceTLL() bool {
	if b.ttl <= 0 {
		return true
	}

	b.ttl--
	if b.ttl > 0 {
		return true
	}

	// filter for events with guaranteed send flags
	events := b.events[:0]
	for _, event := range b.events {
		if event.Guaranteed() {
			events = append(events, event)
		}
	}
	b.events = events

	if len(b.events) > 0 {
		b.ttl = -1 // we need infinite retry for all events left in this batch
		return true
	}

	// all events have been dropped
	return false
}
