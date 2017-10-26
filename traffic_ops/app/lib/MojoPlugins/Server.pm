package MojoPlugins::Server;

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

use Mojo::Base 'Mojolicious::Plugin';
use utf8;
use Carp qw(cluck confess);
use UI::Utils;
use Data::Dumper;
use List::Util qw/shuffle/;

use constant MAX_TRIES => 20;
##To track the active server we want to use
state %active_server_for;

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
			my $response;
			my $active_server = $active_server_for{$schema_result_file};
			my @rs = randomize_online_servers( $self, $schema_result_file );
			if ( defined $active_server ) {
				if ( !grep { $_ eq $active_server } @rs ) {

					# active server no longer listed as available
					undef $active_server;
				}
				else {
					# remove active_server from list
					@rs = grep { $_ ne $active_server } @rs;

					# tack it to the end so it's not reused immediately, but still available if the only one that responds
					push @rs, $active_server;
				}
			}

			for my $server (@rs) {

				# This is the magic!! Dynamically invoke the method on the util to prevent
				# if-then-else
				$helper_class->set_server($server);
				$response = $helper_class->$method_function($self);
				my $status_code = $response->{_rc};
				if ( $response->is_success ) {
					$self->app->log->debug("Using server, $server");
					$active_server = $server;
					last;
				}

				if ( $status_code == 500 ) {
					$self->app->log->warn("Found BAD ONLINE server, $server -- skipping");
				}
				else {
					my $content = $response->{_content};
					$self->app->log->error( "Active Server Severe Error: " . $status_code . " - " . $server . " - " . $content );
				}
			}

			$active_server_for{$schema_result_file} = $active_server;
			if ( !defined $active_server ) {

				# modify response
				my $message =
					"No $server_type servers are available.  Please verify $server_type servers are set to ONLINE and are reachable from Traffic Ops.";
				$response = HTTP::Response->new( 400, undef, HTTP::Headers->new, $message );
			}
			return { response => $response, server => $active_server };
		}
	);
}

sub server_id {
	my $server = shift;
	my $id;
	if ( defined $server ) {
		$id = $server->host_name . '.' . $server->domain_name . ':' . $server->tcp_port;
	}
	return $id;
}

sub randomize_online_servers {
	my $self               = shift;
	my $schema_result_file = shift;

	my @rs = $self->db->resultset($schema_result_file)->search();
	@rs = map { server_id($_) } @rs;

	# if two or more, return shuffled list with current one removed
	return shuffle(@rs);
}

1;
