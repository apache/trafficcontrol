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
	"fmt"
)

func rpbValidateResp(data []byte, expected byte) (err error) {
	if len(data) == 0 {
		err = ErrZeroLength
		return
	}
	if err = rpbEnsureCode(expected, data[0]); err != nil {
		return
	}
	return
}

func rpbEnsureCode(expected byte, actual byte) (err error) {
	if expected != actual {
		err = newClientError(fmt.Sprintf("expected response code %d, got: %d", expected, actual), nil)
	}
	return
}

func rpbBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}
