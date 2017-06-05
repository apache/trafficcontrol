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
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Utils::Helper;

sub delete {
    my $self     	= shift;
    my $ds_id  	 	= $self->param('dsId');
    my $user_id	    = $self->param('userId');

    if ( !&is_oper($self) ) {
        return $self->forbidden();
    }

    my $ds_user = $self->db->resultset('DeliveryserviceTmuser')->search( { deliveryservice => $ds_id, tm_user_id => $user_id }, { prefetch => [ 'deliveryservice', 'tm_user' ] } );
    if ( $ds_user->count == 0 ) {
        return $self->not_found();
    }

    my $row = $ds_user->next;
    my $rs = $ds_user->delete();
    if ($rs) {
        my $msg = "User [ " . $row->tm_user->username . " ] unlinked from deliveryservice [ " . $row->deliveryservice->id . " | " . $row->deliveryservice->xml_id . " ].";
        &log( $self, $msg, "APICHANGE" );
        return $self->success_message($msg);
    }

    return $self->alert( "Failed to unlink user from delivery service." );
}

1;
