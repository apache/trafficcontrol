// Package tcdata provides dynamic loading/unloading of ATC objects to/from a
// Traffic Ops instance.
//
// This should ONLY be imported by tests, that's the library's only purpose.
package tcdata

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/config"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/totest"
)

type TCData struct {
	Config   *config.Config
	TestData *totest.TrafficControl
}

func NewTCData() *TCData {
	return &TCData{
		Config:   &config.Config{},
		TestData: &totest.TrafficControl{},
	}
}
