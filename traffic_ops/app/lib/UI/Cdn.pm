package UI::Cdn;

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
use UI::Parameter;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use UI::ConfigFiles;
use Date::Manip;
use JSON;
use Hash::Merge qw(merge);
use String::CamelCase qw(decamelize);
use DBI;
use Utils::Tenant;

# Yes or no
my %yesno = ( 0 => "no", 1 => "yes", 2 => "no" );

sub index {
    my $self = shift;

    &navbarpage($self);
}

sub add {
    my $self = shift;
    $self->stash( fbox_layout => 1, cdn_data => {} );
    &stash_role($self);
    if ( $self->stash('priv_level') < 30 ) {
        $self->stash( alertmsg => "Insufficient privileges!" );
        $self->redirect_to('/cdns');
    }
}

sub view {
    my $self = shift;
    my $mode = $self->param('mode');
    my $id   = $self->param('id');
    $self->stash( cdn_data => {} );

    my $rs_param = $self->db->resultset('Cdn')->search( { id => $id } );
    my $data = $rs_param->single;

    &stash_role($self);
    $self->stash( fbox_layout => 1, cdn_data => $data );

    if ( $mode eq "edit" and $self->stash('priv_level') > 20 ) {
        $self->render( template => 'cdn/edit' );
    }
    else {
        $self->render( template => 'cdn/view' );
    }
}

sub update {
    my $self       = shift;
    my $id         = $self->param('id');
    my $priv_level = $self->stash('priv_level');
    my $dnssec_enabled = defined($self->param('cdn_data.dnssec_enabled')) ? $self->param('cdn_data.dnssec_enabled') : 0;

    $self->stash(
        id          => $id,
        fbox_layout => 1,
        priv_level  => $priv_level,
        cdn_data    => {
            id   => $id,
            name => $self->param('cdn_data.name'),
            domain_name => $self->param('cdn_data.domain_name'),
            dnssec_enabled => $self->param('cdn_data.dnssec_enabled'),
        }
    );

    if ( !$self->isValidCdn() ) {
        return $self->render( template => 'cdn/edit' );
    }

    my $err = &check_cdn_input($self);
    if ( defined($err) ) {
        $self->flash( alertmsg => $err );
    }
    else {
        my $update = $self->db->resultset('Cdn')->find( { id => $self->param('id') } );
        $update->name( $self->param('cdn_data.name') );
        $update->domain_name( $self->param('cdn_data.domain_name') );
        $update->dnssec_enabled( $dnssec_enabled );
        $update->update();

        # if the update has failed, we don't even get here, we go to the exception page.
    }

    my $msg = "Update Cdn with name:" . $self->param('cdn_data.name') .
              "; domain_name: " . $self->param('cdn_data.domain_name') .
              "; dnssec_enabled: " . $self->param('cdn_data.dnssec_enabled');
    &log( $self, $msg, "UICHANGE" );
    $self->flash( message => "Successfully updated CDN." );
    return $self->redirect_to( '/cdn/edit/' . $id );
}

# Create
sub create {
    my $self = shift;
    my $name = $self->param('cdn_data.name');
    my $dnssec_enabled = defined($self->param('cdn_data.dnssec_enabled')) ? $self->param('cdn_data.dnssec_enabled') : 0;
    my $domain_name = $self->param('cdn_data.domain_name');
    my $data = $self->get_cdns();
    my $cdns = $data->{'cdn'};

    if ( !$self->isValidCdn() ) {
        $self->stash(
            fbox_layout => 1,
            cdn_data    => {
                name => $name,
                domain_name => $domain_name,
                dnssec_enabled => $dnssec_enabled,
            }
        );
        return $self->render('cdn/add');
    }
    if ( exists $cdns->{$name} ) {
        $self->field('cdn_data.name')->is_like( qr/^\/(?!$name\/)/i, "The name exists." );
        $self->stash(
            fbox_layout => 1,
            cdn_data    => {
                name => $name,
                domain_name => $domain_name,
                dnssec_enabled => $dnssec_enabled,
            }
        );
        return $self->render('cdn/add');
    }

    my $new_id = -1;
    my $err    = &check_cdn_input($self);
    if ( defined($err) ) {
        return $self->redirect_to( '/cdn/edit/' . $new_id );
    }
    else {
        my $insert = $self->db->resultset('Cdn')->create( { name => $name, domain_name => $domain_name, dnssec_enabled => $dnssec_enabled } );
        $insert->insert();
        $new_id = $insert->id;
    }
    if ( $new_id == -1 ) {
        my $referer = $self->req->headers->header('referer');
        return $self->redirect_to($referer);
    }
    else {
        &log( $self, "Create cdn with name:" . $self->param('cdn_data.name'), "UICHANGE" );
        $self->flash( message => "Successfully updated CDN." );
        return $self->redirect_to( '/cdn/edit/' . $new_id );
    }
}

