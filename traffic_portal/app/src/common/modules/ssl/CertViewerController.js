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

/** @typedef { import('./CertViewer').Certs } Certs */

let CertViewerController = function ($scope) {
	/** @typedef Certs.Certificate[] */
	this.certChain = [];
	/** @typedef Certs.CertOrder */
	this.certOrder = "";
	/** @typedef string */
	this.chain = "";

	/** @type Certs.Certificate */
	this.nullCert = window.forge.pki.createCertificate();
	this.nullCert.type = "Error";
	this.nullCert.parseError = true;

	this.$onChanges = function() {
		if(this.chain === "") {
			return;
		}
		this.chain = this.chain.replace(/\r\n/g, "\n");
		const parts = this.chain.split("-\n-");
		const certs = new Array(parts.length);
		for (let i = 1; i < parts.length; ++i) {
			parts[i - 1] += "-";
			parts[i] = `-${parts[i]}`;
			certs[i - 1] = this.newCert(parts[i - 1]);
		}
		certs[certs.length - 1] = this.newCert(parts[parts.length - 1]);
		const assignExtraInfo = (c, i) => {
			if (c.parseError) {
				return;
			}
			if (i === 0) {
				c.type = "Root";
			} else if (i === certs.length - 1) {
				c.type = "Client";
			} else {
				c.type = "Intermediate";
			}
			c.sha1 = this.pkiCertToSHA1(c);
			c.sha256 = this.pkiCertToSHA256(c);
			c.parsedIssuer = this.processAttributes(c.issuer.attributes);
			c.parsedSubject = this.processAttributes(c.subject.attributes);
			c.oidName = this.oidToName(c.signatureOid);
		};
		const chain = this.reOrderRootFirst(certs);
		chain.forEach(assignExtraInfo);
		this.certChain = chain;
	}

	/**
	 * Converts a given oid to it's human-readable name
	 */
	this.oidToName = function(oid) {
		if (oid in window.forge.pki.oids) {
			return window.forge.pki.oids[oid];
		}
		return "";
	}

	/**
	 * processAttributes converts attributes into an author
	 */
	this.processAttributes = function(attrs) {
		/** @typedef Certs.Author */
		const a = {commonName: ""};
		for (const attr of attrs) {
			if (attr.name && attr.value) {
				if (typeof attr.value !== "string") {
					console.warn(`Unknown attribute value ${attr.value}`);
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
	 * pkiCertToSHA1 calculates the sha 1 of a cert
	 */
	this.pkiCertToSHA1 = function (cert) {
		const md = window.forge.md.sha1.create();
		md.update(forge.asn1.toDer(window.forge.pki.certificateToAsn1(cert)).getBytes());
		return md.digest().toHex();
	}

	/**
	 * pkiCertToSHA256 calculates the sha 256 of a cert
	 */
	this.pkiCertToSHA256 = function (cert) {
		const md = forge.md.sha256.create();
		md.update(forge.asn1.toDer(window.forge.pki.certificateToAsn1(cert)).getBytes());
		return md.digest().toHex();
	}

	/**
	 * newCert creates a cert from an input string.
	 */
	this.newCert = function (input) {
		try {
			return forge.pki.certificateFromPem(input);
		} catch (e) {
			console.error(`ran into issue creating certificate from input ${input}`, e);
			return this.nullCert;
		}
	}

	/**
	 * reOrderRootFirst sorts a cert chain with the root being first if possible.
	 */
	this.reOrderRootFirst = function (certs) {
		let rootFirst = false;
		let invalid = false;
		for (let i = 1; i < certs.length; ++i) {
			const first = certs[i - 1];
			const next = certs[i];
			if (first.parseError) {
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
				console.error(`Cert chain is invalid, cert ${i - 1} and ${i} are not related`);
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

		if (rootFirst) {
			this.certOrder = "Root -> Client";
			return certs;
		}
		this.certOrder = "Client -> Root";
		certs = certs.reverse();
		return certs;
	}
}

angular.module("trafficPortal.ssl").component("certViewerController", {
	templateUrl: "common/modules/ssl/cert-view.tpl.html",
	controller: CertViewerController,
	bindings: {
		chain: "@"
	}
});

CertViewerController.$inject = ["$scope"];
module.exports = CertViewerController;
