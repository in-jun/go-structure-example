package query

type Option func(*Config)

type Config struct {
	ForUpdate bool
}

func ForUpdate() Option {
	return func(c *Config) { c.ForUpdate = true }
}

func ApplyOptions(opts []Option) Config {
	var c Config
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
