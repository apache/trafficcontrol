// Package config provides tools to load and validate configuration data for the
// t3c tests.
package config

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	log "github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/kelseyhightower/envconfig"
)

// Config reflects the structure of the test-to-api.conf file. It also
// implements github.com/apache/trafficcontrol/v8/lib/go-log.Config for logging.
type Config struct {
	TrafficOps   TrafficOps   `json:"trafficOps"`
	TrafficOpsDB TrafficOpsDB `json:"trafficOpsDB"`
	Default      Default      `json:"default"`
	UseIMS       bool         `json:"use_ims"`
}

// TrafficOps is the section of a Config dealing with Traffic Ops-related
// information.
type TrafficOps struct {
	// URL is the URL of a Traffic Ops instance being used for the tests.
	URL string `json:"URL" envconfig:"TO_URL"`

	// UserPassword is the password of *all* users in the 'Users' property.
	UserPassword string `json:"password" envconfig:"TO_USER_PASSWORD"`

	// Users are the Traffic Ops users to be used in testing.
	Users Users `json:"users"`

	// Insecure instructs tests whether or not to skip certificate verification.
	Insecure bool `json:"sslInsecure" envconfig:"SSL_INSECURE"`
}

// Users structures are the "users" section of the "trafficOps" section of the
// testing configuration file, and are a collection of the usernames of  Traffic
// Ops users that are used in testing.
type Users struct {

	// Disallowed is the username of a "disallowed" Traffic Ops user.
	//
	// Deprecated: This is unused in t3c tests, and may be removed in the
	// future.
	Disallowed string `json:"disallowed" envconfig:"TO_USER_DISALLOWED"`

	// ReadOnly is the username of a Traffic Ops user with "read-only"
	// Permissions.
	//
	// Deprecated: This is unused in t3c tests, and may be removed in the
	// future.
	ReadOnly string `json:"readOnly" envconfig:"TO_USER_READ_ONLY"`

	// Operations is the username of a Traffic Ops user with "operations"-level
	// Permissions.
	//
	// Deprecated: This is unused in t3c tests, and may be removed in the
	// future.
	Operations string `json:"operations" envconfig:"TO_USER_OPERATIONS"`

	// Admin is the username of a Traffic Ops user with the special "admin"
	// Role.
	Admin string `json:"admin" envconfig:"TO_USER_ADMIN"`

	// Portal is the username of a Traffic Ops user with "portal"-level
	// Permissions.
	//
	// Deprecated: This is unused in t3c tests, and may be removed in the
	// future.
	Portal string `json:"portal" envconfig:"TO_USER_PORTAL"`

	// Federation is the username of a Traffic Ops user with
	// "federation"-level Permissions.
	//
	// Deprecated: This is unused in t3c tests, and may be removed in the
	// future.
	Federation string `json:"federation" envconfig:"TO_USER_FEDERATION"`

	// Extension is the username of a Traffic Ops user allowed to manipulate
	// extensions (i.e. creating new serverchecks).
	//
	// These tests currently use Traffic Ops API version 3.0, so the only
	// username that can possibly work is literally "extension". This property
	// MUST be configured to be "extension", or the tests will erroneously fail,
	// no matter what "Priv Level" and/or Permissions you give the user!
	Extension string `json:"extension" envconfig:"TO_USER_EXTENSION"`
}

// TrafficOpsDB is the section of a Config dealing with Traffic Ops DB-related
// information.
type TrafficOpsDB struct {
	// Name is the name of a PostgreSQL database used by the testing Traffic Ops
	// instance.
	Name string `json:"dbname" envconfig:"TODB_NAME"`

	// Hostname is the network hostname where the Traffic Ops database is
	// running.
	Hostname string `json:"hostname" envconfig:"TODB_HOSTNAME"`

	// User is the name of a PostgreSQL user/role that has permissions to
	// manipulate the database identified in Name.
	User string `json:"user" envconfig:"TODB_USER"`

	// Password is the password for the PostgreSQL user/role given in User.
	Password string `json:"password" envconfig:"TODB_PASSWORD"`

	// Port is the port on which PostgreSQL listens for connections to the
	// database used by the testing Traffic Ops instance.
	Port string `json:"port" envconfig:"TODB_PORT"`

	// DBType is a Go database/sql driver name to use when connecting to the
	// Traffic Ops testing instance's database. This MUST be "Pg" (or
	// omitted/blank/null, which uses the default value of "Pg").
	//
	// Deprecated: Since Traffic Ops 3.0, only PostgreSQL databases are
	// supported, so this field has no purpose and will probably be removed at
	// some point.
	DBType string `json:"type" envconfig:"TODB_TYPE"`

	// SSL instructs the database driver that the PostgreSQL instance used by
	// the tests requires SSL-secured connections.
	SSL bool `json:"ssl" envconfig:"TODB_SSL"`

	// Description is a textual description of the database.
	//
	// Deprecated: This unused field serves no purpose, and will likely be
	// removed in the future.
	Description string `json:"description" envconfig:"TODB_DESCRIPTION"`
}

