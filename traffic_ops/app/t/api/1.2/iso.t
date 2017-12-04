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
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
        ->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Should login?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "hostName" => "foo-bar-01",
            "domainName" => "baz.com",
            "rootPass" => "password",
            "dhcp" => "yes",
            "interfaceMtu" => 1500,
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "osversionDir is required" )
    , 'Does the iso generation fail due to missiong osversionDir?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "domainName" => "baz.com",
            "rootPass" => "password",
            "dhcp" => "yes",
            "interfaceMtu" => 1500,
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "hostName is required" )
    , 'Does the iso generation fail due to missiong hostName?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "hostName" => "foo-bar-01",
            "rootPass" => "password",
            "dhcp" => "yes",
            "interfaceMtu" => 1500,
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "domainName is required" )
    , 'Does the iso generation fail due to missiong domainName?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "hostName" => "foo-bar-01",
            "domainName" => "baz.com",
            "dhcp" => "yes",
            "interfaceMtu" => 1500,
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "rootPass is required" )
    , 'Does the iso generation fail due to missiong rootPass?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "hostName" => "foo-bar-01",
            "domainName" => "baz.com",
            "rootPass" => "password",
            "interfaceMtu" => 1500,
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "dhcp is required" )
    , 'Does the iso generation fail due to missiong dhcp?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "hostName" => "foo-bar-01",
            "domainName" => "baz.com",
            "rootPass" => "password",
            "dhcp" => "yes",
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "interfaceMtu is required" )
    , 'Does the iso generation fail due to missiong interfaceMtu?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "hostName" => "foo-bar-01",
            "domainName" => "baz.com",
            "rootPass" => "password",
            "dhcp" => "no",
            "interfaceMtu" => 1500,
            "ipNetmask" => "255.255.255.255",
            "ipGateway" => "10.10.10.10",
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "ipAddress is required if DHCP is no" )
    , 'Does the iso generation fail due to missiong ipAddress if DHCP is no?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "hostName" => "foo-bar-01",
            "domainName" => "baz.com",
            "rootPass" => "password",
            "dhcp" => "no",
            "interfaceMtu" => 1500,
            "ipAddress" => "10.10.10.10",
            "ipGateway" => "10.10.10.10",
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "ipNetmask is required if DHCP is no" )
    , 'Does the iso generation fail due to missiong ipNetmask if DHCP is no?';

ok $t->post_ok('/api/1.2/isos' => {Accept => 'application/json'} => json => {
            "osversionDir" => "centos72-netinstall",
            "hostName" => "foo-bar-01",
            "domainName" => "baz.com",
            "rootPass" => "password",
            "dhcp" => "no",
            "interfaceMtu" => 1500,
            "ipAddress" => "10.10.10.10",
            "ipNetmask" => "255.255.255.255",
            "disk" => "bond0",
        })
        ->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )
        ->json_is( "/alerts/0/text" => "ipGateway is required if DHCP is no" )
    , 'Does the iso generation fail due to missiong ipGateway if DHCP is no?';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
done_testing();

