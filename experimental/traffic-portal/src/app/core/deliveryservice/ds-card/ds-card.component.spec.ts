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
import { HttpClientModule } from "@angular/common/http";
import { type ComponentFixture, TestBed, fakeAsync, tick } from "@angular/core/testing";
import { RouterTestingModule } from "@angular/router/testing";
import { protocolToString, Protocol as LibProtocol } from "trafficops-types";

import { DeliveryServiceService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { DsCardComponent } from "src/app/core/deliveryservice/ds-card/ds-card.component";
import { GeoLimit, GeoProvider, Protocol } from "src/app/models";
import { LinechartDirective } from "src/app/shared/charts/linechart.directive";
import { LoadingComponent } from "src/app/shared/loading/loading.component";

describe("DsCardComponent", () => {
	let component: DsCardComponent;
	let fixture: ComponentFixture<DsCardComponent>;
	let api: DeliveryServiceService;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [
				DsCardComponent,
				LoadingComponent,
				LinechartDirective
			],
			imports: [
				APITestingModule,
				HttpClientModule,
				RouterTestingModule,
			],
		}).compileComponents();
		api = TestBed.inject(DeliveryServiceService);
		fixture = TestBed.createComponent(DsCardComponent);
		component = fixture.componentInstance;
		component.deliveryService = await api.createDeliveryService({
			active: false,
			anonymousBlockingEnabled: false,
			cdnId: -1,
			displayName: "",
			dscp: 0,
			geoLimit: GeoLimit.NONE,
			geoProvider: GeoProvider.MAX_MIND,
			ipv6RoutingEnabled: true,
			logsEnabled: true,
			longDesc: "",
			missLat: 0,
			missLong: 0,
			multiSiteOrigin: false,
			regionalGeoBlocking: false,
			routingName: "",
			tenantId: -1,
			typeId: -1,
			xmlId: ""
		});
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it("renders protocol strings", () => {
		expect(component.protocolString).toBe("");

		component.deliveryService.protocol = Protocol.HTTP;
		expect(component.protocolString).toBe(protocolToString(component.deliveryService.protocol as unknown as LibProtocol));
		component.deliveryService.protocol = Protocol.HTTPS;
		expect(component.protocolString).toBe(protocolToString(component.deliveryService.protocol as unknown as LibProtocol));
		component.deliveryService.protocol = Protocol.HTTP_TO_HTTPS;
		expect(component.protocolString).toBe(protocolToString(component.deliveryService.protocol as unknown as LibProtocol));
		component.deliveryService.protocol = Protocol.HTTP_AND_HTTPS;
		expect(component.protocolString).toBe(protocolToString(component.deliveryService.protocol as unknown as LibProtocol));
	});

	it("toggles its open state, and loads its data", fakeAsync(() => {
		expect(component.open).toBeFalse();
		component.toggle();
		tick();
		expect(component.open).toBeTrue();
		expect(component.graphDataLoaded).toBeTrue();
		component.toggle();
		tick();
		expect(component.open).toBeFalse();
		expect(component.graphDataLoaded).toBeFalse();
	}));
});
