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

package Database;

use warnings;
use strict;

use base qw{ Exporter };
our @EXPORT_OK = qw{ connect };
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

#------------------------------------
sub connect {
    my $databaseConfFile = shift;
    my $todbconf = shift;

    my $conf = InstallUtils::readJson($databaseConfFile);

    # Check if the Postgres db is used and set the admin database to be "postgres"
    my $dbName = $conf->{type};
    if ( $conf->{type} eq "Pg" ) {
        $dbName = "traffic_ops";
    }

    $ENV{PGUSER}     = $conf->{"user"};
    $ENV{PGPASSWORD} = $conf->{"password"};

    my $dsn = sprintf( "DBI:%s:db=%s;host=%s;port=%d", $conf->{type}, $dbName, $conf->{hostname}, $conf->{port} );
    my $dbh = DBI->connect( $dsn, $todbconf->{"user"}, $todbconf->{"password"} );
    if ($dbh) {
        InstallUtils::logger( "Database connection succeeded", "info" );
    }
    else {
        # show error, but don't exit -- let the caller deal with it based on undef $dbh
        InstallUtils::logger( $DBI::errstr, "error" );
    }
    return $dbh;
}


1;
