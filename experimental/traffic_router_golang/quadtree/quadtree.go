package quadtree

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
 *
 */

import (
	"math"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

type ObjT tc.CacheGroupName

type DataT struct {
	Lat float64 // top
	Lon float64 // left
	Obj ObjT
}

const elementSize = 1 // MUST be 1 for the search to work properly

type Quadtree struct {
	root *Node
}

type Node struct {
	left        float64
	right       float64
	top         float64
	bottom      float64
	topLeft     *Node
	topRight    *Node
	bottomLeft  *Node
	bottomRight *Node
	elements    []DataT
}

func (n *Node) hasSplit() bool {
	return n.topLeft != nil // it has split, if any child nodes are non-nil
}

// Create Returns a latitude-longitude Quadtree
// Note this creates a latlon quadtree. In Lat Lon, right > left and top > bottom
// TODO add CreateAndInsert func, for optimisation?
func New() *Quadtree {
	return createGeneric(-180.0, 180.0, 90.0, -90.0)
}

func createGeneric(l, r, t, b float64) *Quadtree {
	return &Quadtree{
		root: &Node{
			left:     l,
			right:    r,
			top:      t,
			bottom:   b,
			elements: make([]DataT, 0, elementSize),
		},
	}
}

// Insert inserts the given data into the quadtree
func (q *Quadtree) Insert(d DataT) {
	q.root.insert(d)
}

func (n *Node) insert(d DataT) {
	if !n.hasSplit() {
		n.elements = append(n.elements, d)
		n.split()
		return
	}

	n.insertAppropriateChild(d)
}

func (n *Node) insertAppropriateChild(d DataT) {
	leftMid := (n.right-n.left)/2 + n.left
	topMid := (n.top-n.bottom)/2 + n.bottom
	left := d.Lon < leftMid
	top := d.Lat > topMid
	if left && top {
		n.topLeft.insert(d)
	} else if top {
		n.topRight.insert(d)
	} else if left {
		n.bottomLeft.insert(d)
	} else {
		n.bottomRight.insert(d)
	}
}

// split splits the node, IFF len(elements) > size. Does NOT split if splitting is unnecessary.
func (n *Node) split() {
	if len(n.elements) <= elementSize {
		return
	}
	if len(n.elements) > 1 && n.elements[len(n.elements)-2].Lat == n.elements[len(n.elements)-1].Lat && n.elements[len(n.elements)-2].Lon == n.elements[len(n.elements)-1].Lon {
		// Allow multiple objects at the same location - the first one inserted will be returned. TODO warn?
		return
	}

	leftMid := (n.right-n.left)/2 + n.left
	topMid := (n.top-n.bottom)/2 + n.bottom
	n.topLeft = &Node{left: n.left, right: leftMid, top: n.top, bottom: topMid}
	n.topRight = &Node{left: leftMid, right: n.right, top: n.top, bottom: topMid}
	n.bottomLeft = &Node{left: n.left, right: leftMid, top: topMid, bottom: n.bottom}
	n.bottomRight = &Node{left: leftMid, right: n.right, top: topMid, bottom: n.bottom}
	for _, e := range n.elements {
		n.insertAppropriateChild(e)
	}
	n.elements = nil
}

func (q *Quadtree) Get(top float64, left float64, bottom float64, right float64) []DataT {
	return q.root.get(top, left, bottom, right)
}

// get returns the data in the given rect
func (n *Node) get(top float64, left float64, bottom float64, right float64) []DataT {
	if top < n.bottom || bottom > n.top || left > n.right || right < n.left {
		return nil
	}

	if !n.hasSplit() {
		return elementsInRect(top, left, bottom, right, n.elements)
	}

	topLeftElements := n.topLeft.get(top, left, bottom, right)
	topRightElements := n.topRight.get(top, left, bottom, right)
	bottomLeftElements := n.bottomLeft.get(top, left, bottom, right)
	bottomRightElements := n.bottomRight.get(top, left, bottom, right)

	all := make([]DataT, 0, len(topLeftElements)+len(topRightElements)+len(bottomLeftElements)+len(bottomRightElements))
	all = append(all, topLeftElements...)
	all = append(all, topRightElements...)
	all = append(all, bottomLeftElements...)
	all = append(all, bottomRightElements...)
	return all
}

func inRect(rtop float64, rleft float64, rbottom float64, rright float64, lat float64, lon float64) bool {
	return lat < rtop && lat > rbottom && lon > rleft && lon < rright
}

func elementsInRect(rtop float64, rleft float64, rbottom float64, rright float64, elements []DataT) []DataT {
	in := []DataT{}
	for _, e := range elements {
		if inRect(rtop, rleft, rbottom, rright, e.Lat, e.Lon) {
			in = append(in, e)
		}
	}
	return in
}

// Nearest returns the nearest object to the given coordinates, and whether an element was found. Returns false iff the tree is empty.
// TODO add NearestN - n nearest objs, e.g. in case no cache in the nearest cachegroup is available.
// TODO add NearestAll - all nearest objs with the same coords, e.g. so multiple cachegroups at the same coordinate can be round-robined
func (q *Quadtree) Nearest(top float64, left float64) (DataT, bool) {
	return q.root.nearest(top, left)
}

func (n *Node) nearest(top float64, left float64) (DataT, bool) {
	candidates := []DataT{}
	if !n.hasSplit() {
		candidates = n.elements
	} else {
		leftMid := (n.right-n.left)/2 + n.left
		topMid := (n.top-n.bottom)/2 + n.bottom
		inLeft := left < leftMid
		inTop := top > topMid
		if inLeft && inTop {
			if nearest, ok := n.topLeft.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
		} else if inTop {
			if nearest, ok := n.topRight.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
		} else if inLeft {
			if nearest, ok := n.bottomLeft.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
		} else {
			if nearest, ok := n.bottomRight.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
		}
		if len(candidates) == 0 {
			// if the quadrant with the given point is empty, check all quadrants
			if nearest, ok := n.topLeft.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
			if nearest, ok := n.topRight.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
			if nearest, ok := n.bottomLeft.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
			if nearest, ok := n.bottomRight.nearest(top, left); ok {
				candidates = append(candidates, nearest)
			}
		}
	}

	if len(candidates) == 0 {
		return DataT{}, false
	}

	best := DataT{}
	bestDist := math.MaxFloat64 // TODO memoize by returning?
	for _, candidate := range candidates {
		if d := distSq(top, left, candidate.Lat, candidate.Lon); d < bestDist {
			bestDist = d
			best = candidate
		}
	}
	return best, true
}

// distSq returns the distance squared between two points. Because calculating the root is expensive, this is useful as an optimisation to compare distances, when the actual distance isn't necessary.
func distSq(atop, aleft, btop, bleft float64) float64 {
	return (btop-atop)*(btop-atop) + (bleft-aleft)*(bleft-aleft)
}
