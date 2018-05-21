package Fixtures::Origin;

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
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
    origin_cdn1 => {
    new   => 'Origin',
        using => {
            id                    => 100,
            name                  => 'test-origin1',
            fqdn                  => 'test-ds1.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 100,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    origin_cdn2 => {
        new   => 'Origin',
        using => {
            id                    => 200,
            name                  => 'test-origin2',
            fqdn                  => 'test-ds2.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 200,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    origin_cdn3 => {
        new   => 'Origin',
        using => {
            id                    => 300,
            name                  => 'test-origin3',
            fqdn                  => 'test-ds3.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 300,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    origin_cdn4 => {
        new   => 'Origin',
        using => {
            id                    => 400,
            name                  => 'test-origin4',
            fqdn                  => 'test-ds4.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 400,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    origin_dns => {
        new   => 'Origin',
        using => {
            id                    => 500,
            name                  => 'test-origin5',
            fqdn                  => 'test-ds5.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 500,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    origin_http_no_cache => {
        new   => 'Origin',
        using => {
            id                    => 600,
            name                  => 'test-origin6',
            fqdn                  => 'test-ds6.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 600,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    steering_origin1 => {
        new   => 'Origin',
        using => {
            id                    => 700,
            name                  => 'test-origin7',
            fqdn                  => 'steering-ds1.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 700,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    steering_origin2 => {
        new   => 'Origin',
        using => {
            id                    => 800,
            name                  => 'test-origin8',
            fqdn                  => 'steering-ds2.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 800,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    new_origin_steering => {
        new   => 'Origin',
        using => {
            id                    => 900,
            name                  => 'test-origin9',
            fqdn                  => 'new-steering-ds.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 900,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    target_origin1 => {
        new   => 'Origin',
        using => {
            id                    => 1000,
            name                  => 'test-origin10',
            fqdn                  => 'target-ds1.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 1000,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    target_origin2 => {
        new   => 'Origin',
        using => {
            id                    => 1100,
            name                  => 'test-origin11',
            fqdn                  => 'target-ds2.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 1100,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    target_origin3 => {
        new   => 'Origin',
        using => {
            id                    => 1200,
            name                  => 'test-origin12',
            fqdn                  => 'target-ds3.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 1200,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    target_origin4 => {
        new   => 'Origin',
        using => {
            id                    => 1300,
            name                  => 'test-origin13',
            fqdn                  => 'target-ds4.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 1300,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    origin_cdn1_root => {
        new   => 'Origin',
        using => {
            id                    => 2100,
            name                  => 'test-origin14',
            fqdn                  => 'test-ds1-root.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 2100,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => 10**9,
        },
    },
    origin_period1 => {
        new   => 'Origin',
        using => {
            id                    => 2200,
            name                  => 'test-origin15',
            fqdn                  => 'foo.bar.edge',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 2200,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
        },
    },
    
);

sub get_definition {
    my ( $self, $name ) = @_;
    return $definition_for{$name};
}

sub all_fixture_names {
    # sort by db id to guarantee insertion order
    return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{id} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
