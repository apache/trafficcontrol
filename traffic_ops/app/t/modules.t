use strict;
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
use warnings;
use Test::More;
use File::Find;
use lib "lib";

# execute this from the traffic_ops/app directory (e.g.: prove t/modules.t)
my $number_of_tests = 0;
my %modules;
my %ignore = ( constant => 1, );

find( { wanted => \&examine, no_chdir => 1 }, ("lib") );

for my $module ( sort( keys(%modules) ) ) {
	$number_of_tests++;

	eval "use $module;";

	is( $@, "", "use $module, found in: " . join( ", ", sort( keys( %{ $modules{$module} } ) ) ) );
}
done_testing($number_of_tests);

sub examine {
	if ( -f "$File::Find::name" ) {
		note("Examining $File::Find::name");
		open( IN, "< $File::Find::name" ) || die("Unable to open $File::Find::name: $!");

		while (<IN>) {
			if ( $_ =~ m/^use (.*?)\;$/ ) {
				if ( !exists( $ignore{$1} ) ) {
					$modules{$1}{$File::Find::name} = 1;
				}
			}
		}

		close(IN);
	}
	else {
		note("Skipping $File::Find::name: $!");
	}
}
