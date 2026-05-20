package query

type Option func(*Config)

type Config struct {
	Page     int
	PageSize int
}

func WithPagination(page, pageSize int) Option {
	return func(c *Config) {
		c.Page = page
		c.PageSize = pageSize
	}
}

func ApplyOptions(opts []Option) Config {
	c := Config{Page: 1, PageSize: 20}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

func (c *Config) Offset() int {
	if c.Page <= 0 {
		return 0
	}
	return (c.Page - 1) * c.PageSize
}
