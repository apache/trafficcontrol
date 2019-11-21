package MojoPlugins::Riak;
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
use Net::Riak;
use Data::Dumper;
use Mojo::UserAgent;
use JSON;
use IO::Socket::SSL qw();
use LWP::UserAgent qw();
use Connection::RiakAdapter;
use File::Slurp;

use constant SERVER_TYPE => 'RIAK';
use constant SCHEMA_FILE => 'RiakHostsOnline';
my $helper_class = eval {'Connection::RiakAdapter'};

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		riak_stats => sub {
			my $self   = shift;
			my $conf   = load_conf($self);
			my $helper = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->stats() }, SCHEMA_FILE );
		}
	);

	$app->renderer->add_helper(
		riak_ping => sub {
			my $self   = shift;
			my $conf   = load_conf($self);
			my $helper = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->ping() }, SCHEMA_FILE );
		}
	);

	$app->renderer->add_helper(
		riak_put => sub {
			my $self         = shift;
			my $bucket       = shift;
			my $key          = shift;
			my $value        = shift;
			my $content_type = shift || "application/json";
			my $conf         = load_conf($self);
			my $helper       = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->put( $bucket, $key, $value, $content_type ) }, SCHEMA_FILE );
		}
	);

	$app->renderer->add_helper(
		riak_get => sub {
			my $self   = shift;
			my $bucket = shift;
			my $key    = shift;
			my $conf   = load_conf($self);
			my $helper = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->get( $bucket, $key ) }, SCHEMA_FILE );
		}
	);

	$app->renderer->add_helper(
		riak_delete => sub {
			my $self   = shift;
			my $bucket = shift;
			my $key    = shift;
			my $conf   = load_conf($self);
			my $helper = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->delete( $bucket, $key ) }, SCHEMA_FILE );
		}
	);

		$app->renderer->add_helper(
		riak_search => sub {
			my $self   = shift;
			my $index = shift;
			my $search_string    = shift;
			my $conf   = load_conf($self);
			my $helper = $helper_class->new( $conf->{user}, $conf->{password} );
			return $self->server_send_request( SERVER_TYPE, $helper, sub { $helper_class->search( $index, $search_string ) }, SCHEMA_FILE );
		}
	);
}

sub load_conf {
	local $/;    #Enable 'slurp' mode
	my $self = shift;
	my $mode = $self->app->mode;
	my $conf = "conf/$mode/riak.conf";
	return Utils::JsonConfig->new($conf);
}

1;
