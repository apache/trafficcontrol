package MojoPlugins::Riak;
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

use Mojo::Base 'Mojolicious::Plugin';
use utf8;
use Carp qw(cluck confess);
use UI::Utils;
use Net::Riak;
use Data::Dumper;
use Mojo::UserAgent;
use JSON;
use Utils::Helper::ResponseHelper;
use IO::Socket::SSL qw();
use LWP::UserAgent qw();
use Utils::Riak;
use File::Slurp;

use constant MAX_TRIES => 30;
##To track the active server we want to use
state $active_server;

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		riak_stats => sub {
			my $self = shift;
			return send_to_online_server( $self, sub { Utils::Riak->stats() } );
		}
	);

	$app->renderer->add_helper(
		riak_ping => sub {
			my $self = shift;
			return send_to_online_server( $self, sub { Utils::Riak->ping() } );
		}
	);

	$app->renderer->add_helper(
		riak_put => sub {
			my $self         = shift;
			my $bucket       = shift;
			my $key          = shift;
			my $value        = shift;
			my $content_type = shift || "application/json";
			return send_to_online_server( $self, sub { Utils::Riak->put( $bucket, $key, $value ) } );
		}
	);

	$app->renderer->add_helper(
		riak_get => sub {
			my $self   = shift;
			my $bucket = shift;
			my $key    = shift;
			return send_to_online_server( $self, sub { Utils::Riak->get( $bucket, $key ) } );
		}
	);

	$app->renderer->add_helper(
		riak_delete => sub {
			my $self   = shift;
			my $bucket = shift;
			my $key    = shift;

			return send_to_online_server( $self, sub { Utils::Riak->delete( $bucket, $key ) } );
		}
	);
}

sub send_to_online_server {
	my $self            = shift;
	my $method_function = shift || confess("Supply a Util::Riak method");
	my $content         = shift;

	my $response;
	my $i           = 0;
	my $status_code = 200;
	my $message;

	my $conf = load_conf($self);
	my $riak_util = Utils::Riak->new( $conf->{user}, $conf->{password} );

	#Find the ONLINE server count
	my $online_count = $self->db->resultset('Server')->search( { status => 2 } )->count();
	my $active_server = find_next_online_server($self);
	if ( defined($active_server) ) {
		$riak_util->set_server($active_server);

		# This logic allows us to move onto the next ONLINE server
		# in the event one of them fails.
		#while ( ( $i == 0 || $status_code == 500 ) && $i <= MAX_TRIES ) {

		while ( ( $status_code <= 500 ) && ( $i <= MAX_TRIES ) ) {

			#$self->app->log->debug("---------------------");

			# This is the magic!! Dynamically invoke the method on the util to prevent
			# if-then-elses
			$response    = $method_function->($self);
			$status_code = $response->{_rc};
			$content     = $response->{_content};

			#$self->app->log->debug( "status_code #-> " . $status_code );
			#$self->app->log->debug( "content #-> " . $content );
			if ( $i >= MAX_TRIES ) {
				$message = "Couldn't connect to Riak servers.  Please make sure they are online!";
				$self->app->log->debug( "message #-> " . $message );
				return { response => $response, server => $active_server };
				last;    #bail
			}

			#$self->app->log->debug( "server #-> " . Dumper($server) );
			#$self->app->log->debug( "response #-> " . Dumper($response) );
			if ( $response->is_success ) {
				return { response => $response, server => $active_server };
			}
			elsif ( $status_code == 500 ) {
				$active_server = find_next_online_server($self);
				$riak_util->set_server($active_server);
				$self->app->log->warn( "Found BAD ONLINE server, skipping: " . $active_server );
			}
			else {
				$self->app->log->error( "Active Server Severe Error: " . $status_code . " - " . $content );
				return { response => $response, server => $active_server };
			}
			$i++;
		}
	}
	else {
		my $message = "No Riak servers are set to ONLINE in the database.  Please verify Riak servers are online and reachable from Traffic Ops.";
		$response = HTTP::Response->new( 400, undef, HTTP::Headers->new, $message );
		return { response => $response, server => $active_server };

		#return $self->alert("No Riak servers are set to ONLINE in the database.  Please verify Riak servers are online and reachable from Traffic Ops.");
	}
}

sub load_conf {
	local $/;    #Enable 'slurp' mode
	my $self = shift;
	my $mode = $self->app->mode;
	my $conf = "conf/$mode/riak.conf";
	return Utils::JsonConfig->new($conf);
}

#MOJOPlugins/Riak
sub find_next_online_server {
	my $self = shift;

	#get servers, if active_server then don't use, unless its the only one online.
	my @rs   = $self->db->resultset('RiakHostsOnline')->search();
	my $size = @rs;
	if ( $size == 1 ) {
		my $server = $rs[0]->host_name . "." . $rs[0]->domain_name . ":" . $rs[0]->tcp_port;
		return $server;
	}
	elsif ( $size > 1 ) {
		my $server_index = int( rand($size) );
		my $server       = $rs[$server_index]->host_name . "." . $rs[$server_index]->domain_name . ":" . $rs[$server_index]->tcp_port;
		my $i            = 0;

		# Keep looking until we find a different server.
		while ( $server eq $active_server && $i < MAX_TRIES ) {
			$server_index = int( rand($size) );
			$server       = $rs[$server_index]->host_name . "." . $rs[$server_index]->domain_name . ":" . $rs[$server_index]->tcp_port;
			$i++;    #safeguard from inifinite loop
		}
		$active_server = $server;
		$self->app->log->debug( "CURRENT active_server #-> " . $active_server );
		state $active_server;
		return $server;
	}
	else {
		$active_server = undef;
		return undef;
	}
}

1;
