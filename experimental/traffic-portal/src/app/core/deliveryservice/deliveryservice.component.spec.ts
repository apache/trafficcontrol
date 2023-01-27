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
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { RouterTestingModule } from "@angular/router/testing";
import {ReplaySubject} from "rxjs";

import { DeliveryServiceService } from "src/app/api";
import { APITestingModule } from "src/app/api/testing";
import { GeoLimit, GeoProvider } from "src/app/models";
import { AlertService } from "src/app/shared/alert/alert.service";
import { LinechartDirective } from "src/app/shared/charts/linechart.directive";
import { CurrentUserService } from "src/app/shared/currentUser/current-user.service";
import { NavigationService } from "src/app/shared/navigation/navigation.service";
import { TpHeaderComponent } from "src/app/shared/navigation/tp-header/tp-header.component";

import { DeliveryserviceComponent } from "./deliveryservice.component";

describe("DeliveryserviceComponent", () => {
	let component: DeliveryserviceComponent;
	let fixture: ComponentFixture<DeliveryserviceComponent>;

	beforeEach(async () => {
		// mock the API
		const mockCurrentUserService = jasmine.createSpyObj(["updateCurrentUser", "hasPermission", "login", "logout"]);
		const navSvc = jasmine.createSpyObj([],{headerHidden: new ReplaySubject<boolean>(), headerTitle: new ReplaySubject<string>()});

		await TestBed.configureTestingModule({
			declarations: [
				DeliveryserviceComponent,
				TpHeaderComponent,
				LinechartDirective
			],
			imports: [
				APITestingModule,
				FormsModule,
				HttpClientModule,
				ReactiveFormsModule,
				RouterTestingModule
			],
			providers: [
				AlertService,
				{ provide: CurrentUserService, useValue: mockCurrentUserService },
				{ provide: NavigationService, useValue: navSvc}
			]
		}).compileComponents();
		const dsService = TestBed.inject(DeliveryServiceService);
		const ds = await dsService.createDeliveryService({
			active: false,
			anonymousBlockingEnabled: false,
			cdnId: -1,
			displayName: "FIZZbuzz",
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
			xmlId: "fizz-buzz"
		});

		fixture = TestBed.createComponent(DeliveryserviceComponent);
		component = fixture.componentInstance;
		component.deliveryservice = ds;
		fixture.detectChanges();
	});

	it("should exist", () => {
		expect(component).toBeTruthy();
	});

	it("loads TPS and bandwidth data charts", fakeAsync(() => {
		component.toDate.setValue("2022-01-01");
		component.fromDate.setValue("2022-01-01");
		component.toTime.setValue("00:00");
		component.fromTime.setValue("00:20");
		component.newDateRange();
		tick();
		expectAsync(component.bandwidthData.toPromise().then(
			ds => {
				expect(ds.length).toBe(1);
			}
		)).toBeResolved();
		expectAsync(component.tpsChartData.toPromise().then(
			ds => {
				expect(ds.length).toBe(5);
			}
		)).toBeResolved();
	}));

});
