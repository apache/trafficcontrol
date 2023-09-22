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

import { TestBed } from "@angular/core/testing";

import {Theme, ThemeManagerService} from "./theme-manager.service";

describe("ThemeManagerService", () => {
	let service: ThemeManagerService;
	const localStorage: Record<string, string> = {};

	beforeEach(() => {
		TestBed.configureTestingModule({});
		service = TestBed.inject(ThemeManagerService);

		spyOn(window.localStorage, "getItem").and.callFake((key) => key in localStorage ? localStorage[key] : null);
		spyOn(window.localStorage, "setItem").and.callFake((key, value) => localStorage[key] = value);
		// eslint-disable-next-line @typescript-eslint/no-dynamic-delete
		spyOn(window.localStorage, "removeItem").and.callFake((key) => delete localStorage[key]);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("init theme manager", () => {
		const theme = {fileName: "dark-default-theme.css", name: "Dark"} as Theme;
		service.themeChanged.subscribe((newTheme: Theme): void => {
			expect(newTheme.fileName).toBe(theme.fileName);
			expect(newTheme.name).toBe(theme.name);
		});

		service.initTheme();
		window.localStorage["current-theme-name"] = JSON.stringify(theme);
		service.initTheme();
	});

	it("set theme", () => {
		const theme = {fileName: "dark-default-theme.css", name: "Dark"} as Theme;
		const sub = service.themeChanged.subscribe((newTheme: Theme): void => {
			expect(newTheme.fileName).toBe(theme.fileName);
			expect(newTheme.name).toBe(theme.name);
		});

		service.loadTheme(theme);
		expect(theme).toEqual(JSON.parse(localStorage["current-theme-name"] ?? ""));

		sub.unsubscribe();

		theme.fileName = undefined;
		service.loadTheme(theme);
		const storedTheme = localStorage["current-theme-name"];
		expect(storedTheme).toBeUndefined();
	});

	it("clear theme", () => {
		const theme = {fileName: "dark-default-theme.css", name: "Dark"} as Theme;

		service.loadTheme(theme);
		expect(theme).toEqual(JSON.parse(localStorage["current-theme-name"] ?? ""));

		service.clearTheme();
		const storedTheme = localStorage["current-theme-name"];
		expect(storedTheme).toBeUndefined();
	});
});
