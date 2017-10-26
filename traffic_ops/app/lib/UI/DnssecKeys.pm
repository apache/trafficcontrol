package UI::DnssecKeys;
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
use API::Cdn;
use Scalar::Util qw(looks_like_number);
use JSON;
use POSIX qw(strftime);
use Date::Parse;

sub index {
	my $self = shift;

	#get a list of cdns from parameters
	&navbarpage($self);
	my @cdns = $self->db->resultset('Cdn')->search( {} )->get_column('name')
		->all();
	$self->stash( cdns => \@cdns, );
}

sub manage {
	my $self = shift;
	&stash_role($self);
	my $cdn_name = $self->param('cdn_name');
	my $ttl;
	my $k_expiry;
	my $algorithm;
	my $digest_type;
	my $digest;

	#get keys for cdn:
	my $keys;
	my $response_container = $self->riak_get( "dnssec", $cdn_name );
	my $get_keys = $response_container->{'response'};
	if ( $get_keys->is_success() ) {
		$keys = decode_json( $get_keys->content );
		my $cdn_ksk = $keys->{$cdn_name}->{ksk};
		foreach my $cdn_krecord (@$cdn_ksk) {
			my $cdn_kstatus = $cdn_krecord->{status};
			if ( $cdn_kstatus eq 'new' )
			{    #ignore anything other than the 'new' record
				my $exp_date = $cdn_krecord->{expirationDate};
				$k_expiry = strftime '%m/%d/%Y %H:%M:%S', gmtime $exp_date;
				$k_expiry .= " GMT";
				$ttl         = $cdn_krecord->{ttl};
				$algorithm   = $cdn_krecord->{dsRecord}->{algorithm};
				$digest_type = $cdn_krecord->{dsRecord}->{digestType};
				$digest      = $cdn_krecord->{dsRecord}->{digest};
			}
		}
	}

	#get active flag
	my $profile_id = $self->get_profile_id_by_cdn($cdn_name);
	my $active     = "false";
	my $cdn_rs = $self->db->resultset('Cdn')->search( { name => $cdn_name } )
		->single();
	my $dnssec_enabled = $cdn_rs->dnssec_enabled;
	if ($dnssec_enabled) {
		$active = "true";
	}

	#stash all the things
	$self->stash(
		msgs   => [],
		dnssec => {
			cdn_name       => $cdn_name,
			ttl            => $ttl,
			k_expiry       => $k_expiry,
			ds_algorithm   => $algorithm,
			ds_digest_type => $digest_type,
			ds_digest      => $digest,
			active         => $active
		},
		fbox_layout => 1
	);
}

sub add {
	my $self           = shift;
	my $cdn_name       = $self->param('cdn_name');
	my $k_expiry       = "365";
	my $z_expiry       = "30";
	my $effective_date = strftime( "%Y-%m-%d %H:%M:%S\n", gmtime(time) );
	my $keys;
	my $existing           = 0;
	my $response_container = $self->riak_get( "dnssec", $cdn_name );
	my $get_keys           = $response_container->{'response'};
	if ( $get_keys->is_success() ) {
		$existing = 1
			; ##change the generate keys dialog based on whether or not keys exist.
		$keys = decode_json( $get_keys->content );
		my $cdn_ksk = $keys->{$cdn_name}->{ksk};
		foreach my $cdn_krecord (@$cdn_ksk) {
			my $cdn_kstatus = $cdn_krecord->{status};
			if ( $cdn_kstatus eq 'new' )
			{    #ignore anything other than the 'new' record
				my $exp_date    = $cdn_krecord->{expirationDate};
				my $incept_date = $cdn_krecord->{inceptionDate};
				$k_expiry = ( $exp_date - $incept_date ) / 86400;
			}
		}
		my $cdn_zsk = $keys->{$cdn_name}->{zsk};
		foreach my $cdn_zrecord (@$cdn_zsk) {
			my $cdn_zstatus = $cdn_zrecord->{status};
			if ( $cdn_zstatus eq 'new' )
			{    #ignore anything other than the 'new' record
				my $exp_date    = $cdn_zrecord->{expirationDate};
				my $incept_date = $cdn_zrecord->{inceptionDate};
				$z_expiry = ( $exp_date - $incept_date ) / 86400;
			}
		}
	}
	&stash_role($self);
	$self->stash(
		msgs   => [],
		dnssec => {
			cdn_name       => $cdn_name,
			k_expiry       => $k_expiry,
			z_expiry       => $z_expiry,
			effective_date => $effective_date,
			existing       => $existing
		},
		fbox_layout => 1
	);
}

