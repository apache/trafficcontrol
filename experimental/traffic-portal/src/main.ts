/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import { enableProdMode } from "@angular/core";
import { platformBrowserDynamic } from "@angular/platform-browser-dynamic";

import { AppModule } from "./app/app.module";
import { environment } from "./environments/environment";

if (environment.production) {
	enableProdMode();
}

document.addEventListener("DOMContentLoaded", () => {
	platformBrowserDynamic().bootstrapModule(AppModule)
		// Bootstrap failures will not be combined with logging service
		// messages, because in that case no logging service could have been
		// initialized. Therefore, consistency is unbroken, and for ease of
		// debugging it's probably best not to mess with the format of Angular
		// framework errors anyhow.
		// eslint-disable-next-line no-console
		.catch(err => console.error(err));
});