# Delete
sub delete {
    my $self = shift;
    my $id   = $self->param('id');

    if ( !&is_admin($self) ) {
        $self->flash( alertmsg => "You must be an ADMIN to perform this operation!" );
    }
    else {
        my $server_count = $self->db->resultset('Server')->search( { cdn_id => $id } )->count();
        if ($server_count > 0) {
            $self->flash( alertmsg => "Failed to delete cdn id = $id has servers." );
            return $self->redirect_to('/close_fancybox.html');
        }
        my $ds_count = $self->db->resultset('Deliveryservice')->search( { cdn_id => $id } )->count();
        if ($ds_count > 0) {
            $self->flash( alertmsg => "Failed to delete cdn id = $id has delivery services." );
            return $self->redirect_to('/close_fancybox.html');
        }
        my $p_name = $self->db->resultset('Cdn')->search( { id => $id } )->get_column('name')->single();
        my $delete = $self->db->resultset('Cdn')->search( { id => $id } );
        $delete->delete();
        &log( $self, "Delete cdn " . $p_name, "UICHANGE" );
    }
    return $self->redirect_to('/close_fancybox.html');
}

sub get_cdns {
    my $self = shift;

    my %data;
    my %cdns;
    my $rs = $self->db->resultset('Cdn');
    while ( my $cdn = $rs->next ) {
        $cdns{ $cdn->name } = $cdn->id;
    }
    %data = ( cdns => \%cdns );

    return \%data;
}

sub check_cdn_input {
    my $self = shift;

    my $sep = "__NEWLINE__";    # the line separator sub that with \n in the .ep javascript
    my $err = undef;

    # First, check permissions
    if ( !&is_oper($self) ) {
        $err .= "You do not have enough privileges to modify this." . $sep;
        return $err;
    }

    return $err;
}

sub isValidCdn {
    my $self = shift;
    $self->field('cdn_data.name')->is_required->is_like( qr/^[0-9a-zA-Z_\.\-]+$/, "Use alphanumeric . or _ ." );

    return $self->valid;
}

