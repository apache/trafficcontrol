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
	ElementCommands,
	EnhancedElementInstance,
	EnhancedPageObject, EnhancedSectionInstance
} from "nightwatch";

interface LoginFormSectionCommands extends EnhancedSectionInstance, ElementCommands {
	fillOut(username: string, password: string): this;
	login(username: string, password: string): this;
}

type LFSection = EnhancedSectionInstance<LoginFormSectionCommands, typeof loginPageObject.sections.loginForm.elements>;

/**
 * Define the type for our PO
 */
//export type LoginPageObject = EnhancedPageObject<{}, {}, LoginSectionInstance>;

/**
 * Define the type for our PO
 */
export type LoginPageObject = EnhancedPageObject<{}, typeof loginPageObject.elements, { loginForm: LFSection }>;
const loginPageObject = {
	elements: {
		snackbarEle: {
			selector: "simple-snack-bar"
		}
	},
	sections: {
		loginForm: {
			commands: {
				fillOut(username: string, password: string): LoginFormSectionCommands  {
					 return this
						.setValue("input#u", username)
						.setValue("@passwordTxt", password);
				},
				login(username: string, password: string): LoginFormSectionCommands {
					 return this.fillOut(username, password)
						.click("@loginBtn");
				},
			} as LoginFormSectionCommands,
			elements: {
				clearBtn: {
					selector: "button[name='clear']"
				},
				loginBtn: {
					selector: "button[name='login']"
				},
				passwordTxt: {
					selector: "input#p"
				},
				resetBtn: {
					selector: "button[name='reset']"
				},
				usernameTxt: {
					selector: "input[id='u']"
				}
			},
			selector: "form[name='login']"
		}
	},
	url(): string {
		return `${this.api.launchUrl}/login`;
	}
} as LoginPageObject;

export default loginPageObject;
