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
use DBIx::Class::Schema::Loader qw/make_schema_at/;

my $usage = "\n"
	. "Usage:  $PROGRAM_NAME [--env (development|test|production|integration)] [arguments]\t\n\n"
	. "Example:  $PROGRAM_NAME --env=test reset\n\n"
	. "Purpose:  This script is used to manage database. The environments are\n"
	. "          defined in the dbconf.yml, as well as the database names.\n\n"
	. "arguments:   \n\n"
	. "create  - Execute db 'create' the database for the current environment.\n"
	. "down  - Roll back a single migration from the current version.\n"
	. "drop  - Execute db 'drop' on the database for the current environment.\n"
	. "redo  - Roll back the most recently applied migration, then run it again.\n"
	. "reset  - Execute db drop, create, load_schema, migrate on the database for the current environment.\n"
	. "reverse_schema  - Reverse engineer the lib/Schema/Result files from the environment database.\n"
	. "seed  - Execute sql from db/seeds.sql for loading static data.\n"
	. "setup  - Execute db drop, create, load_schema, migrate, seed on the database for the current environment.\n"
	. "status  - Print the status of all migrations.\n"
	. "upgrade  - Execute migrate then seed on the database for the current environment.\n";

my $environment = 'development';
my $db_string   = 'tcp:127.0.0.1:3306*twelve_monkeys/twelve_monkeys/twelve';
my $db_protocol;

# This is defaulted to 'to_development' so
# you don't have to specify --env=development for dev workstations
my $db_name     = 'to_development';
my $db_username = 'to_user';
my $db_password = '';
my $host_ip     = '';
my $host_port   = '';
GetOptions( "env=s" => \$environment );
$ENV{'MOJO_MODE'} = $environment;
my $dbh = Schema->database_handle;

sub database_handle_no_dbname {
	my $self    = shift;
	my $db_info = Schema->get_dbinfo();
	our $user = $db_info->{user};
	our $pass = $db_info->{password};
	our $dsn  = get_dsn_no_dbname();

	return DBI->connect( $dsn, $user, $pass, { AutoCommit => 1 } );
}

sub get_dsn_no_dbname {
	my $self = shift;

	my $db_info = Schema->get_dbinfo();
	our $dbname   = $db_info->{dbname};
	our $hostname = $db_info->{hostname};
	our $port     = $db_info->{port};
	our $type     = $db_info->{type};
	our $dsn      = "DBI:$type:host=$hostname;port=$port";
}

my $dbh_no_dbname = database_handle_no_dbname;

my $schema = Schema->connect_to_database;

parse_dbconf_yml_mysql_driver();

STDERR->autoflush(1);
my $argument = shift(@ARGV);
if ( defined($argument) ) {
	if ( $argument eq 'create' ) {
		create();
	}
	elsif ( $argument eq 'drop' ) {
		drop();
	}
	elsif ( $argument eq 'reset' ) {
		drop();
		create();
		load_schema();
		migrate('up');
	}
	elsif ( $argument eq 'upgrade' ) {
		migrate('up');
		seed();
	}
	elsif ( $argument eq 'setup' ) {
		drop();
		create();
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
}
else {
	print $usage;
}

exit(0);

sub parse_dbconf_yml_mysql_driver {
	my ($db_conf) = @_;

	my @files = File::Find::Rule->file()->name('dbconf.yml')->in('.');
	my $meta  = Parse::CPAN::Meta->load_file( $files[0] );

	# Mysql connection string parsing, if db changes added switch here
	# based on the driver: in dbconf.xml
	my $db_connection_string = $meta->{$environment}->{open};

	my @options = split( ':', $db_connection_string );
	$db_protocol = $options[0];

	$host_ip         = $options[1];
	my @port_options = split( '\*', $options[2] );
	$host_port       = $port_options[0];

	my $rest_of_options = $port_options[1];
	my @db_options = split( '/', $rest_of_options );
	$db_name     = $db_options[0];
	$db_username = $db_options[1];
	$db_password = $db_options[2];
}

sub migrate {
	my ($command) = @_;
	print "Migrating database...\n";
	system( 'goose --env=' . $environment . ' ' . $command );
}

sub seed {
	print "Seeding database...\n";
	system("mysql -h $host_ip -P $host_port $db_name -u$db_username -p$db_password < db/seeds.sql");
}

sub load_schema {
	print "Creating database tables...\n";
	system("mysql -h $host_ip -P $host_port $db_name -u$db_username -p$db_password < db/create_tables.sql");
}

sub drop {
	my $q = "drop database " . $db_name;
	print "Executing '" . $q . "'\n";
	if ( defined($dbh) ) {
		my $stmt = $dbh_no_dbname->prepare($q);
		$stmt->execute();
	}
	else {
		print "No database to drop\n";
	}
}

sub create {
	my $q = "create database " . $db_name;
	print "Executing '" . $q . "'\n";
	my $stmt = $dbh_no_dbname->prepare($q);
	$stmt->execute();
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
