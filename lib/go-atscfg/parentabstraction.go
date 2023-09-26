package atscfg

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
	"errors"
	"net"
	"strconv"
	"strings"
)

// ParentAbstraction contains all the data necessary to build either parent.config or strategies.yaml.
type ParentAbstraction struct {
	Services []*ParentAbstractionService
	// Peers is the list of peer proxy caches to be used in a strategy peering group.
	// a cache will only have one set of peers for potential use in all delivery services.
	Peers []*ParentAbstractionServiceParent
}

// ParentAbstractionService represents a single delivery service's parent data.
// For parent.config, this becomes a single dest_domain= line.
// For strategies, this becomes a strategy along with its corresponding groups and hosts.
type ParentAbstractionService struct {
	// Name is a unique name for the service.
	// It can be anything unique, but should probably be the Traffic ops Delivery Service name.
	Name string
	// Comment is a text comment about the service, not including the comment syntax (e.g. # or //).
	// Should be empty if !opt.AddComments.
	Comment string
	// DestDomain is the FQDN of the remap.config target.
	// Becomes parent.config dest_domain directive
	// Becomes strategies.yaml TODO
	DestDomain string

	// Port is the port of the remap.config target,
	// which MUST be valid, and is implicitly 80 for http targets and 443 for https targets.
	// Becomes parent.config port directive
	// Becomes strategies.yaml TODO
	Port int

	// Parents is the list of parents, either parent proxy caches or origins.
	// This is a sorted array. Parents will be inserted into the config file in the order they appear.
	// Becomes parent.config parent= directive members
	// Becomes strategies.yaml TODO
	Parents []*ParentAbstractionServiceParent

	// Parents is the list of secondary parents, either parent proxy caches or origins,
	// to be used if the primary parents fail. See SecondaryMode.
	// Becomes parent.config secondary_parent= directive members
	// Becomes strategies.yaml TODO
	SecondaryParents []*ParentAbstractionServiceParent

	// SecondaryMode is how to try SecondaryParents if primary Parents fail.
	// Becomes parent.config secondary_mode directive
	// Becomes strategies.yaml TODO
	SecondaryMode ParentAbstractionServiceParentSecondaryMode

	// CachePeerResult is used only when the RetryPolicy is set to
	// 'consistent_hash' and the SecondaryMode is set to 'peering'.
	// In the case that it's used and set to 'true', query results
	// from peer caches will not be cached locally.
	CachePeerResult bool

	// GoDirect is whether to go direct to parents via normal HTTP requests.
	// False means to make proxy requests to the parents.
	// Becomes parent.config go_direct and parent_is_proxy directives
	// Becomes strategies.yaml TODO
	GoDirect bool

	// ParentIsProxy A boolean value which indicates if the groups of hosts are proxy caches or origins.
	// true (default) means all the hosts used are Traffic Server caches.
	// false means the hosts are origins.
	// Becomes parent_is_proxy directive.
	// Becomes strategies.yaml TODO
	ParentIsProxy bool

	// IgnoreQueryStringInParentSelection is whether to use the query string of the request
	// when selecting a parent, e.g. via Consistent Hash.
	// Becomes parent.config qstring directive
	// Becomes strategies.yaml TODO
	IgnoreQueryStringInParentSelection bool

	// MarkdownResponseCodes is the list of HTTP response codes from the parent
	// to consider as errors and mark the parent as unhealthy. Typically 5xx codes.
	// Becomes parent.config unavailable_server_retry_responses directive
	// Becomes strategies.yaml TODO
	MarkdownResponseCodes []int

	// ErrorResponseCodes is the list of HTTP response codes from the parent
	// to consider as errors, but NOT mark the parent unhealthy. Typically 4xx codes.
	// Becomes parent.config unavailable_server_retry_responses directive
	// Becomes strategies.yaml TODO
	ErrorResponseCodes []int

	// MaxSimpleRetries is the maximum number of non-markdown errors to attempt
	// before returning the error to the client. See ErrorResponseCodes
	// Becomes parent.config max_simple_retries
	// Becomes strategies.yaml TODO
	MaxSimpleRetries int

	// MaxMarkdownRetries is the maximum number of markdown errors to attempt
	// before returning the error to the client. See MarkdownResponseCodes
	// Becomes parent.config max_unavailable_server_retries
	// Becomes strategies.yaml TODO
	MaxMarkdownRetries int

	// RetryPolicy is how to retry primary versus secondary parents.
	// Becomes parent.config round_robin directive
	// Becomes strategies.yaml TODO
	RetryPolicy ParentAbstractionServiceRetryPolicy

	// Weight is the weight of this parent relative to other parents in consistent hash (and potentially other non-sequential) parent selection. The default is 0.999
	// Becomes parent.config weight directive
	// Becomes strategies.yaml TODO
	Weight float64

	// DS is the delivery service associated with the service
	DS DeliveryService
}

