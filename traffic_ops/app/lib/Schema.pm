use utf8;
package Schema;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

use strict;
use warnings;

use base 'DBIx::Class::Schema';

__PACKAGE__->load_namespaces;


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2015-05-21 13:27:11
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:+I93Laz5+yCNfNrzmlDSow
#
use Cwd;
use JSON;
use Utils::JsonConfig;
use DBI;

sub database_handle {
	my $self    = shift;
	my $db_info = $self->get_dbinfo();
	our $user = $db_info->{user};
	our $pass = $db_info->{password};
	our $dsn  = $self->get_dsn();

	return DBI->connect( $dsn, $user, $pass, { AutoCommit => 1 } );
}

sub connect_to_database {
	my $self = shift;

	my $db_info = $self->get_dbinfo();
	our $user = $db_info->{user};
	our $pass = $db_info->{password};
	our $dsn  = $self->get_dsn();
	return $self->connect( $dsn, $user, $pass, { AutoCommit => 1 } );
}

sub get_dsn {
	my $self = shift;

	my $db_info = $self->get_dbinfo();
	our $dbname   = $db_info->{dbname};
	our $hostname = $db_info->{hostname};
	our $port     = $db_info->{port};
	our $type     = $db_info->{type};
	# add `ssl: true` to database.conf to require ssl
	our $reqssl   = $db_info->{ssl} ? 'require' : 'disable';
	our $dsn      = "DBI:$type:database=$dbname;host=$hostname;port=$port;sslmode=$reqssl;application_name=trafficops-perl";
}

sub get_dbinfo {
	local $/;    #Enable 'slurp' mode

	my $mode = $ENV{MOJO_MODE};
	my $dbconf;
	if ( defined($mode) ) {
		$dbconf = "conf/$mode/database.conf";
	}
	else {
		$dbconf = 'conf/development/database.conf';
	}

	#print( "Using database.conf: " . $dbconf . "\n" );
	return Utils::JsonConfig->new($dbconf);
}

# You can replace this text with custom code or comments, and it will be preserved on regeneration
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
1;
