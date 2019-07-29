package trafficstats

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

import "database/sql"
import "encoding/json"
import "errors"
import "fmt"
import "net/http"
import "strconv"
import "strings"
import "time"

import "github.com/apache/trafficcontrol/lib/go-tc"
import "github.com/apache/trafficcontrol/lib/go-log"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

import influx "github.com/influxdata/influxdb1-client/v2"

const DEFAULT_INTERVAL = tc.TrafficStatsDuration("1m")

var metricTypes = map[string]interface{}{
	"kbps":      struct{}{},
	"tps_total": struct{}{},
	"tps_2xx":   struct{}{},
	"tps_3xx":   struct{}{},
	"tps_4xx":   struct{}{},
	"tps_5xx":   struct{}{},
}

var jsonWithRFCTimestamps = tc.MimeType{
	"application/json",
	map[string]string{"timestamp": "rfc"},
}

var jsonWithUnixTimestamps = tc.MimeType{
	"application/json",
	map[string]string{"timestamp": "unix"},
}

type APIResponse struct {
	Response tc.TrafficStatsResponse `json:"response"`
	// Alerts []tc.Alert `json:"alerts,omitempty"`
}

const dsTenantIDFromXMLIDQuery = `
	SELECT tenant_id
	FROM deliveryservice
	WHERE xml_id = $1`

const xmlidFromIDQuery = `
	SELECT xml_id
	FROM deliveryservice
	WHERE id = $1`

// TODO: Pretty sure all of this could actually be calculated using the fetched data (assuming an
// interval is given). Check to see if that's faster than doing another synchronous HTTP request.
const summaryQuery = `
	SELECT mean(value) AS "average",
	       percentile(value, 5) AS "fifthPercentile",
	       percentile(value, 95) AS "ninetyFifthPercentile",
	       percentile(value, 98) AS "ninetyEighthPercentile",
	       min(value) AS "min",
	       max(value) AS "max",
	       count(value) AS "count"
	FROM "%s"."monthly"."%s.ds.1min"
	WHERE time >= $start
	AND time <= $end
	AND cachegroup = 'total'
	AND deliveryservice = $xmlid`

const seriesQuery = `
	SELECT mean(value) AS "value"
	FROM "%s"."monthly"."%s.ds.1min"
	WHERE cachegroup = 'total'
	AND deliveryservice = $xmlid
	AND time >= $start
	AND time <= $end
	GROUP BY time(%s, %s), cachegroup%s`

type orderable string

const (
	timeOrder  orderable = "time"
	valueOrder orderable = "value"
)

func orderableFromString(v string) *orderable {
	var o orderable
	switch v {
	case "time":
		o = timeOrder
	case "value":
		o = valueOrder
	default:
		return nil
	}
	return &o
}

type config struct {
	DeliveryService string
	End             time.Time
	ExcludeSeries   bool
	ExcludeSummary  bool
	Interval        tc.TrafficStatsDuration
	Limit           *uint64
	MetricType      string
	Offset          *uint64
	OrderBy         *orderable
	Start           time.Time
	Unix            bool
}

