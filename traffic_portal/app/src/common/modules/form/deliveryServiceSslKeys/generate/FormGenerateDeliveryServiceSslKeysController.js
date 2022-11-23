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
 * @param {*} deliveryService
 * @param {*} sslKeys
 * @param {*} sslRequest
 * @param {*} $scope
 * @param {import("../../../../service/utils/angular.ui.bootstrap").IModalService} $uibModal
 * @param {import("../../../../service/utils/LocationUtils")} locationUtils
 * @param {import("../../../../api/DeliveryServiceSslKeysService")} deliveryServiceSslKeysService
 * @param {import("../../../../service/utils/FormUtils")} formUtils
 */
var FormGenerateDeliveryServiceSslKeysController = function(deliveryService, sslKeys, sslRequest, $scope, $uibModal, locationUtils, deliveryServiceSslKeysService, formUtils) {

	var setSSLRequest = function(sslRequest) {
		if (!sslRequest.hostname) {
			console.log('setting default hostname');
			var url = deliveryService.exampleURLs[0],
				defaultHostName = url.split("://")[1];
			if (deliveryService.type.indexOf('HTTP') != -1) {
				var parts = defaultHostName.split(".");
				parts[0] = "*";
				defaultHostName = parts.join(".");
			}
			sslRequest.hostname = defaultHostName;
		}
		return sslRequest;
	};

	var getAcmeProviders = function() {
		deliveryServiceSslKeysService.getAcmeProviders()
			.then(function(result) {
				$scope.acmeProviders = result;
				if (!$scope.acmeProviders.includes('Lets Encrypt')) {
					$scope.acmeProviders.push('Lets Encrypt');
				}
			});
	};

	$scope.loadAcmeProviders = function() {
		if ($scope.useAcme) {
			getAcmeProviders();
		}
	};

	$scope.useAcme = false;
	$scope.acmeProviders = [];
	$scope.acmeProvider = "";
	$scope.hasError = formUtils.hasError;
	$scope.hasPropertyError = formUtils.hasPropertyError;
	$scope.navigateToPath = (path, unsavedChanges) => locationUtils.navigateToPath(path, unsavedChanges);
	$scope.sslRequest = setSSLRequest(sslRequest);

	$scope.hasAcmeProviderError = function() {
		return $scope.acmeProvider === null || $scope.acmeProvider === '';
	};

	$scope.deliveryService = deliveryService;
	$scope.countries = [
		{code:"US", name:"United States (US)"},
		{code:"AD", name:"Andorra (AD)"},
		{code:"AE", name:"United Arab Emirates (AE)"},
		{code:"AF", name:"Afghanistan (AF)"},
		{code:"AG", name:"Antigua and Barbuda (AG)"},
		{code:"AI", name:"Anguilla (AI)"},
		{code:"AL", name:"Albania (AL)"},
		{code:"AM", name:"Armenia (AM)"},
		{code:"AO", name:"Angola (AO)"},
		{code:"AP", name:"Asia/Pacific Region (AP)"},
		{code:"AQ", name:"Antarctica (AQ)"},
		{code:"AR", name:"Argentina (AR)"},
		{code:"AS", name:"American Samoa (AS)"},
		{code:"AT", name:"Austria (AT)"},
		{code:"AU", name:"Australia (AU)"},
		{code:"AW", name:"Aruba (AW)"},
		{code:"AX", name:"Aland Islands (AX)"},
		{code:"AZ", name:"Azerbaijan (AZ)"},
		{code:"BA", name:"Bosnia and Herzegovina (BA)"},
		{code:"BB", name:"Barbados (BB)"},
		{code:"BD", name:"Bangladesh (BD)"},
		{code:"BE", name:"Belgium (BE)"},
		{code:"BF", name:"Burkina Faso (BF)"},
		{code:"BG", name:"Bulgaria (BG)"},
		{code:"BH", name:"Bahrain (BH)"},
		{code:"BI", name:"Burundi (BI)"},
		{code:"BJ", name:"Benin (BJ)"},
		{code:"BL", name:"Saint Bartelemey (BL)"},
		{code:"BM", name:"Bermuda (BM)"},
		{code:"BN", name:"Brunei Darussalam (BN)"},
		{code:"BO", name:"Bolivia (BO)"},
		{code:"BQ", name:"Bonaire, Saint Eustatius and Saba (BQ)"},
		{code:"BR", name:"Brazil (BR)"},
		{code:"BS", name:"Bahamas (BS)"},
		{code:"BT", name:"Bhutan (BT)"},
		{code:"BV", name:"Bouvet Island (BV)"},
		{code:"BW", name:"Botswana (BW)"},
		{code:"BY", name:"Belarus (BY)"},
		{code:"BZ", name:"Belize (BZ)"},
		{code:"CA", name:"Canada (CA)"},
		{code:"CC", name:"Cocos (Keeling) Islands (CC)"},
		{code:"CD", name:"Congo, The Democratic Republic of the (CD)"},
		{code:"CF", name:"Central African Republic (CF)"},
		{code:"CG", name:"Congo (CG)"},
		{code:"CH", name:"Switzerland (CH)"},
		{code:"CI", name:"Cote d'Ivoire (CI)"},
		{code:"CK", name:"Cook Islands (CK)"},
		{code:"CL", name:"Chile (CL)"},
		{code:"CM", name:"Cameroon (CM)"},
		{code:"CN", name:"China (CN)"},
		{code:"CO", name:"Colombia (CO)"},
		{code:"CR", name:"Costa Rica (CR)"},
		{code:"CU", name:"Cuba (CU)"},
		{code:"CV", name:"Cape Verde (CV)"},
		{code:"CW", name:"Curacao (CW)"},
		{code:"CX", name:"Christmas Island (CX)"},
		{code:"CY", name:"Cyprus (CY)"},
		{code:"CZ", name:"Czech Republic (CZ)"},
		{code:"DE", name:"Germany (DE)"},
		{code:"DJ", name:"Djibouti (DJ)"},
		{code:"DK", name:"Denmark (DK)"},
		{code:"DM", name:"Dominica (DM)"},
		{code:"DO", name:"Dominican Republic (DO)"},
		{code:"DZ", name:"Algeria (DZ)"},
		{code:"EC", name:"Ecuador (EC)"},
		{code:"EE", name:"Estonia (EE)"},
		{code:"EG", name:"Egypt (EG)"},
		{code:"EH", name:"Western Sahara (EH)"},
		{code:"ER", name:"Eritrea (ER)"},
		{code:"ES", name:"Spain (ES)"},
		{code:"ET", name:"Ethiopia (ET)"},
		{code:"EU", name:"Europe (EU)"},
		{code:"FI", name:"Finland (FI)"},
		{code:"FJ", name:"Fiji (FJ)"},
		{code:"FK", name:"Falkland Islands (Malvinas) (FK)"},
		{code:"FM", name:"Micronesia, Federated States of (FM)"},
		{code:"FO", name:"Faroe Islands (FO)"},
		{code:"FR", name:"France (FR)"},
		{code:"GA", name:"Gabon (GA)"},
		{code:"GB", name:"United Kingdom (GB)"},
		{code:"GD", name:"Grenada (GD)"},
		{code:"GE", name:"Georgia (GE)"},
		{code:"GF", name:"French Guiana (GF)"},
		{code:"GG", name:"Guernsey (GG)"},
		{code:"GH", name:"Ghana (GH)"},
		{code:"GI", name:"Gibraltar (GI)"},
		{code:"GL", name:"Greenland (GL)"},
		{code:"GM", name:"Gambia (GM)"},
		{code:"GN", name:"Guinea (GN)"},
		{code:"GP", name:"Guadeloupe (GP)"},
		{code:"GQ", name:"Equatorial Guinea (GQ)"},
		{code:"GR", name:"Greece (GR)"},
		{code:"GS", name:"South Georgia and the South Sandwich Islands (GS)"},
		{code:"GT", name:"Guatemala (GT)"},
		{code:"GU", name:"Guam (GU)"},
		{code:"GW", name:"Guinea-Bissau (GW)"},
		{code:"GY", name:"Guyana (GY)"},
		{code:"HK", name:"Hong Kong (HK)"},
		{code:"HM", name:"Heard Island and McDonald Islands (HM)"},
		{code:"HN", name:"Honduras (HN)"},
		{code:"HR", name:"Croatia (HR)"},
		{code:"HT", name:"Haiti (HT)"},
		{code:"HU", name:"Hungary (HU)"},
		{code:"ID", name:"Indonesia (ID)"},
		{code:"IE", name:"Ireland (IE)"},
		{code:"IL", name:"Israel (IL)"},
		{code:"IM", name:"Isle of Man (IM)"},
		{code:"IN", name:"India (IN)"},
		{code:"IO", name:"British Indian Ocean Territory (IO)"},
		{code:"IQ", name:"Iraq (IQ)"},
		{code:"IR", name:"Iran, Islamic Republic of (IR)"},
		{code:"IS", name:"Iceland (IS)"},
		{code:"IT", name:"Italy (IT)"},
		{code:"JE", name:"Jersey (JE)"},
		{code:"JM", name:"Jamaica (JM)"},
		{code:"JO", name:"Jordan (JO)"},
		{code:"JP", name:"Japan (JP)"},
		{code:"KE", name:"Kenya (KE)"},
		{code:"KG", name:"Kyrgyzstan (KG)"},
		{code:"KH", name:"Cambodia (KH)"},
		{code:"KI", name:"Kiribati (KI)"},
		{code:"KM", name:"Comoros (KM)"},
		{code:"KN", name:"Saint Kitts and Nevis (KN)"},
		{code:"KP", name:"Korea, Democratic People's Republic of (KP)"},
		{code:"KR", name:"Korea, Republic of (KR)"},
		{code:"KW", name:"Kuwait (KW)"},
		{code:"KY", name:"Cayman Islands (KY)"},
		{code:"KZ", name:"Kazakhstan (KZ)"},
		{code:"LA", name:"Lao People's Democratic Republic (LA)"},
		{code:"LB", name:"Lebanon (LB)"},
		{code:"LC", name:"Saint Lucia (LC)"},
		{code:"LI", name:"Liechtenstein (LI)"},
		{code:"LK", name:"Sri Lanka (LK)"},
		{code:"LR", name:"Liberia (LR)"},
		{code:"LS", name:"Lesotho (LS)"},
		{code:"LT", name:"Lithuania (LT)"},
		{code:"LU", name:"Luxembourg (LU)"},
		{code:"LV", name:"Latvia (LV)"},
		{code:"LY", name:"Libyan Arab Jamahiriya (LY)"},
		{code:"MA", name:"Morocco (MA)"},
		{code:"MC", name:"Monaco (MC)"},
		{code:"MD", name:"Moldova, Republic of (MD)"},
		{code:"ME", name:"Montenegro (ME)"},
		{code:"MF", name:"Saint Martin (MF)"},
		{code:"MG", name:"Madagascar (MG)"},
		{code:"MH", name:"Marshall Islands (MH)"},
		{code:"MK", name:"Macedonia (MK)"},
		{code:"ML", name:"Mali (ML)"},
		{code:"MM", name:"Myanmar (MM)"},
		{code:"MN", name:"Mongolia (MN)"},
		{code:"MO", name:"Macao (MO)"},
		{code:"MP", name:"Northern Mariana Islands (MP)"},
		{code:"MQ", name:"Martinique (MQ)"},
		{code:"MR", name:"Mauritania (MR)"},
		{code:"MS", name:"Montserrat (MS)"},
		{code:"MT", name:"Malta (MT)"},
		{code:"MU", name:"Mauritius (MU)"},
		{code:"MV", name:"Maldives (MV)"},
		{code:"MW", name:"Malawi (MW)"},
		{code:"MX", name:"Mexico (MX)"},
		{code:"MY", name:"Malaysia (MY)"},
		{code:"MZ", name:"Mozambique (MZ)"},
		{code:"NA", name:"Namibia (NA)"},
		{code:"NC", name:"New Caledonia (NC)"},
		{code:"NE", name:"Niger (NE)"},
		{code:"NF", name:"Norfolk Island (NF)"},
		{code:"NG", name:"Nigeria (NG)"},
		{code:"NI", name:"Nicaragua (NI)"},
		{code:"NL", name:"Netherlands (NL)"},
		{code:"NO", name:"Norway (NO)"},
		{code:"NP", name:"Nepal (NP)"},
		{code:"NR", name:"Nauru (NR)"},
		{code:"NU", name:"Niue (NU)"},
		{code:"NZ", name:"New Zealand (NZ)"},
		{code:"OM", name:"Oman (OM)"},
		{code:"PA", name:"Panama (PA)"},
		{code:"PE", name:"Peru (PE)"},
		{code:"PF", name:"French Polynesia (PF)"},
		{code:"PG", name:"Papua New Guinea (PG)"},
		{code:"PH", name:"Philippines (PH)"},
		{code:"PK", name:"Pakistan (PK)"},
		{code:"PL", name:"Poland (PL)"},
		{code:"PM", name:"Saint Pierre and Miquelon (PM)"},
		{code:"PN", name:"Pitcairn (PN)"},
		{code:"PR", name:"Puerto Rico (PR)"},
		{code:"PS", name:"Palestinian Territory (PS)"},
		{code:"PT", name:"Portugal (PT)"},
		{code:"PW", name:"Palau (PW)"},
		{code:"PY", name:"Paraguay (PY)"},
		{code:"QA", name:"Qatar (QA)"},
		{code:"RE", name:"Reunion (RE)"},
		{code:"RO", name:"Romania (RO)"},
		{code:"RS", name:"Serbia (RS)"},
		{code:"RU", name:"Russian Federation (RU)"},
		{code:"RW", name:"Rwanda (RW)"},
		{code:"SA", name:"Saudi Arabia (SA)"},
		{code:"SB", name:"Solomon Islands (SB)"},
		{code:"SC", name:"Seychelles (SC)"},
		{code:"SD", name:"Sudan (SD)"},
		{code:"SE", name:"Sweden (SE)"},
		{code:"SG", name:"Singapore (SG)"},
		{code:"SH", name:"Saint Helena (SH)"},
		{code:"SI", name:"Slovenia (SI)"},
		{code:"SJ", name:"Svalbard and Jan Mayen (SJ)"},
		{code:"SK", name:"Slovakia (SK)"},
		{code:"SL", name:"Sierra Leone (SL)"},
		{code:"SM", name:"San Marino (SM)"},
		{code:"SN", name:"Senegal (SN)"},
		{code:"SO", name:"Somalia (SO)"},
		{code:"SR", name:"Suriname (SR)"},
		{code:"SS", name:"South Sudan (SS)"},
		{code:"ST", name:"Sao Tome and Principe (ST)"},
		{code:"SV", name:"El Salvador (SV)"},
		{code:"SX", name:"Sint Maarten (SX)"},
		{code:"SY", name:"Syrian Arab Republic (SY)"},
		{code:"SZ", name:"Swaziland (SZ)"},
		{code:"TC", name:"Turks and Caicos Islands (TC)"},
		{code:"TD", name:"Chad (TD)"},
		{code:"TF", name:"French Southern Territories (TF)"},
		{code:"TG", name:"Togo (TG)"},
		{code:"TH", name:"Thailand (TH)"},
		{code:"TJ", name:"Tajikistan (TJ)"},
		{code:"TK", name:"Tokelau (TK)"},
		{code:"TL", name:"Timor-Leste (TL)"},
		{code:"TM", name:"Turkmenistan (TM)"},
		{code:"TN", name:"Tunisia (TN)"},
		{code:"TO", name:"Tonga (TO)"},
		{code:"TR", name:"Turkey (TR)"},
		{code:"TT", name:"Trinidad and Tobago (TT)"},
		{code:"TV", name:"Tuvalu (TV)"},
		{code:"TW", name:"Taiwan (TW)"},
		{code:"TZ", name:"Tanzania, United Republic of (TZ)"},
		{code:"UA", name:"Ukraine (UA)"},
		{code:"UG", name:"Uganda (UG)"},
		{code:"UM", name:"United States Minor Outlying Islands (UM)"},
		{code:"UY", name:"Uruguay (UY)"},
		{code:"UZ", name:"Uzbekistan (UZ)"},
		{code:"VA", name:"Holy See (Vatican City State) (VA)"},
		{code:"VC", name:"Saint Vincent and the Grenadines (VC)"},
		{code:"VE", name:"Venezuela (VE)"},
		{code:"VG", name:"Virgin Islands, British (VG)"},
		{code:"VI", name:"Virgin Islands, U.S. (VI)"},
		{code:"VN", name:"Vietnam (VN)"},
		{code:"VU", name:"Vanuatu (VU)"},
		{code:"WF", name:"Wallis and Futuna (WF)"},
		{code:"WS", name:"Samoa (WS)"},
		{code:"YE", name:"Yemen (YE)"},
		{code:"YT", name:"Mayotte (YT)"},
		{code:"ZA", name:"South Africa (ZA)"},
		{code:"ZM", name:"Zambia (ZM)"},
		{code:"ZW", name: "Zimbabwe (ZW)"}
	];

	$scope.confirmGenerate = function(sslRequest) {
        var params = {
            title: 'Generate New SSL Keys for Delivery Service: ' + deliveryService.xmlId,
            message: ' (replacing any previous keys)'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            deliveryServiceSslKeysService.generateSslKeys(deliveryService, sslKeys, sslRequest).then(
                function() {
                    locationUtils.navigateToPath('/delivery-services/' + deliveryService.id + '/ssl-keys');
                });
        });
    };

    $scope.confirmGenerateAcme = function(sslRequest) {
        var params = {
            title: 'Generate New SSL Keys Using Let\'s Encrypt for Delivery Service: ' + deliveryService.xmlId,
            message: ' (replacing any previous keys)'
        };
        var modalInstance = $uibModal.open({
            templateUrl: 'common/modules/dialog/confirm/dialog.confirm.tpl.html',
            controller: 'DialogConfirmController',
            size: 'md',
            resolve: {
                params: function () {
                    return params;
                }
            }
        });
        modalInstance.result.then(function() {
            sslKeys.authType = $scope.acmeProvider;
            deliveryServiceSslKeysService.generateSslKeysWithAcme(deliveryService, sslKeys, sslRequest).then(
                function() {
                    locationUtils.navigateToPath('/delivery-services/' + deliveryService.id + '/ssl-keys');
                });
        });
    };

};

FormGenerateDeliveryServiceSslKeysController.$inject = ['deliveryService', 'sslKeys', 'sslRequest', '$scope', '$uibModal', 'locationUtils', 'deliveryServiceSslKeysService', 'formUtils'];
module.exports = FormGenerateDeliveryServiceSslKeysController;
