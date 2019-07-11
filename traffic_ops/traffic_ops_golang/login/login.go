package login

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tocookie"

	"github.com/jmoiron/sqlx"
)

func LoginHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		defer r.Body.Close()
		authenticated := false
		form := auth.PasswordForm{}
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			handleErrs(http.StatusBadRequest, err)
			return
		}
		resp := struct {
			tc.Alerts
		}{}
		userAllowed, err, blockingErr := auth.CheckLocalUserIsAllowed(form, db, time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second)
		if blockingErr != nil {
			api.HandleErr(w, r, nil, http.StatusServiceUnavailable, nil, fmt.Errorf("error checking local user password: %s\n", blockingErr.Error()))
			return
		}
		if err != nil {
			log.Errorf("checking local user: %s\n", err.Error())
		}
		if userAllowed {
			authenticated, err, blockingErr = auth.CheckLocalUserPassword(form, db, time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second)
			if blockingErr != nil {
				api.HandleErr(w, r, nil, http.StatusServiceUnavailable, nil, fmt.Errorf("error checking local user password: %s\n", blockingErr.Error()))
				return
			}
			if err != nil {
				log.Errorf("checking local user password: %s\n", err.Error())
			}
			var ldapErr error
			if !authenticated {
				if cfg.LDAPEnabled {
					authenticated, ldapErr = auth.CheckLDAPUser(form, cfg.ConfigLDAP)
					if ldapErr != nil {
						log.Errorf("checking ldap user: %s\n", ldapErr.Error())
					}
				}
			}
			if authenticated {
				expiry := time.Now().Add(time.Hour * 6)
				cookie := tocookie.New(form.Username, expiry, cfg.Secrets[0])
				httpCookie := http.Cookie{Name: "mojolicious", Value: cookie, Path: "/", Expires: expiry, HttpOnly: true}
				http.SetCookie(w, &httpCookie)
				resp = struct {
					tc.Alerts
				}{tc.CreateAlerts(tc.SuccessLevel, "Successfully logged in.")}
			} else {
				resp = struct {
					tc.Alerts
				}{tc.CreateAlerts(tc.ErrorLevel, "Invalid username or password.")}
			}
		} else {
			resp = struct {
				tc.Alerts
			}{tc.CreateAlerts(tc.ErrorLevel, "Invalid username or password.")}
		}
		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		if !authenticated {
			w.WriteHeader(http.StatusUnauthorized)
		}
		fmt.Fprintf(w, "%s", respBts)
	}
}

