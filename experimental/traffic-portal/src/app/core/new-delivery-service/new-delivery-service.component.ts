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
import type { MatStepper } from "@angular/material/stepper";
import { Router } from "@angular/router";

import {CurrentUserService} from "src/app/shared/currentUser/current-user.service";
import {
	bypassable,
	CDN,
	defaultDeliveryService,
	DeliveryService,
	Protocol,
	protocolToString,
	Type
} from "../../models";
import { User } from "../../models/user";
import { CDNService, DeliveryServiceService } from "../../shared/api";

/**
 * A regular expression that matches character strings that are illegal in `xml_id`s
 */
const XML_ID_SANITIZE = /[^a-z0-9\-]+/g;

/**
 * A regular expression that matches a valid xml_id
 */
const VALID_XML_ID = /^[a-z0-9]([a-z0-9\-]*[a-z0-9])?$/;

/* eslint-disable */
/**
 * A regular expression that matches IPv4 addresses
 */
const VALID_IPV4 = /^(1\d\d|2[0-4]\d|25[0-5]|\d\d?)(\.(1\d\d|2[0-4]\d|25[0-5]|\d\d?)){3}$/;
/**
 * A regular expression that matches IPv6 addresses
 * This is huge and ugly, but there's no JS built-in for address parsing afaik.
 */
const VALID_IPV6 = /^((((((([\da-fA-F]{1,4})):){6})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|((::((([\da-fA-F]{1,4})):){5})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|((((([\da-fA-F]{1,4}))):((([\da-fA-F]{1,4})):){4})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,1}(([\da-fA-F]{1,4}))):((([\da-fA-F]{1,4})):){3})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,2}(([\da-fA-F]{1,4}))):((([\da-fA-F]{1,4})):){2})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,3}(([\da-fA-F]{1,4}))):(([\da-fA-F]{1,4})):)((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,4}(([\da-fA-F]{1,4}))):)((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,5}(([\da-fA-F]{1,4}))):)(([\da-fA-F]{1,4})))|(((((([\da-fA-F]{1,4})):){0,6}(([\da-fA-F]{1,4}))):))))$/;
/* eslint-enable */
/**
 * A regular expression that matches a valid hostname
 */
