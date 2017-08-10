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

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func ParsePerlObj(s string) (map[string]interface{}, error) {
	obj, _, err := getObj(s)
	return obj, err
}

func getObj(s string) (map[string]interface{}, string, error) {
	obj := map[string]interface{}{}

	s = strings.TrimSpace(s)
	if len(s) < 1 || s[0] != '{' {
		return obj, "", fmt.Errorf("expected first character '{': %v", s)
	}
	s = s[1:] // strip opening {
	s = strings.TrimSpace(s)

	// read top-level keys
	for {
		s = stripComment(s)
		s = strings.TrimSpace(s)
		// s = stripComment(s)
		if len(s) > 0 && s[0] == '}' {
			return obj, s[1:], nil
		}

		key := ""
		key, s = getKey(s)

		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return obj, "", fmt.Errorf("malformed string after key '%v'", key)
		}

		err := error(nil)
		switch {
		case s[0] == '{':
			v := map[string]interface{}{}
			v, s, err = getObj(s)
			if err != nil {
				return obj, "", fmt.Errorf("Error getting object value after key %v: %v", key, err)
			}
			obj[key] = v
		case s[0] == '\'':
			v := ""
			v, s, err = getStr(s)
			if err != nil {
				return obj, "", fmt.Errorf("Error getting string value after key %v: %v", key, err)
			}
			obj[key] = v
		case unicode.IsDigit(rune(s[0])):
			v := float64(0.0)
			v, s, err = getNum(s)
			if err != nil {
				return obj, "", fmt.Errorf("Error getting numeric value after key %v: %v", key, err)
			}
			obj[key] = v
		case s[0] == '[':
			v := []interface{}{}
			v, s, err = getArr(s)
			if err != nil {
				return obj, "", fmt.Errorf("Error getting array value after key %v: %v", key, err)
			}
			obj[key] = v
		default:
			return obj, "", fmt.Errorf(`malformed string after key "%v"`, key)
		}
		s = strings.TrimSpace(s)
		s = stripComment(s)
		if len(s) > 0 && s[0] == ',' {
			s = s[1:]
			s = strings.TrimSpace(s)
			s = stripComment(s)
			s = strings.TrimSpace(s)
		}
	}
}

func getNum(s string) (float64, string, error) {
	s = strings.TrimSpace(s)
	i := strings.IndexFunc(s, func(r rune) bool { return !unicode.IsDigit(r) && r != '.' })

	numStr := ""
	if i < 0 {
		numStr = s
		s = ""
	} else {
		numStr = s[:i]
		s = s[i:]
	}
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, "", fmt.Errorf("malformed number: %v", err)
	}

	return num, s, nil
}

func getArr(s string) ([]interface{}, string, error) {
	arr := []interface{}{}

	s = strings.TrimSpace(s)
	if len(s) == 0 || s[0] != '[' {
		return nil, "", fmt.Errorf("malformed array, doesn't start with [")
	}
	s = s[1:]
	for {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return nil, "", fmt.Errorf("malformed array, doesn't end with ]")
		}
		if s[0] == ']' {
			return arr, s[1:], nil
		}

		switch {
		case unicode.IsDigit(rune(s[0])) || s[0] == '-' || s[0] == '+' || s[0] == '.':
			num := float64(0.0)
			err := error(nil)
			num, s, err = getNum(s)
			if err != nil {
				return nil, "", fmt.Errorf("malformed number in array: %v", err)
			}
			arr = append(arr, num)
		case s[0] == '\'':
			str := ""
			err := error(nil)
			str, s, err = getStr(s)
			if err != nil {
				return nil, "", fmt.Errorf("malformed string in array: %v", err)
			}
			arr = append(arr, str)
		case s[0] == '[':
			narr := []interface{}{}
			err := error(nil)
			narr, s, err = getArr(s)
			if err != nil {
				return nil, "", fmt.Errorf("malformed array in array: %v", err)
			}
			arr = append(arr, narr)
		case s[0] == '{':
			obj := map[string]interface{}{}
			err := error(nil)
			obj, s, err = getObj(s)
			if err != nil {
				return nil, "", fmt.Errorf("malformed object in array: %v", err)
			}
			arr = append(arr, obj)
		default:
			return nil, "", fmt.Errorf("malformed element in array, unknown initial character: %v", string(s[0]))
		}
		s = strings.TrimSpace(s)
		s = stripComment(s)
		if len(s) > 0 && s[0] == ',' {
			s = s[1:]
			s = strings.TrimSpace(s)
			s = stripComment(s)
			s = strings.TrimSpace(s)
		}
	}
}

func getStr(s string) (string, string, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 || s[0] != '\'' {
		return "", "", fmt.Errorf("malformed string, doesn't start with '")
	}
	s = s[1:]

	str := ""

	escaping := false
	for {
		if len(s) == 0 {
			return "", "", fmt.Errorf("malformed string, doesn't terminate '")
		}

		if !escaping && s[0] == '\'' {
			return str, s[1:], nil
		}

		if escaping {
			str += s[0:1]
			s = s[1:]
			escaping = false
			continue
		}

		if s[0] == '\\' {
			escaping = true
			s = s[1:]
			continue
		}

		str += s[0:1]
		s = s[1:]
		continue
	}

}

func stripComment(s string) string {
	if len(s) == 0 || s[0] != '#' {
		return s
	}
	i := strings.Index(s, "\n")
	if i < 0 {
		return ""
	}
	return s[i:]
}

func getKey(s string) (string, string) {
	i := strings.Index(s, "=>")
	if i < 0 {
		return "", s
	}
	key := s[:i]
	key = strings.TrimSpace(key)

	s = s[i+len("=>"):]
	return key, s
}
