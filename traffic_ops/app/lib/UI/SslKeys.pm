package UI::SslKeys;
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
use UI::DeliveryService;
use Scalar::Util qw(looks_like_number);
use JSON;
use MIME::Base64;

sub add {
	my $self  = shift;
	my $ds_id = $self->param('id');
	my $rs_ds  = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_id }, { prefetch => [ 'cdn', 'type' ] } );
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
				org      => $keys->{org},
				unit     => $keys->{unit},
				hostname => defined( $keys->{hostname} ) ? $keys->{hostname} : $self->get_hostname($ds_id, $data),
				cdn => defined($keys->{cdn}) ? $keys->{cdn} : $data->cdn->name,
				deliveryservice => defined($keys->{deliveryservice}) ? $keys->{deliveryservice} : $xml_id,
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

		$self->stash(
			ssl => {
				version  => $new_version,
				hostname => $self->get_hostname($ds_id, $data),
				cdn => $data->cdn->name,
				deliveryservice => $xml_id
			},
			xml_id      => $xml_id,
			fbox_layout => 1
		);
	}
}

sub update_sslkey {
	my $self = shift;
	my $xml_id = shift;
	my $hostname = shift;
	my $response_container = $self->riak_get( 'ssl', "$xml_id-latest");
	my $response = $response_container->{'response'};

	if ( $response->is_success() ) {
		my $record = decode_json( $response->content );
		$record->{deliveryservice} = $xml_id;
		$record->{hostname} = $hostname;
		my $key = $xml_id;
		my $version = $record->{version};

		$response_container = $self->riak_put( 'ssl', "$key-$version", encode_json($record) );
		$response = $response_container->{'response'};
		if ( !$response->is_success() ) {
			$self->app->log->warn("SSL keys for '$key-$version' could not be updated.  Response was " . $response_container->{_content});
		}
		$response_container = $self->riak_put( 'ssl', "$key-latest", encode_json($record) );
		$response = $response_container->{'response'};
		if ( !$response->is_success() ) {
			$self->app->log->warn("SSL keys for '$key-latest' could not be updated.  Response was " . $response_container->{_content});
		}
	}
}

sub get_hostname {
	my $self = shift;
	my $ds_id = shift;
	my $data = shift;

	my $domain_name = $data->cdn->domain_name;
	my $ds_regexes      = UI::DeliveryService::get_regexp_set( $self, $ds_id );
	my @example_urls    = UI::DeliveryService::get_example_urls( $self, $ds_id, $ds_regexes, $data, $domain_name, $data->protocol );

	#if a DS is https only we want the first example_url
	my $url = $example_urls[0];
	#if a DS is http/https then we want the second one...see https://github.com/Comcast/traffic_control/issues/1268
	if ($data->protocol == 2) {
		$url = $example_urls[1];
	}
	$url =~ s/(https?:\/\/)(.*)//g;
	my $hostname = $2;
	if ($data->type->name =~ m/^HTTP/) {
		#remove routing name and replace with * for wildcard
		my @split_hostname = split(/\./,$hostname);
		$hostname = '*.' . join('.', splice(@split_hostname, 1));
	}
	return $hostname;
}

sub create {
	my $self = shift;
	my $action   = $self->param('ssl.action');
	my $country  = $self->param('ssl.country');
	my $state    = $self->param('ssl.state');
	my $city     = $self->param('ssl.city');
	my $org      = $self->param('ssl.org');
	my $unit     = $self->param('ssl.unit');
	my $hostname = $self->param('ssl.hostname');
	my $version  = $self->param('ssl.version');
	my $cdn  = $self->param('ssl.cdn');
	my $deliveryservice = $self->param('ssl.deliveryservice');

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
		my $record = {
			key => $xml_id,
			version => $version,
			hostname => defined($hostname) ? $hostname : $self->get_hostname($id, $data),
			cdn => defined($cdn) ? $cdn : $data->cdn->name,
			deliveryservice => defined($deliveryservice) ? $deliveryservice : $xml_id
		};
		if ( $action eq "add" ) {
			$record->{certificate}->{crt} = $self->param('ssl.crt');
			$record->{certificate}->{csr} = $self->param('ssl.csr');
			$record->{certificate}->{key} = $self->param('ssl.priv_key');
			$response_container = $self->add_ssl_keys_to_riak( $record );
		}
		else {
			$record->{country} = $country;
			$record->{city} = $city;
			$record->{state} = $state;
			$record->{org} = $org;
			$record->{unit} = $unit;

			$response_container = $self->generate_ssl_keys($record);

		}

		#update version in db
		my $update = $self->db->resultset('Deliveryservice')->find( { id => $id } );
		$update->ssl_key_version($version);
		$update->update();
		my $response = $response_container->{"response"};
		if ( defined($response) && $response->is_success ) {
			&log( $self, "Created ssl keys for Delivery Service $xml_id", "APICHANGE" );
			$self->flash( message => "Successfully created ssl keys for: $xml_id" );
		}
		else {
			$self->app->log->warn("SSL keys for '$xml_id' could not be created.  Response was " . $response_container->{_content});
			$self->flash( alertmsg => "SSL keys for $xml_id could not be created.  Response was: " . $response_container->{_content}  );
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
