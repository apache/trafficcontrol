package auth

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"testing"
)

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

// TODO: Utilize expirimental/certificate_auth/generate_cert.go to create appropriate
// certs on demand for testing, such as expired Bofore/After dates

func TestVerifyClientCertificateSuccess(t *testing.T) {
	rootCertPEMBlock, _ := pem.Decode([]byte(rootCertPEM))
	rootCert, err := x509.ParseCertificate(rootCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for rootCert. err: %s", err)
	}

	rootPool = x509.NewCertPool()
	rootPool.AddCert(rootCert)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Fatal("failed to create request")
	}

	clientCertPEMBlock, _ := pem.Decode([]byte(clientCertPEM))
	clientCert, err := x509.ParseCertificate(clientCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for clientCert. err: %s", err)
	}
	connState := new(tls.ConnectionState)
	connState.PeerCertificates = append(connState.PeerCertificates, clientCert)
	req.TLS = connState

	success, err := VerifyClientCertificate(req, "root/pool/created/above")
	if err != nil {
		t.Fatalf("error attempting failed to verify client certificate: %s", err)
	}
	if !success {
		t.Fatal("failed to verify client certificate")
	}
}

func TestLoadRootCertsSuccess(t *testing.T) {
	rootPool = nil

	err := loadRootCerts("test/success")

	if err != nil {
		t.Fatalf("failed to load certs. err: %s", err)
	}

}

func TestLoadRootCertsEmptyDirPathFail(t *testing.T) {
	rootPool = nil

	err := loadRootCerts("")

	if err == nil {
		t.Fatalf("shoudl have failed to load certs with empty path. err: %s", err)
	}

}

func TestLoadRootCertsFail(t *testing.T) {
	rootPool = nil

	err := loadRootCerts("test/fail")

	if err == nil {
		t.Fatalf("should have failed to load certs attempting to parse PEM key into x509 Cert. err: %s", err)
	}

}

