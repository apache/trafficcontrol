package API::Cache;
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
use Utils::Rascal;

sub get_cache_stats {
	my $self = shift;

	my %response;
	my $master_i = 1;

	my %rascal_host = ();

	my @cdns = $self->db->resultset("Cdn")->search()->get_column('name')->all();

	foreach my $cdn_name (@cdns) {
		$rascal_host{$cdn_name} = $self->get_traffic_monitor_connection( { cdn => $cdn_name } );
	}

	my $big_obj = get_dataserver($self);

	foreach my $cdn ( keys %rascal_host ) {
		$self->get_crstates( $rascal_host{$cdn}, $big_obj );
		$self->add_cache_stats( $rascal_host{$cdn}, $big_obj, \%response, \$master_i );
	}
	$self->add_aggregate_cache_stats( $big_obj, \%response, \$master_i );

	$self->render( json => \%response );
}

sub get_crstates {
	my $self    = shift;
	my $rascal  = shift;
	my $big_obj = shift;

	if ( !defined($rascal) ) {
		$self->app->log->error('No running Traffic Monitor found!');
		return;
	}
	my $cs_hashref = $rascal->get_cr_states();

	if ( defined($cs_hashref) ) {
		foreach my $server ( sort keys %{ $cs_hashref->{'caches'} } ) {
			my $state = $cs_hashref->{'caches'}->{$server}->{'isAvailable'};
			$big_obj->{'caches'}->{$server}->{'healthy'} = defined($state) ? $state : 'false';
			$self->app->log->trace("Setting $server to $state");
		}
	}
}

sub get_dataserver {
	my $self    = shift;
	my $big_obj = {};
	my $data    = UI::Server::getserverdata($self);

	for my $i ( @{$data} ) {
		my $profile = $i->{'profile'};
		my $cdn     = undef;
		if ( defined($profile) && ( $profile =~ m/^EDGE/ || $profile =~ m/^MID/ ) ) {

			$self->app->log->trace( "get_data_server cache: " . $i->{'host_name'} );
			$big_obj->{'caches'}->{ $i->{'host_name'} }->{'cachegroup'}   = $i->{'cachegroup'};
			$big_obj->{'caches'}->{ $i->{'host_name'} }->{'admin_status'} = $i->{'status'};
			$big_obj->{'caches'}->{ $i->{'host_name'} }->{'profile'}      = $i->{'profile'};
			$big_obj->{'caches'}->{ $i->{'host_name'} }->{'domain_name'}  = $i->{'domain_name'};

			$big_obj->{'caches'}->{ $i->{'host_name'} }->{'ip_address'} = $i->{ip_address};
			push( @{ $big_obj->{'cachegroups'}->{ $i->{'cachegroup'} }->{'caches'} }, $i->{'host_name'} );
		}
	}

	return ($big_obj);
}

sub def_or_zero {
	my $val = shift;

	if ( defined($val) ) {
		return $val;
	}
	else {
		return 0;
	}
}

sub int_or_zero {
	my $val = shift;

	if ( defined($val) && ( $val =~ /^\d+?$/ || $val =~ /^\d+\.\d+$/ ) ) {
		return $val;
	}
	else {
		return 0;
	}
}

