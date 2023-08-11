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
import { DOCUMENT } from "@angular/common";
import { Component, type OnInit, ViewChild, Inject } from "@angular/core";
import { FormControl } from "@angular/forms";
import type { MatStepper } from "@angular/material/stepper";
import { Router } from "@angular/router";
import {
	bypassable,
	GeoLimit,
	GeoProvider,
	Protocol,
	protocolToString,
	QStringHandling,
	RangeRequestHandling,
	TypeFromResponse,
	type RequestDeliveryService,
	type ResponseCDN
} from "trafficops-types";

import { CDNService, DeliveryServiceService } from "src/app/api";
import { CurrentUserService } from "src/app/shared/current-user/current-user.service";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { IPV4, IPV6 } from "src/app/utils";

/**
 * A regular expression that matches character strings that are illegal in `xml_id`s
 */
const XML_ID_SANITIZE = /[^a-z0-9-]+/g;

/**
 * A regular expression that matches a valid xml_id
 */
const VALID_XML_ID = /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/;

/**
 * A regular expression that matches a valid hostname
 */
const VALID_HOSTNAME = /^[A-z\d]([A-z0-9-]*[A-z0-9])?(\.[A-z\d]([A-z0-9-]*[A-z0-9])?)*$/;

/**
 * NewDeliveryServiceComponent is the controller for the new Delivery Service
 * creation page.
 */
@Component({
	selector: "app-new-delivery-service",
	styleUrls: ["./new-delivery-service.component.scss"],
	templateUrl: "./new-delivery-service.component.html"
})
export class NewDeliveryServiceComponent implements OnInit {

	/** The Delivery Service being created */
	public deliveryService: RequestDeliveryService = {
		active: true,
		anonymousBlockingEnabled: true,
		cacheurl: null,
		cdnId: 1,
		deepCachingType: "NEVER",
		displayName: "",
		dscp: 0,
		geoLimit: GeoLimit.NONE,
		geoProvider: GeoProvider.MAX_MIND,
		httpBypassFqdn: null,
		infoUrl: null,
		ipv6RoutingEnabled: true,
		logsEnabled: true,
		missLat: 0,
		missLong: 0,
		multiSiteOrigin: false,
		orgServerFqdn: "",
		protocol: Protocol.HTTP,
		qstringIgnore: QStringHandling.USE,
		rangeRequestHandling: RangeRequestHandling.NONE,
		regionalGeoBlocking: false,
		remapText: null,
		routingName: "",
		tenantId: 1,
		typeId: 1,
		xmlId: ""
	};

	/** Allows the user to set 'active' */
	public activeImmediately = new FormControl(false);
	/** Allows the user to set 'bypass*' fields */
	public bypassLoc = new FormControl("");
	/** Allows the user to set 'cdn'/'cdnName' */
	public cdnObject = new FormControl<ResponseCDN | null>(null);
	/** Allows the user to set the 'longDesc' */
	public description = new FormControl("");
	/** Allows the user to set 'ipv6Enabled' */
	public disableIPv6 = new FormControl(false);
	/** Allows the user to set the 'displayName'/'xml_id' */
	public displayName = new FormControl("");
	/** Allows the user to set 'type'/'typeId' */
	public dsType = new FormControl<TypeFromResponse | null>(null);
	/** Allows the user to set 'infoUrl' */
	public infoURL = new FormControl("");
	/** Allows the user to set 'originFqdn' */
	public originURL = new FormControl("");
	/** Allows the user to set 'protocol' */
	public protocol = new FormControl(Protocol.HTTP);

	/** Need This to be a property for template access. */
	public readonly protocolToString = protocolToString;

	/** The available Delivery Service protocols. */
	public readonly protocols = [
		Protocol.HTTP,
		Protocol.HTTPS,
		Protocol.HTTP_AND_HTTPS,
		Protocol.HTTP_TO_HTTPS
	];

	/** The available CDNs from which for the user to choose. */
	public cdns: Array<ResponseCDN> = [];

	/**
	 * The available useInTable=delivery_service Types from which for the user
	 * to choose.
	 */
	public dsTypes: Array<TypeFromResponse> = [];

	/** Need public access to models.bypassable in the template. */
	public bypassable = bypassable;

	/**
	 * A reference to the stepper used in the form.
	 */
	@ViewChild("stepper") public stepper!: MatStepper;

