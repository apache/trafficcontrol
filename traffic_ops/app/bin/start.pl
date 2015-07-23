#!/usr/bin/env perl

#
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
#
use Data::Dumper;

mkdir("log");
my @watch_dirs = "lib";

# Look in the PERL5LIB directories for any TrafficOpsRoutes files.
#print "PERL5LIB: " . Dumper(@INC);
foreach my $dir (@INC) {
	if ( $dir =~ /traffic_ops_extensions/ ) {
		push( @watch_dirs, $dir );
	}
}
my $watch_dirs_arg = join( " -w ", @watch_dirs );
$watch_dirs = join( ":", @watch_dirs );
print "Morbo will restart with changes to-> $watch_dirs\n";

my $cmd = "local/bin/morbo --listen 'http://*:3000' -v script/cdn $watch_dirs_arg";
system($cmd);
