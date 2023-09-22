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

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"

	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
)

func TestLoginWithEmptyCredentials(t *testing.T) {
	testInputs := []string{
		`{"u":"","p":""}`,
		`{"u":"foo","p":""}`,
		`{"u":"","p":"foo"}`,
	}

	for _, testInput := range testInputs {
		w := httptest.NewRecorder()
		body := strings.NewReader(testInput)
		r, err := http.NewRequest(http.MethodPost, "login", body)
		if err != nil {
			t.Error("Error creating new request")
		}
		LoginHandler(nil, config.Config{})(w, r)

		expected := `{"alerts":[{"text":"username and password are required","level":"error"}]}` + "\n"
		if w.Body.String() != expected {
			t.Error("Expected body", expected, "got", w.Body.String())
		}
	}
}

func TestVerifyUrlOnWhiteList(t *testing.T) {
	type TestResult struct {
		Whitelist      []string
		ExpectedResult bool
	}

	completeTestResults := struct {
		Results []TestResult
	}{}

	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{}, ExpectedResult: false})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{""}, ExpectedResult: false})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"*"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong"}, ExpectedResult: false})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.right.com"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"*.right.com"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "test.right.com"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "*.right.*"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "*right*"}, ExpectedResult: true})
	completeTestResults.Results = append(completeTestResults.Results, TestResult{Whitelist: []string{"test.wrong", "*right"}, ExpectedResult: false})

	url := "https://test.right.com/other/parts"

	for _, result := range completeTestResults.Results {
		if matched, _ := VerifyUrlOnWhiteList(url, result.Whitelist); matched != result.ExpectedResult {
			t.Errorf("for whitelist: %v, expected: %v, actual: %v", result.Whitelist, result.ExpectedResult, matched)
		}
	}
}

func TestTemplateRender(t *testing.T) {
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

	f := emailFormatter{
		From:         from,
		To:           to,
		Token:        "test",
		InstanceName: "TO API Unit Tests",
		ResetURL:     "https://example.test/#!/user",
	}

	var tmpl bytes.Buffer
	if err := resetPasswordEmailTemplate.Execute(&tmpl, &f); err != nil {
		t.Fatalf("Failed to render email template: %v", err)
	}
	if tmpl.Len() <= 0 {
		t.Fatalf("Template buffer empty after execution")
	}
	t.Logf("%s", tmpl.String())
}
