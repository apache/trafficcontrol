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
use JSON;
use Data::Dumper;
use Test::TestHelper;
use File::Basename;
use File::Slurp;
use strict;
use warnings;

BEGIN { $ENV{MOJO_MODE} = "integration" }
my $t = Test::Mojo->new('TrafficOps');
no warnings 'once';
use warnings 'all';
use constant TEST_FILE => "/tmp/test.csv";


my $api_version = '1.1';
ok $t->post_ok( '/login', => form => { u => Test::TestHelper::ADMIN_USER, p => Test::TestHelper::ADMIN_USER_PASSWORD } )->status_is(302)
	->or( sub { diag $t->tx->res->content->asset->{content}; } );

my $json = JSON->new->allow_nonref;
$t->get_ok( '/api/' . $api_version . '/servers.json' )->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $servers = $json->decode( $t->tx->res->content->asset->{content} );

$t->get_ok('/datatype')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $type_arr = $json->decode( $t->tx->res->content->asset->{content} );
my %type     = ();
foreach my $tt ( @{$type_arr} ) {
	$type{ $tt->{name} } = $tt->{id};
}

$t->get_ok('/datastatus')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $status_arr = $json->decode( $t->tx->res->content->asset->{content} );
my %status     = ();
foreach my $s ( @{$status_arr} ) {
	$status{ $s->{name} } = $s->{id};
}

$t->get_ok('/dataprofile')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $profile_arr = $json->decode( $t->tx->res->content->asset->{content} );
my %profile     = ();
foreach my $p ( @{$profile_arr} ) {
	$profile{ $p->{name} } = $p->{id};
}

# check the legacy API
$t->get_ok('/datalocation')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

# my $cachegroup_arr = $json->decode( $t->tx->res->content->asset->{content} );
# my %cgroup         = ();
# foreach my $c ( @{$cachegroup_arr} ) {
#   $cgroup{ $c->{name} } = $c->{id};
# }

$t->get_ok('/api/1.1/cachegroups.json')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $response       = $json->decode( $t->tx->res->content->asset->{content} );
my $cachegroup_arr = $response->{response};
my %cgroup         = ();
foreach my $c ( @{$cachegroup_arr} ) {
	$cgroup{ $c->{name} } = $c->{id};
}

