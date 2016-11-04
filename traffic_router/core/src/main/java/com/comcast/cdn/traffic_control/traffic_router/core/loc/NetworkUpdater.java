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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import java.io.File;
import java.io.IOException;

import org.apache.wicket.ajax.json.JSONException;

public class NetworkUpdater extends AbstractServiceUpdater {

	public NetworkUpdater() {
		sourceCompressed = false;
		tmpPrefix = "czf";
		tmpSuffix = ".json";
	}

	@Override
	public boolean loadDatabase() throws IOException, JSONException {
		final File existingDB = databasesDirectory.resolve(databaseName).toFile();

		if (!existingDB.exists() || !existingDB.canRead()) {
			return false;
		}

		return NetworkNode.generateTree(existingDB, false) != null;
	}

	@Override
	public boolean verifyDatabase(final File dbFile) throws IOException, JSONException {
		if (!dbFile.exists() || !dbFile.canRead()) {
			return false;
		}

		return NetworkNode.generateTree(dbFile, true) != null;
	}

}
