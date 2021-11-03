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
import { AfterViewInit, Directive, ElementRef, Input, OnDestroy } from "@angular/core";

import { from, Observable, Subscription } from "rxjs";

import { Chart } from "chart.js"; // TODO: use plotly instead for WebGL-capabale browsers?

import { DataSet } from "../../models/data";

/**
 * LineChartType enumerates the valid types of charts.
 */
export enum LineChartType {
	/**
	 * Plots category proportions.
	 */
	CATEGORY = "category",
	/**
	 * Scatter plots.
	 */
	LINEAR = "linear",
	/**
	 * Logarithmic-scale scatter plots.
	 */
	LOGARITHMIC = "logarithmic",
	/**
	 * Time-series data.
	 */
	TIME = "time"
}

/**
 * LinechartDirective decorates canvases by creating a rendering context for
 * ChartJS charts.
 */
@Directive({
	selector: "[linechart]",
})
export class LinechartDirective implements AfterViewInit, OnDestroy {

	/** The chart context. */
	private ctx: CanvasRenderingContext2D | null = null; // | WebGLRenderingContext;
	/** The Chart.js API object. */
	private chart: Chart | null = null;

	/** The title of the chart. */
	@Input() public chartTitle?: string;
	/** Labels for the datasets. */
	@Input() public chartLabels?: unknown[];
	/** Data to be plotted by the chart. */
	@Input() public chartDataSets: Observable<DataSet[]> = from([]);
	/** The type of the chart. */
	@Input() public chartType?: LineChartType;
	/** A label for the X-axis of the chart. */
	@Input() public chartXAxisLabel?: string;
	/** A label for the Y-axis of the chart. */
	@Input() public chartYAxisLabel?: string;
	/** A callback for the label of each data point, to be optionally provided. */
	@Input() public chartLabelCallback?: (v: unknown, i: number, va: unknown[]) => unknown;
	/** Whether or not to display the chart's legend. */
	@Input() public chartDisplayLegend?: boolean;

	/** A subscription for the chartDataSets input. */
	private subscription: Subscription | null = null;
	/** Chart.js configuration options. */
	private opts: Chart.ChartConfiguration = {};

	/**
	 * Constructor.
	 */
	constructor(private readonly element: ElementRef) { }

	/**
	 * Initializes the chart using the input data.
	 */
	public ngAfterViewInit(): void {
		if (this.element.nativeElement === null) {
			console.warn("Use of DOM directive in non-DOM context!");
			return;
		}

		if (!(this.element.nativeElement instanceof HTMLCanvasElement)) {
			throw new Error("[linechart] Directive can only be used on a canvas!");
		}

		const ctx = this.element.nativeElement.getContext("2d", {alpha: false});
		if (!ctx) {
			throw new Error("Failed to get 2D context for chart canvas");
		}
		this.ctx = ctx;

		if (!this.chartType) {
			this.chartType = LineChartType.LINEAR;
		}

		if (this.chartDisplayLegend === null || this.chartDisplayLegend === undefined) {
			this.chartDisplayLegend = false;
		}

		this.opts = {
			data: {
				datasets: [],
			},
			options: {
				legend: {
					display: true
				},
				scales: {
					xAxes: [{
						display: true,
						scaleLabel: {
							display: this.chartXAxisLabel ? true : false,
							labelString: this.chartXAxisLabel
						},
						type: this.chartType,
					}],
					yAxes: [{
						display: true,
						scaleLabel: {
							display: this.chartYAxisLabel ? true : false,
							labelString: this.chartYAxisLabel
						},
						ticks: {
							suggestedMin: 0
						}
					}]
				},
				title: {
					display: this.chartTitle ? true : false,
					text: this.chartTitle
				},
				tooltips: {
					intersect: false,
					mode: "x"
				}
			},
			type: "line"
		};

		this.subscription = this.chartDataSets.subscribe(
			(data: DataSet[]) => {
				this.dataLoad(data);
			},
			(e: Error) => {
				this.dataError(e);
			}
		);
	}

	/**
	 * Destroys the ChartJS instance and clears the underlying drawing context.
	 */
	private destroyChart(): void {
		if (this.chart) {
			this.chart.clear();
			(this.chart.destroy as () => void)();
			this.chart = null;
			if (this.ctx) {
				this.ctx.clearRect(0, 0, this.element.nativeElement.width, this.element.nativeElement.height);
			}
			this.opts.data = {datasets: [], labels: []};
		}
	}


	/**
	 * Loads a new Chart.
	 *
	 * @param data The new data sets for the new chart.
	 */
	private dataLoad(data: DataSet[]): void {
		this.destroyChart();

		if (data === null || data === undefined || data.some(x => x === null)) {
			this.noData();
			return;
		}

		if (this.opts.data) {
			this.opts.data.datasets = data;
		} else {
			this.opts.data = {datasets: data};
		}

		if (!this.ctx) {
			throw new Error("cannot load data with uninitialized context");
		}

		this.chart = new Chart(this.ctx, this.opts);
	}

	/**
	 * Handles an error when loading data.
	 *
	 * @param e The error that occurred.
	 */
	private dataError(e: Error): void {
		console.error("data error occurred:", e);
		this.destroyChart();
		if (this.ctx) {
			this.ctx.font = "30px serif";
			this.ctx.fillStyle = "black";
			this.ctx.textAlign = "center";
			this.ctx.fillText("Error Fetching Data", this.element.nativeElement.width / 2, this.element.nativeElement.height / 2);
		}
	}

	/**
	 * Handles when there is no data for the chart.
	 */
	private noData(): void {
		this.destroyChart();
		if (this.ctx) {
			this.ctx.font = "30px serif";
			this.ctx.fillStyle = "black";
			this.ctx.textAlign = "center";
			this.ctx.fillText("No Data", this.element.nativeElement.width / 2, this.element.nativeElement.height / 2);
		}
	}

	/** Cleans up chart resources on element destruction. */
	public ngOnDestroy(): void {
		this.destroyChart();

		if (this.subscription) {
			this.subscription.unsubscribe();
		}
	}

}
