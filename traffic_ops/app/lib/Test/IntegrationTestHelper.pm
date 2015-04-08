use utf8;
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
#
#

# Test Helper to allow for simpler Test Cases.
package Test::IntegrationTestHelper;

use strict;
use warnings;
use Test::More;
use Test::Mojo;
use Moose;

use Fixtures::Integration::Deliveryservice;
use Fixtures::Integration::OtherCacheGroup;
use Fixtures::Integration::EdgeCacheGroup;
use Fixtures::Integration::Profile;
use Fixtures::Integration::Parameter;
use Fixtures::Integration::ProfileParameter;
use Fixtures::Integration::Role;
use Fixtures::Integration::Server;
use Fixtures::Integration::Asn;
use Fixtures::Integration::Status;
use Fixtures::Integration::TmUser;
use Fixtures::Integration::Type;
use Fixtures::Integration::ToExtension;
use Fixtures::Integration::Division;
use Fixtures::Integration::Region;
use Fixtures::Integration::PhysLocation;
use Fixtures::Integration::Regex;
use Fixtures::Integration::DeliveryserviceRegex;

use Data::Dumper;

sub load_all_fixtures {
	my $self    = shift;
	my $fixture = shift;

	diag "  " . $fixture->name() . "...\n";
	if ( $fixture->can("gen_data") ) {
		$fixture->gen_data();
	}
	my @fixture_names = $fixture->all_fixture_names;
	foreach my $fixture_name (@fixture_names) {
		$fixture->load($fixture_name);

		#ok $fixture->load($fixture_name), 'Does the ' . $fixture_name . ' load?';
	}
}

sub teardown {
	my $self       = shift;
	my $schema     = shift;
	my $table_name = shift;
	$schema->resultset($table_name)->delete_all;

	#ok $schema->resultset($table_name)->delete_all, 'Does the ' . $table_name . ' teardown?';
}

sub link_servers {
	my $self   = shift;
	my $schema = shift;

	my $rs = $schema->resultset('Server')->search( {} );
	while ( my $server = $rs->next ) {
		my $i = $schema->resultset('Servercheck')->create( { server => $server->id } );
		$i->insert();
		if ( $server->type->name eq 'EDGE' ) {
			if ( $server->profile->name =~ /CDN1/ ) {
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 1 } );
				$i->insert();
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 2 } );
				$i->insert();
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 3 } );
				$i->insert();
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 4 } );
				$i->insert();
			}
			else {
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 11 } );
				$i->insert();
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 12 } );
				$i->insert();
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 13 } );
				$i->insert();
				$i = $schema->resultset('DeliveryserviceServer')->create( { server => $server->id, deliveryservice => 14 } );
				$i->insert();
			}
		}
	}
}

sub load_core_data {
	my $self          = shift;
	my $schema        = shift;
	my $schema_values = { schema => $schema, no_transactions => 1 };
	diag "Initializing DB:";
	$self->load_all_fixtures( Fixtures::Integration::Type->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Role->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::TmUser->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Division->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Region->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::PhysLocation->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Status->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Regex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Parameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Profile->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::ProfileParameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::OtherCacheGroup->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::EdgeCacheGroup->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Asn->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Server->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Deliveryservice->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::DeliveryserviceRegex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::ToExtension->new($schema_values) );

	$self->link_servers($schema);
	diag "Done!";
}

sub delete_cachegroups {
	my $self    = shift;
	my $schema  = shift;
	my $sql     = 'IS NOT NULL';
	my $parents = $schema->resultset('Cachegroup')->search( { parent_cachegroup_id => \$sql } );
	$parents->delete;
	$schema->resultset('Cachegroup')->delete_all;
}

sub unload_core_data {
	my $self   = shift;
	my $schema = shift;
	$self->teardown( $schema, 'Job' );
	$self->teardown( $schema, 'Log' );
	$self->teardown( $schema, 'TmUser' );
	$self->teardown( $schema, 'Role' );
	$self->teardown( $schema, 'Regex' );
	$self->teardown( $schema, 'Deliveryservice' );
	$self->teardown( $schema, 'Server' );
	$self->teardown( $schema, 'Asn' );
	$self->delete_cachegroups($schema);    # cachegroups is special because it refs itself
	$self->teardown( $schema, 'Profile' );
	$self->teardown( $schema, 'Parameter' );
	$self->teardown( $schema, 'ProfileParameter' );
	$self->teardown( $schema, 'ToExtension' );
	$self->teardown( $schema, 'Type' );
	$self->teardown( $schema, 'Status' );
	$self->teardown( $schema, 'PhysLocation' );
	$self->teardown( $schema, 'Region' );
	$self->teardown( $schema, 'Division' );
}

1;
