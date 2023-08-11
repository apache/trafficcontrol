/**
 * @license Apache-2.0
 *
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

import { LoggingService } from "./logging.service";

describe("LoggingService", () => {
	let service: LoggingService;
	const arg = "test";

	beforeEach(() => {
		TestBed.configureTestingModule({});
		service = TestBed.inject(LoggingService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("logs debug messages", () => {
		const debugSpy = spyOn(console, "debug");
		expect(debugSpy).not.toHaveBeenCalled();
		service.debug(arg);
		expect(debugSpy).toHaveBeenCalledTimes(1);
	});

	it("logs error messages", () => {
		const errorSpy = spyOn(console, "error");
		expect(errorSpy).not.toHaveBeenCalled();
		service.error(arg);
		expect(errorSpy).toHaveBeenCalledTimes(1);
	});

	it("logs info messages", () => {
		const infoSpy = spyOn(console, "info");
		expect(infoSpy).not.toHaveBeenCalled();
		service.info(arg);
		expect(infoSpy).toHaveBeenCalledTimes(1);
	});

	it("logs warning messages", () => {
		const warnSpy = spyOn(console, "warn");
		expect(warnSpy).not.toHaveBeenCalled();
		service.warn(arg);
		expect(warnSpy).toHaveBeenCalledTimes(1);
	});
});
