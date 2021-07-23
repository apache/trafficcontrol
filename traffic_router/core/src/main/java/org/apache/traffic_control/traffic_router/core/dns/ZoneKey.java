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

package org.apache.traffic_control.traffic_router.core.dns;

import java.util.Collections;
import java.util.Date;
import java.util.List;

import org.xbill.DNS.Name;
import org.xbill.DNS.Record;

public class ZoneKey implements Comparable<ZoneKey> {
	private Name name;
	protected List<Record> records;
	private int initialHashCode;
	private long timestamp;

	@SuppressWarnings("unchecked")
	public ZoneKey(final Name name, final List<Record> records) {
		/*
		 * Per the canonical format in  RFC 4034, the records must be in order when the RRset is signed;
		 * sort here to ensure consistency with the ZoneKey, which is based on the hashCode of the List<Record>.
		 * Because we want one set of Records per ZoneKey, regardless of whether DNSSEC is enabled, sort in
		 * this constructor, which is inherited by SignedZoneKey.
		 */
		Collections.sort(records);
		this.setName(name);
		this.setRecords(records);
		this.setInitialHashCode(records.hashCode()); // if the records are signed, the hashCode will change
		this.setTimestamp(System.currentTimeMillis());
	}

	public Name getName() {
		return name;
	}

	private void setName(final Name name) {
		this.name = name;
	}

	public List<Record> getRecords() {
		return records;
	}

	private void setRecords(final List<Record> records) {
		this.records = records;
	}

	private int getInitialHashCode() {
		return initialHashCode;
	}

	private void setInitialHashCode(final int initialHashCode) {
		this.initialHashCode = initialHashCode;
	}

	public long getTimestamp() {
		return timestamp;
	}

	private void setTimestamp(final long timestamp) {
		this.timestamp = timestamp;
	}

	public void updateTimestamp() {
		this.timestamp = System.currentTimeMillis();
	}

	public Date getTimestampDate() {
		return new Date(getTimestamp());
	}

	@Override
	public int hashCode() {
		return getName().hashCode() + getInitialHashCode();
	}

	@Override
	public boolean equals(final Object obj) {
		final ZoneKey ozk = (ZoneKey) obj;
		return getName().equals(ozk.getName()) && getInitialHashCode() == ozk.getInitialHashCode() && obj.getClass().equals(this.getClass());
	}

	// this correctly sorts the names such that the superDomains are last
	@Override
	public int compareTo(final ZoneKey zk) {
		final int i = this.name.compareTo(zk.getName());

		if (i < 0) {
			return 1;
		} else if (i > 0) {
			return -1;
		} else {
			return 0;
		}
	}
}