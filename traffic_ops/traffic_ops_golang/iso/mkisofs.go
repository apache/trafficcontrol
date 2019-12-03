package iso

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// bondedRegex matches a bonded device interface name.
var bondedRegex = regexp.MustCompile(`^bond\d+`)

// writeNetworkCfg writes the network.cfg config to w.
func writeNetworkCfg(w io.Writer, r isoRequest, nameservers []string) error {
	var cfg configWriter

	cfg.addIP("IPADDR", r.IPAddr)
	cfg.addIP("NETMASK", r.IPNetmask)
	cfg.addIP("GATEWAY", r.IPGateway)
	isBonded := bondedRegex.MatchString(r.InterfaceName)
	if isBonded {
		cfg.addOpt("BOND_DEVICE", r.InterfaceName)
	} else {
		cfg.addOpt("DEVICE", r.InterfaceName)
	}
	cfg.addOpt("MTU", strconv.Itoa(r.InterfaceMTU))
	cfg.addOpt("NAMESERVER", strings.Join(nameservers, ","))
	cfg.addOpt("HOSTNAME", func() string {
		fqdn := r.HostName
		if r.DomainName != "" {
			fqdn += "." + r.DomainName
		}
		return fqdn
	}())
	cfg.addOpt("NETWORKING_IPV6", "yes")
	cfg.addIP("IPV6ADDR", r.IP6Address)
	cfg.addIP("IPV6_DEFAULTGW", r.IP6Gateway)
	if isBonded {
		cfg.addOpt("BONDING_OPTS", "miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4")
	}
	cfg.addBoolStr("DHCP", r.DHCP)

	_, err := io.Copy(w, &cfg)
	return err
}

// writeMgmtNetworkCfg writes the mgmt_network.cfg config to w.
func writeMgmtNetworkCfg(w io.Writer, r isoRequest) error {
	var cfg configWriter

	// Test if management IP is IPv6
	if r.MgmtIPAddress.To16() != nil && r.MgmtIPAddress.To4() == nil {
		cfg.addIP("IPV6ADDR", r.MgmtIPAddress)
	} else {
		cfg.addIP("IPADDR", r.MgmtIPAddress)
	}
	cfg.addIP("NETMASK", r.MgmtIPNetmask)
	cfg.addIP("GATEWAY", r.MgmtIPGateway)
	cfg.addOpt("DEVICE", r.MgmtInterface)

	_, err := io.Copy(w, &cfg)
	return err
}

// writeDiskCfg writes the disk.cfg config to w.
func writeDiskCfg(w io.Writer, r isoRequest) error {
	var cfg configWriter

	cfg.addOpt("boot_drives", r.Disk)

	_, err := io.Copy(w, &cfg)
	return err
}

// writePasswordCfg writes the password.cfg config to w.
// The salt parameter is optional. If salt is blank, then a
// random 8-character salt will be used.
func writePasswordCfg(w io.Writer, r isoRequest, salt string) error {
	if salt == "" {
		salt = rndSalt(8)
	}

	cryptedPw, err := crypt(r.RootPass, salt)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "rootpw --iscrypted %s\n", cryptedPw)
	return err
}

// configWriter is a helper type to create config files
// of format:
//   OPT="VALUE"
//   OPT="VALUE"
type configWriter struct {
	b    bytes.Buffer
	line int
}

// addOpt adds an option to the config.
func (c *configWriter) addOpt(name, value string) {
	if c.line > 0 {
		fmt.Fprintln(&c.b)
	}
	c.line++
	fmt.Fprintf(&c.b, "%s=%q", name, value)
}

// addIP adds an IP option to the config. It handles the
// case where the given IP is empty/nil.
func (c *configWriter) addIP(name string, ip net.IP) {
	// Avoid using `<nil>`, i.e. net.IP{}.String() = <nil>
	var v string
	if len(ip) > 0 {
		v = ip.String()
	}
	c.addOpt(name, v)
}

// addBoolStr adds a BoolStr option to the config, using
// "yes" / "no" values.
func (c *configWriter) addBoolStr(name string, b boolStr) {
	var v string
	if b.val {
		v = "yes"
	} else {
		v = "no"
	}
	c.addOpt(name, v)
}

// Read satisfies the io.Reader interface, and allows for
// using the configWriter by the io.Copy() function.
func (c *configWriter) Read(p []byte) (int, error) {
	return c.b.Read(p)
}
