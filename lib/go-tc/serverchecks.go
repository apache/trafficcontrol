package tc

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
	"database/sql"
	"errors"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// ServerchecksResponse is a list of Serverchecks as a response.
// swagger:response ServerchecksResponse
// in: body
type ServerchecksResponse struct {
	// in: body
	Response []Servercheck `json:"response"`
	Alerts
}

// CommonCheckFields is a structure containing all of the fields common to both
// Serverchecks and GenericServerChecks.
type CommonCheckFields struct {

	// AdminState is the server's status - called "AdminState" for legacy reasons.
	AdminState string `json:"adminState"`

	// CacheGroup is the name of the Cache Group to which the server belongs.
	CacheGroup string `json:"cacheGroup"`

	// ID is the integral, unique identifier of the server.
	ID int `json:"id"`

	// HostName of the checked server.
	HostName string `json:"hostName"`

	// RevalPending is a flag that indicates if revalidations are pending for the checked server.
	RevalPending bool `json:"revalPending"`

	// Profile is the name of the Profile used by the checked server.
	Profile string `json:"profile"`

	// Type is the name of the server's Type.
	Type string `json:"type"`

	// UpdPending is a flag that indicates if updates are pending for the checked server.
	UpdPending bool `json:"updPending"`
}

// Servercheck is a single Servercheck struct for GET response.
// swagger:model Servercheck
type Servercheck struct {
	CommonCheckFields

	// Various check types
	Checks struct {

		// IPv4 production interface (legacy name)
		Iface10G int `json:"10G"`

		// IPv6 production interface (legacy name)
		Iface10G6 int `json:"10G6"`

		// Cache Disk Usage
		CDU int `json:"CDU"`

		// Cache Hit Ratio
		CHR int `json:"CHR"`

		// DSCP check
		DSCP int `json:"DSCP"`

		// DNS check
		FQDN int `json:"FQDN"`

		// Out-of-band (BMC) interface check
		ILO int `json:"ILO"`

		// IPv4 production interface (new name)
		IPv4 int `json:"IPv4"`

		// IPv6 production interface (new name)
		IPv6 int `json:"IPv6"`

		// MTU check
		MTU int `json:"MTU"`

		// ORT check
		ORT int `json:"ORT"`

		// Traffic Router status for checked server
		RTR int `json:"RTR"`
	} `json:"checks"`
}

// ServercheckPost is a single Servercheck struct for Update and Create to
// depict what changed.
type ServercheckPost struct {

	// The Servercheck data to submit
	//
	// Name of the server check type
	//
	// required: true
	Name string `json:"servercheck_short_name"`

	// ID of the server
	//
	ID int `json:"id"`

	// Name of the server
	HostName string `json:"name" `

	// Value of the check result
	//
	// required: true
	Value int `json:"value"`
}

// ServercheckRequestNullable is a single nullable Servercheck struct for Update
// and Create to depict what changed.
type ServercheckRequestNullable struct {
	Name     *string `json:"servercheck_short_name"`
	ID       *int    `json:"id"`
	Value    *int    `json:"value"`
	HostName *string `json:"host_name"`
}

// Validate implements the
// github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api.ParseValidate
// interface.
func (scp ServercheckRequestNullable) Validate(tx *sql.Tx) error {
	errs := []string{}

	if scp.ID == nil && scp.HostName == nil {
		errs = append(errs, "id or host_name")
	}

	if scp.Name == nil || *scp.Name == "" {
		errs = append(errs, "servercheck_short_name")
	}

	if len(errs) > 0 {
		return util.JoinErrs([]error{errors.New("required fields missing: " + strings.Join(errs, ", "))})
	}
	return nil
}

// ServercheckPostResponse is the response to a Servercheck POST request.
type ServercheckPostResponse struct {
	Alerts []Alert `json:"alerts"`
}

// GenericServerCheck represents a server with some associated meta data presented
// along with its arbitrary "checks". This is unlike a Servercheck in that the
// represented checks are not known before the time a request is made, and checks
// with no value are not presented.
type GenericServerCheck struct {
	CommonCheckFields

	// Checks maps arbitrary checks - up to one per "column" (whatever those mean)
	// done on the server to their values.
	Checks map[string]*int `json:"checks,omitempty"`
}

// ServercheckAPIResponse (not to be confused with ServerchecksResponse) is the
// type of a response from Traffic Ops to a request to its /servercheck
// endpoint (not to be confused with its /servers/checks endpoint).
type ServercheckAPIResponse struct {
	Response []GenericServerCheck `json:"response"`
	Alerts
}

// ServerCheckColumns is a collection of columns associated with a particular
// server's "checks". The meaning of the column names is unknown.
type ServerCheckColumns struct {
	// ID uniquely identifies a servercheck columns row.
	ID int `db:"id"`

	// Server is the ID of the server which is associated with these checks.
	Server int `db:"server"`

	AA *int `db:"aa"`
	AB *int `db:"ab"`
	AC *int `db:"ac"`
	AD *int `db:"ad"`
	AE *int `db:"ae"`
	AF *int `db:"af"`
	AG *int `db:"ag"`
	AH *int `db:"ah"`
	AI *int `db:"ai"`
	AJ *int `db:"aj"`
	AK *int `db:"ak"`
	AL *int `db:"al"`
	AM *int `db:"am"`
	AN *int `db:"an"`
	AO *int `db:"ao"`
	AP *int `db:"ap"`
	AQ *int `db:"aq"`
	AR *int `db:"ar"`
	AT *int `db:"at"`
	AU *int `db:"au"`
	AV *int `db:"av"`
	AW *int `db:"aw"`
	AX *int `db:"ax"`
	AY *int `db:"ay"`
	AZ *int `db:"az"`
	BA *int `db:"ba"`
	BB *int `db:"bb"`
	BC *int `db:"bc"`
	BD *int `db:"bd"`
	BE *int `db:"be"`
	BF *int `db:"bf"`
}
