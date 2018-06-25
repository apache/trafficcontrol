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
package Test::TestHelper;

use strict;
use warnings;
use Test::More;
use Test::Mojo;
use Moose;
use Schema;

use Utils::Tenant;
use Fixtures::Cdn;
use Fixtures::Deliveryservice;
use Fixtures::Origin;
use Fixtures::DeliveryserviceTmuser;
use Fixtures::Asn;
use Fixtures::Cachegroup;
use Fixtures::Coordinate;
use Fixtures::EdgeCachegroup;
use Fixtures::Log;
use Fixtures::Job;
use Fixtures::JobAgent;
use Fixtures::JobStatus;
use Fixtures::Profile;
use Fixtures::Parameter;
use Fixtures::ProfileParameter;
use Fixtures::Role;
use Fixtures::Server;
use Fixtures::Status;
use Fixtures::Tenant;
use Fixtures::TmUser;
use Fixtures::Type;
use Fixtures::Division;
use Fixtures::Region;
use Fixtures::PhysLocation;
use Fixtures::Regex;
use Fixtures::DeliveryserviceRegex;
use Fixtures::DeliveryserviceServer;

use constant ADMIN_USER          => 'admin';
use constant ADMIN_USER_PASSWORD => 'password';

use constant PORTAL_USER          => 'portal';
use constant PORTAL_USER_PASSWORD => 'password';

use constant FEDERATION_USER          => 'federation';
use constant FEDERATION_USER_PASSWORD => 'password';

use constant CODEBIG_USER     => 'codebig';
use constant CODEBIG_PASSWORD => 'password';

use constant STEERING_USER_1 => 'steering1';
use constant STEERING_PASSWORD_1 => 'password';

use constant STEERING_USER_2 => 'steering2';
use constant STEERING_PASSWORD_2 => 'password';

use constant ADMIN_ROOT_USER          => 'admin-root';
use constant ADMIN_ROOT_USER_PASSWORD => 'password';

use constant PORTAL_ROOT_USER          => 'portal-root';
use constant PORTAL_ROOT_USER_PASSWORD => 'password';

use constant READ_ONLY_ROOT_USER          => 'read-only-root';
use constant READ_ONLY_ROOT_USER_PASSWORD => 'password';

sub load_all_fixtures {
	my $self    = shift;
	my $fixture = shift;

	my @fixture_names = $fixture->all_fixture_names;
	foreach my $fixture_name (@fixture_names) {
		$fixture->load($fixture_name);

		#ok $fixture->load($fixture_name), 'Does the ' . $fixture_name . ' load?';
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

	$self->load_all_fixtures( Fixtures::Tenant->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Cdn->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Role->new($schema_values) );
	$self->load_all_fixtures( Fixtures::TmUser->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Status->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Parameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Profile->new($schema_values) );
	$self->load_all_fixtures( Fixtures::ProfileParameter->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Type->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Coordinate->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Cachegroup->new($schema_values) );
	$self->load_all_fixtures( Fixtures::EdgeCachegroup->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Division->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Region->new($schema_values) );
	$self->load_all_fixtures( Fixtures::PhysLocation->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Server->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Asn->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Deliveryservice->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Origin->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Regex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::DeliveryserviceRegex->new($schema_values) );
	$self->load_all_fixtures( Fixtures::DeliveryserviceTmuser->new($schema_values) );
	$self->load_all_fixtures( Fixtures::DeliveryserviceServer->new($schema_values) );
	$self->load_all_fixtures( Fixtures::JobStatus->new($schema_values) );
	$self->load_all_fixtures( Fixtures::JobAgent->new($schema_values) );
	$self->load_all_fixtures( Fixtures::Job->new($schema_values) );
}

sub unload_core_data {
	my $self   = shift;
	my $schema = shift;

	$self->teardown($schema, 'ToExtension');
	$self->teardown($schema, 'Staticdnsentry');
	$self->teardown($schema, 'Job');
	$self->teardown($schema, 'JobAgent');
	$self->teardown($schema, 'JobStatus');
	$self->teardown($schema, 'Log');
	$self->teardown($schema, 'Asn');
	$self->teardown($schema, 'DeliveryserviceTmuser');
	$self->teardown($schema, 'TmUser');
	$self->teardown($schema, 'Role');
	$self->teardown($schema, 'DeliveryserviceRegex');
	$self->teardown($schema, 'Regex');
	$self->teardown($schema, 'DeliveryserviceServer');
	$self->teardown($schema, 'Origin');
	$self->teardown($schema, 'Deliveryservice');
	$self->teardown($schema, 'Server');
	$self->teardown($schema, 'PhysLocation');
	$self->teardown($schema, 'Region');
	$self->teardown($schema, 'Division');
	$self->teardown_cachegroup();
	$self->teardown($schema, 'Coordinate');
	$self->teardown($schema, 'Profile');
	$self->teardown($schema, 'Parameter');
	$self->teardown($schema, 'ProfileParameter');
	$self->teardown($schema, 'Type');
	$self->teardown($schema, 'Status');
	$self->teardown($schema, 'Snapshot');
	$self->teardown($schema, 'Cdn');
	$self->teardown($schema, 'Tenant');
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

# Tearing down the Cachegroup table requires deleting them in a specific order, because
# of the 'parent_cachegroup_id' and nested references.
sub teardown_cachegroup {
	my $self   = shift;

	my $dbh    = Schema->database_handle;
	my $cg = $dbh->prepare("TRUNCATE TABLE cachegroup CASCADE;");
	$cg->execute();
	$cg->finish();
	$dbh->disconnect;
}

1;
