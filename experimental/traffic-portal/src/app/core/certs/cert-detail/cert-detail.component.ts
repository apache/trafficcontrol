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
import { Component, Input, OnChanges } from "@angular/core";
import { AbstractControl, FormControl, ValidationErrors, ValidatorFn } from "@angular/forms";
import { pki, Hex } from "node-forge";

import { oidToName, pkiCertToSHA1, pkiCertToSHA256 } from "src/app/core/certs/cert.util";
import { LoggingService } from "src/app/shared/logging.service";

/**
 * Author contains the information about an author from a cert issuer/subject
 */
export interface Author {
	countryName?: string | undefined;
	stateOrProvince?: string | undefined;
	localityName?: string | undefined;
	orgName?: string | undefined;
	orgUnit?: string | undefined;
	commonName: string;
}

/**
 * Angular validator that checks if the value is either before or after the given time
 *
 * @param before If this is only valid before
 * @param now The time to compare against
 * @returns Validator Function
 */
function createDateValidator(before: boolean, now: Date): ValidatorFn {
	return (control: AbstractControl): ValidationErrors | null => {
		const value = control.value;
		if (!value) {
			return null;
		}

		let valid = false;
		const d = new Date(value);
		if (before) {
			valid = now >= d;
		} else {
			valid = now <= d;
		}

		return valid ? null : {outOfDate: true};
	};
}

/**
 * Controller for the Cert Detail component
 */
@Component({
	selector: "tp-cert-detail",
	styleUrls: ["./cert-detail.component.scss"],
	templateUrl: "./cert-detail.component.html"
})
export class CertDetailComponent implements OnChanges {
	@Input({required: true}) public cert!: pki.Certificate;

	public issuer: Author = {commonName: ""};
	public subject: Author = {commonName: ""};
	public now: Date = new Date();
	public validAfterFormControl = new FormControl<string>("", [createDateValidator(false, this.now)]);
	public validBeforeFormControl = new FormControl<string>("", [createDateValidator(true, this.now)]);

	public sha1: Hex = "";
	public sha256: Hex = "";

	constructor(private readonly log: LoggingService) { }

	/**
	 * processAttributes converts attributes into an author
	 *
	 * @param attrs The attributes to process
	 * @returns The resultant author
	 */
	public processAttributes(attrs: pki.CertificateField[]): Author {
		const a: Author = {commonName: ""};
		for (const attr of attrs) {
			if (attr.name && attr.value) {
				if (typeof attr.value !== "string") {
					this.log.warn(`Unknown attribute value ${attr.value}`);
					continue;
				}
				switch (attr.name) {
					case "commonName":
						a.commonName = attr.value;
						break;
					case "countryName":
						a.countryName = attr.value;
						break;
					case "stateOrProvinceName":
						a.stateOrProvince = attr.value;
						break;
					case "localityName":
						a.localityName = attr.value;
						break;
					case "organizationName":
						a.orgName = attr.value;
						break;
					case "organizationUnitName":
						a.orgUnit = attr.value;
						break;
				}
			}
		}
		return a;
	}

	/**
	 * Calculates certificate details
	 */
	public ngOnChanges(): void {
		this.sha1 = pkiCertToSHA1(this.cert);
		this.sha256 = pkiCertToSHA256(this.cert);
		this.issuer = this.processAttributes(this.cert.issuer.attributes);
		this.subject = this.processAttributes(this.cert.subject.attributes);

		this.validBeforeFormControl.setValue(this.cert.validity.notBefore.toISOString().slice(0, 16));
		this.validAfterFormControl.setValue(this.cert.validity.notAfter.toISOString().slice(0, 16));
		this.validAfterFormControl.markAsTouched();
		this.validBeforeFormControl.markAsTouched();
	}

	protected readonly oidToName = oidToName;
}
