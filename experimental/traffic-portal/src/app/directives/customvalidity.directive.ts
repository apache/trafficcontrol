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
 * CustomvalidityDirective decorates inputs, adding custom validity messages to
 * them.
 */
@Directive({
	selector: "input[customvalidity][valid]"
})
export class CustomvalidityDirective implements AfterViewInit, OnDestroy {
	/**
	 * An Observable that will emit validity messages when the validity state of
	 * the input changes. An empty string signifies "valid", whereas any other
	 * string should be a description of why the input is invalid.
	 */
	@Input() public valid: Observable<string> = of("");

	/** A subscription for the 'validity' input. */
	private subscription: Subscription | null = null;

	/**
	 * Constructor.
	 */
	constructor(private readonly element: ElementRef<HTMLInputElement>) { }

	/** Initializes the validity state of the element. */
	public ngAfterViewInit(): void {
		if (!this.element.nativeElement) {
			console.warn("Use of DOM directive in non-DOM context!");
			return;
		}

		this.subscription = this.valid.subscribe(
			s => {
				if (s) {
					this.element.nativeElement.setCustomValidity(s);
					this.element.nativeElement.reportValidity();
				}
			},
			e => {
				console.error(e);
			}
		);
		this.element.nativeElement.addEventListener("input", () => this.element.nativeElement.setCustomValidity(""));
	}

	/**
	 * Cleans up subscription after the element is destroyed.
	 */
	public ngOnDestroy(): void {
		if (this.subscription) {
			this.subscription.unsubscribe();
		}
	}

}
