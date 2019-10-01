package deliveryservice

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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"html/template"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type DsKey struct {
	XmlId   string
	Version sql.NullInt64
}

type DsExpirationInfo struct {
	XmlId      string
	Version    util.JSONIntStr
	Expiration time.Time
	AuthType   string
	Error      error
}

type ExpirationSummary struct {
	LetsEncryptExpirations []DsExpirationInfo
	SelfSignedExpirations  []DsExpirationInfo
	OtherExpirations       []DsExpirationInfo
}

func RenewCertificates(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusServiceUnavailable, errors.New("the Riak service is unavailable"), errors.New("getting SSL keys from Riak by xml id: Riak is not configured"))
		return
	}

	rows, err := inf.Tx.Tx.Query(`SELECT xml_id, ssl_key_version FROM deliveryservice`)
	if err != nil {
		log.Errorf("querying: %v", err)
		return
	}
	defer rows.Close()

	dses := []DsKey{}
	for rows.Next() {
		ds := DsKey{}
		err := rows.Scan(&ds.XmlId, &ds.Version)
		if err != nil {
			log.Errorf("getting delivery services: %v", err)
			continue
		}
		if ds.Version.Valid && int(ds.Version.Int64) != 0 {
			dses = append(dses, ds)
		}
	}

	keysFound := ExpirationSummary{}
	for _, ds := range dses {
		dsExpInfo := DsExpirationInfo{}
		keyObj, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(ds.XmlId, strconv.Itoa(int(ds.Version.Int64)), inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort)
		if err != nil {
			log.Errorf("getting ssl keys for xmlId: " + ds.XmlId + " and version: " + strconv.Itoa(int(ds.Version.Int64)) + " :" + err.Error())
			dsExpInfo.XmlId = ds.XmlId
			dsExpInfo.Version = util.JSONIntStr(int(ds.Version.Int64))
			dsExpInfo.Error = errors.New("getting ssl keys for xmlId: " + ds.XmlId + " and version: " + strconv.Itoa(int(ds.Version.Int64)) + " :" + err.Error())
			continue
		}
		if !ok {
			log.Errorf("no object found for the specified key with xmlId: " + ds.XmlId + " and version: " + strconv.Itoa(int(ds.Version.Int64)))
			dsExpInfo.XmlId = ds.XmlId
			dsExpInfo.Version = util.JSONIntStr(int(ds.Version.Int64))
			dsExpInfo.Error = errors.New("no object found for the specified key with xmlId: " + ds.XmlId + " and version: " + strconv.Itoa(int(ds.Version.Int64)))
			continue
		}

		err = base64DecodeCertificate(&keyObj.Certificate)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting SSL keys for XMLID '"+ds.XmlId+"': "+err.Error()))
			return
		}
		if keyObj.Expiration.IsZero() {
			expiration, err := parseExpirationFromCert([]byte(keyObj.Certificate.Crt))
			if err != nil {
				api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New(ds.XmlId+": "+err.Error()))
				return
			}
			keyObj.Expiration = expiration
		}

		// Renew only certificates within configured limit plus 3 days
		if keyObj.Expiration.After(time.Now().Add(time.Hour * 24 * time.Duration(inf.Config.ConfigLetsEncrypt.RenewDaysBeforeExpiration)).Add(time.Hour * 24 * 3)) {
			continue
		}

		newVersion := util.JSONIntStr(keyObj.Version.ToInt64() + 1)

		dsExpInfo.XmlId = keyObj.DeliveryService
		dsExpInfo.Version = keyObj.Version
		dsExpInfo.Expiration = keyObj.Expiration
		dsExpInfo.AuthType = keyObj.AuthType

		if keyObj.AuthType == tc.LetsEncryptAuthType {
			req := tc.DeliveryServiceLetsEncryptSSLKeysReq{
				DeliveryServiceSSLKeysReq: tc.DeliveryServiceSSLKeysReq{
					HostName:        &keyObj.Hostname,
					DeliveryService: &keyObj.DeliveryService,
					CDN:             &keyObj.CDN,
					Version:         &newVersion,
				},
			}
			ctx, _ := context.WithTimeout(r.Context(), GetLetsEncryptTimeout())

			if error := GetLetsEncryptCertificates(inf.Config, req, ctx, inf.User); error != nil {
				dsExpInfo.Error = error
			}
			keysFound.LetsEncryptExpirations = append(keysFound.LetsEncryptExpirations, dsExpInfo)

		} else if keyObj.AuthType == tc.SelfSignedCertAuthType {
			if inf.Config.ConfigLetsEncrypt.ConvertSelfSigned {
				req := tc.DeliveryServiceLetsEncryptSSLKeysReq{
					DeliveryServiceSSLKeysReq: tc.DeliveryServiceSSLKeysReq{
						HostName:        &keyObj.Hostname,
						DeliveryService: &keyObj.DeliveryService,
						CDN:             &keyObj.CDN,
						Version:         &newVersion,
					},
				}
				ctx, _ := context.WithTimeout(r.Context(), GetLetsEncryptTimeout())

				if error := GetLetsEncryptCertificates(inf.Config, req, ctx, inf.User); error != nil {
					dsExpInfo.Error = error
				}
			}
			keysFound.SelfSignedExpirations = append(keysFound.SelfSignedExpirations, dsExpInfo)
		} else {
			keysFound.OtherExpirations = append(keysFound.OtherExpirations, dsExpInfo)
		}

	}

	if inf.Config.ConfigSmtp.Enabled && inf.Config.ConfigLetsEncrypt.SendExpEmail {
		err = AlertExpiringCerts(keysFound, *inf.Config)
		if err != nil {
			log.Errorf(err.Error())
			api.HandleErr(w, r, nil, http.StatusInternalServerError, err, nil)
		}
	}

	api.WriteResp(w, r, keysFound)

}

