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
$ENV{PATH}     = "/usr/bin:/usr/local/go/bin:/opt/traffic_ops/install/bin:$ENV{PATH}";

use strict;
use warnings;

use DBI;
use POSIX;
use File::Basename qw{dirname};
use File::Path qw{make_path};
use Crypt::ScryptKDF qw(scrypt_hash);
use Data::Dumper qw(Dumper);
use Scalar::Util qw(looks_like_number);
use Getopt::Long;

use InstallUtils qw{ :all };
use GenerateCert qw{ :all };
use Database qw{ connect };

# paths of the output configuration files
my $databaseConfFile = "/opt/traffic_ops/app/conf/production/database.conf";
my $dbConfFile       = "/opt/traffic_ops/app/db/dbconf.yml";
my $cdnConfFile      = "/opt/traffic_ops/app/conf/cdn.conf";
my $ldapConfFile     = "/opt/traffic_ops/app/conf/ldap.conf";
my $usersConfFile    = "/opt/traffic_ops/install/data/json/users.json";
my $profilesConfFile = "/opt/traffic_ops/install/data/profiles/";
my $opensslConfFile  = "/opt/traffic_ops/install/data/json/openssl_configuration.json";
my $paramConfFile    = "/opt/traffic_ops/install/data/json/profiles.json";

my $custom_profile_dir = $profilesConfFile . "custom";

# stores parameters for traffic ops config
my $parameters;

# location of traffic ops profiles
my $profileDir       = "/opt/traffic_ops/install/data/profiles/";
my $post_install_cfg = "/opt/traffic_ops/install/data/json/post_install.json";

# log file for the installer
my $logFile = "/var/log/traffic_ops/postinstall.log";

# debug mode
my $debug = 1;

# log file for cpan output
my $cpanLogFile = "/var/log/traffic_ops/cpan.log";

# maximum size the uncompressed log file should be before rotating it - rotating it copies the current log
#  file to the same name appended with .bkp replacing the old backup if any is there
my $maxLogSize = 10000000;    #bytes

# whether to create a config file with default values
my $dumpDefaults;

# configuration file output with answers which can be used as input to postinstall
my $outputConfigFile = "/opt/traffic_ops/install/bin/configuration_file.json";

my $inputFile = "";
my $automatic = 0;
my %defaultInputs;

# given a var to the hash of config_var and question, will return the question
sub getConfigQuestion {
    my $var = shift;
    foreach my $key ( keys %{ $var } ) {
        if ( $key ne "hidden" && $key ne "config_var" ) {
            return $key;
        }
    }
}

# question: The question given in the config file
# config_answer: The answer given in the config file - if no config file given will be defaultInput
# hidden: Whether or not the answer should be hidden from the terminal and logs, ex. passwords
#
# Determines if the script is being run in complete interactive mode and prompts user - otherwise
#  returns answer to question in config or defaults

sub getField {
    my $question      = shift;
    my $config_answer = shift;
    my $hidden        = shift;

    # if there is no config file and not in automatic mode prompt for all questions with default answers
    if ( !$inputFile && !$automatic ) {

        # if hidden then dont show password in terminal
        if ($hidden) {
            return InstallUtils::promptPasswordVerify($question);
        }
        else {
            return InstallUtils::promptUser( $question, $config_answer );
        }
    }

    return $config_answer;
}

# userInput: The entire input config file which is either user input or the defaults
# fileName: The name of the output config file given by the input config file
#
# Loops through an input config file and determines answers to each question using getField
#  and returns the hash of answers

sub getConfig {
    my %userInput = %{$_[0]}; shift;
    my $fileName  = shift;

    my %config;

    if ( !defined $userInput{$fileName} ) {
        InstallUtils::logger( "No $fileName found in config", "error" );
    }

    InstallUtils::logger( "===========$fileName===========", "info" );

    foreach my $var ( @{ $userInput{$fileName} } ) {
        my $question = getConfigQuestion($var);
        my $hidden   = $var->{"hidden"} if ( exists $var->{"hidden"} );
        my $answer   = $config{ $var->{"config_var"} } = getField( $question, $var->{$question}, $hidden );

        $config{ $var->{"config_var"} } = $answer;
        if ( !$hidden ) {
            InstallUtils::logger( "$question: $answer", "info" );
        }
    }
    return %config;
}

