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

use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use MojoPlugins::Response;

my $finfo = __FILE__ . ":";

sub getTenantName {
	my $self 		= shift;
	my $tenant_id		= shift;
	return defined($tenant_id) ? $self->db->resultset('Tenant')->search( { id => $tenant_id } )->get_column('name')->single() : "n/a";
}

sub isRootTenant {
	my $self 	= shift;
	my $tenant_id	= shift;
	return !defined($self->db->resultset('Tenant')->search( { id => $tenant_id } )->get_column('parent_id')->single());
}

sub index {
	my $self 	= shift;

	my @data;
	my $orderby = $self->param('orderby') || "name";
	my $rs_data = $self->db->resultset("Tenant")->search( undef, {order_by => 'me.' . $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"active"       => \$row->active,
				"parentId"     => $row->parent_id,
			}
		);
	}
	$self->success( \@data );
}


sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("Tenant")->search( { 'me.id' => $id });
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"name"         => $row->name,
				"active"       => \$row->active,
				"parentId"     => $row->parent_id,
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

	if ( !defined($params) ) {
		return $self->alert("Parameters must be in JSON format.");
	}

	if ( !defined( $params->{name} ) ) {
		return $self->alert("Tenant name is required.");
	}
	
	if ( $params->{name} ne $self->getTenantName($id) ) {
	        my $name = $params->{name};
		my $existing = $self->db->resultset('Tenant')->search( { name => $name } )->get_column('name')->single();
		if ($existing) {
			return $self->alert("A tenant with name \"$name\" already exists.");
		}	
	}	

	if ( !defined( $params->{parentId}) && !$self->isRootTenant($id) ) {
		return $self->alert("Parent Id is required.");
	}
	
	my $is_active = $params->{active};
	
	if ( !$params->{active} && $self->isRootTenant($id)) {
		return $self->alert("Root user cannot be in-active.");
	}
	

	if ( !defined($params->{parentId}) && !isRootTenant($id) ) {
		return $self->alert("Only the \"root\" tenant can have no parent.");
	}
	
	my $values = {
		name      => $params->{name},
		active    => $params->{active},
		parent_id => $params->{parentId}
	};

	my $rs = $tenant->update($values);
	if ($rs) {
		my $response;
		$response->{id}          = $rs->id;
		$response->{name}        = $rs->name;
		$response->{active}      = \$rs->active;
		$response->{parentId}    = $rs->parent_id;
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

	my $name = $params->{name};
	if ( !defined($name) ) {
		return $self->alert("Tenant name is required.");
	}

	my $parent_id = $params->{parentId};
	if ( !defined($parent_id) ) {
		return $self->alert("Parent Id is required.");
	}

	my $existing = $self->db->resultset('Tenant')->search( { name => $name } )->get_column('name')->single();
	if ($existing) {
		return $self->alert("A tenant with name \"$name\" already exists.");
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

	my $insert = $self->db->resultset('Tenant')->create($values);
	my $rs = $insert->insert();
	if ($rs) {
		my $response;
		$response->{id}          	= $rs->id;
		$response->{name}        	= $rs->name;
		$response->{active}        	= \$rs->active;
		$response->{parentId}           = $rs->parent_id;
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
	my $name = $self->db->resultset('Tenant')->search( { id => $id } )->get_column('name')->single();
	
	my $existing_child = $self->db->resultset('Tenant')->search( { parent_id => $id } )->get_column('name')->first();
	if ($existing_child) {
		return $self->alert("Tenant '$name' has children tenant(s): e.g '$existing_child'. Please update these tenants and retry.");
	}

	my $rs = $tenant->delete();
	if ($rs) {
		return $self->success_message("Tenant deleted.");
	} else {
		return $self->alert( "Tenant delete failed." );
	}
}


