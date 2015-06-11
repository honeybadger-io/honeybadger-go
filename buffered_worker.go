package honeybadger

import "fmt"

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
		return fmt.Errorf("the channel is full.")
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
