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
export class DataPoint {
	t?: Date;
	x?: number;
	y: number;
}

export class DataSet {
	label: string;
	data: Array<DataPoint>;
	backgroundColor?: string | Array<string>;
	borderColor?: string | Array<string>;
	borderDash?: number[];
	borderWidth?: number;
	fill?: boolean;
	fillColor?: string;
}

export class DataSetWithSummary {
	dataSet: DataSet;
	min: number;
	max: number;
	fifthPercentile: number;
	ninetyFifthPercentile: number;
	ninetyEighthPercentile: number;
	mean: number;

	public get average(): number {
		return this.mean;
	}
	public set average(a: number) {
		this.mean = a;
	}
}

export class TPSData {
	total: DataSetWithSummary;
	informational?: DataSetWithSummary;
	success: DataSetWithSummary;
	redirection: DataSetWithSummary;
	clientError: DataSetWithSummary;
	serverError: DataSetWithSummary;
}
