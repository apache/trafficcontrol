/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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
import java.util.Collection;
import java.util.List;
import java.util.TreeMap;

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

import com.comcast.cdn.traffic_control.traffic_monitor.health.DsState;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.behaviors.UpdatingAttributeAppender;

public class DsListPanel extends Panel {
//	private static final Logger LOGGER = Logger.getLogger(ServerListPanel.class);
	private static final long serialVersionUID = 1L;

	ListView<String> servers;
	Component[] updateList;
	String dsId;
	
	public DsListPanel(final String id, final Behavior updater, final Component[] updateList) {
		super(id);

		final ModalWindow modal1;
		add(modal1 = new ModalWindow("modal2"));
		modal1.setInitialWidth(1000);
		//      modal1.setCookieName("modal-1");
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
				final int size = DsState.getDsStates().size();
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
						//						if ( cacheState.isError() ) return "error";
						final DsState cs = DsState.get(dsName);
						if ( cs != null && !cs.isAvailable() ) { return "error"; }// "grey"
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
				label = new Label("disabled", new DsStateModel(dsName, "disabledLocations"));
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

	static class DsStateModel extends Model<String> {
		private static final long serialVersionUID = 1L;
		String csName;
		String key;
		public DsStateModel(final String cs, final String key) {
			this.csName = cs;
			this.key = key;
		}
		@Override
		public String getObject( ) {
			final DsState cs = DsState.getState(csName);
			if(cs == null) { return "err"; }
			if("_status_string_".equals(key)) {
				return cs.getStatusString();
			}
			final boolean clearData = cs.getBool("clearData");
			if(clearData) { return "-"; }
			return cs.getLastValue(key);
		}
	}

	public final List<String> getDsList() {
		final Collection<DsState> list = DsState.getDsStates();
		final TreeMap<String, DsState> tmap = new TreeMap<String, DsState>();
		for(DsState cs : list) {
			tmap.put(cs.getId(), cs);
		}
		return new ArrayList<String>(tmap.keySet());
	}

}