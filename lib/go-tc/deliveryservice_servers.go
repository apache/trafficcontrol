package tc

import (
	"time"
)

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

// DeliveryServiceServerResponse ...
type DeliveryServiceServerResponse struct {
	Orderby  string                  `json:"orderby"`
	Response []DeliveryServiceServer `json:"response"`
	Size     int                     `json:"size"`
	Limit    int                     `json:"limit"`
}

type DSSMapResponse struct {
	DsId    int   `json:"dsId"`
	Replace bool  `json:"replace"`
	Servers []int `json:"servers"`
}

type DSSReplaceResponse struct {
	Alerts   []Alert        `json:"alerts"`
	Response DSSMapResponse `json:"response"`
}

type DSServersResponse struct {
	Response DeliveryServiceServers `json:"response"`
}

type DeliveryServiceServers struct {
	ServerNames []string `json:"serverNames"`
	XmlId       string   `json:"xmlId"`
}

// DeliveryServiceServer ...
type DeliveryServiceServer struct {
	Server          *int       `json:"server" db:"server"`
	DeliveryService *int       `json:"deliveryService" db:"deliveryservice"`
	LastUpdated     *TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type Filter int

const (
	Assigned Filter = iota
	Unassigned
	Eligible
)

type DSServersAttrResponse struct {
	Response []DSServer `json:"response"`
}

type DSServer struct {
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
	InterfaceMtu                *int                 `json:"interfaceMtu" db:"interface_mtu"`
	InterfaceName               *string              `json:"interfaceName" db:"interface_name"`
	IP6Address                  *string              `json:"ip6Address" db:"ip6_address"`
	IP6Gateway                  *string              `json:"ip6Gateway" db:"ip6_gateway"`
	IPAddress                   *string              `json:"ipAddress" db:"ip_address"`
	IPGateway                   *string              `json:"ipGateway" db:"ip_gateway"`
	IPNetmask                   *string              `json:"ipNetmask" db:"ip_netmask"`
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