# userInput: The entire input config file which is either user input or the defaults
# dbFileName: The filename of the output config file for the database
# toDBFileName: The filename of the output config file for the Traffic Ops database
#
# Generates a config file for the database based on the questions and answers in the input config file

sub generateDbConf {
    my %userInput = %{$_[0]}; shift;
    my $dbFileName   = shift;
    my $toDBFileName = shift;

    my %dbconf = getConfig( \%userInput, $dbFileName );
    $dbconf{"description"} = "$dbconf{type} database on $dbconf{hostname}:$dbconf{port}";
    make_path( dirname($dbFileName), { mode => 0755 } );
    InstallUtils::writeJson( $dbFileName, \%dbconf );
    InstallUtils::logger( "Database configuration has been saved", "info" );

    # broken out into separate file/config area
    my %todbconf = getConfig( \%userInput, $toDBFileName );

    # Check if the Postgres db is used and set the driver to be "postgres"
    my $dbDriver = $dbconf{type};
    if ( $dbconf{type} eq "Pg" ) {
        $dbDriver = "postgres";
    }

    # No YAML library installed, but this is a simple file..
    open( my $fh, '>', $toDBFileName ) or errorOut("Can't write to $toDBFileName!");
    print $fh "production:\n";
    print $fh "    driver: $dbDriver\n";
    print $fh "    open: host=$dbconf{hostname} port=$dbconf{port} user=$dbconf{user} password=$dbconf{password} dbname=$dbconf{dbname} sslmode=disable\n";
    close $fh;

    return \%todbconf;
}

# userInput: The entire input config file which is either user input or the defaults
# fileName: The filename of the output config file
#
# Generates a config file for the CDN

sub generateCdnConf {
    my %userInput = %{$_[0]}; shift;
    my $fileName  = shift;

    my %cdnConfiguration = getConfig( \%userInput, $fileName );

    # First, read existing one -- already loaded with a bunch of stuff
    my $cdnConf;
    if ( -f $fileName ) {
        $cdnConf = InstallUtils::readJson($fileName) or errorOut("Error loading $fileName: $@");
    }
    if ( lc $cdnConfiguration{genSecret} =~ /^y(?:es)?/ ) {
        my @secrets;
        my $newSecret = InstallUtils::randomWord();

        if (defined($cdnConf->{secrets})) {
            @secrets   = @{ $cdnConf->{secrets} };
            $cdnConf->{secrets} = \@secrets;
            InstallUtils::logger( "Secrets found in cdn.conf file", "debug" );
        } else {
            $cdnConf->{secrets} = \@secrets;
            InstallUtils::logger( "No secrets found in cdn.conf file", "debug" );
        }
        unshift @secrets, InstallUtils::randomWord();
        if ( $cdnConfiguration{keepSecrets} > 0 && $#secrets > $cdnConfiguration{keepSecrets} - 1 ) {

            # Shorten the array to requested length
            $#secrets = $cdnConfiguration{keepSecrets} - 1;
        }
    }
    if (exists $cdnConfiguration{base_url}) {
        $cdnConf->{to}{base_url} = $cdnConfiguration{base_url};
    }
    if (exists $cdnConfiguration{port}) {
        $cdnConf->{"traffic_ops_golang"}{port} = $cdnConfiguration{port};
    }
    $cdnConf->{"traffic_ops_golang"}{"log_location_error"} = "/var/log/traffic_ops/error.log";
    $cdnConf->{"traffic_ops_golang"}{"log_location_event"} = "/var/log/traffic_ops/access.log";

    $cdnConf->{hypnotoad}{workers} = $cdnConfiguration{workers};
    #InstallUtils::logger("cdnConf: " . Dumper($cdnConf), "info" );
    InstallUtils::writeJson( $fileName, $cdnConf );
    InstallUtils::logger( "CDN configuration has been saved", "info" );
}

sub hash_pass {
	my $pass = shift;
	return scrypt_hash($pass, \64, 16384, 8, 1, 64);
}

# userInput: The entire input config file which is either user input or the defaults
# fileName: The filename of the output config file
#
# Generates an LDAP config file

