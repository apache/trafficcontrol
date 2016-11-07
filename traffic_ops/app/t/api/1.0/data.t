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

# These test cases are to ensure the old 1.0 endpoints do not error.
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

$t->get_ok('/datacrans')->status_is(200)->or( sub                               { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datacrans/orderby/asn')->status_is(200)->or( sub                   { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datadeliveryservice')->status_is(200)->or( sub                     { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datadeliveryservice/orderby/id')->status_is(200)->or( sub          { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datadeliveryserviceserver')->status_is(200)->or( sub               { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datadomains')->status_is(200)->or( sub                             { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datahwinfo')->status_is(200)->or( sub                              { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datahwinfo/orderby/id')->status_is(200)->or( sub                   { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datalog')->status_is(200)->or( sub                                 { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datalog/1')->status_is(200)->or( sub                               { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datalinks')->status_is(200)->or( sub                               { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datalinks/orderby/server')->status_is(200)->or( sub                { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datalocation')->status_is(200)->or( sub                            { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datalocation/orderby/id')->status_is(200)->or( sub                 { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datalocationparameter')->status_is(200)->or( sub                   { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataparameter')->status_is(200)->or( sub                           { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataparameter')->status_is(200)->or( sub                           { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataparameter/edge1')->status_is(200)->or( sub                     { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataparameter/orderby/id')->status_is(200)->or( sub                { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataparameter/orderby/id')->status_is(200)->or( sub                { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataphys_location')->status_is(200)->or( sub                       { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataphys_locationtrimmed')->status_is(200)->or( sub                { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataprofile')->status_is(200)->or( sub                             { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataprofile/orderby/id')->status_is(200)->or( sub                  { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataprofileparameter')->status_is(200)->or( sub                    { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataprofileparameter/orderby/profile')->status_is(200)->or( sub    { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataprofiletrimmed')->status_is(200)->or( sub                      { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataregion')->status_is(200)->or( sub                              { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datarole')->status_is(200)->or( sub                                { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datarole/orderby/id')->status_is(200)->or( sub                     { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataserver')->status_is(200)->or( sub                              { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataserver/orderby/id')->status_is(200)->or( sub                   { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/dataserverdetail/select/atlanta-edge-01')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datastaticdnsentry')->status_is(200)->or( sub                      { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datastatus')->status_is(200)->or( sub                              { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datastatus/orderby/id')->status_is(200)->or( sub                   { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datatype')->status_is(200)->or( sub                                { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datatype/orderby/id')->status_is(200)->or( sub                     { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datatypetrimmed')->status_is(200)->or( sub                         { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datauser')->status_is(200)->or( sub                                { diag $t->tx->res->content->asset->{content}; } );
$t->get_ok('/datauser/orderby/id')->status_is(200)->or( sub                     { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
