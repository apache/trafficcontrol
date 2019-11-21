package MojoPlugins::Email;
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
use Data::Dumper;

use UI::Utils;
use POSIX qw(strftime);

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		send_deliveryservice_request => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Email");
			my $email_to = shift || confess("Please provide an email to send this request to.");
			my $details  = shift || confess("Please provide deliveryservice request details.");

			my $instance_name =
				$self->db->resultset('Parameter')->search( { -and => [ name => 'tm.instance_name', config_file => 'global' ] } )->get_column('value')
				->single();
			$self->stash( instance_name => $instance_name );

			my $current_user = $self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->single();
			if ($current_user) {
				$self->stash(
					{
						username  => $current_user->username,
						full_name => $current_user->full_name,
						email     => $current_user->email,
					}
				);
			}
			else {
				$self->stash(
					{
						username  => $self->current_user()->{username},
						full_name => undef,
						email     => undef,
					}
				);
			}

			$self->stash( params => $details );

			my $rc;
			$rc = $self->mail(
				subject  => $instance_name . " Delivery Service Request",
				to       => $email_to,
				template => 'delivery_service/request',
				format   => 'mail'
			);

			return $rc;
		}
	);

	$app->renderer->add_helper(
		send_password_reset_email => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Email");
			my $email_to = shift || confess("Please supply an email address.");
			my $token    = shift || confess("Please supply a GUID token");

			my $portal_pass_reset_url	= $self->config->{'portal'}{'base_url'} . $self->config->{'portal'}{'pass_reset_path'};
			my $portal_email_from		= $self->config->{'portal'}{'email_from'};

			$self->app->log->info( "MOJO_CONFIG: " . $ENV{MOJO_CONFIG} );
			$self->app->log->info( "portal_pass_reset_url: " . $portal_pass_reset_url );

			my $tm_user = {
				email			=> $email_to,
				portal_pass_reset_url	=> $portal_pass_reset_url,
				token			=> $token,
			};

			my $instance_name =
				$self->db->resultset('Parameter')->search( { -and => [ name => 'tm.instance_name', config_file => 'global' ] } )->get_column('value')
				->single();
			$self->stash( instance_name => $instance_name );

			my $rc;
			$self->stash( tm_user => $tm_user, fbox_layout => 1, mode => 'add' );

			if ( defined($email_to) ) {
				if ( defined($portal_email_from) ) {
					$rc = $self->mail(
						subject  => $instance_name . " Password Reset Request",
						from     => $portal_email_from,
						to       => $email_to,
						template => 'user/reset_password',
						format   => 'mail'
					);
				}
				else {
					$rc = $self->mail(
						subject  => $instance_name . " Password Reset Request",
						to       => $email_to,
						template => 'user/reset_password',
						format   => 'mail'
					);
				}

			}

			return $rc;
		}
	);

	$app->renderer->add_helper(
		send_registration_email => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Email");
			my $email_to = shift || confess("Please supply an email address.");
			my $token    = shift || confess("Please supply a GUID token");

			my $portal_user_register_url	= $self->config->{'portal'}{'base_url'} . $self->config->{'portal'}{'user_register_path'};
			my $portal_email_from		= $self->config->{'portal'}{'email_from'};

			$self->app->log->info( "MOJO_CONFIG: " . $ENV{MOJO_CONFIG} );
			$self->app->log->info( "portal_user_register_url: " . $portal_user_register_url );

			my $tm_user = {
				email				=> $email_to,
				portal_user_register_url	=> $portal_user_register_url,
				token				=> $token,
			};

			my $instance_name =
				$self->db->resultset('Parameter')->search( { -and => [ name => 'tm.instance_name', config_file => 'global' ] } )->get_column('value')
				->single();

			$self->stash( instance_name => $instance_name );
			$self->stash( tm_user => $tm_user, fbox_layout => 1, mode => 'add' );
			if ( defined($email_to) ) {
				if ( defined($portal_email_from) ) {
					$self->mail(
						subject  => "Welcome to the " . $instance_name,
						from     => $portal_email_from,
						to       => $email_to,
						template => 'user/registration',
						format   => 'mail'
					);
				}
				else {
					$self->mail( subject => "Welcome to the " . $instance_name, to => $email_to, template => 'user/registration', format => 'mail' );
				}

				my $email_notice = 'Successfully sent registration email to: ' . $email_to;
				$self->app->log->info($email_notice);
				$self->flash( message => $email_notice );
			}
		}
	);

	$app->renderer->add_helper(
		update_user_token => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Email");
			my $email_to = shift;
			my $token    = shift;

			$self->app->log->debug( "Resetting user for email: " . $email_to );

			my $new_id = -1;

			my $dbh = $self->db->resultset('TmUser')->find( { email => $email_to } );
			$dbh->token($token);
			$dbh->update();

			# if the insert has failed, we don't even get here, we go to the exception page.
			#&log( $self, "Reset password for user with email " . $email_to, "UICHANGE" );

		}
	);

	$app->renderer->add_helper(
		create_registration_user => sub {
			my $self     = shift || confess("Call on an instance of MojoPlugins::Email");
			my $email_to = shift;
			my $token    = shift;

			my $new_id = -1;
			my $portal_role_id = $self->db->resultset('Role')->find( { name => 'portal' } );

			my $now = strftime( "%Y-%m-%d %H:%M:%S", gmtime() );

			my $dbh = $self->db->resultset('TmUser')->find( { email => $email_to } );
			if ( defined($dbh) ) {
				$self->app->log->debug( "Updating token, found email: " . $email_to . "\n" );
				$dbh->token($token);
				$dbh->registration_sent($now);
				$dbh->update();
			}
			else {
				$self->app->log->debug("Email not found, creating a new tm_user. \n");

				#NOTE: The token is used as a temp password for user registration and is forced
				#      to be changed after the user logs in, so that mapping is not a mistake.
				$dbh = $self->db->resultset('TmUser')->create(
					{
						email             => $email_to,
						role              => $portal_role_id,
						username          => $token,
						token             => $token,
						registration_sent => $now,
					}
				);
				$dbh->insert();
			}

			# if the insert has failed, we don't even get here, we go to the exception page.
			&log( $self, "Created registration user with email " . $email_to, "UICHANGE" );

			$new_id = $dbh->id;
			if ( $new_id == -1 ) {
				my $referer = $self->req->headers->header('referer');
				return $self->redirect_to($referer);
			}
		}
	);
}

1;
