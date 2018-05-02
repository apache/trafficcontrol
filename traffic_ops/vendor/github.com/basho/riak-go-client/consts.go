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
	"time"
)

const (
	threeSeconds                  = time.Second * 3
	fiveSeconds                   = time.Second * 5
	tenSeconds                    = time.Second * 10
	defaultBucketType             = "default"
	defaultRemotePort             = uint16(8087)
	defaultMinConnections         = uint16(1)
	defaultMaxConnections         = uint16(256)
	defaultIdleExpirationInterval = fiveSeconds
	defaultIdleTimeout            = tenSeconds
	defaultConnectTimeout         = threeSeconds
	defaultRequestTimeout         = fiveSeconds
	defaultHealthCheckInterval    = 125 * time.Millisecond
	defaultExecutionAttempts      = byte(3)
	defaultQueueExecutionInterval = 125 * time.Millisecond
	defaultInitBuffer             = 2048
	defaultTempNetErrorRetries    = uint16(0)
)

var defaultRemoteAddress = fmt.Sprintf("127.0.0.1:%d", defaultRemotePort)
