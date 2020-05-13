package topology

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
 */

import (
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
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

// TOTopology is a type alias on which we can define functions.
type TOTopology struct {
	api.APIInfoImpl `json:"-"`
	tc.Topology
}

// DeleteQueryBase holds a delete query with no WHERE clause and is a
// requirement of the api.GenericOptionsDeleter interface.
func (topology *TOTopology) DeleteQueryBase() string {
	return deleteQueryBase()
}

// ParamColumns maps query parameters to their respective database columns.
func (topology *TOTopology) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"name":        {Column: "t.name"},
		"description": {Column: "t.description"},
		"lastUpdated": {Column: "t.last_updated"},
	}
}

// GenericOptionsDeleter is required by the api.GenericOptionsDeleter interface
// and is called by api.GenericOptionsDelete().
func (topology *TOTopology) DeleteKeyOptions() map[string]dbhelpers.WhereColumnInfo {
	return topology.ParamColumns()
}

func (topology *TOTopology) SetLastUpdated(time tc.TimeNoMod) { topology.LastUpdated = &time }

// GetKeyFieldsInfo is a requirement of the api.Updater interface.
func (topology TOTopology) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"name", api.GetStringKey}}
}

// GetType returns the human-readable type of TOTopology as a string.
func (topology *TOTopology) GetType() string {
	return "topology"
}

// Validate is a requirement of the api.Validator interface.
func (topology *TOTopology) Validate() error {
	nameRule := validation.NewStringRule(tovalidate.IsAlphanumericUnderscoreDash, "must consist of only alphanumeric, dash, or underscore characters.")
	rules := validation.Errors{}
	rules["name"] = validation.Validate(topology.Name, validation.Required, nameRule)

	nodeCount := len(topology.Nodes)
	if nodeCount < 1 {
		rules["length"] = fmt.Errorf("must provide 1 or more node, %v found", nodeCount)
	}
	var (
		cacheGroups      = make([]tc.CacheGroupNullable, nodeCount)
		cacheGroupsExist = true
		err              error
		userErr          error
		sysErr           error
		cacheGroupMap    map[string]tc.CacheGroupNullable
		exists           bool
	)
	_ = err
	cacheGroupNames := make([]string, len(topology.Nodes))
	for index, node := range topology.Nodes {
		rules[fmt.Sprintf("node %v parents size", index)] = validation.Validate(node.Parents, validation.Length(0, 2))
		rules[fmt.Sprintf("node %v duplicate parents", index)] = checkForDuplicateParents(topology.Nodes, index)
		rules[fmt.Sprintf("node %v self parent", index)] = checkForSelfParents(topology.Nodes, index)
		cacheGroupNames[index] = node.Cachegroup
	}
	if cacheGroupMap, userErr, sysErr, _ = cachegroup.GetCacheGroupsByName(cacheGroupNames, topology.APIInfoImpl.ReqInfo.Tx); userErr != nil || sysErr != nil {
		var err error
		message := "could not get cachegroups"
		if userErr != nil {
			err = fmt.Errorf("%s: %s", message, userErr.Error())
		}
		return err
	}
	cacheGroups = make([]tc.CacheGroupNullable, len(topology.Nodes))
	for index, node := range topology.Nodes {
		if cacheGroups[index], exists = cacheGroupMap[node.Cachegroup]; !exists {
			rules[fmt.Sprintf("cachegroup %s not found", node.Cachegroup)] = fmt.Errorf("node %d references nonexistent cachegroup %s", index, node.Cachegroup)
			cacheGroupsExist = false
		}
	}
	rules["duplicate cachegroup name"] = checkUniqueCacheGroupNames(topology.Nodes)

	if cacheGroupsExist {
		for index, node := range topology.Nodes {
			rules[fmt.Sprintf("parent '%v' edge type", node.Cachegroup)] = checkForEdgeParents(topology.Nodes, cacheGroups, index)
		}

		for _, leafMid := range checkForLeafMids(topology.Nodes, cacheGroups) {
			rules[fmt.Sprintf("node %v leaf mid", leafMid.Cachegroup)] = fmt.Errorf("cachegroup %v's type is %v; it cannot be a leaf (it must have at least 1 child)", leafMid.Cachegroup, tc.CacheGroupMidTypeName)
		}
	}
	rules["topology cycles"] = checkForCycles(topology.Nodes)

	errs := tovalidate.ToErrors(rules)
	return util.JoinErrs(errs)
}

