package API::User;
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
use Digest::SHA1 qw(sha1_hex);
use Mojolicious::Validator;
use Mojolicious::Validator::Validation;
use Data::Dumper;
use Test::More;
use Email::Valid;
use Utils::Helper::ResponseHelper;
use Validate::Tiny ':all';
use UI::ConfigFiles;
use UI::Tools;

sub login {
	my $self     = shift;
	my $options  = shift;

	my $u     = $self->req->json->{u};
	my $p     = $self->req->json->{p};
	my $token = $options->{'token'};

	my $result = $self->authenticate( $u, $p, $options );
	if ($result) {
		return $self->success_message("Successfully logged in.");
	}
	elsif ( defined($token) ) {
		return $self->invalid_token;
	}
	else {
		return $self->invalid_username_or_password;
	}
}

sub token_login {
	my $self = shift;

	my $token = $self->req->json->{t};
	return $self->login( { token => $token } );
}

# Read
sub index {
	my $self = shift;
	my @data;
	my $username = $self->param('username');

	my $orderby = "username";
	$orderby = $self->param('orderby') if ( defined $self->param('orderby') );

	my $dbh;
	if ( defined $username ) {
		$dbh = $self->db->resultset("TmUser")->search( { username => $username }, { prefetch => [ 'role', 'tenant' ], order_by => 'me.' . $orderby } );
	}
	else {
		$dbh = $self->db->resultset("TmUser")->search( undef, { prefetch => [ 'role', 'tenant' ], order_by => 'me.' . $orderby } );
	}

	while ( my $row = $dbh->next ) {
		push(
			@data, {
				"addressLine1"    => $row->address_line1,
				"addressLine2"    => $row->address_line2,
				"city"            => $row->city,
				"company"         => $row->company,
				"country"         => $row->country,
				"email"           => $row->email,
				"fullName"        => $row->full_name,
				"gid"             => $row->gid,
				"id"              => $row->id,
				"lastUpdated"     => $row->last_updated,
				"newUser"         => \$row->new_user,
				"phoneNumber"     => $row->phone_number,
				"postalCode"      => $row->postal_code,
				"publicSshKey"    => $row->public_ssh_key,
				"registrationSent"=> \$row->registration_sent,
				"role"            => $row->role->id,
				"rolename"        => $row->role->name,
				"stateOrProvince" => $row->state_or_province,
				"uid"             => $row->uid,
				"username"        => $row->username,
				"tenant"          => defined ($row->tenant) ? $row->tenant->name : undef,
				"tenantId"        => $row->tenant_id
			}
		);
	}
	$self->success( \@data );
}

sub show {
	my $self = shift;
	my $id   = $self->param('id');

	my $rs_data = $self->db->resultset("TmUser")->search( { 'me.id' => $id }, { prefetch => [ 'role' , 'tenant'] } );
	my @data = ();
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"addressLine1"    => $row->address_line1,
				"addressLine2"    => $row->address_line2,
				"city"            => $row->city,
				"company"         => $row->company,
				"country"         => $row->country,
				"email"           => $row->email,
				"fullName"        => $row->full_name,
				"gid"             => $row->gid,
				"id"              => $row->id,
				"lastUpdated"     => $row->last_updated,
				"newUser"         => \$row->new_user,
				"phoneNumber"     => $row->phone_number,
				"postalCode"      => $row->postal_code,
				"publicSshKey"    => $row->public_ssh_key,
				"registrationSent"=> $row->registration_sent,
				"role"            => $row->role->id,
				"rolename"        => $row->role->name,
				"stateOrProvince" => $row->state_or_province,
				"uid"             => $row->uid,
				"username"        => $row->username,
				"tenant"          => defined ($row->tenant) ? $row->tenant->name : undef,
				"tenantId"        => $row->tenant_id
			}
		);
	}
	$self->success( \@data );
}

