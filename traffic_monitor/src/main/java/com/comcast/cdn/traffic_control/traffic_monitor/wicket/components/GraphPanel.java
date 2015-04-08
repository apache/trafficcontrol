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

import java.util.HashMap;
import java.util.Map;

import org.apache.wicket.AttributeModifier;
import org.apache.wicket.Component;
import org.apache.wicket.ajax.json.JSONObject;
import org.apache.wicket.markup.head.IHeaderResponse;
import org.apache.wicket.markup.head.JavaScriptContentHeaderItem;
import org.apache.wicket.markup.html.basic.Label;
import org.apache.wicket.markup.html.panel.Panel;
import org.apache.wicket.model.Model;

public class GraphPanel extends Panel {
	private static final long serialVersionUID = 1L;
	private String graphId;
//	private Component graph;
	private int lineCount = 0;
	private int tickCount;
//	private Map<String, Object> options;
	private Model<Integer> size;
//	private int min;

	public GraphPanel(final String id) {
		super(id);

	}

	public GraphPanel(final String string, final Model<Integer> size) {
		super(string);
		final Component graph = new Label("graph_canvas", "");
		this.add(graph);
		graphId = graph.setOutputMarkupId(true).getMarkupId();
		tickCount = 180;
//		this.min = min;
		this.size = size;
	}
	private Map<String, Object> getOptions(final Model<Integer> serverListSize) {
		final Map<String, Object> options = new HashMap<String, Object>();
		//	graph.Set("chart.background.barcolor1", "white");
		//	graph.Set("chart.background.barcolor2", "white");
			options.put("chart.title.xaxis", "Time >>>");
			options.put("chart.title.yaxis", "Servers Down");
			options.put("chart.title.vpos", new Double(0.5));
			options.put("chart.title", "Servers Marked Down On This Traffic Monitor");
			options.put("chart.title.yaxis.pos", new Double(0.5));
			options.put("chart.title.xaxis.pos", new Double(0.5));
			//obj.Set("chart.ylabels.inside", true);
			options.put("chart.yaxispos", "right");
			options.put("chart.ymax", serverListSize.getObject().intValue());
			options.put("chart.xticks", new Integer(25));
			return options;
		//	graph.Set("chart.filled", true);

		//	var grad = graph.context.createLinearGradient(0, 0, 0, 250);
		//	grad.addColorStop(0, "#efefef");
		//	grad.addColorStop(0.9, "rgba(0,0,0,0)");

		//	graph.Set("chart.fillstyle", [ grad ]);
		
	}

	@Override
	public void renderHead(final IHeaderResponse response) {
		super.renderHead(response);

		//final JSONObject o = new JSONObject(getOptions(size));
		//final String script = "$(function() { createMultiLineGraph('"
		//		+graphId+"', "+lineCount+", "+tickCount+", "+o.toString()
		//		+").Draw(); });";
		//response.render(new JavaScriptContentHeaderItem(script, graphId+"_script" , null));
	}

	public void addDataSource(final Component c) {
		if(lineCount > 4) { return; }
		c.add(new AttributeModifier("data-graph-id", Model.of(graphId)));
		c.add(new AttributeModifier("data-graph-index", Model.of(lineCount)));
		lineCount++;
	}

}