func TestVerifyClientChainSuccess(t *testing.T) {
	rootCertPEMBlock, _ := pem.Decode([]byte(rootCertPEM))
	rootCert, err := x509.ParseCertificate(rootCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for rootCert. err: %s", err)
	}

	rootPool = x509.NewCertPool()
	rootPool.AddCert(rootCert)

	clientCertPEMBlock, _ := pem.Decode([]byte(clientCertPEM))
	clientCert, err := x509.ParseCertificate(clientCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for clientCert. err: %s", err)
	}

	if err = verifyClientRootChain([]*x509.Certificate{clientCert}); err != nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestVerifyClientChainEmptyClientFail(t *testing.T) {
	rootCertPEMBlock, _ := pem.Decode([]byte(rootCertPEM))
	rootCert, err := x509.ParseCertificate(rootCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for rootCert. err: %s", err)
	}

	rootPool = x509.NewCertPool()
	rootPool.AddCert(rootCert)

	if err = verifyClientRootChain([]*x509.Certificate{}); err == nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestVerifyClientChainEmptyRootFail(t *testing.T) {
	clientCertPEMBlock, _ := pem.Decode([]byte(clientCertPEM))
	clientCert, err := x509.ParseCertificate(clientCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for clientCert. err: %s", err)
	}

	if err = verifyClientRootChain([]*x509.Certificate{clientCert}); err == nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestVerifyClientChainWrongCertKeyUsageFail(t *testing.T) {
	rootCertPEMBlock, _ := pem.Decode([]byte(rootCertPEM))
	rootCert, err := x509.ParseCertificate(rootCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for rootCert. err: %s", err)
	}

	rootPool = x509.NewCertPool()
	rootPool.AddCert(rootCert)

	// Server cert contains x509.ExtKeyUsageServerAuth (vs ClientAuth)
	serverCertPEMBlock, _ := pem.Decode([]byte(serverCertPEM))
	serverCert, err := x509.ParseCertificate(serverCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for serverCert. err: %s", err)
	}

	if err = verifyClientRootChain([]*x509.Certificate{serverCert}); err == nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestParseClientCertificateUIDSuccess(t *testing.T) {
	clientCertPEMBlock, _ := pem.Decode([]byte(clientCertPEM))
	clientCert, err := x509.ParseCertificate(clientCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for clientCert. err: %s", err)
	}
	result := ParseClientCertificateUID(clientCert)
	if result != "userid" {
		t.Fatal("failed to parse UID value from certificate")
	}
}

func TestParseClientCertificateUIDFail(t *testing.T) {
	// Server cert does not contain a UID object identifier
	serverCertPEMBlock, _ := pem.Decode([]byte(serverCertPEM))
	serverCert, err := x509.ParseCertificate(serverCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for serverCert. err: %s", err)
	}
	result := ParseClientCertificateUID(serverCert)
	if len(result) > 0 {
		t.Fatal("unexpected UID value from certificate")
	}
}

const rootCertPEM = `-----BEGIN CERTIFICATE-----
MIIFjjCCA3agAwIBAgIIYW10BDhSkLgwDQYJKoZIhvcNAQELBQAwZTELMAkGA1UE
BhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNV
BAoTBkFwYWNoZTEMMAoGA1UECxMDQVRDMRMwEQYDVQQDEwpyb290LmxvY2FsMB4X
DTIyMDkwNjIxMTc1N1oXDTIzMDkwNjIxMTc1N1owZTELMAkGA1UEBhMCVVMxETAP
BgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNVBAoTBkFwYWNo
ZTEMMAoGA1UECxMDQVRDMRMwEQYDVQQDEwpyb290LmxvY2FsMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEA5uKjtuvogDySfFUEIUPluer+1KeH84Vz1NSF
2yShaHSNtPVloPc5sRzs0SHDDMmzY8JPeW5kyFJspuw+hDq0OfOANqt7xkDguaD9
GosxGbJNYU3BCFNXrxB7pOuYFJptlSWH+PWbbxnqEy3mKk/9SA3n7xyEbc+J4Jue
9QtQBBqgyk3updMrWf+bMDHmA6KFnWzfUZV9WrN9GijeqByiZMB7X7fMfM7bNy7d
EbOrBDlFdWmQWSfRgekNc5/dxXk4G3xKDXQysdLLTbMT5HVdxvCy3cgLj9hAwSz2
K2ViyMbcahiYwQ3BVGiszg2wjZr6DoWg1eQqVBHPMSWZG+la5dqugTcX53SGoKVi
YYsX/Mcj+T/HYyzmTLIwVR9hYibgh4tPx6fjqPZRl9BMhSyebWvmbIqnTDzXIZUu
ZiEodweunCztcAG5oQMCwGzB8UTfnqp6juW5pSq+Tuz5wnr2iNZYcpsh7e9TdWEZ
B3/uyBE2i0L7/IadDjFNx8uIY+j9Usgwo9Od6cfhK26e2MKoICxiNx/KfOlDjEli
MaRn+owEaInLGVnu74HcWBm3lkv30k7T64yTeApxw/Qtd+WDR91RWFEE2gdbf2vL
w1WT0n/HFgTrLoiTEukEfcDXfjPkLsWSYydDMAJyymgBuY3GPjL8DmjNwZohe8sc
+hJRe7UCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgGGMA8GA1UdEwEB/wQFMAMBAf8w
HQYDVR0OBBYEFKHzciB4emUO4tqzNwr+U0/Dy0kKMA0GCSqGSIb3DQEBCwUAA4IC
AQDVLadfhKTI/q3hpkIqdaANriMZ8EUSXzcgFu2ockdqh0UjQVg4ZuaIx0GHFkl1
gtgra5L9F1bkEyCCFwiVbheZ99NKBmamEdb/ke3aXkRlsKxFPOWsnOqEEqqLTnjV
5jXv/6D93YbrL/L9rQHV35mYrWHGrEE7qYQfAdo7e9Cy805GuaCKk9BvjfxG+WnI
tmUZjOIIZ8tlcwXibcfKB5T5xUBhNUDaA02LcxYpEQhSANpG4129I6ckz9aOtemb
yHXu3UeMCrh6UOAq7nmTQXp1BCs+zgolHW7GRGBf/UI5IC3AV29LjuM0qs2oLFSP
h87lqobmDZDXgbsHaKY+IaM99sc0z8OtQEHk/b5kqxJGTbCnsDKhQuDICceSG5gS
ZZiC+l9c8BE5pHLGL9omsKQ15QWAo11RoOCeDQHdQ0YjXKa4dGlTomWFWge2qAQx
G4ltnOmj8WggYcYJoZG/XQaQN9iaL5L/0AIu8zFwIHaNCDo7s2Ow6QKpb99PhTML
pB+dlC+T7BZGoicPhcyh4wPyEF3ebNv3eAIuJmciYGy+0YwcxzhNshkzg1V/lIA8
iZfqa7xERCPGS03yy6quR5s37py74osk8xV71jRn/Fp7ogGEsVunNfEObO1Se3O8
EYvdxllGaU6j2bVE4v/8CukTlUNPQx/+ty6EYaNsoIgkBQ==
-----END CERTIFICATE-----
`

const clientCertPEM = `-----BEGIN CERTIFICATE-----
MIIFrjCCA5agAwIBAgIIAOfaLIQ3CvUwDQYJKoZIhvcNAQELBQAwZTELMAkGA1UE
BhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNV
BAoTBkFwYWNoZTEMMAoGA1UECxMDQVRDMRMwEQYDVQQDEwpyb290LmxvY2FsMB4X
DTIyMDkwNjIxMTgwMFoXDTIzMDkwNjIxMTgwMFowfzELMAkGA1UEBhMCVVMxETAP
BgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNVBAoTBkFwYWNo
ZTEMMAoGA1UECxMDQVRDMRUwEwYDVQQDEwxjbGllbnQubG9jYWwxFjAUBgoJkiaJ
k/IsZAEBEwZ1c2VyaWQwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC2
jzAsGVL5CHOt3irZxabD147EKP8d7PRUG6b8MyHjwtJW4gRj72u/CYn3d37kGXCK
vPnhcEIgyeepeD1BmyPUZqPRLGvp7mssSOCBNPJ7ND+WtpqZOPWWEVRp8FVY8ecU
RAPLi3k2DYQj3JNh6p+s7tBrUitHUm2kJ3GJDVJyDEPIfmXXMwfZ9JkaxAZFiRBg
pXq4gI52pz38UM0N1yQUFasSM/HodG0Pk0f4ZN+AoGz9KC7niQfWUucr72E9DSRf
T4L1h0F6yIgshdTvcPrY7/kg3+waPH73sy8WVNvQFVckeGOF0sDHS47hSpOvnInB
EgbCzG9gr5DO2Ik+LNYboMLUbMziyRJBBm4QdVMMbmgc+sLiifo17HVXkgngPl8I
75m47yXUIN3IV8o60OgTxV/NW7Ix/9NEmJKxu2WALJmrMwru/ySp6E/5vmQBaYn3
aczpTUi9FFQbHpkBRCBk0/Kty6fQ98ywYyJc+27cQIhH8qgOThbn4uux2LuJZSym
ek9UdSZVYIvIsZULQojvn8h+YZR94HdhfHccx32tXBqlVPMc+vApZTl7LI6FEa3J
+pAV1gRd9hsMMJ9bS0kYEDaeO5Bgw98IXdtecxuy3+s5Sq8RHIUbhpXPgn1H2T63
T5xlr1g+88BA/S0St6BcVJ8U8E4NkhRlySrr2gnpvwIDAQABo0gwRjAOBgNVHQ8B
Af8EBAMCA4gwEwYDVR0lBAwwCgYIKwYBBQUHAwIwHwYDVR0jBBgwFoAUofNyIHh6
ZQ7i2rM3Cv5TT8PLSQowDQYJKoZIhvcNAQELBQADggIBAIUzKgLHkWQMwo5fyXWH
yP/Tt52lrKLb+CyfhY617XhYFD6H5FHo+qRzNCd0FQqxhaHKAZmsvB5b960eWsAt
4P+1rKXPpaBLVYf5oGKY8/1O2YMKvJYcT7S+mqSLrESXzPLcV820o62ZBUzof/Yn
/H5pxeKZe9j7LcBiQ6dku7sQY8PUiY2sxBTRWm+EKM0GiVS8B2fasI6x8rrz8IJA
cgwyEqCyqzGkrCdmeMUpfeZcefDOAJoQsbMVhhSRkd2h4f0QGJjoLdfNJSjsvfKy
BaiE45qq4YoOfXn0M+WN1pDoxFcutfOoPII2JkQeZ3Avybrfa2+6B0h13QxfA0ww
3M/qrqe1p1dxw6hCWdemfRbiK56KpFNOP/O7AVBkZmRPTiWybmBeHOxNwxP+jz5y
hlKQBneGWZIp0svEt+yFKYxD/RBzVlLektriC1gbowZ5BBipVApjs9hcMfEpF9fu
yTqj6MkfUsEMyLNiYU5vgbE7vR3HjN2yKksOge3BGA59tivwKmMmZs4O+6RWm7lR
aXK9ttZ2bYILS5T7Br8eo3+n3+IKAxnLmUxJ+WaO87ID8yIdjHea0wanxv6ocE1/
kVyluEJ8E2iAf4Ax2d3yAXBuMNZkW/laGWptrSKivBGGUuU5rZcSV76TUh+21zgU
reqlKjsX1XKA3+uPcxYJ/22o
-----END CERTIFICATE-----
`

const serverCertPEM = `-----BEGIN CERTIFICATE-----
MIIFxzCCA6+gAwIBAgIIQ09gD+ZAsDUwDQYJKoZIhvcNAQELBQAwZTELMAkGA1UE
BhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNV
BAoTBkFwYWNoZTEMMAoGA1UECxMDQVRDMRMwEQYDVQQDEwpyb290LmxvY2FsMB4X
DTIyMDkwNjIxMTc1OVoXDTIzMDkwNjIxMTc1OVowZzELMAkGA1UEBhMCVVMxETAP
BgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNVBAoTBkFwYWNo
ZTEMMAoGA1UECxMDQVRDMRUwEwYDVQQDEwxzZXJ2ZXIubG9jYWwwggIiMA0GCSqG
SIb3DQEBAQUAA4ICDwAwggIKAoICAQDB+lLNls2JhyCIEp3UrGYMw+0ZgvXFEIG+
Na/Ru6yw2vVMeOMEN+jrG61BN42+dqx3LyDwbyBAMbESqXHTRB3fNQsLi5oxTr1X
VDfto9Js4Naglzn9awEF+2oTzbqRMzoa5QKZ/Q3haAu3W6hcCYESY+fYPqNIeYgY
1CsDqwMu6E6iEZV0aKBrRx4ZeEPHKYuvDmYa1w2MoQLRjQsLz6VgAUvhdIxBkd+A
eK0zSyYE7XTrx1f5XKKQbqqnyoIbRa0BLHLVUvrIH8PHQEnSXZOG9zs6CLKx5a97
kedN4UhLA1syax2tTtVwW6xYTR7pUlISNRUGzWkpI06b1nXguc85xJr7j01IInLV
/h/pUTN+JoBEouiNS3kNksgIk6Xxh5DXAk9ZHWo3c9/OdJKCnaFfJ/MmJ7JXAAc3
41C1cgUVNUS4z2UMyAzyNbp0KJyOH6edU+yVo2PuAIo7Twkq3uivxdO4D4bt01s4
mbHJ8emfG28MjJj3oeEjW/8JtWN9nYBKcsy+HUYUV1nmV/IzPF/3BCpBOuKPwaOd
8GFQ/hZ79bHADt1J/2hLmQh6M/lv5nf/KSD3wdpXkOxYEPO79i2WaDMMMji1IFlp
cyz5jce2C8jj+kMrDuqN1umHLQhCadN9+jXCNQnbybM0Ryn1gHSB55M7u/K9BIYy
Z9zVeLfK6QIDAQABo3kwdzAOBgNVHQ8BAf8EBAMCA6gwEwYDVR0lBAwwCgYIKwYB
BQUHAwEwHwYDVR0jBBgwFoAUofNyIHh6ZQ7i2rM3Cv5TT8PLSQowLwYDVR0RBCgw
JoIMc2VydmVyLmxvY2FshwR/AAABhxAAAAAAAAAAAAAAAAAAAAABMA0GCSqGSIb3
DQEBCwUAA4ICAQAS0O0IbQnjjZFeIP45VN8XNq6XmZm7BCsWZ/VLJEt4E0Tc2cN1
1F/9HyHR6UCWtI2L6kOFlwPCZEtMVtNsaJ+W8cG2WFpbmKOg2QLQHT1qNiPEZC93
56zSEggMh/cl3+dU4hkm53DkzONdvOcNtkAR6PeSVdxm8rEQJU4d2xFXj2C2G5Zi
Gf3mq23SH3ptMzL7YeY6n7jj7VqQyd03eqUexrl23WLkBbyyLzmdMWE/4c1szKNn
yTH3y1wXpyEyJmiz68mUO9L3DAvcYrhpeABHkYP5PWbfIXYb8Uu9iql6faQun6w9
dWr8dA0hueB8Amc6vnqu5Ym/Hp1iBWHD54AyRbip0d7jjL64CrZ6/bIiX8umaxcK
W7d2cYaSz46oQ3SmGAeT11Ky6Md5yQkfmLXrxkfZ661hmT2vefuh1m/mUl6T0sJW
OIPhgP9S02SASYRT4FHJdZy8EHWio/ZFH6YxYvlXnvzlpv0YgbMsNxkjsry2OXoQ
wjW5/epFn13lc33Uu01EqN0qAcWMFQT/RhbAERa+WyaEy8GP43IICUvRo13YRa7A
BWMeOD4KA+mczxpDLe4KmgzibdRQ2/OJFYoCAdnsDghin+sawnfcqwrCcq8se9Ye
FvhSQSGY7OsQmFg/M3scKScOiNUieyWyJ/4o3Ug35BAJi6o4gOKLI08Dlw==
-----END CERTIFICATE-----
`

//TODO: Create and add intermediate to chain for testing.
const intermediateCertPEM = ``
