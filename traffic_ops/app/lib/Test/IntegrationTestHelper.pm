use utf8;
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

# Test Helper to allow for simpler Test Cases.
package Test::IntegrationTestHelper;

use strict;
use warnings;
use Test::More;
use Test::Mojo;
use Moose;

use Utils::Tenant;
use Fixtures::Integration::Asn;
use Fixtures::Integration::CachegroupParameter;
use Fixtures::Integration::Cachegroup;
use Fixtures::Integration::Cdn;
use Fixtures::Integration::Coordinate;
use Fixtures::Integration::Deliveryservice;
use Fixtures::Integration::DeliveryserviceRegex;
use Fixtures::Integration::DeliveryserviceServer;
use Fixtures::Integration::DeliveryserviceTmuser;
use Fixtures::Integration::Division;
use Fixtures::Integration::FederationDeliveryservice;
use Fixtures::Integration::FederationFederationResolver;
use Fixtures::Integration::Federation;
use Fixtures::Integration::FederationResolver;
use Fixtures::Integration::FederationTmuser;
use Fixtures::Integration::Hwinfo;
use Fixtures::Integration::JobAgent;
use Fixtures::Integration::Job;
use Fixtures::Integration::JobStatus;
use Fixtures::Integration::Log;
use Fixtures::Integration::Origin;
use Fixtures::Integration::Parameter;
use Fixtures::Integration::PhysLocation;
use Fixtures::Integration::ProfileParameter;
use Fixtures::Integration::Profile;
use Fixtures::Integration::Regex;
use Fixtures::Integration::Region;
use Fixtures::Integration::Role;
use Fixtures::Integration::Servercheck;
use Fixtures::Integration::Server;
use Fixtures::Integration::Staticdnsentry;
use Fixtures::Integration::StatsSummary;
use Fixtures::Integration::Status;
use Fixtures::Integration::TmUser;
use Fixtures::Integration::ToExtension;
use Fixtures::Integration::Type;

use Data::Dumper;

sub load_all_fixtures {
	my $self    = shift;
	my $fixture = shift;

	diag "  " . $fixture->name . "...\n";
	if ( $fixture->can("gen_data") ) {
		$fixture->gen_data();
	}
	if ( $fixture->name ne "Cachegroup" ) {
		my @fixture_names = $fixture->all_fixture_names;
		foreach my $fixture_name (@fixture_names) {
			$fixture->load($fixture_name);
		}
	}
	else {
		# these shenanigans are here to first add ORG(36), then MID(7) and then the rest to prevent foreign key failures
		my @fixture_names = $fixture->all_fixture_names;
		foreach my $fixture_name (@fixture_names) {
			my $cg = $fixture->get_definition($fixture_name);
			if ( $cg->{using}->{type} == 36 ) {
				$fixture->load($fixture_name);
			}
		}
		foreach my $fixture_name (@fixture_names) {
			my $cg = $fixture->get_definition($fixture_name);
			if ( $cg->{using}->{type} == 7 ) {
				$fixture->load($fixture_name);
			}
		}
		foreach my $fixture_name (@fixture_names) {
			$fixture->load($fixture_name);
		}
	}
}

sub teardown {
	my $self       = shift;
	my $schema     = shift;
	my $table_name = shift;

	if ($table_name eq 'Tenant') {
		my $tenant_utils = Utils::Tenant->new(undef, 10**9, $schema);
		my $tenants_data = $tenant_utils->create_tenants_data_from_db();
		$tenant_utils->cascade_delete_tenants_tree($tenants_data);
	}
	else {
		$schema->resultset($table_name)->delete_all;
	}
}

## For PSQL sequence to work correctly we cannot hard code
## the id number for an entry in the DB.  So we need to
## reset all primary keys (id) to 1 for consistency in the
## test cases.
sub reset_sequence_id {
	my $self   = shift;
	my $dbh    = Schema->database_handle;

	my $p = $dbh->prepare( "SELECT * FROM pg_class WHERE relkind = 'S';" );
	$p->execute();
	my $foo = $p->fetchall_arrayref( {} );
	$p->finish();


	for my $table ( @$foo ) {
		my $x = $dbh->prepare("ALTER SEQUENCE " . $table->{'relname'} . " RESTART WITH 1");
		$x->execute();
	}
}

sub load_core_data {
	my $self          = shift;
	my $schema        = shift;
	my $schema_values = { schema => $schema, no_transactions => 1 };

	$self->reset_sequence_id();

	diag "Initializing DB:";
	$self->load_all_fixtures( Fixtures::Integration::Cdn->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Type->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Role->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::TmUser->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Division->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Region->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::PhysLocation->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Status->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Coordinate->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Cachegroup->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Regex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Parameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Profile->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::ProfileParameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Asn->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Server->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Deliveryservice->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Origin->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::DeliveryserviceRegex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::DeliveryserviceServer->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::ToExtension->new($schema_values) );

	diag "Done!";
}

sub unload_core_data {
	my $self   = shift;
	my $schema = shift;

	my $dbh    = Schema->database_handle;

	# Suppress NOTICE messages for cascades
	my $nonotice = $dbh->prepare("SET client_min_messages TO WARNING;");
	$nonotice->execute();
	$nonotice->finish();
	for my $source (values %{$schema->source_registrations}) {
		if ( ! $source->isa('DBIx::Class::ResultSource::Table') ) {
			# Skip if it doesn't represent an actual table
			next;
		}
		my $table_name = $source->name;
		my $truncate = $dbh->prepare("TRUNCATE TABLE $table_name CASCADE;");
		$truncate->execute();
		$truncate->finish();
	}
}

1;
