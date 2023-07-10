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
import * as forge from "node-forge";

export namespace Certs {
	type CertType = "Root" | "Client" | "Intermediate" | "Unknown" | "Error";
	type CertOrder = "Client -> Root" | "Root -> Client" | "Unknown" | "Single";

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

	export interface Certificate extends forge.pki.Certificate {
		issuerParsed:  Author;
		subjectParsed: Author;
		type: CertType;
		parseError: boolean;
		sha1: forge.Hex;
		sha256: forge.Hex;
		oidName: string;
		hidden: boolean;
		miscHidden: boolean;
		validityHidden: boolean;
	}
}
