package TrafficOps;
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

use Mojo::Base 'Mojolicious';
use Mojo::Base 'Mojolicious::Controller';
use Mojo::Base 'Mojolicious::Plugin::Config';

use base 'DBIx::Class::Core';
use Schema;
use Data::Dumper;
use Digest::SHA1 qw(sha1_hex);
use JSON;
use Cwd;

use Mojolicious::Plugins;
use Mojolicious::Plugin::Authentication;
use Mojolicious::Plugin::AccessLog;
use Mojolicious::Plugin::FormFields;
use Mojolicious::Plugin::Mail;
use Mojolicious::Static;
use Net::LDAP;
use Data::GUID;
use File::Stat qw/:stat/;
use User::pwent;
use POSIX qw(strftime);
use Utils::JsonConfig;
use MojoX::Log::Log4perl;
use File::Find;
use File::Basename;
use Env qw(PERL5LIB);
use Utils::Helper::Datasource;
use File::Path qw(make_path);
Utils::Helper::Datasource->load_extensions;

use constant SESSION_TIMEOUT => 14400;
my $logging_root_dir;
my $app_root_dir;
my $mode;
my $config;

local $/;    #Enable 'slurp' mode

has schema => sub { return Schema->connect_to_database };
my $to_extensions_lib = $ENV{'TO_EXTENSION_LIB'};
has watch => sub { [qw(lib templates $to_extensions_lib)] };

if ( !defined $ENV{MOJO_CONFIG} ) {
	$ENV{'MOJO_CONFIG'} = 'conf/cdn.conf';
	print("Loading config from /opt/traffic_ops/app/conf/cdn.conf\n");
}
else {
	print( "MOJO_CONFIG overridden: " . $ENV{MOJO_CONFIG} . "\n" );
}

my $ldap_conf_path = 'conf/ldap.conf';
my $ldap_info      = 0;
my $host;
my $admin_dn;
my $admin_pass;
my $search_base;
if ( -e $ldap_conf_path ) {
	$ldap_info   = Utils::JsonConfig->new($ldap_conf_path);
	$host        = $ldap_info->{host};
	$admin_dn    = $ldap_info->{admin_dn};
	$admin_pass  = $ldap_info->{admin_pass};
	$search_base = $ldap_info->{search_base};
}

