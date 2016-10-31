package Utils::Helper::TrafficOpsRoutesLoader;
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
#

use Data::Dumper;
use File::Find;
my $r;
my @loaded_route_packages;

sub new {
	my $self  = {};
	my $class = shift;
	$r = shift;
	return ( bless( $self, $class ) );
}

sub load {
	my $self   = shift;
	my $dashes = "-------------------------------------------------------------\n";
	print $dashes;

	# Look in the PERL5LIB directories for any TrafficOpsRoutes files.
	#print "PERL5LIB: " . Dumper(@INC);
	foreach my $dir (@INC) {
		$self->load_routes($dir);
	}
	print $dashes;
}

sub load_routes {
	my $self     = shift;
	my $root_dir = shift;

	#print "root_dir #-> (" . $root_dir . ")\n";
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
				if ( defined($package_name) ) {

					if ( $package_name =~ /.*TrafficOpsRoutes$/ ) {
						if ( !grep { $_ eq $package_name } @loaded_route_packages ) {
							print "Loading Mojo Routes from package: " . $package_name . "\n";
							eval "use $package_name;";
							my $routes_class = eval {$package_name};
							$routes_class->define($r) || die "Route failed to load from package '" . $package_name . "' interface improperly defined.\n";
							push( @loaded_route_packages, $routes_class );
						}

					}
					close $fn;
				}
			}
		}
	}
}

1;
