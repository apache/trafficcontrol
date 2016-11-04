package MojoPlugins::TrafficMonitorConnection;
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
use Utils::Rascal;

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		get_traffic_monitor_connection => sub {
			my $self = shift;
			my $args = shift;

			if ( !defined($args) || ref($args) ne "HASH" ) {
				confess("Supply a hashref of arguments");
			}

			my $hostname = undef;
			my $port     = undef;

			if ( exists( $args->{cdn} ) ) {
				my $cdn                 = $args->{cdn};
				my $traffic_monitor_row = undef;

				# this is the best query for the job, even though you can't search in it
				my $rs = $self->db->resultset('RascalHostsByCdn')->search();
				while ( my $row = $rs->next ) {
					next unless $cdn eq $row->cdn_name;
					$hostname = $row->host_name . "." . $row->domain_name;
					$port     = $row->tcp_port;
					last;
				}
			}
			elsif ( exists $args->{hostname} ) {
				$hostname = $args->{hostname};
				$port = exists( $args->{port} ) ? $args->{port} : 80;    # port is optional deaults to 80
			}
			else {
				confess("Supply a cdn or host in the argument hashref");
			}

			if ( !$hostname || !$port ) {
				return;
			}

			my $traffic_monitor_connection = new Utils::Rascal( $hostname, $port );
			my $proxy_param =
				$self->db->resultset('Parameter')->search( { -and => [ name => 'tm.traffic_mon_fwd_proxy', config_file => 'global' ] } )->single();
			if ( defined($proxy_param) ) {
				$traffic_monitor_connection->fwd_proxy( $proxy_param->value );
			}

			return $traffic_monitor_connection;
		}
	);

}

1;
