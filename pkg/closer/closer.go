package closer

import (
	"sync"

	"anilibrary-scraper/pkg/logger"
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

func (c *Closers) Close(log logger.Contract) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for i := range c.closers {
		closer := c.closers[len(c.closers)-1-i]

		log.Info("Closing", logger.String("instance", closer.scope))
		if err := closer.callback(); err != nil {
			log.Error("close", logger.String("scope", closer.scope), logger.Error(err))
		}
	}

	c.closers = nil
}
