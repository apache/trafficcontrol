#!/usr/bin/env perl
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
use strict;
use warnings;
use English;
use Getopt::Long;
use FileHandle;
use DBI;
use Cwd;
use Data::Dumper;
use Schema;
use CPAN::Meta;
use File::Find::Rule;

use YAML;
use YAML qw(LoadFile);
use DBIx::Class::Schema::Loader qw/make_schema_at/;

my $usage = "\n"
	. "Usage:  $PROGRAM_NAME [--env (development|test|production|integration)] [--force] --admuser "
    . "(Postgres Admin User) --admpwd (Postgres Admin Password) [arguments]\t\n\n"
	. "Example:  $PROGRAM_NAME --env=test --admuser=postgres --admpwd=postgres123 reset\n\n"
	. "Purpose:  This script is used to manage database. The environments are\n"
	. "          defined in the dbconf.yml, as well as the database names.\n\n"
	. "    --admuser: The Postgres DB admin user that has postgres superuser privledges.\n"
	. "    --admpwd: The password for the Postgres DB admin user.\n"
	. "    --force: Using this flag will proceed to perform tasks without asking.\n"
    . "             The default is to ask to proceed.\n\n"
	. "arguments:   \n\n"
	. "createdb  - Execute db 'createdb' the database for the current environment.\n"
	. "dropdb  - Execute db 'dropdb' on the database for the current environment.\n"
	. "down  - Roll back a single migration from the current version.\n"
	. "createuser  - Execute 'createuser' the user for the current environment.\n"
	. "dropuser  - Execute 'dropuser' the user for the current environment.\n"
	. "showusers  - Execute sql to show all of the user for the current environment.\n"
	. "redo  - Roll back the most recently applied migration, then run it again.\n"
	. "reset  - Execute db 'dropdb', 'createdb', load_schema, migrate on the database for the current environment.\n"
	. "reverse_schema  - Reverse engineer the lib/Schema/Result files from the environment database.\n"
	. "seed  - Execute sql from db/seeds.sql for loading static data.\n"
	. "setup  - Execute db dropdb, createdb, load_schema, migrate, seed on the database for the current environment.\n"
	. "status  - Print the status of all migrations.\n"
	. "upgrade  - Execute migrate then seed on the database for the current environment.\n";

my $environment = 'development';
my $db_protocol;
my $db_admin_username;  # For the Postgres DB Admin User
my $db_admin_password;  # For the Postgres DB Admin User Password.

# supplied command-line "argument", e.g. setup, etc., that requires the admin credentials.
my @requires_admin = qw(createdb dropdb createuser dropuser reset setup);

# supplied command-line "argument", e.g. setup, etc., that we need to ask to proceed.
my @ask_user = qw(createdb dropdb createuser dropuser reset setup upgrade migrate down
	              redo seed load_schema reverse_schema);

# If force is 0; the script will ask to proceed. 1 will not ask and proceed.
# This is to add protection since this script is dangerous and we should
# protect an existing database from a potentially a catastrophic change.
my $force = 0;

# This is defaulted to 'to_development' so
# you don't have to specify --env=development for dev workstations
my $db_conf_path      = 'db/dbconf.yml';
my $db_name           = 'to_development';
my $db_username       = 'to_development';
my $db_password       = '';
my $host_ip           = '';
my $host_port         = '';
GetOptions("env=s" => \$environment,
	       "force" => \$force,
	       "admuser=s" => \$db_admin_username,
	       "admpwd=s" => \$db_admin_password);
$ENV{'MOJO_MODE'} = $environment;

parse_dbconf_yml_pg_driver();

STDERR->autoflush(1);
my $argument = shift(@ARGV);
if ( defined($argument) ) {
	if ($argument ~~ @requires_admin) {
		if (!defined $db_admin_username || !defined $db_admin_password) {
			print "FATAL: The database admin credentials needs to be supplied on the command-line.\n" . $usage;
			exit 1;
		}
	}

	if ($argument ~~ @ask_user) {
		if (!$force) {
			ask_user_to_proceed();
		}
	}

	if ( $argument eq 'createdb' ) {
		createdb();
	}
	elsif ( $argument eq 'dropdb' ) {
		dropdb();
	}
	elsif ( $argument eq 'createuser' ) {
		createuser();
	}
	elsif ( $argument eq 'dropuser' ) {
		dropuser();
	}
	elsif ( $argument eq 'showusers' ) {
		showusers();
	}
	elsif ( $argument eq 'reset' ) {
		dropdb();
		createdb();
		load_schema();
		migrate('up');
	}
	elsif ( $argument eq 'upgrade' ) {
		migrate('up');
		seed();
	}
	elsif ( $argument eq 'setup' ) {
		createuser();
		dropdb();
		createdb();
		load_schema();
		migrate('up');
		seed();
	}
	elsif ( $argument eq 'migrate' ) {
		migrate('up');
	}
	elsif ( $argument eq 'down' ) {
		migrate('down');
	}
	elsif ( $argument eq 'redo' ) {
		migrate('redo');
	}
	elsif ( $argument eq 'status' ) {
		migrate('status');
	}
	elsif ( $argument eq 'dbversion' ) {
		migrate('dbversion');
	}
	elsif ( $argument eq 'seed' ) {
		seed();
	}
	elsif ( $argument eq 'load_schema' ) {
		load_schema();
	}
	elsif ( $argument eq 'reverse_schema' ) {
		reverse_schema();
	}
	else {
		print $usage;
	}
}
else {
	print $usage;
}

