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
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

####################### RW testing - careful with these! #####################################################

#these ids are created by fixtures, so we know them.
my $dsid = 100;
my $atypeid = 20;
my $locid = 1;

# assign some static dns entries to this delivery service
$t->post_ok(
	'/staticdnsentry/' . $dsid . '/update' => form => {
		host_new_0            => 'auto_test_1',
		address_new_0         => '69.241.22.21',
		type_new_0            => $atypeid,
		loc_new_0             => $locid,
		ttl_new_0             => 3600,
		deliveryservice_new_0 => $dsid,
		}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $q   = 'select id from staticdnsentry where host like \'%auto_test_1%\'';
my $get = $dbh->prepare($q);
$get->execute();
my $p = $get->fetchall_arrayref( {} );
$get->finish();
my $sdnsid = $p->[0]->{id};

diag 'sdnsid:' . $sdnsid;

# update it's address
$t->post_ok(
	'/staticdnsentry/' . $dsid . '/update' => form => {
		'host_' . $sdnsid            => 'auto_test_1',
		'address_' . $sdnsid         => '69.241.33.27',
		'type_' . $sdnsid            => $atypeid,
		'loc_' . $sdnsid             => $locid,
		'ttl_' . $sdnsid             => 3600,
		'deliveryservice_' . $sdnsid => $dsid,
		}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# update it's address and a combi
$t->post_ok(
	'/staticdnsentry/' . $dsid . '/update' => form => {
		'host_' . $sdnsid            => 'auto_test_1',
		'address_' . $sdnsid         => '69.241.33.27',
		'type_' . $sdnsid            => $atypeid,
		'loc_' . $sdnsid             => $locid,
		'ttl_' . $sdnsid             => 3600,
		'deliveryservice_' . $sdnsid => $dsid,
		host_new_0                   => 'auto_test_2',
		address_new_0                => '69.241.22.51',
		type_new_0                   => $atypeid,
		loc_new_0                    => $locid,
		ttl_new_0                    => 3600,
		deliveryservice_new_0        => $dsid,
		}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