sub update {
	my $self   = shift;
	my $id     = $self->param('id');
	my $params = $self->req->json;

	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my ( $is_valid, $result ) = $self->is_valid($params);

	if ( !$is_valid ) {
		return $self->alert($result);
	}

	my $user = $self->db->resultset('TmUser')->find( { id => $id } );
	if ( !defined($user) ) {
		return $self->not_found();
	}

	#setting tenant_id to undef if tenant is not set. 
 	my $tenant_id = exists($params->{tenantId}) ? $params->{tenantId} :  undef; 
 	
	my $values = {
		address_line1 			=> $params->{addressLine1},
		address_line2 			=> $params->{addressLine2},
		city 					=> $params->{city},
		company 				=> $params->{company},
		country 				=> $params->{country},
		email 					=> $params->{email},
		full_name 				=> $params->{fullName},
		new_user 				=> ( $params->{newUser} ) ? 1 : 0,
		phone_number 			=> $params->{phoneNumber},
		postal_code 			=> $params->{postalCode},
		public_ssh_key 			=> $params->{publicSshKey},
		role 					=> $params->{role},
		state_or_province 		=> $params->{stateOrProvince},
		username 				=> $params->{username},
		tenant_id 				=> $tenant_id
		
	};

	if ( defined($params->{localPasswd}) && $params->{localPasswd} ne '' ) {
		$values->{"local_passwd"} = sha1_hex($params->{localPasswd});
	}
	if ( defined($params->{confirmLocalPasswd}) && $params->{confirmLocalPasswd} ne '' ) {
		$values->{"confirm_local_passwd"} = sha1_hex($params->{confirmLocalPasswd});
	}

	my $rs = $user->update($values);
	if ($rs) {
		my $response;
		$response->{addressLine1}        	= $rs->address_line1;
		$response->{addressLine2} 			= $rs->address_line2;
		$response->{city} 					= $rs->city;
		$response->{company} 				= $rs->company;
		$response->{country} 				= $rs->country;
		$response->{email} 					= $rs->email;
		$response->{fullName} 				= $rs->full_name;
		$response->{gid}          			= $rs->gid;
		$response->{id}          			= $rs->id;
		$response->{lastUpdated} 			= $rs->last_updated;
		$response->{newUser} 				= \$rs->new_user;
		$response->{phoneNumber} 			= $rs->phone_number;
		$response->{postalCode} 			= $rs->postal_code;
		$response->{publicSshKey} 			= $rs->public_ssh_key;
		$response->{registrationSent} 		= \$rs->registration_sent;
		$response->{role} 					= $rs->role->id;
		$response->{roleName} 				= $rs->role->name;
		$response->{stateOrProvince} 		= $rs->state_or_province;
		$response->{uid} 					= $rs->uid;
		$response->{username} 				= $rs->username;
		$response->{tenantId} 				= $rs->tenant_id;


		&log( $self, "Updated User with username '" . $rs->username . "' for id: " . $rs->id, "APICHANGE" );

		return $self->success( $response, "User update was successful." );
	}
	else {
		return $self->alert("User update failed.");
	}

}

