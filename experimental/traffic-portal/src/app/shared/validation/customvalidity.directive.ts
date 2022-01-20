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
import type { Observable, Subscription } from "rxjs";

/**
 * CustomvalidityDirective decorates inputs, adding custom validity messages to
 * them.
 */
@Directive({
	selector: "input[customvalidity]"
})
export class CustomvalidityDirective implements AfterViewInit, OnDestroy {
	/**
	 * An Observable that will emit validity messages when the validity state of
	 * the input changes. An empty string signifies "valid", whereas any other
	 * string should be a description of why the input is invalid.
	 */
	@Input() public customvalidity!: Observable<string>;

	/** A subscription for the 'validity' input. */
	private subscription!: Subscription;

	/**
	 * Constructor.
	 */
	constructor(private readonly element: ElementRef<HTMLInputElement>) { }

	/** Initializes the validity state of the element. */
	public ngAfterViewInit(): void {
		this.subscription = this.customvalidity.subscribe(
			s => {
				if (s) {
					this.element.nativeElement.setCustomValidity(s);
					this.element.nativeElement.reportValidity();
				} else {
					this.element.nativeElement.setCustomValidity("");
				}
			}
		);
		this.element.nativeElement.addEventListener("input", () => this.element.nativeElement.setCustomValidity(""));
	}

	/**
	 * Cleans up subscription after the element is destroyed.
	 */
	public ngOnDestroy(): void {
		this.subscription.unsubscribe();
	}

}
