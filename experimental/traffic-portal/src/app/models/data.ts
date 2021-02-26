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
export type DataPoint = {
	/**
	 * The y-axis value of the data point.
	 *
	 * Sometimes it's been set to a fixed precision.
	 */
	y: number | string;
} & ({
	/**
	 * Defines a time position for time-series data points.
	 */
	t: Date;
} | {
	/**
	 * The x-axis value for non-time-series data points.
	 */
	x: number;
});

/**
 * A chart.js-compatible data set that includes an array of {@link DataPoint}
 * data as well as chart.js configuration options.
 */
export interface DataSet {
	/** A label describing the data set. */
	label: string;
	/**
	 * The actual data to be plotted.
	 */
	data: Array<DataPoint>;
	/**
	 * See ChartJS documentation for information on how to set the background
	 * color.
	 */
	backgroundColor?: string | Array<string>;
	/**
	 * See ChartJS documentation for information on how to set the border color.
	 */
	borderColor?: string | Array<string>;
	/**
	 * See ChartJS documentation for how to set dash lengths for borders.
	 */
	borderDash?: number[];
	/**
	 * Sets the border width, in px.
	 */
	borderWidth?: number;
	/**
	 * Whether the plotted data points/areas should be filled (true) or merely
	 * outlined (false).
	 */
	fill?: boolean;
	/**
	 * CSS color with which the filling of plotted data points/areas will be
	 * drawn.
	 */
	fillColor?: string;
}

/**
 * Encapsulates a {@link DataSet} with aggregate information as returned by the Traffic Ops API
 */
export interface DataSetWithSummary {
	/** The data set being plotted. */
	dataSet: DataSet;
	/** The minimum value of the data set. */
	min: number;
	/** The maximum value of the data set. */
	max: number;
	/** The number that denotes the fifth percentile of the data set. */
	fifthPercentile: number;
	/** The number that denotes the 95th percentile of the data set. */
	ninetyFifthPercentile: number;
	/** The number that denotes the 98th percentile of the data set. */
	ninetyEighthPercentile: number;
	/** The arithmetic mean - or "average" - of the data set values. */
	mean: number;
}

/**
 * Contains all possible "TPS" data that can be returned by the Traffic Ops API.
 */
export interface TPSData {
	/** A data set describing the total transactions per second. */
	total: DataSetWithSummary;
	/**
	 * A data set describing transactions per second that were served with HTTP
	 * response codes on the interval [100,200).
	 */
	informational?: DataSetWithSummary;
	/**
	 * A data set describing transactions per second that were served with HTTP
	 * response codes on the interval [200,300).
	 */
	success: DataSetWithSummary;
	/**
	 * A data set describing transactions per second that were served with HTTP
	 * response codes on the interval [300,400).
	 */
	redirection: DataSetWithSummary;
	/**
	 * A data set describing transactions per second that were served with HTTP
	 * response codes on the interval [400,500).
	 */
	clientError: DataSetWithSummary;
	/**
	 * A data set describing transactions per second that were served with HTTP
	 * response codes on the interval [500,600).
	 */
	serverError: DataSetWithSummary;
}
