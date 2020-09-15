package MojoPlugins::DnssecKeys;
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
use Net::DNS;
use MIME::Base64;
use Crypt::OpenSSL::RSA;
use Crypt::OpenSSL::Bignum;
use Crypt::OpenSSL::Random;
use Net::DNS::SEC::Private;
use Net::DNS::RR::DS;
use Data::Dumper;
use JSON;
my $TMP_LOCATION = "/var/tmp";

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		generate_store_dnssec_keys => sub {
			my $self          = shift;
			my $key           = shift;
			my $name          = shift;
			my $ttl           = shift;
			my $k_exp_days    = shift;
			my $z_exp_days    = shift;
			my $effectiveDate = shift;
			my $keys          = {};

			my $inception    = time();
			my $z_expiration = time() + ( 86400 * $z_exp_days );
			my $k_expiration = time() + ( 86400 * $k_exp_days );

			#add "." to the end of name if not already there
			if ( ( substr( $name, -1 ) ) ne "." ) {
				$name = $name . ".";
			}

			#get old keys if they exist
			my $old_keys           = {};
			my $response_container = $self->riak_get( "dnssec", $key );
			my $get_keys           = $response_container->{'response'};
			if ( $get_keys->is_success() ) {
				$old_keys = decode_json( $get_keys->content );
			}

			# #create new keys for cdn TLD
			$self->app->log->info("Creating keys for $key.");
			my @zsk = $self->get_dnssec_keys( "zsk", $name, $ttl, $inception,
				$z_expiration, "new", $effectiveDate );
			my @ksk = $self->get_dnssec_keys( "ksk", $name, $ttl, $inception,
				$k_expiration, "new", $effectiveDate, "1" );

			#get old ksk
			my $krecord
				= &get_existing_record( $self, $old_keys, $key, "ksk" );
			if ( defined($krecord) ) {
				$krecord->{status}         = "existing";
				$krecord->{expirationDate} = $effectiveDate;
				push @ksk, $krecord;
			}

			#get old zsk
			my $zrecord
				= &get_existing_record( $self, $old_keys, $key, "zsk" );
			if ( defined($zrecord) ) {
				$zrecord->{status}         = "existing";
				$zrecord->{expirationDate} = $effectiveDate;
				push @zsk, $zrecord;
			}

			#add to keys hash
			$keys->{$key} = { zsk => [@zsk], ksk => [@ksk] };

			#find the cdn's delivery services to generate keys for
			my $cdn = $self->db->resultset('Cdn')->find( { name => $key } );
			my %search = ( cdn_id => $cdn->id );
			my @ds_rs = $self->db->resultset('Deliveryservice')->search( \%search , { prefetch => [ 'type' ] } );
			foreach my $ds (@ds_rs) {
				if (   $ds->type->name !~ m/^HTTP/
					&& $ds->type->name !~ m/^DNS/ )
				{
					next;
				}
				my $xml_id = $ds->xml_id;
				my $ds_id  = $ds->id;

				#create the ds domain name for dnssec keys
				my $cdn_domain_name = $cdn->domain_name;
				my $ds_name = UI::DeliveryService::get_ds_domain_name($self, $ds_id, $xml_id, $cdn_domain_name);

				$self->app->log->info("Creating keys for $xml_id.");
				my @zsk = $self->get_dnssec_keys( "zsk", $ds_name, $ttl,
					$inception, $z_expiration, "new", $effectiveDate );
				my @ksk = $self->get_dnssec_keys( "ksk", $ds_name, $ttl,
					$inception, $k_expiration, "new", $effectiveDate );

				#get old ksk
				my $krecord = &get_existing_record( $self, $old_keys, $xml_id,
					"ksk" );
				if ( defined($krecord) ) {
					$krecord->{status}         = "existing";
					$krecord->{expirationDate} = $effectiveDate;
					push @ksk, $krecord;
				}

				#get old zsk
				my $zrecord = &get_existing_record( $self, $old_keys, $xml_id,
					"zsk" );
				if ( defined($zrecord) ) {
					$zrecord->{status}         = "existing";
					$zrecord->{expirationDate} = $effectiveDate;
					push @zsk, $zrecord;
				}

				#add to keys hash
				$keys->{$xml_id} = { zsk => [@zsk], ksk => [@ksk] };
			}

			my $json_data = encode_json($keys);
			my $response = $self->riak_put( "dnssec", $key, $json_data );

			return $response;
		}
	);
	$app->renderer->add_helper(
		get_profile_id_by_cdn => sub {
			my $self     = shift;
			my $cdn_name = shift;

			my %condition = (
				-and => [
					'cdn.name' => $cdn_name,
					{   -or => [
							'profile.name' => { like => "CCR%" },
							'profile.name' => { like => 'TR%' }
						]
					}
				]
			);
			my $profile_id = $self->db->resultset('Server')->search(
				\%condition,
				{   prefetch => [ 'cdn', 'profile' ],
					select   => 'me.profile',
					distinct => 1
				}
			)->get_column('profile')->single();
			return $profile_id;
		}
	);

	$app->renderer->add_helper(
		get_dnssec_keys => sub {
			my $self          = shift;
			my $type          = shift;
			my $name          = shift;
			my $ttl           = shift;
			my $inception     = shift;
			my $expiration    = shift;
			my $status        = shift;
			my $effectiveDate = shift;
			my $tld           = shift;
			my %keys          = ();
			my %response      = (
				inceptionDate  => $inception,
				expirationDate => $expiration,
				name           => $name,
				ttl            => $ttl,
				status         => $status,
				effectiveDate  => $effectiveDate
			);

			if ( $type eq "zsk" ) {
				%keys = &gen_keys( $self, $name, 0, $ttl );
			}

			else {
				%keys = &gen_keys( $self, $name, 1, $ttl, $tld );
			}

			#add keys to response
			$response{private} = $keys{private_key};
			$response{public}  = $keys{public_key};
			if ($tld) {
				$response{dsRecord} = $keys{ds_record};
			}
			return \%response;
		}
	);

	sub gen_keys {
		my $self = shift;
		my $name = shift;
		my $ksk  = shift;
		my $ttl  = shift;
		my $tld  = shift || 0;
		my $bits = 1024;
		my $flags |= 256;
		my $algorithm = 5
			; # http://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml
		my $protocol = 3;
		my %response = ();

		if ($ksk) {
			$flags |= 1;    # ksk
			$bits *= 2;
		}

		my $keypair
			= Net::DNS::SEC::Private->generate_rsa( $name, $flags, $bits,
			$algorithm );
		my $private_key = encode_base64( $keypair->dump_rsa_priv );

		#trim whitespace
		$private_key =~ s/\s+$//;
		$response{private_key} = $private_key;

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
		$public_key =~ s/\s+$//;
		$response{public_key} = $public_key;

		#create ds record
		if ( $ksk && $tld ) {
			my $ds_rr = create Net::DNS::RR::DS(
				$dnskey_rr,
				digtype => 'SHA-256',
				ttl     => $ttl
			);
			my %ds_record = (
				digest     => $ds_rr->digest,
				digestType => $ds_rr->digtype,
				algorithm  => $ds_rr->algorithm
			);
			$response{ds_record} = \%ds_record;
		}
		return %response;
	}

	sub get_existing_record {
		my $self = shift;
		my $keys = shift;
		my $key  = shift;
		my $type = shift;

		my $existing = $keys->{$key}->{$type};
		foreach my $record (@$existing) {
			if ( $record->{status} eq 'new' ) {
				return $record;
			}
		}
		return undef;
	}

}

1;
