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
use Test::TestHelper;
use Fixtures::Asn;
use Schema;
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#login
ok $t->post_ok( '/api/1.1/user/login', json => { u => 'admin', p => 'password' } )->status_is(200), "Login as admin";

=cut
#table view
ok $t->get_ok('/asns')->status_is(200), "Get table view";

#create
ok $t->post_ok(
	'/asn/create' => form => {
		asn        => '01234',
		cachegroup => '1'
	}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } ), "Create asn";

#get asn id from db
my $q      = 'select id from asn where asn = \'01234\'';
my $get_ds = $dbh->prepare($q);
$get_ds->execute();
my $p = $get_ds->fetchall_arrayref( {} );
$get_ds->finish();
my $asn    = $p->[0];
my $asn_id = $asn->{id};
ok( $asn_id, "does the asn exist?" );

#udpate asn
ok $t->post_ok(
	      '/asns/'
		. $asn_id
		. '/update' => form => {
		asn => 1235,
		}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } ), "update asn with id of $asn_id";
ok my $update_asn = $schema->resultset('Asn')->find( { asn => '1235' } ), "Validate the asn was updated";

#delete asn
ok $t->get_ok("/asns/$asn_id/delete")->status_is(302), "Delete Asn";
my $delete_asn = $schema->resultset('Asn')->find( { asn => '1235' } );
ok !$delete_asn, "Validate the asn was deleted";

=cut

#finish up
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
