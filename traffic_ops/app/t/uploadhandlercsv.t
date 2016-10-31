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
use Test::TestHelper;
use Schema;
use strict;
use warnings;

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

my $q           = 'select * from server where type = 2 limit 1';
my $get_servers = $dbh->prepare($q);
$get_servers->execute();
my $svr = $get_servers->fetchall_arrayref( {} );
$get_servers->finish();
my $test_server_id = $svr->[0]->{id};

# the jsons associated with server
$t->get_ok( '/server_by_id/' . $test_server_id )->status_is(200)->json_is( '/host_name', $svr->[0]->{host_name} )
	->json_is( '/domain_name', $svr->[0]->{domain_name} )->json_is( '/tcp_port', $svr->[0]->{tcp_port} )
	->json_is( '/interface_name', $svr->[0]->{interface_name} )->json_is( '/ip_address', $svr->[0]->{ip_address} )
	->json_is( '/ip_netmask', $svr->[0]->{ip_netmask} )->json_is( '/ip_gateway', $svr->[0]->{ip_gateway} )
	->json_is( '/interface_mtu', $svr->[0]->{interface_mtu} );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
