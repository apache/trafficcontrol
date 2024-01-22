package api

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
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
)

const createChangeLogQuery = `
INSERT INTO log (
	level,
	message,
	tm_user
) VALUES (
	$1,
	$2,
	$3
)
`
const influxServersQuery = `
SELECT (host_name||'.'||domain_name) as fqdn,
       tcp_port,
       https_port
FROM server
WHERE type in ( SELECT id
                FROM type
                WHERE name='INFLUXDB'
              )
AND status=(SELECT id FROM status WHERE name='ONLINE')
`

// Info structures contain all of the information an API route handler needs to
// be able to service a request, including some things that are pre-parsed (e.g.
// query string parameters) for you. It also provides some methods for
// accomplishing common tasks.
type Info struct {
	// Params is a mapping of all query string and path parameters to their
	// respective values. The behavior of this map is not defined when any two
	// query string parameters and/or path parameters share a name. For example,
	// if the route is `cdns/{id}/delivery_services/{id}`, the two cannot be
	// distinguished. Similarly, a request like `GET /api/5.0/cdns?id=1&id=2`
	// will give either an "id" key that maps to "1", or an "id" key that maps
	// to "2". Most convolutedly, for the aforementioned route definition, the
	// request `GET cdns/1/deliveryservices/2?id=3&id=4` gives four possible
	// values for the "id" key. Take care when constructing routes and deciding
	// the parameters they will accept.
	Params map[string]string
	// IntParams is a mapping of all of the declared parameters that are to be
	// parsed as ints to the parsed values of those parameters. No key will
	// appear here that isn't also in Params.
	IntParams map[string]int
	// The currently authenticated user - this may be `nil` on routes that do
	// not require authentication.
	User *auth.CurrentUser
	// A unique identifier for the request.
	ReqID uint64
	// The version of the API being requested. This is a pointer for legacy
	// reasons - all handlers should assume this is not nil (with the possible
	// exception of plugin handlers).
	Version *Version
	// A transaction opened to the Traffic Ops database.
	Tx *sqlx.Tx
	// The cancel function for the request and transaction contexts.
	CancelTx context.CancelFunc
	// The Traffic Vault implementation.
	Vault trafficvault.TrafficVault
	// Config is the Traffic Ops server's current configuration. This is a
	// pointer presumably to save memory; it should and must never be `nil`.
	Config *config.Config

	request *http.Request
	w       http.ResponseWriter
}

// NewInfo get and returns the context info needed by handlers. It also returns
// any user error, any system error, and the status code which should be
// returned to the client if an error occurred.
//
// It is encouraged to call Info.Tx.Tx.Commit() manually when all queries are
// finished, to release database resources early, and also to return an error to
// the user if the commit failed. In practice, though, the `Close` method
// handles this in nearly every case.
//
// NewInfo guarantees the returned Info.Tx is non-`nil` and Info.Tx.Tx is `nil`
// or valid, even if a returned error is not `nil`. Hence, it is safe to pass
// the Tx.Tx to HandleErr when this returns errors.
//
// Close() must be called to free resources, and should be called in a defer
// immediately after NewInfo(), to finish the transaction.
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	  inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
//	  if userErr != nil || sysErr != nil {
//	    api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
//	    return
//	  }
//	  defer inf.Close()
//
//	  respObj, err := finalDatabaseOperation(inf.Tx)
//	  if err != nil {
//	    api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("final db op: %w", err))
//	    return
//	  }
//	  if err := inf.Tx.Tx.Commit(); err != nil {
//	    api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("committing transaction: %w", err))
//	    return
//	  }
//	  api.WriteResp(w, r, respObj)
//	}
func NewInfo(r *http.Request, requiredParams []string, intParamNames []string) (*Info, error, error, int) {
	db, err := GetDB(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, fmt.Errorf("getting db: %w", err), nil, http.StatusInternalServerError
	}
	cfg, err := GetConfig(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, fmt.Errorf("getting config: %w", err), nil, http.StatusInternalServerError
	}
	tv, err := GetTrafficVault(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, fmt.Errorf("getting TrafficVault: %w", err), nil, http.StatusInternalServerError
	}
	reqID, err := getReqID(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, fmt.Errorf("getting reqID: %w", err), nil, http.StatusInternalServerError
	}
	version := GetRequestedAPIVersion(r.URL.Path)

	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, fmt.Errorf("getting user: %w", err), nil, http.StatusInternalServerError
	}
	params, intParams, userErr, sysErr, errCode := AllParams(r, requiredParams, intParamNames)
	if userErr != nil || sysErr != nil {
		return &Info{Tx: &sqlx.Tx{}}, userErr, sysErr, errCode
	}

	// only place we could call cancel here is in Info.Close(), which already
	// will rollback the transaction (which is all cancel will do.)
	// must be last, MUST not return an error if this succeeds, without closing
	// the tx
	dbCtx, cancelTx := context.WithTimeout(r.Context(), time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second)
	tx, err := db.BeginTxx(dbCtx, nil)
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}, CancelTx: cancelTx}, userErr, fmt.Errorf("could not begin transaction: %w", err), http.StatusInternalServerError
	}
	return &Info{
		Config:    cfg,
		ReqID:     reqID,
		Version:   version,
		Params:    params,
		IntParams: intParams,
		User:      user,
		Tx:        tx,
		CancelTx:  cancelTx,
		Vault:     tv,
		request:   r,
	}, nil, nil, http.StatusOK
}

