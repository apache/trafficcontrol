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
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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

func createDB() {
	dbExistsCmd := exec.Command("psql", "-h", HostIP, "-U", DBSuperUser, "-p", HostPort, "-tAc", "SELECT 1 FROM pg_database WHERE datname='"+DBName+"'")
	out, err := dbExistsCmd.Output()
	if err != nil {
		die("unable to check if DB already exists: " + err.Error())
	}
	if len(out) > 0 {
		fmt.Println("Database " + DBName + " already exists")
		return
	}
	createDBCmd := exec.Command("createdb", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-e", "--owner", DBUser, DBName)
	out, err = createDBCmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't create db " + DBName)
	}
}

func dropDB() {
	fmt.Println("Dropping database: " + DBName)
	cmd := exec.Command("dropdb", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-e", "--if-exists", DBName)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't drop db " + DBName)
	}
}

func createUser() {
	fmt.Println("Creating user: " + DBUser)
	userExistsCmd := exec.Command("psql", "-h", HostIP, "-U", DBSuperUser, "-p", HostPort, "-tAc", "SELECT 1 FROM pg_roles WHERE rolname='"+DBUser+"'")
	out, err := userExistsCmd.Output()
	if err != nil {
		die("unable to check if user already exists: " + err.Error())
	}
	if len(out) > 0 {
		fmt.Println("User " + DBUser + " already exists")
		return
	}
	createUserCmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-etAc", "CREATE USER "+DBUser+" WITH LOGIN ENCRYPTED PASSWORD '"+DBPassword+"'")
	out, err = createUserCmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't create user " + DBUser)
	}
}

func dropUser() {
	cmd := exec.Command("dropuser", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-i", "-e", DBUser)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't drop user " + DBUser)
	}
}

func showUsers() {
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-U", DBSuperUser, "-ec", `\du`)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't show users")
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
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-d", DBName, "-U", DBUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(seedsBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+DBPassword)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't patch database w/ required data")
	}
}

func loadSchema() {
	fmt.Println("Creating database tables.")
	schemaBytes, err := ioutil.ReadFile(DBSchemaPath)
	if err != nil {
		die("unable to read '" + DBSchemaPath + "': " + err.Error())
	}
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-d", DBName, "-U", DBUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(schemaBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+DBPassword)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't create database tables")
	}
}

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
	cmd := exec.Command("psql", "-h", HostIP, "-p", HostPort, "-d", DBName, "-U", DBUser, "-e", "-v", "ON_ERROR_STOP=1")
	cmd.Stdin = bytes.NewBuffer(patchesBytes)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+DBPassword)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out)
	if err != nil {
		die("Can't patch database w/ required data")
	}
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
