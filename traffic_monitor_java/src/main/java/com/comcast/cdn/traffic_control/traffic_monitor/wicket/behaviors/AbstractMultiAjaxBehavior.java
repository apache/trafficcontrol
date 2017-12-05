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

import java.util.List;

import org.apache.wicket.Component;
import org.apache.wicket.Page;
import org.apache.wicket.WicketRuntimeException;
import org.apache.wicket.ajax.AjaxChannel;
import org.apache.wicket.ajax.AjaxRequestTarget;
import org.apache.wicket.ajax.IAjaxIndicatorAware;
import org.apache.wicket.ajax.attributes.AjaxRequestAttributes;
import org.apache.wicket.ajax.attributes.CallbackParameter;
import org.apache.wicket.ajax.attributes.IAjaxCallListener;
import org.apache.wicket.ajax.attributes.ThrottlingSettings;
import org.apache.wicket.ajax.attributes.AjaxRequestAttributes.Method;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.json.JsonFunction;
import org.apache.wicket.ajax.json.JsonUtils;
import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.behavior.IBehaviorListener;
import org.apache.wicket.markup.head.IHeaderResponse;
import org.apache.wicket.markup.head.JavaScriptHeaderItem;
import org.apache.wicket.markup.html.IComponentAwareHeaderContributor;
import org.apache.wicket.protocol.http.WebApplication;
import org.apache.wicket.request.Url;
import org.apache.wicket.request.cycle.RequestCycle;
import org.apache.wicket.request.mapper.parameter.PageParameters;
import org.apache.wicket.request.resource.PackageResourceReference;
import org.apache.wicket.request.resource.ResourceReference;
import org.apache.wicket.resource.CoreLibrariesContributor;
import org.apache.wicket.util.lang.Args;
import org.apache.wicket.util.string.Strings;
import org.apache.wicket.util.time.Duration;

/**
 * This class is mostly copied from
 * @see org.apache.wicket.ajax.AbstractDefaultAjaxBehavior
 */
public abstract class AbstractMultiAjaxBehavior extends Behavior implements IBehaviorListener
{
	private static final long serialVersionUID = 1L;
	public static final String FUNC_STR = "function(attrs){%s}";
	protected Component component;

	/** reference to the default indicator gif file. */
	public static final ResourceReference INDICATOR = new PackageResourceReference(AbstractMultiAjaxBehavior.class, "indicator.gif");

	/**
	 * Subclasses should call super.onBind()
	 * 
	 * @see org.apache.wicket.behavior.AbstractAjaxBehavior#onBind()
	 */
	protected void onBind() {
		getComponent().setOutputMarkupId(true);
	}

	/**
	 * @see org.apache.wicket.behavior.AbstractAjaxBehavior#renderHead(Component,
	 *      org.apache.wicket.markup.head.IHeaderResponse)
	 */
	@Override
	public void renderHead(final Component component, final IHeaderResponse response) {
		super.renderHead(component, response);

		CoreLibrariesContributor.contributeAjax(component.getApplication(), response);

		final RequestCycle requestCycle = component.getRequestCycle();
		final Url baseUrl = requestCycle.getUrlRenderer().getBaseUrl();
		final CharSequence ajaxBaseUrl = Strings.escapeMarkup(baseUrl.toString());
		response.render(JavaScriptHeaderItem.forScript("Wicket.Ajax.baseUrl=\"" + ajaxBaseUrl + "\";", "wicket-ajax-base-url"));

		renderExtraHeaderContributors(component, response);
	}

	/**
	 * Renders header contribution by IAjaxCallListener instances which
	 * additionally implement IComponentAwareHeaderContributor interface.
	 * 
	 * @param component
	 *            the component assigned to this behavior
	 * @param response
	 *            the current header response
	 */
	private void renderExtraHeaderContributors(final Component component, final IHeaderResponse response) {
		for (IAjaxCallListener ajaxCallListener : getAttributes().getAjaxCallListeners()) {
			if (ajaxCallListener instanceof IComponentAwareHeaderContributor) {
				final IComponentAwareHeaderContributor contributor = (IComponentAwareHeaderContributor) ajaxCallListener;
				contributor.renderHead(component, response);
			}
		}
	}

	/**
	 * @return the Ajax settings for this behavior
	 * @since 6.0
	 */
	protected final AjaxRequestAttributes getAttributes() {
		return updateAjaxAttributes(new AjaxRequestAttributes());
	}

