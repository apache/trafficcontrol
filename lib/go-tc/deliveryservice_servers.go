package tc

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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// DeliveryServiceServerV5 is the struct used to represent a delivery service server, in the latest minor version for
// API 5.x.
type DeliveryServiceServerV5 = DeliveryServiceServerV50

// DeliveryServiceServerV50 is the type of each entry in the `response` array
// property of responses from Traffic Ops to GET requests made to the
// /deliveryserviceservers endpoint of its API, for api version 5.0.
type DeliveryServiceServerV50 struct {
	Server          *int       `json:"server" db:"server"`
	DeliveryService *int       `json:"deliveryService" db:"deliveryservice"`
	LastUpdated     *time.Time `json:"lastUpdated" db:"last_updated"`
}

// DeliveryServiceServerResponseV5 is the type of a response from Traffic Ops
// to a GET request to the /deliveryserviceserver endpoint, for the latest minor version of api v5.x.
type DeliveryServiceServerResponseV5 = DeliveryServiceServerResponseV50

// DeliveryServiceServerResponseV50 is the type of a response from Traffic Ops
// to a GET request to the /deliveryserviceserver endpoint for API v5.0.
type DeliveryServiceServerResponseV50 struct {
	Orderby  string                    `json:"orderby"`
	Response []DeliveryServiceServerV5 `json:"response"`
	Size     int                       `json:"size"`
	Limit    int                       `json:"limit"`
	Alerts
}

// DeliveryServiceServerResponse is the type of a response from Traffic Ops
// to a GET request to the /deliveryserviceserver endpoint.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DeliveryServiceServerResponse struct {
	Orderby  string                  `json:"orderby"`
	Response []DeliveryServiceServer `json:"response"`
	Size     int                     `json:"size"`
	Limit    int                     `json:"limit"`
	Alerts
}

// DSSMapResponse is the type of the `response` property of a response from
// Traffic Ops to a PUT request made to the /deliveryserviceserver endpoint of
// its API.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSSMapResponse struct {
	DsId    int   `json:"dsId"`
	Replace bool  `json:"replace"`
	Servers []int `json:"servers"`
}

// DSSReplaceResponse is the type of a response from Traffic Ops to a PUT
// request made to the /deliveryserviceserver endpoint of its API.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSSReplaceResponse struct {
	Alerts
	Response DSSMapResponse `json:"response"`
}

// DSServersResponse is the type of a response from Traffic Ops to a POST
// request made to the /deliveryservices/{{XML ID}}/servers endpoint of its
// API.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServersResponse struct {
	Response DeliveryServiceServers `json:"response"`
	Alerts
}

// DeliveryServiceServers structures represent the servers assigned to a
// Delivery Service.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DeliveryServiceServers struct {
	ServerNames []string `json:"serverNames"`
	XmlId       string   `json:"xmlId"`
}

// DeliveryServiceServer is the type of each entry in the `response` array
// property of responses from Traffic Ops to GET requests made to the
// /deliveryserviceservers endpoint of its API.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DeliveryServiceServer struct {
	Server          *int       `json:"server" db:"server"`
	DeliveryService *int       `json:"deliveryService" db:"deliveryservice"`
	LastUpdated     *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

// Filter is a type of unknown purpose - though it appears related to the
// removed APIv1 endpoints regarding the servers which are assigned to a
// Delivery Service, not assigned to a Delivery Service, and eligible for
// assignment to a Delivery Service.
//
// Deprecated: This is not used by any known ATC code, and is likely a remnant
// of legacy behavior that will be removed in the future.
type Filter int

// Enumerated values of the Filter type.
//
// Deprecated: These are not used by any known ATC code, and are likely
// remnants of legacy behavior that will be removed in the future.
const (
	Assigned Filter = iota
	Unassigned
	Eligible
)

