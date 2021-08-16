package tc

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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// EmailTemplate is an html/template.Template for formatting DeliveryServiceRequestRequests into
// text/html email bodies. Its direct use is discouraged, instead use
// DeliveryServiceRequestRequest.Format.
//
// Deprecated: Delivery Services Requests have been deprecated in favor of
// Delivery Service Requests, and will be removed from the Traffic Ops API at
// some point in the future.
var EmailTemplate = template.Must(template.New("Email Template").Parse(`<!DOCTYPE html>
<html lang="en-US">
<head>
<meta charset="utf-8"/>
<title>Delivery Service Request for {{.Customer}}</title>
<style>
aside {
	padding: 0 1em;
	color: #6A737D;
	border-left: .25em solid #DFE2E5;
}
body {
	font-family: sans;
	background-color: white;
}
pre {
	padding: 5px;
	background-color: lightgray;
}
</style>
</head>
<body>
<h1>Delivery Service Request for {{.Customer}}</h1>
<p>{{.ServiceDesc}}</p>
<section>
	<details>
		<summary><h2>Service Description</h2></summary>
		<h3>Content Type</h3>
		<p>{{.ContentType}}</p>
		<h3>Delivery Protocol</h3>
		<p>{{.DeliveryProtocol.String}}</p>
		<h3>Routing Type</h3>
		<p>{{.RoutingType.String}}</p>
	</details>
</section>
<section>
	<details>
		<summary><h2>Traffic &amp; Library Estimates</h2></summary>
		<h3>Peak Bandwidth Estimate</h3>
		<p>{{.PeakBPSEstimate}}Bps</p>
		<h3>Peak Transactions per Second Estimate</h3>
		<p>{{.PeakTPSEstimate}}Tps</p>
		<h3>Max Library Size Estimate</h3>
		<p>{{.MaxLibrarySizeEstimate}}GB</p>
	</details>
</section>
<section>
	<details>
		<summary><h2>Origin Security</h2></summary>
		<h3>Origin Server URL</h3>
		<p><a href="{{.OriginURL}}">{{.OriginURL}}</a></p>
		<h3>Origin Dynamic Remap</h3>
		<p>{{.HasOriginDynamicRemap}}</p>
		<h3>Origin Test File</h3>
		<p>{{.OriginTestFile}}</p>
		<h3>ACL/Whitelist to Access Origin</h3>
		<p>{{.HasOriginACLWhitelist}}</p>
		{{if .OriginHeaders}}<h3>Header(s) to Access Origin</h3>
		<ul>{{range .OriginHeaders}}
			<li>{{.}}</li>{{end}}
		</ul>{{end}}
		<h3>Other Origin Security</h3>
		<p>{{if .OtherOriginSecurity}}{{.OtherOriginSecurity}}{{else}}None{{end}}</p>
	</details>
</section>
<section>
	<details>
		<summary><h2>Core Features</h2></summary>
		<h3>Query String Handling</h3>
		<p>{{.QueryStringHandling}}</p>
		<h3>Range Request Handling</h3>
		<p>{{.RangeRequestHandling}}</p>
		<h3>Signed URLs / URL Tokenization</h3>
		<p>{{.HasSignedURLs}}</p>
		<h3>Negative Caching Customization</h3>
		<p>{{.HasNegativeCachingCustomization}}</p>
		{{if or .HasNegativeCachingCustomization .NegativeCachingCustomizationNote }}<aside>
			<p>{{.NegativeCachingCustomizationNote}}</p>
		</aside>{{else if .HasNegativeCachingCustomization}}<aside>
			<p><b>No instructions given!</b></p>
		</aside>{{end}}
		{{if .ServiceAliases}}<h3>Service Aliases</h3>
		<ul>{{range .ServiceAliases}}
			<li>{{.}}</li>{{end}}
		</ul>{{end}}
	</details>
</section>
{{if or .RateLimitingGBPS .RateLimitingTPS .OverflowService}}<section>
	<details>
		<summary><h2>Service Limits</h2></summary>
		{{if .RateLimitingGBPS}}<h3>Bandwidth Limit</h3>
		<p>{{.RateLimitingGBPS}}GBps</p>{{end}}
		{{if .RateLimitingTPS}}<h3>Transactions per Second Limit</h3>
		<p>{{.RateLimitingTPS}}Tps</p>{{end}}
		{{if .OverflowService}}<h3>Overflow Service</h3>
		<p>{{.OverflowService}}</p>{{end}}
	</details>
</section>{{end}}
{{if or .HeaderRewriteEdge .HeaderRewriteMid .HeaderRewriteRedirectRouter}}<section>
	<details>
		<summary><h2>Header Customization</h2></summary>
		{{if .HeaderRewriteEdge}}<h3>Header Rewrite - Edge Tier</h3>
		<pre>{{.HeaderRewriteEdge}}</pre>{{end}}
		{{if .HeaderRewriteMid}}<h3>Header Rewrite - Mid Tier</h3>
		<pre>{{.HeaderRewriteMid}}</pre>{{end}}
		{{if .HeaderRewriteRedirectRouter}}<h3>Header Rewrite - Router</h3>
		<pre>{{.HeaderRewriteRedirectRouter}}</pre>{{end}}
	</details>
</section>{{end}}
{{if .Notes}}<section>
	<details>
		<summary><h2>Additional Notes</h2></summary>
		<p>{{.Notes}}</p>
	</details>
</section>{{end}}
</body>
</html>
`))

