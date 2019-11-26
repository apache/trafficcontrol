package iso

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
)

// readDefaultUnixResolve reads the /etc/resolv.conf
// file and parses out the nameservers.
func readDefaultUnixResolve() ([]string, error) {
	fd, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return parseResolve(fd)
}

// parseResolve parses r and expects a resolv.conf format.
// It returns all the nameservers found in the file. Any
// formatting or other issues within the file itself are ignored,
// only errors reading from r are returned.
// See following link for more information:
// http://man7.org/linux/man-pages/man5/resolv.conf.5.html
func parseResolve(r io.Reader) ([]string, error) {
	var nameservers []string

	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		if len(l) > 0 && (l[0] == '#' || l[0] == ';') {
			// Ignore comments
			continue
		}
		parts := strings.Fields(l)
		// Look for "nameserver 0.0.0.0" formatted lines
		if len(parts) < 2 || parts[0] != "nameserver" {
			continue
		}
		if net.ParseIP(parts[1]) == nil {
			// Ignore invalid IPs
			continue
		}
		nameservers = append(nameservers, parts[1])
	}

	return nameservers, s.Err()
}
