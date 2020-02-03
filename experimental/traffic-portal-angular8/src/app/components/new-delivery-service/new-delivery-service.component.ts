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
import { Component, ElementRef, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Router } from '@angular/router';
import { first } from 'rxjs/operators';

import { APIService, AuthenticationService } from '../../services';

import { bypassable, CDN, DeliveryService, GeoLimit, GeoProvider, Protocol, QStringHandling, RangeRequestHandling, Type } from '../../models';

/**
 * A regular expression that matches character strings that are illegal in `xml_id`s
*/
const XML_ID_SANITIZE = /[^a-z0-9\-]+/g;

/**
 * A regular expression that matches a valid xml_id
*/
const VALID_XML_ID = /^[a-z0-9]([a-z0-9\-]*[a-z0-9])?$/;

/**
 * A regular expression that matches IPv4 addresses
*/
const VALID_IPv4 = /^(1\d\d|2[0-4]\d|25[0-5]|\d\d?)(\.(1\d\d|2[0-4]\d|25[0-5]|\d\d?)){3}$/;
/** tslint:disable **/
/**
 * A regular expression that matches IPv6 addresses
 * This is huge and ugly, but there's no JS built-in for address parsing afaik.
*/
const VALID_IPv6 = /^((((((([\da-fA-F]{1,4})):){6})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|((::((([\da-fA-F]{1,4})):){5})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|((((([\da-fA-F]{1,4}))):((([\da-fA-F]{1,4})):){4})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,1}(([\da-fA-F]{1,4}))):((([\da-fA-F]{1,4})):){3})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,2}(([\da-fA-F]{1,4}))):((([\da-fA-F]{1,4})):){2})((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,3}(([\da-fA-F]{1,4}))):(([\da-fA-F]{1,4})):)((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,4}(([\da-fA-F]{1,4}))):)((((([\da-fA-F]{1,4})):(([\da-fA-F]{1,4})))|(((((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d]))\.){3}((25[0-5]|([1-9]|1[\d]|2[0-4])?[\d])))))))|(((((([\da-fA-F]{1,4})):){0,5}(([\da-fA-F]{1,4}))):)(([\da-fA-F]{1,4})))|(((((([\da-fA-F]{1,4})):){0,6}(([\da-fA-F]{1,4}))):))))$/;
/** tslint:enable **/
/**
 * A regular expression that matches a valid hostname
*/
const VALID_HOSTNAME = /^[A-z\d]([A-z0-9\-]*[A-z0-9])*(\.[A-z\d]([A-z0-9\-]*[A-z0-9])*)*$/;

@Component({
	selector: 'app-new-delivery-service',
	templateUrl: './new-delivery-service.component.html',
	styleUrls: ['./new-delivery-service.component.scss']
})
export class NewDeliveryServiceComponent implements OnInit {

	/** The Delivery Service being created */
	deliveryService = {} as DeliveryService;

	/** A bunch of form controls */
	activeImmediately = new FormControl();
	bypassLoc = new FormControl('');
	cdnObject = new FormControl('');
	description = new FormControl('');
	disableIPv6 = new FormControl();
	displayName = new FormControl('');
	dsType = new FormControl();
	infoURL = new FormControl('');
	originURL = new FormControl('');
	protocol = new FormControl();

	Protocol = Protocol;

	cdns: Array<CDN>;
	dsTypes: Array<Type>;

	step = 0;

	constructor (private readonly api: APIService, private readonly auth: AuthenticationService, private readonly router: Router) { }

