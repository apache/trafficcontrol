package com.comcast.cdn.traffic_control.traffic_router.core.util;

import java.util.TreeSet;

public class ComparableTreeSet<E> extends TreeSet<E> implements Comparable<ComparableTreeSet<E>> {

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
