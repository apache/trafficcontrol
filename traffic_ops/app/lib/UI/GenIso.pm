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
use API::Iso;

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
		if ( $self->param('stream') ne 'yes' ) {
			$self->stash(iso_file_name => $iso_file_name);
			return $self->render('gen_iso/geniso');
		}
	}
}

sub iso_download {
	my $self = shift;
	my $params = {
		hostName => $self->param('hostname'),
		osversionDir => $self->param('osversion'),
		rootPass => $self->param('rootpass'),
		dhcp => $self->param('dhcp'),
		ipAddress => $self->param('ipaddr'),
		ipNetmask => $self->param('netmask'),
		ipGateway => $self->param('gateway'),
		ip6Address => $self->param('ip6_address'),
		ip6Gateway => $self->param('ip6_gateway'),
		interfaceName => $self->param('dev'),
		interfaceMtu => $self->param('mtu'),
		disk => $self->param('ondisk'),
		mgmtIpAddress => $self->param('mgmt_ip_address'),
		mgmtIpNetmask => $self->param('mgmt_ip_netmask'),
		mgmtIpGateway => $self->param('mgmt_ip_gateway'),
		mgmtInterface => $self->param('mgmt_interface'),
		stream => $self->param('stream')
	};
	my $dl_res = &API::Iso::generate_iso($self, $params);

	if ( $params->{stream} eq 'yes' ) {
		$self->res->headers->content_type("application/download");
		$self->res->headers->content_disposition("attachment; filename=\"$dl_res->{isoName}\"");

		return $self->render( data => $dl_res->{iso} );
	} else {
		# serverselect
		$self->flash(message => "Download ISO here");
		return $dl_res->{isoName};
	}
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
