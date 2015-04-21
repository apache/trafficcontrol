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

# This is the directory we put the configuration files in for kickstart &
# scripts to process: 
my $install_cfg = "ks_scripts"; 

sub geniso {
	my $self = shift;
	print "here\n";
	&navbarpage($self);
	my %serverselect;
	my $rs_server = $self->db->resultset('Server')->search( { undef, { columns => [qw/id host_name domain_name/], orderby => "host_name" } );

	while ( my $row = $rs_server->next ) {
		my $fqdn = $row->host_name . "." . $row->domain_name;
		$serverselect{$fqdn} = $row->id;
	}

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

	# This sets up the "strength" of the hash. So far $1 works (md5). It will produce a sha256 ($5), but it's untested.
	# PROTIP: Do not put the $ in.
	my $digest = "1";

	my $dir;
	my $ksdir = $self->db->resultset('Parameter')->search( { -and => [ name => $ksfiles_parm_name, config_file => $ksfiles_configfile_name ] } )->get_column('value')->single();
	if (defined $ksdir && $ksdir ne "") {
		$dir = $ksdir . "/" . $osversion;
	} else {
		$dir = $filebasedir . "/" . $osversion;
	}
	my $cfg_dir = "$dir/$install_cfg";
	open (STUF,">$cfg_dir/state.out") or die "can't open state"; 
	print STUF "Dir== $dir\n";
	my $cmd = "mkisofs -joliet-long -input-charset utf-8 -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table -R -J -v -T $dir";
	print STUF "$cmd\n";
	close STUF;
	# This constructs the network.cfg file that gets written in the $install_cfg directory
	# in network.cfg
	# This is what we need to create:
	# IP='69.252.248.230'
	# IPV6='...'
	# NETMASK='255.255.255.252'
	# GATEWAY='69.252.248.229'
	# NAMESERVER='69.252.80.80'
	# HOSTNAME='ipcdn-cache-30.cdnlab.comcast.net' 
	# MTU='9000' 
	# BOND_DEVICE='bond0'
	# BONDOPTS='mode=802.3ad,lacp_rate=fast,xmit_hash_policy=layer3+4'
	my $network_string = "IPADDR=\"$ipaddr\"\nNETMASK=\"$netmask\"\nGATEWAY=\"$gateway\"\nBOND_DEVICE=\"$dev\"\nMTU=\"$mtu\"\nNAMESERVER=\"69.252.80.80\"\nHOSTNAME=\"$hostname\"\nNETWORKING_IPV6=\"yes\"\nIPV6ADDR=\"$ip6_address\"\nIPV6_DEFAULTGW=\"$ip6_gateway\"\nBONDING_OPTS=\"miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4\"\nDHCP=\"$dhcp\"";
	# Write out the networking config: 
	open(NF,">$cfg_dir/network.cfg") or die "Could not open network.cfg";
	print NF $network_string;
	close NF;
	my $root_pass_string;
	if ($rootpass eq "") {
		# The following password SHOULD be "Fred". YMMV, you should change this. 
		$root_pass_string = "#No password was passed in." . "\n" . ' rootpw --iscrypted $1$52LoLYxu$AcrXyEZGxiOOv4xp4E0mn/' . "\n";
		} else {
		my @chars = ("A".."Z", "a".."z",0..9) or die 'the @chars thing didn\'t work';
		my $salt;
		$salt .= $chars[rand @chars] for 1..8;
		my $kripted_pw = crypt("$rootpass","\$$digest\$$salt\$") . "\n";
		$root_pass_string = "rootpw --iscrypted $kripted_pw";
		}
	open(PWF, ">$cfg_dir/password.cfg") or die "Could not open password.cfg";
	print PWF "$root_pass_string";
	close PWF;
	
	# This wasn't necessary.
	#if ($ondisk != m/^\s*/) {
	#	$ondisk = '';
	#}

	open (DSK, ">$cfg_dir/disk.cfg") or die "Could not open disk.cfg";
	print DSK "boot_drives=\"$ondisk\"";
	close DSK;

	my $data = `$cmd`;
	$self->render( data => $data );
}
1;
