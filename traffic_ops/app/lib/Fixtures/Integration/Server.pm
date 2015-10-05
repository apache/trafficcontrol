package Fixtures::Integration::Server;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Data::Dumper;

my %definition_for = ();

sub gen_data {
	my @cache_groups = ( 'nyc', 'lax', 'chi', 'hou', 'phl', 'den' );
	my @profiles = ( 45, 45, 45, 45, 26, 26, 26, 26, 26 );
	my $cgr_no = 0;



	my $cdn_profiles->{"45"} = 1;
	$cdn_profiles->{"26"} = 2;

	# EDGES - 8 per cache group
	my $site = 0;
	my $id   = 1;    # so we have predictable id numbers
	foreach my $cg (@cache_groups) {
		my $net1 = '10.10.' . $cgr_no . '.';
		my $net2 = '172.16.' . $cgr_no . '.';
		my $net3 = '2033:D0D0:3300::' . $cgr_no . ":";

		foreach my $i ( 0 .. 8 ) {
			if ( $i == 4 ) {    # half of each cg is in a site
				$site++;
			}
			my $profile_id = $profiles[$i];
			my $cdn_id = $cdn_profiles->{$profile_id};

			my $hostname = 'atsec-' . $cache_groups[$cgr_no] . '-0' . $i;
			$definition_for{$hostname} = {
				new   => 'Server',
				using => {
					id               => $id++,
					host_name        => $hostname,
					domain_name      => $cache_groups[$cgr_no] . '.kabletown.net',
					tcp_port         => 80,
					xmpp_id          => $hostname . '-dummyxmpp',
					xmpp_passwd      => 'X',
					interface_name   => 'bond0',
					ip_address       => $net1 . ( $i + 2 ),
					ip_netmask       => '255.255.255.0',
					ip_gateway       => $net1 . '1',
					ip6_address      => $net3 . ( $i + 2 ) . '/64',
					ip6_gateway      => $net3 . '1',
					interface_mtu    => 9000,
					rack             => 'RR 119.02',
					mgmt_ip_address  => '',
					mgmt_ip_netmask  => '',
					mgmt_ip_gateway  => '',
					ilo_ip_address   => $net2 . ( $i + 6 ),
					ilo_ip_netmask   => '255.255.255.0',
					ilo_ip_gateway   => $net2 . '1',
					ilo_username     => '',
					ilo_password     => '',
					router_host_name => 'rtr-' . $cache_groups[$cgr_no] . '.kabletown.net',
					router_port_name => $cgr_no,
					type             => 1,
					status           => 3,
					profile          => $profile_id,
					cdn_id           => $cdn_id,
					cachegroup       => ( 91 + $cgr_no ),
					phys_location    => ( $site + 1 ),
				}
			};
		}
		$cgr_no++;
		$site++;
	}

	# MIDS  - 8 per cache group
	@cache_groups = ( 'east', 'west' );
	@profiles     = ( 31,     30 );
	$cdn_profiles->{"31"} = 1;
	$cdn_profiles->{"30"} = 2;

	$cgr_no       = 0;
	$site         = 0;
	foreach my $cg (@cache_groups) {
		my $net1 = '10.11.' . $cgr_no . '.';
		my $net2 = '172.17.' . $cgr_no . '.';
		my $net3 = '2033:D0D1:3300::' . $cgr_no . ":";

		foreach my $i ( 0 .. 8 ) {
			if ( $i == 4 ) {    # half of each cg is in a site
				$site++;
			}
			my $hostname = 'atsmid-' . $cache_groups[$cgr_no] . '-0' . $i;

			my $profile_id = $profiles[ ( $i % 2 == 0 ? 1 : 0 ) ];
			my $cdn_id = $cdn_profiles->{$profile_id};

			$definition_for{$hostname} = {
				new   => 'Server',
				using => {
					id               => $id++,
					host_name        => $hostname,
					domain_name      => $cache_groups[$cgr_no] . '.kabletown.net',
					tcp_port         => 80,
					xmpp_id          => $hostname . '-dummyxmpp',
					xmpp_passwd      => 'X',
					interface_name   => 'bond0',
					ip_address       => $net1 . ( $i + 1 ),
					ip_netmask       => '255.255.255.0',
					ip_gateway       => $net1 . '1',
					ip6_address      => $net3 . ( $i + 2 ) . '/64',
					ip6_gateway      => $net3 . '1',
					interface_mtu    => 9000,
					rack             => 'RR 119.02',
					mgmt_ip_address  => '',
					mgmt_ip_netmask  => '',
					mgmt_ip_gateway  => '',
					ilo_ip_address   => $net2 . ( $i + 6 ),
					ilo_ip_netmask   => '255.255.255.0',
					ilo_ip_gateway   => $net2 . '1',
					ilo_username     => '',
					ilo_password     => '',
					router_host_name => 'rtr-' . $cache_groups[$cgr_no] . '.kabletown.net',
					router_port_name => $cgr_no,
					type             => 2,
					status           => 2,
					profile          => $profile_id,
					cdn_id           => $cdn_id,
					cachegroup       => ( $cgr_no + 1 ),
					phys_location    => ( $site + 1 ),
				}
			};
		}
		$cgr_no++;
		$site++;
	}

	# traffic routers
	my $hostname = 'trtr-clw-01';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'clw.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.39.39',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.39.1',
			ip6_address    => '2033:D0D1:3300::333/64',
			ip6_gateway    => '2033:D0D1:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 4,
			status         => 2,
			profile        => 5,
			cdn_id         => 1,
			cachegroup     => 3,
			phys_location  => 100,
		}
	};
	$hostname = 'trtr-clw-02';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'clw.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.39.49',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.39.1',
			ip6_address    => '2033:D0D1:3300::334/64',
			ip6_gateway    => '2033:D0D1:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 4,
			status         => 2,
			profile        => 8,
			cdn_id         => 2,
			cachegroup     => 3,
			phys_location  => 101,
		}
	};
	$hostname = 'trtr-cle-01';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'cle.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.99.39',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.99.1',
			ip6_address    => '2033:D0D1:3300::335/64',
			ip6_gateway    => '2033:D0D1:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 4,
			status         => 2,
			profile        => 5,
			cdn_id         => 1,
			cachegroup     => 5,
			phys_location  => 100,
		}
	};
	$hostname = 'trtr-cle-02';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'cle.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.99.49',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.99.1',
			ip6_address    => '2033:D0D1:3300::336/64',
			ip6_gateway    => '2033:D0D1:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 4,
			status         => 2,
			profile        => 8,
			cdn_id         => 2,
			cachegroup     => 5,
			phys_location  => 101,
		}
	};

	# traffic monitors
	$hostname = 'trmon-clw-01';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'clw.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.29.39',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.29.1',
			ip6_address    => '2033:D021:3300::333/64',
			ip6_gateway    => '2033:D021:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 15,
			status         => 2,
			profile        => 11,
			cdn_id         => 1,
			cachegroup     => 3,
			phys_location  => 100,
		}
	};
	$hostname = 'trmon-clw-02';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'clw.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.29.49',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.29.1',
			ip6_address    => '2033:D021:3300::334/64',
			ip6_gateway    => '2033:D021:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 15,
			status         => 2,
			profile        => 12,
			cdn_id         => 2,
			cachegroup     => 3,
			phys_location  => 101,
		}
	};
	$hostname = 'trmon-cle-01';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'cle.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.19.39',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.19.1',
			ip6_address    => '2033:D011:3300::335/64',
			ip6_gateway    => '2033:D011:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 15,
			status         => 2,
			profile        => 11,
			cdn_id         => 1,
			cachegroup     => 5,
			phys_location  => 100,
		}
	};
	$hostname = 'trmon-cle-02';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id             => $id++,
			host_name      => $hostname,
			domain_name    => 'cle.kabletown.net',
			tcp_port       => 80,
			xmpp_id        => $hostname . '-dummyxmpp',
			xmpp_passwd    => 'X',
			interface_name => 'bond0',
			ip_address     => '172.39.19.49',
			ip_netmask     => '255.255.255.0',
			ip_gateway     => '172.39.19.1',
			ip6_address    => '2033:D011:3300::336/64',
			ip6_gateway    => '2033:D011:3300::1',
			interface_mtu  => 9000,
			rack           => 'RR 119.02',
			type           => 15,
			status         => 2,
			profile        => 12,
			cdn_id         => 2,
			cachegroup     => 5,
			phys_location  => 101,
		}
	};
	$hostname = 'riak1';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id               => $id++,
			host_name        => $hostname,
			domain_name      => 'kabletown.net',
			tcp_port         => 8088,
			xmpp_id          => '',
			xmpp_passwd      => '',
			interface_name   => 'eth1',
			ip_address       => '127.0.0.5',
			ip_netmask       => '255.255.252.0',
			ip_gateway       => '127.0.0.5',
			interface_mtu    => 1500,
			rack             => 'RR 119.02',
			mgmt_ip_address  => '',
			mgmt_ip_netmask  => '',
			mgmt_ip_gateway  => '',
			ilo_ip_address   => '',
			ilo_ip_netmask   => '',
			ilo_ip_gateway   => '',
			ilo_username     => '',
			ilo_password     => '',
			router_host_name => '',
			router_port_name => '',
			type             => 10,
			status           => 2,
			profile          => 47,
			cdn_id           => 2,
			cachegroup       => 1,
			phys_location    => 1,
		},
	},

	$id = 1000;
	$hostname = 'org1';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id               => $id++,
			host_name        => $hostname,
			domain_name      => 'kabletown.net',
			tcp_port         => 80,
			xmpp_id          => '',
			xmpp_passwd      => '',
			interface_name   => 'eth1',
			ip_address       => '10.11.10.2',
			ip_netmask       => '255.255.252.0',
			ip_gateway       => '10.11.10.1',
			interface_mtu    => 1500,
			rack             => 'RR 119.02',
			mgmt_ip_address  => '',
			mgmt_ip_netmask  => '',
			mgmt_ip_gateway  => '',
			ilo_ip_address   => '',
			ilo_ip_netmask   => '',
			ilo_ip_gateway   => '',
			ilo_username     => '',
			ilo_password     => '',
			router_host_name => '',
			router_port_name => '',
			type             => 36,
			status           => 2,
			profile          => 48,
			cdn_id           => 1,
			cachegroup       => 101,
			phys_location    => 1,
		},
	},
	$hostname = 'org2';
	$definition_for{$hostname} = {
		new   => 'Server',
		using => {
			id               => $id++,
			host_name        => $hostname,
			domain_name      => 'kabletown.net',
			tcp_port         => 80,
			xmpp_id          => '',
			xmpp_passwd      => '',
			interface_name   => 'eth1',
			ip_address       => '10.11.12.2',
			ip_netmask       => '255.255.252.0',
			ip_gateway       => '10.11.12.1',
			interface_mtu    => 1500,
			rack             => 'RR 119.02',
			mgmt_ip_address  => '',
			mgmt_ip_netmask  => '',
			mgmt_ip_gateway  => '',
			ilo_ip_address   => '',
			ilo_ip_netmask   => '',
			ilo_ip_gateway   => '',
			ilo_username     => '',
			ilo_password     => '',
			router_host_name => '',
			router_port_name => '',
			type             => 36,
			status           => 2,
			profile          => 49,
			cdn_id           => 1,
			cachegroup       => 102,
			phys_location    => 1,
		},
	},
}

sub name {
	return "Server";
}

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;

1;
