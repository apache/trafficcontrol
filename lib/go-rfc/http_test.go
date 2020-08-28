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
import "net/http"

func ExampleGetHTTPDate() {
	hdrs := http.Header{}
	hdrs.Set("Date", "Tue, 30 Jun 2020 19:07:15 GMT")

	t, ok := GetHTTPDate(hdrs, "Date")
	if !ok {
		fmt.Println("Failed to get date")
		return
	}

	fmt.Printf("%s\n", t)

	// Output: 2020-06-30 19:07:15 +0000 GMT
}

func ExampleParseHTTPDate() {
	t, ok := ParseHTTPDate("Tue, 30 Jun 2020 19:07:15 GMT")
	if !ok {
		fmt.Println("Failed to parse date")
		return
	}

	fmt.Printf("%s\n", t)

	// Output: 2020-06-30 19:07:15 +0000 GMT
}
