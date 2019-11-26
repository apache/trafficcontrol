package iso

import (
	"html/template"
	"io"
	"net"
	"regexp"
	"strings"
)

var (
	bondedRegex = regexp.MustCompile(`^bond\d+`)

	tmplFuncs = template.FuncMap{
		"printIP": func(v net.IP) string {
			if len(v) == 0 {
				return ""
			}
			return v.String()
		},
		"printBoolStr": func(v boolStr) string {
			if v.val {
				return "yes"
			}
			return "no"
		},
	}
)

var networkCfgTmpl = template.Must(template.New("network.cfg").Funcs(tmplFuncs).Parse(
	`IPADDR="{{ .ISORequest.IPAddr | printIP }}"
NETMASK="{{ .ISORequest.IPNetmask | printIP }}"
GATEWAY="{{ .ISORequest.IPGateway | printIP }}"
{{ if .IsBonded }}BOND_DEVICE{{ else }}DEVICE{{ end }}="{{ .ISORequest.InterfaceName }}"
MTU="{{ .ISORequest.InterfaceMTU }}"
NAMESERVER="{{ .Nameservers }}"
HOSTNAME="{{ .FQDN }}"
NETWORKING_IPV6="yes"
IPV6ADDR="{{ .ISORequest.IP6Address | printIP }}"
IPV6_DEFAULTGW="{{ .ISORequest.IP6Gateway | printIP }}"
{{- if .IsBonded }}
BONDING_OPTS="miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4"
{{- end }}
DHCP="{{ .ISORequest.DHCP | printBoolStr }}"`,
))

func writeNetworkCfg(w io.Writer, r isoRequest, nameservers []string) error {
	fqdn := r.HostName
	if r.DomainName != "" {
		fqdn += "." + r.DomainName
	}

	data := struct {
		ISORequest  isoRequest
		Nameservers string
		FQDN        string
		IsBonded    bool
	}{
		ISORequest:  r,
		Nameservers: strings.Join(nameservers, ","),
		FQDN:        fqdn,
		IsBonded:    bondedRegex.MatchString(r.InterfaceName),
	}

	return networkCfgTmpl.Execute(w, data)
}

var mgmtNetworkCfgTmpl = template.Must(template.New("mgmt_network.cfg").Funcs(tmplFuncs).Parse(
	`{{ if .IsIPv6 }}IPV6ADDR{{ else }}IPADDR{{ end }}="{{ .ISORequest.MgmtIPAddress | printIP }}"
NETMASK="{{ .ISORequest.MgmtIPNetmask | printIP }}"
GATEWAY="{{ .ISORequest.MgmtIPGateway | printIP }}"
DEVICE="{{ .ISORequest.MgmtInterface }}"`,
))

func writeMgmtNetworkCfg(w io.Writer, r isoRequest) error {
	data := struct {
		ISORequest isoRequest
		IsIPv6     bool
	}{
		ISORequest: r,
		IsIPv6:     r.MgmtIPAddress.To16() != nil && r.MgmtIPAddress.To4() == nil,
	}

	return mgmtNetworkCfgTmpl.Execute(w, data)
}
