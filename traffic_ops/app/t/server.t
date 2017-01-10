
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
use Mojo::Util qw/squish/;
use DBI;
use strict;
use warnings;
no warnings 'once';
use warnings 'all';
use Test::TestHelper;
use Schema;

#no_transactions=>1 ==> keep fixtures after every execution, beware of duplicate data!
#no_transactions=>0 ==> delete fixtures after every execution

BEGIN { $ENV{MOJO_MODE} = "test" }

my $schema = Schema->connect_to_database;
my $dbh    = Schema->database_handle;
my $t      = Test::Mojo->new('TrafficOps');

#unload data for a clean test
Test::TestHelper->unload_core_data($schema);

#load core test data
Test::TestHelper->load_core_data($schema);

my $q           = 'select * from server where type = 1 limit 1';
my $get_servers = $dbh->prepare($q);
$get_servers->execute();
my $svr = $get_servers->fetchall_arrayref( {} );
$get_servers->finish();
my $test_server_id   = $svr->[0]->{id};
my $profile_name     = $dbh->selectrow_array( 'select name from profile where id=' . $svr->[0]->{profile} );
my $type_name        = $dbh->selectrow_array( 'select name from type where id=' . $svr->[0]->{type} );
my $status_name      = $dbh->selectrow_array( 'select name from status where id=' . $svr->[0]->{status} );
my $location_name    = $dbh->selectrow_array( 'select name from phys_location where id=' . $svr->[0]->{phys_location} );
my $cgroup_name      = $dbh->selectrow_array( 'select name from cachegroup where id=' . $svr->[0]->{cachegroup} );
my $router_port_html = defined( $svr->[0]->{router_port_name} ) ? squish( $svr->[0]->{router_port_name} ) : '';

$t->post_ok( '/login', => form => { u => 'admin', p => 'password' } )->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# /server/view/id page
$t->get_ok("/server/$test_server_id/view")->status_is(200)->text_is( 'td#host_name' => $svr->[0]->{host_name} )
	->text_is( 'td#domain_name' => $svr->[0]->{domain_name} )->text_is( 'td#tcp_port' => $svr->[0]->{tcp_port} )
	->text_is( 'td#interface_name' => $svr->[0]->{interface_name} )->text_is( 'td#ip_address' => $svr->[0]->{ip_address} )
	->text_is( 'td#ip_netmask' => $svr->[0]->{ip_netmask} )->text_is( 'td#ip_gateway' => $svr->[0]->{ip_gateway} )
	->text_is( 'td#interface_mtu' => $svr->[0]->{interface_mtu} )->text_is( 'td#cachegroup' => $cgroup_name )->text_is( 'td#type' => $type_name )
	->text_is( 'td#status' => $status_name )->text_is( 'td#profile' => $profile_name )->text_is( 'td#ilo_ip_address' => $svr->[0]->{ilo_ip_address} )
	->text_is( 'td#ilo_ip_netmask' => $svr->[0]->{ilo_ip_netmask} )->text_is( 'td#ilo_ip_gateway' => $svr->[0]->{ilo_ip_gateway} )
	->text_is( 'td#ilo_username' => $svr->[0]->{ilo_username} )->text_is( 'td#ilo_password' => $svr->[0]->{ilo_password} )
	->text_is( 'td#router_host_name' => defined( $svr->[0]->{router_host_name} ) ? $svr->[0]->{router_host_name} : '' )
	->text_is( 'td#router_port_name' => $router_port_html )->or( sub { diag $t->tx->res->content->asset->{content}; } );

# the jsons associated with server
$t->get_ok( '/server_by_id/' . $test_server_id )->status_is(200)->json_is( '/host_name', $svr->[0]->{host_name} )
	->json_is( '/domain_name', $svr->[0]->{domain_name} )->json_is( '/tcp_port', $svr->[0]->{tcp_port} )
	->json_is( '/interface_name', $svr->[0]->{interface_name} )->json_is( '/ip_address', $svr->[0]->{ip_address} )
	->json_is( '/ip_netmask', $svr->[0]->{ip_netmask} )->json_is( '/ip_gateway', $svr->[0]->{ip_gateway} )
	->json_is( '/interface_mtu', $svr->[0]->{interface_mtu} )->or( sub { diag $t->tx->res->content->asset->{content}; } );

# the jsons associated with server
$t->get_ok("/server_by_id/$test_server_id")->status_is(200)->json_is( '/host_name', $svr->[0]->{host_name} )
	->json_is( '/domain_name', $svr->[0]->{domain_name} )->json_is( '/tcp_port', $svr->[0]->{tcp_port} )
	->json_is( '/interface_name', $svr->[0]->{interface_name} )->json_is( '/ip_address', $svr->[0]->{ip_address} )
	->json_is( '/ip_netmask', $svr->[0]->{ip_netmask} )->json_is( '/ip_gateway', $svr->[0]->{ip_gateway} )
	->json_is( '/interface_mtu', $svr->[0]->{interface_mtu} )->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok('/dataserver')->status_is(200)->json_has('/0/router_port_name')->json_has('/0/ilo_username')->json_has('/0/profile')
	->json_has('/0/interface_mtu')->json_has('/0/status')->json_has('/0/ip_netmask')->json_has('/0/ilo_password')->json_has('/0/mgmt_ip_netmask')
	->json_has('/0/tcp_port')->json_has('/0/id')->json_has('/0/ip_address')->json_has('/0/mgmt_ip_address')->json_has('/0/ilo_ip_address')
	->json_has('/0/interface_name')->json_has('/0/last_updated')->json_has('/0/cachegroup')->json_has('/0/ilo_ip_netmask')->json_has('/0/mgmt_ip_gateway')
	->json_has('/0/ilo_ip_gateway')->json_has('/0/domain_name')->json_has('/0/router_host_name')->json_has('/0/ip_gateway')->json_has('/0/host_name')
	->json_has('/0/type')->or( sub { diag $t->tx->res->content->asset->{content}; } );

