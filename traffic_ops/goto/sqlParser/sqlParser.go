
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlParser

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

var (
	globalDB      sqlx.DB
	colTypeMap    map[string]string
	foreignKeyMap map[string]ForeignKey
	tableMap      map[string][]string
	databaseName  string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*********************************************************************************
 * DB INITIALIZE: Connects given DB creds, creates ColMap FOR SESSION
 ********************************************************************************/
func InitializeDatabase(username string, password string, environment string) sqlx.DB {
	db, err := sqlx.Connect("mysql", username+":"+password+"@tcp(localhost:3306)/"+environment)
	check(err)

	globalDB = *db

	//set global colTypeMap
	databaseName = environment
	tableMap = GetTableMap(databaseName)
	colTypeMap = GetColTypeMap()
	foreignKeyMap = GetForeignKeyMap()
	return *db
}

/*********************************************************************************
 * HELPER FUNCTIONS
 ********************************************************************************/
//if is table, returns 1. else (for example, is view), returns 0.
func IsTable(serverTableName string) bool {
	//check if there is view. else, assume is table
	query := "select exists(select * from information_schema.tables where table_name='" + serverTableName + "' and table_name not in (select table_name from information_schema.views))"
	rows, err := globalDB.Query(query)
	check(err)

	//set up scan interface
	rawBytes := make([]byte, 1)
	scanInterface := make([]interface{}, 1)
	scanInterface[0] = &rawBytes

	//this should only return one row, but Scan panics if not called with Next
	for rows.Next() {
		err := rows.Scan(scanInterface...)
		check(err)
		//if exists as view, delete from view
		if string(rawBytes) == "1" {
			return true
		} else {
			return false
		}
	}

	return false
}

//returns array of table name strings from queried database
func GetTableNames() []string {
	var tableNames []string

	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("SELECT TABLE_NAME FROM information_schema.tables where (table_type='base table' or table_type='view') and table_schema='" + databaseName + "'")
	check(err)

	for rows.Next() {
		err := rows.Scan(tableInterface...)
		check(err)

		tableNames = append(tableNames, string(tableRawBytes))
	}

	return tableNames
}

//returns array of column names from table in database
func GetColumnNames(tableName string) []string {
	colNames := make([]string, 0)
	colNames = append(colNames, tableMap[tableName]...)

	return colNames
}

func GetForeignKeyColumns(tableName string) ([]string, map[string]map[string]interface{}) {
	tableCols := GetColumnNames(tableName)
	foreignKeyRows := make(map[string]map[string]interface{}, 0)
	for idx, col := range tableCols {
		if val, ok := foreignKeyMap[col]; ok {
			tableCols[idx] = val.Alias
			foreignKeyRows[col] = val.ColValues
		}
	}
	return tableCols, foreignKeyRows
}

func GetForeignKeyRows(tableName string) map[string]map[string]interface{} {
	tableCols := GetColumnNames(tableName)
	foreignKeyRows := make(map[string]map[string]interface{}, 0)

	for idx, col := range tableCols {
		if val, ok := foreignKeyMap[col]; ok {
			tableCols[idx] = val.Alias
			foreignKeyRow := val.ColValues
			foreignKeyRows[col] = foreignKeyRow
		}
	}
	return foreignKeyRows

}

/*********************************************************************************
 * DELETE FUNCTIONALITY
 ********************************************************************************/
func Delete(serverTableName string, parameters []string) (bool, error) {
	if !IsTable(serverTableName) {
		return DeleteFromView(serverTableName, parameters)
	} else {
		return false, DeleteFromTable(serverTableName, parameters)
	}
}

//deletes from a table
func DeleteFromTable(tableName string, parameters []string) error {
	return RunDeleteQuery(tableName, parameters)
}

//deletes from a view
func DeleteFromView(viewName string, parameters []string) (bool, error) {
	if len(parameters) == 0 {
		qStr := "drop view " + viewName
		_, err := globalDB.Query(qStr)
		return true, err
	} else {
		return false, RunDeleteQuery(viewName, parameters)
	}
}

//runs query of format "delete from tableName where parameterA=valueA and..."
func RunDeleteQuery(serverTableName string, parameters []string) error {
	//delete from tableName where x = a and y = b
	query := "delete from " + serverTableName

	if len(parameters) > 0 {
		query += " where "

		for _, v := range parameters {
			query += v + " and "
		}
		//removes last "and"
		query = query[:len(query)-4]
	}

	_, err := globalDB.Query(query)
	return err
}

/*********************************************************************************
 * GET FUNCTIONALITY
 ********************************************************************************/

func GetForeignKeyValues(tableName string, colName string) map[string]interface{} {
	//map into an array of type map[colName]value
	query := "select " + colName + ", id from " + tableName
	rows, err := globalDB.Queryx(query)
	check(err)
	//map into an array of type map[colName]value
	rowArray := make(map[string]interface{}, 0)

	for rows.Next() {
		cols, err := rows.SliceScan()
		check(err)

		if val, ok := cols[0].([]byte); ok {
			name := string(val)
			if val2, ok := cols[1].([]byte); ok {
				id := string(val2)
				rowArray[name] = id
			}
		}
	}

	return rowArray
}

func Get(tableName string, joinFKs bool, tableParameters []string) ([]map[string]interface{}, error) {
	regStr := ""
	whereStr := ""
	joinStr := ""
	onStr := ""
	cols := GetColumnNames(tableName)
	for _, col := range cols {
		if val, ok := foreignKeyMap[col]; ok && col != tableName && joinFKs {
			joinCols := GetColumnNames((val.Table))
			for _, joinCol := range joinCols {
				regStr += val.Table + "." + joinCol + " as " + val.Table + "_" + joinCol + ","
			}
			if col == "parent_cachegroup_id" {
				joinStr += "cachegroup2.name as parent_cachegroup,"
				onStr += " join cachegroup as cachegroup2 on cachegroup.parent_cachegroup_id = cachegroup2.id "
			} else {
				joinStr += val.Table + "." + val.Column + " as " + val.Alias + ","
				onStr += " join " + val.Table + " on " + tableName + "." + col + " = " + val.Table + ".id "
			}
		} else {
			if joinFKs {
				regStr += tableName + "." + col + " as " + tableName + "_" + col + ","
			} else {
				regStr += tableName + "." + col + ","
			}
		}
	}

	regStr = regStr[:len(regStr)-1]

	if joinStr != "" {
		joinStr = ", " + joinStr[:len(joinStr)-1]
	}

	sep := "where "
	for _, param := range tableParameters {
		if strings.ContainsAny(param, "=") { // > < and such?
			selectArr := strings.Split(param, "=")
			selectCol := selectArr[0]
			selectVal := selectArr[1]
			if strings.Contains(param, ";") { // prevent SQL injection. The rest will error in SQL
				err := errors.New("Invalid SQL detected:" + param)
				fmt.Println(err)
				return nil, err
			}

			whereStr += sep + selectCol + "=\"" + selectVal + "\""
			sep = " and "
		}
	}

	queryStr := "select " + regStr + joinStr + " from " + tableName + " "

	queryStr += onStr + " " + whereStr

	fmt.Println(queryStr)
	//do the query
	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		return nil, err
	}

	//map into an array of type map[colName]value
	rowArray := make([]map[string]interface{}, 0)

	for rows.Next() {
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			return nil, err
		}

		for k, v := range results {
			//converts the byte array to its correct type
			if b, ok := v.([]byte); ok {
				//if foreign key type conversion
				results[k] = string(b)
				//if val, ok := foreignKeyMap[k]; ok {
				//results[k], err = StringToType(b, colTypeMap[val.Column])
				//} else {
				//results[k], err = StringToType(b, colTypeMap[k])
				//}
				if err != nil {
					return nil, err
				}
			}
		}

		rowArray = append(rowArray, results)
	}
	return rowArray, nil
}

