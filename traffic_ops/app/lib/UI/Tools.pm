package UI::Tools;

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
use UI::Topology;
use UI::Job;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use Mojo::UserAgent;
use POSIX;
use HTTP::Cookies;
use DBI;
use JSON;

sub tools {
    my $self = shift;

    &navbarpage($self);
    my %serverselect;
    my $rs_server = $self->db->resultset('Server')->search( undef, { columns => [qw/id host_name domain_name/], orderby => "host_name" } );

    while ( my $row = $rs_server->next ) {
        my $fqdn = $row->host_name . "." . $row->domain_name;
        $serverselect{$fqdn} = $row->id;
    }

    my %osversions = (
        "CentOS 6.2 Network" => "centos62",
        "CentOS 6.2 Full"    => "full-centos62",
    );

    my $rs_param = $self->db->resultset('Cdn')->search( undef, { columns => 'name' } );
    my @cdn_names;
    while ( my $row = $rs_param->next ) {
        push( @cdn_names, $row->name );
    }

    $self->stash(
        serverselect => \%serverselect,
        osversions   => \%osversions,
        cdn_names    => \@cdn_names,
    );

}

sub snapshot_crconfig {
    my $self = shift;
    &navbarpage($self);

    my @cdn_names =
        $self->db->resultset('Server')->search( { 'type.name' => { -like => 'EDGE%' } }, { prefetch => [ 'cdn', 'type' ], group_by => 'cdn.name' } )
        ->get_column('cdn.name')->all();

    $self->stash( cdn_names => \@cdn_names );
}

sub diff_crconfig_iframe {
    my $self = shift;
    &stash_role($self);
    my $cdn_name = $self->param('cdn_name');

    foreach my $cookie ( @{ $self->req->cookies } ) {
        $self->ua->cookie_jar->add(Mojo::Cookie::Response->new(name => $cookie->{'name'}, value => $cookie->{'value'}, domain => 'localhost', path => '/'));
    }
    my $resp = $self->ua->request_timeout(60)->get('https://localhost:' . $self->config->{'traffic_ops_golang'}{'port'} . '/api/1.2/cdns/' . $cdn_name . '/snapshot/new')->res;
    my $json = undef;
    my $error = undef;
    if ( $resp->code ne '200' ) {
        $error = $resp->message;
    } else {
        $json = decode_json($resp->body)->{'response'};
    }

    my ( @ds_text, @loc_text, @cs_text, @csds_text, @rascal_text, @ccr_text, @cfg_text );
    if ( defined $error ) {
        $self->flash( alertmsg => $error );
    }
    else {

        ( my $ds_text, my $loc_text, my $cs_text, my $csds_text, my $rascal_text, my $ccr_text, my $cfg_text ) =
            UI::Topology::diff_crconfig_json( $self, $json, $cdn_name );
        @ds_text     = @$ds_text;
        @loc_text    = @$loc_text;
        @cs_text     = @$cs_text;
        @csds_text   = @$csds_text;
        @rascal_text = @$rascal_text;
        @ccr_text    = @$ccr_text;
        @cfg_text    = @$cfg_text;
    }
    $self->stash(
        ds_text     => \@ds_text,
        loc_text    => \@loc_text,
        cs_text     => \@cs_text,
        csds_text   => \@csds_text,
        rascal_text => \@rascal_text,
        ccr_text    => \@ccr_text,
        cfg_text    => \@cfg_text,
        cdn         => $cdn_name,
        fbox_layout => 1,
        crconfig_db => $json
    );
}

sub flash_and_close {
    my $self = shift;
		my $msg = $self->param('msg');
		$self->flash( alertmsg => $msg );
		return $self->redirect_to('/utils/close_fancybox');
}

sub queue_updates {
    my $self = shift;
    &stash_role($self);

    my @cdns =
        $self->db->resultset('Server')
        ->search( { 'type.name' => [ { -like => 'EDGE%' }, { -like => 'MID%' } ] }, { prefetch => [ 'cdn', 'type' ], group_by => 'cdn.name' } )
        ->get_column('cdn.name')->all();
    $self->stash( cdns => \@cdns );

    my @cachegroups = $self->db->resultset('Cachegroup')->search( undef, { order_by => "name" } )->get_column('name')->all;
    $self->stash( cachegroups => \@cachegroups );

    &navbarpage($self);
}

sub db_dump {
    my $self = shift;

    my ( $sec, $min, $hour, $day, $month, $year ) = (localtime)[ 0, 1, 2, 3, 4, 5 ];
    $month = sprintf '%02d', $month + 1;
    $day   = sprintf '%02d', $day;
    $hour  = sprintf '%02d', $hour;
    $min   = sprintf '%02d', $min;
    $sec   = sprintf '%02d', $sec;
    $year += 1900;
    my $host = `hostname`;
    chomp($host);

    my $extension = ".pg_dump";
    my $filename = "to-backup-" . $host . "-" . $year . $month . $day . $hour . $min . $sec . $extension;
    $self->stash( filename => $filename );
    &stash_role($self);
    &navbarpage($self);
}

sub invalidate_content {
    my $self = shift;
    &stash_role($self);
    &navbarpage($self);
    my $id         = $self->param('id');
    my %ds         = get_delivery_services( $self, $id );
    my $ttl        = 48;
    my $start_time = strftime( "%Y-%m-%d %H:%M:%S\n", localtime(time) );
    my $regex      = "/foo/.*";
    $self->stash(
        job => {
            ttl        => $ttl,
            start_time => $start_time,
            regex      => $regex
        },
        ds          => \%ds,
        selected_ds => 'default'
    );
}

sub get_delivery_services {
    my $self = shift;
    my $id   = $self->param('id');
    my %delivery_services;
    my $deliveryServicesRs = $self->db->resultset('Deliveryservice');
    while ( my $deliveryService = $deliveryServicesRs->next ) {
        $delivery_services{ $deliveryService->xml_id } = $deliveryService->xml_id;
    }
    return %delivery_services;
}
1;