# This method will run once at server start
sub startup {
	my $self = shift;
	$mode = $self->mode;

	$self->setup_logging($mode);
	$self->validate_cdn_conf();
	$self->setup_mojo_plugins();
	$self->set_secret();

	$self->sessions->default_expiration(SESSION_TIMEOUT);
	my $access_control_allow_origin;
	my $portal_base_url;

	# Set/override app defaults
	$self->defaults( layout => 'jquery' );

	#Static Files
	# TODO: drichardson - use this to potentially put the generate CRConfig.json elsewhere.
	my $static = Mojolicious::Static->new;
	push @{ $static->paths }, 'public';

	# Router
	my $r = $self->routes;

	# This route needs to be at the top to kick in first.
	$r->get('/')->over( authenticated => 1 )->to( 'RascalStatus#health', namespace => 'UI' );

	if ( $mode ne 'test' ) {
		$access_control_allow_origin = $config->{'cors'}{'access_control_allow_origin'};
		if ( defined($access_control_allow_origin) ) {
			$self->app->log->info( "Allowed origins : " . $config->{'cors'}{'access_control_allow_origin'} );
		}

		$portal_base_url = $config->{'portal'}{'base_url'};
		if ( defined($portal_base_url) ) {
			$self->app->log->info( "Portal Base Url : " . $portal_base_url );
		}
	}

	if ($ldap_info) {
		print("Found $ldap_conf_path, LDAP is now enabled.\n");
	}

	$self->hook(
		before_render => sub {
			my ( $self, $args ) = @_;

			# Make sure we are rendering the exception template
			return unless my $template = $args->{template};
			return unless $template eq 'exception';

			$self->app->log->error( $self->stash(" exception ") );

			# Switch to JSON rendering if content negotiation allows it
			$args->{json} = { alerts => [ { "level" => "error", "text" => "An error occurred. Please contact your administrator." } ] }
				if $self->accepts('json');
		}
	);

	if ( defined($access_control_allow_origin) ) {

		# Coors Light header (CORS)
		$self->hook(
			before_dispatch => sub {
				my $self = shift;
				$self->res->headers->header( 'Access-Control-Allow-Origin'      => $config->{'cors'}{'access_control_allow_origin'} );
				$self->res->headers->header( 'Access-Control-Allow-Headers'     => 'Origin, X-Requested-With, Content-Type, Accept' );
				$self->res->headers->header( 'Access-Control-Allow-Methods'     => 'POST,GET,OPTIONS,PUT,DELETE' );
				$self->res->headers->header( 'Access-Control-Allow-Credentials' => 'true' );
				$self->res->headers->header( 'Cache-Control'                    => 'no-cache, no-store, max-age=0, must-revalidate' );
			}
		);
	}

	# -- For CodeBig Verification
	$r->delete(
		'/api/1.1/delete' => sub {
			my $self = shift;
			$self->app->log->error( "DELETE SUCCESS : api_key" . $self->param('api_key') );
			$self->success_message( { response => "DELETE SUCCESS" } );
		}
	);

	# -- For CodeBig Verification
	$r->put(
		'/api/1.1/put' => sub {
			my $self = shift;
			$self->app->log->error( "PUT SUCCESS : api_key" . $self->param('api_key') );
			$self->success_message( { response => "PUT SUCCESS" } );
		}
	);

	# -- Ping - health check for CodeBig
	$r->get(
		'/api/1.1/ping' => sub {
			my $self = shift;
			$self->render(
				json => {
					ping => " pong "
				}
			);
		}
	);

	# ------------------------------------------------------------------------
	# NOTE: Routes should be grouped by their controller
	# ------------------------------------------------------------------------
	# -- About
	$r->get('/help/about')->over( authenticated => 1 )->to( 'Help#about', namespace => 'UI' );
	$r->get('/help/releasenotes')->over( authenticated => 1 )->to( 'Help#releasenotes', namespace => 'UI' );

	# -- Anomaly
	$r->get('/anomaly/:host_name')->to( 'Anomaly#start', namespace => 'UI' );

	# -- BlueImpLoader
	$r->get('/blueimp_uploader')->over( authenticated => 1 )->to( 'blueimp_uploader#blueimp', namespace => 'UI' );

	# -- Cachegroup
	# deprecated - see: /api/1.1/location/:parameter_id/parameter
	# $r->get('/availablelocation/:paramid')->over( authenticated => 1 )->to( 'Cachegroup#availablelocation', namespace => 'UI' );
	$r->get('/misc')->over( authenticated => 1 )->to( 'Cachegroup#index', namespace => 'UI' );
	$r->get('/cachegroups')->over( authenticated => 1 )->to( 'Cachegroup#index', namespace => 'UI' );
	$r->get('/cachegroup/add')->over( authenticated => 1 )->to( 'Cachegroup#add', namespace => 'UI' );
	$r->post('/cachegroup/create')->over( authenticated => 1 )->to( 'Cachegroup#create', namespace => 'UI' );
	$r->get('/cachegroup/:id/delete')->over( authenticated => 1 )->to( 'Cachegroup#delete', namespace => 'UI' );

	# mode is either 'edit' or 'view'.
	$r->route('/cachegroup/:mode/:id')->via('GET')->over( authenticated => 1 )->to( 'Cachegroup#view', namespace => 'UI' );
	$r->post('/cachegroup/:id/update')->over( authenticated => 1 )->to( 'Cachegroup#update', namespace => 'UI' );

	# -- Cdn
	$r->post('/login')->to( 'Cdn#login',         namespace => 'UI' );
	$r->get('/logout')->to( 'Cdn#logoutclicked', namespace => 'UI' );
	$r->get('/loginpage')->to( 'Cdn#loginpage', namespace => 'UI' );
	$r->get('/')->to( 'Cdn#loginpage', namespace => 'UI' );

	# Cdn - Special JSON format for datatables widget
	$r->get('/aadata/:table')->over( authenticated => 1 )->to( 'Cdn#aadata', namespace => 'UI' );
	$r->get('/aadata/:table/:filter/:value')->over( authenticated => 1 )->to( 'Cdn#aadata', namespace => 'UI' );

	# -- Changelog
	$r->get('/log')->over( authenticated => 1 )->to( 'ChangeLog#changelog', namespace => 'UI' );
	$r->post('/create/log')->over( authenticated => 1 )->to( 'ChangeLog#createlog',   namespace => 'UI' );
	$r->get('/newlogcount')->over( authenticated => 1 )->to( 'ChangeLog#newlogcount', namespace => 'UI' );

	# -- Configuredrac - Configure Dell DRAC settings (RAID, BIOS, etc)
	$r->post('/configuredrac')->over( authenticated => 1 )->to( 'Dell#configuredrac', namespace => 'UI' );

	# -- Configfiles
	$r->route('/genfiles/:mode/:id/#filename')->via('GET')->over( authenticated => 1 )->to( 'ConfigFiles#genfiles', namespace => 'UI' );
	$r->route('/genfiles/:mode/byprofile/:profile/CRConfig.xml')->via('GET')->over( authenticated => 1 )
		->to( 'ConfigFiles#genfiles_crconfig_profile', namespace => 'UI' );
	$r->route('/genfiles/:mode/bycdnname/:cdnname/CRConfig.xml')->via('GET')->over( authenticated => 1 )
		->to( 'ConfigFiles#genfiles_crconfig_cdnname', namespace => 'UI' );
	$r->route('/snapshot_crconfig')->via( 'GET', 'POST' )->over( authenticated => 1 )->to( 'ConfigFiles#snapshot_crconfig', namespace => 'UI' );
	$r->post('/upload_ccr_compare')->over( authenticated => 1 )->to( 'ConfigFiles#diff_ccr_xml_file', namespace => 'UI' );

	# -- Asn
	$r->get('/asns')->over( authenticated => 1 )->to( 'Asn#index', namespace => 'UI' );
	$r->get('/asns/add')->over( authenticated => 1 )->to( 'Asn#add', namespace => 'UI' );
	$r->post('/asns/create')->over( authenticated => 1 )->to( 'Asn#create', namespace => 'UI' );
	$r->get('/asns/:id/delete')->over( authenticated => 1 )->to( 'Asn#delete', namespace => 'UI' );
	$r->post('/asns/:id/update')->over( authenticated => 1 )->to( 'Asn#update', namespace => 'UI' );
	$r->route('/asns/:id/:mode')->via('GET')->over( authenticated => 1 )->to( 'Asn#view', namespace => 'UI' );

	# -- CDNs
	$r->get('/cdns/:cdn_name/dnsseckeys/add')->over( authenticated => 1 )->to( 'DnssecKeys#add', namespace => 'UI' );
	$r->get('/cdns/:cdn_name/dnsseckeys/addksk')->over( authenticated => 1 )->to( 'DnssecKeys#addksk', namespace => 'UI' );
	$r->post('/cdns/dnsseckeys/create')->over( authenticated => 1 )->to( 'DnssecKeys#create', namespace => 'UI' );
	$r->post('/cdns/dnsseckeys/genksk')->over( authenticated => 1 )->to( 'DnssecKeys#genksk', namespace => 'UI' );
	$r->get('/cdns/dnsseckeys')->to( 'DnssecKeys#index', namespace => 'UI' );
	$r->get('/cdns/:cdn_name/dnsseckeys/manage')->over( authenticated => 1 )->to( 'DnssecKeys#manage', namespace => 'UI' );
	$r->post('/cdns/dnsseckeys/activate')->over( authenticated => 1 )->to( 'DnssecKeys#activate', namespace => 'UI' );

	# -- Dell - print boxes
	$r->get('/dells')->over( authenticated => 1 )->to( 'Dell#dells', namespace => 'UI' );

	# -- Division
	$r->get('/divisions')->over( authenticated => 1 )->to( 'Division#index', namespace => 'UI' );
	$r->get('/division/add')->over( authenticated => 1 )->to( 'Division#add', namespace => 'UI' );
	$r->post('/division/create')->over( authenticated => 1 )->to( 'Division#create', namespace => 'UI' );
	$r->get('/division/:id/edit')->over( authenticated => 1 )->to( 'Division#edit', namespace => 'UI' );
	$r->post('/division/:id/update')->over( authenticated => 1 )->to( 'Division#update', namespace => 'UI' );
	$r->get('/division/:id/delete')->over( authenticated => 1 )->to( 'Division#delete', namespace => 'UI' );

	# -- DeliverysSrvice
	$r->get('/ds/add')->over( authenticated => 1 )->to( 'DeliveryService#add',  namespace => 'UI' );
	$r->get('/ds/:id')->over( authenticated => 1 )->to( 'DeliveryService#edit', namespace => 'UI' );
	$r->post('/ds/create')->over( authenticated => 1 )->to( 'DeliveryService#create', namespace => 'UI' );
	$r->get('/ds/:id/delete')->over( authenticated => 1 )->to( 'DeliveryService#delete', namespace => 'UI' );
	$r->post('/ds/:id/update')->over( authenticated => 1 )->to( 'DeliveryService#update', namespace => 'UI' );

	# -- Keys - SSL Key management
	$r->get('/ds/:id/sslkeys/add')->to( 'SslKeys#add', namespace => 'UI' );
	$r->post('/ds/sslkeys/create')->over( authenticated => 1 )->to( 'SslKeys#create', namespace => 'UI' );

	# -- Keys - SSL Key management
	$r->get('/ds/:id/urlsigkeys/add')->to( 'UrlSigKeys#add', namespace => 'UI' );

	# JvD: ded route?? # $r->get('/ds_by_id/:id')->over( authenticated => 1 )->to('DeliveryService#ds_by_id', namespace => 'UI' );
	$r->get('/healthdatadeliveryservice')->to( 'DeliveryService#readdeliveryservice', namespace => 'UI' );
	$r->get('/delivery_services')->over( authenticated => 1 )->to( 'DeliveryService#index', namespace => 'UI' );

	# -- DeliveryServiceserver
	$r->post('/dss/:id/update')->over( authenticated => 1 )->to( 'DeliveryServiceServer#assign_servers', namespace => 'UI' )
		;    # update and create are the same... ?
	$r->post('/update/cpdss/:to_server')->over( authenticated => 1 )->to( 'DeliveryServiceServer#clone_server', namespace => 'UI' );
	$r->route('/dss/:id/edit')->via('GET')->over( authenticated => 1 )->to( 'DeliveryServiceServer#edit', namespace => 'UI' );
	$r->route('/cpdssiframe/:mode/:id')->via('GET')->over( authenticated => 1 )->to( 'DeliveryServiceServer#cpdss_iframe', namespace => 'UI' );
	$r->post('/create/dsserver')->over( authenticated => 1 )->to( 'DeliveryServiceServer#create', namespace => 'UI' );

	# -- DeliveryServiceTmuser
	$r->post('/dstmuser')->over( authenticated => 1 )->to( 'DeliveryServiceTmUser#create', namespace => 'UI' );
	$r->get('/dstmuser/:ds/:tm_user_id/delete')->over( authenticated => 1 )->to( 'DeliveryServiceTmUser#delete', namespace => 'UI' );

	# -- Gendbdump - Get DB dump
	$r->get('/dbdump')->over( authenticated => 1 )->to( 'GenDbDump#dbdump', namespace => 'UI' );

	# -- Geniso - From the Tools tab:
	$r->route('/geniso')->via('GET')->over( authenticated => 1 )->to( 'GenIso#geniso', namespace => 'UI' );
	$r->route('/iso_download')->via('GET')->over( authenticated => 1 )->to( 'GenIso#iso_download', namespace => 'UI' );

	# -- Hardware
	$r->get('/hardware')->over( authenticated => 1 )->to( 'Hardware#hardware', namespace => 'UI' );
	$r->get('/hardware/:filter/:byvalue')->over( authenticated => 1 )->to( 'Hardware#hardware', namespace => 'UI' );

	# -- Health - Parameters for rascal
	$r->get('/health')->to( 'Health#healthprofile', namespace => 'UI' );
	$r->get('/healthfull')->to( 'Health#healthfull', namespace => 'UI' );
	$r->get('/health/:cdnname')->to( 'Health#rascal_config', namespace => 'UI' );

	# -- Job - These are for internal/agent job operations
	$r->post('/job/external/new')->to( 'Job#newjob', namespace => 'UI' );
	$r->get('/job/external/view/:id')->to( 'Job#read_job_by_id', namespace => 'UI' );
	$r->post('/job/external/cancel/:id')->to( 'Job#canceljob', namespace => 'UI' );
	$r->get('/job/external/result/view/:id')->to( 'Job#readresult', namespace => 'UI' );
	$r->get('/job/external/status/view/all')->to( 'Job#readstatus', namespace => 'UI' );
	$r->get('/job/agent/viewpendingjobs/:id')->over( authenticated => 1 )->to( 'Job#viewagentjob', namespace => 'UI' );
	$r->post('/job/agent/new')->over( authenticated => 1 )->to( 'Job#newagent', namespace => 'UI' );
	$r->post('/job/agent/result/new')->over( authenticated => 1 )->to( 'Job#newresult', namespace => 'UI' );
	$r->get('/job/agent/statusupdate/:id')->over( authenticated => 1 )->to( 'Job#jobstatusupdate', namespace => 'UI' );
	$r->get('/job/agent/view/all')->over( authenticated => 1 )->to( 'Job#readagent', namespace => 'UI' );
	$r->get('/job/view/all')->over( authenticated => 1 )->to( 'Job#listjob', namespace => 'UI' );
	$r->get('/job/agent/new')->over( authenticated => 1 )->to( 'Job#addagent', namespace => 'UI' );
	$r->get('/job/new')->over( authenticated => 1 )->to( 'Job#addjob', namespace => 'UI' );
	$r->get('/jobs')->over( authenticated => 1 )->to( 'Job#jobs', namespace => 'UI' );

	$r->get('/hardware/:filter/:byvalue')->over( authenticated => 1 )->to( 'Hardware#hardware', namespace => 'UI' );
	$r->get('/custom_charts')->over( authenticated => 1 )->to( 'CustomCharts#custom', namespace => 'UI' );
	$r->get('/custom_charts_single')->over( authenticated => 1 )->to( 'CustomCharts#custom_single_chart', namespace => 'UI' );
	$r->get('/custom_charts_single/cache/#cdn/#cdn_location/:cache/:stat')->over( authenticated => 1 )
		->to( 'CustomCharts#custom_single_chart', namespace => 'UI' );
	$r->get('/custom_charts_single/ds/#cdn/#cdn_location/:ds/:stat')->over( authenticated => 1 )
		->to( 'CustomCharts#custom_single_chart', namespace => 'UI' );
	$r->get('/uploadservercsv')->over( authenticated => 1 )->to( 'UploadServerCsv#uploadservercsv', namespace => 'UI' );
	$r->get('/generic_uploader')->over( authenticated => 1 )->to( 'GenericUploader#generic', namespace => 'UI' );
	$r->post('/upload_handler')->over( authenticated => 1 )->to( 'UploadHandler#upload', namespace => 'UI' );
	$r->post('/uploadhandlercsv')->over( authenticated => 1 )->to( 'UploadHandlerCsv#upload', namespace => 'UI' );

	# -- Cachegroupparameter
	$r->post('/cachegroupparameter/create')->over( authenticated => 1 )->to( 'CachegroupParameter#create', namespace => 'UI' );
	$r->get('/cachegroupparameter/#cachegroup/#parameter/delete')->over( authenticated => 1 )->to( 'CachegroupParameter#delete', namespace => 'UI' );

	# -- Options
	$r->options('/')->to( 'Cdn#options', namespace => 'UI' );
	$r->options('/*')->to( 'Cdn#options', namespace => 'UI' );

	# -- Ort
	$r->route('/ort/:hostname/ort1')->via('GET')->over( authenticated => 1 )->to( 'Ort#ort1', namespace => 'UI' );
	$r->route('/ort/:hostname/packages')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_package_versions', namespace => 'UI' );
	$r->route('/ort/:hostname/package/:package')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_package_version', namespace => 'UI' );
	$r->route('/ort/:hostname/chkconfig')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_chkconfig', namespace => 'UI' );
	$r->route('/ort/:hostname/chkconfig/:package')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_package_chkconfig', namespace => 'UI' );

	# -- Parameter
	$r->post('/parameter/create')->over( authenticated => 1 )->to( 'Parameter#create', namespace => 'UI' );
	$r->get('/parameter/:id/delete')->over( authenticated => 1 )->to( 'Parameter#delete', namespace => 'UI' );
	$r->post('/parameter/:id/update')->over( authenticated => 1 )->to( 'Parameter#update', namespace => 'UI' );
	$r->get('/parameters')->over( authenticated => 1 )->to( 'Parameter#index', namespace => 'UI' );
	$r->get('/parameters/:filter/:byvalue')->over( authenticated => 1 )->to( 'Parameter#index', namespace => 'UI' );
	$r->get('/parameter/add')->over( authenticated => 1 )->to( 'Parameter#add', namespace => 'UI' );
	$r->route('/parameter/:id')->via('GET')->over( authenticated => 1 )->to( 'Parameter#view', namespace => 'UI' );

	# -- PhysLocation
	$r->get('/phys_locations')->over( authenticated => 1 )->to( 'PhysLocation#index', namespace => 'UI' );
	$r->post('/phys_location/create')->over( authenticated => 1 )->to( 'PhysLocation#create', namespace => 'UI' );
	$r->get('/phys_location/add')->over( authenticated => 1 )->to( 'PhysLocation#add', namespace => 'UI' );

	# mode is either 'edit' or 'view'.
	$r->route('/phys_location/:id/edit')->via('GET')->over( authenticated => 1 )->to( 'PhysLocation#edit', namespace => 'UI' );
	$r->get('/phys_location/:id/delete')->over( authenticated => 1 )->to( 'PhysLocation#delete', namespace => 'UI' );
	$r->post('/phys_location/:id/update')->over( authenticated => 1 )->to( 'PhysLocation#update', namespace => 'UI' );

	# -- Profile
	$r->get('/profile/add')->over( authenticated => 1 )->to( 'Profile#add', namespace => 'UI' );
	$r->get('/profile/edit/:id')->over( authenticated => 1 )->to( 'Profile#edit', namespace => 'UI' );
	$r->route('/profile/:id/view')->via('GET')->over( authenticated => 1 )->to( 'Profile#view', namespace => 'UI' );
	$r->route('/cmpprofile/:profile1/:profile2')->via('GET')->over( authenticated => 1 )->to( 'Profile#compareprofile', namespace => 'UI' );
	$r->route('/cmpprofile/aadata/:profile1/:profile2')->via('GET')->over( authenticated => 1 )->to( 'Profile#acompareprofile', namespace => 'UI' );
	$r->post('/profile/create')->over( authenticated => 1 )->to( 'Profile#create', namespace => 'UI' );
	$r->get('/profile/import')->over( authenticated => 1 )->to( 'Profile#import', namespace => 'UI' );
	$r->post('/profile/doImport')->over( authenticated => 1 )->to( 'Profile#doImport', namespace => 'UI' );
	$r->get('/profile/:id/delete')->over( authenticated => 1 )->to( 'Profile#delete', namespace => 'UI' );
	$r->post('/profile/:id/update')->over( authenticated => 1 )->to( 'Profile#update', namespace => 'UI' );

	# select available Profile, DS or Server
	$r->get('/availableprofile/:paramid')->over( authenticated => 1 )->to( 'Profile#availableprofile', namespace => 'UI' );
	$r->route('/profile/:id/export')->via('GET')->over( authenticated => 1 )->to( 'Profile#export', namespace => 'UI' );
	$r->get('/profiles')->over( authenticated => 1 )->to( 'Profile#index', namespace => 'UI' );

	# -- Profileparameter
	$r->post('/profileparameter/create')->over( authenticated => 1 )->to( 'ProfileParameter#create', namespace => 'UI' );
	$r->get('/profileparameter/:profile/:parameter/delete')->over( authenticated => 1 )->to( 'ProfileParameter#delete', namespace => 'UI' );

	# -- Rascalstatus
	$r->get('/edge_health')->over( authenticated => 1 )->to( 'RascalStatus#health', namespace => 'UI' );
	$r->get('/rascalstatus')->over( authenticated => 1 )->to( 'RascalStatus#health', namespace => 'UI' );

	# -- Region
	$r->get('/regions')->over( authenticated => 1 )->to( 'Region#index', namespace => 'UI' );
	$r->get('/region/add')->over( authenticated => 1 )->to( 'Region#add', namespace => 'UI' );
	$r->post('/region/create')->over( authenticated => 1 )->to( 'Region#create', namespace => 'UI' );
	$r->get('/region/:id/edit')->over( authenticated => 1 )->to( 'Region#edit', namespace => 'UI' );
	$r->post('/region/:id/update')->over( authenticated => 1 )->to( 'Region#update', namespace => 'UI' );
	$r->get('/region/:id/delete')->over( authenticated => 1 )->to( 'Region#delete', namespace => 'UI' );

	# -- Server
	$r->post('/server/:name/status/:state')->over( authenticated => 1 )->to( 'Server#rest_update_server_status', namespace => 'UI' );
	$r->get('/server/:name/status')->over( authenticated => 1 )->to( 'Server#get_server_status', namespace => 'UI' );
	$r->get('/server/:key/key')->over( authenticated => 1 )->to( 'Server#get_redis_key', namespace => 'UI' );
	$r->get('/servers')->over( authenticated => 1 )->to( 'Server#index', namespace => 'UI' );
	$r->get('/server/add')->over( authenticated => 1 )->to( 'Server#add', namespace => 'UI' );
	$r->post('/server/:id/update')->over( authenticated => 1 )->to( 'Server#update', namespace => 'UI' );
	$r->get('/server/:id/delete')->over( authenticated => 1 )->to( 'Server#delete', namespace => 'UI' );
	$r->route('/server/:id/:mode')->via('GET')->over( authenticated => 1 )->to( 'Server#view', namespace => 'UI' );
	$r->post('/server/create')->over( authenticated => 1 )->to( 'Server#create', namespace => 'UI' );
	$r->post('/server/updatestatus')->over( authenticated => 1 )->to( 'Server#updatestatus', namespace => 'UI' );

	# -- Serverstatus
	$r->get('/server_check')->to( 'server_check#server_check', namespace => 'UI' );

	# -- Staticdnsentry
	$r->route('/staticdnsentry/:id/edit')->via('GET')->over( authenticated => 1 )->to( 'StaticDnsEntry#edit', namespace => 'UI' );
	$r->post('/staticdnsentry/:dsid/update')->over( authenticated => 1 )->to( 'StaticDnsEntry#update_assignments', namespace => 'UI' );
	$r->get('/staticdnsentry/:id/delete')->over( authenticated => 1 )->to( 'StaticDnsEntry#delete', namespace => 'UI' );

	# -- Status
	$r->post('/status/create')->over( authenticated => 1 )->to( 'Status#create', namespace => 'UI' );
	$r->get('/status/delete/:id')->over( authenticated => 1 )->to( 'Status#delete', namespace => 'UI' );
	$r->post('/status/update/:id')->over( authenticated => 1 )->to( 'Status#update', namespace => 'UI' );

	# -- Tools
	$r->get('/tools')->over( authenticated => 1 )->to( 'Tools#tools', namespace => 'UI' );
	$r->get('/tools/db_dump')->over( authenticated => 1 )->to( 'Tools#db_dump', namespace => 'UI' );
	$r->get('/tools/queue_updates')->over( authenticated => 1 )->to( 'Tools#queue_updates', namespace => 'UI' );
	$r->get('/tools/snapshot_crconfig')->over( authenticated => 1 )->to( 'Tools#snapshot_crconfig', namespace => 'UI' );
	$r->get('/tools/diff_crconfig/:cdn_name')->over( authenticated => 1 )->to( 'Tools#diff_crconfig_iframe', namespace => 'UI' );
	$r->get('/tools/write_crconfig/:cdn_name')->over( authenticated => 1 )->to( 'Tools#write_crconfig', namespace => 'UI' );
	$r->get('/tools/invalidate_content/')->over( authenticated => 1 )->to( 'Tools#invalidate_content', namespace => 'UI' );

	# -- Topology - CCR Config, rewrote in json
	$r->route('/genfiles/:mode/bycdnname/:cdnname/CRConfig.json')->via('GET')->over( authenticated => 1 )->to( 'Topology#ccr_config', namespace => 'UI' );

	$r->get('/types')->over( authenticated => 1 )->to( 'Types#index', namespace => 'UI' );
	$r->route('/types/add')->via('GET')->over( authenticated => 1 )->to( 'Types#add', namespace => 'UI' );
	$r->route('/types/create')->via('POST')->over( authenticated => 1 )->to( 'Types#create', namespace => 'UI' );
	$r->route('/types/:id/update')->over( authenticated => 1 )->to( 'Types#update', namespace => 'UI' );
	$r->route('/types/:id/delete')->over( authenticated => 1 )->to( 'Types#delete', namespace => 'UI' );
	$r->route('/types/:id/:mode')->via('GET')->over( authenticated => 1 )->to( 'Types#view', namespace => 'UI' );

	# -- Update bit - Process updates - legacy stuff.
	$r->get('/update/:host_name')->over( authenticated => 1 )->to( 'Server#readupdate', namespace => 'UI' );
	$r->post('/update/:host_name')->over( authenticated => 1 )->to( 'Server#postupdate', namespace => 'UI' );
	$r->post('/postupdatequeue/:id')->over( authenticated => 1 )->to( 'Server#postupdatequeue', namespace => 'UI' );
	$r->post('/postupdatequeue/:cdn/:cachegroup')->over( authenticated => 1 )->to( 'Server#postupdatequeue', namespace => 'UI' );

	# -- User
	$r->post('/user/register/send')->over( authenticated => 1 )->name('user_register_send')->to( 'User#send_registration', namespace => 'UI' );
	$r->get('/users')->name("user_index")->over( authenticated => 1 )->to( 'User#index', namespace => 'UI' );
	$r->get('/user/:id/edit')->name("user_edit")->over( authenticated => 1 )->to( 'User#edit', namespace => 'UI' );
	$r->get('/user/add')->name('user_add')->over( authenticated => 1 )->to( 'User#add', namespace => 'UI' );
	$r->get('/user/register')->name('user_register')->to( 'User#register', namespace => 'UI' );
	$r->post('/user/:id/reset_password')->name('user_reset_password')->to( 'User#reset_password', namespace => 'UI' );
	$r->post('/user')->name('user_create')->to( 'User#create', namespace => 'UI' );
	$r->post('/user/:id')->name('user_update')->to( 'User#update', namespace => 'UI' );

	# -- Utils
	$r->get('/utils/close_fancybox')->over( authenticated => 1 )->to( 'Utils#close_fancybox', namespace => 'UI' );

	# -- Visualstatus
	$r->get('/visualstatus/:matchstring')->over( authenticated => 1 )->to( 'VisualStatus#graphs', namespace => 'UI' );
	$r->get('/dailysummary')->over( authenticated => 1 )->to( 'VisualStatus#daily_summary', namespace => 'UI' );

	# ------------------------------------------------------------------------
	# API Routes
	# ------------------------------------------------------------------------
	# -- Parameter 1.0 API
	# deprecated - see: /api/1.1/crans
	$r->get('/datacrans')->over( authenticated => 1 )->to( 'Asn#index', namespace => 'UI' );
	$r->get('/datacrans/orderby/:orderby')->over( authenticated => 1 )->to( 'Asn#index', namespace => 'UI' );

	# deprecated - see: /api/1.1/locations
	$r->get('/datalocation')->over( authenticated => 1 )->to( 'Cachegroup#read', namespace => 'UI' );

	# deprecated - see: /api/1.1/locations
	$r->get('/datalocation/orderby/:orderby')->over( authenticated => 1 )->to( 'Cachegroup#read', namespace => 'UI' );
	$r->get('/datalocationtrimmed')->over( authenticated => 1 )->to( 'Cachegroup#readlocationtrimmed', namespace => 'UI' );

	# deprecated - see: /api/1.1/locationparameters
	$r->get('/datalocationparameter')->over( authenticated => 1 )->to( 'CachegroupParameter#index', namespace => 'UI' );

	# deprecated - see: /api/1.1/logs
	$r->get('/datalog')->over( authenticated => 1 )->to( 'ChangeLog#readlog', namespace => 'UI' );
	$r->get('/datalog/:days')->over( authenticated => 1 )->to( 'ChangeLog#readlog', namespace => 'UI' );

	# deprecated - see: /api/1.1/parameters
	$r->get('/dataparameter')->over( authenticated => 1 )->to( 'Parameter#readparameter', namespace => 'UI' );
	$r->get('/dataparameter/#profile_name')->over( authenticated => 1 )->to( 'Parameter#readparameter_for_profile', namespace => 'UI' );
	$r->get('/dataparameter/orderby/:orderby')->over( authenticated => 1 )->to( 'Parameter#readparameter', namespace => 'UI' );

	# deprecated - see: /api/1.1/profiles
	$r->get('/dataprofile')->over( authenticated => 1 )->to( 'Profile#readprofile', namespace => 'UI' );
	$r->get('/dataprofile/orderby/:orderby')->over( authenticated => 1 )->to( 'Profile#readprofile', namespace => 'UI' );
	$r->get('/dataprofiletrimmed')->over( authenticated => 1 )->to( 'Profile#readprofiletrimmed', namespace => 'UI' );

	# deprecated - see: /api/1.1/hwinfo
	$r->get('/datahwinfo')->over( authenticated => 1 )->to( 'HwInfo#readhwinfo', namespace => 'UI' );
	$r->get('/datahwinfo/orderby/:orderby')->over( authenticated => 1 )->to( 'HwInfo#readhwinfo', namespace => 'UI' );

	# deprecated - see: /api/1.1/profileparameters
	$r->get('/dataprofileparameter')->over( authenticated => 1 )->to( 'ProfileParameter#read', namespace => 'UI' );
	$r->get('/dataprofileparameter/orderby/:orderby')->over( authenticated => 1 )->to( 'ProfileParameter#read', namespace => 'UI' );

	# deprecated - see: /api/1.1/deliveryserviceserver
	$r->get('/datalinks')->over( authenticated => 1 )->to( 'DataAll#data_links', namespace => 'UI' );
	$r->get('/datalinks/orderby/:orderby')->over( authenticated => 1 )->to( 'DataAll#data_links', namespace => 'UI' );

	# deprecated - see: /api/1.1/deliveryserviceserver
	$r->get('/datadeliveryserviceserver')->over( authenticated => 1 )->to( 'DeliveryServiceServer#read', namespace => 'UI' );

	# deprecated - see: /api/1.1/cdn/domains
	$r->get('/datadomains')->over( authenticated => 1 )->to( 'DataAll#data_domains', namespace => 'UI' );

	# deprecated - see: /api/1.1/user/:id/deliveryservices/available.json
	$r->get('/availableds/:id')->over( authenticated => 1 )->to( 'DataAll#availableds', namespace => 'UI' );

	# deprecated - see: /api/1.1/deliveryservices.json
	#$r->get('/datadeliveryservice')->over( authenticated => 1 )->to('DeliveryService#read', namespace => 'UI' );
	$r->get('/datadeliveryservice')->to( 'DeliveryService#read', namespace => 'UI' );
	$r->get('/datadeliveryservice/orderby/:orderby')->over( authenticated => 1 )->to( 'DeliveryService#read', namespace => 'UI' );

	# deprecated - see: /api/1.1/deliveryservices.json
	$r->get('/datastatus')->over( authenticated => 1 )->to( 'Status#index', namespace => 'UI' );
	$r->get('/datastatus/orderby/:orderby')->over( authenticated => 1 )->to( 'Status#index', namespace => 'UI' );

	# deprecated - see: /api/1.1/users.json
	$r->get('/datauser')->over( authenticated => 1 )->to( 'User#read', namespace => 'UI' );
	$r->get('/datauser/orderby/:orderby')->over( authenticated => 1 )->to( 'User#read', namespace => 'UI' );

	# deprecated - see: /api/1.1/phys_locations.json
	$r->get('/dataphys_location')->over( authenticated => 1 )->to( 'PhysLocation#readphys_location', namespace => 'UI' );
	$r->get('/dataphys_locationtrimmed')->over( authenticated => 1 )->to( 'PhysLocation#readphys_locationtrimmed', namespace => 'UI' );

	# deprecated - see: /api/1.1/regions.json
	$r->get('/dataregion')->over( authenticated => 1 )->to( 'PhysLocation#readregion', namespace => 'UI' );

	# deprecated - see: /api/1.1/roles.json
	$r->get('/datarole')->over( authenticated => 1 )->to( 'Role#read', namespace => 'UI' );
	$r->get('/datarole/orderby/:orderby')->over( authenticated => 1 )->to( 'Role#read', namespace => 'UI' );

	# deprecated - see: /api/1.1/servers.json and /api/1.1/servers/hostname/:host_name/details.json
	# WARNING: unauthenticated
	#TODO JvD over auth after we have rascal pointed over!!
	$r->get('/dataserver')->to( 'Server#index_response', namespace => 'UI' );
	$r->get('/dataserver/orderby/:orderby')->to( 'Server#index_response', namespace => 'UI' );
	$r->get('/dataserverdetail/select/:select')->over( authenticated => 1 )->to( 'Server#serverdetail', namespace => 'UI' );    # legacy route - rm me later

	# deprecated - see: /api/1.1//api/1.1/staticdnsentries.json
	$r->get('/datastaticdnsentry')->over( authenticated => 1 )->to( 'StaticDnsEntry#read', namespace => 'UI' );

	# -- Type
	# deprecated - see: /api/1.1/types.json
	$r->get('/datatype')->over( authenticated => 1 )->to( 'Types#readtype', namespace => 'UI' );
	$r->get('/datatypetrimmed')->over( authenticated => 1 )->to( 'Types#readtypetrimmed', namespace => 'UI' );
	$r->get('/datatype/orderby/:orderby')->over( authenticated => 1 )->to( 'Types#readtype', namespace => 'UI' );

	# deprecated - see: /api/1.1/servers.json and /api/1.1/servers/hostname/:host_name/details.json
	# duplicate route
	$r->get('/healthdataserver')->to( 'Server#index_response', namespace => 'UI' );

	# deprecated - see: /api/1.1/traffic_monitor/stats.json
	# $r->get('/rascalstatus/getstats')->over( authenticated => 1 )->to( 'RascalStatus#get_host_stats', namespace => 'UI' );

	# deprecated - see: /api/1.1/redis/info/#shortname
	$r->get('/redis/info/#shortname')->over( authenticated => 1 )->to( 'Redis#info', namespace => 'UI' );

	# deprecated - see: /api/1.1/redis/match/#match/start_date/:start
	$r->get('/redis/#match/:start/:end/:interval')->over( authenticated => 1 )->to( 'Redis#stats', namespace => 'UI' );

	# select * from table where id=ID;
	$r->get('/server_by_id/:id')->over( authenticated => 1 )->to( 'Server#server_by_id', namespace => 'UI' );

	#$r->get('/availableserver/:dsid')->over( authenticated => 1 )->to('Server#availableserver', namespace => 'UI' );

	# ------------------------------
	# END - 1.0
	# ------------------------------

	# ------------------------------------------------------------------------
	# START: Version 1.1
	# ------------------------------------------------------------------------

	# -- API DOCS
	$r->get( '/api/1.1/docs' => [ format => [qw(json)] ] )->to( 'ApiDocs#index', namespace => 'API' );

	# -- CACHE GROUPS - #NEW
	# NOTE: any 'trimmed' urls will potentially go away with keys= support
	# -- orderby=key&key=name (where key is the database column)
	# -- query parameter options ?orderby=key&keys=name (where key is the database column)
	$r->get( '/api/1.1/cachegroups'         => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cachegroup#index',         namespace => 'API' );
	$r->get( '/api/1.1/cachegroups/trimmed' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cachegroup#index_trimmed', namespace => 'API' );

	$r->get( '/api/1.1/cachegroup/:parameter_id/parameter' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cachegroup#by_parameter_id', namespace => 'API' );
	$r->get( '/api/1.1/cachegroupparameters' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'CachegroupParameter#index', namespace => 'API' );
	$r->get( '/api/1.1/cachegroups/:parameter_id/parameter/available' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cachegroup#available_for_parameter', namespace => 'API' );

	# -- CHANGE LOG - #NEW
	$r->get( '/api/1.1/logs'            => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ChangeLog#index',       namespace => 'API' );
	$r->get( '/api/1.1/logs/:days/days' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ChangeLog#index',       namespace => 'API' );
	$r->get( '/api/1.1/logs/newcount'   => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ChangeLog#newlogcount', namespace => 'API' );

	# -- CRANS - #NEW
	$r->get( '/api/1.1/asns' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Asn#index', namespace => 'API' );

	# -- HWINFO - #NEW
	# Supports: ?orderby=key
	$r->get('/api/1.1/hwinfo')->over( authenticated => 1 )->to( 'HwInfo#index', namespace => 'API' );

	# -- KEYS
	#ping riak server
	$r->get('/api/1.1/keys/ping')->over( authenticated => 1 )->to( 'Keys#ping_riak', namespace => 'API' );

	$r->get('/api/1.1/riak/ping')->over( authenticated => 1 )->to( 'Riak#ping', namespace => 'API' );

	$r->get( '/api/1.1/riak/bucket/#bucket/key/#key/values' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Riak#get', namespace => 'API' );

	# -- DELIVERY SERVICE
	# USED TO BE - GET /api/1.1/services
	$r->get( '/api/1.1/deliveryservices' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#delivery_services', namespace => 'API' );

	# USED TO BE - GET /api/1.1/services/:id
	$r->get( '/api/1.1/deliveryservices/:id' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#delivery_services', namespace => 'API' );

	# -- DELIVERY SERVICE: Health
	# USED TO BE - GET /api/1.1/services/:id/health
	$r->get( '/api/1.1/deliveryservices/:id/health' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#health', namespace => 'API' );

	# -- DELIVERY SERVICE: Capacity
	# USED TO BE - GET /api/1.1/services/:id/capacity
	$r->get( '/api/1.1/deliveryservices/:id/capacity' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#capacity', namespace => 'API' );

	# -- DELIVERY SERVICE: Routing
	# USED TO BE - GET /api/1.1/services/:id/routing
	$r->get( '/api/1.1/deliveryservices/:id/routing' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#routing', namespace => 'API' );

	# -- DELIVERY SERVICE: State
	# USED TO BE - GET /api/1.1/services/:id/state
	$r->get( '/api/1.1/deliveryservices/:id/state' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#state', namespace => 'API' );

	# -- DELIVERY SERVICE: Metrics
	# USED TO BE - GET /api/1.1/services/:id/summary/:stat/:start/:end/:interval/:window_start/:window_end.json
	$r->get(
		'/api/1.1/deliveryservices/:id/edge/metric_types/:metric/start_date/:start/end_date/:end/interval/:interval/window_start/:window_start/window_end/:window_end'
			=> [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'DeliveryService#get_summary', namespace => 'API' );

	## -- DELIVERY SERVICE: SSL Keys
	## Support for SSL private keys, certs, and csrs
	#gets the latest key by default unless a version query param is provided with ?version=x
	$r->get( '/api/1.1/deliveryservices/xmlId/:xmlid/sslkeys' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'SslKeys#view_by_xml_id', namespace => 'API::DeliveryService' );

	#"pristine hostname"
	$r->get( '/api/1.1/deliveryservices/hostname/#hostname/sslkeys' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'SslKeys#view_by_hostname', namespace => 'API::DeliveryService' );

	#generate new
	$r->post('/api/1.1/deliveryservices/sslkeys/generate')->over( authenticated => 1 )->to( 'SslKeys#generate', namespace => 'API::DeliveryService' );

	#add existing
	$r->post('/api/1.1/deliveryservices/sslkeys/add')->over( authenticated => 1 )->to( 'SslKeys#add', namespace => 'API::DeliveryService' );

	#deletes the latest key by default unless a version query param is provided with ?version=x
	$r->get( '/api/1.1/deliveryservices/xmlId/:xmlid/sslkeys/delete' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'SslKeys#delete', namespace => 'API::DeliveryService' );

	# -- KEYS Url Sig
	$r->post('/api/1.1/deliveryservices/xmlId/:xmlId/urlkeys/generate')->over( authenticated => 1 )
		->to( 'KeysUrlSig#generate', namespace => 'API::DeliveryService' );
	$r->get( '/api/1.1/deliveryservices/xmlId/:xmlId/urlkeys' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'KeysUrlSig#view_by_xmlid', namespace => 'API::DeliveryService' );

	# Supports ?stats=true&data=true
	# USED TO BE - GET /api/1.1/deliveryservices/:id/metrics/:type/:metric/:start/:end.json
	$r->get( '/api/1.1/deliveryservices/:id/server_types/:server_type/metric_types/:metric/start_date/:start/end_date/:end' => [ format => [qw(json)] ] )
		->over( authenticated => 1 )->to( 'DeliveryService#metrics', namespace => 'API' );

	#$r->get( '/api/1.1/deliveryservices/:id/summary/:stat/:start/:end/:interval/:window_start/:window_end' => [ format => [qw(json)] ] )
	#	->over( authenticated => 1 )->to( 'DeliveryService#get_summary', namespace => 'API' );
	# -- DELIVERY SERVICE SERVER - #NEW
	# Supports ?orderby=key
	$r->get('/api/1.1/deliveryserviceserver')->over( authenticated => 1 )->to( 'DeliveryServiceServer#index', namespace => 'API' );

	# -- EXTENSIONS
	$r->get( '/api/1.1/to_extensions' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ToExtension#index', namespace => 'API' );
	$r->post('/api/1.1/to_extensions')->over( authenticated => 1 )->to( 'ToExtension#update', namespace => 'API' );
	$r->post('/api/1.1/to_extensions/:id/delete')->over( authenticated => 1 )->to( 'ToExtension#delete', namespace => 'API' );

	# -- METRICS
	# USED TO BE - GET /api/1.1/metrics/:type/:metric/:start/:end.json
	$r->get( '/api/1.1/metrics/server_types/:server_type/metric_types/:metric/start_date/:start/end_date/:end' => [ format => [qw(json)] ] )
		->over( authenticated => 1 )->to( 'Metrics#index', namespace => 'API' );

	# -- PARAMETER #NEW
	# Supports ?orderby=key
	$r->get( '/api/1.1/parameters'               => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Parameter#index',   namespace => 'API' );
	$r->get( '/api/1.1/parameters/profile/:name' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Parameter#profile', namespace => 'API' );

	# USED TO BE - GET /api/1.1/usage/:ds/:loc/:stat/:start/:end/:interval
	$r->get(
		'/api/1.1/usage/deliveryservices/:ds/cachegroups/:name/metric_types/:metric/start_date/:start_date/end_date/:end_date/interval/:interval' =>
			[ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Usage#deliveryservice', namespace => 'API' );

	# -- PHYS_LOCATION #NEW
	# Supports ?orderby=key
	$r->get('/api/1.1/phys_locations')->over( authenticated => 1 )->to( 'PhysLocation#index', namespace => 'API' );
	$r->get('/api/1.1/phys_locations/trimmed')->over( authenticated => 1 )->to( 'PhysLocation#index_trimmed', namespace => 'API' );

	# -- PROFILES - #NEW
	# Supports ?orderby=key
	$r->get( '/api/1.1/profiles'         => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Profile#index',         namespace => 'API' );
	$r->get( '/api/1.1/profiles/trimmed' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Profile#index_trimmed', namespace => 'API' );

	# -- PROFILE PARAMETERS - #NEW
	# Supports ?orderby=key
	$r->get( '/api/1.1/profileparameters' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ProfileParameter#index', namespace => 'API' );

	# -- REGION #NEW
	# Supports ?orderby=key
	$r->get('/api/1.1/regions')->over( authenticated => 1 )->to( 'Region#index', namespace => 'API' );

	# -- ROLES #NEW
	# Supports ?orderby=key
	$r->get('/api/1.1/roles')->over( authenticated => 1 )->to( 'Role#index', namespace => 'API' );

	# -- SERVER #NEW
	$r->get( '/api/1.1/servers'         => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Server#index',   namespace => 'API' );
	$r->get( '/api/1.1/servers/summary' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Server#summary', namespace => 'API' );
	$r->get( '/api/1.1/servers/hostname/:name/details' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Server#details', namespace => 'API' );
	$r->get( '/api/1.1/servers/checks'     => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ServerCheck#read',   namespace => 'API' );
	$r->get( '/api/1.1/servercheck/aadata' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ServerCheck#aadata', namespace => 'API' );
	$r->post('/api/1.1/servercheck')->over( authenticated => 1 )->to( 'ServerCheck#update', namespace => 'API' );

	# -- STATUS #NEW
	# Supports ?orderby=key
	$r->get('/api/1.1/statuses')->over( authenticated => 1 )->to( 'Status#index', namespace => 'API' );

	# -- STATIC DNS ENTRIES #NEW
	$r->get('/api/1.1/staticdnsentries')->over( authenticated => 1 )->to( 'StaticDnsEntry#index', namespace => 'API' );

	# -- SYSTEM
	$r->get( '/api/1.1/system/info' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'System#get_info', namespace => 'API' );

	# TM Status #NEW #in use # JvD
	$r->get( '/api/1.1/traffic_monitor/stats' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'TrafficMonitor#get_host_stats', namespace => 'API' );

	# -- REDIS #NEW #DR
	$r->get( '/api/1.1/redis/stats' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Redis#get_redis_stats', namespace => 'API' );
	$r->get('/api/1.1/redis/info/:host_name')->over( authenticated => 1 )->to( 'Redis#info', namespace => 'API' );
	$r->get('/api/1.1/redis/match/#match/start_date/:start_date/end_date/:end_date/interval/:interval')->over( authenticated => 1 )
		->to( 'Redis#stats', namespace => 'API' );

	# -- RIAK #NEW
	$r->get('/api/1.1/riak/stats')->over( authenticated => 1 )->to( 'Riak#stats', namespace => 'API' );

	# -- TYPE #NEW
	# Supports ?orderby=key
	$r->get('/api/1.1/types')->over( authenticated => 1 )->to( 'Types#index', namespace => 'API' );
	$r->get('/api/1.1/types/trimmed')->over( authenticated => 1 )->to( 'Types#index_trimmed', namespace => 'API' );

	# --
	# USED TO BE - GET /api/1.1/usage/overview.json
	$r->get( '/api/1.1/cdns/usage/overview' => [ format => [qw(json)] ] )->to( 'Cdn#usage_overview', namespace => 'API' );

	# -- CDN
	# USED TO BE - Nothing, this is new
	$r->get( '/api/1.1/cdns/:name/health' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#health', namespace => 'API' );

	# USED TO BE - GET /api/1.1/health.json
	$r->get( '/api/1.1/cdns/health' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#health', namespace => 'API' );

	# USED TO BE - GET /api/1.1/capacity.json
	$r->get( '/api/1.1/cdns/capacity' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#capacity', namespace => 'API' );

	# USED TO BE - GET /api/1.1/configs/monitoring/:cdn_name
	$r->get( '/api/1.1/cdns/:name/configs/monitoring' => [ format => [qw(json)] ] )->via('GET')->over( authenticated => 1 )
		->to( 'Cdn#configs_monitoring', namespace => 'API' );

	# USED TO BE - GET /api/1.1/routing.json
	$r->get( '/api/1.1/cdns/routing' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#routing', namespace => 'API' );

	#WARNING: this is an intentionally "unauthenticated" route for the Portal Home Page.
	# USED TO BE - GET /api/1.1/metrics/g/:metric/:start/:end/s.json
	$r->get( '/api/1.1/cdns/metric_types/:metric/start_date/:start/end_date/:end' => [ format => [qw(json)] ] )->to( 'Cdn#metrics', namespace => 'API' );

	## -- CDNs: DNSSEC Keys
	## Support for DNSSEC zone signing, key signing, and private keys
	#gets the latest key by default unless a version query param is provided with ?version=x
	$r->get( '/api/1.1/cdns/name/:name/dnsseckeys' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#dnssec_keys', namespace => 'API' );

	#generate new
	$r->post('/api/1.1/cdns/dnsseckeys/generate')->over( authenticated => 1 )->to( 'Cdn#dnssec_keys_generate', namespace => 'API' );

	#delete
	$r->get( '/api/1.1/cdns/name/:name/dnsseckeys/delete' => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cdn#delete_dnssec_keys', namespace => 'API' );

	# -- CDN: Topology
	# USED TO BE - GET /api/1.1/configs/cdns
	$r->get( '/api/1.1/cdns/configs' => [ format => [qw(json)] ] )->via('GET')->over( authenticated => 1 )->to( 'Cdn#get_cdns', namespace => 'API' );

	# USED TO BE - GET /api/1.1/configs/routing/:cdn_name
	$r->get( '/api/1.1/cdns/:name/configs/routing' => [ format => [qw(json)] ] )->via('GET')->over( authenticated => 1 )
		->to( 'Cdn#configs_routing', namespace => 'API' );

	# -- CDN: domains #NEW
	$r->get( '/api/1.1/cdns/domains' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#domains', namespace => 'API' );

	# -- USAGE
	# USED TO BE - GET /api/1.1/daily/usage/:ds/:loc/:stat/:start/:end/:interval
	$r->get( '/api/1.1/cdns/peakusage/:peak_usage_type/deliveryservice/:ds/cachegroup/:name/start_date/:start/end_date/:end/interval/:interval' =>
			[ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#peakusage', namespace => 'API' );

	# -- USERS
	$r->get( '/api/1.1/users' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'User#index', namespace => 'API' );
	$r->post('/api/1.1/user/login')->to( 'User#login', namespace => 'API' );
	$r->get('/api/1.1/user/:id/deliveryservices/available')->over( authenticated => 1 )->to( 'User#get_available_deliveryservices', namespace => 'API' );
	$r->post('/api/1.1/user/login/token')->to( 'User#token_login', namespace => 'API' );
	$r->post('/api/1.1/user/logout')->over( authenticated => 1 )->to( 'Cdn#tool_logout', namespace => 'UI' );

	# TO BE REFACTORED TO /api/1.1/deliveryservices/:id/jobs/keyword/PURGE
	# USED TO BE - GET /api/1.1/user/jobs/purge.json

	# USED TO BE - POST /api/1.1/user/password/reset
	$r->post('/api/1.1/user/reset_password')->to( 'User#reset_password', namespace => 'API' );

	# USED TO BE - GET /api/1.1/user/profile.json
	$r->get( '/api/1.1/user/current' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'User#current', namespace => 'API' );

	# USED TO BE - POST /api/1.1/user/job/purge
	$r->get( '/api/1.1/user/current/jobs' => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Job#index', namespace => 'API' );
	$r->post('/api/1.1/user/current/jobs')->over( authenticated => 1 )->to( 'Job#create', namespace => 'API' );

	# USED TO BE - POST /api/1.1/user/profile.json
	$r->post('/api/1.1/user/current/update')->over( authenticated => 1 )->to( 'User#update_current', namespace => 'API' );

	# ------------------------------------------------------------------------
	# END: Version 1.1
	# ------------------------------------------------------------------------

	# ------------------------------------------------------------------------
	# API Routes 1.2
	# ------------------------------------------------------------------------
	# -- INFLUXDB
	my $api_version = "v12";

	$r->get( '/api/deliveryservices/:dsid/stats' => [ format => [ $api_version . ".json" ] ] )->over( authenticated => 1 )
		->to( 'DeliveryServiceStats#index', namespace => 'API::v12' );
	$r->get( '/api/cache/stats' => [ format => [ $api_version . ".json" ] ] )->over( authenticated => 1 )
		->to( 'CacheStats#index', namespace => 'API::v12' );

	##stats_summary
	$r->get( "/api/$api_version/stats_summary" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'StatsSummary#index', namespace => "API::$api_namespace" );
	$r->post("/api/$api_version/stats_summary/create")->over( authenticated => 1 )->to( 'StatsSummary#create', namespace => "API::$api_namespace" );

	# ------------------------------------------------------------------------
	# END: Version 1.2
	# ------------------------------------------------------------------------

	# -- CATCH ALL
	$r->get('/api/(*everything)')->to( 'Cdn#catch_all', namespace => 'API' );
	$r->post('/api/(*everything)')->to( 'Cdn#catch_all', namespace => 'API' );
	$r->put('/api/(*everything)')->to( 'Cdn#catch_all', namespace => 'API' );
	$r->delete('/api/(*everything)')->to( 'Cdn#catch_all', namespace => 'API' );

	$r->get(
		'/(*everything)' => sub {
			my $self = shift;

			if ( defined( $self->current_user() ) ) {
				$self->render( template => "not_found", status => 404 );
			}
			else {
				$self->flash( login_msg => "Unauthorized. Please log in." );
				$self->render( controller => 'cdn', action => 'loginpage', layout => undef, status => 401 );
			}
		}
	);
}

sub setup_logging {
	my $self = shift;
	my $mode = shift;

	# This check prevents startup from blowing up if no conf/log4perl.conf
	# can be found, the Mojo defaults pattern/appender will kick in.
	if ( $mode eq 'production' ) {
		$logging_root_dir = "/var/log/traffic_ops";
		$app_root_dir     = "/opt/traffic_ops/app";
	}
	else {
		my $pwd = cwd();
		$logging_root_dir = "$pwd/log";
		$app_root_dir     = ".";
		make_path(
			$logging_root_dir, {
				verbose => 1,
			}
		);
	}
	my $log4perl_conf = $app_root_dir . "/conf/$mode/log4perl.conf";
	if ( -e $log4perl_conf ) {
		$self->log( MojoX::Log::Log4perl->new($log4perl_conf) );
	}
	else {
		print( "Warning cannot locate " . $log4perl_conf . ", using defaults\n" );
		$self->log( MojoX::Log::Log4perl->new() );
	}
	print("Reading log4perl config from $log4perl_conf \n");

}

sub setup_mojo_plugins {
	my $self = shift;

	$self->helper( db => sub { $self->schema } );
	$config = $self->plugin('Config');

	$self->plugin(
		'authentication', {
			autoload_user => 1,
			load_user     => sub {
				my ( $app, $username ) = @_;

				my $user_data = $self->db->resultset('TmUser')->search( { username => $username } )->single;
				my $role      = "read-only";
				my $priv      = 10;
				my $local_user;
				if ( defined($user_data) ) {
					$role       = $user_data->role->name;
					$priv       = $user_data->role->priv_level;
					$local_user = $user_data->local_user;
				}

				return {
					'username'   => $username,
					'role'       => $role,
					'priv'       => $priv,
					'local_user' => $local_user,
				};
				return undef;
			},
			validate_user => sub {
				my ( $app, $username, $pass, $options ) = @_;

				my $logged_in_user;
				my $is_authenticated;

				# Check the Token Flow
				my $token = $options->{'token'};
				if ( defined($token) ) {
					$self->app->log->debug("Token was passed, now validating...");
					$logged_in_user = $self->check_token($token);
				}

				# Check the User/Password flow
				else {
					# Check Local User (in the database)
					( $logged_in_user, $is_authenticated ) = $self->check_local_user( $username, $pass );

					# Check LDAP if conf/ldap.conf is defined.
					if ( $ldap_info && ( !$logged_in_user || !$is_authenticated ) ) {
						$logged_in_user = $self->check_ldap_user( $username, $pass );
					}

				}
				return $logged_in_user;
			},
		}
	);

	# Custom TO Plugins
	my $mojo_plugins_dir;
	foreach my $dir (@INC) {
		$mojo_plugins_dir = sprintf( "%s/MojoPlugins", $dir );
		if ( -e $mojo_plugins_dir ) {
			last;
		}
	}
	my $plugins = Mojolicious::Plugins->new;

	my @file_list;
	find(
		sub {
			return unless -f;         #Must be a file
			return unless /\.pm$/;    #Must end with `.pl` suffix
			push @file_list, $File::Find::name;
		},
		$mojo_plugins_dir
	);

	#print join "\n", @file_list;
	foreach my $file (@file_list) {
		open my $fn, '<', $file;
		my $first_line = <$fn>;
		my ( $package_keyword, $package_name ) = ( $first_line =~ m/(package )(.*);/ );
		close $fn;

		#print("Loading:  $package_name\n");
		$plugins->load_plugin($package_name);
		$self->plugin($package_name);
	}

	my $to_email_from = $config->{'to'}{'email_from'};
	if ( defined($to_email_from) ) {

		$self->plugin(
			mail => {
				from => $to_email_from,
				type => 'text/html',
			}
		);

		if ( $mode ne 'test' ) {

			$self->app->log->info("...");
			$self->app->log->info( "Traffic Ops Email From: " . $to_email_from );
		}
	}

	$self->plugin( AccessLog => { log => "$logging_root_dir/access.log" } );

	#FormFields
	$self->plugin('FormFields');

}

sub check_token {
	my $self  = shift;
	my $token = shift;
	$self->app->log->debug( "Locating user with token : " . $token . " \n " );
	my $tm_user = $self->db->resultset('TmUser')->find( { token => $token } );
	if ( defined($tm_user) ) {
		my $token_user = $self->db->resultset('TmUser')->find( { token => $token } );
		my $username = $token_user->username;
		$self->app->log->debug( "Token matched username : " . $username . " \n " );
		return $username;
	}
	else {
		$self->app->log->debug("Failed, could not find a matching token from tm_user. \n ");
		return undef;
	}
}

sub check_ldap_user {
	my $self     = shift;
	my $username = shift;
	my $pass     = shift;
	$self->app->log->debug( "Checking LDAP user: " . $username . "\n" );

	# If user is not found in local tm_user, assume it's an LDAP username, and give RO privs.
	my $user_dn = $self->find_username_in_ldap($username);
	my $is_logged_in = &login_to_ldap( $user_dn, $pass );
	if ( defined($user_dn) && $is_logged_in ) {
		$self->app->log->info( "Successful LDAP logged in : " . $username );
		return $username;
	}
	return undef;
}

sub find_username_in_ldap {
	my $self     = shift;
	my $username = shift;
	my $dn;

	$self->app->log->debug( "Searching LDAP for: " . $username );
	my $ldap = Net::LDAP->new( $host, verify => 'none', timeout => 20 ) or die "$@ ";
	$self->app->log->debug("Binding...");
	my $mesg = $ldap->bind( $admin_dn, password => "$admin_pass" );
	$mesg->code && return undef;
	$mesg = $ldap->search( base => $search_base, filter => "(&(objectCategory=person)(objectClass=user)(sAMAccountName=$username))" );
	$mesg->code && return undef;
	my $entry = $mesg->shift_entry;

	if ($entry) {
		$dn = $entry->dn;
	}
	else {
		$self->app->log->info( "Cannot find " . $username . " in LDAP." );
		return undef;
	}
	$ldap->unbind;
	return $dn;
}

# Lookup user in database
sub check_local_user {
	my $self             = shift;
	my $username         = shift;
	my $pass             = shift;
	my $local_user       = undef;
	my $is_authenticated = 0;

	my $db_user = $self->db->resultset('TmUser')->find( { username => $username } );
	if ( defined($db_user) && defined( $db_user->local_passwd ) ) {
		$self->app->log->info( $username . " was found in the database. " );
		my $db_local_passwd         = $db_user->local_passwd;
		my $db_confirm_local_passwd = $db_user->confirm_local_passwd;
		my $hex_pw_string           = sha1_hex($pass);
		if ( $db_local_passwd eq $hex_pw_string ) {
			$local_user = $username;
			$self->app->log->debug("Password matched.");
			$is_authenticated = 1;
		}
		else {
			$self->app->log->debug("Passwords did not match.");
			$local_user = 0;
		}
	}
	else {
		$self->app->log->info( "Could not find database user : " . $username );
		$local_user = 0;
	}
	return ( $local_user, $is_authenticated );
}

sub login_to_ldap {
	my $ldap;
	my $user_dn = shift;
	my $pass    = shift;
	$ldap = Net::LDAP->new( $host, verify => 'none' ) or die "$@ ";
	my $mesg = $ldap->bind( $user_dn, password => $pass );
	if ( $mesg->code ) {
		$ldap->unbind;
		return 0;
	}
	else {
		$ldap->unbind;
		return 1;
	}
}

# Validates the conf/cdn.conf for certain criteria to
# avoid admin mistakes.
sub validate_cdn_conf {
	my $self = shift;

	my $cdn_conf = $ENV{'MOJO_CONFIG'};

	open( IN, "< $cdn_conf" ) || die("$cdn_conf $!\n");
	local $/;
	my $cdn_info = eval <IN>;
	close(IN);

	my $user;
	if ( exists( $cdn_info->{hypnotoad}{user} ) ) {
		for my $u ( $cdn_info->{hypnotoad}{user} ) {
			$u =~ s/.*?\?(.*)$/$1/;

			$user = $u;
		}
	}

	my $group;
	if ( exists( $cdn_info->{hypnotoad}{group} ) ) {
		for my $g ( $cdn_info->{hypnotoad}{group} ) {
			$g =~ s/.*?\?(.*)$/$1/;

			$group = $g;
		}
	}

	if ( exists( $cdn_info->{hypnotoad}{listen} ) ) {
		for my $listen ( @{ $cdn_info->{hypnotoad}{listen} } ) {
			$listen =~ s/.*?\?(.*)$/$1/;
			if ( $listen !~ /^#/ ) {

				for my $part ( split( /&/, $listen ) ) {
					my ( $k, $v ) = split( /=/, $part );

					if ( $k eq "cert" || $k eq "key" ) {

						my @fstats = stat($v);
						my $uid    = $fstats[4];
						if ( defined($uid) ) {

							my $gid = $fstats[5];

							my $file_owner = getpwuid($uid)->name;

							my $file_group = getgrgid($gid);
							if ( ( $file_owner !~ /$user/ ) || ( $file_group !~ /$group/ ) ) {
								print( "WARNING: " . $v . " is not owned by " . $user . ":" . $group . ".\n" );
							}
						}
					}
				}
			}
		}
	}
}

sub set_secret {
	my $self = shift;

	# Set secret / disable annoying log message
	# The following commit details the change from secret to secrets in 4.63
	# https://github.com/kraih/mojo/commit/57e5129436bf3d717a13e092dd972217938e29b5
	if ( $Mojolicious::VERSION >= 4.63 ) {
		$self->secrets( ['mONKEYDOmONKEYSEE.'] );    # for Mojolicious 4.67, Top Hat
	}
	else {
		$self->secret('MonkeydoMonkeysee.');         # for Mojolicious 3.x
	}

}

1;
