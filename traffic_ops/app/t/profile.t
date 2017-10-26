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
use Schema;
use strict;
use warnings;
use Test::TestHelper;
no warnings 'once';
use warnings 'all';

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

# the jsons
# Note the 3 is the index in the array returned, not the id.  It's safe to assume there are at least 3 profiles.
$t->get_ok('/dataprofile')->status_is(200)->json_has('/0/name')->json_has('/0/description');

####################### RW testing - careful with these! #####################################################

#clean up old crud
&upd_and_del();

# create a new param
$t->post_ok(
	'/profile/create' => form => {
		'profile.name'        => 'JLP_Test',
		'profile.description' => 'JLP Test Host',
		'profile.cdn'         => 100,
		'profile.type'        => 'ATS_PROFILE',
		'profile.routing_disabled' => 0
	}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# modify and delete it
&upd_and_del();

sub upd_and_del() {
	my $q      = 'select id from profile where name = \'JLP_Test\'';
	my $get_ds = $dbh->prepare($q);
	$get_ds->execute();
	my $p = $get_ds->fetchall_arrayref( {} );
	$get_ds->finish();
	my $i = 0;
	while ( defined( $p->[$i] ) ) {
		my $id = $p->[$i]->{id};
		$t->post_ok(
			"/profile/$id/update" => form => {
				'profile.name'        => 'JLP_Test',
				'profile.description' => 'JLP Test Host Updated',
				'profile.cdn'         => 100,
				'profile.type'        => 'ATS_PROFILE',
				'profile.routing_disabled' => 0
			}
		)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
		$t->get_ok("/profile/$id/delete")->status_is(302);
		$i++;
	}
}
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
