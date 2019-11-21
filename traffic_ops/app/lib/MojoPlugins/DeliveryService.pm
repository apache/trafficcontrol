package MojoPlugins::DeliveryService;
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

use Mojo::Base 'Mojolicious::Plugin';
use Carp qw(cluck confess);
use Data::Dumper;
use Utils::Helper::DateHelper;
use JSON;
use HTTP::Date;
use Common::ReturnCodes qw(SUCCESS ERROR);

sub register {
	my ( $self, $app, $conf ) = @_;

	my $no_instance_message = "Call on an instance of MojoPlugins::DeliveryService!";
	
	$app->renderer->add_helper(

		# ensure param returned as a scalar no matter the context,  and allow default value to be provided
		paramAsScalar => sub {
			my $self    = shift;
			my $p       = shift;
			my $default = shift;

			my $v = $self->param($p);
			if ( !defined $v || $v eq '' ) {
				$v = $default;
			}
			return scalar($v);
		},
	);

	$app->renderer->add_helper(
		hr_string_to_mbps => sub {
			my $self = shift;
			my $inp  = shift;

			if    ( !defined($inp) )     { return 0; }                  # default is 0
			elsif ( $inp =~ /^(\d+)T$/ ) { return $1 * 1000000; }
			elsif ( $inp =~ /^(\d+)G$/ ) { return $1 * 1000; }
			elsif ( $inp =~ /^(\d+)M$/ ) { return $1; }
			elsif ( $inp =~ /^(\d+)k$/ ) { return int( $1 / 1000 ); }
			elsif ( $inp =~ /^(\d+)$/  ) { return $1; }
			else                         { return -1; }

		}
	);

	$app->renderer->add_helper(
		hr_string_to_bps => sub {
			my $self = shift;
			my $inp  = shift;

			if    ( !defined($inp) )     { return 0; }                  # default is 0
			elsif ( $inp =~ /^(\d+)T$/ ) { return $1 * 1000000000000; }
			elsif ( $inp =~ /^(\d+)G$/ ) { return $1 * 1000000000; }
			elsif ( $inp =~ /^(\d+)M$/ ) { return $1 * 1000000; }
			elsif ( $inp =~ /^(\d+)k$/ ) { return $1 * 1000; }
			elsif ( $inp =~ /^(\d+)$/  ) { return $1; }
			else                         { return -1; }

		}
	);
	
	$app->renderer->add_helper(
		is_delivery_service_assigned => sub {
			my $self = shift || confess($no_instance_message);
			my $id   = shift || confess("Please supply a delivery service ID");

			my $user_id =
				$self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->get_column('id')->single();
			my @ds_ids = ();

			if ( defined($user_id) ) {
				@ds_ids = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $user_id } )->get_column('deliveryservice')->all();
			}

			my %ds_hash = map { $_ => 1 } @ds_ids;

			# no external user ID = internal; assume authenticated due to route configuration
			if ( !defined($user_id) ) {
				return (1);
			}
			elsif ($user_id) {
				my $result = $self->db->resultset("Deliveryservice")->search( { id => $id } )->single();

				if ( defined($result) && exists( $ds_hash{ $result->id } ) ) {
					return (1);
				}
			}

			return (0);
		}
	);

	$app->renderer->add_helper(
		is_delivery_service_name_assigned => sub {
			my $self    = shift || confess($no_instance_message);
			my $ds_name = shift || confess("Please supply a delivery service name (xml_id)");

			my $user_id =
				$self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->get_column('id')->single();
			my @ds_ids = ();

			if ( defined($user_id) ) {
				@ds_ids = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $user_id } )->get_column('deliveryservice')->all();
			}

			my %ds_hash = map { $_ => 1 } @ds_ids;

			# no external user ID = internal; assume authenticated due to route configuration
			if ( !defined($user_id) ) {
				return (1);
			}
			elsif ($user_id) {
				my $result = $self->db->resultset("Deliveryservice")->search( { xml_id => $ds_name } )->single();

				if ( exists( $ds_hash{ $result->id } ) ) {
					return (1);
				}
			}

			return (0);
		}
	);

	$app->renderer->add_helper(
		is_valid_delivery_service => sub {
			my $self = shift || confess($no_instance_message);
			my $id   = shift || confess("Please supply a delivery service ID");

			my $result = $self->db->resultset("Deliveryservice")->find( { id => $id } );

			if ( defined($result) ) {
				return (1);
			}
			else {
				return (0);
			}
		}
	);

	$app->renderer->add_helper(
		is_valid_delivery_service_name => sub {
			my $self = shift || confess($no_instance_message);
			my $name = shift || confess("Please supply a delivery service 'name' (xml_id)");

			my $result = $self->db->resultset("Deliveryservice")->find( { xml_id => $name } );

			if ( defined($result) ) {
				return (1);
			}
			else {
				return (0);
			}
		}
	);

	$app->renderer->add_helper(
		get_delivery_service_name => sub {
			my $self = shift || confess($no_instance_message);
			my $id   = shift || confess("Please supply a delivery service ID");

			my $result = $self->db->resultset("Deliveryservice")->search( { id => $id } )->single();

			if ( defined($result) ) {
				return ( $result->xml_id );
			}
			else {
				return (0);
			}
		}
	);

	$app->renderer->add_helper(
		find_existing_host_regex => sub {
                        my $self = shift || confess($no_instance_message);
                        my $regex_type = shift || confess("Please supply a regex type");
                        my $host_regex = shift || confess("Please supply a host regular expression");
                        my $cdn_domain = shift || confess("Please supply a cdn domain");
                        my $cdn_id = shift || confess("Please supply a cdn_name");
                        my $ds_id = shift;

                        if ($regex_type ne 'HOST_REGEXP') {
                                return undef;
                        }

                        my $new_regex = $host_regex . $cdn_domain;
                        my %criteria;
                        $criteria{'pattern'} = $host_regex;
                        my $rs_regex = $self->db->resultset('Regex')->search( \%criteria );
                        while ( my $row = $rs_regex->next ) {
                                my $rs_ds_regex = $self->db->resultset('DeliveryserviceRegex')->search( {  regex => $row->id } );
                                while (my $ds_regex_row = $rs_ds_regex->next) {
                                        if (defined($ds_id) && $ds_id == $ds_regex_row->deliveryservice->id ) { # do not compare if it is the same delivery service
                                                next;
                                        }
                                        my $other_cdn_id = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_regex_row->deliveryservice->id })->single->cdn_id;

                                        if (defined($cdn_id) && $other_cdn_id ne $cdn_id) { # do not compare if not the same cdn.
                                                next;
                                        }
                                        return $new_regex; # at this point we know they are the same and are conflicting.
                                }
                        }
                        return undef;
                }
        );

	$app->renderer->add_helper(
		get_cdn_domain_by_ds_id => sub {
			my $self = shift || confess($no_instance_message);
			my $ds_id = shift || confess("Please supply a delivery service id!");

			my $cdn_id = $self->db->resultset('Deliveryservice')->search( { id => $ds_id } )->get_column('cdn_id')->single();
			my $cdn_domain = $self->db->resultset('Cdn')->search( { id =>  $cdn_id } )->get_column('domain_name')->single();
			return $cdn_domain;
		}
	);

	$app->renderer->add_helper(
		get_cdn_domain_by_profile_id => sub {
			my $self = shift || confess($no_instance_message);
			my $profile_id = shift || confess("Please Supply a profile id");

			my $cdn_id = $self->db->resultset('Profile')->search( { id => $profile_id } )->get_column('cdn')->single();
			my $cdn_domain = $self->db->resultset('Cdn')->search( { id =>  $cdn_id } )->get_column('domain_name')->single();
			return $cdn_domain;
		}
	);

	$app->renderer->add_helper(
		get_profile_id_for_name => sub {
			my $self = shift || confess($no_instance_message);
			my $profile_name = shift || confess("Please Supply a profile name");
			return $self->db->resultset('Profile')->search({'me.name' => $profile_name})->single->id;
		}
	);

	$app->renderer->add_helper(
		get_id_for_cdn_name => sub {
			my $self = shift || confess($no_instance_message);
			my $cdn_name = shift || confess("Please Supply a CDN name");
			return $self->db->resultset('Cdn')->search({'me.name' => $cdn_name})->single->id;
		}
	);
}

1;
