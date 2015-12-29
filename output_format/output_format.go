// Copyright 2015 Comcast Cable Communications Management, LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output_format

// {"alerts":[{"level":"success","text":"Successfully logged in."}],"version":"1.1"}
type Result struct {
	Alerts  []Alert
	Version string `json:"version"`
}

type Alert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

type ApiWrapper struct {
	Resp    interface{} `json:"response,omitempty"`
	Error   string      `json:"error,omitempty"`
	Version float64     `json:"version"`
	Alerts  []Alert     `json:"alerts,omitempty"`
}

//wraps the given interface r into a returned Wrapper
//prepped for encoding to stream
func MakeApiResponse(r interface{}, alertString, errString string) ApiWrapper {
	// if alert == "" {
	// 	w := ApiWrapper{
	// 		Resp:    r,
	// 		Error:   err,
	// 		Version: 2.0,
	// 	}

	// } else {
	// 	w := ApiWrapper{
	// 		Resp:    r,
	// 		Error:   err,
	// 		Version: 2.0,
	// 		Alert:   Alert{Level: 1, String: alert},
	// 	}
	// }
	alert := Alert{Level: "success", Text: alertString}
	alerts := make([]Alert, 0, 1)
	alerts = append(alerts, alert)
	w := ApiWrapper{
		Resp:    r,
		Error:   errString,
		Version: 2.0,
		Alerts:  alerts,
	}

	return w
}
