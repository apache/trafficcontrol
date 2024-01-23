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
		address: "input[name='address']",
		city: "input[name='city']",
		comments: "textarea[name='comments']",
		email: "input[name='email']",
		id: "input[name='id']",
		lastUpdated: "input[name='lastUpdated']",
		name: "input[name='name']",
		phone: "input[name='phone']",
		poc: "input[name='poc']",
		region: "mat-select[name='region']",
		saveBtn: "button[type='submit']",
		shortName: "input[name='shortName']",
		state: "input[name='state']",
		zip: "input[name='postalCode']",
	},
};

export default physLocDetailPageObject;
