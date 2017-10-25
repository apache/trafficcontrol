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

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok('/api/1.2/caches/stats')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_has( '/response', 'has a response' )
    ->json_has( '/response/0/profile', 'has a profile key' )
    ->json_has( '/response/0/cachegroup', 'has a cachegroup key' )
    ->json_has( '/response/0/hostname', 'has a hostname key' )
    ->json_has( '/response/0/ip', 'has a ip key' )
    ->json_has( '/response/0/status', 'has a status key' )
    ->json_has( '/response/0/healthy', 'has a healthy key' )
    ->json_has( '/response/0/connections', 'has a connections key' )
    ->json_has( '/response/0/kbps', 'has a kbps key' );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
