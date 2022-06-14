/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var gLogTag map[int]string = map[int]string{}
var gLogLevel = INFO

func init() {
	gLogTag[DEBUG] = "[Debug]"
	gLogTag[INFO] = "[Info]"
	gLogTag[WARN] = "[Warn]"
	gLogTag[ERROR] = "[Error]"
	gLogTag[FATAL] = "[Fatal]"
}

type LogFunc func(level int, format string, args ...interface{})

var Log LogFunc = DefaultLogf

func SetLevel(level int) {
	gLogLevel = level
}

func DummyLog(level int, format string, args ...interface{}) {

}

func DefaultLogf(level int, format string, args ...interface{}) {
	if gLogLevel > level {
		return
	}

	logInfo := fmt.Sprintf(format, args...)
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	log.Printf("%s %s:%d %s", gLogTag[level], shortFile(file), line, logInfo)
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
