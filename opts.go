package toposort

type Option func(*config)

type config struct {
	bp Buffers
}

func WithBuffers(bp Buffers) Option {
	return func(cfg *config) {
		cfg.bp = bp
	}
}