	/**
	 * Gives a chance to the specializations to modify the attributes.
	 * 
	 * @param attributes
	 * @since 6.0
	 */
	protected AjaxRequestAttributes updateAjaxAttributes(final AjaxRequestAttributes attributes) {
		return attributes;
	}

	/**
	 * <pre>
	 * 				{
	 * 					u: 'editable-label?6-1.IBehaviorListener.0-text1-label',  // url
	 * 					m: 'POST',          // method name. Default: 'GET'
	 * 					c: 'label7',        // component id (String) or window for page
	 * 					e: 'click',         // event name
	 * 					sh: [],             // list of success handlers
	 * 					fh: [],             // list of failure handlers
	 * 					pre: [],            // list of preconditions. If empty set default : Wicket.$(settings{c}) !== null
	 * 					ep: {},             // extra parameters
	 * 					async: true|false,  // asynchronous XHR or not
	 * 					ch: 'someName|d',   // AjaxChannel
	 * 					i: 'indicatorId',   // indicator component id
	 * 					ad: true,           // allow default
	 * 				}
	 * </pre>
	 * 
	 * @param component
	 *            the component with that behavior
	 * @return the attributes as string in JSON format
	 */
	protected final CharSequence renderAjaxAttributes(final Component component) {
		final AjaxRequestAttributes attributes = getAttributes();
		final JSONObject attributesJson = new JSONObject();

		try {
			attributesJson.put("u", getCallbackUrl());

			appendIfNotEmpty(attributesJson, "f", attributes.getFormId());
			appendBooleanIf(attributesJson, "mp", attributes.isMultipart(), true);
			appendIfNotEmpty(attributesJson, "sc", attributes
					.getSubmittingComponentName());
			appendIfNotEmpty(attributesJson, "i", findIndicatorId());
			appendListenerAtts(attributesJson, attributes.getAjaxCallListeners());
			appendDynamicExtraParameters(attributes.getDynamicExtraParameters(), attributesJson);
			appendBooleanIf(attributesJson, "async", attributes.isAsynchronous(), false);
			appendBooleanIf(attributesJson, "ad", attributes.isAllowDefault(), true);
			appendBooleanIf(attributesJson, "wr", attributes.isWicketAjaxResponse(), false);

			appendSpecial(attributesJson, attributes);

			final ThrottlingSettings throttlingSettings = attributes.getThrottlingSettings();
			if (throttlingSettings != null) {
				final JSONObject throttlingSettingsJson = new JSONObject();
				throttlingSettingsJson.put("id", throttlingSettings.getId());
				throttlingSettingsJson.put("d", throttlingSettings.getDelay()
						.getMilliseconds());
				appendBooleanIf(throttlingSettingsJson, "p", throttlingSettings.getPostponeTimerOnUpdate(), true);
				attributesJson.put("tr", throttlingSettingsJson);
			}

			postprocessConfiguration(attributesJson, component);
		} catch (JSONException e) {
			throw new WicketRuntimeException(e);
		}

		return attributesJson.toString();
	}

	private void appendSpecial(final JSONObject attributesJson,
			final AjaxRequestAttributes attributes) throws JSONException {
		final Method method = attributes.getMethod();
		if (Method.POST == method) {
			attributesJson.put("m", method);
		}
		if (component instanceof Page == false) {
			final String componentId = component.getMarkupId();
			attributesJson.put("c", componentId);
		}
		final JSONArray extraParameters = JsonUtils.asArray(attributes
				.getExtraParameters());
		if (extraParameters.length() > 0) {
			attributesJson.put("ep", extraParameters);
		}

		final String[] eventNames = attributes.getEventNames();
		if (eventNames.length == 1) {
			attributesJson.put("e", eventNames[0]);
		} else {
			for (String eventName : eventNames) {
				attributesJson.append("e", eventName);
			}
		}
		final AjaxChannel channel = attributes.getChannel();
		if (channel != null) {
			attributesJson.put("ch", channel);
		}
		final Duration requestTimeout = attributes.getRequestTimeout();
		if (requestTimeout != null) {
			attributesJson.put("rt", requestTimeout.getMilliseconds());
		}
		final String dataType = attributes.getDataType();
		if (AjaxRequestAttributes.XML_DATA_TYPE.equals(dataType) == false) {
			attributesJson.put("dt", dataType);
		}
	}

