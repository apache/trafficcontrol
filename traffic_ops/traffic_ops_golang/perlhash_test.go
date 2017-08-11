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
	"reflect"
	"testing"
)

func TestGetStr(t *testing.T) {
	type Strtst struct {
		Input       string
		ExpectedStr string
		ExpectedS   string
	}
	strs := []Strtst{
		Strtst{`'foo'asdf`, `foo`, `asdf`},
		Strtst{`''asdf`, ``, `asdf`},
		Strtst{`'1234'asdf`, `1234`, `asdf`},
		Strtst{`'\\'asdf`, `\`, `asdf`},
		Strtst{`'\''asdf`, `'`, `asdf`},
		Strtst{`'\''asdfqhrewoipjasf`, `'`, `asdfqhrewoipjasf`},
	}

	for _, st := range strs {
		str, s, err := getStr(st.Input)
		if err != nil {
			t.Errorf("expected nil err actual %v", err)
		}
		if str != st.ExpectedStr {
			t.Errorf("expected str '%v' actual '%v'", st.ExpectedStr, str)
		}
		if s != st.ExpectedS {
			t.Errorf("expected s '%v' actual '%v'", st.ExpectedS, s)
		}
	}
}

func TestGetNum(t *testing.T) {
	type Tst struct {
		Input       string
		ExpectedNum float64
		ExpectedS   string
	}
	strs := []Tst{
		Tst{`42asdf`, 42, `asdf`},
		Tst{`42 asdf`, 42, ` asdf`},
		Tst{`42.1asdf`, 42.1, `asdf`},
		Tst{`42.1 asdf`, 42.1, ` asdf`},
		Tst{`42. asdf`, 42, ` asdf`},
		Tst{`.42 asdf`, 0.42, ` asdf`},
		Tst{`0.42 asdf`, 0.42, ` asdf`},
		Tst{`0.42 {asdf`, 0.42, ` {asdf`},
		Tst{`9{asdf`, 9, `{asdf`},
		Tst{`42.1`, 42.1, ``},
		Tst{`9`, 9, ``},
	}

	for _, st := range strs {
		num, s, err := getNum(st.Input)
		if err != nil {
			t.Errorf("expected nil err actual %v", err)
		}
		if num != st.ExpectedNum {
			t.Errorf("expected num '%v' actual '%v'", st.ExpectedNum, num)
		}
		if s != st.ExpectedS {
			t.Errorf("expected s '%v' actual '%v'", st.ExpectedS, s)
		}
	}
}

func TestGetArr(t *testing.T) {
	type Tst struct {
		Input       string
		ExpectedVal []interface{}
		ExpectedS   string
	}
	tsts := []Tst{
		Tst{`[42]asdf`, []interface{}{42.0}, `asdf`},
		Tst{`[42, 'foo']],1]asdf`, []interface{}{42.0, `foo`}, `],1]asdf`},
		Tst{`[
 			'https://[::]:443?cert=/etc/pki/tls/certs/localhost.crt&key=/etc/pki/tls/private/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED'
 		],`, []interface{}{`https://[::]:443?cert=/etc/pki/tls/certs/localhost.crt&key=/etc/pki/tls/private/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED`}, `,`},
		Tst{`[ 'foobar' ],`, []interface{}{"foobar"}, `,`},
		Tst{`[ 'foo' , 'bar', 'baz'   ],`, []interface{}{`foo`, `bar`, `baz`}, `,`},
	}

	for _, st := range tsts {
		val, s, err := getArr(st.Input)
		if err != nil {
			t.Errorf("expected nil err actual %v", err)
		}
		if !reflect.DeepEqual(val, st.ExpectedVal) {
			t.Errorf("expected arr '%v' actual '%v'", st.ExpectedVal, val)
		}
		if s != st.ExpectedS {
			t.Errorf("expected s '%v' actual '%v'", st.ExpectedS, s)
		}
	}
}

func TestGetObj(t *testing.T) {
	type Tst struct {
		Input       string
		ExpectedVal map[string]interface{}
		ExpectedS   string
	}
	tsts := []Tst{
		Tst{`{ a => 'b', c => [42.0, 'd']},asdf`,
			map[string]interface{}{
				`a`: `b`,
				`c`: []interface{}{42.0, `d`},
			},
			`,asdf`},
	}

	for _, st := range tsts {
		val, s, err := getObj(st.Input)
		if err != nil {
			t.Errorf("expected nil err actual %v", err)
		}
		if !reflect.DeepEqual(val, st.ExpectedVal) {
			t.Errorf("expected obj '%+v' actual '%+v'", st.ExpectedVal, val)
		}
		if s != st.ExpectedS {
			t.Errorf("expected s '%v' actual '%v'", st.ExpectedS, s)
		}
	}
}

func TestParsePerlObj(t *testing.T) {
	input := `
{
	hypnotoad => {
		listen => [
			'https://[::]:443?cert=/etc/pki/tls/certs/localhost.crt&key=/etc/pki/tls/private/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED'
		],
		user     => 'trafops',
		group    => 'trafops',
		heartbeat_timeout => 20,
		pid_file => '/var/run/traffic_ops.pid',
		workers  => 96
	},
	cors => {
		access_control_allow_origin => '*'
	},
	to => {
		base_url   => 'http://localhost:3000',                    # this is where traffic ops app resides
		email_from => 'no-reply@traffic-ops-domain.com'           # traffic ops email address
	},
	portal => {
		base_url   => 'http://localhost:8080',                    # this is where the traffic portal resides (a javascript client that consumes the TO API)
		email_from => 'no-reply@traffic-portal-domain.com'        # traffic portal email address
	},

	# 1st secret is used to generate new signatures. Older one kept around for existing signed cookies.
		#  Remove old one(s) when ready to invalidate old cookies.
		secrets => [ 'walrus' ],
	geniso  => {
		iso_root_path => '/opt/traffic_ops/app/public',          # the location where the iso files will be written
	},
	inactivity_timeout => 60
};
`
	expected := map[string]interface{}{
		"hypnotoad": map[string]interface{}{
			"listen":            []interface{}{`https://[::]:443?cert=/etc/pki/tls/certs/localhost.crt&key=/etc/pki/tls/private/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED`},
			"user":              `trafops`,
			"group":             `trafops`,
			"heartbeat_timeout": 20.0,
			"pid_file":          `/var/run/traffic_ops.pid`,
			"workers":           96.0,
		},
		"cors": map[string]interface{}{
			"access_control_allow_origin": `*`,
		},
		"to": map[string]interface{}{
			"base_url":   `http://localhost:3000`,
			"email_from": `no-reply@traffic-ops-domain.com`,
		},
		"portal": map[string]interface{}{
			"base_url":   `http://localhost:8080`,
			"email_from": `no-reply@traffic-portal-domain.com`,
		},
		"secrets": []interface{}{`walrus`},
		"geniso": map[string]interface{}{
			"iso_root_path": `/opt/traffic_ops/app/public`,
		},
		"inactivity_timeout": 60.0,
	}

	val, err := ParsePerlObj(input)
	if err != nil {
		t.Errorf("expected nil err actual %v", err)
	}
	if !reflect.DeepEqual(val, expected) {
		t.Errorf("expected '%+v' actual '%+v'", expected, val)
	}
}
