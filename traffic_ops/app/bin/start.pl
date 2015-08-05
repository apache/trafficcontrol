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
use File::Basename;

mkdir("log");
my @watch_dirs;

# Look in the PERL5LIB directories for any TrafficOpsRoutes files.
#print "PERL5LIB: " . Dumper(@INC);
foreach my $dir (@INC) {
	if ( $dir =~ /traffic_ops_extensions/ ) {
		push( @watch_dirs, $dir );
	}
}

push( @watch_dirs, qw(lib) );

#BEGIN { my $local_lib = "local/lib/perl5" }
#push( @watch_dirs, $local_lib );
my $watch_dirs_arg = join( " -w ", @watch_dirs );
$watch_dirs = join( "\n", @watch_dirs );

print "Morbo will restart with changes to any of the following dirs:\n";
print "(also the order in which Traffic Ops Perl Libraries and Extension modules will be searched)";
print "\n$watch_dirs\n\n";

my $cmd = "export PERL5LIB=lib:local/lib/perl5;local/bin/morbo --listen 'http://*:3000' -v script/cdn -w $watch_dirs_arg";

#print "cmd #-> (" . $cmd . ")\n";
system($cmd);
