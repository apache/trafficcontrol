package request

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
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/jmoiron/sqlx"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var ds = tc.DeliveryServiceV5{
	Active:                   tc.DeliveryServiceActiveState("PRIMED"),
	AnonymousBlockingEnabled: false,
	CCRDNSTTL:                util.IntPtr(20),
	CDNID:                    11,
	CDNName:                  util.StrPtr("testCDN"),
	CheckPath:                util.StrPtr("blah"),
	ConsistentHashRegex:      nil,
	DeepCachingType:          tc.DeepCachingTypeNever,
	DisplayName:              "ds",
	DNSBypassCNAME:           nil,
	DNSBypassIP:              nil,
	DNSBypassIP6:             nil,
	DNSBypassTTL:             nil,
	DSCP:                     0,
	EcsEnabled:               false,
	EdgeHeaderRewrite:        nil,
	FirstHeaderRewrite:       nil,
	GeoLimitRedirectURL:      nil,
	GeoLimit:                 0,
	GeoLimitCountries:        nil,
	GeoProvider:              0,
	GlobalMaxMBPS:            nil,
	GlobalMaxTPS:             nil,
	FQPacingRate:             nil,
	HTTPBypassFQDN:           nil,
	ID:                       util.IntPtr(1),
	InfoURL:                  nil,
	InitialDispersion:        util.IntPtr(1),
	InnerHeaderRewrite:       nil,
	IPV6RoutingEnabled:       util.BoolPtr(true),
	LastHeaderRewrite:        nil,
	LastUpdated:              time.Now(),
	LogsEnabled:              true,
	LongDesc:                 "",
	MaxDNSAnswers:            util.IntPtr(5),
	MaxOriginConnections:     util.IntPtr(2),
	MaxRequestHeaderBytes:    util.IntPtr(0),
	MidHeaderRewrite:         nil,
	MissLat:                  util.FloatPtr(0.0),
	MissLong:                 util.FloatPtr(0.0),
	MultiSiteOrigin:          false,
	OrgServerFQDN:            util.StrPtr("http://1.2.3.4"),
	OriginShield:             nil,
	ProfileID:                util.IntPtr(99),
	ProfileName:              util.StrPtr("profile99"),
	ProfileDesc:              nil,
	Protocol:                 util.IntPtr(1),
	QStringIgnore:            nil,
	RangeRequestHandling:     nil,
	RegexRemap:               nil,
	Regional:                 false,
	RegionalGeoBlocking:      false,
	RemapText:                nil,
	RequiredCapabilities:     nil,
	RoutingName:              "",
	ServiceCategory:          nil,
	SigningAlgorithm:         nil,
	RangeSliceBlockSize:      nil,
	SSLKeyVersion:            nil,
	TenantID:                 100,
	Tenant:                   util.StrPtr("tenant100"),
	TLSVersions:              nil,
	Topology:                 nil,
	TRRequestHeaders:         nil,
	TRResponseHeaders:        nil,
	Type:                     util.StrPtr("type101"),
	TypeID:                   101,
	XMLID:                    "dsXMLID",
}

var dsr = tc.DeliveryServiceRequestNullable{
	AssigneeID:      nil,
	Assignee:        nil,
	AuthorID:        nil,
	Author:          nil,
	ChangeType:      nil,
	CreatedAt:       nil,
	ID:              nil,
	LastEditedBy:    nil,
	LastEditedByID:  nil,
	LastUpdated:     nil,
	DeliveryService: nil,
	Status:          util.Ptr(tc.RequestStatusDraft),
	XMLID:           util.StrPtr("dsXMLID"),
}

