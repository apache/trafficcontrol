package test

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
	"fmt"
	"math/rand"
	"net"
)

func RandStr() *string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	num := 100
	s := ""
	for i := 0; i < num; i++ {
		s += string(chars[rand.Intn(len(chars))])
	}
	return &s
}
func RandStrArray() []string {
	num := 100
	sArray := make([]string, num)
	for i := 0; i < num; i++ {
		sArray[i] = *RandStr()
	}
	return sArray
}

func RandBool() *bool {
	b := rand.Int()%2 == 0
	return &b
}

func RandInt() *int {
	i := rand.Int()
	return &i
}
func RandInt64() *int64 {
	i := rand.Int63()
	return &i
}

func RandUint64() *uint64 {
	i := uint64(rand.Int63())
	return &i
}

func RandUint() *uint {
	i := uint(rand.Int())
	return &i
}

func RandIntn(n int) *int {
	i := rand.Intn(n)
	return &i
}

func RandFloat64() *float64 {
	f := rand.Float64()
	return &f
}

func RandomIPv4() *string {
	first := rand.Int31n(256)
	second := rand.Int31n(256)
	third := rand.Int31n(256)
	fourth := rand.Int31n(256)
	str := fmt.Sprintf("%d.%d.%d.%d", first, second, third, fourth)
	return &str
}

func RandomIPv6() *string {
	ip := net.IP([]byte{
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
	}).String()
	return &ip
}
