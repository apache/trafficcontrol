/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { MatDialogModule, MAT_DIALOG_DATA } from "@angular/material/dialog";

import { CollectionChoiceDialogComponent, type CollectionChoiceDialogData } from "./collection-choice-dialog.component";

describe("CollectionChoiceDialogComponent", () => {
	let component: CollectionChoiceDialogComponent;
	let fixture: ComponentFixture<CollectionChoiceDialogComponent>;
	const data: CollectionChoiceDialogData = {
		collection: [],
		message: "Choose something",
		title: "Choose"
	};

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ CollectionChoiceDialogComponent ],
			imports: [
				MatDialogModule
			],
			providers: [
				{provide: MAT_DIALOG_DATA, useValue: data}
			]
		}).compileComponents();

		fixture = TestBed.createComponent(CollectionChoiceDialogComponent);
		component = fixture.componentInstance;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
	});

	it("gets an appropriate input label", () => {
		expect(component.label).toBe(data.message);
		data.label = "test";
		expect(component.label).toBe(data.label);
	});
});
