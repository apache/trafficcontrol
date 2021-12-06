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
	"fmt"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
)

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

func testCheckReload(t *testing.T, ae argsResults) {
	config, err := json.Marshal(ae.configs)
	if err != nil {
		t.Errorf("failed to encode configs: %v", err)
	}
	out, code := t3cCheckReload(config)
	out = strings.TrimSpace(out)
	if !ae.expectedErr && code != 0 {
		t.Fatalf("expected non-error exit code, actual: %d - output: %s", code, out)
	}
	if ae.expectedErr && code == 0 {
		t.Fatal("expected check-reload to exit with an error, actual: no error")
	}
	if out != ae.expected {
		t.Errorf("expected required action '%s', actual: '%s'", ae.expected, out)
	}
}

func TestCheckReload(t *testing.T) {

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
		testName := fmt.Sprintf("testing check-reload with changed files %+v and installed plugins %+v", ae.configs.ChangedFiles, ae.configs.InstalledPlugins)
		t.Run(testName, func(t *testing.T) { testCheckReload(t, ae) })

	}
}

func t3cCheckReload(configs []byte) (string, int) {
	args := []string{
		"check", "reload",
	}
	stdOut, _, exitCode := t3cutil.DoInput(configs, "t3c", args...)
	return string(stdOut), exitCode
}
