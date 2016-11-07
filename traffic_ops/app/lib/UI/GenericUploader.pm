package UI::GenericUploader;
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
use UI::Server;
use JSON;

sub generic {
	my $self = shift;

	my $match_string = $self->param('matchstring');

	my $location;
	$rs = $self->db->resultset('Location')->search();
	while ( my $row = $rs->next ) {
		$location->{ $row->name } = $row->short_name;
	}

	my $json = JSON->new->allow_nonref;

	my $serverdata = Server::getserverdata($self);
	my $dataserver = $json->pretty->encode($serverdata);

	$self->stash(
		location    => $location,
		graph_page  => 1,
		matchstring => $match_string,
		dataserver  => $dataserver,
	);

	&navbarpage($self);
}

1;
