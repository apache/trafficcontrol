package Utils::DeliveryService;
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

use Mojo::UserAgent;
use utf8;
use Carp qw(cluck confess);
use Data::Dumper;

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	return ( bless( $self, $class ) );
}

sub deliveryservice_lookup_cdn_name_and_ds_name {
	my $self = shift;
	my $dsid = shift || confess("Delivery Service id is required");

	my $cdn_name = "all";
	my $ds_name  = "all";
	if ( $dsid ne "all" ) {
		my $ds = $self->db->resultset('Deliveryservice')->search( { id => $dsid }, {} )->single();
		$ds_name = $ds->xml_id;
		my $param =
			$self->db->resultset('ProfileParameter')
			->search( { -and => [ profile => $ds->profile->id, 'parameter.name' => 'CDN_name' ] }, { prefetch => [ 'parameter', 'profile' ] } )->single();
		$cdn_name = $param->parameter->value;
	}
	return ( $cdn_name, $ds_name );
}

sub build_match {
	my $self            = shift;
	my $cdn_name        = shift;
	my $ds_name         = shift;
	my $cachegroup_name = shift;
	my $peak_usage_type = shift;
	return $cdn_name . ":" . $ds_name . ":" . $cachegroup_name . ":all:" . $peak_usage_type;
}
1;
