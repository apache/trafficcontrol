package UI::VisualStatus;

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
use Data::Dumper;
use JSON;

sub graphs {
	my $self = shift;
	my $match_string = $self->param('matchstring');
	
	my @cdn_names;
	my ( $ds_name, $loc_name, $host_name ) = split( /:/, $match_string );
	if ( $host_name ne 'all' ) {    # we want a specific host, it has to be in only one CDN
		my $server = $self->db->resultset('Server')->search( { host_name => $host_name }, { prefetch => 'cdn' } )->single();
		push( @cdn_names, $server->cdn->name );
	}
	elsif ( $ds_name ne 'all' ) {    # we want a specific DS, it has to be in only one CDN
		my $ds = $self->db->resultset('Deliveryservice')->search( { xml_id => $ds_name }, { prefetch => 'cdn' } )->single();
		push( @cdn_names, $ds->cdn->name );
	}
	else {                                             # we want all the CDNs with edges
		@cdn_names = $self->db->resultset('Server')->search({ 'type.name' => { -like => 'EDGE%' } }, { prefetch => [ 'cdn', 'type' ], group_by => 'cdn.name' } )->get_column('cdn.name')->all();
	}

	my $pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'visual_status_panel_1', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
	my $p1_url = defined($pparam) ? $pparam->parameter->value : undef;
	$pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'visual_status_panel_2', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
	my $p2_url = defined($pparam) ? $pparam->parameter->value : undef;
	$self->stash(
		cdn_names   => \@cdn_names,
		panel_1_url => $p1_url,
		panel_2_url => $p2_url
	);

	&navbarpage($self);
}


sub daily_summary {
	my $self = shift;

	my $pparam =
	$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'daily_bw_url', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
	my $bw_url = defined($pparam) ? $pparam->parameter->value : undef;
	$pparam =
		$self->db->resultset('ProfileParameter')
		->search( { -and => [ 'parameter.name' => 'daily_served_url', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
	my $served_url = defined($pparam) ? $pparam->parameter->value : undef;

	my @cdn_names = $self->db->resultset('Server')->search({ 'type.name' => { -like => 'EDGE%' } }, { prefetch => [ 'cdn', 'type' ], group_by => 'cdn.name' } )->get_column('cdn.name')->all();

	$self->stash(
		daily_bw_url => $bw_url,
		daily_served_url => $served_url,
		cdn_names => \@cdn_names
	);

	&navbarpage($self);
}

1;
