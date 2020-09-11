package main

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
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Development GooseConfig `yaml:"development"`
	Test        GooseConfig `yaml:"test"`
	Integration GooseConfig `yaml:"integration"`
	Production  GooseConfig `yaml:"production"`
}

type GooseConfig struct {
	Driver string `yaml:"driver"`
	Open   string `yaml:"open"`
}

func (conf DBConfig) getGooseConfig(env string) (GooseConfig, error) {
	switch env {
	case EnvDevelopment:
		return conf.Development, nil
	case EnvTest:
		return conf.Test, nil
	case EnvIntegration:
		return conf.Integration, nil
	case EnvProduction:
		return conf.Production, nil
	default:
		return GooseConfig{}, errors.New("invalid environment: " + env)
	}
}

const (
	// the possible environments to use
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvIntegration = "integration"
	EnvProduction  = "production"

	// keys in the goose config's "open" string value
	HostKey     = "host"
	PortKey     = "port"
	UserKey     = "user"
	PasswordKey = "password"
	DBNameKey   = "dbname"

	// available commands
	CmdCreateDB      = "createdb"
	CmdDropDB        = "dropdb"
	CmdCreateUser    = "create_user"
	CmdDropUser      = "drop_user"
	CmdShowUsers     = "show_users"
	CmdReset         = "reset"
	CmdUpgrade       = "upgrade"
	CmdMigrate       = "migrate"
	CmdDown          = "down"
	CmdRedo          = "redo"
	CmdStatus        = "status"
	CmdDBVersion     = "dbversion"
	CmdSeed          = "seed"
	CmdLoadSchema    = "load_schema"
	CmdReverseSchema = "reverse_schema"
	CmdPatch         = "patch"

	// goose commands that don't match the commands for this tool
	GooseUp = "up"

	DBConfigPath       = "db/dbconf.yml"
	DBSeedsPath        = "db/seeds.sql"
	DBSchemaPath       = "db/create_tables.sql"
	DBPatchesPath      = "db/patches.sql"
	DefaultEnvironment = EnvDevelopment
	DefaultDBSuperUser = "postgres"
)

var (
	// globals that are passed in via CLI flags and used in commands
	Environment string

	// globals that are parsed out of DBConfigFile and used in commands
	DBName      string
	DBSuperUser = DefaultDBSuperUser
	DBUser      string
	DBPassword  string
	HostIP      string
	HostPort    string
)

func parseDBConfig() error {
	confBytes, err := ioutil.ReadFile(DBConfigPath)
	if err != nil {
		return errors.New("reading DB conf '" + DBConfigPath + "': " + err.Error())
	}

	dbConfig := DBConfig{}
	err = yaml.Unmarshal(confBytes, &dbConfig)
	if err != nil {
		return errors.New("unmarshalling DB conf yaml: " + err.Error())
	}

	gooseCfg, err := dbConfig.getGooseConfig(Environment)
	if err != nil {
		return errors.New("getting goose config: " + err.Error())
	}

	// parse the 'open' string into a map
	open := make(map[string]string)
	pairs := strings.Split(gooseCfg.Open, " ")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.Split(pair, "=")
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			continue
		}
		open[kv[0]] = kv[1]
	}

	ok := false
	HostIP, ok = open[HostKey]
	if !ok {
		return errors.New("unable to get '" + HostKey + "' for environment '" + Environment + "'")
	}
	HostPort, ok = open[PortKey]
	if !ok {
		return errors.New("unable to get '" + PortKey + "' for environment '" + Environment + "'")
	}
	DBUser, ok = open[UserKey]
	if !ok {
		return errors.New("unable to get '" + UserKey + "' for environment '" + Environment + "'")
	}
	DBPassword, ok = open[PasswordKey]
	if !ok {
		return errors.New("unable to get '" + PasswordKey + "' for environment '" + Environment + "'")
	}
	DBName, ok = open[DBNameKey]
	if !ok {
		return errors.New("unable to get '" + DBNameKey + "' for environment '" + Environment + "'")
	}

	return nil
}

