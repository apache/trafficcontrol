package deliveryservice

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
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// Validate all fields for a delivery service
func Validate(ds tc.DeliveryService) []error {
	var errors []error
	validations := []func(tc.DeliveryService) error{
		ValidActive,
		ValidCacheURL,
		ValidCCRDNSTTL,
		ValidCDNID,
		ValidCheckPath,
		ValidDisplayName,
		ValidDNSBypassCname,
		ValidDNSBypassIP,
		ValidDNSBypassIP6,
		ValidDNSBypassTTL,
		ValidDSCP,
		ValidEdgeHeaderRewrite,
		ValidGeoLimit,
		ValidGeoLimitCountries,
		ValidGeoLimitRedirectURL,
		ValidGeoProvider,
		ValidGlobalMaxMbps,
		ValidGlobalMaxTps,
		ValidHTTPBypassFQDN,
		ValidInfoURL,
		ValidInitialDispersion,
		ValidIPv6RoutingEnabled,
		ValidLogsEnabled,
		ValidLongDesc,
		ValidLongDesc1,
		ValidLongDesc2,
		ValidMaxDNSAnswers,
		ValidMidHeaderRewrite,
		ValidMissLat,
		ValidMissLong,
		ValidMultiSiteOrigin,
		ValidMultiSiteOriginAlgorithm,
		ValidOrgServerFqdn,
		ValidOriginShield,
		ValidProfileID,
		ValidProtocol,
		ValidQstringIgnore,
		ValidRangeRequestHandling,
		ValidRegexRemap,
		ValidRegionalGeoBlocking,
		ValidRemapText,
		ValidRoutingName,
		ValidSigningAlgorithm,
		ValidSslKeyVersion,
		ValidTenantID,
		ValidTrRequestHeaders,
		ValidTrResponseHeaders,
		ValidTypeID,
		ValidXMLID,
	}

	for _, v := range validations {
		if err := v(ds); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// ValidActive ...
func ValidActive(ds tc.DeliveryService) error {
	return nil
}

// ValidCacheURL ,,,
func ValidCacheURL(ds tc.DeliveryService) error {
	return nil
}

// ValidCCRDNSTTL ...
func ValidCCRDNSTTL(ds tc.DeliveryService) error {
	return nil
}

// ValidCDNID ...
func ValidCDNID(ds tc.DeliveryService) error {
	// TODO: validate exists in cdn table
	return nil
}

// ValidCheckPath ...
func ValidCheckPath(ds tc.DeliveryService) error {
	return nil
}

// ValidDNSBypassCname ...
func ValidDNSBypassCname(ds tc.DeliveryService) error {
	return nil
}

// ValidDNSBypassIP ...
func ValidDNSBypassIP(ds tc.DeliveryService) error {
	// TODO: valid IP address
	return nil
}

// ValidDNSBypassIP6 ...
func ValidDNSBypassIP6(ds tc.DeliveryService) error {
	// TODO: valid IP address
	return nil
}

// ValidDNSBypassTTL ...
func ValidDNSBypassTTL(ds tc.DeliveryService) error {
	// TODO: valid TTL range
	return nil
}

// ValidDSCP ...
func ValidDSCP(ds tc.DeliveryService) error {
	return nil
}

// ValidEdgeHeaderRewrite ...
func ValidEdgeHeaderRewrite(ds tc.DeliveryService) error {
	return nil
}

// ValidDisplayName ...
func ValidDisplayName(ds tc.DeliveryService) error {
	if len(ds.DisplayName) > 48 {
		return fmt.Errorf("display name '%s' can be no longer than 48 characters", ds.DisplayName)
	}
	return nil
}

// ValidGeoLimit ...
func ValidGeoLimit(ds tc.DeliveryService) error {
	return nil
}

// ValidGeoLimitCountries ...
func ValidGeoLimitCountries(ds tc.DeliveryService) error {
	return nil
}

// ValidGeoLimitRedirectURL ...
func ValidGeoLimitRedirectURL(ds tc.DeliveryService) error {
	return nil
}

// ValidGeoProvider ...
func ValidGeoProvider(ds tc.DeliveryService) error {
	return nil
}

// ValidGlobalMaxMbps ...
func ValidGlobalMaxMbps(ds tc.DeliveryService) error {
	return nil
}

// ValidGlobalMaxTps ...
func ValidGlobalMaxTps(ds tc.DeliveryService) error {
	return nil
}

// ValidHTTPBypassFQDN ...
func ValidHTTPBypassFQDN(ds tc.DeliveryService) error {
	return nil
}

// ValidInfoURL ...
func ValidInfoURL(ds tc.DeliveryService) error {
	return nil
}

// ValidInitialDispersion ...
func ValidInitialDispersion(ds tc.DeliveryService) error {
	return nil
}

// ValidIPv6RoutingEnabled ...
func ValidIPv6RoutingEnabled(ds tc.DeliveryService) error {
	return nil
}

// ValidLogsEnabled ...
func ValidLogsEnabled(ds tc.DeliveryService) error {
	return nil
}

// ValidLongDesc ...
func ValidLongDesc(ds tc.DeliveryService) error {
	return nil
}

// ValidLongDesc1 ...
func ValidLongDesc1(ds tc.DeliveryService) error {
	return nil
}

// ValidLongDesc2 ...,
func ValidLongDesc2(ds tc.DeliveryService) error {
	return nil
}

// ValidMaxDNSAnswers ...
func ValidMaxDNSAnswers(ds tc.DeliveryService) error {
	return nil
}

// ValidMidHeaderRewrite ...
func ValidMidHeaderRewrite(ds tc.DeliveryService) error {
	return nil
}

// ValidMissLat ...
func ValidMissLat(ds tc.DeliveryService) error {
	if math.Abs(ds.MissLat) > 90 {
		return fmt.Errorf("missLat value %2.0f must not exceed +/- 90.0", ds.MissLat)
	}
	return nil
}

// ValidMissLong ...
func ValidMissLong(ds tc.DeliveryService) error {
	if math.Abs(ds.MissLong) > 90 {
		return fmt.Errorf("missLong value %2.0f must not exceed +/- 90.0", ds.MissLong)
	}
	return nil
}

// ValidMultiSiteOrigin ...
func ValidMultiSiteOrigin(ds tc.DeliveryService) error {
	return nil
}

// ValidMultiSiteOriginAlgorithm ...
func ValidMultiSiteOriginAlgorithm(ds tc.DeliveryService) error {
	return nil
}

// ValidOrgServerFqdn ...
func ValidOrgServerFqdn(ds tc.DeliveryService) error {
	return nil
}

// ValidOriginShield ...
func ValidOriginShield(ds tc.DeliveryService) error {
	return nil
}

// ValidProfileID ...
func ValidProfileID(ds tc.DeliveryService) error {
	// TODO: validate exists in profile table
	return nil
}

// ValidProtocol ...
func ValidProtocol(ds tc.DeliveryService) error {
	return nil
}

// ValidQstringIgnore ...
func ValidQstringIgnore(ds tc.DeliveryService) error {
	return nil
}

// ValidRangeRequestHandling ...
func ValidRangeRequestHandling(ds tc.DeliveryService) error {
	return nil
}

// ValidRegexRemap ...
func ValidRegexRemap(ds tc.DeliveryService) error {
	return nil
}

// ValidRegionalGeoBlocking ...
func ValidRegionalGeoBlocking(ds tc.DeliveryService) error {
	return nil
}

// ValidRemapText ...
func ValidRemapText(ds tc.DeliveryService) error {
	return nil
}

// ValidRoutingName ...
func ValidRoutingName(ds tc.DeliveryService) error {
	if len(ds.RoutingName) > 48 {
		return fmt.Errorf("routing name '%s' can be no longer than 48 characters", ds.RoutingName)
	}
	return nil
}

// ValidSigningAlgorithm ...
func ValidSigningAlgorithm(ds tc.DeliveryService) error {
	return nil
}

// ValidSslKeyVersion ...
func ValidSslKeyVersion(ds tc.DeliveryService) error {
	return nil
}

// ValidTenantID ...
func ValidTenantID(ds tc.DeliveryService) error {
	// TODO: validate exists in tenant table
	return nil
}

// ValidTrRequestHeaders ...
func ValidTrRequestHeaders(ds tc.DeliveryService) error {
	return nil
}

// ValidTrResponseHeaders ...
func ValidTrResponseHeaders(ds tc.DeliveryService) error {
	return nil
}

// ValidTypeID ...
func ValidTypeID(ds tc.DeliveryService) error {
	// TODO: validate exists in type table
	return nil
}

// ValidXMLID ...
func ValidXMLID(ds tc.DeliveryService) error {
	if len(ds.XMLID) == 0 {
		return fmt.Errorf("xmlId is required")
	}
	if strings.IndexFunc(ds.XMLID, unicode.IsSpace) != -1 {
		return fmt.Errorf("xmlId '%s' must not contain spaces", ds.XMLID)
	}
	if len(ds.XMLID) > 48 {
		return fmt.Errorf("xmlId '%s' must be no longer than 48 characters", ds.XMLID)
	}
	return nil
}
