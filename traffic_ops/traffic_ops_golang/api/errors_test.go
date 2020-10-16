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

import "fmt"

func ExampleErrors_String() {
	fmt.Println(NewErrors())

	// Output: Errors(Code=200, SystemError='<nil>', UserError='<nil>')
}

func ExamlpleNewErrors() {
	fmt.Println(NewErrors())
	fmt.Println(NewErrors().Occurred())

	// Output: Errors(Code=200, SystemError='<nil>', UserError='<nil>')
	// false
}

func ExampleErrors_Occurred() {
	err := NewErrors()
	fmt.Println(err.Occurred())

	err.SetSystemError("test")
	fmt.Println(err.Occurred())

	err.SetUserError("test")
	fmt.Println(err.Occurred())

	err.SystemError = nil
	fmt.Println(err.Occurred())

	// Output: false
	// true
	// true
	// true
}
