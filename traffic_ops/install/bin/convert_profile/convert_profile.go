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
 * Convert a Traffic Control Trafficserver Mid/Edge Cache Profile to a newer version
 *
 */
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type InputConfigParams struct {
	InProfile string
	OutFile   string
	Rules     string
	Force     bool
}

// TrafficOps Profile Parsing
type Profile struct {
	Parameters  []Parameter `json:"parameters" yaml:"parameters"`
	Description ProfileDesc `json:"profile" yaml:"profile"`
}

type Parameter struct {
	Name       string `json:"name" yaml:"name"`
	ConfigFile string `json:"config_file" yaml:"config_file"`
	Value      string `json:"value" yaml:"value"`
}

type ProfileDesc struct {
	Description string `json:"description" yaml:"description"`
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	Cdn         string `json:"cdn" yaml:"cdn"`
}

// ConversionPolicy Parsing
type ConversionPolicy struct {
	ValidateParameters []Parameter      `json:"validate_parameters" yaml:"validate_parameters"`
	ReplaceName        ReplaceRule      `json:"replace_name" yaml:"replace_name"`
	ReplaceDescription ReplaceRule      `json:"replace_description" yaml:"replace_description"`
	ConversionRules    []ConversionRule `json:"conversion_actions" yaml:"conversion_actions"`
	AddParameters      []Parameter      `json:"add_parameters" yaml:"add_parameters"`
}

type ReplaceRule struct {
	Old string `json:"old" yaml:"old"`
	New string `json:"new" yaml:"new"`
}

type ConversionRule struct {
	MatchParameter Parameter `json:"match_parameter" yaml:"match_parameter"`
	NewName        string    `json:"new_name" yaml:"new_name"`
	NewConfigFile  string    `json:"new_config_file" yaml:"new_config_file"`
	NewValue       string    `json:"new_value" yaml:"new_value"`
	Action         string    `json:"action" yaml:"action"`
}

func formatParam(p Parameter) string {
	return fmt.Sprintf(`{"%s", "%s", "%s"}`, p.Name, p.ConfigFile, p.Value)
}

// Applies the rule represented by cr to the input parameter.
//
//	Any non-empty string value will be replaced in the input with its new value
//	Additionally an action may indicate a non-replacement operation, such as delete
func (cr ConversionRule) Apply(param Parameter) (Parameter, bool) {
	inParam := formatParam(param)

	if cr.Action == "delete" {
		fmt.Fprintf(os.Stdout, "Deleting parameter %s\n", inParam)
		return param, false
	} else if cr.Action != "" {
		fmt.Fprintf(os.Stderr, "[WARNING] Unknown action %s, skipping action\n", cr.Action)
	}

	if cr.NewName != "" {
		param.Name = cr.NewName
	}

	if cr.NewConfigFile != "" {
		param.ConfigFile = cr.NewConfigFile
	}

	if cr.NewValue != "" {
		param.Value = cr.NewValue
	}
	fmt.Fprintf(os.Stdout, "Updating parameter %s to %s\n", inParam, formatParam(param))

	return param, true
}

func parseArgs() InputConfigParams {
	inputConfig := InputConfigParams{}
	flag.StringVar(&inputConfig.InProfile, "input_profile", "", "Path of input profile")
	flag.StringVar(&inputConfig.Rules, "rules", "", "Path to conversion rules")
	flag.StringVar(&inputConfig.OutFile, "out", "", "Path to write output file to. If not given, uses stdout")
	flag.BoolVar(&inputConfig.Force, "force", false, "Ignore parameter value, making all recommended changes")
	flag.Parse()

	if inputConfig.InProfile == "" {
		fmt.Fprintf(os.Stderr, "[ERROR] Missing required -input_profile parameter\n")
		os.Exit(1)
	}

	if inputConfig.Rules == "" {
		fmt.Fprintf(os.Stderr, "[ERROR] Missing required -rules parameter\n")
		os.Exit(1)
	}

	return inputConfig
}

func readFile(inFile string) []byte {
	file, err := ioutil.ReadFile(inFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Cannot open input file: %s\n", inFile)
		panic(err)
	}

	return file
}

func parseInputProfile(inFile string) *Profile {
	var pt Profile
	err := yaml.Unmarshal(readFile(inFile), &pt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Cannot parse input profile\n")
		panic(err)
	}

	return &pt
}

func parseInputRules(inFile string) *ConversionPolicy {
	var cp ConversionPolicy
	err := yaml.Unmarshal(readFile(inFile), &cp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Cannot parse conversion rules\n")
		panic(err)
	}

	return &cp
}

