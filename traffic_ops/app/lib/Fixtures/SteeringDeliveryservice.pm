package Fixtures::SteeringDeliveryservice;
#
# Copyright 2016 Comcast Cable Communications Management, LLC
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
#
#
use strict;
use warnings FATAL => 'all';

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
    r1 => {
        new => 'Regex',
        using => {
            id => 21001,
            pattern => '.*/force-to-one/.*',
            type => 987,
        },
    },
    r2 => {
        new => 'Regex',
        using => {
            id => 21002,
            pattern => '.*/force-to-one-also/.*',
            type => 987,
        },
    },
    r3 => {
        new => 'Regex',
        using => {
            id => 21003,
            pattern => '.*/use-three/.*',
            type => 987,
        },
    },
    r4 => {
        new => 'Regex',
        using => {
            id => 21004,
            pattern => '.*/go-to-four/.*',
            type => 987,
        },
    },
    steering_ds1 => {
        new   => 'Deliveryservice',
        using => {
            id                    => 10001,
            xml_id                => 'steering-ds1',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => '',
            dns_bypass_ttl        => undef,
            ccr_dns_ttl           => 3600,
            global_max_mbps       => 0,
            global_max_tps        => 0,
            long_desc             => 'steering-ds1 long_desc',
            long_desc_1           => 'steering-ds1 long_desc_1',
            long_desc_2           => 'steering-ds1 long_desc_2',
            max_dns_answers       => 0,
            protocol              => 0,
            org_server_fqdn       => 'http://steering-ds1.edge',
            info_url              => 'http://steering-ds1.edge/info_url.html',
            miss_lat              => '41.881944',
            miss_long             => '-87.627778',
            check_path            => '/crossdomain.xml',
            type                  => 8,
            profile               => 3,
            cdn_id                => 1,
            ipv6_routing_enabled  => 1,
            protocol              => 1,
            display_name          => 'steering-ds1-displayname',
            initial_dispersion    => 1,
            regional_geo_blocking => 1,
        },
    },
    steering_ds2 => {
        new   => 'Deliveryservice',
        using => {
            id                    => 10002,
            xml_id                => 'steering-ds2',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => '',
            dns_bypass_ttl        => undef,
            ccr_dns_ttl           => 3600,
            global_max_mbps       => 0,
            global_max_tps        => 0,
            long_desc             => 'steering-ds2 long_desc',
            long_desc_1           => 'steering-ds2 long_desc_1',
            long_desc_2           => 'steering-ds2 long_desc_2',
            max_dns_answers       => 0,
            protocol              => 0,
            org_server_fqdn       => 'http://steering-ds2.edge',
            info_url              => 'http://steering-ds2.edge/info_url.html',
            miss_lat              => '41.881944',
            miss_long             => '-87.627778',
            check_path            => '/crossdomain.xml',
            type                  => 8,
            profile               => 3,
            cdn_id                => 1,
            ipv6_routing_enabled  => 1,
            protocol              => 1,
            display_name          => 'steering-ds2-displayname',
            initial_dispersion    => 1,
            regional_geo_blocking => 1,
        },
    },
    target_ds1 => {
        new   => 'Deliveryservice',
        using => {
            id                    => 20001,
            xml_id                => 'target-ds1',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => '',
            dns_bypass_ttl        => undef,
            ccr_dns_ttl           => 3600,
            global_max_mbps       => 0,
            global_max_tps        => 0,
            long_desc             => 'target-ds1 long_desc',
            long_desc_1           => 'target-ds1 long_desc_1',
            long_desc_2           => 'target-ds1 long_desc_2',
            max_dns_answers       => 0,
            protocol              => 0,
            org_server_fqdn       => 'http://target-ds1.edge',
            info_url              => 'http://target-ds1.edge/info_url.html',
            miss_lat              => '41.881944',
            miss_long             => '-87.627778',
            check_path            => '/crossdomain.xml',
            type                  => 8,
            profile               => 3,
            cdn_id                => 1,
            ipv6_routing_enabled  => 1,
            protocol              => 1,
            display_name          => 'target-ds1-displayname',
            initial_dispersion    => 1,
            regional_geo_blocking => 1,
        },
    },
    target_ds2 => {
        new   => 'Deliveryservice',
        using => {
            id                    => 20002,
            xml_id                => 'target-ds2',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => '',
            dns_bypass_ttl        => undef,
            ccr_dns_ttl           => 3600,
            global_max_mbps       => 0,
            global_max_tps        => 0,
            long_desc             => 'target-ds2 long_desc',
            long_desc_1           => 'target-ds2 long_desc_1',
            long_desc_2           => 'target-ds2 long_desc_2',
            max_dns_answers       => 0,
            protocol              => 0,
            org_server_fqdn       => 'http://target-ds2.edge',
            info_url              => 'http://target-ds2.edge/info_url.html',
            miss_lat              => '41.881944',
            miss_long             => '-87.627778',
            check_path            => '/crossdomain.xml',
            type                  => 8,
            profile               => 3,
            cdn_id                => 1,
            ipv6_routing_enabled  => 1,
            protocol              => 1,
            display_name          => 'target-ds2-displayname',
            initial_dispersion    => 1,
            regional_geo_blocking => 1,
        },
    },
    target_ds3 => {
        new   => 'Deliveryservice',
        using => {
            id                    => 20003,
            xml_id                => 'target-ds3',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => '',
            dns_bypass_ttl        => undef,
            ccr_dns_ttl           => 3600,
            global_max_mbps       => 0,
            global_max_tps        => 0,
            long_desc             => 'target-ds3 long_desc',
            long_desc_1           => 'target-ds3 long_desc_1',
            long_desc_2           => 'target-ds3 long_desc_2',
            max_dns_answers       => 0,
            protocol              => 0,
            org_server_fqdn       => 'http://target-ds3.edge',
            info_url              => 'http://target-ds3.edge/info_url.html',
            miss_lat              => '41.881944',
            miss_long             => '-87.627778',
            check_path            => '/crossdomain.xml',
            type                  => 8,
            profile               => 3,
            cdn_id                => 1,
            ipv6_routing_enabled  => 1,
            protocol              => 1,
            display_name          => 'target-ds3-displayname',
            initial_dispersion    => 1,
            regional_geo_blocking => 1,
        },
    },
    target_ds4 => {
        new   => 'Deliveryservice',
        using => {
            id                    => 20004,
            xml_id                => 'target-ds4',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => '',
            dns_bypass_ttl        => undef,
            ccr_dns_ttl           => 3600,
            global_max_mbps       => 0,
            global_max_tps        => 0,
            long_desc             => 'target-ds4 long_desc',
            long_desc_1           => 'target-ds4 long_desc_1',
            long_desc_2           => 'target-ds4 long_desc_2',
            max_dns_answers       => 0,
            protocol              => 0,
            org_server_fqdn       => 'http://target-ds4.edge',
            info_url              => 'http://target-ds4.edge/info_url.html',
            miss_lat              => '41.881944',
            miss_long             => '-87.627778',
            check_path            => '/crossdomain.xml',
            type                  => 8,
            profile               => 3,
            cdn_id                => 1,
            ipv6_routing_enabled  => 1,
            protocol              => 1,
            display_name          => 'target-ds4-displayname',
            initial_dispersion    => 1,
            regional_geo_blocking => 1,
        },
    },
);

sub get_definition {
    my ( $self, $name ) = @_;
    return $definition_for{$name};
}

sub all_fixture_names {
    return keys %definition_for;
}

__PACKAGE__->meta->make_immutable;

1;