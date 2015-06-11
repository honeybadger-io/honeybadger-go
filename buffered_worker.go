package honeybadger

import (
	"fmt"
	"log"
)

var (
	WorkerOverflowError = fmt.Errorf("The worker is full; this envelope will be dropped.")
)

func newBufferedWorker() BufferedWorker {
	worker := BufferedWorker{ch: make(chan Envelope, 100)}
	go func() {
		for w := range worker.ch {
			work := func() error {
				defer func() {
					if err := recover(); err != nil {
						worker.log.Printf("worker recovered from panic: %v\n", err)
					}
				}()
				return w()
			}
			if err := work(); err != nil {
				worker.log.Printf("worker processing error: %v\n", err)
			}
		}
	}()
	return worker
}

type BufferedWorker struct {
	ch  chan Envelope
	log *log.Logger
}

func (w BufferedWorker) Push(work Envelope) error {
	select {
	case w.ch <- work:
		return nil
	default:
		return WorkerOverflowError
	}
}

func (w BufferedWorker) Flush() error {
	ch := make(chan bool)
	w.ch <- func() error {
		ch <- true
		return nil
	}
	<-ch
	return nil
}