func connectAndExecute(query string, queryErrPrefix string, returnRows bool) *sql.Rows {
	//todo prepared statements
	//todo inspect results of executing the above query for all possible cases (db exists, db dne, db inaccessible)
	//todo determine test coverage of this portion
	if db, err := dburl.Open(fmt.Sprintf("pg://%s:%s@%s:%s/", DBSuperUser, DBPassword, HostIP, HostPort)); err == nil {
		if !returnRows {
			if _, err := db.Exec(query); err != nil {
				die(queryErrPrefix + ":" + err.Error())
			}
		} else {
			if rows, err := db.Query(query); err != nil {
				die(queryErrPrefix + ":" + err.Error())
			} else {
				return rows
			}
		}
	} else {
		die("couldn't connect to pg at " + HostIP + ": " + err.Error())
	}
	return nil
}

func createDB() {
	fmt.Println("Creating database: " + DBName)
	connectAndExecute(`CREATE DATABASE `+DBName+` OWNER `+DBUser, "Couldn't create database "+DBName, false)
}

func dropDB() {
	fmt.Println("Dropping database: " + DBName)
	connectAndExecute(`DROP DATABASE `+DBName, "Couldn't drop database:"+DBName, false)
}

func createUser() {
	fmt.Println("Creating user: " + DBUser)

	userExistsCmd := connectAndExecute(
		"SELECT 1 FROM pg_roles WHERE rolname='"+DBUser+"'",
		"Unable to query for user existence "+DBUser, true)
	var exists int
	if userExistsCmd.Scan(&exists); exists != 1 {
		fmt.Println("User does not exist, creating: " + DBUser)
		connectAndExecute(
			"CREATE USER "+DBUser+" WITH LOGIN ENCRYPTED PASSWORD '"+DBPassword+"'",
			"Unable to create user",
			false)
	}
}

func dropUser() {
	// TODO do i need to handle owned schema before dropping?
	connectAndExecute(`DROP ROLE `+DBUser, `Couldn't drop user `+DBUser, false)
}

func showUsers() {
	rows := connectAndExecute(`
    SELECT usename AS role_name
	FROM pg_catalog.pg_user
	ORDER BY role_name desc;`, `Couldn't show users`, true)
	// todo cleaner output
	fmt.Println("role_name    |     role_attributes")
	var roleName, roleAttributes string
	for rows.Next() {
		if err := rows.Scan(&roleName, &roleAttributes); err == nil {
			fmt.Printf("%s            | %s", roleName, roleAttributes)
		} else {
			fmt.Printf("Couldn't read row: %s", err.Error())
		}
	}
}

func reset() {
	createUser()
	dropDB()
	createDB()
	loadSchema()
	migrate()
}

func upgrade() {
	goose(GooseUp)
	seed()
	patch()
}

func migrate() {
	goose(GooseUp)
}

func down() {
	goose(CmdDown)
}

func redo() {
	goose(CmdRedo)
}

func status() {
	goose(CmdStatus)
}

func dbVersion() {
	goose(CmdDBVersion)
}

func seed() {
	fmt.Println("Seeding database w/ required data.")
	seedsBytes, err := ioutil.ReadFile(DBSeedsPath)
	if err != nil {
		die("unable to read '" + DBSeedsPath + "': " + err.Error())
	}
	connectAndExecute(string(seedsBytes), "Couldn't seed database", false)
}

func loadSchema() {
	fmt.Println("Creating database tables.")
	schemaBytes, err := ioutil.ReadFile(DBSchemaPath)
	if err != nil {
		die("unable to read '" + DBSchemaPath + "': " + err.Error())
	}
	connectAndExecute(string(schemaBytes), "Couldn't create tables", false)
}

//TODO remove this after perlectomy
func reverseSchema() {
	fmt.Fprintf(os.Stderr, "WARNING: the '%s' command will be removed with Traffic Ops Perl because it will no longer be necessary\n", CmdReverseSchema)
	cmd := exec.Command("db/reverse_schema.pl")
	cmd.Env = append(os.Environ(), "MOJO_MODE="+Environment)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't run `db/reverse_schema.pl`: " + err.Error())
	}
}

