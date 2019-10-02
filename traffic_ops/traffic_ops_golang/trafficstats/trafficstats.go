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
import "github.com/apache/trafficcontrol/lib/go-rfc"
import "github.com/apache/trafficcontrol/lib/go-log"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

import influx "github.com/influxdata/influxdb1-client/v2"

var (
	metricTypes = map[string]interface{}{
		"kbps":      struct{}{},
		"tps_total": struct{}{},
		"tps_2xx":   struct{}{},
		"tps_3xx":   struct{}{},
		"tps_4xx":   struct{}{},
		"tps_5xx":   struct{}{},
	}

	jsonWithRFCTimestamps = rfc.MimeType{
		"application/json",
		map[string]string{"timestamp": "rfc"},
	}

	jsonWithUnixTimestamps = rfc.MimeType{
		"application/json",
		map[string]string{"timestamp": "unix"},
	}
)

const (
	DEFAULT_INTERVAL         = "1m"
	dsTenantIDFromXMLIDQuery = `
		SELECT tenant_id
		FROM deliveryservice
		WHERE xml_id = $1`

	xmlidFromIDQuery = `
		SELECT xml_id
		FROM deliveryservice
		WHERE id = $1`

	// TODO: Pretty sure all of this could actually be calculated using the fetched data (assuming an
	// interval is given). Check to see if that's faster than doing another synchronous HTTP request.
	summaryQuery = `
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

	seriesQuery = `
		SELECT mean(value)
		FROM "%s"."monthly"."%s.ds.1min"
		WHERE cachegroup = 'total'
		AND deliveryservice = $xmlid
		AND time >= $start
		AND time <= $end
		GROUP BY time(%s, %s), cachegroup%s`
)

func ConfigFromRequest(r *http.Request, i *api.APIInfo) (tc.TrafficStatsConfig, error, int) {
	var c tc.TrafficStatsConfig
	var e error
	if accept := r.Header.Get("Accept"); accept != "" {

		mimes, err := rfc.MimeTypesFromAccept(accept)
		if err != nil {
			log.Warnf("Failed to negotiate content, Accept line '%s', error: %v", accept, err)
		} else {

			found := false
			for _, m := range mimes {
				if jsonWithRFCTimestamps.Satisfy(m) {
					found = true
					break
				}

				if jsonWithUnixTimestamps.Satisfy(m) {
					found = true
					c.Unix = true
					break
				}
			}

			if !found {
				e = fmt.Errorf("Failed to negotiate content; cannot produce output satisfying %s", accept)
				return c, e, http.StatusNotAcceptable
			}
		}
	}

	if limit, ok := i.Params["limit"]; ok {
		lim, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			e = errors.New("Invalid limit!")
			return c, e, http.StatusBadRequest
		}
		c.Limit = &lim
	}

	if offset, ok := i.Params["offset"]; ok {
		off, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			e = errors.New("Invalid offset!")
			return c, e, http.StatusBadRequest
		}
		c.Offset = &off
	}

	if orderby, ok := i.Params["orderby"]; ok {
		if c.OrderBy = tc.OrderableFromString(orderby); c.OrderBy == nil {
			e = errors.New("Invalid orderby! Must be 'time' or 'value'")
			return c, e, http.StatusBadRequest
		}
	}

	if c.Start, e = parseTime(i.Params["startDate"]); e != nil {
		log.Errorf("Parsing startDate: %v", e)
		e = errors.New("Invalid startDate!")
		return c, e, http.StatusBadRequest
	}

	if c.End, e = parseTime(i.Params["endDate"]); e != nil {
		log.Errorf("Parsing endDate: %v", e)
		e = errors.New("Invalid endDate!")
		return c, e, http.StatusBadRequest
	}

	if interval, ok := i.Params["interval"]; !ok {
		c.Interval = DEFAULT_INTERVAL
	} else if !tc.TrafficStatsDurationPattern.MatchString(interval) {
		e = errors.New("interval: must be a valid InfluxQL duration literal (resolution no less than minute)")
		return c, e, http.StatusBadRequest
	} else {
		c.Interval = interval
	}

	if ex, ok := i.Params["exclude"]; ok {
		switch tc.ExcludeFromString(ex) {
		case tc.ExcludeSummary:
			c.ExcludeSummary = true
		case tc.ExcludeSeries:
			c.ExcludeSeries = true
		default:
			e = errors.New("Invalid exclude! Must be 'series' or 'summary'")
			return c, e, http.StatusBadRequest
		}
	}

	c.MetricType = i.Params["metricType"]
	if _, ok := metricTypes[c.MetricType]; !ok {
		e = fmt.Errorf("Unknown metric type: %s", c.MetricType)
		return c, e, http.StatusBadRequest
	}

	var ok bool
	if c.DeliveryService, ok = i.Params["deliveryServiceName"]; !ok {
		if c.DeliveryService, ok = i.Params["deliveryService"]; !ok {
			e = errors.New("You must specify deliveryService or deliveryServiceName!")
			return c, e, http.StatusBadRequest
		}

		if dsID, err := strconv.ParseUint(c.DeliveryService, 10, 64); err == nil {
			// sql.ErrNoRows does not *necessarily* mean the DS doesn't exist - an XMLID can simply
			// be numeric, and so it was wrong to treat it as an ID in the first place.
			xmlid := c.DeliveryService
			var exists bool
			if exists, c.DeliveryService, err = getXMLIDFromID(dsID, i.Tx.Tx); err != nil {
				log.Errorf("Converting DSID to XMLID: %v", err)
				e = errors.New("Internal Server Error")
				return c, e, http.StatusInternalServerError
			} else if !exists {
				c.DeliveryService = xmlid
			}
		}
	}

	return c, nil, http.StatusOK
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

	var c tc.TrafficStatsConfig
	if c, userErr, errCode = ConfigFromRequest(r, inf); userErr != nil {
		sysErr = fmt.Errorf("Unable to process deliveryservice_stats request: %v", userErr)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	client, err := inf.CreateInfluxClient()
	if err != nil {
		errCode = http.StatusInternalServerError
		sysErr = err
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	} else if client == nil {
		sysErr = errors.New("Traffic Stats is not configured, but DS stats were requested")
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}
	defer (*client).Close()

	exists, dsTenant, err := dsTenantIDFromXMLID(c.DeliveryService, tx)
	if err != nil {
		sysErr = err
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	} else if !exists {
		userErr = fmt.Errorf("No such Delivery Service: %s", c.DeliveryService)
		errCode = http.StatusNotFound
		api.HandleErr(w, r, tx, errCode, userErr, nil)
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

	resp := struct{ Response tc.TrafficStatsResponse `json:"response"` }{
		Response: tc.TrafficStatsResponse{
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
		log.Debugf("Messages from summary query: %s", tc.MessagesToString(messages))

		if err != nil {
			sysErr = fmt.Errorf("Getting summary response from Influx: %v", err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}

		resp.Response.Summary = &summary
	}

	if !c.ExcludeSeries {
		series, messages, err := getSeries(client, &c, inf.Config.ConfigInflux.DSDBName)
		log.Debugf("Messages from series query: %s", tc.MessagesToString(messages))

		if err != nil {
			sysErr = fmt.Errorf("Getting summary response from Influx: %v", err)
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, sysErr)
			return
		}

		if !c.Unix {
			series.FormatTimestamps()
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
	w.Write(append(respBts, '\n'))
}

func getSummary(client *influx.Client, conf *tc.TrafficStatsConfig, db string) (tc.TrafficStatsSummary, []influx.Message, error) {

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

	value := float64(s.Count*60) * s.Average
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

func parseTime(raw string) (time.Time, error) {
	t, e := time.Parse(time.RFC3339Nano, raw)
	if e == nil {
		return t, nil
	}

	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		t = time.Unix(0, i)
		return t, nil
	}

	return time.Parse(tc.TimeLayout, raw)
}

func dsTenantIDFromXMLID(xmlid string, tx *sql.Tx) (bool, uint, error) {
	row := tx.QueryRow(dsTenantIDFromXMLIDQuery, xmlid)
	var tid uint
	err := row.Scan(&tid)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	return true, tid, err
}

func getXMLIDFromID(id uint64, tx *sql.Tx) (bool, string, error) {
	row := tx.QueryRow(xmlidFromIDQuery, id)
	var xmlid string
	err := row.Scan(&xmlid)
	if err == sql.ErrNoRows {
		return false, "", nil
	}
	return true, xmlid, err
}

func getSeries(client *influx.Client, conf *tc.TrafficStatsConfig, db string) (tc.TrafficStatsSeries, []influx.Message, error) {
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
