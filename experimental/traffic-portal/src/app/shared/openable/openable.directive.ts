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
import { AfterViewInit, Directive, ElementRef, Input, OnDestroy } from "@angular/core";

import { Observable, of, Subscription } from "rxjs";

/**
 * OpenableDirective allows or toggle-able dialogs. This is essentially a
 * polyfill for browsers that don't support true dialog elements.
 */
@Directive({
	selector: "dialog[openable][toggle]"
})
export class OpenableDirective implements AfterViewInit, OnDestroy {
	/** An Observable that emits toggle states for the dialog. */
	@Input() public toggle: Observable<boolean> = of(true);

	/** A subscription for the toggle input. */
	private subscription: Subscription | null = null;

	/**
	 * Constructor.
	 */
	constructor(private readonly element: ElementRef) { }

	/** Initializes toggle listening. */
	public ngAfterViewInit(): void {
		if (!this.element.nativeElement) {
			console.warn("Use of DOM directive in non-DOM context!");
			return;
		}

		this.subscription = this.toggle.subscribe(
			v => {
				if (v) {
					if (!this.element.nativeElement.open) {
						this.element.nativeElement.showModal();
					} else {
						console.warn("Attempted to open dialog that is already open!");
					}
				} else if (this.element.nativeElement.open) {
					this.element.nativeElement.close();
				} else {
					console.warn("Attempted to close dialog that is already closed!");
				}
			},
			e => {
				console.error(e);
			}
		);
	}

	/** cleans up subscriptions on element destruction. */
	public ngOnDestroy(): void {
		if (this.subscription) {
			this.subscription.unsubscribe();
		}
	}
}
