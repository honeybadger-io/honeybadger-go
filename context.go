package honeybadger

// A key/value map used to send extra data to Honeybadger.
type Context hash

// Returns a new Context with values merged.
func (c1 Context) merge(c2 Context) Context {
	for k, v := range c2 {
		c1[k] = v
	}
	return c1
}
