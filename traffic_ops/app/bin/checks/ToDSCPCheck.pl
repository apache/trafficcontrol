#!/usr/bin/perl
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
# DSCP check extension. Populates the 'DSCP' column.
#

# example cron entry
# 0 * * * * root /opt/traffic_ops/app/bin/checks/ToDSCPCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"DSCP\", \"cms_interface\": \"eth0\"}" >> /var/log/traffic_ops/extensionCheck.log 2>&1
# example cron entry with syslog
# 0 * * * * root /opt/traffic_ops/app/bin/checks/ToDSCPCheck.pl -c "{\"base_url\": \"https://localhost\", \"check_name\": \"DSCP\", \"name\": \"Delivery Service\", \"cms_interface\": \"eth0\", \"syslog_facility\": \"local0\"}" > /dev/null 2>&1

use strict;
use warnings;

use Data::Dumper;
use Getopt::Std;
use Log::Log4perl qw(:easy);
use Net::PcapUtils;
use NetAddr::IP;
use NetPacket::Ethernet qw(:strip);
use NetPacket::IP qw(:strip);
use NetPacket::IPv6 qw(:strip);
use NetPacket::TCP;
use JSON;
use Extensions::Helper;
use Sys::Syslog qw(:standard :macros);
use IO::Handle;

# TODO Figure out how to make this match the TC release version
my $VERSION = "0.07";

STDOUT->autoflush(1);

my %args = ();
getopts( "c:f:hl:qs:x:", \%args );

if ( $args{h} ) {
	&help();
	exit();
}

Log::Log4perl->easy_init($ERROR);
if ( defined( $args{l} ) ) {
	if    ( $args{l} == 1 ) { Log::Log4perl->easy_init($FATAL); }
	elsif ( $args{l} == 2 ) { Log::Log4perl->easy_init($ERROR); }
	elsif ( $args{l} == 3 ) { Log::Log4perl->easy_init($WARN); }
	elsif ( $args{l} == 4 ) { Log::Log4perl->easy_init($INFO); }
	elsif ( $args{l} == 5 ) { Log::Log4perl->easy_init($DEBUG); }
	elsif ( $args{l} == 6 ) { Log::Log4perl->easy_init($TRACE); }
	else                    { Log::Log4perl->easy_init($INFO); }
}

DEBUG("Including DEBUG messages in output");
TRACE("Including TRACE messages in output");

if ( !defined( $args{c} ) ) {
	&help();
	exit(1);
}

my $jconf = undef;
eval { $jconf = decode_json( $args{c} ) };
if ($@) {
	FATAL "FATAL Bad json config: $@";
	print "\n\n";
	&help();
	exit(1);
}

my $sslg = undef;
my $chck_lng_nm;
if ( defined( $jconf->{syslog_facility} ) ) {
	$chck_lng_nm = $jconf->{name};
	setlogmask( LOG_UPTO(LOG_INFO) );
	openlog( 'ToChecks', '', $jconf->{syslog_facility} );
	$sslg = 1;
}

my $cms_int = undef;
if ( defined( $jconf->{cms_interface} ) ) {
	$cms_int = $jconf->{cms_interface};
}
else {
	FATAL "FATAL cms_interface must be defined.";
	print "\n\n";
	&help();
	exit(1);
}

my $force = 0;
if ( defined( $args{f} ) ) {
	$force = $args{f};
	DEBUG "DEBUG force: " . $args{f} . "";
}

my $quiet;
if ( $args{q} ) {
	$quiet = 1;
}

my $check_name = $jconf->{check_name};
if ( $check_name ne "DSCP" ) {
	FATAL "FATAL This Check Extension is exclusively for DSCP.";
	print "\n\n";
	&help();
	exit(4);
}

DEBUG "DEBUG " . Dumper($jconf);
my $b_url = $jconf->{base_url};
Extensions::Helper->import();
my $ext = Extensions::Helper->new( { base_url => $b_url, token => '91504CE6-8E4A-46B2-9F9F-FE7C15228498' } );

my %ds_info           = ();
my $jdeliveryservices = $ext->get(Extensions::Helper::DSLIST_PATH);

