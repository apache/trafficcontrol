package API::ToExtension;
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
use POSIX qw(strftime);
use Time::Local;
use Utils::Helper::ResponseHelper;

sub index {
	my $self = shift;
	my @data;
	my $rs = $self->db->resultset('ToExtension')->search( undef, { prefetch => ['type'] } );
	while ( my $row = $rs->next ) {
		next unless ( $row->type->name ne 'CHECK_EXTENSION_OPEN_SLOT' );    # Open slots are not in the list
		push(
			@data, {
				id                     => $row->id + 0,
				name                   => $row->name,
				version                => $row->version,
				info_url               => $row->info_url,
				script_file            => $row->script_file,
				isactive               => $row->isactive,
				additional_config_json => $row->additional_config_json,
				description            => $row->description,
				servercheck_short_name => $row->servercheck_short_name,

				# servercheck_column_name => $row->servercheck_column_name, # Hide col name from the extension developer
				type => $row->type->name
			}
		);
	}

	# Config extensions are driven by the parameter/profile setup, much like the normal config files, using the name => 'location'
	# parameter, and we use the id of the parameter as the id of the extension
	$rs = $self->db->resultset('Parameter')->search( { name => 'location', config_file => { -like => 'to_ext_%.config' } } );
	while ( my $row = $rs->next ) {
		my $file;
		$file = $row->config_file;
		my $subroutine =
			$self->db->resultset('ProfileParameter')
			->search( { -and => [ 'parameter.config_file' => $file, 'parameter.name' => 'SubRoutine' ] },
			{ prefetch => [ 'parameter', 'profile' ] } )->get_column('parameter.value')->single();
		$subroutine =~ s/::[^:]+$/::info/;
		$self->app->log->error( "ToExtDotInfo == " . $subroutine );
		my $info = &{ \&{$subroutine} }();

		push(
			@data, {
				id                     => $row->id + 0,
				name                   => $info->{name},
				version                => $info->{version},
				info_url               => $info->{info_url},
				script_file            => $info->{script_file},
				isactive               => "n/a",
				additional_config_json => "n/a",
				description            => $info->{description},
				servercheck_short_name => "n/a",
				type                   => "CONFIG_EXTENSION",
			}
		);
	}
#
#	$rs = $self->db->resultset('Parameter')->search( { name => 'datasource', config_file => 'global' } );
#	while ( my $row = $rs->next ) {
#		my $source;
#		$source = $row->value;
#		my $ext_hash_ref = &Extensions::DatasourceList::hash_ref();
#		my $subroutine   = $ext_hash_ref->{$source};
#		if ( !defined($subroutine) ) {
#			$self->app->log->error( "No subroutine found for: " . $source );
#		}
#		my $isub;
#		( $isub = $subroutine ) =~ s/::[^:]+$/::info/;
#		my $info = &{ \&{$isub} }();
#		print Dumper($info);
#		push(
#			@data, {
#				id                     => $row->id,
#				name                   => $info->{name},
#				version                => $info->{version},
#				info_url               => $info->{info_url},
#				script_file            => $info->{script_file},
#				isactive               => "n/a",
#				additional_config_json => "n/a",
#				description            => $info->{description},
#				servercheck_short_name => "n/a",
#				type                   => "DATASOURCE_EXTENSION",
#			}
#		);
#	}
	$self->success( \@data );
}

