package UI::Topology;

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
use UI::Render;
use List::Compare;
use JSON;
use Data::Dumper;
use Mojo::Base 'Mojolicious::Controller';
use Time::HiRes qw(gettimeofday);
use File::Basename;
use File::Path;
use Scalar::Util qw(looks_like_number);

sub ccr_config {
    my $self     = shift;
    my $cdn_name = $self->param('cdnname');
    my $json     = &gen_crconfig_json( $self, $cdn_name );
    $self->render( json => $json );
}

# Produces a list of Cdns for traversing child links
sub gen_crconfig_json {
    my $self     = shift;
    my $cdn_name = shift;

    my $data_obj;
    my $cdn_id;
    my %type_to_name;
    my $ccr_domain_name = "";
    my $profile_cache;
    my $cdn_soa_minimum = 30;
    my $cdn_soa_expire  = 604800;
    my $cdn_soa_retry   = 7200;
    my $cdn_soa_refresh = 28800;
    my $cdn_soa_admin   = "traffic_ops";
    my $tld_ttls_soa    = 86400;
    my $tld_ttls_ns     = 3600;

    $SIG{__WARN__} = sub {
        warn $_[0]
            unless $_[0] =~ m/Prefetching multiple has_many rels deliveryservice_servers/;
    };

    $data_obj->{'stats'}->{'CDN_name'}   = $cdn_name;
    $data_obj->{'stats'}->{'date'}       = time();
    $data_obj->{'stats'}->{'tm_version'} = &tm_version();
    $data_obj->{'stats'}->{'tm_path'}    = $self->req->url->path->{'path'};
    $data_obj->{'stats'}->{'tm_host'}    = $self->req->headers->host;
    $data_obj->{'stats'}->{'tm_user'}    = $self->current_user()->{'username'};

    my $rs_cdn_profiles = $self->db->resultset('Server')->search(
        { 'cdn.name' => $cdn_name },
        {
            select => [ 'cdn.id', 'me.profile', 'me.type', 'profile.id', 'type.id' ],
            join   => 'cdn',
            prefetch => [ 'profile', 'type' ],
            distinct => 1,
            group_by => [ 'type.name', 'cdn.id', 'me.profile', 'me.type', 'profile.id', 'type.id'],
        }
    );

    while ( my $row = $rs_cdn_profiles->next ) {
        push( @{ $profile_cache->{ $row->type->name } }, $row->profile->id );
        $cdn_id = defined($cdn_id) ? next : $row->cdn->id;
    }

    my %param_cache;
    my @profile_caches;

    # key is the expression used in the regex below, value is the human readable string
    my $types = {
        "CCR"   => "Traffic Router",
        "^EDGE" => "EDGE",
        "^MID"  => "MID",
    };
    my $found;
    for my $cachetype ( keys %{$types} ) {
        
        for my $this_type ( keys %{$profile_cache} ) {
            if ( $this_type =~ m/$cachetype/ && scalar( @{ $profile_cache->{$this_type} } ) > 0 ) {
                push @profile_caches, @{ $profile_cache->{$this_type} };
                $found = 1;
            }
        }
    }

    if ( !$found ) {
        my $e = Mojo::Exception->throw( "No cache profiles found for CDN: " . $cdn_name );
    }

    my %condition = (
        -and => [
            profile                 => { -in => \@profile_caches, },
            'parameter.config_file' => 'CRConfig.json'
        ]
    );
    my $rs_pp = $self->db->resultset('ProfileParameter')->search( \%condition, { prefetch => [ { 'parameter' => undef }, { 'profile' => undef } ] } );

    #add dnssec.enabled value to config section
    my $cdn_rs = $self->db->resultset('Cdn')->search( { name => $cdn_name } )->single();
    my $dnssec_enabled = "false";
    if ( $cdn_rs->dnssec_enabled == 1 ) {
        $dnssec_enabled = "true";
    }
    $data_obj->{'config'}->{'dnssec.enabled'} = $dnssec_enabled;

    # These params should have consistent values across all profiles used by servers in this CDN:
    my %requested_param_names = (
        'tld.soa.admin'     => 1,
        'tld.soa.expire'    => 1,
        'tld.soa.minimum'   => 1,
        'tld.soa.refresh'   => 1,
        'tld.soa.retry'     => 1,
        'tld.ttls.SOA'      => 1,
        'tld.ttls.NS'       => 1,
        'LogRequestHeaders' => 1,
    );

    # Gather profile/parameter/value for each profile used by servers in this CDN
    while ( my $row = $rs_pp->next ) {
        my $param = $row->parameter->name;

        # cache value of each profile/param for later.
        $param_cache{ $row->profile->id }{$param} = $row->parameter->value;

        if ( $param =~ m/^tld/ ) {
            $param =~ s/tld\.//;
            ( my $top_key, my $second_key ) = split( /\./, $param );
            $data_obj->{'config'}->{$top_key}->{$second_key} = $row->parameter->value;
        }
        elsif ( $param eq 'LogRequestHeaders' ) {
            my $headers;
            foreach my $header ( split( /__RETURN__/, $row->parameter->value ) ) {
                $header = &trim_spaces($header);
                push( @$headers, $header );
            }
            $data_obj->{'config'}->{'requestHeaders'} = $headers;
        }
        elsif ( $param eq 'maxmind.default.override' ) {
            ( my $country_code, my $coordinates ) = split( /\;/, $row->parameter->value );
            ( my $lat, my $long ) = split( /\,/, $coordinates );
            my $geolocation = {
                'countryCode' => "$country_code",
                'lat' => $lat + 0,
                'long' => $long + 0
            };
            if ( !$data_obj->{'config'}->{'maxmindDefaultOverride'} ) {
                @{ $data_obj->{'config'}->{'maxmindDefaultOverride'} } = ();
            }
            push ( @{ $data_obj->{'config'}->{'maxmindDefaultOverride'} }, $geolocation );
        }
        elsif ( !exists $requested_param_names{$param} ) {
            $data_obj->{'config'}->{$param} = $row->parameter->value;
        }
    }

    my ( $param_values, $errors ) = extract_params( [ keys %requested_param_names ], \%param_cache );

    if ( scalar @$errors != 0 ) {
        my $msg = "Errors extracting profile parameters: " . join( '', @$errors );
        return undef, $msg;
    }

    $ccr_domain_name = $self->db->resultset('Cdn')->search( { id => $cdn_id } )->get_column('domain_name')->single();
    $data_obj->{'config'}->{'domain_name'} = $ccr_domain_name;
    $cdn_soa_admin   = $param_values->{'tld.soa.admin'};
    $cdn_soa_expire  = $param_values->{'tld.soa.expire'};
    $cdn_soa_minimum = $param_values->{'tld.soa.minimum'};
    $cdn_soa_refresh = $param_values->{'tld.soa.refresh'};
    $cdn_soa_retry   = $param_values->{'tld.soa.retry'};
    $tld_ttls_soa    = $param_values->{'tld.ttls.SOA'};
    $tld_ttls_ns     = $param_values->{'tld.ttls.NS'};

    my $regex_tracker;
    my $rs_regexes = $self->db->resultset('Regex')->search( {}, { 'prefetch' => 'type' } );
    while ( my $row = $rs_regexes->next ) {
        $regex_tracker->{ $row->id }->{'type'}    = $row->type->name;
        $regex_tracker->{ $row->id }->{'pattern'} = $row->pattern;
    }

    my %cache_tracker;
    my $rs_caches = $self->db->resultset('Server')->search(
        {
            'type.name' => [ { -like => 'EDGE%' }, { -like => 'MID%' }, { -like => 'CCR' }, { -like => 'RASCAL' }, { -like => 'TR' }, { -like => 'TM' } ],
            'me.cdn_id' => $cdn_id
        }, {
            prefetch => [ 'type',      'status',      { 'cachegroup' => 'coordinate' }, 'profile' ],
            columns  => [ 'host_name', 'domain_name', 'tcp_port', 'https_port',   'interface_name', 'ip_address', 'ip6_address', 'id', 'xmpp_id', 'profile.routing_disabled' ]
        }
    );

    while ( my $row = $rs_caches->next ) {

        next
            unless ( $row->status->name eq 'ONLINE'
            || $row->status->name eq 'REPORTED'
            || $row->status->name eq 'ADMIN_DOWN' );

        if ( $row->type->name eq "RASCAL" ) {
            $data_obj->{'monitors'}->{ $row->host_name }->{'fqdn'}      = $row->host_name . "." . $row->domain_name;
            $data_obj->{'monitors'}->{ $row->host_name }->{'status'}    = $row->status->name;
            $data_obj->{'monitors'}->{ $row->host_name }->{'location'}  = $row->cachegroup->name;
            $data_obj->{'monitors'}->{ $row->host_name }->{'port'}      = $row->tcp_port;
            $data_obj->{'monitors'}->{ $row->host_name }->{'httpsPort'} = $row->https_port;
            $data_obj->{'monitors'}->{ $row->host_name }->{'ip'}        = $row->ip_address;
            $data_obj->{'monitors'}->{ $row->host_name }->{'ip6'}       = ( $row->ip6_address || "" );
            $data_obj->{'monitors'}->{ $row->host_name }->{'profile'}   = $row->profile->name;

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
            my $port = ( defined($r) && defined( $r->value ) ) ? $r->value : 80;
            my $pid = $row->profile->id;
            my $weight =
                    defined( $param_cache{$pid}->{'weight'} )
                ? $param_cache{$pid}->{'weight'}
                : 0.999;
            my $weight_multiplier =
                    defined( $param_cache{$pid}->{'weightMultiplier'} )
                ? $param_cache{$pid}->{'weightMultiplier'}
                : 1000;

            $data_obj->{'contentRouters'}->{ $row->host_name }->{'fqdn'}        = $row->host_name . "." . $row->domain_name;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'status'}      = $row->status->name;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'location'}    = $row->cachegroup->name;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'port'}        = $row->tcp_port;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'httpsPort'}   = $row->https_port;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'api.port'}    = $port;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'ip'}          = $row->ip_address;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'ip6'}         = ( $row->ip6_address || "" );
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'profile'}     = $row->profile->name;
            $data_obj->{'contentRouters'}->{ $row->host_name }->{'hashCount'}   = int( $weight * $weight_multiplier );

            # Add Traffic Router cache groups to edgeLocations if valid lat/long is specified (0, 0 is unlikely to be valid)
            # This is necessary to enable localization for Traffic Router related DNS records
            if ( defined( $row->cachegroup->latitude ) && defined( $row->cachegroup->longitude ) &&
                ($row->cachegroup->latitude + 0) != 0 && ($row->cachegroup->longitude + 0) != 0 ) {
                    $data_obj->{'trafficRouterLocations'}->{ $row->cachegroup->name }->{'latitude'}  = $row->cachegroup->latitude + 0;
                    $data_obj->{'trafficRouterLocations'}->{ $row->cachegroup->name }->{'longitude'} = $row->cachegroup->longitude + 0;
            }
        }
        elsif ( $row->type->name =~ m/^EDGE/ || $row->type->name =~ m/^MID/ ) {

            if ( $row->type->name =~ m/^EDGE/ ) {
                $data_obj->{'edgeLocations'}->{ $row->cachegroup->name }->{'latitude'}  = $row->cachegroup->coordinate->latitude + 0;
                $data_obj->{'edgeLocations'}->{ $row->cachegroup->name }->{'longitude'} = $row->cachegroup->coordinate->longitude + 0;
                $data_obj->{'edgeLocations'}->{ $row->cachegroup->name }->{'backupLocations'}->{'fallbackToClosest'} = $row->cachegroup->fallback_to_closest ? "true" : "false";

                my $rs_backups = $self->db->resultset('CachegroupFallback')->search({ primary_cg => $row->cachegroup->id}, {order_by => 'set_order'});
                my $backup_cnt = 0;

                while ( my $backup_row = $rs_backups->next ) {
                    $data_obj->{'edgeLocations'}->{ $row->cachegroup->name }->{'backupLocations'}->{'list'}[$backup_cnt] = $backup_row->backup_cg->name; 
                    $backup_cnt++;
                }
            }

            if ( !exists $cache_tracker{ $row->id } ) {
                $cache_tracker{ $row->id } = $row->host_name;
            }

            my $pid = $row->profile->id;
            my $weight =
                defined( $param_cache{$pid}->{'weight'} )
                ? $param_cache{$pid}->{'weight'}
                : 0.999;
            my $weight_multiplier =
                defined( $param_cache{$pid}->{'weightMultiplier'} )
                ? $param_cache{$pid}->{'weightMultiplier'}
                : 1000;

            $data_obj->{'contentServers'}->{ $row->host_name }->{'locationId'}    = $row->cachegroup->name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'cacheGroup'}    = $row->cachegroup->name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'fqdn'}          = $row->host_name . "." . $row->domain_name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'port'}          = $row->tcp_port;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'httpsPort'}     = $row->https_port;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'interfaceName'} = $row->interface_name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'status'}        = $row->status->name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'ip'}            = $row->ip_address;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'ip6'}           = ( $row->ip6_address || "" );
            $data_obj->{'contentServers'}->{ $row->host_name }->{'profile'}       = $row->profile->name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'type'}          = $row->type->name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'hashId'}        = $row->xmpp_id ? $row->xmpp_id : $row->host_name;
            $data_obj->{'contentServers'}->{ $row->host_name }->{'hashCount'}     = int( $weight * $weight_multiplier );
            $data_obj->{'contentServers'}->{ $row->host_name }->{'routingDisabled'} = $row->profile->routing_disabled;
        }
    }
    my $regexps;
    my $rs_ds = $self->db->resultset('Deliveryservice')->search(
        {
			'me.cdn_id' => $cdn_id,
            'active'     => 1,
            'type.name' => { '!=', [ 'ANY_MAP' ] }
        },
        { prefetch => [ 'deliveryservice_servers', 'deliveryservice_regexes', 'type' ] }
    );

    while ( my $row = $rs_ds->next ) {
        my $protocol;
        if ( $row->type->name =~ m/DNS/ ) {
            $protocol = 'DNS';
        }
        else {
            $protocol = 'HTTP';
        }

        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'routingName'} = $row->routing_name;

        my @server_subrows = $row->deliveryservice_servers->all;
        my @regex_subrows  = $row->deliveryservice_regexes->all;

        my $regex_to_props;
        my %ds_to_remap;
        if ( scalar(@regex_subrows) ) {
            foreach my $subrow (@regex_subrows) {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[ $subrow->set_number ]->{'protocol'} = $protocol;
                $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'} =
                    $regex_tracker->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'};
                $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'set_number'} = $subrow->set_number;
                $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'type'} = $regex_tracker->{ $subrow->{'_column_data'}->{'regex'} }->{'type'};
                if ( $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'type'} eq 'HOST_REGEXP' ) {
                    $ds_to_remap{ $row->xml_id }->[ $subrow->set_number ] = $regex_to_props->{ $subrow->{'_column_data'}->{'regex'} }->{'pattern'};
                }
            }
        }
        my @domains;
        foreach my $regex ( sort keys %{$regex_to_props} ) {
            my $set_number = $regex_to_props->{$regex}->{'set_number'};
            my $pattern    = $regex_to_props->{$regex}->{'pattern'};
            my $type       = $regex_to_props->{$regex}->{'type'};
            if ( $type eq 'HOST_REGEXP' ) {
                push(
                    @{ $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[$set_number]->{'matchlist'} },
                    { 'match-type' => 'HOST', 'regex' => $pattern }
                );
                if ( $set_number == 0 ) {
                    my $host = $pattern;
                    $host =~ s/\\//g;
                    $host =~ s/\.\*//g;
                    $host =~ s/\.//g;
                    push @domains, "$host.$ccr_domain_name";
                }
            }
            elsif ( $type eq 'PATH_REGEXP' ) {
                push(
                    @{ $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[$set_number]->{'matchlist'} },
                    { 'match-type' => 'PATH', 'regex' => $pattern }
                );
            }
            elsif ( $type eq 'HEADER_REGEXP' ) {
                push(
                    @{ $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'matchsets'}->[$set_number]->{'matchlist'} },
                    { 'match-type' => 'HEADER', 'regex' => $pattern }
                );
            }
        }
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'domains'} = \@domains;

        if ( scalar(@server_subrows) ) {

            #my $host_regex = qr/(^(\.)+\*\\\.)(.*)(\\\.(\.)+\*$)/;
            my $host_regex1 = qr/\\|\.\*/;

            #MAT: Have to do this dedup because @server_subrows contains duplicates (* the # of host regexes)
            my %server_subrow_dedup;
            foreach my $subrow (@server_subrows) {
                $server_subrow_dedup{ $subrow->{'_column_data'}->{'server'} } = $subrow->{'_column_data'}->{'deliveryservice'};
            }
            foreach my $server ( keys %server_subrow_dedup ) {

                next if ( !defined( $cache_tracker{$server} ) );
                next if ( $data_obj->{'contentServers'}->{ $cache_tracker{$server} }->{'routingDisabled'} == 1);

                foreach my $host ( @{ $ds_to_remap{ $row->xml_id } } ) {
                    my $remap;
                    if ( $host =~ m/\.\*$/ ) {
                        my $host_copy = $host;
                        $host_copy =~ s/$host_regex1//g;
                        if ( $protocol eq 'DNS' ) {
                            $remap = $row->routing_name . $host_copy . $ccr_domain_name;
                        }
                        else {
                            $remap = $cache_tracker{$server} . $host_copy . $ccr_domain_name;
                        }
                    }
                    else {
                        $remap = $host;
                    }
                    push( @{ $data_obj->{'contentServers'}->{ $cache_tracker{$server} }->{'deliveryServices'}->{ $row->xml_id } }, $remap );
                }
            }
        }

        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'ttl'} = $row->ccr_dns_ttl;
        if ( $protocol ne 'DNS' ) {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'dispersion'} = {
                limit    => int( $row->initial_dispersion ),
                shuffled => 'true'
            };
        }

        my $geo_limit = $row->geo_limit;
        if ( $geo_limit == 0 ) {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'coverageZoneOnly'} = 'false';
        }
        elsif ( $geo_limit == 1 ) {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'coverageZoneOnly'} = 'true';
        }
        else {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'coverageZoneOnly'} = 'false';
            my $geoEnabled = [];
            foreach my $code ( split( ",", $row->geo_limit_countries ) ) {
                push( @$geoEnabled, { 'countryCode' => $code } );
            }
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'geoEnabled'} = $geoEnabled;
        }

        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'deepCachingType'} = $row->deep_caching_type;

        # Default to 'http only'
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'sslEnabled'} = 'false';
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'protocol'}->{'acceptHttps'} = 'false';
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'protocol'}->{'redirectToHttps'} = 'false';

        my $ds_protocol = $row->protocol;

        if (looks_like_number($ds_protocol) && 0 < $ds_protocol && $ds_protocol < 4) {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'sslEnabled'} = 'true';
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'protocol'}->{'acceptHttps'} = 'true';

            # 'https only'
            if ($ds_protocol == 1) {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'protocol'}->{'acceptHttp'} = 'false';

            }

            # 'http to https'
            if ($ds_protocol == 3) {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'protocol'}->{'redirectToHttps'} = 'true';
            }
        }

        my $geo_provider = $row->geo_provider;
        if ( $geo_provider == 1 ) {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'geolocationProvider'} = 'neustarGeolocationService';
        }
        else {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'geolocationProvider'} = 'maxmindGeolocationService';
        }

        if ( defined( $row->max_dns_answers )
            && $row->max_dns_answers ne "" )
        {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'maxDnsIpsForLocation'} = $row->max_dns_answers;
        }

        if ( $protocol =~ m/DNS/ ) {

            #$data_obj->{'deliveryServices'}->{$row->xml_id}->{'matchsets'}->[0]->{'protocol'} = 'DNS';
            if ( defined( $row->dns_bypass_ip ) && $row->dns_bypass_ip ne "" ) {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'DNS'}->{'ip'} = $row->dns_bypass_ip;
            }
            if ( defined( $row->dns_bypass_ip6 )
                && $row->dns_bypass_ip6 ne "" )
            {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'DNS'}->{'ip6'} = $row->dns_bypass_ip6;
            }
            if ( defined( $row->dns_bypass_cname )
                && $row->dns_bypass_cname ne "" )
            {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'DNS'}->{'cname'} = $row->dns_bypass_cname;
            }
            if ( defined( $row->dns_bypass_ttl )
                && $row->dns_bypass_ttl ne "" )
            {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'DNS'}->{'ttl'} = $row->dns_bypass_ttl;
            }
        }
        elsif ( $protocol =~ m/HTTP/ ) {

            #$data_obj->{'deliveryServices'}->{$row->xml_id}->{'matchsets'}->[0]->{'protocol'} = 'HTTP';
            if ( defined( $row->http_bypass_fqdn )
                && $row->http_bypass_fqdn ne "" )
            {
                my $full = $row->http_bypass_fqdn;
                my $fqdn;
                if ( $full =~ m/\:/ ) {
                    my $port;
                    ( $fqdn, $port ) = split( /\:/, $full );
                    # Specify port number only if explicitly set by the DS 'Bypass FQDN' field - issue 1493
                    $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'HTTP'}->{'port'} = $port;
                }
                else {
                    $fqdn = $full;
                }
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'bypassDestination'}->{'HTTP'}->{'fqdn'} = $fqdn;
            }

            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'regionalGeoBlocking'} = $row->regional_geo_blocking ? 'true' : 'false';
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'anonymousBlockingEnabled'} = $row->anonymous_blocking_enabled ? 'true' : 'false';

            if ( defined($row->geo_limit) && $row->geo_limit ne 0 ) {
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'geoLimitRedirectURL'} =
                    defined($row->geolimit_redirect_url) ? $row->geolimit_redirect_url : "";
            }
        }

        if ( defined( $row->tr_response_headers )
            && $row->tr_response_headers ne "" )
        {
            foreach my $header ( split( /__RETURN__/, $row->tr_response_headers ) ) {
                my ( $header_name, $header_value ) = split( /:\s/, $header );
                $header_name                                                                           = &trim_spaces($header_name);
                $header_value                                                                          = &trim_spaces($header_value);
                $header_value                                                                          = &trim_quotes($header_value);
                $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'responseHeaders'}->{$header_name} = $header_value;
            }
        }

        if ( defined( $row->tr_request_headers )
            && $row->tr_request_headers ne "" )
        {
            my $headers;
            foreach my $header ( split( /__RETURN__/, $row->tr_request_headers ) ) {
                $header = &trim_spaces($header);
                push( @$headers, $header );
            }
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'requestHeaders'} = $headers;
        }

        if ( defined( $row->miss_lat ) && $row->miss_lat ne "" ) {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'missLocation'}->{'lat'} = $row->miss_lat + 0;
        }
        if ( defined( $row->miss_long ) && $row->miss_long ne "" ) {
            $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'missLocation'}->{'long'} = $row->miss_long + 0;
        }

        my $ds_ttl = $row->ccr_dns_ttl;
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'ttls'} = {
            'A'    => "$ds_ttl",
            'AAAA' => "$ds_ttl",
            'NS'   => $tld_ttls_ns,
            'SOA'  => $tld_ttls_soa
        };
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'minimum'} = $cdn_soa_minimum;
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'expire'}  = $cdn_soa_expire;
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'retry'}   = $cdn_soa_retry;
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'refresh'} = $cdn_soa_refresh;
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'soa'}->{'admin'}   = $cdn_soa_admin;
        $data_obj->{'deliveryServices'}->{ $row->xml_id }->{'ip6RoutingEnabled'} = $row->ipv6_routing_enabled ? 'true' : 'false';

    }

    my $rs_dns = $self->db->resultset('Staticdnsentry')->search(
        {
            'deliveryservice.active' => 1,
            'deliveryservice.cdn_id' => $cdn_id
        }, {
            prefetch => [ 'deliveryservice', 'type' ],
            columns  => [ 'host',            'type', 'ttl', 'address' ]
        }
    );

    while ( my $dns_row = $rs_dns->next ) {

        my $dns_obj;
        $dns_obj->{'name'}  = $dns_row->host;
        $dns_obj->{'ttl'}   = $dns_row->ttl;
        $dns_obj->{'value'} = $dns_row->address;

        my $type = $dns_row->type->name;
        $type =~ s/\_RECORD//g;
        $dns_obj->{'type'} = $type;

        push( @{ $data_obj->{'deliveryServices'}->{ $dns_row->deliveryservice->xml_id }->{'staticDnsEntries'} }, $dns_obj );
    }

    return ($data_obj);
}

