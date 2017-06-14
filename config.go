package logging

import "errors"
import "fmt"
import "strconv"

type LogConfig struct {
	fileDir      string
	fileName     string
	formatString string
	logLevel     LogLevel
}

// splitType    SplitType
// mu           *sync.Mutex
// maxFileSize  int
// when         string
func GetBasicConfig() (config LogConfig) {
	config = LogConfig{
		fileDir:      ".",
		fileName:     "",
		formatString: "%(dateTime),%(nanoSecond) - [%(fileName) %(lineNo)] %(levelName) %(message)",
		logLevel:     DEBUG,
	}
	return config
}

type Config struct {
	Handlers map[string]map[string]string
	Loggers  map[string][]string
}

var handlerConfig = Config{
	Handlers: map[string]map[string]string{
		"BasicHandler": map[string]string{
			"handlerType":  "BasicHandler",
			"fileDir":      ".",
			"fileName":     "",
			"formatString": "%(dateTime),%(nanoSecond) - [%(fileName) %(lineNo)] %(levelName) %(message)",
			"logLevel":     "DEBUG",
		},
		"RotatingHandler": map[string]string{
			"handlerType":  "RotatingHandler",
			"fileDir":      ".",
			"fileName":     "RotatingHandler.log",
			"formatString": "%(dateTime),%(nanoSecond) - [%(fileName) %(lineNo)] %(levelName) %(message)",
			"logLevel":     "DEBUG",
			"maxFileSize":  "104857600", //1*100*1024*1024
			"backupCount":  "30",
		},
		"TimeRotatingHandler": map[string]string{
			"handlerType":  "TimeRotatingHandler",
			"fileDir":      ".",
			"fileName":     "TimeRotatingHandler.log",
			"formatString": "%(dateTime),%(nanoSecond) - [%(fileName) %(lineNo)] %(levelName) %(message)",
			"logLevel":     "DEBUG",
			"when":         "1d",
			"backupCount":  "30",
		},
	},
	Loggers: map[string][]string{
		"BasicLogger":        []string{"BasicHandler"},
		"RotatingLogger":     []string{"RotatingHandler"},
		"TimeRotatingLogger": []string{"TimeRotatingHandler"},
		"MixedLogger":        []string{"BasicHandler", "RotatingHandler", "TimeRotatingHandler"},
	},
}

var BasicHandlerConfig = map[string]string{
	"handlerType":  "BasicHandler",
	"fileDir":      ".",
	"fileName":     "",
	"formatString": "%(dateTime),%(nanoSecond) - [%(fileName) %(lineNo)] %(levelName) %(message)",
	"logLevel":     "DEBUG",
}

var RotatingHandlerConfig = map[string]string{
	"handlerType":  "RotatingHandler",
	"fileDir":      ".",
	"fileName":     "RotatingHandler.log",
	"formatString": "%(dateTime),%(nanoSecond) - [%(fileName) %(lineNo)] %(levelName) %(message)",
	"logLevel":     "DEBUG",
	"maxFileSize":  "104857600", //1*100*1024*1024
	"backupCount":  "30",
}

var TimeRotatingHandlerConfig = map[string]string{
	"handlerType":  "TimeRotatingHandler",
	"fileDir":      ".",
	"fileName":     "TimeRotatingHandler.log",
	"formatString": "%(dateTime),%(nanoSecond) - [%(fileName) %(lineNo)] %(levelName) %(message)",
	"logLevel":     "DEBUG",
	"when":         "1d",
	"backupCount":  "30",
}

func getBasicHandler(conf map[string]string) (handler1 LogHandler, err error) {
	handler, err := GetBasicHandler(conf["fileDir"], conf["fileName"])
	if err != nil {
		return
	}
	if formatString, ok := conf["formatString"]; ok {
		err = handler.SetFormatString(formatString)
		if err != nil {
			return
		}
	}
	if logLevel, ok := conf["logLevel"]; ok {
		switch logLevel {
		case "DEBUG":
			err = handler.SetLogLevel(DEBUG)
		case "WARNING":
			err = handler.SetLogLevel(WARNING)
		case "ERROR":
			err = handler.SetLogLevel(ERROR)
		default:
			err = errors.New(fmt.Sprintf("err format or logLevel %s", logLevel))
		}
		if err != nil {
			return
		}
	}
	handler1 = handler
	return
}

