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

package GenerateCert;

use strict;

use lib qw(/opt/traffic_ops/install/lib /opt/traffic_ops/lib/perl5 /opt/traffic_ops/app/lib);

use base qw{ Exporter };
our @EXPORT_OK = qw{ createCert };
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

use JSON;
use InstallUtils;
use File::Temp;
use Data::Dumper;
use File::Copy;
use InstallUtils qw{ :all };

my $ca       = "/etc/pki/tls/certs/localhost.ca";
my $csr      = "/etc/pki/tls/certs/localhost.csr";
my $cert     = "/etc/pki/tls/certs/localhost.crt";
my $cdn_conf = "/opt/traffic_ops/app/conf/cdn.conf";
my $key      = "/etc/pki/tls/private/localhost.key";
my $msg      = << 'EOF';

	We're now running a script to generate a self signed X509 SSL certificate.

EOF

sub writeCdn_conf {
    my $cdn_conf = shift;

    # listen param to be inserted
    my $listen_str = "https://[::]:443?cert=${cert}&key=${key}&ca=${ca}&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED";

    # load as perl hash to find string to be replaced
    my $cdnh = do $cdn_conf;
    if ( exists $cdnh->{hypnotoad} ) {
        $cdnh->{hypnotoad}{listen} = [$listen_str];
    }
    else {
        # add the whole hypnotoad config without affecting anything else in the config
        $cdnh->{hypnotoad} = {
            listen   => [$listen_str],
            user     => 'trafops',
            group    => 'trafops',
            pid_file => '/var/run/traffic_ops.pid',
            workers  => 48,
        };
    }

    # dump conf data in compact but readable form
    my $dumper = Data::Dumper->new( [$cdnh] );
    $dumper->Indent(1)->Terse(1)->Quotekeys(0);

    # write whole config to temp file in pwd (keeps in same filesystem)
    my $tmpfile = File::Temp->new( DIR => '.' );
    print $tmpfile $dumper->Dump();
    close $tmpfile;

    # make backup of current file
    my $backup_num = 0;
    my $backup_name;
    do {
        $backup_num++;
        $backup_name = "$cdn_conf.backup$backup_num";
    } while ( -e $backup_name );
    rename( $cdn_conf, $backup_name ) or die("rename(): $!");

    # rename temp file to cdn.conf and set ownership/permissions same as backup
    my @stats = stat($backup_name);
    my ( $uid, $gid, $perm ) = @stats[ 4, 5, 2 ];
    move( "$tmpfile", $cdn_conf ) or die("move(): $!");

    chown $uid, $gid, $cdn_conf;
    chmod $perm, $cdn_conf;
}

# execOpenssl takes a description of the command being done, and an array of arguments to OpenSSL,
# and tries to execute the command, on failure prompting the user to retry.
# The description should be capitalized, but not terminated with punctuation.
# Returns the OpenSSL exit code.
sub execOpenssl {
    my ( $description, @args ) = @_;
    InstallUtils::logger( $description, "info" );
    my $result = 1;
    while ( $result != 0 ) {
        $result = InstallUtils::execCommand( "openssl", @args );
        if ( $result != 0 ) {
            my $ans = "";
            while ( $ans !~ /^[yY]/ && $ans !~ /^[nN]/ ) {
                $ans = InstallUtils::promptUser( $description . " failed. Try again (y/n)", "y" );
            }
            if ( $ans =~ /^[nN]/ ) {
                return $result;
            }
        }
    }
    return $result;
}

