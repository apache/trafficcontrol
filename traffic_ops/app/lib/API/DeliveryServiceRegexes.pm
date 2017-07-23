package API::DeliveryServiceRegexes;
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
use Utils::Tenant;
use UI::DeliveryService;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Common::ReturnCodes qw(SUCCESS ERROR);
use Validate::Tiny ':all';

sub all {
	my $self = shift;

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();

	my $rs;
	if ( &is_privileged($self) or $tenant_utils->use_tenancy()) {
		$rs = $self->db->resultset('Deliveryservice')->search( undef, { prefetch => [ 'cdn', { 'deliveryservice_regexes' => { 'regex' => 'type' } } ], order_by => 'xml_id' } );

		my @regexes;
		while ( my $ds = $rs->next ) {
			if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
				next;
			}
			my $re_rs = $ds->deliveryservice_regexes;
			my @matchlist;
			while ( my $re_row = $re_rs->next ) {
				push(
					@matchlist, {
						type      => $re_row->regex->type->name,
						pattern   => $re_row->regex->pattern,
						setNumber => $re_row->set_number,
					}
				);
			}
			my $delivery_service->{dsName} = $ds->xml_id;
			$delivery_service->{regexes} = \@matchlist;
			push( @regexes, $delivery_service );
		}

		return $self->success( \@regexes );
	}
	else {
		return $self->forbidden("Forbidden. Insufficent privileges.");
	}

}

sub index {
	my $self  = shift;
	my $ds_id = $self->param('dsId');

	my $ds = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
	if ( !defined($ds) ) {
		return $self->not_found();
	}

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();
	if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
		return $self->forbidden("Forbidden. The delivery-service is not available to the user's tenant");
	}

	my %criteria;
	$criteria{'deliveryservice'} = $ds_id;

	my $rs_data = $self->db->resultset("DeliveryserviceRegex")->search( \%criteria, { prefetch => [ { 'regex' => 'type' } ], order_by => 'me.set_number' } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"        => $row->regex->id,
				"pattern"   => $row->regex->pattern,
				"type"      => $row->regex->type->id,
				"typeName"  => $row->regex->type->name,
				"setNumber" => $row->set_number,
			}
		);
	}
	$self->success( \@data );
}

sub show {
	my $self     = shift;
	my $ds_id    = $self->param('dsId');
	my $regex_id = $self->param('id');

	my $ds_regex = $self->db->resultset('DeliveryserviceRegex')->search( { deliveryservice => $ds_id, regex => $regex_id } );
	if ( !defined($ds_regex) ) {
		return $self->not_found();
	}

	my $ds = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
	if ( !defined($ds) ) {
		return $self->not_found();
	}

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();
	if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
		return $self->forbidden("Forbidden. The delivery-service is not available to the user's tenant");
	}

	my %criteria;
	$criteria{'deliveryservice'} = $ds_id;
	$criteria{'regex'}           = $regex_id;

	my $rs_data = $self->db->resultset("DeliveryserviceRegex")->search( \%criteria, { prefetch => [ { 'regex' => 'type' } ] } );
	my @data    = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"        => $row->regex->id,
				"pattern"   => $row->regex->pattern,
				"type"      => $row->regex->type->id,
				"typeName"  => $row->regex->type->name,
				"setNumber" => $row->set_number,
			}
		);
	}
	$self->success( \@data );
}

sub create {
	my $self   = shift;
	my $ds_id  = $self->param('dsId');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $ds = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
	if ( !defined($ds) ) {
		return $self->not_found();
	}
	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();
	if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
        #unlike other places, here we return 403 and not 400, as the path itself (including the DS) is forbidden
		return $self->forbidden("Forbidden. The delivery-service is not available to the user's tenant");
	}

	my ( $is_valid, $result ) = $self->is_regex_valid($params);

	if ( !$is_valid ) {
		return $self->alert($result);
	}

	my $values = {
		pattern => $params->{pattern},
		type    => $params->{type},
	};

	my $rs_regex = $self->db->resultset('Regex')->create($values)->insert();
	if ($rs_regex) {

		# now insert the regex into the deliveryservice_regex table along with set number
		my $rs_ds_regex = $self->db->resultset('DeliveryserviceRegex')
			->create( { deliveryservice => $ds_id, regex => $rs_regex->id, set_number => $params->{setNumber} } )->insert();

		my $response;
		$response->{id}        = $rs_regex->id;
		$response->{pattern}   = $rs_regex->pattern;
		$response->{type}      = $rs_regex->type->id;
		$response->{typeName}  = $rs_regex->type->name;
		$response->{setNumber} = $rs_ds_regex->set_number;

		&log( $self, "Regex created [ " . $rs_regex->pattern . " ] for deliveryservice: " . $ds_id, "APICHANGE" );

		return $self->success( $response, "Delivery service regex creation was successful." );
	}
	else {
		return $self->alert("Delivery service regex creation failed.");
	}

}

