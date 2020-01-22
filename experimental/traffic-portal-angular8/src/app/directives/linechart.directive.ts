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
import { AfterViewInit, Directive, ElementRef, OnDestroy } from '@angular/core';

import { Observable, Subscription } from 'rxjs';

import { Chart } from 'chart.js'; // TODO: use plotly instead for WebGL-capabale browsers?

import { DataSet } from '../models/data';

export enum LineChartType {
	Category = 'category',
	Linear = 'linear',
	Logarithmic = 'logarithmic',
	Time = 'time'
}

@Directive({
	selector: '[linechart]',
	inputs: [
		'chartTitle',
		'chartLabels',
		'chartDataSets',
		'chartType',
		'chartXAxisLabel',
		'chartYAxisLabel',
		'chartLabelCallback',
		'chartDisplayLegend'
	]
})
export class LinechartDirective implements AfterViewInit, OnDestroy {

	ctx: CanvasRenderingContext2D; // | WebGLRenderingContext;
	chart: Chart;

	chartTitle?: string;
	chartLabels?: any[];
	chartDataSets: Observable<DataSet[]>;
	chartType?: LineChartType;
	chartXAxisLabel?: string;
	chartYAxisLabel?: string;
	chartLabelCallback?: (v: any, i: number, va: any[]) => any;
	chartDisplayLegend?: boolean;

	private subscription: Subscription;
	private opts: any;

	constructor (private readonly element: ElementRef) { }

	ngAfterViewInit () {
		if (this.element.nativeElement === null) {
			console.warn('Use of DOM directive in non-DOM context!');
			return;
		}

		if (!(this.element.nativeElement instanceof HTMLCanvasElement)) {
			throw new Error('[linechart] Directive can only be used on a canvas!');
		}

		this.ctx = (this.element.nativeElement as HTMLCanvasElement).getContext('2d', {alpha: false});

		if (!this.chartType) {
			this.chartType = LineChartType.Linear;
		}

		if (this.chartDisplayLegend === null || this.chartDisplayLegend === undefined) {
			this.chartDisplayLegend = false;
		}

		this.opts = {
			type: 'line',
			data: {
				datasets: [],
			},
			options: {
				legend: {
					display: true
				},
				tooltips: {
					intersect: false,
					mode: 'x'
				},
				title: {
					display: this.chartTitle ? true : false,
					text: this.chartTitle
				},
				scales: {
					xAxes: [{
						display: true,
						type: this.chartType,
						scaleLabel: {
							display: this.chartXAxisLabel ? true : false,
							labelString: this.chartXAxisLabel
						}
					}],
					yAxes: [{
						display: true,
						ticks: {
							suggestedMin: 0
						},
						scaleLabel: {
							display: this.chartYAxisLabel ? true : false,
							labelString: this.chartYAxisLabel
						}
					}]
				}
			}
		};

		this.subscription = this.chartDataSets.subscribe((data: DataSet[]) => {this.dataLoad(data); },
		                                                 (e: Error) => {this.dataError(e); });
	}

	private destroyChart () {
		if (this.chart) {
			this.chart.clear();
			this.chart.destroy();
			this.chart = null;
			this.ctx.clearRect(0, 0, this.element.nativeElement.width, this.element.nativeElement.height);
			this.opts.data = {datasets: [], labels: []};
		}
	}


	private dataLoad (data: DataSet[]) {
		this.destroyChart();

		if (data === null || data === undefined || data.some(x => x === null)) {
			this.noData();
			return;
		}

		this.opts.data.datasets = data;

		this.chart = new Chart(this.ctx, this.opts);
	}

	private dataError (e: Error) {
		this.destroyChart();
		this.ctx.font = '30px serif';
		this.ctx.fillStyle = 'black';
		this.ctx.textAlign = 'center';
		this.ctx.fillText('Error Fetching Data', this.element.nativeElement.width / 2., this.element.nativeElement.height / 2.);
	}

	private noData () {
		this.destroyChart();
		this.ctx.font = '30px serif';
		this.ctx.fillStyle = 'black';
		this.ctx.textAlign = 'center';
		this.ctx.fillText('No Data', this.element.nativeElement.width / 2., this.element.nativeElement.height / 2.);
	}

	ngOnDestroy () {
		this.destroyChart();

		if (this.subscription) {
			this.subscription.unsubscribe();
		}
	}

}
