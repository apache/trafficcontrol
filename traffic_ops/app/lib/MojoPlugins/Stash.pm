package MojoPlugins::Stash;
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
use Carp qw(cluck confess);

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(

		# stash the cdn array for a form select.. optional second arg is the one that is selected.
		stash_cdn_selector => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Stash");
			my $selected = shift || -1;

			my $rs = $self->db->resultset('Cdn')->search(undef);
			my @cdns;
			while ( my $row = $rs->next ) {
				if ( $row->id == $selected ) {
					push( @cdns, [ $row->name => $row->id, selected => 'true' ] );
				}
				else {
					push( @cdns, [ $row->name => $row->id ] );
				}
			}
			$self->stash( cdns => \@cdns );
		}
	);

	$app->renderer->add_helper(
		# stash the profile_type array for a form select. optional second arg is the one that is selected.
		stash_profile_type_selector => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Stash");
			my $selected = shift || -1;

			my $enum_possible = $self->enum_values("profile_type");
			my @types;
			foreach my $val ( @{$enum_possible} ) {
				if ( $val eq $selected ) {
					push( @types, [ $val => $val, selected => 'true' ] );
				}
				else {
					push( @types, [ $val => $val ] );
				}
			}
			$self->stash( profile_types => \@types );
		}
	);

	$app->renderer->add_helper(
		# stash the profile array for a form select.. optional second arg is the one that is selected.
		stash_profile_selector => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Stash");
			my $ptype    = shift || confess("Please supply a profile type");
			my $selected = shift || -1;

			my @profiles;
			if ( $selected == -1 ) {
				push( @profiles, [ "No Profile" => -1, selected => 'true' ] );
			}
			my $rs = $self->db->resultset('Profile')->search( { type => $ptype } );
			while ( my $row = $rs->next ) {
				if ( $row->id == $selected ) {
					push( @profiles, [ $row->name => $row->id, selected => 'true' ] );
				}
				else {
					push( @profiles, [ $row->name => $row->id ] );
				}
			}
			$self->stash( profiles => \@profiles );
		}
	);
}

1;
