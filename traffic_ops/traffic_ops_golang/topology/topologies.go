package topology

import (
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/cachegroup"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
	"net/http"
)

type TOTopology struct {
	api.APIInfoImpl `json:"-"`
	tc.Topology
}

func (topology *TOTopology) DeleteQueryBase() string {
	return deleteQueryBase()
}

func (topology *TOTopology) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":        dbhelpers.WhereColumnInfo{"t.name", nil},
		"description": dbhelpers.WhereColumnInfo{"t.description", nil},
		"lastUpdated": dbhelpers.WhereColumnInfo{"t.last_updated", nil},
	}
}

func (topology *TOTopology) DeleteKeyOptions() map[string]dbhelpers.WhereColumnInfo {
	return topology.ParamColumns()
}

func (topology *TOTopology) SetLastUpdated(time tc.TimeNoMod) { topology.LastUpdated = &time }

func (topology TOTopology) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"name", api.GetStringKey}}
}

func (topology *TOTopology) GetType() string {
	return "topology"
}

func (topology *TOTopology) Validate() error {
	nameRule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters.")
	rules := validation.Errors{}
	rules["name"] = validation.Validate(topology.Name, validation.Required, nameRule)

	nodeCount := len(*topology.Nodes)
	rules["length"] = validation.Validate(nodeCount, validation.Min(1))
	cacheGroups := make([]*tc.CacheGroupNullable, nodeCount)
	var err error
	for index := 0; index < nodeCount; index++ {
		node := (*topology.Nodes)[index]
		rules[fmt.Sprintf("node %v parents size", index)] = validation.Validate((*node).Parents, validation.Length(0, 2))
		rules[fmt.Sprintf("node %v duplicate parents", index)] = checkForDuplicateParents(topology.Nodes, index)
		if cacheGroups[index], err = cachegroup.GetCacheGroupByName((*node).Cachegroup, &topology.APIInfoImpl); err != nil {
			rules[fmt.Sprintf("node %v parents size", index)] = fmt.Errorf("error getting cachegroup %v: %v", (*node).Cachegroup, err.Error())
		}
	}
	rules["duplicate cachegroup name"] = checkUniqueCacheGroupNames(topology.Nodes)

	for index := 0; index < nodeCount; index++ {
		rules["parents edge type"] = checkForEdgeParents(topology.Nodes, &cacheGroups, index)
	}

	rules["topology cycles"] = checkForCycles(topology.Nodes)

	errs := tovalidate.ToErrors(rules)
	return util.JoinErrs(errs)
}

// Implementation of the Identifier, Validator interface functions
func (topology TOTopology) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{"name": topology.Name}, true
}

func (topology *TOTopology) SetKeys(keys map[string]interface{}) {
	topology.Name, _ = keys["name"].(string)
}

func (topology *TOTopology) GetAuditName() string {
	return topology.Name
}

func (topology *TOTopology) Create() (error, error, int) {
	tx := topology.APIInfo().Tx.Tx
	err := tx.QueryRow(insertQuery(), topology.Name, topology.Description).Scan(&topology.Name, &topology.Description, &topology.LastUpdated)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		return userErr, sysErr, errCode
	}

	if userErr, sysErr, errCode := topology.addNodes(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	if userErr, sysErr, errCode := topology.addParents(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	return nil, nil, 0
}

func (t *TOTopology) Read() ([]interface{}, error, error, int) {
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(t.ReqInfo.Params, t.ParamColumns())
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}
	query := selectQuery() + where + orderBy + pagination
	rows, err := t.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("topology read: querying: " + err.Error()), http.StatusInternalServerError
	}
	defer rows.Close()

	interfaces := []interface{}{}
	topologies := map[string]*tc.Topology{}
	indices := map[int]int{}
	for index := 0; rows.Next(); index++ {
		var name, description string
		var lastUpdated tc.TimeNoMod
		topologyNode := tc.TopologyNode{}
		topologyNode.Parents = []int{}
		var parents pq.Int64Array
		if err = rows.Scan(
			&name,
			&description,
			&lastUpdated,
			&topologyNode.Id,
			&topologyNode.Cachegroup,
			&parents,
		); err != nil {
			return nil, nil, errors.New("topology read: scanning: " + err.Error()), http.StatusInternalServerError
		}
		for _, id := range parents {
			topologyNode.Parents = append(topologyNode.Parents, int(id))
		}
		indices[topologyNode.Id] = index
		if _, exists := topologies[name]; !exists {
			topology := tc.Topology{Nodes: &[]*tc.TopologyNode{}}
			topologies[name] = &topology
			topology.Name = name
			topology.Description = description
			topology.LastUpdated = &lastUpdated
		}
		*topologies[name].Nodes = append(*topologies[name].Nodes, &topologyNode)
	}

	for _, topology := range topologies {
		nodes := topology.Nodes
		nodeCount := len(*topology.Nodes)
		nodeMap := map[int]int{}
		for index := 0; index < nodeCount; index++ {
			nodeMap[(*nodes)[index].Id] = index
		}
		for _, node := range *nodes {
			for parentIndex := 0; parentIndex < len((*node).Parents); parentIndex++ {
				(*node).Parents[parentIndex] = nodeMap[(*node).Parents[parentIndex]]
			}
		}
		interfaces = append(interfaces, *topology)
	}
	return interfaces, nil, nil, http.StatusOK
}

