package UI::GenIso;
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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';

my $filebasedir = "/var/www/files";
my $ksfiles_parm_name = "kickstart.files.location";
my $ksfiles_configfile_name = "mkisofs";

sub geniso {
	my $self = shift;

	&navbarpage($self);
	my %serverselect;
	my $rs_server = $self->db->resultset('Server')->search( undef, { columns => [qw/id host_name domain_name/], orderby => "host_name" } );

	while ( my $row = $rs_server->next ) {
		my $fqdn = $row->host_name . "." . $row->domain_name;
		$serverselect{$fqdn} = $row->id;
	}

	my $osversionsdir;
	my $ksdir = $self->db->resultset('Parameter')->search( { -and => [ name => $ksfiles_parm_name, config_file => $ksfiles_configfile_name ] } )->get_column('value')->single();
	if (defined $ksdir && $ksdir ne "") {
		$osversionsdir = $ksdir;
	} else {
		$osversionsdir = $filebasedir;
	}

	my %osversions;

	{
		open(CFG, "<$osversionsdir/osversions.cfg") || die("$osversionsdir/osversions.cfg:$!");
		local $/;
		eval <CFG>;
		close CFG;
	}

	$self->stash(
		serverselect => \%serverselect,
		osversions   => \%osversions,
	);
}

sub iso_download {
	my $self = shift;
	my $hostname = $self->param('hostname');
	my $osversion = $self->param('osversion');
	my $rootpass = $self->param('rootpass');
	my $dhcp = $self->param('dhcp');
	my $ipaddr = $self->param('ipaddr');
	my $netmask = $self->param('netmask');
	my $gateway = $self->param('gateway');
	my $ip6_address = $self->param('ip6_address');
	my $ip6_gateway = $self->param('ip6_gateway');
	my $dev = $self->param('dev');
	my $mtu = $self->param('mtu');
	my $ondisk = $self->param('ondisk');
	my $lacp;
	$lacp = 1 if ($dev =~ m/^bond0$/);

	$self->res->headers->content_type("application/download");
	$self->res->headers->content_disposition("attachment; filename=\"$hostname-$osversion.iso\"");

	my $dir;
	my $ksdir = $self->db->resultset('Parameter')->search( { -and => [ name => $ksfiles_parm_name, config_file => $ksfiles_configfile_name ] } )->get_column('value')->single();
	if (defined $ksdir && $ksdir ne "") {
		$dir = $ksdir . "/" . $osversion;
	} else {
		$dir = $filebasedir . "/" . $osversion;
	}

	my $cmd = "mkisofs -input-charset utf-8 -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table -R -J -v -T $dir";

	open(IN, "<$dir/ks.src") || die("$dir/ks.src:$!");
	my @ks = <IN>;
	close IN;

	open (OUT, ">$dir/ks.cfg") || die("$dir/ks.cfg:$!");
	for my $line(@ks) {
		# Create Network line
		if ($line =~ m/^network/) {
			my $net_line;

			# Use ip6 gateway only if ipv4 gateway is not defined
			$gateway = $ip6_gateway if ($gateway =~ m/^\s*$/);

			if ($dhcp eq 'yes') {
				$net_line = "network --bootproto=dhcp --hostname=$hostname";
			} else {
				$net_line = "network --bootproto=static --ip=$ipaddr --netmask=$netmask --gateway=$gateway --nameserver=69.252.80.80 --hostname=$hostname --mtu=$mtu";
			}

			# Fill in a default value for the network device if one isn't specified
			if ($dev =~ m/^\s*$/) {
				$net_line .= " --device=link";
			} else {
				$net_line .= " --device=$dev";
			}

			# IPV6 stuff
			if ($ip6_address =~ m/^\s*$/) {
				$net_line .= " --noipv6";
			} else {
				$net_line .= " --ipv6=$ip6_address";
				if ($ip6_address =~ m/^\s*$/) {
					$net_line .= " --ipv6gateway=$ip6_gateway"; 
				}
				# There should probably be some sort of error thrown if there is an ipv6 
				# address without a gateway. 
			}

			$line = $net_line;
		}

		# Set the disk to use
		if ($line =~ m/^ignoredisk/) {
			if ($ondisk !~ m/^\s?$/) {
				$line = "ignoredisk --only-use=$ondisk\n";
			} else {
				$line = undef;
			}
		}

		# Set rootpass
		if ($rootpass =~ m/^\S+$/) {
			if ($line =~ m/^rootpw/) {
				$line = "rootpw $rootpass\n";
			}
		}

		# Place additional stuff at the bottom of the file here
		if ($line =~ m/^eject$/) {
			if ($lacp) {
				my $string;
				if ($ip6_address =~ m/^\s*$/) {
					$string = "echo -e 'DEVICE=\"$dev\"\\nBOOTPROTO=\"static\"\\nDNS1=\"69.252.80.80\"\\nIPADDR=\"$ipaddr\"\\nNETMASK=\"$netmask\"\\nGATEWAY=\"$gateway\"\\nIPV6INIT=\"no\"\\nMTU=\"$mtu\"\\nONBOOT=\"yes\"\\nBONDING_OPTS=\"miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4\"' >> /etc/sysconfig/network-scripts/ifcfg-$dev";
				}
				else {
					$string = "echo -e 'DEVICE=\"$dev\"\\nBOOTPROTO=\"static\"\\nDNS1=\"69.252.80.80\"\\nIPADDR=\"$ipaddr\"\\nNETMASK=\"$netmask\"\\nGATEWAY=\"$gateway\"\\nIPV6ADDR=\"$ip6_address\"\\nIPV6INIT=\"yes\"\\nIPV6_DEFAULTGW=\"$ip6_gateway\"\\nNETWORKING_IPV6=\"yes\"\\nMTU=\"$mtu\"\\nONBOOT=\"yes\"\\nBONDING_OPTS=\"miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4\"' >> /etc/sysconfig/network-scripts/ifcfg-$dev";
				}
				my $setupslavestring = "perl /var/tmp/scripts/detect10ginterfaces.pl";
				$line = $string . "\n" . $setupslavestring . "\n" . $line;
			}
		}

		print OUT $line;
	}

	close OUT;

	my $data = `$cmd`;
	$self->render( data => $data );
}

1;
