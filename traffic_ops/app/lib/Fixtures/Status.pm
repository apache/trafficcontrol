package Fixtures::Status;
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
	status_offline => {
		new   => 'Status',
		using => {
			id   => 1,
			name => 'OFFLINE',
			description =>
				'Edge: Puts server in CCR config file in this state, but CCR will never route traffic to it. Mid: Server will not be included in parent.config files for its edge caches',
		},
	},
	status_online => {
		new   => 'Status',
		using => {
			id   => 2,
			name => 'ONLINE',
			description =>
				'Edge: Puts server in CCR config file in this state, and CCR will always route traffic to it. Mid: Server will be included in parent.config files for its edges',
		},
	},
	status_reported => {
		new   => 'Status',
		using => {
			id          => 3,
			name        => 'REPORTED',
			description => 'Edge: Puts server in CCR config file in this state, and CCR will adhere to the health protocol. Mid: N/A for now',
		},
	},
	status_admin_down => {
		new   => 'Status',
		using => {
			id   => 4,
			name => 'ADMIN_DOWN',
			description =>
				'Temporary down. Edge: XMPP client will send status OFFLINE to CCR, otherwise similar to REPORTED. Mid: Server will not be included in parent.config files for its edge caches',
		},
	},
	status_ccr_ignore => {
		new   => 'Status',
		using => {
			id          => 5,
			name        => 'CCR_IGNORE',
			description => 'Edge: 12M will not include caches in this state in CCR config files. Mid: N/A for now',
		},
	},
    status_pre_prod => {
        new   => 'Status',
        using => {
            id          => 6,
            name        => 'PRE_PROD',
            description => 'Pre Production. Not active in any configuration.',
        },
    },
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db name to guarantee insertion order
	return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{name} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
