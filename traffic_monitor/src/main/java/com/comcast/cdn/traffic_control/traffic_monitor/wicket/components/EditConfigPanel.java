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

import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.AjaxRequestTarget;
import org.apache.wicket.ajax.json.JSONArray;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.ajax.markup.html.AjaxLink;
import org.apache.wicket.ajax.markup.html.form.AjaxButton;
import org.apache.wicket.extensions.ajax.markup.html.AjaxEditableChoiceLabel;
import org.apache.wicket.extensions.ajax.markup.html.AjaxEditableLabel;
import org.apache.wicket.markup.html.WebMarkupContainer;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.form.Form;
import org.apache.wicket.markup.html.list.ListItem;
import org.apache.wicket.markup.html.list.ListView;
import org.apache.wicket.markup.html.panel.Panel;
import org.apache.wicket.model.Model;
import org.apache.wicket.model.util.WildcardListModel;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.health.TmWatcher;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Fetcher;

public class EditConfigPanel extends Panel {
	private static final Logger LOGGER = Logger.getLogger(EditConfigPanel.class);
	private static final long serialVersionUID = 1L;
	
	private static final String TM_HOSTNAME_KEY = "tm.hostname";
	private static final String CDN_NAME_KEY = "cdnName";
	AjaxEditableChoiceLabel<String> cdnName;

	public EditConfigPanel(final String id) {
		super(id);

		final MonitorConfig config = ConfigHandler.getConfig();

		if (config == null || !config.allowConfigEdit()) {
			return;
		}

		final WebMarkupContainer c = new WebMarkupContainer("configeditpanel");
		add(c);

		this.setOutputMarkupId(true);
		final Form<MonitorConfig>myeditform = new Form<MonitorConfig>("configform");
		setForm(myeditform, ConfigHandler.getConfig().getBaseProps());
		c.add(myeditform);
	}

	private final void setForm(final Form<MonitorConfig> editform, final Map<String, String> baseconfig) {
		final MonitorConfig config = ConfigHandler.getConfig();

		if (config == null || !config.allowConfigEdit()) {
			return;
		}

		final List<String> keys = sort(baseconfig.keySet());
		keys.remove(TM_HOSTNAME_KEY);
		keys.remove(CDN_NAME_KEY);

		final List<String> TM_HOSTS = Arrays.asList(new String[] {
				"tm.company.net" 
			});
		final ArrayList<String> cdnList = new ArrayList<String>();
		final String hostName = baseconfig.get(TM_HOSTNAME_KEY);
		setCdnList(cdnList, hostName);
		final AjaxEditableChoiceLabel<String> tmHost = new AjaxEditableChoiceLabel<String>(
				"tmHost", new BaseConfigModel(TM_HOSTNAME_KEY, baseconfig), TM_HOSTS) {
			private static final long serialVersionUID = 1L;
			@Override
			protected void onSubmit(final AjaxRequestTarget target) {
				super.onSubmit(target);
				final String hostName = baseconfig.get(TM_HOSTNAME_KEY);
				setCdnList(cdnList, hostName);
				target.add(cdnName);
			}
		};
//		listSites
		editform.add(tmHost);

		cdnName = new AjaxEditableChoiceLabel<String>(
				CDN_NAME_KEY, new BaseConfigModel(CDN_NAME_KEY, baseconfig),  new WildcardListModel<String>(cdnList));
//		listSites
		editform.add(cdnName);

		final ListView<String> propView = new ListView<String>("propList", keys) {
			private static final long serialVersionUID = 1L;

			@Override
			protected void populateItem(final ListItem<String> item) {
				final String key = item.getModelObject();

				final Label label = new Label("key", key);
				item.add(label);

				//				editform.add(new TextField("value", new Model<String>(baseconfig.get(key)))) ;
				item.add(new AjaxEditableLabel<String>("value", new BaseConfigModel(key, baseconfig)));

				item.add(label);
			}
		};
		editform.add(propView);

		editform.add(new AjaxLink<Object>("cancel") {
			private static final long serialVersionUID = 1L;
			@Override
			public void onClick(final AjaxRequestTarget target) {
				final Map<String, String> oldProps = ConfigHandler.getConfig().getBaseProps();
				baseconfig.putAll(oldProps);
				final String hostName = baseconfig.get(TM_HOSTNAME_KEY);
				setCdnList(cdnList, hostName);
				if(target != null) {
					target.add(EditConfigPanel.this);
				}
			}
		});
		editform.add(new AjaxButton("submit", editform) {
			private static final long serialVersionUID = 1L;
			@Override
			public void onSubmit(final AjaxRequestTarget target, final Form<?> form) {
				try {
					ConfigHandler.saveBaseConfig(baseconfig);
					TmWatcher.getInstance().refresh();
				} catch (JSONException e) {
					LOGGER.warn(e,e);
				} catch (IOException e) {
					LOGGER.warn(e,e);
				}
//				EditConfigPanel.this.showForm(false);
				if(target != null) {
					target.add(EditConfigPanel.this);
				}
			}
		});
	}

	protected final void setCdnList(final List<String> cdnList, final String hostName) {
		final MonitorConfig config = ConfigHandler.getConfig();

		if (config == null || !config.allowConfigEdit()) {
			return;
		}

		cdnList.clear();
		try {
			final File file = Fetcher.downloadTM("https://"+hostName+"/dataparameter", config.getAuthUrl(), config.getAuthUsername(), config.getAuthPassword(), 500);
			final String str = IOUtils.toString(new FileReader(file));
			file.delete();
			final JSONArray ja = new JSONArray(str);
//			LOGGER.warn(ja.toString(2));
			for(int i = 0; i < ja.length(); i++) {
				// "name": "CDN_name",
				final JSONObject jo = ja.getJSONObject(i);
				if("CDN_name".equals(jo.optString("name"))) {
					cdnList.add(jo.optString("value"));
				}
			}
		} catch (java.net.SocketTimeoutException e) {
			cdnList.add(ConfigHandler.getConfig().getBaseProps().get(CDN_NAME_KEY));
			LOGGER.warn("TM timeout");
		} catch (Exception e) {
			cdnList.add(ConfigHandler.getConfig().getBaseProps().get(CDN_NAME_KEY));
			LOGGER.warn(e,e);
		}
	}

	private List<String> sort(final Set<String> props) {
		final TreeSet<String> set = new TreeSet<String>(props);
 		return new ArrayList<String>(set);
	}

	static class BaseConfigModel extends Model<String> {
		private static final long serialVersionUID = 1L;
		final private Map<String, String> config;
		final String key;
		public BaseConfigModel(final String key, final Map<String, String> config) {
			this.key = key;
			this.config = config;
		}
		@Override
		public String getObject( ) {
			if(config == null) { return "[no config]"; }
			final String r = config.get(key);
			if(r == null) { return "[null]"; }
			return r;
		}
		@Override
		public void setObject(final String val ) {
			config.put(key,val);
		}
	}
}

