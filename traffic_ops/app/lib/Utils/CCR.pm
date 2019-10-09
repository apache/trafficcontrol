package Utils::CCR;
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

use Mojo::UserAgent;
use utf8;
use Carp qw(cluck confess);
use JSON;
use Data::Dumper;

use constant {
	CCR_CACHE             => 1,
	CCR_DELIVERY_SERVICE  => 2,
	CCR_CACHE_UNAVAILABLE => "cache unavailable",
	CCR_CACHE_AVAILABLE   => "cache available",
	CCR_CACHE_UNKNOWN     => "cache unknown",
};

sub new {
	my $self     = {};
	my $class    = shift;
	my $ccr_host = shift;
	my $ccr_port = shift || 80;
	my $is_secure_port = shift || 0;

	if ( !defined($ccr_host) ) {
		confess("First constructor argument must be the CCR host");
	}
	elsif ( $ccr_port !~ m/^\d+$/ ) {
		confess("Second constructor argument must be the port number of the CCR host");
	}

	$self->{CCR_HOST}          = $ccr_host;
	$self->{CCR_PORT}          = $ccr_port;
	$self->{USER_AGENT}        = Mojo::UserAgent->new;
	$self->{FWD_PROXY}         = undef;
	$self->{IS_SECURE_PORT}    = $is_secure_port;

	return ( bless( $self, $class ) );
}

sub fwd_proxy {
	my $self = shift || confess("Call on an instance of Utils::Rascal");
	my $fwd_proxy = shift;

	if ( defined($fwd_proxy) ) {
		$self->{USER_AGENT}->proxy->http($fwd_proxy);
	}

	return ( $self->{FWD_PROXY} );
}

sub ua {
	my $self = shift || confess("Call on an instance of Utils::Rascal");
	return ( $self->{USER_AGENT} );
}

sub get_host {
	my $self = shift || confess("Call on an instance of Utils::CCR");
	return ( $self->{CCR_HOST} );
}

sub get_port {
	my $self = shift || confess("Call on an instance of Utils::CCR");
	return ( $self->{CCR_PORT} );
}

sub get_is_secure_port {
	my $self = shift || confess("Call on an instance of Utils::CCR");
	return ( $self->{IS_SECURE_PORT} );
}

sub get_url {
	my $self = shift || confess("Call on an instance of Utils::CCR");

	my $protocol = "http";
	if ( $self->get_is_secure_port() ) {
	    $protocol = "https";
	}

	my $url = "$protocol://" . $self->get_host() . ":" . $self->get_port();
	return ( $url );
}

sub get_location {
	my $self = shift || confess("Call on an instance of Utils::CCR");
	return ( $self->{LOCATION} );
}

sub set_location {
	my $self = shift || confess("Call on an instance of Utils::CCR");
	$self->{LOCATION} = shift || confess("Supply a location");
}

sub get_caches_by_location {
	my $self     = shift || confess("Call on an instance of Utils::CCR");
	my $location = shift || confess("Supply a location");
	my $url      = $self->get_url() . "/crs/locations/$location/caches";
	my $result   = $self->ua->get($url)->res->content->asset->slurp;
	my $content  = $response->content->asset->slurp;
	return Utils::Helper::ResponseHelper->handle_response( $response, $content );

	if ( defined($result) ) {
		return ( decode_json($result) );
	}
	else {
		return (undef);
	}
}

sub get_cache_state {
	my $self  = shift || confess("Call on an instance of Utils::CCR");
	my $cache = shift || confess("Supply a cache name");

	if ( !defined( $self->get_location() ) ) {
		confess("In order to call this method you must set the location via set_location");
	}

	return ( $self->get_cache_state_by_location( $cache, $self->get_location() ) );
}

sub get_cache_state_by_location {
	my $self     = shift || confess("Call on an instance of Utils::CCR");
	my $cache    = shift || confess("Supply a cache name");
	my $location = shift || confess("Supply a location");

	my $caches = $self->get_caches_by_location($location);

	if ( !defined($caches) ) {
		confess( "Unable to retrieve cache list from " . $self->get_host() );
	}

	for my $cache_ref ( @{$caches} ) {
		if ( $cache_ref->{cacheId} eq $cache ) {
			if ( $cache_ref->{cacheOnline} eq "true" ) {
				return (CCR_CACHE_AVAILABLE);
			}
			else {
				return (CCR_CACHE_UNAVAILABLE);
			}
		}
	}

	return (CCR_CACHE_UNKNOWN);
}

sub get_state {
	my $self = shift || confess("Call on an instance of Utils::Rascal");
	my $type = shift || confess("Supply a state to retrieve from CrStates");
	my $what = shift || confess("Supply a what to check");

	if ( $type == CCR_CACHE ) {
		return ( $self->get_cache_state($what) );
	}

	return (undef);
}

sub get_crs_stats {
	my $self     = shift || confess("Call on an instance of Utils::CCR");
	my $url      = $self->get_url() . "/crs/stats";
	my $response = $self->ua->get($url)->res;
	my $result   = $response->content->asset->slurp;

	if ( ( $response->code eq '200' ) && defined($result) && $result ne "" ) {
		return ( decode_json($result) );
	}
	else {
		return (undef);
	}
}

1;