	private void appendIfNotEmpty(final JSONObject attributesJson, final String key,
			final String str) throws JSONException {
		if (Strings.isEmpty(str) == false) {
			attributesJson.put(key, str);
		}
	}

	private void appendBooleanIf(final JSONObject attributesJson,
			final String key, final boolean myB, final boolean match) throws JSONException {
		if (myB == match) {
			attributesJson.put(key, match);
		}
	}

	private static void appendDynamicExtraParameters(
			final List<CharSequence> dynamicExtraParameters, final JSONObject attributesJson) throws JSONException {
		if (dynamicExtraParameters == null) { return; }
		for (CharSequence dynamicExtraParameter : dynamicExtraParameters) {
			attributesJson.append("dep", new JsonFunction(
					String.format(
							"function(attrs){%s}",
							dynamicExtraParameter)
					));
		}
	}

	private void appendListenerAtts(final JSONObject attributesJson, 
			final List<IAjaxCallListener> ajaxCallListeners) throws JSONException {
		for (IAjaxCallListener ajaxCallListener : ajaxCallListeners) {
			if (ajaxCallListener == null) { continue; }
			appendListenerHandler(ajaxCallListener.getBeforeHandler(component), attributesJson, "bh",FUNC_STR);
			appendListenerHandler(ajaxCallListener.getBeforeSendHandler(component), attributesJson,
					"bsh", "function(attrs, jqXHR, settings){%s}");
			appendListenerHandler(ajaxCallListener.getAfterHandler(component), attributesJson, "ah", FUNC_STR);
			appendListenerHandler(ajaxCallListener.getSuccessHandler(component), attributesJson, "sh",
					"function(attrs, jqXHR, data, textStatus){%s}");
			appendListenerHandler(ajaxCallListener.getFailureHandler(component), attributesJson, "fh",
					"function(attrs, jqXHR, errorMessage, textStatus){%s}");
			appendListenerHandler(ajaxCallListener.getCompleteHandler(component), attributesJson,
					"coh", "function(attrs, jqXHR, textStatus){%s}");
			appendListenerHandler(ajaxCallListener.getPrecondition(component), attributesJson, "pre", FUNC_STR);
		}
	}

	private static void appendListenerHandler(final CharSequence handler,
			final JSONObject attributesJson, final String propertyName,
			final String functionTemplate) throws JSONException {
		if (Strings.isEmpty(handler) == false) {
			final JsonFunction function;
			if (handler instanceof JsonFunction) {
				function = (JsonFunction) handler;
			} else {
				final String func = String.format(functionTemplate, handler);
				function = new JsonFunction(func);
			}
			attributesJson.append(propertyName, function);
		}
	}

	/**
	 * Gives a chance to modify the JSON attributesJson that is going to be used
	 * as attributes for the Ajax call.
	 * 
	 * @param attributesJson
	 *            the JSON object created by #renderAjaxAttributes()
	 * @param component
	 *            the component with the attached Ajax behavior
	 * @throws JSONException
	 *             thrown if an error occurs while modifying
	 *             {@literal attributesJson} argument
	 */
	protected void postprocessConfiguration(final JSONObject attributesJson, final Component component) throws JSONException {
	}

	/**
	 * @return javascript that will generate an ajax GET request to this
	 *         behavior with its assigned component
	 */
	public CharSequence getCallbackScript() {
		return getCallbackScript(getComponent());
	}

	/**
	 * @param component
	 *            the component to use when generating the attributes
	 * @return script that can be used to execute this Ajax behavior.
	 */
	// 'protected' because this method is intended to be called by other
	// Behavior methods which
	// accept the component as parameter
	protected CharSequence getCallbackScript(final Component component) {
		final CharSequence ajaxAttributes = renderAjaxAttributes(component);
		return "Wicket.Ajax.ajax(" + ajaxAttributes + ");";
	}

