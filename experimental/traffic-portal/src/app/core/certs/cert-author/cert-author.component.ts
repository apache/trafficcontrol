import { Component, Input } from "@angular/core";

import { Author } from "src/app/core/certs/cert-detail/cert-detail.component";

/**
 * CertAuthorComponent is the controller used for displaying a cert author.
 */
@Component({
	selector: "tp-cert-author",
	styleUrls: ["./cert-author.component.scss"],
	templateUrl: "./cert-author.component.html"
})
export class CertAuthorComponent {
	@Input() public author: Author | undefined;

}
