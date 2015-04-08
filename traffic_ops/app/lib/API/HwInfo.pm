package API::HwInfo;
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

sub index {
	my $self = shift;
	my $orderby = $self->param('orderby') || "serverid";
	my @data;

	# get list of servers in one query
	my $rs_data = $self->db->resultset("Hwinfo")->search( undef, { prefetch => [ { 'serverid' => undef, } ], order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		my $id = $row->id;
		push(
			@data, {
				"serverId"       => $row->serverid->id,
				"serverHostName" => $row->serverid->host_name,
				"description"    => $row->description,
				"val"            => $row->val,
				"lastUpdated"    => $row->last_updated,
			}
		);
	}

	$self->success( \@data );
}

1;