sub generateLdapConf {
    my %userInput = %{$_[0]}; shift;
    my $fileName  = shift;
    my %ldapInput = %{@{$userInput{$fileName}}[0]};
    my $useLdap = $ldapInput{"Do you want to set up LDAP?"};

    if ( !lc $useLdap =~ /^y(?:es)?/ ) {
        InstallUtils::logger( "Not setting up ldap", "info" );
        return;
    }

    my %ldapConf = getConfig( \%userInput, $fileName );
    # convert any deprecated keys to the correct key name
    my %keys_converted = ( password => 'admin_pass', hostname => 'host' );
    for my $key (keys %ldapConf) {
        if ( exists $keys_converted{$key} ) {
            $ldapConf{ $keys_converted{$key} } = delete $ldapConf{$key};
        }
    }

    my @requiredKeys = qw{ host admin_dn admin_pass search_base search_query insecure ldap_timeout_secs };
    for my $k (@requiredKeys) {
        if (! exists $ldapConf{$k} ) {
            errorOut("$k is a required key in $fileName");
        }
    }

    delete $ldapConf{setupLdap};

    # do a very loose check of form -- 'host' must be hostname:port
    if ( $ldapConf{ host } !~ /^\S+:\d+$/ ) {
        errorOut("host in $fileName must be of form 'hostname:port'");
    }

    make_path( dirname($fileName), { mode => 0755 } );
    InstallUtils::writeJson( $fileName, \%ldapConf );
}

sub generateUsersConf {
    my %userInput = %{$_[0]}; shift;
    my $fileName  = shift;

    my %user = ();
    my %config = getConfig( \%userInput, $fileName );

    $user{username} = $config{tmAdminUser};
    $user{password} = hash_pass( $config{tmAdminPw} );

    InstallUtils::writeJson( $fileName, \%user );
    $user{password} = $config{tmAdminPw};
    return \%user;
}

sub generateProfilesDir {
    my %userInput = %{$_[0]}; shift;
    my $fileName  = shift;

    my $userIn = $userInput{$fileName};
}

sub generateOpenSSLConf {
    my %userInput = %{$_[0]}; shift;
    my $fileName  = shift;

    my %config = getConfig( \%userInput, $fileName );
    return \%config;
}

sub generateParamConf {
    my %userInput = %{$_[0]}; shift;
    my $fileName  = shift;

    my %config = getConfig( \%userInput, $fileName );
    InstallUtils::writeJson( $fileName, \%config );
    return \%config;
}

# check default values for missing config_var parameter
sub sanityCheckDefaults {
    foreach my $file ( ( keys %defaultInputs ) ) {
        foreach my $defaultValue ( @{ $defaultInputs{$file} } ) {
            my $question = getConfigQuestion(\%$defaultValue);

            my %defaultValueHash = %$defaultValue;
            if ( !defined $defaultValueHash{"config_var"}
                || $defaultValueHash{"config_var"} eq "" )
            {
                errorOut("Question '$question' in file '$file' has no config_var");
            }
        }
    }
}

# userInput: The entire input config file which is either user input or the defaults
#
# Checks the input config file against the default inputs. If there is a question located in the default inputs which
#  is not located in the input config file it will output a warning message.

sub sanityCheckConfig {
    my %userInput = %{$_[0]}; shift;
    my $diffs     = 0;

    foreach my $file ( ( keys %defaultInputs ) ) {
        if ( !defined $userInput{$file} ) {
            InstallUtils::logger( "File '$file' found in defaults but not config file", "warn" );
            @{$userInput{$file}} = [];
        }

        foreach my $defaultValue ( @{ $defaultInputs{$file} } ) {

            my $found = 0;
            foreach my $configValue ( @{ $userInput{$file} } ) {
                if ( $defaultValue->{"config_var"} eq $configValue->{"config_var"} ) {
                    $found = 1;
                }
            }

            # if the question is not found in the config file add it from defaults
            if ( !$found ) {
                my $question = getConfigQuestion($defaultValue);
                InstallUtils::logger( "Question '$question' found in defaults but not in '$file'", "warn" );

                my %temp;
                my $answer;
                my $hidden = exists $defaultValue->{"hidden"} && $defaultValue->{"hidden"} ? 1 : 0;

                # in automatic mode add the missing question with default answer
                if ($automatic) {
                    $answer = $defaultValue->{$question};
                    InstallUtils::logger( "Adding question '$question' with default answer " . ( $hidden ? "" : "'$answer'" ), "info" );
                }

                # in interactive mode prompt the user for answer to missing question
                else {
                    InstallUtils::logger( "Prompting user for answer", "info" );
                    if ($hidden) {
                        $answer = InstallUtils::promptPasswordVerify($question);
                    }
                    else {
                        $answer = InstallUtils::promptUser( $question, $defaultValue->{$question} );
                    }
                }

                %temp = (
                    "config_var" => $defaultValue->{"config_var"},
                    $question    => $answer
                );

                if ($hidden) {
                    $temp{"hidden"} .= "true";
                }

                push @{ $userInput{$file} }, \%temp;

                $diffs++;
            }
        }
    }

    InstallUtils::logger( "File sanity check complete - found $diffs difference" . ( $diffs == 1 ? "" : "s" ), "info" );
}

