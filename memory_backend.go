package honeybadger

import (
	"fmt"
	"reflect"
	"sync"
)

// MemoryBackend is a Backend that writes error notices to a slice.  The
// MemoryBackend is mainly useful for testing and will cause issues if used in
// production. MemoryBackend is thread safe but order can't be guaranteed.
type MemoryBackend struct {
	Notices []*Notice
	mu      sync.Mutex
}

// NewMemoryBackend creates a new MemoryBackend.
func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		Notices: make([]*Notice, 0),
	}
}

// Notify adds the given payload (if it is a Notice) to Notices.
func (b *MemoryBackend) Notify(_ Feature, payload Payload) error {
	notice, ok := payload.(*Notice)
	if !ok {
		return fmt.Errorf("memory backend does not support payload of type %q", reflect.TypeOf(payload))
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.Notices = append(b.Notices, notice)

	return nil
}

// Reset clears the set of Notices
func (b *MemoryBackend) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Notices = b.Notices[:0]
}
