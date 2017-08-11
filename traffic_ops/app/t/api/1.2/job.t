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
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

# get first one sorted by start_time DESC
$t->get_ok("/api/1.2/jobs")->status_is(200)->json_is( "/response/0/id", 300 )
    ->json_is( "/response/0/assetUrl", "http://cdn2.edge/job3/.*" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

# get first one sorted by start_time and filtered by DS
$t->get_ok("/api/1.2/jobs?dsId=200")->status_is(200)->json_is( "/response/0/id", 200 )
    ->json_is( "/response/0/assetUrl", "http://cdn2.edge/job2/.*" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

# get first one sorted by start_time and filtered by User
$t->get_ok("/api/1.2/jobs?userId=200")->status_is(200)->json_is( "/response/0/id", 300 )
    ->json_is( "/response/0/assetUrl", "http://cdn2.edge/job3/.*" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

# get specific job
$t->get_ok("/api/1.2/jobs/100")->status_is(200)->json_is( "/response/0/id", 100 )
    ->json_is( "/response/0/keyword", "PURGE" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
done_testing();

