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

import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.EOFException;
import java.io.IOException;
import java.io.InputStream;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.SocketTimeoutException;
import java.nio.channels.Channels;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.xbill.DNS.Message;
import org.xbill.DNS.WireParseException;

public class TCP extends AbstractProtocol {
    private static final Logger LOGGER = LogManager.getLogger(TCP.class);
    private int readTimeout = 3000; // default

    private ServerSocket serverSocket;

    /**
     * Gets serverSocket.
     * 
     * @return the serverSocket
     */
    public ServerSocket getServerSocket() {
        return serverSocket;
    }

    @Override
    public void run() {
        while (!isShutdownRequested()) {
            final TCPSocketHandler handler;
            try {
                handler = new TCPSocketHandler(getServerSocket().accept());
                submit(handler);
            } catch (final IOException e) {
				LOGGER.warn("error: " + e);
            }
        }
    }

    /**
     * Sets serverSocket.
     * 
     * @param serverSocket
     *            the serverSocket to set
     */
    public void setServerSocket(final ServerSocket serverSocket) {
        this.serverSocket = serverSocket;
    }

    @Override
    protected int getMaxResponseLength(final Message request) {
        return Integer.MAX_VALUE;
    }

    /**
     * This class is package private for unit testing purposes.
     */
    class TCPSocketHandler implements SocketHandler {
        private final Socket socket;
        private boolean cancel;

        /**
         * This method is package private for unit testing purposes.
         * 
         * @param socket
         */
        TCPSocketHandler(final Socket socket) {
            this.socket = socket;
        }

        @Override
        @SuppressWarnings("PMD.EmptyCatchBlock")
        public void run() {
            InetAddress client = null;
            if (cancel) {
                cleanup();
                return;
            }

            try (InputStream iis = Channels.newInputStream(Channels.newChannel(socket.getInputStream()));
                 DataInputStream is = new DataInputStream(iis);
                 DataOutputStream os = new DataOutputStream(socket.getOutputStream())
            ) {
                socket.setSoTimeout(getReadTimeout());
                client = socket.getInetAddress();

                final int length = is.readUnsignedShort();
                final byte[] request = new byte[length];
                is.readFully(request);

                final byte[] response = query(client, request);
                os.writeShort(response.length);
                os.write(response);
            } catch (final WireParseException e) {
                // This is already recorded in the access log
            } catch (final SocketTimeoutException e) {
                String hostAddress = "unknown";
                if (client != null) {
                    hostAddress = client.getHostAddress();
                }
                LOGGER.error("The socket with the Client at: " +
                        hostAddress + " has timed out. Error: " + e.getMessage());
            } catch (final EOFException e) {
                String hostAddress = "unavailable";
                if (client != null) {
                    hostAddress = client.getHostAddress();
                }
                LOGGER.error("The client at " + hostAddress +
                        " has closed the connection prematurely. Error: " + e.getMessage());
            } catch (final Exception e) {
                LOGGER.error(e.getMessage(), e);
            } finally {
                cleanup();
            }
        }

        @Override
        public void cleanup() {
            if (socket == null) {
                return;
            }

            try {
                socket.close();
            } catch (final IOException e) {
                LOGGER.debug(e.getMessage(), e);
            }
        }

        @Override
        public void cancel() {
            this.cancel = true;
        }
    }

    @Override
    public void shutdown() {
    	super.shutdown();
    	try {
			serverSocket.close();
		} catch (IOException e) {
			LOGGER.warn("error on shutdown", e);
		}
	}


	public int getReadTimeout() {
		return readTimeout;
	}

	public void setReadTimeout(final int readTimeout) {
		this.readTimeout = readTimeout;
	}


}