func AlertExpiringCerts(certsFound ExpirationSummary, config config.Config) error {
	email := strings.Join(config.ConfigSmtp.ToEmail, ",")
	if config.ConfigLetsEncrypt.Email != "" {
		email = config.ConfigLetsEncrypt.Email
	}
	header := "From: " + config.ConfigSmtp.FromEmail + "\n" +
		"To: " + email + "\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"Subject: Certificate Expiration Summary\n\n"

	error := SendEmail(config, header, certsFound)
	if error != nil {
		return error
	}

	return nil
}

func SendEmail(config config.Config, header string, data interface{}) error {
	var auth smtp.Auth
	if config.ConfigSmtp.User != "" {
		auth = LoginAuth("", config.ConfigSmtp.User, config.ConfigSmtp.Password, strings.Split(config.ConfigSmtp.Address, ":")[0])
	}

	email := config.ConfigSmtp.ToEmail
	if config.ConfigLetsEncrypt.Email != "" {
		email = []string{config.ConfigLetsEncrypt.Email}
	}

	msgBodyBuffer, err := parseTemplate("/opt/traffic_ops/app/templates/send_mail/autorenewcerts_mail.ep", data)
	if err != nil {
		return err
	}
	msg := append([]byte(header), msgBodyBuffer.Bytes()...)

	error := smtp.SendMail(config.ConfigSmtp.Address, auth, config.ConfigSmtp.FromEmail, email, []byte(msg))
	if error != nil {
		return errors.New("Failed to send email: " + error.Error())
	}
	return nil
}

func parseTemplate(templateFileName string, data interface{}) (*bytes.Buffer, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf, nil
}

type loginAuth struct {
	identity, username, password string
	host                         string
}

func LoginAuth(identity, username, password, host string) smtp.Auth {
	return &loginAuth{identity, username, password, host}
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "LOGIN", resp, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)

	if more {
		if command == "username" {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if command == "password" {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}
