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

sub get_tenant_heirarchy_depth {
	#return "undef" in case of error
	#a root tenant is of depth 0
	my $self = shift;
	my $tenants_data = shift;
	my $tenant_id = shift;

	if (!defined($tenants_data->{tenants_dict}{$tenant_id})) {
		return undef; #tenant does not exists #TODO -ask jeremy how to log
	}

	my $iter_id = $tenant_id;

	my $depth = 0; 
	while (defined($iter_id)) {
		$iter_id = $tenants_data->{tenants_dict}{$iter_id}{parent};
		$depth++; 
		if ($depth > $self->max_heirarchy_limit()) 		
		{
			return undef; #heirarchy limit #TODO -ask jeremy how to log
		}		
	}
	
	return $depth-1;
}

sub get_tenant_heirarchy_height {
	#return "undef" in case of error
	#a leaf tenant is of height 0
	my $self 	= shift;
	my $tenants_data = shift;
	my $tenant_id	= shift;

	if (!defined($tenants_data->{tenants_dict}{$tenant_id})) {
		return undef; #tenant does not exists #TODO -ask jeremy how to log
	}

	#calc tenant height
	my @tenants_list = reverse($self->get_hierarchic_tenants_list($tenants_data, $tenant_id));
	my %tenants_height = {};
	
	foreach my $tenant_row (@tenants_list) {
		my $tid = $tenant_row->id;
		$tenants_height{$tid} = 0;
	}
	
	foreach my $tenant_row (@tenants_list) {
		my $tid = $tenant_row->id;
		my $par_id = $tenant_row->parent_id;
		if (($tenants_height{$par_id}) < ($tenants_height{$tid}+1)) {
			$tenants_height{$par_id} = $tenants_height{$tid}+1;
		}
	}	 
	
	return $tenants_height{$tenant_id}; 
}


sub is_anchestor_of {
	#return "undef" in case of error
	my $self = shift;
	my $tenants_data = shift;
	my $anchestor_id = shift;
	my $descendant_id = shift;

	if (!defined($anchestor_id)) {
		return undef; #anchestor tenant is not defined #TODO -ask jeremy how to log
	}
	
	if (!defined($tenants_data->{tenants_dict}{$anchestor_id})) {
		return undef; #anchestor tenant does not exists #TODO -ask jeremy how to log
	}

	if (!defined($descendant_id)) {
		return undef; #descendant tenant is not defined #TODO -ask jeremy how to log
	}
	
	if (!defined($tenants_data->{tenants_dict}{$descendant_id})) {
		return undef; #descendant tenant does not exists #TODO -ask jeremy how to log
	}

	my $iter_id = $descendant_id;

	my $descendant_depth = 0; 
	while (defined($iter_id)) {
		if ($anchestor_id == $iter_id)
		{
			return 1;
		}
		$iter_id = $tenants_data->{tenants_dict}{$iter_id}{parent};
		$descendant_depth++; 
		if ($descendant_depth > $self->max_heirarchy_limit()) 		
		{#recursion limit
			return undef; #TODO -ask jeremy how to log 
		}		
	}
	
	return 0;
}

sub max_heirarchy_limit {
	my $self = shift;
	return 100;
}




##############################################################

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
	my $is_user_tenat_parent_of_resource = $self->is_anchestor_of($tenants_data, $user_tenant, $resource_tenant);
	if (!defined($is_user_tenat_parent_of_resource)) {
		#error - give access only to root tenant (so it can fix the problem)
		return $self->is_root_tenant($tenants_data, $user_tenant);
	}
	return $is_user_tenat_parent_of_resource;
}

1;
