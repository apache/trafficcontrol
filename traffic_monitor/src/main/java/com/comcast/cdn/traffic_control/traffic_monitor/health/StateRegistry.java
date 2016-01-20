package com.comcast.cdn.traffic_control.traffic_monitor.health;

import com.comcast.cdn.traffic_control.traffic_monitor.KeyValue;
import com.comcast.cdn.traffic_control.traffic_monitor.StateKeyValue;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
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
		return states.put(abstractState.getId(), abstractState);
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
		return get(stateId).getLastValue(key);
	}

	public List<String> getUniqueIds() {
		final Set<String> idSet = new TreeSet<String>();

		for(AbstractState state : states.values()) {
			idSet.add(state.getId());
		}

		final List<String> idList = new ArrayList<String>();
		idList.addAll(idSet);
		return idList;
	}

	public List<KeyValue> getModelList(final String hostname) {
		final List<KeyValue> modelList = new ArrayList<KeyValue>();

		AbstractState cacheState = states.get(hostname);

		for(String key : cacheState.getStatisticsKeys()) {
			modelList.add(new StateKeyValue(key, cacheState, this));
		}

		return modelList;
	}

	protected AbstractState createState(final String id) {
		return null;
	}
}
