package Utils::Tenant;
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
#
# A usage example - examing if a tenant can update a delivery-service:
# my $tenant_utils = Utils::Tenant->new($self);
# my $tenants_data = $tenant_utils->create_tenants_data_from_db();
# if (!$tenant_utils->is_ds_writeable($tenants_data, <resource_tenant>)) {
# 	return $self->forbidden(); #Parent tenant is not under user's tenancy
# }
#

use Data::Dumper;
use UI::Utils;

sub new {
    my $class               = shift;
    my $context             = shift;
    my $current_user_tenant = shift
      ; #optional - allowing the user tenancy to be set from outside, for testing capabilities
    my $dbh = shift
      ; #optional - allowing the DB handle to be set from outside, for testing capabilities

    if ( !defined($current_user_tenant) ) {

# For now, until the current user tenant ID will come from the jwt, the current user tenant is taken from the DB.
        $current_user_tenant =
          $context->db->resultset('TmUser')
          ->search( { username => $context->current_user()->{username} } )
          ->get_column('tenant_id')->single();
    }

    if ( !defined($dbh) ) {
        $dbh = $context->db;
    }

    #For the reviewer - this is should probably do + lock the "GLOBAL" profile to the root tenant only.
    #Otherwise, anyone can disable tenancy
    #However, As UTs does not have a "GLOBAL" prfile currently, and I'm not sure it is wise to add one, I'm
    #leaving this code out
    #
    #my $global_profile_id = $dbh->resultset("Parameter")->search( { name => 'GLOBAL'} )->get_column('id')->single();
    #my $global_profile_parameters = $dbh->resultset('ProfileParameter')->search( { profile => $global_profile_id } );
    #my $use_tenancy_value = $dbh->resultset("Parameter")->search( { id => { -in => $global_profile_parameters->get_column('parameter')->as_query },
    #                                                                config_file => 'global', name => 'use_tenancy' } )
    #    ->get_column('value')->single();

    # tenancy is always enabled
    my $use_tenancy = 1;

    my $self = {
        dbh     => $dbh,
        context => $context, #saving the context - use it only for log please...

# In order to reduce the number of calls from the DB, the current user tenant is taken in the class creation.
# the below parameters are held temporarily until the info is taken from the jwt
        current_user_tenant => $current_user_tenant,
        use_tenancy => $use_tenancy,
    };
    bless $self, $class;
    return $self;
}

sub create_tenants_data_from_db {
    my $self = shift;
    my $orderby = shift || "name";    #some default

    my $tenants_data = {
        tenants_dict         => undef,
        root_tenants         => undef,
        order_by             => -1,
        ordered_tenants_list => undef,
    };
    $tenants_data->{order_by} = $orderby;

    # read the data from the DB
    my $tenants_table = $self->{dbh}->resultset("Tenant")
      ->search( undef, { order_by => $tenants_data->{order_by} } );

    # build the tenants dict and list. tenants list is kept ordered
    $tenants_data->{ordered_tenants_list} = ();
    $tenants_data->{tenants_dict}         = {};
    while ( my $row = $tenants_table->next ) {
        push( @{ $tenants_data->{ordered_tenants_list} }, $row->id );
        $tenants_data->{tenants_dict}->{ $row->id } = {
            row      => $row,
            parent   => $row->parent_id,
            children => (),
        };
    }

    #build the root and children tenants lists, ordered by the orderby
    $tenants_data->{root_tenants} = ();
    foreach my $key ( @{ $tenants_data->{ordered_tenants_list} } ) {
        my $value  = $tenants_data->{tenants_dict}->{$key};
        my $parent = $value->{parent};
        if ( !defined($parent) ) {
            push @{ $tenants_data->{root_tenants} }, $key;
        }
        else {
            push @{ $tenants_data->{tenants_dict}->{$parent}{children} }, $key;
        }
    }

    return $tenants_data;
}

sub current_user_tenant {
    my $self = shift;
    return $self->{current_user_tenant};
}

sub get_tenant_by_id {
    my $self         = shift;
    my $tenants_data = shift;
    my $tenant_id    = shift;

    return $tenants_data->{tenants_dict}->{$tenant_id}{row};
}