	constructor(
		private readonly dsAPI: DeliveryServiceService,
		private readonly cdnAPI: CDNService,
		private readonly auth: CurrentUserService,
		private readonly router: Router,
		private readonly navSvc: NavigationService,
		@Inject(DOCUMENT) private readonly document: Document,
		private readonly log: LoggingService,
	) { }

	/**
	 * Initializes all of the extra data needed to construct a Delivery Service
	 * (types, cdns, etc).
	 */
	public async ngOnInit(): Promise<void> {
		const success = await this.auth.updateCurrentUser();
		if (!success || this.auth.currentUser === null) {
			return;
		}
		this.navSvc.headerTitle.next("New Delivery Service");

		this.deliveryService.tenantId = this.auth.currentUser.tenantId;
		const typeP = await this.dsAPI.getDSTypes().then(
			types => {
				this.dsTypes = types;
				for (const t of types) {
					if (t.name === "HTTP") {
						this.deliveryService.typeId = t.id;
						this.dsType.setValue(t);
						break;
					}
				}
			}
		);
		if (!this.auth.currentUser || !this.auth.currentUser.tenantId) {
			this.log.error("Cannot set default CDN - user has no tenant");
			return typeP;
		}
		const dsP = this.dsAPI.getDeliveryServices().then(
			d => {
				const cdnsInUse = new Map<number, number>();
				for (const ds of d) {
					if (ds.tenantId === this.auth.currentUser?.tenantId) {
						const usedCDNs = cdnsInUse.get(ds.tenantId);
						if (!usedCDNs) {
							cdnsInUse.set(ds.tenantId, 1);
						} else {
							cdnsInUse.set(ds.tenantId, usedCDNs + 1);
						}
					}
				}

				let most = -Infinity;
				let mostId = -1;
				cdnsInUse.forEach( (v: number, k: number) => {
					if (v > most) {
						most = v;
						mostId = k;
					}
				});
				this.setDefaultCDN(mostId);
			}
		);
		await Promise.all([typeP, dsP]);
	}

	/**
	 * Sets the default CDN based on the passed integral, unique identifier.
	 *
	 * @param id The integral, unique identifier of the CDN which is assumed to
	 * be the CDN which contains more of the current user's tenant's Delivery
	 * Services than any other - unless said tenant has no Delivery Services, in
	 * which case it should be `-1` and the selected CDN will be the first CDN
	 * in lexigraphical order by name.
	 */
	private async setDefaultCDN(id: number): Promise<void> {
		const cdns = await this.cdnAPI.getCDNs();
		if (!cdns) {
			throw new Error("no CDNs found in the API");
		}
		this.cdns = [];
		let def;
		for (const cdn of cdns) {
			// this is a special, magic-value CDN that can't have any DSes
			if (cdn.name !== "ALL") {
				this.cdns.push(cdn);
				if (id > 0) {
					if (cdn.id === id) {
						def = cdn;
					}
				} else if (!def || cdn.name < def.name) {
					def = cdn;
				}
			}
		}
		if (this.cdns.length < 1) {
			throw new Error("the only CDN is 'ALL', which cannot be used for Delivery Services");
		}

		if (!def) {
			def = this.cdns[0];
		}
		this.deliveryService.cdnId = def.id;
		this.cdnObject.setValue(def);
	}

	/**
	 * Updates the header text based on the status of the current delivery service
	 */
	public updateDisplayName(): void {
		const value = this.displayName.value ?? "";
		this.navSvc.headerTitle.next(value === "" ? "New Delivery Service" : value);
	}

	/**
	 * When a user submits their origin URL, this parses that out into the
	 * related DS fields.
	 */
	public setOriginURL(): void {
		let url;
		try {
			url = new URL(this.originURL.value ?? "");
		} catch (e) {
			this.log.error("invalid origin URL:", e);
			return;
		}
		this.deliveryService.orgServerFqdn = url.origin;
		if (url.pathname) {
			this.deliveryService.checkPath = url.pathname;
		}

		switch (url.protocol) {
			case "http:":
				this.deliveryService.protocol = Protocol.HTTP_AND_HTTPS;
				this.protocol.setValue(Protocol.HTTP_AND_HTTPS);
				break;
			case "https:":
				this.deliveryService.protocol = Protocol.HTTP_TO_HTTPS;
				this.protocol.setValue(Protocol.HTTP_TO_HTTPS);
				break;
			default:
				this.deliveryService.protocol = Protocol.HTTP_AND_HTTPS;
				this.protocol.setValue(Protocol.HTTP_AND_HTTPS);
				break;
		}

		if (this.activeImmediately.dirty) {
			this.deliveryService.active = this.activeImmediately.value ?? false;
		}

		this.deliveryService.displayName = `Delivery Service for ${url.hostname}`;
		this.displayName.setValue(this.deliveryService.displayName);
		this.stepper.next();
	}

