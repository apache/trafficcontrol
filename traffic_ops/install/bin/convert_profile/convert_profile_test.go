/**
 *
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *
 * Unit Tests for Profile Conversion Utility
 */
package main

import (
	_ "fmt"
	"testing"
)

func TestValidateParameter(t *testing.T) {
	// Test happy days validation
	profile := Profile{[]Parameter{{"trafficserver", "package", "622"}}, ProfileDesc{}}

	validators := []Parameter{{"trafficserver", "package", "622"}}

	if !ValidateParameters(&profile, validators) {
		t.Error("Failed to validate parameter")
	}

	profile.Parameters[0].Value = "6.2.2"
	if ValidateParameters(&profile, validators) {
		t.Error("Failed to catch value mismatch")
	}

	// Test name not found
	profile.Parameters[0].Value = "622"
	profile.Parameters[0].Name = "ats"
	if ValidateParameters(&profile, validators) {
		t.Error("Failed to catch missing parameter")
	}
}

func TestDeleteParameter(t *testing.T) {
	profile := Profile{[]Parameter{{"proxy.cluster.a", "records.config", "INT 1"},
		{"proxy.cluster.b", "records.config", "STRING abc"}}, ProfileDesc{}}

	deleteRules := []ConversionRule{
		{MatchParameter: Parameter{"proxy\\.cluster\\..+", "records\\.config", ".*"},
			Action: "delete"}}

	ConvertProfile(&profile, deleteRules, false)

	if len(profile.Parameters) != 0 {
		t.Error("Failed to delete parameter")
	}
}

func TestModifyConfigFile(t *testing.T) {
	profile := Profile{[]Parameter{{"LogFormat.Format", "logs_xml.config", "<abcdef> <ghi>>"},
		Parameter{"LogFormat.Name", "logs_xml.config", "custom_ats"}}, ProfileDesc{}}

	modifyRules := []ConversionRule{
		{MatchParameter: Parameter{"LogFormat\\..*", "logs_xml\\.config", ".*"},
			NewConfigFile: "logging.config"}}

	ConvertProfile(&profile, modifyRules, false)

	for _, param := range profile.Parameters {
		if param.ConfigFile != "logging.config" {
			t.Error("Failed to update config file")
		}
	}
}

func TestModifyValueForce(t *testing.T) {
	profile := Profile{[]Parameter{{"proxy.config.hostdb.timeout", "records.config", "INT 1440"}}, ProfileDesc{}}

	modifyRules := []ConversionRule{
		{MatchParameter: Parameter{"proxy\\.config\\.hostdb\\.timeout", "records\\.config", "INT 1440"},
			NewValue: "INT 86400"}}

	ConvertProfile(&profile, modifyRules, true)

	if profile.Parameters[0].Value != "INT 86400" {
		t.Error("Failed to update value")
	}
}

func TestModifyValueSkip(t *testing.T) {
	profile := Profile{[]Parameter{{"proxy.config.hostdb.timeout", "records.config", "INT 5000"}}, ProfileDesc{}}

	modifyRules := []ConversionRule{
		{MatchParameter: Parameter{"proxy\\.config\\.hostdb\\.timeout", "records\\.config", "INT 1440"},
			NewValue: "INT 86400"}}

	ConvertProfile(&profile, modifyRules, false)

	if profile.Parameters[0].Value != "INT 5000" {
		t.Error("Incorrectly updated value")
	}
}

func TestModifyNameValue(t *testing.T) {
	profile := Profile{[]Parameter{{"proxy.config.log.xml_config_file", "records.config", "logs_xml.config"}}, ProfileDesc{}}

	modifyRules := []ConversionRule{
		{MatchParameter: Parameter{"proxy\\.config\\.log\\.xml_config_file", "records\\.config", "logs_xml\\.config"},
			NewName:  "proxy.config.log.config.filename",
			NewValue: "logging.config"}}

	ConvertProfile(&profile, modifyRules, false)
	if profile.Parameters[0].Name != "proxy.config.log.config.filename" {
		t.Error("Failed to update parameter name")
	}

	if profile.Parameters[0].Value != "logging.config" {
		t.Error("Failed to update parameter value")
	}
}
