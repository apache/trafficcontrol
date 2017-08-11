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
use strict;
use warnings;
use Schema;
use Test::TestHelper;
use Fixtures::TmUser;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok(
	'/user',
	=> form => {
		'tm_user.full_name'            => 'fullname',
		'tm_user.username'             => 'testcase',
		'tm_user.public_ssh_key'	   => 'ssh-key',
		'tm_user.phone_number'         => 'phone_number',
		'tm_user.email'                => 'email@email.com',
		'tm_user.local_passwd'         => 'password',
		'tm_user.confirm_local_passwd' => 'password',
		'tm_user.role'                 => 1,
		'tm_user.company'              => 'ABC Company',
	}
)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Can a user be created?';

ok $t->get_ok('/datauser')->status_is(200)->json_is( '/0/username', 'admin' )->json_is( '/0/role', 4 ), 'Does the admin username exist?';

ok $t->get_ok('/datauser/orderby/role')->status_is(200)->json_has('/0/rolename')->json_has('/0/username')->json_has('/0/id')->json_has('/0/role'),
	'Does the user sort by role?';

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
done_testing();
