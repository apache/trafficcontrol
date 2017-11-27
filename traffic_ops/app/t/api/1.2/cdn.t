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
	->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

$t->get_ok("/api/1.2/cdns")->status_is(200)->json_is( "/response/0/id", 100 )
    ->json_is( "/response/0/name", "cdn1" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/cdns/100")->status_is(200)->json_is( "/response/0/id", 100 )
    ->json_is( "/response/0/name", "cdn1" )->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/cdns/100/snapshot')
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    , 'Is the cdn snapshotted?';

$t->get_ok("/api/1.2/cdns/cdn1/snapshot")->status_is(200)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/cdns/cdn1/snapshot/new")->status_is(200)
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/cdns/100/queue_update' => {Accept => 'application/json'} => json => {
            "action" => "queue" })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/cdnId" => 100 )
        ->json_is( "/response/action" => "queue" )
    , 'Does the cdn details return?';

$t->get_ok("/api/1.2/servers?cdnId=100")->status_is(200)
    ->json_is( "/response/0/updPending", 1 )
    ->json_is( "/response/1/updPending", 1 )
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/cdns/100/queue_update' => {Accept => 'application/json'} => json => {
            "action" => "dequeue" })
        ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/response/cdnId" => 100 )
        ->json_is( "/response/action" => "dequeue" )
    , 'Does the cdn details return?';

$t->get_ok("/api/1.2/servers?cdnId=100")->status_is(200)
    ->json_is( "/response/0/updPending", 0 )
    ->json_is( "/response/1/updPending", 0 )
    ->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.2/cdns' => {Accept => 'application/json'} => json => {
        "name" => "cdn_test", "dnssecEnabled" => \1, "domainName" => "testcdn.kabletown.net" })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/name" => "cdn_test" )
    ->json_is( "/alerts/0/level" => "success" )
    ->json_is( "/alerts/0/text" => "cdn was created." )
            , 'Do the cdn queue update details return?';

my $cdn_id = &get_cdn_id('cdn_test');

ok $t->put_ok('/api/1.2/cdns/' . $cdn_id  => {Accept => 'application/json'} => json => {
        "name" => "cdn_test2", "domainName" => "testcdn2.kabletown.net", "dnssecEnabled" => \1
        })
    ->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
    ->json_is( "/response/name" => "cdn_test2" )
    ->json_is( "/alerts/0/level" => "success" )
            , 'Does the cdn details return?';

ok $t->delete_ok('/api/1.2/cdns/' . $cdn_id)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->delete_ok('/api/1.2/cdns/name/cdn3')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->put_ok('/api/1.2/cdns/' . $cdn_id  => {Accept => 'application/json'} => json => {
        "name" => "cdn_test3"
        })
    ->status_is(404)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

sub get_cdn_id {
    my $name = shift;
    my $q    = "select id from cdn where name = \'$name\'";
    my $get_svr = $dbh->prepare($q);
    $get_svr->execute();
    my $p = $get_svr->fetchall_arrayref( {} );
    $get_svr->finish();
    my $id = $p->[0]->{id};
    return $id;
}
