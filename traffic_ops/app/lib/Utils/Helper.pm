package Utils::Helper;
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

use Carp qw(cluck confess);
use Data::Dumper;
use Digest::SHA1 qw(sha1_hex);
use Crypt::ScryptKDF qw(scrypt_hash scrypt_hash_verify);

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	$self->{MOJO} = $args->{mojo} || confess("The mojo argument is required");

	return ( bless( $self, $class ) );
}

sub mojo {
	my $self = shift || confess("Call on an instance of Utils::Helper");
	return ( $self->{MOJO} );
}

sub escape_influxql {
	my $query = shift;
	$query=~s#\\#\\\\#g;
	$query=~s#'#\\'#g;
	$query=~s#"#\\"#g;
	$query=~s#\n#\\n#g;
	return $query
}

sub is_valid_delivery_service {
	my $self = shift || confess("Call on an instance of Utils::Helper");
	my $id   = shift || confess("Please supply a delivery service ID");

	my $result = $self->mojo->db->resultset("Deliveryservice")->find( { id => $id } );

	if ( defined($result) ) {
		return (1);
	}
	else {
		return (0);
	}
}

sub get_delivery_service_name {
	my $self = shift || confess("Call on an instance of Utils::Helper");
	my $id   = shift || confess("Please supply a delivery service ID");

	my $result = $self->mojo->db->resultset("Deliveryservice")->search( { id => $id } )->single();

	if ( defined($result) ) {
		return ( $result->xml_id );
	}
	else {
		return (0);
	}
}

sub is_delivery_service_assigned {
	my $self = shift || confess("Call on an instance of Utils::Helper");
	my $id   = shift || confess("Please supply a delivery service ID");

	my $user_id = $self->mojo->db->resultset('TmUser')->search( { username => $self->mojo->current_user()->{username} } )->get_column('id')->single();
	my @ds_ids = ();

	if ( defined($user_id) ) {
		@ds_ids = $self->mojo->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $user_id } )->get_column('deliveryservice')->all();
	}

	my %ds_hash = map { $_ => 1 } @ds_ids;

	# no external user ID = internal; assume authenticated due to route configuration
	if ( !defined($user_id) ) {
		return (1);
	}
	elsif ($user_id) {
		my $result = $self->mojo->db->resultset("Deliveryservice")->search( { id => $id } )->single();

		if ( exists( $ds_hash{ $result->id } ) ) {
			return (1);
		}
	}

	return (0);
}

sub not_found {
	my $self = shift || confess("Call on an instance of Utils::Helper");

	$self->mojo->render(
		status => 404,
		json   => {
			message => {
				type    => "error",
				content => "Resource not found"
			}
		}
	);
}

sub forbidden {
	my $self = shift || confess("Call on an instance of Utils::Helper");

	$self->mojo->render(
		status => 403,
		json   => {
			message => {
				type    => "error",
				content => "Forbidden"
			}
		}
	);
}

sub error {
	my $self = shift || confess("Call on an instance of Utils::Helper");

	$self->mojo->render(
		status => 500,
		json   => {
			message => {
				type    => "error",
				content => "Internal server error"
			}
		}
	);
}

sub hash_pass {
	my $pass = shift;
	return scrypt_hash($pass, \64, 16384, 8, 1, 64);
}

# verify_pass verifies that the given plaintext password matches the hashed password from the database.
# It first checks if it's scrypt, and if that verification fails, it checks if it's an old insecure sha1 hash. The insecure sha1 fallback is deprecated, and will be removed in the next major version.
sub verify_pass {
	my $pass        = shift;
	my $hashed_pass = shift;
	if ( scrypt_hash_verify($pass, $hashed_pass) ) {
		return 1;
	}
	return $hashed_pass eq sha1_hex($pass); # DEPRECATED - remove in next major version
}

1;
