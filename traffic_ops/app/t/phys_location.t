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
use DBI;
use strict;
use warnings;
use Test::TestHelper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $t   = Test::Mojo->new('TrafficOps');
my $dbh = Schema->database_handle;
my $schema = Schema->connect_to_database;

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => 'portal', p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $q      = 'select * from phys_location limit 1';
my $get_ds = $dbh->prepare($q);
$get_ds->execute();
my $p = $get_ds->fetchall_arrayref( {} );
$get_ds->finish();

$t->get_ok('/dataphys_location')->status_is(200)->json_has('/0/name')->json_has('/0/short_name')->json_has('/0/address')->json_has('/0/city')
	->json_has('/0/zip')->json_has('/0/phone')->json_has('/0/state')->json_has('/0/email')->or( sub { diag $t->tx->res->content->asset->{content}; } );

# $t->get_ok('/datalocation')->status_is(200)->json_has('/0/type_name')->json_has('/0/longitude')->json_has('/0/short_name')
# ->json_has('/0/parent_location_id')->json_has('/0/name')->json_has('/0/type_id')->json_has('/0/latitude')->json_has('/0/parent_location_name');

####################### RW testing - careful with these! #####################################################

#clean up old crud
&upd_and_del();

# create a new loc
$t->post_ok(
	'/phys_location/create' => form => {
		'location.name'       => 'jlp-test-location',
		'location.short_name' => 'jlp',
		'location.address'    => '4100 East Dry Creek Rd.',
		'location.city'       => 'Centennial',
		'location.zip'        => '80184',
		'location.state'      => 'CO',
		'location.phone'      => '',
		'location.poc'        => 'Jan van Doorn',
		'location.email'      => 'jvd@comcast.com',
		'location.comments'   => 'boo',
		'location.region'     => '100',
	}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# test strange email address
$t->post_ok(
	'/phys_location/create' => form => {
		'location.name'       => 'try-bad-email',
		'location.short_name' => 'jlp',
		'location.address'    => '4100 East Dry Creek Rd.',
		'location.city'       => 'Centennial',
		'location.zip'        => '80184',
		'location.state'      => 'CO',
		'location.phone'      => '',
		'location.poc'        => 'Bubba',
		'location.email'      => 'Louie was here',
		'location.comments'   => 'boo',
		'location.region'     => '100',
	}
)->status_is(200)->message( 'invalid email' );

# modify and delete it
&upd_and_del( );

sub upd_and_del() {
	my %overrides = @_;
	my $q      = 'select id from phys_location where name = \'jlp-test-location\'';
	my $get_ds = $dbh->prepare($q);
	$get_ds->execute();
	my $p = $get_ds->fetchall_arrayref( {} );
	$get_ds->finish();
	my $i = 0;
	while ( defined( $p->[$i] ) ) {
		my $id = $p->[$i]->{id};
		$t->post_ok(
			"/phys_location/$id/update" => form => {
				'location.name'       => 'jlp-test-location',
				'location.short_name' => 'jlp',
				'location.address'    => '4100 East Dry Creek Rd.',
				'location.city'       => 'Centennial',
				'location.zip'        => '80184',
				'location.state'      => 'CO',
				'location.poc'        => 'Bubba',
				'location.phone'      => '800-334-5545',
				'location.email'      => 'jvd@comcast.com',
				'location.comments'   => 'boo',
				'location.region'     => '100',
				}
		)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
		$t->get_ok( "/phys_location/$id/delete" )->status_is(302);
		$i++;
	}
}
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
