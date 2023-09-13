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

// RandStr returns, as a string, a random 100 character alphanumeric, including '-' and '_'.
func RandStr() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	num := 100
	s := ""
	for i := 0; i < num; i++ {
		s += string(chars[rand.Intn(len(chars))])
	}
	return s
}

// RandStrArray returns, as an array of 100 strings, with 100 character random alphanumeric, including '-' and '_'.
func RandStrArray() []string {
	num := 100
	sArray := make([]string, num)
	for i := 0; i < num; i++ {
		sArray[i] = RandStr()
	}
	return sArray
}

// RandBool returns a random boolean value.
func RandBool() bool {
	b := rand.Int()%2 == 0
	return b
}

// RandInt returns a random int.
func RandInt() int {
	i := rand.Int()
	return i
}

// RandIntForActive returns a random int 0-5
func RandIntForActive() int {
	mini := 0
	maxi := 5
	return rand.Intn(maxi-mini) + mini
}

// RandInt64 returns a random signed 64-bit int.
func RandInt64() int64 {
	i := rand.Int63()
	return i
}

// RandUint64 returns a random unsigned 64-bit int.
func RandUint64() uint64 {
	i := uint64(rand.Int63())
	return i
}

// RandUint returns a random unsigned int.
func RandUint() uint {
	i := uint(rand.Int())
	return i
}

// RandIntn returns a random int in the half-open interval [0,n).
func RandIntn(n int) int {
	i := rand.Intn(n)
	return i
}

// RandFloat64 returns a random float64.
func RandFloat64() float64 {
	f := rand.Float64()
	return f
}

// RandomIPv4 returns, as a string, a random IP address.
func RandomIPv4() string {
	first := rand.Int31n(256)
	second := rand.Int31n(256)
	third := rand.Int31n(256)
	fourth := rand.Int31n(256)
	str := fmt.Sprintf("%d.%d.%d.%d", first, second, third, fourth)
	return str
}

// RandomIPv6 returns, as a string, a random IPv6 address.
func RandomIPv6() string {
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
	return ip
}
