package topology

import (
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

type TOTopology struct {
	api.APIInfoImpl `json:"-"`
	tc.Topology
}

func (topology *TOTopology) SetLastUpdated(time tc.TimeNoMod) { topology.LastUpdated = &time }

func (topology TOTopology) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"name", api.GetStringKey}}
}

func (topology *TOTopology) GetType() string {
	return "topology"
}

func (topology *TOTopology) Validate() error {
	return nil
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
	nodeCount := len(topology.Nodes)
	for index := 0; index < nodeCount; index++ {
		node := &topology.Nodes[index]
		err := tx.QueryRow(nodeInsertQuery(), topology.Name, node.Cachegroup).Scan(&node.Id, &topology.Name, &node.Cachegroup, &node.LastUpdated)
		if err != nil {
			userErr, sysErr, errCode := api.ParseDBError(err)
			return userErr, sysErr, errCode
		}
	}

	for index := 0; index < nodeCount; index++ {
		node := &topology.Nodes[index]
		for rank := 1; rank <= len(node.Parents); rank++ {
			parent := topology.Nodes[node.Parents[rank-1]]
			err := tx.QueryRow(nodeParentInsertQuery(), node.Id, parent.Id, &rank).Scan(&node.Id, &parent.Id, &rank)
			if err != nil {
				userErr, sysErr, errCode := api.ParseDBError(err)
				return userErr, sysErr, errCode
			}
		}
	}

	return nil, nil, 0
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
VALUES ($1, $2)
RETURNING id, topology, cachegroup, last_updated
`
	return query
}

func nodeParentInsertQuery() string {
	query := `
INSERT INTO topology_cachegroup_parents (child, parent, rank)
VALUES ($1, $2, $3)
RETURNING child, parent, rank
`
	return query
}
