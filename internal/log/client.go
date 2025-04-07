package log

import "github.com/withmandala/go-log"

type ClientLogger struct {
	logger *log.Logger
}

func (c *ClientLogger) Init() {
	c.logger = InitLogger()
}

func (c *ClientLogger) Info(message ...interface{}) {
	c.logger.Info(message)
}

func (c *ClientLogger) Warn(message ...interface{}) {
	c.logger.Warn(message)
}

func (c *ClientLogger) Error(message ...interface{}) {
	c.logger.Error(message)
}

func (c *ClientLogger) Debug(message ...interface{}) {
	c.logger.Debug(message)
}

func (c *ClientLogger) Trace(message ...interface{}) {
	c.logger.Trace(message)
}

func (c *ClientLogger) Fatal(message ...interface{}) {
	c.logger.Fatal(message)
}
