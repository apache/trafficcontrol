package main

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

 type Alert struct {
	Text string `json:"text"`
	Level string `json:"level"`
 }

 type Alerts struct {
 	Alerts []Alert `json:"alerts"`
 }

 const ErrorLevel = "error"
 const SuccessLevel = "success"
 const InfoLevel = "info"
 const WarnLevel = "warn"

 func CreateErrorAlerts(errs ...error) Alerts {
 	alerts := []Alert{}
	for _ , err := range errs {
		alerts = append(alerts,Alert{err.Error(),ErrorLevel})
	}
	return Alerts{alerts}
 }

 func CreateAlerts(level string, messages ...string) Alerts {
	 alerts := []Alert{}
	 for _ , message := range messages {
		 alerts = append(alerts,Alert{message,level})
	 }
	 return Alerts{alerts}
}