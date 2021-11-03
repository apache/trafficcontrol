package util

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
	"fmt"
)

func ExampleJoinErrsStr() {
	errs := []error{
		errors.New("test"),
		errors.New("quest"),
	}

	fmt.Println(JoinErrsStr(errs))
	fmt.Println(JoinErrsStr(nil))
	// Output: test, quest
	//

}

func ExampleErrsToStrs() {
	errs := []error{
		errors.New("test"),
		errors.New("quest"),
	}
	strs := ErrsToStrs(errs)
	fmt.Println(strs[0])
	fmt.Println(strs[1])
	// Output: test
	// quest
}

func ExampleJoinErrsSep() {
	errs := []error{
		errors.New("test"),
		errors.New("quest"),
	}

	fmt.Println(JoinErrsSep(errs, "\n"))

	// Output: test
	// quest
}

func ExampleCamelToSnakeCase() {
	camel := "camelCase"
	fmt.Println(CamelToSnakeCase(camel))
	camel = "PascalCase"
	fmt.Println(CamelToSnakeCase(camel))
	camel = "IPIsAnInitialismForInternetProtocol"
	fmt.Println(CamelToSnakeCase(camel))

	// Output: camel_case
	// pascal_case
	// ipis_an_initialism_for_internet_protocol
}
