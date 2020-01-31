package API::ServerCheck;

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
use Data::Dumper;

# sub index {
# 	my $self = shift;
# 	my $data = getservercheckdata($self);
# 	$self->render( json => $data );
# }

# source for the datatable
sub aadata {
	my $self  = shift;
	my %data  = ( "aaData" => undef );
	my $table = $self->param('table');

	# # and create the 'inside out' aadata table
	my %condition = ( 'name' => [ { -like => 'MID%' }, { -like => 'EDGE%' } ] );
	my $rs_type = $self->db->resultset('Type')->search( \%condition );
	my $rs =
		$self->db->resultset('Server')
		->search( { 'me.type' => { -in => $rs_type->get_column('id')->as_query } }, { prefetch => [ 'servercheck', 'status', 'profile' ]} );
	while ( my $server = $rs->next ) {
		if ( !defined $server || !defined $server->servercheck ) {
			next;
		}
		my @line = (
			$server->id,              $server->host_name,       $server->profile->name,   $server->status->name,    $server->upd_pending,
			$server->servercheck->aa, $server->servercheck->ab, $server->servercheck->ac, $server->servercheck->ad, $server->servercheck->ae,
			$server->servercheck->af, $server->servercheck->ag, $server->servercheck->ah, $server->servercheck->ai, $server->servercheck->aj,
			$server->servercheck->ak, $server->servercheck->al, $server->servercheck->am, $server->servercheck->an, $server->servercheck->ao,
			$server->servercheck->ap, $server->servercheck->aq, $server->servercheck->ar, $server->servercheck->at, $server->servercheck->au,
			$server->servercheck->av, $server->servercheck->aw, $server->servercheck->ax, $server->servercheck->ay, $server->servercheck->az,
			$server->servercheck->ba, $server->servercheck->bb, $server->servercheck->bc, $server->servercheck->bd, $server->servercheck->bd,
			$server->servercheck->be, $server->servercheck->bf,
		);
		push( @{ $data{'aaData'} }, \@line );
	}
	return $self->deprecation_with_no_alternative(200, \%data);
}

# read for not crazy Datatables
sub read {
	my $self = shift;

	my %condition = ( 'type.name' => [ { -like => 'MID%' }, { -like => 'EDGE%' } ] );

	my $rs_extensions = $self->db->resultset('ToExtension')->search(undef, { prefetch => [ 'type' ] });
	my %mapping       = ();
	while ( my $ext = $rs_extensions->next ) {
		next unless ( $ext->type->name eq "CHECK_EXTENSION_BOOL" || $ext->type->name eq "CHECK_EXTENSION_NUM" );
		$mapping{ $ext->servercheck_column_name } = $ext->servercheck_short_name;
	}
	my $rs = $self->db->resultset('Server')->search( \%condition, { prefetch => [ 'servercheck', 'status', 'profile', 'cachegroup', 'type' ], order_by => 'me.host_name ASC' } );
	my @data;
	while ( my $server = $rs->next ) {
		my $v;
		$v->{id}		= $server->id + 0;
		$v->{hostName}		= $server->host_name;
		$v->{profile}		= $server->profile->name;
		$v->{adminState}	= $server->status->name;
		$v->{cacheGroup}	= $server->cachegroup->name;
		$v->{type}		= $server->type->name;
		$v->{updPending}	= \$server->upd_pending;
		$v->{revalPending}	= \$server->reval_pending;
		foreach my $col (qw/aa ab ac ad ae af ag ah ai aj ak al am an ao ap aq ar at au av aw ax ay az ba bb bc bd be bf/) {
			if ( defined( $mapping{$col} ) && defined( $server->servercheck ) ) {
				my $server_check_val = $server->servercheck->$col();
				if ( defined($server_check_val) ) {
					$server_check_val = $server_check_val + 0; # coerce to a number, if it's not null
				}
				$v->{checks}->{ $mapping{$col} } = $server_check_val;
			}
		}
		push( @data, $v );
	}
	return $self->success( \@data );
}

# update - we _should_ never have to do a create, so if there's an update for a servercheck line that doesn't exist, just create it.
sub update {
	my $self = shift;

	my $server_host_name       = $self->req->json->{host_name};
	my $server_id              = $self->req->json->{id};
	my $servercheck_short_name = $self->req->json->{servercheck_short_name};
	my $value                  = $self->req->json->{value};

	if ( $self->current_user()->{username} ne "extension" ) {
		return $self->alert( { error => "Invalid user for this API. Only the \"extension\" user can use this." } );
	}

	if ( !defined($server_id) || $server_id eq "" ) {
		$server_id = $self->db->resultset('Server')->search( { host_name => $server_host_name } )->get_column('id')->single();
	}

	if ( !defined($server_id) || $server_id eq "" ) {
		return $self->alert( { error => "Server not found" } );
	}

	my $column_name =
		$self->db->resultset('ToExtension')->search( { servercheck_short_name => $servercheck_short_name } )->get_column('servercheck_column_name')
		->single();
	if ( !defined($column_name) || $column_name eq "" ) {
		return $self->alert( { error => "Server Check Extension " . $servercheck_short_name . " not found - Do you need to install it?" } );
	}

	my $update = $self->db->resultset('Servercheck')->search( { server => $server_id } );
	$update->update_or_create( { $column_name => $value, } );

	return $self->success_message("Server Check was successfully updated.");
}
1;