// ParentAbstractionServices implements sort.Interface.
type ParentAbstractionServices []*ParentAbstractionService

// Len implements part of sort.Interface.
func (ps ParentAbstractionServices) Len() int { return len(ps) }

// Swap implements part of sort.Interface.
func (ps ParentAbstractionServices) Swap(i, j int) { ps[i], ps[j] = ps[j], ps[i] }

// Less implements part of sort.Interface.
func (ps ParentAbstractionServices) Less(i, j int) bool {
	if ps[i].DestDomain != ps[j].DestDomain {
		return ps[i].DestDomain < ps[j].DestDomain
	}
	return ps[i].Port < ps[j].Port
}

// ParentAbstractionServiceParentSecondaryMode is the "secondary parent mode" of
// a Delivery Service parenting abstraction. Only certain values are allowed.
type ParentAbstractionServiceParentSecondaryMode string

// The allowable values of a ParentAbstractionServiceParentSecondaryMode (with
// the exception of ParentAbstractionServiceParentSecondaryModeInvalid, which
// does not represent a valid ParentAbstractionServiceParentSecondaryMode).
const (
	ParentAbstractionServiceParentSecondaryModeExhaust   = ParentAbstractionServiceParentSecondaryMode("exhaust")
	ParentAbstractionServiceParentSecondaryModeAlternate = ParentAbstractionServiceParentSecondaryMode("alternate")
	ParentAbstractionServiceParentSecondaryModePeering   = ParentAbstractionServiceParentSecondaryMode("peering")
	ParentAbstractionServiceParentSecondaryModeInvalid   = ParentAbstractionServiceParentSecondaryMode("")
)

// ParentAbstractionServiceParentSecondaryModeDefault is the "secondary parent
// mode" that is used in parenting abstraction if one is not explicitly
// configured.
const ParentAbstractionServiceParentSecondaryModeDefault = ParentAbstractionServiceParentSecondaryModeAlternate

// ToParentDotConfigVal returns the ATS parent.config secondary_mode= value for the enum.
// If the mode is invalid, ParentAbstractionServiceParentSecondaryModeDefault is returned without error.
func (mo ParentAbstractionServiceParentSecondaryMode) ToParentDotConfigVal() string {
	switch mo {
	case ParentAbstractionServiceParentSecondaryModeExhaust:
		return "2"
	case ParentAbstractionServiceParentSecondaryModeAlternate:
		return "1"
	default:
		return ParentAbstractionServiceParentSecondaryModeDefault.ToParentDotConfigVal()
	}
}

// A ParentAbstractionServiceRetryPolicy is a "retry policy" that will be used
// by Delivery Service parenting.
type ParentAbstractionServiceRetryPolicy string

// These are the valid value of a ParentAbstractionServiceRetryPolicy - with the
// exception of ParentAbstractionServiceRetryPolicyInvalid, which does not
// represent a valid ParentAbstractionServiceRetryPolicy.
const (
	ParentAbstractionServiceRetryPolicyRoundRobinIP     = ParentAbstractionServiceRetryPolicy("round_robin_ip")
	ParentAbstractionServiceRetryPolicyRoundRobinStrict = ParentAbstractionServiceRetryPolicy("round_robin_strict")
	ParentAbstractionServiceRetryPolicyFirst            = ParentAbstractionServiceRetryPolicy("first")
	ParentAbstractionServiceRetryPolicyLatched          = ParentAbstractionServiceRetryPolicy("latched")
	ParentAbstractionServiceRetryPolicyConsistentHash   = ParentAbstractionServiceRetryPolicy("consistent_hash")
	ParentAbstractionServiceRetryPolicyInvalid          = ParentAbstractionServiceRetryPolicy("")
)

