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
	};
	return bless $self, $class;
}

sub get_hierarchic_tenants_list {
	my $self = shift;
	my $tree_root = shift;	
	
	$self->_init_tenants_if_needed();

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

##############################################################

sub _init_tenants {
	my $self = shift;
	my $tenants_table = $self->{context}->db->resultset("Tenant")->search( undef, { order_by => "id" });
	
	my @ordered_by_id = ();
	$self->{tenants_dict} = {};
	while ( my $row = $tenants_table->next ) {
		push (@ordered_by_id, $row->id);
		$self->{tenants_dict}->{$row->id} = {
			row => $row,
			parent => $row->parent_id,
			children => (),
		}
	}
	
	$self->{root_tenants} = ();
	foreach my $key (@ordered_by_id) {
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
	if (!defined($self->{tenants_dict})) {
		$self->_init_tenants();
	}
}

1;
