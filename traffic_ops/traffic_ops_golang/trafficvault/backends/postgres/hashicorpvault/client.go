package hashicorpvault

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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
)

const (
	defaultTimeout   = 30 * time.Second
	userAgent        = "TrafficOps/6.0"
	vaultTokenHeader = "X-Vault-Token"
)

type Client struct {
	address    string
	roleID     string
	secretID   string
	token      string
	httpClient *http.Client
	loginPath  string
	secretPath string
}

func NewClient(address, roleID, secretID, loginPath, secretPath string, timeout time.Duration, insecure bool) *Client {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	res := Client{
		address:  address,
		roleID:   roleID,
		secretID: secretID,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: insecure, MinVersion: tls.VersionTLS12},
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		loginPath:  loginPath,
		secretPath: secretPath,
	}
	return &res
}

type appRoleLoginRequest struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

type appRoleLoginResponse struct {
	Auth   auth     `json:"auth"`
	Errors []string `json:"errors"`
}

type auth struct {
	ClientToken string `json:"client_token"`
}

func (c *Client) Login() error {
	data := appRoleLoginRequest{
		RoleID:   c.roleID,
		SecretID: c.secretID,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return errors.New("marshalling login request body: " + err.Error())
	}
	requestURL := c.getURL(c.loginPath)
	resp, remoteAddr, err := c.doRequest(http.MethodPost, requestURL, body)
	if err != nil {
		return fmt.Errorf("doing login HTTP request (addr = %s): %s", remoteAddr, err.Error())
	}
	defer log.Close(resp.Body, "closing HashiCorp Vault login response body")
	loginResp := appRoleLoginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return fmt.Errorf("decoding HashCorp Vault login response body (addr = %s): %s", remoteAddr, err.Error())
	}
	if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
		errs := strings.Join(loginResp.Errors, ", ")
		return fmt.Errorf("login attempt (addr = %s) returned status code: %s, errors: %s", remoteAddr, resp.Status, errs)
	}
	if loginResp.Auth.ClientToken == "" {
		return fmt.Errorf("login response body contained empty auth.client_token (addr = %s)", remoteAddr)
	}
	c.token = loginResp.Auth.ClientToken
	log.Infof("successfully authenticated to HashiCorp Vault (addr = %s)", remoteAddr)
	return nil
}

type secretResponse struct {
	Data   secretData `json:"data"`
	Errors []string   `json:"errors"`
}

type secretData struct {
	Data secretKeyValue `json:"data"`
}

type secretKeyValue struct {
	TrafficVaultKey string `json:"traffic_vault_key"`
}

func (c *Client) GetSecret() (string, error) {
	requestURL := c.getURL(c.secretPath)
	resp, remoteAddr, err := c.doRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("doing secret HTTP request (addr = %s): %s", remoteAddr, err.Error())
	}
	defer log.Close(resp.Body, "closing HashiCorp Vault secret response body")
	secretResp := secretResponse{}
	err = json.NewDecoder(resp.Body).Decode(&secretResp)
	if err != nil {
		return "", fmt.Errorf("decoding HashCorp Vault secret response body (addr = %s): %s", remoteAddr, err.Error())
	}
	if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
		errs := strings.Join(secretResp.Errors, ", ")
		return "", fmt.Errorf("attempting to get secret (addr = %s) returned status code: %s, errors: %s", remoteAddr, resp.Status, errs)
	}
	if secretResp.Data.Data.TrafficVaultKey == "" {
		return "", fmt.Errorf("secret response body contained empty traffic_vault_key (addr = %s)", remoteAddr)
	}
	log.Infof("successfully retrieved secret traffic_vault_key from HashiCorp Vault (addr = %s)", remoteAddr)
	return secretResp.Data.Data.TrafficVaultKey, nil
}

func (c *Client) doRequest(method, url string, body []byte) (*http.Response, string, error) {
	remoteAddr := ""
	var resp *http.Response
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, "", errors.New("creating http request: " + err.Error())
		}
		req.Header.Set(rfc.ContentType, rfc.ApplicationJSON)
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, "", errors.New("creating http request: " + err.Error())
		}
	}
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			remoteAddr = connInfo.Conn.RemoteAddr().String()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	req.Header.Set(rfc.UserAgent, userAgent)
	if c.token != "" {
		req.Header.Set(vaultTokenHeader, c.token)
	}
	resp, err = c.httpClient.Do(req)
	return resp, remoteAddr, err
}

func (c *Client) getURL(path string) string {
	return strings.TrimSuffix(c.address, "/") + "/" + strings.TrimPrefix(path, "/")
}
