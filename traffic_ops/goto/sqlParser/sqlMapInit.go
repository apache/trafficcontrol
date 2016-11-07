package sqlParser

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
	"fmt"
	"strings"
)

type ForeignKey struct {
	Table     string
	Column    string
	Alias     string
	ColValues map[string]interface{}
	//ColValues map[string]int
}

func MakeForeignKey(table string, column string, alias string) ForeignKey {
	colValues := GetForeignKeyValues(table, column)
	key := ForeignKey{table, column, alias, colValues}
	return key
}

func GetForeignKeyMap() map[string]ForeignKey {
	ForeignKeyMap := make(map[string]ForeignKey)

	//cachegroup.name
	key := MakeForeignKey("cachegroup", "name", "cachegroup_name")
	ForeignKeyMap["cachegroup"] = key

	//deliveryservice.xml_id
	key = MakeForeignKey("deliveryservice", "xml_id", "deliveryservice_name")
	ForeignKeyMap["deliveryservice"] = key
	ForeignKeyMap["job_deliveryservice"] = key

	//division.name
	key = MakeForeignKey("division", "name", "division_name")
	ForeignKeyMap["division"] = key

	//parameter.name
	key = MakeForeignKey("parameter", "name", "parameter_name")
	ForeignKeyMap["parameter"] = key

	//parent.cachegroup
	key = MakeForeignKey("cachegroup", "name", "parent_cachegroup")
	ForeignKeyMap["parent_cachegroup_id"] = key
	//phys_location.name!!
	key = MakeForeignKey("phys_location", "name", "phys_location_name")
	ForeignKeyMap["phys_location"] = key

	//profile.name
	key = MakeForeignKey("profile", "name", "profile_name")
	ForeignKeyMap["profile"] = key

	//regex.pattern
	key = MakeForeignKey("regex", "pattern", "regex_pattern")
	ForeignKeyMap["regex"] = key

	//region.name
	key = MakeForeignKey("region", "name", "region_name")
	ForeignKeyMap["region"] = key

	//status.name
	key = MakeForeignKey("status", "name", "status_name")
	ForeignKeyMap["status"] = key

	//server.host_name
	key = MakeForeignKey("server", "host_name", "server_name")
	ForeignKeyMap["serverid"] = key
	ForeignKeyMap["server"] = key

	//tm_user.username
	key = MakeForeignKey("tm_user", "username", "tm_user_username")
	ForeignKeyMap["tm_user"] = key
	ForeignKeyMap["tm_user_id"] = key
	ForeignKeyMap["job_user"] = key

	//type.name
	key = MakeForeignKey("type", "name", "type_name")
	ForeignKeyMap["type"] = key
	return ForeignKeyMap
}

//returns a map of each column name in table to its appropriate GoLang tpye (name string)
func GetColTypeMap() map[string]string {
	colMap := make(map[string]string, 0)

	cols, err := globalDB.Queryx("SELECT DISTINCT COLUMN_NAME, COLUMN_TYPE FROM information_schema.columns")
	check(err)

	for cols.Next() {
		var colName string
		var colType string

		err = cols.Scan(&colName, &colType)
		//split because SQL type returns are sometimes ex. int(11)
		colMap[colName] = strings.Split(colType, "(")[0]
	}

	return colMap
}

func GetTableMap(environment string) map[string][]string {
	var tableNames []string
	var tableMap = make(map[string][]string)

	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	fmt.Println("db", environment)
	rows, err := globalDB.Query("SELECT TABLE_NAME FROM information_schema.tables where (table_type='base table' or table_type='view') and table_schema='" + environment + "'")
	check(err)

	for rows.Next() {
		err := rows.Scan(tableInterface...)
		check(err)

		tableNames = append(tableNames, string(tableRawBytes))
	}

	for _, table := range tableNames {

		query := "SELECT column_name from information_schema.columns where table_name='" + table + "' and table_schema='" + environment + "' ORDER BY column_name asc"
		fmt.Println(query)
		rows, err = globalDB.Query(query)
		check(err)

		colMap := make([]string, 0)

		for rows.Next() {
			err = rows.Scan(tableInterface...)
			check(err)

			colMap = append(colMap, string(tableRawBytes))
		}

		tableMap[table] = colMap
	}
	return tableMap
}
