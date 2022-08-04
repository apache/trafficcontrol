// Package tc provides structures, constants, and functions that are used
// throughout the components of Apache Traffic Control.
//
// In general, the symbols defined here should be used by more than one
// component of Traffic Control - otherwise it can just appear in the only
// component that needs it. Most often this means that the symbols herein
// defined are referring to objects and/or concepts exposed through the Traffic
// Ops API, and usually serve to define payloads of HTTP requests and
// responses.
//
// # Enumerated Types
//
// Enumerated types - which should typically go in enum.go - should be treated
// as enumerables, and MUST NOT be cast as anything else (integer, strings,
// etc.). Enums MUST NOT be compared to strings or integers via casting.
// Enumerable data SHOULD be stored as the enumeration, not as a string or
// number. The *only* reason they are internally represented as strings, is to
// make them implicitly serialize to human-readable JSON. They should not be
// treated as strings. Casting or storing strings or numbers defeats the
// purpose of enum safety and conveniences.
//
// An example of an enumerated string type 'Foo' and an enumerated integral (or
// arbitrary) type 'Bar' are shown below:
//
//	type Foo string
//	const (
//	    FooA Foo = "A"
//	    FooB Foo = "B"
//	)
//
//	type Bar int
//	const (
//	    BarTest Bar = iota
//	    BarQuest
//	)
//
// Note the way each member of the "enum" is prefixed with the type name, to
// help avoid collisions. Also note how, for string enumerations, the type must
// be repeated with each assignment, whereas for an arbitrary enumeration, you
// can make use of the special 'iota' language constant to start with some
// arbitrary value and then just make names on subsequent lines to implicitly
// increment from there, while maintaining proper type.
//
// Enumerables that need to be serialized and deserialized in JSON should use
// strings to make them easiest to understand and work with (serialization will
// work out-of-the-box that way). One way to implement type safety for
// enumerated types that require serialization support is to implement the
// encoding/json.Unmarshaler interface and return an error for unsupported
// values. However, note that this causes the unmarshaling to halt immediately,
// which is fast but if there are other errors in the JSON document that would
// be encountered later, they will not be reported. Therefore, types used in
// structures that the Traffic Ops API unmarshals should generally instead use
// an "invalid" value that indicates the problem but allows parsing to
// continue. The way this is normally done is with a 'FromString' method like
// so:
//
//	type Foo string
//	const (
//	    FooA Foo = "A"
//	    FooB Foo = "B"
//	    FooInvalid = ""
//	)
//
//	func FooFromString(foo string) Foo {
//	    switch foo {
//	    case FooA:
//	        fallthrough
//	    case FooB:
//	        return Foo(foo)
//	    }
//	    return FooInvalid
//	}
//
// However, this requires postprocessing after deserialization, so one might
// instead implement encoding/json.Unmarshaler:
//
//	import "errors"
//
//	type Foo string
//	const (
//	    FooA Foo = "A"
//	    FooB Foo = "B"
//	    FooInvalid = ""
//	)
//
//	func (f *Foo) UnmarshalJSON(data []byte) error {
//	    if string(data) == "null" {
//	        return errors.New("'null' is not a valid 'Foo'")
//	    }
//	    s := strings.Trim(string(data), "\"")
//	    switch s {
//	    case FooA:
//	        fallthrough
//	    case FooB:
//	        *f = Foo(s)
//	        return
//	    }
//	    // This is an invalid *value* not a parse error, so we don't return
//	    // an error and instead just make it clear that the value was
//	    // invalid.
//	    *f = FooInvalid
//	    return nil
//	}
//
// Though in this case the original, raw, string value of the 'Foo' is lost. Be
// aware of the limitations of these two designs during implementation.
//
// When storing enumumerable data in memory, it SHOULD be converted to and
// stored as an enum via the corresponding `FromString` function, checked
// whether the conversion failed and Invalid values handled, and valid data
// stored as the enum. This guarantees stored data is valid, and catches
// invalid input as soon as possible.
//
// Conversion functions, whether they be a 'FromString' function or an
// UnmarshalJSON method, should not be case-insensitive.
//
// When adding new enum types, enums should be internally stored as strings, so
// they implicitly serialize as human-readable JSON, unless the performance or
// memory of integers is necessary (it almost certainly isn't). Enums should
// always have the "invalid" value as the empty string (or 0), so
// default-initialized enums are invalid.
//
// Enums should always have a String() method, that way they implement
// fmt.Stringer and can be easily represented in a human-readable way.
//
// # When to Use Custom Types
//
// A type should be defined whenever there is a need to represent a data
// *structure* (in a 'struct'), or whenever the type has some enforced
// semantics. For example, enumerated types have the attached semantic that
// there is a list of valid values, and generally there are methods and/or
// functions associated with those types that act on that limitation. On the
// otherhand, if there is a type 'Foo' that has some property 'Bar' that
// contains arbitrary textual data with no other semantics or limitations,
// then 'string' is the appropriate type; there is no need to make a type
// 'FooBar' to express the relationship between a Foo and its Bar - that's the
// job of the Go language itself.
//
// Similarly, try to avoid duplication. For example, Parameters can be assigned
// to Profiles or (in legacy code that may have been removed by the time this
// is being read) Cache Groups. It was not necessary to create a type
// CacheGroupParameter that contains all of the same data as a regular
// Parameter simply because of this relationship.
//
// # Versioning Structures
//
// Structures used by the Traffic Ops API may change over time as the API
// itself changes. Typically, each new major API version will cause some
// breaking change to some structure in this package. This means that a new
// structure will need to be created to avoid breaking old clients that are
// deprecated but still supported. The naming convention for a new version of a
// structure that represents a Foo object in the new version - say, 2.0 - is to
// call it FooV20. Then, a more general type alias should be made for the API
// major version for use in client methods (so that they can be silently
// upgraded for non-breaking changes) - in this case: FooV2. A deprecation
// notice should then be added to the legacy type (which is hopefully but not
// necessarily named FooV1/FooV11 etc.) indicating the new structure. For
// example:
//
//	// FooV11 represents a Foo in TO APIv1.1.
//	//
//	// Deprecated: TO APIv1.1 is deprecated; upgrade to FooV2.
//	type FooV11 struct {
//	    Bar string `json:"bar"`
//	}
//	// FooV1 represents a Foo in the latest minor version of TO APIv1.
//	//
//	// Deprecated: TO APIv1 is deprecated; upgrade to FooV2.
//	type FooV1 = FooV11
//
//	// FooV20 represents a Foo in TO APIv2.0.
//	type FooV20 struct {
//	    Bar string `json:"bar"`
//	}
//	// FooV2 represents a Foo in the latest minor version of TO APIv2.
//	type FooV2 = FooV20
//
// Note that there is no type alias simply named "Foo" - that is exactly how it
// should be.
//
// Legacy types may not include version information - newer versions should add
// it as described above. Legacy types may also include the word "Nullable"
// somewhere - strip this out! That word included in a type name is only to
// differentiate between older structures that were not "nullable" and the
// newer ones. Newly added types should only have one structure that is as
// "nullable" or "non-nullable" as it needs to be, so there is no need to
// specify.
package tc

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