# Create
sub create {
	my $self = shift;
	my $params = $self->req->json;
	
	if ( !&is_oper($self) ) {
		return $self->forbidden();
	}

	my $name = $params->{username};
	if ( !defined($name) ) {
		return $self->alert("Username is required.");
	}
	
	my $existing = $self->db->resultset('TmUser')->search( { username => $name } )->single();
	if ($existing) {
		return $self->alert("A user with username \"$name\" already exists.");
	}


	if ( !defined($params->{localPassword}) ) {
		return $self->alert("local-password is required.");
	}

	if ( !defined($params->{confirmLocalPassword}) ) {
		return $self->alert("confirm-local-password is required.");
	}

	if ( !defined($params->{role}) ) {
		return $self->alert("role is required.");
	}
	
	#setting tenant_id to the user's tenant if tenant is not set. TODO(nirs): remove when tenancy is no longer optional in the API
	my $tenant_id = exists($params->{tenantId}) ? $params->{tenantId} :  $self->current_user_tenant();

	my $values = {
		address_line1 			=> $params->{addressLine1},
		address_line2 			=> $params->{addressLine2},
		city 				=> $params->{city},
		company 			=> $params->{company},
		country 			=> $params->{country},
		email 				=> $params->{email},
		full_name 			=> $params->{fullName},
		new_user            		=> 0,		
		phone_number 			=> $params->{phoneNumber},
		postal_code 			=> $params->{postalCode},
		public_ssh_key 			=> $params->{publicSshKey},
		registration_sent 		=> undef,
		role 				=> $params->{role},
		state_or_province 		=> $params->{stateOrProvince},
		username 			=> $params->{username},
		uid                  		=> 0,		
		gid                  		=> 0,
		local_passwd         		=> sha1_hex($params->{localPassword} ),
		confirm_local_passwd 		=> sha1_hex($params->{confirmLocalPassword} ),
		tenant_id			=> $tenant_id,		

	};
	
	my ( $is_valid, $result ) = $self->is_valid($values);

	if ( !$is_valid ) {
		return $self->alert($result);
	}
	
	my $insert = $self->db->resultset('TmUser')->create($values);
	my $rs = $insert->insert();

	if ($rs) {
		my $response;
		$response->{addressLine1}        	= $rs->address_line1;
		$response->{addressLine2} 		= $rs->address_line2;
		$response->{city} 			= $rs->city;
		$response->{company} 			= $rs->company;
		$response->{country} 			= $rs->country;
		$response->{email} 			= $rs->email;
		$response->{fullName} 			= $rs->full_name;
		$response->{gid}          		= $rs->gid;
		$response->{id}          		= $rs->id;
		$response->{lastUpdated} 		= $rs->last_updated;
		$response->{newUser} 			= \$rs->new_user;
		$response->{phoneNumber} 		= $rs->phone_number;
		$response->{postalCode} 		= $rs->postal_code;
		$response->{publicSshKey} 		= $rs->public_ssh_key;
		$response->{registrationSent} 		= \$rs->registration_sent;
		$response->{role} 			= $rs->role->id;
		$response->{roleName} 			= $rs->role->name;
		$response->{stateOrProvince} 		= $rs->state_or_province;
		$response->{uid} 			= $rs->uid;
		$response->{username} 			= $rs->username;
		$response->{tenantId} 			= $rs->tenant_id;

		&log( $self, "Adding User with username '" . $rs->username . "' for id: " . $rs->id, "APICHANGE" );

		return $self->success( $response, "User creation was successful." );
	}
	else {
		return $self->alert("User creation failed.");
	}
}


# Reset the User Profile password
sub reset_password {
	my $self     = shift;
	my $email_to = $self->req->json->{email};
	my $dbh      = $self->db->resultset('TmUser')->find( { email => $email_to } );
	if ( defined($dbh) ) {

		my $email_notice = 'Successfully sent reset password to: ' . $email_to;
		$self->app->log->info($email_notice);

		my $token = Data::GUID->new;
		if ( $self->send_password_reset_email( $email_to, $token ) ) {
			$self->update_user_token( $email_to, $token );
		}

		return $self->success_message( "Successfully sent password reset to email '" . $email_to . "'" );
	}
	else {
		return $self->alert( { "Email not found " => "'" . $email_to . "'" } );
	}

}

sub get_available_deliveryservices {
	my $self = shift;
	my @data;
	my $id = $self->param('id');
	my %dsids;
	my %takendsids;

	my $rs_takendsids = undef;
	$rs_takendsids = $self->db->resultset("DeliveryserviceTmuser")->search( { 'tm_user_id' => $id } );

	while ( my $row = $rs_takendsids->next ) {
		$takendsids{ $row->deliveryservice->id } = undef;
	}

	my $rs_links = $self->db->resultset("Deliveryservice")->search( undef, { order_by => "xml_id" } );
	while ( my $row = $rs_links->next ) {
		if ( !exists( $takendsids{ $row->id } ) ) {
			push( @data, { "id" => $row->id, "xmlId" => $row->xml_id } );
		}
	}

	$self->success( \@data );
}