func patch() {
	fmt.Println("Patching database with required data fixes.")
	patchesBytes, err := ioutil.ReadFile(DBPatchesPath)
	if err != nil {
		die("unable to read '" + DBPatchesPath + "': " + err.Error())
	}
	connectAndExecute(string(patchesBytes), "Couldn't patch database", false)
}

func goose(arg string) {
	fmt.Println("Running goose " + arg + "...")
	cmd := exec.Command("goose", "--env="+Environment, arg)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't run goose: " + err.Error())
	}
}

func die(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func usage() string {
	programName := os.Args[0]
	home := "$HOME"
	home = os.Getenv("HOME")
	return `
Usage:  ` + programName + ` [--env (development|test|production|integration)] [arguments]

Example:  ` + programName + ` --env=test reset

Purpose:  This script is used to manage database. The environments are
          defined in the dbconf.yml, as well as the database names.

NOTE:
Postgres Superuser: The 'postgres' superuser needs to be created to run ` + programName + ` and setup databases.
If the 'postgres' superuser has not been created or password has not been set then run the following commands accordingly.

Create the 'postgres' user as a super user (if not created):

     $ createuser postgres --superuser --createrole --createdb --login --pwprompt

Modify your ` + home + `/.pgpass file which allows for easy command line access by defaulting the user and password for the database
without prompts.

 Postgres .pgpass file format:
 hostname:port:database:username:password

 ----------------------
 Example Contents
 ----------------------
 *:*:*:postgres:your-postgres-password
 *:*:*:traffic_ops:the-password-in-dbconf.yml
 ----------------------

 Save the following example into this file ` + home + `/.pgpass with the permissions of this file
 so only your user can read and write.

     $ chmod 0600 ` + home + `/.pgpass

===================================================================================================================
` + programName + ` arguments:

createdb  - Execute db 'createdb' the database for the current environment.
create_user  - Execute 'create_user' the user for the current environment (traffic_ops).
dropdb  - Execute db 'dropdb' on the database for the current environment.
down  - Roll back a single migration from the current version.
drop_user  - Execute 'drop_user' the user for the current environment (traffic_ops).
patch  - Execute sql from db/patches.sql for loading post-migration data patches.
redo  - Roll back the most recently applied migration, then run it again.
reset  - Execute db 'dropdb', 'createdb', load_schema, migrate on the database for the current environment.
reverse_schema  - Reverse engineer the lib/Schema/Result files from the environment database.
seed  - Execute sql from db/seeds.sql for loading static data.
show_users  - Execute sql to show all of the user for the current environment.
status  - Print the status of all migrations.
upgrade  - Execute migrate, seed, and patches on the database for the current environment.
migrate  - Execute migrate (without seeds or patches) on the database for the current environment.
`
}

func main() {
	flag.StringVar(&Environment, "env", DefaultEnvironment, "The environment to use (defined in "+DBConfigPath+").")
	flag.Parse()
	if len(flag.Args()) != 1 || flag.Arg(0) == "" {
		die(usage())
	}
	if Environment == "" {
		die(usage())
	}
	if err := parseDBConfig(); err != nil {
		die(err.Error())
	}
	commands := make(map[string]func())

	commands[CmdCreateDB] = createDB
	commands[CmdDropDB] = dropDB
	commands[CmdCreateUser] = createUser
	commands[CmdDropUser] = dropUser
	commands[CmdShowUsers] = showUsers
	commands[CmdReset] = reset
	commands[CmdUpgrade] = upgrade
	commands[CmdMigrate] = migrate
	commands[CmdDown] = down
	commands[CmdRedo] = redo
	commands[CmdStatus] = status
	commands[CmdDBVersion] = dbVersion
	commands[CmdSeed] = seed
	commands[CmdLoadSchema] = loadSchema
	commands[CmdReverseSchema] = reverseSchema
	commands[CmdPatch] = patch

	userCmd := flag.Arg(0)
	if cmd, ok := commands[userCmd]; ok {
		cmd()
	} else {
		fmt.Println(usage())
		die("invalid command: " + userCmd)
	}
}