// ValidateParameters will verify that all parameters in paramsToValidate appear
//
//	exactly in profile
func ValidateParameters(profile *Profile,
	paramsToValidate []Parameter) bool {

	for _, validate := range paramsToValidate {
		found := false

		for _, param := range profile.Parameters {
			if param.Name == validate.Name && param.ConfigFile == validate.ConfigFile {
				found = true

				if param.Value != validate.Value {
					fmt.Fprintf(os.Stderr, "[ERROR] Parameter %s does not match value\n", param.Name)
					fmt.Fprintf(os.Stderr, "[ERROR]   Actual Value: %s Expected Value: %s\n", param.Value, validate.Value)
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// There is currently no check if the name, config_file or value is already present
func AddParameters(profile *Profile, params []Parameter) {
	filteredParams := profile.Parameters
	for _, param := range params {
		// fmt.Fprintf(os.Stdout, "Adding parameter %s\n", add.Name)
		fmt.Fprintf(os.Stdout, "Adding parameter %s\n", formatParam(param))
		filteredParams = append(filteredParams, param)
	}
	profile.Parameters = filteredParams
}

// ConvertProfile will modify paramaters as described by matching entries in conversionActions
// If ignoreValue is set to true, the Value field in matcher will be ignored, effectively matching
// all values
func ConvertProfile(profile *Profile,
	rules []ConversionRule,
	ignoreValue bool) {
	filteredParams := profile.Parameters[:0]

	for _, param := range profile.Parameters {

		matched := false
		//fmt.Fprintf(os.Stdout, "ConvertProfile: %s\n", param.Name)
		for _, rule := range rules {
			if paramsMatch(rule.MatchParameter, param, ignoreValue) {
				matched = true

				updatedParam, keep := rule.Apply(param)
				if keep {
					filteredParams = append(filteredParams, updatedParam)
				}

				break
			}
		}

		// If there is no matching rule for a parameter, it automatically falls through unmodified
		if !matched {
			filteredParams = append(filteredParams, param)
		}
	}

	profile.Parameters = filteredParams
}

// paramsMatch returns true when param fulfills all matching critera in matcher
func paramsMatch(matcher Parameter, param Parameter, ignoreValue bool) bool {
	nameRe := regexp.MustCompile(matcher.Name)
	cfgRe := regexp.MustCompile(matcher.ConfigFile)
	valueRe := regexp.MustCompile(matcher.Value)

	if nil != nameRe.FindStringIndex(param.Name) &&
		nil != cfgRe.FindStringIndex(param.ConfigFile) {

		if ignoreValue || nil != valueRe.FindStringIndex(param.Value) {
			return true

		} else {
			fmt.Fprintf(os.Stderr, "[ACTION REQUIRED] Found modified value. Skip modifying {\"%s\", \"%s\", \"%s\"}. Please update manually\n",
				param.Name, param.ConfigFile, param.Value)
		}
	}
	return false
}

func UpdateDetails(p *Profile, rules *ConversionPolicy) {
	p.Description.Name = strings.Replace(p.Description.Name, rules.ReplaceName.Old, rules.ReplaceName.New, -1)
	p.Description.Description = strings.Replace(p.Description.Description, rules.ReplaceDescription.Old, rules.ReplaceDescription.New, -1)
}

func main() {
	cfgParam := parseArgs()
	fmt.Fprintf(os.Stderr, "Traffic Control Profile Conversion Utility\n")
	fmt.Fprintf(os.Stderr, "Input Profile: %s\n", cfgParam.InProfile)
	fmt.Fprintf(os.Stderr, "Conversion Rules: %s\n", cfgParam.Rules)
	if cfgParam.Force {
		fmt.Fprintf(os.Stderr, "[WARNING] Ignoring existing parameter values in comparisons, making all suggested changes\n")
	}

	inProfile := parseInputProfile(cfgParam.InProfile)
	rules := parseInputRules(cfgParam.Rules)

	if !ValidateParameters(inProfile, rules.ValidateParameters) {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to validate required parameters in profile\n")
		os.Exit(-1)
	}
	ConvertProfile(inProfile, rules.ConversionRules, cfgParam.Force)
	UpdateDetails(inProfile, rules)
	AddParameters(inProfile, rules.AddParameters)

	// Can't use the standard JSON Marshaller because it forces HTML escape
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(inProfile); err != nil {
		panic(err)
	}

	indentedBuffer := new(bytes.Buffer)
	if err := json.Indent(indentedBuffer, buf.Bytes(), "", "    "); err != nil {
		panic(err)
	}

	if cfgParam.OutFile != "" {
		err := ioutil.WriteFile(cfgParam.OutFile, indentedBuffer.Bytes(), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Cannot write output file")
			panic(err)
		}
	} else {
		fmt.Printf("%s", indentedBuffer.String())
	}
}
