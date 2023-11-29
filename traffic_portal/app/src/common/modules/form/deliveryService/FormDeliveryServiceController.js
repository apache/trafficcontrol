/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * This is the parent controller for all kinds of Delivery Service forms - edit,
 * creation, request, etc.
 *
 * @param {import("../../../api/DeliveryServiceService").DeliveryService} deliveryService
 * @param {import("../../../api/DeliveryServiceService").DeliveryService | undefined} dsCurrent
 * @param {unknown} origin
 * @param {unknown[]} topologies
 * @param {string} type
 * @param {{name: string}[]} types
 * @param {import("angular").IScope & Record<PropertyKey, any>} $scope
 * @param {import("../../../service/utils/FormUtils")} formUtils
 * @param {import("../../../service/utils/TenantUtils")} tenantUtils
 * @param {import("../../../service/utils/DeliveryServiceUtils")} deliveryServiceUtils
 * @param {import("../../../api/DeliveryServiceService")} deliveryServiceService
 * @param {import("../../../api/CDNService")} cdnService
 * @param {import("../../../api/ProfileService")} profileService
 * @param {import("../../../api/TenantService")} tenantService
 * @param {import("../../../models/PropertiesModel")} propertiesModel
 * @param {import("../../../models/UserModel")} userModel
 * @param {import("../../../api/ServerCapabilityService")} serverCapabilityService
 * @param {import("../../../api/ServiceCategoryService")} serviceCategoryService
 */
