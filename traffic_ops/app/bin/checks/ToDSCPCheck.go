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

/* ToDSCPCheck.go
   This app scans all REPORTED or ADMIN_DOWN cache nodes for expected DSCP
   marks on each delivery service.
   NOTE: if a particular delivery service DOES NOT have a check path
   configured, then it WILL BE SKIPPED by this tool.
*/

package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/romana/rlog"
)

// Traffic Ops connection params
const AllowInsecureConnections = false
const UserAgent = "go/tc-dscp-monitor"
const UseClientCache = false
const TrafficOpsRequestTimeout = time.Second * time.Duration(10)

var (
	protocol     int
	sslflag      bool
	host_header  string
	confInt      *string
	http4        *http.Transport
	http6        *http.Transport
	https4       *http.Transport
	https6       *http.Transport
	httpClient4  *http.Client
	httpClient6  *http.Client
	httpsClient4 *http.Client
	httpsClient6 *http.Client
	http_port    string
	https_port   string
	eth_layer    layers.Ethernet
	ip4_layer    layers.IPv4
	ip6_layer    layers.IPv6
	tcp_layer    layers.TCP
	tls_layer    layers.TLS
	payload      gopacket.Payload
)

var connect_timeout = 2000 * time.Millisecond
var http_timeout = 1500 * time.Millisecond
var pcap_timeout = 250 * time.Millisecond

var parser *gopacket.DecodingLayerParser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth_layer, &tcp_layer, &tls_layer, &ip4_layer, &ip6_layer, &payload)

