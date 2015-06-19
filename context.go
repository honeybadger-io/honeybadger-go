package honeybadger

// A key/value map used to send extra data to Honeybadger.
type Context hash

// Updates the context object.
func (target Context) Update(context Context) {
	for k, v := range context {
		target[k] = v
	}
}