var FormDeliveryServiceController = function(deliveryService, dsCurrent, origin, topologies, type, types, $scope, formUtils, tenantUtils, deliveryServiceUtils, locationUtils, deliveryServiceService, cdnService, profileService, tenantService, propertiesModel, userModel, serverCapabilityService, serviceCategoryService) {

	/**
	 * This is used to cache TLS version settings when the checkbox is toggled.
	 * @type null | [string, ...string[]]
	 */
	let cachedTLSVersions = null;

	$scope.exposeInactive = !!(propertiesModel.properties.deliveryServices?.exposeInactive);

	$scope.showSensitive = false;

	const knownVersions = new Set(["1.0", "1.1", "1.2", "1.3"]);
	/**
	 * Checks if a TLS version is unknown.
	 * @param {string} v
	 */
	$scope.tlsVersionUnknown = v => v && !knownVersions.has(v);

	const insecureVersions = new Set(["1.0", "1.1"]);
	/**
	 * Checks if a TLS version is known to be insecure.
	 * @param {string} v
	 */
	$scope.tlsVersionInsecure = v => v && insecureVersions.has(v);

	/**
	 * This toggles whether TLS versions are restricted for the Delivery
	 * Service.
	 *
	 * It uses cachedTLSVersions to cache TLS version restrictions, so that the
	 * DS is always ready to submit without manipulation, but the UI "remembers"
	 * the TLS versions that existed on toggling restrictions off.
	 *
	 * This is called when the checkbox's 'change' event fires - that event is
	 * not handled here.
	 */
	function toggleTLSRestrict() {
		if ($scope.restrictTLS) {
			if (cachedTLSVersions instanceof Array && cachedTLSVersions.length > 0) {
				deliveryService.tlsVersions = cachedTLSVersions;
			} else {
				deliveryService.tlsVersions = [""];
			}
			cachedTLSVersions = null;
			return;
		}
		if (deliveryService.tlsVersions instanceof Array && deliveryService.tlsVersions.length > 0) {
			cachedTLSVersions = deliveryService.tlsVersions;
		} else {
			cachedTLSVersions = null;
		}

		deliveryService.tlsVersions = null;
	}
	$scope.toggleTLSRestrict = toggleTLSRestrict;

	$scope.hasGeoLimitCountries = function(ds) {
		return ds !== undefined && (ds.geoLimit === 1 || ds.geoLimit === 2);
	}

	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);

	$scope.loadGeoLimitCountriesRaw = function (ds) {
		if($scope.hasGeoLimitCountries(ds)) {
			ds.geoLimitCountriesRaw = (ds.geoLimitCountries ?? []).join(",");
		} else {
			ds.geoLimitCountriesRaw = "";
		}
	}

	$scope.loadGeoLimitCountries = function (ds) {
		if($scope.hasGeoLimitCountries(ds)) {
			ds.geoLimitCountries = ds.geoLimitCountriesRaw.split(",");
		} else {
			ds.geoLimitCountriesRaw = "";
			ds.geoLimitCountries = [];
		}
	}

	/**
	 * Removes a TLS version at the given index.
	 * @param {number} index
	 */
	$scope.removeTLSVersion = function(index) {
		deliveryService.tlsVersions?.splice(index, 1);
	};

	/**
	 * Adds a TLS version at the given index.
	 * @param {number} index
	 */
	$scope.addTLSVersion = function(index) {
		deliveryService.tlsVersions?.splice(index+1, 0, "");
	};

	/** Compare Arrays
	 *
	 * @template T extends number[] | boolean[] | bigint[] | string[]
	 *
	 * @param {T} a
	 * @param {T} b
	 * @returns `false` if the arrays are equal, `true` otherwise.
	 */
	function arrayCompare (a, b) {
		if (a === b) return false;
		if (a.length !== b.length) return true;

		for (let i = 0; i < a.length; i++) {
			if (a[i] !== b[i]) return true;
		}
		return false;
	};
	$scope.arrayCompare = arrayCompare;

	/**
	 * This function is called when capability is updated on a DSR
	 */
	function capabilityChange() {
		const cap = [];
		for (const [key, value] of Object.entries($scope.selectedCapabilities)) {
			if (value) {
				cap.push(key);
			}
		}
		deliveryService.requiredCapabilities = cap;
	}
	$scope.capabilityChange = capabilityChange;

	/**
	 * This function is called on 'change' events for any and all TLS Version
	 * inputs, and sets validity states of duplicates.
	 *
	 * This can't use a normal validator because it depends on a value checking
	 * against a list containing itself. AngularJS sets values that fail
	 * validation to `undefined`, so if there's a set of TLS versions
	 * `["1.3", "1.3"]`, then the validator will set one of them to `undefined`.
	 * Now the set is `["1.3", undefined]`, so there are no more duplicates, so
	 * the set is marked as valid.
	 */
	function validateTLS() {
		if (!$scope.generalConfig || !($scope.deliveryService.tlsVersions instanceof Array)) {
			return;
		}

		const verMap = new Map();
		for (let i = 0; i < $scope.deliveryService.tlsVersions.length; ++i) {
			const propName = `tlsVersion${i+1}`;
			if (propName in $scope.generalConfig) {
				$scope.generalConfig[propName].$setValidity("duplicates", true);
			}

			const ver = $scope.deliveryService.tlsVersions[i];
			if (ver === undefined) {
				continue;
			}
			const current = verMap.get(ver);
			if (current) {
				current.count++;
				current.indices.push(i);
			} else {
				verMap.set(ver, {
					count: 1,
					indices: [i]
				});
			}
		}

		for (const index of Array.from(verMap).filter(v=>v[1].count>1).flatMap(v=>v[1].indices)) {
			const propName = `tlsVersion${index+1}`;
			if (propName in $scope.generalConfig) {
				$scope.generalConfig[propName].$setValidity("duplicates", false);
			}
		}
	}
	$scope.validateTLS = validateTLS;

	async function getSteeringTargets() {
		if(type.indexOf("HTTP") > -1)  {
			const configs = await deliveryServiceService.getSteering();
			const dsTargets = deliveryServiceUtils.getSteeringTargetsForDS([deliveryService.xmlId], configs);
			$scope.steeringTargetsFor = Array.from(dsTargets[deliveryService.xmlId]);
		}
	}

	/**
	 * Updates the CDNs on the $scope.
	 * @returns {Promise<void>}
	 */
	async function getCDNs() {
		$scope.cdns = await cdnService.getCDNs();
	}

	/**
	 * Updates the Profiles on the $scope.
	 * @returns {Promise<void>}
	 */
	async function getProfiles() {
		/** @type {{type: string}[]} */
		const result = await profileService.getProfiles({ orderby: "name" });
		$scope.profiles = result.filter(p => p.type === "DS_PROFILE");
	}

	/**
	 * Updates the Tenants on the $scope.
	 * @returns {Promise<void>}
	 */
	async function getTenants() {
		const tenants = await tenantService.getTenants();
		const tenant = tenants.find(t => t.id === userModel.user.tenantId);
		$scope.tenants = tenantUtils.hierarchySort(tenantUtils.groupTenantsByParent(tenants), tenant?.parentId, []);
		tenantUtils.addLevels($scope.tenants);
	}

	$scope.selectedCapabilities = {};
	/**
	 * Updates the server Capabilities on the $scope.
	 * @returns {Promise<void>}
	 */
	async function getRequiredCapabilities() {
		$scope.requiredCapabilities = await serverCapabilityService.getServerCapabilities();
		$scope.selectedCapabilities = Object.fromEntries($scope.requiredCapabilities.map(dsc => [dsc.name, $scope.deliveryService.requiredCapabilities.includes(dsc.name)]))
	}

	/**
	 * Updates the Service Categories on the $scope.
	 * @returns {Promise<void>}
	 */
	async function getServiceCategories() {
		$scope.serviceCategories = await serviceCategoryService.getServiceCategories({dsId: deliveryService.id })
	}

	/**
	 * Formats the 'dsCurrent' active flag into a human-readable string. Returns
	 * an empty string if dsCurrent isn't defined.
	 *
	 * @returns {string}
	 */
	function formatCurrentActive() {
		if (!dsCurrent) {
			return "";
		}
		let {active} = dsCurrent;
		if (!propertiesModel.properties.deliveryServices?.exposeInactive && active !== "ACTIVE") {
			active = "INACTIVE";
		}
		return active.split(" ").map(w => w[0].toUpperCase() + w.substring(1).toLowerCase()).join(" ");
	}

	$scope.formatCurrentActive = formatCurrentActive;

	$scope.deliveryService = deliveryService;

	$scope.showGeneralConfig = true;

	$scope.showCacheConfig = true;

	$scope.showRoutingConfig = true;

	$scope.dsCurrent = dsCurrent; // this ds is used primarily for showing the diff between a ds request and the current DS

	$scope.origin = Array.isArray(origin) ? origin[0] : origin;

	$scope.topologies = topologies;

	$scope.showChartsButton = !!(propertiesModel.properties.deliveryServices?.charts?.customLink?.show);

	$scope.openCharts = ds => deliveryServiceUtils.openCharts(ds);

	$scope.dsRequestsEnabled = !!(propertiesModel.properties.dsRequests?.enabled);

	/**
	 * Gods have mercy.
	 *
	 * @param {import("../../../api/DeliveryServiceService").DeliveryService} ds
	 * @returns {string | undefined} An absolutely unsafe direct HTML segment.
	 */
	$scope.edgeFQDNs = function(ds) {
		return ds.exampleURLs?.join("<br/>");
	};

	$scope.DRAFT = 0;
	$scope.SUBMITTED = 1;
	$scope.REJECTED = 2;
	$scope.PENDING = 3;
	$scope.COMPLETE = 4;

	// these may be overriden in a child class. i.e. FormEditDeliveryServiceController
	$scope.saveable = () => true;
	$scope.deletable = () => true;

	$scope.types = types.filter(currentType => {
		let category;
		if (type.includes("ANY_MAP")) {
			category = "ANY_MAP";
		} else if (type.includes("DNS")) {
			category = "DNS";
		} else if (type.includes("HTTP")) {
			category = "HTTP";
		} else if (type.includes("STEERING")) {
			category = 'STEERING';
		} else {
			throw new Error(`unrecognized type: '${type}'`);
		}
		return currentType.name.includes(category);
	});

	$scope.clientSteeringType = types.find(t => t.name === "CLIENT_STEERING");

	/**
	 * Checks if a given Delivery Service uses the "Client Steering" flavor of
	 * Steering-based routing.
	 *
	 * @param {import("../../../api/DeliveryServiceService").DeliveryService} ds The Delivery Service in question.
	 * @returns {boolean} `true` if `ds` uses
	 */
	$scope.isClientSteering = function(ds) {
		if (ds.typeId == $scope.clientSteeringType.id) {
			return true;
		} else {
			ds.trResponseHeaders = "";
			return false;
		}
	};

	$scope.signingAlgos = [
		{ value: null, label: 'None' },
		{ value: 'url_sig', label: 'URL Signature Keys' },
		{ value: 'uri_signing', label: 'URI Signing Keys' }
	];

	$scope.protocols = [
		{ value: 0, label: 'HTTP' },
		{ value: 1, label: 'HTTPS' },
		{ value: 2, label: 'HTTP AND HTTPS' },
		{ value: 3, label: 'HTTP TO HTTPS' }
	];

	$scope.qStrings = [
		{ value: 0, label: 'Use query parameter strings in cache key and pass in upstream requests' },
		{ value: 1, label: 'Do not use query parameter strings in cache key, but do pass in upstream requests' },
		{ value: 2, label: 'Neither use query parameter strings in cache key, nor pass in upstream requests' }
	];

	$scope.geoLimits = [
		{ value: 0, label: 'None' },
		{ value: 1, label: 'Coverage Zone File only' },
		{ value: 2, label: 'Coverage Zone File and Country Code(s)' }
	];

	$scope.geoProviders = [
		{ value: 0, label: 'Maxmind' },
		{ value: 1, label: 'Neustar' }
	];

	$scope.dscps = [
		{ value: 0, label: '0 - Best Effort' },
		{ value: 10, label: '10 - AF11' },
		{ value: 12, label: '12 - AF12' },
		{ value: 14, label: '14 - AF13' },
		{ value: 18, label: '18 - AF21' },
		{ value: 20, label: '20 - AF22' },
		{ value: 22, label: '22 - AF23' },
		{ value: 26, label: '26 - AF31' },
		{ value: 28, label: '28 - AF32' },
		{ value: 30, label: '30 - AF33' },
		{ value: 34, label: '34 - AF41' },
		{ value: 36, label: '36 - AF42' },
		{ value: 37, label: '37 - ' },
		{ value: 38, label: '38 - AF43' },
		{ value: 8, label: '8 - CS1' },
		{ value: 16, label: '16 - CS2' },
		{ value: 24, label: '24 - CS3' },
		{ value: 32, label: '32 - CS4' },
		{ value: 40, label: '40 - CS5' },
		{ value: 48, label: '48 - CS6' },
		{ value: 56, label: '56 - CS7' }
	];

	$scope.rrhs = [
		{ value: 0, label: "Don't cache Range Requests" },
		{ value: 1, label: "Use the background_fetch plugin" },
		{ value: 2, label: "Use the cache_range_requests plugin" },
		{ value: 3, label: "Use the slice plugin" }
	];

	$scope.msoAlgos = [
		{ value: 0, label: "0 - Consistent Hash" },
		{ value: 1, label: "1 - Primary/Backup" },
		{ value: 2, label: "2 - Strict Round Robin" },
		{ value: 3, label: "3 - IP-based Round Robin" },
		{ value: 4, label: "4 - Latch on Failover" }
	];

	/**
	 * Handles changes to the set signing algorithm used by the Delivery Service
	 * by updating the legacy 'signed' property accordingly.
	 *
	 * @param {null|string} signingAlgorithm
	 */
	$scope.changeSigningAlgorithm = function(signingAlgorithm) {
		if (signingAlgorithm === null) {
			deliveryService.signed = false;
		} else {
			deliveryService.signed = true;
		}
	};

	/**
	 * Encodes the given regular expression into $scope.encodedRegex.
	 * @param {string} consistentHashRegex
	 */
	$scope.encodeRegex = function(consistentHashRegex) {
		if (consistentHashRegex !== undefined) {
			$scope.encodedRegex = encodeURIComponent(consistentHashRegex);
		} else {
			$scope.encodedRegex = "";
		}
	};

	/**
	 * Adds a blank consistent hashing query string parameter to the Delivery
	 * Service.
	 */
	$scope.addQueryParam = () => $scope.deliveryService.consistentHashQueryParams.push("");

	/**
	 * Removes a consistent hashing query string parameter from the Delivery
	 * Service at the given index.
	 *
	 * @param {number} index
	 */
	$scope.removeQueryParam = function(index) {
		if ($scope.deliveryService.consistentHashQueryParams.length > 1) {
			$scope.deliveryService.consistentHashQueryParams.splice(index, 1);
		} else {
			// if only one query param is left, don't remove the item from the array. instead, just blank it out
			// so the dynamic form widget will still be visible. empty strings get stripped out on save anyhow.
			$scope.deliveryService.consistentHashQueryParams[index] = "";
		}
		$scope.deliveryServiceForm.$pristine = false; // this enables the 'update' button in the ds form
	};

	$scope.hasError = input => formUtils.hasError(input);

	/**
	 * Checks if a TLS Version has a specific error.
	 *
	 * @param {number} index The index of the TLS Version to check into the
	 * form's Delivery Service's `tlsVersions` array.
	 * @param {string} property The name of the error to check.
	 * @returns {boolean} Whether or not the indicated TLS Version has the given
	 * error.
	 */
	function tlsVersionHasPropertyError(index, property) {
		if (!$scope.generalConfig) {
			return false;
		}
		const propName = `tlsVersion${index+1}`;
		if (!(propName in $scope.generalConfig)) {
			return false;
		}
		return formUtils.hasPropertyError($scope.generalConfig[propName], property);
	}
	$scope.tlsVersionHasPropertyError = tlsVersionHasPropertyError;

	this.$onInit = function() {
		$scope.loadGeoLimitCountriesRaw(deliveryService);
		$scope.loadGeoLimitCountriesRaw(dsCurrent);
	}

	/**
	 * Checks if a TLS Version has any error.
	 *
	 * @param {number} index The index of the TLS Version to check into the
	 * form's Delivery Service's `tlsVersions` array.
	 * @returns {boolean} Whether or not the indicated TLS Version has an error.
	 */
	 function tlsVersionHasError(index) {
		if (!$scope.generalConfig) {
			return false;
		}
		const propName = `tlsVersion${index+1}`;
		if (!(propName in $scope.generalConfig)) {
			return false;
		}
		return formUtils.hasError($scope.generalConfig[propName]);
	}
	$scope.tlsVersionHasError = tlsVersionHasError;

	$scope.hasPropertyError = (input, property) => formUtils.hasPropertyError(input, property);

	$scope.rangeRequestSelected = function() {
		if ($scope.deliveryService.rangeRequestHandling != 3) {
			$scope.deliveryService.rangeSliceBlockSize = null;
		}
	};

	getCDNs();
	getProfiles();
	getTenants();
	getRequiredCapabilities();
	getServiceCategories();
	getSteeringTargets();
	if (!deliveryService.consistentHashQueryParams || deliveryService.consistentHashQueryParams.length < 1) {
		// add an empty one so the dynamic form widget is visible. empty strings get stripped out on save anyhow.
		$scope.deliveryService.consistentHashQueryParams = [ "" ];
	}
	if (deliveryService.lastUpdated) {
		// TS checkers hate him for this one weird trick:
		// @ts-ignore
		deliveryService.lastUpdated = new Date(deliveryService.lastUpdated.replace("+00", "Z"));
		// ... the right way to do this is with an interceptor, but nobody
		// wants to put in that kinda work on a legacy product.
	}

	if (!$scope.exposeInactive && deliveryService.active === "INACTIVE") {
		deliveryService.active = "PRIMED";
	}
};

FormDeliveryServiceController.$inject = ["deliveryService", "dsCurrent", "origin", "topologies", "type", "types", "$scope", "formUtils", "tenantUtils", "deliveryServiceUtils", "locationUtils", "deliveryServiceService", "cdnService", "profileService", "tenantService", "propertiesModel", "userModel", "serverCapabilityService", "serviceCategoryService"];
module.exports = FormDeliveryServiceController;