sub write_crconfig_json_to_db {
    my $self          = shift;
    my $cdn_name      = shift;
    my $crconfig_db   = shift;
    my $crconfig_json = encode_json($crconfig_db);

    my $snapshot = $self->db->resultset('Snapshot')->find( { cdn => $cdn_name } );
    if ( defined($snapshot) ) {
        $snapshot->update({ content => $crconfig_json });
    } else {
        my $insert = $self->db->resultset('Snapshot')->create( { cdn => $cdn_name, content => $crconfig_json } );
        $insert->insert();
    }

}

sub diff_crconfig_json {
    my $self     = shift;
    my $json     = shift;
    my $cdn_name = shift;

    my $current_snapshot = $self->db->resultset('Snapshot')->search( { cdn => $cdn_name } )->get_column('content')->single();

    if ( !defined($current_snapshot) )
    {
        my @err = ();
        $err[0] = "There is no existing CRConfig for " . $cdn_name . " to diff against... Is this the first snapshot???";
        my @caution = ();
        $caution[0] = "If you are not sure why you are getting this message, please do not proceed!";
        my @proceed = ();
        $proceed[0] = "To proceed writing the snapshot anyway click the 'Write CRConfig' button below.";
        my @dummy = ();
        return ( \@err, \@dummy, \@caution, \@dummy, \@dummy, \@proceed, \@dummy );
    }

    $current_snapshot = decode_json($current_snapshot);

    (
        my $ds_strings,
        my $loc_strings,
        my $cs_strings,
        my $csds_strings,
        my $rascal_strings,
        my $ccr_strings,
        my $cfg_strings
    ) = &crconfig_strings($current_snapshot);
    my @ds_strings     = @$ds_strings;
    my @loc_strings    = @$loc_strings;
    my @cs_strings     = @$cs_strings;
    my @csds_strings   = @$csds_strings;
    my @rascal_strings = @$rascal_strings;
    my @ccr_strings    = @$ccr_strings;
    my @cfg_strings    = @$cfg_strings;

    ( my $db_ds_strings, my $db_loc_strings, my $db_cs_strings, my $db_csds_strings, my $db_rascal_strings, my $db_ccr_strings, my $db_cfg_strings ) =
        &crconfig_strings($json);
    my @db_ds_strings     = @$db_ds_strings;
    my @db_loc_strings    = @$db_loc_strings;
    my @db_cs_strings     = @$db_cs_strings;
    my @db_csds_strings   = @$db_csds_strings;
    my @db_rascal_strings = @$db_rascal_strings;
    my @db_ccr_strings    = @$db_ccr_strings;
    my @db_cfg_strings    = @$db_cfg_strings;

    my @ds_text     = &compare_lists( \@db_ds_strings,     \@ds_strings,     "Section: Delivery Services" );
    my @loc_text    = &compare_lists( \@db_loc_strings,    \@loc_strings,    "Section: Edge Cachegroups" );
    my @cs_text     = &compare_lists( \@db_cs_strings,     \@cs_strings,     "Section: Traffic Servers" );
    my @csds_text   = &compare_lists( \@db_csds_strings,   \@csds_strings,   "Section: Traffic Server - Delivery Services" );
    my @rascal_text = &compare_lists( \@db_rascal_strings, \@rascal_strings, "Section: Traffic Monitors" );
    my @ccr_text    = &compare_lists( \@db_ccr_strings,    \@ccr_strings,    "Section: Traffic Routers" );
    my @cfg_text    = &compare_lists( \@db_cfg_strings,    \@cfg_strings,    "Section: CDN Configs" );

    return ( \@ds_text, \@loc_text, \@cs_text, \@csds_text, \@rascal_text, \@ccr_text, \@cfg_text );
}

