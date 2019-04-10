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
    "github.com/lunny/log"
    "os"
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
    LOG_TAG[DEBUG] = "Debug"
    LOG_TAG[INFO] = "Info"
    LOG_TAG[WARN] = "Warn"
    LOG_TAG[ERROR] = "Error"
    LOG_TAG[FATAL] = "Fatal"
}

type LogFunc func(level int, format string, args ...interface{})

var Log LogFunc = DefaultLogf

func DefaultLogf(level int, format string, args ...interface{}) {
    logInfo := fmt.Sprintf(format, args)
    log.Printf("%s: %s", LOG_TAG[level], logInfo)
    if level >= FATAL {
        os.Exit(-1)
    }
}

func Debug(format string, args ...interface{}) {
    Log(DEBUG, format, args)
}

func Info(format string, args ...interface{}) {
    Log(DEBUG, format, args)
}
func Warn(format string, args ...interface{}) {
    Log(DEBUG, format, args)
}
func Error(format string, args ...interface{}) {
    Log(DEBUG, format, args)
}
func Fatal(format string, args ...interface{}) {
    Log(DEBUG, format, args)
}
