package MojoPlugins::Validation;
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
use Data::Dumper;
use Data::Dump qw(dump);
use Email::Valid;

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		is_username_in_db => sub {
			my $self     = shift;
			my $username = shift;
			my $is_in_db;
			my @users = $self->db->resultset('TmUser')->search( undef, { columns => [qw/username/] } );
			my $fcns ||= MojoPlugins::Validation::Functions->new( $self, @_ );
			return $fcns->is_username_in_db_list( @users, $username );
		}
	);
	$app->renderer->add_helper(
		is_username_taken => sub {
			my $self     = shift;
			my $username = shift;

			my $db_username = $self->db->resultset('TmUser')->find( { username => $username } );
			if ( defined($db_username) ) {
				$self->field('tm_user.username')->is_like( qr/$db_username/, "Username is already taken" );
			}
		}
	);

	$app->renderer->add_helper(
		is_email_taken => sub {
			my $self = shift;

			my $email = $self->param('tm_user.email');
			$self->app->log->debug( "email #-> " . Dumper($email) );
			my $dbh = $self->db->resultset('TmUser')->search( { email => $email } );
			my $count = $dbh->count();
			if ( $count > 0 ) {

				#$self->app->log->debug( "email #-> " . Dumper($email) );
				my $db_email = $dbh->single()->email;
				$self->field('tm_user.email')->is_like( qr/ . $db_email . /, "Email is already taken" );
			}
		}
	);

	$app->renderer->add_helper(
		is_email_format_valid => sub {
			my $self  = shift;
			my $email = $self->param('tm_user.email');

			#$self->app->log->debug( "valid email #-> " . Dumper($email) );
			if ( defined($email) ) {
				unless ( Email::Valid->address( -address => $email, -mxcheck => 1 ) ) {
					$self->field('tm_user.email')->is_like( qr/ . $email . /, "Email is not a valid format" );
				}

			}
		}
	);

	$app->renderer->add_helper(
		is_password_uncommon => sub {
			my $self  = shift;
			my $pass = $self->param('tm_user.local_passwd');
			my $blacklist = $self->app->{invalid_passwords};
			if ( defined($pass) && defined($blacklist->{$pass}) ) {
				$self->field('tm_user.local_passwd')->is_like( qr/ . $pass . /, "Password is too common." );
			}
		}
	);
}

package MojoPlugins::Validation::Functions;

use Mojo::Base -strict;
use Scalar::Util;
use Carp ();
use Validate::Tiny;

sub new {
	my $class = shift;
	my ( $c, $object ) = @_;
	my $self = bless {
		c      => $c,
		object => $object
	}, $class;

	Scalar::Util::weaken $self->{c};
	return $self;
}

sub is_username_in_db_list {
	my $self             = shift;
	my $inbound_username = shift;
	my @users            = shift;
	my $is_username_in_db_list;

	# This was a very tricky validation because there was no support in the Mojolicious::Plugin::FormFields for
	# if a string 'is in' list (the documented 'is_in' should really be called 'is_not_in' because that is
	# what it tests.
	my @user_names;
	foreach my $user_name (@users) {
		push( @user_names, $user_name->username );
	}

	foreach my $user (@users) {
		my $username = $user->username;
		if ( $username eq $inbound_username ) {
			$is_username_in_db_list = 'true';
			last;
		}
	}
	return $is_username_in_db_list;
}

1;