func (topology *TOTopology) removeParents() error {
	_, err := topology.ReqInfo.Tx.Exec(deleteParentsQuery(), topology.Name)
	if err != nil {
		return errors.New("topology update: error deleting old parents: " + err.Error())
	}
	return nil
}

func (topology *TOTopology) removeNodes(cachegroups *[]string) error {
	_, err := topology.ReqInfo.Tx.Exec(deleteNodesQuery(), topology.Name, pq.Array(*cachegroups))
	if err != nil {
		return errors.New("topology update: error removing old unused nodes: " + err.Error())
	}
	return nil
}

func (topology *TOTopology) addNodes() (error, error, int) {
	nodeCount := len(*topology.Nodes)
	cachegroups := make([]string, nodeCount)
	for index := 0; index < nodeCount; index++ {
		cachegroups[index] = (*topology.Nodes)[index].Cachegroup
	}

	rows, err := topology.ReqInfo.Tx.Query(nodeInsertQuery(), topology.Name, pq.Array(cachegroups))
	if err != nil {
		return nil, errors.New("error adding nodes: " + err.Error()), http.StatusInternalServerError
	}
	for _, node := range *topology.Nodes {
		rows.Next()
		err = rows.Scan(&node.Id, &topology.Name, &node.Cachegroup)
		if err != nil {
			userErr, sysErr, errCode := api.ParseDBError(err)
			return userErr, sysErr, errCode
		}
	}
	err = rows.Close()
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		return userErr, sysErr, errCode
	}
	return nil, nil, http.StatusOK
}

func (topology *TOTopology) addParents() (error, error, int) {
	children := []int{}
	parents := []int{}
	ranks := []int{}
	nodeCount := len(*topology.Nodes)
	for index := 0; index < nodeCount; index++ {
		node := (*topology.Nodes)[index]
		for rank := 1; rank <= len(node.Parents); rank++ {
			parent := (*topology.Nodes)[node.Parents[rank-1]]
			children = append(children, node.Id)
			parents = append(parents, parent.Id)
			ranks = append(ranks, rank)
		}
	}
	rows, err := topology.ReqInfo.Tx.Query(nodeParentInsertQuery(), pq.Array(children), pq.Array(parents), pq.Array(ranks))
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		return userErr, sysErr, errCode
	}
	for index := 0; index < nodeCount; index++ {
		node := (*topology.Nodes)[index]
		for rank := 1; rank <= len(node.Parents); rank++ {
			rows.Next()
			parent := (*topology.Nodes)[node.Parents[rank-1]]
			err = rows.Scan(&node.Id, &parent.Id, &rank)
			if err != nil {
				userErr, sysErr, errCode := api.ParseDBError(err)
				return userErr, sysErr, errCode
			}
		}
	}
	err = rows.Close()
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		return userErr, sysErr, errCode
	}
	return nil, nil, http.StatusOK
}

func (topology *TOTopology) setDescription() (error, error, int) {
	rows, err := topology.ReqInfo.Tx.Query(updateQuery(), topology.Description, topology.Name)
	if err != nil {
		return nil, fmt.Errorf("topology update: error setting the description for topology %v: %v", topology.Name, err.Error()), http.StatusInternalServerError
	}
	for rows.Next() {
		err = rows.Scan(&topology.Name, &topology.Description, &topology.LastUpdated)
		if err != nil {
			userErr, sysErr, errCode := api.ParseDBError(err)
			return userErr, sysErr, errCode
		}
	}
	return nil, nil, http.StatusOK
}