sub get_tenants_list {
    my $self         = shift;
    my $tenants_data = shift;

    my @result = ();
    foreach my $tenant_id ( @{ $tenants_data->{ordered_tenants_list} } ) {
        push @result, $tenants_data->{tenants_dict}->{$tenant_id}{row};
    }

    return @result;
}

sub get_hierarchic_tenants_list {
    my $self         = shift;
    my $tenants_data = shift;
    my $tree_root    = shift;

#building an heirarchic list via standard DFS.
#First - adding to the stack the root nodes under which we want to get the tenats
    my @stack = ();
    if ( defined($tree_root) ) {
        if (exists($tenants_data->{tenants_dict}->{$tree_root})){
            push( @stack, $tree_root );
        }
    }
    else {
# root is not set, putting all roots, using "reverse" as we push it into a stack
# (from which we pop) and we want to keep the original order
        push( @stack, reverse( @{ $tenants_data->{root_tenants} } ) );
    }

#starting the actual DFS, poping from the stack, and pushing the poped node children
    my @result = ();
    while (@stack) {
        my $tenant_id = pop @stack;
        push( @result, $tenants_data->{tenants_dict}->{$tenant_id}{row} );

      # pushing the children in a reverse order, as we working with stack and we
      # pop from the end (but want to keep the overall order)
        push(
            @stack,
            reverse(
                @{ $tenants_data->{tenants_dict}->{$tenant_id}{children} }
            )
        );
    }

    return @result;
}

sub is_root_tenant {
    my $self         = shift;
    my $tenants_data = shift;
    my $tenant_id    = shift;

    if ( !defined($tenant_id) ) {
        return 0;
    }

    #root <==> parent is undef
    return !( defined( $tenants_data->{tenants_dict}{$tenant_id}{parent} ) );
}

sub cascade_delete_tenants_tree {
    #assuming all relevant tenants are not in use
    my $self = shift;
    my $tenants_data = shift;
    my $tree_root = shift;

    my @tenants = reverse $self->get_hierarchic_tenants_list($tenants_data, $tree_root);
    foreach my $tenant (@tenants) {
        $tenant->delete();
    }
}

sub is_tenant_resource_accessible {
    my $self             = shift;
    my $tenants_data     = shift;
    my $resource_tenancy = shift;

    return $self->_is_resource_accessable( $tenants_data, $resource_tenancy);
}

sub is_user_resource_accessible {
    my $self             = shift;
    my $tenants_data     = shift;
    my $resource_tenancy = shift;

    return $self->_is_resource_accessable( $tenants_data, $resource_tenancy);
}

sub is_ds_resource_accessible {
    my $self             = shift;
    my $tenants_data     = shift;
    my $resource_tenancy = shift;

    return $self->_is_resource_accessable( $tenants_data, $resource_tenancy);
}

sub is_ds_resource_accessible_to_tenant {
    my $self             = shift;
    my $tenants_data     = shift;
    my $resource_tenancy = shift;
    my $user_tenancy = shift;

    return $self->_is_resource_accessable_to_tenant( $tenants_data, $resource_tenancy, $user_tenancy);
}

sub use_tenancy {
    # tenancy is always enabled
    return 1;
}

sub get_tenant_heirarchy_depth {

    #return "undef" in case of error
    #a root tenant is of depth 0
    my $self         = shift;
    my $tenants_data = shift;
    my $tenant_id    = shift;

    if ( !defined( $tenants_data->{tenants_dict}{$tenant_id} ) ) {
        $self->_error(
            "Check tenancy depth - tenant $tenant_id does not exists");
        return undef;
    }

    my $iter_id = $tenant_id;

    my $depth = 0;
    while ( defined($iter_id) ) {
        $iter_id = $tenants_data->{tenants_dict}{$iter_id}{parent};
        $depth++;
        if ( $depth > $self->max_heirarchy_limit() ) {
            $self->_error(
"Check tenancy depth for tenant $tenant_id - reached heirarchy limit"
            );
            return undef;
        }
    }

    return $depth - 1;
}

