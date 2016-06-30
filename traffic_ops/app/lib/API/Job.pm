package API::Job;

#
# Copyright 2015 Comcast Cable Communications Management, LLC
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;

use Mojo::Base 'Mojolicious::Controller';
use UI::Utils;
use Digest::SHA1 qw(sha1_hex);
use Mojolicious::Validator;
use Mojolicious::Validator::Validation;
use Mojo::JSON;
use Time::Local;
use LWP;
use Email::Valid;
use MojoPlugins::Response;
use MojoPlugins::Job;
use Utils::Helper::ResponseHelper;
use Validate::Tiny ':all';
use UI::ConfigFiles;
use UI::Tools;
use Data::Dumper;

sub index {
	my $self = shift;

	my $response = [];
	my $username = $self->current_user()->{username};
	my $ds_id    = $self->param('dsId');
	my $keyword  = $self->param('keyword') || 'PURGE';

	my $dbh;
	if ( defined($ds_id) ) {
		$dbh = $self->db->resultset('Job')->search(
			{ keyword  => $keyword, 'job_user.username'     => $username, 'job_deliveryservice.id' => $ds_id },
			{ prefetch => [         { 'job_deliveryservice' => undef } ], join                     => 'job_user' }
		);
		my $row_count = $dbh->count();
		if ( defined($dbh) && ( $row_count > 0 ) ) {
			my @data = $self->job_ds_data($dbh);
			my $rh   = new Utils::Helper::ResponseHelper();
			$response = $rh->camelcase_response_keys(@data);
		}
	}
	else {
		$dbh =
			$self->db->resultset('Job')
			->search( { keyword => $keyword, 'job_user.username' => $username }, { prefetch => [ { 'job_user' => undef } ], join => 'job_user' } );
		my $row_count = $dbh->count();
		if ( defined($dbh) && ( $row_count > 0 ) ) {
			my @data = $self->job_data($dbh);
			my $rh   = new Utils::Helper::ResponseHelper();
			$response = $rh->camelcase_response_keys(@data);
		}
	}

	return $self->success($response);
}

# Creates a purge job based upon the Deliveryservice (ds_id) instead
# of the ds_xml_id like the UI does.
sub create {
	my $self = shift;

	my $ds_id      = $self->req->json->{dsId};
	my $agent      = $self->req->json->{agent};
	my $keyword    = $self->req->json->{keyword};
	my $regex      = $self->req->json->{regex};
	my $ttl        = $self->req->json->{ttl};
	my $start_time = $self->req->json->{startTime};
	my $asset_type = $self->req->json->{assetType};

	if ( !&is_admin($self) && !&is_oper($self) ) {

		# not admin or operations -- only an assigned user can purge
		my $tm_user = $self->db->resultset('TmUser')->search( { username => $self->current_user()->{username} } )->single();
		my $tm_user_id = $tm_user->id;

		if ( defined($ds_id) ) {

			# select deliveryservice from deliveryservice_tmuser where deliveryservice=$ds_id
			my $dbh = $self->db->resultset('DeliveryserviceTmuser')->search( { deliveryservice => $ds_id, tm_user_id => $tm_user_id }, { id => 1 } );
			my $count = $dbh->count();

			if ( $count == 0 ) {
			    $self->forbidden("Forbidden. Delivery service not assigned to user.");
				return;
			}
		}
	}

	# Just pass "true" in the urgent key to make it urgent.
	my $urgent = $self->req->json->{urgent};

	my ( $is_valid, $result ) = $self->is_valid( { dsId => $ds_id, regex => $regex, startTime => $start_time, ttl => $ttl } );
	if ($is_valid) {
		my $new_id = $self->create_new_job( $ds_id, $regex, $start_time, $ttl, 'PURGE', $urgent );
		if ($new_id) {
			my $saved_job = $self->db->resultset("Job")->find( { id => $new_id } );
			my $asset_url = $saved_job->asset_url;
			&log( $self, "Invalidate content request submitted for " . $asset_url, "APICHANGE" );
			return $self->success_message( "Invalidate content request submitted for: " . $asset_url . " (" . $saved_job->parameters . ")" );
		}
		else {
			return $self->alert( { "Error creating invalidate content request" . $ds_id } );
		}
	}
	else {
		return $self->alert($result);
	}
}

sub is_valid {
	my $self = shift;
	my $job  = shift;

	my $rules = {
		fields => [qw/dsId regex startTime ttl/],

		# Checks to perform on all fields
		checks => [

			# All of these are required
			[qw/regex startTime ttl dsId/] => is_required("is required"),

			ttl => sub {
				my $value  = shift;
				my $params = shift;
				if ( defined( $params->{'ttl'} ) ) {
					return $self->is_ttl_in_range($value);
				}
			},

			startTime => sub {
				my $value  = shift;
				my $params = shift;
				if ( defined( $params->{'ttl'} ) ) {
					return $self->is_valid_date_format($value);
				}
			},
			startTime => sub {
				my $value  = shift;
				my $params = shift;
				if ( defined( $params->{'ttl'} ) ) {
					return $self->is_more_than_two_days($value);
				}
			},

		]
	};

	# Validate the input against the rules
	my $result = validate( $job, $rules );

	if ( $result->{success} ) {

		#print "success: " . dump( $result->{data} );
		return ( 1, $result->{data} );
	}
	else {

		#print "failed " . Dumper( $result->{error} );
		return ( 0, $result->{error} );
	}

}

sub is_valid_date_format {
	my $self  = shift;
	my $value = shift;
	if ( !defined $value or $value eq '' ) {
		return undef;
	}

	if (
		( $value ne '' )
		&& ( $value !~
			qr/^((((19|[2-9]\d)\d{2})[\/\.-](0[13578]|1[02])[\/\.-](0[1-9]|[12]\d|3[01])\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((19|[2-9]\d)\d{2})[\/\.-](0[13456789]|1[012])[\/\  .-](0[1-9]|[12]\d|30)\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((19|[2-9]\d)\d{2})[\/\.-](02)[\/\.-](0[1-9]|1\d|2[0-8])\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]))|(((1[  6-9]|[2-9]\d)(0[48]|[2468][048]|[13579][26])|((16|[2468][048]|[3579][26])00))[\/\.-](02)[\/\.-](29)\s(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])))$/
		)
		)
	{
		return "has an invalidate date format, should be in the form of YYYY-MM-DD HH:MM:SS";
	}

	return undef;
}

sub is_ttl_in_range {
	my $self  = shift;
	my $value = shift;
	my $min_hours = 1;
	my $max_days =
		$self->db->resultset('Parameter')->search( { name => "maxRevalDurationDays" }, { config_file => "regex_revalidate.config" } )->get_column('value')->first;
	my $max_hours = $max_days * 24;

	if ( !defined $value or $value eq '' ) {
		return undef;
	}

	if ( ( $value ne '' ) && ( $value < $min_hours || $value > $max_hours ) ) {
		return "should be between " . $min_hours . " and " . $max_hours;
	}

	return undef;
}

sub is_more_than_two_days {
	my $self  = shift;
	my $value = shift;
	if ( !defined $value or $value eq '' ) {
		return undef;
	}

	my $dh               = new Utils::Helper::DateHelper();
	my $start_time_epoch = $dh->date_to_epoch($value);
	my $date_range       = abs( $start_time_epoch - time() );
	if ( ( $value ne '' ) && ( $date_range > 172800 ) ) {
		return "needs to be within two days from now.";
	}

	return undef;
}

1;
