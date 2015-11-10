package com.comcast.cdn.traffic_control.traffic_router.core.util;

public class ComparableStringByLength implements Comparable<ComparableStringByLength> {
	final private String string;

	public ComparableStringByLength(final String string) {
		this.string = string;
	}

	@Override
	public int compareTo(final ComparableStringByLength other) {
		if (string.length() == other.string.length()) {
			return string.compareTo(other.string);
		}

		return (string.length() > other.string.length()) ? -1 : 1;
	}

	@Override
	public String toString() {
		return string;
	}

	@Override
	@SuppressWarnings("PMD.IfStmtsMustUseBraces")
	public boolean equals(final Object other) {
		if (this == other) return true;

		if (other == null || (getClass() != other.getClass() && String.class != other.getClass())) return false;

		if (String.class == other.getClass()) {
			return string.equals(other);
		}

		final ComparableStringByLength that = (ComparableStringByLength) other;

		return !(string != null ? !string.equals(that.string) : that.string != null);

	}

	@Override
	public int hashCode() {
		return string != null ? string.hashCode() : 0;
	}
}