sub add_cache_stats {
	my $self     = shift;
	my $rascal   = shift;
	my $big_obj  = shift;
	my $resp_ref = shift;
	my $master_i = shift;

	if ( !defined($rascal) ) {
		$self->app->log->error('No running Rascal server found!');
		return;
	}

	my $args = { hc => 1, stats => "ats\.proxy\.process\.http\.current\_client\_connections\,bandwidth" };
	my $bigstats_hashref = $rascal->get_cache_stats($args);

	foreach my $server ( sort keys %{ $bigstats_hashref->{'caches'} } ) {
		if ( !defined( $big_obj->{'caches'}->{$server}->{'profile'} ) ) { next; }
		my $server_obj = $bigstats_hashref->{'caches'}->{$server};
		if ( exists( $server_obj->{'bandwidth'} ) ) {
			$self->app->log->trace("Processing server: $server");
			my $err_string = "";
			if ( defined( $server_obj->{'err -string'} ) ) {
				$err_string = $server_obj->{'err -string'}->[ $#{ $server_obj->{'err -string'} } ]->{'value'};
			}
			$big_obj->{'caches'}->{$server}->{'mbps_out'} = $server_obj->{'bandwidth'}->[0]->{'value'};
			$big_obj->{'caches'}->{$server}->{'connections'} =
				$server_obj->{'ats.proxy.process.http.current_client_connections'}->[0]->{'value'};

			$resp_ref->{'response'}->[ ${$master_i} ]->{'profile'}     = "$big_obj->{'caches'}->{$server}->{'profile'}";
			$resp_ref->{'response'}->[ ${$master_i} ]->{'hostname'}    = "$server";
			$resp_ref->{'response'}->[ ${$master_i} ]->{'cachegroup'}  = "$big_obj->{'caches'}->{$server}->{'cachegroup'}";
			$resp_ref->{'response'}->[ ${$master_i} ]->{'healthy'}     = &def_or_zero( $big_obj->{'caches'}->{$server}->{'healthy'} ) ? \1 : \0;
			$resp_ref->{'response'}->[ ${$master_i} ]->{'status'}      = "$big_obj->{'caches'}->{$server}->{'admin_status'}";
			$resp_ref->{'response'}->[ ${$master_i} ]->{'connections'} = 0 + &int_or_zero( $big_obj->{'caches'}->{$server}->{'connections'} );
			$resp_ref->{'response'}->[ ${$master_i} ]->{'kbps'}        = 0 + &int_or_zero( $big_obj->{'caches'}->{$server}->{'mbps_out'} );
			$resp_ref->{'response'}->[ ${$master_i} ]->{'ip'}          = "$big_obj->{'caches'}->{$server}->{'ip_address'}";
			${$master_i}++;
		}
	}
}

sub add_aggregate_cache_stats {
	my $self      = shift;
	my $big_obj   = shift;
	my $resp_ref  = shift;
	my $master_i  = shift;
	my $all_bw    = 0;
	my $all_conns = 0;
	foreach my $cachegroup ( sort keys %{ $big_obj->{'cachegroups'} } ) {
		my $total_bw    = 0;
		my $total_conns = 0;
		foreach my $cache ( @{ $big_obj->{'cachegroups'}->{$cachegroup}->{'caches'} } ) {
			$self->app->log->trace("Processing cache: $cache");
			if ( exists( $big_obj->{'caches'}->{$cache}->{'mbps_out'} ) ) {
				$total_bw += &int_or_zero( $big_obj->{'caches'}->{$cache}->{'mbps_out'} );
			}

			if ( exists( $big_obj->{'caches'}->{$cache}->{'connections'} ) ) {
				$total_conns += &int_or_zero( $big_obj->{'caches'}->{$cache}->{'connections'} );
			}
		}
		$self->app->log->trace("For cachegroup: $cachegroup, I found $total_bw Mbps, and $total_conns connections");
		$resp_ref->{'response'}->[ ${$master_i} ]->{'profile'}     = "ALL";
		$resp_ref->{'response'}->[ ${$master_i} ]->{'hostname'}    = "ALL";
		$resp_ref->{'response'}->[ ${$master_i} ]->{'cachegroup'}  = "$cachegroup";
		$resp_ref->{'response'}->[ ${$master_i} ]->{'healthy'}     = \1;
		$resp_ref->{'response'}->[ ${$master_i} ]->{'status'}      = "ALL";
		$resp_ref->{'response'}->[ ${$master_i} ]->{'connections'} = $total_conns;
		$resp_ref->{'response'}->[ ${$master_i} ]->{'kbps'}        = $total_bw;
		$resp_ref->{'response'}->[ ${$master_i} ]->{'ip'}          = undef;
		$all_bw    += $total_bw;
		$all_conns += $total_conns;
		${$master_i}++;
	}
	$resp_ref->{'response'}->[0]->{'profile'}     = "ALL";
	$resp_ref->{'response'}->[0]->{'hostname'}    = "ALL";
	$resp_ref->{'response'}->[0]->{'cachegroup'}  = "ALL";
	$resp_ref->{'response'}->[0]->{'healthy'}     = \1;
	$resp_ref->{'response'}->[0]->{'status'}      = "ALL";
	$resp_ref->{'response'}->[0]->{'connections'} = $all_conns;
	$resp_ref->{'response'}->[0]->{'kbps'}        = $all_bw;
	$resp_ref->{'response'}->[0]->{'ip'}          = undef;
}

1;
