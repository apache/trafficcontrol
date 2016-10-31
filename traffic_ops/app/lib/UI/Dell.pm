package UI::Dell;
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
use Mojo::Base 'Mojolicious::Controller';

sub dells {
	my $self = shift;
	# ASSumes the Dell profiles are 1 and 3
	my $rs_dells = $self->db->resultset("Server")->search( {
		-and => [
			ilo_ip_address => {'!=', undef }
			],
		-or => [
			profile => 1,
			profile => 3
			],
		}, {order_by => "host_name", columns => [qw/id host_name domain_name/]});
	my @data;
	while (my $row = $rs_dells->next) {
		push(@data, {
				"id" => $row->id,
				"fqdn" => $row->host_name . "." . $row->domain_name
		});
	}

	$self->render( json => \@data );
}

sub configuredrac {
	my $self = shift;
	my $serverid = $self->param('serverid');
	my $alertmsg;

	my $iloip = $self->db->resultset("Server")->search({id => $serverid})->get_column("ilo_ip_address")->single;
	if ($iloip) {
		$alertmsg = "Job successfully submitted";
		# Kick off job here $iloip
	} else {
		$alertmsg = "No ILO IP address found, aborting";
	}
	$self->flash(alertmsg => $alertmsg);
	return $self->redirect_to('/#tabs=5');
}

1;
