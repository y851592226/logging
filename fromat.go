package logging

import "fmt"
import "runtime"
import "path"
import "path/filepath"
import "strconv"
import "time"
import "errors"

type formatter struct {
	name       string // Name of the logger
	levelName  string
	pathName   string
	fileName   string
	funcName   string
	lineNo     string
	date       string
	unixTime   string
	dateTime   string
	weekday    string
	nanoSecond string
	ascTime    string
	message    string
}

func baseFormatter() (fm formatter) {
	fm = formatter{
		name:       "!@#",
		levelName:  "!@#",
		pathName:   "!@#",
		fileName:   "!@#",
		funcName:   "!@#",
		lineNo:     "!@#",
		date:       "!@#",
		unixTime:   "!@#",
		dateTime:   "!@#",
		weekday:    "!@#",
		nanoSecond: "!@#",
		ascTime:    "!@#",
		message:    "!@#",
	}
	return
}

func (fm *formatter) getTimeInfo() {
	t := time.Now()
	fm.ascTime = t.Format("2006-01-02 15:04:05,") + fmt.Sprintf("%09s", strconv.Itoa(t.Nanosecond()))
	fm.date = fm.ascTime[:10]
	fm.dateTime = fm.ascTime[:19]
	fm.nanoSecond = fm.ascTime[20:]
	fm.unixTime = strconv.FormatInt(t.Unix(), 10)
	fm.weekday = t.Weekday().String()
}

func (fm *formatter) getDate() string {
	if fm.date == "!@#" {
		fm.getTimeInfo()
	}
	return fm.date
}

func (fm *formatter) getDateTime() string {
	if fm.dateTime == "!@#" {
		fm.getTimeInfo()
	}
	return fm.dateTime
}

func (fm *formatter) getNanoSecond() string {
	if fm.nanoSecond == "!@#" {
		fm.getTimeInfo()
	}
	return fm.nanoSecond
}

func (fm *formatter) getAscTime() string {
	if fm.ascTime == "!@#" {
		fm.getTimeInfo()
	}
	return fm.ascTime
}

func (fm *formatter) getUnixTime() string {
	if fm.unixTime == "!@#" {
		fm.getTimeInfo()
	}
	return fm.unixTime
}

func (fm *formatter) getWeekday() string {
	if fm.weekday == "!@#" {
		fm.getTimeInfo()
	}
	return fm.weekday
}

func (fm *formatter) getPathName() string {
	if fm.pathName == "!@#" {
		fm.getCallInfo()

	}
	return fm.pathName
}

func (fm *formatter) getFileName() string {
	if fm.fileName == "!@#" {
		fm.getCallInfo()

	}
	return fm.fileName
}

func (fm *formatter) getFuncName() string {
	if fm.funcName == "!@#" {
		fm.getCallInfo()

	}
	return fm.funcName
}

func (fm *formatter) getLevelName() string {
	return fm.levelName
}

func (fm *formatter) setLevelName(logLevel LogLevel) {
	switch logLevel {
	case DEBUG:
		fm.levelName = "DEBUG"
	case WARNING:
		fm.levelName = "WARNING"
	case ERROR:
		fm.levelName = "ERROR"
	default:
		fm.levelName = "ERROR"
	}
}

func (fm *formatter) getName() string {
	return fm.name
}

func (fm *formatter) setName(name string) {
	fm.name = name
}

func (fm *formatter) getCallInfo() {
	pc, file, line, ok := runtime.Caller(5)
	if ok {
		fm.pathName = file
		fm.fileName = filepath.Base(file)
		fm.lineNo = strconv.Itoa(line)
		fm.funcName = path.Base(runtime.FuncForPC(pc).Name())
	} else {
		fm.pathName = "???"
		fm.fileName = "???"
		fm.lineNo = "???"
		fm.funcName = "???"
	}
}

func (fm *formatter) getLineNo() string {
	if fm.lineNo == "!@#" {
		fm.getCallInfo()
	}
	return fm.lineNo
}

func (fm *formatter) getMessage() string {
	return fm.message
}

func (fm *formatter) setMessage(message string) {
	fm.message = message
}

func (fm *formatter) getNil() string {
	return "error format had been setted"
}

func (handler *BasicHandler) setFormatFunc(format string) (err error) {
	switch format {
	case "name":
		handler.formatFunc = append(handler.formatFunc, handler.getName)
	case "levelName":
		handler.formatFunc = append(handler.formatFunc, handler.getLevelName)
	case "pathName":
		handler.formatFunc = append(handler.formatFunc, handler.getPathName)
	case "fileName":
		handler.formatFunc = append(handler.formatFunc, handler.getFileName)
	case "funcName":
		handler.formatFunc = append(handler.formatFunc, handler.getFuncName)
	case "lineNo":
		handler.formatFunc = append(handler.formatFunc, handler.getLineNo)
	case "date":
		handler.formatFunc = append(handler.formatFunc, handler.getDate)
	case "unixTime":
		handler.formatFunc = append(handler.formatFunc, handler.getUnixTime)
	case "nanoSecond":
		handler.formatFunc = append(handler.formatFunc, handler.getNanoSecond)
	case "ascTime":
		handler.formatFunc = append(handler.formatFunc, handler.getAscTime)
	case "dateTime":
		handler.formatFunc = append(handler.formatFunc, handler.getDateTime)
	case "weekday":
		handler.formatFunc = append(handler.formatFunc, handler.getWeekday)
	case "message":
		handler.formatFunc = append(handler.formatFunc, handler.getMessage)
	default:
		err = errors.New("error formatName %(" + format + ")")
	}
	return
}

func (handler *BasicHandler) setFormatter() (err error) {
	formatString := handler.logConfig.formatString
	begin := false
	format := []byte{}
	for i := 0; i < len(formatString); i++ {
		if begin {
			if formatString[i] == ')' {
				begin = false
				err = handler.setFormatFunc(string(format[1:]))
				if err != nil {
					return
				}
			} else {
				format = append(format, formatString[i])
			}
		} else if formatString[i] == '%' {
			if i+1 < len(formatString) && formatString[i+1] == '(' {
				begin = true
				format = format[:0]
				handler.formatString += "%s"
			}

		} else {
			handler.formatString += string(formatString[i])
		}
	}
	handler.formatString += "\n"
	if begin {
		handler.formatString = ""
		err = errors.New("error format \"" + handler.logConfig.formatString + "\"")
	}
	return
}