// Implementation of the Identifier, Validator interface functions
func (topology TOTopology) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{"name": topology.Name}, true
}

// SetKeys is a requirement of the api.Updater interface and is called by
// api.UpdateHandler().
func (topology *TOTopology) SetKeys(keys map[string]interface{}) {
	topology.Name, _ = keys["name"].(string)
}

// GetAuditName is a requirement of the api.Identifier interface.
func (topology *TOTopology) GetAuditName() string {
	return topology.Name
}

// Create is a requirement of the api.Creator interface.
func (topology *TOTopology) Create() (error, error, int) {
	tx := topology.APIInfo().Tx.Tx
	err := tx.QueryRow(insertQuery(), topology.Name, topology.Description).Scan(&topology.Name, &topology.Description, &topology.LastUpdated)
	if err != nil {
		return api.ParseDBError(err)
	}

	if userErr, sysErr, errCode := topology.addNodes(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	if userErr, sysErr, errCode := topology.addParents(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	return nil, nil, 0
}

// Read is a requirement of the api.Reader interface and is called by api.ReadHandler().
func (topology *TOTopology) Read() ([]interface{}, error, error, int) {
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(topology.ReqInfo.Params, topology.ParamColumns())
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}
	query := selectQuery() + where + orderBy + pagination
	rows, err := topology.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("topology read: querying: " + err.Error()), http.StatusInternalServerError
	}
	defer log.Close(rows, "unable to close DB connection")

	var interfaces []interface{}
	topologies := map[string]*tc.Topology{}
	indices := map[int]int{}
	for index := 0; rows.Next(); index++ {
		var (
			name, description string
			lastUpdated       tc.TimeNoMod
		)
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
			topology := tc.Topology{Nodes: []tc.TopologyNode{}}
			topologies[name] = &topology
			topology.Name = name
			topology.Description = description
			topology.LastUpdated = &lastUpdated
		}
		topologies[name].Nodes = append(topologies[name].Nodes, topologyNode)
	}

	for _, topology := range topologies {
		nodeMap := map[int]int{}
		for index, node := range topology.Nodes {
			nodeMap[node.Id] = index
		}
		for _, node := range topology.Nodes {
			for parentIndex := 0; parentIndex < len(node.Parents); parentIndex++ {
				node.Parents[parentIndex] = nodeMap[node.Parents[parentIndex]]
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
	var cachegroupsToInsert []string
	var indices = make([]int, 0)
	for index, node := range topology.Nodes {
		if node.Id == 0 {
			cachegroupsToInsert = append(cachegroupsToInsert, node.Cachegroup)
			indices = append(indices, index)
		}
	}
	if len(cachegroupsToInsert) == 0 {
		return nil, nil, http.StatusOK
	}
	rows, err := topology.ReqInfo.Tx.Query(nodeInsertQuery(), topology.Name, pq.Array(cachegroupsToInsert))
	if err != nil {
		return nil, errors.New("error adding nodes: " + err.Error()), http.StatusInternalServerError
	}
	defer log.Close(rows, "unable to close DB connection")
	for _, index := range indices {
		rows.Next()
		err = rows.Scan(&topology.Nodes[index].Id, &topology.Name, &topology.Nodes[index].Cachegroup)
		if err != nil {
			return api.ParseDBError(err)
		}
	}
	return nil, nil, http.StatusOK
}

func (topology *TOTopology) addParents() (error, error, int) {
	var (
		children []int
		parents  []int
		ranks    []int
	)
	for _, node := range topology.Nodes {
		for rank := 1; rank <= len(node.Parents); rank++ {
			parent := topology.Nodes[node.Parents[rank-1]]
			children = append(children, node.Id)
			parents = append(parents, parent.Id)
			ranks = append(ranks, rank)
		}
	}
	rows, err := topology.ReqInfo.Tx.Query(nodeParentInsertQuery(), pq.Array(children), pq.Array(parents), pq.Array(ranks))
	if err != nil {
		return api.ParseDBError(err)
	}
	defer log.Close(rows, "unable to close DB connection")
	for _, node := range topology.Nodes {
		for rank := 1; rank <= len(node.Parents); rank++ {
			rows.Next()
			parent := topology.Nodes[node.Parents[rank-1]]
			err = rows.Scan(&node.Id, &parent.Id, &rank)
			if err != nil {
				return api.ParseDBError(err)
			}
		}
	}
	return nil, nil, http.StatusOK
}

func (topology *TOTopology) setDescription() (error, error, int) {
	rows, err := topology.ReqInfo.Tx.Query(updateQuery(), topology.Description, topology.Name)
	if err != nil {
		return nil, fmt.Errorf("topology update: error setting the description for topology %v: %v", topology.Name, err.Error()), http.StatusInternalServerError
	}
	defer log.Close(rows, "unable to close DB connection")
	for rows.Next() {
		err = rows.Scan(&topology.Name, &topology.Description, &topology.LastUpdated)
		if err != nil {
			return api.ParseDBError(err)
		}
	}
	return nil, nil, http.StatusOK
}

// Update is a requirement of the api.Updater interface.
func (topology *TOTopology) Update() (error, error, int) {
	topologies, userErr, sysErr, errCode := topology.Read()
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if len(topologies) != 1 {
		return fmt.Errorf("cannot find exactly 1 topology with the query string provided"), nil, http.StatusBadRequest
	}
	oldTopology := TOTopology{APIInfoImpl: topology.APIInfoImpl, Topology: topologies[0].(tc.Topology)}
	if userErr, sysErr, errCode := topology.setDescription(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	if err := oldTopology.removeParents(); err != nil {
		return nil, err, http.StatusInternalServerError
	}
	var oldNodes, newNodes = map[string]int{}, map[string]int{}
	for index, node := range oldTopology.Nodes {
		oldNodes[node.Cachegroup] = index
	}
	for index, node := range topology.Nodes {
		newNodes[node.Cachegroup] = index
	}
	var toRemove []string
	for cachegroupName := range oldNodes {
		if _, exists := newNodes[cachegroupName]; !exists {
			toRemove = append(toRemove, cachegroupName)
		} else {
			topology.Nodes[newNodes[cachegroupName]].Id = oldTopology.Nodes[oldNodes[cachegroupName]].Id
		}
	}
	if len(toRemove) > 0 {
		if err := oldTopology.removeNodes(&toRemove); err != nil {
			return nil, err, http.StatusInternalServerError
		}
	}
	if userErr, sysErr, errCode := topology.addNodes(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if userErr, sysErr, errCode := topology.addParents(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	return nil, nil, http.StatusOK
}

// Delete is unused and simply satisfies the Deleter interface
// (although TOTOpology is used as an OptionsDeleter)
func (topology *TOTopology) Delete() (error, error, int) {
	return nil, nil, 0
}

// OptionsDelete is a requirement of the OptionsDeleter interface.
func (topology *TOTopology) OptionsDelete() (error, error, int) {
	topologies, userErr, sysErr, errCode := topology.Read()
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if len(topologies) != 1 {
		return fmt.Errorf("cannot find exactly 1 topology with the query string provided"), nil, http.StatusBadRequest
	}
	topology.Topology = topologies[0].(tc.Topology)

	var cachegroups []string
	for _, node := range topology.Nodes {
		cachegroups = append(cachegroups, node.Cachegroup)
	}
	if err := topology.removeNodes(&cachegroups); err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if err := topology.removeParents(); err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return api.GenericOptionsDelete(topology)
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
