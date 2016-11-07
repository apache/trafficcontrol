package UI::GenIso;
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
use Mojolicious::Plugins;
use File::Find;
use File::Basename;
use File::Path qw(make_path);
use Data::Dumper;
use Common::ReturnCodes qw(SUCCESS ERROR);
use Mojolicious::Plugin::Config;

my $filebasedir = "/var/www/files";
my $ksfiles_parm_name = "kickstart.files.location";
my $ksfiles_configfile_name = "mkisofs";

# This is the directory we put the configuration files in for kickstart &
# scripts to process: 
my $install_cfg = "ks_scripts"; 

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
	# my $ksdir = $self->db->resultset('Parameter')->search( {  and => [ name => $ksfiles_parm_name, config_file => $ksfiles_configfile_name ] } )->get_column('value')->single();
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

	my $hostname = $self->param('hostname');
	if (defined($hostname)){
		my $iso_file_name = $self->iso_download();
		$self->stash( iso_file_name => $iso_file_name);
		#return $self->redirect_to('/geniso');
                #return $self->redirect_to("/$iso_dir/" . $iso_file_name);
		return $self->render('gen_iso/geniso');
	}
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

	# This sets up the "strength" of the hash. So far $1 works (md5). It will produce a sha256 ($5), but it's untested.
	# PROTIP: Do not put the $ in.
	my $digest = "1";

	# Read /etc/resolv.conf and get the nameservers. This is (supposedly) a hack
	# until we get a reasonable UI set up. 
	my ($nameservers, $line, $nsip) = "";
	open(RESOLV, '< /etc/resolv.conf') || die ("/etc/resolv.conf: $!");
	while ($line = <RESOLV>) {
		if ($line =~ /^nameserver /) {
			$nsip = (split(" ", $line))[1];
			$nameservers = "$nsip $nameservers";
		}
	}

	$nameservers =~ s/ /,/g;
	$nameservers =~ s/,$//;

	my $dir;
	my $ksdir = $self->db->resultset('Parameter')->search( { -and => [ name => $ksfiles_parm_name, config_file => $ksfiles_configfile_name ] } )->get_column('value')->single();

	if (defined $ksdir && $ksdir ne "") {
		$dir = $ksdir . "/" . $osversion;
	} else {
		$dir = $filebasedir . "/" . $osversion;
	}

	my $cfg_dir = "$dir/$install_cfg";
	$self->app->log->info("cfg_dir: " . $cfg_dir);

	# This constructs the network.cfg file that gets written in the $install_cfg directory
	# in network.cfg
	# This is what we need to create:
	# IP='192.168.0.2'
	# IPV6='...'
	# NETMASK='255.255.255.0'
	# GATEWAY='192.168.0.1'
	# NAMESERVER='8.8.8.8,8.8.4.4'
	# HOSTNAME='mid-cache-01.mycdn.mydomain.net' 
	# MTU='1500' 
	# BOND_DEVICE='em0'
	# BONDOPTS='mode=802.3ad,lacp_rate=fast,xmit_hash_policy=layer3+4'
	my $network_string = "IPADDR=\"$ipaddr\"\nNETMASK=\"$netmask\"\nGATEWAY=\"$gateway\"\nBOND_DEVICE=\"$dev\"\nMTU=\"$mtu\"\nNAMESERVER=\"$nameservers\"\nHOSTNAME=\"$hostname\"\nNETWORKING_IPV6=\"yes\"\nIPV6ADDR=\"$ip6_address\"\nIPV6_DEFAULTGW=\"$ip6_gateway\"\nBONDING_OPTS=\"miimon=100 mode=4 lacp_rate=fast xmit_hash_policy=layer3+4\"\nDHCP=\"$dhcp\"";
	# Write out the networking config: 
	open(NF, "> $cfg_dir/network.cfg") or die "$cfg_dir/network.cfg: $!";
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

	open(PWF, "> $cfg_dir/password.cfg") or die "$cfg_dir/password.cfg: $!";
	print PWF "$root_pass_string";
	close PWF;

	open (DSK, "> $cfg_dir/disk.cfg") or die "$cfg_dir/disk.cfg: $!";
	print DSK "boot_drives=\"$ondisk\"";
	close DSK;

	my $config_file = $ENV{'MOJO_CONFIG'};
	my $fh;
	open($fh, "<", $config_file) or die "$config_file: $!";

	my $iso_dir = "iso";
	my $config = $self->app->config;
	my $iso_root_path = $config->{'geniso'}{'iso_root_path'}; 

	my $iso_dir_path = join("/", $iso_root_path, $iso_dir);
	make_path($iso_dir_path);

	my $iso_file_name = "$hostname-$osversion.iso";
	my $iso_file_path = join("/", $iso_dir_path, $iso_file_name);

	my $cmd = "mkisofs -o $iso_file_path -joliet-long -input-charset utf-8 -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table -R -J -v -T $dir";
	my $type = "default";
	my $custom_cmd = sprintf("%s/%s", $dir, "generate");

	if (-f $custom_cmd && -x $custom_cmd) {
		$cmd = sprintf("%s %s", $custom_cmd, $iso_file_path);
		$type = "custom";
	} elsif (-f $custom_cmd && ! -x $custom_cmd) {
		$self->app->log->warn("$custom_cmd exists but is not executable; using $type ISO generation command");
	}

	$self->app->log->info("Using $type ISO generation command: " . $cmd);

	# This just writes the string we're going to use to generate the ISO. You
	# won't need it unless you're debugging stuff, but it doesn't really hurt.
	open(STUF,"> $cfg_dir/state.out") or die "$cfg_dir/state.out: $!";
	print STUF "Dir== $dir\n";
	print STUF "$cmd\n";
	close STUF;

	$self->app->log->info("Writing ISO: " . $iso_file_path);
	my $output = `$cmd 2>&1` || die("Error executing $cmd:");
	$self->app->log->info($output);

	# serverselect
	$self->flash( message => "Download ISO here" );
	return $iso_file_name;
}

sub find_conf_path {
    my $self = shift;
    my $req_conf  = shift;
    $self->app->log->info("req_conf: " . $req_conf);
    $self->app->log->info("package: " . __PACKAGE__ );
    #$self->app->log->info("INC: " . Dumper(\%INC) );
    my $p =  __PACKAGE__;
    $p =~ s/::/\//g;
    my $mod_path  = $INC{ $p . '.pm' };
    $self->app->log->info("mod_path: " . $mod_path);
    my $conf_path = join( '/', dirname( dirname( dirname($mod_path) )), 'conf', $req_conf );
    $self->app->log->info("conf_path: " . $conf_path);
    return $conf_path;
}
1;
