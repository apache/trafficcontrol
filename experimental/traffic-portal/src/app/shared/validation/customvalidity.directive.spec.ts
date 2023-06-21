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
import { Component, type DebugElement, ElementRef } from "@angular/core";
import { type ComponentFixture, TestBed } from "@angular/core/testing";
import { By } from "@angular/platform-browser";
import { BehaviorSubject } from "rxjs";

import { CustomvalidityDirective } from "./customvalidity.directive";

/**
 * This component is used to exercise all the ways the Customvalidity directive
 * can be used.
 */
@Component({
	template: `
		<input type="checkbox" [customvalidity]="validity"/>
		<input type="color" [customvalidity]="validity"/>
		<input type="date" [customvalidity]="validity"/>
		<input type="datetime-local" [customvalidity]="validity"/>
		<input type="email" [customvalidity]="validity"/>
		<input type="month" [customvalidity]="validity"/>
		<input type="number" [customvalidity]="validity"/>
		<input type="password" [customvalidity]="validity"/>
		<input type="range" [customvalidity]="validity"/>
		<input type="tel" [customvalidity]="validity"/>
		<input type="text" [customvalidity]="validity"/>
		<input type="time" [customvalidity]="validity"/>
		<input type="url" [customvalidity]="validity"/>
		<input type="week" [customvalidity]="validity"/>
		<input type="text" id="no-directive-applied"/>
	`
})
class TestComponent {
	public validity = new BehaviorSubject("");
}

describe("CustomvalidityDirective", () => {
	let fixture: ComponentFixture<TestComponent>;
	let customInputs: Array<DebugElement>;
	let normalInput: DebugElement;

	beforeEach(() => {
		fixture = TestBed.configureTestingModule({
			declarations: [CustomvalidityDirective, TestComponent]
		}).createComponent(TestComponent);

		fixture.detectChanges();

		customInputs = fixture.debugElement.queryAll(By.directive(CustomvalidityDirective));
		normalInput = fixture.debugElement.query(By.css("#no-directive-applied"));
	});

	it("should create an instance", () => {
		const directive = new CustomvalidityDirective(new ElementRef(document.createElement("input")));
		expect(directive).toBeTruthy();
	});

	it("binds to only the elements with its directive", () => {
		expect(customInputs.length).toBe(14);
		expect(normalInput).toBeTruthy();
		expect(normalInput.properties.customvalidity).toBeUndefined();
	});

	it("updates input validity when the input emits a change", () => {
		for (const input of customInputs) {
			const rawInput = input.nativeElement as HTMLInputElement;
			expect(rawInput.willValidate).toBeTrue();
			expect(rawInput.validationMessage).toBe("");
			const {customError, valid} = rawInput.validity;
			expect(customError).toBeFalse();
			expect(valid).toBeTrue();
			expect(rawInput.reportValidity()).toBeTrue();
		}
		let rawNormalInput = normalInput.nativeElement as HTMLInputElement;
		expect(rawNormalInput.willValidate).toBeTrue();
		expect(rawNormalInput.validationMessage).toBe("");
		expect(rawNormalInput.validity.customError).toBeFalse();
		expect(rawNormalInput.validity.valid).toBeTrue();
		expect(rawNormalInput.reportValidity()).toBeTrue();

		const customMessage = "my custom validity error message";
		fixture.componentInstance.validity.next(customMessage);
		fixture.detectChanges();

		for (const input of customInputs) {
			const rawInput = input.nativeElement as HTMLInputElement;
			expect(rawInput.willValidate).toBeTrue();
			expect(rawInput.validationMessage).toBe(customMessage);
			const {customError, valid} = rawInput.validity;
			expect(customError).toBeTrue();
			expect(valid).toBeFalse();
			expect(rawInput.reportValidity()).toBeFalse();
		}

		rawNormalInput = normalInput.nativeElement as HTMLInputElement;
		expect(rawNormalInput.willValidate).toBeTrue();
		expect(rawNormalInput.validationMessage).toBe("");
		expect(rawNormalInput.validity.customError).toBeFalse();
		expect(rawNormalInput.validity.valid).toBeTrue();
		expect(rawNormalInput.reportValidity()).toBeTrue();

		fixture.componentInstance.validity.next("");
		fixture.detectChanges();

		for (const input of customInputs) {
			const rawInput = input.nativeElement as HTMLInputElement;
			expect(rawInput.willValidate).toBeTrue();
			expect(rawInput.validationMessage).toBe("");
			const {customError, valid} = rawInput.validity;
			expect(customError).toBeFalse();
			expect(valid).toBeTrue();
			expect(rawInput.reportValidity()).toBeTrue();
		}
		rawNormalInput = normalInput.nativeElement as HTMLInputElement;
		expect(rawNormalInput.willValidate).toBeTrue();
		expect(rawNormalInput.validationMessage).toBe("");
		expect(rawNormalInput.validity.customError).toBeFalse();
		expect(rawNormalInput.validity.valid).toBeTrue();
		expect(rawNormalInput.reportValidity()).toBeTrue();
	});

	it("clears validity on input", () => {
		fixture.componentInstance.validity.next("invalid");

		for (const input of customInputs) {
			const rawInput = input.nativeElement as HTMLInputElement;
			expect(rawInput.validity.valid).toBeFalse();
			rawInput.dispatchEvent(new InputEvent("input", {data: "some text"}));
			expect(rawInput.validity.valid).toBeTrue();
		}
	});
});
