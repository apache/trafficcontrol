/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import * as forge from "node-forge";

import { oidToName, pkiCertToSHA1, pkiCertToSHA256 } from "src/app/core/certs/cert.util";

const certPEM = `
-----BEGIN CERTIFICATE-----
MIIDeTCCAmECFDWSnKTtkcoRnoTz6ChHqTuvCUPHMA0GCSqGSIb3DQEBCwUAMHkx
CzAJBgNVBAYTAlVTMQswCQYDVQQIDAJDTzEPMA0GA1UEBwwGRGVudmVyMQ8wDQYD
VQQKDAZBcGFjaGUxGDAWBgNVBAsMD1RyYWZmaWMgQ29udHJvbDEhMB8GA1UEAwwY
dHJhZmZpY29wcy5kZXYuY2lhYi50ZXN0MB4XDTIzMDYwNTE0MDY0OFoXDTI0MDYw
NDE0MDY0OFoweTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNPMQ8wDQYDVQQHDAZE
ZW52ZXIxDzANBgNVBAoMBkFwYWNoZTEYMBYGA1UECwwPVHJhZmZpYyBDb250cm9s
MSEwHwYDVQQDDBh0cmFmZmljb3BzLmRldi5jaWFiLnRlc3QwggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDKAcQK9fe9w7p7eMnygnlV0rlbUdVr9DEQpKym
Ul7zGj9/Ta3n0h8xWrmmMi2ZJnIUI4AV7HKaYXiAke1rbEx2jAdvXdjNm/S7RORy
M0piJc8Si4/EJI1sZU17kZ7howXJvAMCQBqcI+hG93ATlUIOoYuluX7wSNIMw1Np
lT5bcmVDf5nVQGnrPw22mCGjH5JBxW5i1DjCoNovHfFgNmwP6y8C1jygoMPL+rxl
sq8fyUE/+qtcEkjUrr4oi9kjTESDqHghrkejKk6NPlPi97SDz2Ffdagoq2aqBhw9
P86JgplPVHHMWOLXBww0wPAClqY8H7CIt5rgZzoWmoR0DjjNAgMBAAEwDQYJKoZI
hvcNAQELBQADggEBAMFz7k+egg+hP86ylEAuUfcy/beO3Pf3Fn7oMh5MDENfOzON
IFqZOQ8pN1zfoAx0rRTzYHcg/AZs2AA4oh+WyEKHDrmICGfsF481b6A0EarZ/cRy
MF3Vh5rTd8ujWT4V9GP3Hc/I3F5tUKxPWiVEKTVRr6wzjwtXctOnhcbB3FeRtGDY
CfVBYMSEDJmAyMchfST/GwdG46Ak2TSaMpOf6tL5aMw+xfmDI68JGwG0LNliyEoW
xOHRCtWd5Q+Sn3rgx4h6nzdZOGHw3HwDbsX/y/dZNc7luUImEWwTyhohnO9XqaBX
EsdMDJmBaoVum+sR6ch08TsqrTHAfdB3xJF37Wc=
-----END CERTIFICATE-----`;
describe("Cert Utilities Test", () => {
	it("oid to name", () => {
		expect(oidToName("thisisnotanoid")).toBe("");
		expect(oidToName("1.2.840.113549.1.1.12")).toBe("sha384WithRSAEncryption");
	});

	it("sha cert digest", () => {
		const cert = forge.pki.certificateFromPem(certPEM);
		expect(pkiCertToSHA1(cert)).toBe("4afd0f20041efe4474ef44a125d0715cacc269e5");
		expect(pkiCertToSHA256(cert)).toBe("531d5eca87cc038077c6403615e0df646c0eac604a6292b79db1f8a014a7fdf8");

		expect(() => pkiCertToSHA1(forge.pki.createCertificate())).toThrowError();
		expect(() => pkiCertToSHA256(forge.pki.createCertificate())).toThrowError();
	});
});
