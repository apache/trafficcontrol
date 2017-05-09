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
#    ->json_is( "/response/hostname" => "http://10.75.168.91")
#    ->json_is( "/response/ds_assigned/1" => "ds1")
#            , 'Does the deliveryservice details return?';
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
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

ok $t->post_ok('/api/1.2/cachegroups/300/deliveryservices' => {Accept => 'application/json'} => json => {
        "deliveryServices" => [ 100 ]})
     ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
     ->json_is( "/response/id" => 300 )
     ->json_is( "/response/deliveryServices/0" => 100 )
     ->json_is( "/alerts/0/level" => "success" )
     ->json_is( "/alerts/0/text" => "Delivery services successfully assigned to all the servers of cache group 300" )
            , 'Does the delivery services assign details return?';

ok $t->get_ok('/api/1.2/deliveryserviceserver')
     ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
     ->json_is( "/response/0/deliveryService" => "100" )
     ->json_is( "/response/0/server" => 100 )
     ->json_is( "/response/1/deliveryService" => 100 )
     ->json_is( "/response/1/server" => 300 )
     ->json_is( "/response/2/deliveryService" => 100 )
     ->json_is( "/response/2/server" => 600 )
            , 'Does the delivery services servers details return?';

ok $t->delete_ok('/api/1.2/deliveryservice_server/100/100')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
ok $t->delete_ok('/api/1.2/deliveryservice_server/100/100')->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
