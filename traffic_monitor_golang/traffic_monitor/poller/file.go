package poller

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
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
)

// FilePoller starts a goroutine polling the given file for changes. When changes occur, including an initial read, the result callback is called asynchronously. Returns a kill chan, which will kill the file poller when written to.
func File(filename string, result func([]byte, error)) (chan<- struct{}, error) {
	die := make(chan struct{})

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file '%v': %v", filename, err)
	}
	go result(contents, nil)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	go func() {
		watcher.Add(filename)
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					go result(ioutil.ReadFile(filename))
				}
			case err := <-watcher.Errors:
				go result(nil, err)
			case <-die:
				return
			}
		}
	}()

	return die, nil
}
