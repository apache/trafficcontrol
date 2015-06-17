package Utils::Redis;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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

use Carp qw(cluck confess);
use JSON;
use Data::Dumper;
use Mojo::UserAgent;
use File::Find;

use constant { TIMEOUT => 30, };

my $mojo;

sub new {
	my $self  = {};
	my $class = shift;
	$mojo = shift;
	return ( bless( $self, $class ) );
}

sub redis_connection_string {
	my $self = shift;
	my $rs   = $mojo->db->resultset('RedisHostsOnline')->search()->single();
	if ( defined($rs) ) {
		my $redis_db_host = $rs->host_name . "." . $rs->domain_name . ":" . $rs->tcp_port . " - " . $rs->status_name;
		return $rs->host_name . "." . $rs->domain_name . ":" . $rs->tcp_port;
	}
	else {
		$mojo->app->log->error("Could not find an ONLINE instance of Redis in the 'Servers'");
		return undef;
	}
}

sub redis_connect {
	my $self = shift;

	my $redis_connection_string = redis_connection_string();

	my $rm = Common::RedisFactory->new( $mojo, $redis_connection_string );
	return $rm->connection();
}

1;
