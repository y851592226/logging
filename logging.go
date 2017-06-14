package logging

import "sync"

type SplitType int

const (
	SplitBySize SplitType = iota
	SplitByTime
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	WARNING
	ERROR
)

type FileLogger struct {
	name       string
	mu         *sync.Mutex
	logHandler []*LogHandler
}

var globalLogMap = make(map[string]*FileLogger)
var mutex = new(sync.Mutex)
var handlerId = 1

func (fl *FileLogger) Debug(format string, v ...interface{}) {
	for _, value := range fl.logHandler {
		if (*value).getLogLevel() <= DEBUG {
			(*value).writeLog(fl.name, DEBUG, format, v...)
		}
	}
}

func (fl *FileLogger) Warning(format string, v ...interface{}) {
	for _, value := range fl.logHandler {
		if (*value).getLogLevel() <= WARNING {
			(*value).writeLog(fl.name, WARNING, format, v...)
		}
	}
}

func (fl *FileLogger) Error(format string, v ...interface{}) {
	for _, value := range fl.logHandler {
		if (*value).getLogLevel() <= ERROR {
			(*value).writeLog(fl.name, ERROR, format, v...)
		}
	}
}

func (fl *FileLogger) Close() {
	for _, value := range fl.logHandler {
		(*value).Close()
	}
}

func (fl *FileLogger) AddHandler(handler LogHandler) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.logHandler = append(fl.logHandler, &handler)
}

func (fl *FileLogger) RemoveHandler(handler LogHandler) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	logHandler := []*LogHandler{}
	for _, _handler := range fl.logHandler {
		if (*_handler).getId() != handler.getId() {
			logHandler = append(logHandler, _handler)
		}
	}
	fl.logHandler = logHandler
}

func GetLogger(logname string) (logger *FileLogger) {
	mutex.Lock()
	defer mutex.Unlock()
	logger, ok := globalLogMap[logname]
	if !ok {
		logger = &FileLogger{name: logname, mu: new(sync.Mutex), logHandler: []*LogHandler{}}
		globalLogMap[logname] = logger
	}
	return
}