$t->get_ok('/dataphys_location')->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );
my $ploc_arr = $json->decode( $t->tx->res->content->asset->{content} );
my %ploc     = ();
foreach my $c ( @{$ploc_arr} ) {
	$ploc{ $c->{name} } = $c->{id};
}
my %type_done = ();
foreach my $server ( @{ $servers->{response} } ) {

	if ( !$type_done{ $server->{type} } ) {

		#diag Dumper($server);
		diag "Testing type " . $server->{type} . " with " . $server->{hostName};

		# these are optional, show up as undef in json and '' in html
		my $ilo_ip_address   = defined( $server->{iloIpAddress} )   ? $server->{iloIpAddress}   : '';
		my $ilo_ip_netmask   = defined( $server->{iloIpNetmask} )   ? $server->{iloIpNetmask}   : '';
		my $ilo_ip_gateway   = defined( $server->{iloIpGateway} )   ? $server->{iloIpGateway}   : '';
		my $ilo_username     = defined( $server->{iloUsername} )    ? $server->{iloUsername}    : '';
		my $router_host_name = defined( $server->{routerHostName} ) ? $server->{routerHostName} : '';
		my $router_port_name = defined( $server->{routerPortName} ) ? $server->{routerPortName} : '';
		$t->get_ok( '/server/' . $server->{id} . '/view' )->status_is(200)->text_is( 'td#host_name' => $server->{hostName} )
			->text_is( 'td#domain_name' => $server->{domainName} )->text_is( 'td#tcp_port' => $server->{tcpPort} )
			->text_is( 'td#interface_name' => $server->{interfaceName} )->text_is( 'td#ip_address' => $server->{ipAddress} )
			->text_is( 'td#ip_netmask' => $server->{ipNetmask} )->text_is( 'td#ip_gateway' => $server->{ipGateway} )
			->text_is( 'td#interface_mtu' => $server->{interfaceMtu} )->text_is( 'td#cachegroup' => $server->{cachegroup} )
			->text_is( 'td#type' => $server->{type} )->text_is( 'td#status' => $server->{status} )->text_is( 'td#profile' => $server->{profile} )
			->text_is( 'td#ilo_ip_address' => $ilo_ip_address )->text_is( 'td#ilo_ip_netmask' => $ilo_ip_netmask )
			->text_is( 'td#ilo_ip_gateway' => $ilo_ip_gateway )->text_is( 'td#ilo_username'   => $ilo_username )
			->text_is( 'td#router_host_name' => $router_host_name )->text_is( 'td#router_port_name' => $router_port_name);

		# the jsons associated with server
		$t->get_ok( '/server_by_id/' . $server->{id} )->status_is(200)->json_is( '/host_name', $server->{hostName} )
			->json_is( '/domain_name', $server->{domainName} )->json_is( '/tcp_port', $server->{tcpPort} )
			->json_is( '/interface_name', $server->{interfaceName} )->json_is( '/ip_address', $server->{ipAddress} )
			->json_is( '/ip_netmask', $server->{ipNetmask} )->json_is( '/ip_gateway', $server->{ipGateway} )
			->json_is( '/interface_mtu', $server->{interfaceMtu} );

		# flip some status stuff
		$t->post_ok( '/server/' . $server->{hostName} . '/status/REPORTED' => form => {} )->json_is( '/result' => 'SUCCESS' );
		$t->get_ok( '/dataserverdetail/select/' . $server->{hostName} )->status_is(200)->json_is( '/0/status' => 'REPORTED' )
			->or( sub { diag $t->tx->res->content->asset->{content}; } );
		$t->post_ok( '/server/' . $server->{hostName} . '/status/OFFLINE' => form => {} )->json_is( '/result' => 'SUCCESS' );
		$t->get_ok( '/dataserverdetail/select/' . $server->{hostName} )->status_is(200)->json_is( '/0/status' => 'OFFLINE' )
			->or( sub { diag $t->tx->res->content->asset->{content}; } );
		$t->post_ok( '/server/' . $server->{hostName} . '/status/ONLINE' => form => {} )->json_is( '/result' => 'SUCCESS' );
		$t->get_ok( '/dataserverdetail/select/' . $server->{hostName} )->status_is(200)->json_is( '/0/status' => 'ONLINE' )
			->or( sub { diag $t->tx->res->content->asset->{content}; } );

		$server->{router_host_name} = "UPDATED";
		#diag '/server/' . $server->{id};
		$t->post_ok(
			      '/server/'
				. $server->{id}
				. '/update' => form => {
				xmpp_id          => $server->{xmppId},
				xmpp_passwd      => $server->{xmpp_Passwd},
				interface_name   => $server->{interfaceName},
				ip_address       => $server->{ipAddress},
				ip_netmask       => $server->{ipNetmask},
				ip_gateway       => $server->{ipGateway},
				ip6_address      => $server->{ip6Address},
				ip6_gateway      => $server->{ip6Gateway},
				interface_mtu    => $server->{interfaceMtu},
				phys_location    => $server->{physLocation},
				cachegroup       => $server->{cachegroup},
				type             => $server->{type},
				profile          => $server->{profile},
				mgmt_ip_address  => $server->{mgmtIpAddress},
				mgmt_ip_netmask  => $server->{mgmtIpNetmask},
				mgmt_ip_gateway  => $server->{mgmtIpGateway},
				ilo_ip_address   => $server->{iloIpAddress},
				ilo_ip_netmask   => $server->{iloIpNetmask},
				ilo_ip_gateway   => $server->{iloIpGateway},
				ilo_username     => $server->{iloUsername},
				ilo_password     => $server->{iloPassword},
				router_host_name => $server->{routerHostName},
				router_port_name => $server->{routerPortName},
				https_port       => $server->{httpsPort},
				status           => $server->{status},
				}
		)->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );

		$type_done{ $server->{type} } = 1;
	}
}

sub build_tmpfile {
	my ($contents) = @_;
	unlink TEST_FILE;
	write_file( TEST_FILE, $contents );
}

# Header
my $header =
	"host,domain,int,ip4,subnet,gw,ip6,gw6,mtu,cdn,cachegroup,phys_loc,rack,type,prof,port,1g_ip,1g_subnet,1g_gw,ilo_ip,ilo_subnet,ilo_gw,ilo_user,ilo_pwd,r_host,r_port,https_port,offline_reason";

