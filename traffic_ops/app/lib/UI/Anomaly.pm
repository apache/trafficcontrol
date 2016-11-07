package UI::Anomaly;
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
use UI::Server;
use JSON;
use WWW::Curl::Easy;
use Data::Dumper;
use Mojo::UserAgent;
use Digest::MD5 qw(md5_hex);


my $debug = 0;

sub start {
	my $self = shift;
	my $host_name = $self->param('host_name');
	my $rs_server = $self->db->resultset('Server')->search( { 'host_name' => $host_name }, { prefetch => 'type' } );
	my $row = $rs_server->next;
	my $type_name = $row->type->name;
	my $json;
	if ( $type_name eq 'RASCAL' ) {
		$json = &rascal($self, $host_name);
	}
	else {
		$json = { 'Message' => 'Unsupported type' };
	}
	$self->render( json => $json );
}

sub rascal {
	my $self = shift;
	my $host_name = shift;
	my $json;
	my $tm_crconfig_url;

	my $rs_server = $self->db->resultset('Server')->search( { 'host_name' => $host_name }, { prefetch => 'profile' } )->single();
	my $domain_name = $rs_server->domain_name;
	my $profile_id = $rs_server->profile->id;
	my $cdn_name = $rs_server->cdn->name;

	my %condition = ( 'me.name' => 'tm.crConfig.polling.url', 'profile_parameters.profile' => $profile_id );
	my $rs_param = $self->db->resultset('Parameter')->search( \%condition, { prefetch => 'profile_parameters' } );
	while (my $row = $rs_param->next) {
		if ($row->name eq 'tm.crConfig.polling.url') {
			$tm_crconfig_url = $row->value;
		}
	}
	$tm_crconfig_url =~ s/\$\{cdnName\}/$json->{'CDN_name'}/;
	$tm_crconfig_url =~ s/http(s)*\:\/\///;
	$tm_crconfig_url =~ s/\$\{tmHostname\}//;
	my $path = $json->{'CDN_name'} . "/CRConfig.xml";
	my $ua = Mojo::UserAgent->new;
	my $tm_crconfig = $self->ua->get($tm_crconfig_url)->res->content->asset->{'content'};
	my $rascal_crconfig = $self->ua->get("http://" . $host_name . "." . $domain_name . "/publish/CrConfig")->res->content->asset->{'content'};
	my $rascal_md5 = md5_hex($rascal_crconfig);
	my $tm_md5 = md5_hex($tm_crconfig);
	if ($rascal_md5 ne $tm_md5) {
		$json->{'CRConfig'}->{'State'} = 'Error'; 
		$json->{'CRConfig'}->{'Message'} = 'md5 from 12M: ' . $tm_md5 . ', md5 from ' . $host_name . ': ' . $rascal_md5;
	}
	return $json;
}