// IDNoMod type is used to suppress JSON unmarshalling
type IDNoMod int

// DeliveryServiceRequestRequest is a literal request to make a Delivery Service.
//
// Deprecated: Delivery Services Requests have been deprecated in favor of
// Delivery Service Requests, and will be removed from the Traffic Ops API at
// some point in the future.
type DeliveryServiceRequestRequest struct {
	// EmailTo is the email address that is ultimately the destination of a formatted DeliveryServiceRequestRequest.
	EmailTo string `json:"emailTo"`
	// Details holds the actual request in a data structure.
	Details DeliveryServiceRequestDetails `json:"details"`
}

// DeliveryServiceRequestDetails holds information about what a user is trying
// to change, with respect to a delivery service.
//
// Deprecated: Delivery Services Requests have been deprecated in favor of
// Delivery Service Requests, and will be removed from the Traffic Ops API at
// some point in the future.
type DeliveryServiceRequestDetails struct {
	// ContentType is the type of content to be delivered, e.g. "static", "VOD" etc.
	ContentType string `json:"contentType"`
	// Customer is the requesting customer - typically this is a Tenant.
	Customer string `json:"customer"`
	// DeepCachingType represents whether or not the Delivery Service should use Deep Caching.
	DeepCachingType *DeepCachingType `json:"deepCachingType"`
	// Delivery Protocol is the protocol clients should use to connect to the Delivery Service.
	DeliveryProtocol *Protocol `json:"deliveryProtocol"`
	// HasNegativeCachingCustomization indicates whether or not the resulting Delivery Service should
	// customize the use of negative caching. When this is `true`, NegativeCachingCustomizationNote
	// should be consulted for instructions on the customization.
	HasNegativeCachingCustomization *bool `json:"hasNegativeCachingCustomization"`
	// HasOriginACLWhitelist indicates whether or not the Origin has an ACL whitelist. When this is
	// `true`, Notes should ideally contain the actual whitelist (or viewing instructions).
	HasOriginACLWhitelist *bool `json:"hasOriginACLWhitelist"`
	// Has OriginDynamicRemap indicates whether or not the OriginURL can dynamically map to multiple
	// different actual origin servers.
	HasOriginDynamicRemap *bool `json:"hasOriginDynamicRemap"`
	// HasSignedURLs indicates whether or not the resulting Delivery Service should sign its URLs.
	HasSignedURLs *bool `json:"hasSignedURLs"`
	// HeaderRewriteEdge is an optional HeaderRewrite rule to apply at the Edge tier.
	HeaderRewriteEdge *string `json:"headerRewriteEdge"`
	// HeaderRewriteMid is an optional HeaderRewrite rule to apply at the Mid tier.
	HeaderRewriteMid *string `json:"headerRewriteMid"`
	// HeaderRewriteRedirectRouter is an optional HeaderRewrite rule to apply at routing time by
	// the Traffic Router.
	HeaderRewriteRedirectRouter *string `json:"headerRewriteRedirectRouter"`
	// MaxLibrarySizeEstimate is an estimation of the total size of content that will be delivered
	// through the resulting Delivery Service.
	MaxLibrarySizeEstimate string `json:"maxLibrarySizeEstimate"`
	// NegativeCachingCustomizationNote is an optional note describing the customization to be
	// applied to Negative Caching. This should never be `nil` (or empty) if
	// HasNegativeCachingCustomization is `true`, but in that case the recipient ought to contact
	// Customer for instructions.
	NegativeCachingCustomizationNote *string `json:"negativeCachingCustomizationNote"`
	// Notes is an optional set of extra information supplied to describe the requested Delivery
	// Service.
	Notes *string `json:"notes"`
	// OriginHeaders is an optional list of HTTP headers that must be sent in requests to the Origin. When
	// parsing from JSON, this field can be either an actual array of headers, or a string containing
	// a comma-delimited list of said headers.
	OriginHeaders *OriginHeaders `json:"originHeaders"`
	// OriginTestFile is the path to a file on the origin that can be requested to test the server's
	// operational readiness, e.g. '/test.xml'.
	OriginTestFile string `json:"originTestFile"`
	// OriginURL is the URL of the origin server that has the content to be served by the requested
	// Delivery Service.
	OriginURL string `json:"originURL"`
	// OtherOriginSecurity is an optional note about any and all other Security employed by the origin
	// server (beyond an ACL whitelist, which has its own field: HasOriginACLWhitelist).
	OtherOriginSecurity *string `json:"otherOriginSecurity"`
	// OverflowService is an optional IP Address or URL to which clients should be redirected when
	// the requested Delivery Service exceeds its operational capacity.
	OverflowService *string `json:"overflowService"`
	// PeakBPSEstimate is an estimate of the bytes per second expected at peak operation.
	PeakBPSEstimate string `json:"peakBPSEstimate"`
	// PeakTPSEstimate is an estimate of the transactions per second expected at peak operation.
	PeakTPSEstimate string `json:"peakTPSEstimate"`
	// QueryStringHandling describes the manner in which the CDN should handle query strings in client
	// requests. Generally one of "use", "drop", or "ignore-in-cache-key-and-pass-up".
	QueryStringHandling string `json:"queryStringHandling"`
	// RangeRequestHandling describes the manner in which HTTP requests are handled.
	RangeRequestHandling string `json:"rangeRequestHandling"`
	// RateLimitingGBPS is an optional rate limit for the requested Delivery Service in gigabytes per
	// second.
	RateLimitingGBPS *uint `json:"rateLimitingGBPS"`
	// RateLimitingTPS is an optional rate limit for the requested Delivery Service in transactions
	// per second.
	RateLimitingTPS *uint `json:"rateLimitingTPS"`
	// RoutingName is the top-level DNS label under which the Delivery Service should be requested.
	RoutingName string `json:"routingName"`
	// RoutingType is the type of routing Traffic Router should perform for the requested Delivery
	// Service.
	RoutingType *DSType `json:"routingType"`
	// ServiceAliases is an optional list of alternative names for the requested Delivery Service.
	ServiceAliases []string `json:"serviceAliases"`
	// ServiceDesc is a basic description of the requested Delivery Service.
	ServiceDesc string `json:"serviceDesc"`
}

