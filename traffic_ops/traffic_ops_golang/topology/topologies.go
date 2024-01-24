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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/cachegroup"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/topology/topology_validation"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/util/ims"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

// TOTopology is a type alias on which we can define functions.
type TOTopology struct {
	api.APIInfoImpl `json:"-"`
	Alerts          tc.Alerts `json:"-"`
	RequestedName   string    `json:"-"`
	tc.Topology
}

// GetAlerts implements the AlertsResponse interface.
func (topology *TOTopology) GetAlerts() tc.Alerts {
	return topology.Alerts
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
	return []api.KeyFieldInfo{{Field: "name", Func: api.GetStringKey}}
}

// GetType returns the human-readable type of TOTopology as a string.
func (topology *TOTopology) GetType() string {
	return "topology"
}

// DowngradeTopologyNodes downgrades v5 topology nodes into legacy topology node structures.
func DowngradeTopologyNodes(nodes []tc.TopologyNodeV5) []tc.TopologyNode {
	legacyNodes := make([]tc.TopologyNode, len(nodes))
	for i, n := range nodes {
		var legacyNode tc.TopologyNode
		legacyNode.Id = n.Id
		legacyNode.Cachegroup = n.Cachegroup
		legacyNode.Parents = make([]int, len(n.Parents))
		for i, p := range n.Parents {
			legacyNode.Parents[i] = p
		}
		if n.LastUpdated != nil {
			legacyNode.LastUpdated = tc.TimeNoModFromTime(*n.LastUpdated)
		}
		legacyNodes[i] = legacyNode
	}
	return legacyNodes
}

// ValidateTopology validates a v5 topology to make sure that the supplied fields are valid.
func ValidateTopology(topology tc.TopologyV5, reqInfo *api.Info) (tc.Alerts, error, error) {
	var alertsObject tc.Alerts
	currentTopoName := reqInfo.Params["name"]
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

	legacyNodes := DowngradeTopologyNodes(topology.Nodes)
	cacheGroupNames := make([]string, len(topology.Nodes))
	for index, node := range topology.Nodes {
		rules[fmt.Sprintf("node %v parents size", index)] = validation.Validate(node.Parents, validation.Length(0, 2))
		rules[fmt.Sprintf("node %v duplicate parents", index)] = checkForDuplicateParents(legacyNodes, index)
		rules[fmt.Sprintf("node %v self parent", index)] = checkForSelfParents(legacyNodes, index)
		cacheGroupNames[index] = node.Cachegroup
	}
	if cacheGroupMap, userErr, sysErr, _ = cachegroup.GetCacheGroupsByName(cacheGroupNames, reqInfo.Tx); userErr != nil || sysErr != nil {
		return alertsObject, userErr, sysErr
	}
	cacheGroups = make([]tc.CacheGroupNullable, len(topology.Nodes))
	for index, node := range topology.Nodes {
		if cacheGroups[index], exists = cacheGroupMap[node.Cachegroup]; !exists {
			rules[fmt.Sprintf("cachegroup %s not found", node.Cachegroup)] = fmt.Errorf("node %d references nonexistent cachegroup %s", index, node.Cachegroup)
			cacheGroupsExist = false
		}
	}
	rules["duplicate cachegroup name"] = checkUniqueCacheGroupNames(legacyNodes)
	if !cacheGroupsExist {
		return alertsObject, util.JoinErrs(tovalidate.ToErrors(rules)), nil
	}

	for index, node := range topology.Nodes {
		alerts, err := checkForEdgeParents(topology, cacheGroups, index)
		if len(alerts.Alerts) != 0 {
			alertsObject = alerts
		}
		rules[fmt.Sprintf("parent '%v' edge type", node.Cachegroup)] = err
	}

	cacheGroupIds := make([]int, len(cacheGroupNames))
	for index, cacheGroup := range cacheGroups {
		cacheGroupIds[index] = *cacheGroup.ID
	}
	dsCDNs, err := dbhelpers.GetDeliveryServiceCDNsByTopology(reqInfo.Tx.Tx, currentTopoName)
	if err != nil {
		return alertsObject, errors.New("unable to validate topology"), fmt.Errorf("validating Topology: %w", err)
	}
	rules["empty cachegroups"] = topology_validation.CheckForEmptyCacheGroups(reqInfo.Tx, cacheGroupIds, dsCDNs, false, nil)
	//Get current Topology-CG for the requested change.
	topoCachegroupNames := getCachegroupNames(legacyNodes)
	rules["required capabilities"] = validateDSRequiredCapabilities(reqInfo.Tx.Tx, currentTopoName, topoCachegroupNames)

	userErr, sysErr, _ = dbhelpers.CheckTopologyOrgServerCGInDSCG(reqInfo.Tx.Tx, dsCDNs, currentTopoName, topoCachegroupNames)
	if userErr != nil || sysErr != nil {
		return alertsObject, userErr, sysErr
	}

	/* Only perform further checks if everything so far is valid */
	if err = util.JoinErrs(tovalidate.ToErrors(rules)); err != nil {
		return alertsObject, err, nil
	}

	for _, leafMid := range checkForLeafMids(legacyNodes, cacheGroups) {
		rules[fmt.Sprintf("node %v leaf mid", leafMid.Cachegroup)] = fmt.Errorf("cachegroup %v's type is %v; it cannot be a leaf (it must have at least 1 primary or secondary child cache group)", leafMid.Cachegroup, tc.CacheGroupMidTypeName)
	}
	_, rules["topology cycles"] = checkForCycles(legacyNodes)
	rules["super-topology cycles"] = checkForCyclesAcrossTopologies(reqInfo, legacyNodes, topology.Name)

	return alertsObject, util.JoinErrs(tovalidate.ToErrors(rules)), nil
}

