package topology

import (
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
)

func checkUniqueCacheGroupNames(nodes []tc.TopologyNode) error {
	nodeCount := len(nodes)
	cacheGroupNames := map[string]bool{}
	for index := 0; index < nodeCount; index++ {
		if _, exists := cacheGroupNames[nodes[index].Cachegroup]; exists {
			return fmt.Errorf("cachegroup %v cannot be used more than once in the topology.", nodes[index].Cachegroup)
		}
		cacheGroupNames[nodes[index].Cachegroup] = true
	}
	return nil
}

func checkForDuplicateParents(nodes *[]tc.TopologyNode, index int) error {
	parents := (*nodes)[index].Parents
	if len(parents) != 2 || parents[0] != parents[1] {
		return nil
	}
	return fmt.Errorf("cachegroup %v cannot be both a primary and secondary parent of cachegroup %v.", (*nodes)[parents[0]].Cachegroup, (*nodes)[index].Cachegroup)
}

func checkForEdgeParents(nodes *[]tc.TopologyNode, cachegroups *[]*tc.CacheGroupNullable, nodeIndex int) error {
	node := &(*nodes)[nodeIndex]
	parentsLength := len(node.Parents)
	errs := make([]error, parentsLength)
	for parentIndex := 0; parentIndex < parentsLength; parentIndex++ {
		cacheGroupType := (*cachegroups)[node.Parents[parentIndex]].Type
		if *cacheGroupType == tc.EdgeCacheGroupType {
			errs[parentIndex] = fmt.Errorf("cachegroup %v's type is %v; it cannot be a parent of %v.", (*nodes)[parentIndex].Cachegroup, tc.EdgeCacheGroupType, node.Cachegroup)
		}
	}
	return util.JoinErrs(errs)
}
