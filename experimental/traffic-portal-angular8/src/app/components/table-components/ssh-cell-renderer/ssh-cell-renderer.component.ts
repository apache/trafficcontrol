import { Component, OnInit, OnDestroy } from "@angular/core";

import { Subscription } from "rxjs";

import { AuthenticationService } from "../../../services";

@Component({
	selector: "ssh-cell-renderer",
	styleUrls: ["./ssh-cell-renderer.component.scss"],
	templateUrl: "./ssh-cell-renderer.component.html"
})
export class SSHCellRendererComponent implements OnInit, OnDestroy {

	private user: Subscription;
	private username = "";

	/** The IP address or hostname to which the SSH link will point. */
	public value = "";

	/** The SSH URL to use. */
	public get href(): string {
		return `ssh://${this.username}@${this.value}`;
	}

	constructor(private readonly auth: AuthenticationService) { }

	public init(params: {value: string}): void {
		this.value = params.value;
	}

	public agInit(params) {
		console.info(params);
	}

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

	public ngOnDestroy(): void {
		this.user.unsubscribe();
	}

}
