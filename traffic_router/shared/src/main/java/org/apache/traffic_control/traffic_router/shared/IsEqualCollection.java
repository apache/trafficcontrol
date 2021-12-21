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

package org.apache.traffic_control.traffic_router.shared;

import org.hamcrest.Description;
import org.hamcrest.Factory;
import org.hamcrest.Matcher;
import org.hamcrest.core.IsEqual;

import java.util.Collection;

final public class IsEqualCollection<T> extends IsEqual<T> {
	private final Object expectedValue;

	private IsEqualCollection(final T equalArg) {
		super(equalArg);
		expectedValue = equalArg;
	}

	private void describeItems(final Description description, final Object value) {
		if (value instanceof Collection) {
			final Object[] items = ((Collection) value).toArray();

			description.appendText("\n{");
			for (final Object item : items) {
				description.appendText("\n\t");
				description.appendText(item.toString());
			}
			description.appendText("\n}");
		}
	}

	@Override
	public void describeTo(final Description description) {
		if (expectedValue instanceof Collection) {
			description.appendText("all of the following in order\n");
			describeItems(description,expectedValue);
			return;
		}

		super.describeTo(description);
	}

	@Override
	public void describeMismatch(final Object actualValue, final Description mismatchDescription) {
		if (actualValue instanceof Collection) {
			mismatchDescription.appendText("had the items\n");
			describeItems(mismatchDescription, actualValue);
			return;
		}

		super.describeMismatch(actualValue, mismatchDescription);
	}

	@Factory
	public static <T> Matcher<T> equalTo(final T operand) {
		return new IsEqualCollection<>(operand);
	}
}
