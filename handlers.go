package logging

import "fmt"
import "sync"
import "os"
import "path"
import "io"
import "errors"
import "strconv"
import "syscall"
import "time"
import "path/filepath"
import "sort"
import "regexp"
import "reflect"

type LogHandler interface {
	writeLog(name string, logLevel LogLevel, format string, v ...interface{})
	getLogLevel() LogLevel
	getId() int
    Close()
}

type BasicHandler struct {
	mu        *sync.Mutex
	logConfig *LogConfig
	out       io.ReadWriteCloser
	id        int
	formatter
	formatString string
	formatFunc   []func() string
}

func GetBasicHandler(fileDir, fileName string) (basicHandler *BasicHandler, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	basicHandler = new(BasicHandler)
	logConfig := GetBasicConfig()
	logConfig.fileName = fileName
	logConfig.fileDir = fileDir
	basicHandler.logConfig = &logConfig
	basicHandler.mu = new(sync.Mutex)
	basicHandler.id = handlerId
	handlerId++
	basicHandler.out = os.Stdout
	err = basicHandler.setOut()
	if err != nil {
		return
	}
	basicHandler.formatter = baseFormatter()
	basicHandler.setFormatter()
	return
}

func (handler *BasicHandler) setOut() (err error) {
	if handler.out != os.Stdout {
		handler.out.Close()
	}
	if handler.logConfig.fileDir == "" {
		handler.logConfig.fileDir = "."
	}
	if handler.logConfig.fileName == "" {
		handler.out = os.Stdout
	} else {
		if handler.logConfig.fileName != path.Base(handler.logConfig.fileName) {
			err = errors.New(handler.logConfig.fileName + " is not a file name")
			return
		}
		filepath := path.Join(handler.logConfig.fileDir, handler.logConfig.fileName)
		handler.out, err = os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	}
	return
}

func (handler *BasicHandler) SetFormatString(format string) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.formatter = baseFormatter()
	handler.logConfig.formatString = format
	handler.formatString = ""
	handler.formatFunc = []func() string{}
	err = handler.setFormatter()
	return
}

func (handler *BasicHandler) GetFormatString() string {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	return handler.logConfig.formatString
}

func (handler *BasicHandler) SetFilePath(fileDir, fileName string) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.logConfig.fileDir = fileDir
	handler.logConfig.fileName = fileName
	err = handler.setOut()
	return
}

func (handler *BasicHandler) GetFilePath() string {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	return path.Join(handler.logConfig.fileDir, handler.logConfig.fileName)
}

func (handler *BasicHandler) SetLogLevel(logLevel LogLevel) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.logConfig.logLevel = logLevel
	return
}

func (handler *BasicHandler) getLogLevel() LogLevel {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	return handler.logConfig.logLevel
}

func (handler *BasicHandler) Close() {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	if handler.out != os.Stdout {
		handler.out.Close()
	}
}

func (handler *BasicHandler) getId() int {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	return handler.id
}

func (handler *BasicHandler) writeLog(name string, logLevel LogLevel, format string, v ...interface{}) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.formatter = baseFormatter()
	handler.formatter.setName(name)
	handler.formatter.setLevelName(logLevel)
	handler.formatter.setMessage(fmt.Sprintf(format, v...))
	value := []interface{}{}
	for _, fun := range handler.formatFunc {
		value = append(value, fun())
	}
	fmt.Fprintf(handler.out, handler.formatString, value...)
}

type RotatingHandler struct {
	BasicHandler
	splitType       SplitType
	maxFileSize     int64
	backupCount     int
	currentFileSize int64
}

func GetRotatingHandler(fileDir, fileName string) (rotatingHandler *RotatingHandler, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	rotatingHandler = new(RotatingHandler)
	logConfig := GetBasicConfig()
	logConfig.fileName = fileName
	logConfig.fileDir = fileDir
	rotatingHandler.maxFileSize = 1 * 100 * 1024 * 1024
	rotatingHandler.splitType = SplitBySize
	rotatingHandler.backupCount = 30
	rotatingHandler.logConfig = &logConfig
	rotatingHandler.mu = new(sync.Mutex)
	rotatingHandler.id = handlerId
	handlerId++
	rotatingHandler.out = os.Stdout
	err = rotatingHandler.setOut()
	if err != nil {
		return
	}
	rotatingHandler.formatter = baseFormatter()
	rotatingHandler.setFormatter()
	return
}

