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

package com.comcast.cdn.traffic_control.traffic_monitor.health;

import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

import org.apache.commons.io.FileUtils;
import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.apache.wicket.ajax.json.JSONException;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.model.Model;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;
import com.comcast.cdn.traffic_control.traffic_monitor.config.MonitorConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.publish.CrConfig;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Fetcher;
import com.comcast.cdn.traffic_control.traffic_monitor.util.PeriodicResourceUpdater;
import com.comcast.cdn.traffic_control.traffic_monitor.util.Updatable;

public class TmWatcher {
	private static final Logger LOGGER = Logger.getLogger(TmWatcher.class);
	static TmWatcher instance;
	List<TmListener> tmlisteners = new ArrayList<TmListener>();
	private final HealthDeterminer hd;
	private final static String CFG_KEY_SUFFIX = "-config";

	public TmWatcher(final HealthDeterminer hd) {
		this.hd = hd;
	}


	public void addTmListener(final TmListener tl) {
		synchronized(tmlisteners) {
			tmlisteners.add(tl);
		}
	}
//	public Updatable getDataServerHandler() {
//		return new Updatable() {
//			@Override
//			public boolean update(final File newDB) {
//				LOGGER.debug("enter: "+newDB);
//				try {
//					final String str = IOUtils.toString(new FileReader(newDB));
//					final JSONArray o = new JSONArray(str);
//					//			LOGGER.warn(o.toString(2));
//					LOGGER.debug("array size: "+o.length());
//					synchronized(tmlisteners) {
//						for(TmListener l : tmlisteners) {
//							try {
//								l.handleServerList(o);
//							} catch(Exception e) {
//								LOGGER.error(e.toString(), e);
//							}
//						}
//					}
//				} catch (Exception e) {
//					LOGGER.warn("error on update: "+newDB, e);
//					return false;
//				}
//				return true;
//			}
//			@Override
//			public boolean update(final JSONObject jsonObject) throws JSONException {
//				return false;
//			}
//		};
//	}
	public Updatable getConfigHandler(final MonitorConfig config, final Updatable updateHandler) {
		return new Updatable() {
			@Override
			public boolean update(final File newDB) {
				LOGGER.debug("enter: "+newDB);
				try {
					final String str = IOUtils.toString(new FileReader(newDB));
					final JSONObject o = new JSONObject(str);
					String cfgKey = null;

					@SuppressWarnings("unchecked")
					final Iterator<String> it = o.keys();
					while (it.hasNext()) {
						final String key = it.next();
						LOGGER.info("KEY -> " + key);

						if (key.endsWith(CFG_KEY_SUFFIX)) {
							cfgKey = key;
							break;
						}
					}

					if (cfgKey != null) {
						config.update(o.getJSONObject(cfgKey));
						updateHandler.update(o);
						return true;
					} else {
						LOGGER.fatal("Unable to find configuration key in health JSON; must end with " + CFG_KEY_SUFFIX);
						return false;
					}
				} catch (Exception e) {
					LOGGER.warn("error on update: "+newDB, e);
					return false;
				}
			}
			@Override
			public boolean update(final JSONObject jsonObject) throws JSONException {
				return false;
			}

		};
	}

	private Updatable getCrConfigUpdateHandler() {
		return new Updatable() {
			@Override
			public boolean update(final File newDB) {
				try {
					final String jsonStr = FileUtils.readFileToString(newDB);
					final JSONObject jo = new JSONObject(jsonStr);
					for(TmListener l : tmlisteners) {
						try {
							l.handleCrConfig(jo);
						} catch(Exception e) {
							LOGGER.error(e.toString(), e);
						}
					}
					return true;
				} catch (IOException e) {
					LOGGER.warn(e,e);
				} catch (JSONException e) {
					LOGGER.warn(e,e);
				}
				return false;
			}
			@Override
			public boolean update(final JSONObject jsonObject) throws JSONException {
				return false;
			}
		};
	}

	public static TmWatcher getInstance() {
		return instance;
	}

	public void refresh() {
		tmUpdater.forceUpdate();
	}

	PeriodicResourceUpdater tmUpdater;
	public void init() {
		synchronized(LOGGER) {
			if(instance == null) {
				instance = this;
			}
		}
		final MonitorConfig config = ConfigHandler.getInstance().getConfig();
		tmUpdater = new PeriodicResourceUpdater(
				new Model<Long>() {
					private static final long serialVersionUID = 1L;
					@Override
					public Long getObject( ) {
						return config.getTmFrequency();
					}
				}) {
			@Override
			protected File fetchFile(final String url) throws IOException {
				return Fetcher.downloadTM(url, config.getAuthUrl(), config.getAuthUsername(), config.getAuthPassword(), config.getConnectionTimeout());
			}
		};
		tmUpdater.add(this.getConfigHandler(config, hd.getUpdateHandler()), new Model<String>() {
			private static final long serialVersionUID = 1L;
			@Override
			public String getObject( ) {
				return config.getHeathUrl();
			}
		}, "health-params.js");

		this.addTmListener(CrConfig.getCrConfigListener());
		tmUpdater.add(getCrConfigUpdateHandler(), new Model<String>() {
			private static final long serialVersionUID = 1L;
			@Override
			public String getObject( ) {
				return config.getCrConfigUrl();
			}
		}, "cr-config.json");

		tmUpdater.init();

	}



	public void destroy() {
		tmUpdater.destroy();
	}

}
