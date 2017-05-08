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

use lib qw(/opt/traffic_ops/install/lib /opt/traffic_ops/app/lib /opt/traffic_ops/app/local/lib/perl5);

$ENV{PERL5LIB} = "/opt/traffic_ops/install/lib:/opt/traffic_ops/app/lib:/opt/traffic_ops/app/local/lib/perl5:$ENV{PERL5LIB}";
$ENV{PATH}     = "/usr/bin:/opt/traffic_ops/go/bin:/usr/local/go/bin:/opt/traffic_ops/install/bin:$ENV{PATH}";
$ENV{GOPATH} = "/opt/traffic_ops/go";

use strict;
use warnings;

use Safe;
use DBI;
use POSIX;
use Digest::SHA1 qw(sha1_hex);
use Data::Dumper qw(Dumper);
use Getopt::Long;
use Cwd;
my $pwd = getcwd;
my $dbh;

use InstallUtils qw{ :all };
use Database qw{ connect };

# paths of the output configuration files
my $databaseConfFile = "/opt/traffic_ops/app/conf/production/database.conf";
my $dbConfFile       = "/opt/traffic_ops/app/db/dbconf.yml";
# log file for the installer
my $logFile = "/var/log/traffic_ops/export_profiles.log";

# debug mode
my $debug = 0;
my $maxLogSize = 10000000;    #bytes
my $outputFilePath;


