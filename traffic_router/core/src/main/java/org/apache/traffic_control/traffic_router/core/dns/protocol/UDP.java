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

package org.apache.traffic_control.traffic_router.core.dns.protocol;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.xbill.DNS.Message;
import org.xbill.DNS.OPTRecord;
import org.xbill.DNS.WireParseException;

public class UDP extends AbstractProtocol {
    private static final Logger LOGGER = LogManager.getLogger(UDP.class);

    private static final int UDP_MSG_LENGTH = 512;

    private DatagramSocket datagramSocket;

    /**
     * Gets datagramSocket.
     * 
     * @return the datagramSocket
     */
    public DatagramSocket getDatagramSocket() {
        return datagramSocket;
    }

    @Override
    public void run() {
        while (!isShutdownRequested()) {
            try {
                final byte[] buffer = new byte[UDP_MSG_LENGTH];
                final DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                datagramSocket.receive(packet);
                submit(new UDPPacketHandler(packet));
            } catch (final IOException e) {
				LOGGER.warn("error: " + e);
            }
        }
    }
    @Override
    public void shutdown() {
    	super.shutdown();
    	datagramSocket.close();
    }

    /**
     * Sets datagramSocket.
     * 
     * @param datagramSocket
     *            the datagramSocket to set
     */
    public void setDatagramSocket(final DatagramSocket datagramSocket) {
        this.datagramSocket = datagramSocket;
    }

    @Override
    protected int getMaxResponseLength(final Message request) {
        int result = UDP_MSG_LENGTH;
        if ((request != null) && (request.getOPT() != null)) {
            final OPTRecord opt = request.getOPT();
            result = opt.getPayloadSize();
        }
        return result;
    }

    /**
     * This class is package private for unit testing purposes.
     */
    class UDPPacketHandler implements SocketHandler {
        private final DatagramPacket packet;
        private boolean cancel;

        /**
         * This method is package private for unit testing purposes.
         * 
         * @param packet
         */
        UDPPacketHandler(final DatagramPacket packet) {
            this.packet = packet;
        }

        @Override
        @SuppressWarnings("PMD.EmptyCatchBlock")
        public void run() {
            if (cancel) {
                cleanup();
                return;
            }

            try {
                final InetAddress client = packet.getAddress();
                final byte[] request = new byte[packet.getLength()];
                System.arraycopy(packet.getData(), 0, request, 0, request.length);
                final byte[] response = query(client, request);

                final DatagramPacket outPacket = new DatagramPacket(response, response.length,
                        packet.getSocketAddress());
                getDatagramSocket().send(outPacket);
            } catch (final WireParseException e) {
                // This is already recorded in the access log
            } catch (final Exception e) {
                LOGGER.error(e.getMessage(), e);
            }
        }

        @Override
        public void cleanup() {
            // noop for UDP
        }

        @Override
        public void cancel() {
            this.cancel = true;
        }
    }

}
