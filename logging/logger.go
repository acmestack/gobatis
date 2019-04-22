/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description: 
 */

package logging

import (
    "fmt"
    "os"
    "runtime"
    "log"
)

const (
    DEBUG = iota
    INFO
    WARN
    ERROR
    FATAL
)

var LOG_TAG map[int]string = map[int]string{}

func init() {
    LOG_TAG[DEBUG] = "[Debug]"
    LOG_TAG[INFO] = "[Info]"
    LOG_TAG[WARN] = "[Warn]"
    LOG_TAG[ERROR] = "[Error]"
    LOG_TAG[FATAL] = "[Fatal]"
}

type LogFunc func(level int, format string, args ...interface{})

var Log LogFunc = DefaultLogf

func DefaultLogf(level int, format string, args ...interface{}) {
    logInfo := fmt.Sprintf(format, args...)
    var file string
    var line int
    var ok bool
    _, file, line, ok = runtime.Caller(2)
    if !ok {
        file = "???"
        line = 0
    }
    log.Printf("%s %s:%d %s", LOG_TAG[level], shortFile(file), line, logInfo)
    if level >= FATAL {
        os.Exit(-1)
    }
}

func shortFile(file string) string {
    short := file
    for i := len(file) - 1; i > 0; i-- {
        if file[i] == '/' {
            short = file[i+1:]
            break
        }
    }
    return short
}

func Debug(format string, args ...interface{}) {
    Log(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
    Log(INFO, format, args...)
}
func Warn(format string, args ...interface{}) {
    Log(WARN, format, args...)
}
func Error(format string, args ...interface{}) {
    Log(ERROR, format, args...)
}
func Fatal(format string, args ...interface{}) {
    Log(FATAL, format, args...)
}