sub addksk {
	my $self     = shift;
	my $cdn_name = $self->param('cdn_name');
	my $k_expiry = "365";
	my $keys;
	my $effective_date = strftime( "%Y-%m-%d %H:%M:%S\n", gmtime(time) );
	my $response_container = $self->riak_get( "dnssec", $cdn_name );
	my $get_keys = $response_container->{'response'};
	if ( $get_keys->is_success() ) {
		$keys = decode_json( $get_keys->content );
		my $cdn_ksk = $keys->{$cdn_name}->{ksk};
		foreach my $cdn_krecord (@$cdn_ksk) {
			my $cdn_kstatus = $cdn_krecord->{status};
			if ( $cdn_kstatus eq 'new' )
			{    #ignore anything other than the 'new' record
				my $exp_date    = $cdn_krecord->{expirationDate};
				my $incept_date = $cdn_krecord->{inceptionDate};
				$k_expiry = ( $exp_date - $incept_date ) / 86400;
			}
		}
	}
	&stash_role($self);
	$self->stash(
		msgs   => [],
		dnssec => {
			cdn_name       => $cdn_name,
			k_expiry       => $k_expiry,
			z_expiry       => "1",               ##for is_valid purposes only.
			effective_date => $effective_date,
		},
		fbox_layout => 1
	);
}

sub activate {
	my $self     = shift;
	my $cdn_name = $self->param('dnssec.cdn_name');
	my $active   = $self->param('dnssec.active_flag');

	$active eq 'true'
		? $self->db->resultset('Cdn')->search( { name => $cdn_name } )
		->update( { dnssec_enabled => 1 } )
		: $self->db->resultset('Cdn')->search( { name => $cdn_name } )
		->update( { dnssec_enabled => 0 } );

	&stash_role($self);
	$self->stash(
		msgs        => [],
		dnssec      => {},
		fbox_layout => 1
	);
	$self->flash(
		message => "Active flag for $cdn_name was set to $active." );
	return $self->redirect_to("/cdns/$cdn_name/dnsseckeys/manage");
}

sub create {
	my $self           = shift;
	my $cdn_name       = $self->param('dnssec.cdn_name');
	my $z_expiry       = $self->param('dnssec.z_expiry');
	my $k_expiry       = $self->param('dnssec.k_expiry');
	my $ttl            = "60";
	my $effective_date = $self->param('dnssec.effective_date');
	if ( !defined($effective_date) ) {
		$effective_date = time();
	}
	else {
		$effective_date = str2time($effective_date);
	}
	&stash_role($self);

	if ( !&is_admin($self) ) {
		$self->flash( alertmsg => "Keys can only be generated by admins!" );
		return $self->redirect_to("/cdns/dnsseckeys/add");
	}

	if ( $self->is_valid() ) {

		#get profile_id for cdn
		my $profile_id = $self->get_profile_id_by_cdn($cdn_name);
		my %condition  = (
			'parameter.name' => 'tld.ttls.DNSKEY',
			'profile.name'   => $profile_id
		);
		my $rs_pp = $self->db->resultset('ProfileParameter')->search(
			\%condition,
			{   prefetch =>
					[ { 'parameter' => undef }, { 'profile' => undef } ]
			}
		)->single;
		if ($rs_pp) {
			$ttl = $rs_pp->parameter->value;
		}

		#create keys
		my $profile = $self->db->resultset('Profile')->search( { 'me.id' => $profile_id }, { prefetch => ['cdn'] } )->single();
		my $domain_name = $profile->cdn->domain_name;

		my $response_container = $self->riak_ping();
		my $ping_response      = $response_container->{response};
		if ( $ping_response->is_success ) {
			my $response_container
				= $self->generate_store_dnssec_keys( $cdn_name, $domain_name, $ttl,
				$k_expiry, $z_expiry, $effective_date );
			my $response = $response_container->{response};
			if ( $response->is_success ) {
				&log( $self, "Created dnssec keys for CDN $cdn_name",
					"UICHANGE" );
				$self->flash( message =>
						"Successfully created dnssec keys for: $cdn_name" );
			}
			else {
				$self->flash(
					{   "DNSSEC keys for $domain_name could not be created.  Response was"
							. $response->{_content}
					}
				);
			}
		}
		else {
			my @cdns = $self->db->resultset('Cdn')->search( {} )
				->get_column('name')->all();

			$self->stash(
				dnssec => {
					cdn_name => $cdn_name,
					ttl      => $ttl,
					k_expiry => $k_expiry,
					z_expiry => $z_expiry,
				},
				cdns        => \@cdns,
				fbox_layout => 1
			);
			my @msgs;
			push( @msgs, $ping_response->{_content} );
			$self->stash( msgs => \@msgs );
		}
		&stash_role($self);
		return $self->redirect_to("/cdns/$cdn_name/dnsseckeys/add");
	}
	else {
		&stash_role($self);
		$self->build_stash();
		$self->stash( msgs => [] );
		$self->render("dnssec_keys/add");
	}
}

