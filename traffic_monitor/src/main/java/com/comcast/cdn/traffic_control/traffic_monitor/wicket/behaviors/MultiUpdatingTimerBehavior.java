/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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
	protected final List<Component> components = new ArrayList<Component>();

	public MultiUpdatingTimerBehavior(final Duration updateInterval) {
		setUpdateInterval(updateInterval);
	}

	protected  void	onPostProcessTarget(final AjaxRequestTarget target) {
		final JSONObject jo = new JSONObject();
//		final StringBuilder sb = new StringBuilder();
//		sb.append("[");
		for(Component c : components) {
//			sb.append("\"").append(c.getMarkupId()).append("\",");
			final JSONObject jo2 = new JSONObject();
			try {
				boolean hasB = false;
				final List<UpdatingAttributeAppender> blist = c.getBehaviors(UpdatingAttributeAppender.class);
				for(UpdatingAttributeAppender b : blist) {
					jo2.put(b.getAttribute(), b.getAttributeValue(c));
					hasB = true;
				}
				if(!hasB) {
					jo2.put("v", String.valueOf(c.getDefaultModel().getObject())); 
				}
				jo.put(c.getMarkupId(), jo2);
			} catch (JSONException e) {
				LOGGER.warn(e,e);
			}
		}
		target.appendJavaScript("updateAjaxComponents("+jo.toString()+");");
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

	protected final void onTimer(final AjaxRequestTarget target) {
		target.add(getComponent());
		onPostProcessTarget(target);
	}

	
	
	
	
	
	
	private static final String WICKET_TIMERS_ID = MultiUpdatingTimerBehavior.class.getSimpleName() + "-timers";
	/** The update interval */
	private Duration updateInterval;

	private boolean stopped = false;

	private boolean headRendered = false;

	/**
	 * Sets the update interval duration. This method should only be called within the
	 * {@link #onTimer(AjaxRequestTarget)} method.
	 * 
	 * @param updateInterval
	 */
	protected final void setUpdateInterval(final Duration updateInterval) {
		if (updateInterval == null || updateInterval.getMilliseconds() <= 0) {
			throw new IllegalArgumentException("Invalid update interval");
		}
		this.updateInterval = updateInterval;
	}

	/**
	 * Returns the update interval
	 * 
	 * @return The update interval
	 */
	public final Duration getUpdateInterval() {
		return updateInterval;
	}

	@Override
	public void renderHead(final Component component, final IHeaderResponse response) {
		super.renderHead(component, response);

		response.render(JavaScriptHeaderItem.forScript("if (typeof(Wicket.TimerHandles) === 'undefined') {Wicket.TimerHandles = {}}",
				WICKET_TIMERS_ID));

		final WebRequest request = (WebRequest) component.getRequest();

		if (!isStopped() && (!headRendered || !request.isAjax())) {
			headRendered = true;
			response.render(OnLoadHeaderItem.forScript(getJsTimeoutCall(updateInterval)));
		}
	}

	/**
	 * @param updateInterval
	 *            Duration between AJAX callbacks
	 * @return JS script
	 */
	protected final String getJsTimeoutCall(final Duration updateInterval) {
		final CharSequence js = JavaScriptUtils.escapeQuotes(getCallbackScript());

		final String timeoutHandle = getTimeoutHandle();
		// this might look strange, but it is necessary for IE not to leak :(
		return timeoutHandle+" = setTimeout('" + js + "', " +
				updateInterval.getMilliseconds() + ')';
	}

	/**
	 * @return the name of the handle that is used to stop any scheduled timer
	 */
	private String getTimeoutHandle() {
		return "Wicket.TimerHandles['"+getComponent().getMarkupId() + "']";
	}

	/**
	 * 
	 * @see org.apache.wicket.ajax.AbstractDefaultAjaxBehavior#respond(AjaxRequestTarget)
	 */
	@Override
	protected final void respond(final AjaxRequestTarget target) {
		if (!isStopped() && isEnabled(getComponent()))
		{
			onTimer(target);

			target.getHeaderResponse().render(
					OnLoadHeaderItem.forScript(getJsTimeoutCall(updateInterval)));
		}
	}

	/**
	 * @return {@code true} if has been stopped via {@link #stop(AjaxRequestTarget)}
	 */
	public final boolean isStopped() {
		return stopped;
	}

	/**
	 * Re-enables the timer if already stopped
	 *
	 * @param target
	 */
	public final void restart(final AjaxRequestTarget target) {
		if (isStopped()) {
			stopped = false;
			headRendered = false;
			target.add(getComponent());
		}
	}

	/**
	 * Stops the timer
	 */
	public final void stop(final AjaxRequestTarget target) {
		stopped = true;
		final String timeoutHandle = getTimeoutHandle();
		target.prependJavaScript("clearTimeout("+timeoutHandle+"); delete "+timeoutHandle+";");
	}

}
