package API::ProfileParameter;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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

# Read
sub index {
	my $self = shift;
	my @data;
	my $orderby = "profile";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("ProfileParameter")->search( undef, { order_by => $orderby } );
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"profile"     => $row->profile->name,
				"parameter"   => $row->parameter->id,
				"lastUpdated" => $row->last_updated,
			}
		);
	}
	$self->success( \@data );
}

sub create {
	my $self = shift;
	my $profile_id = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $profile = $self->db->resultset('Profile')->find( { id => $profile_id } );
	if ( !defined($profile) ) {
		return $self->not_found();
	}

	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format.");
	}
	if ( !defined($params->{parametersId}) ) {
		return $self->alert("parameter parametersId is must.");
	}
	if ( ref($params->{parametersId}) ne 'ARRAY' ) {
		return $self->alert("parameter parametersId must be array.");
	}
	if ( scalar(@{ $params->{parametersId} }) == 0 ) {
		return $self->alert("parametersId array length is 0.");
	}

	my @param_ids = ();
	foreach my $param (@{ $params->{parametersId} }) {
		my $param_exist = $self->db->resultset('Parameter')->find( { id => $param } );
		if ( !defined($param_exist) ) {
			return $self->alert("parameter with id: $param isn't existed.");
		}
		my $profile_param_exist = $self->db->resultset('ProfileParameter')->find( { parameter => $param_exist->id, profile => $profile->id } );
		if ( defined($profile_param_exist) ) {
			return $self->alert("parameter: $param already associated with profile: $profile_id");
		}
		push(@param_ids, $param_exist->id);
	}

	foreach my $param_id (@param_ids) {
		$self->db->resultset('ProfileParameter')->create( { parameter => $param_id, profile => $profile->id } )->insert();
	}

	&log( $self, "Associate new parameters to profile: $profile_id", "APICHANGE" );

	my @new_params = ();
	my $rs = $self->db->resultset('ProfileParameter')->search( { profile => $profile->id } );
	while ( my $row = $rs->next ) {
		push(@new_params, $row->parameter->id)
	}

	my $response;
	$response->{id} = $profile_id;
	$response->{parametersId} = \@new_params;
	return $self->success($response, "Parameters were associated to profile: " . $profile_id);
}

sub delete {
	my $self = shift;
	my $profile_id = $self->param('profile_id');
	my $parameter_id = $self->param('parameter_id');

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $profile = $self->db->resultset('Profile')->find( { id => $profile_id } );
	if ( !defined($profile) ) {
		return $self->not_found();
	}

	my $parameter = $self->db->resultset('Parameter')->find( { id => $parameter_id } );
	if ( !defined($parameter) ) {
		return $self->not_found();
	}

	my $delete = $self->db->resultset('ProfileParameter')->find( { parameter => $parameter->id, profile => $profile->id } );
	if ( !defined($delete) ) {
		return $self->alert("parameter: $parameter_id isn't associated with profile: $profile_id.");
	}

	$delete->delete();

	&log( $self, "Delete profile parameter " . $profile->name . " <-> " . $parameter->name, "APICHANGE" );

	return $self->success_message("Profile parameter association was deleted.");
}

1;
