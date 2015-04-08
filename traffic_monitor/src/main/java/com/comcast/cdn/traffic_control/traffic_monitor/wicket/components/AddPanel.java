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

import org.apache.log4j.Logger;
import org.apache.wicket.ajax.AjaxRequestTarget;
import org.apache.wicket.ajax.markup.html.AjaxFallbackLink;
import org.apache.wicket.ajax.markup.html.AjaxLink;
import org.apache.wicket.ajax.markup.html.form.AjaxButton;
import org.apache.wicket.markup.html.form.Form;
import org.apache.wicket.markup.html.form.TextField;
import org.apache.wicket.markup.html.link.Link;
import org.apache.wicket.markup.html.panel.Panel;
import org.apache.wicket.model.Model;

import com.comcast.cdn.traffic_control.traffic_monitor.config.ConfigHandler;

public class AddPanel extends Panel {
	private static final Logger LOGGER = Logger.getLogger(AddPanel.class);
	private static final long serialVersionUID = 1L;
	final Link<Object> addlink;
	final Form<Object> editform;
	private final TextField<String> field;
	
	public AddPanel(final String id) {
		super(id);
		
		this.setOutputMarkupId(true);
		editform = new Form<Object>("editform");
		addlink = new AjaxFallbackLink<Object>("add") {
			private static final long serialVersionUID = 1L;

			@Override
			public void onClick(final AjaxRequestTarget target) {
				AddPanel.this.showForm(true);
				if(target != null) {
					target.add(AddPanel.this);
				}
			}
		};
		add(addlink);
		field = new TextField<String>("field", new Model<String>(""));
		editform.add(field);
		editform.add(new AjaxLink<Object>("cancel") {
			private static final long serialVersionUID = 1L;

			@Override
			public void onClick(final AjaxRequestTarget target) {
				AddPanel.this.showForm(false);
				if(target != null) {
					target.add(AddPanel.this);
				}
			}
		});
		editform.add(new AjaxButton("submit", editform) {
			private static final long serialVersionUID = 1L;

			@Override
			public void onSubmit(final AjaxRequestTarget target, final Form<?> form) {
				addResponse(field.getModelObject().toString(), target);
				AddPanel.this.showForm(false);
				if(target != null) {
					target.add(AddPanel.this);
				}
			}
		});
		this.add(editform);
		showForm(false);
	}

	protected final void showForm(final boolean b) {
		field.setModelObject("");
		if(b) {
			addlink.setVisible(false);
			editform.setVisible(true);
		} else {
			addlink.setVisible(true);
			editform.setVisible(false);
		}
	}
	public void addResponse(final String str, final AjaxRequestTarget target) {
		LOGGER.debug(field.getModelObject());
	}
	public void usageExample() {
		add(new AddPanel("addcontainer") {
			private static final long serialVersionUID = 1L;

			@Override
			public void addResponse(final String str, final AjaxRequestTarget target) {
				LOGGER.debug("finally: "+str);
//				AtsWatcherConfig config = 
						ConfigHandler.getConfig();
//				config.addAtsServer(str);
				ConfigHandler.saveConfig();
//				servers.modelChanged();
//				servers.removeAll();
//				if(target != null) {
//					target.add(container);
//				}
			}
		});

	}
}
