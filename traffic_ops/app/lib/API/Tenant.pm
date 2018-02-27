package API::Tenant;
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
use Utils::Tenant;

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use MojoPlugins::Response;
use Validate::Tiny ':all';

my $finfo = __FILE__ . ":";

sub getTenantName {
	my $self 		= shift;
	my $tenant_id		= shift;
	return defined($tenant_id) ? $self->db->resultset('Tenant')->search( { id => $tenant_id } )->get_column('name')->single() : "n/a";
}

sub index {
	my $self 	= shift;	
	my $orderby = $self->param('orderby') || "name";

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db($orderby);

	my @data = ();
	my @tenants_list = $tenant_utils->get_hierarchic_tenants_list($tenants_data, undef);
	foreach my $row (@tenants_list) {
		if (!$tenant_utils->is_tenant_resource_accessible($tenants_data, $row->id)) {
            next;
        }
        push(
            @data, {
                "id"             => $row->id,
                "name"           => $row->name,
                "active"         => \$row->active,
                "parentId"       => $row->parent_id,
                "parentName"     => ( defined $row->parent_id ) ? $tenant_utils->get_tenant_by_id($tenants_data, $row->parent_id)->name : undef,
            }
        );
	}
	$self->success( \@data );
}


sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db(undef);

	my @data = ();
	my $rs_data = $self->db->resultset("Tenant")->search( { 'me.id' => $id });
	while ( my $row = $rs_data->next ) {
		if (!$tenant_utils->is_tenant_resource_accessible($tenants_data, $row->id)) {
            return $self->forbidden();
        }
        push(
            @data, {
                "id"           => $row->id,
                "name"         => $row->name,
                "active"       => \$row->active,
                "parentId"     => $row->parent_id,
                "parentName"   => ( defined $row->parent_id ) ? $tenant_utils->get_tenant_by_id($tenants_data, $row->parent_id)->name : undef,
            }
        );
	}
	$self->success( \@data );
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $tenant = $self->db->resultset('Tenant')->find( { id => $id } );
	if ( !defined($tenant) ) {
		return $self->not_found();
	}

	my ( $is_valid, $result ) = $self->is_tenant_valid($params);

	if ( !$is_valid ) {
		return $self->alert($result);
	}

	if ( $params->{name} ne $self->getTenantName($id) ) {
	        my $name = $params->{name};
		my $existing = $self->db->resultset('Tenant')->search( { name => $name } )->get_column('name')->single();
		if ($existing) {
			return $self->alert("A tenant with name \"$name\" already exists.");
		}	
	}	

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db(undef);

    if (!$tenant_utils->is_tenant_resource_accessible($tenants_data, $id)) {
        return $self->forbidden(); #Current owning tenant is not under user's tenancy
    }

	if ( !$params->{active} && $tenant_utils->is_root_tenant($tenants_data, $id)) {
		return $self->alert("Root tenant cannot be in-active.");
	}

	if ($params->{parentId} != $tenant->parent_id) {
		#parent replacement
		if (!$tenant_utils->is_tenant_resource_accessible($tenants_data, $tenant->parent_id)) {
			#Current owning tenant is not under user's tenancy
			return $self>alert("Invalid parent tenant change. The current tenant parent is not avaialble for you to edit");
		}
		if (!defined($tenant_utils->get_tenant_by_id($tenants_data, $params->{parentId}))) {
			return $self->alert("Parent tenant does not exists.");
		}
		if (!$tenant_utils->is_tenant_resource_accessible($tenants_data, $params->{parentId})) {
			#Parent tenant to be set is not under user's tenancy
			return $self->alert("Invalid parent tenant. This tenant is not available to you for parent assignment.");
		}
		my $parent_depth = $tenant_utils->get_tenant_heirarchy_depth($tenants_data, $params->{parentId});
		if (!defined($parent_depth))
		{
			return $self->alert("Failed to retrieve parent tenant depth.");
		}

		my $tenant_height = $tenant_utils->get_tenant_heirarchy_height($tenants_data, $id);
		if (!defined($tenant_height))
		{
			return $self->alert("Failed to retrieve tenant height.");
		}
	
		if ($parent_depth+$tenant_height+1 > $tenant_utils->max_heirarchy_limit())
		{
			return $self->alert("Parent tenant is invalid: heirarchy limit reached.");
		}
	
		if ($params->{parentId} == $id){
			return $self->alert("Parent tenant is invalid: same as updated tenant.");
		}

		my $is_tenant_achestor_of_parent = $tenant_utils->is_anchestor_of($tenants_data, $id, $params->{parentId});
		if (!defined($is_tenant_achestor_of_parent))
		{
			return $self->alert("Failed to check tenant and parent current relations.");
		}

		if ($is_tenant_achestor_of_parent)
		{
			return $self->alert("Parent tenant is invalid: a child of the updated tenant.");
		}		
	
	}


	#operation	
	my $values = {
		name      => $params->{name},
		active    => $params->{active},
		parent_id => $params->{parentId}
	};

	#$tenants_data is about to become outdated
	my $rs = $tenant->update($values);
	if ($rs) {
		my %idnames;
		my $response;

		my $rs_idnames = $self->db->resultset("Tenant")->search( undef, { columns => [qw/id name/] } );
		while ( my $row = $rs_idnames->next ) {
			$idnames{ $row->id } = $row->name;
		}

		$response->{id}          = $rs->id;
		$response->{name}        = $rs->name;
		$response->{active}      = $rs->active;
		$response->{parentId}    = $rs->parent_id;
		$response->{parentName}  = ( defined $rs->parent_id ) ? $idnames{ $rs->parent_id } : undef;
		$response->{lastUpdated} = $rs->last_updated;
		&log( $self, "Updated Tenant name '" . $rs->name . "' for id: " . $rs->id, "APICHANGE" );
		return $self->success( $response, "Tenant update was successful." );
	}
	else {
		return $self->alert("Tenant update failed.");
	}

}


