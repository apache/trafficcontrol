package auth

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
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"

	"github.com/go-ldap/ldap/v3"
)

var defaultSet bool

func setLdapTimeoutDefault(duration time.Duration) {
	if !defaultSet {
		ldap.DefaultTimeout = duration
		defaultSet = true
	}
}

const (
	LDAPWithTLS = "ldaps://"
	LDAPNoTLS   = "ldap://"
)

func ConnectToLDAP(cfg *config.ConfigLDAP) (*ldap.Conn, error) {
	setLdapTimeoutDefault(time.Duration(cfg.LDAPTimeoutSecs) * time.Second)
	host := strings.ToLower(cfg.Host)
	var l *ldap.Conn
	var err error
	if strings.HasPrefix(host, LDAPWithTLS) {
		host = strings.TrimPrefix(host, LDAPWithTLS)
		l, err = ldap.DialTLS("tcp", host, &tls.Config{InsecureSkipVerify: cfg.Insecure, ServerName: strings.Split(host, ":")[0]})
		if err != nil {
			log.Errorln("error dialing tls")
			return nil, err
		}
	} else if strings.HasPrefix(host, LDAPNoTLS) {
		host = strings.TrimPrefix(host, LDAPNoTLS)
		l, err = ldap.Dial("tcp", host)
		if err != nil {
			log.Errorln("error dialing")
			return nil, err
		}
	}
	return l, nil
}

func LookupUserDN(username string, cfg *config.ConfigLDAP) (string, bool, error) {
	l, err := ConnectToLDAP(cfg)
	if err != nil {
		log.Errorln("unable to connect to ldap to lookup user")
		return "", false, err
	}
	defer l.Close()
	// Bind with admin user
	err = l.Bind(cfg.AdminDN, cfg.AdminPass)
	if err != nil {
		log.Errorln("error binding admin user")
		return "", false, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		cfg.SearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(cfg.SearchQuery, ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Errorln("error issuing search: ", err)
		return "", false, err
	}

	if len(sr.Entries) < 1 {
		return "", false, errors.New("User does not exist")
	} else if len(sr.Entries) > 1 {
		return "", false, errors.New("too many user entries returned")
	}
	userDN := sr.Entries[0].DN
	return userDN, true, nil
}

func AuthenticateUserDN(userDN string, password string, cfg *config.ConfigLDAP) (bool, error) {
	l, err := ConnectToLDAP(cfg)
	if err != nil {
		log.Errorln("unable to connect to ldap to authenticate user")
		return false, err
	}
	defer l.Close()

	// Bind as the user to verify their password
	err = l.Bind(userDN, password)
	if err != nil {
		log.Errorf("unable to bind as user: %+v\n", userDN)
		return false, err
	}
	return true, nil
}
