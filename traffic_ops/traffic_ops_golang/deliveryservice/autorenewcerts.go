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
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
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

const emailTemplateFile = "/opt/traffic_ops/app/templates/send_mail/autorenewcerts_mail.html"

func RenewCertificates(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, errors.New("the Riak service is unavailable"), errors.New("getting SSL keys from Riak by xml id: Riak is not configured"))
		return
	}

	rows, err := inf.Tx.Tx.Query(`SELECT xml_id, ssl_key_version FROM deliveryservice`)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, err)
		return
	}
	defer rows.Close()

	keysFound := ExpirationSummary{}
	for rows.Next() {
		ds := DsKey{}
		err := rows.Scan(&ds.XmlId, &ds.Version)
		if err != nil {
			log.Errorf("getting delivery services: %v", err)
			continue
		}
		if ds.Version.Valid && ds.Version.Int64 != 0 {
			continue
		}

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

		expiration, err := parseExpirationFromCert([]byte(keyObj.Certificate.Crt))
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New(ds.XmlId+": "+err.Error()))
			return
		}

		// Renew only certificates within configured limit plus 3 days
		if expiration.After(time.Now().Add(time.Hour * 24 * time.Duration(inf.Config.ConfigLetsEncrypt.RenewDaysBeforeExpiration)).Add(time.Hour * 24 * 3)) {
			continue
		}

		newVersion := util.JSONIntStr(keyObj.Version.ToInt64() + 1)

		dsExpInfo.XmlId = keyObj.DeliveryService
		dsExpInfo.Version = keyObj.Version
		dsExpInfo.Expiration = expiration
		dsExpInfo.AuthType = keyObj.AuthType

		if keyObj.AuthType == tc.LetsEncryptAuthType || (keyObj.AuthType == tc.SelfSignedCertAuthType && inf.Config.ConfigLetsEncrypt.ConvertSelfSigned) {
			req := tc.DeliveryServiceLetsEncryptSSLKeysReq{
				DeliveryServiceSSLKeysReq: tc.DeliveryServiceSSLKeysReq{
					HostName:        &keyObj.Hostname,
					DeliveryService: &keyObj.DeliveryService,
					CDN:             &keyObj.CDN,
					Version:         &newVersion,
				},
			}
			ctx, _ := context.WithTimeout(r.Context(), LetsEncryptTimeout)

			if error := GetLetsEncryptCertificates(inf.Config, req, ctx, inf.User); error != nil {
				dsExpInfo.Error = error
			}
			keysFound.LetsEncryptExpirations = append(keysFound.LetsEncryptExpirations, dsExpInfo)

		} else if keyObj.AuthType == tc.SelfSignedCertAuthType {
			keysFound.SelfSignedExpirations = append(keysFound.SelfSignedExpirations, dsExpInfo)
		} else {
			keysFound.OtherExpirations = append(keysFound.OtherExpirations, dsExpInfo)
		}

	}

	if inf.Config.SMTP.Enabled && inf.Config.ConfigLetsEncrypt.SendExpEmail {
		errCode, userErr, sysErr := AlertExpiringCerts(keysFound, *inf.Config)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
			return
		}

	}

	api.WriteResp(w, r, keysFound)

}

func AlertExpiringCerts(certsFound ExpirationSummary, config config.Config) (int, error, error) {
	header := "From: " + config.ConfigTO.EmailFrom.String() + "\r\n" +
		"To: " + config.ConfigLetsEncrypt.Email + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"Subject: Certificate Expiration Summary\r\n\r\n"

	return api.SendEmailFromTemplate(config, header, certsFound, emailTemplateFile, config.ConfigLetsEncrypt.Email)
}
