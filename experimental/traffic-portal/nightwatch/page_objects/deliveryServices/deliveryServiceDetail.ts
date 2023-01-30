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
import {
	EnhancedPageObject, EnhancedSectionInstance
} from "nightwatch";

/**
 * Define the type for our PO
 */
export type DeliveryServiceDetailPageObject = EnhancedPageObject<{}, typeof deliveryServiceDetailPageObject.elements,
{ dateInputForm: EnhancedSectionInstance<{}, typeof deliveryServiceDetailPageObject.sections.dateInputForm.elements> }>;

const deliveryServiceDetailPageObject = {
	elements: {
		bandwidthChart: {
			selector: "canvas#bandwidthData"
		},
		invalidateJobs: {
			selector: "a#invalidate"
		},
		tpsChart: {
			selector: "canvas#tpsChartData"
		},
	},
	sections: {
		dateInputForm: {
			elements: {
				fromDate: {
					selector: "input[name='fromdate']"
				},
				fromTime: {
					selector: "input[name='fromtime']"
				},
				refreshBtn: {
					selector: "button[name='timespanRefresh']"
				},
				steeringIcon: {
					selector: "div.actions > mat-icon"
				},
				toDate: {
					selector: "input[name='todate']"
				},
				toTime: {
					selector: "input[name='totime']"
				}
			},
			selector: "form[name='timespan']"
		}
	}
};

export default deliveryServiceDetailPageObject;
