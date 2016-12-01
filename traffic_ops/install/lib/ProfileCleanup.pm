
package ProfileCleanup;

use InstallUtils qw{ :all };
use WWW::Curl::Easy;
use LWP::UserAgent;

use base qw{ Exporter };
our @EXPORT_OK = qw{ replace_profile_templates import_profiles profiles_exist };
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

sub profile_replace {
    my ($profile) = @_;
    my $profile_bak = $profile . ".bak";
    rename( $profile, $profile_bak ) or die("rename(): $!");
    open( my $fh,  '<', $profile_bak ) or die("open(): $!");
    open( my $ofh, '>', $profile )     or die("open(): $!");
    while (<$fh>) {
        s/{{.TmUrl}}/$::parameters->{'tm.url'}/g;
        s/{{.TmInfoUrl}}/$::parameters->{"tminfo.url"}/g;
        s/{{.TmInstanceName}}/$::parameters->{"cdnname"}/g;
        s/{{.GeolocationPollingUrl}}/$::parameters->{"geolocation.polling.url"}/g;
        s/{{.Geolocation6PollingUrl}}/$::parameters->{"geolocation6.polling.url"}/g;
        s/{{.TmUrl}}/$::parameters->{'tm.url'}/g;
        s/{{.TmToolName}}/Traffic Ops/g;
        s/{{.HealthPollingInterval}}/$::parameters->{"health.polling.interval"}/g;
        s/{{.CoveragezonePollingUrl}}/$::parameters->{"coveragezone.polling.url"}/g;
        s/{{.DomainName}}/$::parameters->{"domainname"}/g;
        s/{{.TldSoaAdmin}}/$::parameters->{"tld.soa.admin"}/g;
        s/{{.DrivePrefix}}/$::parameters->{"Drive_Prefix"}/g;
        s/{{.HealthThresholdLoadavg}}/$::parameters->{"health.threshold.loadavg"}/g;
        s/{{.HealthThresholdAvailableBandwidthInKbps}}/$::parameters->{"health.threshold.availableBandwidthInKbps"}/g;
        s/{{.RAMDrivePrefix}}/$::parameters->{"RAM_Drive_Prefix"}/g;
        s/{{.RAMDriveLetters}}/$::parameters->{"RAM_Drive_Letters"}/g;
        s/{{.HealthConnectionTimeout}}/$::parameters->{"health.connection.timeout"}/g;
        s#{{.CronOrtSyncds}}#*/15 * * * * root /opt/ort/traffic_ops_ort.pl syncds warn $::parameters->{'tm.url'} $tmAdminUser:$tmAdminPw > /tmp/ort/syncds.log 2>&1#g;
        print $ofh $_;
    }
    close $fh;
    close $ofh;
    unlink $profile_bak;
}

sub replace_profile_templates {
    my $conf = shift;

    $::parameters->{'tm.url'}                                    = $conf->{"tm.url"};
    $::parameters->{"tminfo.url"}                                = "$::parameters->{'tm.url'}/info";
    $::parameters->{"cdnname"}                                   = $conf->{"cdn_name"};
    $::parameters->{"geolocation.polling.url"}                   = "$::parameters->{'tm.url'}/routing/GeoIP2-City.mmdb.gz";
    $::parameters->{"geolocation6.polling.url"}                  = "$::parameters->{'tm.url'}/routing/GeoIP2-Cityv6.mmdb.gz";
    $::parameters->{"health.polling.interval"}                   = $conf->{"health_polling_int"};
    $::parameters->{"coveragezone.polling.url"}                  = "$::parameters->{'tm.url'}/routing/coverage-zone.json";
    $::parameters->{"domainname"}                                = $conf->{"dns_subdomain"};
    $::parameters->{"tld.soa.admin"}                             = $conf->{"soa_admin"};
    $::parameters->{"Drive_Prefix"}                              = $conf->{"driver_prefix"};
    $::parameters->{"RAM_Drive_Prefix"}                          = $conf->{"ram_drive_prefix"};
    $::parameters->{"RAM_Drive_Letters"}                         = $conf->{"ram_drive_letters"};
    $::parameters->{"health.threshold.loadavg"}                  = $conf->{"health_thresh_load_avg"};
    $::parameters->{"health.threshold.availableBandwidthInKbps"} = $conf->{"health_thresh_kbps"};
    $::parameters->{"health.connection.timeout"}                 = $conf->{"health_connect_timeout"};

    profile_replace( $::profile_dir . "profile.global.traffic_ops" );
    profile_replace( $::profile_dir . "profile.traffic_monitor.traffic_ops" );
    profile_replace( $::profile_dir . "profile.traffic_router.traffic_ops" );
    profile_replace( $::profile_dir . "profile.trafficserver_edge.traffic_ops" );
    profile_replace( $::profile_dir . "profile.trafficserver_mid.traffic_ops" );
    writeJson( $::post_install_cfg, $::parameters );
}

