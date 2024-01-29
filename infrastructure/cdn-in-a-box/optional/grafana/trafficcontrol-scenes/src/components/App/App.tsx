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

import { AppRootProps } from "@grafana/data";
import { Routes } from "components/Routes";
import React, { ReactElement } from "react";
import { PluginPropsContext } from "utils/utils.plugin";

/**
 * Renders the component by providing the PluginPropsContext to its children.
 *
 * @returns {ReactElement} The rendered component.
 */
export class App extends React.PureComponent<AppRootProps> {
	/**
	 * Renders the component and returns the JSX element.
	 *
	 * @returns The JSX element representing the rendered component.
	 */
	public render(): ReactElement<AppRootProps> {
		return (
			<PluginPropsContext.Provider value={this.props}>
				<Routes/>
			</PluginPropsContext.Provider>
		);
	}
}
