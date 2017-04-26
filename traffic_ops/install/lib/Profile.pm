#!/usr/bin/perl

#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

package Profile;

use warnings;
use strict;

use InstallUtils qw{ :all };
use Data::Dumper qw(Dumper);

use base qw{ Exporter };
our @EXPORT_OK = qw{ setup };
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

sub setup {
    my $dbh        = shift;
    my $paramconf = shift;

    my $insert_stmt;
    InstallUtils::logger( "paramconf " . Dumper($paramconf), "info" );

    setup_parameters($dbh, $paramconf);

    setup_profiles($dbh, $paramconf);

}

sub setup_parameters {
    my $dbh = shift;
    my $paramconf = shift;

    InstallUtils::logger( "=========== Setting up parameters", "info" );
    my $insert_stmt = <<"QUERY";
    insert into parameter (name, config_file, value) values ('tm.url', 'global', '$paramconf->{"tm.url"}') ON CONFLICT (name, config_file, value) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    my $insert_stmt = <<"QUERY";
    insert into parameter (name, config_file, value) values ('tm.toolname', 'global', 'Traffic Ops') ON CONFLICT (name, config_file, value) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);
}

sub setup_profiles {
    my $dbh = shift;
    my $paramconf = shift;

    InstallUtils::logger( "=========== Setting up profiles", "info" );
    my $insert_stmt = <<"QUERY";
    insert into profile (name, description, type) values ('GLOBAL', 'Global Traffic Ops profile, DO NOT DELETE', 'UNK_PROFILE') ON CONFLICT (name) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.url' and config_file = 'global' and value = '$paramconf->{"tm.url"}') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.toolname' and config_file = 'global' and value = 'Traffic Ops') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.infourl' and config_file = 'global' and value = '$paramconf->{"tm.url"}/doc') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);
 
    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.logourl' and config_file = 'global' and value = '/images/tc_logo.png') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.instance_name' and config_file = 'global' and value = 'Traffic Ops CDN') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.traffic_mon_fwd_proxy' and config_file = 'global' and value = '$paramconf->{"tm.url"}:81') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation.polling.url' and config_file = 'CRConfig.json' and value = 'http://$paramconf->{"dns_subdomain"}/cdn/MaxMind/GeoLiteCity.dat.gz') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);

    $insert_stmt = <<"QUERY";
    insert into profile_parameter (profile, parameter) values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation6.polling.url' and config_file = 'CRConfig.json' and value = 'http://$paramconf->{"dns_subdomain"}/cdn/MaxMind/GeoLiteCityv6.dat.gz') ) ON CONFLICT (profile, parameter) DO NOTHING;
QUERY
    doInsert($dbh, $insert_stmt);
}

sub doInsert {
    my $dbh = shift;
    my $insert_stmt = shift;

    InstallUtils::logger( "\n" . $insert_stmt, "info" );
    my $stmt = $dbh->prepare($insert_stmt);
    $stmt->execute();
}

1;
