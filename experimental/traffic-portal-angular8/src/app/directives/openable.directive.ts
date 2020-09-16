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
	selector: "dialog[openable]"
})
export class OpenableDirective implements AfterViewInit, OnDestroy {
	@Input('toggle') toggle: Observable<boolean>;

	private subscription: Subscription;

	constructor (private readonly element: ElementRef) { }

	ngAfterViewInit () {
		if (!this.element.nativeElement) {
			console.warn('Use of DOM directive in non-DOM context!');
			return;
		}

		this.subscription = this.toggle.subscribe(
			v => {
				if (v) {
					if (!this.element.nativeElement.open) {
						this.element.nativeElement.showModal();
					} else {
						console.warn('Attempted to open dialog that is already open!');
					}
				} else if (this.element.nativeElement.open) {
					this.element.nativeElement.close();
				} else {
					console.warn('Attempted to close dialog that is already closed!');
				}
			},
			e => {
				console.error(e);
			}
		);
	}

	ngOnDestroy () {
		this.subscription.unsubscribe();
	}
}
