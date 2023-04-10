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
import {
	EnhancedElementInstance,
	EnhancedPageObject, EnhancedSectionInstance, NightwatchAPI
} from "nightwatch";

/**
 * Defines the commands for the loginForm section
 */
interface LoginFormSectionCommands extends EnhancedSectionInstance, EnhancedElementInstance<EnhancedPageObject> {
	fillOut(username: string, password: string): this;
	login(username: string, password: string): this;
	loginAndWait(username: string, password: string): this;
}

/**
 * Defines the loginForm section
 */
type LoginFormSection = EnhancedSectionInstance<LoginFormSectionCommands, typeof loginPageObject.sections.loginForm.elements>;

/**
 * Define the type for our PO
 */
export type LoginPageObject = EnhancedPageObject<{}, {}, { loginForm: LoginFormSection }>;

const loginPageObject = {
	api: {} as NightwatchAPI,
	sections: {
		loginForm: {
			commands: {
				fillOut(username: string, password: string): LoginFormSectionCommands {
					return this
						.setValue("@usernameTxt", username)
						.setValue("@passwordTxt", password);
				},
				login(username: string, password: string): LoginFormSectionCommands {
					return this.fillOut(username, password)
						.click("@loginBtn");
				},
				loginAndWait(username: string, password: string): LoginFormSectionCommands {
					const ret = this.login(username, password);
					browser.page.common()
						.assert.textContains("@snackbarEle", "Success");
					return ret;
				}
			} as LoginFormSectionCommands,
			elements: {
				clearBtn: {
					selector: "button[name='clear']"
				},
				loginBtn: {
					selector: "button[name='login']"
				},
				passwordTxt: {
					selector: "input[name='p']"
				},
				resetBtn: {
					selector: "button[name='reset']"
				},
				usernameTxt: {
					selector: "input[name='u']"
				}
			},
			selector: "form[name='login']"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/login`;
	}
};

export default loginPageObject;
