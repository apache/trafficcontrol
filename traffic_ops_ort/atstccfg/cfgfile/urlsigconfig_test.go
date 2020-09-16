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
	"testing"
)

func TestGetDSFromURLSigConfigFileName(t *testing.T) {
	expecteds := map[string]string{
		"url_sig_foo.config":                        "foo",
		"url_sig_.config":                           "",
		"url_sig.config":                            "",
		"url_sig_foo.conf":                          "",
		"url_sig_foo.confi":                         "",
		"url_sig_foo_bar_baz.config":                "foo_bar_baz",
		"url_sig_url_sig_foo_bar_baz.config.config": "url_sig_foo_bar_baz.config",
	}

	for fileName, expected := range expecteds {
		actual := GetDSFromURLSigConfigFileName(fileName)
		if expected != actual {
			t.Errorf("GetDSFromURLSigConfigFileName('%v') expected '%v' actual '%v'\n", fileName, expected, actual)
		}
	}
}
