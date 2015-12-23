package todb

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type Asn struct {
	Id          int64     `db:"id" json:"id"`
	Asn         int64     `db:"asn" json:"asn"`
	Cachegroup  int64     `db:"cachegroup" json:"cachegroup"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Cachegroup struct {
	Id                 int64      `db:"id" json:"id"`
	Name               string     `db:"name" json:"name"`
	ShortName          string     `db:"short_name" json:"shortName"`
	Latitude           null.Float `db:"latitude" json:"latitude"`
	Longitude          null.Float `db:"longitude" json:"longitude"`
	ParentCachegroupId null.Int   `db:"parent_cachegroup_id" json:"parentCachegroupId"`
	Type               int64      `db:"type" json:"type"`
	LastUpdated        time.Time  `db:"last_updated" json:"lastUpdated"`
}

type CachegroupParameter struct {
	Cachegroup  int64     `db:"cachegroup" json:"cachegroup"`
	Parameter   int64     `db:"parameter" json:"parameter"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Cdn struct {
	Id            int64       `db:"id" json:"id"`
	Name          null.String `db:"name" json:"name"`
	LastUpdated   time.Time   `db:"last_updated" json:"lastUpdated"`
	DnssecEnabled null.Int    `db:"dnssec_enabled" json:"dnssecEnabled"`
}

type Deliveryservice struct {
	Id                   int64       `db:"id" json:"id"`
	XmlId                string      `db:"xml_id" json:"xmlId"`
	Active               int64       `db:"active" json:"active"`
	Dscp                 int64       `db:"dscp" json:"dscp"`
	Signed               null.Int    `db:"signed" json:"signed"`
	QstringIgnore        null.Int    `db:"qstring_ignore" json:"qstringIgnore"`
	GeoLimit             null.Int    `db:"geo_limit" json:"geoLimit"`
	HttpBypassFqdn       null.String `db:"http_bypass_fqdn" json:"httpBypassFqdn"`
	DnsBypassIp          null.String `db:"dns_bypass_ip" json:"dnsBypassIp"`
	DnsBypassIp6         null.String `db:"dns_bypass_ip6" json:"dnsBypassIp6"`
	DnsBypassTtl         null.Int    `db:"dns_bypass_ttl" json:"dnsBypassTtl"`
	OrgServerFqdn        null.String `db:"org_server_fqdn" json:"orgServerFqdn"`
	Type                 int64       `db:"type" json:"type"`
	Profile              int64       `db:"profile" json:"profile"`
	CdnId                null.Int    `db:"cdn_id" json:"cdnId"`
	CcrDnsTtl            null.Int    `db:"ccr_dns_ttl" json:"ccrDnsTtl"`
	GlobalMaxMbps        null.Int    `db:"global_max_mbps" json:"globalMaxMbps"`
	GlobalMaxTps         null.Int    `db:"global_max_tps" json:"globalMaxTps"`
	LongDesc             null.String `db:"long_desc" json:"longDesc"`
	LongDesc1            null.String `db:"long_desc_1" json:"longDesc1"`
	LongDesc2            null.String `db:"long_desc_2" json:"longDesc2"`
	MaxDnsAnswers        null.Int    `db:"max_dns_answers" json:"maxDnsAnswers"`
	InfoUrl              null.String `db:"info_url" json:"infoUrl"`
	MissLat              null.Float  `db:"miss_lat" json:"missLat"`
	MissLong             null.Float  `db:"miss_long" json:"missLong"`
	CheckPath            null.String `db:"check_path" json:"checkPath"`
	LastUpdated          time.Time   `db:"last_updated" json:"lastUpdated"`
	Protocol             null.Int    `db:"protocol" json:"protocol"`
	SslKeyVersion        null.Int    `db:"ssl_key_version" json:"sslKeyVersion"`
	Ipv6RoutingEnabled   null.Int    `db:"ipv6_routing_enabled" json:"ipv6RoutingEnabled"`
	RangeRequestHandling null.Int    `db:"range_request_handling" json:"rangeRequestHandling"`
	EdgeHeaderRewrite    null.String `db:"edge_header_rewrite" json:"edgeHeaderRewrite"`
	OriginShield         null.String `db:"origin_shield" json:"originShield"`
	MidHeaderRewrite     null.String `db:"mid_header_rewrite" json:"midHeaderRewrite"`
	RegexRemap           null.String `db:"regex_remap" json:"regexRemap"`
	Cacheurl             null.String `db:"cacheurl" json:"cacheurl"`
	RemapText            null.String `db:"remap_text" json:"remapText"`
	MultiSiteOrigin      null.Int    `db:"multi_site_origin" json:"multiSiteOrigin"`
	DisplayName          string      `db:"display_name" json:"displayName"`
	TrResponseHeaders    null.String `db:"tr_response_headers" json:"trResponseHeaders"`
	InitialDispersion    null.Int    `db:"initial_dispersion" json:"initialDispersion"`
	DnsBypassCname       null.String `db:"dns_bypass_cname" json:"dnsBypassCname"`
	TrRequestHeaders     null.String `db:"tr_request_headers" json:"trRequestHeaders"`
}

type DeliveryserviceRegex struct {
	Deliveryservice int64    `db:"deliveryservice" json:"deliveryservice"`
	Regex           int64    `db:"regex" json:"regex"`
	SetNumber       null.Int `db:"set_number" json:"setNumber"`
}

type DeliveryserviceServer struct {
	Deliveryservice int64     `db:"deliveryservice" json:"deliveryservice"`
	Server          int64     `db:"server" json:"server"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
}

type DeliveryserviceTmuser struct {
	Deliveryservice int64     `db:"deliveryservice" json:"deliveryservice"`
	TmUserId        int64     `db:"tm_user_id" json:"tmUserId"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
}

type Division struct {
	Id          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Federation struct {
	Id          int64       `db:"id" json:"id"`
	Cname       string      `db:"cname" json:"cname"`
	Description null.String `db:"description" json:"description"`
	Ttl         int64       `db:"ttl" json:"ttl"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type FederationDeliveryservice struct {
	Federation      int64     `db:"federation" json:"federation"`
	Deliveryservice int64     `db:"deliveryservice" json:"deliveryservice"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
}

type FederationFederationResolver struct {
	Federation         int64     `db:"federation" json:"federation"`
	FederationResolver int64     `db:"federation_resolver" json:"federationResolver"`
	LastUpdated        time.Time `db:"last_updated" json:"lastUpdated"`
}

type FederationResolver struct {
	Id          int64     `db:"id" json:"id"`
	IpAddress   string    `db:"ip_address" json:"ipAddress"`
	Type        int64     `db:"type" json:"type"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type FederationTmuser struct {
	Federation  int64     `db:"federation" json:"federation"`
	TmUser      int64     `db:"tm_user" json:"tmUser"`
	Role        int64     `db:"role" json:"role"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type GooseDbVersion struct {
	Id        int64     `db:"id" json:"id"`
	VersionId int64     `db:"version_id" json:"versionId"`
	IsApplied int64     `db:"is_applied" json:"isApplied"`
	Tstamp    time.Time `db:"tstamp" json:"tstamp"`
}

type Hwinfo struct {
	Id          int64     `db:"id" json:"id"`
	Serverid    int64     `db:"serverid" json:"serverid"`
	Description string    `db:"description" json:"description"`
	Val         string    `db:"val" json:"val"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Job struct {
	Id                 int64       `db:"id" json:"id"`
	Agent              null.Int    `db:"agent" json:"agent"`
	ObjectType         null.String `db:"object_type" json:"objectType"`
	ObjectName         null.String `db:"object_name" json:"objectName"`
	Keyword            string      `db:"keyword" json:"keyword"`
	Parameters         null.String `db:"parameters" json:"parameters"`
	AssetUrl           string      `db:"asset_url" json:"assetUrl"`
	AssetType          string      `db:"asset_type" json:"assetType"`
	Status             int64       `db:"status" json:"status"`
	StartTime          time.Time   `db:"start_time" json:"startTime"`
	EnteredTime        time.Time   `db:"entered_time" json:"enteredTime"`
	JobUser            int64       `db:"job_user" json:"jobUser"`
	LastUpdated        time.Time   `db:"last_updated" json:"lastUpdated"`
	JobDeliveryservice null.Int    `db:"job_deliveryservice" json:"jobDeliveryservice"`
}

type JobAgent struct {
	Id          int64       `db:"id" json:"id"`
	Name        null.String `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	Active      int64       `db:"active" json:"active"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type JobResult struct {
	Id          int64       `db:"id" json:"id"`
	Job         int64       `db:"job" json:"job"`
	Agent       int64       `db:"agent" json:"agent"`
	Result      string      `db:"result" json:"result"`
	Description null.String `db:"description" json:"description"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type JobStatus struct {
	Id          int64       `db:"id" json:"id"`
	Name        null.String `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type Log struct {
	Id          int64       `db:"id" json:"id"`
	Level       null.String `db:"level" json:"level"`
	Message     string      `db:"message" json:"message"`
	TmUser      int64       `db:"tm_user" json:"tmUser"`
	Ticketnum   null.String `db:"ticketnum" json:"ticketnum"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type Parameter struct {
	Id          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	ConfigFile  string    `db:"config_file" json:"configFile"`
	Value       string    `db:"value" json:"value"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type PhysLocation struct {
	Id          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	ShortName   string      `db:"short_name" json:"shortName"`
	Address     string      `db:"address" json:"address"`
	City        string      `db:"city" json:"city"`
	State       string      `db:"state" json:"state"`
	Zip         string      `db:"zip" json:"zip"`
	Poc         null.String `db:"poc" json:"poc"`
	Phone       null.String `db:"phone" json:"phone"`
	Email       null.String `db:"email" json:"email"`
	Comments    null.String `db:"comments" json:"comments"`
	Region      int64       `db:"region" json:"region"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type Profile struct {
	Id          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type ProfileParameter struct {
	Profile     int64     `db:"profile" json:"profile"`
	Parameter   int64     `db:"parameter" json:"parameter"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Regex struct {
	Id          int64     `db:"id" json:"id"`
	Pattern     string    `db:"pattern" json:"pattern"`
	Type        int64     `db:"type" json:"type"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Region struct {
	Id          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Division    int64     `db:"division" json:"division"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Role struct {
	Id          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	PrivLevel   int64       `db:"priv_level" json:"privLevel"`
}

type Server struct {
	Id             int64       `db:"id" json:"id"`
	HostName       string      `db:"host_name" json:"hostName"`
	DomainName     string      `db:"domain_name" json:"domainName"`
	TcpPort        null.Int    `db:"tcp_port" json:"tcpPort"`
	XmppId         null.String `db:"xmpp_id" json:"xmppId"`
	XmppPasswd     null.String `db:"xmpp_passwd" json:"xmppPasswd"`
	InterfaceName  string      `db:"interface_name" json:"interfaceName"`
	IpAddress      string      `db:"ip_address" json:"ipAddress"`
	IpNetmask      string      `db:"ip_netmask" json:"ipNetmask"`
	IpGateway      string      `db:"ip_gateway" json:"ipGateway"`
	Ip6Address     null.String `db:"ip6_address" json:"ip6Address"`
	Ip6Gateway     null.String `db:"ip6_gateway" json:"ip6Gateway"`
	InterfaceMtu   int64       `db:"interface_mtu" json:"interfaceMtu"`
	PhysLocation   int64       `db:"phys_location" json:"physLocation"`
	Rack           null.String `db:"rack" json:"rack"`
	Cachegroup     int64       `db:"cachegroup" json:"cachegroup"`
	Type           int64       `db:"type" json:"type"`
	Status         int64       `db:"status" json:"status"`
	UpdPending     int64       `db:"upd_pending" json:"updPending"`
	Profile        int64       `db:"profile" json:"profile"`
	CdnId          null.Int    `db:"cdn_id" json:"cdnId"`
	MgmtIpAddress  null.String `db:"mgmt_ip_address" json:"mgmtIpAddress"`
	MgmtIpNetmask  null.String `db:"mgmt_ip_netmask" json:"mgmtIpNetmask"`
	MgmtIpGateway  null.String `db:"mgmt_ip_gateway" json:"mgmtIpGateway"`
	IloIpAddress   null.String `db:"ilo_ip_address" json:"iloIpAddress"`
	IloIpNetmask   null.String `db:"ilo_ip_netmask" json:"iloIpNetmask"`
	IloIpGateway   null.String `db:"ilo_ip_gateway" json:"iloIpGateway"`
	IloUsername    null.String `db:"ilo_username" json:"iloUsername"`
	IloPassword    null.String `db:"ilo_password" json:"iloPassword"`
	RouterHostName null.String `db:"router_host_name" json:"routerHostName"`
	RouterPortName null.String `db:"router_port_name" json:"routerPortName"`
	LastUpdated    time.Time   `db:"last_updated" json:"lastUpdated"`
}

type Servercheck struct {
	Id          int64     `db:"id" json:"id"`
	Server      int64     `db:"server" json:"server"`
	Aa          null.Int  `db:"aa" json:"aa"`
	Ab          null.Int  `db:"ab" json:"ab"`
	Ac          null.Int  `db:"ac" json:"ac"`
	Ad          null.Int  `db:"ad" json:"ad"`
	Ae          null.Int  `db:"ae" json:"ae"`
	Af          null.Int  `db:"af" json:"af"`
	Ag          null.Int  `db:"ag" json:"ag"`
	Ah          null.Int  `db:"ah" json:"ah"`
	Ai          null.Int  `db:"ai" json:"ai"`
	Aj          null.Int  `db:"aj" json:"aj"`
	Ak          null.Int  `db:"ak" json:"ak"`
	Al          null.Int  `db:"al" json:"al"`
	Am          null.Int  `db:"am" json:"am"`
	An          null.Int  `db:"an" json:"an"`
	Ao          null.Int  `db:"ao" json:"ao"`
	Ap          null.Int  `db:"ap" json:"ap"`
	Aq          null.Int  `db:"aq" json:"aq"`
	Ar          null.Int  `db:"ar" json:"ar"`
	As          null.Int  `db:"as" json:"as"`
	At          null.Int  `db:"at" json:"at"`
	Au          null.Int  `db:"au" json:"au"`
	Av          null.Int  `db:"av" json:"av"`
	Aw          null.Int  `db:"aw" json:"aw"`
	Ax          null.Int  `db:"ax" json:"ax"`
	Ay          null.Int  `db:"ay" json:"ay"`
	Az          null.Int  `db:"az" json:"az"`
	Ba          null.Int  `db:"ba" json:"ba"`
	Bb          null.Int  `db:"bb" json:"bb"`
	Bc          null.Int  `db:"bc" json:"bc"`
	Bd          null.Int  `db:"bd" json:"bd"`
	Be          null.Int  `db:"be" json:"be"`
	LastUpdated time.Time `db:"last_updated" json:"lastUpdated"`
}

type Staticdnsentry struct {
	Id              int64     `db:"id" json:"id"`
	Host            string    `db:"host" json:"host"`
	Address         string    `db:"address" json:"address"`
	Type            int64     `db:"type" json:"type"`
	Ttl             int64     `db:"ttl" json:"ttl"`
	Deliveryservice int64     `db:"deliveryservice" json:"deliveryservice"`
	Cachegroup      null.Int  `db:"cachegroup" json:"cachegroup"`
	LastUpdated     time.Time `db:"last_updated" json:"lastUpdated"`
}

type StatsSummary struct {
	Id                  int64     `db:"id" json:"id"`
	CdnName             string    `db:"cdn_name" json:"cdnName"`
	DeliveryserviceName string    `db:"deliveryservice_name" json:"deliveryserviceName"`
	StatName            string    `db:"stat_name" json:"statName"`
	StatValue           float64   `db:"stat_value" json:"statValue"`
	SummaryTime         time.Time `db:"summary_time" json:"summaryTime"`
	StatDate            time.Time `db:"stat_date" json:"statDate"`
}

type Status struct {
	Id          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}

type TmUser struct {
	Id                 int64       `db:"id" json:"id"`
	Username           null.String `db:"username" json:"username"`
	Role               null.Int    `db:"role" json:"role"`
	Uid                null.Int    `db:"uid" json:"uid"`
	Gid                null.Int    `db:"gid" json:"gid"`
	LocalPasswd        null.String `db:"local_passwd" json:"localPasswd"`
	ConfirmLocalPasswd null.String `db:"confirm_local_passwd" json:"confirmLocalPasswd"`
	LastUpdated        time.Time   `db:"last_updated" json:"lastUpdated"`
	Company            null.String `db:"company" json:"company"`
	Email              null.String `db:"email" json:"email"`
	FullName           null.String `db:"full_name" json:"fullName"`
	NewUser            int64       `db:"new_user" json:"newUser"`
	AddressLine1       null.String `db:"address_line1" json:"addressLine1"`
	AddressLine2       null.String `db:"address_line2" json:"addressLine2"`
	City               null.String `db:"city" json:"city"`
	StateOrProvince    null.String `db:"state_or_province" json:"stateOrProvince"`
	PhoneNumber        null.String `db:"phone_number" json:"phoneNumber"`
	PostalCode         null.String `db:"postal_code" json:"postalCode"`
	Country            null.String `db:"country" json:"country"`
	LocalUser          int64       `db:"local_user" json:"localUser"`
	Token              null.String `db:"token" json:"token"`
	RegistrationSent   time.Time   `db:"registration_sent" json:"registrationSent"`
}

type ToExtension struct {
	Id                    int64       `db:"id" json:"id"`
	Name                  string      `db:"name" json:"name"`
	Version               string      `db:"version" json:"version"`
	InfoUrl               string      `db:"info_url" json:"infoUrl"`
	ScriptFile            string      `db:"script_file" json:"scriptFile"`
	Isactive              int64       `db:"isactive" json:"isactive"`
	AdditionalConfigJson  null.String `db:"additional_config_json" json:"additionalConfigJson"`
	Description           null.String `db:"description" json:"description"`
	ServercheckShortName  null.String `db:"servercheck_short_name" json:"servercheckShortName"`
	ServercheckColumnName null.String `db:"servercheck_column_name" json:"servercheckColumnName"`
	Type                  int64       `db:"type" json:"type"`
	LastUpdated           time.Time   `db:"last_updated" json:"lastUpdated"`
}

type Type struct {
	Id          int64       `db:"id" json:"id"`
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	UseInTable  null.String `db:"use_in_table" json:"useInTable"`
	LastUpdated time.Time   `db:"last_updated" json:"lastUpdated"`
}