#----------------------------
# Good Test
my $content = join( "\n",
	$header,
	"good-host,chi.kabletown.net,bond0,10.10.2.200,255.255.255.0,10.10.2.254,2033:D0D0:3300::2:1A/64,2033:D0D0:3300::2:1,9000,CDN1,us-il-chicago,plocation-chi-1,rack33,EDGE,1,80,10.10.33.1,255.255.255.0,10.10.33.44,10.254.254.12,255.255.255.0,10.254.254.1,user,passwd,router_33,port_66,443,N/A\n"
);

&build_tmpfile($content);
my $asset = Mojo::Asset::File->new( path => TEST_FILE );

$t->post_ok(
	'/uploadhandlercsv' => form => { 'file-0' => { name => 'file-0', asset => $asset, filename => basename(TEST_FILE), content => $asset->slurp } } )
	->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#----------------------------
# Bad 'Type' look for -BAD
$content = join( "\n",
	$header,
	"atsec-chi-09,chi.kabletown.net,bond0,10.10.2.200,255.255.255.0,10.10.2.254,2033:D0D0:3300::2:1A/64,2033:D0D0:3300::2:1,9000,CDN1,us-il-chicago,plocation-chi-1,rack33,EDGE-BAD,1,80,10.10.33.1,255.255.255.0,10.10.33.44,10.254.254.12,255.255.255.0,10.254.254.1,user,passwd,router_33,port_66,443,N/A\n"
);

&build_tmpfile($content);
$t->post_ok(
	'/uploadhandlercsv' => form => { 'file-0' => { name => 'file-0', asset => $asset, filename => basename(TEST_FILE), content => $asset->slurp } } )
	->json_has("[EXCEPTION_ERROR]")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#----------------------------
# Bad 'Profile' look for -BAD
$content = join( "\n",
	$header,
	"atsec-chi-09,chi.kabletown.net,bond0,10.10.2.200,255.255.255.0,10.10.2.254,2033:D0D0:3300::2:1A/64,2033:D0D0:3300::2:1,9000,CDN2,us-il-chicago,plocation-chi-1,rack33,EDGE,1-BAD,80,10.10.33.1,255.255.255.0,10.10.33.44,10.254.254.12,255.255.255.0,10.254.254.1,user,passwd,router_33,port_66,443,N/A\n"
);

&build_tmpfile($content);
$t->post_ok(
	'/uploadhandlercsv' => form => { 'file-0' => { name => 'file-0', asset => $asset, filename => basename(TEST_FILE), content => $asset->slurp } } )
	->json_has("[EXCEPTION_ERROR]")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

#----------------------------
# Bad 'Cache Group' look for -BAD
$content = join( "\n",
	$header,
	"atsec-chi-09,chi.kabletown.net,bond0,10.10.2.200,255.255.255.0,10.10.2.254,2033:D0D0:3300::2:1A/64,2033:D0D0:3300::2:1,9000,CDN2,us-il-chicago-BAD,plocation-chi-1,rack33,EDGE,1,80,10.10.33.1,255.255.255.0,10.10.33.44,10.254.254.12,255.255.255.0,10.254.254.1,user,passwd,router_33,port_66,443,N/A\n"
);

#----------------------------
# Bad 'Physical Location' look for -BAD
$content = join( "\n",
	$header,
	"atsec-chi-09,chi.kabletown.net,bond0,10.10.2.200,255.255.255.0,10.10.2.254,2033:D0D0:3300::2:1A/64,2033:D0D0:3300::2:1,9000,CDN1,us-il-chicago,plocation-chi-1-BAD,rack33,EDGE,1,80,10.10.33.1,255.255.255.0,10.10.33.44,10.254.254.12,255.255.255.0,10.254.254.1,user,passwd,router_33,port_66,443,N/A\n"
);

&build_tmpfile($content);
$t->post_ok(
	'/uploadhandlercsv' => form => { 'file-0' => { name => 'file-0', asset => $asset, filename => basename(TEST_FILE), content => $asset->slurp } } )
	->json_has("[EXCEPTION_ERROR]")->status_is(200)->or( sub { diag $t->tx->res->content->asset->{content}; } );

ok $t->get_ok('/logout')->status_is(302)->or( sub { diag $t->tx->res->content->asset->{content}; } );
done_testing();
