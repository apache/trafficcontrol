package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"net/url"
	"os"
	"strconv"

	"golang.org/x/net/publicsuffix"

	"github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/client"
)

// GetClient returns a TO Client, using a cached cookie if it exists, or logging in otherwise
func GetClient(toURL string, toUser string, toPass string, tempDir string) (*toclient.Session, error) {
	cookies, err := GetCookiesFromFile(tempDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG failed to get cookies from cache file (trying real TO): "+err.Error()+"\n")
		cookies = ""
	}

	if cookies == "" {
		err := error(nil)
		cookies, err = GetCookiesFromTO(toURL, toUser, toPass, tempDir)
		if err != nil {
			return nil, errors.New("getting cookies from Traffic Ops: " + err.Error())
		}
		fmt.Fprintf(os.Stderr, "DEBUG using cookies from TO\n")
	} else {
		fmt.Fprintf(os.Stderr, "DEBUG using cookies from cache file\n")
	}

	useCache := false
	toClient := toclient.NewNoAuthSession(toURL, TOInsecure, UserAgent, useCache, TOTimeout)
	toClient.UserName = toUser
	toClient.Password = toPass

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.New("making cookie jar: " + err.Error())
	}
	toClient.Client.Jar = jar

	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return nil, errors.New("parsing Traffic Ops URL '" + toURL + "': " + err.Error())
	}

	toClient.Client.Jar.SetCookies(toURLParsed, StringToCookies(cookies))
	return toClient, nil
}

// GetCookies gets the cookies from logging in to Traffic Ops.
// If this succeeds, it also writes the cookies to TempSubdir/TempCookieFileName.
func GetCookiesFromTO(toURL string, toUser string, toPass string, tempDir string) (string, error) {
	toURLParsed, err := url.Parse(toURL)
	if err != nil {
		return "", errors.New("parsing Traffic Ops URL '" + toURL + "': " + err.Error())
	}

	toUseCache := false
	toClient, toIP, err := toclient.LoginWithAgent(toURL, toUser, toPass, TOInsecure, UserAgent, toUseCache, TOTimeout)
	if err != nil {
		toIPStr := ""
		if toIP != nil {
			toIPStr = toIP.String()
		}
		return "", errors.New("logging in to Traffic Ops IP '" + toIPStr + "': " + err.Error())
	}

	cookiesStr := CookiesToString(toClient.Client.Jar.Cookies(toURLParsed))
	WriteCookiesToFile(cookiesStr, tempDir)

	return cookiesStr, nil
}

// TrafficOpsRequest makes a request to Traffic Ops for the given method, url, and body.
// If it gets an Unauthorized or Forbidden, it tries to log in again and makes the request again.
func TrafficOpsRequest(toClient **toclient.Session, cfg Cfg, method string, url string, body []byte) (string, error) {
	resp, toIP, err := rawTrafficOpsRequest(*toClient, method, url, body)
	if err != nil {
		toIPStr := ""
		if toIP != nil {
			toIPStr = toIP.String()
		}
		return "", errors.New("requesting from Traffic Ops '" + toIPStr + "': " + err.Error())
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		resp.Body.Close()

		fmt.Fprintf(os.Stderr, "DEBUG TrafficOpsRequest got unauthorized/forbidden, logging in again\n")
		fmt.Fprintf(os.Stderr, "DEBUG TrafficOpsRequest url '%v' user '%v' pass '%v'\n", (*toClient).URL, (*toClient).UserName, (*toClient).Password)

		useCache := false
		newTOClient, toIP, err := toclient.LoginWithAgent((*toClient).URL, (*toClient).UserName, (*toClient).Password, TOInsecure, UserAgent, useCache, TOTimeout)
		if err != nil {
			toIPStr := ""
			if toIP != nil {
				toIPStr = toIP.String()
			}
			return "", errors.New("logging in to Traffic Ops IP '" + toIPStr + "': " + err.Error())
		}
		*toClient = newTOClient

		resp, toIP, err = rawTrafficOpsRequest(*toClient, method, url, body)
		if err != nil {
			toIPStr := ""
			if toIP != nil {
				toIPStr = toIP.String()
			}
			return "", errors.New("requesting from Traffic Ops '" + toIPStr + "': " + err.Error())
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			bts = []byte("(read failure)") // if it's a non-200 and the body read fails, don't error, just note the read fail in the error
		}
		return "", errors.New("Traffic Ops returned non-200 code '" + strconv.Itoa(resp.StatusCode) + "' body '" + string(bts) + "'")
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		toIPStr := ""
		if toIP != nil {
			toIPStr = toIP.String()
		}
		return "", errors.New("reading body from Traffic Ops '" + toIPStr + "': " + err.Error())
	}

	return string(bts), nil
}

// rawTrafficOpsRequest makes a request to Traffic Ops for the given method, url, and body.
// If it gets an Unauthorized or Forbidden, it tries to log in again and makes the request again.
func rawTrafficOpsRequest(toClient *toclient.Session, method string, url string, body []byte) (*http.Response, net.Addr, error) {
	bodyReader := io.Reader(nil)
	if len(body) > 0 {
		bodyReader = bytes.NewBuffer(body)
	}

	remoteAddr := net.Addr(nil)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, remoteAddr, err
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr()
		},
	}))

	req.Header.Set("User-Agent", toClient.UserAgentStr)

	resp, err := toClient.Client.Do(req)
	if err != nil {
		return nil, remoteAddr, err
	}

	return resp, remoteAddr, nil
}

// MaybeIPStr returns the Traffic Ops IP string if it isn't nil, or the empty string if it is.
func MaybeIPStr(reqInf toclient.ReqInf) string {
	if reqInf.RemoteAddr != nil {
		return reqInf.RemoteAddr.String()
	}
	return ""
}

// TCParamsToParamsWithProfiles unmarshals the Profiles that the tc struct doesn't.
func TCParamsToParamsWithProfiles(tcParams []tc.Parameter) ([]ParameterWithProfiles, error) {
	params := make([]ParameterWithProfiles, 0, len(tcParams))
	for _, tcParam := range tcParams {
		param := ParameterWithProfiles{Parameter: tcParam}

		profiles := []string{}
		if err := json.Unmarshal(tcParam.Profiles, &profiles); err != nil {
			return nil, errors.New("unmarshalling JSON from parameter '" + strconv.Itoa(param.ID) + "': " + err.Error())
		}
		param.ProfileNames = profiles
		param.Profiles = nil
		params = append(params, param)
	}
	return params, nil
}

type ParameterWithProfiles struct {
	tc.Parameter
	ProfileNames []string
}