// DSServerBase contains the base information for a Delivery Service Server.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerBase struct {
	Cachegroup                  *string              `json:"cachegroup" db:"cachegroup"`
	CachegroupID                *int                 `json:"cachegroupId" db:"cachegroup_id"`
	CDNID                       *int                 `json:"cdnId" db:"cdn_id"`
	CDNName                     *string              `json:"cdnName" db:"cdn_name"`
	DeliveryServices            *map[string][]string `json:"deliveryServices,omitempty"`
	DomainName                  *string              `json:"domainName" db:"domain_name"`
	FQDN                        *string              `json:"fqdn,omitempty"`
	FqdnTime                    time.Time            `json:"-"`
	GUID                        *string              `json:"guid" db:"guid"`
	HostName                    *string              `json:"hostName" db:"host_name"`
	HTTPSPort                   *int                 `json:"httpsPort" db:"https_port"`
	ID                          *int                 `json:"id" db:"id"`
	ILOIPAddress                *string              `json:"iloIpAddress" db:"ilo_ip_address"`
	ILOIPGateway                *string              `json:"iloIpGateway" db:"ilo_ip_gateway"`
	ILOIPNetmask                *string              `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	ILOPassword                 *string              `json:"iloPassword" db:"ilo_password"`
	ILOUsername                 *string              `json:"iloUsername" db:"ilo_username"`
	LastUpdated                 *TimeNoMod           `json:"lastUpdated" db:"last_updated"`
	MgmtIPAddress               *string              `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIPGateway               *string              `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIPNetmask               *string              `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason               *string              `json:"offlineReason" db:"offline_reason"`
	PhysLocation                *string              `json:"physLocation" db:"phys_location"`
	PhysLocationID              *int                 `json:"physLocationId" db:"phys_location_id"`
	Profile                     *string              `json:"profile" db:"profile"`
	ProfileDesc                 *string              `json:"profileDesc" db:"profile_desc"`
	ProfileID                   *int                 `json:"profileId" db:"profile_id"`
	Rack                        *string              `json:"rack" db:"rack"`
	RouterHostName              *string              `json:"routerHostName" db:"router_host_name"`
	RouterPortName              *string              `json:"routerPortName" db:"router_port_name"`
	Status                      *string              `json:"status" db:"status"`
	StatusID                    *int                 `json:"statusId" db:"status_id"`
	TCPPort                     *int                 `json:"tcpPort" db:"tcp_port"`
	Type                        string               `json:"type" db:"server_type"`
	TypeID                      *int                 `json:"typeId" db:"server_type_id"`
	UpdPending                  *bool                `json:"updPending" db:"upd_pending"`
	ServerCapabilities          []string             `json:"-" db:"server_capabilities"`
	DeliveryServiceCapabilities []string             `json:"-" db:"deliveryservice_capabilities"`
}