# Get all the Delivery Services
foreach my $ds ( @{$jdeliveryservices} ) {
	$ds_info{ $ds->{id} } = $ds;

	# Get the DS details and add the matchList
	my $ds_details = $ext->get( '/api/4.0/deliveryservices/' . $ds->{id} );

	#DEBUG "Debug matchList before push: " . Dumper($ds_details->[0]{matchList});
	# The matchList was removed at one point and now it's back.
	# Leaving this here for when that happens again.
	#push( @{ $ds_info{ $ds->{id} }{matchList} }, @{ $ds_details->[0]{matchList} } );
	#DEBUG "Debug matchList after push: " . Dumper($ds_info{ $ds->{id} }{matchList});
}

my $jdataserver = ();
if ( defined( $args{s} ) ) {
	DEBUG "DEBUG Getting single server";

	# TODO check for a successful return
	# This returns a reference to a hash
	$jdataserver = $ext->get( '/api/4.0/servers/hostname/' . $args{s} . '/details' );

	# create an array
	my @tmp = ();

	# put the hash reference in the array
	$tmp[0] = $jdataserver;

	# clear $jdataserver
	$jdataserver = ();

	# create a reference to the array
	$jdataserver = \@tmp;

}
else {
	# This is a reference to an array
	$jdataserver = $ext->get(Extensions::Helper::SERVERLIST_PATH);
}

my %cdns  = ();
my $jcdns = $ext->get('/api/4.0/cdns');

# convert the cdns to a hash
foreach my $cdn ( @{$jcdns} ) {
	DEBUG "DEBUG cdn_name: " . $cdn->{name};
	$cdns{ $cdn->{name} } = $cdn;
	DEBUG "DEBUG domain_name: " . $cdns{ $cdn->{name} }->{domainName};
}

