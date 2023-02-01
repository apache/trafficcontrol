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

import { Component, type OnInit } from "@angular/core";
import { FormControl } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import type { ColDef } from "ag-grid-community";
import { BehaviorSubject, Subscribable } from "rxjs";

import { NavigationService } from "../navigation/navigation.service";

import type { ContextMenuActionEvent, ContextMenuItem } from "./generic-table.component";

/**
 * AbstractTableComponent subclasses should be wrappers for a generic table.
 * They should all also use the inheritable template
 * "abstract-table.component.html" from this directory.
 */
@Component({
	template: ""
})
export abstract class AbstractTableComponent<T extends object> implements OnInit {
	/** The data to be put into the table. */
	public abstract data: T[] | Promise<T[]> | Subscribable<T[]>;
	/** Definitions for the columns of the table. */
	public abstract columnDefs: ColDef<T>[];
	/** A list of items to be displayed when the user opens a context menu */
	public abstract contextMenuItems: ContextMenuItem<T>[];
	/**
	 * The table's unique context; this is used to save settings, so it should
	 * be unique for every unique data set e.g. "servers" is different from
	 * "deliveryservice-servers".
	 */
	public abstract readonly context: string;
	/** A human-friendly name for the table. */
	public abstract readonly tableName: string;

	/**
	 * A subject that child components can subscribe to for access to the fuzzy
	 * search query text.
	 */
	public fuzzySubject = new BehaviorSubject("");
	/** Form controller for the user search input. */
	public fuzzControl = new FormControl("", {nonNullable: true});

	/**
	 * The type of an optional floating-action button for the table. If it's a
	 * "link", then the inserted FAB will be an "A" tag and it will link to the
	 * {@link fabLink}. If it's a "button", the inserted FAB will be a "BUTTON"
	 * tag and clicks will trigger {@link handleFAB}. If it's `null`, no FAB
	 * will be inserted at all.
	 */
	public fabType: "link" | "button" | null = null;
	/**
	 * The icon to use for the FAB. The default is what most pages probably
	 * want, a simple "+".
	 */
	public fabIcon = "add";
	/** The title of the FAB (if any). Pages with FABs should override this. */
	public fabTitle = "Add new";
	/** A link for "link"-type FABs. It's relative to the current URL. */
	public fabLink = "#";

	constructor(protected readonly route: ActivatedRoute, private readonly navSvc: NavigationService) {
	}

	/**
	 * A handler for context menu items.
	 *
	 * @param a The action selected by the user.
	 */
	public abstract handleContextMenu(a: ContextMenuActionEvent<T>): void | PromiseLike<void>;

	/** Angular lifecycle hook. */
	public ngOnInit(): void {
		this.navSvc.headerTitle.next(this.tableName);
		this.route.queryParamMap.subscribe(
			m => {
				const search = m.get("search");
				if (search) {
					this.fuzzControl.setValue(decodeURIComponent(search));
					this.updateURL();
				}
			}
		);
	}

	/** Update the URL's 'search' query parameter for the user's search input. */
	public updateURL(): void {
		this.fuzzySubject.next(this.fuzzControl.value);
	}

	/**
	 * This method is expected to return a value indicating that a user is or
	 * isn't allowed to use the page's FAB. A table with a FAB should override
	 * this, as the default behavior is to always disallow.
	 *
	 * @returns `true` if the user is allowed, `false` otherwise.
	 */
	public fabPermission(): boolean {
		return false;
	}

	/**
	 * Handles a floating action button click. The default implementation does
	 * nothing but cancel propagation. Tables with a FAB should override this.
	 *
	 * @param e The click event, in case it's needed.
	 */
	public handleFAB(e: MouseEvent): void | PromiseLike<void> {
		e.stopPropagation();
	}
}
