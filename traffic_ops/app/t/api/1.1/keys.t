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
use POSIX ();
use Mojo::Base -strict;
use Test::More;
use Test::Mojo;
use Test::TestHelper;
use Test::MockModule;
use Test::MockObject;
use strict;
use warnings;
use JSON;
use Data::Dumper;
use warnings 'all';
use Test::TestHelper;
use Test::MockModule;
no warnings 'once';
use Schema;


#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }
my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

#load data so we can login
#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

#login
ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(200),
	'Log into the admin user?';
	
my $fake_lwp = new Test::MockModule( 'LWP::UserAgent', no_auto => 1 );
my $fake_get = HTTP::Response->new( 200, undef, HTTP::Headers->new, "OK");
$fake_lwp->mock( 'get', sub { return $fake_get } );
my $fake_put = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'put', sub { return $fake_put } );
my $fake_delete = HTTP::Response->new( 204, undef, HTTP::Headers->new, undef );
$fake_lwp->mock( 'delete', sub { return $fake_delete } );

#ping
ok $t->get_ok("/api/1.1/keys/ping.json")
->json_is("/response", "OK")
->status_is(200)
->or( sub { diag $t->tx->res->content->asset->{content}; } );

# logout
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
done_testing();
