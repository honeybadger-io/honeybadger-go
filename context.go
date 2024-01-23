package honeybadger

import "context"

// Context is used to send extra data to Honeybadger.
type Context hash

// ctxKey is use in WithContext and FromContext to store and load the
// honeybadger.Context into a context.Context.
type ctxKey struct{}

// Update applies the values in other Context to context.
func (c Context) Update(other Context) {
	for k, v := range other {
		c[k] = v
	}
}

// WithContext adds the honeybadger.Context to the given context.Context and
// returns the new context.Context.
func (c Context) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// FromContext retrieves a honeybadger.Context from the context.Context.
// FromContext will return nil if no Honeybadger context exists in ctx.
func FromContext(ctx context.Context) Context {
	if c, ok := ctx.Value(ctxKey{}).(Context); ok {
		return c
	}

	return nil
}
