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

# NOTE:
#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $dbh    = Schema->database_handle;
my $schema = Schema->connect_to_database;
my $t      = Test::Mojo->new('TrafficOps');

Test::TestHelper->unload_core_data($schema);
Test::TestHelper->load_core_data($schema);

my $fixture_name;

#$fixture_name = 'server_edge1';
#ok $ds->load($fixture_name), 'Does the ' . $fixture_name . ' load?';
ok $schema->resultset('Cdn')->find( { name => 'cdn1' } ), 'cdn1 parameter exists?';
ok $schema->resultset('Profile')->find( { name => 'EDGE1' } ), 'Profile edge1 exists?';

ok $schema->resultset('Deliveryservice')->find( { xml_id => 'test-ds1' } ), 'Deliveryservice test-ds1 exists?';
$t->post_ok( '/login', => form => { u => 'admin', p => 'password' } )->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $q = "SELECT deliveryservice.id, 
           deliveryservice.xml_id, 
           (SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
               FROM origin o
               WHERE o.deliveryservice = deliveryservice.id
               AND o.is_primary) as org_server_fqdn,
           deliveryservice.type,
           profile.id AS profile_id, 
           cdn.name AS cdn_name 
     FROM deliveryservice 
     JOIN profile ON profile.id = deliveryservice.profile 
     JOIN cdn ON cdn.id = deliveryservice.cdn_id
     WHERE deliveryservice.active = 'true' ORDER BY RANDOM() LIMIT 1";

my $get_ds = $dbh->prepare($q);
$get_ds->execute();
my $p = $get_ds->fetchall_arrayref( {} );
$get_ds->finish();
my $ds    = $p->[0];
my $ds_id = $ds->{id};

my $host     = $ds->{org_server_fqdn};
my $cdn_name = $ds->{cdn_name};
my $ds_name  = $ds->{xml_id};

