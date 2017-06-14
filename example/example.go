package main

import "github.com/y851592226/logging"

var handlerConfig = logging.Config{
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

func exampleMapConfig(){
   err := logging.MapConfig(handlerConfig)
    if err!=nil{
        panic(err.Error())
    }
    log := logging.GetLogger("MixedLogger")
    defer log.Close()
    log.Debug("%s","exampleMapConfig")
}

func exampleBasicHandler(){
    basicHandler, err := logging.GetBasicHandler(".", "")
    if err != nil {
        panic(err.Error())
    }
    defer basicHandler.Close()
    log := logging.GetLogger("log1")
    err = basicHandler.SetLogLevel(logging.WARNING)
    if err != nil {
        panic(err.Error())
    }
    err = basicHandler.SetFormatString("%(dateTime) - [%(fileName) <%(funcName)> %(lineNo)] %(levelName) %(message)")
    if err != nil {
        panic(err.Error())
    }
    err = basicHandler.SetFilePath(".","log1.log")
    if err != nil {
        panic(err.Error())
    }
    log.AddHandler(basicHandler)
    log.Debug("%s", "this is the first log")
    log.Warning("%s", "this is the second log")
    log.Error("%s", "this is the third log")
    log.RemoveHandler(basicHandler)
}

func exampleRotatingHandler() {
    rotatingHandler, err := logging.GetRotatingHandler(".", "log2.log")
    if err != nil {
        panic(err.Error())
    }
    defer rotatingHandler.Close()
    //设置最大文件大小100M
    err=rotatingHandler.SetMaxFileSize(1*100*1024*1024)
    if err != nil {
        panic(err.Error())
    }
    //设置最多备份30个日志文件，多余的日志文件会自动被删除 0代表只进行日志切分，不删除日志文件
    err=rotatingHandler.SetBackupCount(30)
    if err != nil {
        panic(err.Error())
    }
    log := logging.GetLogger("log2")
    log.AddHandler(rotatingHandler)
    log.Debug("%s", "this is the first log")
    log.Warning("%s", "this is the second log")
    log.Error("%s", "this is the third log")
}

func exampleTimeRotatingHandler() {
    timeRotatingHandler, err := logging.GetTimeRotatingHandler(".", "log3.log")
    if err != nil {
        panic(err.Error())
    }
    defer timeRotatingHandler.Close()
    //设置日志文件切分时间 例如 1000s、2h、1d  按照小时会在整点时间进行切分 按照天会在0点进行切分
    err=timeRotatingHandler.SetWhen("1d")
    if err != nil {
        panic(err.Error())
    }
    //设置最多备份30个日志文件，最早的日志文件会自动被删除 0代表只进行日志切分，不删除日志文件
    err=timeRotatingHandler.SetBackupCount(30)
    if err != nil {
        panic(err.Error())
    }
    log := logging.GetLogger("log3")
    log.AddHandler(timeRotatingHandler)
    log.Debug("%s", "this is the first log")
    log.Warning("%s", "this is the second log")
    log.Error("%s", "this is the third log")
}

func main() {
    exampleMapConfig()
    exampleBasicHandler()
    exampleRotatingHandler()
    exampleTimeRotatingHandler()
}
