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

import "fmt"
import "testing"

const TEST_MIME = "tEXt/vND.plaIN+cRLf;    charset=utf-8; q=2.2;foo=bar"
const PRETTY_MIME = "text/vnd.plain+crlf; charset=utf-8; foo=bar; q=2.2"

func ExampleNewMimeType() {
	m, err := NewMimeType("text/plain;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Name:", m.Name, "Parameters:", m.Parameters)
	// Output: Name: text/plain Parameters: map[charset:utf-8]
}

func ExampleMimeType_Quality() {
	m, err := NewMimeType("text/plain;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	q := m.Quality()
	fmt.Print(q)

	m.Parameters["q"] = "0.9"
	q = m.Quality()
	fmt.Println("", q)

	// Output: 1 0.9
}

func ExampleMimeType_Charset() {
	m, err := NewMimeType("text/plain;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c := m.Charset() // it's okay to ignore this error, because I was a good boy and used NewMimeType
	fmt.Println(c)
	// Output: utf-8
}

func ExampleMimeType_Type() {
	m, err := NewMimeType("text/plain;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m.Type())
	// Output: text
}

func ExampleMimeType_SubType() {
	m, err := NewMimeType("text/vnd.plain+crlf;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m.SubType())
	// Output: vnd.plain+crlf
}

func ExampleMimeType_Facet() {
	m, err := NewMimeType("text/vnd.plain+crlf;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m.Facet())
	// Output: vnd
}

func ExampleMimeType_Syntax() {
	m, err := NewMimeType("text/vnd.plain+crlf;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m.Syntax())
	// Output: crlf
}

func ExampleMimeType_String() {
	m, err := NewMimeType("text/plain;foo=bar;charset=utf-8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(m)
	// Output: text/plain; charset=utf-8; foo=bar
}

func ExampleMimeType_Satisfy() {
	m, err := NewMimeType("text/plain")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	o, err := NewMimeType("text/*")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(m.Satisfy(o), o.Satisfy(m))
	// Output: true false
}

func ExampleMimeType_Less() {
	one, err := NewMimeType("text/plain")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	two, err := NewMimeType("text/*")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	three, err := NewMimeType("text/plain;q=0.9")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(one.Less(two), two.Less(three), one.Less(three))
	// Output: false false false
}

func ExampleMimeTypesFromAccept() {
	const acceptLine = "text/html,text/xml;q=0.9,text/*;q=0.9,*/*"
	mimes, err := MimeTypesFromAccept(acceptLine)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, m := range mimes {
		fmt.Printf("%s, ", m)
	}
	fmt.Println()
	// Output: text/html, */*, text/xml; q=0.9, text/*; q=0.9,
}

func ExampleSortMimeTypes() {
	// Normally don't do this, but for the sake of brevity in an example I will
	mimes := []MimeType{
		MimeType{
			"text/html",
			map[string]string{},
		},
		MimeType{
			"text/xml",
			map[string]string{"q": "0.9"},
		},
		MimeType{
			"text/*",
			map[string]string{"q": "0.9"},
		},
		MimeType{
			"*/*",
			map[string]string{},
		},
	}

	SortMimeTypes(mimes)
	for _, m := range mimes {
		fmt.Printf("%s, ", m)
	}
	fmt.Println()
	// Output: text/html, */*, text/xml; q=0.9, text/*; q=0.9,
}

func TestMimeType(t *testing.T) {
	m, err := NewMimeType(TEST_MIME)
	if err != nil {
		t.Fatalf("Failed to construct a MimeType from TEST_MIME: %v", err)
	}

	if m.Name != "text/vnd.plain+crlf" {
		t.Errorf("Incorrect MIME name, expected 'text/vnd.plain+crld' but got '%s'", m.Name)
	}

	if m.Type() != "text" {
		t.Errorf("Incorrect 'Type', expected 'text' but got '%s'", m.Type())
	}

	if m.SubType() != "vnd.plain+crlf" {
		t.Errorf("Incorrect 'SubType', expected 'vnd.plain+crlf' but got '%s'", m.SubType())
	}

	if m.Facet() != "vnd" {
		t.Errorf("Incorrect 'Facet', expected 'vnd' but got '%s'", m.Facet())
	}

	if m.Syntax() != "crlf" {
		t.Errorf("Incorrect 'syntax suffix', expected 'crlf' but got '%s'", m.Syntax())
	}

	if m.String() != PRETTY_MIME {
		t.Errorf("Incorrect string representation, expected '%s' but got '%s'", PRETTY_MIME, m.String())
	}

	if len(m.Parameters) != 3 {
		t.Errorf("Incorrect number of Parameters, expected 3 but got %d", len(m.Parameters))
	}

	if q := m.Quality(); q != 2.2 {
		t.Errorf("Incorrect quality, expected 2.2 but got %g", q)
	}

	if c := m.Charset(); c != "utf-8" {
		t.Errorf("Incorrect charset, expected 'utf-8', but got '%s'", c)
	}

}

func TestMimeType_Satisfy(t *testing.T) {
	m, err := NewMimeType(TEST_MIME)
	if err != nil {
		t.Fatalf("Failed to construct a MimeType from TEST_MIME: %v", err)
	}

	o, err := NewMimeType("*/*")
	if err != nil {
		t.Fatalf("Failed to construct a MimeType from '*/*': %v", err)
	}

	if !m.Satisfy(o) {
		t.Errorf("Expected %s to satisfy %s, but it did not", m, o)
	}

	if o.Satisfy(m) {
		t.Errorf("Expected %s to not satisfy %s, but it did", o, m)
	}

	if o, err = NewMimeType("text/*"); err != nil {
		t.Fatalf("Failed to construct a MimeType from 'text/*': %v", err)
	}

	if !m.Satisfy(o) {
		t.Errorf("Expected %s to satisfy %s, but it did not", m, o)
	}

	if o.Satisfy(m) {
		t.Errorf("Expected %s to not satisfy %s, but it did", o, m)
	}

	if o, err = NewMimeType("text/vnd.plain+crlf;q=2.1"); err != nil {
		t.Fatalf("Failed to construct a MimeType from 'text/vnd.plain+crlf;q=2.1': %v", err)
	}

	if !m.Satisfy(o) {
		t.Errorf("Expected %s to satisfy %s, but it did not", m, o)
	}

	if o.Satisfy(m) {
		t.Errorf("Expected %s to not satisfy %s, but it did", o, m)
	}
}

func TestMimeType_Less(t *testing.T) {
	m, err := NewMimeType("text/*")
	if err != nil {
		t.Fatalf("Failed to construct MimeType from 'text/*': %v", err)
	}

	o, err := NewMimeType("*/*")
	if err != nil {
		t.Fatalf("Failed to construct MimeType from '*/*': %v", err)
	}

	const less = "Expected %s to be less than %s, but it was not"
	const notLess = "Expected %s to not be less than %s, but it was"

	if m.Less(o) {
		t.Errorf(notLess, m, o)
	}

	if !o.Less(m) {
		t.Errorf(less, o, m)
	}

	if o, err = NewMimeType("text/plain"); err != nil {
		t.Fatalf("Failed to construct MimeType from 'text/plain': %v", err)
	}

	if !m.Less(o) {
		t.Errorf(less, m, o)
	}

	if o.Less(m) {
		t.Errorf(notLess, o, m)
	}

	if m, err = NewMimeType("text/plain;foo=bar;q=1.0"); err != nil {
		t.Fatalf("Failed to construct MimeType from 'text/plain;foo=bar;q=1.0': %v", err)
	}

	if m.Less(o) {
		t.Errorf(notLess, m, o)
	}

	if !o.Less(m) {
		t.Errorf(less, o, m)
	}

	if o, err = NewMimeType("text/plain;q=1.1"); err != nil {
		t.Fatalf("Failed to construct MimeType from 'text/plain;q=1.1': %v", err)
	}

	if !m.Less(o) {
		t.Errorf(less, m, o)
	}

	if o.Less(m) {
		t.Errorf(notLess, o, m)
	}

	if o, err = NewMimeType("text/plain;fizz=buzz;q=1.0"); err != nil {
		t.Fatalf("Failed to construct MimeType from 'text/plain;fizz=buzz;q=1.0': %v", err)
	}

	if m.Less(o) {
		t.Errorf(notLess, m, o)
	}

	if o.Less(m) {
		t.Errorf(notLess, o, m)
	}

}
