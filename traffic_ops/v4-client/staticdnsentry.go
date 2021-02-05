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
)

const (
	API_STATIC_DNS_ENTRIES = apiBase + "/staticdnsentries"
)

func staticDNSEntryIDs(to *Session, sdns *tc.StaticDNSEntry) error {
	if sdns.CacheGroupID == 0 && sdns.CacheGroupName != "" {
		p, _, err := to.GetCacheGroupNullableByNameWithHdr(sdns.CacheGroupName, nil)
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
		dses, _, err := to.GetDeliveryServiceByXMLIDNullableWithHdr(sdns.DeliveryService, nil)
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
		types, _, err := to.GetTypeByNameWithHdr(sdns.Type, nil)
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

// CreateStaticDNSEntry creates a Static DNS Entry.
func (to *Session) CreateStaticDNSEntry(sdns tc.StaticDNSEntry) (tc.Alerts, ReqInf, error) {
	// fill in missing IDs from names
	var alerts tc.Alerts
	err := staticDNSEntryIDs(to, &sdns)
	if err != nil {
		return alerts, ReqInf{CacheHitStatus: CacheHitStatusMiss}, err
	}
	reqInf, err := to.post(API_STATIC_DNS_ENTRIES, sdns, nil, &alerts)
	return alerts, reqInf, err
}

func (to *Session) UpdateStaticDNSEntryByIDWithHdr(id int, sdns tc.StaticDNSEntry, header http.Header) (tc.Alerts, ReqInf, int, error) {
	// fill in missing IDs from names
	var alerts tc.Alerts
	err := staticDNSEntryIDs(to, &sdns)
	if err != nil {
		return alerts, ReqInf{CacheHitStatus: CacheHitStatusMiss}, 0, err
	}
	route := fmt.Sprintf("%s?id=%d", API_STATIC_DNS_ENTRIES, id)
	reqInf, err := to.put(route, sdns, header, &alerts)
	return tc.Alerts{}, reqInf, reqInf.StatusCode, err
}

// UpdateStaticDNSEntryByID updates a Static DNS Entry by ID.
// Deprecated: UpdateStaticDNSEntryByID will be removed in 6.0. Use UpdateStaticDNSEntryByIDWithHdr.
func (to *Session) UpdateStaticDNSEntryByID(id int, sdns tc.StaticDNSEntry) (tc.Alerts, ReqInf, int, error) {
	return to.UpdateStaticDNSEntryByIDWithHdr(id, sdns, nil)
}

func (to *Session) GetStaticDNSEntriesWithHdr(header http.Header) ([]tc.StaticDNSEntry, ReqInf, error) {
	var data tc.StaticDNSEntriesResponse
	reqInf, err := to.get(API_STATIC_DNS_ENTRIES, header, &data)
	return data.Response, reqInf, err
}

// GetStaticDNSEntries returns a list of Static DNS Entrys.
// Deprecated: GetStaticDNSEntries will be removed in 6.0. Use GetStaticDNSEntriesWithHdr.
func (to *Session) GetStaticDNSEntries() ([]tc.StaticDNSEntry, ReqInf, error) {
	return to.GetStaticDNSEntriesWithHdr(nil)
}

func (to *Session) GetStaticDNSEntryByIDWithHdr(id int, header http.Header) ([]tc.StaticDNSEntry, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_STATIC_DNS_ENTRIES, id)
	var data tc.StaticDNSEntriesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetStaticDNSEntryByID GETs a Static DNS Entry by the Static DNS Entry's ID.
// Deprecated: GetStaticDNSEntryByID will be removed in 6.0. Use GetStaticDNSEntryByIDWithHdr.
func (to *Session) GetStaticDNSEntryByID(id int) ([]tc.StaticDNSEntry, ReqInf, error) {
	return to.GetStaticDNSEntryByIDWithHdr(id, nil)
}

func (to *Session) GetStaticDNSEntriesByHostWithHdr(host string, header http.Header) ([]tc.StaticDNSEntry, ReqInf, error) {
	route := fmt.Sprintf("%s?host=%s", API_STATIC_DNS_ENTRIES, url.QueryEscape(host))
	var data tc.StaticDNSEntriesResponse
	reqInf, err := to.get(route, header, &data)
	return data.Response, reqInf, err
}

// GetStaticDNSEntriesByHost GETs a Static DNS Entry by the Static DNS Entry's host.
// Deprecated: GetStaticDNSEntriesByHost will be removed in 6.0. Use GetStaticDNSEntriesByHostWithHdr.
func (to *Session) GetStaticDNSEntriesByHost(host string) ([]tc.StaticDNSEntry, ReqInf, error) {
	return to.GetStaticDNSEntriesByHostWithHdr(host, nil)
}

// DeleteStaticDNSEntryByID DELETEs a Static DNS Entry by ID.
func (to *Session) DeleteStaticDNSEntryByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_STATIC_DNS_ENTRIES, id)
	var alerts tc.Alerts
	reqInf, err := to.del(route, nil, &alerts)
	return alerts, reqInf, err
}
