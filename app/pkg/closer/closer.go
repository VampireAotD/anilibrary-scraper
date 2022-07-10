package closer

import (
	"fmt"
	"io"
	"sync"

	"anilibrary-request-parser/app/pkg/logger"
	"go.uber.org/zap"
)

type Closers struct {
	logger  logger.Logger
	mutex   sync.Mutex
	closers []io.Closer
}

func New(logger logger.Logger) Closers {
	return Closers{logger: logger}
}

func (c *Closers) Add(closer io.Closer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.closers = append(c.closers, closer)
}

func (c *Closers) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, closer := range c.closers {
		if err := closer.Close(); err != nil {
			c.logger.Error(fmt.Sprintf("close %T", closer), zap.Error(err))
		}
	}

	c.closers = nil
}
