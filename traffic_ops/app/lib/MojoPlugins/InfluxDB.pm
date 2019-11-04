package MojoPlugins::InfluxDB;
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
use Data::Dumper;
use Mojo::UserAgent;
use JSON;
use IO::Socket::SSL qw();
use LWP::UserAgent qw();
use File::Slurp;

use Extensions::TrafficStats::Connection::InfluxDBAdapter;

use constant SERVER_TYPE             => 'INFLUXDB';
use constant SCHEMA_FILE             => 'InfluxDBHostsOnline';
use constant INFLUXDB_CONF_FILE_NAME => 'influxdb.conf';
my $helper_class = eval {'Extensions::TrafficStats::Connection::InfluxDBAdapter'};

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		influxdb_write => sub {
			my $self         = shift;
			my $write_point  = shift || confess("Supply an InfluxDB 'write_point'");
			my $content_type = shift || "application/json";
			my $mode         = $self->app->mode;
			my $conf         = Utils::JsonConfig->load_conf( $mode, INFLUXDB_CONF_FILE_NAME );
			my $helper       = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->write( $write_point, $content_type ) }, SCHEMA_FILE );
		}
	);

	$app->renderer->add_helper(
		influxdb_query => sub {
			my $self    = shift;
			my $db_name = shift;
			my $query   = shift;
			my $mode    = $self->app->mode;
			my $conf    = Utils::JsonConfig->load_conf( $mode, INFLUXDB_CONF_FILE_NAME );
			my $helper  = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->query( $db_name, $query ) }, SCHEMA_FILE );
		}
	);

}
1;