func configFromRequest(r *http.Request, i *api.APIInfo) (c config, e error, code int) {
	code = http.StatusBadRequest // this is so common I may as well do it right away

	if accept := r.Header.Get("Accept"); accept != "" {

		mimes, err := tc.MimeTypesFromAccept(accept)
		if err != nil {
			log.Warnf("Failed to negotiate content, Accept line '%s', error: %v", accept, err)
		} else {

			found := false
			for _, m := range mimes {
				if jsonWithRFCTimestamps.Equal(m) {
					found = true
					break
				}

				if jsonWithUnixTimestamps.Equal(m) {
					found = true
					c.Unix = true
					break
				}
			}

			if !found {
				e = fmt.Errorf("Failed to negotiate content; cannot produce output satisfying %s", accept)
				code = http.StatusNotAcceptable
				return
			}
		}
	}

	if limit, ok := i.Params["limit"]; ok {
		lim, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			e = errors.New("Invalid limit!")
			return
		}
		c.Limit = &lim
	}

	if offset, ok := i.Params["offset"]; ok {
		off, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			e = errors.New("Invalid offset!")
			return
		}
		c.Offset = &off
	}

	if orderby, ok := i.Params["orderby"]; ok {
		if c.OrderBy = orderableFromString(orderby); c.OrderBy == nil {
			e = errors.New("Invalid orderby!")
			return
		}
	}

	if c.Start, e = parseTime(i.Params["startDate"]); e != nil {
		log.Errorf("Parsing startDate: %v", e)
		e = errors.New("Invalid startDate!")
		return
	}

	if c.End, e = parseTime(i.Params["endDate"]); e != nil {
		log.Errorf("Parsing endDate: %v", e)
		e = errors.New("Invalid endDate!")
		return
	}

	if interval, ok := i.Params["interval"]; !ok {
		c.Interval = DEFAULT_INTERVAL
	} else if c.Interval = tc.TrafficStatsDurationFromString(interval); c.Interval == tc.InvalidDuration {
		log.Errorf("Error parsing 'interval' query parameter: %v", e)
		e = errors.New("Invalid interval!")
		return
	}

	if ex, ok := i.Params["exclude"]; ok {
		switch ex {
		case "series":
			c.ExcludeSeries = true
		case "summary":
			c.ExcludeSummary = true
		default:
			e = errors.New("Invalid exclude! Must be 'series' or 'summary'")
			return
		}
	}

	c.MetricType = i.Params["metricType"]
	if _, ok := metricTypes[c.MetricType]; !ok {
		e = fmt.Errorf("Unknown metric type: %s", c.MetricType)
		return
	}

	var ok bool
	if c.DeliveryService, ok = i.Params["deliveryServiceName"]; !ok {
		if c.DeliveryService, ok = i.Params["deliveryService"]; !ok {
			e = errors.New("You must specify deliveryService or deliveryServiceName!")
			return
		}

		if dsID, err := strconv.ParseUint(c.DeliveryService, 10, 64); err != nil {
			// sql.ErrNoRows does not *necessarily* mean the DS doesn't exist - an XMLID can simply
			// be numeric, and so it was wrong to treat it as an ID in the first place.
			xmlid := c.DeliveryService
			if c.DeliveryService, err = getXMLIDFromID(dsID, i.Tx.Tx); err != nil && err != sql.ErrNoRows {
				log.Errorf("Converting DSID to XMLID: %v", err)
				e = errors.New("Internal Server Error")
				code = http.StatusInternalServerError
				return
			} else if err == sql.ErrNoRows {
				c.DeliveryService = xmlid
			}
		}
	}

	code = http.StatusOK
	return
}

