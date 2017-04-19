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

use Env;
use Env qw(HOME);

my $usage = "\n"
	. "Usage:  $PROGRAM_NAME [--env (development|test|production|integration)] [arguments]\t\n\n"
	. "Example:  $PROGRAM_NAME --env=test reset\n\n"
	. "Purpose:  This script is used to manage database. The environments are\n"
	. "          defined in the dbconf.yml, as well as the database names.\n\n"
	. "NOTE: \n"
	. "Postgres Superuser: The 'postgres' superuser needs to be created to run $PROGRAM_NAME and setup databases.\n"
	. "If the 'postgres' superuser has not been created or password has not been set then run the following commands accordingly. \n\n"
	. "Create the 'postgres' user as a super user (if not created):\n\n"
	. "     \$ createuser postgres --superuser --createrole --createdb --login --pwprompt\n\n"
	. "Modify your $HOME/.pgpass file which allows for easy command line access by defaulting the user and password for the database\n"
	. "without prompts.\n\n"
	. " Postgres .pgpass file format:\n"
	. " hostname:port:database:username:password\n\n"
	. " ----------------------\n"
	. " Example Contents\n"
	. " ----------------------\n"
	. " *:*:*:postgres:your-postgres-password \n"
	. " *:*:*:traffic_ops:the-password-in-dbconf.yml \n"
	. " ----------------------\n\n"
	. " Save the following example into this file $HOME/.pgpass with the permissions of this file\n"
	. " so only $USER can read and write.\n\n"
	. "     \$ chmod 0600 $HOME/.pgpass\n\n"
	. "===================================================================================================================\n"
	. "$PROGRAM_NAME arguments:   \n\n"
	. "createdb  - Execute db 'createdb' the database for the current environment.\n"
	. "create_user  - Execute 'create_user' the user for the current environment (traffic_ops).\n"
	. "dropdb  - Execute db 'dropdb' on the database for the current environment.\n"
	. "down  - Roll back a single migration from the current version.\n"
	. "drop_user  - Execute 'drop_user' the user for the current environment (traffic_ops).\n"
	. "redo  - Roll back the most recently applied migration, then run it again.\n"
	. "reset  - Execute db 'dropdb', 'createdb', load_schema, migrate on the database for the current environment.\n"
	. "reverse_schema  - Reverse engineer the lib/Schema/Result files from the environment database.\n"
	. "seed  - Execute sql from db/seeds.sql for loading static data.\n"
	. "show_users  - Execute sql to show all of the user for the current environment.\n"
	. "status  - Print the status of all migrations.\n"
	. "upgrade  - Execute migrate then seed on the database for the current environment.\n";

my $environment = 'development';
my $db_protocol;

# This is defaulted to 'to_development' so
# you don't have to specify --env=development for dev workstations
my $db_name     = 'to_development';
my $db_super_user = 'postgres';
my $db_user = '';
my $db_password = '';
my $host_ip     = '';
my $host_port   = '';
GetOptions( "env=s" => \$environment );
$ENV{'MOJO_MODE'} = $environment;

parse_dbconf_yml_pg_driver();

STDERR->autoflush(1);
my $argument = shift(@ARGV);
if ( defined($argument) ) {
	if ( $argument eq 'createdb' ) {
		createdb();
	}
	elsif ( $argument eq 'dropdb' ) {
		dropdb();
	}
	elsif ( $argument eq 'create_user' ) {
		create_user();
	}
	elsif ( $argument eq 'drop_user' ) {
		drop_user();
	}
	elsif ( $argument eq 'show_users' ) {
		show_users();
	}
	elsif ( $argument eq 'reset' ) {
		create_user();
		dropdb();
		createdb();
		load_schema();
		migrate('up');
	}
	elsif ( $argument eq 'upgrade' ) {
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

sub parse_dbconf_yml_pg_driver {
	my $db_conf       = LoadFile('db/dbconf.yml');
	my $db_connection = $db_conf->{$environment};
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
	$db_user = $hash->{user};
	$db_password = $hash->{password};
	$db_name     = $hash->{dbname};
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
	local $ENV{PGPASSWORD} = $db_password;
	if ( system("psql -h $host_ip -p $host_port -d $db_name -U $db_user -e < db/seeds.sql") != 0 ) {
		die "Can't seed database\n";
	}
}

sub load_schema {
	print "Creating database tables.\n";
	local $ENV{PGPASSWORD} = $db_password;
	if ( system("psql -h $host_ip -p $host_port -d $db_name -U $db_user -e < db/create_tables.sql") != 0 ) {
		die "Can't create database tables\n";
	}
}

sub dropdb {
	print "Dropping database: $db_name\n";
	if ( system("dropdb -h $host_ip -p $host_port -U $db_super_user -e --if-exists $db_name;") != 0 ) {
		die "Can't drop db $db_name\n";
	}
}

sub createdb {
	my $db_exists = `psql -h $host_ip -U $db_super_user -p $host_port -tAc "SELECT 1 FROM pg_database WHERE datname='$db_name'"`;
	if ($db_exists) {
		print "Database $db_name already exists\n";
		return;
	}
	my $cmd = "createdb -h $host_ip -p $host_port -U $db_super_user -e --owner $db_user $db_name;";
	if ( system($cmd) != 0 ) {
		die "Can't create db $db_name\n";
	}

}

sub create_user {
	print "Creating user: $db_user\n";
	my $user_exists = `psql -h $host_ip -p $host_port -U $db_super_user -tAc "SELECT 1 FROM pg_roles WHERE rolname='$db_user'"`;

	if (!$user_exists) {
		my $cmd = "CREATE USER $db_user WITH LOGIN ENCRYPTED PASSWORD '$db_password'";
		if ( system(qq{psql -h $host_ip -p $host_port -U $db_super_user -etAc "$cmd"}) != 0 ) {
			die "Can't create user $db_user\n";
		}
	}
}

sub drop_user {
	if ( system("dropuser -h $host_ip -p $host_port -U $db_super_user-i -e $db_user;") != 0 ) {
		die "Can't drop user $db_user\n";
	}
}

sub show_users {
	if ( system("psql -h $host_ip -p $host_port -U $db_super_user -ec '\\du';") != 0 ) {
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

