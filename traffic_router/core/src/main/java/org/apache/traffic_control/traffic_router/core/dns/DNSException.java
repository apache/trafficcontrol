/*
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

package org.apache.traffic_control.traffic_router.core.dns;

/**
 * An {@link Exception} that relates to a specified DNS RCODE.
 */
public class DNSException extends Exception {

    private static final long serialVersionUID = 1L;

    private final int rcode;

    /**
     * @param rcode
     *            the DNS RCODE associated with this exception.
     */
    public DNSException(final int rcode) {
        this.rcode = rcode;
    }

    /**
     * @param rcode
     *            the DNS RCODE associated with this exception.
     * @param message
     *            a human readable message associated with the exception
     */
    public DNSException(final int rcode, final String message) {
        super(message);
        this.rcode = rcode;
    }

    /**
     * @param rcode
     *            the DNS RCODE associated with this exception.
     * @param message
     *            a human readable message associated with the exception
     * @param cause
     *            a chained throwable that caused this exception
     */
    public DNSException(final int rcode, final String message, final Throwable cause) {
        super(message, cause);
        this.rcode = rcode;
    }

    /**
     * @param rcode
     *            the DNS RCODE associated with this exception.
     * @param cause
     *            a chained {@link Throwable} that caused this exception
     */
    public DNSException(final int rcode, final Throwable cause) {
        super(cause);
        this.rcode = rcode;
    }

    /**
     * Gets rcode.
     * 
     * @return the rcode
     */
    public int getRcode() {
        return rcode;
    }

}