sub genksk {
	my $self           = shift;
	my $cdn_name       = $self->param('dnssec.cdn_name');
	my $k_expiry       = $self->param('dnssec.k_expiry');
	my $ttl            = 60;
	my $effective_date = $self->param('dnssec.effective_date');
	if ( !defined($effective_date) ) {
		$effective_date = time();
	}
	else {
		$effective_date = str2time($effective_date);
	}
	&stash_role($self);

	if ( !&is_admin($self) ) {
		$self->flash( alertmsg => "Keys can only be generated by admins!" );
		return $self->redirect_to("/cdns/dnsseckeys/addksk");
	}

	if ( $self->is_valid() ) {

		# get profile_id for cdn
		my $profile_id = $self->get_profile_id_by_cdn($cdn_name);
		my %condition  = (
			'parameter.name' => 'tld.ttls.DNSKEY',
			'profile.name'   => $profile_id
		);
		my $rs_pp = $self->db->resultset('ProfileParameter')->search(
			\%condition,
			{   prefetch =>
					[ { 'parameter' => undef }, { 'profile' => undef } ]
			}
		)->single;
		if ($rs_pp) {
			$ttl = $rs_pp->parameter->value;
		}

		#get effective multiplier
		my $effective_multiplier;
		%condition = (
			'parameter.name' => 'DNSKEY.effective.multiplier',
			'profile.name'   => $profile_id
		);
		$rs_pp = $self->db->resultset('ProfileParameter')->search(
			\%condition,
			{   prefetch =>
					[ { 'parameter' => undef }, { 'profile' => undef } ]
			}
		)->single;
		if ($rs_pp) {
			$effective_multiplier = $rs_pp->parameter->value;
		}
		else {
			$effective_multiplier = '2';
		}
		my $inception = time();
		my $k_expiration = time() + ( 86400 * $k_expiry );
		my $keys;
		my $response_container = $self->riak_get( "dnssec", $cdn_name );
		my $get_keys = $response_container->{'response'};
		if ( $get_keys->is_success() ) {
			$keys = decode_json( $get_keys->content );
		}
		my $new_key
			= API::Cdn::regen_expired_keys( $self, "ksk", $cdn_name, $keys,
			$effective_date, 1, 1 );
		$keys->{$cdn_name} = $new_key;
		my $json_data = encode_json($keys);
		$response_container
			= $self->riak_put( "dnssec", $cdn_name, $json_data );
		my $response = $response_container->{"response"};
		if ( $response->is_success ) {
			&log( $self, "Generate KSK for CDN $cdn_name", "UICHANGE" );
			$self->flash(
				message => "Successfully generated KSK for: $cdn_name" );
		}
		else {
			my @msgs;
			push( @msgs, $response->{_content} );
			$self->stash( msgs => \@msgs );
			$self->flash(
				{   "KSK for $cdn_name could not be created.  Response was"
						. $response->{_content}
				}
			);
		}
		&stash_role($self);
		return $self->redirect_to("/cdns/$cdn_name/dnsseckeys/addksk");
	}
	else {
		&stash_role($self);
		$self->build_stash();
		$self->stash( msgs => [] );
		$self->render("dnssec_keys/addksk");
	}
}

sub build_stash {
	my $self     = shift;
	my $cdn_name = $self->param('dnssec.cdn_name');
	my $ttl      = $self->param('dnssec.ttl');
	my $z_expiry = $self->param('dnssec.z_expiry');
	my $k_expiry = $self->param('dnssec.k_expiry');
	my @cdns = $self->db->resultset('Cdn')->search( {} )->get_column('name')
		->all();
	&navbarpage($self);
	$self->stash(
		dnssec => {
			cdn_name => $cdn_name,
			ttl      => $ttl,
			k_expiry => $k_expiry,
			z_expiry => $z_expiry,
		},
		cdns        => \@cdns,
		fbox_layout => 1
	);
}

sub is_valid {
	my $self           = shift;
	my $cdn_name       = $self->param('dnssec.cdn_name');
	my $z_expiry       = $self->param('dnssec.z_expiry');
	my $k_expiry       = $self->param('dnssec.k_expiry');
	my $effective_date = $self->param('dnssec.effective_date');

	if ( $cdn_name eq "default" ) {
		$self->field('dnssec.cdn_name')
			->is_equal( "", "Please choose a CDN" );
	}
	if ( $z_expiry eq "" || !looks_like_number($z_expiry) || $z_expiry < 1 ) {
		$self->field('dnssec.z_expiry')
			->is_equal( "", "$z_expiry is not a number greater than 0" );
	}
	if ( $k_expiry eq "" || !looks_like_number($k_expiry) || $k_expiry < 1 ) {
		$self->field('dnssec.k_expiry')
			->is_equal( "", "$k_expiry is not a number greater than 0" );
	}

	if ($effective_date) {
		$self->field('dnssec.effective_date')->is_like(
			qr/^((((19|[2-9]\d)\d{2})[\/\.-](0[13578]|1[02])[\/\.-](0[1-9]|[12]\d|3[01])\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((19|[2-9]\d)\d{2})[\/\.-](0[13456789]|1[012])[\/\.-](0[1-9]|[12]\d|30)\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((19|[2-9]\d)\d{2})[\/\.-](02)[\/\.-](0[1-9]|1\d|2[0-8])\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((1[6-9]|[2-9]\d)(0[48]|[2468][048]|[13579][26])|((16|[2468][048]|[3579][26])00))[\/\.-](02)[\/\.-](29)\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])))$/,
			"Effective Date is not a valid dateTime!  Should be in the format of YYYY-MM-DD HH:MM:SS"
		);
	}

	return $self->valid;
}

1;