// DefaultParentAbstractionServiceRetryPolicy is the "retry policy" that will be
// used by Delivery Service parenting if one is not explicitly configured.
const DefaultParentAbstractionServiceRetryPolicy = ParentAbstractionServiceRetryPolicyConsistentHash

// ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy converts a parent
// selection algorithm Parameter Value to a generic
// ParentAbstractionServiceRetryPolicy.
func ParentSelectAlgorithmToParentAbstractionServiceRetryPolicy(alg string) ParentAbstractionServiceRetryPolicy {
	switch strings.TrimSpace(strings.ToLower(alg)) {
	case "true":
		return ParentAbstractionServiceRetryPolicyRoundRobinIP
	case "strict":
		return ParentAbstractionServiceRetryPolicyRoundRobinStrict
	case "false":
		return ParentAbstractionServiceRetryPolicyFirst
	case "consistent_hash":
		return ParentAbstractionServiceRetryPolicyConsistentHash
	case "latched":
		return ParentAbstractionServiceRetryPolicyLatched
	default:
		return ParentAbstractionServiceRetryPolicyInvalid
	}
}

// ToParentDotConfigFormat returns the ATS parent.config round_robin= value for the policy.
// If the policy is invalid, the default is returned without error.
func (po ParentAbstractionServiceRetryPolicy) ToParentDotConfigFormat() string {
	switch po {
	case ParentAbstractionServiceRetryPolicyRoundRobinIP:
		return "true"
	case ParentAbstractionServiceRetryPolicyRoundRobinStrict:
		return "strict"
	case ParentAbstractionServiceRetryPolicyFirst:
		return "false"
	case ParentAbstractionServiceRetryPolicyLatched:
		return "latched"
	case ParentAbstractionServiceRetryPolicyConsistentHash:
		return "consistent_hash"
	default:
		return "consistent_hash"
	}
}

// ParentSelectParamQStringHandlingToBool returns whether the param is to use the query string in the parent select algorithm or not.
// If the parameter value is not valid, returns nil.
func ParentSelectParamQStringHandlingToBool(paramVal string) *bool {
	switch strings.TrimSpace(strings.ToLower(paramVal)) {
	case "consider":
		v := true
		return &v
	case "ignore":
		v := false
		return &v
	}
	return nil
}

// ParentAbstractionServiceParent represents a single "parent" as an abstracted
// concept.
type ParentAbstractionServiceParent struct {
	// FQDN is the parent FQDN that ATS will use. Note this may be an IP.
	FQDN   string
	Port   int
	Weight float64
}

// Key returns a unique key that can be used to compare parents for equality.
func (sp ParentAbstractionServiceParent) Key() string {
	return sp.FQDN + ":" + strconv.Itoa(sp.Port)
}

type peersSort []*ParentAbstractionServiceParent

func (a peersSort) Len() int           { return len(a) }
func (a peersSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a peersSort) Less(i, j int) bool { return a[i].Key() < a[j].Key() }

// RemoveParentDuplicates returns all values in the input list that have unique
// outputs for their Key method. Earlier duplicates are used while later
// occurrences of degenerate "Key"s are discarded.
func RemoveParentDuplicates(inputs []*ParentAbstractionServiceParent, seens map[string]struct{}) ([]*ParentAbstractionServiceParent, map[string]struct{}) {
	if seens == nil {
		seens = make(map[string]struct{})
	}
	uniques := []*ParentAbstractionServiceParent{}
	for _, input := range inputs {
		key := input.Key()
		if _, ok := seens[key]; !ok {
			uniques = append(uniques, input)
			seens[key] = struct{}{}
		}
	}
	return uniques, seens
}

