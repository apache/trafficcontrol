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

import static org.junit.Assert.assertEquals;

import org.junit.Test;
import org.xbill.DNS.Rcode;

public class DNSExceptionTest {

    @Test
    public void testDNSExceptionInt() {
        final int rcode = Rcode.NXDOMAIN;
        final DNSException ex = new DNSException(rcode);
        assertEquals(rcode, ex.getRcode());
    }

    @Test
    public void testDNSExceptionIntString() {
        final int rcode = Rcode.NXDOMAIN;
        final String msg = "message";
        final DNSException ex = new DNSException(rcode, msg);
        assertEquals(rcode, ex.getRcode());
        assertEquals(msg, ex.getMessage());
    }

    @Test
    public void testDNSExceptionIntStringThrowable() {
        final int rcode = Rcode.NXDOMAIN;
        final String msg = "message";
        final Exception cause = new Exception();
        final DNSException ex = new DNSException(rcode, msg, cause);
        assertEquals(rcode, ex.getRcode());
        assertEquals(msg, ex.getMessage());
        assertEquals(cause, ex.getCause());
    }

    @Test
    public void testDNSExceptionIntThrowable() {
        final int rcode = Rcode.NXDOMAIN;
        final Exception cause = new Exception();
        final DNSException ex = new DNSException(rcode, cause);
        assertEquals(rcode, ex.getRcode());
        assertEquals(cause, ex.getCause());
    }

}
