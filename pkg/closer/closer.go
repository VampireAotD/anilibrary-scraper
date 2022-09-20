package closer

import (
	"sync"

	"go.uber.org/zap"
)

type Closer struct {
	scope    string
	callback func() error
}

type Closers struct {
	mutex   sync.Mutex
	closers []Closer
}

func (c *Closers) Add(scope string, callback func() error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	closer := Closer{
		scope:    scope,
		callback: callback,
	}

	c.closers = append(c.closers, closer)
}

func (c *Closers) Close(logger *zap.Logger) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.closers {
		closer := c.closers[len(c.closers)-1-i]

		logger.Info("Closing", zap.String("instance", closer.scope))
		if err := closer.callback(); err != nil {
			logger.Error("close", zap.String("scope", closer.scope), zap.Error(err))
		}
	}

	c.closers = nil
}
