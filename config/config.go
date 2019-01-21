package config

type Option interface {
	apply(*option)
}

type Config interface {
}

type ConfigImp struct {
	A     int
	opts  *option
}

func NewConfig(a int, opts ...Option) ConfigImp {
	cfg := ConfigImp{a, newDefaultOpt()}
	for _, opt := range opts {
		opt.apply(cfg.opts)
	}
	return cfg
}

func newDefaultOpt() *option {
	return new(option)
}

type option struct {
	Timeout  bool
	C        int
}

type funcOption struct {
	f func(*option)
}

func (fo *funcOption) apply(opt *option) {
	fo.f(opt)
}

func newFuncOption(f func(*option)) Option {
	return &funcOption{f}
}

func WithTimeout() Option {
	return newFuncOption(func(opt *option) {
		opt.Timeout = true
	})
}

func WithC(c int) Option {
	return newFuncOption(func(opt *option) {
		opt.C = c
	})
}

