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

/**
 * A single point of data compatible with chart.js's configuration and API
 */
export interface DataPoint {

	/**
	 * Defines a time position - {@link DataPoint.x} should not exist on @{link Datapoint}s that define 't'
	 */
	t?: Date;
	x?: number;
	y: number;
}

/**
 * A chart.js-compatible data set that includes an array of {@link DataPoint} data as well as chart.js configuration options
 */
export interface DataSet {
	label: string;
	data: Array<DataPoint>;
	backgroundColor?: string | Array<string>;
	borderColor?: string | Array<string>;
	borderDash?: number[];
	borderWidth?: number;
	fill?: boolean;
	fillColor?: string;
}

/**
 * Encapsulates a {@link DataSet} with aggregate information as returned by the Traffic Ops API
 */
export interface DataSetWithSummary {
	dataSet: DataSet;
	min: number;
	max: number;
	fifthPercentile: number;
	ninetyFifthPercentile: number;
	ninetyEighthPercentile: number;
	mean: number;
}

/**
 * Contains all possible "TPS" data that can be returned by the Traffic Ops API
 */
export interface TPSData {
	total: DataSetWithSummary;
	informational?: DataSetWithSummary;
	success: DataSetWithSummary;
	redirection: DataSetWithSummary;
	clientError: DataSetWithSummary;
	serverError: DataSetWithSummary;
}