// DSServerBaseV4 contains the base information for a Delivery Service Server.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerBaseV4 struct {
	Cachegroup                  *string              `json:"cachegroup" db:"cachegroup"`
	CachegroupID                *int                 `json:"cachegroupId" db:"cachegroup_id"`
	CDNID                       *int                 `json:"cdnId" db:"cdn_id"`
	CDNName                     *string              `json:"cdnName" db:"cdn_name"`
	DeliveryServices            *map[string][]string `json:"deliveryServices,omitempty"`
	DomainName                  *string              `json:"domainName" db:"domain_name"`
	FQDN                        *string              `json:"fqdn,omitempty"`
	FqdnTime                    time.Time            `json:"-"`
	GUID                        *string              `json:"guid" db:"guid"`
	HostName                    *string              `json:"hostName" db:"host_name"`
	HTTPSPort                   *int                 `json:"httpsPort" db:"https_port"`
	ID                          *int                 `json:"id" db:"id"`
	ILOIPAddress                *string              `json:"iloIpAddress" db:"ilo_ip_address"`
	ILOIPGateway                *string              `json:"iloIpGateway" db:"ilo_ip_gateway"`
	ILOIPNetmask                *string              `json:"iloIpNetmask" db:"ilo_ip_netmask"`
	ILOPassword                 *string              `json:"iloPassword" db:"ilo_password"`
	ILOUsername                 *string              `json:"iloUsername" db:"ilo_username"`
	LastUpdated                 *TimeNoMod           `json:"lastUpdated" db:"last_updated"`
	MgmtIPAddress               *string              `json:"mgmtIpAddress" db:"mgmt_ip_address"`
	MgmtIPGateway               *string              `json:"mgmtIpGateway" db:"mgmt_ip_gateway"`
	MgmtIPNetmask               *string              `json:"mgmtIpNetmask" db:"mgmt_ip_netmask"`
	OfflineReason               *string              `json:"offlineReason" db:"offline_reason"`
	PhysLocation                *string              `json:"physLocation" db:"phys_location"`
	PhysLocationID              *int                 `json:"physLocationId" db:"phys_location_id"`
	ProfileNames                []string             `json:"profileNames" db:"profile_name"`
	Rack                        *string              `json:"rack" db:"rack"`
	Status                      *string              `json:"status" db:"status"`
	StatusID                    *int                 `json:"statusId" db:"status_id"`
	TCPPort                     *int                 `json:"tcpPort" db:"tcp_port"`
	Type                        string               `json:"type" db:"server_type"`
	TypeID                      *int                 `json:"typeId" db:"server_type_id"`
	UpdPending                  *bool                `json:"updPending" db:"upd_pending"`
	ServerCapabilities          []string             `json:"-" db:"server_capabilities"`
	DeliveryServiceCapabilities []string             `json:"-" db:"deliveryservice_capabilities"`
}

// DSServerV11 contains the legacy format for a Delivery Service Server.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerV11 struct {
	DSServerBase
	LegacyInterfaceDetails
}

// DSServer contains information for a Delivery Service Server.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServer struct {
	DSServerBase
	ServerInterfaces *[]ServerInterfaceInfo `json:"interfaces" db:"interfaces"`
}

// DSServerResponseV30 is the type of a response from Traffic Ops to a request
// for servers assigned to a Delivery Service - in API version 3.0.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerResponseV30 struct {
	Response []DSServer `json:"response"`
	Alerts
}

// DSServerV4 contains information for a V4.x Delivery Service Server.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerV4 struct {
	DSServerBaseV4
	ServerInterfaces *[]ServerInterfaceInfoV40 `json:"interfaces" db:"interfaces"`
}

// DSServerResponseV40 is the type of a response from Traffic Ops to a request
// for servers assigned to a Delivery Service - in API version 4.0.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerResponseV40 struct {
	Response []DSServerV4 `json:"response"`
	Alerts
}

// DSServerResponseV4 is the type of a response from Traffic Ops to a request
// for servers assigned to a Delivery Service - in the latest minor version of
// API version 4.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerResponseV4 = DSServerResponseV40

// DSServerV5 is an alias for the latest minor version of the major version 5.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerV5 = DSServerV50

// DSServerV50 contains information about a server associated with some Delivery
// Service Server.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerV50 struct {
	ServerV5
	DeliveryServices            map[string][]string `json:"deliveryServices,omitempty"`
	ServerCapabilities          []string            `json:"-" db:"server_capabilities"`
	DeliveryServiceCapabilities []string            `json:"-" db:"deliveryservice_capabilities"`
}

// DSServerResponseV5 is an alias for the latest minor version of the major version 5.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerResponseV5 = DSServerResponseV50

// DSServerResponseV50 is response from Traffic Ops to a request for servers assigned to a Delivery Service - in  the latest minor version APIv50.
//
// Deprecated: Direct server-to-DS assignments are deprecated in favor of
// Topology usage.
type DSServerResponseV50 struct {
	Response []DSServerV50 `json:"response"`
	Alerts
}

