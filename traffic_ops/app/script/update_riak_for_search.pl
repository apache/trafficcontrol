#!/usr/bin/perl
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

use strict;
use warnings;
use Data::Dumper;
use Getopt::Long;
use JSON;
use LWP::UserAgent;
use MIME::Base64;

my $to_url     = 'http://localhost:3000';
my $to_un = "admin";
my $to_pw = "password";

GetOptions( "to_url=s" => \$to_url, "to_un=s" => \$to_un, "to_pw=s" => \$to_pw );

my $ua = LWP::UserAgent->new;
$ua->timeout(30);
$ua->ssl_opts(verify_hostname => 0);

my $cookie = &get_cookie($ua, $to_url, $to_un, $to_pw);
$ua->default_header('Cookie' => $cookie);
my $dss = &get_deliveryservices($to_url, $ua);

foreach my $ds (@$dss) {
	if ($ds->{protocol} > 0) {
		my $xml_id = $ds->{xmlId};
		my $cdn = $ds->{cdnName};
		print "Updating record for: $xml_id\n";
		my $record = &get_riak_record($xml_id, $to_url, $ua);
		if (!defined($record)) {
			next;
		}
		$record->{deliveryservice} = $xml_id;
		$record->{cdn} = $cdn;
		$record->{certificate}->{crt} = decode_base64($record->{certificate}->{crt});
		$record->{certificate}->{csr} = decode_base64($record->{certificate}->{csr});
		$record->{certificate}->{key} = decode_base64($record->{certificate}->{key});
		if (!defined($record->{hostname})) {  #add the hostname if it's not there
			my $hostname = $ds->{exampleURLs}[0];
			$hostname =~ /(https?:\/\/)(.*)/;
			$record->{hostname} = $2;
		}

		if (!defined($record->{key})) {
			$record->{key} = $xml_id;
		}

		#send back to riak
		&send_riak_record($record, $to_url, $ua);
	}

}

sub get_cookie {
	my $ua = shift;
	my $to_host = shift;
	my $u = shift;
	my $p = shift;

	my $url = $to_host . "/api/2.0/user/login";
	my $json = encode_json({u => $u, p => $p});
	my $req = HTTP::Request->new( 'POST', $url );
	$req->header( 'Content-Type' => 'application/json' );
	$req->content( $json );
	my $response = $ua->request( $req );

	if(!$response->is_success() || $response->code() > 400) {
		print "Could not login to traffic_ops!  Response was ". $response->{_rc} . " - " . $response->{_msg} . "\n";
		exit 1;
	}

	my $cookie;
	if ( $response->header('Set-Cookie') ) {
		($cookie) = split(/\;/, $response->header('Set-Cookie'));
	}

	if ( $cookie =~ m/mojolicious/ ) {
		return $cookie;
	}
	else {
		print "FATAL mojolicious cookie not found from Traffic Ops!\n";
		exit 1;
	}
}

sub get_deliveryservices {
	my $to_host = shift;
	my $ua = shift;

	my $url = $to_host . "/api/2.0/deliveryservices";
	my $response = $ua->get( $url );

	if(!$response->is_success() || $response->code() > 400) {
		print "Could not get deliveryservices!  Response was ". $response->{_rc} . " - " . $response->{_msg} . "\n";
		exit 1;
	}

	my $content = decode_json($response->{_content});
	return $content->{"response"};
}

sub get_riak_record {
	my $xml_id = shift;
	my $to_url = shift;
	my $ua = shift;
	my $url = $to_url . "/api/2.0/deliveryservices/xmlId/$xml_id/sslkeys";
	my $response = $ua->get( $url );

	if(!$response->is_success() || $response->code() > 400) {
		print "Could not get ssl record for $xml_id from riak!  Response was ". $response->{_rc} . " - " . $response->{_msg} . "Skipping...\n";
		return;
	}

	my $content = decode_json($response->{_content});
	return $content->{"response"};
}

sub send_riak_record {
	my $record = shift;
	my $to_url = shift;
	my $ua = shift;

	my $url = $to_url . "/api/2.0/deliveryservices/sslkeys/add";
	my $req = HTTP::Request->new( 'POST', $url );
	$req->header( 'Content-Type' => 'application/json' );
	$req->content( encode_json($record) );
	my $response = $ua->request( $req );

	if(!$response->is_success()) {
		my $key = $record->{key};
		print "Could not send riak record for $key!  Response was ". $response->{_rc} . " - " . $response->{_msg} . "\n";
	}

}

