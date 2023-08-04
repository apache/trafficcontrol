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
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { RouterTestingModule } from "@angular/router/testing";
import * as forge from "node-forge";

import { APITestingModule } from "src/app/api/testing";
import { AppUIModule } from "src/app/app.ui.module";

import { CertViewerComponent } from "./cert-viewer.component";

/**
 * Creates any type of cert in a chain
 *
 * @param validStart Not valid before time
 * @param validEnd Not valid after time
 * @param issuer The issuer attributes
 * @param subject The subject attributes
 * @param extensions The cert extensions
 * @returns The cert
 */
export function createCert(validStart: Date, validEnd: Date, issuer: forge.pki.CertificateField[],
	subject: forge.pki.CertificateField[], extensions: unknown[]): forge.pki.Certificate {
	const kp = forge.pki.rsa.generateKeyPair(2048);
	const cert = forge.pki.createCertificate();
	cert.publicKey = kp.publicKey;
	cert.privateKey = kp.privateKey;
	cert.validity.notBefore = validStart;
	cert.validity.notAfter = validEnd;
	cert.setSubject(subject);
	cert.setIssuer(issuer);
	cert.setExtensions(extensions);
	cert.sign(kp.privateKey, forge.md.sha256.create());
	return cert;
}

/**
 * Creates a certificate and signs it using a CA
 *
 * @param validStart Not valid before time
 * @param validEnd Not valid after time
 * @param subject Cert subject
 * @param ca The CA that issued this cert
 * @returns The cert
 */
export function createCertAndSign(validStart: Date, validEnd: Date, subject: forge.pki.CertificateField[],
	ca: forge.pki.Certificate): forge.pki.Certificate {
	const cert = createCert(validStart, validEnd, ca.issuer.attributes, subject, []);
	cert.sign(ca.privateKey, forge.md.sha256.create());
	return cert;
}

/**
 * Creates a Certificate Authority
 *
 * @param validStart Not valid before time
 * @param validEnd Not valid after time
 * @param attrs Both subject & issuer attributes
 * @returns The CA
 */
export function createCA(validStart: Date, validEnd: Date, attrs: forge.pki.CertificateField[]): forge.pki.Certificate {
	return createCert(validStart, validEnd, attrs, attrs,
		[{
			cA: true,
			name: "basicConstraints"
		}, {
			keyCertSign: true,
			name: "keyUsage"
		}, {
			name: "nsCertType",
			server: true
		}]);
}

describe("CertViewerComponent", () => {
	let component: CertViewerComponent;
	let fixture: ComponentFixture<CertViewerComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [CertViewerComponent],
			imports: [APITestingModule, AppUIModule, NoopAnimationsModule, RouterTestingModule]
		})
			.compileComponents();

		fixture = TestBed.createComponent(CertViewerComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
		expect(component.dsCert).toBe(false);
	});

	it("new cert", () => {
		component.inputCert = "invalid";
		component.process(true);
		expect(component.certChain.length).toBe(1);
		expect(component.certChain[0].parseError).toBeTrue();
		expect(component.certChain[0].type).toBe("Error");
		expect(component.certOrder).toBe("Unknown");
	});

	it("good cert", () => {
		const today = new Date();
		component.inputCert = forge.pki.certificateToPem(createCert(today, today, [{
			name: "commonName",
			value: "Test"
		}], [{
			name: "commonName",
			value: "Test"
		}], []));
		component.process(true);
		expect(component.certChain.length).toBe(1);
		expect(component.certChain[0].parseError).toBeFalsy();
		expect(component.certChain[0].type).toBe("Root");
		expect(component.certOrder).toBe("Single");

	});

	it("root chain", () => {
		const today = new Date();
		const ca = createCA(today, today, [{
			name: "commonName",
			value: "Test"
		}]);
		const cert = createCertAndSign(today, today, [{
			name: "commonName",
			value: "Test2"
		}], ca);

		component.inputCert = `${forge.pki.certificateToPem(ca)}${forge.pki.certificateToPem(cert)}`;
		component.process(true);
		expect(component.certChain.length).toBe(2);
		expect(component.certChain.some(c => c.parseError)).toBeFalse();
		expect(component.certChain[0].type).toBe("Root");
		expect(component.certChain[1].type).toBe("Client");
		expect(component.certOrder).toBe("Root -> Client");

		component.inputCert = `${forge.pki.certificateToPem(cert)}${forge.pki.certificateToPem(ca)}`;
		component.process(true);
		expect(component.certChain.length).toBe(2);
		expect(component.certChain.some(c => c.parseError)).toBeFalse();
		expect(component.certChain[0].type).toBe("Root");
		expect(component.certChain[1].type).toBe("Client");
		expect(component.certOrder).toBe("Client -> Root");
	});
});