// Validate is a requirement of the api.Validator interface.
func (topology *TOTopology) Validate() (error, error) {
	currentTopoName := topology.APIInfoImpl.ReqInfo.Params["name"]
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
		return userErr, sysErr
	}
	cacheGroups = make([]tc.CacheGroupNullable, len(topology.Nodes))
	for index, node := range topology.Nodes {
		if cacheGroups[index], exists = cacheGroupMap[node.Cachegroup]; !exists {
			rules[fmt.Sprintf("cachegroup %s not found", node.Cachegroup)] = fmt.Errorf("node %d references nonexistent cachegroup %s", index, node.Cachegroup)
			cacheGroupsExist = false
		}
	}
	rules["duplicate cachegroup name"] = checkUniqueCacheGroupNames(topology.Nodes)
	if !cacheGroupsExist {
		return util.JoinErrs(tovalidate.ToErrors(rules)), nil
	}

	for index, node := range topology.Nodes {
		rules[fmt.Sprintf("parent '%v' edge type", node.Cachegroup)] = topology.checkForEdgeParents(cacheGroups, index)
	}

	cacheGroupIds := make([]int, len(cacheGroupNames))
	for index, cacheGroup := range cacheGroups {
		cacheGroupIds[index] = *cacheGroup.ID
	}
	dsCDNs, err := dbhelpers.GetDeliveryServiceCDNsByTopology(topology.ReqInfo.Tx.Tx, currentTopoName)
	if err != nil {
		return errors.New("unable to validate topology"), fmt.Errorf("validating Topology: %w", err)
	}
	rules["empty cachegroups"] = topology_validation.CheckForEmptyCacheGroups(topology.ReqInfo.Tx, cacheGroupIds, dsCDNs, false, nil)

	//Get current Topology-CG for the requested change.
	topoCachegroupNames := getCachegroupNames(topology.Nodes)
	rules["required capabilities"] = validateDSRequiredCapabilities(topology.APIInfo().Tx.Tx, currentTopoName, topoCachegroupNames)

	userErr, sysErr, _ = dbhelpers.CheckTopologyOrgServerCGInDSCG(topology.ReqInfo.Tx.Tx, dsCDNs, currentTopoName, topoCachegroupNames)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr
	}

	/* Only perform further checks if everything so far is valid */
	if err = util.JoinErrs(tovalidate.ToErrors(rules)); err != nil {
		return err, nil
	}

	for _, leafMid := range checkForLeafMids(topology.Nodes, cacheGroups) {
		rules[fmt.Sprintf("node %v leaf mid", leafMid.Cachegroup)] = fmt.Errorf("cachegroup %v's type is %v; it cannot be a leaf (it must have at least 1 primary or secondary child cache group)", leafMid.Cachegroup, tc.CacheGroupMidTypeName)
	}
	_, rules["topology cycles"] = checkForCycles(topology.Nodes)
	rules["super-topology cycles"] = checkForCyclesAcrossTopologies(topology.APIInfo(), topology.Nodes, topology.Name)

	return util.JoinErrs(tovalidate.ToErrors(rules)), nil
}

