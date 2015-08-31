package UI::Utils;
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
use NetAddr::IP;
use Data::Dumper;

# !!!
# If you are using this module, put the use line before these:
# use Mojo::Base 'Mojolicious::Controller';
# use Data::Dumper;
# !!!

# Note: this is just a temp scheme, we may change this later
# Version is $major.$minor[.$micro]  - $micro is optional
# A release gets cut with just a $major.$minor
# The presence of a $micro means this version (branch) has been patched and released with that patch.
# Lowest $micro number, when present is 1.
my $version = "1.1.5-dev";

require Exporter;
our @ISA = qw(Exporter);

use constant READ  => 10;
use constant OPER  => 20;
use constant ADMIN => 30;

our %EXPORT_TAGS = (
	'all' => [
		qw(trim_whitespace is_admin is_oper log is_ipaddress is_ip6address is_netmask in_same_net is_hostname admin_status_id type_id
			profile_id profile_ids tm_version tm_url name_version_string is_regexp stash_role navbarpage rascal_hosts_by_cdn)
	]
);
our @EXPORT_OK = ( @{ $EXPORT_TAGS{all} } );
our @EXPORT    = ( @{ $EXPORT_TAGS{all} } );

sub tm_version {
	my $self = shift;

	return $version;
}

sub tm_url {
	my $self = shift;

	return $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.url', config_file => 'global' ] } )->get_column('value')->single();
}

sub trim_whitespace() {
	my $param = shift;

	if (ref($param) eq 'HASH') {
		foreach my $key (keys %{$param}) {
			${$param}{$key} =~ s/^\s+|\s+$//g;
		}
	} elsif (ref($param) eq 'ARRAY') {
		for ($i=0; $i <= $#{$param}; $i++) {
		   $param->[$i] =~ s/^\s+|\s+$//g;
		}
	} else {
		$param =~ s/^\s+|\s+$//g;
	}

	return $param;

}

sub name_version_string {
	my $self = shift;

	my $nv_string =
		$self->db->resultset('Parameter')->search( { -and => [ name => 'tm.toolname', config_file => 'global' ] } )->get_column('value')->single();
	$nv_string .= " ("
		. $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.url', config_file => 'global' ] } )->get_column('value')->single() . ")";
	return $nv_string;
}

sub is_regexp() {
	my $regexp = shift;

	eval {qr/$regexp/};

	return $@;
}

sub admin_status_id() {
	my $self        = shift;
	my $stat_string = shift;

	return $self->db->resultset('Status')->search( { name => $stat_string } )->get_column('id')->single();
}

sub type_id() {
	my $self        = shift;
	my $type_string = shift;

	return $self->db->resultset('Type')->search( { name => $type_string } )->get_column('id')->single();
}

sub profile_id() {
	my $self           = shift;
	my $profile_string = shift;

	return $self->db->resultset('Profile')->search( { name => $profile_string } )->get_column('id')->single();
}

sub profile_ids() {
	my $self           = shift;
	my $profile_string = shift;

	my @ids = $self->db->resultset('Profile')->search( { name => { -like => $profile_string } } )->get_column('id')->all();

	#print join("FF", @ids) . ":::::::::::::::::::::::::::::\n";
	return @ids;
}

sub is_hostname() {
	my $string = shift;

	if ( $string =~ /^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$/ ) {
		return 1;
	}
	else {
		return 0;
	}
}

sub is_ipaddress() {
	my $string = shift();

	my $ip = new NetAddr::IP($string);

	if ( defined($ip) ) {
		return 1;
	}
	else {
		return 0;
	}
}

sub is_ip6address() {
	my $string = shift();

	my $ip = new6 NetAddr::IP($string);

	if ( defined($ip) ) {
		return 1;
	}
	else {
		return 0;
	}
}

# yeah, yeah, lame check. also check with is_ipaddress.
sub is_netmask() {
	my $string = shift();

	if ( $string =~ /^255\.255\./ ) {
		return 1;
	}
	else {
		return 0;
	}
}

