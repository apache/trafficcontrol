package API::Keys;
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
use MojoPlugins::Response;
use Utils::Helper::ResponseHelper;

sub ping_riak {
	my $self = shift;
	my $response_container;
	my $response;

	if ( !&is_admin($self) ) {
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		$response_container = $self->riak_ping();
		$response           = $response_container->{"response"};
	}
	$self->success( $response->content );
}

1;
