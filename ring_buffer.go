package honeybadger

type ringBuffer struct {
	buf        []*eventPayload
	head, tail int // head: next pop, tail: next push
	size       int
	cap        int
}

func newRingBuffer(limit int) *ringBuffer {
	return &ringBuffer{buf: make([]*eventPayload, limit), cap: limit}
}

func (q *ringBuffer) push(it *eventPayload) bool {
	if q.size == q.cap {
		return false
	} // full
	q.buf[q.tail] = it
	q.tail = (q.tail + 1) % q.cap
	q.size++
	return true
}

func (q *ringBuffer) pop() *eventPayload {
	if q.size == 0 {
		return nil
	}
	it := q.buf[q.head]
	q.buf[q.head] = nil
	q.head = (q.head + 1) % q.cap
	q.size--
	return it
}

func (q *ringBuffer) drain() []*eventPayload {
	if q.size == 0 {
		return nil
	}

	out := make([]*eventPayload, q.size)
	if q.head < q.tail {
		copy(out, q.buf[q.head:q.tail])
	} else {
		n := copy(out, q.buf[q.head:])
		copy(out[n:], q.buf[:q.tail])
	}

	for i := range q.buf {
		q.buf[i] = nil
	}

	q.head, q.tail, q.size = 0, 0, 0

	return out
}

func (q *ringBuffer) len() int { return q.size }
