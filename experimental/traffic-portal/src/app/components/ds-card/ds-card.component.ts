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
import { Component, Input, OnInit } from "@angular/core";

import { Subject } from "rxjs";
import { first } from "rxjs/operators";

import { DataPoint, DataSet, DeliveryService, GeoProvider, protocolToString } from "../../models";
import { DeliveryServiceService } from "../../services/api";


/**
 * DsCardComponent is a component for displaying information about a Delivery
 * Service in an expand-able card.
 */
@Component({
	selector: "ds-card[deliveryService]",
	styleUrls: ["./ds-card.component.scss"],
	templateUrl: "./ds-card.component.html"
})
export class DsCardComponent implements OnInit {

	/** The Delivery Service being described by this component. */
	@Input() public deliveryService: DeliveryService;

	/**
	 * The date to use as the 'current' date/time - the end of the date/time
	 * range for the chart data.
	 *
	 * The default is the time of the component's creation.
	 */
	@Input() public now: Date = new Date();

	/**
	 * The date to use as the 'beginning of the current day' - the start of the
	 * date/time range for the chart data.
	 *
	 * The default is the time of the component's creation (which is the same
	 * as'now', so you really ought to specify!)
	 */
	@Input() public today: Date = new Date();

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
	public chartData: Subject<Array<DataSet | null>>;

	/** Bandwidth data at the Mid-tier level. */
	private readonly midBandwidthData: DataSet;
	/** Bandwidth data at the Edge-tier level. */
	private readonly edgeBandwidthData: DataSet;

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

	/**
	 * Constructor.
	 */
	constructor(private readonly dsAPI: DeliveryServiceService) {
		this.deliveryService = {
			active: true,
			anonymousBlockingEnabled: false,
			cdnId: -1,
			displayName: "",
			dscp: -1,
			geoLimit: 0,
			geoProvider: GeoProvider.MAX_MIND,
			ipv6RoutingEnabled: true,
			logsEnabled: true,
			longDesc: "",
			missLat: 0,
			missLong: 0,
			multiSiteOrigin: false,
			regionalGeoBlocking: false,
			routingName: "",
			typeId: -1,
			xmlId: "-1"
		};

		this.available = 100;
		this.maintenance = 0;
		this.utilized = 0;
		this.loaded = false;
		this.open = false;
		this.chartData = new Subject<Array<DataSet | null>>();

		this.edgeBandwidthData = {
			borderColor: "#3CBA9F",
			data: new Array<DataPoint>(),
			fill: true,
			fillColor: "#3CBA9F",
			label: "Edge-tier Bandwidth"
		};

		this.midBandwidthData = {
			borderColor: "#BA3C57",
			data: new Array<DataPoint>(),
			fill: true,
			fillColor: "#BA3C57",
			label: "Mid-tier Bandwidth"
		};

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
	 * @param e A DOM Event caused then the open/close state of a `<details>` element changes.
	 *
	 * this will only fetch health and capacity information once per page load,
	 * but will update the Bandwidth graph every time the details panel is
	 * opened. Bandwidth data is calculated using 60s intervals starting at
	 * 00:00 UTC the current day and ending at the current date/time.
	 */
	public toggle(e: Event): void {
		if (!this.deliveryService.id) {
			console.error("Toggling DS card for DS with no ID:", this.deliveryService);
			return;
		}
		if ((e.target as HTMLDetailsElement).open) {
			if (!this.loaded) {
				this.loaded = true;
				this.dsAPI.getDSCapacity(this.deliveryService.id).pipe(first()).subscribe(
					r => {
						if (r) {
							this.available = r.availablePercent;
							this.maintenance = r.maintenancePercent;
							this.utilized = r.utilizedPercent;
						}
					}
				);
				this.dsAPI.getDSHealth(this.deliveryService.id).pipe(first()).subscribe(
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
			this.chartData.next([]);
		}
	}

	/**
	 * Requests new data for the charts from the API and loads it.
	 */
	private loadChart(): void {
		const xmlID = this.deliveryService.xmlId;
		this.dsAPI.getDSKBPS(xmlID, this.today, this.now, "1m", false, true).pipe(first()).subscribe(
			data => {
				for (const d of data) {
					if (d.y === null) {
						continue;
					}
					this.edgeBandwidthData.data.push(d);
				}
				this.chartData.next([this.edgeBandwidthData, this.midBandwidthData]);
				this.graphDataLoaded = true;
			},
			(e: Error) => {
				this.graphDataLoaded = true;
				this.chartData.next([null, this.midBandwidthData]);
				console.error(`Failed getting edge KBPS for DS '${xmlID}':`, e);
			}
		);

		this.dsAPI.getDSKBPS(xmlID, this.today, this.now, "1m", true, true).pipe(first()).subscribe(
			(data: Array<DataPoint>) => {
				for (const d of data) {
					if (d.y === null) {
						continue;
					}
					this.midBandwidthData.data.push(d);
				}
				this.chartData.next([this.edgeBandwidthData, this.midBandwidthData]);
				this.graphDataLoaded = true;
			},
			(e: Error) => {
				this.chartData.next([this.edgeBandwidthData, null]);
				this.graphDataLoaded = true;
				console.error(`Failed getting mid KBPS for DS '${xmlID}':`, e);
			}
		);
	}
}
