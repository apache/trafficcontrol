package API::DeliveryServiceServer;
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
	#update to be ints
	$limit += 0;
	$page += 0;
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

sub assign_ds_to_cachegroup {
	my $self   = shift;
	my $cg_id  = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $cachegroup = $self->db->resultset('Cachegroup')->search( { id => $cg_id } )->single();
	if ( !defined($cachegroup) ) {
		return $self->not_found();
	}

	if ( ($cachegroup->type->name ne "EDGE_LOC") and ($cachegroup->type->name ne "ORG_LOC") ) {
		return $self->alert("cachegroup should be type EDGE_LOC or ORG_LOC.");
	}

	if ( !defined($params) ) {
		return $self->alert("parameters should in json format.");
	}

	if ( !defined($params->{deliveryServices}) ) {
		return $self->alert("parameter deliveryServices is must.");
	}

	if ( ref($params->{deliveryServices}) ne 'ARRAY' ) {
		return $self->alert("parameter deliveryServices must be array.");
	}

	my $cdn = "";
	my $servers = $self->db->resultset('Server')->search(
		{
			cachegroup => $cg_id,
			'type.name' => { -in => ['EDGE', 'ORG'] }
		},
		{ prefetch => ['type'] }
	);
	while ( my $server = $servers->next ) {
		if ($cdn eq "") {
			$cdn = $server->cdn_id;
		} elsif ($cdn ne $server->cdn_id) {
			return $self->alert("servers do not belong to same cdn.");
		}
	}

	my $deliveryservice_IDs = "";
	foreach my $ds_id (@{ $params->{deliveryServices} }) {
		my $ds = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
		if ( !defined($ds) ) {
			return $self->alert("deliveryservice with id $ds_id does not existed");
		}
		if ($cdn eq "") {
			$cdn = $ds->cdn_id;
		} elsif ($cdn ne $ds->cdn_id) {
			return $self->alert("servers/deliveryservices do not belong to same cdn.");
		}
		$deliveryservice_IDs = $deliveryservice_IDs . " " .  $ds_id;
	}

	$servers = $self->db->resultset('Server')->search(
		{
			cachegroup => $cg_id,
			'type.name' => { -in => ['EDGE', 'ORG'] }
		},
		{ prefetch => ['type'] }
	);

	my @server_names = ();
	while ( my $server = $servers->next ) {
		push(@server_names, $server->host_name);
		foreach my $ds_id (@{ $params->{deliveryServices} }) {
			my $find = $self->db->resultset('DeliveryserviceServer')->find(
				{
					deliveryservice => $ds_id,
					server          => $server->id
				}
			);

			if (!defined($find)) {
				my $insert = $self->db->resultset('DeliveryserviceServer')->create(
					{
						deliveryservice => $ds_id,
						server          => $server->id
					}
				);
				$insert->insert();

				if ($server->type->name eq 'EDGE') {
					my $ds = $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->single();
					&UI::DeliveryService::header_rewrite( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->edge_header_rewrite, "edge" );
					&UI::DeliveryService::regex_remap( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->regex_remap );
					&UI::DeliveryService::cacheurl( $self, $ds->id, $ds->profile, $ds->xml_id, $ds->cacheurl );
				}
				$self->app->log->info("assign server " . $server->id . " to ds " . $ds_id);
			}
		}
	}

	&log( $self, "assign servers in cache group $cg_id to deliveryservices $deliveryservice_IDs", "APICHANGE" );

	my $response;
	$response->{id} = $cg_id;
	$response->{serverNames} = \@server_names;
	$response->{deliveryServices} = $params->{deliveryServices};
	$self->success( $response, "Delivery services successfully assigned to all the servers of cache group $cg_id" );
}

1;
