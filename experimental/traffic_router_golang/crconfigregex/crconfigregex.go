package crconfigregex

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
 *
 */

import (
	"errors"
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// TODO remove duplication with Traffic Monitor

// Regexes maps Delivery Service Regular Expressions to delivery services.
// For performance, we categorize Regular Expressions into 3 categories:
// 1. Direct string matches, with no regular expression matching characters
// 2. .*\.foo\..* expressions, where foo is a direct string match with no regular expression matching characters
// 3. Everything else
// This allows us to do a cheap match on 1 and 2, and only regex match the uncommon case.
type Regexes struct {
	DirectMatches                      map[string]tc.DeliveryServiceName
	DotStartSlashDotFooSlashDotDotStar map[string]tc.DeliveryServiceName
	RegexMatch                         map[*regexp.Regexp]tc.DeliveryServiceName
}

// DeliveryService returns the delivery service which matches the given fqdn, or false.
func (d Regexes) DeliveryService(domain, subdomain, subsubdomain string) (tc.DeliveryServiceName, bool) {
	if ds, ok := d.DotStartSlashDotFooSlashDotDotStar[subdomain]; ok {
		return ds, true
	}
	fqdn := subsubdomain + "." + subdomain + "." + domain
	if ds, ok := d.DirectMatches[fqdn]; ok {
		return ds, true
	}
	for regex, ds := range d.RegexMatch {
		if regex.MatchString(fqdn) {
			return ds, true
		}
	}
	return "", false
}

// NewRegexes constructs a new Regexes object, initializing internal pointer members.
func new() Regexes {
	return Regexes{
		DirectMatches:                      map[string]tc.DeliveryServiceName{},
		DotStartSlashDotFooSlashDotDotStar: map[string]tc.DeliveryServiceName{},
		RegexMatch:                         map[*regexp.Regexp]tc.DeliveryServiceName{},
	}
}

// getDeliveryServiceRegexes gets the regexes of each delivery service, for the given CDN, from Traffic Ops.
// Returns a map[deliveryService][]regex.
func Get(crc *tc.CRConfig) (Regexes, error) {
	dsRegexes := map[tc.DeliveryServiceName][]string{}

	for dsNameStr, dsData := range crc.DeliveryServices {
		dsName := tc.DeliveryServiceName(dsNameStr)
		if len(dsData.MatchSets) < 1 {
			return Regexes{}, errors.New("CRConfig missing regex for '" + string(dsName) + "'")
		}
		for _, matchset := range dsData.MatchSets {
			if len(matchset.MatchList) < 1 {
				return Regexes{}, errors.New("CRConfig missing Regex for '" + string(dsName) + "'")
			}
			dsRegexes[dsName] = append(dsRegexes[dsName], matchset.MatchList[0].Regex)
		}
	}

	return createRegexes(dsRegexes)
}

// TODO precompute, move to TOData; call when we get new delivery services, instead of every time we create new stats
func createRegexes(dsToRegex map[tc.DeliveryServiceName][]string) (Regexes, error) {
	dsRegexes := new()

	for ds, regexStrs := range dsToRegex {
		for _, regexStr := range regexStrs {
			prefix := `.*\.`
			suffix := `\..*`
			if strings.HasPrefix(regexStr, prefix) && strings.HasSuffix(regexStr, suffix) {
				matchStr := regexStr[len(prefix) : len(regexStr)-len(suffix)]
				if otherDs, ok := dsRegexes.DotStartSlashDotFooSlashDotDotStar[matchStr]; ok {
					return dsRegexes, errors.New("duplicate regex " + regexStr + " (" + matchStr + ") in " + string(ds) + " and " + string(otherDs))
				}
				dsRegexes.DotStartSlashDotFooSlashDotDotStar[matchStr] = ds
				continue
			}
			if !strings.ContainsAny(regexStr, `[]^\:{}()|?+*,=%@<>!'`) {
				if otherDs, ok := dsRegexes.DirectMatches[regexStr]; ok {
					return dsRegexes, errors.New("duplicate Regex " + regexStr + " in " + string(ds) + " and " + string(otherDs))
				}
				dsRegexes.DirectMatches[regexStr] = ds
				continue
			}
			// TODO warn? regex matches are unusual
			r, err := regexp.Compile(regexStr)
			if err != nil {
				return dsRegexes, errors.New("regex " + regexStr + " failed to compile: " + err.Error())
			}
			dsRegexes.RegexMatch[r] = ds
		}
	}
	return dsRegexes, nil
}