func (handler *RotatingHandler) setOut() (err error) {
	if handler.out != os.Stdout {
		handler.out.Close()
	}
	if handler.logConfig.fileDir == "" {
		handler.logConfig.fileDir = "."
	}
	if handler.logConfig.fileName == "" {
		err = errors.New("fileName has not been set")
		return
	} else {
		if handler.logConfig.fileName != path.Base(handler.logConfig.fileName) {
			err = errors.New(handler.logConfig.fileName + " is not a file name")
			return
		}
		filepath := path.Join(handler.logConfig.fileDir, handler.logConfig.fileName)
		handler.out, err = os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
		stat, err1 := os.Stat(filepath)
		if err1 != nil {
			err = err1
			return
		}
		handler.currentFileSize = stat.Size()
	}
	return
}

func (handler *RotatingHandler) SetFilePath(fileDir, fileName string) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.logConfig.fileDir = fileDir
	handler.logConfig.fileName = fileName
	err = handler.setOut()
	return
}

func (handler *RotatingHandler) SetMaxFileSize(size int64) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	if size <= 0 {
		err = errors.New("size must be a positive number")
	}
	handler.maxFileSize = size
	return
}

func (handler *RotatingHandler) SetBackupCount(count int) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	if count < 0 {
		err = errors.New("count can't be a negative number")
	}
	handler.backupCount = count
	return
}

func (handler *RotatingHandler) writeLog(name string, logLevel LogLevel, format string, v ...interface{}) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.formatter = baseFormatter()
	handler.formatter.setName(name)
	handler.formatter.setLevelName(logLevel)
	handler.formatter.setMessage(fmt.Sprintf(format, v...))
	value := []interface{}{}
	for _, fun := range handler.formatFunc {
		value = append(value, fun())
	}
	s := fmt.Sprintf(handler.formatString, value...)
	if int64(len(s))+handler.currentFileSize > handler.maxFileSize {
		handler.doRorate()
	}
	handler.currentFileSize += int64(len(s))
	io.WriteString(handler.out, s)
}

func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getMaxLogNum(fileName string) (num int) {
	num = 1
	for ; ; num++ {
		result, _ := IsPathExists(fileName + "." + strconv.Itoa(num))
		if !result {
			break
		}
	}
	return
}

func min(a, b int) int {
	if a < b && a != 0 {
		return a
	} else {
		return b
	}
}

func (handler *RotatingHandler) doRorate() {
	handler.out.Close()
	filepath := path.Join(handler.logConfig.fileDir, handler.logConfig.fileName)
	for i := min(handler.backupCount, getMaxLogNum(filepath)); i >= 1; i-- {
		sfn := ""
		if i == 1 {
			sfn = filepath
		} else {
			sfn = filepath + "." + strconv.Itoa(i-1)
		}
		dfn := filepath + "." + strconv.Itoa(i)
		exist, _ := IsPathExists(sfn)
		if exist {
			os.Rename(sfn, dfn)
		}
	}
	handler.out, _ = os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	stat, _ := os.Stat(filepath)
	handler.currentFileSize = stat.Size()
}

type TimeRotatingHandler struct {
	BasicHandler
	splitType   SplitType
	backupCount int
	when        string
	createTime  time.Time
	rotateTime  time.Time
	fileTag     string
}

func GetTimeRotatingHandler(fileDir, fileName string) (timerotatingHandler *TimeRotatingHandler, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	timerotatingHandler = new(TimeRotatingHandler)
	logConfig := GetBasicConfig()
	logConfig.fileName = fileName
	logConfig.fileDir = fileDir
	timerotatingHandler.splitType = SplitByTime
	timerotatingHandler.backupCount = 30
	timerotatingHandler.when = "1d"
	timerotatingHandler.logConfig = &logConfig
	timerotatingHandler.mu = new(sync.Mutex)
	timerotatingHandler.id = handlerId
	handlerId++
	timerotatingHandler.out = os.Stdout
	err = timerotatingHandler.setOut()
	if err != nil {
		return
	}
	timerotatingHandler.formatter = baseFormatter()
	timerotatingHandler.setFormatter()
	timerotatingHandler.rotateTime = getRotateTime(timerotatingHandler.createTime, timerotatingHandler.when)
	timerotatingHandler.fileTag = getFileTag(timerotatingHandler.createTime, timerotatingHandler.when)
	return
}

func (handler *TimeRotatingHandler) SetBackupCount(count int) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	if count < 0 {
		err = errors.New("count can't be a negative number")
	}
	handler.backupCount = count
	return
}

func (handler *TimeRotatingHandler) SetWhen(when string) (err error) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	reg := regexp.MustCompile(`^\d+[shd]$`)
	if !reg.MatchString(when) {
		err = errors.New("error format of when:" + when)
		return
	}
	handler.when = when
	handler.rotateTime = getRotateTime(handler.createTime, handler.when)
	handler.fileTag = getFileTag(handler.createTime, handler.when)
	return
}

