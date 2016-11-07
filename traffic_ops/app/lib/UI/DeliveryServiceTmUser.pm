package UI::DeliveryServiceTmUser;
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

# Create
sub create {
	my $self    = shift;
	my $ds      = $self->param('deliveryservice');
	my $tm_user_id = $self->param('tm_user_id');

	my $new_id = -1;
	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $ds_name = $self->db->resultset('Deliveryservice')->search( { id => $ds } )->get_column('xml_id')->single();
		my $insert = $self->db->resultset('DeliveryserviceTmuser')->create( { deliveryservice => $ds, tm_user_id => $tm_user_id } );

		$insert->insert();
		$new_id = $insert->id;
		&log( $self, "Create ds_tm_user " . $ds_name . " <-> " . $tm_user_id, "UICHANGE" );

	}
	my $referer = $self->req->headers->header('referer');
	return $self->redirect_to($referer);
}

# Delete
sub delete {
	my $self    = shift;
	my $ds      = $self->param('ds');
	my $tm_user_id = $self->param('tm_user_id');

	if ( !&is_oper($self) ) {
		$self->flash( alertmsg => "No can do. Get more privs." );
	}
	else {
		my $ds_name = $self->db->resultset('Deliveryservice')->search( { id => $ds } )->get_column('xml_id')->single();
		my $delete = $self->db->resultset('DeliveryserviceTmuser')->search( { deliveryservice => $ds, tm_user_id => $tm_user_id } );
		$delete->delete();
		&log( $self, "Delete delivery_service_tmuser " . $ds_name . " <-> " . $tm_user_id, "UICHANGE" );
	}
	my $referer = $self->req->headers->header('referer');
	return $self->redirect_to($referer);
}

1;
