package Extensions::TrafficStats::Connection::InfluxDBAdapter;
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
#
#

use utf8;
use Carp qw(cluck confess);
use UI::Utils;
use Data::Dumper;
use Mojo::UserAgent;
use JSON;
use IO::Socket::SSL qw();
use File::Slurp;
use URI::Escape;
use Mojolicious::Types;

# This Perl Module was needed to better support SSL
use LWP::UserAgent qw();
use constant APPLICATION_JSON => 'application/json';

# The purpose of this class is to provide for a wrapper
# and 'mock' to TrafficStats
my $ua;
my $influxdb_server;
my $influxdb_db_name;
my $username;
my $password;

sub new {
	my $class = shift;
	$username = shift;
	$password = shift;
	my $self = bless {
		c               => $class,
		influxdb_server => $influxdb_server,
		username        => $username,
		password        => $password,
	}, $class;

	$ua = LWP::UserAgent->new();

	# timeout is in seconds
	$ua->timeout(20);
	$ua->ssl_opts( verify_hostname => 0, SSL_verify_mode => 0x00 );

	return $self;
}

sub set_db_name {
	my $self = shift;
	$influxdb_db_name = shift;
	$self->{'influxdb_db_name'} = $influxdb_db_name;
}

sub set_server {
	my $self = shift;
	$influxdb_server = shift;
	$self->{'influxdb_server'} = $influxdb_server;
}

sub write {
	my $self         = shift;
	my $write_point  = shift || confess("Supply a write_point in the form ofa hash.");
	my $content_type = shift || APPLICATION_JSON;

	$ua->default_header( 'Content-Type' => $content_type );

	my $fqdn = $self->get_url("/write");

	my $write_point_data;
	if ( $content_type eq APPLICATION_JSON ) {
		$write_point_data = to_json($write_point);
	}
	else {
		confess("Only 'application/json' 'Content-Type' allowed\n");
	}
	return $ua->post( $fqdn, Content => $write_point_data );
}

sub query {
	my $self = shift;

	# db name should not be included when create influxdb databases
	my $db_name = shift;
	my $query   = shift || confess("Supply a key");
	my $pretty  = shift;

	my @uri;
	if ( defined($db_name) ) {
		push( @uri, "db=" . $db_name );
	}
	if ( defined($pretty) ) {
		push( @uri, "pretty=true" );
	}

	push( @uri, "q=" . uri_escape($query) );
	my $uri = join( "&", @uri );
	my $fqdn = $self->get_url( "/query?" . $uri );
	return $ua->get($fqdn);
}

sub get_url {
	my $self = shift;
	my $uri = shift || "";

	my $url;
	my $base_url = "http://";

	if ( !defined($influxdb_server) ) {
		confess("Please specify an influxdb_server with set_server()");
	}

	return $base_url . $influxdb_server . $uri;
}

sub get_url_https {
	my $self = shift;
	my $uri = shift || "";

	my $url;
	my $base_url = "https://$username:$password@";

	if ( $uri !~ m/^\// ) {
		$uri = "/" . $uri;
	}

	return $base_url . $influxdb_server . $uri;
}

1;