exit(0);

sub ask_user_to_proceed {
	print "\nWARNING: This action is making changes to your database.\n\nAre you sure you want to "
	      . "proceed? ('Yes' to proceed): ";
	my $proceed = <STDIN>;
	chomp($proceed);
	if ($proceed ne "Yes") {
		# If not "Yes", exit.
		print "Exiting.\n";
		exit 1;
	}
}

sub parse_dbconf_yml_pg_driver {
	my $db_conf       = LoadFile($db_conf_path);
	my $db_connection = $db_conf->{$environment};

	# let's make sure the environment specified or set is actually configured in the confiration file.
	if (!defined $db_connection) {
		die "An invalid environment [$environment] has been specified.  This enviornment is not configured " .
				"in the configuration file [$db_conf_path].";
	}

	$db_protocol = $db_connection->{driver};
	my $open = $db_connection->{open};

	# Goose requires the 'open' line in the dbconf file to be a scalar.
	# example:
	#		open: host=127.0.0.1 port=5432 user=to_user password=twelve dbname=to_development sslmode=disable
	# We need access to these values for db connections so I am manipulating the 'open'
	# line so that it can be loaded into a hash.
	$open = join "\n", map { s/=/ : /; $_ } split " ", $open;
	my $hash = Load $open;

	$host_ip     = $hash->{host};
	$host_port   = $hash->{port};
	$db_name     = $hash->{dbname};
	$db_username = $hash->{user};
	$db_password = $hash->{password};
}

sub set_default {
	my $variable = shift;
	my $default = shift;

	my $retval = $default;
	if (defined $variable) {
		$retval = $variable;
	}

	return $retval;
}

sub get_psql_uri {
	my $dbuser = set_default(shift, $db_username);
	my $dbpasswd = set_default(shift, $db_password);
	my $dbname = set_default(shift, $db_name);
    my $dbhost = set_default(shift, $host_ip);
	my $dbport = set_default(shift, $host_port);

	my $uri = sprintf 'postgresql://%s:%s@%s:%s', $dbuser, $dbpasswd, $dbhost, $dbport;

	if ( defined $dbname ) {
		$uri .= "/$dbname";
	}

	return $uri;
}

sub migrate {
	my ($command) = @_;

	print "Migrating database...\n";
	if ( system("goose --env=$environment $command") != 0 ) {
		die "Can't run goose\n";
	}
}

sub seed {
	print "Seeding database.\n";
	my $uri = get_psql_uri();
	if ( system("psql $uri -e < db/seeds.sql") != 0 ) {
		die "Can't seed database\n";
	}
}

sub load_schema {
	print "Creating database tables.\n";
	my $uri = get_psql_uri();
	if ( system("psql $uri -e < db/create_tables.sql") != 0 ) {
		die "Can't create database tables\n";
	}
}

sub dropdb {
	my $uri = get_psql_uri($db_admin_username, $db_admin_password, 'postgres');
	my $cmd = "DROP DATABASE IF EXISTS $db_name;";
	if ( system(qq{psql $uri -tAec "$cmd"}) != 0 ) {
		die "Can't drop db $db_name\n";
	}
}

sub createdb {
	createuser();
	my $uri = get_psql_uri($db_admin_username, $db_admin_password, 'postgres');
	my $db_exists = `psql $uri -tAc "SELECT 1 FROM pg_database WHERE datname='$db_name'"`;
	if ($db_exists) {
		print "Database $db_name already exists\n";
		return;
	}

	my $cmd = "CREATE DATABASE $db_name;";
	if ( system(qq{psql $uri -tAec "$cmd"}) != 0 ) {
		die "Can't create db $db_name\n";
	}
}

sub createuser {
	my $uri = get_psql_uri($db_admin_username, $db_admin_password, 'postgres');
	my $user_exists = `psql $uri -tAc "SELECT 1 FROM pg_roles WHERE rolname='$db_username'"`;
	if ($user_exists) {
		print "Role $db_username already exists\n";
		return;
	}

	my $cmd = "CREATE USER $db_username WITH SUPERUSER CREATEROLE CREATEDB ENCRYPTED PASSWORD '$db_password'";
	if ( system(qq{psql $uri -tAc "$cmd"}) != 0 ) {
		die "Can't create user $db_username\n";
	}
}

sub dropuser {
	my $uri = get_psql_uri($db_admin_username, $db_admin_password, 'postgres');
	my $cmd = "DROP ROLE $db_username;";
	if ( system(qq{psql $uri -tAec "$cmd"}) != 0 ) {
		die "Can't drop user $db_username\n";
	}
}

sub showusers {
	my $uri = get_psql_uri();
	if ( system("psql $uri -ec '\\du';") != 0 ) {
		die "Can't show users";
	}
}

sub reverse_schema {

	my $db_info = Schema->get_dbinfo();
	my $user    = $db_info->{user};
	my $pass    = $db_info->{password};
	my $dsn     = Schema->get_dsn();
	make_schema_at(
		'Schema', {
			debug                   => 1,
			dump_directory          => './lib',
			overwrite_modifications => 1,
		},
		[ $dsn, $user, $pass ],
	);
}
