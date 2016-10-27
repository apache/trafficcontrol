package com.comcast.cdn.traffic_control.traffic_router.core;

import org.hamcrest.Description;
import org.hamcrest.Factory;
import org.hamcrest.Matcher;
import org.hamcrest.core.IsEqual;

import java.util.Collection;

public class IsEqualCollection<T> extends IsEqual<T> {
	private final Object expectedValue;

	private IsEqualCollection(T equalArg) {
		super(equalArg);
		expectedValue = equalArg;
	}

	private void describeItems(Description description, Object value) {
		if (value instanceof Collection) {
			Object[] items = ((Collection) value).toArray();

			description.appendText("\n{");
			for (Object item : items) {
				description.appendText("\n\t");
				description.appendText(item.toString());
			}
			description.appendText("\n}");
		}
	}

	@Override
	public void describeTo(Description description) {
		if (expectedValue instanceof Collection) {
			description.appendText("all of the following in order\n");
			describeItems(description,expectedValue);
			return;
		}

		super.describeTo(description);
	}

	@Override
	public void describeMismatch(Object actualValue, Description mismatchDescription) {
		if (actualValue instanceof Collection) {
			mismatchDescription.appendText("had the items\n");
			describeItems(mismatchDescription, actualValue);
			return;
		}

		super.describeMismatch(actualValue, mismatchDescription);
	}

	@Factory
	public static <T> Matcher<T> equalTo(T operand) {
		return new IsEqualCollection<>(operand);
	}
}