// Default represents the default values of a set of options that can be
// overridden by command-line options.
type Default struct {
	Session Session   `json:"session"`
	Log     Locations `json:"logLocations"`
	// IncludeSystemTests has no effect or known purpose.
	//
	// Deprecated: This field has no effect or known purpose.
	IncludeSystemTests bool `json:"includeSystemTests"`
}

// Session contains default configuration options for authenticated sessions
// with the Traffic Ops API.
type Session struct {
	TimeoutInSecs int `json:"timeoutInSecs" envconfig:"SESSION_TIMEOUT_IN_SECS"`
}

// Locations is a set of logging locations as defined by the
// github.com/apache/trafficcontrol/v8/lib/go-log package.
type Locations struct {
	Debug   string `json:"debug"`
	Event   string `json:"event"`
	Error   string `json:"error"`
	Info    string `json:"info"`
	Warning string `json:"warning"`
}

// LoadConfig reads the given configuration file into a Config.
func LoadConfig(confPath string) (Config, error) {
	var cfg Config

	confBytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read CDN configuration: %w", err)
	}

	err = json.Unmarshal(confBytes, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse configuration from '%s': %w", confPath, err)
	}

	if err := validate(confPath, cfg); err != nil {
		return cfg, fmt.Errorf("failed to validate configuration: %w", err)
	}

	if err := envconfig.Process("traffic-ops-client-tests", &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse configuration from environment: %w", err)
	}

	return cfg, nil
}

type multiError []error

func (me multiError) Error() string {
	var sb strings.Builder
	for _, e := range me {
		fmt.Fprintln(&sb, e)
	}
	return sb.String()
}

// validate all required fields in the config.
func validate(confPath string, config Config) error {
	var errs multiError

	var f string
	f = "TrafficOps"
	toTag, ok := getStructTag(config, f)
	if !ok {
		errs = append(errs, fmt.Errorf("'%s' must be configured in %s", toTag, confPath))
	}

	if config.TrafficOps.URL == "" {
		f = "URL"
		tag, ok := getStructTag(config.TrafficOps, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	if config.TrafficOps.Users.Disallowed == "" {
		f = "Disallowed"
		tag, ok := getStructTag(config.TrafficOps.Users, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	if config.TrafficOps.Users.ReadOnly == "" {
		f = "ReadOnly"
		tag, ok := getStructTag(config.TrafficOps.Users, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	if config.TrafficOps.Users.Operations == "" {
		f = "Operations"
		tag, ok := getStructTag(config.TrafficOps.Users, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	if config.TrafficOps.Users.Admin == "" {
		f = "Admin"
		tag, ok := getStructTag(config.TrafficOps.Users, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	if config.TrafficOps.Users.Portal == "" {
		f = "Portal"
		tag, ok := getStructTag(config.TrafficOps.Users, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	if config.TrafficOps.Users.Federation == "" {
		f = "Federation"
		tag, ok := getStructTag(config.TrafficOps.Users, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func getStructTag(thing interface{}, fieldName string) (string, bool) {
	var tag string
	var ok bool
	t := reflect.TypeOf(thing)
	if t != nil {
		if f, ok := t.FieldByName(fieldName); ok {
			tag = f.Tag.Get("json")
			return tag, ok
		}
	}
	return tag, ok
}

// ErrorLog provides the location to which error-level messages should be
// logged.
func (c Config) ErrorLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Error)
}

// WarningLog provides the location to which warning-level messages should be
// logged.
func (c Config) WarningLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Warning)
}

// InfoLog provides the location to which info-level messages should be
// logged.
func (c Config) InfoLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Info)
}

// DebugLog provides the location to which debug-level messages should be
// logged.
func (c Config) DebugLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Debug)
}

// EventLog provides the location to which event-level messages should be
// logged.
func (c Config) EventLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Event)
}
