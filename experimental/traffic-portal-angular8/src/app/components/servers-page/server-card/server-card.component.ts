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
import { Component, Input } from "@angular/core";

import { first } from "rxjs/operators";

import { Server, Servercheck } from "../../../models";
import { ServerService } from "../../../services/api";

import { faTimesCircle, faCheckCircle, faClock, faQuestionCircle, IconDefinition } from "@fortawesome/free-solid-svg-icons";

/**
 * ServerCardComponent is the controller for a single server card.
 */
@Component({
	selector: "server-card",
	styleUrls: ["./server-card.component.scss"],
	templateUrl: "./server-card.component.html",
})
export class ServerCardComponent {
	/** The server described by this card */
	@Input() public server: Server;

	/** The "server checks" associated with this card's server. */
	public checks: Map<string, number | boolean>;
	/** Whether or not the card's contents are showing. */
	public open: boolean;
	/** The icon used for checks with unknown values. */
	public unknownIcon = faQuestionCircle;

	constructor(private readonly api: ServerService) {
		this.open = false;
		this.checks = new Map();
	}

	/** Returns whether or not the server is considered "down" */
	public down(): boolean {
		return this.server.status === "ADMIN_DOWN" || this.server.status === "OFFLINE";
	}

	/**
	 * cacheServer returns 'true' if this component's server is a 'cache server', 'false' otherwise.
	 */
	public cacheServer(): boolean {
		return this.server.type !== undefined && (this.server.type.indexOf("EDGE") === 0 || this.server.type.indexOf("MID") === 0);
	}

	/**
	 * Returns the icon to be used to display for the server's "Updates Pending" field.
	 */
	public updPendingIcon(): IconDefinition {
		return this.server.updPending ? faClock : faCheckCircle;
	}

	/**
	 * Returns the appropriate title to describe the server's "Updates Pending" state.
	 */
	public updPendingTitle(): string {
		return this.server.updPending ? "Updates are pending" : "No updates are pending";
	}

	/**
	 * Returns the appropriate icon to use for the server's "Revalidations Pending" field.
	 */
	public revalPendingIcon(): IconDefinition {
		return this.server.revalPending ? faClock : faCheckCircle;
	}

	/**
	 * Returns the appropriate title to describe the server's "Revalidations Pending" state.
	 */
	public revalPendingTitle(): string {
		return this.server.revalPending ? "Revalidations are pending" : "No revalidations are pending";
	}

	/** Returns the appropriate icon for the given check. */
	public checkIcon(check: string): IconDefinition {
		if (!this.checks.has(check)) {
			return faQuestionCircle;
		}
		const val = this.checks.get(check);
		if (typeof val !== "boolean") {
			console.error(`Expected boolean value for server #${this.server.id} '${check}' check, got '${val}' (${typeof val})`);
			return faQuestionCircle;
		}
		return val ? faCheckCircle : faTimesCircle;
	}

	/**
	 * Returns the appropriate title for the given check.
	 */
	public checkTitle(check: string): string {
		if (!this.checks.has(check)) {
			return "Unknown";
		}
		const val = this.checks.get(check);
		if (typeof val === "boolean") {
			return val ? "Check successful" : "Check failed";
		}
		return String(val);
	}

	/**
	 * Returns the CSS class to use for the given check based on its value.
	 */
	public checkClass(check: string): string {
		if (!this.checks.has(check)) {
			return "pending";
		}
		const val = this.checks.get(check);
		if (typeof val === "boolean") {
			return val ? "" : "bad";
		}
		return "";
	}

	/** Returns the value of the given check */
	public checkValue(check: string): number | undefined | boolean {
		return this.checks.get(check);
	}

	/**
	 * Handler for when the card's open/close state is toggled.
	 */
	public toggle(e: Event): void {
		this.open = !this.open;
		if (this.open && this.cacheServer()) {
			if (!this.server.id) {
				console.error("Server has no ID - cannot load checks");
				return;
			}
			this.api.getServerChecks(this.server.id).pipe(first()).subscribe(
				(s: Servercheck) => {
					this.checks = s.checkMap();
				}
			);
		}
	}
}
