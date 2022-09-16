package closer

import (
	"io"
	"sync"

	"anilibrary-request-parser/pkg/logger"
)

type Closer struct {
	scope    string
	callback io.Closer
}

type Closers struct {
	mutex   sync.Mutex
	closers []Closer
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

func (c *Closers) Close(log logger.Logger) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.closers {
		closer := c.closers[len(c.closers)-1-i]

		log.Info("Closing", logger.String("instance", closer.scope))
		if err := closer.callback.Close(); err != nil {
			log.Error("close", logger.String("scope", closer.scope), logger.Error(err))
		}
	}

	c.closers = nil
}
