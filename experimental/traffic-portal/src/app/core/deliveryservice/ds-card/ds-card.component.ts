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
import { trigger, style, animate, transition } from "@angular/animations";
import { Component, Input, type OnInit } from "@angular/core";
import { Subject } from "rxjs";
import { protocolToString, type ResponseDeliveryService } from "trafficops-types";

import { DeliveryServiceService } from "src/app/api";
import type {
	DataPoint,
	DataSet,
} from "src/app/models";
import { LoggingService } from "src/app/shared/logging.service";

/**
 * DsCardComponent is a component for displaying information about a Delivery
 * Service in an expand-able card.
 */
@Component({
	animations: [
		trigger(
			"enterAnimation", [
				transition(":enter", [
					style({opacity: 0, transform: "translateY(-100%)"}),
					animate("200ms", style({opacity: 1, transform: "translateY(0)"}))
				]),
				transition(":leave", [
					style({opacity: 1, transform: "translateY(0)"}),
					animate("200ms", style({opacity: 0, transform: "translateY(-100%)"}))
				])
			]
		)
	],
	selector: "ds-card[deliveryService]",
	styleUrls: ["./ds-card.component.scss"],
	templateUrl: "./ds-card.component.html"
})
export class DsCardComponent implements OnInit {

	/** The Delivery Service being described by this component. */
	@Input() public deliveryService!: ResponseDeliveryService;

	/** Whether or not this is the first DS Card in a list. Affects styling. */
	@Input() public first = false;
	/** Whether or not this is the last DS Card in a list. Affects styling. */
	@Input() public last = false;

	/**
	 * The date to use as the 'current' date/time - the end of the date/time
	 * range for the chart data.
	 *
	 * If either `now` or `today` is not given, `now` will be set to the moment
	 * of the component's creation, and `today` will be based on that.
	 */
	@Input() public now!: Date;

	/**
	 * The date to use as the 'beginning of the current day' - the start of the
	 * date/time range for the chart data.
	 *
	 * If either `now` or `today` is not given, `now` will be set to the moment
	 * of the component's creation, and `today` will be based on that.
	 */
	@Input() public today!: Date;

	/**
	 * The number of cache servers available to serve traffic in this Delivery
	 * Service.
	 */
	public available: number;
	/**
	 * The number of cache servers within this Delivery Service currently
	 * undergoing maintenance.
	 */
	public maintenance: number;
	/**
	 * The amount of cache server bandwidth being utilized within this Delivery
	 * Service.
	 */
	public utilized: number;

	/** Health measured as a percent of Cache Groups that are healthy */
	public healthy = 0;

	/** Bandwidth chart data. */
	public chartData: Subject<[DataSet]>;

	/** This must be a member to have access in the template. */
	public protocolToString = protocolToString;

	/** Describes whether or not the card is expanded. */
	public open: boolean;

	/** Describes whether or not the card's data has been loaded. */
	private loaded: boolean;

	/**
	 * Describes whether or not the card's data specific to charts has been
	 * loaded.
	 */
	public graphDataLoaded: boolean;

	/** The Protocol of the Delivery Service as a string. */
	public get protocolString(): string {
		if (this.deliveryService.protocol !== undefined) {
			return protocolToString(this.deliveryService.protocol);
		}
		return "";
	}

	constructor(private readonly dsAPI: DeliveryServiceService, private readonly log: LoggingService) {
		this.available = 100;
		this.maintenance = 0;
		this.utilized = 0;
		this.loaded = false;
		this.open = false;
		this.chartData = new Subject<[DataSet]>();

		this.graphDataLoaded = false;
	}

	/**
	 * Runs initialization, setting 'now' and 'today' if they weren't given as
	 * input.
	 */
	public ngOnInit(): void {
		if (!this.now || !this.today) {
			this.now = new Date();
			this.now.setUTCMilliseconds(0);
			this.today = new Date(this.now.getFullYear(), this.now.getMonth(), this.now.getDate());
		}
	}

	/**
	 * Triggered when the details element is opened or closed, and fetches the
	 * latest stats.
	 *
	 * this will only fetch health and capacity information once per page load,
	 * but will update the Bandwidth graph every time the details panel is
	 * opened. Bandwidth data is calculated using 60s intervals starting at
	 * 00:00 UTC the current day and ending at the current date/time.
	 */
	public toggle(): void {
		if (!this.open) {
			if (!this.loaded) {
				this.loaded = true;
				this.dsAPI.getDSCapacity(this.deliveryService.id).then(
					r => {
						if (r) {
							this.available = r.availablePercent;
							this.maintenance = r.maintenancePercent;
							this.utilized = r.utilizedPercent;
						}
					}
				);
				this.dsAPI.getDSHealth(this.deliveryService.id).then(
					r => {
						if (r) {
							if (r.totalOnline === 0) {
								this.healthy = 0;
							} else {
								this.healthy = Number(r.totalOnline) / (Number(r.totalOnline) + Number(r.totalOffline));
							}
						}
					}
				);
			}
			this.open = true;
			this.loadChart();
		} else {
			this.open = false;
			this.graphDataLoaded = false;
			this.loaded = false;
			this.chartData.next([{
				borderColor: "#3CBA9F",
				data: new Array<DataPoint>(),
				fill: true,
				fillColor: "#3CBA9F",
				label: "Edge-tier Bandwidth"
			}]);
		}
	}

	/**
	 * Requests new data for the charts from the API and loads it.
	 */
	private async loadChart(): Promise<void> {
		const xmlID = this.deliveryService.xmlId;
		try {
			const data = await this.dsAPI.getDSKBPS(xmlID, this.today, this.now, "1m", false, true);
			const chartData = {
				borderColor: "#3CBA9F",
				data: new Array<DataPoint>(),
				fill: true,
				fillColor: "#3CBA9F",
				label: "Edge-tier Bandwidth"
			};
			for (const d of data) {
				if (d.y === null) {
					continue;
				}
				chartData.data.push(d);
			}
			this.chartData.next([chartData]);
			this.graphDataLoaded = true;
		} catch (e) {
			this.graphDataLoaded = true;
			this.chartData.next([{
				borderColor: "#3CBA9F",
				data: new Array<DataPoint>(),
				fill: true,
				fillColor: "#3CBA9F",
				label: "Edge-tier Bandwidth"
			}]);
			this.log.error(`Failed getting edge KBPS for DS '${xmlID}':`, e);
		}
	}
}
