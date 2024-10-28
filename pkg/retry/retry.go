package retry

import "errors"

const DefaultMaxRetries = 2

type Option struct {
	maxRetries int
}

func WithMaxRetries(maxRetries int) Option {
	if maxRetries < 0 {
		return Option{}
	}
	return Option{maxRetries: maxRetries}
}

// Do retries the function until it returns nil or the max retries is reached.
// It returns the last error if the max retries is reached.
//
// NOTE: default max retries is 2.
func Do(fn func() error, opts Option) error {
	if opts.maxRetries <= 0 {
		opts.maxRetries = DefaultMaxRetries
	}

	var err error
	for i := 0; i < opts.maxRetries; i++ {
		fnerr := fn()
		if fnerr != nil {
			err = errors.Join(err, fnerr)
		}
	}

	return err
}
