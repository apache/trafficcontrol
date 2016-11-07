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
use Schema;
use Fixtures::TmUser;
use Test::TestHelper;
use Data::Dumper;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#login
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#create dsserver
ok $t->post_ok(
	'/create/dsserver',
	=> form => {
		server          => '1',
		deliveryservice => '2'
	}
)->status_is(302), "create deliveryservice_server";

# validate ds_server was created
ok $t->get_ok('/datadeliveryserviceserver')->status_is(200)->json_is( '/1/deliveryservice' => 'test-ds1' )->json_is( '/1/server' => '1' ),
	"validate deliveryservice_server was added";

# validate edit route
ok $t->get_ok('/dss/1/edit')->status_is(200), "validate edit screen";

#assign_servers
ok $t->post_ok(
	'/dss/1/update' => form => {
		'id'         => '1',
		'serverid_2' => 'on',
		'serverid_1' => 'off'
	}
)->status_is(302), "assign server to ds";

#clone_server
ok $t->post_ok(
	'/update/cpdss/2' => form => {
		'from_server' => '1',
		'to_server'   => '2',
	}
)->status_is(302), "clone server";

#validate clone
ok $t->get_ok('/datadeliveryserviceserver')->status_is(200)->json_is( '/1/deliveryservice' => 'steering-ds2' )->json_is( '/1/server' => '2' ),
	"validate deliveryservice was cloned";

#validate cp dss view
ok $t->get_ok('/cpdssiframe/view/1')->status_is(200), "cpdss iframe";

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