sub crconfig_strings {
    my $config      = shift;
    my $config_json = $config;

    my @ds_strings;
    my @loc_strings;
    my @cs_strings;
    my @csds_strings;
    my @config_strings;
    my @rascal_strings;
    my @ccr_strings;

    foreach my $ds ( sort keys %{ $config_json->{'deliveryServices'} } ) {
        my $return = &stringify_ds( $config_json->{'deliveryServices'}->{$ds} );
        push( @ds_strings, "|DS:$ds$return" );
    }
    foreach my $cachegroup ( sort keys %{ $config_json->{'edgeLocations'} } ) {
        my $return = &stringify_cachegroup( $config_json->{'edgeLocations'}->{$cachegroup} );
        push( @loc_strings, "|edge-cachegroup:$cachegroup|$return" );
    }
    foreach my $server ( sort keys %{ $config_json->{'contentServers'} } ) {
        my $return = &stringify_content_server( $config_json->{'contentServers'}->{$server} );
        push( @cs_strings, $return );
        my @return = &stringify_cs_ds( $config_json->{'contentServers'}->{$server}->{'deliveryServices'}, $server );
        push( @csds_strings, @return );
    }
    foreach my $cfg ( sort keys %{ $config_json->{'config'} } ) {
        my $string;
        if ( $cfg eq 'ttls' || $cfg eq 'soa' ) {
            $string = "|param:$cfg";
            foreach my $key ( sort keys %{ $config_json->{'config'}->{$cfg} } ) {
                $string .= "|$key:" . $config_json->{'config'}->{$cfg}->{$key};
            }
            push( @config_strings, $string );
        }
        elsif ( $cfg eq 'maxmindDefaultOverride' ) {
            foreach my $element ( @{ $config_json->{'config'}->{$cfg} } ) {
                $string = "|param:$cfg";
                foreach my $key ( sort keys %{ $element } ) {
                    $string .= "|$key:" . $element->{$key};
                }
                push( @config_strings, $string );
            }
        }
        else {
            $string = "|param:$cfg|value:" . $config_json->{'config'}->{$cfg} . "|";
            push( @config_strings, $string );
        }
    }
    foreach my $rascal ( sort keys %{ $config_json->{'monitors'} } ) {
        my $return = &stringify_rascal( $config_json->{'monitors'}->{$rascal} );
        push( @rascal_strings, $return );
    }
    foreach my $ccr ( sort keys %{ $config_json->{'contentRouters'} } ) {
        my $return = &stringify_ccr( $config_json->{'contentRouters'}->{$ccr} );
        push( @ccr_strings, $return );
    }
    foreach my $trafficRouterLocation ( sort keys %{ $config_json->{'trafficRouterLocations'} } ) {
        my $return = &stringify_cachegroup( $config_json->{'trafficRouterLocations'}->{$trafficRouterLocation} );
        push( @ccr_strings, $return );
    }

    return ( \@ds_strings, \@loc_strings, \@cs_strings, \@csds_strings, \@rascal_strings, \@ccr_strings, \@config_strings );

}

