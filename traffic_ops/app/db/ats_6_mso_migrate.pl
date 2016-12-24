#!/usr/bin/perl
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

# Script to create profiles / parameters from the header_rewrite settings to move to ATS 6.2
# support.
#
# Please be careful using this script, it was created for a very specific usecase, and only
# tested our env - JvD

use strict;
use warnings;
use DBI;

my $driver   = "Pg";
my $database = $ARGV[0];
my $userid   = $ARGV[1];
my $password = $ARGV[2];

my $dsn = "DBI:$driver:dbname=$database;host=127.0.0.1;port=5432";
my $dbh = DBI->connect( $dsn, $userid, $password, { RaiseError => 1 } )
	or die $DBI::errstr;

print "Opened database successfully\n";

my $sql = 'select id,xml_id,mid_header_rewrite from deliveryservice where multi_site_origin=true;';

my $sth = $dbh->prepare_cached($sql);
$sth->execute || die "Couldn't execute statement: " . $sth->errstr;
while ( my @data = $sth->fetchrow_array() ) {
	if ( !defined( $data[2] ) ) {
		next;
	}
	my @lines = split( /__RETURN__/, $data[2] );
	my %profile_ids;
	foreach my $line (@lines) {
		if ( !defined( $profile_ids{ $data[1] } ) ) {
			my $insp = $dbh->prepare('INSERT INTO PROFILE ("name", "description") VALUES(?, ?);');
			$insp->bind_param( 1, "MIDMSO_" . $data[1] );
			$insp->bind_param( 2, "Profile for " . $data[1] . " MSO settings" );
			$insp->execute();
			my $profile_id = $dbh->last_insert_id( undef, undef, "profile", undef );
			$profile_ids{ $data[1] } = $profile_id;
		}
		if ( $line =~ /set-config proxy.config.http.parent_origin/ ) {
			my $setting = $line;
			$setting =~ s/set-config //;
			print $data[1] . " ->" . $setting . "\n";
			$setting =~ s/^ *//;
			my ( $name, $value ) = split( /\s+/, $setting );
			print $name . " -> " . $value . "\n";
			my $insh = $dbh->prepare('INSERT INTO PARAMETER ("name", "config_file", "value") VALUES (?, ?, ?);');
			$insh->bind_param( 1, $name );
			$insh->bind_param( 2, 'parent.config' );
			$insh->bind_param( 3, $value );
			$insh->execute();
			my $param_id = $dbh->last_insert_id( undef, undef, "parameter", undef );
			print "Last inserted: " . $param_id;
			my $inspp = $dbh->prepare('INSERT INTO PROFILE_PARAMETER ("parameter", "profile") VALUES (?, ?);');
			$inspp->bind_param( 1, $param_id );
			$inspp->bind_param( 2, $profile_ids{ $data[1] } );
			$inspp->execute();
		}
	}
}

$dbh->disconnect();
