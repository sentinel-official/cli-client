package context

type ServerContext struct {
	home string
}

func NewServerContext() ServerContext {
	return ServerContext{}
}

func (c ServerContext) WithHome(v string) ServerContext {
	c.home = v
	return c
}

func (c ServerContext) Home() string {
	return c.home
}
