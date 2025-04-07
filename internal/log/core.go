package log

import "github.com/withmandala/go-log"

type CoreLogger struct {
	logger *log.Logger
}

func (c *CoreLogger) Init() {
	c.logger = InitLogger()
}

func (c *CoreLogger) Info(message ...interface{}) {
	c.logger.Info(message)
}

func (c *CoreLogger) Warn(message ...interface{}) {
	c.logger.Warn(message)
}

func (c *CoreLogger) Error(message ...interface{}) {
	c.logger.Error(message)
}

func (c *CoreLogger) Debug(message ...interface{}) {
	c.logger.Debug(message)
}

func (c *CoreLogger) Trace(message ...interface{}) {
	c.logger.Trace(message)
}

func (c *CoreLogger) Fatal(message ...interface{}) {
	c.logger.Fatal(message)
}
