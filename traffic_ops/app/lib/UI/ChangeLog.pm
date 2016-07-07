package UI::ChangeLog;
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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use DBI;

sub changelog {
	my $self = shift;

	my $date_string = `date "+%Y-%m-%d% %H:%M:%S"`;
	chomp($date_string);
	$self->session( last_seen_log => $date_string );
	&navbarpage($self);
}

sub readlog {
	my $self    = shift;
	my $numdays = defined( $self->param('days') ) ? $self->param('days') : 30;
	my $rows    = defined( $self->param('days') ) ? 1000000 : 1000;              # all of them gets to be too much

	my @data;
	my $interval = "> now() - interval '" . $numdays . " day'";                  # postgres
	my $rs = $self->db->resultset('Log')->search( { 'me.last_updated' => \$interval },
		{ prefetch => [ { 'tm_user' => undef } ], order_by => { -desc => 'me.last_updated' }, rows => $rows } );
	while ( my $row = $rs->next ) {
		push(
			@data, {
				"id"           => $row->id,
				"level"        => $row->level,
				"message"      => $row->message,
				"user"         => $row->tm_user->username,
				"ticketnum"    => $row->ticketnum,
				"last_updated" => $row->last_updated,
			}
		);
	}

	# setting cookie in the lib/Cdn/alog sub - this will be cached
	# my $date_string = `date "+%Y-%m-%d% %H:%M:%S"`;
	# chomp($date_string);
	# $self->session( last_seen_log => $date_string );
	$self->render( json => \@data );
}

sub newlogcount {
	my $self   = shift;
	my $cookie = $self->cookie('last_seen_log');
	my $user   = $self->current_user()->{userid};

	my $count = 0;
	if ( !defined($cookie) ) {
		my $date_string = `date "+%Y-%m-%d% %H:%M:%S"`;
		chomp($date_string);
		$self->cookie( last_seen_log => $date_string, { path => "/", max_age => 604800 } );    # expires in a week.
	}
	else {
		my $since_string = "> \'" . $cookie . "\'";
		$count = $self->db->resultset('Log')->search( { -and => [ { tm_user => { '!=' => $user } }, { last_updated => \$since_string } ] }, )->count();
	}
	my $jdata = { newlogcount => $count };
	$self->render( json => $jdata );
}

sub createlog {
	my $self = shift;

	&log( $self, $self->param('message'), $self->param('level') );

	return $self->redirect_to('/misc');
}

1;
