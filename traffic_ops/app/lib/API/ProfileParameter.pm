package API::ProfileParameter;
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

# Read
sub index {
	my $self = shift;
	my @data;
	my $orderby = "profile";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );
	my $rs_data = $self->db->resultset("ProfileParameter")->search( undef, { prefetch => [ 'profile' ], order_by => $orderby } );
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
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	if ( !defined($params) ) {
		return $self->alert("parameters must be in JSON format.");
	}
	if ( ref($params) ne 'ARRAY' ) {
		my @temparry;
		push(@temparry, $params);
		$params = \@temparry;
	}
	if ( scalar(@{ $params }) == 0 ) {
		return $self->alert("parameters array length is 0.");
	}

	$self->db->txn_begin();
	foreach my $param (@{ $params }) {
		my $profile_exist = $self->db->resultset('Profile')->find( { id => $param->{profileId} } );
		if ( !defined($profile_exist) ) {
			$self->db->txn_rollback();
			return $self->alert("profile with id: " . $param->{profileId} . " isn't existed.");
		}
		my $param_exist = $self->db->resultset('Parameter')->find( { id => $param->{parameterId} } );
		if ( !defined($param_exist) ) {
			$self->db->txn_rollback();
			return $self->alert("parameter with id: " . $param->{parameterId} . " isn't existed.");
		}
		my $profile_param_exist = $self->db->resultset('ProfileParameter')->find( { parameter => $param_exist->id, profile => $profile_exist->id } );
		if ( defined($profile_param_exist) ) {
			$self->db->txn_rollback();
			return $self->alert("parameter: " . $param->{parameterId} . " already associated with profile: " . $param->{profileId});
		}
		$self->db->resultset('ProfileParameter')->create( { parameter => $param_exist->id, profile => $profile_exist->id } )->insert();
	}
	$self->db->txn_commit();

	&log( $self, "New profile parameter associations were created.", "APICHANGE" );

	my $response = $params;
	return $self->success($response, "Profile parameter associations were created.");
}

sub assign_params_to_profile {
	my $self 		= shift;
	my $params 		= $self->req->json;
	my $profile_id	= $params->{profileId};
	my $param_ids	= $params->{paramIds};
	my $replace 	= $params->{replace};
	my $count		= 0;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $profile = $self->db->resultset('Profile')->find( { id => $profile_id } );
	if ( !defined($profile) ) {
		return $self->not_found();
	}

	if ( ref($param_ids) ne 'ARRAY' ) {
		return $self->alert("Parameters must be an array");
	}

	if ( $replace ) {
		# start fresh and delete existing profile/parameter associations
		my $delete = $self->db->resultset('ProfileParameter')->search( { profile => $profile_id } );
		$delete->delete();
	}

	my @values = ( [ qw( profile parameter ) ]); # column names are required for 'populate' function

	foreach my $param_id (@{ $param_ids }) {
		push(@values, [ $profile_id, $param_id ]);
		$count++;
	}

	$self->db->resultset("ProfileParameter")->populate(\@values);

	my $msg = $count . " parameters were assigned to the " . $profile->name . " profile";
	&log( $self, $msg, "APICHANGE" );

	my $response = $params;
	return $self->success($response, $msg);
}

