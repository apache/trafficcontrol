import { NgModule } from "@angular/core";
import { CommonModule } from "@angular/common";
import { ServerService } from "../server.service";
import { ServerService as TestingServerService } from "./server.service";

export * from "./server.service";

/**
 * The API Testing Module provides mock services that allow components to use
 * the Traffic Ops API without actually requiring a running Traffic Ops.
 */
@NgModule({
	declarations: [],
	imports: [
		CommonModule
	],
	providers: [{provide: ServerService, useClass: TestingServerService}]
})
export class APITestingModule { }
