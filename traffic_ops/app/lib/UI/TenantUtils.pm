package UI::TenantUtils;
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


#
# This class provide utilities to examine different tenancy aspects, 
# and specifically tenancy based access restrictions.
#
# The class itself is almost* stateless. However in order to reduce calls to DB, its
# API allows the user to create the "DATA" object (using "create_tenants_data_from_db") 
# and run utility functions over it.
#
# For now, until the current user tenant ID will come from the jwt, the current user tenant is taken from the DB.
# In order to reduce the number of calls from the DB, the current user tenant is taken in the class creation
#

use Data::Dumper;
use UI::Utils;


sub new {
	my $class = shift;
	my $context = shift;
	# For now, until the current user tenant ID will come from the jwt, the current user tenant is taken from the DB.
	my $current_user_tenant = $context->db->resultset('TmUser')->search( { username => $context->current_user()->{username} } )->get_column('tenant_id')->single();
	my $dbh = $context->db; 
	my $self  = {		
	        dbh => $dbh,
		# In order to reduce the number of calls from the DB, the current user tenant is taken in the class creation.
		# the below parameters are held temporarily until the info is taken from the jwt
	        current_user_tenant => $current_user_tenant,
	        is_ldap => $context->is_ldap(),
	};
	bless $self, $class;
	return $self;
}

sub create_tenants_data_from_db {
	my $self = shift;
	my $tenants_data = {
		tenants_dict => undef,
		root_tenants => undef,
		order_by => -1,
		ordered_by => undef,
	};
	$tenants_data->{order_by} = shift || "name";#some default
	
	my $tenants_table = $self->{dbh}->resultset("Tenant")->search( undef, { order_by => $tenants_data->{order_by} });
	
	$tenants_data->{ordered_by} = ();
	$tenants_data->{tenants_dict} = {};
	while ( my $row = $tenants_table->next ) {
		push (@{ $tenants_data->{ordered_by} }, $row->id);
		$tenants_data->{tenants_dict}->{$row->id} = {
			row => $row,
			parent => $row->parent_id,
			children => (),
		}
	}
	
	$tenants_data->{root_tenants} = ();
	foreach my $key (@{ $tenants_data->{ordered_by} }) {
		my $value = $tenants_data->{tenants_dict}->{$key};
		my $parent = $value->{parent};
		if (!defined($parent))
		{
			push @{ $tenants_data->{root_tenants} }, $key;
		}
		else{
			push @{ $tenants_data->{tenants_dict}->{$parent}{children} }, $key;
		}
	}
	
	return $tenants_data;
}

sub current_user_tenant {
	my $self = shift;
	return $self->{current_user_tenant};
}

sub get_tenant {
	my $self = shift;
	my $tenants_data = shift;
	my $tenant_id = shift;	
	
	return $tenants_data->{tenants_dict}->{$tenant_id}{row};
}

sub get_tenants_list {
	my $self = shift;
	my $tenants_data = shift;
	my $order_by = shift;	
	
	my @result = ();
	foreach my $tenant_id (@{ $tenants_data->{ordered_by} }) {
		push @result, $tenants_data->{tenants_dict}->{$tenant_id}{row};
	}

	return @result;	
}

sub get_hierarchic_tenants_list {
	my $self = shift;
	my $tenants_data = shift;
	my $tree_root = shift;	
	my $order_by = shift;	
	
	my @stack = ();
	if (defined($tree_root)){
		push (@stack, $tree_root);
	}
	else {
		push (@stack, reverse(@{$tenants_data->{root_tenants}}));
	}

	my @result = ();
	while (@stack) {
		my $tenant_id = pop @stack;
		push (@result, $tenants_data->{tenants_dict}->{$tenant_id}{row});
		push (@stack, reverse(@{$tenants_data->{tenants_dict}->{$tenant_id}{children}}));
	}
		
	return @result;	
}

sub is_root_tenant {
	my $self = shift;
	my $tenants_data = shift;
	my $tenant_id = shift;
	
	if (!defined($tenant_id)) {
		return 0;
	}
	
	if (defined($tenants_data->{tenants_dict})) {
		return !(defined($tenants_data->{tenants_dict}{$tenant_id}{parent}));
	}
	return !defined($self->{dbh}->resultset('Tenant')->search( { id => $tenant_id } )->get_column('parent_id')->single()); 
}

sub is_tenant_resource_readable {
	my $self = shift;
	my $tenants_data = shift;
	my $resource_tenancy = shift;
    
	return $self->_is_resource_accessable ($tenants_data, $resource_tenancy, "r");
}

sub is_tenant_resource_writeable {
	my $self = shift;
	my $tenants_data = shift;
	my $resource_tenancy = shift;
    
	return $self->_is_resource_accessable ($tenants_data, $resource_tenancy, "w");
}



##############################################################

sub _tenancy_heirarchy_limit {
	my $self = shift;
	return 100;
}

sub _is_resource_accessable {
	my $self = shift;
	my $tenants_data = shift;
	my $resource_tenant = shift;
	my $operation = shift;

	if (!defined($resource_tenant)) {
		#the object has no tenancy - opened for all
	        return 1;
    	}

	if ($self->{is_ldap}) {
		if ($operation eq "r") {
			#ldap user, can read all tenants - temporary for now as an LDAP user as no tenant and is part of the TC operator.
			# should be removed when LDAP is gone
			return 1;
		}
		#ldap user, has no tenancy, cannot write anything
		return 0;
	}
    	
    	my $user_tenant = $self->current_user_tenant();
	if (!defined($user_tenant)) {
		#the user has no tenancy, - cannot approach items with tenancy
		return 0;
	}

	my $tenant_record = $tenants_data->{tenants_dict}->{$user_tenant};
	my $is_active_tenant = $tenant_record->{row}->active;
	if (! $is_active_tenant) {
		#user tenant is in-active - cannot do any operation
		return 0;
	}

	if ($user_tenant == $resource_tenant) {
	    #resource has same tenancy of the user, operations are allowed
	    return 1;
	}

	#checking if the user tenant is an ancestor of the resource tenant
	for (my $depth = 0; $depth < $self->_tenancy_heirarchy_limit(); $depth++) {
	
		if (!defined($resource_tenant)){
			#reached top tenant, resource is not under the user tenancy
			return 0;
		}

        	if ($user_tenant == $resource_tenant) {
		    #resource has child tenancy of the user, operations are allowed
        	    return 1;
		}
		
		$resource_tenant =  $tenants_data->{tenants_dict}->{$resource_tenant}->{parent};
	};
	
	#not found - recursion limit, give only access to root tenant
	return $self->is_root_tenant($tenants_data, $user_tenant);
}

1;
