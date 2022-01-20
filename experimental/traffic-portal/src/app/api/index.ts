import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { ServerService } from "./server.service";

export * from "./server.service";

/**
 * The API Module contains all logic used to access the Traffic Ops API.
 */
@NgModule({
	declarations: [],
	imports: [
		CommonModule
	],
	providers: [ServerService]
})
export class APIModule { }
