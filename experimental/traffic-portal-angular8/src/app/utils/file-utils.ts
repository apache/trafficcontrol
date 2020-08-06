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

/**
 * Causes the browser to initiate a download
 *
 * @param content Should be a File, data Blob, raw data string, or some arbitrary object. Arbitrary objects will be
 * converted to JSON via JSON.stringify.
 * @param filename Will be the name of the file to download. If not given and the {@link content} is not a File,
 * will default to 'download'.
 * @param type The MIME type of the download - used as a hint for file extensions. If not given and {@link content}
 * is not a File or Blob, this will default to 'text/plain' for strings, and 'application/json' for all others.
 */
export function download (content: Blob | File | string | any, filename?: string, type?: string): void {
	if (content instanceof File) {
		const url = URL.createObjectURL(content);
		window.location.assign(url);
		URL.revokeObjectURL(url);
		return;
	}

	let f: File;
	if (!filename) {
		filename = "download";
	}

	if (content instanceof Blob) {
		f = new File([content], filename);
	} else if (typeof(content) === "string") {
		if (!type) {
			type = "text/plain";
		}
		f = new File([content], filename, {type: type});
	} else {
		if (!type) {
			type = "application/json";
		}
		f = new File([JSON.stringify(content)], filename, {type: type});
	}

	const exportURL = URL.createObjectURL(f);
	window.location.assign(exportURL);
	URL.revokeObjectURL(exportURL);
}
