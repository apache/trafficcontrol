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
import { FocusMonitor } from "@angular/cdk/a11y";
import { BooleanInput, coerceBooleanProperty } from "@angular/cdk/coercion";
import { Component, ElementRef, HostBinding, Input, type OnDestroy, Optional, Self, Output, EventEmitter } from "@angular/core";
import { type ControlValueAccessor, NgControl, FormControl } from "@angular/forms";
import { MatFormField, MatFormFieldControl } from "@angular/material/form-field";
import { Subject } from "rxjs";

import { AutocompleteValue } from "src/app/utils";

/**
 * An ObscuredTextInputComponent implements a form control compatible with
 * mat-form-field (but not `matSuffix`!) that hides its data as a "password"
 * field by default, but provides a toggle-able "reveal" button that allows
 * displaying it in plain text. Note that this form control MUST have a name,
 * regardless of context.
 *
 * This can be bound to a form control or you can do a 2-way binding to the
 * `value` property.
 *
 * @example <caption>Using with a controller's FormControl property</caption>
 * <tp-obscured-text-input [name]="password" [formControl]="passwordControl">
 * </tp-obscured-text-input>
 * @example <caption>Using with a controller's FormGroup property</caption>
 * <tp-obscured-text-input [name]="password" [formControlName]="password">
 * </tp-obscured-text-input>
 * @example <caption>Binding directly to value</caption>
 * <tp-obscured-text-input [name]="password" [(value)]="password">
 * </tp-obscured-text-input>
 */
@Component({
	providers: [{provide: MatFormFieldControl, useExisting: ObscuredTextInputComponent}],
	selector: "tp-obscured-text-input[name]",
	styleUrls: ["./obscured-text-input.component.scss"],
	templateUrl: "./obscured-text-input.component.html"
})
export class ObscuredTextInputComponent implements OnDestroy, MatFormFieldControl<string>, ControlValueAccessor {

	/** The next unique ID for an `ObscuredTextInputComponent` instance. */
	public static nextID = 0;

	/** autocomplete */
	private autoc = AutocompleteValue.OFF;
	/** placeholder */
	private phd = "";
	/** focused */
	private foc = false;
	/** required */
	private req = false;
	/** disabled */
	private dis = false;
	/** maxlength */
	private max = -1;
	/** minlength */
	private min = -1;
	/** register-able trigger that fires on model changes*/
	private onChanged: ((_: unknown) => void) | null = null;
	/** register-able trigger that fires when the control is "touched" */
	private onTouched: (() => void) | null = null;

	/**
	 * The type of the input field of the control, which can be toggled between
	 * an "obscured" `"password"` state (default) and a "revealed" `"text"`
	 * state.
	 */
	public type: "text" | "password" = "password";
	/**
	 * A subject that emits (nothing) whenever the state of the form control
	 * changes.
	 */
	public stateChanges = new Subject<void>();
	/**
	 * The controller for the underlying `<input>` element.
	 */
	public readonly control = new FormControl<string>("");
	/** Tracks whether the user has "touched" the control. */
	public touched = false;
	/** A name Angular uses to track unique form control types. */
	public readonly controlType = "tp-obscured-text-input";

	/** `true` if there is no value entered, `false` otherwise. */
	public get empty(): boolean {
		return !this.value;
	}

	/** `true` if the control is invalid, `false` otherwise. */
	public get errorState(): boolean {
		return this.touched && this.control.invalid;
	}

	/**
	 * `true` if the Angular Material floating label should be in its floating
	 * state, `false` if it should instead be placed inside the control.
	 */
	@HostBinding("class.floating")
	public get shouldLabelFloat(): boolean {
		return !this.empty || this.focused;
	}

	/** A unique ID for the form control. */
	@HostBinding()
	public id = `tp-obscured-text-input-${ObscuredTextInputComponent.nextID++}`;

	// This is necessary to avoid clashing with consumer-set aria-describedby
	// IDs.
	// eslint-disable-next-line @angular-eslint/no-input-rename
	@Input("aria-describedby")
	public userDescribedBy = "";

	/** The form control's `aria-describedby` attribute value. */
	public describedBy = "";

	/** A name for the form control. */
	@Input() public name!: string;

	/** An autocomplete setting for the form control. */
	@Input()
	public get autocomplete(): AutocompleteValue {
		return this.autoc;
	}
	public set autocomplete(a: AutocompleteValue) {
		this.autoc = a;
		this.stateChanges.next();
	}

	/**
	 * The value of the form control. If the control is invalid, the value will
	 * be `null`.
	 */
	@Input()
	public get value(): string | null {
		if (this.control.valid) {
			return this.control.value;
		}
		return null;
	}
	public set value(v: string | null) {
		this.control.setValue(v ?? "");
		this.valueChange.emit(v);
		this.stateChanges.next();
	}

	/**
	 * Emits the value whenever it changes.
	 */
	@Output() public valueChange = new EventEmitter<string | null>();