func nodesInOtherTopologies(info *api.Info, topologyNodes []tc.TopologyNode) ([]tc.TopologyNode, map[string][]string, error) {
	currentTopoName := info.Params["name"]
	baseError := errors.New("unable to verify that there are no cycles across all topologies")
	where := `WHERE name != :topology_name`
	query := selectQueryWithParentNames() + where + `
		UNION ` + selectNonTopologyCacheGroupsQuery() + `
		UNION ` + selectNonTopologyParentCacheGroupsQuery()

	parameters := map[string]interface{}{
		"topology_name":    currentTopoName,
		"edge_type_prefix": strings.ToLower(tc.EdgeTypePrefix) + "%",
		"mid_type_prefix":  strings.ToLower(tc.MidTypePrefix) + "%",
	}
	rows, err := info.Tx.NamedQuery(query, parameters)
	if err != nil {
		return nil, nil, baseError
	}
	defer log.Close(rows, "unable to close DB connection")

	parentMapMap := map[string]map[string]bool{}
	topologiesMapByCacheGroup := map[string]map[string]bool{}
	for index := 0; rows.Next(); index++ {
		var (
			topologyName string
			cacheGroup   string
			parents      []string
		)
		topologyNode := tc.TopologyNode{}
		topologyNode.Parents = []int{}
		if err = rows.Scan(
			&topologyName,
			&cacheGroup,
			pq.Array(&parents),
		); err != nil {
			return nil, nil, baseError
		}
		if _, exists := parentMapMap[cacheGroup]; !exists {
			parentMapMap[cacheGroup] = map[string]bool{}
		}
		if _, exists := topologiesMapByCacheGroup[cacheGroup]; !exists {
			topologiesMapByCacheGroup[cacheGroup] = map[string]bool{}
		}
		for _, parent := range parents {
			parentMapMap[cacheGroup][parent] = true
		}
		topologiesMapByCacheGroup[cacheGroup][topologyName] = true
	}

	topologiesByCacheGroup := map[string][]string{}
	// Build the list of topologies containing each cache group
	for cacheGroup, topologiesMap := range topologiesMapByCacheGroup {
		var topologies []string
		for topology, _ := range topologiesMap {
			topologies = append(topologies, topology)
		}
		topologiesByCacheGroup[cacheGroup] = topologies
	}

	// Add nodes for the topology we are validating
	for _, node := range topologyNodes {
		if _, exists := parentMapMap[node.Cachegroup]; !exists {
			parentMapMap[node.Cachegroup] = map[string]bool{}
		}
		for _, parentCacheGroupIndex := range node.Parents {
			parentMapMap[node.Cachegroup][topologyNodes[parentCacheGroupIndex].Cachegroup] = true
		}
	}

	indexByCachegroup := map[string]int{}
	var cacheGroups []string
	index := 0
	// Get an index for each cachegroup
	for cacheGroup, _ := range parentMapMap {
		cacheGroups = append(cacheGroups, cacheGroup)
		indexByCachegroup[cacheGroup] = index
		index++
	}

	var nodes []tc.TopologyNode
	// Reduce parentMapMap to an array of TopologyNodes
	// We can't rely on iterating through parentMapMap in the same order twice, so we iterate through
	// cacheGroups instead
	for _, cacheGroup := range cacheGroups {
		parentMap := parentMapMap[cacheGroup]
		node := tc.TopologyNode{Cachegroup: cacheGroup}
		for parent, _ := range parentMap {
			node.Parents = append(node.Parents, indexByCachegroup[parent])
		}
		nodes = append(nodes, node)
	}

	return nodes, topologiesByCacheGroup, nil
}