// ParseRetryResponses parses a raw Parameter Value containing HTTP response
// codes for scenarios when parents should be "retried" into a list of the
// actual numeric codes.
func ParseRetryResponses(resp string) ([]int, error) {
	resp = strings.TrimSpace(resp)
	if len(resp) > 2 && resp[0] == '"' {
		resp = resp[1 : len(resp)-1]
	}
	codes := []int{}
	codeStrs := strings.Split(resp, ",")
	for _, codeStr := range codeStrs {
		codeStr = strings.TrimSpace(codeStr)
		if codeStr == "" {
			continue
		}
		code, err := strconv.Atoi(codeStr)
		if err != nil {
			return nil, errors.New("malformed")
		}
		codes = append(codes, code)
	}
	return codes, nil
}

// DefaultSimpleRetryCodes is the set of HTTP response codes that are used to
// indicate a parent should be "retried" if none are explicitly configured.
var DefaultSimpleRetryCodes = []int{404}

// DefaultUnavailableServerRetryCodes is the set of HTTP response codes that are
// used to indicate a parent is "unavailable" and should be "retried" if none
// are explicitly configured.
var DefaultUnavailableServerRetryCodes = []int{503}

// DefaultIgnoreQueryStringInParentSelection is used to decide whether a
// request's query string should be used or dropped during selecting a parent
// when that behavior is not explicitly configured.
const DefaultIgnoreQueryStringInParentSelection = false

func parentAbstractionToParentDotConfig(pa *ParentAbstraction, opt *ParentConfigOpts, atsMajorVersion uint) (string, []string, error) {
	warnings := []string{}
	txt := ""

	// parent.config dest_domain directives must be unique.
	// This is the "duplicate origin problem"
	processedOriginsToDSNames := map[string]string{}

	for _, svc := range pa.Services {
		if existingDS, ok := processedOriginsToDSNames[svc.DestDomain]; ok {
			warnings = append(warnings, "duplicate origin! DS '"+svc.Name+"' and '"+existingDS+"' share origin '"+svc.DestDomain+"': skipping '"+svc.Name+"'!")
			continue
		}

		svcLine, svcWarns, err := svc.ToParentDotConfigLine(opt, atsMajorVersion)
		warnings = append(warnings, svcWarns...)
		if err != nil {
			// TODO add DS name
			// TODO don't error? No single delivery service should be able to break others.
			return "", warnings, errors.New("creating parent.config line from service: " + err.Error())
		}
		txt += svcLine + "\n"

		processedOriginsToDSNames[svc.DestDomain] = svc.Name
	}
	return txt, warnings, nil
}

