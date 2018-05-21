package UI::Job;
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
use LWP;
use UI::ConfigFiles;
use UI::Tools;
use MojoPlugins::Job;
use Utils::Helper::ResponseHelper;
use Utils::Helper::DateHelper;

my $fileInfo = __FILE__ . ":";

# Table view
sub jobs {
	my $self = shift;

	&navbarpage($self);
}

sub newjob {
	my $self = shift;
	my %response;
	my $agent       = $self->param('job.agent') || 1;
	my $object_type = $self->param('job.object_type');
	my $object_name = $self->param('job.object_name');
	my $keyword     = $self->param('job.keyword');
	my $parameters  = $self->param('job.parameters');
	my $ds_xml_id   = $self->param('job.ds_xml_id');
	my $ttl         = $self->param('job.ttl');
	my $regex       = $self->param('job.regex');
	my $asset_type  = $self->param('job.asset_type');
	my $start_time  = $self->param('job.start_time');
	my $urgent      = $self->param('job.urgent');

	my $status = 1;
	my $entered_time = strftime( "%Y-%m-%d %H:%M:%S", gmtime() );
	my %err;
	my $job_user;
	my $user;

	my $message = "Invalidate content job submitted successfully!";

	if ( !defined($parameters) || $parameters eq "" ) {
		if ( defined($ttl) && $ttl =~ m/^\d/ ) {
			$parameters = "TTL:" . $ttl . 'h';
		}
	}

	if ( $self->is_valid_job() ) {
		my $ds = $self->db->resultset("Deliveryservice")->search( { xml_id => $ds_xml_id }, { prefetch => [ 'type', 'profile', 'cdn' ] } )->single();
		my $org_server_fqdn;
		if ( $ds->type->name eq 'ANY_MAP' ) {
			$org_server_fqdn = $ds->remap_text;
			$org_server_fqdn =~ s/^\S+\s+(\S+)\s+.*/$1/;    # get the thing after (re)map
			$org_server_fqdn =~ s/\/$//;
		}
		else {
			$org_server_fqdn = UI::DeliveryService::compute_org_server_fqdn($self, $ds->id);
		}
		my $ds_id = $ds->id;

		$user = $self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->get_column('id')->single();

		if ( !defined($user) ) {
			$err{"status"} = "failure";

			$err{"message"} = "Permission denied.";
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

		if ( $start_time !~ m/^\d{10}$/ ) {
			$err{"status"}  = "failure";
			$err{"message"} = "Invalid time or format";
		}

		# right now, only PURGE is allowed
		if ( !defined($keyword) || $keyword ne "PURGE" ) {
			$err{"status"}  = "failure";
			$err{"message"} = "keyword " . $keyword . " is not supported";
		}
		elsif ( $keyword eq "PURGE" ) {
			if ( !defined($org_server_fqdn) || !defined($asset_type) ) {
				$err{"status"}  = "failure";
				$err{"message"} = "Missing parameters, need at least org_server_fqdn and asset_type";
			}
			if ( $self->check_job_auth( $user, $org_server_fqdn ) == 0 ) {
				$err{"status"}  = "failure";
				$err{"message"} = "Insufficient permissions to act on this object";
			}
		}

		# ttl parameter must have a TTL:xxh where xx is the TTL of the regex_reval
		# 48 hours is the default ttl.
		if ( $keyword eq "PURGE" ) {
			if ( !defined($parameters) || $parameters eq "" ) {
				$parameters = "Using default TTL:48h";
			}
			elsif ( $parameters !~ /TTL:(\d+)h/ ) {
				$err{"status"}  = "failure";
				$err{"message"} = "Invalid TTL parameter - expecting TTL:xxh";
			}
		}

		my $start_time_gmt = strftime( "%Y-%m-%d %H:%M:%S", gmtime($start_time) );
		if ( !%err ) {
			if ( $regex =~ m/(^\/.+)/ ) {
				$org_server_fqdn = $org_server_fqdn . "$regex";
			}
			else {
				$org_server_fqdn = $org_server_fqdn . "/$regex";
			}
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
					job_user            => $user,
					start_time          => $start_time_gmt,
					job_deliveryservice => $ds->id,
				}
			);
			$insert->insert();
			$response{"job"} = $insert->id;

			&log( $self, "UI entry " . $response{job} . " forced new regex_revalidate.config snapshot", "UICHANGE" );

			# my $ds_id =
			# 	$self->db->resultset('Deliveryservice')->search( { xml_id => $ds_xml_id }, { prefetch => ['profile'] } )->get_column('id')->single();
			my $rs       = $self->db->resultset('Deliveryservice')->search( { 'me.xml_id' => $ds_xml_id }, { prefetch => 'cdn' } )->single;
			my $ds_id    = $rs->id;

			$self->set_update_server_bits($ds_id);

			if ( defined $response{"job"} ) {
				$response{"status"} = "success";
				$self->flash( message => $message );
				$self->render( json => \%response );
				return $self->redirect_to('/tools/invalidate_content');
			}
			else {
				$message = "Job not successfully submitted to DB";
				my $id = $self->param('id');
				$err{input_params} = $self->req->params->to_hash;
				my %delivery_service = UI::Tools::get_delivery_services( $self, $id );
				$self->stash(
					job => { regex => $regex, start_time => $start_time, ttl => $ttl },
					ds  => \%delivery_service
				);
				&navbarpage($self);
				$self->flash( message => $message );
				return $self->redirect_to('/tools/invalidate_content');
			}
		}
		else {
			my $id               = $self->param('id');
			my %delivery_service = UI::Tools::get_delivery_services( $self, $id );
			my $selected_ds      = $ds_xml_id;
			$err{input_params} = $self->req->params->to_hash;
			$self->stash( job => { regex => $regex, start_time => $start_time, ttl => $ttl }, ds => \%delivery_service, selected_ds => $selected_ds );
			&navbarpage($self);
			$self->flash( message => $err{'message'} );
			return $self->redirect_to('/tools/invalidate_content');
		}
	}
	else {
		my $id               = $self->param('id');
		my %delivery_service = UI::Tools::get_delivery_services( $self, $id );
		my $selected_ds      = $ds_xml_id;
		$self->stash( job => { regex => $regex, start_time => $start_time, ttl => $ttl }, ds => \%delivery_service, selected_ds => $selected_ds );
		&navbarpage($self);
		$self->render('tools/invalidate_content');
	}
}

