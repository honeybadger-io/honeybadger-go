package honeybadger

type Envelope func() error

type Worker interface {
	Push(Envelope) error
	Flush() error
}
