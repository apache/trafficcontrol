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
use DBI;
use strict;
use warnings;
use Data::Dumper;
use warnings 'all';
use Schema;
use Test::TestHelper;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }
my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');
no warnings 'once';

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $query = "select distinct profile.name 
	        from profile 
	        left join profile_parameter on profile.id=profile_parameter.profile 
			left join parameter on parameter.id=profile_parameter.parameter 
			where parameter.config_file='rascal.properties'";
my $select = $dbh->prepare($query);
$select->execute();
my @checked_profiles;
while ( my @row = $select->fetchrow_array ) {
	push( @checked_profiles, $row[0] );
}

$query  = "select * from parameter where config_file='rascal.properties'";
$select = $dbh->prepare($query);
$select->execute();
my $lines;
while ( my @row = $select->fetchrow_array ) {
	push( @{ $lines->{ $row[0] } }, @row );
}

$query  = "select xml_id, global_max_mbps, global_max_tps from deliveryservice where active=true";
$select = $dbh->prepare($query);
$select->execute();
my $ds;
while ( my @row = $select->fetchrow_array ) {
	push( @{ $ds->{ $row[0] } }, @row );
}

$t->get_ok('/health')->status_is(200);
foreach my $profile (@checked_profiles) {
	$t->json_has($profile);
}
foreach my $parameter ( sort keys %{$lines} ) {
	foreach my $item ( @{ $lines->{$parameter} } ) {
		$t->json_has($item);
	}
}
foreach my $xml_id ( sort keys %{$ds} ) {
	foreach my $item ( @{ $ds->{$xml_id} } ) {
		if ( defined($item) ) {
			$t->json_has($item);
		}
	}
}

$t->get_ok('/healthfull')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
foreach my $profile (@checked_profiles) {
	$t->json_has($profile);
}
foreach my $parameter ( sort keys %{$lines} ) {
	foreach my $item ( @{ $lines->{$parameter} } ) {
		$t->json_has($item);
	}
}

$query  = "select value from parameter where name='CDN_name'";
$select = $dbh->prepare($query);
$select->execute();
my @cdn_names;
while ( my @row = $select->fetchrow_array ) {
	push( @cdn_names, $row[0] );
}
$query  = "select name, value from parameter where config_file='rascal-config.txt' and name != 'location'";
$select = $dbh->prepare($query);
$select->execute();
$lines = ();
while ( my @row = $select->fetchrow_array ) {
	push( @{ $lines->{ $row[0] } }, @row );
}

$t->get_ok("/health/cdn2.json")->status_is(200)->json_is( "/profiles/MID/MID1/health.threshold.loadavg", "25.0" )
	->json_is( "/profiles/MID/MID1/history.count", "30" )->json_is( "/deliveryServices/test-ds5/status", "REPORTED" )
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

$dbh->disconnect;

done_testing();
