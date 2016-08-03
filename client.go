package stats

import "io"

type Client interface {
	io.Closer

	Gauge(Opts) Gauge

	Counter(Opts) Counter

	Histogram(Opts) Histogram
}

type Config struct {
	Backend Backend
	Scope   string
	Tags    Tags
}

func NewClient(scope string, backend Backend, tags ...Tag) Client {
	return NewClientWith(Config{
		Backend: backend,
		Scope:   scope,
		Tags:    tags,
	})
}

func NewClientWith(config Config) Client {
	return client{
		backend: config.Backend,
		scope:   config.Scope,
		tags:    config.Tags.Copy(),
	}
}

type client struct {
	backend Backend
	scope   string
	tags    Tags
}

func (c client) Close() error {
	return c.backend.Close()
}

func (c client) Gauge(opts Opts) Gauge {
	return NewGauge(c.opts(opts), c.backend.Set)
}

func (c client) Counter(opts Opts) Counter {
	return NewCounter(c.opts(opts), c.backend.Add)
}

func (c client) Histogram(opts Opts) Histogram {
	return NewHistogram(c.opts(opts), c.backend.Observe)
}

func (c client) opts(opts Opts) Opts {
	if len(opts.Scope) == 0 {
		opts.Scope = c.scope
	}

	n1 := len(c.tags)
	n2 := len(opts.Tags)

	tags := make(Tags, n1+n2)
	copy(tags, c.tags)
	copy(tags[n1:], opts.Tags)
	opts.Tags = tags

	return opts
}