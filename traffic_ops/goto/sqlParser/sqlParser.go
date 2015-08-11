// Copyright 2015 Comcast Cable Communications Management, LLC

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

/**********************************************************************************
 * DIRECTORY:
 * 1. DB INITIALIZE
 	a. InitializeDatabase(username string, password string, environment string) sqlx.DB{}
	   Initializes the database and builds a column type map for reference in get query type conversions.
	b. GetColMap() map[string]string{}
	   Map of each column name in table to its appropriate GoLang type.
 * 2. HELPER FUNCTIONS
 	a. IsTable(serverTableName string) int{}
	   If is a table in the DB, returns 1. Else, returns 0. (Useful for differentiating between
	   views and tables.)
 * 3. DELETE
 	a. Delete(serverTableName string, parameters []string){}
	   Handles delete request by differentiating between view/table.
	b. DeleteFromTable(tableName string, parameters []string)
	   Deletes rows from tableName given parameters.
	   Note that this deletes rows, does not drop tables.
	c. DeleteFromView(viewName string, parameters []string){}
	   Deletes rows from viewName given parameters.
	   Note that this deletes rows but ALSO COULD DROP TABLES.
	d. RunDeleteQuery(serverTableName string, parameters []string){}
	   Constructs and runs general delete query.
 * 4. GET
 	a. Get(tableName string, tableParams []string) []map[string]interface(){}
	   Constructs and runs general query, returning results in an array of maps
	   (each map represents a row, with column name key and actual value in value.
 * 5. POST
    a. Post(tableName string, fileName string){}
	   Handles post request by differentiating between view/table.
	b. AddRow(newRow itnerface{}, tableName string){}
	   Adds a new row to the table by constructing and querying an insert statement.
	c. AddRows(newRows []interface{}, tableName string){}
	   Adds multiple rows to table.
	d. PostRows(tableName string, fileName string){}
	   Parses the JSON-represented rows from given file and adds them to table.
	e. type View struct{}
	   Views are added via POSTs of the view name and the query.
	f. func PostViews(fileName string){}
	   Parses the JSON-represented view and adds the new view to the table.
	g. func MakeView(viewName string, view string){}
	   Constructs and queries the statement needed to add a new view to the database.
 * 6. PUT
 	a. Put(tableName string, parameters []string, fileName string){}
	   Parses the JSON-represented rows in the given file and UPDATES the rows specified.
	b. UpdateRow(newRow itnerface{}, tableName string, parameters []string){}
	   Constructs and queries the statement needed to update rows.
 *********************************************************************************/

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

var (
	globalDB sqlx.DB
	colMap   map[string]string
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

	//set global colMap
	colMap = GetColMap()
	return *db
}

//returns a map of each column name in table to its appropriate GoLang tpye (name string)
func GetColMap() map[string]string {
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

	rows, err := globalDB.Query("SELECT TABLE_NAME FROM information_schema.tables where table_type='base table' or table_type='view'")
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
	var colNames []string

	colRawBytes := make([]byte, 1)
	colInterface := make([]interface{}, 1)

	colInterface[0] = &colRawBytes

	rows, err := globalDB.Query("SELECT COLUMN_NAME FROM information_schema.columns WHERE TABLE_NAME='" + tableName + "' ORDER BY column_name asc")
	check(err)

	for rows.Next() {
		err := rows.Scan(colInterface...)
		check(err)

		colNames = append(colNames, string(colRawBytes))
	}
	return colNames
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
//returns interface from given table OR view from queried database
func Get(tableName string, tableParams []string) ([]map[string]interface{}, error) {
	//if where exists, append
	whereStmt := ""
	if len(tableParams) > 0 {
		whereStmt += " where "

		for _, v := range tableParams {
			whereStmt += v + " and "
		}

		whereStmt = whereStmt[:len(whereStmt)-4]
	}

	//do the query
	rows, err := globalDB.Queryx("SELECT * from " + tableName + whereStmt)
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
				results[k], err = StringToType(b, colMap[k])
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

	return err
}

func AddRows(newRows []interface{}, tableName string) error {
	for _, row := range newRows {
		err := AddRow(row, tableName)
		if err != nil {
			return err
		}
	}

	return nil
}

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
		return err
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
	return err
}

/*********************************************************************************
 * PUT FUNCTIONALITY
 ********************************************************************************/
func Put(tableName string, parameters []string, jsonByte []byte) error {
	//unmarshals the json into an interface
	var f []interface{}
	err := json.Unmarshal(jsonByte, &f)
	if err != nil {
		return err
	}
	//adds the interface row to table in database
	return UpdateRows(f, tableName, parameters)
}

func UpdateRows(newRows []interface{}, tableName string, parameters []string) error {
	for _, row := range newRows {
		err := UpdateRow(row, tableName, parameters)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateRow(newRow interface{}, tableName string, parameters []string) error {
	query := "update " + tableName

	updateParameters := newRow.(map[string]interface{})
	//new changes
	if len(updateParameters) > 0 {
		query += " set "

		for k, v := range updateParameters {
			typeStr, err := TypeToString(v)
			if err != nil {
				return err
			}
			query += k + "='" + typeStr + "', "
		}

		query = query[:len(query)-2]
	}

	//where
	if len(parameters) > 0 {
		query += " where "

		for _, v := range parameters {
			query += v + " and "
		}

		query = query[:len(query)-4]
	}

	_, err := globalDB.Query(query)
	return err
}
