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
use Test::TestHelper;
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#login
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

#NOTE: These route parameters may cause problems until Mock/Fixture data can be setup accordingly.
#ok $t->get_ok('/api/1.1/metrics/g/kbps/1423343701/1423602901.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();
