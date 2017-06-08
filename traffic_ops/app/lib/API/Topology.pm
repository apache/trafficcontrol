package API::Topology;
#
##
## Licensed under the Apache License, Version 2.0 (the "License");
## you may not use this file except in compliance with the License.
## You may obtain a copy of the License at
##
##     http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software
## distributed under the License is distributed on an "AS IS" BASIS,
## WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and
## limitations under the License.
##
##
##
#
## JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
#


use Mojo::Base 'Mojolicious::Controller';
use JSON;
use MojoPlugins::Response;
use UI::Utils;
use UI::Topology;
use Data::Dumper;

sub SnapshotCRConfig {
    my $self = shift;
    my $cdn_id = $self->param('id');
    my $cdn_name = $self->param('cdn_name');
    my $cdn;

    if ( !&is_oper($self) ) {
        return $self->forbidden("You must be an ADMIN or OPER to perform this operation!");
    }

    if ( defined $cdn_id ) {
        $cdn = $self->db->resultset("Cdn")->find( { id => $cdn_id } );
        $cdn_name = $cdn->name if defined $cdn;
    }

    if ( !defined $cdn ) {
        $cdn = $self->db->resultset('Cdn')->find( { name => $cdn_name } );
        if ( !defined($cdn) ) {
            return $self->not_found();
        }
    }

    my @cdn_names = $self->db->resultset('Server')->search({ 'type.name' => 'EDGE' }, { prefetch => [ 'cdn', 'type' ], group_by => 'cdn.name' } )->get_column('cdn.name')->all();
    my $num = grep /^$cdn_name$/, @cdn_names;
    if ($num <= 0) {
        return $self->alert("CDN_name [" . $cdn_name. "] is not found in edge server cdn");
    }

    my $json = &UI::Topology::gen_crconfig_json($self, $cdn_name);
    &UI::Topology::write_crconfig_json_to_db($self, $cdn_name, $json);
    &UI::Utils::log($self, "Snapshot of CRConfig performed for $cdn_name", "APICHANGE");
    return $self->success("SUCCESS");
}

sub get_snapshot {
    my $self        = shift;
    my $cdn_name    = $self->param('name');

    if ( !&is_oper($self) ) {
        return $self->forbidden();
    }

    my $cdn = $self->db->resultset('Cdn')->find( { name => $cdn_name } );
    if ( !defined($cdn) ) {
        return $self->not_found();
    }

    my $snapshot = $self->db->resultset('Snapshot')->search( { cdn => $cdn_name } )->get_column('content')->single();
    if ( !defined($snapshot) ) {
        return $self->success( {} );
    }

    $self->success( decode_json($snapshot) );
}

sub get_new_snapshot {
    my $self        = shift;
    my $cdn_name    = $self->param('name');

    if ( !&is_oper($self) ) {
        return $self->forbidden();
    }

    my $cdn = $self->db->resultset('Cdn')->find( { name => $cdn_name } );
    if ( !defined($cdn) ) {
        return $self->not_found();
    }

    my $json = &UI::Topology::gen_crconfig_json($self, $cdn_name);

    $self->success( $json );
}

1;
