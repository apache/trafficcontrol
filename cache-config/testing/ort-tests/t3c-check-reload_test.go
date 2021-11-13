package orttest

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
	"encoding/json"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
)

func TestCheckReload(t *testing.T) {
	type ChangedCfg struct {
		ChangedFiles     string `json:"changed_files"`
		InstalledPlugins string `json:"installed_plugins"`
	}

	type argsResults struct {
		configs     ChangedCfg
		mode        string
		expected    string
		expectedErr bool
	}

	argsExpected := []argsResults{
		{
			configs: ChangedCfg{
				ChangedFiles:     "/etc/trafficserver/remap.config,/etc/trafficserver/parent.config",
				InstalledPlugins: "",
			},
			expected: "reload",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/etc/trafficserver/anything.foo",
				InstalledPlugins: "",
			},
			expected: "reload",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/opt/trafficserver/etc/trafficserver/anything.foo",
				InstalledPlugins: "",
			},
			expected: "reload",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/foo/bar/hdr_rw_foo.config",
				InstalledPlugins: "",
			},
			expected: "reload",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/foo/bar/uri_signing_dsname.config",
				InstalledPlugins: "",
			},
			expected: "reload",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/foo/bar/url_sig_dsname.config,foo",
				InstalledPlugins: "",
			},
			expected: "reload",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "plugin.config,foo",
				InstalledPlugins: "",
			},
			expected: "restart",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/etc/trafficserver/anything.foo",
				InstalledPlugins: "anything",
			},
			expected: "restart",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "",
				InstalledPlugins: "anything",
			},
			expected: "restart",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "",
				InstalledPlugins: "anything,anythingelse",
			},
			expected: "restart",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/foo/bar/ssl_multicert.config",
				InstalledPlugins: "",
			},
			expected: "reload",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "foo",
				InstalledPlugins: "",
			},
			expected: "",
		},
		{
			configs: ChangedCfg{
				ChangedFiles:     "/foo/bar/baz.config",
				InstalledPlugins: "",
			},
			expected: "",
		},
	}

	for _, ae := range argsExpected {
		config, err := json.Marshal(ae.configs)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		out, code := t3cCheckReload(config)
		out = strings.TrimSpace(out)
		if !ae.expectedErr && code != 0 {
			t.Errorf("expected configs %+v packages %+v would not error, actual: code %v output '%v'",
				ae.configs.ChangedFiles, ae.configs.InstalledPlugins, code, out)
			continue
		} else if ae.expectedErr && code == 0 {
			t.Errorf("expected configs %+v packages %+v would error, actual: no error",
				ae.configs.ChangedFiles, ae.configs.InstalledPlugins)
			continue
		}
		if out != ae.expected {
			t.Errorf("expected configs %+v packages %+v would need '%v', actual: '%v'",
				ae.configs.ChangedFiles, ae.configs.InstalledPlugins, ae.expected, out)
		}
	}
}

func t3cCheckReload(configs []byte) (string, int) {
	args := []string{
		"check", "reload",
	}
	stdOut, _, exitCode := t3cutil.DoInput(configs, "t3c", args...)
	return string(stdOut), exitCode
}
