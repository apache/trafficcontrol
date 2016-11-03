#!/usr/bin/perl
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

use LWP::Simple;
use JSON qw( decode_json );
use Data::Dumper;
use strict;
use warnings;

my $requrl = 'http://localhost:8080/request';

my $json = get( $requrl );
die "Could not get $requrl!" unless defined $json;
my $tables = decode_json( $json );
print Dumper($tables);
foreach my $table ( @{ $tables } ) {
	print $table . " => ";
	my $url = 'http://localhost:8080/api/' . $table . '?format=moosefixture&join=no';
	my $t = get($url);
	my $filename;
	if ($t =~ /package Fixtures::Integration::(\S+);/) {
		$filename = $1 . '.pm';
	}
	print $filename;
	open(FH, ">$filename");
	print FH $t;
	close(FH);
	print "\n";
}


