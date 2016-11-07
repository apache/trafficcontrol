package Test::TrafficServerConfig;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


use strict;
use warnings;
use Carp qw/cluck/;

BEGIN {
	use Exporter;
	our @EXPORT_OK = qw{ loadConfigFile loadConfig };
}

my $keyval_re = qr/
	(\w+)=               # key=
	(
		"[^"]*"  |       # quoted string
		[^"\s]*          # unquoted value (no spaces)
	)
	(?:\s+|$)            # white space or end-of-line
	/x;

sub parseConfigLine {
	my $line = shift;
	my %h;

	$line =~ s/^\s+//;
	$line =~ s/\s+$//;

	while ( $line =~ /${keyval_re}\s*/g ) {
		my ( $k, $v ) = ( $1, $2 );

		# remove surrounding quotes if there
		$v =~ s/^"(.*)"$/$1/;
		if ( $k =~ /parent/ ) {
			$h{$k} = [ split /;/, $v ];
		}
		else {
			$h{$k} = $v;
		}
	}
	return \%h;
}

sub parseConfig {
	my $lines = shift;
	my $uq    = qr/(?<!\\)"/;
	my @config;
	for my $line ( split /\n/, $lines ) {
		next if $line =~ /^\s*#/;
		push @config, parseConfigLine($line);
	}
	return \@config;
}

sub loadConfigFile {
	my $cf = shift;
	if ( !-f $cf ) {
		return {};
	}

	open my $cfh, '<', $cf or return {};

	local $/;    # slurp mode
	my $txt = <$cfh>;
	close $cfh;

	return parseConfig($txt);
}

1;
