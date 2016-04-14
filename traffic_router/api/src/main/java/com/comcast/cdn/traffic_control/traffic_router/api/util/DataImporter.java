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

package com.comcast.cdn.traffic_control.traffic_router.api.util;

import java.lang.management.ManagementFactory;

import javax.management.MBeanServer;
import javax.management.ObjectName;

public class DataImporter {
	private final ObjectName objectName;
	private final MBeanServer mbeanServer;

	public DataImporter(final String mbeanName) throws DataImporterException {
		try {
			this.objectName = new ObjectName(mbeanName);
			this.mbeanServer = ManagementFactory.getPlatformMBeanServer();
		} catch (Exception ex) {
			throw new DataImporterException(ex);
		}
	}

	public Object invokeOperation(final String operation) throws DataImporterException {
		return invokeOperation(operation, new Object[] {});
	}

	public <T> Object invokeOperation(final String operation, final T... parameters) throws DataImporterException {
		String[] signature = new String[parameters.length];

		for (int i = 0; i < parameters.length; i++) {
			signature[i] = parameters[i].getClass().getName();
		}

		try {
			return mbeanServer.invoke(objectName, operation, parameters, signature);
		} catch (Exception ex) {
			throw new DataImporterException(ex);
		}
	}
}
