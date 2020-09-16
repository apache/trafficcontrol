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
import { AfterViewInit, Directive, ElementRef, Input, OnDestroy } from '@angular/core';

import { Observable, Subscription } from 'rxjs';

@Directive({
	selector: 'input[customvalidity]'
})
export class CustomvalidityDirective implements AfterViewInit, OnDestroy {
	@Input('validity') validity: Observable<string>;

	private subscription: Subscription;

	constructor (private readonly element: ElementRef<HTMLInputElement>) { }

	ngAfterViewInit () {
		if (!this.element.nativeElement) {
			console.warn('Use of DOM directive in non-DOM context!');
			return;
		}

		this.subscription = this.validity.subscribe(
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
		this.element.nativeElement.addEventListener('input', unused_e => this.element.nativeElement.setCustomValidity(''));
	}

	ngOnDestroy () {
		this.subscription.unsubscribe();
	}

}
