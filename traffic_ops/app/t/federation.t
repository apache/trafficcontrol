package main;

#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
use Test::TestHelper;
use strict;
use warnings;
use Schema;
use Fixtures::TmUser;
use Test::TestHelper;
use Fixtures::Federation;
use Fixtures::FederationDeliveryservice;
use Fixtures::FederationResolver;
use Fixtures::FederationFederationResolver;
use Fixtures::FederationTmuser;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');
my $t3_id;

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);
Test::TestHelper->teardown( $schema, 'FederationResolver' );
Test::TestHelper->teardown( $schema, 'Federation' );

#load core test data
Test::TestHelper->load_core_data($schema);

my $schema_values = { schema => $schema, no_transactions => 1 };

#
# FederationResolver
#
my $fr = Fixtures::FederationResolver->new($schema_values);
Test::TestHelper->load_all_fixtures($fr);
#
# Federation
#
my $fed = Fixtures::Federation->new($schema_values);
Test::TestHelper->load_all_fixtures($fed);

# FederationDeliveryservice
#
my $fds = Fixtures::FederationDeliveryservice->new($schema_values);
Test::TestHelper->load_all_fixtures($fds);

my $ffr = Fixtures::FederationFederationResolver->new($schema_values);
Test::TestHelper->load_all_fixtures($ffr);

my $ft = Fixtures::FederationTmuser->new($schema_values);
Test::TestHelper->load_all_fixtures($ft);

#login
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok('/federation')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();