sub create {
	my $self   = shift;
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my ( $is_valid, $result ) = $self->is_tenant_valid($params);

	if ( !$is_valid ) {
		return $self->alert($result);
	}

	my $parent_id = $params->{parentId};

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db(undef);
	
	if (!defined($tenant_utils->get_tenant_by_id($tenants_data, $params->{parentId}))) {
		return $self->alert("Parent tenant does not exists.");
	}

	if (!$tenant_utils->is_tenant_resource_accessible($tenants_data, $params->{parentId})) {
		return $self->alert("Invalid parent tenant. This tenant is not available to you for parent assignment.");
	}

    my $parent_depth = $tenant_utils->get_tenant_heirarchy_depth($tenants_data, $params->{parentId});

	if (!defined($parent_depth))
	{
		return $self->alert("Failed to retrieve parent tenant depth.");
	}
	
	if ($parent_depth+1 > $tenant_utils->max_heirarchy_limit()-1)
	{
		return $self->alert("Parent tenant is invalid: heirarchy limit reached.");
	}
	
	my $existing = $self->db->resultset('Tenant')->search( { name => $params->{name} } )->get_column('name')->single();
	if ($existing) {
		return $self->alert("A tenant with name " . $params->{name} . " already exists.");
	}

	my $is_active = exists($params->{active})? $params->{active} : 0; #optional, if not set use default
	
	if ( !$is_active && !defined($parent_id)) {
		return $self->alert("Root user cannot be in-active.");
	}
	
	my $values = {
		name 		=> $params->{name} ,
		active		=> $is_active,
		parent_id 	=> $params->{parentId}
	};

	#$tenants_data is about to become outdated
	my $insert = $self->db->resultset('Tenant')->create($values);
	my $rs = $insert->insert();
	if ($rs) {
		my %idnames;
		my $response;

		my $rs_idnames = $self->db->resultset("Tenant")->search( undef, { columns => [qw/id name/] } );
		while ( my $row = $rs_idnames->next ) {
			$idnames{ $row->id } = $row->name;
		}

		$response->{id}          	= $rs->id;
		$response->{name}        	= $rs->name;
		$response->{active}        	= $rs->active;
		$response->{parentId}		= $rs->parent_id;
		$response->{parentName}  	= ( defined $rs->parent_id ) ? $idnames{ $rs->parent_id } : undef;
		$response->{lastUpdated} 	= $rs->last_updated;

		&log( $self, "Created Tenant name '" . $rs->name . "' for id: " . $rs->id, "APICHANGE" );

		return $self->success( $response, "Tenant create was successful." );
	}
	else {
		return $self->alert("Tenant create failed.");
	}

}


sub delete {
	my $self = shift;
	my $id     = $self->param('id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $tenant = $self->db->resultset('Tenant')->find( { id => $id } );
	if ( !defined($tenant) ) {
		return $self->not_found();
	}	

	my $tenant_utils = Utils::Tenant->new($self);
	my $tenants_data = $tenant_utils->create_tenants_data_from_db(undef);
	
	if (!$tenant_utils->is_tenant_resource_accessible($tenants_data, $id)) {
		return $self->forbidden(); #tenant is not under user's tenancy
	}

	my $name = $tenant->name;
	
	my $existing_child = $self->db->resultset('Tenant')->search( { parent_id => $id }, {order_by => 'me.name' } )->get_column('name')->first();
	if ($existing_child) {
		return $self->alert("Tenant '$name' has children tenant(s): e.g '$existing_child'. Please update these tenants and retry.");
	}

	#The order of the below tests is intentional
	my $existing_ds = $self->db->resultset('Deliveryservice')->search( { tenant_id => $id }, {order_by => 'me.xml_id' })->get_column('xml_id')->first();
	if ($existing_ds) {
		return $self->alert("Tenant '$name' is assign with delivery-services(s): e.g. '$existing_ds'. Please update/delete these delivery-services and retry.");
	}

	my $existing_user = $self->db->resultset('TmUser')->search( { tenant_id => $id }, {order_by => 'me.username' })->get_column('username')->first();
	if ($existing_user) {
		return $self->alert("Tenant '$name' is assign with user(s): e.g. '$existing_user'. Please update these users and retry.");
	}

	#$tenants_data is about to become outdated
	my $rs = $tenant->delete();
	if ($rs) {
		return $self->success_message("Tenant deleted.");
	} else {
		return $self->alert( "Tenant delete failed." );
	}
}

sub is_tenant_valid {
	my $self   	= shift;
	my $params 	= shift;

	my $rules = {
		fields => [
			qw/name active parentId/
		],

		# Validation checks to perform
		checks => [
			name		=> [ is_required("is required") ],
			active		=> [ is_required("is required") ],
			parentId	=> [ is_required("is required"), is_like( qr/^\d+$/, "must be a positive integer" ) ],
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



