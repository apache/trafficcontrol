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

package org.apache.traffic_control.traffic_router.core.util;

import java.util.TreeSet;

public class ComparableTreeSet<E> extends TreeSet<E> implements Comparable<ComparableTreeSet<E>> {
    private static final long serialVersionUID = 1L;

    @Override
    public int compareTo(final ComparableTreeSet<E> o) {
        if (isEmpty() && !o.isEmpty()) {
            return 1;
        }
        else if (o.isEmpty()) {
            return -1;
        }

        if (this.equals(o)) {
            return 0;
        }

        if (containsAll(o)) {
            // this comes first because it is a superset??????
            return -1;
        }

        if (o.containsAll(this)) {
            return 1;
        }

        final Object item = first();
        final Object otherItem = o.first();
        if (item instanceof Comparable) {
            return ((Comparable) item).compareTo(otherItem);
        }

        return item.hashCode() - otherItem.hashCode();
    }
}
