package com.comcast.cdn.traffic_control.traffic_monitor.wicket.components;

import com.comcast.cdn.traffic_control.traffic_monitor.StatisticModel;
import com.comcast.cdn.traffic_control.traffic_monitor.MonitorPage;
import com.comcast.cdn.traffic_control.traffic_monitor.health.StateRegistry;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.MultiUpdatingTimerBehavior;
import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.util.time.Duration;

import java.util.ArrayList;
import java.util.List;

public abstract class StateDetailsPage extends MonitorPage {

	public StateDetailsPage(final String id, final String label) {
		final Behavior updater = new MultiUpdatingTimerBehavior(Duration.seconds(1));
		final Label wicketLabel = new Label(label, id);
		wicketLabel.add(updater);
		add(wicketLabel);

		final List<StatisticModel> statisticModels = new ArrayList<StatisticModel>();

		for (String key : getStateRegistry().get(id).getStatisticsKeys()) {
			statisticModels.add(new StatisticModel(key) {
				@Override
				public String getObject() {
					return (getStateRegistry().has(id)) ? getStateRegistry().get(id, getKey()) : "";
				}
			});
		}

		final ListView<StatisticModel> statesListView = new ListView<StatisticModel>("params", statisticModels) {
			private static final long serialVersionUID = 1L;
			@Override
			protected void populateItem(final ListItem<StatisticModel> item) {
				final StatisticModel statisticModel = item.getModelObject();

				final Label keyLabel = new Label("key", statisticModel.getKey());
				final Label valueLabel = new Label("value", statisticModel);

				valueLabel.add(updater);

				item.add(keyLabel);
				item.add(valueLabel);
			}
		};

		add(statesListView);
	}

	protected abstract StateRegistry getStateRegistry();
}
