package MojoPlugins::Daemonize;
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
use POSIX qw(:sys_wait_h close setsid);

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		fork_and_daemonize => sub {
			my $self = shift;
			my $method = shift;

			# reap any finished child processes
			my $kid;
			do {
				$kid = waitpid(-1, WNOHANG);
				$self->app->log->debug("Reaping PID $kid");
			} while $kid > 0;

			my $pid  = fork();

			if ( !defined($pid) ) {
				$self->app->log->fatal("Unable to fork: $!");
				return (-1);
			}

			if ( $pid == 0 ) {
				$self->inactivity_timeout(0);
				POSIX::setsid();
			}

			return $pid;
		}
	);

}

1;
