/**
 * @license Apache-2.0
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
import { DOCUMENT } from "@angular/common";
import { TestBed } from "@angular/core/testing";

import { FileUtilsService } from "./file-utils.service";

describe("FileUtilsService", () => {
	let service: FileUtilsService;
	let spy: jasmine.Spy<typeof globalThis.open>;
	let urlSpy: jasmine.Spy<typeof URL.createObjectURL>;

	beforeEach(() => {
		const window = {open: (): void => {
			// do nothing
		}};
		spy = spyOn(window, "open").and.callThrough();
		urlSpy = spyOn(URL, "createObjectURL").and.callThrough();
		TestBed.configureTestingModule({
			providers: [
				FileUtilsService,
				{provide: DOCUMENT, useValue: {defaultView: window}}
			]
		});
		expect(spy).not.toHaveBeenCalled();
		service = TestBed.inject(FileUtilsService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("downloading string-data files", () => {
		const data = "data";
		it("downloads files when everything is specified", () => {
			const f = new File([data], "myfilename", {type: "mytype"});
			service.download(data, f.name, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only filename is specified", () => {
			const f = new File([data], "myfilename", {type: FileUtilsService.TEXT_CONTENT_TYPE});
			service.download(data, f.name);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only content type is specified", () => {
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.txt`, {type: "mytype"});
			service.download(data, undefined, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files without any details", () => {
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.txt`, {type: FileUtilsService.TEXT_CONTENT_TYPE});
			service.download(data);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
	});

	describe("downloading ArrayBuffer-data files", () => {
		const data = new ArrayBuffer(27);
		it("downloads files when everything is specified", () => {
			const f = new File([data], "myfilename", {type: "mytype"});
			service.download(data, f.name, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only filename is specified", () => {
			const f = new File([data], "myfilename", {type: FileUtilsService.BINARY_DATA_CONTENT_TYPE});
			service.download(data, f.name);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only content type is specified", () => {
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.bin`, {type: "mytype"});
			service.download(data, undefined, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files without any details", () => {
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.bin`, {type: FileUtilsService.BINARY_DATA_CONTENT_TYPE});
			service.download(data);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
	});

	describe("downloading ArrayBufferView-data files", () => {
		const data = new Uint8Array();
		it("downloads files when everything is specified", () => {
			const f = new File([data], "myfilename", {type: "mytype"});
			service.download(data, f.name, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only filename is specified", () => {
			const f = new File([data], "myfilename", {type: FileUtilsService.BINARY_DATA_CONTENT_TYPE});
			service.download(data, f.name);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only content type is specified", () => {
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.bin`, {type: "mytype"});
			service.download(data, undefined, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files without any details", () => {
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.bin`, {type: FileUtilsService.BINARY_DATA_CONTENT_TYPE});
			service.download(data);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
	});

	describe("downloading arbitrary object-data files", () => {
		const data = {data: "object"};
		it("downloads files when everything is specified", () => {
			const f = new File([JSON.stringify(data)], "myfilename.bin", {type: "mytype"});
			service.download(data, f.name, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only filename is specified", () => {
			const f = new File([JSON.stringify(data)], "myfilename", {type: FileUtilsService.JSON_DATA_CONTENT_TYPE});
			service.download(data, f.name);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only content type is specified", () => {
			const f = new File([JSON.stringify(data)], `${FileUtilsService.DEFAULT_FILE_NAME}.json`, {type: "mytype"});
			service.download(data, undefined, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files without any details", () => {
			const f = new File(
				[JSON.stringify(data)],
				`${FileUtilsService.DEFAULT_FILE_NAME}.json`,
				{type: FileUtilsService.JSON_DATA_CONTENT_TYPE}
			);
			service.download(data);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
	});

	describe("downloading blobs as files", () => {
		it("downloads files when everything is specified", () => {
			const data = new Blob(["data"], {type: "myType"});
			const f = new File([data], "myfilename.bin", {type: data.type});
			service.download(data, f.name, f.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only filename is specified", () => {
			const data = new Blob(["data"]);
			const f = new File([JSON.stringify(data)], "myfilename.bin");
			service.download(data, f.name);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only content type is specified", () => {
			const data = new Blob(["data"], {type: "myType"});
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.bin`, {type: "mytype"});
			service.download(data, undefined, data.type);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files without any details", () => {
			const data = new Blob(["data"]);
			const f = new File([data], `${FileUtilsService.DEFAULT_FILE_NAME}.bin`);
			service.download(data);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
	});

	describe("downloading pre-constructed files", () => {
		it("downloads files when everything is specified", () => {
			const f = new File(["data"], "myfile", {type: "mytype"});
			service.download(f);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only filename is specified", () => {
			const f = new File(["data"], "myfile");
			service.download(f);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files when only content type is specified", () => {
			const f = new File(["data"], "", {type: "mytype"});
			service.download(f);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
		it("downloads files without any details", () => {
			const f = new File(["data"], "");
			service.download(f);
			expect(spy).toHaveBeenCalledTimes(1);
			expect(urlSpy).toHaveBeenCalledOnceWith(f);
		});
	});

});