/*********************************************************************************
 * POST FUNCTIONALITY
 ********************************************************************************/
func Post(tableName string, jsonByte []byte) (string, error) {
	if IsTable(tableName) {
		err := PostRows(tableName, jsonByte)
		return tableName, err
	} else {
		return PostViews(jsonByte)
	}
}

//adds new row to table
func AddRow(newRow interface{}, tableName string) error {
	m := newRow.(map[string]interface{})
	//insert into table (colA, colB) values (valA, valB);
	query := "INSERT INTO " + tableName + " ("
	keyStr := ""
	valueStr := ""

	for k, v := range m {
		keyStr += k + ","
		typeStr, err := TypeToString(v)
		if err != nil {
			return err
		}

		valueStr += "'" + typeStr + "',"
	}

	keyStr = keyStr[:len(keyStr)-1]
	valueStr = valueStr[:len(valueStr)-1]

	query += keyStr + ") VALUES ( " + valueStr + " );"
	_, err := globalDB.Query(query)
	fmt.Println(query)
	return err
}

func AddRows(newRows []interface{}, tableName string) error {
	for _, row := range newRows {
		err := AddRow(row, tableName)
		if err != nil {
			return err
		}
	}
	foreignKeyMap = GetForeignKeyMap()
	return nil
}

/*********************************************************************************
 * POST FUNCTIONALITY
 ********************************************************************************/
//adds JSON from FILENAME to TABLE
//CURRENTLY ONLY ONE ROW
func PostRows(tableName string, jsonByte []byte) error {
	var f []interface{}

	err := json.Unmarshal(jsonByte, &f)
	if err != nil {
		return err
	}

	err2 := AddRows(f, tableName)
	if err2 != nil {
		return err2
	}

	return nil
}

//view details are marshalled into this View struct
type View struct {
	Name  string
	Query string
}

//adds JSON from FILENAME to TABLE
func PostViews(jsonByte []byte) (string, error) {
	var views []View

	err := json.Unmarshal(jsonByte, &views)
	if err != nil {
		return "", err
	}

	var viewName string
	for _, view := range views {
		viewName = view.Name
		err = MakeView(view.Name, view.Query)
		if err != nil {
			return viewName, err
		}
	}

	return viewName, nil
}

func MakeView(viewName string, view string) error {
	qStr := "create view " + viewName + " as " + view
	_, err := globalDB.Query(qStr)
	tableMap = GetTableMap(databaseName)
	return err
}

/*********************************************************************************
 * PUT FUNCTIONALITY
 ********************************************************************************/
func Put(tableName string, jsonByte []byte) error {
	//unmarshals the json into an interface
	var f []interface{}
	err := json.Unmarshal(jsonByte, &f)
	if err != nil {
		return err
	}
	//adds the interface row to table in database
	return UpdateRows(f, tableName)
}

func UpdateRows(newRows []interface{}, tableName string) error {
	for _, row := range newRows {
		err := UpdateRow(row, tableName)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateRow(newRow interface{}, tableName string) error {
	query := "update " + tableName

	updateParameters := newRow.(map[string]interface{})

	var idString string
	//new changes
	if len(updateParameters) > 0 {
		query += " set "

		for k, v := range updateParameters {
			typeStr, err := TypeToString(v)
			if err != nil {
				return err
			}
			if k == "last_updated" {
				idString += k + "='" + typeStr + "' "
			} else {
				query += k + "='" + typeStr + "', "
			}
		}

		query = query[:len(query)-2] + " where " + idString
	}

	fmt.Println(query)
	_, err := globalDB.Query(query)
	return err
}
