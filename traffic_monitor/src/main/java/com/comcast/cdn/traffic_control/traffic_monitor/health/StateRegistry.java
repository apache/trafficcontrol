package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.Collection;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class StateRegistry {
	protected final Map<String, AbstractState> states = new ConcurrentHashMap<String, AbstractState>();

	public AbstractState get(final String id) {
		synchronized(states) {
			return states.get(id);
		}
	}

	public Collection<AbstractState> getAll() {
		synchronized (states) {
			return states.values();
		}
	}

	public AbstractState getOrCreate(final String id) {
		synchronized (states) {
			AbstractState abstractState = states.get(id);

			if (abstractState != null) {
				return abstractState;
			}

			return put(createState(id));
		}
	}

	public AbstractState put(AbstractState abstractState) {
		states.put(abstractState.getId(), abstractState);
		return abstractState;
	}

	public int size() {
		synchronized (states) {
			return states.size();
		}
	}

	public boolean has(final String id) {
		return states.containsKey(id);
	}

	public String get(final String stateId, final String key) {
		if (!has(stateId)) {
			return "";
		}

		return get(stateId).getLastValue(key);
	}

	public long getSumOfLongStatistic(final String key) {
		long sum = 0;
		for(AbstractState state : states.values()) {
			sum += state.getDouble(key);
		}
		return sum;
	}

	protected AbstractState createState(final String id) {
		return null;
	}
}
