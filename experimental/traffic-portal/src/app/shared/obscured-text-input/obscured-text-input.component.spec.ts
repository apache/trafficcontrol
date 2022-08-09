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
import { ComponentFixture, TestBed } from "@angular/core/testing";

import { AutocompleteValue } from "src/app/utils";

import { ObscuredTextInputComponent } from "./obscured-text-input.component";

const name = "property";

describe("ObscuredTextInputComponent", () => {
	let component: ObscuredTextInputComponent;
	let fixture: ComponentFixture<ObscuredTextInputComponent>;

	beforeEach(async () => {
		await TestBed.configureTestingModule({
			declarations: [ ObscuredTextInputComponent ]
		})
			.compileComponents();

		fixture = TestBed.createComponent(ObscuredTextInputComponent);
		component = fixture.componentInstance;
		component.name = name;
		fixture.detectChanges();
	});

	it("should create", () => {
		expect(component).toBeTruthy();
		expect(component.name).toBe(name);
	});

	it("gets the error state correctly", () => {
		expect(component.errorState).toBeFalse();
		component.control.setErrors({customError: true});
		expect(component.errorState).toBeFalse();
		component.touched = true;
		expect(component.errorState).toBeTrue();
		component.control.setErrors(null);
		expect(component.errorState).toBeFalse();
	});

	it("gets and sets the value", () => {
		expect(component.value).toBe("");
		component.value = null as string | null;
		expect(component.value).toBe("");
		const value = "foo";
		component.value = value;
		expect(component.value).toBe(value);
		component.control.setErrors({customError: true});
		expect(component.value).toBeNull();
	});

	it("emits state changes", () => {
		const spyName = "subscriber";
		const spy = jasmine.createSpy(spyName);
		const subscription = component.stateChanges.subscribe(spy);
		expect(spy).not.toHaveBeenCalled();
		component.value = "foo";
		expect(spy).toHaveBeenCalled();
		component.autocomplete = AutocompleteValue.ON;
		expect(spy).toHaveBeenCalledTimes(2);
		component.placeholder = "placeholder";
		expect(spy).toHaveBeenCalledTimes(3);
		component.disabled = true;
		expect(spy).toHaveBeenCalledTimes(4);
		component.required = true;
		expect(spy).toHaveBeenCalledTimes(5);
		component.maxLength = 10;
		expect(spy).toHaveBeenCalledTimes(6);
		component.minLength = 5;
		expect(spy).toHaveBeenCalledTimes(7);
		component.focused = true;
		expect(spy).toHaveBeenCalledTimes(8);
		component.focus();
		expect(spy).toHaveBeenCalledTimes(8);
		component.blur(new FocusEvent("blur"));
		expect(spy).toHaveBeenCalledTimes(8);
		if (!(component.elementRef.nativeElement instanceof HTMLElement)) {
			return fail("element ref not a reference to an element");
		}
		component.blur(new FocusEvent("blur", {relatedTarget: component.elementRef.nativeElement.parentElement}));
		expect(spy).toHaveBeenCalledTimes(9);
		component.focus();
		expect(spy).toHaveBeenCalledTimes(10);
		component.writeValue("testquest");
		expect(spy).toHaveBeenCalledTimes(11);
		component.onChange(new Event("change"));
		expect(spy).toHaveBeenCalledTimes(12);
		component.setDisabledState(false);
		expect(spy).toHaveBeenCalledTimes(13);
		subscription.unsubscribe();
		component.ngOnDestroy();
		expect(component.stateChanges.isStopped).toBeTrue();
	});

	it("gets and sets its inputs with the right default values", () => {
		// Defaults
		expect(component.minLength).toBe(-1);
		expect(component.maxLength).toBe(-1);
		expect(component.touched).toBeFalse();
		expect(component.autocomplete).toBe(AutocompleteValue.OFF);
		expect(component.userDescribedBy).toBe("");
		expect(component.required).toBeFalse();
		expect(component.placeholder).toBe("");

		// Changes
		component.minLength = 5;
		expect(component.minLength).toBe(5);
		component.maxLength = 5;
		expect(component.maxLength).toBe(5);
		component.touched = true;
		expect(component.touched).toBe(true);
		component.autocomplete = AutocompleteValue.ON;
		expect(component.autocomplete).toBe(AutocompleteValue.ON);
		component.userDescribedBy = "element-id another-element-id";
		expect(component.userDescribedBy).toBe("element-id another-element-id");
		component.required = true;
		expect(component.required).toBeTrue();
		component.placeholder = "placeholder";
		expect(component.placeholder).toBe("placeholder");
	});

	it("sets describedby from Angular forms", () => {
		expect(component.describedBy).toBe("");
		component.setDescribedByIds(["element-id", "another-element-id"]);
		expect(component.describedBy).toBe("element-id another-element-id");
	});

	it("registers onChange callbacks", () => {
		if (!(component.elementRef.nativeElement instanceof HTMLElement)) {
			return fail("element ref not set to a reference");
		}
		const spy = jasmine.createSpy("changeSpy");
		component.registerOnChange(spy);
		expect(spy).not.toHaveBeenCalled();
		const input = component.elementRef.nativeElement.querySelector("input");
		input?.dispatchEvent(new Event("change"));
		expect(spy).toHaveBeenCalled();
	});

	it("registers onTouch callbacks", () => {
		if (!(component.elementRef.nativeElement instanceof HTMLElement)) {
			return fail("element ref not set to a reference");
		}
		const spy = jasmine.createSpy("touchSpy");
		component.registerOnTouched(spy);
		expect(spy).not.toHaveBeenCalled();
		component.blur(new FocusEvent("blur", {relatedTarget: component.elementRef.nativeElement.parentElement}));
		expect(spy).toHaveBeenCalled();
	});

	it("toggles state", () => {
		expect(component.type).toBe("password");
		component.toggle();
		expect(component.type).toBe("text");
		component.toggle();
		expect(component.type).toBe("password");
	});

	it("focuses the input element on container clicks", () => {
		if (!(component.elementRef.nativeElement instanceof HTMLElement)) {
			return fail("element ref not set to a reference");
		}
		const input = component.elementRef.nativeElement.querySelector("input");
		if (!input) {
			return fail("component doesn't contain an input element");
		}
		const button = component.elementRef.nativeElement.querySelector("button");
		if (!button) {
			return fail("component doesn't contain a button element");
		}
		expect(document.activeElement).not.toBe(input);
		component.onContainerClick({target: button} as unknown as MouseEvent);
		expect(document.activeElement).not.toBe(input);
		component.onContainerClick({target: input} as unknown as MouseEvent);
		expect(document.activeElement).not.toBe(input);
		component.onContainerClick({target: component.elementRef.nativeElement} as unknown as MouseEvent);
		expect(document.activeElement).toBe(input);
	});
});
