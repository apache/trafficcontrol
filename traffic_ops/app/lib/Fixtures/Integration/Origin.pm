package Fixtures::Integration::Origin;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;

my %definition_for = (
	## id => 1
	'0' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin1',
            fqdn                  => 'cdl.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 1,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
	## id => 2
	'1' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin2',
            fqdn                  => 'games.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 2,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
	## id => 3
	'2' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin3',
            fqdn                  => 'images.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 3,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
	## id => 4
	'3' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin4',
            fqdn                  => 'movies.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 4,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
	## id => 5
	'4' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin5',
            fqdn                  => 'movies.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 5,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
	## id => 6
	'5' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin6',
            fqdn                  => 'games.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 6,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
	## id => 7
	'6' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin7',
            fqdn                  => 'national-tv.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 7,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
	## id => 8
	'7' => {
		new      => 'Origin',
		using => {
            name                  => 'test-origin8',
            fqdn                  => 'cc.origin.kabletown.net',
            protocol              => 'http',
            is_primary            => 1,
            port                  => undef,
            ip_address            => undef,
            ip6_address           => undef,
            deliveryservice       => 8,
            coordinate            => undef,
            profile               => undef,
            cachegroup            => undef,
            tenant                => undef,
		},
	},
);

sub name {
	return "Origin";
}

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{name} cmp $definition_for{$b}{using}{name} } keys %definition_for);
}
__PACKAGE__->meta->make_immutable;
1;
