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

import java.util.List;
import java.util.concurrent.ExecutorService;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.context.support.ClassPathXmlApplicationContext;

import org.apache.traffic_control.traffic_router.core.dns.protocol.Protocol;

public final class NameServerMain {
    private static final Logger LOGGER = LogManager.getLogger(NameServerMain.class);

    private ExecutorService protocolService;
    private List<Protocol> protocols;

    /**
     * Shuts down all configured protocols.
     */
    public void destroy() {
        for (final Protocol protocol : getProtocols()) {
            protocol.shutdown();
        }
        getProtocolService().shutdownNow();
    }

    /**
     * Gets protocols.
     * 
     * @return the protocols
     */
    public List<Protocol> getProtocols() {
        return protocols;
    }

    /**
     * Gets protocolService.
     * 
     * @return the protocolService
     */
    public ExecutorService getProtocolService() {
        return protocolService;
    }

    /**
     * Initializes the configured protocols.
     */
    public void init() {
        for (final Protocol protocol : getProtocols()) {
            getProtocolService().submit(protocol);
        }
    }

    /**
     * Sets protocols.
     * 
     * @param protocols
     *            the protocols to set
     */
    public void setProtocols(final List<Protocol> protocols) {
        this.protocols = protocols;
    }

    /**
     * Sets protocolService.
     * 
     * @param protocolService
     *            the protocolService to set
     */
    public void setProtocolService(final ExecutorService protocolService) {
        this.protocolService = protocolService;
    }

    /**
     * @param args
     */
    public static void main(final String[] args) {
        try (ClassPathXmlApplicationContext ctx = new ClassPathXmlApplicationContext("/dns-traffic-router.xml")) {
            ctx.getBean("NameServerMain");
            LOGGER.info("PROCESS_SUCCEEDED");
        } catch (final Exception e) {
            LOGGER.fatal("PROCESS_FAILED");
            LOGGER.fatal(e.getMessage(), e);
            System.exit(1);
        }
    }

}
