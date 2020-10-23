// Copyright 2015-present Basho Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package riak

import (
	"strings"
	"testing"
)

func TestSplitRemoteAddress(t *testing.T) {
	s := strings.SplitN(defaultRemoteAddress, ":", 2)
	if expected, actual := "127.0.0.1", s[0]; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
	if expected, actual := "8087", s[1]; expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestNewClientWithInvalidData(t *testing.T) {
	opts := &NewClientOptions{
		RemoteAddresses: []string{
			"FOO:BAR:BAZ",
		},
	}
	c, err := NewClient(opts)
	if err == nil {
		t.Errorf("expected non-nil error, %v", c)
	}

	opts = &NewClientOptions{
		RemoteAddresses: []string{
			"127.0.0.1:FRAZZLE",
		},
	}
	c, err = NewClient(opts)
	if err == nil {
		t.Errorf("expected non-nil error, %v", c)
	}
}