# Read the current user profile and produce the result
sub current {
	my $self = shift;
	my @data;
	my $current_username = $self->current_user()->{username};

	if ( &is_ldap($self) ) {
		my $role = $self->db->resultset('Role')->search( { name => "read-only" } )->get_column('id')->single;

		push(
			@data, {
				"id"              => "0",
				"username"        => $current_username,
				"tenantId"	  => undef,
				"tenant"          => undef,
				"publicSshKey"  => "",
				"role"            => $role,
				"uid"             => "0",
				"gid"             => "0",
				"company"         => "",
				"email"           => "",
				"fullName"        => "",
				"newUser"         => \0,
				"localUser"       => \0,
				"addressLine1"    => "",
				"addressLine2"    => "",
				"city"            => "",
				"stateOrProvince" => "",
				"phoneNumber"     => "",
				"postalCode"      => "",
				"country"         => "",
			}
		);

		return $self->success( @data );
	}
	else {
		my $dbh = $self->db->resultset('TmUser')->search( { username => $current_username } , { prefetch => [ 'role' , 'tenant' ] } );
		while ( my $row = $dbh->next ) {
			push(
				@data, {
					"id"              => $row->id,
					"username"        => $row->username,
					"publicSshKey"    => $row->public_ssh_key,
					"role"            => $row->role->id,
					"uid"             => $row->uid,
					"gid"             => $row->gid,
					"company"         => $row->company,
					"email"           => $row->email,
					"fullName"        => $row->full_name,
					"newUser"         => \$row->new_user,
					"localUser"       => \1,
					"addressLine1"    => $row->address_line1,
					"addressLine2"    => $row->address_line2,
					"city"            => $row->city,
					"stateOrProvince" => $row->state_or_province,
					"phoneNumber"     => $row->phone_number,
					"postalCode"      => $row->postal_code,
					"country"         => $row->country,
					"tenant"          => defined ($row->tenant) ? $row->tenant->name : undef,
					"tenantId"        => $row->tenant_id,
				}
			);
		}
		return $self->success(@data);
	}
}

# Update
sub update_current {
	my $self = shift;

	my $user = $self->req->json->{user};
	if ( &is_ldap($self) ) {
		return $self->alert("Profile cannot be updated because '" . $user->{username} ."' is logged in as LDAP.");
	}

	my $db_user;

	# Prevent these from getting updated
	# Do not modify the localPasswd if it comes across as blank.
	my $local_passwd = $user->{"localPasswd"};
	if ( defined($local_passwd) && ( $local_passwd eq '' ) ) {
		delete( $user->{"localPasswd"} );
	}

	# Do not modify the confirmLocalPasswd if it comes across as blank.
	my $confirm_local_passwd = $user->{"confirmLocalPasswd"};
	if ( defined($confirm_local_passwd) && ( $confirm_local_passwd eq '' ) ) {
		delete( $user->{"confirmLocalPasswd"} );
	}

	my ( $is_valid, $result ) = $self->is_valid($user);

	if ($is_valid) {
		my $username = $self->current_user()->{username};
		my $dbh = $self->db->resultset('TmUser')->find( { username => $username } );

		# Updating a user implies it is no longer new
		$db_user->{"new_user"} = 0;

		# These if "defined" checks allow for partial user updates, otherwise the entire
		# user would need to be passed through.
		if ( defined($local_passwd) && $local_passwd ne '' ) {
			$db_user->{"local_passwd"} = sha1_hex($local_passwd);
		}
		if ( defined($confirm_local_passwd) && $confirm_local_passwd ne '' ) {
			$db_user->{"confirm_local_passwd"} = sha1_hex($confirm_local_passwd);
		}
		if ( defined( $user->{"id"} ) ) {
			$db_user->{"id"} = $user->{"id"};
		}
		if ( defined( $user->{"username"} ) ) {
			$db_user->{"username"} = $user->{"username"};
		}
		if ( defined( $user->{"tenantId"} ) ) {
 			$db_user->{"tenant_id"} = $user->{"tenantId"};
 		}
		if ( defined( $user->{"public_ssh_key"} ) ) {
			$db_user->{"public_ssh_key"} = $user->{"public_ssh_key"};
		}
		if ( &is_admin($self) && defined( $user->{"role"} ) ) {
			$db_user->{"role"} = $user->{"role"};
		}
		if ( defined( $user->{"uid"} ) ) {
			$db_user->{"uid"} = $user->{"uid"};
		}
		if ( defined( $user->{"gid"} ) ) {
			$db_user->{"gid"} = $user->{"gid"};
		}
		if ( defined( $user->{"company"} ) ) {
			$db_user->{"company"} = $user->{"company"};
		}
		if ( defined( $user->{"email"} ) ) {
			$db_user->{"email"} = $user->{"email"};
		}
		if ( defined( $user->{"fullName"} ) ) {
			$db_user->{"full_name"} = $user->{"fullName"};
		}
		if ( defined( $user->{"newUser"} ) ) {
			$db_user->{"new_user"} = $user->{"newUser"};
		}
		if ( defined( $user->{"addressLine1"} ) ) {
			$db_user->{"address_line1"} = $user->{"addressLine1"};
		}
		if ( defined( $user->{"addressline2"} ) ) {
			$db_user->{"address_line2"} = $user->{"addressLine2"};
		}
		if ( defined( $user->{"city"} ) ) {
			$db_user->{"city"} = $user->{"city"};
		}
		if ( defined( $user->{"stateOrProvince"} ) ) {
			$db_user->{"state_or_province"} = $user->{"stateOrProvince"};
		}
		if ( defined( $user->{"phoneNumber"} ) ) {
			$db_user->{"phone_number"} = $user->{"phoneNumber"};
		}
		if ( defined( $user->{"postalCode"} ) ) {
			$db_user->{"postal_code"} = $user->{"postalCode"};
		}
		if ( defined( $user->{"country"} ) ) {
			$db_user->{"country"} = $user->{"country"};
		}
		$dbh->update($db_user);
		return $self->success_message("UserProfile was successfully updated.");
	}
	else {
		return $self->alert($result);
	}
}