	/**
	 * Generates a javascript function that can take parameters and performs an
	 * AJAX call which includes these parameters. The generated code looks like
	 * this:
	 * 
	 * <pre>
	 * function(param1, param2) {
	 *    var attrs = attrsJson;
	 *    var params = {'param1': param1, 'param2': param2};
	 *    attrs.ep = jQuery.extend(attrs.ep, params);
	 *    Wicket.Ajax.ajax(attrs);
	 * }
	 * </pre>
	 * 
	 * @param extraParameters
	 * @return A function that can be used as a callback function in javascript
	 */
	public CharSequence getCallbackFunction(
			final CallbackParameter... extraParameters) {
		final StringBuilder sb = new StringBuilder();
		sb.append("function (");
		boolean first = true;
		for (CallbackParameter curExtraParameter : extraParameters) {
			if (curExtraParameter.getFunctionParameterName() != null) {
				if (!first)
					{ sb.append(','); }
				else
					{ first = false; }
				sb.append(curExtraParameter.getFunctionParameterName());
			}
		}
		sb.append(") {\n");
		sb.append(getCallbackFunctionBody(extraParameters));
		sb.append("}\n");
		return sb;
	}

	/**
	 * Generates the body the
	 * {@linkplain #getCallbackFunction(CallbackParameter...) callback function}
	 * . To embed this code directly into a piece of javascript, make sure any
	 * context parameters are available as local variables, global variables or
	 * within the closure.
	 * 
	 * @param extraParameters
	 * @return The body of the
	 *         {@linkplain #getCallbackFunction(CallbackParameter...) callback
	 *         function}.
	 */
	public CharSequence getCallbackFunctionBody(final CallbackParameter... extraParameters) {
		final CharSequence attrsJson = renderAjaxAttributes(getComponent());
		final StringBuilder sb = new StringBuilder();
		sb.append("var attrs = ");
		sb.append(attrsJson);
		sb.append(";\n");
		sb.append("var params = {");
		boolean first = true;
		for (CallbackParameter curExtraParameter : extraParameters) {
			if (curExtraParameter.getAjaxParameterName() != null) {
				if (!first)
					{ sb.append(','); }
				else
					{ first = false; }
				sb.append('\'')
						.append(curExtraParameter.getAjaxParameterName())
						.append("': ")
						.append(curExtraParameter.getAjaxParameterCode());
			}
		}
		sb.append("};\n");
		if (getAttributes().getExtraParameters().isEmpty()) {
			sb.append("attrs.ep = params;\n");
		} else {
			sb.append("attrs.ep = Wicket.merge(attrs.ep, params);\n");
		}
		sb.append("Wicket.Ajax.ajax(attrs);\n");
		return sb;
	}

	/**
	 * Finds the markup id of the indicator. The default search order is:
	 * component, behavior, component's parent hierarchy.
	 * 
	 * @return markup id or <code>null</code> if no indicator found
	 */
	protected String findIndicatorId() {
		if (getComponent() instanceof IAjaxIndicatorAware) {
			return ((IAjaxIndicatorAware) getComponent())
					.getAjaxIndicatorMarkupId();
		}

		if (this instanceof IAjaxIndicatorAware) {
			return ((IAjaxIndicatorAware) this).getAjaxIndicatorMarkupId();
		}

		Component parent = getComponent().getParent();
		while (parent != null) {
			if (parent instanceof IAjaxIndicatorAware) {
				return ((IAjaxIndicatorAware) parent)
						.getAjaxIndicatorMarkupId();
			}
			parent = parent.getParent();
		}
		return null;
	}

	/**
	 * @see org.apache.wicket.behavior.IBehaviorListener#onRequest()
	 */
	@Override
	public final void onRequest() {
		final WebApplication app = (WebApplication) getComponent().getApplication();
		final AjaxRequestTarget target = app.newAjaxRequestTarget(getComponent()
				.getPage());

		final RequestCycle requestCycle = RequestCycle.get();
		requestCycle.scheduleRequestHandlerAfterCurrent(target);

		respond(target);
	}

	protected abstract void respond(AjaxRequestTarget target);

	protected final Component getComponent() {
		return component;
	}

	@Override
	public void bind(final Component hostComponent) {
		Args.notNull(hostComponent, "hostComponent");

		if (component != null) {
			throw new IllegalStateException(
					"this kind of handler cannot be attached to "
							+ "multiple components; it is already attached to component "
							+ component + ", but component " + hostComponent
							+ " wants to be attached too");
		}

		component = hostComponent;

		// call the callback
		onBind();
	}

	/**
	 * @return the url that references this handler
	 */
	public CharSequence getCallbackUrl() {
		if (getComponent() == null) {
			throw new IllegalArgumentException("Behavior must be bound to a component to create the URL");
		}

		return getComponent().urlFor(this, IBehaviorListener.INTERFACE, new PageParameters());
	}
}