func TestValidateV5(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("opening mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()

	dsrV5 := dsr.Upgrade().Upgrade().Upgrade()
	userErr, sysErr := validateV5(dsrV5, db.MustBegin().Tx)

	if sysErr != nil {
		t.Fatalf("expected no error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Fatalf("expected userErr because change type is absent, but got nothing")
	}

	dsrV5.ChangeType = tc.DSRChangeTypeCreate
	mock.ExpectBegin()
	userErr, sysErr = validateV5(dsrV5, db.MustBegin().Tx)
	if sysErr != nil {
		t.Fatalf("expected no error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Fatalf("expected userErr because requested is absent for changetype 'change', but got nothing")
	}

	dsrV5.Requested = &ds
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{
		"name",
		"use_in_table",
	})
	rows.AddRow("type101", "server")
	mock.ExpectQuery("SELECT name, use_in_table*").WillReturnRows(rows)
	userErr, sysErr = validateV5(dsrV5, db.MustBegin().Tx)
	if sysErr != nil {
		t.Fatalf("expected no error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Fatalf("expected userErr because use_in_table is not deliveryservice, but got nothing")
	}

	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{
		"name",
		"use_in_table",
	})
	rows.AddRow("type101", "deliveryservice")
	mock.ExpectQuery("SELECT name, use_in_table*").WillReturnRows(rows)
	userErr, sysErr = validateV5(dsrV5, db.MustBegin().Tx)
	if userErr != nil || sysErr != nil {
		t.Fatalf("no error expected, but got usererr: %v, sysErr: %v", userErr, sysErr)
	}

	dsrV5.ChangeType = tc.DSRChangeTypeDelete
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{
		"name",
		"use_in_table",
	})
	rows.AddRow("type101", "deliveryservice")
	userErr, sysErr = validateV5(dsrV5, db.MustBegin().Tx)
	if sysErr != nil {
		t.Fatalf("expected no error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Fatalf("expected userErr because original is not present for changetype 'delete', but got nothing")
	}

	dsrV5.Requested = nil
	dsrV5.Original = &ds
	mock.ExpectBegin()
	userErr, sysErr = validateV5(dsrV5, db.MustBegin().Tx)
	if userErr != nil || sysErr != nil {
		t.Fatalf("no error expected, but got usererr: %v, sysErr: %v", userErr, sysErr)
	}

	dsrV5.Assignee = util.StrPtr("testUser")
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{
		"id",
	})
	rows.AddRow(10)
	mock.ExpectQuery("SELECT id FROM tm_user*").WillReturnRows(rows)
	userErr, sysErr = validateV5(dsrV5, db.MustBegin().Tx)
	if userErr != nil || sysErr != nil {
		t.Fatalf("no error expected, but got usererr: %v, sysErr: %v", userErr, sysErr)
	}
}

func TestValidateLegacy(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("opening mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	mock.ExpectBegin()

	dsV30 := ds.Downgrade().DowngradeToV31()
	dsr.DeliveryService = &dsV30
	// expect error because ChangeType is absent
	userErr, sysErr := validateLegacy(dsr, db.MustBegin().Tx)
	if sysErr != nil {
		t.Fatalf("expected no error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Fatalf("expected userErr because change type is absent, but got nothing")
	}
	dsr.ChangeType = util.StrPtr(string(tc.DSRChangeTypeCreate))
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{
		"name",
		"use_in_table",
	})
	rows.AddRow("type101", "server")
	mock.ExpectQuery("SELECT name, use_in_table*").WillReturnRows(rows)
	userErr, sysErr = validateLegacy(dsr, db.MustBegin().Tx)
	if sysErr != nil {
		t.Fatalf("expected no error, but got sysErr: %v", sysErr)
	}
	if userErr == nil {
		t.Fatalf("expected userErr because use_in_table is not deliveryservice, but got nothing")
	}

	mock.ExpectBegin()
	rows.AddRow("type101", "deliveryservice")
	mock.ExpectQuery("SELECT name, use_in_table*").WillReturnRows(rows)
	userErr, sysErr = validateLegacy(dsr, db.MustBegin().Tx)
	if userErr != nil || sysErr != nil {
		t.Fatalf("no error expected, but got usererr: %v, sysErr: %v", userErr, sysErr)
	}

	dsr.ID = util.IntPtr(1)
	dsr.Status = util.Ptr(tc.RequestStatusSubmitted)
	mock.ExpectBegin()

	rows2 := sqlmock.NewRows([]string{
		"status",
	})
	rows2.AddRow([]byte("submitted"))
	mock.ExpectQuery("SELECT status*").WillReturnRows(rows2)

	rows.AddRow("type101", "deliveryservice")
	mock.ExpectQuery("SELECT name, use_in_table*").WillReturnRows(rows)
	userErr, sysErr = validateLegacy(dsr, db.MustBegin().Tx)
	if userErr != nil || sysErr != nil {
		t.Fatalf("no error expected, but got usererr: %v, sysErr: %v", userErr, sysErr)
	}
}
