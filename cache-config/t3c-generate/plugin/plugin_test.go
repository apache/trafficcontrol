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
	"testing"

	"github.com/apache/trafficcontrol/v8/cache-config/t3c-generate/config"
	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
)

func TestPlugin(t *testing.T) {
	AddPlugin(10000, Funcs{
		modifyFiles: func(d ModifyFilesData) []t3cutil.ATSConfigFile {
			if d.TOData.Server == nil || d.TOData.Server.HostName == "" || d.TOData.Server.HostName != "testplugin" {
				return d.Files
			}
			fi := t3cutil.ATSConfigFile{}
			fi.Text = "testfile\n"
			fi.ContentType = "text/plain"
			fi.LineComment = ""
			fi.Name = "testfile.txt"
			fi.Path = "/opt/trafficserver/etc/trafficserver/"
			d.Files = append(d.Files, fi)
			return d.Files
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

	modifyFilesData := ModifyFilesData{
		Cfg:    cfg,
		Files:  []t3cutil.ATSConfigFile{},
		TOData: &t3cutil.ConfigData{},
	}

	newFiles := plugins.ModifyFiles(modifyFilesData)
	if len(newFiles) > 0 {
		t.Error("Expected server '' to be unhandled by a plugin, actual: handled")
	}

	if modifyFilesData.TOData.Server == nil {
		modifyFilesData.TOData.Server = &atscfg.Server{}
	}
	modifyFilesData.TOData.Server.HostName = "testplugin"
	newFiles = plugins.ModifyFiles(modifyFilesData)
	if len(newFiles) == 0 {
		t.Error("Expected server 'testplugin' to be handled by plugin, actual: unhandled")
	}
	fi := newFiles[0]
	if fi.Text != "testfile\n" {
		t.Errorf(`Expected plugin text 'testfile\n', actual %v`, fi.Text)
	}
}