// Format formats the DeliveryServiceRequestDetails into the text/html body of an email. The template
// used is EmailTemplate.
func (d DeliveryServiceRequestDetails) Format() (string, error) {
	b := &strings.Builder{}

	if err := EmailTemplate.Execute(b, d); err != nil {
		return "", fmt.Errorf("Failed to apply template: %w", err)
	}
	return b.String(), nil
}

// Validate validates that the delivery service request has all of the required fields. In some cases,
// e.g. the top-level EmailTo field, the format is also checked for correctness.
func (d *DeliveryServiceRequestRequest) Validate() error {
	errs := make([]error, 0, 2)

	err := validation.ValidateStruct(d,
		validation.Field(&d.EmailTo, validation.Required, is.Email),
	)
	if err != nil {
		errs = append(errs, err)
	}

	details := d.Details
	err = validation.ValidateStruct(&details,
		validation.Field(&details.ContentType, validation.Required),
		validation.Field(&details.Customer, validation.Required, validation.Match(regexp.MustCompile(`^[\w@!#$%^&\*\(\)\[\]\. -]+$`))),
		validation.Field(&details.DeepCachingType, validation.By(
			func(t interface{}) error {
				if t != (*DeepCachingType)(nil) && *t.(*DeepCachingType) == DeepCachingTypeInvalid {
					return errors.New("deepCachingType: invalid Deep Caching Type")
				}
				return nil
			})),
		validation.Field(&details.DeliveryProtocol, validation.By(
			func(p interface{}) error {
				if p == (*Protocol)(nil) {
					return errors.New("deliveryProtocol: required")
				}
				if *p.(*Protocol) == ProtocolInvalid {
					return errors.New("deliveryProtocol: invalid Protocol")
				}
				return nil
			})),
		validation.Field(&details.HasNegativeCachingCustomization, validation.By(
			func(h interface{}) error {
				if h == (*bool)(nil) {
					return errors.New("hasNegativeCachingCustomization: required")
				}
				return nil
			})),
		validation.Field(&details.HasOriginACLWhitelist, validation.By(
			func(h interface{}) error {
				if h == (*bool)(nil) {
					return errors.New("hasNegativeCachingCustomization: required")
				}
				return nil
			})),
		validation.Field(&details.HasOriginDynamicRemap, validation.By(
			func(h interface{}) error {
				if h == (*bool)(nil) {
					return errors.New("hasNegativeCachingCustomization: required")
				}
				return nil
			})),
		validation.Field(&details.HasSignedURLs, validation.By(
			func(h interface{}) error {
				if h == (*bool)(nil) {
					return errors.New("hasNegativeCachingCustomization: required")
				}
				return nil
			})),
		validation.Field(&details.MaxLibrarySizeEstimate, validation.Required),
		validation.Field(&details.OriginHeaders, validation.By(
			func(h interface{}) error {
				if h == (*OriginHeaders)(nil) {
					return nil
				}
				if len(*h.(*OriginHeaders)) < 1 {
					return errors.New("originHeaders: cannot be an empty list (use 'null' if none)")
				}
				return nil
			})),
		validation.Field(&details.OriginTestFile, validation.Required),
		validation.Field(&details.OriginURL, validation.Required, is.URL),
		validation.Field(&details.PeakBPSEstimate, validation.Required),
		validation.Field(&details.PeakTPSEstimate, validation.Required),
		validation.Field(&details.QueryStringHandling, validation.Required),
		validation.Field(&details.RangeRequestHandling, validation.Required),
		validation.Field(&details.RoutingType, validation.By(
			func(t interface{}) error {
				if t == (*DSType)(nil) || *(t.(*DSType)) == "" {
					return errors.New("routingType: required")
				}
				*t.(*DSType) = DSTypeFromString(string(*t.(*DSType)))
				if *t.(*DSType) == DSTypeInvalid {
					return errors.New("routingType: invalid Routing Type")
				}
				return nil
			})),
		validation.Field(&details.ServiceDesc, validation.Required),
	)

	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}

