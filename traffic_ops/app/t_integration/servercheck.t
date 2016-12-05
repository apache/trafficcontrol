package main;
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
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use Mojo::Util qw/squish/;
use DBI;
use JSON;
use Data::Dumper;
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "integration" }
my $t = Test::Mojo->new('TrafficOps');
no warnings 'once';
use warnings 'all';

my $api_version = '1.1';
$t->post_ok( '/login', => form => { u => 'admin', p => 'password' } )->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $json = JSON->new->allow_nonref;
$t->get_ok( '/api/' . $api_version . '/servers.json' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $servers = $json->decode( $t->tx->res->content->asset->{content} );

my $test_server    = "atsec-hou-05";
my $test_server_id = "33";
my %type_done      = ();

# The web - page. Check to see it has all the headers
$t->get_ok('/server_check')->status_is(200)->text_is( 'th#col1' => 'Hostname' )->text_is( 'th#col2' => 'Profile' )->text_is( 'th#col3' => 'ADMIN' )
	->text_is( 'th#col4' => 'UPD' )->text_is( 'th#col5' => 'ILO' )->text_is( 'th#col6' => '10G' )->text_is( 'th#col7' => 'FQDN' )
	->text_is( 'th#col8' => 'DSCP' )->text_is( 'th#col9' => '' )->text_is( 'th#col10' => '' )->text_is( 'th#col11' => '10G6' )->text_is( 'th#col12' => '' )
	->text_is( 'th#col13' => 'STAT' )->text_is( 'th#col14' => '' )->text_is( 'th#col15' => 'MTU' )->text_is( 'th#col16' => 'TRTR' )
	->text_is( 'th#col17' => 'TRMO' )->text_is( 'th#col18' => 'CHR' )->text_is( 'th#col19' => 'CDU' )->text_is( 'th#col20' => 'ORT' )
	->text_is( 'th#col21' => '' )->text_is( 'th#col22' => '' );

# The json that populates it... Everything should be default / NULL
$t->get_ok('/api/1.1/servercheck/aadata.json')->status_is(200);
my $jdata = $json->decode( $t->tx->res->content->asset->{content} );
foreach my $line ( @{ $jdata->{aaData} } ) {
	if ( $line->[1] eq $test_server ) {
		ok defined( $line->[0] ) && defined( $line->[1] ) && defined( $line->[2] ) && defined( $line->[3] ), "Are the madatory fields defined?";
		my $i = 5;
		ok $line->[3] eq "REPORTED", "Is " . $test_server . " status REPORTED? found " . $line->[3];
		ok $line->[4] == 0, "Is " . $test_server . " upd_pending 0? found " . $line->[4];
		while ( $i < scalar( @{$line} ) ) {
			my $check = $line->[ $i++ ];
			next unless defined($check);
			ok $check == 0, "Are all other fields 0? Found " . $check . " for index " . $i;
		}
	}
}

$t->post_ok( '/server/' . $test_server . '/status/OFFLINE' => form => {} )->json_is( '/result' => 'SUCCESS' );
$t->get_ok( '/dataserverdetail/select/' . $test_server )->status_is(200)->json_is( '/0/status' => 'OFFLINE' );
$t->post_ok( '/server/' . $test_server . '/status/REPORTED' => form => {} )->json_is( '/result' => 'SUCCESS' );    # for next test
$t->get_ok( '/dataserverdetail/select/' . $test_server )->status_is(200)->json_is( '/0/status' => 'REPORTED' );

# login with the extension user token, and set some values for $test_server.
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );              # logout as admin first
my $token = '91504CE6-8E4A-46B2-9F9F-FE7C15228498';
my $path  = '/api/1.1/user/login/token';
$t->post_ok( $path => json => { t => $token } )->status_is(200)->json_is( '/alerts/0/text' => "Successfully logged in." )
	->json_is( '/alerts/0/level' => "success" );

# set some 1's
for my $test (qw/ILO 10G FQDN DSCP 10G6 STAT MTU TRTR TRMO/) {
	$t->post_ok( '/api/1.1/servercheck' => json => { id => $test_server_id, servercheck_short_name => $test, value => 1 } )->status_is(200)
		->json_is( '/alerts/0/text' => "Server Check was successfully updated." )->json_is( '/alerts/0/level' => "success" );
}

# and some ints
for my $test (qw/CHR CDU ORT/) {
	$t->post_ok( '/api/1.1/servercheck' => json => { id => $test_server_id, servercheck_short_name => $test, value => 99 } )->status_is(200)
		->json_is( '/alerts/0/text' => "Server Check was successfully updated." )->json_is( '/alerts/0/level' => "success" );
}

$t->get_ok('/api/1.1/servercheck/aadata.json')->status_is(200);
$jdata = $json->decode( $t->tx->res->content->asset->{content} );
foreach my $line ( @{ $jdata->{aaData} } ) {
	if ( $line->[1] eq $test_server ) {
		ok defined( $line->[0] ) && defined( $line->[1] ) && defined( $line->[2] ) && defined( $line->[3] ), "Are the madatory fields defined?";
		my $i = 5;
		ok $line->[3] eq "REPORTED", "Is " . $test_server . " status REPORTED? found " . $line->[3];
		ok $line->[4] == 0, "Is " . $test_server . " upd_pending 0? found " . $line->[4];
		while ( $i < scalar( @{$line} ) ) {

			# diag Dumper($line);
			my $check = $line->[ $i++ ];
			if ( $i == 10 || $i == 11 || $i == 13 || $i == 15 || $i > 21 ) {

				# ok !defined($check) || $check != 0, "Are all the right fields undef? Found not true for index " . $i;
			}
			elsif ( $i >= 19 && $i <= 21 ) {
				ok $check == 99, "Are all the right fields 99? Found " . $check . " for index " . $i;
			}
			else {
				ok $check == 1, "Are all other fields 1? Found " . $check . " for index " . $i;
			}
		}
	}
}

# set some 0's
for my $test (qw/ILO 10G FQDN DSCP 10G6 STAT MTU TRTR TRMO CHR CDU ORT/) {
	$t->post_ok( '/api/1.1/servercheck' => json => { id => $test_server_id, servercheck_short_name => $test, value => 0 } )->status_is(200)
		->json_is( '/alerts/0/text' => "Server Check was successfully updated." )->json_is( '/alerts/0/level' => "success" );
}

$t->get_ok('/api/1.1/servercheck/aadata.json')->status_is(200);
$jdata = $json->decode( $t->tx->res->content->asset->{content} );
foreach my $line ( @{ $jdata->{aaData} } ) {
	if ( $line->[1] eq $test_server ) {
		ok defined( $line->[0] ) && defined( $line->[1] ) && defined( $line->[2] ) && defined( $line->[3] ), "Are the madatory fields defined?";
		my $i = 5;
		ok $line->[3] eq "REPORTED", "Is " . $test_server . " status REPORTED? found " . $line->[3];
		ok $line->[4] == 0, "Is " . $test_server . " upd_pending 0? found " . $line->[4];
		while ( $i < scalar( @{$line} ) ) {
			my $check = $line->[ $i++ ];
			next unless defined($check);
			ok $check == 0, "Are all other fields 0? Found " . $check . " for index " . $i;
		}
	}
}

done_testing();