func (newTopology *TOTopology) Update() (error, error, int) {
	topologies, userErr, sysErr, errCode := newTopology.Read()
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if len(topologies) != 1 {
		return fmt.Errorf("cannot find exactly 1 topology with the query string provided."), nil, http.StatusBadRequest
	}
	topology := TOTopology{APIInfoImpl: newTopology.APIInfoImpl, Topology: topologies[0].(tc.Topology)}
	if userErr, sysErr, errCode := newTopology.setDescription(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	if err := topology.removeParents(); err != nil {
		return nil, err, http.StatusInternalServerError
	}
	var oldNodes, newNodes = map[string]int{}, map[string]int{}
	var oldNodesLength, newNodesLength = len(*topology.Nodes), len(*newTopology.Nodes)
	for index := 0; index < oldNodesLength; index++ {
		node := (*topology.Nodes)[index]
		oldNodes[(*node).Cachegroup] = index
	}
	for index := 0; index < newNodesLength; index++ {
		node := (*newTopology.Nodes)[index]
		newNodes[(*node).Cachegroup] = index
	}
	var toRemove, toAdd = []string{}, []*tc.TopologyNode{}
	for cachegroupName, _ := range oldNodes {
		if _, exists := newNodes[cachegroupName]; !exists {
			toRemove = append(toRemove, cachegroupName)
		} else {
			(*newTopology.Nodes)[newNodes[cachegroupName]].Id = (*topology.Nodes)[oldNodes[cachegroupName]].Id
		}
	}
	for cachegroupName, index := range newNodes {
		if _, exists := oldNodes[cachegroupName]; !exists {
			toAdd = append(toAdd, (*newTopology.Nodes)[index])
		}
	}
	if err := topology.removeNodes(&toRemove); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	topology.Nodes = &toAdd
	if userErr, sysErr, errCode := topology.addNodes(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	nodeCount := len(*topology.Nodes)
	for index := 0; index < nodeCount; index++ {
		(*newTopology.Nodes)[newNodes[(*topology.Nodes)[index].Cachegroup]] = (*topology.Nodes)[index]
	}

	if userErr, sysErr, errCode := newTopology.addParents(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	return nil, nil, http.StatusOK
}

// Delete is unused and simply satisfies the Deleter interface
// (although TOTOpology is used as an OptionsDeleter)
func (topology *TOTopology) Delete() (error, error, int) {
	return nil, nil, 0
}

func (topology *TOTopology) OptionsDelete() (error, error, int) {
	topologies, userErr, sysErr, errCode := topology.Read()
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if len(topologies) != 1 {
		return fmt.Errorf("cannot find exactly 1 topology with the query string provided."), nil, http.StatusBadRequest
	}
	topology = &TOTopology{APIInfoImpl: topology.APIInfoImpl, Topology: topologies[0].(tc.Topology)}

	var cachegroups []string
	for _, node := range *topology.Nodes {
		cachegroups = append(cachegroups, node.Cachegroup)
	}
	if err := topology.removeNodes(&cachegroups); err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if err := topology.removeParents(); err != nil {
		return nil, err, http.StatusInternalServerError
	}
	api.GenericOptionsDelete(topology)
	return nil, nil, http.StatusOK
}

func insertQuery() string {
	query := `
INSERT INTO topology (name, description)
VALUES ($1, $2)
RETURNING name, description, last_updated
`
	return query
}

func nodeInsertQuery() string {
	query := `
INSERT INTO topology_cachegroup (topology, cachegroup)
VALUES ($1, unnest($2::text[]))
RETURNING id, topology, cachegroup
`
	return query
}

func nodeParentInsertQuery() string {
	query := `
INSERT INTO topology_cachegroup_parents (child, parent, rank)
VALUES (unnest($1::int[]), unnest($2::int[]), unnest($3::int[]))
RETURNING child, parent, rank
`
	return query
}

func selectQuery() string {
	query := `
SELECT t.name, t.description, t.last_updated,
tc.id, tc.cachegroup,
	(SELECT COALESCE (ARRAY_AGG (CAST (tcp.parent as INT) ORDER BY tcp.rank ASC)) AS parents
	FROM topology_cachegroup tc2
	INNER JOIN topology_cachegroup_parents tcp ON tc2.id = tcp.child
	WHERE tc2.topology = tc.topology
	AND tc2.cachegroup = tc.cachegroup
	)
FROM topology t
JOIN topology_cachegroup tc on t.name = tc.topology
`
	return query
}

func deleteQueryBase() string {
	query := `
DELETE FROM topology t
`
	return query
}

func deleteNodesQuery() string {
	query := `
DELETE FROM topology_cachegroup tc
WHERE tc.topology = $1
AND tc.cachegroup = ANY ($2::text[])
`
	return query
}

func deleteParentsQuery() string {
	query := `
DELETE FROM topology_cachegroup_parents tcp
WHERE tcp.child IN
    (SELECT tc.id
    FROM topology t
    JOIN topology_cachegroup tc on t.name = tc.topology
    WHERE t.name = $1)
`
	return query
}

func updateQuery() string {
	query := `
UPDATE topology t SET
description = $1
WHERE t.name = $2
RETURNING t.name, t.description, t.last_updated
`
	return query
}

func nodeUpdateQuery() string {
	query := `
UPDATE topology_cachegroup tc SET
tc.topology = $1, tc.cachegroup = unnest($2::text[])
WHERE tc.id = unnest($3::int[])
RETURNING tc.id, tc.topology, tc.cachegroup
`
	return query
}
