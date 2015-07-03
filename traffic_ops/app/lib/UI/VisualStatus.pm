package UI::VisualStatus;

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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';

sub graphs {
	my $self = shift;

	my $pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'visual_status_panel_1', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
	my $p1_url = $pparam->parameter->value;
	$pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'visual_status_panel_2', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
	my $p2_url = $pparam->parameter->value;
	$self->stash(
		panel_1_url => $p1_url,
		panel_2_url => $p2_url
	);

	&navbarpage($self);
}

sub daily_summary {
	my $self = shift;

	my @cdn_names;
	my $rs = $self->db->resultset('Parameter')->search( { name => 'CDN_name' } );
	while ( my $row = $rs->next ) {
		push( @cdn_names, $row->value );
	}

	my $tool_instance =
		$self->db->resultset('Parameter')->search( { -and => [ name => 'tm.instance_name', config_file => 'global' ] } )->get_column('value')->single();

	$self->stash(
		cdn_names     => \@cdn_names,
		tool_instance => $tool_instance,
		graph_page    => 1,
	);

	&navbarpage($self);
}

1;
