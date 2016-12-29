package API::DeliveryService::SslKeys;
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
use MojoPlugins::Response;
use JSON;
use MIME::Base64;
use UI::DeliveryService;
use Data::Dumper;

sub add {
	my $self     = shift;
	my $key_type = "ssl";
	my $key      = $self->req->json->{key};
	my $version  = $self->req->json->{version};
	my $crt      = $self->req->json->{certificate}->{crt};
	my $csr      = $self->req->json->{certificate}->{csr};
	my $priv_key = $self->req->json->{certificate}->{key};
	my $hostname = $self->req->json->{hostname};
	my $cdn = $self->req->json->{cdn};
	my $deliveryservice = $self->req->json->{deliveryservice};

	if ( !&is_admin($self) ) {
		return $self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	my $record = {
		key => $key,
		version => $version,
		hostname => $hostname,
		cdn => $cdn,
		deliveryservice => $deliveryservice,
		certificate => {
			crt => $crt,
			csr => $csr,
			key => $priv_key
		}
	};

	my $response_container = $self->add_ssl_keys_to_riak( $record );
	my $response = $response_container->{"response"};
	if ( $response->is_success() ) {
		&log( $self, "Added ssl keys for Delivery Service $key", "APICHANGE" );
		return $self->success("Successfully added ssl keys for $key");
	}
	else {
		return $self->alert( $response->{_content} );
	}
}

#named like this because there is a plugin called generate_ssl_keys in Mojoplugins/SslKeys.pm
sub generate {
	my $self         = shift;
	my $key_type     = "ssl";
	my $key          = $self->req->json->{key};
	my $version      = $self->req->json->{version};        # int
	my $hostname     = $self->req->json->{hostname};
	my $country      = $self->req->json->{country};
	my $state        = $self->req->json->{state};
	my $city         = $self->req->json->{city};
	my $org          = $self->req->json->{organization};
	my $unit         = $self->req->json->{businessUnit};
	my $cdn = $self->req->json->{cdn};
	my $deliveryservice = $self->req->json->{deliveryservice};
	my $tmp_location = "/var/tmp";

	if ( !&is_admin($self) ) {
		return $self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}

	#generate the cert:
	my $record = {
		key => $key,
		version => $version,
		hostname => $hostname,
		country => $country,
		city => $city,
		state => $state,
		org => $org,
		unit => $unit,
		cdn => $cdn,
		deliveryservice => $deliveryservice,
	};

	my $response_container = $self->generate_ssl_keys( $record );
	my $response = $response_container->{"response"};
	if ( $response->is_success() ) {
		&log( $self, "Created ssl keys for Delivery Service $key", "APICHANGE" );
		return $self->success("Successfully created ssl keys for $key");
	}
	else {
		return $self->alert( $response->{_content} );
	}
}

sub view_by_xml_id {
	my $self    = shift;
	my $key     = $self->param('xmlid');
	my $version = $self->param('version');
	if ( !&is_admin($self) ) {
		return $self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		if ( !$version ) {
			$version = 'latest';
		}
		$key = "$key-$version";
		my $response_container = $self->riak_get( "ssl", $key );
		my $response = $response_container->{"response"};
		$response->is_success()
			? $self->success( decode_json( $response->content ) )
			: $self->alert( { Error => " - A record for ssl key $key could not be found.  Response was: " . $response->content } );
	}
}

sub view_by_hostname {
	my $self    = shift;
	my $key     = $self->param('hostname');
	my $version = $self->param('version');

	if ( !&is_admin($self) ) {
		return $self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		#use hostname to get hostname regex
		my @split_url = split( /\./, $key );
		my $host_regex = $split_url[1];
		my $domain_name;

		for ( my $i = 2; $i < $#split_url; $i++ ) {
			$domain_name .= $split_url[$i] . ".";
		}
		$domain_name .= $split_url[$#split_url];

		$host_regex = '.*\.' . $host_regex . '\..*';

		if ( !$host_regex || !$domain_name ) {
			return $self->alert( { Error => " - $key does not contain a valid delivery service." } ) if !$host_regex;
			return $self->alert( { Error => " - $key does not contain a valid domain name." } )      if !$domain_name;
		}

		my @ds_ids_regex = $self->db->resultset('Deliveryservice')
			->search( { 'regex.pattern' => "$host_regex" }, { join => { deliveryservice_regexes => { regex => undef } } } )->get_column('id')->all();

		# TODO JvD - test this with online riak servers!
		my $cdn_id = $self->db->resultset('Cdn')->search( { domain_name => $domain_name } )->get_column('id')->single();
		my@domain_profiles = $self->db->resultset('Profile')->search( { cdn => $cdn_id } )->get_column('id')->all();

		my $rs_ds = $self->db->resultset('Deliveryservice')->search( { 'profile' => { -in => \@domain_profiles } }, {} );

		my $xml_id;
		my %ds_ids_regex = map { $_ => undef } @ds_ids_regex;

		while ( my $row = $rs_ds->next ) {
			if ( exists( $ds_ids_regex{ $row->id } ) ) {
				$xml_id = $row->xml_id;
			}
		}

		if ( !$version ) {
			$version = 'latest';
		}
		$key = "$xml_id-$version";
		my $response_container = $self->riak_get( "ssl", $key );
		my $response = $response_container->{"response"};
		$response->is_success()
			? $self->success( decode_json( $response->content ) )
			: $self->alert( { Error => " - A record for ssl key $key could not be found.  Response was: " . $response->content } );
	}
}

sub delete {
	my $self    = shift;
	my $key     = $self->param('xmlid');
	my $version = $self->param('version');
	my $response_container;
	my $response;
	if ( !&is_admin($self) ) {
		return $self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		if ($version) {
			$key = $key . "-" . $version;
			$self->app->log->info("deleting key_type = ssl, key = $key");
			$response_container = $self->riak_delete( "ssl", $key );
			$response = $response_container->{"response"};
		}
		else {
			#TODO figure out riak searching so we dont have to hardcode "latest"
			$key = "$key-latest";
			$self->app->log->info("deleting key_type = ssl, key = $key");
			$response_container = $self->riak_delete( "ssl", $key );
			$response = $response_container->{"response"};
		}

		# $self->app->log->info("delete rc = $rc");
		if ( $response->is_success() ) {
			&log( $self, "Deleted ssl keys for Delivery Service $key", "APICHANGE" );
			return $self->success("Successfully deleted ssl keys for $key");
		}
		else {
			return $self->alert( $response->content );
		}
	}
}

1;