sub in_same_net() {
	my $ipstr1 = shift;
	my $ipstr2 = shift;

	#print $ipstr1 . " And " . $ipstr2 . "\n";
	my $ip1 = new NetAddr::IP($ipstr1);
	my $ip2 = new NetAddr::IP($ipstr2);

	if ( !defined($ip1) || !defined($ip2) ) { return 0; }

	my $network = $ip1->network();
	return $network->contains($ip2);
}

# log message to the log table
sub log() {
	my $self    = shift;
	my $message = shift;
	my $level   = shift;

	# my $user    = $self->tx->req->env->{REMOTE_USER};

	# For testing on local morbo
	#	if ( !defined($user) ) { $user = "jvando001"; }
	my $user;
	if ( $level eq 'CODEBIG' ) {
		$user = "codebig";
	}
	else {
		$user = $self->current_user()->{username};
	}

	$user = $self->db->resultset('TmUser')->search( { username => $user } )->get_column('id')->single;
	my $insert = $self->db->resultset('Log')->create(
		{
			tm_user => 0 + $user,    # the 0 + forces it to be treated as a number, and no ''
			message => $message,
			level   => $level,
		}
	);

}

# returns true if the user in $self has operations privs
sub is_oper() {
	my $self = shift;

	return &has_priv( $self, OPER );
}

# returns true if the user in $self has admin privs
sub is_admin() {
	my $self = shift;

	return &has_priv( $self, ADMIN );
}

## not exported ##

sub has_priv() {
	my $self     = shift;
	my $checkval = shift;

	my $user      = $self->current_user()->{username};
	my $user_data = $self->db->resultset('TmUser')->search( { username => $user } )->single;
	my $role      = "read-only";
	my $priv      = 10;
	if ( defined($user_data) ) {
		$role = $user_data->role->name;
		$priv = $user_data->role->priv_level;
	}
	return ( $priv >= $checkval );
}

sub stash_role {
	my $self = shift;

	if ( !defined( $self->current_user() ) ) {
		return $self->redirect_to('/login.html');
	}

	my $user = $self->current_user()->{username};
	my $role = $self->current_user()->{role};
	my $priv = $self->current_user()->{priv};

	$self->stash(
		user       => $user,
		role       => $role,
		priv_level => $priv,
	);
}

sub devmode {
	my $self = shift;

	my $devmode = undef;
	if ( $self->tx->local_port == 3000 ) {
		$devmode = 1;
	}
	return $devmode;
}

sub navbarpage {
	my $self = shift;

	&stash_role($self);

	$self->stash(
		devmode    => &devmode($self),
		navbarpage => 1,
	);

	# my @profiles = ();
	# my $rs_p = $self->db->resultset('Profile')->search( undef, { order_by => "name" } );
	# while ( my $row = $rs_p->next ) {
	# 	push( @profiles, $row->name );
	# }
	# $self->stash( profiles => \@profiles, );

	$self->res->headers->expires( Mojo::Date->new( time + 3600 ) );
}

sub rascal_hosts_by_cdn {
	my $self     = shift;
	my $params   = shift;
	my $cdn_name = exists( $params->{'cdn_name'} ) ? $params->{'cdn_name'} : undef;
	my $status   = exists( $params->{'status'} ) ? $params->{'status'} : undef;
	my $rascals;

	my $rs = $self->db->resultset('RascalHostsByCdnAll')->search();
	while ( my $row = $rs->next ) {
		if (   ( !$status && !$cdn_name )
			|| ( !$status                    && $cdn_name eq $row->cdn_name )
			|| ( !$cdn_name                  && $status eq $row->status )
			|| ( $cdn_name eq $row->cdn_name && $status eq $row->status ) )
		{
			push( @{ $rascals->{ $row->cdn_name }->{ $row->status } }, $row->host_name . "." . $row->domain_name );
		}
	}

	return $rascals;
}

sub exec_command {
	my ( $class, $command, @args ) = @_;
	my $pid    = fork();
	my $result = 0;

	if ( $pid == 0 ) {
		exec( $command, @args );
		exit 0;
	}
	else {
		wait;
		$result = $?;
		if ( $result != 0 ) {
			print "ERROR executing: $commands,  args: " . join( ' ', @args ) . "\n";
		}
	}
	return $result;
}

1;
