package honeybadger

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Batch struct {
	events   []*eventPayload
	attempts int
}

type EventsWorker struct {
	backend         Backend
	batchSize       int
	throttleWait    time.Duration
	timeout         time.Duration
	maxQueueSize    int
	maxRetries      int
	dropLogInterval time.Duration
	logger          Logger

	ticker      *time.Ticker
	dropTicker  *time.Ticker
	queue       *ringBuffer
	queueSize   int
	batches     []*Batch
	throttling atomic.Bool
	dropped    atomic.Int64
	lastDropLog time.Time

	in         chan *eventPayload
	flushCh    chan chan struct{}
	shutdownCh chan struct{}

	wg   sync.WaitGroup
	once sync.Once
}

func NewEventsWorker(cfg *Configuration) *EventsWorker {
	ctx := cfg.Context
	if ctx == nil {
		ctx = context.Background()
	}

	w := &EventsWorker{
		backend:         cfg.Backend,
		batchSize:       cfg.EventsBatchSize,
		timeout:         cfg.EventsTimeout,
		maxQueueSize:    cfg.EventsMaxQueueSize,
		maxRetries:      cfg.EventsMaxRetries,
		throttleWait:    cfg.EventsThrottleWait,
		dropLogInterval: cfg.EventsDropLogInterval,
		logger:          cfg.Logger,
		// +1 so we can push before checking flush threshold without dropping an event.
		queue:      newRingBuffer(cfg.EventsBatchSize + 1),
		queueSize:  0,
		batches:    make([]*Batch, 0),
		in:         make(chan *eventPayload, cfg.EventsMaxQueueSize),
		flushCh:    make(chan chan struct{}, 1),
		shutdownCh: make(chan struct{}),
	}
	w.wg.Add(1)
	go w.run(ctx)
	return w
}

func (w *EventsWorker) Push(e *eventPayload) {
	select {
	case w.in <- e:
	default:
		w.dropped.Add(1)
	}
}

func (w *EventsWorker) Flush() {
	// Check if already stopped
	select {
	case <-w.shutdownCh:
		return
	default:
	}

	done := make(chan struct{})
	select {
	case w.flushCh <- done:
		<-done
	case <-w.shutdownCh:
	}
}

func (w *EventsWorker) Stop() {
	w.once.Do(func() {
		close(w.shutdownCh)
		w.wg.Wait()
	})
}

func (w *EventsWorker) logDropSummary() {
	dropped := w.dropped.Swap(0)
	if dropped > 0 {
		w.logger.Printf("events worker dropped %d events due to full queue (capacity: %d, current size: %d)\n", dropped, w.maxQueueSize, w.queueSize)
		w.lastDropLog = time.Now()
	}
}

func (w *EventsWorker) AttemptSend() bool {
	events := w.queue.drain()
	if len(events) > 0 {
		w.batches = append(w.batches, &Batch{events: events, attempts: 0})
	}
	w.queue = newRingBuffer(w.batchSize + 1)

	for len(w.batches) > 0 {
		if w.throttling.Load() {
			break
		}

		batch := w.batches[0]
		if batch.attempts > w.maxRetries {
			w.logger.Printf("events worker dropping batch after %d failed attempts\n", batch.attempts)
			w.batches = w.batches[1:]
			w.queueSize -= len(batch.events)
			continue
		}

		err := w.backend.Event(batch.events)

		if err == ErrRateExceeded {
			w.logger.Printf("events worker received rate limit; throttling for %v\n", w.throttleWait)
			w.throttling.Store(true)
			go func() {
				time.Sleep(w.throttleWait)
				w.throttling.Store(false)
				w.logger.Printf("events worker throttle window expired; resuming sends\n")
				w.Flush()
			}()
			break
		} else if err != nil {
			batch.attempts++
			w.logger.Printf("events worker send error: %v\n", err)
			break
		} else {
			w.batches = w.batches[1:]
			w.queueSize -= len(batch.events)
		}
	}

	return len(w.batches) > 0
}

func (w *EventsWorker) run(ctx context.Context) {
	defer w.wg.Done()

	w.ticker = time.NewTicker(w.timeout)
	defer w.ticker.Stop()

	var dropTickerCh <-chan time.Time
	if w.dropLogInterval > 0 {
		w.dropTicker = time.NewTicker(w.dropLogInterval)
		defer w.dropTicker.Stop()
		dropTickerCh = w.dropTicker.C
	}

	flush := func() {
		if w.queue.len() == 0 && len(w.batches) == 0 {
			return
		}

		hasPendingBatches := w.AttemptSend()
		if hasPendingBatches && !w.throttling.Load() {
			w.ticker.Reset(w.timeout)
		}
	}

	for {
		select {
		case <-ctx.Done():
			flush()
			w.logDropSummary()
			return

		case <-w.shutdownCh:
			flush()
			w.logDropSummary()
			return

		case <-w.ticker.C:
			flush()

		case <-dropTickerCh:
			w.logDropSummary()

		case done := <-w.flushCh:
			// Drain pending events from input channel before flushing
		drainLoop:
			for {
				select {
				case e := <-w.in:
					w.queue.push(e)
					if w.queueSize >= w.maxQueueSize {
						if w.queue.len() > 0 {
							w.queue.pop()
							w.dropped.Add(1)
						}
					} else {
						w.queueSize++
					}
				default:
					break drainLoop
				}
			}
			flush()
			close(done)

		case e := <-w.in:
			w.queue.push(e)

			if w.queueSize >= w.maxQueueSize {
				// Drop oldest if at capacity
				if w.queue.len() > 0 {
					w.queue.pop()
					w.dropped.Add(1)
				}
			} else {
				w.queueSize++
			}

			if w.queue.len() >= w.batchSize {
				flush()
				w.ticker.Reset(w.timeout)
			}
		}
	}
}
