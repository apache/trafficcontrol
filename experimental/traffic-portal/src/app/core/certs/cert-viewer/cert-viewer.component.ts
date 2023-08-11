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
import { Component, OnInit, ViewChild } from "@angular/core";
import { FormControl } from "@angular/forms";
import { MatTabGroup } from "@angular/material/tabs";
import { ActivatedRoute, Router } from "@angular/router";
import { pki } from "node-forge";
import { type ResponseDeliveryServiceSSLKey } from "trafficops-types";

import { DeliveryServiceService } from "src/app/api";
import { LoggingService } from "src/app/shared/logging.service";

/**
 * What type of cert is it
 */
type CertType = "Root" | "Client" | "Intermediate" | "Unknown" | "Error";

/**
 * Detected order of the cert chain
 */
type CertOrder = "Client -> Root" | "Root -> Client" | "Unknown" | "Single";

/**
 * Wrapper around Certificate that contains additional fields
 */
export interface AugmentedCertificate extends pki.Certificate {
	type: CertType;
	parseError: boolean;
}

export const NULL_CERT = pki.createCertificate() as AugmentedCertificate;
NULL_CERT.type = "Error";
NULL_CERT.parseError = true;

/**
 * Controller for the Cert Viewer component.
 */
@Component({
	selector: "tp-cert-viewer",
	styleUrls: ["./cert-viewer.component.scss"],
	templateUrl: "./cert-viewer.component.html"
})
export class CertViewerComponent implements OnInit {
	public cert!: ResponseDeliveryServiceSSLKey;
	public inputCert = "";
	public dsCert = false;

	public certChain: Array<AugmentedCertificate> = [];
	public certOrder: CertOrder | undefined;
	public privateKeyFormControl = new FormControl("");

	@ViewChild("matTab") public matTab!: MatTabGroup;
	constructor(
		private readonly route: ActivatedRoute,
		private readonly dsAPI: DeliveryServiceService,
		private readonly router: Router,
		private readonly log: LoggingService,
	) { }

	/**
	 * newCert creates a cert from an input string.
	 *
	 * @param input The text to read as a cert
	 * @private
	 * @returns Resultant Cert
	 */
	private newCert(input: string): AugmentedCertificate {
		try {
			return pki.certificateFromPem(input) as AugmentedCertificate;
		} catch (e) {
			this.log.error(`ran into issue creating certificate from input ${input}`, e);
			return NULL_CERT;
		}
	}

	/**
	 * process takes the Cert Chain text input and parses it.
	 *
	 * @param uploaded if the certificate was uploaded by the client.
	 */
	public process(uploaded: boolean = false): void {
		this.inputCert = this.inputCert.replace(/\r\n/g, "\n");
		const parts = this.inputCert.split("-\n-");
		const certs = new Array<AugmentedCertificate>(parts.length);
		for(let i = 1; i < parts.length; ++i) {
			parts[i-1] += "-";
			parts[i] = `-${parts[i]}`;
			certs[i-1] = this.newCert(parts[i - 1]);
		}
		certs[certs.length-1] = this.newCert(parts[parts.length - 1]);
		const assignType = (c: AugmentedCertificate, i: number): void => {
			if(c.parseError) {
				return;
			}
			if (i === 0) {
				c.type = "Root";
			} else if (i === certs.length - 1) {
				c.type = "Client";
			} else {
				c.type = "Intermediate";
			}
		};
		const chain = this.reOrderRootFirst(certs);
		chain.forEach(assignType);
		this.certChain = chain;

		if(this.matTab && uploaded) {
			this.matTab.selectedIndex = 1;
		}
	}

	/**
	 * reOrderRootFirst sorts a cert chain with the root being first if possible.
	 *
	 * @param certs The list of certs to reorder
	 * @returns The processed certs
	 */
	public reOrderRootFirst(certs: Array<AugmentedCertificate>): Array<AugmentedCertificate> {
		let rootFirst = false;
		let invalid = false;
		for(let i = 1; i < certs.length; ++i){
			const first = certs[i-1];
			const next = certs[i];
			if(first.parseError) {
				invalid = true;
				continue;
			} else if (next.parseError) {
				invalid = true;
				continue;
			}
			if (first.issued(next)) {
				rootFirst = true;
			} else if (next.issued(first)) {
				rootFirst = false;
			} else {
				invalid = true;
				this.log.error(`Cert chain is invalid, cert ${i-1} and ${i} are not related`);
			}
		}

		if (certs.length === 1) {
			if (certs[0].parseError) {
				invalid = true;
			} else {
				this.certOrder = "Single";
				return certs;
			}
		}
		if (invalid) {
			this.certOrder = "Unknown";
			return certs;
		}

		if(rootFirst) {
			this.certOrder = "Root -> Client";
			return certs;
		}
		this.certOrder = "Client -> Root";
		certs = certs.reverse();
		return certs;
	}

	/**
	 * Checks if we are a DS cert or any user provided cert.
	 */
	public async ngOnInit(): Promise<void> {
		const ID = this.route.snapshot.paramMap.get("xmlId");
		if (ID === null) {
			this.dsCert = false;
			return;
		}
		try {
			this.cert = await this.dsAPI.getSSLKeys(ID);
		} catch {
			await this.router.navigate(["/core/certs/ssl/"]);
			return;
		}
		this.dsCert = true;
		this.inputCert = this.cert.certificate.crt;
		this.privateKeyFormControl.setValue(this.cert.certificate.key);
		this.process();
	}

}
