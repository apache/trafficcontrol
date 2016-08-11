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
use Data::Dumper;
use DBI;
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "integration" }
my $t = Test::Mojo->new('TrafficOps');
no warnings 'once';
use warnings 'all';
my $api_version = '1.1';

my $token = '91504CE6-8E4A-46B2-9F9F-FE7C15228498';
my $path  = '/api/1.1/user/login/token';
$t->post_ok( $path => json => { t => $token } )->status_is(200)->json_is( '/alerts/0/text' => "Successfully logged in." )
	->json_is( '/alerts/0/level' => "success" );

my $json = JSON->new->allow_nonref;

my @etypes = ( "CHECK_EXTENSION_BOOL", "CHECK_EXTENSION_NUM" );
foreach my $num ( 1 .. 36 ) {
	$t->get_ok('/api/1.1/to_extensions.json')->status_is(200);
	my $extlist = $json->decode( $t->tx->res->content->asset->{content} );

	if ( scalar( @{ $extlist->{response} } ) < 31 ) {
		$t->post_ok(
			'/api/1.1/to_extensions' => json => {
				type                   => $etypes[ $num % 2 ],
				name                   => "X" . $num . "_TESTiNG",
				servercheck_short_name => "X" . $num,
				additional_config_json => "{ \"select\": \"ilo_ip_address\", \"cron\": \"9 * * * *\" }",
				version                => "1.0.0",
				isactive               => "1",
				description            => "description",
				script_file            => "ping",
				info_url               => "http://foo.com/bar.html"
			}
			)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( '/alerts/0/text' => "Check Extension Loaded." )
			->json_is( '/alerts/0/level' => "success" );
	}
	else {
		$t->post_ok(
			'/api/1.1/to_extensions' => json => {
				type                   => "CHECK_EXTENSION_BOOL",
				name                   => "X" . $num . "_TESTiNG",
				servercheck_short_name => "X" . $num,
				additional_config_json =>
					"{ \"path\": \"/api/1.1/servers.json\",  \"match\": { \"type\": \"EDGE\"}, \"select\": \"ilo_ip_address\", \"cron\": \"9 * * * *\" }",
				version     => "1.0.0",
				isactive    => "1",
				description => "description",
				script_file => "ping",
				info_url    => "http://foo.com/bar.html"
			}
			)->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
			->json_is( '/alerts/0/text' => "error No open slots left for checks, delete one first." )->json_is( '/alerts/0/level' => "error" );
	}
}

# check to see if the server checks page looks good.
$t->get_ok('/server_check')->status_is(200)->text_is( 'th#col1' => 'Hostname' )->text_is( 'th#col2' => 'Profile' )->text_is( 'th#col3' => 'ADMIN' )
	->text_is( 'th#col4' => 'UPD' )->text_is( 'th#col5' => 'ILO' )->text_is( 'th#col6' => '10G' )->text_is( 'th#col7' => 'FQDN' )
	->text_is( 'th#col8' => 'DSCP' )->text_is( 'th#col9' => 'X1' )->text_is( 'th#col10' => 'X2' )->text_is( 'th#col11' => '10G6' )
	->text_is( 'th#col12' => 'X3' )->text_is( 'th#col13' => 'STAT' )->text_is( 'th#col14' => 'X4' )->text_is( 'th#col15' => 'MTU' )
	->text_is( 'th#col16' => 'TRTR' )->text_is( 'th#col17' => 'TRMO' )->text_is( 'th#col18' => 'CHR' )->text_is( 'th#col19' => 'CDU' )
	->text_is( 'th#col20' => 'ORT' )->text_is( 'th#col21' => 'X5' )->text_is( 'th#col22' => 'X6' )->text_is( 'th#col23' => 'X7' )
	->text_is( 'th#col24' => 'X8' )->text_is( 'th#col25' => 'X9' )->text_is( 'th#col26' => 'X10' )->text_is( 'th#col27' => 'X11' )
	->text_is( 'th#col28' => 'X12' )->text_is( 'th#col29' => 'X13' )->text_is( 'th#col30' => 'X14' );

# post stome status
my $test_server_id = "23";
for my $test (qw/ILO 10G FQDN DSCP 10G6 STAT MTU TRTR TRMO X1 X4 X6 X9/) {
	$t->post_ok( '/api/1.1/servercheck' => json => { id => $test_server_id, servercheck_short_name => $test, value => 1 } )->status_is(200)
		->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( '/alerts/0/text' => "Server Check was successfully updated." )
		->json_is( '/alerts/0/level' => "success" );
}

# clean up and test "delete"
$t->get_ok('/api/1.1/to_extensions.json')->status_is(200);
my $extlist = $json->decode( $t->tx->res->content->asset->{content} );
foreach my $ext ( @{ $extlist->{response} } ) {
	if ( $ext->{name} =~ /_TESTiNG$/ ) {
		$t->post_ok( '/api/1.1/to_extensions/' . $ext->{id} . '/delete' )->status_is(200)->json_is( '/alerts/0/text' => "Extension deleted." )
			->json_is( '/alerts/0/level' => "success" );
	}
}
done_testing();