foreach my $server ( @{$jdataserver} ) {
	next unless $server->{type} =~ m/^EDGE/;    # We know this is DSCP, so we know we want edges only
	if ( $server->{status} ne 'REPORTED' ) {
		INFO "INFO Skipping server: $server->{hostName} because status is: $server->{status}";
		next;
	}
	my @ip_addrs    = ( $server->{ipAddress}, $server->{ip6Address} );
	my $host_name   = trim( $server->{hostName} );
	my $fqdn        = $host_name . "." . trim( $server->{domainName} );
	my $interface   = trim( $server->{interfaceName} );
	my $status      = $server->{status};
	my $port        = $server->{tcpPort};                                                    # this is for both http and https
	my $protocol    = "http";
	my $ip_protocol = "ipv4";
	my $domain_name = undef;
	my $details     = $ext->get( '/api/4.0/servers/hostname/' . $host_name . '/details' );
	my $successful  = 1;                                                                     # assume all is good

	if ( ( defined( $ip_addrs[1] ) ) && ( $ip_addrs[1] =~ m/:/ ) ) {
		$ip_addrs[1] = new NetAddr::IP( $ip_addrs[1] );
		$ip_addrs[1] =~ s/\/\d+$//;
		$ip_addrs[1] = lc( $ip_addrs[1] );
	}

	DEBUG "DEBUG Checking: " . $host_name;

	# Loop thru each delivery service associated with this server
	foreach my $dsid ( @{ $details->{deliveryservices} } ) {
		my $ds = $ds_info{$dsid};
		if ( defined( $args{x} ) ) {

			# TODO figure out how to not loop through all DSs
			next unless trim( $ds->{xmlId} ) eq $args{x};
		}
		$domain_name = $cdns{ $ds->{cdnName} }->{domainName};
		DEBUG "DEBUG xmlId: "
			. $ds->{xmlId}
			. " dsid: "
			. $dsid
			. " cdnName: "
			. $ds->{cdnName}
			. " domain_name: "
			. $domain_name
			. " active: "
			. $ds->{active}
			. " checkpath: "
			. ( defined( $ds->{checkPath} ) ? $ds->{checkPath} : "not defined" )
			. " protocol: "
			. $ds->{protocol};

		if ( $ds->{protocol} == 1 || $ds->{protocol} == 3 ) {
			$protocol = "https";
			$port     = trim( $server->{httpsPort} );
		}

		if ( $ds->{active} && defined( $ds->{checkPath} ) && $ds->{checkPath} ne "" ) {
			my $target = $ds->{xmlId} . ":" . $ip_protocol;
			my $header = "";
			if ( !defined( $ds->{matchList} ) || ( $ds->{matchList} eq "" ) ) {
				$successful = 0;
				ERROR "ERROR Failed deliveryService xmlId: " . $ds->{xmlId} . " matchList undefined";
				if ($sslg) {
					my @tmp = ( $fqdn, $check_name, $chck_lng_nm, 'FAIL', $status, $target, $header, "matchList undefined" );
					syslog( LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s status=%s target=%s url=%s msg=\"%s\"", @tmp );
				}
			}
			else {
				foreach my $match ( @{ $ds->{matchList} } ) {
					TRACE "TRACE checking match: " . $match->{pattern};
					my $prefix;
					if ( $match->{type} eq 'HOST_REGEXP' ) {
						if ( $match->{pattern} =~ /\*/ ) {
							my $tmp = $match->{pattern};
							$tmp =~ s/\\//g;
							$tmp =~ s/\.\*//g;
							$header .= $tmp;
							if ( ( $ds->{type} =~ /^DNS/ ) && ( $match->{pattern} =~ m/^\.\*/ ) ) {
								DEBUG "DEBUG setting prefix to edge";
								$prefix = 'edge';
							}
							elsif ( $ds->{type} =~ /^DNS/ ) {
								DEBUG "DEBUG setting prefix to ''";
								$prefix = '';
							}
							else {
								# this is the hostname, not the fqdn
								$prefix = $host_name;
							}

							# TODO There may still be a bug here when a DS has a regex
							# that is a FQDN and ds->{type} of HTTP.
							$header .= $domain_name;
							$header = $prefix . $header;
						}
						else {
							# TODO How to handle the other regex patterns
							TRACE "TRACE Not a Host Regex. Going to next regex.";
							next;
						}
					}
					foreach my $ip (@ip_addrs) {
						TRACE "TRACE checking " . $ip;
						if ( !defined($ip) ) {
							TRACE "TRACE IP not defined. Should go to next ip";
							next;
						}
						elsif ( $ip =~ m/\./ ) {
							DEBUG "DEBUG ipv4: " . $ip;
							$ip_protocol = "ipv4";
						}
						elsif ( $ip =~ m/:/ ) {
							$ip_protocol = "ipv6";
							$ip          = new NetAddr::IP($ip);
							$ip =~ s/\/\d+$//;
							$ip = lc($ip);
							DEBUG "DEBUG ipv6: " . $ip;
						}
						else {
							TRACE "TRACE Not an IPv4 or IPv6 IP. Should go to next IP.";
							next;
						}
						my $checkpath = $ds->{checkPath};
						$checkpath =~ s/^\///;
						my $url = $protocol . "://" . $ip . ":" . $port . "/" . $checkpath;

						DEBUG "DEBUG About to check header: " . $header . " url: " . $url . " ip_protocol: " . $ip_protocol;

						my $dscp_found;
						if ( $force == 0 ) {
							$dscp_found = get_dscp( $url, $ip, $cms_int, $header, $ip_protocol, $port );
							if ( ( $dscp_found == -1 ) || ( $dscp_found != $ds->{dscp} ) ) {
								sleep 1;
								DEBUG "DEBUG About to check again header: " . $header . " url: " . $url . " ip_protocol: " . $ip_protocol;
								$dscp_found = &get_dscp( $url, $ip, $cms_int, $header, $ip_protocol, $port );
							}
						}
						elsif ( $force == 1 ) {
							$dscp_found = -1;
						}
						elsif ( $force == 2 ) {
							$dscp_found = $ds->{dscp};
						}
						elsif ( $force == 3 ) {
							$dscp_found = $ds->{dscp} + 1;
						}
						if ( $dscp_found == -1 ) {
							$successful = 0;
							ERROR "ERROR Failed xmlId: " . $ds->{xmlId};
							if ($sslg) {
								my @tmp = ( $fqdn, $check_name, $chck_lng_nm, 'FAIL', $status, $target, $header, "Unable to connect to edge server" );
								syslog( LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s status=%s target=%s url=%s msg=\"%s\"", @tmp );
							}
						}
						elsif ( $dscp_found == $ds->{dscp} ) {
							INFO "INFO Success xmlId: " . $ds->{xmlId};
							if ($sslg) {
								my @tmp = ( $fqdn, $check_name, $chck_lng_nm, 'OK', $status, $target, $header );
								syslog( LOG_INFO, "hostname=%s check=%s name=\"%s\" result=%s status=%s target=%s url=%s msg=\"\"", @tmp );
							}
						}
						else {
							$successful = 0;
							ERROR "ERROR Failed xmlId: " . $ds->{xmlId};
							if ($sslg) {
								$successful = 0;
								my @tmp = (
									$fqdn, $check_name, $chck_lng_nm, 'FAIL', $status,
									$target, $header, "Expected DSCP value of $ds->{dscp} got $dscp_found"
								);
								syslog( LOG_ERR, "hostname=%s check=%s name=\"%s\" result=%s status=%s target=%s url=%s msg=\"%s\"", @tmp );
							}
						}
					}
				}
			}
			DEBUG "DEBUG Finished checking";
		}
	}
	if ($successful) {
		$ext->post_result( $server->{id}, $check_name, 1 ) if ( !$quiet );
	}
	else {
		$ext->post_result( $server->{id}, $check_name, 0 ) if ( !$quiet );
	}
}