// CheckPrecondition checks a request's "preconditions" - its If-Match and
// If-Unmodified-Since headers versus the last updated time of the requested
// object(s), and returns (in order), an HTTP response code appropriate for the
// precondition check results, a user-safe error that should be returned to
// clients, and a server-side error that should be logged.
// Callers must pass in a query that will return one row containing one column
// that is the representative date/time of the last update of the requested
// object(s), and optionally any values for placeholder arguments in the query.
func (inf Info) CheckPrecondition(query string, args ...interface{}) (int, error, error) {
	if inf.request == nil {
		return http.StatusInternalServerError, nil, NilRequestError
	}

	ius := inf.request.Header.Get(rfc.IfUnmodifiedSince)
	etag := inf.request.Header.Get(rfc.IfMatch)
	if ius == "" && etag == "" {
		return http.StatusOK, nil, nil
	}

	if inf.Tx == nil || inf.Tx.Tx == nil {
		return http.StatusInternalServerError, nil, NilTransactionError
	}

	var lastUpdated time.Time
	if err := inf.Tx.Tx.QueryRow(query, args...).Scan(&lastUpdated); err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("scanning for lastUpdated: %w", err)
	}

	if etag != "" {
		if et, ok := rfc.ParseETags(strings.Split(etag, ",")); ok {
			if lastUpdated.After(et) {
				return http.StatusPreconditionFailed, ResourceModifiedError, nil
			}
		}
	}

	if ius == "" {
		return http.StatusOK, nil, nil
	}

	if tm, ok := rfc.ParseHTTPDate(ius); ok {
		if lastUpdated.After(tm) {
			return http.StatusPreconditionFailed, ResourceModifiedError, nil
		}
	}

	return http.StatusOK, nil, nil
}

// Close implements the io.Closer interface. It should be called in a defer immediately after NewInfo().
//
// Close will commit the transaction, if it hasn't been rolled back.
func (inf *Info) Close() {
	defer inf.CancelTx()
	if err := inf.Tx.Tx.Commit(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Errorln("committing transaction: " + err.Error())
	}
}

// WriteOKResponse writes a 200 OK response with the given object as the
// 'response' property of the response body.
//
// This CANNOT be used by any Info that wasn't constructed for the caller by
// Wrap - ing a Handler (yet).
func (inf Info) WriteOKResponse(resp any) (int, error, error) {
	WriteResp(inf.w, inf.request, resp)
	return http.StatusOK, nil, nil
}

// WriteOKResponseWithSummary writes a 200 OK response with the given object as
// the 'response' property of the response body, and the given count as the
// `count` property of the response's summary.
//
// This CANNOT be used by any Info that wasn't constructed for the caller by
// Wrap - ing a Handler (yet).
//
// Deprecated: Summary sections on responses were intended to cover up for a
// deficiency in jQuery-based tables on the front-end, so now that we aren't
// using those anymore it serves no purpose.
func (inf Info) WriteOKResponseWithSummary(resp any, count uint64) (int, error, error) {
	WriteRespWithSummary(inf.w, inf.request, resp, count)
	return http.StatusOK, nil, nil
}

// WriteNotModifiedResponse writes a 304 Not Modified response with the given
// time as the last modified time in the headers.
//
// This CANNOT be used by any Info that wasn't constructed for the caller by
// Wrap - ing a Handler (yet).
func (inf Info) WriteNotModifiedResponse(lastModified time.Time) (int, error, error) {
	inf.w.Header().Set(rfc.LastModified, FormatLastModified(lastModified))
	inf.w.WriteHeader(http.StatusNotModified)
	setRespWritten(inf.request)
	return http.StatusNotModified, nil, nil
}

// WriteSuccessResponse writes the given response object as the `response`
// property of the response body, with the accompanying message as a
// success-level Alert.
func (inf Info) WriteSuccessResponse(resp any, message string) (int, error, error) {
	WriteAlertsObj(inf.w, inf.request, http.StatusOK, tc.CreateAlerts(tc.SuccessLevel, message), resp)
	return http.StatusOK, nil, nil
}

// WriteCreatedResponse writes the given response object as the `response`
// property of the response body of a 201 created response, with the
// accompanying message as a success-level Alert. It also sets the Location
// header to the given path. This will be automatically prefaced with the
// correct path to the API version the client requested.
func (inf Info) WriteCreatedResponse(resp any, message, path string) (int, error, error) {
	inf.w.Header().Set(rfc.Location, strings.Join([]string{"/api", inf.Version.String(), strings.TrimPrefix(path, "/")}, "/"))
	inf.w.WriteHeader(http.StatusCreated)
	WriteAlertsObj(inf.w, inf.request, http.StatusCreated, tc.CreateAlerts(tc.SuccessLevel, message), resp)
	return http.StatusCreated, nil, nil
}

