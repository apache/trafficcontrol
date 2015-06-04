package Utils::Helper::ExtensionsHelper;
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
use Carp qw(cluck confess);
use Data::Dumper;
use JSON;

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	return ( bless( $self, $class ) );
}

sub prepend_extensions {
	my $to_ext_lib = $ENV{"TO_EXTENSIONS_LIB"};
	if ( defined($to_ext_lib) ) {
		if ( -e $to_ext_lib ) {
			unshift( @INC, $to_ext_lib );
			print "Found TO_EXTENSIONS_LIB prepending to library path: $to_ext_lib\n";
		}
		else {
			print "\nWARNING TO_EXTENSIONS_LIB environment variable is defined as $to_ext_lib but does not exist.\n\n";
		}
	}
}

1;