	/**
	 * Sets the metadata of the new Delivery Service. This will attempt to make
	 * an `xml_id` out of the submitted name by casefolding to lower and
	 * replacing special characters with `-`. If that replacement somehow fails,
	 * it will trigger a form validation error and will convert the submitted
	 * name value to a placeholder while making the actual value an empty
	 * string.
	 */
	public setMetaInformation(): void {
		this.deliveryService.displayName = this.displayName.value ?? "";
		this.deliveryService.xmlId = this.deliveryService.displayName.toLocaleLowerCase().replace(XML_ID_SANITIZE, "-");
		if (! VALID_XML_ID.test(this.deliveryService.xmlId)) {
			// According to https://stackoverflow.com/questions/39642547/is-it-possible-to-get-native-element-for-formcontrol
			// this could instead be implemented with a Directive, but that doesn't really seem like
			// any less of a hack to me.
			const nativeDisplayNameElement = this.document.getElementById("displayName") as HTMLInputElement;
			nativeDisplayNameElement.setCustomValidity(
				"Failed to create a unique key from Delivery Service name. Try using less special characters."
			);
			nativeDisplayNameElement.reportValidity();
			nativeDisplayNameElement.value = "";
			nativeDisplayNameElement.setCustomValidity("");
			return;
		}
		this.deliveryService.longDesc = this.description.value;

		if (this.infoURL.value) {
			this.deliveryService.infoUrl = this.infoURL.value;
		}
		this.stepper.next();
	}

	/**
	 * Currently merely sets the CDN to which the Delivery Service shall belong.
	 * In the future, this should probably handle queuing cache assignments and
	 * the like, possibly with advanced controls for things like Traffic Router
	 * DNS and redirects
	 */
	public setInfrastructureInformation(): void {
		this.deliveryService.cdnId = this.cdnObject.value?.id ?? -1;
		if (this.dsType.dirty) {
			this.deliveryService.typeId = this.dsType.value?.id ?? -1;
		}

		if (this.protocol.dirty) {
			this.deliveryService.protocol = this.protocol.value;
		}

		if (this.disableIPv6.dirty) {
			this.deliveryService.ipv6RoutingEnabled = !this.disableIPv6.value;
		}

		if (this.bypassLoc.dirty) {
			switch (this.dsType.value?.name) {
				case "DNS":
				case "DNS_LIVE":
				case "DNS_LIVE_NATNL":
					try {
						this.setDNSBypass(this.bypassLoc.value ?? "");
					} catch (e) {
						this.log.error("failed to set DNS bypass:", e);
						const nativeBypassElement = this.document.getElementById("bypass-loc") as HTMLInputElement;
						nativeBypassElement.setCustomValidity(e instanceof Error ? e.message : String(e));
						nativeBypassElement.reportValidity();
						nativeBypassElement.value = "";
						nativeBypassElement.setCustomValidity("");
						return;
					}
					break;
				default:
					this.deliveryService.httpBypassFqdn = this.bypassLoc.value;
					break;
			}
		}

		this.dsAPI.createDeliveryService(this.deliveryService).then(
			v => {
				this.router.navigate(["/"], {queryParams: {search: encodeURIComponent(v.displayName)}});
			}
		);
	}

	/**
	 * Sets the appropriate bypass location for the new Delivery Service,
	 * assuming it is DNS-routed.
	 *
	 * @param v Represents a Bypass value - either an IP(v4/v6) address or a hostname.
	 * @throws {Error} if `v` is not a valid Bypass value
	 */
	public setDNSBypass(v: string): void {
		if (IPV6.test(v)) {
			this.deliveryService.dnsBypassIp6 = v;
		} else if (IPV4.test(v)) {
			this.deliveryService.dnsBypassIp = v;
		} else if (VALID_HOSTNAME.test(v)) {
			this.deliveryService.dnsBypassCname = v;
		} else {
			throw new Error(`"${v} is not a valid IPv4/IPv6 address or hostname`);
		}
	}

	/**
	 * Allows a user to return to the previous step.
	 */
	public previous(): void {
		this.stepper.previous();
	}
}