sub is_valid {
	my $self = shift;
	my $user = shift;

	my $rules = {
		fields => [
			qw/fullName username public_ssh_key email role uid gid localPasswd confirmLocalPasswd company newUser addressLine1 addressLine2 city stateOrProvince phoneNumber postalCode country/
		],

		# Checks to perform on all fields
		checks => [

			# All of these are required
			[qw/full_name username email/] => is_required("is required"),

			# pass2 must be equal to pass
			localPasswd => sub {
				my $value  = shift;
				my $params = shift;
				if ( defined( $params->{'localPasswd'} ) ) {
					return $self->is_good_password( $value, $params );
				}
			},

			# email must be unique
			email => sub {
				my $value  = shift;
				my $params = shift;
				if ( defined( $params->{'email'} ) ) {
					return $self->is_email_taken( $value, $params );
				}
			},

			# custom sub validates an email address
			email => sub {
				my ( $value, $params ) = @_;
				Email::Valid->address($value) ? undef : 'email is not a valid format';
			},

			# username must be unique
			username => sub {
				my $value  = shift;
				my $params = shift;
				if ( defined( $params->{'username'} ) ) {
					return $self->is_username_taken( $value, $params );
				}
			},
			
		]
	};

	# Validate the input against the rules
	my $result = validate( $user, $rules );

	if ( $result->{success} ) {

		#print "success: " . dump( $result->{data} );
		return ( 1, $result->{data} );
	}
	else {
		#print "failed " . Dumper( $result->{error} );
		return ( 0, $result->{error} );
	}

}

sub is_username_taken {
	my $self     = shift;
	my $username = shift;

	my $user_with_username = $self->db->resultset('TmUser')->search( { username => $username } )->single;
	if ( defined($user_with_username) ) {
		my $user_id = $user_with_username->id;

		my $current_user = $self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->single;
		my $current_userid = $current_user->id;

		my %condition = ( -and => [ { username => $username }, { id => { 'not in' => [ $current_userid, $user_id ] } } ] );
		my $count = $self->db->resultset('TmUser')->search( \%condition )->count();

		if ( $count > 0 ) {
			return "is already taken";
		}
	}

	return undef;
}

sub is_email_taken {
	my $self   = shift;
	my $email  = shift;

	my $user_with_email = $self->db->resultset('TmUser')->search( { email => $email } )->single;
	if ( defined($user_with_email) ) {
		my $user_id = $user_with_email->id;

		my $current_user = $self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->single;
		my $current_userid = $current_user->id;

		my %condition = ( -and => [ { email => $email }, { id => { 'not in' => [ $current_userid, $user_id ] } } ] );
		my $count = $self->db->resultset('TmUser')->search( \%condition )->count();

		if ( $count > 0 ) {
			return "is already taken";
		}
	}

	return undef;
}

sub is_good_password {
	my $self   = shift;
	my $value  = shift;
	my $params = shift;
	if ( !defined $value or $value eq '' ) {
		return undef;
	}

	if ( $value ne $params->{'confirmLocalPasswd'} ) {
		return "Your 'New Password' must match the 'Confirm New Password'.";
	}

	if ( $value eq $params->{'username'} ) {
		return "Your password cannot be the same as your username.";
	}

	if ( length($value) < 8 ) {
		return "Password must be greater than 7 chars.";
	}

	if ( defined($self->app->{invalid_passwords}->{$value}) ) {
		return "Password is too common.";
	}

	# At this point we're happy with the password
	return undef;
}

1;
