package retry

type Option struct {
	MaxRetries int
}

func WithMaxRetries(maxRetries int) Option {
	return Option{MaxRetries: maxRetries}
}

// Do retries the function until it returns nil or the max retries is reached.
// It returns the last error if the max retries is reached.
//
// NOTE: default max retries is 2.
func Do(fn func() error, opts Option) error {
	if opts.MaxRetries == 0 {
		opts.MaxRetries = 2
	}

	var err error
	for i := 0; i < opts.MaxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
	}
	return err
}
