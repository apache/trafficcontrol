package com.comcast.cdn.traffic_control.traffic_router.core.hash;

import java.util.List;

public interface Hashable <E> {
	Hashable<E> generateHashes(String hashId, int hashCount);
	double getClosestHash(double hash);
	List<Double> getHashValues();
}