const VALID_HOSTNAME = /^[A-z\d]([A-z0-9\-]*[A-z0-9])*(\.[A-z\d]([A-z0-9\-]*[A-z0-9])*)*$/;

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
	public deliveryService: DeliveryService = {...defaultDeliveryService};

	/** Allows the user to set 'active' */
	public activeImmediately = new FormControl();
	/** Allows the user to set 'bypass*' fields */
	public bypassLoc = new FormControl("");
	/** Allows the user to set 'cdn'/'cdnName' */
	public cdnObject = new FormControl("");
	/** Allows the user to set the 'longDesc' */
	public description = new FormControl("");
	/** Allows the user to set 'ipv6Enabled' */
	public disableIPv6 = new FormControl();
	/** Allows the user to set the 'displayName'/'xml_id' */
	public displayName = new FormControl("");
	/** Allows the user to set 'type'/'typeId' */
	public dsType = new FormControl();
	/** Allows the user to set 'infoUrl' */
	public infoURL = new FormControl("");
	/** Allows the user to set 'originFqdn' */
	public originURL = new FormControl("");
	/** Allows the user to set 'protocol' */
	public protocol = new FormControl();

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
	public cdns: Array<CDN> = [];

	/**
	 * The available useInTable=delivery_service Types from which for the user
	 * to choose.
	 */
	public dsTypes: Array<Type> = [];

	/** Need public access to models.bypassable in the template. */
	public bypassable = bypassable;

	/**
	 * A reference to the stepper used in the form - this is never actually
	 * undefined, but the value is set by the decorator, so to satisfy the
	 * compiler we need to mark it as optional.
	 */
	@ViewChild("stepper") public stepper?: MatStepper;

	/**
	 * Constructor.
	 */
	constructor(
		private readonly dsAPI: DeliveryServiceService,
		private readonly cdnAPI: CDNService,
		private readonly auth: CurrentUserService,
		private readonly router: Router
	) { }

	/**
	 * Initializes all of the extra data needed to construct a Delivery Service
	 * (types, cdns, etc).
	 */
	public ngOnInit(): void {
		this.auth.updateCurrentUser().then( success => {
			if (!success || this.auth.currentUser === null) {
				return;
			}
			this.deliveryService.tenant = this.auth.currentUser.tenant;
			this.deliveryService.tenantId = this.auth.currentUser.tenantId;
			this.dsAPI.getDSTypes().then(
				(types: Array<Type>) => {
					this.dsTypes = types;
					for (const t of types) {
						if (t.name === "HTTP") {
							this.deliveryService.type = t.name;
							this.deliveryService.typeId = t.id;
							this.dsType.setValue(t);
							break;
						}
					}
				}
			);
			if (!this.auth.currentUser || !this.auth.currentUser.tenantId) {
				console.error("Cannot set default CDN - user has no tenant");
				return;
			}
			this.dsAPI.getDeliveryServices().then(
				d => {
					const cdnsInUse = new Map<number, number>();
					for (const ds of d) {
						if (ds.tenantId === undefined) {
							console.warn("Delivery Service has no tenant:", ds);
							continue;
						}
						if (ds.tenantId === (this.auth.currentUser as User).tenantId) {
							const usedCDNs = cdnsInUse.get(ds.tenantId);
							if (!usedCDNs) {
								cdnsInUse.set(ds.tenantId, 1);
							} else {
								cdnsInUse.set(ds.tenantId, usedCDNs + 1);
							}
						}
					}

					let most = 0;
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

		});
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
	private setDefaultCDN(id: number): void {
		this.cdnAPI.getCDNs().then(
			cdns => {
				if (!cdns) {
					console.warn("No CDNs found in the API");
					return;
				}
				this.cdns = new Array<CDN>();
				let def: CDN | null = null;
				cdns.forEach( (c: CDN) => {

					// this is a special, magic-value CDN that can't have any DSes
					if (c.name !== "ALL") {
						this.cdns.push(c);
						if (id > 0) {
							if (c.id === id) {
								def = c;
							}
						} else if (!def || c.name < def.name) {
							def = c;
						}
					}
				});
				if (!def) {
					def = this.cdns[0];
				}
				this.deliveryService.cdnId = def.id;
				this.deliveryService.cdnName = def.name;
				this.cdnObject.setValue(def);
			}
		);
	}

	/**
	 * When a user submits their origin URL, this parses that out into the
	 * related DS fields.
	 *
	 * @param stepper The stepper controlling the form flow - used to advance to the next step.
	 */
	public setOriginURL(): void {
		const parser = document.createElement("a");
		parser.href = this.originURL.value;
		this.deliveryService.orgServerFqdn = parser.origin;
		if (parser.pathname) {
			this.deliveryService.checkPath = parser.pathname;
		}

		switch (parser.protocol) {
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
			this.deliveryService.active = this.activeImmediately.value;
		}

		this.deliveryService.displayName = `Delivery Service for ${parser.hostname}`;
		this.displayName.setValue(this.deliveryService.displayName);
		this.stepper?.next();
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
		this.deliveryService.displayName = this.displayName.value;
		this.deliveryService.xmlId = this.deliveryService.displayName.toLocaleLowerCase().replace(XML_ID_SANITIZE, "-");
		if (! VALID_XML_ID.test(this.deliveryService.xmlId)) {
			// According to https://stackoverflow.com/questions/39642547/is-it-possible-to-get-native-element-for-formcontrol
			// this could instead be implemented with a Directive, but that doesn't really seem like
			// any less of a hack to me.
			const nativeDisplayNameElement = document.getElementById("displayName") as HTMLInputElement;
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
		this.stepper?.next();
	}

	/**
	 * Currently merely sets the CDN to which the Delivery Service shall belong.
	 * In the future, this should probably handle queuing cache assignments and
	 * the like, possibly with advanced controls for things like Traffic Router
	 * DNS and redirects
	 */
	public setInfrastructureInformation(): void {
		this.deliveryService.cdnName = this.cdnObject.value.name;
		this.deliveryService.cdnId = this.cdnObject.value.id;
		if (this.dsType.dirty) {
			this.deliveryService.typeId = this.dsType.value.id;
			this.deliveryService.type = this.dsType.value.name;
		}

		if (this.protocol.dirty) {
			this.deliveryService.protocol = this.protocol.value;
		}

		if (this.disableIPv6.dirty) {
			this.deliveryService.ipv6RoutingEnabled = !this.disableIPv6.value;
		}

		if (this.bypassLoc.dirty) {
			switch (this.deliveryService.type) {
				case "DNS":
				case "DNS_LIVE":
				case "DNS_LIVE_NATNL":
					try {
						this.setDNSBypass(this.bypassLoc.value);
					} catch (e) {
						console.error(e);
						const nativeBypassElement = document.getElementById("bypass-loc") as HTMLInputElement;
						nativeBypassElement.setCustomValidity(e.message);
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
				if (v) {
					console.log("New Delivery Service '%s' created", this.deliveryService.displayName);
					this.router.navigate(["/"], {queryParams: {search: encodeURIComponent(this.deliveryService.displayName)}});
				} else {
					console.error("Failed to create deliveryService: ", this.deliveryService);
				}
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
		if (VALID_IPV6.test(v)) {
			this.deliveryService.dnsBypassIp6 = v;
		} else if (VALID_IPV4.test(v)) {
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
		this.stepper?.previous();
	}
}