sub createCert {

    # the file used for ssl configuration
    my $opensslconf = shift;

    if ( !defined $opensslconf ) {
        InstallUtils::logger( "No input file - running openssl configuration in interactive mode", "info" );
    }

    InstallUtils::logger( $msg, "info" );

    InstallUtils::logger( "Postinstall SSL Certificate Creation", "info" );

    my $params;
    my $passphrase;

    # load the parameters for the certificate
    if ( defined $opensslconf ) {
        my $config = InstallUtils::readJson($opensslconf);
        if ( defined $config->{country} ) {

            # the parameters to auto generate the certificate
            $params = "/C=$config->{country}/ST=$config->{state}/L=$config->{locality}/O=$config->{company}/OU=$config->{org_unit}/CN=$config->{common_name}/";

            $passphrase = $config->{rsaPassword};
        }
    }

    if ( execOpenssl( "Generating an RSA Private Server Key", "genrsa", "-des3", "-out", "server.key", "-passout", "pass:$passphrase", "1024" ) != 0 ) {
        exit 1;
    }
    InstallUtils::logger( "The server key has been generated", "info" );

    if ($params) {
        if ( execOpenssl( "Creating a Certificate Signing Request (CSR)", "req", "-new", "-key", "server.key", "-out", "server.csr", "-passin", "pass:$passphrase", "-subj", $params ) != 0 ) {
            exit 1;
        }
    }
    else {
        if ( execOpenssl( "Creating a Certificate Signing Request (CSR)", "req", "-new", "-key", "server.key", "-out", "server.csr", "-passin", "pass:$passphrase" ) != 0 ) {
            exit 1;
        }
    }

    InstallUtils::logger( "The Certificate Signing Request has been generated", "info" );

    InstallUtils::execCommand( "/bin/mv", "server.key", "server.key.orig" );

    if ( execOpenssl( "Removing the pass phrase from the server key", "rsa", "-in", "server.key.orig", "-out", "server.key", "-passin", "pass:$passphrase" ) != 0 ) {
        exit 1;
    }
    InstallUtils::logger( "The pass phrase has been removed from the server key", "info" );

    if ( execOpenssl( "Generating a Self-signed certificate", "x509", "-req", "-days", "365", "-in", "server.csr", "-signkey", "server.key", "-out", "server.crt" ) != 0 ) {
        exit 1;
    }
    InstallUtils::logger( "A server key and self signed certificate has been generated", "info" );

    InstallUtils::logger( "Installing the server key and server certificate", "info" );

    my $result = InstallUtils::execCommand( "/bin/cp", "server.key", "$key" );
    if ( $result != 0 ) {
        errorOut("Failed to install the private server key");
    }
    $result = InstallUtils::execCommand( "/bin/chmod", "600",             "$key" );
    $result = InstallUtils::execCommand( "/bin/chown", "trafops:trafops", "$key" );

    if ( $result != 0 ) {
        errorOut("Failed to install the private server key");
    }

    InstallUtils::logger( "The private key has been installed",     "info" );
    InstallUtils::logger( "Installing the self signed certificate", "info" );

    $result = InstallUtils::execCommand( "/bin/cp", "server.crt", "$cert" );

    if ( $result != 0 ) {
        errorOut("Failed to install the self signed certificate");
    }

    $result = InstallUtils::execCommand( "/bin/chmod", "600",             "$cert" );
    $result = InstallUtils::execCommand( "/bin/chown", "trafops:trafops", "$cert" );

    if ( $result != 0 ) {
        errorOut("Failed to install the self signed certificate");
    }

    InstallUtils::logger( "Saving the self signed csr", "info" );
    $result = InstallUtils::execCommand( "/bin/cp", "server.csr", "$csr" );

    if ( $result != 0 ) {
        errorOut("Failed to save the self signed csr");
    }
    $result = InstallUtils::execCommand( "/bin/chmod", "664",             "$csr" );
    $result = InstallUtils::execCommand( "/bin/chown", "trafops:trafops", "$csr" );

    writeCdn_conf($cdn_conf);

    my $msg = << 'EOF';

        The self signed certificate has now been installed. 

        You may obtain a certificate signed by a Certificate Authority using the
        server.csr file saved in the current directory.  Once you have obtained
        a signed certificate, copy it to /etc/pki/tls/certs/localhost.crt and
        restart Traffic Ops.

EOF

    InstallUtils::logger( $msg, "info" );

    return 0;
}
