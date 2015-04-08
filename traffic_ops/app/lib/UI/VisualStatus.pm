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

	my $match_string = $self->param('matchstring');

	my @cdn_names;
	my $ds_capacity = 0;
	my ( $ds_name, $loc_name, $host_name ) = split( /:/, $match_string );
	if ( $host_name ne 'all' ) {    # we want a specific host, it has to be in only one CDN
		my $server = $self->db->resultset('Server')->search( { host_name => $host_name } )->single();
		my $param =
			$self->db->resultset('ProfileParameter')
			->search( { -and => [ 'parameter.name' => 'CDN_name', 'parameter.name' => 'CDN_name', 'me.profile' => $server->profile->id ] },
			{ prefetch => [ 'parameter', 'profile' ] } )->single();
		my $cdn_name = $param->parameter->value;
		push( @cdn_names, $cdn_name );
	}
	elsif ( $ds_name ne 'all' ) {    # we want a specific DS, it has to be in only one CDN
		my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_name } )->single();
		my $param =
			$self->db->resultset('ProfileParameter')
			->search( { -and => [ 'parameter.name' => 'CDN_name', 'parameter.name' => 'CDN_name', 'me.profile' => $ds->profile->id ] },
			{ prefetch => [ 'parameter', 'profile' ] } )->single();
		my $cdn_name = $param->parameter->value;
		push( @cdn_names, $cdn_name );
		$ds_capacity = $ds->global_max_mbps / 1000;    # everything is in kbps in the stats
	}
	else {                                             # we want all the CDNs
		my $rs = $self->db->resultset('Parameter')->search( { name => 'CDN_name' } );
		while ( my $row = $rs->next ) {
			push( @cdn_names, $row->value );
		}
	}
	$self->stash(
		cdn_names   => \@cdn_names,
		graph_page  => 1,
		matchstring => $match_string,
		ds_capacity => $ds_capacity,
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

	my $tool_instance = $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.instance_name', config_file => 'global' ] } )
		->get_column('value')->single();

	$self->stash(
		cdn_names  => \@cdn_names,
		tool_instance => $tool_instance,
		graph_page => 1,
	);

	&navbarpage($self);
}

1;
