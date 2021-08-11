package rfc

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
	"mime"
	"sort"
	"strconv"
	"strings"
)

/*
MimeType represents a "Media Type" as defined by RFC6838, along with some ease-of-use functionality.

Note that this structure is in no way guaranteed to represent a *real* MIME Type, only one that
is *syntactically valid*. The hope is that it will make content negotiation easier for developers,
and should not be considered a security measure by any standard.
*/
type MimeType struct {
	// Name is the full name of the MIME Type, e.g. 'application/json'.
	// Usually for printing, it's better to call MimeType.String
	Name string
	// Parameters contains a map of provided parameter names to corresponding values. Note that for
	// MimeTypes constructed with NewMimeType, this will always be initialized, even when empty.
	Parameters map[string]string
}

/*
Quality retrieves and parses the "quality" parameter of a MimeType.

As specified in RFC7231, the quality parameter's name is "q", not actually "quality". To obtain
a literal "quality" parameter value, access MimeType.Parameters directly.
MimeTypes with no "q" parameter implicitly have a "quality" of 1.0.
*/
func (m MimeType) Quality() float64 {
	if m.Parameters == nil {
		return 1
	}

	fs, ok := m.Parameters["q"]
	if !ok {
		return 1
	}

	ret, err := strconv.ParseFloat(fs, 64)
	if err != nil {
		return 1
	}
	return ret
}

/*
Charset retrieves the "charset" parameter of a MimeType.

Returns an empty string if no charset exists in the parameters, or if the parameters themselves are
not initialized.
*/
func (m MimeType) Charset() string {
	if m.Parameters == nil {
		return ""
	}

	c, ok := m.Parameters["charset"]
	if !ok {
		return ""
	}
	return c
}

// Type returns only the "main" type of a MimeType.
func (m MimeType) Type() string {
	return strings.SplitN(m.Name, "/", 2)[0]
}

// SubType returns only the "sub" type of a MimeType.
func (m MimeType) SubType() string {
	s := strings.SplitN(m.Name, "/", 2)
	if len(s) != 2 {
		return ""
	}
	return s[1]
}

// Facet returns the MimeType's "facet" if one exists, otherwise an empty string.
func (m MimeType) Facet() string {
	s := m.SubType()
	if fx := strings.SplitN(s, ".", 2); len(fx) == 2 {
		return fx[0]
	}
	return ""
}

// Syntax returns the MimeType's "syntax suffix" if one exists, otherwise an empty string.
func (m MimeType) Syntax() string {
	s := m.SubType()
	if fx := strings.Split(s, "+"); len(fx) > 1 {
		return fx[len(fx)-1]
	}
	return ""
}

// String implements the Stringer interface using mime.FormatMediaType.
func (m MimeType) String() string {
	return mime.FormatMediaType(m.Name, m.Parameters)
}

// Satisfy checks whether or not the MimeType "satisfies" some other MimeType, o.
//
// Note that this does not check if the two are literally the *same*. Specifically, if the Type or
// SubType of the given MimeType o is the special '*' name, then this will instead check whether or
// not this MimeType can *satisfy* the other according to RFC7231. This means that this satisfaction
// check is NOT associative - that is a.Satisfy(b) does not imply b.Satisfy(a).
//
// See Also: The MDN documentation on the Accept Header: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept
func (m MimeType) Satisfy(o MimeType) bool {
	// literally anything will satisfy this
	if o.Type() == "*" && o.SubType() == "*" {
		return true
	}

	// it's not syntactically valid to have a "*" type and a non-"*" subtype, e.g. "*/foo", so we're
	// done here
	if o.Type() != "*" && o.SubType() == "*" {
		return o.Type() == m.Type()
	}

	if o.Type() != m.Type() || o.SubType() != m.SubType() {
		return false
	}

	for k, v := range o.Parameters {
		if k == "q" {
			continue
		}

		if mv, ok := m.Parameters[k]; !ok || mv != v {
			return false
		}
	}

	return true
}

/*
Less checks whether or not this MimeType is "less than" some other MimeType, o.

This is done using a comparison of "quality value" of the two MimeTypes, as specified in RFC7231.

See Also: The MDN documentation on "quality value" comparisons: https://developer.mozilla.org/en-US/docs/Glossary/Quality_values
*/
func (m MimeType) Less(o MimeType) bool {
	mq := m.Quality()
	oq := o.Quality()

	if mq < oq {
		return true
	} else if oq < mq {
		return false
	}

	mType := m.Type()
	mSub := m.SubType()
	oType := o.Type()
	oSub := o.SubType()

	if mSub == "*" {
		if oSub != "*" {
			return true
		}

		if mType == "*" {
			if oType != "*" {
				return true
			}
		} else if oType == "*" {
			return false
		}

	} else if oSub == "*" {
		return false
	}

	return len(m.Parameters) < len(o.Parameters)
}

/*
NewMimeType creates a new MimeType, initializing its Parameters map.

If v cannot be parsed as a valid MIME (by mime.ParseMediaType), returns an error.
*/
func NewMimeType(v string) (m MimeType, err error) {
	m.Name, m.Parameters, err = mime.ParseMediaType(v)
	return m, err
}

/*
MimeTypesFromAccept constructs a *sorted* list of MimeTypes from the provided text, which is assumed
to be from an HTTP 'Accept' header. The list is sorted using to SortMimeTypes.

If a is an empty string, this will return an empty slice and no error.
*/
func MimeTypesFromAccept(a string) ([]MimeType, error) {
	mimes := []MimeType{}
	if a == "" {
		return mimes, nil
	}

	for _, raw := range strings.Split(a, ",") {
		m, err := NewMimeType(raw)
		if err != nil {
			return mimes, err
		}
		mimes = append(mimes, m)
	}
	SortMimeTypes(mimes)
	return mimes, nil
}

/*
SortMimeTypes sorts the passed MimeTypes according to their "quality value". See MimeType.Less for
more information on how MimeTypes are compared.
*/
func SortMimeTypes(m []MimeType) {
	// using !Less because default sort order is ascending
	sort.SliceStable(m, func(i, j int) bool { return !m[i].Less(m[j]) })
}

// MIME_JSON is a pre-defined MimeType for JSON data.
var MIME_JSON = MimeType{
	Name:       "application/json",
	Parameters: map[string]string{},
}

// MIME_PLAINTEXT is a pre-defined MimeType for plain text data.
var MIME_PLAINTEXT = MimeType{
	Name:       "text/plain",
	Parameters: map[string]string{"charset": "utf-8"},
}

// MIME_HTML is a pre-defined MimeType for HTML data.
var MIME_HTML = MimeType{
	Name:       "text/html",
	Parameters: map[string]string{"charset": "utf-8"},
}

// MIME_CSS is a pre-defined MimeType for CSS data.
var MIME_CSS = MimeType{
	Name:       "text/css",
	Parameters: map[string]string{"charset": "utf-8"},
}

// MIME_JS is a pre-defined MimeType for JavaScript data.
var MIME_JS = MimeType{
	Name:       "text/javascript",
	Parameters: map[string]string{"charset": "utf-8"},
}
