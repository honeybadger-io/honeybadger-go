package honeybadger

import "fmt"

var (
	WorkerOverflowError = fmt.Errorf("The worker is full; this envelope will be dropped.")
)

func newBufferedWorker() Worker {
	worker := make(BufferedWorker, 100)
	go func() {
		for w := range worker {
			work := func() error {
				defer func() {
					if err := recover(); err != nil {
						fmt.Printf("worker recovered from panic: %v\n", err)
					}
				}()
				return w()
			}
			if err := work(); err != nil {
				fmt.Printf("worker processing error: %v\n", err)
			}
		}
	}()
	return worker
}

type BufferedWorker chan Envelope

func (worker BufferedWorker) Push(work Envelope) error {
	select {
	case worker <- work:
		return nil
	default:
		return WorkerOverflowError
	}
}

func (worker BufferedWorker) Flush() error {
	ch := make(chan bool)
	worker <- func() error {
		ch <- true
		return nil
	}
	<-ch
	return nil
}