func validateDSRequiredCapabilities(tx *sql.Tx, currentTopoName string, cachegroups []string) error {
	baseError := errors.New("unable to verify that delivery service required capabilities are satisfied")
	dsRequiredCapabilities, dsCDNs, err := getDSRequiredCapabilitiesByTopology(currentTopoName, tx)
	if err != nil {
		log.Errorf("validating delivery service required capabilities for topology %s: %v", currentTopoName, err)
		return baseError
	}
	if len(dsRequiredCapabilities) == 0 {
		return nil
	}
	cdnMap := make(map[int]struct{})
	for _, cdn := range dsCDNs {
		cdnMap[cdn] = struct{}{}
	}
	CDNs := []int{}
	for cdn := range cdnMap {
		CDNs = append(CDNs, cdn)
	}
	q := `
SELECT
  s.id,
  s.cdn_id,
  c.name,
  ARRAY_REMOVE(ARRAY_AGG(ssc.server_capability ORDER BY ssc.server_capability), NULL) AS capabilities
FROM server s
LEFT JOIN server_server_capability ssc ON ssc.server = s.id
JOIN cachegroup c ON c.id = s.cachegroup
WHERE
  c.name = ANY($1)
  AND s.cdn_id = ANY($2)
  AND c.type != (SELECT id FROM type WHERE name = '` + tc.CacheGroupOriginTypeName + `')
GROUP BY s.id, s.cdn_id, c.name
`
	rows, err := tx.Query(q, pq.Array(cachegroups), pq.Array(CDNs))
	if err != nil {
		log.Errorf("querying server capabilities in topology.validateDSRequiredCapabilities: %v", err)
		return baseError
	}
	cachegroupServers, serverCapabilities, serverCDNs, err := dbhelpers.ScanCachegroupsServerCapabilities(rows)
	if err != nil {
		log.Errorf("validating delivery service required capabilities for topology %s: %v", currentTopoName, err)
		return baseError
	}

	cdnCachegroupServers := make(map[int]map[string][]int)
	for _, cdn := range dsCDNs {
		if _, ok := cdnCachegroupServers[cdn]; !ok {
			cdnCachegroupServers[cdn] = make(map[string][]int)
		}
	}
	for cg, servers := range cachegroupServers {
		for _, s := range servers {
			cdnCachegroupServers[serverCDNs[s]][cg] = append(cdnCachegroupServers[serverCDNs[s]][cg], s)
		}
	}

	invalidDSes := []string{}
	for ds, dsReqCaps := range dsRequiredCapabilities {
		invalidCachegroups := deliveryservice.GetInvalidCachegroupsForRequiredCapabilities(cdnCachegroupServers[dsCDNs[ds]], serverCapabilities, dsReqCaps)
		if len(invalidCachegroups) > 0 {
			invalidDSes = append(invalidDSes, fmt.Sprintf("%s: cachegroups [%s] do not meet required capabilities", ds, strings.Join(invalidCachegroups, ", ")))
		}
	}
	if len(invalidDSes) > 0 {
		return errors.New("cannot update topology. The following delivery services would not be satisfied: " + strings.Join(invalidDSes, "; "))
	}

	return nil
}

// getDSRequiredCapabilitiesByTopology returns a map of DS xml_id to required capabilities,
// a map of xml_id to cdn_id, and an error (if one occurs).
func getDSRequiredCapabilitiesByTopology(name string, tx *sql.Tx) (map[string][]string, map[string]int, error) {
	q := `
SELECT
  d.xml_id,
  d.cdn_id,
  d.required_capabilities
FROM deliveryservice d
WHERE
  d.topology = $1
GROUP BY d.xml_id, d.cdn_id, d.required_capabilities
`
	rows, err := tx.Query(q, name)
	if err != nil {
		return nil, nil, fmt.Errorf("querying delivery service required capabilities by topology: %v", err)
	}
	defer log.Close(rows, "closing rows in getDSRequiredCapabilitiesByTopology")

	requiredCapabilities := make(map[string][]string)
	dsCdnIDs := make(map[string]int)
	for rows.Next() {
		xmlID := ""
		cdnID := 0
		reqCaps := []string{}
		if err := rows.Scan(&xmlID, &cdnID, pq.Array(&reqCaps)); err != nil {
			return nil, nil, fmt.Errorf("scanning delivery service required capabilities by topology: %v", err)
		}
		requiredCapabilities[xmlID] = reqCaps
		dsCdnIDs[xmlID] = cdnID
	}
	return requiredCapabilities, dsCdnIDs, nil
}

func getCachegroupNames(nodes []tc.TopologyNode) []string {
	cgSet := make(map[string]struct{})
	for _, n := range nodes {
		cgSet[n.Cachegroup] = struct{}{}
	}
	cachegroups := make([]string, 0, len(cgSet))
	for c := range cgSet {
		cachegroups = append(cachegroups, c)
	}
	return cachegroups
}

// Implementation of the Identifier, Validator interface functions
func (topology TOTopology) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{"name": topology.Name}, true
}

// SetKeys is a requirement of the api.Updater interface and is called by
// api.UpdateHandler().
func (topology *TOTopology) SetKeys(keys map[string]interface{}) {
	topology.RequestedName = topology.Name
	topology.Name, _ = keys["name"].(string)
}

