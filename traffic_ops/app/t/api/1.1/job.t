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
use Data::Dumper;
use strict;
use warnings;
use Schema;
use Test::TestHelper;
use Fixtures::DeliveryserviceTmuser;
use Fixtures::DeliveryserviceServer;
use Fixtures::JobAgent;
use Fixtures::JobStatus;
use Fixtures::Job;
use Fixtures::Parameter;
use Fixtures::Server;
use POSIX qw(strftime);

BEGIN { $ENV{MOJO_MODE} = "test" }

# NOTE:
#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

Test::TestHelper->teardown( $schema, 'JobAgent' );
my $jobagent = Fixtures::JobAgent->new( { schema => $schema, no_transactions => 1 } );
Test::TestHelper->load_all_fixtures($jobagent);

Test::TestHelper->teardown( $schema, 'JobStatus' );
my $jobstatus = Fixtures::JobStatus->new( { schema => $schema, no_transactions => 1 } );
Test::TestHelper->load_all_fixtures($jobstatus);

Test::TestHelper->teardown( $schema, 'Job' );
my $jobs = Fixtures::Job->new( { schema => $schema, no_transactions => 1 } );
Test::TestHelper->load_all_fixtures($jobs);

ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::ADMIN_ROOT_USER, p => Test::TestHelper::ADMIN_ROOT_USER_PASSWORD } )->status_is(200),
	'Log in as admin user?';

#change the use_tenancy parameter to 0 (id from Parameters fixture) to test assigned dses table
ok $t->put_ok('/api/1.2/parameters/67' => {Accept => 'application/json'} => json =>
		{
			'value'       => '0',
		}
	)->status_is(200);

ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->post_ok( '/api/1.1/user/login', json => { u => Test::TestHelper::PORTAL_USER, p => Test::TestHelper::PORTAL_USER_PASSWORD } )->status_is(200),
	'Log into the portal user?';

ok $schema->resultset('Cdn')->find( { name => 'cdn1' } ), 'cdn1 parameter exists?';

ok $schema->resultset('Profile')->find( { name => 'EDGE1' } ), 'Profile edge1 exists?';

ok $schema->resultset('Deliveryservice')->find( { xml_id => 'test-ds1' } ), 'Deliveryservice test-ds1 exists?';

my $now = strftime( "%Y-%m-%d %H:%M:%S", gmtime() );

ok $t->get_ok('/api/1.1/user/current/jobs.json')->status_is(200)->json_has( '/response', 'has a response' ), 'No jobs returns successfully?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 100,
		regex     => '/foo1/.*',
		ttl       => 48,
		startTime => $now,
	}
	)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Create the first purge job?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 100,
		regex     => '/foo1/.*',
		ttl       => 0,
		startTime => $now,
	}
	)->status_is(400)->json_is( '/alerts', [ { level => "error", text => "ttl should be between 1 and 72" } ] )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Will not create a purge job without the ttl in the proper low range?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 100,
		regex     => '/foo1/.*',
		ttl       => 3000,
		startTime => $now,
	}
	)->status_is(400)->json_is( '/alerts', [ { level => "error", text => "ttl should be between 1 and 72" } ] )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Will not create a purge job without the ttl in the proper high range?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 100,
		regex     => '/foo2/.*',
		ttl       => 49,
		startTime => $now,
	}
)->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'Create a second purge job?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 800,
		regex     => '/foo2/.*',
		ttl       => 49,
		startTime => $now,
	}
)->status_is(403)->or( sub { diag $t->tx->res->content->asset->{content}; } ), 'purge job for this service not authorized for this user';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 100,
		ttl       => 49,
		startTime => $now,
	}
	)->status_is(400)->json_is( '/alerts', [ { level => "error", text => "regex is required" } ] )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Will not create a purge job without the regex field?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId  => 100,
		regex => '/foo2/.*',
		ttl   => 49,
	}
	)->status_is(400)->json_is( '/alerts', [ { level => "error", text => "startTime is required" } ] )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Will not create a purge job without the startTime field?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 100,
		regex     => '/foo2/.*',
		ttl       => 49,
		startTime => '2015-01-09',
	}
	)->status_is(400)
	->json_is( '/alerts', [ { level => "error", text => "startTime has an invalidate date format, should be in the form of YYYY-MM-DD HH:MM:SS" } ] )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Will not create a purge job without a properly formatted startTime field?';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		dsId      => 100,
		regex     => '/foo2/.*',
		startTime => $now,
	}
	)->status_is(400)->json_is( '/alerts', [ { level => "error", text => "ttl is required" } ] )
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Will not create a purge job without the ttl field?';

ok $t->get_ok('/api/1.1/user/current/jobs.json')->status_is(200)->json_has('test-ds1.edge/foo1/')->json_has('TTL:48h')
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the first purge job exist?';

ok $t->get_ok('/api/1.1/user/current/jobs.json')->status_is(200)->json_has('test-ds1.edge/foo2/')->json_has('TTL:49h')
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the second purge job exist?';

ok $t->get_ok('/api/1.1/user/current/jobs.json?keyword=PURGE')->status_is(200)->json_has('test-ds1.edge/foo1/')->json_has('TTL:48h')
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the correct job return with keyword=PURGE?';

ok $t->get_ok('/api/1.1/user/current/jobs.json?keyword=PURGE&dsId=1')->status_is(200)->json_has('test-ds1.edge/foo1/')->json_has('TTL:48h')
	->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Does the correct job return with keyword=PURGE? and dsId=1';

ok $t->post_ok(
	'/api/1.1/user/current/jobs',
	json => {
		regex     => '/foo1/.*',
		ttl       => 48,
		startTime => $now,
	}
	)->status_is(400)->or( sub { diag $t->tx->res->content->asset->{content}; } ),
	'Will not create a purge job without the dsId?';

ok $t->post_ok('/api/1.1/user/logout')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect();
done_testing();
