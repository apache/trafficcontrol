import { Component, OnInit, OnDestroy } from "@angular/core";
import { DomSanitizer, SafeUrl } from "@angular/platform-browser";

import { Subscription } from "rxjs";

import { AuthenticationService } from "../../../services";

interface InitParams {
	/** the cell's value */
	value: string;
}

/**
 * SSHCellRendererComponent is an AG-Grid cell renderer that provides ssh:// links as content.
 */
@Component({
	selector: "ssh-cell-renderer",
	styleUrls: ["./ssh-cell-renderer.component.scss"],
	templateUrl: "./ssh-cell-renderer.component.html"
})
export class SSHCellRendererComponent implements OnInit, OnDestroy {

	private user: Subscription;
	private username = "";

	/** The IP address or hostname to which the SSH link will point. */
	public get value(): string {
		return this.val;
	}
	private val = "";

	/** The SSH URL to use. */
	public get href(): SafeUrl {
		const url = `ssh://${this.username}@${this.value}`;
		return this.sanitizer.bypassSecurityTrustUrl(url);
	}

	constructor(private readonly auth: AuthenticationService, private readonly sanitizer: DomSanitizer) { }

	/** Called by the AG-Grid API at initalization */
	public init(params: InitParams): void {
		this.val = params.value;
	}

	/** called after ag-grid is initalized */
	public agInit(params: unknown): void {
		console.log(params);
	}

	/** Called when the Angular view is initialized, sets up the username for the SSH links */
	public ngOnInit(): void {
		this.user = this.auth.currentUser.subscribe(
			u => {
				if (!u) {
					return;
				}
				this.username = u.username || "";
			}
		);
	}

	/** Cleans up resources when the component is no longer rendered. */
	public ngOnDestroy(): void {
		this.user.unsubscribe();
	}

}