sub stringify_ds {
    my $ds = shift;
    my $string;
    foreach my $matchset ( @{ $ds->{'matchsets'} } ) {
        $string .= "|<br>&emsp;protocol:" . $matchset->{'protocol'};
        foreach my $matchlist ( @{ $matchset->{'matchlist'} } ) {
            $string .= "|regex:" . $matchlist->{'regex'};
            $string .= "|match-type:" . $matchlist->{'match-type'};
        }
    }
    $string .= "|<br>&emsp;CZF Only:" . $ds->{'coverageZoneOnly'};
    if ( defined( $ds->{'geoEnabled'} ) ) {
        $string .= "|Geo Limit: true; Countries: ";
        foreach my $country ( @{ $ds->{'geoEnabled'} } ) {
            $string .= $country->{'countryCode'} . " ";
        }
    }
    if ( defined( $ds->{'missLocation'} ) ) {
        $string .= "|GeoMiss: " . $ds->{'missLocation'}->{'lat'} . "," . $ds->{'missLocation'}->{'long'};
    }
    if (defined( $ds->{'deepCachingType'} ) ) {
        $string .= "|deepCachingType: " . $ds->{'deepCachingType'};
    }
    if ( defined( $ds->{'bypassDestination'} ) ) {
        $string .= "<br>|BypassDest:";
        if ( defined( $ds->{'bypassDestination'}->{'DNS'}->{'ip'} ) ) {
            $string .= " -ip:" . $ds->{'bypassDestination'}->{'DNS'}->{'ip'};
        }
        if ( defined( $ds->{'bypassDestination'}->{'DNS'}->{'ip6'} ) ) {
            $string .= " -ip6:" . $ds->{'bypassDestination'}->{'DNS'}->{'ip6'};
        }
        if ( defined( $ds->{'bypassDestination'}->{'DNS'}->{'cname'} ) ) {
            $string .= " -cname:" . $ds->{'bypassDestination'}->{'DNS'}->{'cname'};
        }
        if ( defined( $ds->{'bypassDestination'}->{'DNS'}->{'ttl'} ) ) {
            $string .= " -ttl:" . $ds->{'bypassDestination'}->{'DNS'}->{'ttl'};
        }
        if ( defined( $ds->{'bypassDestination'}->{'HTTP'}->{'fqdn'} ) ) {
            $string .= " -fqdn:" . $ds->{'bypassDestination'}->{'HTTP'}->{'fqdn'};
        }
    }
    if ( defined( $ds->{'ip6RoutingEnabled'} ) ) {
        $string .= "|ip6RoutingEnabled: " . $ds->{'ip6RoutingEnabled'};
    }
    if ( defined( $ds->{'routingName'} ) ) {
        $string .= "|routingName: " . $ds->{'routingName'};
    }
    if ( defined( $ds->{'maxDnsIpsForLocation'} ) ) {
        $string .= "|maxDnsIpsForLocation:" . $ds->{'maxDnsIpsForLocation'};
    }
    if ( defined( $ds->{'responseHeaders'} ) ) {
        foreach my $header ( sort keys %{ $ds->{'responseHeaders'} } ) {
            $string .= "|responseHeader:$header:" . $ds->{'responseHeaders'}->{$header};
        }
    }
    if ( defined( $ds->{'dispersion'} ) ) {
        $string .= "|dispersion: limit=" . $ds->{'dispersion'}->{'limit'} . ", shuffled=" . $ds->{'dispersion'}->{'shuffled'};
    }
    if ( defined( $ds->{'geoLocationProvider'} ) ) {
        $string .= "|GeoLocation_Provider:" . $ds->{'geoLocationProvider'};
    }
    if ( defined( $ds->{'regionalGeoBlocking'} ) ) {
        $string .= "|Regional_Geoblocking:" . $ds->{'regionalGeoBlocking'};
    }
    if ( defined( $ds->{'anonymousBlockingEnabled'} ) ) {
        $string .= "|Anonymous_Blocking:" . $ds->{'anonymousBlockingEnabled'};
    }
    if ( defined( $ds->{'geoLimitRedirectURL'}) ) {
		$string .= "|Geolimit_Redirect_URL:" . $ds->{'geoLimitRedirectURL'};
	}
    $string .= "|<br>&emsp;DNS TTLs: A:" . $ds->{'ttls'}->{'A'} . " AAAA:" . $ds->{'ttls'}->{'AAAA'} . "|";
    foreach my $dns ( @{ $ds->{'staticDnsEntries'} } ) {
        $string .= "|<br>&emsp;staticDns: |name:" . $dns->{'name'} . "|type:" . $dns->{'type'} . "|ttl:" . $dns->{'ttl'} . "|addr:" . $dns->{'value'} . "|";
    }

    if (defined($ds->{'protocol'})) {
        $string .= "|protocol: ";

        if (defined($ds->{'protocol'}->{'acceptHttp'})) {
            $string .= " acceptHttp=" . $ds->{'protocol'}->{'acceptHttp'};
        }

        if (defined($ds->{'protocol'}->{'acceptHttps'})) {
            $string .= " acceptHttps=" . $ds->{'protocol'}->{'acceptHttps'};
        }

        if (defined($ds->{'protocol'}->{'redirectToHttps'})) {
            $string .= " redirectToHttps=" . $ds->{'protocol'}->{'redirectToHttps'};
        }
    }

    return $string;
}

