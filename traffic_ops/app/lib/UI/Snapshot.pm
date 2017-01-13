package UI::Snapshot;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;

sub get_cdn_snapshot {
    my $self = shift;
    my $cdn_name   = $self->param('cdn_name');

    my $snapshot = $self->db->resultset('Snapshot')->search( { cdn => $cdn_name } )->get_column('content')->single();
    if ( !defined($snapshot) ) {
        return $self->not_found();
    }

    $self->res->headers->content_type("application/download");
    $self->res->headers->content_disposition("attachment; filename=\"CRConfig.json\"");
    $self->render( text => $snapshot, format => 'json' );

}

1;
