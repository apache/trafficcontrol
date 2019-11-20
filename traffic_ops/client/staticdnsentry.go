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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

const (
	API_v13_StaticDNSEntries = "/api/1.3/staticdnsentries"
)

func staticDNSEntryIDs(to *Session, sdns *tc.StaticDNSEntry) error {
	if sdns.CacheGroupID == 0 && sdns.CacheGroupName != "" {
		p, _, err := to.GetCacheGroupByName(sdns.CacheGroupName)
		if err != nil {
			return err
		}
		if len(p) == 0 {
			return errors.New("no CacheGroup named " + sdns.CacheGroupName)
		}
		sdns.CacheGroupID = p[0].ID
	}

	if sdns.DeliveryServiceID == 0 && sdns.DeliveryService != "" {
		dses, _, err := to.GetDeliveryServiceByXMLID(sdns.DeliveryService)
		if err != nil {
			return err
		}
		if len(dses) == 0 {
			return errors.New("no deliveryservice with name " + sdns.DeliveryService)
		}
		sdns.DeliveryServiceID = dses[0].ID
	}

	if sdns.TypeID == 0 && sdns.Type != "" {
		types, _, err := to.GetTypeByName(sdns.Type)
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

// Create a StaticDNSEntry
func (to *Session) CreateStaticDNSEntry(sdns tc.StaticDNSEntry) (tc.Alerts, ReqInf, error) {
	// fill in missing IDs from names
	staticDNSEntryIDs(to, &sdns)
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(sdns)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	resp, remoteAddr, err := to.request(http.MethodPost, API_v13_StaticDNSEntries, reqBody)
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}

// Update a StaticDNSEntry by ID
func (to *Session) UpdateStaticDNSEntryByID(id int, sdns tc.StaticDNSEntry) (tc.Alerts, ReqInf, int, error) {
	// fill in missing IDs from names
	staticDNSEntryIDs(to, &sdns)
	var remoteAddr net.Addr
	reqBody, err := json.Marshal(sdns)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, 0, err
	}
	route := fmt.Sprintf("%s?id=%d", API_v13_StaticDNSEntries, id)
	resp, remoteAddr, errClient := to.RawRequest(http.MethodPut, route, reqBody)
	if resp != nil {
		defer resp.Body.Close()
		var alerts tc.Alerts
		if err = json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
			return alerts, reqInf, resp.StatusCode, err
		}
		return alerts, reqInf, resp.StatusCode, errClient
	}
	return tc.Alerts{}, reqInf, 0, errClient
}

// Returns a list of StaticDNSEntrys
func (to *Session) GetStaticDNSEntries() ([]tc.StaticDNSEntry, ReqInf, error) {
	resp, remoteAddr, err := to.request(http.MethodGet, API_v13_StaticDNSEntries, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.StaticDNSEntriesResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Response, reqInf, nil
}

// GET a StaticDNSEntry by the StaticDNSEntry ID
func (to *Session) GetStaticDNSEntryByID(id int) ([]tc.StaticDNSEntry, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_v13_StaticDNSEntries, id)
	resp, remoteAddr, err := to.request(http.MethodGet, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.StaticDNSEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// GET a StaticDNSEntry by the StaticDNSEntry hsot
func (to *Session) GetStaticDNSEntriesByHost(host string) ([]tc.StaticDNSEntry, ReqInf, error) {
	url := fmt.Sprintf("%s?host=%s", API_v13_StaticDNSEntries, host)
	resp, remoteAddr, err := to.request(http.MethodGet, url, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return nil, reqInf, err
	}
	defer resp.Body.Close()

	var data tc.StaticDNSEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, reqInf, err
	}

	return data.Response, reqInf, nil
}

// DELETE a StaticDNSEntry by ID
func (to *Session) DeleteStaticDNSEntryByID(id int) (tc.Alerts, ReqInf, error) {
	route := fmt.Sprintf("%s?id=%d", API_v13_StaticDNSEntries, id)
	resp, remoteAddr, err := to.request(http.MethodDelete, route, nil)
	reqInf := ReqInf{CacheHitStatus: CacheHitStatusMiss, RemoteAddr: remoteAddr}
	if err != nil {
		return tc.Alerts{}, reqInf, err
	}
	defer resp.Body.Close()
	var alerts tc.Alerts
	err = json.NewDecoder(resp.Body).Decode(&alerts)
	return alerts, reqInf, nil
}
