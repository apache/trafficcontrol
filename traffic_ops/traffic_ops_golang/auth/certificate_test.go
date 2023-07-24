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
// certs on demand for testing, such as expired Before/After dates

func TestVerifyClientCertificate_Success(t *testing.T) {
	rootPool = nil // ensure root pool is empty

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
	intermediateCertPEMBlock, _ := pem.Decode([]byte(intermediateCertPEM))
	intermediateCert, err := x509.ParseCertificate(intermediateCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for intermediateCert. err: %s", err)
	}
	connState := new(tls.ConnectionState)
	connState.PeerCertificates = append(connState.PeerCertificates, clientCert)
	connState.PeerCertificates = append(connState.PeerCertificates, intermediateCert)
	req.TLS = connState

	err = VerifyClientCertificate(req, "root/pool/created/above", false)
	if err != nil {
		t.Fatalf("error failed to verify client certificate: %s", err)
	}
}

func TestVerifyClientCertificate_NoIntermediate_Fail(t *testing.T) {
	rootPool = nil // ensure root pool is empty

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

	err = VerifyClientCertificate(req, "root/pool/created/above", false)
	if err == nil {
		t.Fatalf("should have failed without intermediate certificate: %s", err)
	}
}

func TestLoadRootCerts_EmptyDirPath_Fail(t *testing.T) {
	rootPool = nil

	err := loadRootCerts("")

	if err == nil {
		t.Fatalf("should have failed to load certs with empty path. err: %s", err)
	}

}