$t->get_ok('/dataserver/orderby/cachegroup')->status_is(200)->json_has('/0/router_port_name')->json_has('/0/ilo_username')->json_has('/0/profile')
	->json_has('/0/interface_mtu')->json_has('/0/status')->json_has('/0/ip_netmask')->json_has('/0/ilo_password')->json_has('/0/mgmt_ip_netmask')
	->json_has('/0/tcp_port')->json_has('/0/id')->json_has('/0/ip_address')->json_has('/0/mgmt_ip_address')->json_has('/0/ilo_ip_address')
	->json_has('/0/interface_name')->json_has('/0/last_updated')->json_has('/0/cachegroup')->json_has('/0/ilo_ip_netmask')->json_has('/0/mgmt_ip_gateway')
	->json_has('/0/ilo_ip_gateway')->json_has('/0/domain_name')->json_has('/0/router_host_name')->json_has('/0/ip_gateway')->json_has('/0/host_name')
	->json_has('/0/type')->or( sub { diag $t->tx->res->content->asset->{content}; } );

# $t->get_ok('/dataserverdetail/select/1')->status_is(200)->json_has('/1/host_name');

#clean up old crud
&upd_and_del();

# create a new server
$t->post_ok(
	'/server/create' => form => {
		host_name        => 'test-01',
		domain_name      => 'kabletown.net',
		tcp_port         => '80',
		xmpp_id          => 'test-01@kabletown.net',
		xmpp_passwd      => 'very_secret_passwd',
		interface_name   => 'bond0',
		ip_address       => '3.3.3.3',
		ip_netmask       => '255.255.255.0',
		ip_gateway       => '3.3.3.9',
		ip6_address      => '2009:334:333::2/64',
		ip6_gateway      => '2009:334:333::1',
		interface_mtu    => '9000',
		phys_location    => 100,
		cachegroup       => 100,
		type             => 1,
		profile          => 100,
		cdn              => 100,
		mgmt_ip_address  => '192.168.1.1',
		mgmt_ip_gateway  => '192.168.1.2',
		mgmt_ip_netmask  => '255.255.255.0',
		ilo_ip_address   => '3.9.9.3',
		ilo_ip_netmask   => '255.255.255.0',
		ilo_ip_gateway   => '3.9.9.9',
		ilo_username     => 'user',
		ilo_password     => 'tt',
		router_host_name => 'ur091.home.net',
		router_port_name => 'ae99.99',
		https_port       => '443',
		offline_reason   => 'N/A'
	}
)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# modify and delete it
&upd_and_del();

sub upd_and_del() {
	my $q          = 'select id, host_name from server where host_name = \'test-01\'';
	my $get_server = $dbh->prepare($q);
	$get_server->execute();
	my $p = $get_server->fetchall_arrayref( {} );
	$get_server->finish();
	my $i = 0;
	while ( defined( $p->[$i] ) ) {

		diag $p->[$i]->{host_name};
		my $id = $p->[$i]->{id};
		$t->post_ok(
			"/server/$id/update" => form => {
				xmpp_id          => 'test-01@cdn.net',
				xmpp_passwd      => 'very_secretother_passwd',
				interface_name   => 'bond0',
				ip_address       => '3.3.3.9',
				ip_netmask       => '255.255.255.0',
				ip_gateway       => '3.3.3.1',
				ip6_address      => '2009:334:333::2/64',
				ip6_gateway      => '2009:334:333::1',
				interface_mtu    => '9000',
				phys_location    => 100,
				cachegroup       => 100,
				type             => 1,
				profile          => 100,
				cdn_id           => 100,
				mgmt_ip_address  => '192.168.3.1',
				mgmt_ip_netmask  => '192.168.3.2',
				mgmt_ip_gateway  => '255.255.255.0',
				ilo_ip_address   => '3.9.3.3',
				ilo_ip_netmask   => '255.255.255.0',
				ilo_ip_gateway   => '3.9.3.9',
				ilo_username     => 'user',
				ilo_password     => 'tt',
				router_host_name => '',
				router_port_name => '',
				https_port       => '443',
			}
		)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

		$t->post_ok( '/server/test-01/status/REPORTED' => form => {} )->json_is( '/result' => 'SUCCESS' );
		$t->get_ok('/dataserverdetail/select/test-01')->status_is(200)->json_is( '/0/status' => 'REPORTED' );
		$t->post_ok( '/server/test-01/status/OFFLINE' => form => {} )->json_is( '/result' => 'SUCCESS' );
		$t->get_ok('/dataserverdetail/select/test-01')->status_is(200)->json_is( '/0/status' => 'OFFLINE' );
		$t->post_ok( '/server/test-01/status/ONLINE' => form => {} )->json_is( '/result' => 'SUCCESS' );
		$t->get_ok('/dataserverdetail/select/test-01')->status_is(200)->json_is( '/0/status' => 'ONLINE' );

		diag $id;
		$t->get_ok("/server/$id/delete")->status_is(302);
		$i++;
	}
}

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
$dbh->disconnect();
done_testing();
