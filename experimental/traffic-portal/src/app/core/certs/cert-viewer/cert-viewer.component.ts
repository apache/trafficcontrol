import { Component, OnInit, ViewChild } from "@angular/core";
import { MatTabGroup } from "@angular/material/tabs";
import { ActivatedRoute } from "@angular/router";
import * as forge from "node-forge";
import { ResponseDeliveryServiceSSLKey } from "trafficops-types";

import { DeliveryServiceService } from "src/app/api";

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
export interface AugmentedCertificate extends forge.pki.Certificate {
	type: CertType;
	parseError: boolean;
}

export const NULL_CERT = forge.pki.createCertificate() as AugmentedCertificate;
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

	@ViewChild("matTab") public matTab!: MatTabGroup;
	constructor(
		private readonly route: ActivatedRoute,
		private readonly dsAPI: DeliveryServiceService) {
	}

	/**
	 * newCert creates a cert from an input string.
	 *
	 * @param input The text to read as a cert
	 * @private
	 * @returns Resultant Cert
	 */
	private newCert(input: string): AugmentedCertificate {
		try {
			return forge.pki.certificateFromPem(input) as AugmentedCertificate;
		} catch (e) {
			console.error(`ran into issue creating certificate from input ${input}`, e);
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
				console.error(`Cert chain is invalid, cert ${i-1} and ${i} are not related`);
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
		this.cert = await this.dsAPI.getSSLKeys(ID);
		this.dsCert = true;
		this.inputCert = this.cert.certificate.crt;
		this.process();
	}

}