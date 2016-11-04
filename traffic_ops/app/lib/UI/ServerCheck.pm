package UI::ServerCheck;
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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

sub server_check {
	my $self = shift;

	my @type_ids = &type_ids( $self, 'CHECK_EXT%' );
	my $rs_extensions =
		$self->db->resultset('ToExtension')->search( { type => { -in => \@type_ids } }, { prefetch => ['type'], order_by => ["servercheck_column_name"] } );
	my @extensions;
	while ( my $row = $rs_extensions->next ) {
		push(
			@extensions, {
				id       => $row->id,
				col_name => $row->servercheck_short_name,
				column   => $row->servercheck_column_name,
				type     => $row->type->name,
				isactive => $row->isactive
			}
		);
	}
	$self->stash( extensions => \@extensions );

	&navbarpage($self);
}

1;
