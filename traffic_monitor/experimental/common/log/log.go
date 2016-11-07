// Inspired by https://www.goinggo.net/2013/11/using-log-package-in-go.html
package log

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


import (
	"fmt"
	"io"
	"log"
	"time"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(debugW, infoW, warnW, errW io.Writer) {
	Debug = log.New(debugW, "DEBUG: ", log.Lshortfile)
	Info = log.New(infoW, "INFO: ", log.Lshortfile)
	Warning = log.New(warnW, "WARNING: ", log.Lshortfile)
	Error = log.New(errW, "ERROR: ", log.Lshortfile)
}

const timeFormat = time.RFC3339Nano

func Errorf(format string, v ...interface{}) {
	Error.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Errorln(v ...interface{}) {
	Error.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}
func Warnf(format string, v ...interface{}) {
	Warning.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Warnln(v ...interface{}) {
	Warning.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}
func Infof(format string, v ...interface{}) {
	Info.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Infoln(v ...interface{}) {
	Info.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}
func Debugf(format string, v ...interface{}) {
	Debug.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}
func Debugln(v ...interface{}) {
	Debug.Output(3, time.Now().Format(timeFormat)+": "+fmt.Sprintln(v...))
}
