package iso

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
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

// ISOs handler is responsible for generating and returning an ISO image.
//
// Response types:
//
// Error:
//   HTTP 400
//   {
//     "alerts": [
//       {"level":"error","text":"hostName is required"},
//       {"level":"error","text":"disk is required"},
//       ...,
//     ]
//   }
//
// Success (streaming = false):
//   HTTP 200
//   {
//     "alerts": [
//       {"level":"success","text":"Generate ISO was successful."}
//     ],
//     "response": {
//       "isoURL":"https:\/\/trafficops-perl.infra.ciab.test\/iso\/db.infra.ciab.test-centos72.iso",
//       "isoName":"db.infra.ciab.test-centos72.iso"
//     }
//   }
//
// Success (streaming = true):
//   HTTP 200
//   Content-Disposition: attachment; filename="db.infra.ciab.test-centos72.iso"
//   Content-Type: application/download
func ISOs(w http.ResponseWriter, req *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(req, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, req, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var ir isoRequest
	if err := json.NewDecoder(req.Body).Decode(&ir); err != nil {
		api.HandleErr(w, req, inf.Tx.Tx, http.StatusBadRequest, errors.New("unable to process request"), err)
		return
	}

	level, message, resp, err := isos(w, req)
	if err != nil {
		api.HandleErr(w, req, inf.Tx.Tx, 500, errors.New("user error"), errors.New("sys error"))
		return
	}
	api.WriteRespAlertObj(w, req, level, message, resp)
}

func isos(w http.ResponseWriter, req *http.Request) (tc.AlertLevel, string, interface{}, error) {
	return tc.SuccessLevel, "this is the message", new(interface{}), nil
}

type isoRequest struct {
	OSVersionDir  string  `json:"osversionDir"`
	HostName      string  `json:"hostName"`
	DomainName    string  `json:"domainName"`
	RootPass      string  `json:"rootPass"`
	DHCP          boolStr `json:"dhcp"`
	IPAddr        net.IP  `json:"ipAddress"`
	IPNetmask     net.IP  `json:"ipNetmask"`
	IPGateway     net.IP  `json:"ipGateway"`
	IP6Address    net.IP  `json:"ip6Address"`
	IP6Gateway    net.IP  `json:"ip6Gateway"`
	InterfaceName string  `json:"interfaceName"`
	InterfaceMTU  int     `json:"interfaceMtu"`
	Disk          string  `json:"disk"`
	MgmtIPAddress net.IP  `json:"mgmtIpAddress"`
	MgmtIPNetmask net.IP  `json:"mgmtIpNetmask"`
	MgmtIPGateway net.IP  `json:"mgmtIpGateway"`
	MgmtInterface string  `json:"mgmtInterface"`
	Stream        boolStr `json:"stream"`
}

func (i *isoRequest) validate() []error {
	/*
		# Validation checks to perform
		checks => [
			osversionDir  => [ is_required("is required") ],
			hostName      => [ is_required("is required") ],
			domainName    => [ is_required("is required") ],
			rootPass      => [ is_required("is required") ],
			dhcp          => [ is_required("is required") ],
			interfaceMtu  => [ is_required("is required") ],
			disk          => [ is_required("is required") ],
			mgmtInterface => [ is_required_if((defined($mgmtIpAddress) && $mgmtIpAddress ne ""), "- Management interface is required when Management IP is provided") ],
			mgmtIpGateway => [ is_required_if((defined($mgmtIpAddress) && $mgmtIpAddress ne ""), "- Management gateway is required when Management IP is provided") ],
			ipAddress     => is_required_if(
				sub {
					my $params = shift;
					return $params->{dhcp} eq 'no';
				},
				"is required if DHCP is no"
			),
			ipNetmask => is_required_if(
				sub {
					my $params = shift;
					return $params->{dhcp} eq 'no';
				},
				"is required if DHCP is no"
			),
			ipGateway => is_required_if(
				sub {
					my $params = shift;
					return $params->{dhcp} eq 'no';
				},
				"is required if DHCP is no"
			),
		]
	*/
	var errs []error
	addErr := func(msg string) { errs = append(errs, errors.New(msg)) }

	if i.OSVersionDir == "" {
		addErr("osversionDir is required")
	}
	if i.HostName == "" {
		addErr("hostName is required")
	}
	if i.DomainName == "" {
		addErr("domainName is required")
	}
	if i.RootPass == "" {
		addErr("rootPass is required")
	}
	if !i.DHCP.isSet {
		addErr("dhcp is required")
	}
	if i.InterfaceMTU == 0 {
		addErr("interfaceMtu is required")
	}
	if i.Disk == "" {
		addErr("disk is required")
	}
	if len(i.MgmtIPAddress) > 0 {
		if i.MgmtInterface == "" {
			addErr("mgmtInterface is required when mgmtIpAddress is provided")
		}
		if len(i.MgmtIPGateway) == 0 {
			addErr("mgmtIpGateway is required when mgmtIpAddress is provided")
		}
	}
	if i.DHCP.val == false {
		if len(i.IPAddr) == 0 {
			addErr("ipAddress is required if DHCP is no")
		}
		if len(i.IPNetmask) == 0 {
			addErr("ipNetmask is required if DHCP is no")
		}
		if len(i.IPGateway) == 0 {
			addErr("ipGateway is required if DHCP is no")
		}
	}

	return errs
}

// boolStr is used to decode boolean strings (e.g. "yes") as
// part of a JSON response. Part of the /isos JSON request
// generated by TrafficPortal uses this format.
// If an unrecognize or empty string is given, then
// the 'val' and 'isSet' fields will be false. Otherwise,
// 'isSet' will be true.
type boolStr struct {
	isSet bool // false if UnmarshalText is given an unrecognized value
	val   bool
}

// UnmarshalText decodes strings representing boolean values.
// It nevers returns an error to allow for all validation errors
// to be grouped together.
func (b *boolStr) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "yes", "true", "1":
		b.val = true
		b.isSet = true
	case "no", "false", "0":
		b.val = false
		b.isSet = true
	}
	return nil
}