func (handler *TimeRotatingHandler) setOut() (err error) {
	if handler.out != os.Stdout {
		handler.out.Close()
	}
	if handler.logConfig.fileDir == "" {
		handler.logConfig.fileDir = "."
	}
	if handler.logConfig.fileName == "" {
		err = errors.New("fileName has not been set")
		return
	} else {
		if handler.logConfig.fileName != path.Base(handler.logConfig.fileName) {
			err = errors.New(handler.logConfig.fileName + " is not a file name")
			return
		}
		filepath := path.Join(handler.logConfig.fileDir, handler.logConfig.fileName)
		handler.out, err = os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		handler.createTime,err = GetBirthtime(filepath)
		if err != nil {
			return
		}
	}
	return
}

func getField(stat syscall.Stat_t,name string) (result syscall.Timespec,ok bool){
    t := reflect.TypeOf(stat)
    _,ok = t.FieldByName(name)
    if !ok{
        return
    }
    v := reflect.ValueOf(stat)
    v = v.FieldByName(name)
    result = (v.Interface()).(syscall.Timespec)
    return
}

func GetBirthtime(fileName string)(t time.Time, err error){
    fileInfo, err := os.Stat(fileName)
    if err!=nil{
        return
    }
    stat := fileInfo.Sys().(*syscall.Stat_t)
    timeName :=[]string{"Birthtimespec","Ctim","Ctimespec"}
    for _,name := range timeName{
        result,ok:=getField(*stat,name)
        if ok{
            t=time.Unix(result.Sec,0)
            return
        }
    }
    t=time.Now()
    return
}

func (handler *TimeRotatingHandler) writeLog(name string, logLevel LogLevel, format string, v ...interface{}) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.formatter = baseFormatter()
	handler.formatter.setName(name)
	handler.formatter.setLevelName(logLevel)
	handler.formatter.setMessage(fmt.Sprintf(format, v...))
	value := []interface{}{}
	for _, fun := range handler.formatFunc {
		value = append(value, fun())
	}
	s := fmt.Sprintf(handler.formatString, value...)
	if handler.checkRorate() {
		handler.doRorate()
	}
	io.WriteString(handler.out, s)
}

func getFileTag(begin time.Time, when string) (fileTag string) {
	timeType := when[len(when)-1:]
	switch timeType {
	case "s":
		fileTag = begin.Format("2006-01-02 15:04:05")
	case "h":
		fileTag = begin.Format("2006-01-02 15")
	case "d":
		fileTag = begin.Format("2006-01-02")
	default:
		fileTag = begin.Format("2006-01-02")
	}
	return fileTag
}

func getRotateTime(begin time.Time, when string) (end time.Time) {
	timeType := when[len(when)-1:]
	_value, _ := strconv.ParseInt(when[:len(when)-1], 10, 64)
	value := time.Duration(_value)
	switch timeType {
	case "s":
		end = begin.Add(value * time.Second)
	case "h":
		begin, _ = time.Parse("2006-01-02 03", begin.Format("2006-01-02 03"))
		end = begin.Add(value * time.Hour)
	case "d":
		begin, _ = time.Parse("2006-01-02", begin.Format("2006-01-02"))
		end = begin.Add(value * 24 * time.Hour)
	default:
		begin, _ = time.Parse("2006-01-02", begin.Format("2006-01-02"))
		end = begin.Add(value * 24 * time.Hour)
	}
	return end
}

func (handler *TimeRotatingHandler) checkRorate() (result bool) {
	t := time.Now()
	result = false
	if t.After(handler.rotateTime) {
		result = true
	}
	return
}

func WalkDir(dirPth, prefix string, when string) (files []string, err error) {
	if dirPth == "" {
		dirPth = "."
	}
	restring := `^$`
	switch when[len(when)-1:] {
	case "s":
		restring = `^` + prefix + `\.` + `\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\d` + `$`
	case "h":
		restring = `^` + prefix + `\.` + `\d\d\d\d-\d\d-\d\d \d\d` + `$`
	case "d":
		restring = `^` + prefix + `\.` + `\d\d\d\d-\d\d-\d\d` + `$`
	default:
		restring = `^$`
	}
	reg := regexp.MustCompile(restring)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if reg.MatchString(filename) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func (handler *TimeRotatingHandler) doRorate() {
	handler.out.Close()
	sfn := path.Join(handler.logConfig.fileDir, handler.logConfig.fileName)
	dfn := path.Join(handler.logConfig.fileDir, handler.logConfig.fileName+"."+handler.fileTag)
	os.Rename(sfn, dfn)
	files, _ := WalkDir(handler.logConfig.fileDir, handler.logConfig.fileName, handler.when)
	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	for i := range files {
		if i >= handler.backupCount {
			os.Remove(files[i])
		}
	}
	handler.out, _ = os.OpenFile(sfn, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	handler.createTime = time.Now()
	handler.rotateTime = getRotateTime(handler.createTime, handler.when)
	handler.fileTag = getFileTag(handler.createTime, handler.when)
}