// ToDSServerBaseV4 upgrades the DSServerBase to the structure used by the
// latest minor version of version 4 of Traffic Ops's API.
func (oldBase DSServerBase) ToDSServerBaseV4() DSServerBaseV4 {
	var dsServerBaseV4 DSServerBaseV4
	dsServerBaseV4.Cachegroup = oldBase.Cachegroup
	dsServerBaseV4.CachegroupID = oldBase.CachegroupID
	dsServerBaseV4.CDNID = oldBase.CDNID
	dsServerBaseV4.CDNName = oldBase.CDNName
	dsServerBaseV4.DeliveryServices = oldBase.DeliveryServices
	dsServerBaseV4.DomainName = oldBase.DomainName
	dsServerBaseV4.FQDN = oldBase.FQDN
	dsServerBaseV4.FqdnTime = oldBase.FqdnTime
	dsServerBaseV4.GUID = oldBase.GUID
	dsServerBaseV4.HostName = oldBase.HostName
	dsServerBaseV4.HTTPSPort = oldBase.HTTPSPort
	dsServerBaseV4.ID = oldBase.ID
	dsServerBaseV4.ILOIPAddress = oldBase.ILOIPAddress
	dsServerBaseV4.ILOIPGateway = oldBase.ILOIPGateway
	dsServerBaseV4.ILOIPNetmask = oldBase.ILOIPNetmask
	dsServerBaseV4.ILOPassword = oldBase.ILOPassword
	dsServerBaseV4.ILOUsername = oldBase.ILOUsername
	dsServerBaseV4.LastUpdated = oldBase.LastUpdated
	dsServerBaseV4.MgmtIPAddress = oldBase.MgmtIPAddress
	dsServerBaseV4.MgmtIPGateway = oldBase.MgmtIPGateway
	dsServerBaseV4.MgmtIPNetmask = oldBase.MgmtIPNetmask
	dsServerBaseV4.OfflineReason = oldBase.OfflineReason
	dsServerBaseV4.PhysLocation = oldBase.PhysLocation
	dsServerBaseV4.PhysLocationID = oldBase.PhysLocationID
	dsServerBaseV4.ProfileNames = []string{*oldBase.Profile}
	dsServerBaseV4.Rack = oldBase.Rack
	dsServerBaseV4.Status = oldBase.Status
	dsServerBaseV4.StatusID = oldBase.StatusID
	dsServerBaseV4.TCPPort = oldBase.TCPPort
	dsServerBaseV4.Type = oldBase.Type
	dsServerBaseV4.TypeID = oldBase.TypeID
	dsServerBaseV4.UpdPending = oldBase.UpdPending
	dsServerBaseV4.ServerCapabilities = oldBase.ServerCapabilities
	dsServerBaseV4.DeliveryServiceCapabilities = oldBase.DeliveryServiceCapabilities
	return dsServerBaseV4
}

