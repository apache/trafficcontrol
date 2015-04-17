package MojoPlugins::Server;
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
use Data::Dumper;
use Mojo::UserAgent;
use JSON;
use IO::Socket::SSL qw();
use LWP::UserAgent qw();
use File::Slurp;

use constant MAX_TRIES => 30;
##To track the active server we want to use
state $active_server = "NOT FOUND";

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(

		# This subroutine serves as a Retry Delegate for the specified helper by managing calls to the 'active' server.
		# The servers it connects to are defined in the 'Server' table with a status of 'ONLINE".
		#
		# If a remote server has connectivity issues, it will attempt to find the next available 'ONLINE' server
		# to communicate with, if it cannot find "any" it will then return an error.
		server_send_request => sub {
			my $self               = shift;
			my $server_type        = shift;
			my $helper_class       = shift || confess("Supply a Helper 'class'");
			my $method_function    = shift || confess("Supply a Helper class 'method'");
			my $schema_result_file = shift || confess("Supply a schema result file, ie: 'InfluxDBHostsOnline'");

			my $content;
			my $response;
			my $i           = 0;
			my $status_code = 200;
			my $message;

			my $active_server = activate_next_online_server( $self, $schema_result_file );
			if ( defined($active_server) ) {
				$helper_class->set_server($active_server);

				while ( ( $status_code <= 500 ) && ( $i <= MAX_TRIES ) ) {

					# This is the magic!! Dynamically invoke the method on the util to prevent
					# if-then-else
					$response    = $helper_class->$method_function($self);
					$status_code = $response->{_rc};
					$content     = $response->{_content};

					if ( $i >= MAX_TRIES ) {
						$message = "Couldn't connect to any " . $server_type . " servers.  Please make sure they are online!";
						$self->app->log->error( "Error: " . $message );
						return { response => $response, server => $active_server };
						last;    #bail
					}

					if ( $response->is_success ) {
						return { response => $response, server => $active_server };
					}
					elsif ( $status_code == 500 ) {
						$active_server = activate_next_online_server( $self, "InfluxDBHostsOnline" );
						$helper_class->set_server($active_server);
						if ( defined($active_server) ) {
							$self->app->log->warn( "Found BAD ONLINE server, skipping: " . $active_server );
						}
						else {
							$self->app->log->warn("No active server defined");
						}
					}
					else {
						$self->app->log->error( "Active Server Severe Error: " . $status_code . " - " . $content );
						return { response => $response, server => $active_server };
					}
					$i++;
				}
			}
			else {
				my $message =
					  "No "
					. $server_type
					. " servers are set to ONLINE in the database.  Please verify "
					. $server_type
					. " servers are online and reachable from Traffic Ops.";
				$response = HTTP::Response->new( 400, undef, HTTP::Headers->new, $message );
				return { response => $response, server => $active_server };
			}
		}
	);

	$app->renderer->add_helper();

}

# This subroutine only handles looking in the database (randomly) for
# 'ONLINE' servers based upon the specified '$schema_result_file' then
# making the discovered server the 'active' server.
sub activate_next_online_server {
	my $self               = shift;
	my $schema_result_file = shift;

	#get servers, if active_server then don't use, unless its the only one online.
	my @rs   = $self->db->resultset($schema_result_file)->search();
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

		# Perls way of maintain a 'state' variable.
		state $active_server;

		return $server;
	}
	else {
		$active_server = undef;
		return undef;
	}
}

1;
