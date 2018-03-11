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

package com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors;

import java.util.ArrayList;
import java.util.List;

import org.apache.log4j.Logger;
import org.apache.wicket.Component;
import org.apache.wicket.ajax.AjaxRequestTarget;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.core.util.string.JavaScriptUtils;
import org.apache.wicket.markup.head.IHeaderResponse;
import org.apache.wicket.markup.head.JavaScriptHeaderItem;
import org.apache.wicket.markup.head.OnLoadHeaderItem;
import org.apache.wicket.request.http.WebRequest;
import org.apache.wicket.util.lang.Args;
import org.apache.wicket.util.time.Duration;

public class MultiUpdatingTimerBehavior extends AbstractMultiAjaxBehavior {
	private static final Logger LOGGER = Logger.getLogger(MultiUpdatingTimerBehavior.class);
	private static final long serialVersionUID = 1L;
	private static final String WICKET_TIMERS_ID = MultiUpdatingTimerBehavior.class.getSimpleName() + "-timers";

	protected final List<Component> components = new ArrayList<Component>();
	private Duration updateInterval;
	private boolean stopped = false;
	private boolean headRendered = false;

	public MultiUpdatingTimerBehavior(final Duration updateInterval) {
		setUpdateInterval(updateInterval);
	}

	@Override
	public void bind(final Component hostComponent) {
		Args.notNull(hostComponent, "hostComponent");

		if (component == null) {
			component = hostComponent;
		}

		hostComponent.setOutputMarkupId(true);
		components.add(hostComponent);
	}

	@Override
	public void renderHead(final Component component, final IHeaderResponse response) {
		super.renderHead(component, response);

		response.render(JavaScriptHeaderItem.forScript("if (typeof(Wicket.TimerHandles) === 'undefined') {Wicket.TimerHandles = {}}", WICKET_TIMERS_ID));

		final WebRequest request = (WebRequest) component.getRequest();

		if (!isStopped() && (!headRendered || !request.isAjax())) {
			headRendered = true;
			response.render(OnLoadHeaderItem.forScript(getJsTimeoutCall(updateInterval)));
		}
	}

	/**
	 *
	 * @see org.apache.wicket.ajax.AbstractDefaultAjaxBehavior#respond(AjaxRequestTarget)
	 */
	@Override
	protected final void respond(final AjaxRequestTarget target) {
		if (isStopped() || !isEnabled(getComponent())) {
			return;
		}

		target.add(getComponent());
		final JSONObject behaviorComponentsJson = getBehaviorComponents();

		target.appendJavaScript("updateAjaxComponents(" + behaviorComponentsJson.toString() + ");");
		target.getHeaderResponse().render(OnLoadHeaderItem.forScript(getJsTimeoutCall(updateInterval)));
	}

	private JSONObject getBehaviorComponents() {
		final JSONObject jsonObject = new JSONObject();

		for (Component component : components) {
			final JSONObject jsonBehaviors = new JSONObject();

			try {
				final List<UpdatingAttributeAppender> behaviors = component.getBehaviors(UpdatingAttributeAppender.class);

				for (UpdatingAttributeAppender behavior : behaviors) {
					jsonBehaviors.put(behavior.getAttribute(), behavior.getAttributeValue(component));
				}

				if (behaviors.isEmpty()) {
					jsonBehaviors.put("v", String.valueOf(component.getDefaultModel().getObject()));
				}

				jsonObject.put(component.getMarkupId(), jsonBehaviors);
			} catch (JSONException e) {
				LOGGER.warn(e,e);
			}
		}
		return jsonObject;
	}

	protected final void setUpdateInterval(final Duration updateInterval) {
		if (updateInterval == null || updateInterval.getMilliseconds() <= 0) {
			throw new IllegalArgumentException("Invalid update interval");
		}
		this.updateInterval = updateInterval;
	}

	/**
	 * @param updateInterval
	 *            Duration between AJAX callbacks
	 * @return JS script
	 */
	protected final String getJsTimeoutCall(final Duration updateInterval) {
		final CharSequence js = JavaScriptUtils.escapeQuotes(getCallbackScript());

		// this might look strange, but it is necessary for IE not to leak :(
		return "Wicket.TimerHandles['" + getComponent().getMarkupId() + "']" + " = setTimeout('" + js + "', " + updateInterval.getMilliseconds() + ')';
	}

	public final boolean isStopped() {
		return stopped;
	}
}
