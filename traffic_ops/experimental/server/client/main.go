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

package main

import (
	"./client"
	"fmt"
	"os"
)

func main() {
	client, err := client.Login(os.Args[1], os.Args[2], os.Args[3], true)
	if err != nil {
		fmt.Println("err 00", err)
	}

	resp, err := client.GetText("/api/2.0/asn")
	if err != nil {
		fmt.Println("err 11:", err)
	}
	fmt.Println(resp)
	ret, e := client.PostJson("/api/2.0/asn", []byte("{\"asn\":45454, \"cachegroup\":28}"))
	if e != nil {
		fmt.Println("err 22:", e)
	}
	fmt.Println(string(ret))

	ret, e = client.PutJson("/api/2.0/asn/51", []byte("{\"asn\":33933, \"cachegroup\":28}"))
	if e != nil {
		fmt.Println("err 32:", e)
	}
	fmt.Println(string(ret))

	ret, e = client.Delete("/api/2.0/asn/52")
	if e != nil {
		fmt.Println("err 42:", e)
	}
	fmt.Println(string(ret))

	ret, e = client.Delete("/api/2.0/cachegroup/3")
	if e != nil {
		fmt.Println("err 52:", e)
	}
	fmt.Println(string(ret))
}
