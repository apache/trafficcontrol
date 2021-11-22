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

package org.apache.traffic_control.traffic_router.core.http;

import org.springframework.web.util.ContentCachingResponseWrapper;

import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class BufferedResponse extends ContentCachingResponseWrapper {
	protected int contentLength;
	protected ServletResponse response;

	public BufferedResponse(final HttpServletResponse response) {
		super(response);
		this.response = response;
	}

	@Override
	public void setContentLength(final int len) {
		contentLength = len;
		super.setContentLength(len);
	}

	@Override
	public void setContentLengthLong(final long len) {
		contentLength = (int) len;
		super.setContentLengthLong(len);
	}

	@Override
	public void copyBodyToResponse() throws IOException {
		if (this.getContentSize() == 0) {
			response.setContentLength(contentLength);
		} else {
			// When the content size is greater than 0, copyBodyToResponse()
			// will set Content-Length.
			super.copyBodyToResponse();
		}
	}
}
