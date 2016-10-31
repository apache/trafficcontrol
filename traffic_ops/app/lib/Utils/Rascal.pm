package Utils::Rascal;
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
use File::Slurp;
use Utils::Helper::ResponseHelper;

use constant {
	RASCAL_CACHE                        => 1,
	RASCAL_DELIVERY_SERVICE             => 2,
	RASCAL_CACHE_UNAVAILABLE            => "cache unavailable",
	RASCAL_CACHE_AVAILABLE              => "cache available",
	RASCAL_CACHE_UNKNOWN                => "cache unknown",
	RASCAL_DELIVERY_SERVICE_UNAVAILABLE => "delivery service unavailable",
	RASCAL_DELIVERY_SERVICE_AVAILABLE   => "delivery service available",
	RASCAL_DELIVERY_SERVICE_UNKNOWN     => "delivery service unknown",
};

sub new {
	my $self        = {};
	my $class       = shift;
	my $rascal_host = shift;
	my $rascal_port = shift || 80;

	if ( !defined($rascal_host) ) {
		confess("First constructor argument must be the Rascal host");
	}
	elsif ( $rascal_port !~ m/^\d+$/ ) {
		confess("Second constructor argument must be the port number of the Rascal host");
	}

	$self->{RASCAL_HOST} = $rascal_host;
	$self->{RASCAL_PORT} = $rascal_port;
	$self->{USER_AGENT}  = Mojo::UserAgent->new;
	$self->{FWD_PROXY}   = undef;

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
	my $self = shift || confess("Call on an instance of Utils::Rascal");
	return ( $self->{RASCAL_HOST} );
}

sub get_port {
	my $self = shift || confess("Call on an instance of Utils::Rascal");
	return ( $self->{RASCAL_PORT} );
}

sub get_url {
	my $self = shift || confess("Call on an instance of Utils::Rascal");
	my $uri  = shift || "";

	if ( $uri !~ m/^\// ) {
		$uri = "/" . $uri;
	}

	my $url = "http://" . $self->get_host() . ":" . $self->get_port() . $uri;
	return ($url);
}

sub get_stat {
	my $self = shift || confess("Call on an instance of UI::Rascal");
	my $stat = shift || confess("Supply a stat name");

	my $stats = $self->get_stats();

	if ( exists( $stats->{$stat} ) ) {
		return ( $stats->{$stat} );
	}
	else {
		return (undef);
	}
}

sub get_cr_config {
	my $self     = shift || confess("Call on an instance of UI::Rascal");
	my $url      = $self->get_url("/publish/CrConfig\?json");
	my $response = $self->ua->get($url)->res;
	my $content  = $response->content->asset->slurp;

	return Utils::Helper::ResponseHelper->handle_response( $response, $content );
}

sub get_stats {
	my $self          = shift || confess("Call on an instance of UI::Rascal");
	my $url           = $self->get_url("/publish/Stats");
	my $response      = $self->ua->get($url)->res;
	my $content       = $response->content->asset->{'content'};
	my $json_response = Utils::Helper::ResponseHelper->handle_response( $response, $content );

	if ( defined($json_response) ) {

		# kept this because of the inner json response
		return ( $json_response->{stats} );
	}
	else {
		return (undef);
	}
}

sub get_cr_states {
	my $self = shift || confess("Call on an instance of UI::Rascal");
	my $raw  = shift;
	my $url  = $self->get_url("/publish/CrStates");

	if ( defined($raw) && $raw ) {
		$url .= "?raw";
	}

	my $response = $self->ua->get($url)->res;
	my $content  = $self->ua->get($url)->res->content->asset->slurp;

	return Utils::Helper::ResponseHelper->handle_response( $response, $content );
}

sub get_states {
	my $self  = shift || confess("Call on an instance of UI::Rascal");
	my $state = shift || confess("Supply a state to retrieve from CrStates");
	my $raw   = shift;

	my $cr_states = $self->get_cr_states($raw);

	#write_file( "/tmp/states.log", "cr_states: " . Dumper($cr_states) . "\n" );

	if ( defined($cr_states) && exists( $cr_states->{$state} ) ) {
		return ( $cr_states->{$state} );
	}

	return (undef);
}

sub get_state {
	my $self = shift || confess("Call on an instance of UI::Rascal");
	my $type = shift || confess("Supply a state to retrieve from CrStates");
	my $what = shift || confess("Supply a what to check");

	if ( $type == RASCAL_CACHE ) {
		return ( $self->get_cache_state($what) );
	}
	elsif ( $type == RASCAL_DELIVERY_SERVICE ) {
		return ( $self->get_delivery_service_state($what) );
	}

	return (undef);
}

sub get_cache_states {
	my $self = shift || confess("Call on an instance of UI::Rascal");
	return ( $self->get_states("caches") );
}

sub get_cache_state {
	my $self  = shift || confess("Call on an instance of UI::Rascal");
	my $cache = shift || confess("Supply a cache name");

	my $cache_states = $self->get_cache_states();

	if ( !defined($cache_states) ) {
		confess( "Unable to retrieve cache states from " . $self->get_host() );
	}

	if ( exists( $cache_states->{$cache} ) ) {
		if ( $cache_states->{$cache}->{isAvailable} eq "true" ) {
			return (RASCAL_CACHE_AVAILABLE);
		}
		else {
			return (RASCAL_CACHE_UNAVAILABLE);
		}
	}

	return (RASCAL_CACHE_UNKNOWN);
}

sub get_delivery_service_states {
	my $self = shift || confess("Call on an instance of UI::Rascal");
	return ( $self->get_states("deliveryServices") );
}

sub get_delivery_service_state {
	my $self             = shift || confess("Call on an instance of UI::Rascal");
	my $delivery_service = shift || confess("Supply a delivery service");

	my $delivery_service_states = $self->get_delivery_service_states();

	if ( !defined($delivery_service_states) ) {
		confess( "Unable to retrieve delivery service states from " . $self->get_host() );
	}

	if ( exists( $delivery_service_states->{$delivery_service} ) ) {
		if ( $delivery_service_states->{$delivery_service}->{isAvailable} eq "true" ) {
			return (RASCAL_DELIVERY_SERVICE_AVAILABLE);
		}
		else {
			return (RASCAL_DELIVERY_SERVICE_UNAVAILABLE);
		}
	}

	return (RASCAL_DELIVERY_SERVICE_UNKNOWN);
}

sub get_cache_stats {
	my $self = shift || confess("Call on an instance of UI::Rascal");
	my $args = shift;

	my $url = $self->get_url("/publish/CacheStats?hc=1");

	if ( defined($args) && ref($args) eq "HASH" ) {
		for my $key ( keys( %{$args} ) ) {
			$url .= "&";
			$url .= $key;
			$url .= "=";
			$url .= $args->{$key};
		}
	}

	my $response = $self->ua->get($url)->res;
	my $content  = $response->content->asset->slurp;
	return Utils::Helper::ResponseHelper->handle_response( $response, $content );

}

sub get_cache_stats_by_host {
	my $self = shift || confess("Call on an instance of UI::Rascal");
	my $host = shift || confess("Supply a host");
	my $url  = $self->get_url("/publish/CacheStats/$host/?hc=1");
	my $result   = $self->ua->get($url)->res->content->asset->{'content'};
	my $response = $self->ua->get($url)->res;
	my $content  = $response->content->asset->slurp;
	return Utils::Helper::ResponseHelper->handle_response( $response, $content );
}

1;
