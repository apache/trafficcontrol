package login

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

import "bytes"
import "net/mail"
import "testing"

import "github.com/apache/trafficcontrol/v8/lib/go-rfc"

func TestRegistrationTemplateRender(t *testing.T) {
	to := rfc.EmailAddress{
		Address: mail.Address{
			Address: "em@i.l",
			Name:    "",
		},
	}
	from := rfc.EmailAddress{
		Address: mail.Address{
			Address: "no-reply@test.quest",
			Name:    "",
		},
	}

	f := registrationEmailFormatter{
		From:         from,
		InstanceName: "test",
		RegisterURL:  "http://localhost/#!/user",
		To:           to,
		Token:        "token",
	}

	var tmpl bytes.Buffer
	if err := registrationEmailTemplate.Execute(&tmpl, &f); err != nil {
		t.Fatalf("Failed to render email template: %v", err)
	}
	if tmpl.Len() <= 0 {
		t.Fatalf("Template buffer empty after execution")
	}
	t.Logf("%s", tmpl.String())
}
