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
no warnings 'once';
use warnings 'all';
use Test::TestHelper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $q      = 'select * from status limit 1';
my $get_ds = $dbh->prepare($q);
$get_ds->execute();
my $p = $get_ds->fetchall_arrayref( {} );
$get_ds->finish();

# the jsons
# Note the 3 is the index in the array returned, not the id.  It's safe to assume there are at least 3 statuses.
$t->get_ok('/datastatus')->status_is(200)->json_has('/3/name')->json_has('/3/description');

#clean up old crud
&upd_and_del();

# create a new param
$t->post_ok(
	'/status/create' => form => {
		name        => 'JLP_TEST_STATUS',
		description => 'This is a test status'

	}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# modify and delete it
&upd_and_del();

sub upd_and_del() {
	my $q      = 'select id from status where name = \'JLP_TEST_STATUS\'';
	my $get_ds = $dbh->prepare($q);
	$get_ds->execute();
	my $p = $get_ds->fetchall_arrayref( {} );
	$get_ds->finish();
	my $i = 0;
	while ( defined( $p->[$i] ) ) {
		my $id = $p->[$i]->{id};
		$t->post_ok(
			'/status/update/'
				. $id => form => {
				name        => 'JLP_TEST_STATUS',
				description => 'This is a test status'
				}
		)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
		$t->get_ok( '/status/delete/' . $id )->status_is(302);
		$i++;
	}
}
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
