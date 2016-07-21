package Fixtures::SteeringDeliveryservice;
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
#
#
use strict;
use warnings FATAL => 'all';

use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
    target_filter_1 => {
        new => 'Regex',
        using => {
            pattern => '.*/force-to-one/.*',
            type => 987,
        },
    },
    target_filter_1_2 => {
        new => 'Regex',
        using => {
            pattern => '.*/force-to-one-also/.*',
            type => 987,
        },
    },
    target_filter_3 => {
        new => 'Regex',
        using => {
            pattern => '.*/use-three/.*',
            type => 987,
        },
    },
    target_filter_4 => {
        new => 'Regex',
        using => {
            pattern => '.*/go-to-four/.*',
            type => 987,
        },
    },
    hr_steering_1 => {
        new => 'Regex',
        using => {
            pattern => '.*\.steering-ds1\..*',
            type => 19,
        },
    },
    hr_steering_2 => {
        new => 'Regex',
        using => {
            pattern => '.*\.steering-ds2\..*',
            type => 19,
        },
    },
    hr_target_1 => {
        new => 'Regex',
        using => {
            pattern => '.*\.target-ds1\..*',
            type => 19,
        },
    },
    hr_target_2 => {
        new => 'Regex',
        using => {
            pattern => '.*\.target-ds2\..*',
            type => 19,
        },
    },
    hr_target_3 => {
        new => 'Regex',
        using => {
            pattern => '.*\.target-ds3\..*',
            type => 19,
        },
    },
    hr_target_4 => {
        new => 'Regex',
        using => {
            pattern => '.*\.target-ds4\..*',
            type => 19,
        },
    },
    hr_new_steering => {
        new => 'Regex',
        using => {
            pattern => '.*\.new-steering-ds\..*',
            type => 19,
        },
    },
    steering_ds1 => {
        new   => 'Deliveryservice',
        using => {
            xml_id                => 'steering-ds1',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => 'hokeypokey',
            dns_bypass_ttl        => 10,
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
            xml_id                => 'steering-ds2',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => 'hokeypokey',
            dns_bypass_ttl        => 10,
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
            xml_id                => 'target-ds1',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => 'hokeypokey',
            dns_bypass_ttl        => 10,
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
            xml_id                => 'target-ds2',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => 'hokeypokey',
            dns_bypass_ttl        => 10,
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
            xml_id                => 'target-ds3',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => 'hokeypokey',
            dns_bypass_ttl        => 10,
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
            xml_id                => 'target-ds4',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => 'hokeypokey',
            dns_bypass_ttl        => 10,
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
    new_steering => {
        new   => 'Deliveryservice',
        using => {
            xml_id                => 'new-steering-ds',
            active                => 1,
            dscp                  => 40,
            signed                => 0,
            qstring_ignore        => 0,
            geo_limit             => 0,
            http_bypass_fqdn      => '',
            dns_bypass_ip         => 'hokeypokey',
            dns_bypass_ttl        => 10,
            ccr_dns_ttl           => 3600,
            global_max_mbps       => 0,
            global_max_tps        => 0,
            long_desc             => 'new-steering-ds long_desc',
            long_desc_1           => 'new-steering-ds long_desc_1',
            long_desc_2           => 'new-steering-ds long_desc_2',
            max_dns_answers       => 0,
            protocol              => 0,
            org_server_fqdn       => 'http://new-steering-ds.edge',
            info_url              => 'http://new-steering-ds.edge/info_url.html',
            miss_lat              => '41.881944',
            miss_long             => '-87.627778',
            check_path            => '/crossdomain.xml',
            type                  => 8,
            profile               => 3,
            cdn_id                => 1,
            ipv6_routing_enabled  => 1,
            protocol              => 1,
            display_name          => 'new-steering-ds-displayname',
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
