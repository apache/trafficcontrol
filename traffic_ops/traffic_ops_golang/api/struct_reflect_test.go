package api

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
	"testing"
)

type Dog struct {
	Name string
}

func (d Dog) String() string {
	return "bark!"
}

type IdentifierImpl struct {
	ID int `info:"key"`
}

type Person struct {
	Age  int
	Name *string `info:"test"`

	IdentifierImpl
	Pet *Dog
}

type PrivatePerson struct {
	name string
	age  int
}

// good for debugging
func printArgs(dest ...interface{}) error {
	for _, elem := range dest {

		switch v := elem.(type) {
		case *int:
			fmt.Printf("%v -> %v\n", v, *v)
		case *string:
			fmt.Printf("%v -> %v\n", v, *v)
		case **int:
			fmt.Printf("%v -> -> %v\n", v, **v)
		case **string:
			fmt.Printf("%v -> -> %v\n", v, **v)
		default:
			return fmt.Errorf("unknown type")
		}
	}
	return nil
}

// TODO:
//	Test Scanner
//	Test Double Pointer
//	Test Non-pointer
//	Test Pointer to item
//	Test Not All Fields exported
//	Test MustGetFieldInfo? no..
// Current coverage: 78.2% (this file)

// Test cases:
// case 0:	S   (EmbeddedStruct)
// case 1:
//		a)	*S  (EmbeddedStructPtr)
//		b)	**T (UnsupportedField)
//		c)	*T  (NullableStructField)
// case 2:	T   (NonNullableStructField)
//
func TestGetFieldInfo(t *testing.T) {

	name := "Matthew"
	newName := "Matt"

	key := IdentifierImpl{
		ID: 118908,
	}

	p := &Person{
		Age:            20,
		Name:           &name,
		IdentifierImpl: key,
		Pet:            &Dog{Name: "Daisy"},
	}

	fields, err := GetFieldInfo(p)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	ptrs := fields[AllFields]

	// Changing the struct will update `ptrs`
	p.Age = 21
	p.Name = &newName

	// Increment all ints by 1 (via ptrs)
	for _, elem := range ptrs {
		if intPtr, ok := elem.(*int); ok {
			*intPtr = *intPtr + 1
			continue
		}
		if intPtr, ok := elem.(**int); ok {
			**intPtr = **intPtr + 1
		}
	}

	if p.Age != 22 {
		t.Errorf("Person does not have expected age!")
	}

	// TODO: Change this..?
	if **fields["test"][0].(**string) != newName {
		t.Errorf("Person does not have expected name!")
	}

}
