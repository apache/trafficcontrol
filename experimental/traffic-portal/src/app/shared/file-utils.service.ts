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
import { Inject, Injectable } from "@angular/core";

import { isArrayBufferView } from "../utils";

/**
 * FileUtilsService provides utilities dealing with files, uploads, downloads,
 * etc.
 */
@Injectable()
export class FileUtilsService {

	/** The default MIME-Type for string data downloads. */
	public static readonly TEXT_CONTENT_TYPE = "text/plain+x-traffic-ops;charset=UTF-8";
	/** The default MIME-Type for raw binary data downloads. */
	public static readonly BINARY_DATA_CONTENT_TYPE = "application/octet-stream";
	/** The default MIME-Type for arbitrary object downloads. */
	public static readonly JSON_DATA_CONTENT_TYPE = "application/json+x-traffic-ops";

	/**
	 * The file name that will be used for downloads, if one is not provided.
	 */
	public static readonly DEFAULT_FILE_NAME = "download";

	/** A pre-compiled expression that matches a '.' at the start of a string */
	private static readonly EXT_PATTERN = /^\./;

	private readonly window: Window;

	constructor(@Inject(DOCUMENT) document: Document) {
		const {defaultView} = document;
		if (!defaultView) {
			throw new Error("global root document has no default view; cannot access required functionality");
		}
		this.window = defaultView;
	}

	/**
	 * Builds a file name for {@link FileUtilsService.download}.
	 *
	 * @param passed The file name as passed to the service.
	 * @param defaultExt The default file extension, to be used if the passed
	 * file name doesn't have one - or if no file name was passed.
	 * @returns The constructed file name.
	 */
	private static constructFileName(passed: string | undefined, defaultExt: string): string {
		const basename = passed || this.DEFAULT_FILE_NAME;
		if (passed && basename.includes(".")) {
			return basename;
		}
		return `${basename}.${defaultExt.replace(this.EXT_PATTERN, "")}`;
	}

	/**
	 * Initiates a download of a file created from the passed data.
	 *
	 * @param data The data to be contained in the downloaded file. If this is
	 * a string or a buffer of binary data, it's treated as the literal contents
	 * of the file. If it's some arbitrary, unrecognized object, it'll be
	 * encoded using `JSON.stringify`.
	 * @param fileName The suggested name of the file in the browser's "Save"
	 * dialog. If this isn't provided (or is an empty string),
	 * {@link FileUtilsService.DEFAULT_FILE_NAME} will be used.
	 * @param contentType The MIME content type of the data stored in the file,
	 * which is used by some clients to determine how to handle the data (for
	 * example most browsers display PDFs in-browser instead of initiating a
	 * download). If this isn't provided,
	 * {@link FileUtilsService.BINARY_DATA_CONTENT_TYPE} will be assumed for
	 * ArrayBuffer and ArrayBufferView data,
	 * {@link FileUtilsService.TEXT_CONTENT_TYPE} will be assumed for string
	 * data, and {@link FileUtilsService.JSON_DATA_CONTENT_TYPE} will be used
	 * when the passed data is encoded as JSON.
	 */
	public download(data: ArrayBuffer | ArrayBufferView | string | object, fileName?: string, contentType?: string): void;
	/**
	 * Initiates a download of a file created from the passed content.
	 *
	 * @param content The data to be contained in the downloaded file. If the
	 * blob has a `type`, it will be used as the file's MIME content type -
	 * otherwise, {@link FileUtilsService.BINARY_DATA_CONTENT_TYPE} will be
	 * used.
	 * @param fileName The suggested name of the file in the browser's "Save"
	 * dialog. If this isn't provided (or is an empty string),
	 * {@link FileUtilsService.DEFAULT_FILE_NAME} will be used.
	 */
	public download(content: Blob, fileName?: string): void;
	/**
	 * Initiates a download of a file.
	 *
	 * @param file The file to be downloaded.
	 */
	public download(file: File): void;
	/**
	 * Initiates a download of a file created from the passed data, or passed
	 * directly as a pre-constructed content blob or file object.
	 *
	 * @param file The data to be contained in the downloaded file. If this is
	 * a string or a buffer of binary data, it's treated as the literal contents
	 * of the file. If it's some arbitrary, unrecognized object, it'll be
	 * encoded using `JSON.stringify`. If it's a blob, then it is treated as the
	 * data to be contained in the downloaded file. If the blob has a `type`, it
	 * will be used as the file's MIME content type - otherwise,
	 * {@link FileUtilsService.BINARY_DATA_CONTENT_TYPE} will be used (the
	 * `contentType` parameter of this service method is ignored in that case).
	 * Lastly, if this is a pre-constructed `File`, then it is assumed to be
	 * fully ready for download; `fileName` and `contentType` are both ignored.
	 * @param fileName The suggested name of the file in the browser's "Save"
	 * dialog. If this isn't provided (or is an empty string),
	 * {@link FileUtilsService.DEFAULT_FILE_NAME} will be used - unless `file`
	 * was a `File`, in which case it is totally ignored whether passed or not.
	 * @param contentType The MIME content type of the data stored in the file,
	 * which is used by some clients to determine how to handle the data (for
	 * example most browsers display PDFs in-browser instead of initiating a
	 * download). This is ignored if the data passed in `file` is a `Blob` or
	 * `File`, and their respective `type` is used instead (if a `Blob` is
	 * missing a type, then it will be set to
	 * {@link FileUtilsService.BINARY_DATA_CONTENT_TYPE}).
	 */
	public download(file: File | Blob | ArrayBuffer | ArrayBufferView | string | object, fileName?: string, contentType?: string): void {
		let f: File;
		if (typeof(file) === "string") {
			const fname = FileUtilsService.constructFileName(fileName, ".txt");
			f = new File([file], fname, {type: contentType || FileUtilsService.TEXT_CONTENT_TYPE});
		} else if (file instanceof(ArrayBuffer) || isArrayBufferView(file)) {
			const fname = FileUtilsService.constructFileName(fileName, ".bin");
			f = new File([file], fname, {type: contentType || FileUtilsService.BINARY_DATA_CONTENT_TYPE});
		} else if (file instanceof File) {
			f = file;
		} else if (file instanceof Blob) {
			const fname = FileUtilsService.constructFileName(fileName, ".bin");
			f = new File([file], fname, {type: file.type || FileUtilsService.BINARY_DATA_CONTENT_TYPE});
		} else {
			const content = JSON.stringify(file);
			const fname = FileUtilsService.constructFileName(fileName, ".json");
			f = new File([content], fname, {type: contentType || FileUtilsService.JSON_DATA_CONTENT_TYPE});
		}

		const url = URL.createObjectURL(f);
		this.window.open(url, "_blank");
		URL.revokeObjectURL(url);
	}
}
