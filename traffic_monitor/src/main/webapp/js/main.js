/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * 
 *   http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

	function updateAjaxComponents(arg) {
		for ( var key in arg) {
			if (arg.hasOwnProperty(key)) {
				var $e = $("#"+key);
				var o = arg[key];
				for (var key2 in o) {
					if(key2 === "v") {
						$e.text(o["v"]);
						var graphId = $e.attr("data-graph-id");
						if(graphId != null) {
							var index = $e.attr("data-graph-index");
						}
					} else {
						$e.attr(key2, o[key2]);
					}
				}
			}
		}
	}

