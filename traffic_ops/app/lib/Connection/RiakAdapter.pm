package Connection::RiakAdapter;
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
use Net::Riak;
use Data::Dumper;
use Mojo::UserAgent;
use JSON;
use IO::Socket::SSL qw();
use File::Slurp;

# This Perl Module was needed to better support SSL for the 'Vault'
use LWP::UserAgent qw();
use LWP::ConnCache;

use constant RIAK_ROOT_URI => "riak";

# The purpose of this class is to provide for an easy method
# to 'mock' Riak calls.
my $ua;
my $conn_cache;
my $riak_server;
my $username;
my $password;

sub new {
	my $class = shift;
	$username = shift;
	$password = shift;
	my $self = bless {
		c           => $class,
		riak_server => $riak_server,
		username    => $username,
		password    => $password,
	}, $class;

	$ua = LWP::UserAgent->new();
	$ua->timeout(20);
	$ua->ssl_opts( verify_hostname => 0, SSL_verify_mode => 0x00 );
	if (!defined $conn_cache) {
	  $conn_cache = LWP::ConnCache->new( { total_capacity => 4096 } );
	}
	$ua->conn_cache($conn_cache);

	return $self;
}

sub set_server {
	my $self = shift;
	$riak_server = shift;
	$self->{'riak_server'} = $riak_server;
}

sub get_key_uri {
	my $self   = shift;
	my $bucket = shift || confess("Supply a bucket");
	my $key    = shift || confess("Supply a key");

	my @uri = ( RIAK_ROOT_URI, $bucket, $key );
	return File::Spec->join( "/", @uri );
}

sub ping {
	my $self = shift;
	my $fqdn = $self->get_url("/ping");
	return $ua->get($fqdn);
}

sub stats {
	my $self = shift;
	my $fqdn = $self->get_url("/stats");
	return $ua->get($fqdn);
}

sub put {
	my $self         = shift;
	my $bucket       = shift || confess("Supply a bucket");
	my $key          = shift || confess("Supply a key");
	my $value        = shift || confess("Supply a value");
	my $content_type = shift || 'application/json';

	my $key_uri = $self->get_key_uri( $bucket, $key );
	my $fqdn = $self->get_url($key_uri);

	return $ua->put( $fqdn, Content => $value, 'Content-Type'=> $content_type );
}

sub delete {
	my $self    = shift;
	my $bucket  = shift || confess("Supply a bucket");
	my $key     = shift || confess("Supply a key");
	my $key_uri = $self->get_key_uri( $bucket, $key );
	my $key_ctx = $self->get_url($key_uri);
	return $ua->delete( $key_ctx );
}

sub get {
	my $self   = shift;
	my $bucket = shift || confess("Supply a bucket");
	my $key    = shift || confess("Supply a key");

	my $key_uri = $self->get_key_uri( $bucket, $key );
	my $fqdn = $self->get_url($key_uri);
	return $ua->get($fqdn);
}

sub search {
	my $self   = shift;
	my $index = shift || confess("Supply a search index");
	my $search_string    = shift || confess("Supply a search string");

	my $key_uri = "/search/query/$index?wt=json&" . $search_string;
	my $fqdn = $self->get_url($key_uri);
	return $ua->get($fqdn);
}

#MOJOPlugins/Riak
sub get_url {
	my $self = shift;
	my $uri = shift || "";

	my $url;
	my $base_url = "https://$username:$password@";

	if ( $uri !~ m/^\// ) {
		$uri = "/" . $uri;
	}

	return $base_url . $riak_server . $uri;
}

1;
