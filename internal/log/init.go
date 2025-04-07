package log

var coreLogger *CoreLogger
var clientLogger *ClientLogger

func initLoggerSystem() {
	coreLogger = &CoreLogger{}
	coreLogger.Init()

	clientLogger = &ClientLogger{}
	clientLogger.Init()
}

func GetCoreInstance() *CoreLogger {
	return coreLogger
}

func GetClientInstance() *ClientLogger {
	return clientLogger
}

func init() {
	initLoggerSystem()
}