sub stringify_cachegroup {
    my $loc    = shift;
    my $string = "longitude:" . $loc->{'longitude'} . "|latitude:" . $loc->{'latitude'} . "|";
    return $string;
}

sub stringify_content_server {
    my $cs = shift;
    my $string =
          "&emsp;|fqdn:"
        . $cs->{'fqdn'} . "|ip:"
        . $cs->{'ip'} . "|ip6:"
        . $cs->{'ip6'}
        . "|port:"
        . $cs->{'port'}
        . "|interfaceName: "
        . $cs->{'interfaceName'}
        . "|<br>&emsp;&emsp;|cacheGroup:"
        . (
        defined( $cs->{'cacheGroup'} )
        ? $cs->{'cacheGroup'}
        : $cs->{'locationId'}
        )
        . "|profile:"
        . $cs->{'profile'}
        . "|status:"
        . $cs->{'status'}
        . "|type: "
        . $cs->{'type'}
        . "|hashId: "
        . ( $cs->{'hashId'} || "" )
        . "|hashCount: "
        . ( $cs->{'hashCount'} || "" ) . "|";
    return $string;
}

sub stringify_cs_ds {
    my $csds   = shift;
    my $server = shift;
    my @strings;
    foreach my $ds ( sort keys %{$csds} ) {
        if ( ref( $csds->{$ds} ) eq 'ARRAY' ) {
            foreach my $map ( @{ $csds->{$ds} } ) {
                next if !defined $map;
                push( @strings, "|ds:" . $ds . "|server:" . $server . "|mapped:" . $map . "|" );
            }
        }
        else {
            push( @strings, "|ds:" . $ds . "|server:" . $server . "|mapped:" . $csds->{$ds} . "|" );
        }
    }
    return @strings;
}