sub aprofileparameter {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs;
    if ( defined( $self->param('filter') ) ) {
        my $col = $self->param('filter');
        my $val = $self->param('value');

        # print "col: $col and val: $val \n";
        my $p_id = &profile_id( $self, $val );
        $rs = $self->db->resultset('Parameter')->search(
            { $col => $p_id },
            {
                join        => [ { 'profile_parameters' => 'parameter' }, { 'profile_parameters' => 'profile' }, ],
                '+select'   => ['profile.name'],
                '+as'       => ['profile_name'],
                '+order_by' => ['profile.name'],
                distinct    => 1,
            }
        );
    }
    else {
        $rs = $self->db->resultset('Parameter')->search(
            undef, {
                join        => [ { 'profile_parameters' => 'parameter' }, { 'profile_parameters' => 'profile' }, ],
                '+select'   => ['profile.name'],
                '+as'       => ['profile_name'],
                '+order_by' => ['profile.name'],
                distinct    => 1,
            }
        );

    }

    while ( my $row = $rs->next ) {
        my @line;
        @line = [ $row->id, $row->{_column_data}->{profile_name}, $row->name, $row->config_file, $row->value ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub aparameter {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $col = undef;
    my $val = undef;

    if ( defined( $self->param('filter') ) ) {
        $col = $self->param('filter');
        $val = $self->param('value');
    }

    my $rs = undef;
    if ( $col eq 'profile' and $val eq 'ORPHANS' ) { # Used with 'Parameters > Orphaned Parameters' menu item
        my $lindked_profile_rs    = $self->db->resultset('ProfileParameter')->search(undef);
        my $lindked_cachegroup_rs = $self->db->resultset('CachegroupParameter')->search(undef);
        $rs = $self->db->resultset('Parameter')->search(
            {
                -and => [
                    id => {
                        -not_in => $lindked_profile_rs->get_column('parameter')->as_query
                    },
                    id => {
                        -not_in => $lindked_cachegroup_rs->get_column('parameter')->as_query
                    }
                ]
            }
        );
        while ( my $row = $rs->next ) {
            my $secure = "no";
            if ( $row->secure == 1 ) {
                $secure = "yes";
            }
            my $value = $row->value;
            &UI::Parameter::conceal_secure_parameter_value( $self, $row->secure, \$value );
            my @line = [ $row->id, "NONE", $row->name, $row->config_file, $value, $secure, "profile" ];
            push( @{ $data{'aaData'} }, @line );
        }
        $rs = undef;
    }
    elsif ( $col eq 'profile' && $val ne 'all' ) { # Used with 'Parameters > Global Profile' menu item
        my $p_id = &profile_id( $self, $val );
        $rs = $self->db->resultset('ProfileParameter')->search( { $col => $p_id }, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );
    }
    elsif ( !defined($col) || ( $col eq 'profile' && $val eq 'all' ) ) { # Used with 'Parameters > All Profiles' menu item
        $rs = $self->db->resultset('ProfileParameter')->search( undef, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );
    }

    if ( defined($rs) ) {
        while ( my $row = $rs->next ) {
            my $secure = "no";
            if ( $row->parameter->secure == 1 ) {
                $secure = "yes";
            }
            my $value = $row->parameter->value;
            &UI::Parameter::conceal_secure_parameter_value( $self, $row->parameter->secure, \$value );
            my @line = [ $row->parameter->id, $row->profile->name, $row->parameter->name, $row->parameter->config_file, $value, $secure ];
            push( @{ $data{'aaData'} }, @line );
        }
    }

    $rs = undef;
    if ( $col eq 'cachegroup' && $val ne 'all' ) {
        my $l_id = $self->db->resultset('Cachegroup')->search( { short_name => $val } )->get_column('id')->single();
        $rs = $self->db->resultset('CachegroupParameter')
            ->search( { $col => $l_id }, { prefetch => [ { 'parameter' => undef }, { 'cachegroup' => undef } ] } );
    }
    elsif ( !defined($col) || ( $col eq 'cachegroup' && $val eq 'all' ) ) { # Used with 'Parameters > All Cache Groups' menu item
        $rs = $self->db->resultset('CachegroupParameter')->search( undef, { prefetch => [ { 'parameter' => undef }, { 'cachegroup' => undef } ] } );
    }

    if ( defined($rs) ) {
        while ( my $row = $rs->next ) {
            my $secure = "no";
            if ( $row->parameter->secure == 1 ) {
                $secure = "yes";
            }
            my $value = $row->parameter->value;
            &UI::Parameter::conceal_secure_parameter_value( $self, $row->parameter->secure, \$value );
            my @line = [ $row->parameter->id, $row->cachegroup->name, $row->parameter->name, $row->parameter->config_file, $row->parameter->value, $secure ];
            push( @{ $data{'aaData'} }, @line );
        }
    }

    $self->render( json => \%data );
}

sub aserver {
    my $self          = shift;
    my $server_select = shift;
    my %data          = ( "aaData" => [] );
    my $pparam =
        $self->db->resultset('ProfileParameter')
        ->search( { -and => [ 'parameter.name' => 'server_graph_url', 'profile.name' => 'GLOBAL' ] }, { prefetch => [ 'parameter', 'profile' ] } )
        ->single();
    my $srvg_url = defined($pparam) ? $pparam->parameter->value : '';

    my $rs = $self->db->resultset('Server')->search( undef, { prefetch => [ 'cdn', 'cachegroup', 'type', 'profile', 'status', 'phys_location' ] } );
    while ( my $row = $rs->next ) {
        my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";
        my @line;
        if ($server_select) {
            @line = [ $row->id, $row->host_name, $row->domain_name, $row->ip_address, $row->type->name, $row->profile->name, $cdn_name ];
        }
        else {
            my $aux_url = "";
            my $img     = "";

            if ( $row->type->name =~ m/^MID/ || $row->type->name =~ m/^EDGE/ ) {
                $aux_url = $srvg_url . $row->host_name;
                $img     = "graph.png";
            }
            elsif ( $row->type->name eq "CCR" ) {
                my $rs_param = $self->db->resultset('Parameter')->search(
                    {
                        'profile_parameters.profile' => $row->profile->id,
                        'name'                       => 'api.port'
                    },
                    { join => 'profile_parameters' }
                );
                my $r = $rs_param->single;
                my $port =
                    ( defined($r) && defined( $r->value ) )
                    ? $r->value
                    : 80;
                $aux_url = "http://" . $row->host_name . "." . $row->domain_name . ":" . $port . "/crs/stats";
                $img     = "info.png";
            }
            elsif ( $row->type->name eq "RASCAL" ) {
                $aux_url = "http://" . $row->host_name . "." . $row->domain_name . "/";
                $img     = "info.png";
            }

            my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";
            @line = [
                $row->id,              $row->host_name,        $row->domain_name,         "dummy",
                $cdn_name,             $row->cachegroup->name, $row->phys_location->name, $row->ip_address,
                $row->ip6_address,     $row->status->name,     $row->profile->name,       $row->ilo_ip_address,
                $row->mgmt_ip_address, $row->type->name,       $aux_url,                  $img,
		        $row->offline_reason
            ];
        }
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub aasn {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('Asn')->search( undef, { prefetch => [ { 'cachegroup' => 'cachegroups' }, ] } );

    while ( my $row = $rs->next ) {

        my @line = [ $row->id, $row->cachegroup->name, $row->asn, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub aphys_location {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('PhysLocation')->search( undef, { prefetch => ['region'] } );

    while ( my $row = $rs->next ) {

        next if $row->short_name eq 'UNDEF';

        my @line = [ $row->id, $row->name, $row->short_name, $row->address, $row->city, $row->state, $row->region->name, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub adeliveryservice {
    my $self       = shift;
    my %data       = ( "aaData" => [] );
    my %geo_limits = ( 0 => "none", 1 => "CZF", 2 => "CZF + Countries" );
    my %protocol   = ( 0 => "http", 1 => "https", 2 => "http/https", 3 => "http to https" );

    my $rs = $self->db->resultset('Deliveryservice')->search(
        {},
        {
            prefetch => [ 'type', 'cdn', { profile => { profile_parameters => 'parameter' } } ],
            join     => { profile => { profile_parameters => 'parameter' } },
            distinct => 1
        }
    );

    my $tenant_utils = Utils::Tenant->new($self);
    my $tenants_data = $tenant_utils->create_tenants_data_from_db();

    while ( my $row = $rs->next ) {
        if (!$tenant_utils->is_user_resource_accessible($tenants_data, $row->tenant_id)) {
            next;
        }

        my $cdn_name = defined( $row->cdn_id ) ? $row->cdn->name : "";

        # This will be undefined for 'Steering' delivery services
        my $org_server_fqdn = UI::DeliveryService::compute_org_server_fqdn($self, $row->id) // "";

        my $ptext = defined($row->profile) ? $row->profile->name : "-";
        my $line = [
            $row->id,                       $row->xml_id,                $org_server_fqdn,                "dummy",
            $cdn_name,                      $ptext,                      $row->ccr_dns_ttl,                    $yesno{ $row->active },
            $row->type->name,               $row->dscp,                  $row->signing_algorithm,              $row->qstring_ignore,
            $geo_limits{ $row->geo_limit }, $protocol{ $row->protocol }, $yesno{ $row->ipv6_routing_enabled }, $row->range_request_handling,
            $row->http_bypass_fqdn,         $row->dns_bypass_ip,         $row->dns_bypass_ip6,                 $row->dns_bypass_ttl,
            0.0 + $row->miss_lat,           0.0 + $row->miss_long,
        ];
        push( @{ $data{'aaData'} }, $line );
    }
    $self->render( json => \%data );
}

sub hwinfo {
    my $self           = shift;
    my $idisplay_start = $self->param("iDisplayStart") || 0;
    my $sort_order     = $self->param("sSortDir_0") || "asc";
    my $search_field   = $self->param("sSearch") || '';
    my $sort_column    = $self->param("iSortCol_0") || "id";
    my $echo           = $self->param("sEcho");

    # NOTE: If changes are made to send additional columns then this mapping has to be updated
    # to match the Column Number coming from datatables to it's name
    # Unfortunately, this is a short coming with the jquery datatables ui widget in that it expects
    # an array arrays instead of an array of hashes
    my $sort_direction        = sprintf( "-%s", $sort_order );
    my @column_number_to_name = qw{ serverid.host_name description val last_updated  };
    my $column_name           = $column_number_to_name[ $sort_column - 1 ] || "serverid";

    my $idisplay_length = $self->param("iDisplayLength") || 10;
    my %data = ( "data" => [] );

    my %condition;
    my %attrs;
    my %nolimit;
    my %nolimit_attrs;

    %condition = (
        -or => {
            'me.description'     => { -like => '%' . $search_field . '%' },
            'me.val'             => { -like => '%' . $search_field . '%' },
            'serverid.host_name' => { -like => '%' . $search_field . '%' }
        }
    );

    my $total_count = $self->db->resultset('Hwinfo')->search()->count();
    my $filtered_count;
    my $dbh;
    if ( $search_field eq '' ) {

        # if no filtering has occurred obviouslly these would be equal
        $filtered_count = $total_count;
    }
    else {
        # Now count the filtered records
        my %filtered_attrs = ( attrs => [ { 'serverid' => undef } ], join => 'serverid', undef );
        my $filtered = $self->db->resultset('Hwinfo')->search( \%condition, \%filtered_attrs );
        $filtered_count = $filtered->count();
    }
    my $page = $idisplay_start + 1;

    #my %limit = ( page => 100, rows => 100, order_by => { $sort_direction => $column_name } );

    my %limit = ( offset => $idisplay_start, rows => $idisplay_length, order_by => { $sort_direction => $column_name } );
    %attrs = ( attrs => [ { 'serverid' => undef } ], join => 'serverid', %limit );

    $dbh = $self->db->resultset('Hwinfo')->search( \%condition, \%attrs );

    # Now load up the rows
    while ( my $row = $dbh->next ) {
        my @line = [ $row->serverid->id, $row->serverid->host_name . "." . $row->serverid->domain_name, $row->description, $row->val, $row->last_updated ];
        push( @{ $data{'data'} }, @line );
    }

    %data = %{ merge( \%data, { recordsTotal    => $total_count } ) };
    %data = %{ merge( \%data, { recordsFiltered => $filtered_count } ) };

    $self->render( json => \%data );
}

sub ajob {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('Job')->search(
        undef, {
            prefetch => [       { 'job_user' => undef }, { agent => undef }, { status => undef } ],
            order_by => { -desc => 'me.entered_time' }
        }
    );

    while ( my $row = $rs->next ) {

        my @line = [ $row->id, $row->job_user->username, $row->asset_url, $row->asset_type, $row->entered_time, $row->status->name, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub alog {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $interval = "> now() - interval '30 day'";    # postgres
    my $rs = $self->db->resultset('Log')->search(
        { 'me.last_updated' => \$interval },
        {
            prefetch => [       { 'tm_user' => undef } ],
            order_by => { -desc => 'me.last_updated' },
            rows     => 1000
        }
    );

    while ( my $row = $rs->next ) {

        my @line = [ $row->last_updated, $row->level, $row->message, $row->tm_user->username, $row->ticketnum ];
        push( @{ $data{'aaData'} }, @line );
    }

    # setting cookie here, because the HTML page is often cached.
    my $date_string = `date "+%Y-%m-%d% %H:%M:%S"`;
    chomp($date_string);
    $self->cookie(
        last_seen_log => $date_string,
        { path => "/", max_age => 604800 }
    );    # expires in a week.
    $self->render( json => \%data );
}

sub acdn {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my %id_to_name = ();
    my $rs         = $self->db->resultset('Cdn')->search(undef);
    while ( my $row = $rs->next ) {
        $id_to_name{ $row->id } = $row->name;
    }

    $rs = $self->db->resultset('Cdn')->search(undef);
    while ( my $row = $rs->next ) {
        my @line = [ $row->id, $row->name, $row->domain_name, $yesno{ $row->dnssec_enabled }, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub acachegroup {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my %id_to_name = ();
    my $rs = $self->db->resultset('Cachegroup')->search( undef, { prefetch => [ { 'type' => undef } ] } );
    while ( my $row = $rs->next ) {
        $id_to_name{ $row->id } = $row->name;
    }

    $rs = $self->db->resultset('Cachegroup')->search( undef, { prefetch => [ { 'type' => undef }, 'coordinate' ] } );

    while ( my $row = $rs->next ) {
        my @line = [
            $row->id, $row->name, $row->short_name, $row->type->name,
            defined( $row->coordinate ) ? 0.0 + $row->coordinate->latitude : undef,
            defined( $row->coordinate ) ? 0.0 + $row->coordinate->longitude: undef,
            defined( $row->parent_cachegroup_id )
            ? $id_to_name{ $row->parent_cachegroup_id }
            : undef,
            $row->last_updated
        ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub auser {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('TmUser')->search( undef, { prefetch => [ { 'role' => undef } ] } );

    my $tenant_utils = Utils::Tenant->new($self);
    my $tenants_data = $tenant_utils->create_tenants_data_from_db();

    while ( my $row = $rs->next ) {
        if (!$tenant_utils->is_user_resource_accessible($tenants_data, $row->tenant_id)) {
            next;
        }
        my @line = [
            $row->id,           $row->username, $row->role->name, $row->full_name, $row->company,   $row->email,
            $row->phone_number, $row->uid,      $row->gid,        \1,              \$row->new_user, $row->last_updated
        ];

        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub afederation {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my @line;
    my $feds = $self->db->resultset('Federation')->search(undef);

    if ( $feds->count > 0 ) {
        while ( my $f = $feds->next ) {
            my $fed_id = $f->id;
            my $xml_id;
            my $user;

            # An assumption is being made that there is currently only a 1-1 relationship of the CNAME to the DeliveryService
            # Even though the datamodel supports multiples (at the moment)
            my $fed_dses = $f->federation_deliveryservices;
            while ( my $fd = $fed_dses->next ) {
                $xml_id = $fd->deliveryservice->xml_id;
            }

            my $tm_users = $f->federation_tmusers;
            while ( my $u = $tm_users->next ) {
                $user = $u->tm_user;
            }

            my $full_name = "";
            my $username  = "";
            my $company   = "";
            if ( defined($user) ) {
                $full_name = $user->full_name;
                $username  = $user->username;
                $company   = $user->company;
            }
            @line = [ $f->id, $f->cname, $xml_id, $f->description, $f->ttl, $full_name, $username, $company ];
            push( @{ $data{'aaData'} }, @line );
        }
    }

    $self->render( json => \%data );
}

sub aprofile {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('Profile')->search(undef, { prefetch => ['cdn'] } );

    while ( my $row = $rs->next ) {
        my $ctext = defined( $row->cdn ) ? $row->cdn->name : "-";
        my $routing_text = "No";
        if ( $row->routing_disabled == 1 ) {
            $routing_text = "Yes";
        }
        my @line = [ $row->id, $row->name, $row->name, $row->description, $row->type, $ctext, $routing_text, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub atype {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('Type')->search(undef);

    while ( my $row = $rs->next ) {
        my @line = [ $row->id, $row->name, $row->description, $row->use_in_table, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub adivision {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('Division')->search(undef);

    while ( my $row = $rs->next ) {
        my @line = [ $row->id, $row->name, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

sub aregion {
    my $self = shift;
    my %data = ( "aaData" => [] );

    my $rs = $self->db->resultset('Region')->search( undef, { prefetch => [ { 'division' => undef } ] } );

    while ( my $row = $rs->next ) {
        my @line = [ $row->id, $row->name, $row->division->name, $row->last_updated ];
        push( @{ $data{'aaData'} }, @line );
    }
    $self->render( json => \%data );
}

# TODO JvD: should really make all these lower case URLs. Mixed case URLs suck.
sub aadata {
    my $self  = shift;
    my $table = $self->param('table');

    if ( $table eq 'Serverstatus' ) {
        &aserverstatus($self);
    }
    elsif ( $table eq 'ProfileParameter' ) {
        &aprofileparameter($self);
    }
    elsif ( $table eq 'Server' ) {
        &aserver( $self, 0 );
    }
    elsif ( $table eq 'Asn' ) {
        &aasn($self);
    }
    elsif ( $table eq 'Deliveryservice' ) {
        &adeliveryservice($self);
    }
    elsif ( $table eq 'Hwinfo' ) {
        &hwinfo($self);
    }
    elsif ( $table eq 'Federation' ) {
        &afederation($self);
    }
    elsif ( $table eq 'ServerSelect' ) {
        &aserver( $self, 1 );
    }
    elsif ( $table eq 'Log' ) {
        &alog($self);
    }
    elsif ( $table eq 'Job' ) {
        &ajob($self);
    }
    elsif ( $table eq 'Cdn' ) {
        &acdn($self);
    }
    elsif ( $table eq 'Cachegroup' ) {
        &acachegroup($self);
    }
    elsif ( $table eq 'Type' ) {
        &atype($self);
    }
    elsif ( $table eq 'User' ) {
        &auser($self);
    }
    elsif ( $table eq 'Profile' ) {
        &aprofile($self);
    }
    elsif ( $table eq 'Parameter' ) {
        &aparameter($self);
    }
    elsif ( $table eq 'Physlocation' ) {
        &aphys_location($self);
    }
    elsif ( $table eq 'Division' ) {
        &adivision($self);
    }
    elsif ( $table eq 'Region' ) {
        &aregion($self);
    }

    else {
        $self->render( text => "Traffic Ops error, something is not configured properly." );
    }
}

#### JvD Start new UI stuff
sub loginpage {
    my $self = shift;
    $self->render( layout => undef );
}

# don't call this logout... recurses
sub logoutclicked {
    my $self = shift;

    $self->logout();
    return $self->redirect_to('/loginpage');
}

sub login {
    my $self = shift;

    my ( $u, $p ) = ( $self->req->param('u'), $self->req->param('p') );
    my $result = $self->authenticate( $u, $p );

    if ($result) {
        my $referer = $self->req->headers->header('referer');
        if ( !defined($referer) ) {
            $referer = '/';
        }
        if ( $referer =~ /\/login/ ) {
            if ( &UI::Utils::is_ldap($self) ) {
                $referer = '/dailysummary'; # LDAP-only users can't see edge_health
            } else {
                $referer = '/edge_health';
            }
        }
        return $self->redirect_to($referer);
    }
    else {
        $self->flash( login_msg => "Invalid username or password, please try again." );
        return $self->redirect_to('/loginpage');
    }
}

sub options {
    my $self = shift;

    # this essentially serves a blank page; options are in the HTTP header in Cdn.pm
    $self->res->headers->content_type("text/plain");
    $self->render(
        template => undef,
        layout   => undef,
        text     => "",
        status   => 200
    );
}

1;
