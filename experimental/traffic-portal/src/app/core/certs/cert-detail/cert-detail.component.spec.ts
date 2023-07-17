/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { pki } from "node-forge";

import { pkiCertToSHA1, pkiCertToSHA256 } from "src/app/core/certs/cert.util";

import { CertDetailComponent } from "./cert-detail.component";

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

describe("CertDetailComponent", () => {
	let component: CertDetailComponent;
	let cert: pki.Certificate;
	let fixture: ComponentFixture<CertDetailComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [CertDetailComponent]
		})
			.compileComponents();

		cert = pki.certificateFromPem(certPEM);
		cert.setIssuer([
			{
				name: "countryName",
				shortName: "C",
				type: "2.5.4.6",
				value: "C",
				valueTagClass: 19,
			},
			{
				name: "stateOrProvinceName",
				shortName: "ST",
				type: "2.5.4.8",
				value: "ST",
				valueTagClass: 19,
			},
			{
				name: "localityName",
				shortName: "L",
				type: "2.5.4.7",
				value: "L",
				valueTagClass: 19,
			},
			{
				name: "organizationName",
				shortName: "O",
				type: "2.5.4.10",
				value: "O",
				valueTagClass: 19,
			},
			{
				name: "commonName",
				shortName: "CN",
				type: "2.5.4.3",
				value: "CN",
				valueTagClass: 19,
			},
			{
				name: "madeUp",
				shortName: "idksomething",
				type: "127.0.0.1",
				value: "doesntmatter",
			}
		]);
		cert.setSubject([
			{
				name: "countryName",
				shortName: "C",
				type: "2.5.4.6",
				value: "C",
				valueTagClass: 19,
			},
			{
				name: "stateOrProvinceName",
				shortName: "ST",
				type: "2.5.4.8",
				value: "ST",
				valueTagClass: 19,
			},
			{
				name: "localityName",
				shortName: "L",
				type: "2.5.4.7",
				value: "L",
				valueTagClass: 19,
			},
			{
				name: "organizationName",
				shortName: "O",
				type: "2.5.4.10",
				value: "O",
				valueTagClass: 19,
			},
			{
				name: "commonName",
				shortName: "CN",
				type: "2.5.4.3",
				value: "CN",
				valueTagClass: 19,
			}

		]);

		fixture = TestBed.createComponent(CertDetailComponent);
		component = fixture.componentInstance;
		component.cert = cert;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("Fields calculated correctly", async () => {
		component.ngOnChanges();
		expect(component.issuer).toEqual({
			commonName: "CN",
			countryName: "C",
			localityName: "L",
			orgName: "O",
			stateOrProvince: "ST",
		});
		expect(component.subject).toEqual({
			commonName: "CN",
			countryName: "C",
			localityName: "L",
			orgName: "O",
			stateOrProvince: "ST",
		});
		expect(component.sha1).toBe(pkiCertToSHA1(component.cert));
		expect(component.sha256).toBe(pkiCertToSHA256(component.cert));
	});
});
