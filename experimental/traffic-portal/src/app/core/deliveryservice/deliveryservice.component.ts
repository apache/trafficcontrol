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
import { Component, type OnInit } from "@angular/core";
import { UntypedFormControl } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { Subject } from "rxjs";
import { AlertLevel, ResponseDeliveryService } from "trafficops-types";

import { DeliveryServiceService } from "src/app/api";
import type { DataPoint, DataSet } from "src/app/models";
import { AlertService } from "src/app/shared/alert/alert.service";
import { LoggingService } from "src/app/shared/logging.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";

/**
 * DeliveryserviceComponent is the controller for a single Delivery Service's
 * details page.
 */
@Component({
	selector: "tp-deliveryservice",
	styleUrls: ["./deliveryservice.component.scss"],
	templateUrl: "./deliveryservice.component.html"
})
export class DeliveryserviceComponent implements OnInit {

	/** The Delivery Service described by this component. */
	public deliveryservice = {} as ResponseDeliveryService;

	/** Data for the bandwidth chart. */
	public bandwidthData = new Subject<[DataSet]>();

	/** Data for the transactions per second chart. */
	public tpsChartData: Subject<Array<DataSet>>;

	/** End date for charts. */
	private to: Date = new Date();
	/** Start date for charts. */
	private from: Date = new Date();

	/**
	 * Form controller for the user's date selector for the end of the time
	 * range.
	 */
	public fromDate: UntypedFormControl = new UntypedFormControl();

	/**
	 * Form controller for the user's time selector for the end of the time
	 * range.
	 */
	public fromTime: UntypedFormControl = new UntypedFormControl();

	/**
	 * Form controller for the user's date selector for the beginning of the
	 * time range.
	 */
	public toDate: UntypedFormControl = new UntypedFormControl();

	/**
	 * Form controller for the user's date selector for the beginning of the
	 * time range.
	 */
	public toTime: UntypedFormControl = new UntypedFormControl();

	/* Contains the DS xmlIds that this DS is this steering target for. */
	public steeringTargetsFor = new Set<string>();

	/** The size of each single interval for data grouping, in seconds. */
	private bucketSize = 300;

	constructor(
		private readonly route: ActivatedRoute,
		private readonly api: DeliveryServiceService,
		private readonly alerts: AlertService,
		private readonly navSvc: NavigationService,
		private readonly log: LoggingService,
	) {
		this.bandwidthData.next([{
			backgroundColor: "#BA3C57",
			borderColor: "#BA3C57",
			data: new Array<DataPoint>(),
			fill: false,
			label: "Edge-Tier"
		}]);
		this.tpsChartData = new Subject<Array<DataSet>>();
	}

	/**
	 * Runs initialization, including setting up date/time range controls and
	 * fetching data.
	 */
	public ngOnInit(): void {
		this.to.setUTCMilliseconds(0);
		this.from = new Date(this.to.getFullYear(), this.to.getMonth(), this.to.getDate());

		const dateStr = String(this.from.getFullYear()).padStart(4, "0").concat(
			"-", String(this.from.getMonth() + 1).padStart(2, "0").concat(
				"-", String(this.from.getDate()).padStart(2, "0")));

		this.fromDate = new UntypedFormControl(dateStr);
		this.fromTime = new UntypedFormControl("00:00");
		this.toDate = new UntypedFormControl(dateStr);
		const timeStr = String(this.to.getHours()).padStart(2, "0").concat(":", String(this.to.getMinutes()).padStart(2, "0"));
		this.toTime = new UntypedFormControl(timeStr);

		const DSID = this.route.snapshot.paramMap.get("id");
		if (!DSID) {
			this.log.error("Missing route 'id' parameter");
			return;
		}

		this.api.getDeliveryServices(parseInt(DSID, 10)).then(
			d => {
				this.deliveryservice = d;
				this.loadBandwidth();
				this.loadTPS();
				this.navSvc.headerTitle.next(d.displayName);

				this.api.getSteering().then(configs => {
					configs.forEach(config => {
						config.targets.forEach(target => {
							if (target.deliveryService === this.deliveryservice.xmlId) {
								this.steeringTargetsFor.add(config.deliveryService);
							}
						});
					});
				});
			}
		);
	}

	/**
	 * Returns the tooltip text for the steering target displays.
	 *
	 * @returns Tooltip text.
	 */
	public steeringTargetDisplay(): string {
		return `Steering target for: ${Array.from(this.steeringTargetsFor).join(", ")}`;
	}

	/**
	 * Runs when a new date/time range is selected by the user, updating the
	 * chart data accordingly.
	 */
	public newDateRange(): void {
		this.to = new Date(this.toDate.value.concat("T", this.toTime.value));
		this.from = new Date(this.fromDate.value.concat("T", this.fromTime.value));

		// I need to set these explicitly, just in case - the API can't handle millisecond precision
		this.to.setUTCMilliseconds(0);
		this.from.setUTCMilliseconds(0);

		// This should set it to the number of seconds needed to bucket 500 datapoints
		this.bucketSize = (this.to.getTime() - this.from.getTime()) / 30000000;
		this.loadBandwidth();
		this.loadTPS();
	}

	/**
	 * Loads new data for the bandwidth chart.
	 */
	private async loadBandwidth(): Promise<void> {
		let interval: string;
		if (this.bucketSize < 1) {
			interval = "1m";
		} else {
			interval = `${Math.round(this.bucketSize)}m`;
		}

		const xmlID = this.deliveryservice.xmlId;

		let data;
		try {
			data = await this.api.getDSKBPS(xmlID, this.from, this.to, interval, false, true);
		} catch (e) {
			this.alerts.newAlert(AlertLevel.WARNING, "Edge-Tier bandwidth data not found!");
			this.log.error(`Failed to get edge KBPS data for '${xmlID}':`, e);
			return;
		}

		const chartData = {
			backgroundColor: "#BA3C57",
			borderColor: "#BA3C57",
			data,
			fill: false,
			label: "Edge-Tier"
		};
		this.bandwidthData.next([chartData]);
	}

	/**
	 * Loads new data for the TPS chart.
	 */
	private loadTPS(): void {
		let interval: string;
		if (this.bucketSize < 1) {
			interval = "1m";
		} else {
			interval = `${Math.round(this.bucketSize)}m`;
		}

		this.api.getAllDSTPSData(this.deliveryservice.xmlId, this.from, this.to, interval, false).then(
			data => {
				data.total.dataSet.label = "Total";
				data.total.dataSet.borderColor = "#3C96BA";
				data.success.dataSet.label = "Successful Responses";
				data.success.dataSet.borderColor = "#3CBA5F";
				data.redirection.dataSet.label = "Redirection Responses";
				data.redirection.dataSet.borderColor = "#9f3CBA";
				data.clientError.dataSet.label = "Client Error Responses";
				data.clientError.dataSet.borderColor = "#BA9E3B";
				data.serverError.dataSet.label = "Server Error Responses";
				data.serverError.dataSet.borderColor = "#BA3C57";

				this.tpsChartData.next([
					data.total.dataSet,
					data.success.dataSet,
					data.redirection.dataSet,
					data.clientError.dataSet,
					data.serverError.dataSet
				]);
			},
			e => {
				this.log.error(`Failed to get edge TPS data for '${this.deliveryservice.xmlId}':`, e);
				this.alerts.newAlert(AlertLevel.WARNING, "Edge-Tier transaction data not found!");
			}
		);
	}

}