// GetAuditName is a requirement of the api.Identifier interface.
func (topology *TOTopology) GetAuditName() string {
	return topology.Name
}

// Create is the handler for creating a new topology entity in api version 5.0.
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx
	var topology tc.TopologyV5
	if err := json.NewDecoder(r.Body).Decode(&topology); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	alerts, userErr, sysErr := ValidateTopology(topology, inf)
	if userErr != nil || sysErr != nil {
		code := http.StatusBadRequest
		if sysErr != nil {
			code = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, code, userErr, sysErr)
		return
	}

	legacyNodes := DowngradeTopologyNodes(topology.Nodes)
	userErr, sysErr, statusCode := checkIfTopologyCanBeAlteredByCurrentUser(inf, legacyNodes)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}
	err := tx.QueryRow(insertQuery(), topology.Name, topology.Description).Scan(&topology.Name, &topology.Description, &topology.LastUpdated)
	if err != nil {
		userErr, sysErr, statusCode = api.ParseDBError(err)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}
	topology.LastUpdated, err = util.ConvertTimeFormat(*topology.LastUpdated, time.RFC3339)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("could not convert last updated time into rfc3339 format: %s", err.Error()))
		return
	}

	if userErr, sysErr, statusCode = addNodes(tx, topology.Name, &topology); userErr != nil || sysErr != nil {
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	if userErr, sysErr, statusCode = addParents(tx, topology.Nodes); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}

	alertsObject := tc.CreateAlerts(tc.SuccessLevel, "topology was created.")
	if len(alerts.Alerts) != 0 {
		alertsObject.AddAlerts(alerts)
	}
	api.WriteAlertsObj(w, r, http.StatusOK, alertsObject, topology)

	changeLogMsg := fmt.Sprintf("Topology: %s, ACTION: Created topology, keys: {name: %s }", topology.Name, topology.Name)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Create is a requirement of the api.Creator interface.
func (topology *TOTopology) Create() (error, error, int) {
	tx := topology.APIInfo().Tx.Tx
	userErr, sysErr, statusCode := checkIfTopologyCanBeAlteredByCurrentUser(topology.APIInfo(), topology.Nodes)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
	}
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

func readTopologies(r *http.Request, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, &maxTime
	}
	defer inf.Close()

	interfaces := make([]interface{}, 0)

	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"name":        {Column: "t.name"},
		"description": {Column: "t.description"},
		"lastUpdated": {Column: "t.last_updated"},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, &maxTime
	}

	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(inf.Tx, r.Header, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return []interface{}{}, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}

	// Case where we need to run the second query
	query := selectQuery() + where + orderBy + pagination
	rows, err := inf.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("topology read: querying: " + err.Error()), http.StatusInternalServerError, &maxTime
	}
	defer log.Close(rows, "unable to close DB connection")

	topologies := map[string]*tc.TopologyV5{}
	indices := map[int]int{}
	for index := 0; rows.Next(); index++ {
		var (
			name, description string
			lastUpdated       time.Time
		)
		topologyNode := tc.TopologyNodeV5{}
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
			return nil, nil, errors.New("topology read: scanning: " + err.Error()), http.StatusInternalServerError, &maxTime
		}
		for _, id := range parents {
			topologyNode.Parents = append(topologyNode.Parents, int(id))
		}
		indices[topologyNode.Id] = index
		if _, exists := topologies[name]; !exists {
			topology := tc.TopologyV5{Nodes: []tc.TopologyNodeV5{}}
			topologies[name] = &topology
			topology.Name = name
			topology.Description = description
			topology.LastUpdated, err = util.ConvertTimeFormat(lastUpdated, time.RFC3339)
			if err != nil {
				return nil, nil, errors.New("couldn't convert last updated time to rfc3339 format: " + err.Error()), http.StatusInternalServerError, &maxTime
			}
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
	return interfaces, nil, nil, http.StatusOK, &maxTime
}

// Read is the handler for reading topologies in api version 5.0.
func Read(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, statusCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, nil, statusCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	interfaces, userErr, sysErr, statusCode, maxTime := readTopologies(r, inf.Config.UseIMS)
	if statusCode == http.StatusNotModified {
		api.AddLastModifiedHdr(w, *maxTime)
		w.WriteHeader(http.StatusNotModified)
		return
	}
	if api.SetLastModifiedHeader(r, inf.Config.UseIMS) {
		date := maxTime.Format(rfc.LastModifiedFormat)
		w.Header().Add(rfc.LastModified, date)
	}
	api.WriteResp(w, r, interfaces)
	return
}

