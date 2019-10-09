package MojoPlugins::TrafficRouterConnection;
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

use Mojo::Base 'Mojolicious::Plugin';
use Carp qw(cluck confess);
use Data::Dumper;
use Utils::CCR;

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		get_traffic_router_connection => sub {
			my $self = shift;
			my $args = shift;

			if ( !defined($args) || ref($args) ne "HASH" ) {
				confess("Supply a hashref of arguments");
			}

			my $hostname = undef;
			my $port     = undef;
			my $secure_port = undef;

			if ( exists( $args->{cdn} ) ) {
				my $cdn                = $args->{cdn};
				my $traffic_router_row = undef;

				# TODO !!
				
				if ( !defined($traffic_router_row) ) {
					confess("No TrafficRouter servers found for: $cdn");
				}
				$hostname = $traffic_router_row->host_name . "." . $traffic_router_row->domain_name;
			}
			elsif ( exists $args->{hostname} ) {
				$hostname = $args->{hostname};
			}
			else {
				confess("Supply a cdn or host in the argument hashref");
			}

			my $traffic_router_connection = undef;

			if ( defined( $args->{port} ) ) {
				$port = $args->{port};
			}
			else {
				my $hostonly = ( split( /\./, $hostname ) )[0];
				my $server = $self->db->resultset('Server')->search( { host_name => $hostonly } )->single();

				my $pp_secure_api_port =
					$self->db->resultset('ProfileParameter')
						->search( { -and => [ 'profile.id' => $server->profile->id, 'parameter.name' => 'secure.api.port', 'parameter.config_file' => 'server.xml' ] },
						{ prefetch => [ 'parameter', 'profile' ] } )->single();
					if ( defined($pp_secure_api_port) ) {
						$secure_port = $pp_secure_api_port->parameter->value;
					}

				if ( defined( $secure_port ) ) {
					$port = $secure_port;
					$traffic_router_connection = new Utils::CCR( $hostname, $port, 1 );
				}
				else {
					my $pp_api_port =
						$self->db->resultset('ProfileParameter')
							->search( { -and => [ 'profile.id' => $server->profile->id, 'parameter.name' => 'api.port', 'parameter.config_file' => 'server.xml' ] },
							{ prefetch => [ 'parameter', 'profile' ] } )->single();
						$port = $pp_api_port->parameter->value;
					$traffic_router_connection = new Utils::CCR( $hostname, $port);
				}
			}
			my $proxy_param =
				$self->db->resultset('Parameter')->search( { -and => [ name => 'tm.traffic_rtr_fwd_proxy', config_file => 'global' ] } )->single();
			if ( defined($proxy_param) ) {
				$traffic_router_connection->fwd_proxy( $proxy_param->value );
			}

			return $traffic_router_connection;
		}
	);

}

1;