# A function which returns the default inputs data structure. These questions and answers will be used if there is no
#  user input config file or if there are questions in the input config file which do not have answers

sub getDefaults {
    return (
        $databaseConfFile => [
            {
                "Database type" => "Pg",
                "config_var"    => "type"
            },
            {
                "Database name" => "traffic_ops",
                "config_var"    => "dbname"
            },
            {
                "Database server hostname IP or FQDN" => "localhost",
                "config_var"                          => "hostname"
            },
            {
                "Database port number" => "5432",
                "config_var"           => "port"
            },
            {
                "Traffic Ops database user" => "traffic_ops",
                "config_var"                => "user"
            },
            {
                "Password for Traffic Ops database user" => "",
                "config_var"                             => "password",
                "hidden"                                 => "true"
            }
        ],
        $dbConfFile => [
            {
                "Database server root (admin) user" => "postgres",
                "config_var"                        => "pgUser"
            },
            {
                "Password for database server admin" => "",
                "config_var"                         => "pgPassword",
                "hidden"                             => "true"
            }
        ],
        $cdnConfFile => [
            {
                "Generate a new secret?" => "yes",
                "config_var"             => "genSecret"
            },
            {
                "Number of secrets to keep?" => "1",
                "config_var"                 => "keepSecrets"
            },
            {
                "Port to serve on?"          => "443",
                "config_var"                 => "port"
            },
            {
                "Number of workers?" => "12",
                "config_var"         => "workers"
            },
            {
                "Traffic Ops url?"   => "http://localhost:3000",
                "config_var"         => "base_url"
            },
            {
                "ldap.conf location? (default is /opt/traffic_ops/app/conf/ldap.conf)" => "",
                "config_var"         => "ldap_conf_location"
            }
        ],
        $ldapConfFile => [
            {
                "Do you want to set up LDAP?" => "no",
                "config_var"                  => "setupLdap"
            },
            {
                "LDAP server hostname" => "",
                "config_var"           => "host"
            },
            {
                "LDAP Admin DN" => "",
                "config_var"    => "admin_dn"
            },
            {
                "LDAP Admin Password" => "",
                "config_var"          => "admin_pass",
                "hidden"              => "true"
            },
            {
                "LDAP Search Base" => "",
                "config_var"       => "search_base"
            },
            {
                "LDAP Search Query" => "",
                "config_var"       => "search_query"
            },
            {
                "LDAP Skip TLS verify" => "",
                "config_var"       => "insecure"
            },
            {
                "LDAP Timeout Seconds" => "",
                "config_var"       => "ldap_timeout_secs"
            }
        ],
        $usersConfFile => [
            {
                "Administration username for Traffic Ops" => "admin",
                "config_var"                              => "tmAdminUser"
            },
            {
                "Password for the admin user" => "",
                "config_var"                  => "tmAdminPw",
                "hidden"                      => "true"
            }
        ],
        $profilesConfFile => [
            {
                "Add custom profiles?" => "no",
                "config_var"           => "custom_profiles"
            }
        ],
        $opensslConfFile => [
            {
                "Do you want to generate a certificate?" => "yes",
                "config_var"                             => "genCert"
            },
            {
                "Country Name (2 letter code)" => "",
                "config_var"                   => "country"
            },
            {
                "State or Province Name (full name)" => "",
                "config_var"                         => "state"
            },
            {
                "Locality Name (eg, city)" => "",
                "config_var"               => "locality"
            },
            {
                "Organization Name (eg, company)" => "",
                "config_var"                      => "company"
            },
            {
                "Organizational Unit Name (eg, section)" => "",
                "config_var"                             => "org_unit"
            },
            {
                "Common Name (eg, your name or your server's hostname)" => "",
                "config_var"                                            => "common_name"
            },
            {
                "RSA Passphrase" => "CHANGEME!!",
                "config_var"     => "rsaPassword",
                "hidden"         => "true"
            }
        ],
        $paramConfFile => [
            {
                "Traffic Ops url" => "https://localhost",
                "config_var"      => "tm.url"
            },
            {
                "Human-readable CDN Name.  (No whitespace, please)" => "kabletown_cdn",
                "config_var"                                        => "cdn_name"
            },
            {
                "DNS sub-domain for which your CDN is authoritative" => "cdn1.kabletown.net",
                "config_var"                                         => "dns_subdomain"
            }
        ],
    );
}

