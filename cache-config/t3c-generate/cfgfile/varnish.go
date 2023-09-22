package cfgfile

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
	"errors"
	"path/filepath"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/varnishcfg"
)

// GetVarnishConfigs returns varnish configuration files
// TODO: add varnishncsa and hitch configs
func GetVarnishConfigs(toData *t3cutil.ConfigData, cfg config.Cfg) ([]t3cutil.ATSConfigFile, error) {
	vclBuilder := varnishcfg.NewVCLBuilder(toData)
	vcl, warnings, err := vclBuilder.BuildVCLFile()
	logWarnings("Generating varnish configuration files: ", warnings)

	configs := make([]t3cutil.ATSConfigFile, 0)
	// TODO: should be parameterized and generated from varnishcfg
	configs = append(configs, t3cutil.ATSConfigFile{
		Name:        "default.vcl",
		Text:        vcl,
		Path:        cfg.Dir,
		ContentType: "text/plain; charset=us-ascii",
		LineComment: "//",
		Secure:      false,
	})
	txt, hitchWarnings := varnishcfg.GetHitchConfig(toData.DeliveryServices, filepath.Join(cfg.Dir, "ssl/"))
	warnings = append(warnings, hitchWarnings...)
	logWarnings("Generating hitch configuration files: ", hitchWarnings)

	configs = append(configs, t3cutil.ATSConfigFile{
		Name:        "hitch.conf",
		Text:        txt,
		Path:        cfg.Dir,
		ContentType: "text/plain; charset=us-ascii",
		LineComment: "//",
		Secure:      false,
	})

	sslConfigs, err := GetSSLCertsAndKeyFiles(toData)
	if err != nil {
		return nil, errors.New("getting ssl key and cert config files: " + err.Error())
	}
	for i := range sslConfigs {
		// path changed manually because GetSSLCertsAndKeyFiles hardcodes the directory certs and keys are written to.
		// will be removed once GetSSLCertsAndKeyFiles uses proxy.config.ssl.server.cert.path parameter.
		sslConfigs[i].Path = filepath.Join(cfg.Dir, "ssl/")
	}
	configs = append(configs, sslConfigs...)

	return configs, nil
}
