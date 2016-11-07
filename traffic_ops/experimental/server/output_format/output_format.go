
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

import (
	"log"
	// "github.com/jmoiron/sqlx"
	"database/sql"
	"reflect"
	"strconv"
)

const APIVERSION = 2.0

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

func MakeAlert(alertTxt string, alertLevel string) []Alert {
	alert := Alert{Level: alertLevel, Text: alertTxt}
	var alerts []Alert
	alerts = append(alerts, alert)
	return alerts
}

// wraps the given interface r into a returned Wrapper
// prepped for encoding to stream
func MakeApiResponse(r interface{}, alerts []Alert, err error) ApiWrapper {
	var w ApiWrapper
	if err != nil {
		alerts = append(alerts, Alert{Level: "error", Text: "Internal error: " + err.Error()})
		w = ApiWrapper{
			Version: APIVERSION,
			Alerts:  alerts,
		}
	} else {
		w = ApiWrapper{
			Version: APIVERSION,
			Alerts:  alerts,
		}
		if r != nil {
			rType := reflect.TypeOf(r)
			if rType.Kind() == reflect.Slice {
				rows := reflect.ValueOf(r)
				numRows := strconv.Itoa(rows.Len())
				alerts = append(alerts, Alert{Level: "success", Text: numRows + " rows returned."})
				w = ApiWrapper{
					Resp:    r,
					Version: APIVERSION,
					Alerts:  alerts,
				}
			} else if rType.Kind() == reflect.Struct {
				// lastInserted doesn't work w pq
				// lastInserted, err := r.(sql.Result).LastInsertId()
				// if err != nil {
				// 	log.Println("error on LastInsertedId")
				// }
				rowsAffected, err := r.(sql.Result).RowsAffected()
				if err != nil {
					log.Println("error on RowsAffected()")
					alerts = append(alerts, Alert{Level: "error", Text: "Internal error:" + err.Error()})
				} else {
					alerts = append(alerts, Alert{Level: "success", Text: strconv.FormatInt(rowsAffected, 10) + " rows affected."})
				}
				w = ApiWrapper{
					Version: APIVERSION,
					Alerts:  alerts,
				}

			}
		}
	}

	return w
}
