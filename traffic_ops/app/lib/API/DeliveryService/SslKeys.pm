package API::DeliveryService::SslKeys;
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
use Utils::Helper::Datasource;
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

	if ( !&is_admin($self) ) {
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}

	else {

	}
	my $response_container = $self->add_ssl_keys_to_riak( $key, $version, $crt, $csr, $priv_key );
	my $response = $response_container->{"response"};
	if ( $response->is_success() ) {
		&log( $self, "Added ssl keys for Delivery Service $key", "APICHANGE" );
		return $self->success("Successfully added ssl keys for $key");
	}
	else {
		return $self->alert( $response->{_content} );
	}
}

#named like this cause there is a plugin called generate_ssl_keys
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
	my $tmp_location = "/var/tmp";

	if ( !&is_admin($self) ) {
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
	}
	else {
		#generate the cert:
		my $response_container = $self->generate_ssl_keys( $hostname, $country, $city, $state, $org, $unit, $version, $key );
		my $response = $response_container->{"response"};
		if ( $response->is_success() ) {
			&log( $self, "Created ssl keys for Delivery Service $key", "APICHANGE" );
			return $self->success("Successfully created ssl keys for $key");
		}
		else {
			return $self->alert( $response->{_content} );
		}
	}
}

sub view_by_xml_id {
	my $self    = shift;
	my $key     = $self->param('xmlid');
	my $version = $self->param('version');
	if ( !&is_admin($self) ) {
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
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
		if ( !$host_regex ) {
			return $self->alert( { Error => " - $key is not a valid hostname." } );
		}
		my $xml_id = $self->db->resultset('Deliveryservice')->search(
			{ -and => [ 'regex.pattern' => [ { like => "%$host_regex%" } ] ] },
			{
				join     => { deliveryservice_regexes => { regex => undef } },
				distinct => 1
			}
		)->get_column('xml_id')->single();

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
		$self->alert( { Error => " - You must be an ADMIN to perform this operation!" } );
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
