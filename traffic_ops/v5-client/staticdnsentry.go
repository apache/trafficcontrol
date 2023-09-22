package client

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

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

// apiStaticDNSEntries is the full path to the /staticdnsentries API
// endpoint.
const apiStaticDNSEntries = "/staticdnsentries"

func staticDNSEntryIDsV5(to *Session, sdns *tc.StaticDNSEntryV5) error {
	if sdns == nil {
		return errors.New("cannot resolve names to IDs for nil StaticDNSEntry")
	}
	if (sdns.CacheGroupID == nil || *sdns.CacheGroupID == 0) && (sdns.CacheGroupName != nil && *sdns.CacheGroupName != "") {
		opts := NewRequestOptions()
		opts.QueryParameters.Set("name", *sdns.CacheGroupName)
		p, _, err := to.GetCacheGroups(opts)
		if err != nil {
			return err
		}
		if len(p.Response) == 0 {
			return errors.New("no CacheGroup named " + *sdns.CacheGroupName)
		}
		if p.Response[0].ID == nil {
			return errors.New("CacheGroup named " + *sdns.CacheGroupName + " has a nil ID")
		}
		sdns.CacheGroupID = p.Response[0].ID
	}

	if (sdns.DeliveryServiceID == nil || *sdns.DeliveryServiceID == 0) && (sdns.DeliveryService != nil && *sdns.DeliveryService != "") {
		opts := NewRequestOptions()
		opts.QueryParameters.Set("xmlId", *sdns.DeliveryService)
		dses, _, err := to.GetDeliveryServices(opts)
		if err != nil {
			return err
		}
		if len(dses.Response) == 0 {
			return errors.New("no deliveryservice with name " + *sdns.DeliveryService)
		}
		if dses.Response[0].ID == nil {
			return errors.New("Deliveryservice with name " + *sdns.DeliveryService + " has a nil ID")
		}
		sdns.DeliveryServiceID = dses.Response[0].ID
	}

	if (sdns.TypeID == nil || *sdns.TypeID == 0) && (sdns.Type != nil && *sdns.Type != "") {
		opts := NewRequestOptions()
		opts.QueryParameters.Set("name", *sdns.Type)
		types, _, err := to.GetTypes(opts)
		if err != nil {
			return err
		}
		if len(types.Response) == 0 {
			return errors.New("no type with name " + *sdns.Type)
		}
		sdns.TypeID = &types.Response[0].ID
	}

	return nil
}

// CreateStaticDNSEntry creates the given Static DNS Entry.
func (to *Session) CreateStaticDNSEntry(sdns tc.StaticDNSEntryV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	// fill in missing IDs from names
	var alerts tc.Alerts
	err := staticDNSEntryIDsV5(to, &sdns)
	if err != nil {
		return alerts, toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}, err
	}
	reqInf, err := to.post(apiStaticDNSEntries, opts, sdns, &alerts)
	return alerts, reqInf, err
}

// UpdateStaticDNSEntry replaces the Static DNS Entry identified by 'id' with
// the one provided.
func (to *Session) UpdateStaticDNSEntry(id int, sdns tc.StaticDNSEntryV5, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	// fill in missing IDs from names
	var alerts tc.Alerts
	err := staticDNSEntryIDsV5(to, &sdns)
	if err != nil {
		return alerts, toclientlib.ReqInf{CacheHitStatus: toclientlib.CacheHitStatusMiss}, err
	}
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	reqInf, err := to.put(apiStaticDNSEntries, opts, sdns, &alerts)
	return alerts, reqInf, err
}

// GetStaticDNSEntries retrieves all Static DNS Entries stored in Traffic Ops.
func (to *Session) GetStaticDNSEntries(opts RequestOptions) (tc.StaticDNSEntriesResponseV5, toclientlib.ReqInf, error) {
	var data tc.StaticDNSEntriesResponseV5
	reqInf, err := to.get(apiStaticDNSEntries, opts, &data)
	return data, reqInf, err
}

// DeleteStaticDNSEntry deletes the Static DNS Entry with the given ID.
func (to *Session) DeleteStaticDNSEntry(id int, opts RequestOptions) (tc.Alerts, toclientlib.ReqInf, error) {
	if opts.QueryParameters == nil {
		opts.QueryParameters = url.Values{}
	}
	opts.QueryParameters.Set("id", strconv.Itoa(id))
	var alerts tc.Alerts
	reqInf, err := to.del(apiStaticDNSEntries, opts, &alerts)
	return alerts, reqInf, err
}
