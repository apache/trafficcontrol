package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.util.Collection;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

public class StateRegistry<T extends AbstractState> {
	protected final Map<String, T> states = new ConcurrentHashMap<String, T>();

	public T get(final String id) {
		synchronized(states) {
			return states.get(id);
		}
	}

	public Collection<T> getAll() {
		synchronized (states) {
			return states.values();
		}
	}

	public T getOrCreate(final String id) {
		synchronized (states) {
			T abstractState = states.get(id);

			if (abstractState != null) {
				return abstractState;
			}

			return put(createState(id));
		}
	}

	public T put(T abstractState) {
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

	public void removeAllBut(final List<T> states) {
		final Set<String> stateIds = new HashSet<String>();

		for (T state : states) {
			stateIds.add(state.getId());
		}

		removeAllBut(stateIds);
	}

	protected T createState(final String id) {
		return null;
	}

	public void removeAllBut(Set<String> stateIds) {
		synchronized (states) {
			for (String key : states.keySet()) {
				if (!stateIds.contains(key)) {
					states.remove(key);
				}
			}
		}
	}
}
