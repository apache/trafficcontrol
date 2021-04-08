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

package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/toclientlib"
)

const (
	// APIStaticDNSEntries is the full path to the /staticdnsentries API
	// endpoint.
	APIStaticDNSEntries = "/staticdnsentries"
)

func staticDNSEntryIDs(to *Session, sdns *tc.StaticDNSEntry) error {
	if sdns.CacheGroupID == 0 && sdns.CacheGroupName != "" {
		p, _, err := to.GetCacheGroupByName(sdns.CacheGroupName, nil)
		if err != nil {
			return err
		}
		if len(p) == 0 {
			return errors.New("no CacheGroup named " + sdns.CacheGroupName)
		}
		if p[0].ID == nil {
			return errors.New("CacheGroup named " + sdns.CacheGroupName + " has a nil ID")
		}
		sdns.CacheGroupID = *p[0].ID
	}

	if sdns.DeliveryServiceID == 0 && sdns.DeliveryService != "" {
		dses, _, err := to.GetDeliveryServiceByXMLID(sdns.DeliveryService, nil)
		if err != nil {
			return err
		}
		if len(dses) == 0 {
			return errors.New("no deliveryservice with name " + sdns.DeliveryService)
		}
		if dses[0].ID == nil {
			return errors.New("Deliveryservice with name " + sdns.DeliveryService + " has a nil ID")
		}
		sdns.DeliveryServiceID = *dses[0].ID
	}

	if sdns.TypeID == 0 && sdns.Type != "" {
		types, _, err := to.GetTypeByName(sdns.Type, nil)
		if err != nil {
			return err
		}
		if len(types) == 0 {
			return errors.New("no type with name " + sdns.Type)
		}
		sdns.TypeID = types[0].ID
	}

	return nil
}

// CreateStaticDNSEntry creates the given Static DNS Entry.
func (to *Session) CreateStaticDNSEntry(sdns tc.StaticDNSEntry) (tc.Alerts, toclientlib.ReqInf, error) {
	// fill in missing IDs from names
	var alerts tc.Alerts
	err := staticDNSEntryIDs(to, &sdns)
	if err != nil {
		return alerts, toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}, err
	}
	reqInf, err := to.post(APIStaticDNSEntries, sdns, nil, &alerts)
	return alerts, reqInf, err
}

// UpdateStaticDNSEntry replaces the Static DNS Entry identified by 'id' with
// the one provided.
func (to *Session) UpdateStaticDNSEntry(id int, sdns tc.StaticDNSEntry, header http.Header) (tc.Alerts, toclientlib.ReqInf, int, error) {
	// fill in missing IDs from names
	var alerts tc.Alerts
	err := staticDNSEntryIDs(to, &sdns)
	if err != nil {
		return alerts, toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}, 0, err
	}
	route := fmt.Sprintf("%s?id=%d", APIStaticDNSEntries, id)
	reqInf, err := to.put(route, sdns, header, &alerts)
	return tc.Alerts{}, reqInf, reqInf.StatusCode, err
}

// GetStaticDNSEntries retrieves all Static DNS Entries stored in Traffic Ops.
func (to *Session) GetStaticDNSEntries(header http.Header) ([]tc.StaticDNSEntry, toclientlib.ReqInf, error) {
	var data tc.StaticDNSEntriesResponse
	reqInf, err := to.get(APIStaticDNSEntries, header, &data)
	return data.Response, reqInf, err
}

// GetStaticDNSEntryByID retrieves the Static DNS Entry with the given ID.
func (to *Session) GetStaticDNSEntryByID(id int, header http.Header) ([]tc.StaticDNSEntry, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIStaticDNSEntries, id)
	var data tc.StaticDNSEntriesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetStaticDNSEntriesByHost retrieves all Static DNS Entries stored in Traffic Ops
// with the given Host.
func (to *Session) GetStaticDNSEntriesByHost(host string, header http.Header) ([]tc.StaticDNSEntry, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?host=%s", APIStaticDNSEntries, url.QueryEscape(host))
	var data tc.StaticDNSEntriesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// DeleteStaticDNSEntry deletes the Static DNS Entry with the given ID.
func (to *Session) DeleteStaticDNSEntry(id int) (tc.Alerts, toclientlib.ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", APIStaticDNSEntries, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