	/**
	 * Placeholder text for the form control.
	 */
	@Input()
	public get placeholder(): string {
		return this.phd;
	}
	public set placeholder(p: string) {
		this.phd = p;
		this.stateChanges.next();
	}

	/**
	 * Whether the form control should be invalid when empty.
	 */
	@Input()
	public get required(): boolean {
		return this.req;
	}
	public set required(r: BooleanInput) {
		this.req = coerceBooleanProperty(r);
		this.stateChanges.next();
	}

	/**
	 * Whether the form control should be disabled.
	 */
	@Input()
	public get disabled(): boolean {
		return this.dis;
	}
	public set disabled(d: BooleanInput) {
		this.dis = coerceBooleanProperty(d);
		this.stateChanges.next();
	}

	/**
	 * A maximum allowed length.
	 */
	@Input()
	public get maxLength(): number {
		return this.max;
	}
	public set maxLength(n: number) {
		this.max = n;
		this.stateChanges.next();
	}

	/**
	 * A minimum required length.
	 */
	@Input()
	public get minLength(): number {
		return this.min;
	}
	public set minLength(n: number) {
		this.min = n;
		this.stateChanges.next();
	}

	/**
	 * Whether the form control currently has user focus.
	 */
	public get focused(): boolean {
		return this.foc;
	}
	public set focused(f: boolean) {
		this.foc = f;
		this.stateChanges.next();
	}

	constructor(
		private readonly focusMonitor: FocusMonitor,
		public readonly elementRef: ElementRef,
		@Optional() @Self() public ngControl: NgControl | null,
		@Optional() public parentFormField: MatFormField | null,
	) {
		if (this.ngControl !== null) {
			this.ngControl.valueAccessor = this;
		}
	}

	/**
	 * Angular lifecycle hook; cleans up persistent resources.
	 */
	public ngOnDestroy(): void {
		this.stateChanges.complete();
		this.focusMonitor.stopMonitoring(this.elementRef);
	}

	/**
	 * Toggles the obscured/revealed state of the form control's text.
	 */
	public toggle(): void {
		if (this.type === "password") {
			this.type = "text";
		} else {
			this.type = "password";
		}
	}

	/**
	 * Event handler for events where the form control gains focus.
	 */
	public focus(): void {
		if (!this.focused) {
			this.focused = true;
		}
	}

	/**
	 * Event handler for events where the form control loses focus.
	 *
	 * @param e The focus event in question.
	 */
	public blur(e: FocusEvent): void {
		const ref = this.elementRef.nativeElement;
		if ((e.relatedTarget !== null && !ref.contains(e.relatedTarget)) || ref.contains(e.target)) {
			this.focused = false;
			if (!this.touched) {
				this.touched = true;
				if (this.onTouched) {
					this.onTouched();
				}
			}
		}
	}

	/**
	 * Adds extra `aria-describedby` labels to the form control.
	 *
	 * @param ids Any added `aria-describedby` labels added by a Reactive form
	 * group.
	 */
	public setDescribedByIds(ids: string[]): void {
		this.describedBy = ids.join(" ");
	}

	/**
	 * Handles a user clicking on the component's container/host element,
	 * whenever that click can't be directly associated with some element that
	 * it contains.
	 *
	 * @param event The click event in question.
	 */
	public onContainerClick(event: MouseEvent): void {
		if (event.target instanceof HTMLElement) {
			const tag = event.target.tagName.toLowerCase();
			switch(tag) {
				case "button":
				case "input":
					break;
				default:
					this.elementRef.nativeElement.querySelector("input").focus();
			}
		}
	}

	/**
	 * Used by Angular to set the control's value when used in a Form Group.
	 *
	 * @param obj The value being set.
	 */
	public writeValue(obj: string): void {
		this.value = obj;
	}

	/**
	 * Handles changes to the control.
	 *
	 * @param e The input change event.
	 */
	public onChange(e: Event): void {
		this.valueChange.emit(this.value);
		if (this.onChanged && e.target instanceof HTMLInputElement) {
			this.onChanged(e.target.value);
		}
		this.stateChanges.next();
	}

	/**
	 * Registers a callback that will be called with the control's new value
	 * whenever that value changes.
	 *
	 * @param fn A function that will be called when the control's value
	 * changes.
	 */
	public registerOnChange(fn: (_: unknown) => void): void {
		this.onChanged = fn;
	}

	/**
	 * Registers a callback that will be called whenever the user "touches" the
	 * control.
	 *
	 * @param fn A function that will be called when the control is touched.
	 */
	public registerOnTouched(fn: () => void): void {
		this.onTouched = fn;
	}

	/**
	 * Used by Angular to set whether the form control is disabled when it is
	 * used in a Reactive Form group.
	 *
	 * @param disabled `true` if the control should be disabled, `false`
	 * otherwise.
	 */
	public setDisabledState(disabled: boolean): void {
		this.disabled = disabled;
	}
}
