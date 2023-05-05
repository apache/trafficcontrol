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
import { Component, type OnInit } from "@angular/core";
import { FormControl, FormGroup, ValidationErrors, Validators } from "@angular/forms";
import type { ISORequest } from "trafficops-types";

import { MiscAPIsService } from "src/app/api";
import { FileUtilsService } from "src/app/shared/file-utils.service";
import { IPV4, IPV6, IPV6_WITH_CIDR } from "src/app/utils";

/**
 * The controller for a form that can be used to generate ISOs.
 *
 * Ideally this will be removed in the future, but it is not yet deprecated.
 */
@Component({
	selector: "tp-isogeneration-form",
	styleUrls: ["../../styles/form.page.scss", "./isogeneration-form.component.scss"],
	templateUrl: "./isogeneration-form.component.html",
})
export class ISOGenerationFormComponent implements OnInit {
	public osVersions: [string, string][] = [];

	/**
	 * This is the reactive form controller for the page. Due to some complex
	 * validation, as well as spotty browser support for the absolutely MASSIVE
	 * regular expressions some the pattern validators use, this can't reliably
	 * done in a template-driven form.
	 */
	public form = new FormGroup({
		disk: new FormControl("", {nonNullable: true}),
		fqdn: new FormControl("", {
			nonNullable: true,
			validators: Validators.pattern(/^[a-zA-Z0-9][a-zA-Z0-9-]+\.([a-zA-Z0-9][a-zA-Z0-9-]+)$/)
		}),
		interfaceName: new FormControl(""),
		ipv4Address: new FormControl("", {nonNullable: true, validators: Validators.pattern(IPV4)}),
		ipv4Gateway: new FormControl("", {nonNullable: true, validators: Validators.pattern(IPV4)}),
		ipv4Netmask: new FormControl("", {nonNullable: true, validators: Validators.pattern(IPV4)}),
		ipv6Address: new FormControl("", Validators.pattern(IPV6_WITH_CIDR)),
		ipv6Gateway: new FormControl("", Validators.pattern(IPV6)),
		mtu: new FormControl(1500, {nonNullable: true}),
		osVersion: new FormControl("", {nonNullable: true}),
		rootPass: new FormControl("", {nonNullable: true}),
		rootPassConfirm: new FormControl("", {nonNullable: true}),
		useDHCP: new FormControl(true, {nonNullable: true}),
	});

	/** `true` if IPv4 will be dynamic, `false` otherwise. */
	public get useDHCP(): boolean {
		return this.form.controls.useDHCP.value;
	}

	constructor(private readonly api: MiscAPIsService, private readonly fileService: FileUtilsService) {
		this.form.controls.rootPassConfirm.addValidators((ctrl): ValidationErrors | null => {
			if (this.form.controls.rootPass.value !== ctrl.value) {
				return {
					mismatch: true
				};
			}
			return null;
		});
	}

	/** Angular lifecycle hook. */
	public async ngOnInit(): Promise<void> {
		const osVersions = await this.api.getISOOSVersions();
		this.osVersions = Object.entries(osVersions);
	}

	/**
	 * Checks if the warning about an unusual MTU should be hidden.
	 *
	 * @returns `true` if the MTU warning should be hidden, `false` otherwise.
	 */
	public hideMTUWarning(): boolean {
		return this.form.controls.mtu.value === 1500 || this.form.controls.mtu.value === 9000;
	}

	/**
	 * Handles form submission.
	 *
	 * @param event The DOM form submission event.
	 */
	public async submit(event: Event): Promise<void> {
		event.preventDefault();
		event.stopPropagation();
		if (this.form.invalid) {
			return;
		}

		const fqdn = this.form.controls.fqdn.value;
		const [hostName, domainName] = fqdn.split(".", 2);

		let req: ISORequest;
		if (this.useDHCP) {
			req = {
				dhcp: "yes",
				disk: this.form.controls.disk.value,
				domainName,
				hostName,
				interfaceMtu: this.form.controls.mtu.value,
				interfaceName: this.form.controls.interfaceName.value,
				ip6Address: this.form.controls.ipv6Address.value,
				ip6Gateway: this.form.controls.ipv6Gateway.value,
				osVersionDir: this.form.controls.osVersion.value,
				rootPass: this.form.controls.rootPass.value
			};
		} else {
			req = {
				dhcp: "no",
				disk: this.form.controls.disk.value,
				domainName,
				hostName,
				interfaceMtu: this.form.controls.mtu.value,
				interfaceName: this.form.controls.interfaceName.value,
				ip6Address: this.form.controls.ipv6Address.value,
				ip6Gateway: this.form.controls.ipv6Gateway.value,
				ipAddress: this.form.controls.ipv4Address.value,
				ipGateway: this.form.controls.ipv4Gateway.value,
				ipNetmask: this.form.controls.ipv4Netmask.value,
				osVersionDir: this.form.controls.osVersion.value,
				rootPass: this.form.controls.rootPass.value
			};
		}

		const response = await this.api.generateISO(req);
		this.fileService.download(response, `${fqdn}-${this.form.controls.osVersion.value}.iso`);
	}

}
