/**
 * @license Apache-2.0
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
import { Injectable } from "@angular/core";
import type { OSVersions } from "trafficops-types";

/**
 * This service implements APIs that aren't specific to any given ATC object.
 * They can provide things like system information, access to the wacky external
 * tooling provided by TO, basically anything that doesn't fit in a different
 * API service.
 */
@Injectable()
export class MiscAPIsService {
	/** Some static mock OS versions. */
	public readonly osVersions = {
		// eslint-disable-next-line @typescript-eslint/naming-convention
		"CentOS 7": "centos7",
		// eslint-disable-next-line @typescript-eslint/naming-convention
		"Rocky Linux 8": "rocky8"
	};

	/**
	 * Retrieves the operating system versions that can be used to generate
	 * system images through the Traffic Ops API.
	 *
	 * @returns A mapping of human-friendly operating system names to
	 * machine-readable OS IDs that can be used in subsequent requests to
	 * {@link MiscAPIsService.generateISO}.
	 */
	public async getISOOSVersions(): Promise<OSVersions> {
		return this.osVersions;
	}

	/**
	 * A mock call for generating system images.
	 *
	 * @returns In tests, this returns an empty data blob no matter what you
	 * pass it.
	 */
	public async generateISO(): Promise<Blob> {
		return new Blob();
	}
}
