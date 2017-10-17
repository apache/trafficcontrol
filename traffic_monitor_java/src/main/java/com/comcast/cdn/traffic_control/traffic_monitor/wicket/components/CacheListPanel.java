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
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheState;
import com.comcast.cdn.traffic_control.traffic_monitor.health.CacheStateRegistry;
import com.comcast.cdn.traffic_control.traffic_monitor.wicket.models.CacheStateModel;
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

public class CacheListPanel extends Panel {
	private static final long serialVersionUID = 1L;
	ListView<String> servers;
	Component[] updateList;
	String hostname;

	public CacheListPanel(final String id, final Behavior updater, final Component[] updateList) {
		super(id);

		final ModalWindow modal1;
		add(modal1 = new ModalWindow("modal1"));
		modal1.setInitialWidth(1000);
		modal1.setPageCreator(new ModalWindow.PageCreator() {
			private static final long serialVersionUID = 1L;
			public Page createPage() {
				return new CacheDetailsPage(hostname);
			}
		});

		this.updateList = updateList;
		final WebMarkupContainer container = new WebMarkupContainer("listpanel");
		container.setOutputMarkupId(true);
		add(container);
		servers = createServerListView(updater, modal1);
		servers.setOutputMarkupId(true);
		container.setOutputMarkupId(true);
		container.add(servers);

		add(new AbstractAjaxTimerBehavior(Duration.seconds(1)) {
			private static final long serialVersionUID = 1L;
			int serverCount = 0;
			@Override
			protected final void onTimer(final AjaxRequestTarget target) {
				final int size = CacheStateRegistry.getInstance().size();
				if(serverCount != size) {
					serverCount = size;
					servers.setList(getServerList());
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

	private ListView<String> createServerListView(final Behavior updater, final ModalWindow modalWindow) {
		return new ListView<String>("servers", getServerList()) {
			private static final long serialVersionUID = 1L;

			@Override
			protected void populateItem(final ListItem<String> item) {
				final String cacheName = item.getModelObject();
				final CacheState cacheState = CacheStateRegistry.getInstance().get(cacheName);

				item.add(new UpdatingAttributeAppender("class", new Model<String>("") {
					@Override
					public String getObject() {
						return CacheStateRegistry.getInstance().getStatusString(cacheName);
					}
				}, " "));

				item.add(updater);

				Label label = new Label("type", cacheState.getCache().getType());
				label.add(updater);
				item.add(label);
				label = new Label("status", new CacheStateModel(cacheName, "_status_string_"));
				label.add(updater);
				item.add(label);
				label = new Label("loadavg", new CacheStateModel(cacheName, "loadavg"));
				label.add(updater);
				item.add(label);
				label = new Label("queryTime", new CacheStateModel(cacheName, "queryTime"));
				label.add(updater);
				item.add(label);
				label = new Label("kbps", new CacheStateModel(cacheName, "kbps"));
				label.add(updater);
				item.add(label);
				label = new Label("maxKbps", new CacheStateModel(cacheName, "maxKbps"));
				label.add(updater);
				item.add(label);
				label = new Label("current_client_connections", new CacheStateModel(cacheName, "ats.proxy.process.http.current_client_connections"));
				label.add(updater);
				item.add(label);

				final AjaxLink<Void> link = new AjaxLink<Void>("fulldetails") {
					private static final long serialVersionUID = 1L;
					@Override
					public void onClick(final AjaxRequestTarget target) {
						hostname = cacheName;
						modalWindow.show(target);
					}
				};
				link.setBody(new Model<String>(cacheName));
				item.add(link);
			}
		};
	}

	public final List<String> getServerList() {
		final TreeMap<String, AbstractState> tmap = new TreeMap<String, AbstractState>();
		for(AbstractState cs : CacheStateRegistry.getInstance().getAll()) {
			tmap.put(cs.getId(), cs);
		}
		return new ArrayList<String>(tmap.keySet());
	}

}