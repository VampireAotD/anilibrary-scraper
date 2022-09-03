package closer

import (
	"io"
	"sync"

	logger2 "anilibrary-request-parser/pkg/logger"
)

type Closer struct {
	scope    string
	callback io.Closer
}

type Closers struct {
	logger  logger2.Logger
	mutex   sync.Mutex
	closers []Closer
}

func New(logger logger2.Logger) Closers {
	return Closers{logger: logger}
}

func (c *Closers) Add(scope string, callback io.Closer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	closer := Closer{
		scope:    scope,
		callback: callback,
	}

	c.closers = append(c.closers, closer)
}

func (c *Closers) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.closers {
		closer := c.closers[len(c.closers)-1-i]

		if err := closer.callback.Close(); err != nil {
			c.logger.Error("close", logger2.String("scope", closer.scope), logger2.Error(err))
		}
	}

	c.closers = nil
}