sub stringify_rascal {
    my $rascal = shift;
    my $string =
          "|fqdn:"
        . $rascal->{'fqdn'} . "|ip:"
        . $rascal->{'ip'} . "|ip6:"
        . $rascal->{'ip6'}
        . "|port:"
        . $rascal->{'port'}
        . "|location:"
        . $rascal->{'location'}
        . "|status:"
        . $rascal->{'status'}
        . "|profile:"
        . $rascal->{'profile'} . "|";
    return $string;
}

sub stringify_ccr {
    my $ccr = shift;
    my $string =
          "|fqdn:"
        . $ccr->{'fqdn'} . "|ip:"
        . $ccr->{'ip'} . "|ip6:"
        . $ccr->{'ip6'}
        . "|port:"
        . $ccr->{'port'}
        . "|api.port:"
        . $ccr->{'api.port'}
        . "|location:"
        . $ccr->{'location'}
        . "|status:"
        . $ccr->{'status'}
        . "|profile:"
        . $ccr->{'profile'} . "|";
    return $string;
}

sub compare_lists {
    my $list_db   = shift;
    my $list_disk = shift;
    my $text      = shift;
    my @compare_text;

    my $list_compare_obj = List::Compare->new( \@$list_db, \@$list_disk );

    my @db_only = $list_compare_obj->get_Lonly;
    if ( $#db_only >= 0 ) {
        push( @compare_text, "    " . $text . " only in Traffic Ops:" );
        foreach my $ds_string (@db_only) {
            push( @compare_text, "        " . $ds_string );
        }
    }
    my @disk_only = $list_compare_obj->get_Ronly;
    if ( $#disk_only >= 0 ) {
        push( @compare_text, "    " . $text . " only in old file:" );
        foreach my $ds_string (@disk_only) {
            push( @compare_text, "        " . $ds_string );
        }
    }
    if ( $#disk_only < 0 && $#db_only < 0 ) {
        push( @compare_text, "    " . $text . " are the same." );
    }
    return @compare_text;
}