	ngOnInit () {
		this.auth.updateCurrentUser().subscribe( success => {
			if (!success || this.auth.currentUserValue === null) {
				return;
			}
			this.deliveryService.active = false;
			this.deliveryService.anonymousBlockingEnabled = false;
			this.deliveryService.deepCachingType = 'NEVER';
			this.deliveryService.dscp = 0;
			this.deliveryService.geoLimit = GeoLimit.None;
			this.deliveryService.geoProvider = GeoProvider.MaxMind;
			this.deliveryService.logsEnabled = true;
			this.deliveryService.initialDispersion = 1;
			this.deliveryService.ipv6RoutingEnabled = true;
			this.deliveryService.missLat = 0.0;
			this.deliveryService.missLong = 0.0;
			this.deliveryService.multiSiteOrigin = false;
			this.deliveryService.qstringIgnore = QStringHandling.USE;
			this.deliveryService.rangeRequestHandling = RangeRequestHandling.NONE;
			this.deliveryService.regionalGeoBlocking = false;
			this.deliveryService.tenant = this.auth.currentUserValue.tenant;
			this.deliveryService.tenantId = this.auth.currentUserValue.tenantId;
			this.api.getDSTypes().pipe(first()).subscribe(
				(types: Array<Type>) => {
					this.dsTypes = types;
					for (const t of types) {
						if (t.name === 'HTTP') {
							this.deliveryService.type = t.name;
							this.deliveryService.typeId = t.id;
							this.dsType.setValue(t);
							break;
						}
					}
				}
			);
			this.api.getDeliveryServices().pipe(first()).subscribe(
				d => {
					const cdnsInUse = new Map<number, number>();
					for (const ds of d) {
						if (ds.tenantId === this.auth.currentUserValue.tenantId) {
							if (!cdnsInUse.get(ds.tenantId)) {
								cdnsInUse.set(ds.tenantId, 1);
							} else {
								cdnsInUse.set(ds.tenantId, cdnsInUse.get(ds.tenantId) + 1);
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
	 * @param id The integral, unique identifier of the CDN which is assumed to be the CDN which
	 * contains more of the current user's tenant's Delivery Services than any other - unless said
	 * tenant has no Delivery Services, in which case it should be `-1` and the selected CDN will be
	 * the first CDN in lexigraphical order by name.
	*/
	private setDefaultCDN (id: number) {
		this.api.getCDNs().pipe(first()).subscribe(
			(cdns: Map<string, CDN>) => {
				this.cdns = new Array<CDN>();
				let def: CDN;
				cdns.forEach( (c: CDN, name: string) => {

					// this is a special, magic-value CDN that can't have any DSes
					if (c.name !== 'ALL') {
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

	public bypassable = bypassable;

	/**
	 * When a user submits their origin URL, this parses that out into the related DS fields
	*/
	setOriginURL () {
		const parser = document.createElement('a') as HTMLAnchorElement;
		parser.href = this.originURL.value;
		this.deliveryService.orgServerFqdn = parser.origin;
		if (parser.pathname) {
			this.deliveryService.checkPath = parser.pathname;
		}

		switch (parser.protocol) {
			case 'http:':
				this.deliveryService.protocol = Protocol.HTTP_AND_HTTPS;
				this.protocol.setValue(Protocol.HTTP_AND_HTTPS);
				break;
			case 'https:':
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

		this.deliveryService.displayName = 'Delivery Service for ' + parser.hostname;
		this.displayName.setValue(this.deliveryService.displayName);
		++this.step;
	}

	/**
	 * Sets the metadata of the new Delivery Service. This will attempt to make an `xml_id` out of
	 * the submitted name by casefolding to lower and replacing special characters with `-`. If that
	 * replacement somehow fails, it will trigger a form validation error and will convert the
	 * submitted name value to a placeholder while making the actual value an empty string.
	*/
	setMetaInformation () {
		this.deliveryService.displayName = this.displayName.value;
		this.deliveryService.xmlId = this.deliveryService.displayName.toLocaleLowerCase().replace(XML_ID_SANITIZE, '-');
		if (! VALID_XML_ID.test(this.deliveryService.xmlId)) {
			// According to https://stackoverflow.com/questions/39642547/is-it-possible-to-get-native-element-for-formcontrol
			// this could instead be implemented with a Directive, but that doesn't really seem like
			// any less of a hack to me.
			const nativeDisplayNameElement = document.getElementById('displayName') as HTMLInputElement;
			nativeDisplayNameElement.setCustomValidity(
				'Failed to create a unique key from Delivery Service name. Try using less special characters.'
			);
			nativeDisplayNameElement.reportValidity();
			nativeDisplayNameElement.value = '';
			nativeDisplayNameElement.setCustomValidity('');
			return;
		}
		this.deliveryService.longDesc = this.description.value;

		if (this.infoURL.value) {
			this.deliveryService.infoUrl = this.infoURL.value;
		}
		++this.step;
	}

	/**
	 * Currently merely sets the CDN to which the Delivery Service shall belong. In the future, this
	 * should probably handle queuing cache assignments and the like, possibly with advanced
	 * controls for things like Traffic Router DNS and redirects
	*/
	setInfrastructureInformation () {
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
				case 'DNS':
				case 'DNS_LIVE':
				case 'DNS_LIVE_NATNL':
					try {
						this.setDNSBypass(this.bypassLoc.value);
					} catch (e) {
						console.error(e);
						const nativeBypassElement = document.getElementById('bypass-loc') as HTMLInputElement;
						nativeBypassElement.setCustomValidity(e.message);
						nativeBypassElement.reportValidity();
						nativeBypassElement.value = '';
						nativeBypassElement.setCustomValidity('');
						return;
					}
					break;
				default:
					this.deliveryService.httpBypassFqdn = this.bypassLoc.value;
					break;
			}
		}

		this.api.createDeliveryService(this.deliveryService).pipe(first()).subscribe(
			v => {
				if (v) {
					console.log("New Delivery Service '%s' created", this.deliveryService.displayName);
					this.router.navigate(['/'], {queryParams: {search: encodeURIComponent(this.deliveryService.displayName)}});
				} else {
					console.error('Failed to create deliveryService: ', this.deliveryService);
				}
			}
		);
	}

	/**
	 * Sets the appropriate bypass location for the new Delivery Service, assuming it is DNS-routed
	 * @param v Represents a Bypass value - either an IP(v4/v6) address or a hostname.
	 * @throws {Error} if `v` is not a valid Bypass value
	 */
	setDNSBypass (v: string) {
		if (VALID_IPv6.test(v)) {
			this.deliveryService.dnsBypassIp6 = v;
		} else if (VALID_IPv4.test(v)) {
			this.deliveryService.dnsBypassIp = v;
		} else if (VALID_HOSTNAME.test(v)) {
			this.deliveryService.dnsBypassCname = v;
		} else {
			throw new Error("'" + v + "' is not a valid IPv4/IPv6 address or hostname!");
		}
	}

	/**
	 * Allows a user to return to the previous step
	**/
	previous () {
		--this.step;
	}
}
