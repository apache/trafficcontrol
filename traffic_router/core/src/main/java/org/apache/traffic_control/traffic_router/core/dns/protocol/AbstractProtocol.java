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

import org.apache.traffic_control.traffic_router.core.dns.DNSAccessEventBuilder;
import org.apache.traffic_control.traffic_router.core.dns.DNSAccessRecord;
import org.apache.traffic_control.traffic_router.core.dns.NameServer;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.xbill.DNS.Message;
import org.xbill.DNS.Rcode;
import org.xbill.DNS.Section;
import org.xbill.DNS.WireParseException;

import java.net.InetAddress;
import java.util.concurrent.*;

@SuppressWarnings("PMD.MoreThanOneLogger")
public abstract class AbstractProtocol implements Protocol {
    private static final Logger ACCESS = LogManager.getLogger("org.apache.traffic_control.traffic_router.core.access");
    private static final Logger LOGGER = LogManager.getLogger(AbstractProtocol.class);

    private static final int NUM_SECTIONS = 4;
    protected boolean shutdownRequested;
    private ThreadPoolExecutor executorService;
    private ExecutorService cancelService;
    private NameServer nameServer;
    private int taskTimeout = 5000; // default
    private int queueDepth = 1000; // default

    /**
     * Gets executorService.
     * 
     * @return the executorService
     */
    public ThreadPoolExecutor getExecutorService() {
        return executorService;
    }

    /**
     * Gets nameServer.
     * 
     * @return the nameServer
     */
    public NameServer getNameServer() {
        return nameServer;
    }

    /**
     * Sets executorService.
     * 
     * @param executorService
     *            the executorService to set
     */
    public void setExecutorService(final ThreadPoolExecutor executorService) {
        this.executorService = executorService;
    }

    /**
     * Sets nameServer.
     * 
     * @param nameServer
     *            the nameServer to set
     */
    public void setNameServer(final NameServer nameServer) {
        this.nameServer = nameServer;
    }

    @Override
    public void shutdown() {
        shutdownRequested = true;
        executorService.shutdownNow();
        cancelService.shutdownNow();
    }

    /**
     * Returns the maximum length of the response.
     * 
     * @param request
     *
     * @return the maximum length in bytes
     */
    protected abstract int getMaxResponseLength(Message request);

    /**
     * Gets shutdownRequested.
     * 
     * @return the shutdownRequested
     */
    protected boolean isShutdownRequested() {
        return shutdownRequested;
    }

    /**
     * Queries the DNS nameServer and returns the response.
     * 
     * @param client
     *            the IP address of the client
     * @param request
     *            the DNS request in wire format
     * @return the DNS response in wire format
     */
    protected byte[] query(final InetAddress client, final byte[] request) throws WireParseException {
        Message query = null;
        Message response = null;
        final long queryTimeMillis = System.currentTimeMillis();
        final DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(queryTimeMillis, client);
        DNSAccessRecord dnsAccessRecord = builder.build();

        try {
            query = new Message(request);
            dnsAccessRecord = builder.dnsMessage(query).build();
            response = getNameServer().query(query, client, builder);
            dnsAccessRecord = builder.dnsMessage(response).build();

            ACCESS.info(DNSAccessEventBuilder.create(dnsAccessRecord));
        } catch (final WireParseException e) {
            ACCESS.info(DNSAccessEventBuilder.create(dnsAccessRecord, e));
            throw e;
        } catch (final Exception e) {
            ACCESS.info(DNSAccessEventBuilder.create(dnsAccessRecord, e));
            response = createServerFail(query);
        }


        return response.toWire(getMaxResponseLength(query));
    }

    /**
     * Submits a request handler to be executed.
     * 
     * @param job
     *            the handler to be executed
     */
	protected void submit(final SocketHandler job) {
		final int queueLength = executorService.getQueue().size();
		Future<?> handler;

		if ((queueDepth > 0 && queueLength >= queueDepth) || (queueDepth == 0 && queueLength > 0)) {
			LOGGER.warn(
				String.format("%s request thread pool full and queue depth limit reached (%d >= %d); discarding request",
				this.getClass().getSimpleName(), queueLength, queueDepth)
			);

			// causes the underlying SocketHandler inner class of each implementing protocol to call a cleanup() method
			job.cancel();

			// add to the cancellation thread pool instead of the task executor pool
			handler = cancelService.submit(job);
		} else {
			handler = executorService.submit(job);
		}

		cancelService.submit(getCanceler(handler));
	}

	private Runnable getCanceler(final Future<?> handler) {
		return new Runnable() {
			public void run() {
				try {
					handler.get(getTaskTimeout(), TimeUnit.MILLISECONDS);
				} catch (InterruptedException | ExecutionException | TimeoutException e) {
					handler.cancel(true);
				}
			}
		};
	}

    private Message createServerFail(final Message query) {
        final Message response = new Message();
        if (query != null) {
            response.setHeader(query.getHeader());
            // This has the side effect of clearing counts out of the header
            for (int i = 0; i < NUM_SECTIONS; i++) {
                response.removeAllRecords(i);
            }
            response.addRecord(query.getQuestion(), Section.QUESTION);
        }
        response.getHeader().setRcode(Rcode.SERVFAIL);
        return response;
    }

	public int getTaskTimeout() {
		return taskTimeout;
	}

	public void setTaskTimeout(final int taskTimeout) {
		this.taskTimeout = taskTimeout;
	}

	public void setQueueDepth(final int queueDepth) {
		this.queueDepth = queueDepth;
	}

	public ExecutorService getCancelService() {
		return cancelService;
	}

	public void setCancelService(final ExecutorService cancelService) {
		this.cancelService = cancelService;
	}
}
