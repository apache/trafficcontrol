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
		server          => '100',
		deliveryservice => '200'
	}
)->status_is(302), "create deliveryservice_server";

# validate ds_server was created
ok $t->get_ok('/datadeliveryserviceserver')->status_is(200)->json_is( '/1/deliveryservice' => 'test-ds1' )->json_is( '/1/server' => '300' ),
	"validate deliveryservice_server was added";

# validate edit route
ok $t->get_ok('/dss/200/edit')->status_is(200), "validate edit screen";

#assign_servers
ok $t->post_ok(
	'/dss/100/update' => form => {
		'serverid_200' => 'on',
		'serverid_100' => 'off'
	}
)->status_is(302), "assign server to ds";

#clone_server
ok $t->post_ok(
	'/update/cpdss/200' => form => {
		'from_server' => '100',
		'to_server'   => '200',
	}
)->status_is(302), "clone server";

#validate clone
ok $t->get_ok('/datadeliveryserviceserver')->status_is(200)->json_is( '/8/deliveryservice' => 'steering-ds1' )->json_is( '/8/server' => '900' ),
	"validate deliveryservice was cloned";

#validate cp dss view
ok $t->get_ok('/cpdssiframe/view/100')->status_is(200), "cpdss iframe";

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