closelog();

sub get_dscp {
	my $url      = shift;
	my $ip       = shift;
	my $dev      = shift;
	my $header   = shift;
	my $protocol = shift;
	my $port     = shift;

	my $tos     = undef;
	my $max_len = 0;

	my $src_port = int( rand( 65535 - 1024 ) ) + 1024;
	TRACE "TRACE In sub get_dscp";
	DEBUG "DEBUG get_dscp ip:" . $ip . " url:" . $url . " dev:" . $dev . " port:" . $src_port . " ip protocol: " . $protocol;

	# Use curl to get some traffic from the URL, but send the command to the background, so the capture that follows
	# is while traffic is being returned
	my $curl = 'curl --connect-timeout 2 --local-port ' . $src_port . ' --' . $protocol . ' -s -H "Host: ' . $header . '" ' . $url . ' 2>&1 > /dev/null';
	DEBUG "DEBUG curl: " . $curl;
	if ( $protocol eq 'ipv6' ) {
		DEBUG "DEBUG running ip6";
		system( "(sleep 1; " . $curl . " || ping6 -c 10 $ip 2>&1 > /dev/null)  &" );
	}
	else {
		DEBUG "DEBUG running ip4";
		system( "(sleep 1; " . $curl . " || ping -c 10 $ip 2>&1 > /dev/null)  &" );
	}

	Net::PcapUtils::loop(
		sub {
			my ( $user, $hdr, $pkt ) = @_;
			my $ip_obj;
			$ip_obj->{src_ip} = new NetAddr::IP( $ip_obj->{src_ip} );
			if ( $protocol eq "ipv4" ) {
				$ip_obj = NetPacket::IP->decode( eth_strip($pkt) );
				DEBUG "DEBUG <=> $ip_obj->{src_ip} -> $ip_obj->{dest_ip} proto: $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
				my $tcp_obj = NetPacket::TCP->decode( $ip_obj->{data} );
				DEBUG
					"DEBUG TCP1 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} proto: $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
				if ( $ip_obj->{src_ip} eq $ip && $ip_obj->{len} > $max_len && $ip_obj->{proto} == 6 ) {
					my $tcp_obj = NetPacket::TCP->decode( $ip_obj->{data} );
					DEBUG
						"DEBUG TCP2 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
					if ( ( $tcp_obj->{src_port} == $port ) && ( $tcp_obj->{dest_port} == $src_port ) ) {
						DEBUG
							"DEBUG TCP3 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} $ip_obj->{proto} tos $ip_obj->{tos} len $ip_obj->{len}\n";
						$max_len = $ip_obj->{len};
						$tos     = $ip_obj->{tos};
					}
				}
			}
			elsif ( $protocol eq "ipv6" ) {
				$ip_obj = NetPacket::IPv6->decode( eth_strip($pkt) );
				DEBUG "DEBUG <=> $ip_obj->{src_ip} -> $ip_obj->{dest_ip} proto: $ip_obj->{nxt} tos $ip_obj->{class} len $ip_obj->{plen}\n";
				my $tcp_obj = NetPacket::TCP->decode( $ip_obj->{data} );

				# DEBUG "DEBUG Comparison of IPs ip_obj->{src_ip}: $ip_obj->{src_ip}, ip: $ip\n"
				DEBUG
					"DEBUG TCP1 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} proto: $ip_obj->{nxt} tos $ip_obj->{class} len $ip_obj->{plen}\n";
				if ( $ip_obj->{src_ip} eq $ip && $ip_obj->{plen} > $max_len && $ip_obj->{nxt} == 6 ) {
					my $tcp_obj = NetPacket::TCP->decode( $ip_obj->{data} );
					DEBUG
						"DEBUG TCP2 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} $ip_obj->{nxt} tos $ip_obj->{class} len $ip_obj->{plen}\n";
					if ( $tcp_obj->{src_port} == $port && $tcp_obj->{dest_port} == $src_port ) {
						DEBUG
							"DEBUG TCP3 $ip_obj->{src_ip}:$tcp_obj->{src_port} -> $ip_obj->{dest_ip}:$tcp_obj->{dest_port} $ip_obj->{nxt} tos $ip_obj->{class} len $ip_obj->{plen}\n";
						$max_len = $ip_obj->{plen};
						$tos     = $ip_obj->{class};
					}
				}
			}
		},
		FILTER     => 'host ' . $ip,
		DEV        => $dev,
		NUMPACKETS => 7,
		TIMEOUT    => 10
	);

	my $dscp;
	if ( defined($tos) ) {
		$dscp = $tos >> 2;
	}
	else {
		$dscp = -1;
	}
	TRACE "TRACE exiting sub get_dscp returning $dscp for dscp";
	return $dscp;
}

