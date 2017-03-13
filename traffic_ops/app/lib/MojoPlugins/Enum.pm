package MojoPlugins::Enum;
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

use Mojo::Base 'Mojolicious::Plugin';

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(

		# Returns an array with the possible values for the enum.
		# Note that for now this is postgres specific; if we need to support other databases, we need to add support here.
		enum_values => sub {
			my $self = shift;
			my $enum_name = shift;

			print ">>> " . $enum_name . "\n";
			my %views = (  # to add more enums, just add the key, val pair here.
				'profile_type', 'ProfileTypeValue'
			 ); 
			my @possible = $self->db->resultset( $views{$enum_name })->search(undef)->get_column('value')->all();
			return \@possible;
		}
	);
}

1;