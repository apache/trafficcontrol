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
use Schema;
use strict;
use warnings;
use Test::TestHelper;
use Fixtures::TmUser;
use Fixtures::Federation;
use Fixtures::FederationDeliveryservice;
use Fixtures::FederationResolver;
use Fixtures::FederationFederationResolver;
use Fixtures::FederationTmuser;
use Data::Dumper;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $t      = Test::Mojo->new("TrafficOps");
my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $schema_values = { schema => $schema, no_transactions => 1 };

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_all_fixtures( Fixtures::Tenant->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Cdn->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Role->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::TmUser->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Type->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Profile->new($schema_values) );
Test::TestHelper->load_all_fixtures( Fixtures::Deliveryservice->new($schema_values) );

Test::TestHelper->teardown( $schema, "Federation" );
Test::TestHelper->teardown( $schema, "FederationDeliveryservice" );
Test::TestHelper->teardown( $schema, "FederationFederationResolver" );
Test::TestHelper->teardown( $schema, "FederationResolver" );

my $federation = Fixtures::Federation->new($schema_values);
Test::TestHelper->load_all_fixtures($federation);

# FederationDeliveryservice
#
my $fmd = Fixtures::FederationDeliveryservice->new($schema_values);
Test::TestHelper->load_all_fixtures($fmd);

#my $federation_federation_resolver = Fixtures::FederationFederationResolver->new($schema_values);
#Test::TestHelper->load_all_fixtures($federation_federation_resolver);

my $ft = Fixtures::FederationTmuser->new($schema_values);
Test::TestHelper->load_all_fixtures($ft);

####### Federation User ###########################################################################
ok $t->post_ok(
	"/login",
	=> form => {
		u => Test::TestHelper::FEDERATION_USER,
		p => Test::TestHelper::FEDERATION_USER_PASSWORD
	}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->post_ok( "/api/1.2/federations", json => { federations => [ { deliveryService => "test-ds1" } ] } )->status_is(400)
	->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text/", "mappings is required" );

####### Add API #########################################################################
$t->post_ok(
	"/api/1.2/federations",
	json => {
		federations => [ { mappings => { resolve4 => ["127.0.0.1/32"] } } ]
	}
)->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } )->json_is( "/alerts/0/text/", "deliveryService is required" );

$t->post_ok(
	"/api/1.2/federations",
	json => {
		federations => [
			{
				deliveryService => "test-ds1",
				mappings        => {
					resolve4 => ["127.1.1.1/32"],
					resolve6 => ["fd06:d8c6:14:eeee/123"]
				}
			}
		]
	}
)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/federations.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/deliveryService", "test-ds1" )->json_is( "/response/0/mappings/0/cname", "cname1." )
	->json_is( "/response/0/mappings/0/ttl", "86400" )->json_is( "/response/0/mappings/0/resolve4/0", "127.1.1.1/32" );


ok $t->get_ok("/api/1.2/federations.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/deliveryService", "test-ds1" )->json_is( "/response/0/mappings/0/cname", "cname1." )
	->json_is( "/response/0/mappings/0/ttl", "86400" )->json_is( "/response/0/mappings/0/resolve4/0", "127.1.1.1\/32" )
	->json_is( "/response/0/mappings/0/resolve4/0", "127.1.1.1/32" )

	->json_is( "/response/1/deliveryService", "test-ds2" )->json_is( "/response/1/mappings/0/cname", "cname2." )
	->json_is( "/response/1/mappings/0/ttl", "86400" );

ok $t->delete_ok("/api/1.2/federations")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok("/api/1.2/federations.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_hasnt( "/response/0/mappings/0/resolve6/0", "FE80::0202:B3FF:FE1E:8329/128" )->json_hasnt( "/response/0/mappings/0/resolve4/0", "127.0.0.1/32" )

	->json_hasnt( "/response/1/mappings/0/resolve6/0", "FE80::0202:B3FF:FE1E:8330/128" )->json_hasnt( "/response/1/mappings/0/resolve4/0", "127.0.0.2/32" );

####### Update API ######################################################################
$t->put_ok(
	"/api/1.2/federations",
	json => {
		federations => [
			{
				deliveryService => "test-ds1",
				mappings        => { resolve4 => ["127.4.4.4/32"] }
			}
		]
	}
)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/federations.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/deliveryService", "test-ds1" )->json_is( "/response/0/mappings/0/cname", "cname1." )
	->json_is( "/response/0/mappings/0/ttl", "86400" )->json_is( "/response/0/mappings/0/resolve4/0", "127.4.4.4/32" );

$t->put_ok(
	"/api/1.2/federations",
	json => {
		federations => [
			{
				deliveryService => "test-ds1",
				mappings        => { resolve4 => ["255.255.255.255/32"] }
			}
		]
	}
)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/federations.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( "/response/0/deliveryService", "test-ds1" )->json_is( "/response/0/mappings/0/cname", "cname1." )
	->json_is( "/response/0/mappings/0/ttl", "86400" )->json_is( "/response/0/mappings/0/resolve4/0", "255.255.255.255/32" );

ok $t->get_ok("/logout")->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

####### Admin User ################################################################################
ok $t->post_ok( "/login", => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok("/api/1.2/federations.json")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok("/logout")->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );


####### Cleanup DB ######################################################################
#Test::TestHelper->teardown( $schema, "Federation" );
#Test::TestHelper->teardown( $schema, "FederationDeliveryservice" );
#Test::TestHelper->teardown( $schema, "FederationFederationResolver" );
#Test::TestHelper->teardown( $schema, "FederationResolver" );

$dbh->disconnect();
done_testing();