sub get_tenant_heirarchy_height {

    #return "undef" in case of error
    #a leaf tenant is of height 0
    my $self         = shift;
    my $tenants_data = shift;
    my $tenant_id    = shift;

    if ( !defined( $tenants_data->{tenants_dict}{$tenant_id} ) ) {
        $self->_error(
            "Check tenancy height - tenant $tenant_id does not exists");
        return undef;
    }

    #calc tenant height
    my @tenants_list = reverse(
        $self->get_hierarchic_tenants_list( $tenants_data, $tenant_id ) );
    my %tenants_height = {};

    foreach my $tenant_row (@tenants_list) {
        my $tid = $tenant_row->id;
        $tenants_height{$tid} = 0;
    }

    foreach my $tenant_row (@tenants_list) {
        my $tid    = $tenant_row->id;
        my $par_id = $tenant_row->parent_id;
        if ( ( $tenants_height{$par_id} ) < ( $tenants_height{$tid} + 1 ) ) {
            $tenants_height{$par_id} = $tenants_height{$tid} + 1;
        }
    }

    return $tenants_height{$tenant_id};
}

sub is_anchestor_of {

    #return "undef" in case of error
    my $self          = shift;
    my $tenants_data  = shift;
    my $anchestor_id  = shift;
    my $descendant_id = shift;

    if ( !defined($anchestor_id) ) {
        $self->_error("Check tenants relations - got undef anchestor");
        return undef;
    }

    if ( !defined( $tenants_data->{tenants_dict}{$anchestor_id} ) ) {
        $self->_error(
            "Check tenants relations - tenant $anchestor_id does not exists");
        return undef;
    }

    if ( !defined($descendant_id) ) {
        $self->_error("Check tenants relations - got undef descendant");
        return undef;
    }

    if ( !defined( $tenants_data->{tenants_dict}{$descendant_id} ) ) {
        $self->_error(
            "Check tenants relations - tenant $descendant_id does not exists");
        return undef;
    }

    my $iter_id = $descendant_id;

    my $descendant_depth = 0;
    while ( defined($iter_id) ) {
        if ( $anchestor_id == $iter_id ) {
            return 1;
        }
        $iter_id = $tenants_data->{tenants_dict}{$iter_id}{parent};
        $descendant_depth++;
        if ( $descendant_depth > $self->max_heirarchy_limit() )
        {    #recursion limit
            $self->_error(
"Tenants relation failed for tenants $anchestor_id / $descendant_id - reached heirarchy limit"
            );
            return undef;
        }
    }

    return 0;
}

sub max_heirarchy_limit {
    my $self = shift;
    return 100;
}

##############################################################

sub _error {
    my $self    = shift;
    my $message = shift;

    $context = $self->{context};
    if ( defined($context) ) {
        $context->app->log->error($message);
    }
    else {
        print "Error: ", $message, "\n";
    }
}

sub _is_resource_accessable {
    my $self = shift;
    my $tenants_data = shift;
    my $resource_tenant = shift;
    my $user_tenant = $self->current_user_tenant();
    return $self->_is_resource_accessable_to_tenant($tenants_data, $resource_tenant, $user_tenant)
}

sub _is_resource_accessable_to_tenant {
    my $self            = shift;
    my $tenants_data    = shift;
    my $resource_tenant = shift;
    my $user_tenant     = shift;

    if (!$self->{use_tenancy}) {
        #mechanisem disabled
        return 1;
    }


    if ( defined($user_tenant) ) {
        my $tenant_record    = $tenants_data->{tenants_dict}->{$user_tenant};
        my $is_active_tenant = $tenant_record->{row}->active;
        if ( !$is_active_tenant ) {

            #user tenant is in-active - cannot do any operation
            return 0;
        }
    }

    if ( !defined($resource_tenant) ) {

        #the object has no tenancy - opened for all
        return 1;
    }

    if ( !defined($user_tenant) ) {

        #the user has no tenancy, - cannot approach items with tenancy
        return 0;
    }

    if ( $user_tenant == $resource_tenant ) {

        #resource has same tenancy of the user, operations are allowed
        return 1;
    }

    #checking if the user tenant is an ancestor of the resource tenant
    my $is_user_tenat_parent_of_resource =
      $self->is_anchestor_of( $tenants_data, $user_tenant, $resource_tenant );
    if ( !defined($is_user_tenat_parent_of_resource) ) {

        #error - give access only to root tenant (so it can fix the problem)
        return $self->is_root_tenant( $tenants_data, $user_tenant );
    }
    return $is_user_tenat_parent_of_resource;
}

1;
