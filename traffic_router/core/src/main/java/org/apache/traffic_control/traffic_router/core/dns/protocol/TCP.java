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
import java.io.IOException;
import java.io.InputStream;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.nio.channels.Channels;

import org.apache.log4j.Logger;
import org.xbill.DNS.Message;
import org.xbill.DNS.WireParseException;

public class TCP extends AbstractProtocol {
    private static final Logger LOGGER = Logger.getLogger(TCP.class);
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
            try {
                final Socket socket = getServerSocket().accept();
                final TCPSocketHandler handler = new TCPSocketHandler(socket);
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
            if (cancel) {
                cleanup();
                return;
            }

            try {
                socket.setSoTimeout(getReadTimeout());
                final InetAddress client = socket.getInetAddress();

                final InputStream iis = Channels.newInputStream(Channels.newChannel(socket.getInputStream()));
                final DataInputStream is = new DataInputStream(iis);
                final DataOutputStream os = new DataOutputStream(socket.getOutputStream());
                final int length = is.readUnsignedShort();
                final byte[] request = new byte[length];
                is.readFully(request);

                final byte[] response = query(client, request);
                os.writeShort(response.length);
                os.write(response);
            } catch (final WireParseException e) {
                // This is already recorded in the access log
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
