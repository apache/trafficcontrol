package Fixtures::Integration::PhysLocation;
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

my %definition_for = ();

sub gen_data {
	my @cache_groups = ( 'nyc', 'lax', 'chi', 'hou', 'phl', 'den' );
	my @regions      = ( 1,     2,     3,     3,     1,     2 );
	my @states       = ( 'NY',  'CA',  'IL',  'TX',  'PA',  'CO' );

	my $cgrno = 0;

	# each cache group has 2 phys locations
	my $counter = 1;
	foreach my $loc (@cache_groups) {
		my $site     = 1;
		my $plocname = "plocation-" . $cache_groups[$cgrno] . "-" . $site;
		$definition_for{$plocname} = {
			new   => 'PhysLocation',
			using => {
				id         => $counter,
				name       => $plocname,
				short_name => $cache_groups[$cgrno] . "-" . $site,
				address    => $counter . ' Main Street',
				city       => $cache_groups[$cgrno],
				state      => $states[$cgrno],
				zip        => '12345',
				region     => $regions[$cgrno],
			},
		};
		$counter++;
		$site                      = 2;
		$plocname                  = "plocation-" . $cache_groups[$cgrno] . "-" . $site;
		$definition_for{$plocname} = {
			new   => 'PhysLocation',
			using => {
				id         => $counter,
				name       => $plocname,
				short_name => $cache_groups[$cgrno] . "-" . $site,
				address    => $counter . ' Broadway',
				city       => $cache_groups[$cgrno],
				state      => $states[$cgrno],
				zip        => '12345',
				region     => $regions[$cgrno],
			},
		};
		$counter++;
		$cgrno++;
	}
	my $plocname = 'cloud-east';
	$definition_for{$plocname} = {
		new   => 'PhysLocation',
		using => {
			id         => 100,
			name       => $plocname,
			short_name => 'clw',
			address    => '-',
			city       => '-',
			state      => '-',
			zip        => '-',
			region     => 1,
		},
	};
	$plocname = 'cloud-west';
	$definition_for{$plocname} = {
		new   => 'PhysLocation',
		using => {
			id         => 101,
			name       => $plocname,
			short_name => 'cle',
			address    => '-',
			city       => '-',
			state      => '-',
			zip        => '-',
			region     => 2, 
		},
	};

}
sub name {
	return "PhysLocation";
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
