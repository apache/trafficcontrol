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

import { type AfterViewInit, Directive, ElementRef, Input, type OnDestroy } from "@angular/core";
import { Chart } from "chart.js"; // TODO: use plotly instead for WebGL-capabale browsers?
import { from, type Observable, type Subscription } from "rxjs";

import type { DataSet } from "src/app/models/data";

import { LoggingService } from "../logging.service";

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
	@Input() public chartDataSets: Observable<Array<DataSet | null> | null> = from([]);
	/** The type of the chart. */
	@Input() public chartType?: "category" | "linear" | "logarithmic" | "time";
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

	constructor(private readonly element: ElementRef, private readonly log: LoggingService) { }

	/**
	 * Initializes the chart using the input data.
	 */
	public ngAfterViewInit(): void {
		if (!(this.element.nativeElement instanceof HTMLCanvasElement)) {
			throw new Error("[linechart] Directive can only be used on a canvas in a context where DOM access is allowed");
		}

		const ctx = this.element.nativeElement.getContext("2d", {alpha: false});
		if (!ctx) {
			throw new Error("Failed to get 2D context for chart canvas");
		}
		this.ctx = ctx;

		if (!this.chartType) {
			this.chartType = "linear";
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
			data => {
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
	private dataLoad(data: Array<DataSet | null> | null): void {
		this.destroyChart();

		const hasNoNulls = (arr: Array<DataSet | null>): arr is Array<DataSet> => !arr.some(x=>x===null);

		if (data === null || data === undefined || !hasNoNulls(data)) {
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
		this.log.error("data error occurred:", e);
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
