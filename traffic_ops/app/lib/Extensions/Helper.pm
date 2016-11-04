package Extensions::Helper;
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

use IO::Socket::SSL;
use Mojo::UserAgent;
use Data::Dumper;

use Log::Log4perl qw(:easy);
use JSON;
use strict;
use warnings;

use constant {
	API_VERSION     => "1.2",
	DEFAULT_TOKEN   => '91504CE6-8E4A-46B2-9F9F-FE7C15228498',
	SERVERLIST_PATH => '/api/1.2/servers.json',
	DSLIST_PATH     => '/api/1.2/deliveryservices.json',
	PARAMETER_PATH  => '/api/1.2/parameters'
};

require Exporter;
our @ISA         = qw(Exporter);
our %EXPORT_TAGS = ( 'all' => [qw( )] );
our @EXPORT_OK   = ( @{ $EXPORT_TAGS{'all'} } );
our @EXPORT      = qw(

);

our $ua;
our $b_url;
our $token;

sub new {
	my $proto = shift;
	my $self  = {};

	bless( $self, $proto );
	$self->_init(@_);

	return ($self);
}

sub _init {
	my $self   = shift;
	my $config = shift;

	$b_url = $config->{base_url};
	if ( defined( $config->{token} ) ) {
		$token = $config->{token};
	}
	else {
		$token = DEFAULT_TOKEN;
	}
	$self->_session();
}

## create a traffic ops sesssion
sub _session {
	my $self = shift;

	my $url = $b_url . '/api/' . API_VERSION . '/user/login/token';
	DEBUG "session " . $url;

	IO::Socket::SSL::set_defaults(
		verify_hostname => 0,
		SSL_verify_mode => SSL_VERIFY_NONE
	);
	$ua = Mojo::UserAgent->new;
	my $tx = $ua->post( $url => json => { t => $token } );
	if ( my $res = $tx->success ) { TRACE Dumper( $res->body ) }
	else {
		my $err = $tx->error;
		ERROR Dumper($err);
		die "$err->{code} response: $err->{message}" if $err->{code};
		die "Connection error: $err->{message}";
	}
	return 1;
}

sub get {
	my $self = shift;
	my $path = shift;

	DEBUG "get " . $b_url . $path;
	my $tx = $ua->get( $b_url . $path );
	if ( my $res = $tx->success ) {
		my $jresp = JSON::decode_json( $res->body );
		TRACE "Success (" . length( $res->body ) . " bytes)\n";
		return $jresp->{response};
	}
	else {
		ERROR Dumper( $tx->error );
	}
}

# post json object to any path
sub post_json {
	my $self = shift;
	my $path = shift;
	my $json = shift;

	my $url = $b_url . $path;
	DEBUG "Post: " . $url . " data:" . Dumper($json);
	my $tx = $ua->post( $url => json => $json );
	if ( my $res = $tx->success ) {
		my $jresp = JSON::decode_json( $res->body );
		TRACE "Success: " . Dumper( $res->body );
		return $res;
	}
	else {
		ERROR Dumper( $tx->error );
		return $tx->error;
	}
}

# post the check results
sub post_result {
	my $self       = shift;
	my $server_id  = shift;
	my $check_name = shift;
	my $result     = shift;

	my $r = { id => $server_id, servercheck_short_name => $check_name, value => $result };
	my $path = "/api/" . API_VERSION . "/servercheck";
	return $self->post_json( $path, $r );
}

1;

