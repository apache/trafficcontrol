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
use POSIX qw(close setsid);

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(
		# Note: Calling fork_and_daemonize() returns twice: Once for the parent and the other for the daemon (=child).
		# Caller should check return value:
		#  <0 means an error occured
		#  0 means you are the daemon (=child)
		#  1 means you are the parent (= original process)
		fork_and_daemonize => sub {
			my $self = shift;
			my $pid  = fork();

			if ( !defined($pid) ) {
				$self->app->log->fatal("fork_and_daemonize(): Parent unable to fork: $!");
				return -1;
			}

			if ( $pid == 0 ) {
				# This is the first child
				$self->inactivity_timeout(0);
				POSIX::setsid();
				# First child forks daemon and exits with a value that signals the parent how the fork went
    			my $pid  = fork();

	    		if ( !defined($pid) ) {
					$self->app->log->fatal("fork_and_daemonize(): Child unable to fork: $!");
					# Exit with -1 to let parent know that the fork failed
					exit(-1);
				}
				if ($pid > 0) {
					# First child: Fork was OK, exit with 0
					exit(0);
				}
				# This is the daemon, return 0 to caller
				return 0;
			}

			# This is the parent. Wait for first child to exit
			my $rc = waitpid($pid, 0);
			if ($rc != $pid) {
				$self->app->log->fatal("fork_and_daemonize(): Parent waitpid($pid) returned $rc, expecting $pid. $!");
				return -1;
			}
			$rc = ${^CHILD_ERROR_NATIVE};
			if ($rc) {
				$self->app->log->fatal("fork_and_daemonize(): First child exited with $rc. $!");
				return -1;
			}

			# Parent: Do not return $pid as this is the pid of the first child, which is not interesting
			# Return 1 to signal caller that this is the parent
			return 1;
		}
	);

}

1;