// RequestHeaders returns the headers sent by the client in the API request.
func (inf Info) RequestHeaders() http.Header {
	return inf.request.Header
}

// SetLastModified sets the "last modified" header on the response writer.
//
// This CANNOT be used by any Info that wasn't constructed for the caller by
// Wrap - ing a Handler (yet).
func (inf Info) SetLastModified(t time.Time) {
	inf.w.Header().Set(rfc.LastModified, FormatLastModified(t))
}

// DecodeBody reads the client request's body and attempts to decode it into the
// provided reference.
func (inf Info) DecodeBody(ref any) error {
	return json.NewDecoder(inf.request.Body).Decode(ref)
}

// SendMail is a convenience method used to call SendMail using an Info
// structure's configuration.
func (inf *Info) SendMail(to rfc.EmailAddress, msg []byte) (int, error, error) {
	return SendMail(to, msg, inf.Config)
}

// IsResourceAuthorizedToCurrentUser is a convenience method used to call
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/tenant.IsResourceAuthorizedToUserTx
// using an Info structure to provide the current user and database transaction.
func (inf *Info) IsResourceAuthorizedToCurrentUser(resourceTenantID int) (bool, error) {
	return tenant.IsResourceAuthorizedToUserTx(resourceTenantID, inf.User, inf.Tx.Tx)
}

// CreateChangeLog creates a new changelog message at the APICHANGE level for
// the current user.
func (inf Info) CreateChangeLog(msg string) {
	_, err := inf.Tx.Tx.Exec(createChangeLogQuery, ApiChange, msg, inf.User.ID)
	if err != nil {
		log.Errorf("Inserting chage log level '%s' message '%s' for user '%s': %v", ApiChange, msg, inf.User.UserName, err)
	}
}

// UseIMS returns whether or not If-Modified-Since constraints should be used to
// service the given request.
func (inf Info) UseIMS() bool {
	if inf.request == nil || inf.Config == nil {
		return false
	}
	return inf.Config.UseIMS && inf.request.Header.Get(rfc.IfModifiedSince) != ""
}

// CreateInfluxClient constructs and returns an InfluxDB HTTP client, if enabled
// and when possible. The error this returns should not be exposed to the user;
// it's for logging purposes only.
//
// If Influx connections are not enabled, this will return `nil` - but also no
// error. It is expected that the caller will handle this situation
// appropriately.
func (inf *Info) CreateInfluxClient() (*influx.Client, error) {
	if !inf.Config.InfluxEnabled || inf.Config.ConfigInflux == nil {
		return nil, nil
	}

	var fqdn string
	var tcpPort uint
	var httpsPort sql.NullInt64 // this is the only one that's optional

	row := inf.Tx.Tx.QueryRow(influxServersQuery)
	if e := row.Scan(&fqdn, &tcpPort, &httpsPort); e != nil {
		return nil, fmt.Errorf("failed to create influx client: %w", e)
	}

	host := "http%s://%s:%d"
	if inf.Config.ConfigInflux.Secure != nil && *inf.Config.ConfigInflux.Secure {
		if !httpsPort.Valid {
			log.Warnf("INFLUXDB Server %s has no secure ports, assuming default of 8086!", fqdn)
			httpsPort = sql.NullInt64{Int64: 8086, Valid: true}
		}

		p := httpsPort.Int64
		if p <= 0 || p > 65535 {
			log.Warnf("INFLUXDB Server %s has invalid port, assuming default of 8086!", fqdn)
			p = 8086
		}

		host = fmt.Sprintf(host, "s", fqdn, p)
	} else if tcpPort > 0 && tcpPort <= 65535 {
		host = fmt.Sprintf(host, "", fqdn, tcpPort)
	} else {
		log.Warnf("INFLUXDB Server %s has invalid port, assuming default of 8086!", fqdn)
		host = fmt.Sprintf(host, "", fqdn, 8086)
	}

	config := influx.HTTPConfig{
		Addr:      host,
		Username:  inf.Config.ConfigInflux.User,
		Password:  inf.Config.ConfigInflux.Password,
		UserAgent: fmt.Sprintf("TrafficOps/%s (Go)", inf.Config.Version),
		Timeout:   time.Duration(float64(inf.Config.ReadTimeout)/2.1) * time.Second,
	}

	var client influx.Client
	client, e := influx.NewHTTPClient(config)
	if e != nil {
		return nil, fmt.Errorf("failed to create influx client: %w", e)
	}
	if client == nil {
		return nil, errors.New("failed to create influx client: client was nil")
	}
	return &client, e
}

// DefaultSort sets the `orderby` query string parameter to the given value, as
// though the client had set it, should it be missing.
func (inf Info) DefaultSort(param string) {
	if _, ok := inf.Params["orderby"]; !ok {
		inf.Params["orderby"] = param
	}
}
