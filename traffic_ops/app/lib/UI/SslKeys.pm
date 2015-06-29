package UI::SslKeys;
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
use UI::DeliveryService;
use Scalar::Util qw(looks_like_number);
use JSON;
use MIME::Base64;

sub add {
	my $self  = shift;
	my $ds_id = $self->param('id');
	my $rs_ds  = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_id } );
	my $data   = $rs_ds->single;
	my $xml_id = $data->xml_id;
	&stash_role($self);

	#get key data from keystore
	my $response_container = $self->riak_get( 'ssl', "$xml_id-latest");
	my $get_keys = $response_container->{'response'};
	if ( $get_keys->is_success() ) {
		my $keys = decode_json( $get_keys->content );
		my $version = $keys->{version} + 1;
		$self->stash(
			ssl => {
				country  => $keys->{country},
				state    => $keys->{state},
				city     => $keys->{city},
				org      => $keys->{organization},
				unit     => $keys->{businessUnit},
				hostname => $keys->{hostname},
				csr		 => decode_base64($keys->{certificate}->{csr}),
				crt		 => decode_base64($keys->{certificate}->{crt}),
				priv_key => decode_base64($keys->{certificate}->{key}),
				version  => $version
			},
			xml_id      => $xml_id,
			fbox_layout => 1
		);
	}
	else {
		my $ssl_key_version = $data->ssl_key_version;
		my $new_version     = $ssl_key_version + 1;
		my $domain_name     = UI::DeliveryService::get_cdn_domain( $self, $ds_id );
		my $ds_regexes      = UI::DeliveryService::get_regexp_set( $self, $ds_id );
		my @example_urls    = UI::DeliveryService::get_example_urls( $self, $ds_id, $ds_regexes, $data, $domain_name, $data->protocol );

		#first one is the one we want
		my $hostname = $example_urls[0];
		$hostname =~ /(https?:\/\/)(.*)/;

		$self->stash(
			ssl => {
				version  => $new_version,
				hostname => $2,
			},
			xml_id      => $xml_id,
			fbox_layout => 1
		);
	}
}

sub create {
	my $self = shift;
	##Check to see if we are adding existing keys or generating new ones.
	my $action   = $self->param('ssl.action');
	my $country  = $self->param('ssl.country');
	my $state    = $self->param('ssl.state');
	my $city     = $self->param('ssl.city');
	my $org      = $self->param('ssl.org');
	my $unit     = $self->param('ssl.unit');
	my $hostname = $self->param('ssl.hostname');
	my $version  = $self->param('ssl.version');

	# get ds info
	my $xml_id = $self->param('xml_id');
	my $rs_ds =
		$self->db->resultset('Deliveryservice')->search( { 'me.xml_id' => $xml_id }, { prefetch => [ { 'type' => undef }, { 'profile' => undef } ] } );
	my $data = $rs_ds->single;
	my $id   = $data->id;

	if ( !&is_admin($self) ) {
		$self->flash( alertmsg => "Keys can only be added by admins!" );
		return $self->redirect_to("/ds/$id/sslkeys/add");
	}

	if ( $self->is_valid() ) {
		my $response_container;
		if ( $action eq "add" ) {

			#add existing keys to keystore
			my $csr      = $self->param('ssl.csr');
			my $crt      = $self->param('ssl.crt');
			my $priv_key = $self->param('ssl.priv_key');
			$response_container = $self->add_ssl_keys_to_riak( $xml_id, $version, $crt, $csr, $priv_key );
		}
		else {
			#generate keys
			#add to keystore
			$response_container = $self->generate_ssl_keys( $hostname, $country, $city, $state, $org, $unit, $version, $xml_id );
		}

		#update version in db
		my $update = $self->db->resultset('Deliveryservice')->find( { id => $id } );
		$update->ssl_key_version($version);
		$update->update();
		my $response = $response_container->{"response"};
		if ( $response->is_success ) {
			&log( $self, "Created ssl keys for Delivery Service $xml_id", "APICHANGE" );
			$self->flash( message => "Successfully created ssl keys for: $xml_id" );
		}
		else {
			$self->flash( { Error => " - SSL keys for '$xml_id' could not be created.  Response was" . $response->{_content} } );
		}
		return $self->redirect_to("/ds/$id/sslkeys/add");
	}
	else {
		&stash_role($self);
		$self->stash(
			ssl => {
				country    => $country,
				state      => $state,
				city       => $city,
				org        => $org,
				unit       => $unit,
				domainName => $hostname,
				version    => $version,
			},
			xml_id      => $xml_id,
			fbox_layout => 1
		);
		return $self->render("ssl_keys/add");
	}
}

sub is_valid {
	my $self   = shift;
	my $action = $self->param('ssl.action');

	if ( $action eq "add" ) {
		$self->field('ssl.csr')->is_required("Certificate Signing Request cannot be empty.");
		$self->field('ssl.crt')->is_required("Certificate cannot be empty.");
		$self->field('ssl.priv_key')->is_required("Private Key cannot be empty.");
	}
	else {
		my $country  = $self->param('ssl.country');
		my $state    = $self->param('ssl.state');
		my $city     = $self->param('ssl.city');
		my $org      = $self->param('ssl.org');
		my $unit     = $self->param('ssl.unit');
		my $hostname = $self->param('ssl.hostname');
		my $xml_id   = $self->param('xml_id');
		my $version  = $self->param('ssl.version');

		$self->field('ssl.country')->is_required("Country cannot be empty");
		$self->field('ssl.state')->is_required("State cannot be empty");
		$self->field('ssl.city')->is_required("City cannot be empty");
		$self->field('ssl.org')->is_required("Organization cannot be empty");
		$self->field('ssl.unit')->is_required("Unit cannot be empty");

		if ( length($country) != 2 ) {
			$self->field('ssl.country')->is_equal( "", "Country code must be 2 characters only!" );
		}

		#strip off http(s)://
		if ( !&is_hostname($hostname) ) {
			$self->field('ssl.hostname')->is_equal( "", "$hostname is not a valid hostname!" );
		}

		my $rs_ds =
			$self->db->resultset('Deliveryservice')->search( { 'me.xml_id' => $xml_id }, { prefetch => [ { 'type' => undef }, { 'profile' => undef } ] } );
		my $data            = $rs_ds->single;
		my $ssl_key_version = $data->ssl_key_version;

		if ( $version < $ssl_key_version ) {
			$self->field('ssl.version')->is_equal( "", "Version must be greater than the latest version $ssl_key_version" );
		}
	}
	return $self->valid;
}

1;
