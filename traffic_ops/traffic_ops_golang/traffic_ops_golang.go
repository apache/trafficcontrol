package main

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
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const Version = "0.1"

const DefaultConfigPath = "/opt/traffic_ops/traffic_ops_golang.config"

const OldConfig = true
const OldConfigCDNConfPath = "/opt/traffic_ops/app/conf/cdn.conf"
const OldConfigDBConfPath = "/opt/traffic_ops/app/conf/production/database.conf"

func main() {
	configFileName := flag.String("cfg", "", "The config file path")
	oldConfig := flag.Bool("oldcfg", true, "Whether to look for old Perl Traffic Ops config files")
	flag.Parse()
	if *configFileName == "" {
		*configFileName = DefaultConfigPath
	}

	cfg := Config{}
	err := error(nil)
	if !*oldConfig {
		if cfg, err = LoadConfig(*configFileName); err != nil {
			fmt.Println("Error loading config '" + *configFileName + "': " + err.Error())
			return
		}
	} else {
		if cfg, err = GetPerlConfigs(OldConfigCDNConfPath, OldConfigDBConfPath); err != nil {
			fmt.Println("Error loading old configs '" + OldConfigCDNConfPath + "' and '" + OldConfigDBConfPath + "': " + err.Error())
			return
		}
	}

	if err := log.InitCfg(cfg); err != nil {
		fmt.Printf("Error initializing loggers: %v\n", err)
		return
	}

	log.Infof(`Using Config values:
		Port:                 %s
		Db Server:            %s
		Db User:              %s
		Db Name:              %s
		Db Ssl:               %t
		Max Db Connections:   %d
		TO URL:               %s
		Insecure:             %t
		Cert Path:            %s
		Key Path:             %s
		Proxy Timeout:        %v
		Proxy KeepAlive:      %v
		Proxy tls handshake:  %v
		Proxy header timeout: %v
		Read Timeout:         %v
		Read Header Timeout:  %v
		Write Timeout:        %v
		Idle Timeout:         %v
		Error Log:            %s
		Warn Log:             %s
		Info Log:             %s
		Debug Log:            %s
		Event Log:            %s`, cfg.HTTPPort, cfg.DBServer, cfg.DBUser, cfg.DBDB, cfg.DBSSL, cfg.MaxDBConnections, cfg.TOURLStr, cfg.Insecure, cfg.CertPath, cfg.KeyPath, time.Duration(cfg.ProxyTimeout)*time.Second, time.Duration(cfg.ProxyKeepAlive)*time.Second, time.Duration(cfg.ProxyTLSTimeout)*time.Second, time.Duration(cfg.ProxyReadHeaderTimeout)*time.Second, time.Duration(cfg.ReadTimeout)*time.Second, time.Duration(cfg.ReadHeaderTimeout)*time.Second, time.Duration(cfg.WriteTimeout)*time.Second, time.Duration(cfg.IdleTimeout)*time.Second, cfg.LogLocationError, cfg.LogLocationWarning, cfg.LogLocationInfo, cfg.LogLocationDebug, cfg.LogLocationEvent)

	sslStr := "require"
	if !cfg.DBSSL {
		sslStr = "disable"
	}

	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.DBUser, cfg.DBPass, cfg.DBServer, cfg.DBDB, sslStr))
	if err != nil {
		log.Errorf("opening database: %v\n", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.MaxDBConnections)

	if err := RegisterRoutes(ServerData{DB: db, Config: cfg}); err != nil {
		log.Errorf("registering routes: %v\n", err)
		return
	}

	log.Infof("Listening on " + cfg.HTTPPort)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		TLSConfig:         &tls.Config{InsecureSkipVerify: cfg.Insecure},
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeout) * time.Second,
	}
	server.SetKeepAlivesEnabled(true)

	log.Debugf("our server struct: %++v \n", server)
	if err := server.ListenAndServeTLS(cfg.CertPath, cfg.KeyPath); err != nil {
		log.Errorf("stopping server: %v\n", err)
		return
	}
}
