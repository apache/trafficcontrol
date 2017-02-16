/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package main

import (
	"bytes"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type License string

type Licenses []License

func (lics Licenses) Len() int {
	return len(lics)
}
func (lics Licenses) Swap(i, j int) {
	lics[i], lics[j] = lics[j], lics[i]
}
func (lics Licenses) Less(i, j int) bool {
	return lics[i] < lics[j]
}

func Uniq(lics []License) []License {
	if len(lics) == 0 {
		return nil
	}
	sortLics := make(Licenses, len(lics))
	copy(sortLics, lics)
	sort.Sort(sortLics)
	var uniqLics []License
	for _, lic := range sortLics {
		if len(uniqLics) == 0 || lic != uniqLics[len(uniqLics)-1] {
			uniqLics = append(uniqLics, lic)
		}
	}
	return uniqLics
}

func Remove(lics []License, rmLic License) []License {
	var rmLics []License
	for _, lic := range lics {
		if lic != rmLic {
			rmLics = append(rmLics, lic)
		}
	}
	return rmLics
}

func Has(lics []License, lic License) bool {
	for _, l := range lics {
		if l == lic {
			return true
		}
	}
	return false
}

func Collide(lics []License) []License {
	var toRm []License
	for _, lic := range lics {
		if string(lic)[0] == '!' {
			toRm = append(toRm, lic)
			toRm = append(toRm, License(string(lic)[1:]))
		}
	}
	newLics := lics
	for _, rm := range toRm {
		newLics = Remove(newLics, rm)
	}

	if Has(newLics, License("GoBSD")) {
		newLics = Remove(newLics, "BSD")
	}

	if len(newLics) > 1 {
		newLics = Remove(newLics, "Docs")
	}
	if len(newLics) > 1 {
		newLics = Remove(newLics, "Generated")
	}
	return newLics
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func wordMatcher(words []string) func(out *bool) (in chan<- string, done <-chan struct{}) {
	return func(out *bool) (in chan<- string, done <-chan struct{}) {
		ch := make(chan string, 32)
		doneCh := make(chan struct{})
		go func() {
			defer func() { close(doneCh) }()
			defer func() {
				for _ = range ch {
				}
			}()
			if len(words) == 0 {
				*out = true
				return
			}

			i := 0
			for word := range ch {
				//				fmt.Println(getGID(), "Comparing "+words[i]+" to "+word)
				if words[i] == word {
					i++
				} else {
					i = 0
				}

				if i == len(words) {
					*out = true
					return
				}
			}
		}()
		return ch, doneCh
	}
}

type multiMatcher []struct {
	ch      chan<- string
	done    <-chan struct{}
	value   *bool
	license License
}

func newMultiMatcher(in <-chan string) []License {
	var mm multiMatcher
	mmAppend := func(words []string, license License) {
		value := new(bool)
		ch, done := wordMatcher(words)(value)
		mm = append(mm, struct {
			ch      chan<- string
			done    <-chan struct{}
			value   *bool
			license License
		}{
			ch:      ch,
			done:    done,
			value:   value,
			license: license,
		})
	}

	mmAppend(wordsApache, License("Apache"))
	mmAppend(wordsApache2, License("Apache"))
	mmAppend(wordsBSD, License("BSD"))
	mmAppend(wordsBSD2, License("BSD"))
	mmAppend(wordsMIT, License("MIT"))
	mmAppend(wordsMIT2, License("MIT"))
	mmAppend(wordsGoBSD, License("GoBSD"))
	mmAppend(wordsISC, License("ISC"))
	mmAppend(wordsGen, License("Generated"))
	mmAppend(wordsX11, License("X11"))
	mmAppend(wordsWTFPL, License("WTFPL"))
	mmAppend(wordsGPL, License("GPL/LGPL"))
	mmAppend(wordsGPL2, License("GPL/LGPL"))
	mmAppend(wordsGPL3, License("GPL/LGPL"))
	mmAppend(wordsGPL4, License("GPL/LGPL"))
	mmAppend(wordsLGPL, License("GPL/LGPL"))
	mmAppend(wordsLGPL2, License("GPL/LGPL"))
	mmAppend(wordsLGPL3, License("GPL/LGPL"))
	mmAppend(wordsLGPL4, License("GPL/LGPL"))

	for word := range in {
		for _, m := range mm {
			m.ch <- word
		}
	}
	for _, m := range mm {
		close(m.ch)
	}
	for _, m := range mm {
		<-m.done
	}

	var licenses []License
doMatcher:
	for _, m := range mm {
		if *m.value {
			for _, lic := range licenses {
				if lic == m.license {
					continue doMatcher
				}
			}
			licenses = append(licenses, m.license)
		}
	}
	return licenses
}

func stripPunc(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return r
		}
		if unicode.IsSpace(r) {
			return r
		}
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, s)
}

func makeWords(s string) []string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "\n", ` `, -1)
	s = stripPunc(s)
	return strings.Split(s, ` `)
}

var (
	wordsApache  = makeWords(`Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements.`)
	wordsApache2 = makeWords(`Licensed under the Apache License`)
	wordsBSD     = makeWords(`Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:`)
	wordsBSD2    = makeWords(`BSD`)
	wordsMIT     = makeWords(`Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files`)
	wordsMIT2    = makeWords(`MIT`)
	wordsGoBSD   = makeWords(`The Go Authors. All rights reserved.`)
	wordsISC     = makeWords(`Permission to use, copy, modify, and distribute this software for any purpose with or without fee is hereby granted, provided that the above copyright notice and this permission notice appear in all copies`)
	wordsGen     = makeWords(`DO NOT MODIFY THE FIRST PART OF THIS FILE`)
	wordsX11     = makeWords(`X11`)
	wordsWTFPL   = makeWords(`WTFPL`)
	wordsGPL     = makeWords(`GNU General Public License`)
	wordsGPL2    = makeWords(`GPL`)
	wordsGPL3    = makeWords(`GPLv2`)
	wordsGPL4    = makeWords(`GPLv3`)
	wordsLGPL    = makeWords(`GNU Lesser General Public License`)
	wordsLGPL2   = makeWords(`LGPL`)
	wordsLGPL3   = makeWords(`LGPLv2`)
	wordsLGPL4   = makeWords(`LGPLv3`)
)