// OriginHeaders represents a list of the headers that must be sent to the Origin.
type OriginHeaders []string

// UnmarshalJSON implements the json.Unmarshaler interface.
func (o *OriginHeaders) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = OriginHeaders([]string{})
		return nil
	}

	headers := []string{}
	if err := json.Unmarshal(data, headers); err == nil {
		*o = OriginHeaders(headers)
		return nil
	}

	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("%s does not represent Origin Headers: %v", string(data), err)
	}

	*o = OriginHeaders(strings.Split(s, ","))
	return nil
}

// DeliveryServiceRequest is used as part of the workflow to create,
// modify, or delete a delivery service.
type DeliveryServiceRequest struct {
	AssigneeID      int             `json:"assigneeId,omitempty"`
	Assignee        string          `json:"assignee,omitempty"`
	AuthorID        IDNoMod         `json:"authorId"`
	Author          string          `json:"author"`
	ChangeType      string          `json:"changeType"`
	CreatedAt       *TimeNoMod      `json:"createdAt"`
	ID              int             `json:"id"`
	LastEditedBy    string          `json:"lastEditedBy,omitempty"`
	LastEditedByID  IDNoMod         `json:"lastEditedById,omitempty"`
	LastUpdated     *TimeNoMod      `json:"lastUpdated"`
	DeliveryService DeliveryService `json:"deliveryService"` // TODO version DeliveryServiceRequest
	Status          RequestStatus   `json:"status"`
	XMLID           string          `json:"-" db:"xml_id"`
}

