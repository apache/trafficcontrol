package plugin

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
	"net/url"
	"testing"

	"github.com/apache/trafficcontrol/traffic_ops/ort/atstccfg/config"
)

func TestPlugin(t *testing.T) {
	handledPath := "/__test"
	AddPlugin(10000, Funcs{
		onRequest: func(d OnRequestData) IsRequestHandled {
			if d.Cfg.TOURL.Path != handledPath {
				return RequestUnhandled
			}
			// a real plugin would print and exit here
			return RequestHandled
		},
	})

	pluginNames := List()

	expectedPluginName := `plugin_test` // plugin names are the file name. This must be the name of this file.

	found := false
	for _, pluginName := range pluginNames {
		if pluginName == expectedPluginName {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected plugin named '%v', actual %+v", expectedPluginName, pluginNames)
	}

	cfg := config.Cfg{}
	plugins := Get(cfg)

	tcCfg := config.TCCfg{Cfg: cfg}
	onReqData := OnRequestData{Cfg: tcCfg}

	nonPluginURL, err := url.Parse("https://example.net/should-not-be-a-known-config")
	if err != nil {
		t.Fatal(err)
	}
	onReqData.Cfg.TOURL = nonPluginURL
	handled := plugins.OnRequest(onReqData)
	if handled {
		t.Errorf("Expected url %s to be unhandled by a plugin, actual: handled", nonPluginURL)
	}

	pluginURL, err := url.Parse("https://example.net" + handledPath)
	if err != nil {
		t.Fatal(err)
	}
	onReqData.Cfg.TOURL = pluginURL
	handled = plugins.OnRequest(onReqData)
	if !handled {
		t.Errorf("Expected url %s to be handled by plugin, actual: handled", pluginURL)
	}
}
