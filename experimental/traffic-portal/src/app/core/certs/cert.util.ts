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
import { pki, md, asn1 } from "node-forge";

/**
 * Converts a given oid to it's human-readable name
 *
 * @param oid The oid to translate
 * @returns The human-readable oid
 */
export function oidToName(oid: string): string {
	if (oid in pki.oids) {
		return pki.oids[oid];
	}
	return "";
}

/**
 * Calculate the SHA-1 of a given certificate
 *
 * @param cert The certificate to hash
 * @returns SHA-1 of the cert
 */
export function pkiCertToSHA1(cert: pki.Certificate): string {
	const md1 = md.sha1.create();
	md1.update(asn1.toDer(pki.certificateToAsn1(cert)).getBytes());
	return md1.digest().toHex();
}

/**
 * Calculate the SHA-256 of a given certificate
 *
 * @param cert The certificate to hash
 * @returns SHA-256 of the cert
 */
export function pkiCertToSHA256(cert: pki.Certificate): string {
	const md256 = md.sha256.create();
	md256.update(asn1.toDer(pki.certificateToAsn1(cert)).getBytes());
	return md256.digest().toHex();
}
