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

//
// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --

/**
 * Logs in as the configured testing user. This ends with the URL at /core.
 */
function login(): void {
	cy.visit("login");
	cy.fixture("login").then(
		({username, password}: {username: string; password: string}) => {
			cy.get("input").first().type(username);
			cy.get("input").last().type(password);
			cy.contains("button", "Login").click();
			cy.window().its("location.pathname").should("eq", "/core");
		}
	);
}

// This is how Cypress is organized; we don't have control over it.
// eslint-disable-next-line @typescript-eslint/no-namespace
declare namespace Cypress {
	// All declarations in the namespace must match in type, and Cypress itself
	// uses `any`, so we are powerless to do something better here.
	/* eslint-disable @typescript-eslint/no-explicit-any */
	// This type parameter is necessary because it makes the interface match the
	// expanded definition provided by Cypress itself.
	/* eslint-disable @typescript-eslint/no-unused-vars */
	// We shouldn't document Cypress interfaces.
	// eslint-disable-next-line jsdoc/require-jsdoc
	interface Chainable<Subject = any> {
		login(): typeof login;
	}
	/* eslint-enable @typescript-eslint/no-explicit-any */
	/* eslint-enable @typescript-eslint/no-unused-vars */
}
Cypress.Commands.add("login", login);
