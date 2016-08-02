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

use Fixtures::Integration::Asn;
use Fixtures::Integration::CachegroupParameter;
use Fixtures::Integration::Cachegroup;
use Fixtures::Integration::Cdn;
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
use Fixtures::Integration::GooseDbVersion;
use Fixtures::Integration::Hwinfo;
use Fixtures::Integration::JobAgent;
use Fixtures::Integration::Job;
use Fixtures::Integration::JobResult;
use Fixtures::Integration::JobStatus;
use Fixtures::Integration::Log;
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
	$schema->resultset($table_name)->delete_all;
}

sub reset_sequence_id {
	my $self   = shift;
	my $dbh    = Schema->database_handle;

	my @table_names = qw(
		asn
		cachegroup
		cdn
		deliveryservice
		division
		federation
		federation_resolver
		hwinfo
		job_agent
		job_status
		log
		parameter
		phys_location
		profile
		regex
		region
		role
		server
		staticdnsentry
		stats_summary
		status
		tm_user
		type );
	foreach my $name (@table_names) {
		my $p = $dbh->prepare("ALTER SEQUENCE " . $name . "_id_seq RESTART WITH 1");
		$p->execute();
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
	$self->load_all_fixtures( Fixtures::Integration::Cachegroup->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Regex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Parameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Profile->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::ProfileParameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Asn->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Server->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::Deliveryservice->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::DeliveryserviceRegex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::DeliveryserviceServer->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Integration::ToExtension->new($schema_values) );

	diag "Done!";
}

# Tearing down the Cachegroup table requires deleting them in a specific order, because
# of the 'parent_cachegroup_id' and nested references.
sub delete_cachegroups {
	my $self   = shift;

	my $dbh    = Schema->database_handle;
	my $cg = $dbh->prepare("TRUNCATE TABLE cachegroup CASCADE;");
	$cg->execute();
	$cg->finish();
	$dbh->disconnect;
}

sub unload_core_data {
	my $self   = shift;
	my $schema = shift;
	$self->teardown( $schema, 'Job' );
	$self->teardown( $schema, 'Log' );
	$self->teardown( $schema, 'TmUser' );
	$self->teardown( $schema, 'Role' );
	$self->teardown( $schema, 'Regex' );
	$self->teardown( $schema, 'DeliveryserviceServer' );
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
	$self->teardown( $schema, 'Cdn' );
}

1;
