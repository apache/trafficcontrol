package test

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
	"reflect"
	"strings"
)

// Extract the tag annotations from a struct into a string array
func ColsFromStructByTag(tagName string, thing interface{}) []string {
	cols := []string{}
	t := reflect.TypeOf(thing)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if (strings.Compare(tagName, "db") == 0) && (tagName != "") {
			// Get the field tag value
			tag := field.Tag.Get(tagName)
			if tag != "" {
				cols = append(cols, tag)
			}
		}
	}
	return cols
}
