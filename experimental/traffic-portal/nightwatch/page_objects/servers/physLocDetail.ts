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

import { EnhancedPageObject } from "nightwatch";

/**
 * Defines the PageObject for PhysLoc Details.
 */
export type PhysLocDetailPageObject = EnhancedPageObject<{}, typeof physLocDetailPageObject.elements>;

const physLocDetailPageObject = {
	elements: {
		address: {
			selector: "input[name='address']"
		},
		city: {
			selector: "input[name='city']"
		},
		comments: {
			selector: "textarea[name='comments']"
		},
		email: {
			selector: "input[name='email']"
		},
		id: {
			selector: "input[name='id']"
		},
		lastUpdated: {
			selector: "input[name='lastUpdated']"
		},
		name: {
			selector: "input[name='name']"
		},
		phone: {
			selector: "input[name='phone']"
		},
		poc: {
			selector: "input[name='poc']"
		},
		region: {
			selector: "mat-select[name='region']"
		},
		saveBtn: {
			selector: "button[type='submit']"
		},
		shortName: {
			selector: "input[name='shortName']"
		},
		state: {
			selector: "input[name='state']"
		},
		zip: {
			selector: "input[name='postalCode']"
		}
	},
};

export default physLocDetailPageObject;