// Delete is the handler for removing a topology entity in api version 5.0.
func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, statusCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	topologies, userErr, sysErr, statusCode, _ := readTopologies(r, false)
	if len(topologies) != 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("cannot find exactly 1 topology with the query string provided"), nil)
		return
	}

	topology := topologies[0].(tc.TopologyV5)
	nodes := DowngradeTopologyNodes(topology.Nodes)
	userErr, sysErr, statusCode = checkIfTopologyCanBeAlteredByCurrentUser(inf, nodes)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	where, _, _, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(inf.Params, map[string]dbhelpers.WhereColumnInfo{
		"name":        {Column: "t.name"},
		"description": {Column: "t.description"},
		"lastUpdated": {Column: "t.last_updated"},
	})
	if len(errs) > 0 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, util.JoinErrs(errs), nil)
		return
	}

	query := deleteQueryBase() + where
	result, err := inf.Tx.NamedExec(query, queryValues)
	if err != nil {
		userErr, sysErr, statusCode = api.ParseDBError(err)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
			return
		}
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("deleting topology: getting rows affected: "+err.Error()))
		return
	} else if rowsAffected < 1 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusNotFound, errors.New("no topology with that key found"), nil)
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, fmt.Errorf("topology delete affected too many rows: %d", rowsAffected))
		return
	}

	alertsObject := tc.CreateAlerts(tc.SuccessLevel, "topology was deleted.")
	api.WriteAlertsObj(w, r, http.StatusOK, alertsObject, topology)

	changeLogMsg := fmt.Sprintf("TOPOLOGY: %s, ACTION: Deleted topology, keys: {name: %s }", topology.Name, topology.Name)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Update is the handler for modifying a topology entity in api version 5.0.
func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, statusCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, statusCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	tx := inf.Tx.Tx

	topologies, userErr, sysErr, statusCode, _ := readTopologies(r, false)
	if len(topologies) != 1 {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("cannot find exactly 1 topology with the query string provided"), nil)
		return
	}
	var topology tc.TopologyV5
	if err := json.NewDecoder(r.Body).Decode(&topology); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	alerts, userErr, sysErr := ValidateTopology(topology, inf)
	if userErr != nil || sysErr != nil {
		code := http.StatusBadRequest
		if sysErr != nil {
			code = http.StatusInternalServerError
		}
		api.HandleErr(w, r, inf.Tx.Tx, code, userErr, sysErr)
		return
	}

	requestedName := inf.Params["name"]
	// check if the entity was already updated
	userErr, sysErr, statusCode = api.CheckIfUnModifiedByName(r.Header, inf.Tx, requestedName, "topology")
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}
	nodes := DowngradeTopologyNodes(topology.Nodes)
	userErr, sysErr, statusCode = checkIfTopologyCanBeAlteredByCurrentUser(inf, nodes)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	oldTopology := topologies[0].(tc.TopologyV5)

	if err := removeParents(tx, oldTopology.Name); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
		return
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
		if err := removeNodes(inf.Tx.Tx, oldTopology.Name, &toRemove); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, err)
			return
		}
	}

	if userErr, sysErr, statusCode = setTopologyDetails(tx, &topology, oldTopology.Name); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	if userErr, sysErr, statusCode = addNodes(tx, topology.Name, &topology); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}
	if userErr, sysErr, statusCode = addParents(tx, topology.Nodes); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, statusCode, userErr, sysErr)
		return
	}

	alertsObject := tc.CreateAlerts(tc.SuccessLevel, "topology was updated.")
	if len(alerts.Alerts) != 0 {
		alertsObject.AddAlerts(alerts)
	}
	api.WriteAlertsObj(w, r, http.StatusOK, alertsObject, topology)

	changeLogMsg := fmt.Sprintf("TOPOLOGY: %s, ACTION: Updated topology, keys: {name: %s }", topology.Name, topology.Name)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Read is a requirement of the api.Reader interface and is called by api.ReadHandler().