type Config struct {
	URL    string `json:"to_url"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type Server struct {
	id        int
	cdn       string
	name      string
	fqdn      string
	ip4       string
	ip6       string
	httpPort  string
	httpsPort string
	failcount int
	status    string
}

func NewServer(id int, name string, status string, f int) Server {
	server := Server{}
	server.id = id
	server.name = name
	server.status = status
	server.failcount = f
	return server
}

func LoadConfig(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}

func capture(ctx context.Context, s Server, iface *string, ch_dscp chan uint8, ch_ready chan uint8, ip string, ssl bool) {
	var pcap_port string
	if ssl {
		pcap_port = s.httpsPort
	} else {
		pcap_port = s.httpPort
	}

	pcap_filter := "tcp and src " + ip + " and port " + pcap_port + " and (tcp[tcpflags] & tcp-push != 0 or ip6[53] & 8 != 0)"
	rlog.Debugf("capture() filter='%s'", pcap_filter)
	if handle, err := pcap.OpenLive(*iface, 1400, false, pcap_timeout); err != nil {
		rlog.Error("capture() pcap.OpenLive() error:", err)
	} else if err := handle.SetBPFFilter(pcap_filter); err != nil {
		rlog.Error("capture() handle.SetBPFFilter() error:", err)
	} else if err := handle.SetDirection(pcap.DirectionIn); err != nil {
		rlog.Error("capture() handle.SetDirection() error:", err)
	} else {
		defer handle.Close()
		decodedLayers := []gopacket.LayerType{}
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		timer := time.NewTimer(2000 * time.Millisecond)
		ch_ready <- 1 // signal that capture is ready to proceed
		pktCount := 0
		for {
			select {
			case <-ctx.Done():
				handle.Close() // without this, serious fh leak
				rlog.Debug("capture() context cancelled")
				return
			case <-timer.C:
				rlog.Debug("capture() timed out before packets received")
				ch_dscp <- 254
				return
			default:
				pktCount++
				packet, err := packetSource.NextPacket()
				if err == io.EOF {
					break
				} else if err != nil {
					rlog.Error("capture() Error:", err)
					continue
				}
				if sslflag == true && pktCount < 6 {
					// skip the TLS handshake packets - they may not provide real DSCP value
					rlog.Tracef(1, "Packet #%d: %s", pktCount, packet)
					continue
				}
				rlog.Tracef(1, "Packet #%d: %s", pktCount, packet)
				err = parser.DecodeLayers(packet.Data(), &decodedLayers)
				if err != nil {
					rlog.Warn(err)
				}
				for _, typ := range decodedLayers {
					switch typ {
					case layers.LayerTypeIPv4:
						ch_dscp <- ip4_layer.TOS
						return
					case layers.LayerTypeIPv6:
						ch_dscp <- ip6_layer.TrafficClass
						return
					}
				}
			}
		}
	}
}

func protocol_picker(s Server, check_ip string, host_header string, check_path string, v6flag bool) (cap_dscp string) {
	if protocol == 0 || protocol == 2 {
		// if prot is HTTP and HTTPS, just check HTTP - both is overkill
		// do HTTP stuff
		sslflag = false
	}
	if protocol == 1 || protocol == 3 {
		// if a DS is HTTPS *only*...
		// do HTTPS stuff
		sslflag = true
	}
	rlog.Debugf("protocol_picker() ssl=%t", sslflag)
	cap_dscp = request(confInt, s, host_header, check_ip, check_path, v6flag, sslflag)
	return
}

func request(iface *string, s Server, host_header string, ip string, check_path string, v6 bool, ssl bool) string {
	var cap_dscp uint8
	var cap_dscp2 string
	var url string
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // make sure all paths cancel the context to avoid context leak
	ch_dscp := make(chan uint8)
	ch_ready := make(chan uint8)
	go capture(ctx, s, iface, ch_dscp, ch_ready, ip, ssl)
	if ssl {
		url = "https://" + host_header + check_path
	} else {
		url = "http://" + host_header + check_path
	}
	rlog.Infof("request() url=%s ip=%s", url, ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		rlog.Error("request() http.NewRequest() error:", err)
	}
	req.Close = true
	req.Header.Add("Cache-Control", "only-if-cached")
	req.Header.Add("User-Agent", UserAgent)
	ready := <-ch_ready
	if ready == 1 {
		rlog.Debug("request() received go signal from capture()")
	}
	if v6 && ssl == false {
		resp, err := httpClient6.Do(req)
		if err != nil {
			rlog.Error("request() httpClient6.Do() error:", err)
		} else {
			defer resp.Body.Close()
			io.Copy(ioutil.Discard, resp.Body)
		}
	} else if v6 == false && ssl == false {
		resp, err := httpClient4.Do(req)
		if err != nil {
			rlog.Error("request() httpClient4.Do() error:", err)
		} else {
			defer resp.Body.Close()
			io.Copy(ioutil.Discard, resp.Body)
		}
	} else if v6 && ssl {
		resp, err := httpsClient6.Do(req)
		if err != nil {
			rlog.Error("request() httpClient6.Do() error:", err)
		} else {
			defer resp.Body.Close()
			io.Copy(ioutil.Discard, resp.Body)
		}
	} else {
		resp, err := httpsClient4.Do(req)
		if err != nil {
			rlog.Error("request() httpClient4.Do() error:", err)
		} else {
			defer resp.Body.Close()
			io.Copy(ioutil.Discard, resp.Body)
		}
	}
	cap_tos, more := <-ch_dscp
	if more {
		rlog.Debugf("request() received tos=%d", cap_tos)
		cancel() // cancel context to prevent goroutine leak!
	} else {
		rlog.Debug("request() received all dscp values")
		cancel() // cancel context to prevent goroutine leak!
	}
	if cap_tos == 254 {
		rlog.Error("request() no valid DSCP mark received")
		cap_dscp2 = "-1"
	} else {
		cap_dscp = cap_tos >> 2
		cap_dscp2 = strconv.Itoa(int(cap_dscp))
		rlog.Debugf("request() received ipv6=%t dscp=%s", v6, cap_dscp2)
	}
	return cap_dscp2
}

func check_result(want string, have string) bool {
	if want == have {
		rlog.Debugf("check_result() success want=%s got=%s", want, have)
		return true
	} else if have == "-1" {
		rlog.Debugf("check_result() undetermined (ignoring) want=%s got=CAPTURE_TIMEOUT", want)
		return true
	} else {
		rlog.Debugf("check_result() failure want=%s got=%s", want, have)
		return false
	}
}

func main() {
	var (
		cpath_new  string
		cap_dscp   string
		conf_dscp  string
		xmlID      string
		check_ip4  string
		check_ip6  string
		dialer_ip4 string
		dialer_ip6 string
	)

	jobStart := time.Now()

	// define default config file path
	cpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		rlog.Error("Config error:", err)
		os.Exit(1)
	}
	cpath_new = strings.Replace(cpath, "/bin/checks", "/conf/check-config.json", 1)

	// command-line flags
	confPtr := flag.String("conf", cpath_new, "Config file path")
	confInt = flag.String("iface", "undef", "Network interface for packet capture")
	confName := flag.String("name", "DSCP", "Check name to pass to TO, e.g. 'DSCP'")
	confInclude := flag.String("host", "undef", "Specific host or regex to include (optional)")
	confSyslog := flag.Bool("syslog", false, "Log check results to syslog")
	confCdn := flag.String("cdn", "all", "Check specific CDN by name")
	confExclude := flag.String("exclude", "undef", "Hostname regex to exclude")
	confReset := flag.Bool("reset", false, "Reset check values in TO to 'blank' state")
	confQuiet := flag.Bool("q", false, "Do not send updates to TO")
	flag.Parse()

	// configure syslog logger
	if *confSyslog == true {
		logwriter, err := syslog.New(syslog.LOG_INFO, os.Args[0])
		if err == nil {
			log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
			log.SetOutput(logwriter)
		}
	}

	if *confInt == "undef" {
		rlog.Error("Must specify network interface for packet capture")
		os.Exit(1)
	}
	if *confName == "undef" {
		rlog.Error("Must specify check name for update to send to TO")
		os.Exit(1)
	}

	// load config json
	config, err := LoadConfig(*confPtr)
	if err != nil {
		rlog.Error("Error loading config:", err)
		os.Exit(1)
	}

	// connect to TO
	session, _, err := toclient.LoginWithAgent(
		config.URL,
		config.User,
		config.Passwd,
		AllowInsecureConnections,
		UserAgent,
		UseClientCache,
		TrafficOpsRequestTimeout)
	if err != nil {
		rlog.Criticalf("An error occurred while logging in: %v\n", err)
		os.Exit(1)
	}

	// Make TO API call for server details
	var servers tc.ServersV3Response
	servers, _, err = session.GetServersWithHdr(nil, nil)
	if err != nil {
		rlog.Criticalf("An error occurred while getting servers: %v\n", err)
		os.Exit(1)
	}

	// Make TO API call for delivery service details
	var deliveryservices []tc.DeliveryServiceNullableV30
	deliveryservices, _, err = session.GetDeliveryServicesV30WithHdr(nil, nil)
	if err != nil {
		rlog.Criticalf("An error occurred while getting delivery services: %v\n", err)
		os.Exit(1)
	}

	// Make TO API call for cdn details
	var cdns []tc.CDN
	cdns, _, err = session.GetCDNs()
	if err != nil {
		rlog.Criticalf("An error occurred while getting cdns: %v\n", err)
		os.Exit(1)
	}

	// map cdn to domain name
	cdn_map := make(map[string]string)
	for _, cdn := range cdns {
		cdn_map[cdn.Name] = cdn.DomainName
	}

	// map ds info
	ds_matchlist := make(map[string][]tc.DeliveryServiceMatch)
	ds_types := make(map[string]tc.DSType)
	for _, ds := range deliveryservices {
		ds_matchlist[*ds.XMLID] = *ds.MatchList
		ds_types[*ds.XMLID] = *ds.Type
	}

	for _, server := range servers.Response {
		re, err := regexp.Compile("^EDGE.*")
		if err != nil {
			rlog.Error("supplied exclusion regex does not compile:", err)
			os.Exit(1)
		}
		if re.Match([]byte(server.Type)) {
			serverStart := time.Now()
			if *confInclude != "undef" {
				re_inc, err := regexp.Compile(*confInclude)
				if err != nil {
					rlog.Error("supplied exclusion regex does not compile:", err)
					os.Exit(1)
				}
				if !re_inc.MatchString(*server.HostName) {
					rlog.Debugf("%s does not match the provided include regex, skipping", server.HostName)
					continue
				}
			}
			if *confCdn != "all" && *confCdn != *server.CDNName {
				rlog.Debugf("%s is not assinged to the specified CDN '%s', skipping", server.HostName, *confCdn)
				continue
			}
			if *confExclude != "undef" {
				re, err := regexp.Compile(*confExclude)
				if err != nil {
					rlog.Error("supplied exclusion regex does not compile:", err)
					os.Exit(1)
				}
				if re.MatchString(*server.HostName) {
					rlog.Debugf("%s matches the provided exclude regex, skipping", server.HostName)
					continue
				}
			}
			s := NewServer(*server.ID, *server.HostName, *server.Status, -1)
			doV4 := false //default
			doV6 := false //default
			defaulStatusValue := -1
			var statusData tc.ServercheckRequestNullable
			statusData.ID = &s.id
			statusData.Name = confName
			statusData.HostName = &s.name
			statusData.Value = &defaulStatusValue
			s.fqdn = s.name + "." + *server.DomainName
			rlog.Infof("Next server=%s status=%s", s.fqdn, s.status)
			if *confSyslog {
				log.Printf("Next server=%s status=%s", s.fqdn, s.status)
			}
			if (s.status == "REPORTED" || s.status == "ADMIN_DOWN") && *confReset != true {
				s.failcount = 0
				s.cdn = *server.CDNName
				s.httpPort = strconv.Itoa(*server.TCPPort)
				s.httpsPort = strconv.Itoa(*server.HTTPSPort)
				for _, interf := range server.Interfaces {
					for _, addr := range interf.IPAddresses {
						if s.ip4 == "" && strings.Count(addr.Address, ":") == 0 {
							s.ip4 = strings.Split(addr.Address, "/")[0]
						}
						if s.ip6 == "" && strings.Count(addr.Address, ":") > 0 {
							s.ip6 = strings.Split(addr.Address, "/")[0]
						}
					}
				}
				rlog.Debugf("Ports for %s: http=%s https=%s", s.name, http_port, https_port)
				services, _, err := session.GetDeliveryServicesByServerV30WithHdr(s.id, nil)
				if err != nil {
					rlog.Error("Error getting delivery services from TO:", err)
					os.Exit(1)
				}

				dialer := &net.Dialer{
					Timeout:       connect_timeout,
					KeepAlive:     -1,
					DualStack:     false,
					FallbackDelay: -1,
				}
				if s.ip4 != "" {
					doV4 = true
					check_ip4 = s.ip4
					dialer_ip4 = s.ip4
				}
				if s.ip6 != "" {
					doV6 = true
					check_ip6 = s.ip6
					dialer_ip6 = "[" + s.ip6 + "]"
				}

				// it is necessary to define custom Transports in order to support
				// TLS SNI in go. Otherwise, we could have just had the http client
				// connect to the IPv4 or IPv6 address, and set a custom Host header
				// to ID the test target. SNI must be done this way, however, so in
				// order to be consistent, just set up Transports for all protocol
				// combinations.
				http4 = &http.Transport{
					DisableKeepAlives: true,
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						addr = dialer_ip4 + ":" + s.httpPort
						return dialer.DialContext(ctx, network, addr)
					},
				}
				http6 = &http.Transport{
					DisableKeepAlives: true,
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						addr = dialer_ip6 + ":" + s.httpPort
						return dialer.DialContext(ctx, network, addr)
					},
				}
				https4 = &http.Transport{
					DisableKeepAlives: true,
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						addr = dialer_ip4 + ":" + s.httpsPort
						return dialer.DialContext(ctx, network, addr)
					},
				}
				https6 = &http.Transport{
					DisableKeepAlives: true,
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
					DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
						addr = dialer_ip6 + ":" + s.httpsPort
						return dialer.DialContext(ctx, network, addr)
					},
				}

				// The Client isntances are tied to each Transport instance, and are also
				// configured to prevent following any HTTP redirects.
				httpClient4 = &http.Client{
					Timeout:   http_timeout,
					Transport: http4,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}
				httpClient6 = &http.Client{
					Timeout:   http_timeout,
					Transport: http6,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}
				httpsClient4 = &http.Client{
					Timeout:   http_timeout,
					Transport: https4,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}
				httpsClient6 = &http.Client{
					Timeout:   http_timeout,
					Transport: https6,
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}

				for _, service := range services {
					xmlID = *service.XMLID
					if *service.Active == false {
						rlog.Infof("Skipping ds=%s active=false", xmlID)
						continue
					} else if *service.DSCP == 0 {
						// routers may override with default mark in this case
						rlog.Infof("Skipping ds=%s dscp=0", xmlID)
						continue
					} else if *service.CheckPath == "" {
						rlog.Infof("Skipping ds=%s no check path set", xmlID)
						continue
					}
					if matched, _ := regexp.Match(`^/`, []byte(*service.CheckPath)); matched == false {
						//prepend leading slash if missing
						*service.CheckPath = "/" + *service.CheckPath
					}
					protocol = *service.Protocol
					conf_dscp = strconv.Itoa(*service.DSCP)
					check_path := service.CheckPath
					routing_name := service.RoutingName
					rlog.Infof("checking ds=%s server=%s cdn=%s dscp=%s", xmlID, s.fqdn, s.cdn, conf_dscp)
					for _, match := range ds_matchlist[xmlID] {
						if match.Type == "HOST_REGEXP" {
							if matched, err := regexp.MatchString(`\*`, match.Pattern); err != nil {
								rlog.Error(err)
							} else if matched == true {
								re := regexp.MustCompile(`(\\|\.\*)`)
								host_header = re.ReplaceAllString(match.Pattern, "")
								matched, err = regexp.MatchString(`^DNS.*`, string(ds_types[xmlID]))
								if err != nil {
									rlog.Error(err)
								}
								if matched == true {
									host_header = *routing_name + host_header + cdn_map[s.cdn]
								} else {
									host_header = s.name + host_header + cdn_map[s.cdn]
								}
							} else {
								host_header = match.Pattern
							}
						}
					}
					var v6flag bool
					if doV4 {
						// do IPv4 stuff
						v6flag = false
						cap_dscp = protocol_picker(s, check_ip4, host_header, *check_path, v6flag)
						success := check_result(conf_dscp, cap_dscp)
						if success == false {
							// retry to be sure - something like out-of-order packets may have been an issue
							rlog.Info("first IPv4 check failed - retrying")
							cap_dscp = protocol_picker(s, check_ip4, host_header, *check_path, v6flag)
							success = check_result(conf_dscp, cap_dscp)
						}
						if success == false {
							rlog.Infof("result=failure type=ipv4 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip4, conf_dscp, cap_dscp)
							if *confSyslog {
								log.Printf("result=failure type=ipv4 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip4, conf_dscp, cap_dscp)
							}
							s.failcount++
						} else {
							if cap_dscp == "-1" {
								cap_dscp = "CAPTURE_TIMEOUT (IGNORING)"
							}
							rlog.Infof("result=success type=ipv4 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip4, conf_dscp, cap_dscp)
							if *confSyslog {
								log.Printf("result=success type=ipv4 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip4, conf_dscp, cap_dscp)
							}
						}
					}
					if doV6 {
						// do IPv6 stuff
						v6flag = true
						cap_dscp = protocol_picker(s, check_ip6, host_header, *check_path, v6flag)
						success := check_result(conf_dscp, cap_dscp)
						if success == false {
							// retry to be sure - something like out-of-order packets may have been an issue
							rlog.Info("first IPv6 check failed - retrying")
							cap_dscp = protocol_picker(s, check_ip6, host_header, *check_path, v6flag)
							success = check_result(conf_dscp, cap_dscp)
						}
						if success == false {
							rlog.Infof("result=failure type=ipv6 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip6, conf_dscp, cap_dscp)
							if *confSyslog {
								log.Printf("result=failure type=ipv6 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip6, conf_dscp, cap_dscp)
							}
							s.failcount++
						} else {
							if cap_dscp == "-1" {
								cap_dscp = "CAPTURE_TIMEOUT (IGNORING)"
							}
							rlog.Infof("result=success type=ipv6 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip6, conf_dscp, cap_dscp)
							if *confSyslog {
								log.Printf("result=success type=ipv6 server=%s cdn=%s ds=%s ip=%s want=%s got=%s", s.fqdn, s.cdn, xmlID, check_ip6, conf_dscp, cap_dscp)
							}
						}
					}
				}
			}
			// send status update to TO
			if s.failcount == -1 {
				// server not checked
				*statusData.Value = -1
			} else if s.failcount > 0 {
				// server had failures
				*statusData.Value = 0
			} else {
				// server looks OK
				*statusData.Value = 1
			}
			serverElapsed := time.Since(serverStart)
			rlog.Infof("Finished checking server=%s result=%d cdn=%s elapsed=%s", s.fqdn, *statusData.Value, s.cdn, serverElapsed)
			if *confSyslog {
				log.Printf("Finished checking server=%s result=%d cdn=%s elapsed=%s", s.fqdn, *statusData.Value, s.cdn, serverElapsed)
			}
			if *confQuiet == false {
				rlog.Debug("Sending update to TO")
				_, _, err := session.InsertServerCheckStatus(statusData)
				if err != nil {
					rlog.Error("Error updating server check status with TO:", err)
				}
			} else {
				rlog.Debug("Skipping update to TO")
			}
		}
	}
	jobElapsed := time.Since(jobStart)
	rlog.Info("Job complete", jobElapsed)
	if *confSyslog {
		log.Print("Job complete totaltime=", jobElapsed)
	}
	os.Exit(0)
}
