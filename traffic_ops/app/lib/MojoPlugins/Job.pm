package MojoPlugins::Job;

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
use Mojo::Base 'Mojolicious::Plugin';
use Data::Dumper;
use Carp qw(cluck confess);
use Data::Dumper;
use POSIX qw(strftime);
use UI::Utils;
use File::Path qw(make_path);

use constant PENDING      => 1;
use constant PROGRESS     => 2;
use constant COMPLETED    => 3;
use constant CANCELLED    => 4;
use constant REGEX_CONFIG => 'regex_revalidate.config';

sub register {
	my ( $self, $app, $conf ) = @_;

	$app->renderer->add_helper(

		# set the update bit for all the Caches in the CDN of this delivery service.
		set_update_server_bits => sub {
			my $self  = shift;
			my $ds_id = shift;

			my $cdn_id = $self->db->resultset('Deliveryservice')->search( { 'me.id' => $ds_id } )->get_column('cdn_id')->single();

			my @offstates;
			my $offline = $self->db->resultset('Status')->search( { 'name' => 'OFFLINE' } )->get_column('id')->single();
			if ($offline) {
				push( @offstates, $offline );
			}
			my $pre_prod = $self->db->resultset('Status')->search( { 'name' => 'PRE_PROD' } )->get_column('id')->single();
			if ($pre_prod) {
				push( @offstates, $pre_prod );
			}

			# Only queue updates for servers that have profiles with the regex_revalidate.config file location parameter.
			# If that parameter is not there, other mechanisms (ansible? script?) must be used to copy the
			# regex_revalidate.config to the caches.
			my @profiles = $self->db->resultset('ProfileParameter')->search(
				{
					-and => [
						'parameter.name'        => 'location',
						'parameter.config_file' => 'regex_revalidate.config'
					]
				},
				{ prefetch => [qw{ parameter profile }] }
			)->get_column('profile')->all();

			my $update_server_bit_rs = $self->db->resultset('Server')->search(
				{
					'me.cdn_id' => $cdn_id,
					-and        => { status => { 'not in' => \@offstates }, profile => { 'in' => \@profiles } }
				}
			);

			my $use_reval_pending = $self->db->resultset('Parameter')->search( { -and => [ 'name' => 'use_reval_pending', 'config_file' => 'global' ] } )->get_column('value')->single;

			if ( defined($use_reval_pending) && $use_reval_pending ne '0' ) {
				my $result = $update_server_bit_rs->update( { reval_pending => 1 } );
				&log( $self, "Set reval_pending = 1 for all applicable caches", "OPER" );
			}
			else {
				my $result = $update_server_bit_rs->update( { upd_pending => 1 } );
				&log( $self, "Set upd_pending = 1 for all applicable caches", "OPER" );
			}
		}
	);

	$app->renderer->add_helper(
		job_data => sub {
			my $self = shift;
			my $dbh  = shift;

			my @data;
			while ( my $row = $dbh->next ) {
				push(
					@data, {
						id           => $row->id,
						agent        => $row->agent->name,
						object_type  => $row->object_type,
						object_name  => $row->object_name,
						entered_time => $row->entered_time,
						keyword      => $row->keyword,
						parameters   => $row->parameters,
						asset_url    => $row->asset_url,
						asset_type   => $row->asset_type,
						status       => $row->status->name,
						username     => $row->job_user->username,
						start_time   => $row->start_time,
					}
				);
			}
			return \@data;
		}
	);

	$app->renderer->add_helper(
		job_ds_data => sub {
			my $self = shift;
			my $dbh  = shift;

			my @data;
			while ( my $row = $dbh->next ) {
				push(
					@data, {
						id           => $row->id,
						agent        => $row->agent->name,
						object_type  => $row->object_type,
						object_name  => $row->object_name,
						entered_time => $row->entered_time,
						keyword      => $row->keyword,
						parameters   => $row->parameters,
						asset_url    => $row->asset_url,
						asset_type   => $row->asset_type,
						status       => $row->status->name,
						username     => $row->job_user->username,
						start_time   => $row->start_time,
						ds_id        => $row->job_deliveryservice->id,
						ds_xml_id    => $row->job_deliveryservice->xml_id,
					}
				);
			}
			return \@data;
		}
	);

	$app->renderer->add_helper(
		create_new_job => sub {
			my $self       = shift;
			my $ds_id      = shift;
			my $regex      = shift;
			my $start_time = shift;
			my $ttl        = shift || '';
			my $keyword    = shift || 'PURGE';
			my $urgent     = shift;

			# Defaulted parameters
			my $parameters  = shift;
			my $asset_type  = shift || 'file';
			my $status      = shift || 1;
			my $object_type = shift;
			my $object_name = shift;

			if ( !defined($parameters) || $parameters eq "" ) {
				if ( defined($ttl) && $ttl =~ m/^\d/ ) {
					$parameters = "TTL:" . $ttl . 'h';
				}
			}

			## Calculate start time
			# Convert to unix time and give a default value if not specified
			if ( !defined($start_time) || $start_time eq "" ) {
				$start_time = time();
			}
			else {
				my $dh = new Utils::Helper::DateHelper();
				$start_time = $dh->date_to_epoch($start_time);
			}

			# add 60s if not urgent
			if ( !defined $urgent ) {
				$start_time = $start_time + 60;
			}
			my $start_time_gmt = strftime( "%Y-%m-%d %H:%M:%S", gmtime($start_time) );
			my $entered_time   = strftime( "%Y-%m-%d %H:%M:%S", gmtime() );
			my $org_server_fqdn = UI::DeliveryService::compute_org_server_fqdn($self, $ds_id);

			my $tm_user_id = $self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->get_column('id')->single();

			$regex =~ m/(^\/.+)/ ? $org_server_fqdn = $org_server_fqdn . "/$regex" : $org_server_fqdn = $org_server_fqdn . "$regex";
			my $insert = $self->db->resultset('Job')->create(
				{
					agent               => 1,
					object_type         => $object_type,
					object_name         => $object_name,
					entered_time        => $entered_time,
					keyword             => $keyword,
					parameters          => $parameters,
					asset_url           => $org_server_fqdn,
					asset_type          => $asset_type,
					status              => $status,
					job_user            => $tm_user_id,
					start_time          => $start_time_gmt,
					job_deliveryservice => $ds_id,
				}
			);

			my $new_record = $insert->insert();

			&log( $self, "Created new Purge Job " . $ds_id . " forced new " . REGEX_CONFIG . " snapshot", "APICHANGE" );

			$self->set_update_server_bits($ds_id);
			return $new_record->id;
		}
	);

	$app->renderer->add_helper(
		check_job_auth => sub {
			my $self       = shift;
			my $tm_user_id = shift;
			my $asset      = shift;
			if ( &is_admin($self) ) {
				return 1;
			}
			else {

				my $rs_ds_ids = $self->db->resultset('DeliveryserviceTmuser')->search( { tm_user_id => $tm_user_id } );
				my $rs_ds = $self->db->resultset('Deliveryservice')->search( { id => { -in => $rs_ds_ids->get_column('deliveryservice')->as_query } } );

				my ( $scheme, $asset_hostname, $path, $query, $fragment ) = $asset =~ m|(?:([^:/?#]+):)?(?://([^/?#]*))?([^?#]*)(?:\?([^#]*))?(?:#(.*))?|;

				while ( my $ds_row = $rs_ds->next ) {
					my $org_server_fqdn = UI::DeliveryService::compute_org_server_fqdn($self, $ds_row->id);
					if ( defined($org_server_fqdn) && $asset =~ /$org_server_fqdn/ ) {
						return 1;    # Success
					}
				}
				return 0;            # Fail
			}
		}
	);
}

1;
