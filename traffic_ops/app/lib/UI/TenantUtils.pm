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
use NetAddr::IP;
use Data::Dumper;
use Switch;
use UI::Utils;


sub new {
	my $class = shift;
	my $self  = {
	        context => shift,
		user_tenant_id => -1,		
		tenants_dict => undef,
		root_tenants => undef,
		order_by => -1,
		ordered_by => undef,
	};
	return bless $self, $class;
}


sub current_user_tenant {
	my $self = shift;
	if ($self->{user_tenant_id} == -1) 
	{
		$self->{user_tenant_id} = $self->{context}->db->resultset('TmUser')->search( { username => $self->{context}->current_user()->{username} } )->get_column('tenant_id')->single();
	}
	return $self->{user_tenant_id};
}

sub get_hierarchic_tenants_list {
	my $self = shift;
	my $tree_root = shift;	
	my $order_by = shift;	
	
	$self->_init_tenants_if_needed($order_by);

	my @stack = ();
	if (defined($tree_root)){
		push (@stack, $tree_root);
	}
	else {
		push (@stack, reverse(@{$self->{root_tenants}}));
	}

	my @result = ();
	while (@stack) {
		my $tenant_id = pop @stack;
		push (@result, $self->{tenants_dict}->{$tenant_id}{row});
		push (@stack, reverse(@{$self->{tenants_dict}->{$tenant_id}{children}}));
	}
		
	return @result;	
}

sub is_root_tenant {
	my $self = shift;
	my $tenant_id = shift;
	
	if (!defined($tenant_id)) {
		return 0;
	}
	
	if (defined($self->{tenants_dict})) {
		return !(defined($self->{tenants_dict}{$tenant_id}{parent}));
	}
	return !defined($self->{context}->db->resultset('Tenant')->search( { id => $tenant_id } )->get_column('parent_id')->single()); 
}

sub is_tenant_resource_readable {
    my $self = shift;
    my $resource_tenancy = shift;
    
    return _is_resource_accessable ($self, $resource_tenancy, "r");
}

sub is_tenant_resource_writeable {
    my $self = shift;
    my $resource_tenancy = shift;
    
    return _is_resource_accessable ($self, $resource_tenancy, "w");
}



##############################################################

sub _init_tenants {
	my $self = shift;
	$self->{order_by} = shift || "name";#some default
	my $tenants_table = $self->{context}->db->resultset("Tenant")->search( undef, { order_by => $self->{order_by} });
	
	$self->{ordered_by} = ();
	$self->{tenants_dict} = {};
	while ( my $row = $tenants_table->next ) {
		push (@{ $self->{ordered_by} }, $row->id);
		$self->{tenants_dict}->{$row->id} = {
			row => $row,
			parent => $row->parent_id,
			children => (),
		}
	}
	
	$self->{root_tenants} = ();
	foreach my $key (@{ $self->{ordered_by} }) {
		my $value = $self->{tenants_dict}->{$key};
		my $parent = $value->{parent};
		if (!defined($parent))
		{
			push @{ $self->{root_tenants} }, $key;
		}
		else{
			push @{ $self->{tenants_dict}->{$parent}{children} }, $key;
		}
	}
}

sub _init_tenants_if_needed {
	my $self = shift;
	my $order_by = shift;
	if (($self->{order_by} == -1) || (defined($order_by) && $order_by != $self->{order_by})) {
		## first run to build the list OR (the order is important AND is not the current order)
		$self->_init_tenants($order_by);
	}
}

sub _max_tenancy_heirarchy {
	my $self = shift;
	return 100;
}

sub _is_resource_accessable {
	my $self = shift;
	my $resource_tenant = shift;
	my $operation = shift;

	if (!defined($resource_tenant)) {
		#the object has no tenancy - opened for all
	        return 1;
    	}

	if (&is_ldap($self->{context})) {
		if ($operation eq "r") {
			#ldap user, can read all tenants - temporary for now as an LDAP user as no tenant and is part of the TC operator.
			# should be removed when LDAP is gone
			return 1;
		}
		#ldap user, has no tenancy, cannot write anything
		return 0;
	}
    	
    	my $user_tenant = current_user_tenant($self);
	if (!defined($user_tenant)) {
		#the user has no tenancy, - cannot approach items with tenancy
		return 0;
	}

	$self->_init_tenants_if_needed(undef);
	my $tenant_record = $self->{tenants_dict}->{$user_tenant};
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
	for (my $depth = 0; $depth < $self->_max_tenancy_heirarchy(); $depth++) {
	
		if (!defined($resource_tenant)){
			#reached top tenant, resource is not under the user tenancy
			return 0;
		}

        	if ($user_tenant == $resource_tenant) {
		    #resource has child tenancy of the user, operations are allowed
        	    return 1;
		}
		
		$resource_tenant =  $self->{tenants_dict}->{$resource_tenant}->{parent};
	};
	
	#not found - recursion limit, give only access to root tenant
	return $self->is_root_tenant(current_user_tenant($self));
}

1;
