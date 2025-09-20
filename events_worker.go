package honeybadger

import (
	"context"
	"sync"
	"time"
)

type EventsWorker struct {
	backend   Backend
	batchSize int
	timeout   time.Duration

	in      chan *eventPayload
	flushCh chan struct{}

	cancel context.CancelFunc
	wg     sync.WaitGroup
	once   sync.Once
}

func NewEventsWorker(cfg *Configuration) *EventsWorker {
	ctx, cancel := context.WithCancel(context.Background())

	w := &EventsWorker{
		backend:   cfg.Backend,
		batchSize: cfg.EventsBatchSize,
		timeout:   cfg.EventsTimeout,
		in:        make(chan *eventPayload),
		flushCh:   make(chan struct{}, 1),
		cancel:    cancel,
	}
	w.wg.Add(1)
	go w.run(ctx)
	return w
}

func (w *EventsWorker) Push(e *eventPayload) { w.in <- e }

func (w *EventsWorker) Flush() {
	select {
	case w.flushCh <- struct{}{}:
	default:
	}
}

func (w *EventsWorker) Stop() {
	w.once.Do(func() {
		w.cancel()
		w.wg.Wait()
	})
}

func (w *EventsWorker) run(ctx context.Context) {
	defer w.wg.Done()

	batch := make([]*eventPayload, 0, w.batchSize)
	ticker := time.NewTicker(w.timeout)
	defer ticker.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		w.backend.Event(batch)
		batch = batch[:0]
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			return

		case <-ticker.C:
			flush()

		case <-w.flushCh:
			flush()

		case e := <-w.in:
			batch = append(batch, e)
			if len(batch) >= w.batchSize {
				flush()
				ticker.Reset(w.timeout)
			}
		}
	}
}