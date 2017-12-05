#!/usr/bin/env perl

#
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
#
use Data::Dumper;
use File::Basename;
use Env qw/PERL5LIB/;
use Getopt::Long qw(GetOptions);

mkdir("log");
my @watch_dirs;

GetOptions('secure' => \$secure);

# Look in the PERL5LIB directories for any TrafficOpsRoutes files.
foreach my $dir (@INC) {
	if ( $dir =~ /traffic_ops_extensions/ ) {
		push( @watch_dirs, $dir );
	}
}

my $to_lib = dirname($0) . "/../lib";
push( @watch_dirs, "templates" );
push( @watch_dirs, $to_lib );

my $watch_dirs_arg = join( " -w ", @watch_dirs );
$watch_dirs = join( "\n", @watch_dirs );

print "Morbo will restart with changes to any of the following dirs:\n";
print "(also the order in which Traffic Ops Perl Libraries and Extension modules will be searched)";
print "\n$watch_dirs\n\n";

my $local_dir = dirname($0) . "/../local";
my $export    = 'export PERL5LIB=$PERL5LIB:' . $local_dir . '/lib/perl5/:' . $to_lib;
my $cmd;
if ($secure ) {
   $cmd       = "$export && " . $local_dir . "/bin/morbo --listen 'https://*:60443' -v $local_dir/../script/cdn -w $watch_dirs_arg";
 }
else {
   $cmd       = "$export && " . $local_dir . "/bin/morbo --listen 'http://*:3000' -v $local_dir/../script/cdn -w $watch_dirs_arg";
}

system($cmd);
