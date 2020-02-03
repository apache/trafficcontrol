package API::DeliveryServiceUser;
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
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Utils::Helper;

sub delete {
    my $self     	= shift;
    my $ds_id  	 	= $self->param('dsId');
    my $user_id	    = $self->param('userId');

    if ( !&is_oper($self) ) {
        return $self->with_deprecation_with_no_alternative("Forbidden", "error", 403);
    }

    my $user = $self->db->resultset('TmUser')->find( { id => $user_id } );
    if ( !defined($user) ) {
        return $self->with_deprecation_with_no_alternative("Resource not found.", "error", 404);
    }
    my $tenant_utils = Utils::Tenant->new($self);
    my $tenants_data = $tenant_utils->create_tenants_data_from_db();
    if (!$tenant_utils->is_user_resource_accessible($tenants_data, $user->tenant_id)) {
        #no access to resource tenant
        return $self->with_deprecation_with_no_alternative("Forbidden. User tenant is not available to the working user.", "error", 403);
    }

    my $ds = $self->db->resultset('Deliveryservice')->find( { id => $ds_id } );
    if ( !defined($ds) ) {
        return $self->with_deprecation_with_no_alternative("Resource not found.", "error", 404);
    }
    if (!$tenant_utils->is_ds_resource_accessible($tenants_data, $ds->tenant_id)) {
        return $self->with_deprecation_with_no_alternative("Forbidden. Delivery-service tenant is not available to the user.", "error", 403);
    }

    #not checking DS tenancy on deletion - we manage the user here - we remove permissions to touch a DS

    my $ds_user = $self->db->resultset('DeliveryserviceTmuser')->search( { deliveryservice => $ds_id, tm_user_id => $user_id }, { prefetch => [ 'deliveryservice', 'tm_user' ] } );
    if ( $ds_user->count == 0 ) {
        return $self->with_deprecation_with_no_alternative("Resource not found.", "error", 404);
    }

    my $row = $ds_user->next;
    my $rs = $ds_user->delete();
    if ($rs) {
        my $msg = "User [ " . $row->tm_user->username . " ] unlinked from deliveryservice [ " . $row->deliveryservice->id . " | " . $row->deliveryservice->xml_id . " ].";
        &log( $self, $msg, "APICHANGE" );
        return $self->with_deprecation_with_no_alternative($msg, "success", 200);
    }
    return $self->with_deprecation_with_no_alternative("Failed to unlink user from delivery service.", "error", 400);
}

1;