sub assign_profiles_to_param {
	my $self 		= shift;
	my $params 		= $self->req->json;
	my $param_id	= $params->{paramId};
	my $profile_ids	= $params->{profileIds};
	my $replace 	= $params->{replace};
	my $count		= 0;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $parameter = $self->db->resultset('Parameter')->find( { id => $param_id } );
	if ( !defined($parameter) ) {
		return $self->not_found();
	}

	if ( ref($profile_ids) ne 'ARRAY' ) {
		return $self->alert("Profiles must be an array");
	}

	if ( $replace ) {
		# start fresh and delete existing parameter/profile associations
		my $delete = $self->db->resultset('ProfileParameter')->search( { parameter => $param_id } );
		$delete->delete();
	}

	my @values = ( [ qw( profile parameter ) ]); # column names are required for 'populate' function

	foreach my $profile_id (@{ $profile_ids }) {
		push(@values, [ $profile_id, $param_id ]);
		$count++;
	}

	$self->db->resultset("ProfileParameter")->populate(\@values);

	my $msg = $count . " profiles were assigned to the " . $parameter->name . " parameter";
	&log( $self, $msg, "APICHANGE" );

	my $response = $params;
	return $self->success($response, $msg);
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

sub create_param_for_profile_name {
	my $self = shift;
	my $json = $self->req->json;
	my $profileName = $self->param('name');
	if ( !&is_oper($self) ) {
		return $self->forbidden("You must be an admin or oper to perform this operation!");
	}

    my $profile_find = $self->db->resultset('Profile')->find({ name => $profileName });
    if ( !defined($profile_find) ){
        return $self->not_found("profile ". $profileName. " does not exist.");
    }

    #return &addex($self, $profileName, $json, $profile_find); 
    return $self->addex($profileName, $json, $profile_find); 
}

sub create_param_for_profile_id {
	my $self = shift;
	my $json = $self->req->json;
	my $profileId = $self->param('id');
	if ( !&is_oper($self) ) {
		return $self->forbidden("You must be an admin or oper to perform this operation!");
	}

    my $profile_find = $self->db->resultset('Profile')->find({ id => $profileId });
    if ( !defined($profile_find) ){
        return $self->not_found("profile with id ". $profileId. " does not exist.");
    }

    return $self->addex($profile_find->name, $json, $profile_find); 
}

sub addex {
	my $self = shift;
    my $profileName = shift;
    my $json = shift;
    my $profile_find = shift;

	if ( ref($json) ne 'ARRAY' ) {
		my @temparry;
		push(@temparry, $json);
		$json = \@temparry;
	}
	if ( scalar(@{ $json }) == 0 ) {
		return $self->alert("parameters array length is 0.");
	}

    my @new_parameters = ();
    $self->db->txn_begin();
    foreach my $param (@{ $json }) {
        if ( !defined($param->{name}) ) {
            $self->db->txn_rollback();
            return $self->alert("there is parameter name does not provide , configFile:".$param->{configFile}." , value:".$param->{value});
        }
        if ( !defined($param->{configFile}) ) {
            $self->db->txn_rollback();
            return $self->alert("there is parameter configFile does not provide , name:".$param->{name}." , value:".$param->{value});
        }
        if ( !defined($param->{value}) ) {
            $self->db->txn_rollback();
            return $self->alert("there is parameter value does not provide , name:".$param->{name}." , configFile:".$param->{configFile});
        }
        if ( !defined($param->{secure}) ) {
            $param->{secure} = 0
        } else {
            if (($param->{secure} ne '0') && ($param->{secure} ne '1')) {
                $self->db->txn_rollback();
                return $self->alert("secure must 0 or 1, parameter [name:".$param->{name}." , configFile:".$param->{configFile}." , value:".$param->{value}." , secure:".$param->{secure}."]");
            }
            $param->{secure} = 0 if ($param->{secure} eq '0' );
            $param->{secure} = 1 if ($param->{secure} eq '1' );
            if ( $param->{secure} != 0 && !&is_admin($self)) {
                $self->db->txn_rollback();
                return $self->forbidden("Parameter[name:".$param->{name}." , configFile:".$param->{configFile}." , value:".$param->{value}."] secure=1, You must be an admin to perform this operation!");
            }
        }

        my $param_find = $self->db->resultset('Parameter')->find({ 
                name => $param->{name},
                config_file => $param->{configFile},
                value => $param->{value},
            } ) ;
        my $param_id = undef;
        if ( !defined($param_find) ){
            my $insert = $self->db->resultset('Parameter')->create(
                {
                    name            => $param->{name},
                    config_file     => $param->{configFile},
                    value           => $param->{value},
                    secure          => $param->{secure}
                }
            );
            $insert->insert();
            $param_id = $insert->id;
            push(@new_parameters, {
                    'id'            => $insert->id,
                    'name'          => $insert->name,
                    'configFile'    => $insert->config_file,
                    'value'         => $insert->value,
                    'secure'        => $insert->secure
                })
        } else {
            $param_id = $param_find->id;
            push(@new_parameters, {
                    'id'            => $param_find->id,
                    'name'          => $param_find->name,
                    'configFile'    => $param_find->config_file,
                    'value'         => $param_find->value,
                    'secure'        => $param_find->secure
                })
        }

        my $profile_parameter_find = $self->db->resultset("ProfileParameter")->find( {
                profile     => $profile_find->id,
                parameter   => $param_id
            } );
        if ( !defined($profile_parameter_find) ) {
            my $insert = $self->db->resultset('ProfileParameter')->create(
                {
                    profile     => $profile_find->id,
                    parameter   => $param_id
                }
            );
            $insert->insert();
        } else {
            $self->app->log->warn("parameter [name:".$param_find->name." , configFile:".$param_find->config_file." , value:".$param_find->value."] has already assigned to profile ". $profileName);
        }
    }
    $self->db->txn_commit();

    my $response;
    $response->{profileName} = $profile_find->name;
    $response->{profileId}   = $profile_find->id;
    $response->{parameters}  = \@new_parameters;
    $self->success($response, "Assign parameters successfully to profile ". $response->{profileName});
}

1;
