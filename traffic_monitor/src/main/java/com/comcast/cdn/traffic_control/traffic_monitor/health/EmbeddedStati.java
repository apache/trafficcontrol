package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.HashMap;
import java.util.Map;

public class EmbeddedStati implements java.io.Serializable {
	private static final long serialVersionUID = 1L;
	private DsStati currentDtati;
	private final String id;

	public EmbeddedStati(final String base, final String id, final String delimiter) {
		final StringBuilder statId = new StringBuilder();

		if (base != null) {
			statId.append(base);
			statId.append(delimiter);
		}

		statId.append(id);

		this.id = statId.toString();
	}

	public EmbeddedStati(final String base, final String id) {
		this(base, id, ".");
	}

	public void accumulate(final DsStati stati) {
		if (currentDtati == null) {
			currentDtati = new DsStati(stati);
		} else {
			currentDtati.accumulate(stati);
		}
	}

	public Map<String, String> completeRound() {
		if (currentDtati == null) {
			return null;
		}

		final Map<String, String> r = new HashMap<String, String>();

		r.putAll(currentDtati.getStati(this.getId()));
		currentDtati = null;

		return r;
	}

	public String getId() {
		return id;
	}
}
