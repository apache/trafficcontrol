package API::FederationUser;
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

sub index {
    my $self    = shift;
    my $fed_id  = $self->param('fedId');

    my @data;
    my $rs_data = $self->db->resultset("FederationTmuser")->search( { 'federation' => $fed_id }, { prefetch => [ 'tm_user' ] } );
    while ( my $row = $rs_data->next ) {
        push(
            @data, {
                "company"   => $row->tm_user->company,
                "email"     => $row->tm_user->email,
                "fullName"  => $row->tm_user->full_name,
                "id"        => $row->tm_user->id,
                "role"      => $row->tm_user->role->name,
                "username"  => $row->tm_user->username,
            }
        );
    }
    $self->success( \@data );
}

sub assign_users_to_federation {
    my $self        = shift;
    my $fed_id      = $self->param('fedId');
    my $params      = $self->req->json;
    my $user_ids    = $params->{userIds};
    my $replace     = $params->{replace};
    my $count       = 0;

    if ( !&is_admin($self) ) {
        return $self->forbidden();
    }

    my $fed = $self->db->resultset('Federation')->find( { id => $fed_id } );
    if ( !defined($fed) ) {
        return $self->not_found();
    }

    if ( ref($user_ids) ne 'ARRAY' ) {
        return $self->alert("User IDs must be an array");
    }

    if ( $replace ) {
        # start fresh and delete existing fed/user associations
        my $delete = $self->db->resultset('FederationTmuser')->search( { federation => $fed_id } );
        $delete->delete();
    }

    my @values = ( [ qw( federation tm_user ) ]); # column names are required for 'populate' function

    foreach my $user_id (@{ $user_ids }) {
        push(@values, [ $fed_id, $user_id ]);
        $count++;
    }

    $self->db->resultset("FederationTmuser")->populate(\@values);

    my $msg = $count . " user(s) were assigned to the " . $fed->cname . " federation";
    &log( $self, $msg, "APICHANGE" );

    my $response = $params;
    return $self->success($response, $msg);
}

sub delete {
    my $self        = shift;
    my $fed_id      = $self->param('fedId');
    my $user_id     = $self->param('userId');

    if ( !&is_admin($self) ) {
        return $self->forbidden();
    }

    my $fed_user = $self->db->resultset("FederationTmuser")->search( { 'federation.id' => $fed_id, 'tm_user' => $user_id }, { prefetch => [ 'federation', 'tm_user' ] } );
    if ( !defined($fed_user) ) {
        return $self->not_found();
    }

    my $row = $fed_user->next;
    my $rs = $fed_user->delete();
    if ($rs) {
        my $msg = "Removed user [ " . $row->tm_user->username . " ] from federation [ " . $row->federation->cname . " ]";
        &log( $self, $msg, "APICHANGE" );
        return $self->success_message($msg);
    }

    return $self->alert( "Failed to remove user from federation." );
}

1;
