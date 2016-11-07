package Utils::Helper::ResponseHelper;
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

use Carp qw(cluck confess);
use Data::Dumper;
use String::CamelCase qw(camelize);
use Scalar::Util 'reftype';
use JSON;

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	return ( bless( $self, $class ) );
}

#
#  Takes 12M Snake_case hash keys and converts them camelCase keys
#  for exposing that data to the API
#
sub camelcase_response_keys {
	my $self = shift || confess("Call on an instance of Utils::Helper");
	my @data = shift;
	my @camelcase_response;

	foreach my $data_element (@data) {
		foreach my $element (@$data_element) {
			my %camel_case_hash;
			my $reftype = reftype $element;
			if ( $reftype eq 'HASH' ) {
				foreach my $key ( keys %{$element} ) {
					my $k = lcfirst( camelize($key) );
					my $v = $element->{$key};
					$camel_case_hash{$k} = $v;
				}
			}

			push( @camelcase_response, \%camel_case_hash );

		}
	}

	return \@camelcase_response;
}

sub handle_response {
	my $self     = shift;
	my $response = shift;
	my $content  = shift;
	if ( ( $response->code eq '200' ) && defined($content) && $content ne "" ) {
		return ( decode_json($content) );
	}
	else {
		return (undef);
	}
}

1;
