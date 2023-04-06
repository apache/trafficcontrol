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
	"errors"
	"reflect"
	"sort"
	"strings"
)

// ColsFromStructByTag extracts the tag annotations from a struct into a string array.
func ColsFromStructByTag(tagName string, thing interface{}) []string {
	return ColsFromStructByTagExclude(tagName, thing, nil)
}

// InsertAtStr inserts insertMap (string to insert at -> []insert names) into cols non-destructively.
func InsertAtStr(cols []string, insertMap map[string][]string) []string {
	if insertMap == nil {
		return cols
	}
	if cols == nil {
		return nil
	}

	colLen := len(cols)
	insertLen := 0
	for _, val := range insertMap {
		insertLen += len(val)
	}
	newColumns := make([]string, colLen+insertLen)
	oldIndex := 0
	for newIndex := 0; newIndex < len(newColumns); newIndex++ {
		newColumns[newIndex] = (cols)[oldIndex]
		if inserts, ok := insertMap[newColumns[newIndex]]; ok {
			for j, insert := range inserts {
				newColumns[newIndex+j+1] = insert
			}
			newIndex += len(inserts)
		}
		oldIndex++
	}
	return newColumns
}

// ColsFromStructByTagExclude extracts the tag annotations from a struct into a string array except for excludedColumns.
func ColsFromStructByTagExclude(tagName string, thing interface{}, excludeColumns []string) []string {
	var cols []string
	var excludeMap map[string]bool
	if excludeColumns != nil {
		excludeMap = make(map[string]bool, len(excludeColumns))
		for _, col := range excludeColumns {
			excludeMap[col] = true
		}
	}
	t := reflect.TypeOf(thing)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tagName != "" {
			// Get the field tag value
			tag := field.Tag.Get(tagName)
			if _, ok := excludeMap[tag]; !ok && tag != "" {
				cols = append(cols, tag)
			}
		}
	}
	return cols
}

// sortableErrors provides ordering a list of errors for easier comparison with an expected list
type sortableErrors []error

func (s sortableErrors) Len() int {
	return len(s)
}
func (s sortableErrors) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortableErrors) Less(i, j int) bool {
	return s[i].Error() < s[j].Error()
}

// SortErrors sorts the list of errors lexically
func SortErrors(p []error) []error {
	if p == nil {
		return p
	}
	sort.Sort(sortableErrors(p))
	return p
}

func SplitErrors(err error) []error {
	if err == nil {
		return []error{}
	}
	strs := strings.Split(err.Error(), ", ")
	errs := []error{}
	for _, str := range strs {
		errs = append(errs, errors.New(str))
	}
	return errs
}
