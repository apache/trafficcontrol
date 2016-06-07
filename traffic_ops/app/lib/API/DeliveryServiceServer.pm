package API::DeliveryServiceServer;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Utils::Helper;

sub index {
	my $self = shift;
	my @data;
	my $orderby = $self->param('orderby') || "deliveryservice";

	# defaulted pagination and limits because there are 38129 rows in this table and counting...
	my $page  = $self->param('page')  || 1;
	my $limit = $self->param('limit') || 20;
	my $rs_data = $self->db->resultset("DeliveryserviceServer")->search( undef, { page => $page, rows => $limit, order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"deliveryService" => $row->deliveryservice->id,
				"server"          => $row->server->id,
				"lastUpdated"     => $row->last_updated,
			}
		);
	}
	$self->success( \@data, undef, $orderby, $limit, $page );
}

sub domains {
	my $self = shift;
	my @data;

	my @ccrprofs = $self->db->resultset('Profile')->search( { name => { -like => 'CCR%' } } )->get_column('id')->all();
	my $rs_pp =
		$self->db->resultset('ProfileParameter')
		->search( { profile => { -in => \@ccrprofs }, 'parameter.name' => 'domain_name', 'parameter.config_file' => 'CRConfig.json' },
		{ prefetch => [ 'parameter', 'profile' ] } );
	while ( my $row = $rs_pp->next ) {
		push(
			@data, {
				"domainName"         => $row->parameter->value,
				"parameterId"        => $row->parameter->id,
				"profileId"          => $row->profile->id,
				"profileName"        => $row->profile->name,
				"profileDescription" => $row->profile->description,
			}
		);

	}
	$self->success( \@data );
}
1;