# update creates if there is no id in the json.
sub update {
	my $self = shift;
	my $msg  = "Error";

	my $new_id = 1;
	my $jdata  = $self->req->json;

	if ( $self->current_user()->{username} ne "extension" ) {
		return $self->alert( { error => "Invalid user for this API. Only the \"extension\" user can use this." } );
	}

	if ( defined( $jdata->{id} ) ) {
		return $self->alert( { error => "ToExtension update not supported; delete and re-add." } );
	}
	else {
		# we are creating.
		# print Dumper($jdata);
		my $type_id = &type_id( $self, $jdata->{type} );
		if (   !defined($type_id)
			|| !( $jdata->{type} =~ /^CHECK_EXTENSION_/ || $jdata->{type} =~ /^CONFIG_EXTENSION$/ || $jdata->{type} =~ /^STATISTIC_EXTENSION$/ ) )
		{
			return $self->alert( { error => "Invalid Extension type: " . $jdata->{type} } );
		}

		if ( $jdata->{type} =~ /CHECK_EXTENSION_/ ) {
			foreach my $f (qw/name servercheck_short_name/) {
				my $exists = $self->db->resultset('ToExtension')->search( { $f => $jdata->{$f} } )->single();
				if ( defined($exists) ) {
					return $self->alert( { error => "A Check extension is already loaded with " . $f . " = " . $jdata->{$f} } );
				}
			}

			# check extensions go in an open slot in the extensions table, first check if there's an open slot.
			my $open_type = &type_id( $self, 'CHECK_EXTENSION_OPEN_SLOT' );
			my $slot = $self->db->resultset('ToExtension')->search( { type => $open_type }, { rows => 1, order_by => ["servercheck_column_name"] } )->single();
			if ( !defined($slot) ) {
				return $self->alert( { error => "No open slots left for checks, delete one first." } );
			}
			$slot->update(
				{
					name                   => $jdata->{name},
					version                => $jdata->{version},
					info_url               => $jdata->{info_url},
					script_file            => $jdata->{script_file},
					isactive               => $jdata->{isactive},
					additional_config_json => $jdata->{additional_config_json},
					description            => $jdata->{description},
					servercheck_short_name => $jdata->{servercheck_short_name},
					type                   => $type_id
				}
			);

			# set all values in servercheck to 0
			my $clear = $self->db->resultset('Servercheck')->search( {} );    # all
			$clear->update( { $slot->servercheck_column_name => 0 } );    #

			return $self->success_message( "Check Extension Loaded.", { id => $slot->id } );
		}

		# Should not get here for CHECK_EXTENSION_* type, already returned the success above.
		my $insert = $self->db->resultset('ToExtension')->create(
			name                    => $jdata->{name},
			version                 => $jdata->{version},
			info_url                => $jdata->{info_url},
			script_file             => $jdata->{script_file},
			isactive                => $jdata->{isactive},
			additional_config_json  => $jdata->{additional_config_json},
			description             => $jdata->{description},
			servercheck_short_name  => $jdata->{servercheck_short_name},
			servercheck_column_name => $jdata->{servercheck_column_name},
			type                    => $type_id
		);
		$insert->insert();
		$new_id = $insert->id;
		if ( !defined($new_id) ) {
			return $self->alert( { error => "Unknown database error when inserting Extension." } );
		}
	}

	return $self->success_message( "Extension loaded.", { id => $new_id } );
}

sub delete {
	my $self = shift;
	my $msg  = "Error";

	my $new_id = 1;
	my $id     = $self->param('id');
	my $alt = "DELETE /servercheck/extensions/:id";

	# print Dumper($self->req);
	if ( $self->current_user()->{username} ne "extension" ) {
		return $self->with_deprecation("Invalid user for this API. Only the \"extension\" user can use this.", "error", 400, $alt);
	}

	if ( !defined($id) ) {
		return $self->with_deprecation("ToExtension delete requires an id.", "error", 400, $alt);
	}
	else {
		my $delete = $self->db->resultset('ToExtension')->search( { id => $id } )->single();
		if ( !defined($delete) ) {
			return $self->with_deprecation("ToExtension with id " . $id . " not found.", "error", 400, $alt);
		}
		if ( $delete->type->name =~ /^CHECK_EXTENSION_/ ) {
			my $open_type_id = &type_id( $self, 'CHECK_EXTENSION_OPEN_SLOT' );
			$delete->update(
				{
					name                   => 'OPEN',
					version                => '0',
					info_url               => '',
					script_file            => '',
					isactive               => '0',
					additional_config_json => '',
					servercheck_short_name => '',
					type                   => $open_type_id,
				}
			);
		}
		else {
			$delete->delete();
		}
	}
	return $self->with_deprecation("Extension deleted.", "success", 200, $alt);
}

1;
