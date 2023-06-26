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
import { MatDialog } from "@angular/material/dialog";
import { serviceAddresses, type ISORequest, type ResponseServer } from "trafficops-types";

import { MiscAPIsService, ServerService } from "src/app/api";
import {
	CollectionChoiceDialogComponent,
	CollectionChoiceDialogData
} from "src/app/shared/dialogs/collection-choice-dialog/collection-choice-dialog.component";
import { FileUtilsService } from "src/app/shared/file-utils.service";
import { IPV4, IPV6, IPV6_WITH_CIDR } from "src/app/utils";

/**
 * The controller for a form that can be used to generate ISOs.
 *
 * Ideally this will be removed in the future, but it is not yet deprecated.
 */
@Component({
	selector: "tp-isogeneration-form",
	styleUrls: ["./isogeneration-form.component.scss"],
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
		mgmtInterface: new FormControl(""),
		mgmtIpAddress: new FormControl("", Validators.pattern(IPV4)),
		mgmtIpGateway: new FormControl("", Validators.pattern(IPV4)),
		mgmtIpNetmask: new FormControl("", Validators.pattern(IPV4)),
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

	constructor(
		private readonly api: MiscAPIsService,
		private readonly fileService: FileUtilsService,
		private readonly dialog: MatDialog,
		private readonly serverAPI: ServerService
	) {
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
				mgmtInterface: this.form.controls.mgmtInterface.value,
				mgmtIpAddress: this.form.controls.mgmtIpAddress.value,
				mgmtIpGateway: this.form.controls.mgmtIpGateway.value,
				mgmtIpNetmask: this.form.controls.mgmtIpNetmask.value,
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
				mgmtInterface: this.form.controls.mgmtInterface.value,
				mgmtIpAddress: this.form.controls.mgmtIpAddress.value,
				mgmtIpGateway: this.form.controls.mgmtIpGateway.value,
				mgmtIpNetmask: this.form.controls.mgmtIpNetmask.value,
				osVersionDir: this.form.controls.osVersion.value,
				rootPass: this.form.controls.rootPass.value
			};
		}

		const response = await this.api.generateISO(req);
		this.fileService.download(response, `${fqdn}-${this.form.controls.osVersion.value}.iso`);
	}

	/**
	 * Copies information from a server to fill in form information.
	 *
	 * @param server The server from which to copy data.
	 */
	private copyServerData(server: ResponseServer): void {
		const inf = ServerService.getServiceInterface(server);
		this.form.controls.useDHCP.setValue(false);
		this.form.controls.fqdn.setValue(`${server.hostName}.${server.domainName}`);
		this.form.controls.interfaceName.setValue(inf.name);
		this.form.controls.mgmtIpAddress.setValue(server.mgmtIpAddress);
		this.form.controls.mgmtIpGateway.setValue(server.mgmtIpGateway);
		this.form.controls.mgmtIpNetmask.setValue(server.mgmtIpNetmask);
		if (inf.mtu) {
			this.form.controls.mtu.setValue(inf.mtu);
		}

		const [ipv4, ipv6] = serviceAddresses([inf]);
		if (ipv4) {
			if (ipv4.gateway) {
				this.form.controls.ipv4Gateway.setValue(ipv4.gateway);
			}

			const [addr, mask] = ServerService.extractNetmask(ipv4);
			this.form.controls.ipv4Address.setValue(addr);
			if (mask) {
				this.form.controls.ipv4Netmask.setValue(mask);
			}
		}

		if (ipv6) {
			if (ipv6.gateway) {
				this.form.controls.ipv6Gateway.setValue(ipv6.gateway);
			}
			this.form.controls.ipv6Address.setValue(ipv6.address);
		}
	}

	/**
	 * Opens a dialog for copying the information stored on a Traffic Ops Server
	 * object to fill in form fields.
	 */
	public async openCopyDialog(): Promise<void> {
		const collection = (await this.serverAPI.getServers()).map(
			s => ({
				label: `${s.hostName}.${s.domainName}`,
				value: s
			})
		);

		const data = {
			collection,
			message: "Select a server from which to copy ISO generation information.",
			title: "Copy Attributes from Server"
		};
		const d = this.dialog.open<
		CollectionChoiceDialogComponent<ResponseServer>,
		CollectionChoiceDialogData<ResponseServer>,
		ResponseServer
		>(
			CollectionChoiceDialogComponent, {data}
		);

		const selected = await d.afterClosed().toPromise();
		if (selected) {
			this.copyServerData(selected);
		}
	}

}
