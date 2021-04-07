package auth

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
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// SCRYPTComponents the input parameters to the Scrypt encryption key format
type SCRYPTComponents struct {
	Algorithm string // The SCRYPT algorithm prefix
	N         int    // CPU/memory cost parameter (logN)
	R         int    // block size parameter (octets)
	P         int    // parallelization parameter (positive int)
	Salt      []byte // salt value
	SaltLen   int    // bytes to use as salt (octets)
	DK        []byte // derived key value
	DKLen     int    // length of the derived key (octets)
}

const KEY_DELIM = ":"

// The SCRYPT functionality defined in this package is derived based upon the following
// references:
// https://pkg.go.dev/golang.org/x/crypto/scrypt
// https://www.tarsnap.com/scrypt/scrypt.pdf
var DefaultParams = SCRYPTComponents{
	Algorithm: "SCRYPT",
	N:         16384,
	R:         8,
	P:         1,
	SaltLen:   16,
	DKLen:     64}

// DerivePassword uses the https://pkg.go.dev/golang.org/x/crypto/scrypt package to
// return an encrypted password that is compatible with the
// Perl CPAN library Crypt::ScryptKDF for backward compatibility
// to authenticate through the Perl API the same way.
// See: http://cpansearch.perl.org/src/MIK/Crypt-ScryptKDF-0.010/lib/Crypt/ScryptKDF.pm
func DerivePassword(password string) (string, error) {
	var salt []byte
	var err error
	salt, err = generateSalt(DefaultParams.DKLen)
	if err != nil {
		return "", errors.New("generating salt: " + err.Error())
	}
	key, err := scrypt.Key([]byte(password), salt, DefaultParams.N, DefaultParams.R, DefaultParams.P, DefaultParams.DKLen)
	if err != nil {
		return "", err
	}
	nStr := strconv.Itoa(DefaultParams.N)
	if err != nil {
		return "", errors.New("converting N: " + err.Error())
	}
	rStr := strconv.Itoa(DefaultParams.R)
	pStr := strconv.Itoa(DefaultParams.P)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	keyBase64 := base64.StdEncoding.EncodeToString(key)

	// The SCRYPT prefix is added because the Mojolicious Perl library adds this as a
	// prefix to every password in the database.  So it's added for compatibility.
	scryptPass := []string{DefaultParams.Algorithm, nStr, rStr, pStr, saltBase64, keyBase64}

	return strings.Join(scryptPass, KEY_DELIM), nil
}

// VerifySCRYPTPassword parses the original Derived Key (DK) from the SCRYPT password
// so that it can compare that with the password/scriptPassword param
func VerifySCRYPTPassword(password string, scryptPassword string) error {

	scomp, err := parseScrypt(scryptPassword)
	if err != nil {
		return err
	}

	keylenBytes := len(scryptPassword) - DefaultParams.DKLen
	if keylenBytes < 1 {
		return errors.New("Invalid scryptPassword length")
	}
	// scrypt the cleartext password with the same parameters and salt
	tmpDK, err := scrypt.Key([]byte(password),
		[]byte(scomp.Salt),
		scomp.N, // Must be a power of 2 greater than 1
		scomp.R,
		scomp.P, // r*p must be < 2^30
		DefaultParams.DKLen)
	if err != nil {
		return err
	}

	// Compare the Derived Key from the SCRYPT password
	if subtle.ConstantTimeCompare(scomp.DK, tmpDK) != 1 {
		return errors.New("invalid password")
	}

	return err
}

func parseScrypt(scryptPassword string) (SCRYPTComponents, error) {

	var err error
	var scomp SCRYPTComponents

	if scryptPassword == "" {
		return scomp, errors.New("scrypt password is required")
	}

	sh := strings.Split(scryptPassword, ":")

	// Algorithm
	scomp.Algorithm = sh[0]
	if scomp.Algorithm == "" {
		return scomp, errors.New("Algorithm was not defined")
	}
	if scomp.Algorithm != DefaultParams.Algorithm {
		return scomp, fmt.Errorf("Algorithm defined is not %s", DefaultParams.Algorithm)
	}

	// N
	n := sh[1]
	if n == "" {
		return scomp, errors.New("N was not defined")
	}
	var nInt int
	nInt, err = strconv.Atoi(n)
	if err != nil {
		return scomp, fmt.Errorf("%v i=%d, type: %T", err, nInt, nInt)
	}
	scomp.N = nInt

	// R
	r := sh[2]
	if r == "" {
		return scomp, errors.New("r was not defined")
	}

	scomp.R, err = strconv.Atoi(r)
	if err != nil {
		return scomp, errors.New(fmt.Sprintf("i=%d, type: %T\n", scomp.R, scomp.R))
	}

	// P
	p := sh[3]
	if p == "" {
		return scomp, errors.New("p was not defined")
	}
	scomp.P, err = strconv.Atoi(p)
	if err != nil {
		return scomp, errors.New(fmt.Sprintf("i=%d, type: %T\n", scomp.P, scomp.P))
	}

	// Salt
	saltBase64 := sh[4]

	scomp.Salt, err = base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return scomp, errors.New("salt cannot be decoded")
	}
	scomp.SaltLen = len(scomp.Salt)
	if len(scomp.Salt) == 0 {
		return scomp, errors.New("salt length cannot be zero")
	}

	// Salt
	dkBase64 := sh[5]
	scomp.DK, err = base64.StdEncoding.DecodeString(dkBase64)
	if err != nil {
		return scomp, errors.New("key cannot be decoded")
	}
	return scomp, err
}

// generateSalt returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateSalt(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
