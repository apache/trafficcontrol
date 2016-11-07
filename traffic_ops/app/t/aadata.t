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

BEGIN { $ENV{MOJO_MODE} = "test" }
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#login
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok('/aadata/Server')->status_is(200)->json_has('/aaData')->json_has('atlanta-mid-01');

$t->get_ok('/aadata/ProfileParameter')->status_is(200)->json_has('/aaData')->json_has('domain_name');

$t->get_ok('/aadata/ServerSelect')->status_is(200)->json_has('/aaData')->json_has('atlanta-edge-01');

$t->get_ok('/aadata/Hwinfo')->status_is(200)->json_has('/data')->json_has('atlanta-edge-01.ga.atlanta.kabletown.net');

$t->get_ok('/aadata/Log')->status_is(200)->json_has('/aaData');

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

done_testing();

sub get_server_name {

}