sub getLatestProfiles {
    my $result;

    my $q=sprintf("SELECT * FROM profile as p 
                            INNER JOIN profile_parameter as pp ON p.id = pp.profile 
                            INNER JOIN parameter as param ON pp.parameter = param.id 
                            WHERE param.config_file = 'build_setup' order by param.name");
    my $stmt = $dbh->prepare($q);
    $stmt->execute(); 
    return $stmt;
}

sub parameterInsert {

    my $srcProfileName = shift;
    my $targetProfileName = shift;
    my $fileName = shift;

    my $excludes = { 
             #NOTE: Prefix begins with ^, 
             #      Suffix ends with \$
              name => { prefix   => "^visual_status_panel|^latest_" ,
                        suffix   => "^.dnssec.inception\$|_fw_proxy\$|_graph_url\$", 
                        contains => "^allow_ip|allow_ip6|purge_allow_ip.*|.*ramdisk_size.*", 
                      },
             value => {
                        contains => "comcast", 
                      },
       config_file => { prefix   => "^teak|^url_sig_|^regex_remap_|^hdr_rw_" ,
                        suffix   => ".dnssec.inception\$|_fw_proxy$|_graph_url\$", 
                        contains => "^dns\.zone\$|^http-log4j\.properties\$|^cacheurl_voice-guidance-tts.config\$|^cacheurl_col-coam-ads-jitp.config\$|^cacheurl_cloudtv-web-comp.config\$", 
                      },
    };
    my $nameExcludes=sprintf("%s|%s|%s", $excludes->{name}{prefix}, $excludes->{name}{suffix}, $excludes->{name}{contains});
    my $valueExcludes=sprintf("%s", $excludes->{value}{contains});
    my $configFileExcludes=sprintf("%s|%s|%s", $excludes->{config_file}{prefix}, $excludes->{config_file}{suffix}, $excludes->{config_file}{contains});

    my $q=sprintf("SELECT * FROM profile as p 
                            INNER JOIN profile_parameter as pp ON p.id = pp.profile 
                            INNER JOIN parameter as param ON pp.parameter = param.id 
                            WHERE p.name = '%s' AND 
                                  param.name !~ '%s' AND 
                                  param.value !~ '%s' AND 
                                  param.config_file !~ '%s' 
                                  order by param.name", $srcProfileName, $nameExcludes, $valueExcludes, $configFileExcludes);
    my $stmt = $dbh->prepare($q);
    $stmt->execute();
    my $insert;

    while ( my $row = $stmt->fetchrow_hashref ) {
       
       $insert=sprintf("INSERT INTO parameter (name, config_file, value) VALUES ('%s','%s','%s') ON CONFLICT (name, config_file, value) DO NOTHING;", $row->{name}, $row->{config_file}, scrub_value($row->{value}));

       appendToFile($fileName, $insert);
       profileParameterInsert($srcProfileName, $targetProfileName, $row->{name}, scrub_value($row->{value}), $row->{config_file}, $fileName);
    }
}

sub scrub_value {
    my $value = shift;
	$value =~ s/ xmt="%\<\{X-MoneyTrace\}cqh>"//g;
	return $value;
}

sub profileParameterInsert {

    my $srcProfileName = shift;
    my $targetProfileName = shift;
    my $targetParameterName = shift;
    my $targetParameterValue = shift;
    my $configFile = shift;
    my $fileName = shift;

    my $insert=sprintf("INSERT INTO profile_parameter (profile, parameter) VALUES ( (select id from profile where name = '%s'), (select id from parameter where name = '%s' and config_file = '%s' and value = '%s') )  ON CONFLICT (profile, parameter) DO NOTHING;\n", $targetProfileName, $targetParameterName, $configFile, $targetParameterValue );

    #InstallUtils::logger( "insert: $insert", "debug" );
    appendToFile($fileName, $insert);
}


sub profileInsert {

    my $srcProfileName = shift;
    my $targetProfileName = shift;
    my $targetProfileDesc = shift;
    my $fileName = shift;

    my $q=sprintf("SELECT * FROM profile WHERE name = '%s' order by name;", $srcProfileName, $fileName);

    my $stmt = $dbh->prepare($q);
    $stmt->execute();
    my $name;
    my $insert;

    while ( my $row = $stmt->fetchrow_hashref ) {
       
       #InstallUtils::logger( Dumper($row), "debug" );
       $insert=sprintf("INSERT INTO profile (name, description, type) VALUES ('%s','%s','%s') ON CONFLICT (name) DO NOTHING;", $targetProfileName, $targetProfileDesc, $row->{type});
       #InstallUtils::logger( "insert $insert", "debug" );
       appendToFile($fileName, $insert);
    }

}

sub generateInserts {

    my $srcProfileName = shift;
    my $targetProfileName = shift;
    my $targetProfileDesc = shift;
    my $configFile = shift;
    my $fileName = shift;

    profileInsert($srcProfileName, $targetProfileName, $targetProfileDesc, $fileName);
    parameterInsert($srcProfileName, $targetProfileName, $fileName);
}

sub appendToFile {

    my $filename = shift;
    my $contents = shift;
    open(my $fh, '>>', $outputFilePath) or die "Could not open file '$outputFilePath' $!";
    say $fh $contents;
    close $fh;
}

# -cfile     - Input File:       The input config file used to ask and answer questions
#                                will look to the defaults for the answer. If the answer is not in the defaults
#                                the script will exit
# -defaults  - Defaults:         Writes out a configuration file with defaults which can be used as input
# -debug     - Debug Mode:       More output to the terminal
# -h         - Help:             Basic command line help menu

sub main {
    my $help = 0;

    my $inputFile;

    # help string
    my $usageString = "Usage: $0 [-debug]\n";

    GetOptions(
        "cfile=s"     => \$inputFile,
        "debug"       => \$debug,
        "help"        => \$help
    ) or die($usageString);

    if ($help) {
        print $usageString;
        return;
    }

    # check if the user running postinstall is root
    if ( $ENV{USER} ne "root" ) {
        errorOut("You must run this script as the root user");
    }

    InstallUtils::initLogger( $debug, $logFile );

    InstallUtils::logger( "Starting export", "info" );

    if ($debug) {
      InstallUtils::logger( "Debug is on", "info" );
    }

    InstallUtils::rotateLog($logFile);

    if ( -s $logFile > $maxLogSize ) {
        InstallUtils::logger( "Postinstall log above max size of $maxLogSize bytes - rotating", "info" );
        rotateLog($logFile);
    }

    $dbh = Database::connect($databaseConfFile);
    my $latestProfiles  = getLatestProfiles();
    my $name;

    my $lookupTable = { 
       latest_traffic_monitor => { profile => { name => "TM_PROFILE" ,
                                               description => "Traffic Monitor" }, 
                                 },
                                 { parameter => { config_file => "rascal-config.txt",
                                                  description => "Traffic Monitor" }, 
                                 },
       latest_traffic_router  => { profile => { name => "TR_PROFILE" ,
                                               description => "Traffic Router" }, 
                                 },
                                 { parameter => { config_file => "CRConfig.json",
                                                  description => "Traffic Router" }, 
                                 },
       latest_traffic_stats  => { profile => { name => "TS_PROFILE" ,
                                               description => "Traffic Stats" }, 
                                 },
       latest_traffic_vault  => { profile => { name => "TV_PROFILE" ,
                                               description => "Traffic Vault" }, 
                                 },
       latest_trafficserver_edge  => { profile => { name => "EDGE_PROFILE" ,
                                               description => "Edge Cache" }, 
                                 },
       latest_trafficserver_mid  => { profile => { name => "MID_PROFILE" ,
                                               description => "Mid Cache" }, 
                                 },
    };

    InstallUtils::logger( "Looking for latest profiles", "info" );
    while ( my $row = $latestProfiles->fetchrow_hashref ) {
       
       my $name = $row->{'name'};
       my $fileName = $row->{name}. ".sql";
       my $srcProfileName = $row->{value};
       $fileName =~ s/latest_//g;

       if (exists $lookupTable->{$name}) {
          my $profile = $lookupTable->{$name};
          my $targetProfileName = $profile->{profile}{name};
          my $targetProfileDesc = $profile->{profile}{description};
          my $targetParameterConfigFile = $profile->{parameter}{config_file};
          $outputFilePath=sprintf("%s/%s", $pwd, $fileName);
          unlink($outputFilePath);
          generateInserts($srcProfileName, $targetProfileName, $targetProfileDesc, $targetParameterConfigFile, $fileName);
          InstallUtils::logger( "Wrote exported file: $outputFilePath\n", "debug" );
      }
    }


   # Success!
    $dbh->disconnect();
}

main;

# vi:syntax=perl