func GetDSStats(w http.ResponseWriter, r *http.Request) {
	// Perl didn't require "interval", but it would only return summary data if it was not given
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"metricType", "startDate", "endDate"}, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var c config
	if c, userErr, errCode = configFromRequest(r, inf); userErr != nil {
		sysErr = fmt.Errorf("Unable to process deliveryservice_stats request: %v", userErr)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	client, err := inf.CreateInfluxClient()
	if err != nil {
		if err == sql.ErrNoRows {
			errCode = http.StatusServiceUnavailable
			userErr = errors.New("No InfluxDB servers available!")
			sysErr = userErr
		} else {
			errCode = http.StatusInternalServerError
			userErr = nil
			sysErr = err
		}
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	} else if client == nil {
		userErr = errors.New("Traffic Stats is not configured!")
		sysErr = userErr
		errCode = http.StatusServiceUnavailable
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer (*client).Close()

	var dsTenant uint
	if dsTenant, err = dsTenantIDFromXMLID(c.DeliveryService, tx); err != nil {
		sysErr = err

		if err == sql.ErrNoRows {
			userErr = fmt.Errorf("No such Delivery Service: %s", c.DeliveryService)
			errCode = http.StatusNotFound
		} else {
			errCode = http.StatusInternalServerError
		}

		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	authorized, err := tenant.IsResourceAuthorizedToUserTx(int(dsTenant), inf.User, tx)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
	} else if !authorized {
		// If the Tenant is not authorized to use the resource, then we DON'T tell them that.
		// Instead, we don't disclose that such a Delivery Service exists at all - in keeping with
		// the behavior of /deliveryservices
		// This is different from what Perl used to do, but then again Perl didn't check tenancy at
		// all.
		userErr = fmt.Errorf("No such Delivery Service: %s", c.DeliveryService)
		sysErr = fmt.Errorf("GetDSStats: unauthorized Tenant (#%d) access", inf.User.TenantID)
		errCode = http.StatusNotFound
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	var errBuilder strings.Builder
	resp := APIResponse{
		tc.TrafficStatsResponse{
			Source:  tc.TRAFFIC_STATS_SOURCE,
			Version: tc.TRAFFIC_STATS_VERSION,
			Series:  nil,
			Summary: nil,
		},
	}

	// TODO: as above, this could be done on TO itself, thus sending only one synchronous request
	// per hit on this endpoint, rather than the current two. Not sure if that's worth it for large
	// data sets, though.
	if !c.ExcludeSummary {
		summary, messages, err := getSummary(client, &c, inf.Config.ConfigInflux.DSDBName)

		errBuilder := strings.Builder{}
		if messages != nil && len(messages) > 0 {
			errBuilder.Write([]byte("Messages: ["))
			for _, m := range messages {
				errBuilder.WriteString(m.Level)
				errBuilder.WriteRune(':')
				errBuilder.WriteString(m.Text)
				errBuilder.Write([]byte(", "))
			}
			errBuilder.WriteRune(']')
		}

		if err != nil {
			sysErr = fmt.Errorf("Getting summary response from Influx: %v - %s", err, errBuilder.String())
			api.HandleErr(w, r, tx, http.StatusBadGateway, nil, sysErr)
			return
		}

		log.Debugf("Messages from summary query: %s", errBuilder.String())
		errBuilder.Reset()
		messages = nil

		resp.Response.Summary = &summary
	}

	if !c.ExcludeSeries {
		series, messages, err := getSeries(client, &c, inf.Config.ConfigInflux.DSDBName)
		if messages != nil && len(messages) > 0 {
			errBuilder.Write([]byte("Messages: ["))
			for _, m := range messages {
				errBuilder.WriteString(m.Level)
				errBuilder.WriteRune(':')
				errBuilder.WriteString(m.Text)
				errBuilder.Write([]byte(", "))
			}
			errBuilder.WriteRune(']')
		}

		if err != nil {
			sysErr = fmt.Errorf("Getting summary response from Influx: %v - %s", err, errBuilder.String())
			api.HandleErr(w, r, tx, http.StatusBadGateway, nil, sysErr)
			return
		}

		log.Debugf("Messages from series query: %s", errBuilder.String())
		errBuilder.Reset()
		messages = nil

		if !c.Unix {
			for i, v := range series.Values {
				if len(v) != 2 {
					log.Warnf("Malformed series data point: %v", v)
					continue
				}

				// TODO: model the data better so this isn't as scary (possible?)
				switch t := v[0].(type) {
				case int64:
					series.Values[i][0] = time.Unix(0, v[0].(int64)).Format(time.RFC3339)
				case float64:
					series.Values[i][0] = time.Unix(0, int64(v[0].(float64))).Format(time.RFC3339)
				case json.Number:
					val, err := v[0].(json.Number).Int64()
					if err != nil {
						log.Warnf("Error encountered trying to coerce %v to an int64: %v", v, err)
					} else {
						series.Values[i][0] = time.Unix(0, val).Format(time.RFC3339)
					}
				default:
					log.Warnf("Invalid type %T for data point", t)
				}
			}
		}

		resp.Response.Series = &series
	}

	respBts, err := json.Marshal(resp)
	if err != nil {
		sysErr = fmt.Errorf("Marshalling response: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	if c.Unix {
		w.Header().Set(tc.ContentType, jsonWithUnixTimestamps.String())
	} else {
		w.Header().Set(tc.ContentType, jsonWithRFCTimestamps.String())
	}
	w.Header().Set(http.CanonicalHeaderKey("vary"), http.CanonicalHeaderKey("Accept"))
	w.Write(respBts)
}

func getSummary(client *influx.Client, conf *config, db string) (tc.TrafficStatsSummary, []influx.Message, error) {

	msgs := []influx.Message{}
	s := tc.TrafficStatsSummary{}
	qStr := fmt.Sprintf(summaryQuery, db, conf.MetricType)
	q := influx.NewQueryWithParameters(qStr,
		db,
		"rfc3339", //this doesn't actually seem to have any effect...
		map[string]interface{}{
			"xmlid":    conf.DeliveryService,
			"start":    conf.Start,
			"end":      conf.End,
			"interval": string(conf.Interval),
		})
	q.RetentionPolicy = "monthly"

	log.Debugf("InfluxDB SummaryQuery: %+v", q)

	resp, err := (*client).Query(q)
	if err != nil {
		return s, msgs, err
	}

	if resp.Results != nil && len(resp.Results) == 1 {
		r := resp.Results[0]
		if r.Messages != nil {
			for _, m := range r.Messages {
				if m != nil {
					msgs = append(msgs, *m)
				}
			}
		}

		if len(r.Series) != 1 {
			return s, msgs, fmt.Errorf("Improper number of series: %d", len(r.Series))
		}

		series := r.Series[0]

		if len(series.Values) != 1 {
			return s, msgs, fmt.Errorf("Improper number of returned rows: %d", len(r.Series[0].Values))
		}

		vals := series.Values[0]
		if len(vals) != 8 || len(series.Columns) != 8 {
			return s, msgs, fmt.Errorf("Improper number of returned values in row: %d (%d cols)", len(vals), len(series.Columns))
		}

		mappedValues := map[string]interface{}{}
		for i, v := range vals {
			mappedValues[series.Columns[i]] = v
		}

		var err error
		if s.Average, err = extractFloat64("average", mappedValues); err != nil {
			return s, msgs, err
		}

		if s.Count, err = extractUInt("count", mappedValues); err != nil {
			return s, msgs, err
		}

		if s.FifthPercentile, err = extractFloat64("fifthPercentile", mappedValues); err != nil {
			return s, msgs, err
		}

		if s.Max, err = extractFloat64("max", mappedValues); err != nil {
			return s, msgs, err
		}

		if s.Min, err = extractFloat64("min", mappedValues); err != nil {
			return s, msgs, err
		}

		if s.Max, err = extractFloat64("max", mappedValues); err != nil {
			return s, msgs, err
		}

		if s.NinetyEighthPercentile, err = extractFloat64("ninetyEighthPercentile", mappedValues); err != nil {
			return s, msgs, err
		}

		if s.NinetyFifthPercentile, err = extractFloat64("ninetyFifthPercentile", mappedValues); err != nil {
			return s, msgs, err
		}

	} else {
		log.Debugf("InfluxDB summary response: %+v", resp)
		return s, msgs, errors.New("'results' missing or improper!")
	}

	if resp.Error() != nil {
		log.Debugf("response error, summary object was: %+v", s)
		return s, msgs, resp.Error()
	}

	value := float64(s.Count * 60) * s.Average
	if conf.MetricType == "kbps" {
		value /= 1000
		s.TotalBytes = &value
	} else {
		s.TotalTransactions = &value
	}

	return s, msgs, nil
}

func extractUInt(k string, m map[string]interface{}) (uint, error) {
	tmp, ok := m[k]
	if !ok {
		return 0, fmt.Errorf("response has no value for column %s", k)
	}
	switch t := tmp.(type) {
	case float64:
		return uint(tmp.(float64)), nil
	case json.Number:
		ret, err := tmp.(json.Number).Int64()
		if err != nil {
			return 0, fmt.Errorf("Error parsing value for column '%s' as an int64: %v", k, err)
		}
		return uint(ret), nil
	default:
		return 0, fmt.Errorf("invalid type for column '%s' - expected unsigned integer, got %T (%v)", k, t, tmp)
	}
}

func extractFloat64(k string, m map[string]interface{}) (float64, error) {
	tmp, ok := m[k]
	if !ok {
		return 0, fmt.Errorf("response has no value for column %s", k)
	}

	switch t := tmp.(type) {
	case float64:
		return tmp.(float64), nil
	case json.Number:
		ret, err := tmp.(json.Number).Float64()
		if err != nil {
			return 0, fmt.Errorf("Error parsing value for column '%s' as a float64: %v", k, err)
		}
		return ret, nil
	case nil:
		// This is the only field that can be nil - because sometimes there isn't enough data for
		// 5% of it to be below a value. Could probably coalesce this in the query at some point, though.
		if k == "fifthPercentile" {
			return 0, nil
		}
		return 0, fmt.Errorf("column '%s' was null/blank", k)
	default:
		return 0, fmt.Errorf("invalid type for column '%s' - expected 'float64' or 'json.Number', got %T (%v)", k, t, tmp)
	}
}

func parseTime(raw string) (t time.Time, e error) {
	if t, e = time.Parse(time.RFC3339Nano, raw); e == nil {
		return
	}

	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		t = time.Unix(0, i)
		return
	}

	t, e = time.Parse(tc.TimeLayout, raw)
	return
}

func dsTenantIDFromXMLID(xmlid string, tx *sql.Tx) (tid uint, err error) {
	row := tx.QueryRow(dsTenantIDFromXMLIDQuery, xmlid)
	err = row.Scan(&tid)
	return
}

func getXMLIDFromID(id uint64, tx *sql.Tx) (xmlid string, err error) {
	row := tx.QueryRow(xmlidFromIDQuery, id)
	err = row.Scan(&xmlid)
	return
}

// This is a stupid, dirty hack to try to convince Influx to not give back data that's outside of the
// range in a WHERE clause. It doesn't work, but it helps.
// (https://github.com/influxdata/influxdb/issues/8010)
func (c *config) OffsetString() string {
	return fmt.Sprintf("%ds", int64(c.Start.Sub(time.Unix(0, 0))/time.Second)%c.Interval.Seconds())
}

func getSeries(client *influx.Client, conf *config, db string) (tc.TrafficStatsSeries, []influx.Message, error) {
	s := tc.TrafficStatsSeries{}
	msgs := []influx.Message{}
	extraClauses := strings.Builder{}
	if conf.OrderBy != nil {
		extraClauses.Write([]byte(" ORDER BY "))
		extraClauses.WriteString(string(*conf.OrderBy))
	}

	if conf.Limit != nil {
		extraClauses.WriteString(fmt.Sprintf(" LIMIT %d", *conf.Limit))
	}

	if conf.Offset != nil {
		extraClauses.WriteString(fmt.Sprintf(" OFFSET %d", *conf.Offset))
	}

	qStr := fmt.Sprintf(seriesQuery, db, conf.MetricType, conf.Interval, conf.OffsetString(), extraClauses.String())
	q := influx.NewQueryWithParameters(qStr,
		db,
		"rfc3339", // this doesn't seem to do anything...
		map[string]interface{}{
			"xmlid": conf.DeliveryService,
			"start": conf.Start,
			"end":   conf.End,
		})
	q.RetentionPolicy = "monthly"

	log.Debugf("InfluxDB series query: %+v", q)

	resp, err := (*client).Query(q)
	if err != nil {
		return s, msgs, err
	}

	if resp.Results != nil && len(resp.Results) == 1 {
		r := resp.Results[0]
		if r.Messages != nil {
			for _, m := range r.Messages {
				if m != nil {
					msgs = append(msgs, *m)
				}
			}
		}

		if len(r.Series) != 1 {
			return s, msgs, fmt.Errorf("Improper number of series: %d", len(r.Series))
		}

		series := r.Series[0]

		s = tc.TrafficStatsSeries{
			Name:    series.Name,
			Tags:    series.Tags,
			Values:  series.Values,
			Columns: series.Columns,
			Count:   uint(len(series.Values)),
		}

	} else {
		log.Debugf("InfluxDB series response: %+v", resp)
		return s, msgs, errors.New("'results' missing or improper!")
	}

	if resp.Error() != nil {
		log.Debugf("response error, series object was %+v", s)
		return s, msgs, resp.Error()
	}

	return s, msgs, nil
}