sub is_valid_job {
	my $self = shift;

	# check if ttl is between min and max
	my $min_hours = 1;
	my $max_days =
		$self->db->resultset('Parameter')->search( { name => "maxRevalDurationDays" }, { config_file => "regex_revalidate.config" } )->get_column('value')->first;
	my $max_hours = $max_days * 24;	
	my $ttl = $self->param('job.ttl');
	if ( $ttl eq '' || $ttl < $min_hours || $ttl > $max_hours ) {
		$self->field('job.ttl')->is_like( qr/^\//,, "Must be between " . $min_hours . " and " . $max_hours )
			;    # hack: this will fail and trigger error message
	}

	$self->field('job.start_time')->is_required->is_like(
		qr/^((((19|[2-9]\d)\d{2})[\/\.-](0[13578]|1[02])[\/\.-](0[1-9]|[12]\d|3[01])\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((19|[2-9]\d)\d{2})[\/\.-](0[13456789]|1[012])[\/\.-](0[1-9]|[12]\d|30)\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((19|[2-9]\d)\d{2})[\/\.-](02)[\/\.-](0[1-9]|1\d|2[0-8])\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((1[6-9]|[2-9]\d)(0[48]|[2468][048]|[13579][26])|((16|[2468][048]|[3579][26])00))[\/\.-](02)[\/\.-](29)\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])))$/,
		"Not a valid dateTime!  Should be in the format of YYYY-MM-DD HH:MM:SS"
	);
	$self->field('job.regex')->is_required;
	my $ds_xml_id = $self->param('job.ds_xml_id');
	if ( $ds_xml_id eq 'default' ) {
		$self->field('job.ds_xml_id')->is_like( qr/^\/(?!default\/)/,, "Please choose a Delivery Service!" );
	}

	# check for start_date too far in the future (>2 days?)
	my $now        = time();
	my $start_time = $self->param('job.start_time');
	my $dh         = new Utils::Helper::DateHelper();
	$start_time = $dh->date_to_epoch($start_time);

	# check for start_date too far in the future (>2 days?)
	if ( abs( $start_time - $now ) > 172800 ) {
		$self->field('job.start_time')->is_like( qr/^\//,, "The start time is too far in the future, > 2 days" );
	}

	return $self->valid;
}

sub read_job_by_id {

	my $self = shift;

	my @data = $self->read_job_by_column_name_and_value( 'id', $self->param('id') );

	my $rh = new Utils::Helper::ResponseHelper();

	$self->render( json => \@data );
}

sub canceljob {

	my $self = shift;
	my $id   = $self->param('id');
	my $user = $self->param('user');
	my %response;

	# TODO validate user is authorized to add jobs
	my $update = $self->db->resultset('Job')->find( { id => $id } );
	$update->status(4);
	$update->update();

	# TODO do actual validation if update ran
	$response{"status"} = "success";

	$self->render( json => \%response );
}

sub jobstatusupdate {

	#        $r->get('/job/agent/statusupdate/:id')->to('job#jobstatusupdate');
	my $self = shift;
	my $id   = $self->param('id');
	my %response;

	my $update = $self->db->resultset('Job')->find( { id => $id } );
	$update->status(2);
	$update->update();

	# TODO do actual validation if update ran
	$response{"status"} = "success";

	$self->render( json => \%response );
}

