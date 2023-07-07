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

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.TreeMap;

import org.apache.traffic_control.traffic_router.core.request.Request;
import org.apache.traffic_control.traffic_router.core.request.RequestMatcher;

public class DeliveryServiceMatcher implements Comparable<DeliveryServiceMatcher> {
	public enum Type {
		HOST, HEADER, PATH
	}

	private DeliveryService deliveryService;
	final private List<RequestMatcher> requestMatchers = new ArrayList<RequestMatcher>();

	public DeliveryServiceMatcher(final DeliveryService ds) {
		this.deliveryService = ds;
	}

	public DeliveryService getDeliveryService() {
		return deliveryService;
	}

	public void setDeliveryService(final DeliveryService deliveryService) {
		this.deliveryService = deliveryService;
	}

	public void addMatch(final Type type, final String string, final String target) {
		requestMatchers.add(new RequestMatcher(type, string, target));
	}

	public List<RequestMatcher> getRequestMatchers() {
		return new ArrayList<>(this.requestMatchers);
	}

	public boolean matches(final Request request) {
		for (final RequestMatcher matcher : requestMatchers) {
			if (!matcher.matches(request)) {
				return false;
			}
		}

		return !requestMatchers.isEmpty();
	}

	@Override
	@SuppressWarnings("PMD.IfStmtsMustUseBraces")
	public boolean equals(final Object deliveryServiceMatcher) {
		if (this == deliveryServiceMatcher) return true;
		if (deliveryServiceMatcher == null || getClass() != deliveryServiceMatcher.getClass()) return false;
		final DeliveryServiceMatcher that = (DeliveryServiceMatcher) deliveryServiceMatcher;

		if (deliveryService != null ? !deliveryService.equals(that.deliveryService) : that.deliveryService != null) return false;
		return !(requestMatchers != null ? !requestMatchers.equals(that.requestMatchers) : that.requestMatchers != null);
	}

	@Override
	public int hashCode() {
		int result = deliveryService != null ? deliveryService.hashCode() : 0;
		result = 31 * result + (requestMatchers != null ? requestMatchers.hashCode() : 0);
		return result;
	}


	@SuppressWarnings("PMD.NPathComplexity")
	@Override
	public int compareTo(final DeliveryServiceMatcher that) {
		if (this == that || this.equals(that)) {
			return 0;
		}

		final Set<RequestMatcher> uniqueToThis = new HashSet<RequestMatcher>();
		uniqueToThis.addAll(this.requestMatchers);

		final Set<RequestMatcher> uniqueToThat = new HashSet<RequestMatcher>();
		uniqueToThat.addAll(that.requestMatchers);

		for (final RequestMatcher myRequestMatcher : requestMatchers) {
			if (uniqueToThat.remove(myRequestMatcher)) {
				uniqueToThis.remove(myRequestMatcher);
			}
		}

		final TreeMap<RequestMatcher, DeliveryServiceMatcher> map = new TreeMap<RequestMatcher, DeliveryServiceMatcher>();

		for (final RequestMatcher thisMatcher : uniqueToThis) {
			map.put(thisMatcher, this);
		}

		for (final RequestMatcher thatMatcher : uniqueToThat) {
			map.put(thatMatcher, that);
		}

		if (map.isEmpty()) {
			return 0;
		}

		return (this == map.firstEntry().getValue()) ? -1 : 1;
	}

	@Override
	public String toString() {
		if (requestMatchers.size() > 1) {
			return "DeliveryServiceMatcher{" +
				"deliveryService=" + deliveryService +
				", requestMatchers=" + requestMatchers +
				'}';
		}

		return "DeliveryServiceMatcher{" +
			"deliveryService=" + deliveryService +
			", requestMatcher=" + requestMatchers.get(0) +
			'}';

	}
}
