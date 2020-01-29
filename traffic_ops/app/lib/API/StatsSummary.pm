package API::StatsSummary;
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

# JvD Note: you always want to put Utils as the first use. Sh*t don't work if it's after the Mojo lines.
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
my $builder;
use constant SUCCESS => 0;
use constant ERROR   => 1;

sub index {
	my $self      = shift;
	my $cdn_name  = $self->param('cdnName');
	my $ds_name   = $self->param('deliveryServiceName');
	my $stat_name = $self->param('statName');

	# TODO: Implement start_date and end_date
	# my $start_date  = $self->param('startDate');
	# my $end_date    = $self->param('endDate');
	my $last_summary_date = $self->param('lastSummaryDate');    ##Boolean.  Used by traffic stats to determine if summary_data needs to be written.
	my %q;

	if ($last_summary_date) {
		if ($stat_name) {
			$self->app->log->debug("statName -> $stat_name");
			%q = ( 'stat_name' => $stat_name );
		}
		my $response = $self->db->resultset('StatsSummary')->search( \%q )->get_column('summary_time')->max();
		return $self->success( { "summaryTime" => $response } );
	}
	##add name and delivery_service to search
	if ($cdn_name) {
		%q = ( 'cdn_name' => $cdn_name );
	}
	if ($ds_name) {
		%q = ( %q, 'deliveryservice_name' => $ds_name );
	}
	if ($stat_name) {
		%q = ( %q, 'stat_name' => $stat_name );
	}
	my $rs_data = $self->db->resultset('StatsSummary')->search( \%q );
	my @data;
	while ( my $row = $rs_data->next ) {
		push(
			@data, {
				"cdnName"        		=> $row->cdn_name,
				"deliveryServiceName"   => $row->deliveryservice_name,
				"statName"       		=> $row->stat_name,
				"statValue" 	 		=> $row->stat_value,
				"summaryTime" 	 		=> $row->summary_time,
				"statDate" 	 			=> $row->stat_date,
		}
		);
	}
	return $self->success( \@data );
}

sub create {
	my $self         = shift;
	my $cdn_name     = $self->req->json->{cdnName} || "all";
	my $ds_name      = $self->req->json->{deliveryServiceName} || "all";
	my $stat_name    = $self->req->json->{statName};
	my $stat_value   = $self->req->json->{statValue};
	my $summary_time = $self->req->json->{summaryTime};
	my $stat_date    = $self->req->json->{statDate};

	my $alternative = "POST /stats_summary";

	if ( !defined($stat_name) || !defined($stat_value) || !defined($stat_date) ) {
		return $self->with_deprecation("Please provide a stat name, value, and date", "error", 400, $alternative);
	}

	my $insert = $self->db->resultset('StatsSummary')->create(
		{
			cdn_name             => $cdn_name,
			deliveryservice_name => $ds_name,
			stat_name            => $stat_name,
			stat_value           => $stat_value,
			summary_time         => $summary_time,
			stat_date            => $stat_date,
		}
	);
	$insert->insert();
	return $self->with_deprecation("Successfully added stats summary record", "success", 200, $alternative);
}

1;
