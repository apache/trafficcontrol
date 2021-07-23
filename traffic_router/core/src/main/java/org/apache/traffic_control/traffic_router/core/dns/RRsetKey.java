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

import org.xbill.DNS.RRset;
import org.xbill.DNS.Record;

import java.util.Iterator;

public class RRsetKey {
    final private RRset rrset;

    public RRsetKey(final RRset rrset) {
        this.rrset = rrset;
    }

    @Override
    public boolean equals(final Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        final RRsetKey that = (RRsetKey) o;
        final Iterator thisIterator = rrset.rrs(false);
        final Iterator thatIterator = that.rrset.rrs(false);
        while (thisIterator.hasNext() && thatIterator.hasNext()) {
            if (!thisIterator.next().equals(thatIterator.next())) {
                return false;
            }
        }
        return !thisIterator.hasNext() && !thatIterator.hasNext();
    }

    @Override
    public int hashCode() {
        int hashCode = 1;
        final Iterator it = rrset.rrs(false);
        while (it.hasNext()) {
            final Record r = (Record) it.next();
            hashCode = 31*hashCode + (r==null? 0 : r.hashCode());
        }
        return hashCode;
    }
}
