package UI::HwInfo;
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

sub readhwinfo {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "serverid";
	my $limit   = $self->param('limit')   || 1000;

	# get list of servers in one query
	my $rs_data = $self->db->resultset("Hwinfo")->search( undef, { prefetch => [ 'serverid' ], rows => $limit } );

	while ( my $row = $rs_data->next ) {
		my $id = $row->id;
		push(
			@data, {
				"serverid"     => $row->serverid->host_name,
				"description"  => $row->description,
				"val"          => $row->val,
				"last_updated" => $row->last_updated,
			}
		);
	}

	#$self->deprecate( \@data, { see => '/api/1.1/locations.json' } );
	$self->render( json => \@data );
}

1;
