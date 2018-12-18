package UI::GenDbDump;
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	 http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
#
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use UI::Utils;
use IO::Compress::Gzip qw(gzip $GzipError);

sub dbdump {
	my $self = shift;
	my $filename = $self->param('filename');

	if ( !&is_admin($self) ) {
		$self->internal_server_error( { Error => "Insufficient permissions for DB Dump. Admin or Operations access is required." } );	
		return;
	}

	my ($db_name, $host, $port) = $Schema::dsn =~ /:database=([^;]*);host=([^;]+);port=(\d+)/;
	my $db_user = $Schema::user;
	my $db_pass = $Schema::pass;

	my $ok = open my $fh, '-|', "PGPASSWORD=\"$db_pass\" pg_dump -b -Fc --no-owner -h $host -p $port -U $db_user -d $db_name";
	if (! $ok ) {
		$self->internal_server_error( { Error => "Error dumping database" } );	
		return;
	}

	# slurp it in..
	local $/;
	my $data = <$fh>;


	$self->res->headers->content_type("application/download");
	$self->res->headers->content_disposition( "attachment; filename=\"" . $filename . "\"" );
	$self->render( data => $data );
	close $fh;
}

1;
