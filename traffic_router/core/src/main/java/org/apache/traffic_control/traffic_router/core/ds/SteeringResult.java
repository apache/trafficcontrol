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

package org.apache.traffic_control.traffic_router.core.ds;

import org.apache.traffic_control.traffic_router.core.edge.Cache;

public class SteeringResult {
    private SteeringTarget steeringTarget;
    private DeliveryService deliveryService;
    private Cache cache;

    public SteeringResult(final SteeringTarget steeringTarget, final DeliveryService deliveryService) {
        this.steeringTarget = steeringTarget;
        this.deliveryService = deliveryService;
    }

    public SteeringTarget getSteeringTarget() {
        return steeringTarget;
    }

    public void setSteeringTarget(final SteeringTarget steeringTarget) {
        this.steeringTarget = steeringTarget;
    }

    public DeliveryService getDeliveryService() {
        return deliveryService;
    }

    public void setDeliveryService(final DeliveryService deliveryService) {
        this.deliveryService = deliveryService;
    }

    public Cache getCache() {
        return cache;
    }

    public void setCache(final Cache cache) {
        this.cache = cache;
    }

}
