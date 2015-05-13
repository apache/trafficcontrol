package MojoPlugins::DnssecKeys;
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

use Mojo::Base 'Mojolicious::Plugin';
use MIME::Base64;
use Net::DNS;
use MIME::Base64;
use Crypt::OpenSSL::RSA;
use Crypt::OpenSSL::Bignum;
use Crypt::OpenSSL::Random;
use Net::DNS::SEC::Private;
my $TMP_LOCATION = "/var/tmp";

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		generate_store_dnssec_keys => sub {
			my $key_type   = "dnssec";
			my $self       = shift;
			my $key        = shift;
			my $name       = shift;
			my $ttl        = shift;
			my $k_exp_days = shift;
			my $z_exp_days = shift;

			my $inception    = time();
			my $z_expiration = time() + ( 86400 * $z_exp_days );
			my $k_expiration = time() + ( 86400 * $k_exp_days );

			my $json = JSON->new;
			my %keys = ();

			#add "." to the end of name if not already there
			if ( ( substr( $name, -1 ) ) ne "." ) {
				$name = $name . ".";
			}

			# #create keys for cdn TLD
			$self->app->log->info("Creating keys for $key.");
			my $zsk = $self->get_dnssec_keys( "zsk", $name, $ttl, $inception, $z_expiration );
			my $ksk = $self->get_dnssec_keys( "ksk", $name, $ttl, $inception, $k_expiration );

			#add to keys hash
			$keys{$key} = { zsk => $zsk, ksk => $ksk };

			#get delivery services
			#first get profile_id
			my $profile_id = $self->get_profile_id_by_cdn($key);

			#then get deliveryservices
			my %search = ( profile => $profile_id );
			my @ds_rs = $self->db->resultset('Deliveryservice')->search( \%search );
			foreach my $ds (@ds_rs) {
				my $xml_id = $ds->xml_id;
				my $ds_id  = $ds->id;

				#create the ds domain name for dnssec keys
				my $domain_name = UI::DeliveryService::get_cdn_domain( $self, $ds_id );
				my $ds_regexes = UI::DeliveryService::get_regexp_set( $self, $ds_id );
				my $rs_ds =
					$self->db->resultset('Deliveryservice')
					->search( { 'me.xml_id' => $xml_id }, { prefetch => [ { 'type' => undef }, { 'profile' => undef } ] } );
				my $data = $rs_ds->single;
				my @example_urls = UI::DeliveryService::get_example_urls( $self, $ds_id, $ds_regexes, $data, $domain_name, $data->protocol );

				#first one is the one we want.  period at end for dnssec, substring off stuff we dont want
				my $ds_name = $example_urls[0] . ".";
				my $length = length($ds_name) - index( $ds_name, "." );
				$ds_name = substr( $ds_name, index( $ds_name, "." ) + 1, $length );
				$self->app->log->info("Creating keys for $xml_id.");
				my $zsk = $self->get_dnssec_keys( "zsk", $ds_name, $ttl, $inception, $z_expiration );
				my $ksk = $self->get_dnssec_keys( "ksk", $ds_name, $ttl, $inception, $k_expiration );

				#add to keys hash
				$keys{$xml_id} = { zsk => $zsk, ksk => $ksk };
			}

			#add a param to the database to track changes
			#check to see if param already exists
			my $param_id =
				$self->db->resultset('Parameter')->search( { name => $key . ".dnssec.inception", config_file => "CRConfig.json" } )->get_column('id')
				->single();

			#if exists, update
			if ( defined($param_id) ) {
				my $param_update = $self->db->resultset('Parameter')->find( { id => $param_id } );
				$param_update->value($inception);
				$param_update->update();
			}

			#else insert param
			else {
				my $param_insert = $self->db->resultset('Parameter')->create(
					{
						name        => $key . ".dnssec.inception",
						config_file => "CRConfig.json",
						value       => $inception,
					}
				);
				$param_insert->insert();
				$param_id = $param_insert->id();

				#insert into profile_param
				my $pp_insert = $self->db->resultset('ProfileParameter')->create(
					{
						profile   => $profile_id,
						parameter => $param_id,
					}
				);
				$pp_insert->insert();

			}

			my $json_data = $json->encode( \%keys );
			my $response = $self->riak_put( $key_type, $key, $json_data );

			return $response;
		}
	);
	$app->renderer->add_helper(
		get_profile_id_by_cdn => sub {
			my $self       = shift;
			my $cdn_name   = shift;
			my $profile_id = $self->db->resultset('Profile')->search(
				{ -and => [ 'parameter.name' => 'CDN_Name', 'parameter.value' => $cdn_name, 'me.name' => { -like => 'CCR%' } ] },
				{
					join => { profile_parameters => { parameter => undef } },
				}
			)->get_column('id')->single();
			return $profile_id;
		}
	);

	$app->renderer->add_helper(
		get_dnssec_keys => sub {
			my $self       = shift;
			my $type       = shift;
			my $name       = shift;
			my $ttl        = shift;
			my $inception  = shift;
			my $expiration = shift;
			my %keys       = ();

			if ( $type eq "zsk" ) {
				%keys = &gen_keys( $self, $name, 0, $ttl );
			}

			else {
				%keys = &gen_keys( $self, $name, 1, $ttl );
			}

			#store in hash for response
			my %response = (
				private        => $keys{private_key},
				public         => $keys{public_key},
				inceptionDate  => $inception,
				expirationDate => $expiration,
				name           => $name,
				ttl            => $ttl
			);
			return \%response;
		}
	);

	sub gen_keys {
		my $self = shift;
		my $name = shift;
		my $ksk  = shift;
		my $ttl  = shift;
		my $bits = 2048;
		my $flags |= 256;
		my $algorithm = 5;    # http://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml
		my $protocol  = 3;

		if ($ksk) {
			$flags |= 1;      # ksk
			$bits *= 2;
		}

		my $keypair = Net::DNS::SEC::Private->generate_rsa( $name, $flags, $bits, $algorithm );
		my $private_key = encode_base64( $keypair->dump_rsa_priv );

		my $dnskey_rr = new Net::DNS::RR(
			name      => $name,
			type      => "DNSKEY",
			flags     => $flags,
			protocol  => $protocol,
			algorithm => $algorithm,
			publickey => $keypair->dump_rsa_pub,
			ttl       => $ttl
		);
		my $public_key = encode_base64( $dnskey_rr->plain );

		#trim whitespace
		$private_key =~ s/\s+$//;
		$public_key =~ s/\s+$//;

		my %response = (
			private_key => $private_key,
			public_key  => $public_key,
		);
		return %response;
	}

}

1;
