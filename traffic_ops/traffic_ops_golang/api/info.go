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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficvault"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
)

// Info is a structure produced from a client's request that provides (nearly)
// all of the information and functionality needed to service that request.
type Info struct {
	// Params is a map of request "parameters" to their values. Request
	// parameters are those found either in the query string - such as 'foo'
	// with a value of 'bar' in /api/4.0/servers?foo=bar - or as required route
	// parameters - such as 'test' with a value of 'quest' in
	// /api/4.0/servers/quest (assuming the route definition was
	// /api/4.0/servers/{test}).
	Params map[string]string
	// IntParams is a map of request "parameters" to their values - but ONLY if
	// those values are integers. Which parameters should be integers is
	// typically determined by the arguments to 'NewInfo'.
	IntParams map[string]int
	// The currently authenticated user, if and when the client is
	// authenticated. For routes that require authentication, this *should* be
	// non-nil, assuming the Info was properly generated.
	User *auth.CurrentUser
	// ReqID is a unique ID for the request to which this Info belongs.
	ReqID uint64
	// Version specifies the API version requested by the client.
	Version Version
	// Tx is a reference to an open database transaction built to service the
	// request. It will be closed with the Info itself.
	Tx *sqlx.Tx
	// Config is a reference to the Traffic Ops server's configuration.
	Config *config.Config
	// Vault implements the interaction interface for Traffic Vault.
	Vault     trafficvault.TrafficVault
	request   *http.Request
	writer    http.ResponseWriter
	ctxCancel context.CancelFunc
}

// NewInfo constructs Info needed by handlers from a client request. It also
// returns any user error, any system error, and the status code which should
// be returned to the client if an error occurred. The Info pointer returned is
// guaranteed to not be 'nil'.
//
// It is encouraged to call Info.Tx.Tx.Commit() manually when all queries are
// finished, to release database resources early, and also to return an error
// to the user if the commit failed.
//
// NewInfo guarantees the returned Info.Tx is non-nil and Info.Tx.Tx is nil or
// valid, even if a returned error is not nil. Hence, it is safe to pass the
// Tx.Tx to HandleErr when this returns errors.
//
// Close() must be called to free resources, and should be called in a defer
// immediately after NewInfo(), to finish the transaction.
//
// Example:
//  func handler(w http.ResponseWriter, r *http.Request) {
//    inf, userErr, sysErr, errCode := api.NewInfo(w, r, nil, nil)
//    if userErr != nil || sysErr != nil {
//      api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
//      return
//    }
//    defer inf.Close()
//
//    respObj, err := finalDatabaseOperation(inf.Tx)
//    if err != nil {
//      api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("final db op: " + err.Error()))
//      return
//    }
//    if err := inf.Tx.Tx.Commit(); err != nil {
//      api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("committing transaction: " + err.Error()))
//      return
//    }
//    api.WriteResp(w, r, respObj)
//  }
//
func NewInfo(w http.ResponseWriter, r *http.Request, requiredParams, intParamNames []string) (*Info, error, error, int) {
	db, err := GetDB(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, errors.New("getting db: " + err.Error()), nil, http.StatusInternalServerError
	}
	cfg, err := GetConfig(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, errors.New("getting config: " + err.Error()), nil, http.StatusInternalServerError
	}
	reqID, err := getReqID(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, errors.New("getting reqID: " + err.Error()), nil, http.StatusInternalServerError
	}
	version := getRequestedAPIVersion(r.URL.Path)

	user, err := auth.GetCurrentUser(r.Context())
	if err != nil {
		return &Info{Tx: &sqlx.Tx{}}, errors.New("getting user: " + err.Error()), nil, http.StatusInternalServerError
	}
	params, intParams, userErr, sysErr, errCode := AllParams(r, requiredParams, intParamNames)
	if userErr != nil || sysErr != nil {
		return &Info{Tx: &sqlx.Tx{}}, userErr, sysErr, errCode
	}
	dbCtx, ctxCancel := context.WithTimeout(r.Context(), time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second) //only place we could call cancel here is in Info.Close(), which already will rollback the transaction (which is all cancel will do.)
	tx, err := db.BeginTxx(dbCtx, nil)                                                                         // must be last, MUST not return an error if this succeeds, without closing the tx
	if err != nil {
		ctxCancel()
		return &Info{Tx: &sqlx.Tx{}}, userErr, errors.New("could not begin transaction: " + err.Error()), http.StatusInternalServerError
	}
	return &Info{
		Config:    cfg,
		ReqID:     reqID,
		Version:   version,
		Params:    params,
		IntParams: intParams,
		User:      user,
		Tx:        tx,
		request:   r,
		writer:    w,
		ctxCancel: ctxCancel,
	}, nil, nil, http.StatusOK
}

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
		return http.StatusInternalServerError, nil, fmt.Errorf("scanning for lastUpdated: %v", err)
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

// Close implements the io.Closer interface. It should be called in a defer
// immediately after NewInfo().
//
// Close will commit the transaction, if it hasn't been rolled back.
func (inf *Info) Close() {
	defer inf.ctxCancel()
	if err := inf.Tx.Tx.Commit(); err != nil && err != sql.ErrTxDone {
		log.Errorln("committing transaction: " + err.Error())
	}
}

// SendMail is a convenience method used to call SendMail using an Info
// structure's configuration.
func (inf *Info) SendMail(to rfc.EmailAddress, msg []byte) (int, error, error) {
	return SendMail(to, msg, inf.Config)
}

// IsResourceAuthorizedToCurrentUser is a convenience method used to call
// github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant.IsResourceAuthorizedToUserTx
// using an Info structure to provide the current user and database transaction.
func (inf *Info) IsResourceAuthorizedToCurrentUser(resourceTenantID int) (bool, error) {
	return tenant.IsResourceAuthorizedToUserTx(resourceTenantID, inf.User, inf.Tx.Tx)
}

// CreateInfluxClient constructs and returns an InfluxDB HTTP client, if enabled and when possible.
// The error this returns should not be exposed to the user; it's for logging purposes only.
//
// If Influx connections are not enabled, this will return `nil` - but also no error. It is expected
// that the caller will handle this situation appropriately.
func (inf *Info) CreateInfluxClient() (*influx.Client, error) {
	if !inf.Config.InfluxEnabled {
		return nil, nil
	}

	var fqdn string
	var tcpPort uint
	var httpsPort sql.NullInt64 // this is the only one that's optional

	row := inf.Tx.Tx.QueryRow(influxServersQuery)
	if e := row.Scan(&fqdn, &tcpPort, &httpsPort); e != nil {
		return nil, fmt.Errorf("Failed to create influx client: %v", e)
	}

	host := "http%s://%s:%d"
	if inf.Config.ConfigInflux != nil && *inf.Config.ConfigInflux.Secure {
		if !httpsPort.Valid {
			log.Warnf("INFLUXDB Server %s has no secure ports, assuming default of 8086!", fqdn)
			httpsPort = sql.NullInt64{Int64: 8086, Valid: true}
		}
		port, err := httpsPort.Value()
		if err != nil {
			return nil, fmt.Errorf("Failed to create influx client: %v", err)
		}

		p := port.(int64)
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
	if client == nil {
		return nil, fmt.Errorf("Failed to create influx client (client was nil): %v", e)
	}
	return &client, e
}