// ToParentDotConfigLine constructs a line in the parent.config Apache Traffic
// Server configuration file for the abstraction with the given options and for
// the given major version of Apache Traffic Server. It returns the line, any
// warnings to be issued, and any error that occurred during generation.
func (svc *ParentAbstractionService) ToParentDotConfigLine(opt *ParentConfigOpts, atsMajorVersion uint) (string, []string, error) {
	warnings := []string{}
	txt := ""
	if opt.AddComments && svc.Comment != "" {
		txt += LineCommentParentDotConfig + " " + svc.Comment + "\n"
	}

	// if the domain is an IP, we have to use dest_ip.
	// Using an IP in dest_domain will be silently ignored by ATS, causing the remap/DS to use the fallthrough 'dest_domain=.' rule
	if domainIsIP := net.ParseIP(svc.DestDomain) != nil; domainIsIP {
		txt += `dest_ip=` + svc.DestDomain
	} else {
		txt += `dest_domain=` + svc.DestDomain
	}

	if svc.Port != 0 {
		txt += ` port=` + strconv.Itoa(svc.Port)
	}

	if atsMajorVersion >= 6 && svc.RetryPolicy == ParentAbstractionServiceRetryPolicyConsistentHash && len(svc.SecondaryParents) > 0 {
		// TODO add quotes
		if len(svc.Parents) > 0 {
			txt += ` parent="` + ParentAbstractionServiceParentsToParentDotConfigLine(svc.Parents) + `"`
		}
		if len(svc.SecondaryParents) > 0 {
			txt += ` secondary_parent="` + ParentAbstractionServiceParentsToParentDotConfigLine(svc.SecondaryParents) + `"`
		}
		txt += ` secondary_mode=` + svc.SecondaryMode.ToParentDotConfigVal()
	} else {
		parents := []*ParentAbstractionServiceParent{}
		for _, pa := range svc.Parents {
			parents = append(parents, pa)
		}
		for _, pa := range svc.SecondaryParents {
			parents = append(parents, pa)
		}
		if len(parents) > 0 {
			txt += ` parent="` + ParentAbstractionServiceParentsToParentDotConfigLine(parents) + `"`
		}
	}

	txt += ` round_robin=` + svc.RetryPolicy.ToParentDotConfigFormat()
	txt += ` go_direct=` + strconv.FormatBool(svc.GoDirect)
	txt += ` qstring=`
	if !svc.IgnoreQueryStringInParentSelection {
		txt += `consider`
	} else {
		txt += `ignore`
	}
	txt += ` parent_is_proxy=` + strconv.FormatBool(svc.ParentIsProxy)

	if svc.MaxSimpleRetries > 0 && svc.MaxMarkdownRetries > 0 {
		txt += ` parent_retry=both`
	} else if svc.MaxSimpleRetries > 0 {
		txt += ` parent_retry=simple_retry`
	} else if svc.MaxMarkdownRetries > 0 {
		txt += ` parent_retry=unavailable_server_retry`
	}

	if svc.MaxSimpleRetries > 0 {
		txt += ` max_simple_retries=` + strconv.Itoa(svc.MaxSimpleRetries)
	}
	if svc.MaxMarkdownRetries > 0 {
		txt += ` max_unavailable_server_retries=` + strconv.Itoa(svc.MaxMarkdownRetries)
	}

	if len(svc.ErrorResponseCodes) > 0 {
		if atsMajorVersion >= 9 {
			txt += ` simple_server_retry_responses="` + strings.Join(intsToStrs(svc.ErrorResponseCodes), `,`) + `"`
		} else {
			warnings = append(warnings, "Service '"+svc.Name+"' had simple retry codes '"+strings.Join(intsToStrs(svc.ErrorResponseCodes), ",")+"' but ATS version "+strconv.FormatUint(uint64(atsMajorVersion), 10)+" < 9 does not support custom simple retry codes, omitting!")
		}
	}

	if len(svc.MarkdownResponseCodes) > 0 {
		txt += ` unavailable_server_retry_responses="` + strings.Join(intsToStrs(svc.MarkdownResponseCodes), `,`) + `"`
	}

	// 	txt += ` unavailable_server_retry_responses=` + unavailableServerRetryResponses

	return txt, warnings, nil
}

func intsToStrs(is []int) []string {
	strs := []string{}
	for _, i := range is {
		strs = append(strs, strconv.Itoa(i))
	}
	return strs
}

// ToParentDotConfigFormat converts the abstracted parent into a concrete piece
// of a configuration file line for the parent.config parent configuration
// implementation.
func (pa *ParentAbstractionServiceParent) ToParentDotConfigFormat() string {
	return pa.FQDN + ":" + strconv.Itoa(pa.Port) + "|" + strconv.FormatFloat(pa.Weight, 'f', -1, 64)
}

// ParentDotConfigParentSeparator is the string used to delimit multiple parents
// on a line in a parent.config Apache Traffic Server configuration file.
const ParentDotConfigParentSeparator = `;`

// ParentAbstractionServiceParentsToParentDotConfigLine creates a line in the
// parent.config implementation of a parenting configuration file given the set
// of parents it should configure.
func ParentAbstractionServiceParentsToParentDotConfigLine(parents []*ParentAbstractionServiceParent) string {
	parentStrs := []string{}
	for _, parent := range parents {
		parentStrs = append(parentStrs, parent.ToParentDotConfigFormat())
	}
	return strings.Join(parentStrs, ParentDotConfigParentSeparator)
}