sub ltrim { my $s = shift; $s =~ s/^\s+//;       return $s }
sub rtrim { my $s = shift; $s =~ s/\s+$//;       return $s }
sub trim  { my $s = shift; $s =~ s/^\s+|\s+$//g; return $s }

sub help() {
	print
		"ToDSCPCheck.pl -c \"{\\\"base_url\\\": \\\"https://localhost\\\", \\\"check_name\\\": \\\"DSCP\\\", \\\"cms_interface\\\": \\\"eth0\\\"[, \\\"name\\\": \\\"DSCP Service Check\\\", \\\"syslog_facility\\\": \\\"local0\\\"]}\" [-f <1-3>] [-l <1-6>] [-q] [-s <hostname>] [-x <xmlId>]\n";
	print "\n";
	print "-c   json formatted list of variables\n";
	print "     base_url: required\n";
	print "        URL of the Traffic Ops server.\n";
	print "     check_name: required\n";
	print "        The name of this check.\n";
	print "     cms_interface: required\n";
	print "        Interface used to communicate with edges.\n";
	print "     name: optional\n";
	print "        The long name of this check. used in conjuction with syslog_facility.\n";
	print "     syslog_facility: optional\n";
	print "        The syslog facility to send messages. Requires the \"name\" option to\n";
	print "        be set.\n";
	print "-f   Force a FAIL or OK message\n";
	print "        1: FAIL Unable to connect to edge server.\n";
	print "        2: OK\n";
	print "        3: FAIL DSCP values didn't match.\n";
	print "-h   Print this message\n";
	print "-l   Logging level. 1 - 6. 1 being least (FATAL). 6 being most (TRACE). Default\n";
	print "     is 2 (INFO).\n";
	print "-q   Don't post results to Traffic Ops.\n";
	print "-s   Check a specific server. Host name as it appears in Traffic Ops. Not FQDN.\n";
	print "-x   Check a specific xmlId.\n";
	print "================================================================================\n";

	# the above line of equal signs is 80 columns
	print "\n";
}
