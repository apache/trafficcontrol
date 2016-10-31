
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// started from https://github.com/asdf072/struct-create

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
)

var defaults = Configuration{
	DbUser:     os.Args[1],
	DbPassword: os.Args[2],
	DbName:     os.Args[3],
	PkgName:    "todb",
	TagLabel:   "db",
}

var config Configuration

type Configuration struct {
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
	// PkgName gives name of the package using the stucts
	PkgName string `json:"pkg_name"`
	// TagLabel produces tags commonly used to match database field names with Go struct members
	TagLabel string `json:"tag_label"`
}

type ColumnSchema struct {
	TableName              string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnKey              string
}

func writeGetters(schemas []ColumnSchema) (int, error) {
	file, err := os.Create("dbgetters.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	header := "package todb\n\n"
	header += "import (\n"
	// header += "\"gopkg.in/guregu/null.v3\"\n"
	header += "\"fmt\"\n"
	header += ")\n"

	out := ""
	currentTable := ""
	for _, cs := range schemas {
		if cs.TableName != currentTable {
			out += "func get" + formatName(cs.TableName) + "()([]" + formatName(cs.TableName) + ", error) {\n"
			out += "	ret := []" + formatName(cs.TableName) + "{}\n"
			out += "	queryStr := \"select * from " + cs.TableName + "\"\n"
			out += "	err := globalDB.Select(&ret, queryStr)\n"
			out += "	if err != nil {\n"
			out += "		fmt.Println(err)\n"
			out += "		return nil, err\n"
			out += "	}\n"
			out += "	return ret, nil\n"
			out += "}\n\n"

		}
		currentTable = cs.TableName
	}

	currentTable = ""
	out += "func GetTable(tableName string) (interface{}, error) {\n"
	for _, cs := range schemas {
		if cs.TableName != currentTable {
			out += "	if tableName == \"" + cs.TableName + "\" {\n"
			out += "		return get" + formatName(cs.TableName) + "()\n"
			out += "	}\n"
		}
		currentTable = cs.TableName
	}
	out += "	return nil, nil\n"
	out += "}\n\n"

	totalBytes, err := fmt.Fprint(file, header+out)
	if err != nil {
		log.Fatal(err)
	}
	return totalBytes, nil
}

func writeStructs(schemas []ColumnSchema) (int, error) {
	file, err := os.Create("dbstructs.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	currentTable := ""

	out := ""
	for _, cs := range schemas {

		if cs.TableName != currentTable {
			if currentTable != "" {
				out = out + "}\n\n"
			}
			out = out + "type " + formatName(cs.TableName) + " struct{\n"
		}

		goType, _, err := goType(&cs)

		if err != nil {
			log.Fatal(err)
		}
		out = out + "\t" + formatName(cs.ColumnName) + " " + goType
		if len(config.TagLabel) > 0 {
			out = out + "\t`" + config.TagLabel + ":\"" + cs.ColumnName + "\" json:\"" + formatNameLower(cs.ColumnName) + "\"`"
		}
		out = out + "\n"
		currentTable = cs.TableName

	}
	out = out + "}"

	header := "package " + config.PkgName + "\n\n"
	header += "import (\n"
	header += "\"gopkg.in/guregu/null.v3\"\n"
	header += "\"time\"\n"
	header += ")\n\n"

	totalBytes, err := fmt.Fprint(file, header+out)
	if err != nil {
		log.Fatal(err)
	}
	return totalBytes, nil
}

func getSchema() []ColumnSchema {
	conn, err := sql.Open("mysql", config.DbUser+":"+config.DbPassword+"@/information_schema")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	q := "SELECT TABLE_NAME, COLUMN_NAME, IS_NULLABLE, DATA_TYPE, " +
		"CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, COLUMN_TYPE, " +
		"COLUMN_KEY FROM COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION"
	rows, err := conn.Query(q, config.DbName)
	if err != nil {
		log.Fatal(err)
	}
	columns := []ColumnSchema{}
	for rows.Next() {
		cs := ColumnSchema{}
		err := rows.Scan(&cs.TableName, &cs.ColumnName, &cs.IsNullable, &cs.DataType,
			&cs.CharacterMaximumLength, &cs.NumericPrecision, &cs.NumericScale,
			&cs.ColumnType, &cs.ColumnKey)
		if err != nil {
			log.Fatal(err)
		}
		columns = append(columns, cs)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return columns
}

func formatName(name string) string {
	parts := strings.Split(name, "_")
	newName := ""
	for _, p := range parts {
		if len(p) < 1 {
			continue
		}
		newName = newName + strings.Replace(p, string(p[0]), strings.ToUpper(string(p[0])), 1)
	}
	return newName
}

func formatNameLower(name string) string {
	newName := formatName(name)
	newName = strings.Replace(newName, string(newName[0]), strings.ToLower(string(newName[0])), 1)
	return newName
}

func goType(col *ColumnSchema) (string, string, error) {
	requiredImport := ""
	if col.IsNullable == "YES" {
		requiredImport = "database/sql"
	}
	var gt string = ""
	switch col.DataType {
	case "char", "varchar", "enum", "text", "longtext", "mediumtext", "tinytext":
		if col.IsNullable == "YES" {
			gt = "null.String"
		} else {
			gt = "string"
		}
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		gt = "[]byte"
	case "date", "time", "datetime", "timestamp":
		gt, requiredImport = "time.Time", "time"
	case "tinyint", "smallint", "int", "mediumint", "bigint":
		if col.IsNullable == "YES" {
			gt = "null.Int"
		} else {
			gt = "int64"
		}
	case "float", "decimal", "double":
		if col.IsNullable == "YES" {
			gt = "null.Float"
		} else {
			gt = "float64"
		}
	}
	if gt == "" {
		n := col.TableName + "." + col.ColumnName
		return "", "", errors.New("No compatible datatype (" + col.DataType + ") for " + n + " found")
	}
	return gt, requiredImport, nil
}

var configFile = flag.String("json", "", "Config file")

func main() {
	flag.Parse()

	if len(*configFile) > 0 {
		f, err := os.Open(*configFile)
		if err != nil {
			log.Fatal(err)
		}
		err = json.NewDecoder(f).Decode(&config)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		config = defaults
	}

	columns := getSchema()
	bytes, err := writeStructs(columns)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ok %d\n", bytes)
	bytes, err = writeGetters(columns)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ok %d\n", bytes)
}
