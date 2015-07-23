package Utils::Helper::TrafficOpsRoutesLoader;
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
#

use Data::Dumper;
use File::Find;
my $r;

sub new {
	my $self  = {};
	my $class = shift;
	$r = shift;
	return ( bless( $self, $class ) );
}

sub load {
	my $self = shift;

	$self->load_routes(".");
	my $to_extensions_lib = $ENV{"TO_EXTENSIONS_LIB"};
	$self->load_routes($to_extensions_lib);
}

sub load_routes {
	my $self     = shift;
	my $root_dir = shift;
	if ( defined($root_dir) ) {
		if ( -e $root_dir ) {
			my @file_list;
			find(
				sub {
					return unless -f;         #Must be a file
					return unless /\.pm$/;    #Must end with `.pm` suffix
					push @file_list, $File::Find::name;
				},
				$root_dir
			);

			foreach my $file (@file_list) {
				open my $fn, '<', $file;
				my $first_line = <$fn>;
				my ( $package_keyword, $package_name ) = ( $first_line =~ m/(package )(.*);/ );
				if ( $package_name =~ /TrafficOpsRoutes$/ ) {
					print "Loading Mojo routes from package: " . $package_name . "\n";
					eval "use $package_name;";
					my $routes_class = eval {$package_name};
					$routes_class->define($r) || die "Route failed to load from package '" . $package_name . "' interface improperly defined.\n";
				}
				close $fn;
			}
		}
	}
}

1;
