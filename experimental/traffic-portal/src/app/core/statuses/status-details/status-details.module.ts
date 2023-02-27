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
import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { ReactiveFormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatFormFieldModule } from "@angular/material/form-field";
import {MatGridListModule} from "@angular/material/grid-list";
import { MatInputModule } from "@angular/material/input";
import { RouterModule } from "@angular/router";

import { StatusesService } from "src/app/api/statuses.service";
import { SharedModule } from "src/app/shared/shared.module";

import { StatusDetailsComponent } from "./status-details.component";

const StatusDetailRouting = RouterModule.forChild([
	{
		path: "",
		component: StatusDetailsComponent
	}
]);

@NgModule({
	declarations: [
		StatusDetailsComponent
	],
	imports: [
		CommonModule,
		StatusDetailRouting,
		ReactiveFormsModule,
		MatFormFieldModule,
		MatInputModule,
		MatGridListModule,
		MatCardModule,
		MatButtonModule,
		SharedModule
	],
	providers:[
		StatusesService
	]
})
export class StatusDetailsModule { }
