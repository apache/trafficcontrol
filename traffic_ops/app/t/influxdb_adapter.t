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
use Test::MockModule;
use Extensions::TrafficStats::Connection::InfluxDBAdapter;
use Data::Dumper;

use constant SERVER                => "radon-influxdb.local:8086";
use constant DB_NAME               => "mydb";
use constant SERIES_NAME           => "myseries";
use constant RETENTION_POLICY_NAME => "mypolicy";

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

#GET
my $influxdb_util = Extensions::TrafficStats::Connection::InfluxDBAdapter->new( Test::TestHelper::ADMIN_USER, Test::TestHelper::ADMIN_USER_PASSWORD );
$influxdb_util->set_db_name(DB_NAME);
$influxdb_util->set_server(SERVER);

my $fake_answer = "OK";
my $fake_lwp    = new Test::MockModule( 'LWP::UserAgent', no_auto => 1 );
my $fake_header = HTTP::Headers->new;
$fake_header->header( 'Content-Type' => 'application/json' );    # set
my $fake_response = HTTP::Response->new( 200, undef, $fake_header, $fake_answer );
$fake_lwp->mock( 'post', sub { return $fake_response } );
$fake_lwp->mock( 'get',  sub { return $fake_response } );

my $query = "DROP DATABASE " . DB_NAME;
my $response = $influxdb_util->query( undef, $query );
is( $response->{_rc},      200 );
is( $response->{_content}, $fake_answer );

$query = "CREATE DATABASE " . DB_NAME;
$response = $influxdb_util->query( undef, $query );
is( $response->{_rc},      200 );
is( $response->{_content}, $fake_answer );

$query = "CREATE RETENTION POLICY " . RETENTION_POLICY_NAME . " ON " . DB_NAME . " DURATION 365d REPLICATION 1 DEFAULT";
$response = $influxdb_util->query( DB_NAME, $query );
is( $response->{_rc},      200 );
is( $response->{_content}, $fake_answer );

my $write_point = {
	"database"        => DB_NAME,
	"retentionPolicy" => "mypolicy",
	"points"          => [
		{
			"name"   => SERIES_NAME,
			"fields" => { value => 11111 }
		}, {
			"name"   => SERIES_NAME,
			"fields" => { value => 22222 }
		}, {
			"name"   => SERIES_NAME,
			"fields" => { value => 33333 }
		}
	]
};
$response = $influxdb_util->write($write_point);
is( $response->{_rc},      200 );
is( $response->{_content}, $fake_answer );

done_testing();