# carried over from old postinstall
#
# todbconf: The database configuration to be used
# opensslconf: The openssl configuration if any

sub setupDatabaseData {
    my $dbh = shift;
    my $adminconf = shift;
    my $paramconf = shift;
    InstallUtils::logger( "paramconf " . Dumper($paramconf), "info" );

    my $result;

    my $q = <<"QUERY";
    select exists(select 1 from pg_tables where schemaname = 'public' and tablename = 'tm_user')
QUERY

    my $stmt = $dbh->prepare($q);
    $stmt->execute();

    InstallUtils::logger( "Setting up the database data", "info" );
    my $tables_found;
    while ( my $row = $stmt->fetch() ) {
       $tables_found = $row->[0];
    }
    if ($tables_found) {
       InstallUtils::logger( "Found existing tables skipping table creation", "info" );
    } else  {
       invoke_db_admin_pl("load_schema");
    }
    invoke_db_admin_pl("migrate");
    invoke_db_admin_pl("seed");
    invoke_db_admin_pl("patch");

    # Skip the insert if the admin 'username' is already there.
    my $hashed_passwd = hash_pass( $adminconf->{"password"} );
    my $insert_admin = <<"ADMIN";
    insert into tm_user (username, tenant_id, role, local_passwd, confirm_local_passwd)
                values  ('$adminconf->{"username"}',
                        (select id from tenant where name = 'root'),
                        (select id from role where name = 'admin'),
                         '$hashed_passwd',
                        '$hashed_passwd' )
                        ON CONFLICT (username) DO NOTHING;
ADMIN
    $dbh->do($insert_admin);

    insert_cdn($dbh, $paramconf);
    insert_parameters($dbh, $paramconf);
    insert_profiles($dbh, $paramconf);


}

sub invoke_db_admin_pl {
    my $action    = shift;

    chdir("/opt/traffic_ops/app");
    my $result = InstallUtils::execCommand( "db/admin", "--env=production", $action );

    if ( $result != 0 ) {
        errorOut("Database $action failed");
    }
    else {
        InstallUtils::logger( "Database $action succeeded", "info" );
    }

    return $result;
}

sub setupCertificates {
    my $opensslconf      = shift;

    my $result;

    if ( lc $opensslconf->{"genCert"} =~ /^y(?:es)?/ ) {
        if ( -x "/usr/bin/openssl" ) {
            InstallUtils::logger( "Installing SSL Certificates", "info" );
            $result = GenerateCert::createCert($opensslconf);

            if ( $result != 0 ) {
                errorOut("SSL Certificate Installation failed");
            }
            else {
                InstallUtils::logger( "SSL Certificates have been installed", "info" );
            }
        }
        else {
            InstallUtils::logger( "Unable to install SSL certificates as openssl is not installed",                                     "error" );
            InstallUtils::logger( "Install openssl and then run /opt/traffic_ops/install/bin/generateCert to install SSL certificates", "error" );
            exit 4;
        }
    }
    else {
        InstallUtils::logger( "Not generating openssl certification", "info" );
    }
}

