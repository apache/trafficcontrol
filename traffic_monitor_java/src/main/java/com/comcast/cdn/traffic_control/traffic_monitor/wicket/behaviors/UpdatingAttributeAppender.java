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

import org.apache.wicket.Component;
import org.apache.wicket.behavior.AttributeAppender;
import org.apache.wicket.markup.MarkupNotFoundException;
import org.apache.wicket.model.IModel;
import org.apache.wicket.util.value.ValueMap;

public class UpdatingAttributeAppender extends AttributeAppender {
	private static final long serialVersionUID = 1L;

	public UpdatingAttributeAppender(final String attribute, final IModel<?> replaceModel, final String s) {
		super(attribute, replaceModel, s);
	}

	public IModel<?> getModel() {
		return super.getReplaceModel();
	}

	public String getAttributeValue(final Component component) {
		String value = null;

		try {
			final ValueMap atts = component.getMarkupAttributes();
			value = toStringOrNull(atts.get(getAttribute()));
		} catch (MarkupNotFoundException e) {
			// Ignore
		}

		return newValue(value, toStringOrNull(getReplaceModel().getObject()));
	}

	private String toStringOrNull(final Object replacementValue) {
		return (replacementValue != null) ? replacementValue.toString() : null;
	}
}