// OauthLoginHandler accepts a JSON web token previously obtained from an OAuth provider, decodes it, validates it, authorizes the user against the database, and returns the login result as either an error or success message
func OauthLoginHandler(db *sqlx.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)
		defer r.Body.Close()
		authenticated := false
		resp := struct {
			tc.Alerts
		}{}

		form := auth.PasswordForm{}
		parameters := struct {
			AuthCodeTokenUrl string `json:"authCodeTokenUrl"`
			Code             string `json:"code"`
			ClientId         string `json:"clientId"`
			ClientSecret     string `json:"clientSecret"`
			RedirectUri      string `json:"redirectUri"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&parameters); err != nil {
			handleErrs(http.StatusBadRequest, err)
			return
		}

		data := url.Values{}
		data.Add("code", parameters.Code)
		data.Add("client_id", parameters.ClientId)
		data.Add("client_secret", parameters.ClientSecret)
		data.Add("grant_type", "authorization_code") // Required by RFC6749 section 4.1.3
		data.Add("redirect_uri", parameters.RedirectUri)

		req, err := http.NewRequest("POST", parameters.AuthCodeTokenUrl, bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			log.Errorf("obtaining token using code from oauth provider\n%s", err.Error())
			return
		}

		client := http.Client{}
		response, err := client.Do(req)
		if err != nil {
			log.Errorf("getting an http client\n%s", err.Error())
			return
		}
		defer response.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		encodedToken := ""

		var result map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
			log.Warnf("Error parsing JSON response from oAuth: %s", err.Error())
			encodedToken = buf.String()
		} else if _, ok := result["access_token"]; !ok {
			sysErr := fmt.Errorf("Missing access token in response:\n%s\n", buf.String())
			usrErr := errors.New("Bad response from OAuth2.0 provider")
			api.HandleErr(w, r, nil, http.StatusBadGateway, usrErr, sysErr)
			return
		} else {
			switch t := result["access_token"].(type) {
			case string:
				encodedToken = result["access_token"].(string)
			default:
				sysErr := fmt.Errorf("Incorrect type of access_token! Expected 'string', got '%v'\n", t)
				usrErr := errors.New("Bad response from OAuth2.0 provider")
				api.HandleErr(w, r, nil, http.StatusBadGateway, usrErr, sysErr)
				return
			}
		}

		if encodedToken == "" {
			log.Errorf("Token not found in request but is required")
			handleErrs(http.StatusBadRequest, errors.New("Token not found in request but is required"))
			return
		}

		decodedToken, err := jwt.Parse(encodedToken, func(unverifiedToken *jwt.Token) (interface{}, error) {
			publicKeyUrl := unverifiedToken.Header["jku"].(string)
			publicKeyId := unverifiedToken.Header["kid"].(string)

			matched, err := VerifyUrlOnWhiteList(publicKeyUrl, cfg.ConfigTrafficOpsGolang.WhitelistedOAuthUrls)
			if err != nil {
				return nil, err
			}
			if !matched {
				return nil, errors.New("Key URL from token is not included in the whitelisted urls. Received: " + publicKeyUrl)
			}

			keys, err := jwk.FetchHTTP(publicKeyUrl)
			if err != nil {
				return nil, errors.New("Error fetching JSON key set with message: " + err.Error())
			}

			keyById := keys.LookupKeyID(publicKeyId)
			if len(keyById) == 0 {
				return nil, errors.New("No public key found for id: " + publicKeyId + " at url: " + publicKeyUrl)
			}

			selectedKey, err := keyById[0].Materialize()
			if err != nil {
				return nil, errors.New("Error materializing key from JSON key set with message: " + err.Error())
			}

			return selectedKey, nil
		})
		if err != nil {
			handleErrs(http.StatusInternalServerError, errors.New("Error decoding token with message: "+err.Error()))
			log.Errorf("Error decoding token: %s\n", err.Error())
			return
		}

		authenticated = decodedToken.Valid

		userId := decodedToken.Claims.(jwt.MapClaims)["sub"].(string)
		form.Username = userId

		userAllowed, err, blockingErr := auth.CheckLocalUserIsAllowed(form, db, time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second)
		if blockingErr != nil {
			api.HandleErr(w, r, nil, http.StatusServiceUnavailable, nil, fmt.Errorf("error checking local user password: %s\n", blockingErr.Error()))
			return
		}
		if err != nil {
			log.Errorf("checking local user: %s\n", err.Error())
			return
		}

		if userAllowed && authenticated {
			expiry := time.Now().Add(time.Hour * 6)
			cookie := tocookie.New(userId, expiry, cfg.Secrets[0])
			httpCookie := http.Cookie{Name: "mojolicious", Value: cookie, Path: "/", Expires: expiry, HttpOnly: true}
			http.SetCookie(w, &httpCookie)
			resp = struct {
				tc.Alerts
			}{tc.CreateAlerts(tc.SuccessLevel, "Successfully logged in.")}
		} else {
			resp = struct {
				tc.Alerts
			}{tc.CreateAlerts(tc.ErrorLevel, "Invalid username or password.")}
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}
		w.Header().Set(tc.ContentType, tc.ApplicationJson)
		if !authenticated {
			w.WriteHeader(http.StatusUnauthorized)
		}
		fmt.Fprintf(w, "%s", respBts)

	}
}

func VerifyUrlOnWhiteList(urlString string, whiteListedUrls []string) (bool, error) {

	for _, listing := range whiteListedUrls {
		if listing == "" {
			continue
		}

		urlParsed, err := url.Parse(urlString)
		if err != nil {
			return false, err
		}

		matched, err := filepath.Match(listing, urlParsed.Hostname())
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}
	return false, nil
}