#------------------------------------
sub insert_cdn {

    my $dbh = shift;
    my $paramconf = shift;

    InstallUtils::logger( "=========== Setting up cdn", "info" );

    # Enable multiple inserts into one commit
    $dbh->{pg_server_prepare} = 0;

	my $cdn_name = $paramconf->{"cdn_name"};
	my $dns_subdomain = $paramconf->{"dns_subdomain"};

    my $insert_stmt = <<INSERTS;

    -- global parameters
    insert into cdn (name, domain_name, dnssec_enabled)
                values ('$cdn_name', '$dns_subdomain', false)
                ON CONFLICT (name) DO NOTHING;

INSERTS
    doInsert($dbh, $insert_stmt);
}

#------------------------------------
sub insert_parameters {
    my $dbh = shift;
    my $paramconf = shift;

    InstallUtils::logger( "=========== Setting up parameters", "info" );

    # Enable multiple inserts into one commit
    $dbh->{pg_server_prepare} = 0;

	my $tm_url = $paramconf->{"tm.url"};

    my $insert_stmt = <<INSERTS;
    -- global parameters
    insert into parameter (name, config_file, value)
                values ('tm.url', 'global', '$tm_url')
                ON CONFLICT (name, config_file, value) DO NOTHING;

    insert into parameter (name, config_file, value)
                values ('tm.infourl', 'global', '$tm_url/doc')
                ON CONFLICT (name, config_file, value) DO NOTHING;

    -- CRConfig.json parameters
    insert into parameter (name, config_file, value)
                values ('geolocation.polling.url', 'CRConfig.json', '$tm_url/routing/GeoLite2-City.mmdb.gz')
                ON CONFLICT (name, config_file, value) DO NOTHING;

    insert into parameter (name, config_file, value)
                values ('geolocation6.polling.url', 'CRConfig.json', '$tm_url/routing/GeoLiteCityv6.dat.gz')
                ON CONFLICT (name, config_file, value) DO NOTHING;

INSERTS
    doInsert($dbh, $insert_stmt);
}

#------------------------------------
sub insert_profiles {
    my $dbh = shift;
    my $paramconf = shift;

    InstallUtils::logger( "\n=========== Setting up profiles", "info" );
	my $tm_url = $paramconf->{"tm.url"};

    my $insert_stmt = <<INSERTS;

    -- global parameters
    insert into profile (name, description, type, cdn)
                values ('GLOBAL', 'Global Traffic Ops profile, DO NOT DELETE', 'UNK_PROFILE',  (SELECT id FROM cdn WHERE name='ALL'))
                ON CONFLICT (name) DO NOTHING;

    insert into profile_parameter (profile, parameter)
                values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.url' and config_file = 'global' and value = '$tm_url') )
                ON CONFLICT (profile, parameter) DO NOTHING;

    insert into profile_parameter (profile, parameter)
                values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'tm.infourl' and config_file = 'global' and value = '$tm_url/doc') )
                ON CONFLICT (profile, parameter) DO NOTHING;

    insert into profile_parameter (profile, parameter)
                values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation.polling.url' and config_file = 'CRConfig.json' and value = '$tm_url/routing/GeoLite2-City.mmdb.gz') )
                ON CONFLICT (profile, parameter) DO NOTHING;

    insert into profile_parameter (profile, parameter)
                values ( (select id from profile where name = 'GLOBAL'), (select id from parameter where name = 'geolocation6.polling.url' and config_file = 'CRConfig.json' and value = '$tm_url/routing/GeoLiteCityv6.dat.gz') )
                ON CONFLICT (profile, parameter) DO NOTHING;

INSERTS
    doInsert($dbh, $insert_stmt);
}

#------------------------------------
sub doInsert {
    my $dbh = shift;
    my $insert_stmt = shift;

    InstallUtils::logger( "\n" . $insert_stmt, "info" );
    my $stmt = $dbh->prepare($insert_stmt);
    $stmt->execute();
}



# -cfile     - Input File:       The input config file used to ask and answer questions
# -a         - Automatic mode:   If there are questions in the config file which do not have answers, the script
#                                will look to the defaults for the answer. If the answer is not in the defaults
#                                the script will exit
# -defaults  - Defaults:         Writes out a configuration file with defaults which can be used as input
# -debug     - Debug Mode:       More output to the terminal
# -h         - Help:             Basic command line help menu