// DeliveryServiceRequestNullable is used as part of the workflow to create,
// modify, or delete a delivery service.
type DeliveryServiceRequestNullable struct {
	AssigneeID      *int                        `json:"assigneeId,omitempty" db:"assignee_id"`
	Assignee        *string                     `json:"assignee,omitempty"`
	AuthorID        *IDNoMod                    `json:"authorId" db:"author_id"`
	Author          *string                     `json:"author"`
	ChangeType      *string                     `json:"changeType" db:"change_type"`
	CreatedAt       *TimeNoMod                  `json:"createdAt" db:"created_at"`
	ID              *int                        `json:"id" db:"id"`
	LastEditedBy    *string                     `json:"lastEditedBy"`
	LastEditedByID  *IDNoMod                    `json:"lastEditedById" db:"last_edited_by_id"`
	LastUpdated     *TimeNoMod                  `json:"lastUpdated" db:"last_updated"`
	DeliveryService *DeliveryServiceNullableV30 `json:"deliveryService" db:"deliveryservice"`
	Status          *RequestStatus              `json:"status" db:"status"`
	XMLID           *string                     `json:"-" db:"xml_id"`
}

// UnmarshalJSON implements the json.Unmarshaller interface to suppress unmarshalling for IDNoMod
func (a *IDNoMod) UnmarshalJSON([]byte) error {
	return nil
}

// RequestStatus captures where in the workflow this request is
type RequestStatus string

const (
	// RequestStatusInvalid -- invalid state
	RequestStatusInvalid = RequestStatus("invalid")
	// RequestStatusDraft -- newly created; not ready to be reviewed
	RequestStatusDraft = RequestStatus("draft")
	// RequestStatusSubmitted -- newly created; ready to be reviewed
	RequestStatusSubmitted = RequestStatus("submitted")
	// RequestStatusRejected -- reviewed, but problems found
	RequestStatusRejected = RequestStatus("rejected")
	// RequestStatusPending -- reviewed and locked; ready to be implemented
	RequestStatusPending = RequestStatus("pending")
	// RequestStatusComplete -- implemented and locked
	RequestStatusComplete = RequestStatus("complete")
)

// RequestStatuses -- user-visible string associated with each of the above
var RequestStatuses = []RequestStatus{
	// "invalid" -- don't list here..
	"draft",
	"submitted",
	"rejected",
	"pending",
	"complete",
}

// UnmarshalJSON implements json.Unmarshaller
func (r *RequestStatus) UnmarshalJSON(b []byte) error {
	u, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	// just check to see if the string represents a valid requeststatus
	_, err = RequestStatusFromString(u)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, (*string)(r))
}

// MarshalJSON implements json.Marshaller
func (r RequestStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(r))
}

// Value implements driver.Valuer
func (r *RequestStatus) Value() (driver.Value, error) {
	v, err := json.Marshal(r)
	log.Debugf("value is %v; err is %v", v, err)
	v = []byte(strings.Trim(string(v), `"`))
	return v, err
}

// Scan implements sql.Scanner
func (r *RequestStatus) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected requeststatus in byte array form; got %T", src)
	}
	b = []byte(`"` + string(b) + `"`)
	return json.Unmarshal(b, r)
}

// RequestStatusFromString gets the status enumeration from a string
func RequestStatusFromString(rs string) (RequestStatus, error) {
	if rs == "" {
		return RequestStatusDraft, nil
	}
	for _, s := range RequestStatuses {
		if string(s) == rs {
			return s, nil
		}
	}
	return RequestStatusInvalid, errors.New(rs + " is not a valid RequestStatus name")
}

// ValidTransition returns nil if the transition is allowed for the workflow, an error if not
func (r RequestStatus) ValidTransition(to RequestStatus) error {
	if r == RequestStatusRejected || r == RequestStatusComplete {
		// once rejected or completed,  no changes allowed
		return errors.New(string(r) + " request cannot be changed")
	}

	if r == to {
		// no change -- always allowed
		return nil
	}

	// indicate if valid transitioning to this RequestStatus
	switch to {
	case RequestStatusDraft:
		// can go back to draft if submitted or rejected
		if r == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusSubmitted:
		// can go be submitted if draft or rejected
		if r == RequestStatusDraft {
			return nil
		}
	case RequestStatusRejected:
		// only submitted can be rejected
		if r == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusPending:
		// only submitted can move to pending
		if r == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusComplete:
		// only submitted or pending requests can be completed
		if r == RequestStatusSubmitted || r == RequestStatusPending {
			return nil
		}
	}
	return errors.New("invalid transition from " + string(r) + " to " + string(to))
}