func TestVerifyClientChainSuccess(t *testing.T) {
	rootPool = nil // ensure root pool is empty

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
	intermediateCertPEMBlock, _ := pem.Decode([]byte(intermediateCertPEM))
	intermediateCert, err := x509.ParseCertificate(intermediateCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for intermediateCert. err: %s", err)
	}

	if err = verifyClientRootChain([]*x509.Certificate{clientCert, intermediateCert}, false); err != nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestVerifyClientChain_EmptyClient_Fail(t *testing.T) {
	rootPool = nil // ensure root pool is empty

	rootCertPEMBlock, _ := pem.Decode([]byte(rootCertPEM))
	rootCert, err := x509.ParseCertificate(rootCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for rootCert. err: %s", err)
	}

	rootPool = x509.NewCertPool()
	rootPool.AddCert(rootCert)

	if err = verifyClientRootChain([]*x509.Certificate{}, false); err == nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestVerifyClientChain_EmptyRoot_Fail(t *testing.T) {
	rootPool = nil // ensure root pool is empty

	clientCertPEMBlock, _ := pem.Decode([]byte(clientCertPEM))
	clientCert, err := x509.ParseCertificate(clientCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for clientCert. err: %s", err)
	}

	if err = verifyClientRootChain([]*x509.Certificate{clientCert}, false); err == nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestVerifyClientChain_WrongCertKeyUsage_Fail(t *testing.T) {
	rootPool = nil // ensure root pool is empty

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

	if err = verifyClientRootChain([]*x509.Certificate{serverCert}, false); err == nil {
		t.Fatalf("failed to verify certificate chain with valid certs. err: %s", err)
	}
}

func TestParseClientCertificateUID_Success(t *testing.T) {
	rootPool = nil // ensure root pool is empty

	clientCertPEMBlock, _ := pem.Decode([]byte(clientCertPEM))
	clientCert, err := x509.ParseCertificate(clientCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for clientCert. err: %s", err)
	}
	result, _ := ParseClientCertificateUID(clientCert)
	if result != "userid" {
		t.Fatal("failed to parse UID value from certificate")
	}
}

func TestParseClientCertificateUID_Fail(t *testing.T) {
	rootPool = nil // ensure root pool is empty

	// Server cert does not contain a UID object identifier
	serverCertPEMBlock, _ := pem.Decode([]byte(serverCertPEM))
	serverCert, err := x509.ParseCertificate(serverCertPEMBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to extract x509 from PEM string for serverCert. err: %s", err)
	}
	result, _ := ParseClientCertificateUID(serverCert)
	if len(result) > 0 {
		t.Fatal("unexpected UID value from certificate")
	}
}

const rootCertPEM = `-----BEGIN CERTIFICATE-----
MIIFjjCCA3agAwIBAgIIKk/S4uUM2nIwDQYJKoZIhvcNAQELBQAwZTELMAkGA1UE
BhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNV
BAoTBkFwYWNoZTEMMAoGA1UECxMDQVRDMRMwEQYDVQQDEwpyb290LmxvY2FsMB4X
DTIyMTAwNDIwNDIxNVoXDTM3MTAwNDIwNDIxNVowZTELMAkGA1UEBhMCVVMxETAP
BgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNVBAoTBkFwYWNo
ZTEMMAoGA1UECxMDQVRDMRMwEQYDVQQDEwpyb290LmxvY2FsMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAyJtV6lUQ8ecqI9D9HQKnCaD2gjU7CfKRNZe8
FEHXlA1rQlU+rDpPmafHZMNaXXJusOxIN70nGEnlTn9ZL+8TMCsKyeq5Y6Diqubw
Ws6kgVpsG73T2X2/gdcow3poCcSAOO0JZypVK3vFlVoB/fBdvB2f3CusV2qYmphf
ffUcykKSSWV6lbeAZYwwOwuKy+eWmgedEJQIQqGqfNAal/UEiGeiqvrsfzu/DzBF
0VXcljTJnXLkESgxESIUHwhIDjcM5sFS5NW/Dru4lodfUPDMW8B9qrW7j7ocDWLK
gbw2ct34HKVBwXC7dYosnawZJ9IVeKa+lMQDRGb5N+Rw6j/iX4JOk5m16bqSEJnh
U4vAk502IfXGFULLDCbm0ju84Hul4oq7I6rPrnTinWGMUCkzyKjhs/7aBvfOsmFr
VyGCnaLw+rEdOr8pPWYP5hfBfggjIoFHb25DWTIbJeu2wr0+F33/60w33RnXoKCl
zmR5Bsfxqaayxd8FcisKignaeibOUtcd+I0xunu/VjXzX7MEA8qrEdFKjiAmj+U8
WiQkf2u3v37vj1mA+qp65qudaAHwvDhJmVdviri1OGBqF3zYM8xrWxXxQpdQFZh4
63XxipzF0sTAoDgfcvDsCeKwfwXBvisx4dGHa7a72YvQUD8Kts7XHgwfUBrTWc3m
RGTKc08CAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgGGMA8GA1UdEwEB/wQFMAMBAf8w
HQYDVR0OBBYEFIcDMB0+S1Glsa3cEWW8sa396avzMA0GCSqGSIb3DQEBCwUAA4IC
AQBDApE0YcWX1MY9MII1ddheaeAuA5DuwPtSpWS2RWu2NEOZxuQZSrlyaJS0J16L
7OKElgI5M1tHFO2/3ogukIZzrawEYCm/70lYpR9IrSwAEtrkE01D258+d2YpUwlV
uHdYV/rr5pMqXh8hL8ZS6m3CPuMz/w0mytgMiAPLCQb6n2EOuJdbp3EC2FNPa6r2
5w/MZ+Xrih2fYipVt4oNandKsLdnKeJcvr9h17z6A+QQgA6+BjMGp1eLT/oz3vjo
+4i0Jf8LvkyN9JCmG7zBEMnFxHBebgQeY9TfPS/wOYxpD97UrmEWa7xHi3g9xbnr
3LBoi0rrWkXzYq/CpnEWrMVQo4Z3AGm7z5k+d6KwOoyRntWZA94YUTfGz4Jln4cD
4s4soK5hv87LOcithnajbYcujwhm/YPIMQxh3C0Ziu9qWvYNGS9nnoJvITvQZQOc
YD2Htd/PQTRbqoL93Xdv32/f8zAxHru/4xjf/CkBQ63HvQcY/0z7FW91Bulp31qT
B/bQ1a6PAaCkYC0SXCNnnnyUYx4xS0ggSBFk0ruJpo6NC2+8fDiwxGUZ54VhACw3
ASix7tHQ/yOcqj2gf6NpZpcdQY+WTgpEWr6orGEuNSp+5pEmD2Qi40TD8KccYaMe
upDLy1qU3S7/C6TXdiIXb5ZX12wreew33AKTc/EqNoEZOQ==
-----END CERTIFICATE-----
`

const intermediateCertPEM = `-----BEGIN CERTIFICATE-----
MIIF2zCCA8OgAwIBAgIISnLu/F5oKSAwDQYJKoZIhvcNAQELBQAwZTELMAkGA1UE
BhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNV
BAoTBkFwYWNoZTEMMAoGA1UECxMDQVRDMRMwEQYDVQQDEwpyb290LmxvY2FsMB4X
DTIyMTAwNDIwNDIxNloXDTMyMTAwNDIwNDIxNlowbTELMAkGA1UEBhMCVVMxETAP
BgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNVBAoTBkFwYWNo
ZTEMMAoGA1UECxMDQVRDMRswGQYDVQQDExJpbnRlcm1lZGlhdGUubG9jYWwwggIi
MA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC/77BSHSnuXFAnqu/FH2ecke/E
UvhfHzcF/01qXPRK4tXTJfA0whtYoJ2qIpdBDH5UcnMfHyHWHXnay/4OYbwFM8Fi
EpexX1ecgRxph9S8KTvh1pzkI3axfQoz55xoQQNFcJZ70QxgCs9WinCqY2Y+9SLo
P1rFZSRhCSYAuveyDfDVDzU21vDYFC9uLYZvolt5G/cBPHOOTF+KTgrk6Xg3XVYn
XvId6gva1guuxRzIIRDq1Lh6zLLH2Ox632OkQs/OwrkkfwFoczvNvAOIxMf7NmTd
9j1MVBH5Agu6RwZNQQ9JZg4VugugcHN94REiyg01ypQ3yfXZGATDVbp9qNTAyH0Q
lOeqqRxJjwS/9l99b28uRDE9bQFn4+uCU/pGJadAYAEg399Vp5/77o+f+kYnaeqr
NomCXJv5nGBTaqc7EwbV0/pjqzelfwy/O0wGSSqNddQWGdwQ2Sm6cAqlz4uB5Zpu
5yBDEToaQ0yhmNanoqhQq3pcEDeccnTP2aGa3P4SlDrjSDFrmJdffHg5xWeaslxI
mZfxH2ChiE+y0wtuQ91A/h/fIA0IRQIUpog2d8LNeKtsboErXZRORO7L6x5GFAYM
NVnWgw/eF+zl5fkk+OI2UK3PFMosjjMjtHTcWyKFuXpE2v0wj5vcXMTmuOGClOx3
yq0sB7LWDo0yYkfo1QIDAQABo4GGMIGDMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUE
FjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNV
HQ4EFgQUdTW8W+avOivYqcxslZgJwkfR8CUwHwYDVR0jBBgwFoAUhwMwHT5LUaWx
rdwRZbyxrf3pq/MwDQYJKoZIhvcNAQELBQADggIBAMf0zJYiMyQgwrXEKKSChzKr
kZzoOxz/9Jey2IfWi6exsalj0lguX6IBE5oriiNHOOb56IT84EjBt4uHKolGdzKl
KL+RCFGoe7h40KigY/I9pBkUNm30N21QFV6lHh2fXhjkLCBExpgP0VxZJKwZn5uH
fLvXZSFxgvGKHiuA58eW41S8xE7jPsTC3eprLTXIpUF1Yh5en33bhVtVdtaw3YSu
lbKY5y+kZLgJIlcediz4IqZvZ5MiEaD1e+tNwmtN44yayA7JMihk64qUE4R/j20l
JYHNIyjnujYRKxHbU8oQq2XTvQbTLl3MlYYzlXXQ/g86pdIPiricPP4tpBC02ASz
8iMcdpFtkC96M4lf/n9GZkFBIyXStcoJXSFxcEDDzpW2FYzu8SOr5Lj7YFmAVw/q
p4B56wOWEEvRTcM9v7+uP1AbH95KoAr1hd2z/tp5JFtQQrnmSxIRcomrXbp5RBpt
dMugMKmkZJfI9waXLqRn8WhEQ6/MloyLNjfAUn2IoK+araSBbfusMlfc2skJK4eD
AO71dMq6EVQr6TrTzfL0pUjPJDRvff6DPIj2mNcXvRPPEjO8VQzi6NOEitssYbwL
QOlf+M07UmfY8RjqhbTQgB6nAtAu1g4y7oHf7XTMZlRc0EZkJiHXQ5EoHlz4D+3J
oXi7IrLNF0+1/RFnYmKY
-----END CERTIFICATE-----
`

const clientCertPEM = `-----BEGIN CERTIFICATE-----
MIIFtjCCA56gAwIBAgIILNSPo4amjqowDQYJKoZIhvcNAQELBQAwbTELMAkGA1UE
BhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNV
BAoTBkFwYWNoZTEMMAoGA1UECxMDQVRDMRswGQYDVQQDExJpbnRlcm1lZGlhdGUu
bG9jYWwwHhcNMjIxMDA0MjA0MjE4WhcNMjcxMDA0MjA0MjE4WjB/MQswCQYDVQQG
EwJVUzERMA8GA1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRlbnZlcjEPMA0GA1UE
ChMGQXBhY2hlMQwwCgYDVQQLEwNBVEMxFTATBgNVBAMTDGNsaWVudC5sb2NhbDEW
MBQGCgmSJomT8ixkAQETBnVzZXJpZDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBAKIjAo5kSO+MYvTLc8yLjzNMvTmOIncYFWd6qUFh4io1aTO3CiD+bEWp
3m5Q+8tACFAIBfLLf06cM0PW2bTwJet4Ol+Wp6v5ieKWDGtz5Ae6piGfu+C3TCW5
mljdEwDVEXMxftufhgfdc4qVj/plg8iuQRyj4AEhmFtXEvOssWixAPoyk5afUBOB
qaJwqwDY2xja6lE9ECvHULJIGviKFHsZ69TzkMzq0af2/dYRrjG+Zj1ACm9WO6y0
Y/v/ztcnFtA6Z4EEqDfMA7+BJqDdhjzJ7ISymboNLEtq2sYWrJe+5DPsACSpoTly
FPqn3omy4ax8c4lAnF2Ud/KUqouctGzBwP3lz48aW/Nj5anXh6HoI/h1qGXrRRHI
fTAf0rcbSELkIRiDHpCzIMLBjEZefPiBqjukvZHstR34+TwItnCqNhfGl/+dLGbZ
W5xw2JGcOu+ccTNqUI7bs0wxC8zppAEZ5epXmifhMEwsrQPB66vQIR8lWZyyczOm
rllE36gwDoMjcjSJLuOv7iHzmJVK1S0lccZ1gfoCAhV1cK+YN/BE0Qbtq1ehmYV2
R0B6kIqzlG3L+s7rNWTj814YbPh8WgMlnApHE523OdeTaScw8ukGz/p4apP/VRsX
fuLzUwB3xabIsllFClNfr8MaZgirMzYVDDuvTbzezSorbplU4I5nAgMBAAGjSDBG
MA4GA1UdDwEB/wQEAwIDiDATBgNVHSUEDDAKBggrBgEFBQcDAjAfBgNVHSMEGDAW
gBR1Nbxb5q86K9ipzGyVmAnCR9HwJTANBgkqhkiG9w0BAQsFAAOCAgEAd/kOLQyr
5PNMirK1EfoYnO2lme/QyO44Wr3kZZQ4X6ZsBKYciuC09nVmBc2VDQS/YOWP1vLu
6UbH8pho5xdqoj+KvDsJtRhKeN+0LxHgJ0u6kmq2Fid3GG5kwxeSEz/7LOJd5Qp9
E0MbRtm1eu19IVl+3XMSyNLvA0vfAAawpFa85E/rDZfUKWa3JgrjweYXCpz8EgmB
mGroZO4uBc40gLcJOcxGqEB0YioQ6WqLDpOGhWZRxQRdzC1kp64fYkiI7wtX7fBp
VnUwJmi1+A3gLvrT67zkauO8W9niLrorvu3naBDgtRoZhxTsRCjx26NV4aQ9Kx1p
4c5H8RfrD5b8vo+QXbodE9Zj2IZfJZew3/xM9W2GQYT/gPk7AInWKnefUbes3WuE
gzdVaS8FpQCbP7+VJNTIutG5AdvZMz66nYtJOMDb1NB8w/oAI9DBKhNyGwL8AFGq
FbYSDEWpwnu/bBq4uXMRyEUVcbPdRRaabRG4XeowVA40d/VeDkbWlihTDEg49dLd
DIkfv7MZSwYYo35pbGyqzpTEotAZCqeXGSVwwMsjNfSwjb7JTohnnL+0aLfgpxH6
6v5GpzUawbJrwqvj3/BsBjJAa/95Dyh8IaCB/Jk6pgGLP8+s3SPNpE/JVHDDcKtD
83Q6ELF+UUDc7BsykJBtsxddSkFI9u9Hm7s=
-----END CERTIFICATE-----
`

const serverCertPEM = `-----BEGIN CERTIFICATE-----
MIIFzzCCA7egAwIBAgIIJ/N/XohlywcwDQYJKoZIhvcNAQELBQAwbTELMAkGA1UE
BhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxDzANBgNV
BAoTBkFwYWNoZTEMMAoGA1UECxMDQVRDMRswGQYDVQQDExJpbnRlcm1lZGlhdGUu
bG9jYWwwHhcNMjIxMDA0MjA0MjE3WhcNMjcxMDA0MjA0MjE3WjBnMQswCQYDVQQG
EwJVUzERMA8GA1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRlbnZlcjEPMA0GA1UE
ChMGQXBhY2hlMQwwCgYDVQQLEwNBVEMxFTATBgNVBAMTDHNlcnZlci5sb2NhbDCC
AiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAL3I9nBo7Bjd4aDc7415fghn
nzSXG43908hHXAstg8dZ4nUf3mzzEqSjmiqpJpx+Mr92jJAhpemVEi26WXneMXLW
PSimIRApX84veK3FLxRCpOebSx2QaBgy08eK3015wJ2s7faxLxuVNjdKHSRbZ7yU
vPGvSrjYQAb8XRwj8PKnmEFOd08U0O6QN2ib+5sAfWgbgab03Hv+xRMy1qgOqVsv
LC+dSUXO0HzxQoNd3cMqi8EWU43SAFYfa8DsDoObCTwl1qvLCRiRdrgkplgsPBW2
42Axf9qYZfrFF2IZubawY5jbo6WCvI5Nr/tOHWrSuEUFKSweD0+s+yOlZxaQ/NJo
B0pCiul/8JCKTENHh18UlQB5eVRjt3g0IJnlZK1J6vVXP7jgt3LbX6vXSPgmOInK
Pdo7xr37wxshjazE+rsEPZEjVmjbLDMLTrYRvBkfaATB7ht6oiGCP50MjPBRykoi
ASC0I/MQJXB7GksssSpYR8NZ5mnVRf9D+BeaL1/Nhjj/Kb7wcescV+uhaCiI5E+W
CluFZAD4vyO4uMBiRR/DmT83l0g6XUbaKvpiY3hIGbMjMqb8sVM1PEQsea1VzbBq
6GoMQEtfu/07q4nUP2mgUCQVGSluGjicjOvoikq8aLzj7WMEWbvWF3iPTXsebAJc
9Sa2hTZr7zZjhj1BrC4lAgMBAAGjeTB3MA4GA1UdDwEB/wQEAwIDqDATBgNVHSUE
DDAKBggrBgEFBQcDATAfBgNVHSMEGDAWgBR1Nbxb5q86K9ipzGyVmAnCR9HwJTAv
BgNVHREEKDAmggxzZXJ2ZXIubG9jYWyHBH8AAAGHEAAAAAAAAAAAAAAAAAAAAAEw
DQYJKoZIhvcNAQELBQADggIBAIheFqEz1OTG38/8N+r2gKMLoj7W7EsuJzfkXSgD
eSZFNkFf2R5Nx5U5SC+LlHqV5VtZGxqYMSusAB0mdaIqut1BBpgAOkmuVble0yN8
U2QUPfRbihZRbqwl/FhhzbHAHfW5F+rAQ7VmAsYPVDIA3Vnz8vODVrhS6AT+1Zhd
dhmvlRVxsE8qlckeGaw5FS2rCUiybhmUclGjvrhQ1UJSMbe03V+yYM3xD2cRVHi5
9ycDvudVo8lhxVcMTkWIdft325F41Ra3BfRZqyq08OfSn75ny5w+GENiAI5Tei+n
GK9csiHzBz+EgjSfZ/6zwm+1dXXzYwo7pHFjM4tfylyv3V+k5HgWPQgmVfkViuvn
sIXyG2/wZQhWDYO17gxpQd4RS6tDc+Jf6T6T5uifjZIAsZAT8BxMoqzNleuNA15t
tYfOdfCUon0ZvPZvF8yPAQpRnWahCJdb++Obu1ftUjkTSTSSSHqgGrOa628U6ZNj
f1v1rGfY10dNB4hAbiMrSvLDRF+h6YT9ldbyh8S19JETBzfyGweDGLi7+MglK047
arn0qSijWuR4Zfy6B8muObxYqp80GOXsHhoLm6azKKv00yhND3Z+16QSYHDuDwNf
BbzBP0SXV643tyWjUZOrVXXw3bv5X/RRjD69x0bw2HeFFeKhGzgaBkK6WsAA0uyT
qrCR
-----END CERTIFICATE-----
`