func getRotatingHandler(conf map[string]string) (handler1 LogHandler, err error) {
	handler, err := GetRotatingHandler(conf["fileDir"], conf["fileName"])
	if err != nil {
		return
	}
	if formatString, ok := conf["formatString"]; ok {
		err = handler.SetFormatString(formatString)
		if err != nil {
			return
		}
	}
	if logLevel, ok := conf["logLevel"]; ok {
		switch logLevel {
		case "DEBUG":
			err = handler.SetLogLevel(DEBUG)
		case "WARNING":
			err = handler.SetLogLevel(WARNING)
		case "ERROR":
			err = handler.SetLogLevel(ERROR)
		default:
			err = errors.New(fmt.Sprintf("err format or logLevel %s", logLevel))
		}
		if err != nil {
			return
		}
	}
	if maxFileSize, ok := conf["maxFileSize"]; ok {
		size, err1 := strconv.ParseInt(maxFileSize, 10, 64)
		if err1 != nil {
			err = err1
			return
		}
		err = handler.SetMaxFileSize(size)
		if err != nil {
			return
		}

	}
	if backupCount, ok := conf["backupCount"]; ok {
		count, err1 := strconv.Atoi(backupCount)
		if err1 != nil {
			err = err1
			return
		}
		err = handler.SetBackupCount(count)
		if err != nil {
			return
		}
	}
	handler1 = handler
	return
}

func getTimeRotatingHandler(conf map[string]string) (handler1 LogHandler, err error) {
	handler, err := GetTimeRotatingHandler(conf["fileDir"], conf["fileName"])
	if err != nil {
		return
	}
	if formatString, ok := conf["formatString"]; ok {
		err = handler.SetFormatString(formatString)
		if err != nil {
			return
		}
	}
	if logLevel, ok := conf["logLevel"]; ok {
		switch logLevel {
		case "DEBUG":
			err = handler.SetLogLevel(DEBUG)
		case "WARNING":
			err = handler.SetLogLevel(WARNING)
		case "ERROR":
			err = handler.SetLogLevel(ERROR)
		default:
			err = errors.New(fmt.Sprintf("err format or logLevel %s", logLevel))
		}
		if err != nil {
			return
		}
	}
	if when, ok := conf["when"]; ok {
		err = handler.SetWhen(when)
		if err != nil {
			return
		}
	}
	if backupCount, ok := conf["backupCount"]; ok {
		count, err1 := strconv.Atoi(backupCount)
		if err1 != nil {
			err = err1
			return
		}
		err = handler.SetBackupCount(count)
		if err != nil {
			return
		}
	}
	handler1 = handler
	return
}

func getHandler(conf map[string]string) (handler LogHandler, err error) {
	switch conf["handlerType"] {
	case "BasicHandler":
		return getBasicHandler(conf)
	case "RotatingHandler":
		return getRotatingHandler(conf)
	case "TimeRotatingHandler":
		return getTimeRotatingHandler(conf)
	default:
		return nil, errors.New(fmt.Sprintf("err format of handlerType %s", conf["handlerType"]))
	}
	return
}

func MapConfig(config Config) (err error) {
	handlers := map[string]LogHandler{}
	for k, v := range config.Handlers {
		handler, err1 := getHandler(v)
		if err1 != nil {
			err = err1
			return
		}
		handlers[k] = handler
	}
	for k, v := range config.Loggers {
		logger := GetLogger(k)
		for _, handlerName := range v {
			if handler, ok := handlers[handlerName]; ok {
				logger.AddHandler(handler)
			} else {
				err = errors.New(fmt.Sprintf("handlerName:%s not exists", handlerName))
				return
			}
		}
	}
	return
}
