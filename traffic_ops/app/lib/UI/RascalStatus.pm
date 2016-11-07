package UI::RascalStatus;

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

sub health {
	my $self = shift;
	my $pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'all_graph_url', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
	my $ag_url = defined($pparam) ? $pparam->parameter->value : undef;
	$pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'cachegroup_graph_url', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )
		->single();
	my $cgg_url = defined($pparam) ? $pparam->parameter->value : undef;
	$pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'server_graph_url', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )
		->single();
	my $srvg_url = defined($pparam) ? $pparam->parameter->value : undef;
	$self->stash(
		cachegroup_graph_url => $cgg_url,
		all_graph_url        => $ag_url,
		server_graph_url     => $srvg_url,
	);

	&navbarpage($self);
}

1;
