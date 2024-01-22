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

describe("Login Page", () => {
	beforeEach(() => {
		cy.visit("/login");
	});
	it("Clears the form when the 'Clear' button is clicked", () => {
		const usernameInput = cy.get("input").first();
		usernameInput.type("test");

		const passwordInput = cy.get("input").last();
		passwordInput.type("asdf");

		cy.contains("button", "Clear").click();

		usernameInput.should("have.value", "");
		passwordInput.should("have.value", "");
	});

	it("Rejects incorrect passwords", () => {
		cy.get("input").first().type("test");
		cy.get("input").last().type("asdf");

		cy.contains("button", "Login").click();

		cy.contains("Invalid username or password");
	});

	it("Logs in", () => {
		cy.login();
	});
});