func (topology *TOTopology) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	interfaces := make([]interface{}, 0)
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(topology.ReqInfo.Params, topology.ParamColumns())
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}
	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(topology.ReqInfo.Tx, h, queryValues, selectMaxLastUpdatedQuery(where))
		if !runSecond {
			log.Debugln("IMS HIT")
			return interfaces, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	// Case where we need to run the second query
	query := selectQuery() + where + orderBy + pagination
	rows, err := topology.ReqInfo.Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, errors.New("topology read: querying: " + err.Error()), http.StatusInternalServerError, nil
	}
	defer log.Close(rows, "unable to close DB connection")
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
			return nil, nil, errors.New("topology read: scanning: " + err.Error()), http.StatusInternalServerError, nil
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
	return interfaces, nil, nil, http.StatusOK, &maxTime
}

func removeParents(tx *sql.Tx, name string) error {
	_, err := tx.Exec(deleteParentsQuery(), name)
	if err != nil {
		return errors.New("topology update: error deleting old parents: " + err.Error())
	}
	return nil
}

func removeNodes(tx *sql.Tx, name string, cachegroups *[]string) error {
	_, err := tx.Exec(deleteNodesQuery(), name, pq.Array(*cachegroups))
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

func addNodes(tx *sql.Tx, name string, topology *tc.TopologyV5) (error, error, int) {
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
	rows, err := tx.Query(nodeInsertQuery(), name, pq.Array(cachegroupsToInsert))
	if err != nil {
		return nil, errors.New("error adding nodes: " + err.Error()), http.StatusInternalServerError
	}
	defer log.Close(rows, "unable to close DB connection")
	for _, index := range indices {
		rows.Next()
		err = rows.Scan(&topology.Nodes[index].Id, &name, &topology.Nodes[index].Cachegroup)
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

func addParents(tx *sql.Tx, nodes []tc.TopologyNodeV5) (error, error, int) {
	var (
		children []int
		parents  []int
		ranks    []int
	)
	for _, node := range nodes {
		for rank := 1; rank <= len(node.Parents); rank++ {
			parent := nodes[node.Parents[rank-1]]
			children = append(children, node.Id)
			parents = append(parents, parent.Id)
			ranks = append(ranks, rank)
		}
	}
	rows, err := tx.Query(nodeParentInsertQuery(), pq.Array(children), pq.Array(parents), pq.Array(ranks))
	if err != nil {
		return api.ParseDBError(err)
	}
	defer log.Close(rows, "unable to close DB connection")
	for _, node := range nodes {
		for rank := 1; rank <= len(node.Parents); rank++ {
			rows.Next()
			parent := nodes[node.Parents[rank-1]]
			err = rows.Scan(&node.Id, &parent.Id, &rank)
			if err != nil {
				return api.ParseDBError(err)
			}
		}
	}
	return nil, nil, http.StatusOK
}

func setTopologyDetails(tx *sql.Tx, topology *tc.TopologyV5, oldTopologyName string) (error, error, int) {
	rows, err := tx.Query(updateQuery(), topology.Name, topology.Description, oldTopologyName)
	if err != nil {
		return nil, fmt.Errorf("topology update: error setting the name and/or description for topology %v: %v", topology.Name, err.Error()), http.StatusInternalServerError
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

func (topology *TOTopology) setTopologyDetails() (error, error, int) {
	rows, err := topology.ReqInfo.Tx.Query(updateQuery(), topology.RequestedName, topology.Description, topology.Name)
	if err != nil {
		return nil, fmt.Errorf("topology update: error setting the name and/or description for topology %v: %v", topology.Name, err.Error()), http.StatusInternalServerError
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
func (topology *TOTopology) Update(h http.Header) (error, error, int) {
	topologies, userErr, sysErr, errCode, _ := topology.Read(h, false)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if len(topologies) != 1 {
		return fmt.Errorf("cannot find exactly 1 topology with the query string provided"), nil, http.StatusBadRequest
	}

	// check if the entity was already updated
	userErr, sysErr, errCode = api.CheckIfUnModifiedByName(h, topology.ReqInfo.Tx, topology.Name, "topology")
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	userErr, sysErr, statusCode := checkIfTopologyCanBeAlteredByCurrentUser(topology.APIInfo(), topology.Nodes)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
	}
	oldTopology := TOTopology{APIInfoImpl: topology.APIInfoImpl, Topology: topologies[0].(tc.Topology)}

	if err := removeParents(oldTopology.APIInfoImpl.APIInfo().Tx.Tx, oldTopology.Name); err != nil {
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
		if err := removeNodes(oldTopology.APIInfoImpl.APIInfo().Tx.Tx, oldTopology.Name, &toRemove); err != nil {
			return nil, err, http.StatusInternalServerError
		}
	}
	if userErr, sysErr, errCode := topology.setTopologyDetails(); userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
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
	userErr, sysErr, statusCode := checkIfTopologyCanBeAlteredByCurrentUser(topology.APIInfo(), topology.Nodes)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
	}
	return nil, nil, 0
}

// OptionsDelete is a requirement of the OptionsDeleter interface.
func (topology *TOTopology) OptionsDelete() (error, error, int) {
	topologies, userErr, sysErr, errCode, _ := topology.Read(nil, false)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}
	if len(topologies) != 1 {
		return fmt.Errorf("cannot find exactly 1 topology with the query string provided"), nil, http.StatusBadRequest
	}
	topology.Topology = topologies[0].(tc.Topology)
	userErr, sysErr, statusCode := checkIfTopologyCanBeAlteredByCurrentUser(topology.APIInfo(), topology.Nodes)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
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

func selectQueryWithParentNames() string {
	query := `
SELECT t.name, tc.cachegroup,
	(SELECT COALESCE (ARRAY_AGG (tcpc.cachegroup)) AS parents
	FROM topology_cachegroup tc2
	INNER JOIN topology_cachegroup_parents tcp ON tc2.id = tcp.child
	INNER JOIN topology_cachegroup tcpc ON tcp.parent = tcpc.id
	WHERE tc2.topology = tc.topology
	AND tc2.cachegroup = tc.cachegroup
	)
FROM topology t
JOIN topology_cachegroup tc on t.name = tc.topology
`
	return query
}

func selectNonTopologyCacheGroupsQuery() string {
	query := `
SELECT 'non-topology cachegroups' AS name, c."name" AS cachegroup,
	(SELECT COALESCE (ARRAY_AGG (pc."name")) AS parents
	FROM cachegroup pc
	WHERE pc.id = c.parent_cachegroup_id
	OR pc.id = c.secondary_parent_cachegroup_id
	)
FROM cachegroup c
JOIN "type" t ON c."type" = t.id
WHERE (LOWER(t.name) LIKE :edge_type_prefix
OR LOWER(t.name) LIKE :mid_type_prefix)
AND (c.parent_cachegroup_id IS NOT NULL
OR c.secondary_parent_cachegroup_id IS NOT NULL)
`
	return query
}

func selectNonTopologyParentCacheGroupsQuery() string {
	query := `
SELECT 'non-topology cachegroups' AS name, pc2."name" AS cachegroup,
	(SELECT COALESCE (ARRAY_AGG (pc."name")) AS parents
	FROM cachegroup pc
	WHERE pc.id = pc2.parent_cachegroup_id
	OR pc.id = pc2.secondary_parent_cachegroup_id
	)
FROM cachegroup c
JOIN "type" t ON c."type" = t.id
JOIN cachegroup pc2 ON c.parent_cachegroup_id = pc2.id
	OR c.secondary_parent_cachegroup_id = pc2.id
WHERE (LOWER(t.name) LIKE :edge_type_prefix
OR LOWER(t.name) LIKE :mid_type_prefix)
AND (c.parent_cachegroup_id IS NOT NULL
OR c.secondary_parent_cachegroup_id IS NOT NULL)
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
name = $1,
description = $2
WHERE t.name = $3
RETURNING t.name, t.description, t.last_updated
`
	return query
}

func selectMaxLastUpdatedQuery(where string) string {
	return `SELECT max(ti) from (
		SELECT max(t.last_updated) as ti from topology t JOIN topology_cachegroup tc on t.name = tc.topology` + where +
		` UNION ALL
	select max(last_updated) as ti from last_deleted l where l.table_name='topology') as res`
}

func checkIfTopologyCanBeAlteredByCurrentUser(info *api.Info, nodes []tc.TopologyNode) (error, error, int) {
	cachegroups := getCachegroupNames(nodes)
	serverIDs, err := dbhelpers.GetServerIDsFromCachegroupNames(info.Tx.Tx, cachegroups)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	cdns, err := dbhelpers.GetCDNNamesFromServerIds(info.Tx.Tx, serverIDs)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	userErr, sysErr, statusCode := dbhelpers.CheckIfCurrentUserCanModifyCDNs(info.Tx.Tx, cdns, info.User.UserName)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, statusCode
	}
	return nil, nil, http.StatusOK
}