sub update {
	my $self     = shift;
	my $ds_id    = $self->param('dsId');
	my $regex_id = $self->param('id');
	my $params   = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $ds = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
	if ( !defined($ds) ) {
		return $self->not_found();
	}
	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db();
	if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
		return $self->forbidden("Forbidden. The delivery-service is not available to the user's tenant");
	}

	my ( $is_valid, $result ) = $self->is_regex_valid($params);

	if ( !$is_valid ) {
		return $self->alert($result);
	}

	my $ds_regex = $self->db->resultset('DeliveryserviceRegex')->search( { deliveryservice => $ds_id, regex => $regex_id } );
	if ( !defined($ds_regex) ) {
		return $self->not_found();
	}

	my $values = {
		pattern => $params->{pattern},
		type    => $params->{type},
	};

	my $regex = $self->db->resultset('Regex')->find( { id => $regex_id } )->update($values);
	if ($regex) {

		# now update the set_number in the deliveryservice_regex table
		$ds_regex->update( { set_number => $params->{setNumber} } );

		my $response;
		$response->{id}        = $regex->id;
		$response->{pattern}   = $regex->pattern;
		$response->{type}      = $regex->type->id;
		$response->{typeName}  = $regex->type->name;
		$response->{setNumber} = $params->{setNumber};

		&log( $self, "Regex updated [ " . $regex->pattern . " ] for deliveryservice: " . $ds_id, "APICHANGE" );

		return $self->success( $response, "Delivery service regex update was successful." );
	}
	else {
		return $self->alert("Delivery service regex update failed.");
	}

}

sub delete {
	my $self     = shift;
	my $ds_id    = $self->param('dsId');
	my $regex_id = $self->param('id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $ds = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
	if ( !defined($ds) ) {
        return $self->not_found();
	}
	else{
		my $tenant_utils = Utils::Tenant->new($self);
		my $tenants_data = $tenant_utils->create_tenants_data_from_db();
		if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
			return $self->forbidden("Forbidden. The delivery-service is not available to the user's tenant");
		}
	}

	my $ds_regex = $self->db->resultset('DeliveryserviceRegex')->search( { deliveryservice => $ds_id, regex => $regex_id } );
	if ( !defined($ds_regex) ) {
		return $self->not_found();
	}

	my $count = $self->db->resultset('RegexesForDeliveryService')->search( {}, { bind => [$ds_id] } )->count;
	if ( $count < 2 ) {
		return $self->alert("A delivery service must have at least one regex.");
	}

	my $regex = $self->db->resultset('Regex')->find( { id => $regex_id } )->delete();
	if ($regex) {

		# now delete the entry in the deliveryservice_regex table
		$ds_regex->delete();

		&log( $self, "Regex deleted [ " . $regex->pattern . " ] for deliveryservice: " . $ds_id, "APICHANGE" );

		return $self->success_message("Delivery service regex delete was successful.");
	}
	else {
		return $self->alert("Delivery service regex delete failed.");
	}

}

sub is_regex_valid {
	my $self   = shift;
	my $params = shift;

	if ( !$self->is_valid_regex_type( $params->{type} ) ) {
		return ( 0, "Invalid regex type" );
	}

	my $rules = {
		fields => [qw/pattern type setNumber/],

		# Validation checks to perform
		checks => [
			pattern   => [ is_required("is required") ],
			type      => [ is_required("is required") ],
			setNumber => [ is_required("is required") ],
		]
	};

	# Validate the input against the rules
	my $result = validate( $params, $rules );

	if ( $result->{success} ) {
		return ( 1, $result->{data} );
	}
	else {
		return ( 0, $result->{error} );
	}
}

sub is_valid_regex_type {
	my $self    = shift;
	my $type_id = shift;

	my $rs = $self->db->resultset("Type")->find( { id => $type_id } );
	if ( defined($rs) && ( $rs->use_in_table eq "regex" ) ) {
		return 1;
	}
	return 0;
}

1;
