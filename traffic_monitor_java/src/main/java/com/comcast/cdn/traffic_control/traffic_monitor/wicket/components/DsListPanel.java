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

package com.comcast.cdn.traffic_control.traffic_monitor.wicket.components;

import java.util.ArrayList;
import java.util.List;
import java.util.TreeMap;

import com.comcast.cdn.traffic_control.traffic_monitor.health.AbstractState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DeliveryServiceStateRegistry;
import com.comcast.cdn.traffic_control.traffic_monitor.health.DsState;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.models.DsStateModel;
import org.apache.wicket.Component;
import org.apache.wicket.Page;
import org.apache.wicket.ajax.AbstractAjaxTimerBehavior;
import org.apache.wicket.ajax.AjaxRequestTarget;
import org.apache.wicket.ajax.markup.html.AjaxLink;
import org.apache.wicket.behavior.Behavior;
import org.apache.wicket.extensions.ajax.markup.html.modal.ModalWindow;
import org.apache.wicket.markup.html.WebMarkupContainer;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.markup.html.panel.Panel;
import org.apache.wicket.model.Model;
import org.apache.wicket.util.time.Duration;

import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.UpdatingAttributeAppender;

public class DsListPanel extends Panel {
	private static final long serialVersionUID = 1L;

	ListView<String> servers;
	Component[] updateList;
	String dsId;
	
	public DsListPanel(final String id, final Behavior updater, final Component[] updateList) {
		super(id);

		final ModalWindow modal1;
		add(modal1 = new ModalWindow("modal2"));
		modal1.setInitialWidth(1000);
		modal1.setPageCreator(new ModalWindow.PageCreator() {
			private static final long serialVersionUID = 1L;
			public Page createPage() {
				return new DsDetailsPage(dsId);
			}
		});


		this.updateList = updateList;
		final WebMarkupContainer container = new WebMarkupContainer("listpanel");
		container.setOutputMarkupId(true);
		add(container);
		servers = createDsListView(updater, modal1);
		servers.setOutputMarkupId(true);
		container.setOutputMarkupId(true);
		container.add(servers);

		add(new AbstractAjaxTimerBehavior(Duration.seconds(1)) {
			private static final long serialVersionUID = 1L;
			int serverCount = 0;
			@Override
			protected final void onTimer(final AjaxRequestTarget target) {
				final int size = DeliveryServiceStateRegistry.getInstance().size();
				if(serverCount != size) {
					serverCount = size;
					servers.setList(getDsList());
					target.add(container);
					if(updateList!=null) {
						for(Component c : updateList) {
							target.add(c);
						}
					}
				}
			}
		});
	}

	private ListView<String> createDsListView(final Behavior updater,
			final ModalWindow modalWindow) {
		return new ListView<String>("ds", getDsList()) {
			private static final long serialVersionUID = 1L;

			@Override
			protected void populateItem(final ListItem<String> item) {
				final String dsName = item.getModelObject();

				item.add(new UpdatingAttributeAppender("class", new Model<String>("") {
					private static final long serialVersionUID = 1L;

					@Override
					public String getObject( ) {
						final AbstractState state = DeliveryServiceStateRegistry.getInstance().get(dsName);
						if ( state != null && !state.isAvailable() ) { return "error"; }
						else { return " "; }
					}
				}, " "));
				item.add(updater);

				Label label = new Label("status", new DsStateModel(dsName, "_status_string_"));
				label.add(updater);
				item.add(label);
				label = new Label("kbps", new DsStateModel(dsName, "total.kbps"));
				label.add(updater);
				item.add(label);
				label = new Label("tps", new DsStateModel(dsName, "total.tps_total"));
				label.add(updater);
				item.add(label);
				label = new Label("tps_2xx", new DsStateModel(dsName, "total.tps_2xx"));
				label.add(updater);
				item.add(label);
				label = new Label("tps_3xx", new DsStateModel(dsName, "total.tps_3xx"));
				label.add(updater);
				item.add(label);
				label = new Label("tps_4xx", new DsStateModel(dsName, "total.tps_4xx"));
				label.add(updater);
				item.add(label);
				label = new Label("tps_5xx", new DsStateModel(dsName, "total.tps_5xx"));
				label.add(updater);
				item.add(label);
				label = new Label("disabled", new DsStateModel(dsName, DsState.DISABLED_LOCATIONS));
				label.add(updater);
				item.add(label);

				label = new Label("caches-reporting", new DsStateModel(dsName, "caches-reporting"));
				label.add(updater);
				item.add(label);
				label = new Label("caches-available", new DsStateModel(dsName, "caches-available"));
				label.add(updater);
				item.add(label);
				label = new Label("caches-configured", new DsStateModel(dsName, "caches-configured"));
				label.add(updater);
				item.add(label);

				final AjaxLink<Void> link = new AjaxLink<Void>("fulldetails") {
					private static final long serialVersionUID = 1L;
					@Override
		            public void onClick(final AjaxRequestTarget target) {
						dsId = dsName;
						modalWindow.show(target);
		            }
		        };
				link.setBody(new Model<String>(dsName));
				item.add(link);
			}
		};
	}

	public final List<String> getDsList() {
		final TreeMap<String, AbstractState> tmap = new TreeMap<String, AbstractState>();
		for(AbstractState state : DeliveryServiceStateRegistry.getInstance().getAll()) {
			tmap.put(state.getId(), state);
		}
		return new ArrayList<String>(tmap.keySet());
	}

}