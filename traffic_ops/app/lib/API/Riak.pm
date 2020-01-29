package API::Riak;
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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use UI::Server;
use JSON;
use Data::Dumper;

sub stats {
	my $self = shift;

	my $response = $self->riak_stats()->{'response'};
	my $content  = $response->content;
	return $self->deprecation_with_no_alternative(200, decode_json($content));
}

sub ping {
	my $self = shift;

	#my $response = $self->riak_stats();
	my $ping     = $self->riak_ping();
	my $server   = $ping->{'server'};
	my $response = $ping->{'response'};
	my $content  = $response->content;
	return $self->success( { server => $server, status => $content } );
}

sub get {
	my $self = shift;

	my $bucket = $self->param("bucket");
	my $key    = $self->param("key");

	my $riak_get = $self->riak_get( $bucket, $key );
	my $response = $riak_get->{'response'};
	my $content  = $response->{_content};
	if ( $response->is_success() ) {
		return $self->success( decode_json($content) );
	}
	else {
		$self->app->log->debug( "riak_get #-> " . Dumper($riak_get) );
		my $rc = $response->{_rc};
		return $self->alert( $content, $rc );
	}
}

1;