// ToDSServerBase downgrades the DSServerBaseV4 to the structure used by the
// Traffic Ops API in versions earlier than 4.0.
func (baseV4 DSServerBaseV4) ToDSServerBase(routerHostName, routerPort, pDesc *string, pID *int) DSServerBase {
	var dsServerBase DSServerBase
	dsServerBase.Cachegroup = baseV4.Cachegroup
	dsServerBase.CachegroupID = baseV4.CachegroupID
	dsServerBase.CDNID = baseV4.CDNID
	dsServerBase.CDNName = baseV4.CDNName
	dsServerBase.DeliveryServices = baseV4.DeliveryServices
	dsServerBase.DomainName = baseV4.DomainName
	dsServerBase.FQDN = baseV4.FQDN
	dsServerBase.FqdnTime = baseV4.FqdnTime
	dsServerBase.GUID = baseV4.GUID
	dsServerBase.HostName = baseV4.HostName
	dsServerBase.HTTPSPort = baseV4.HTTPSPort
	dsServerBase.ID = baseV4.ID
	dsServerBase.ILOIPAddress = baseV4.ILOIPAddress
	dsServerBase.ILOIPGateway = baseV4.ILOIPGateway
	dsServerBase.ILOIPNetmask = baseV4.ILOIPNetmask
	dsServerBase.ILOPassword = baseV4.ILOPassword
	dsServerBase.ILOUsername = baseV4.ILOUsername
	dsServerBase.LastUpdated = baseV4.LastUpdated
	dsServerBase.MgmtIPAddress = baseV4.MgmtIPAddress
	dsServerBase.MgmtIPGateway = baseV4.MgmtIPGateway
	dsServerBase.MgmtIPNetmask = baseV4.MgmtIPNetmask
	dsServerBase.OfflineReason = baseV4.OfflineReason
	dsServerBase.PhysLocation = baseV4.PhysLocation
	dsServerBase.PhysLocationID = baseV4.PhysLocationID
	dsServerBase.Profile = &baseV4.ProfileNames[0]
	dsServerBase.ProfileDesc = pDesc
	dsServerBase.ProfileID = pID
	dsServerBase.Rack = baseV4.Rack
	dsServerBase.Status = baseV4.Status
	dsServerBase.StatusID = baseV4.StatusID
	dsServerBase.TCPPort = baseV4.TCPPort
	dsServerBase.Type = baseV4.Type
	dsServerBase.TypeID = baseV4.TypeID
	dsServerBase.UpdPending = baseV4.UpdPending
	dsServerBase.ServerCapabilities = baseV4.ServerCapabilities
	dsServerBase.DeliveryServiceCapabilities = baseV4.DeliveryServiceCapabilities
	dsServerBase.RouterHostName = routerHostName
	dsServerBase.RouterPortName = routerPort
	return dsServerBase
}

// ToDSServerV5 convert DSServerV4 lastUpdated time format to RFC3339 for DSServerV5
// and also assign V4 values to V5
func (server DSServerV4) Upgrade() DSServerV5 {
	var t time.Time
	if server.LastUpdated != nil {
		t = server.LastUpdated.Time
	}

	var dses map[string][]string
	if server.DeliveryServices != nil {
		dses = *server.DeliveryServices
	}

	return DSServerV5{
		DeliveryServices: dses,
		ServerV5: ServerV5{
			CacheGroup:         util.CoalesceToDefault(server.Cachegroup),
			CacheGroupID:       util.CoalesceToDefault(server.CachegroupID),
			CDNID:              util.CoalesceToDefault(server.CDNID),
			CDN:                util.CoalesceToDefault(server.CDNName),
			DomainName:         util.CoalesceToDefault(server.DomainName),
			GUID:               server.GUID,
			HostName:           util.CoalesceToDefault(server.HostName),
			HTTPSPort:          server.HTTPSPort,
			ID:                 util.CoalesceToDefault(server.ID),
			ILOIPAddress:       server.ILOIPAddress,
			ILOIPGateway:       server.ILOIPGateway,
			ILOIPNetmask:       server.ILOIPNetmask,
			ILOPassword:        server.ILOPassword,
			ILOUsername:        server.ILOUsername,
			LastUpdated:        t,
			MgmtIPAddress:      server.MgmtIPAddress,
			MgmtIPGateway:      server.MgmtIPGateway,
			MgmtIPNetmask:      server.MgmtIPNetmask,
			OfflineReason:      server.OfflineReason,
			Profiles:           server.ProfileNames,
			PhysicalLocation:   util.CoalesceToDefault(server.PhysLocation),
			PhysicalLocationID: util.CoalesceToDefault(server.PhysLocationID),
			Rack:               server.Rack,
			Status:             util.CoalesceToDefault(server.Status),
			StatusID:           util.CoalesceToDefault(server.StatusID),
			TCPPort:            server.TCPPort,
			Type:               server.Type,
			TypeID:             util.CoalesceToDefault(server.TypeID),
			Interfaces:         *server.ServerInterfaces,
		},
		ServerCapabilities:          server.ServerCapabilities,
		DeliveryServiceCapabilities: server.DeliveryServiceCapabilities,
	}
}