sub extract_params {

    # array of param names to look for
    my $param_names = shift;

    # hash of profile id/param name/param value pulled from db
    my $param_cache = shift;

    my %errors;
    my %return_params;

    # ensure each param has exactly one value
    for my $param_name (@$param_names) {
        my $param_val;
        for my $profile_id ( keys %$param_cache ) {
            if ( !exists $param_cache->{$profile_id}{$param_name} ) {
                next;
            }
            my $new_val = $param_cache->{$profile_id}{$param_name};
            if ( defined $param_val && $new_val ne $param_val ) {

                # ERROR!!
                push @{ $errors{$param_name} }, $param_val, $new_val;
            }
            $param_val = $new_val;
        }
        $return_params{$param_name} = $param_val;
    }

    # Create a single error message for each parameter with inconsistent values
    my @errors;
    for my $param_name ( keys %errors ) {

        # filter out dups
        my %seen;
        my @values = grep { !$seen{$_}++ } @{ $errors{$param_name} };
        push @errors, "Parameter $param_name has multiple values (", join( ', ', @values ) . ") from profiles associated with servers in this CDN. ";
    }
    return \%return_params, \@errors;
}

sub trim_spaces {
    my $text = shift;
    $text =~ s/^\s+//g;
    $text =~ s/\s+$//g;
    return $text;
}

sub trim_quotes {
    my $text = shift;
    $text =~ s/^\"//g;
    $text =~ s/\"$//g;
    return $text;
}
1;
