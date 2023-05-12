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

# Check the cdn.conf for the cert and key file references -- abort if they don't match what's defined here
# This normally wouldn't happen unless the user modified the cdn.conf to reference different file names, and in that
# case, they're probably generating certs outside of this anyway: this check is just here for safety..
sub checkCdnConf {
	my $cdn_conf = shift;
	my $conf;
	# load cdn.conf
	{
		local $/;  # slurp mode
		open my $fh, '<', $cdn_conf or die "Cannot load $cdn_conf\n";
		$conf = decode_json(scalar <$fh>);
	}

	my $key_conf = $conf->{key};
	my $cert_conf = $conf->{cert};
	my $msg;

	if (!defined cert_conf) {
		my $msg = <<"EOF";
	The "cert" portion of $cdn_conf is missing from $cdn_conf.
	Please ensure it contains the same structure as the one originally installed.
EOF
	}

	if (!defined $key_conf) {
		my $msg = <<"EOF";
	The "key" portion of $cdn_conf is missing from $cdn_conf.
	Please ensure it contains the same structure as the one originally installed.
EOF
	}

	if ($cert_conf !~ m@cert=$cert@ || $key_conf !~ m@key=$key@) {
		$msg = << "EOF";
	The "cert and key" portion of $cdn_conf is:
	$cert_conf $key_conf
	and does not reference the same "cert=" and "key=" values as are created here.
	Please modify $cdn_conf to add the following as fields:
	cert: $cert, key: $key
EOF
	}

	return $msg;
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


# creates a certificate with parameters used from postinstall passed into $config
sub createCert {

    # the file used for ssl configuration
    my $config = shift;

    InstallUtils::logger( $msg, "info" );

    InstallUtils::logger( "Postinstall SSL Certificate Creation", "info" );

    my $params;
    my $passphrase;

    # create the string of parameters
    $params = "/C=$config->{country}/ST=$config->{state}/L=$config->{locality}/O=$config->{company}/OU=$config->{org_unit}/CN=$config->{common_name}/";

    $passphrase = $config->{rsaPassword};

    InstallUtils::logger( "The server key has been generated", "info" );

    if ( execOpenssl( "Generating an RSA Private Server Key", "genrsa", "-des3", "-out", "server.key", "-passout", "pass:$passphrase", "1024" ) != 0 ) {
        exit 1;
    }
    if ( execOpenssl( "Creating a Certificate Signing Request (CSR)", "req", "-new", "-key", "server.key", "-out", "server.csr", "-passin", "pass:$passphrase", "-subj", $params ) != 0 ) {
        exit 1;
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


    my $msg = << 'EOF';

        The self signed certificate has now been installed. 

        You may obtain a certificate signed by a Certificate Authority using the
        server.csr file saved in the current directory.  Once you have obtained
        a signed certificate, copy it to /etc/pki/tls/certs/localhost.crt and
        restart Traffic Ops.

EOF

    InstallUtils::logger( $msg, "info" );
    my $error = checkCdnConf($cdn_conf);
    if ($error) {
	    errorOut( $error, "error ");
	    exit 1;
    }

    return 0;
}
