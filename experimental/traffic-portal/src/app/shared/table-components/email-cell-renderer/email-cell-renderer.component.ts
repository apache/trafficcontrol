import { Component } from "@angular/core";
import { DomSanitizer, type SafeUrl } from "@angular/platform-browser";
import type { ICellRendererParams } from "ag-grid-community";

/**
 * A table cell renderer for cells containing email addresses - it link-ifies
 * the email into a `mailto` link.
 */
@Component({
	selector: "tp-email-cell-renderer",
	styleUrls: ["./email-cell-renderer.component.scss"],
	templateUrl: "./email-cell-renderer.component.html"
})
export class EmailCellRendererComponent {

	/** The email address which the mailto link will point. */
	public get value(): string {
		return this.val;
	}

	/** The raw value for the email address. */
	private val = "";

	/** The mailto URL to use. */
	public get href(): SafeUrl {
		const url = `mailto:${this.val}`;
		return this.sanitizer.bypassSecurityTrustUrl(url);
	}

	constructor(private readonly sanitizer: DomSanitizer) {
	}

	/**
	 * Called when the value changes - I don't think this will ever happen.
	 *
	 * @param params The AG-Grid cell-rendering parameters (refer to AG-Grid docs).
	 * @returns 'true' if the component could be refreshed (always for this simple component) or 'false' if the component must be recreated.
	 */
	public refresh(params: ICellRendererParams): true {
		this.val = params.value;
		return true;
	}

	/**
	 * Called after ag-grid is initalized.
	 *
	 * @param params The AG-Grid cell-rendering parameters (refer to AG-Grid docs).
	 */
	public agInit(params: ICellRendererParams): void {
		this.val = params.value;
	}
}
