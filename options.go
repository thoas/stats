package stats

// Options are stats options.
type Options struct {
	statusCode *int
	recorder   ResponseWriter
}

func (o Options) StatusCode() int {
	if o.recorder != nil {
		return o.recorder.Status()
	}

	return *o.statusCode
}

// Option represents a stats option.
type Option func(*Options)

// WithStatusCode sets the status code to use in stats.
func WithStatusCode(statusCode int) Option {
	return func(o *Options) {
		o.statusCode = &statusCode
	}
}

// WithRecorder sets the recorder to use in stats.
func WithRecorder(recorder ResponseWriter) Option {
	return func(o *Options) {
		o.recorder = recorder
	}
}

// newOptions takes functional options and returns options.
func newOptions(options ...Option) *Options {
	opts := &Options{}
	for _, o := range options {
		o(opts)
	}
	return opts
}