# Takes the Traffic Ops URI, user, and password.
# Returns the cookie, or the empty string on error
sub get_traffic_ops_cookie {
    my ( $uri, $user, $pass ) = @_;

    my $loginUri = "/api/1.2/user/login";

    my $curl          = WWW::Curl::Easy->new;
    my $response_body = "";
    open( my $fileb, ">", \$response_body );
    my $loginData = JSON::encode_json( { u => $user, p => $pass } );
    $curl->setopt( WWW::Curl::Easy::CURLOPT_URL,            $uri . $loginUri );
    $curl->setopt( WWW::Curl::Easy::CURLOPT_SSL_VERIFYPEER, 0 );
    $curl->setopt( WWW::Curl::Easy::CURLOPT_HEADER,         1 );                  # include header in response
    $curl->setopt( WWW::Curl::Easy::CURLOPT_NOBODY,         1 );                  # disclude body in response
    $curl->setopt( WWW::Curl::Easy::CURLOPT_POST,           1 );
    $curl->setopt( WWW::Curl::Easy::CURLOPT_POSTFIELDS,     $loginData );
    $curl->setopt( WWW::Curl::Easy::CURLOPT_WRITEDATA,      $fileb );             # put response in this var
    $curl->perform();

    my $cookie = $response_body;
    if ( $cookie =~ /mojolicious=(.*); expires/ ) {
        $cookie = $1;
    }
    else {
        $cookie = "";
    }
    return $cookie;
}

# Takes the filename of a Traffic Ops (TO) profile to import, the TO URI, and the TO login cookie
sub profile_import_single {
    my ( $profileFilename, $uri, $trafficOpsCookie ) = @_;
    logger( "Importing Profiles with: " . "curl -v -k -X POST -H \"Cookie: mojolicious=$trafficOpsCookie\" -F \"filename=$profileFilename\" -F \"profile_to_import=\@$profileFilename\" $uri/profile/doImport", "info" );
    my $rc = execCommand("curl -v -k -X POST -H \"Cookie: mojolicious=$trafficOpsCookie\" -F \"filename=$profileFilename\" -F \"profile_to_import=\@$profileFilename\" $uri/profile/doImport");
    if ( $rc != 0 ) {
        logger( "Failed to import Traffic Ops profile, check the console output and rerun postinstall once you've resolved the error", "error" );
    }
}

sub import_profiles {
    my $config = shift;
    logger( "Importing profiles...", "info" );

    my $toUri  = $::parameters->{'tm.url'};
    my $toUser = $config->{"username"};
    my $toPass = $config->{"password"};

    my $toCookie = get_traffic_ops_cookie( $toUri, $toUser, $toPass );

    logger( "Got cookie: " . $toCookie, "info" );

    # \todo use an array?
    logger( "Importing Global profile...", "info" );
    profile_import_single( $::profile_dir . "profile.global.traffic_ops", $toUri, $toCookie );
    logger( "Importing Traffic Monitor profile...", "info" );
    profile_import_single( $::profile_dir . "profile.traffic_monitor.traffic_ops", $toUri, $toCookie );
    logger( "Importing Traffic Router profile...", "info" );
    profile_import_single( $::profile_dir . "profile.traffic_router.traffic_ops", $toUri, $toCookie );
    logger( "Importing TrafficServer Edge profile...", "info" );
    profile_import_single( $::profile_dir . "profile.trafficserver_edge.traffic_ops", $toUri, $toCookie );
    logger( "Importing TrafficServer Mid profile...", "info" );
    profile_import_single( $::profile_dir . "profile.trafficserver_mid.traffic_ops", $toUri, $toCookie );
    logger( "Finished Importing Profiles.", "info" );
}

sub profiles_exist {
    my $config = shift;
    my $tmurl  = shift;

    if ( -f $::reconfigure_defaults ) {
        logger( "Default profiles were previously created. Remove " . $::reconfigure_defaults . " to create again", "warn" );
        return 1;
    }

    $::parameters->{'tm.url'} = $tmurl;

    my $uri = $::parameters->{'tm.url'};
    my $toCookie = get_traffic_ops_cookie( $::parameters->{'tm.url'}, $config->{"username"}, $config->{"password"} );

    my $profileEndpoint = "/api/1.2/profiles.json";

    my $ua = LWP::UserAgent->new;
    $ua->ssl_opts( verify_hostname => 0, SSL_verify_mode => 0x00 );
    my $req = HTTP::Request->new( GET => $uri . $profileEndpoint );
    $req->header( 'Cookie' => "mojolicious=" . $toCookie );
    my $resp = $ua->request($req);

    if ( !$resp->is_success ) {
        logger( "Error checking if profiles exist: " . $resp->status_line, "error" );
        return 1;    # return true, so we don't attempt to create profiles
    }
    my $message = $resp->decoded_content;

    my $profiles = JSON->new->utf8->decode($message);
    if (   ( !defined $profiles->{"response"} )
        || ( ref $profiles->{"response"} ne 'ARRAY' ) )
    {
        logger( "Error checking if profiles exist: invalid JSON: $message", "error" );
        return 1;    # return true, so we don't attempt to create profiles
    }

    my $num_profiles = scalar( @{ $profiles->{"response"} } );
    logger( "Existing Profile Count: $num_profiles", "info" );

    my %initial_profiles = (
        "INFLUXDB"      => 1,
        "RIAK_ALL"      => 1,
        "TRAFFIC_STATS" => 1
    );

    my $profiles_response = $profiles->{"response"};
    foreach my $profile (@$profiles_response) {
        if ( !exists $initial_profiles{ $profile->{"name"} } ) {
            logger( "Found existing profile (" . $profile->{"name"} . ")", "info" );
            open( my $reconfigure_defaults_file, '>', $::reconfigure_defaults ) or die("Failed to open() $reconfigure_defaults: $!");
            close($reconfigure_defaults_file);
            return 1;
        }
    }
    return 0;
}