sub listjob {

	#$r->get('/job/view/all')->to('job#listjob');
	my $self = shift;
	my @data;
	my $orderby = "id";

	my $rs_data = $self->db->resultset('Job')->search( undef, { order_by => 'me.' . $orderby, } );
	while ( my $row = $rs_data->next ) {
		my %hash = (
			id           => $row->id,
			agent        => $row->agent->name,
			object_type  => $row->object_type,
			object_name  => $row->object_name,
			keyword      => $row->keyword,
			parameters   => $row->parameters,
			asset_url    => $row->asset_url,
			asset_type   => $row->asset_type,
			status       => $row->status->name,
			start_time   => $row->start_time,
			entered_time => $row->entered_time,
			last_updated => $row->last_updated,
		);
		push( @data, \%hash );
	}

	$self->render( json => \@data );
}

#### Status subs
sub readstatus {

	#        $r->get('/job/external/status/view/all')->to('job#readstatus');
	my $self = shift;
	my @data;
	my $rs_data = $self->db->resultset('JobStatus')->search();
	while ( my $row = $rs_data->next ) {
		my %hash = (
			id           => $row->id,
			name         => $row->name,
			description  => $row->description,
			last_updated => $row->last_updated,
		);
		push( @data, \%hash );
	}

	if ( !@data ) {
		push( @data, { result => "No data found" } );
	}

	$self->render( json => \@data );
}

#### Agent subs
sub viewagentjob {

	#        $r->get('/job/agent/list/:id')->to('job#listagentjob');
	# See if there's a pending job for this agent
	my $self    = shift;
	my $id      = $self->param('id');
	my $orderby = "id";
	my @data;

	my $rs_data;
	if ( $id eq "all" ) {
		$rs_data = $self->db->resultset('Job')->search( { status => 1 }, { prefetch => [ 'agent', 'status', 'job_user' ], order_by => 'me.' . $orderby, } );
	}
	else {
		$rs_data = $self->db->resultset('Job')->search( { agent => $id, status => 1 }, { prefetch => [ 'agent', 'status', 'job_user' ], order_by => 'me.' . $orderby, } );
	}

	while ( my $row = $rs_data->next ) {

		my %hash = (
			id           => $row->id,
			agent_name   => $row->agent->name,
			object_type  => $row->object_type,
			object_name  => $row->object_name,
			keyword      => $row->keyword,
			parameters   => $row->parameters,
			asset_url    => $row->asset_url,
			asset_type   => $row->asset_type,
			status       => $row->status->name,
			start_time   => $row->start_time,
			entered_time => $row->entered_time,
			username     => $row->job_user->username,
			company      => $row->job_user->company,
			last_updated => $row->last_updated,
		);

		push( @data, \%hash );
	}

	if ( !@data ) {
		push( @data, { result => "No pending jobs for agent ID $id" } );
	}

	$self->render( json => \@data );
}

sub addagent {
	my $self = shift;

	#        $r->get('/job/agent/new')->to('job#addagent');
}

sub newagent {

	#        $r->post('/job/agent/new')->to('job#newagent');
	my $self        = shift;
	my $description = $self->param('description');
	my $name        = $self->param('name');
	my $active      = $self->param('active');

	my %err;
	my %response;

	# TODO Check for required fields
	if ( $name =~ m/^\s*$/ )        { $err{"Error"} = "Required field missing" }
	if ( $description =~ m/^\s*$/ ) { $err{"Error"} = "Required field missing" }
	if ( $active =~ m/^\s*$/ )      { $active       = 0 }

	if (%err) {

		# Set error headers and gen message
		return $self->render( json => \%err );
	}
	else {
		my $insert = $self->db->resultset('JobAgent')->create(
			{
				name        => $name,
				description => $description,
				active      => $active,
			}
		);
		$insert->insert();
		$response{"result"} = $insert->id;
	}

	if ( defined $response{"result"} ) {
		$response{"status"} = "success";
		&log( $self, "Created new job agent " . $name, "UICHANGE" );
	}
	else {
		$response{"status"} = "failure";
	}

	$self->render( json => \%response );

}

sub readagent {

	#        $r->get('/job/agent/view/all')->to('job#readagent');
	my $self = shift;
	my @data;
	my $rs_data = $self->db->resultset('JobAgent')->search();

	while ( my $row = $rs_data->next ) {
		my %hash = (
			id           => $row->id,
			name         => $row->name,
			description  => $row->description,
			active       => $row->active,
			last_updated => $row->last_updated,
		);
		push( @data, \%hash );
	}

	if ( !@data ) {
		push( @data, { result => "No data found" } );
	}

	$self->render( json => \@data );
}

1;