sub main {
    my $help = 0;

    # help string
    my $usageString = "Usage: postinstall [-a] [-debug] [-defaults[=<outfile]] [-r] -cfile=[config_file]\n";

    GetOptions(
        "cfile=s"     => \$inputFile,
        "automatic"   => \$automatic,
        "defaults:s"  => \$dumpDefaults,
        "debug"       => \$debug,
        "help"        => \$help
    ) or die($usageString);

    # stores the default questions and answers
    %defaultInputs = getDefaults();

    if ($help) {
        print $usageString;
        return;
    }

    # check if the user running postinstall is root
    if ( $> != 0 ) {
        errorOut("You must run this script as the root user");
    }

    InstallUtils::initLogger( $debug, $logFile );

    print("unzipping log\n");
    if ( -f "$logFile.gz" ) {
        InstallUtils::execCommand( "/bin/gunzip", "-f", "$logFile.gz" );
    }

    InstallUtils::logger( "Starting postinstall", "info" );

    InstallUtils::logger( "Debug is on", "info" );

    if ($automatic) {
        InstallUtils::logger( "Running in automatic mode", "info" );
    }

    if (defined $dumpDefaults) {
        # -defaults flag provided.
        if ($dumpDefaults ne "") {
	    # -defaults=<filename>  -- if -defaults without a file name, use the default.
	    # dumpDefaults with value -- use that as output file name
	    $outputConfigFile = $dumpDefaults;
        }
        InstallUtils::logger( "Writing default configuration to $outputConfigFile", "info" );
        InstallUtils::writeJson( $outputConfigFile, %defaultInputs );
        return;
    }

    InstallUtils::rotateLog($cpanLogFile);

    if ( -s $logFile > $maxLogSize ) {
        InstallUtils::logger( "Postinstall log above max size of $maxLogSize bytes - rotating", "info" );
        rotateLog($logFile);
    }

    # used to store the questions and answers provided by the user
    my %userInput;

    # if no input file provided use the defaults
    if ( $inputFile eq "" ) {
        InstallUtils::logger( "No input file given - using defaults", "info" );
        %userInput = %defaultInputs;
    }
    else {
        InstallUtils::logger( "Using input file $inputFile", "info" );

        # check if the input file exists
        errorOut("File '$inputFile' not found") if ( !-f $inputFile );

        # read and store the input file
        %userInput = %{InstallUtils::readJson($inputFile)};
    }

    # sanity check the defaults if running them automatically
    sanityCheckDefaults();

    # check the input config file against the defaults to check for missing questions
    sanityCheckConfig(\%userInput) if ( $inputFile ne "" );

    chdir("/opt/traffic_ops/install/bin");

    # The generator functions handle checking input/default/automatic mode
    # todbconf will be used later when setting up the database
    my $todbconf = generateDbConf( \%userInput, $databaseConfFile, $dbConfFile );
    generateLdapConf( \%userInput, $ldapConfFile );
    my $adminconf = generateUsersConf( \%userInput, $usersConfFile );
    my $custom_profile = generateProfilesDir( \%userInput, $profilesConfFile );
    my $opensslconf = generateOpenSSLConf( \%userInput, $opensslConfFile );
    my $paramconf = generateParamConf( \%userInput, $paramConfFile );

    if ( !-f $post_install_cfg ) {
        InstallUtils::writeJson( $post_install_cfg, {} );
    }

    setupCertificates( $opensslconf );
    generateCdnConf( \%userInput, $cdnConfFile );

    my $dbh = Database::connect($databaseConfFile, $todbconf);
    if (!$dbh) {
        InstallUtils::logger("Can't connect to the database.  Use the script `/opt/traffic_ops/install/bin/todb_bootstrap.sh` on the db server to create it and run `postinstall` again.", "error");
        exit(-1);
    }

    setupDatabaseData( $dbh, $adminconf, $paramconf );

    InstallUtils::logger("Starting Traffic Ops", "info" );
    InstallUtils::execCommand("/sbin/service traffic_ops restart");

    InstallUtils::logger("Waiting for Traffic Ops to restart", "info" );

    InstallUtils::logger("Success! Postinstall complete.");

    #InstallUtils::logger("Zipping up $logFile to $logFile.gz");
    #InstallUtils::execCommand( "/bin/gzip", "$logFile" );

   # Success!
    $dbh->disconnect();
}

main;

# vi:syntax=perl