diag "Testing " . $ds_name . " (" . $host . ") in " . $cdn_name;
$q = 'SELECT DISTINCT server.host_name,
                      cdn.name as cdn_name 
					  FROM server 
	  JOIN profile ON profile.id = server.profile
	  JOIN cdn ON cdn.id = server.cdn_id
      WHERE cdn.name!=\'' . $cdn_name . '\' 
            and server.type in (select id from type where name =\'EDGE\' or name =\'MID\')
            and server.status in (select id from status where name =\'REPORTED\' or name =\'ONLINE\')';

my $get_servers = $dbh->prepare($q);
$get_servers->execute();
$p = $get_servers->fetchall_arrayref( {} );
$get_servers->finish();
my $slist;
my $j = 0;
while ( defined( $p->[$j] ) ) {
	$slist->{ $p->[$j]->{host_name} } = $p->[$j]->{cdn_name};

	# diag $p->[$j]->{host_name};
	$j++;
}
diag "Other servers " . $j;

# get all edges assicated with this cdn name
$q = 'SELECT DISTINCT server.host_name FROM server
JOIN profile ON profile.id = server.profile
JOIN cdn ON cdn.id = server.cdn_id
WHERE cdn.name=\'' . $cdn_name . '\' and server.type in (select id from type where name =\'EDGE\' or name =\'MID\')
and server.status in (select id from status where name =\'REPORTED\' or name =\'ONLINE\')';

# diag $q ;
$get_servers = $dbh->prepare($q);
$get_servers->execute();
$p = $get_servers->fetchall_arrayref( {} );
$get_servers->finish();
$j = 0;
while ( defined( $p->[$j] ) ) {
	$slist->{ $p->[$j]->{host_name} } = $cdn_name;

	# diag $p->[$j]->{host_name};
	$j++;
}
diag "Servers in " . $cdn_name . " CDN: " . $j;

$q = 'select id, username from tm_user where username = \'testuser\'';
my $get_user = $dbh->prepare($q);
$get_user->execute();
$p = $get_user->fetchall_arrayref( {} );
$get_user->finish();
my $juser_id = $p->[0]->{id};

# user is there, and it has a ds it can mod - let's post a job.
my $test_string = '/pa1/pa2/_asjauymejshqka_dedbf_339933.ism';
my $url         = 'http://' . $host . $test_string . '/.*';

my $job_id;

# make sure without X-Codebig-Principal header it's no good
$t->post_ok(
	'/job/external/new' => form => {
		'job.ds_xml_id'  => 'ds-cdn1',
		'job.agent'      => 1,
		'job.keyword'    => 'PURGE',
		'job.regex'      => '/foo/.*',
		'job.ttl'        => '48',
		'job.start_time' => '2014-12-12 14:57:45',
		'job.asset_type' => 'file',
	}
)->status_is(200);

done_testing();
exit();

$t->get_ok( '/job/external/result/view/' . $job_id )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( '/0/result' => 'No data found for this job ID' );

$t->post_ok(
	'/job/external/new' => form => {
		keyword    => 'PURGE',
		asset_url  => $url,
		parameters => 'TTL:48h',
		asset_type => 'SMOOTH'
	}
)->status_is(200)->json_is( '/status' => 'success' )->or( sub { diag $t->tx->res->content->asset->{content} } );

#  Get our job id
$job_id = $t->tx->res->json->{job};
diag "our job id is " . $job_id;

# make sure an invalid keyword fails
$t->post_ok(
	'/job/external/new' => form => {
		keyword    => 'POOP',
		asset_url  => $url,
		asset_type => 'SMOOTH'
	}
)->status_is(200)->json_is( '/status' => 'failure' );

# make sure an invalid host fails
$t->post_ok(
	'/job/external/new' => form => {
		keyword    => 'PURGE',
		asset_url  => 'http://foo.bar.net' . $test_string,
		asset_type => 'SMOOTH'
	}
)->status_is(200)->json_is( '/status' => 'failure' );

# check to see if it got posted
$t->get_ok( '/job/external/view/' . $job_id )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( '/0/keyword' => 'PURGE' )->json_is( '/0/status' => 'PENDING' )->json_is( '/0/asset_type' => 'SMOOTH' )->json_is( '/0/asset_url' => $url );

$t->get_ok( '/job/external/result/view/' . $job_id )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( '/0/result' => 'No data found for this job ID' );

# let's see if the agent can find it
$t->get_ok('/job/agent/viewpendingjobs/all')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $json = JSON->new->allow_nonref;
my $jdec = $json->decode( $t->tx->res->content->asset->{content} );
foreach my $job ( @{$jdec} ) {
	if ( $job->{id} == $job_id ) {
		ok( $job->{status} eq "PENDING" );
	}
}

# progress it
$t->get_ok( '/job/agent/statusupdate/' . $job_id )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# shouldn't show up anymore here
$t->get_ok('/job/agent/viewpendingjobs/all')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$json = JSON->new->allow_nonref;
$jdec = $json->decode( $t->tx->res->content->asset->{content} );
foreach my $job ( @{$jdec} ) {
	if ( $job->{id} == $job_id ) {
		ok( $job->{status} eq "PENDING" );
	}
}

# check to see if it got progressed
$t->get_ok( '/job/external/view/' . $job_id )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( '/0/keyword' => 'PURGE' )->json_is( '/0/status' => 'IN_PROGRESS' )->json_is( '/0/asset_type' => 'SMOOTH' )
	->json_is( '/0/asset_url' => $url );

# finish it
$t->post_ok(
	'/job/agent/result/new' => form => {
		job         => $job_id,
		agent       => 1,
		result      => 'COMPLETED',
		description => 'Test progressed successfully',
	}
)->status_is(200)->json_is( '/status' => 'success' )->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok( '/job/external/result/view/' . $job_id )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } )
	->json_is( '/0/agent' => 1 )->json_is( '/0/result' => 'COMPLETED' )->json_is( '/0/description' => 'Test progressed successfully' )
	->json_is( '/0/job' => $job_id );

diag "spot checking /genfiles...";
my $rand = 10 + int( rand(10) );    # only test every x server; there's too many, this gets too slow.
my $k    = 0;
foreach my $cache ( keys %{$slist} ) {
	$k++;
	next unless ( $k % $rand == 0 );
	note "checking /genfiles " . $cache . "(" . $k . ") for " . $slist->{$cache};
	my $response = $t->ua->get( '/genfiles/view/' . $cache . '/regex_revalidate.config' );
	my $content  = $response->res->content->asset->{content};
	if ( $slist->{$cache} eq $cdn_name ) {
		ok( $content =~ qr/$test_string/m, $cache . ' should have test_string in regex_revalidate.config' );
	}
	else {
		# JvD note: there are mutliple delivery services that are in both CDNs but have the same origin.
		# if ( $ds_name !~ /omg-.*/ ) {    #
		# 	ok( $content !~ qr/$test_string/m, $cache . ' should not have test_string in regex_revalidate.config' );
		# }
	}
}
ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();

