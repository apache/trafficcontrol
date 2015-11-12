package TrafficOpsRoutes;
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

sub new {
	my $self  = {};
	my $class = shift;
	return ( bless( $self, $class ) );
}

sub define {
	my $self = shift;
	my $r    = shift;

	$self->ui_routes($r);

	my $namespace = "API";

	# 1.0 Routes
	$self->api_1_0_routes( $r, "UI" );

	# $version Routes
	my $version = "1.1";
	$self->api_routes( $r, $version, $namespace );

	# 1.2 Routes
	$version = "1.2";
	$self->api_routes( $r, $version, $namespace );

	# Traffic Stats Extension
	$self->traffic_stats_routes( $r, $version );
	$self->catch_all( $r, $namespace );
}

sub ui_routes {
	my $self      = shift;
	my $r         = shift;
	my $namespace = "UI";

	# This route needs to be at the top to kick in first.
	$r->get('/')->over( authenticated => 1 )->to( 'RascalStatus#health', namespace => $namespace );

	# ------------------------------------------------------------------------
	# NOTE: Routes should be grouped by their controller
	# ------------------------------------------------------------------------
	# -- About
	$r->get('/help/about')->over( authenticated => 1 )->to( 'Help#about', namespace => $namespace );
	$r->get('/help/releasenotes')->over( authenticated => 1 )->to( 'Help#releasenotes', namespace => $namespace );

	# -- Anomaly
	$r->get('/anomaly/:host_name')->to( 'Anomaly#start', namespace => $namespace );

	# -- BlueImpLoader
	$r->get('/blueimp_uploader')->over( authenticated => 1 )->to( 'blueimp_uploader#blueimp', namespace => $namespace );

	# -- Cachegroup
	# deprecated - see: /api/$version/location/:parameter_id/parameter
	# $r->get('/availablelocation/:paramid')->over( authenticated => 1 )->to( 'Cachegroup#availablelocation', namespace => $namespace );
	$r->get('/misc')->over( authenticated => 1 )->to( 'Cachegroup#index', namespace => $namespace );
	$r->get('/cachegroups')->over( authenticated => 1 )->to( 'Cachegroup#index', namespace => $namespace );
	$r->get('/cachegroup/add')->over( authenticated => 1 )->to( 'Cachegroup#add', namespace => $namespace );
	$r->post('/cachegroup/create')->over( authenticated => 1 )->to( 'Cachegroup#create', namespace => $namespace );
	$r->get('/cachegroup/:id/delete')->over( authenticated => 1 )->to( 'Cachegroup#delete', namespace => $namespace );

	# mode is either 'edit' or 'view'.
	$r->route('/cachegroup/:mode/:id')->via('GET')->over( authenticated => 1 )->to( 'Cachegroup#view', namespace => $namespace );
	$r->post('/cachegroup/:id/update')->over( authenticated => 1 )->to( 'Cachegroup#update', namespace => $namespace );

	# -- Cdn
	$r->post('/login')->to( 'Cdn#login',         namespace => $namespace );
	$r->get('/logout')->to( 'Cdn#logoutclicked', namespace => $namespace );
	$r->get('/loginpage')->to( 'Cdn#loginpage', namespace => $namespace );
	$r->get('/')->to( 'Cdn#loginpage', namespace => $namespace );

	# Cdn - Special JSON format for datatables widget
	$r->get('/aadata/:table')->over( authenticated => 1 )->to( 'Cdn#aadata', namespace => $namespace );
	$r->get('/aadata/:table/:filter/:value')->over( authenticated => 1 )->to( 'Cdn#aadata', namespace => $namespace );

	# -- Changelog
	$r->get('/log')->over( authenticated => 1 )->to( 'ChangeLog#changelog', namespace => $namespace );
	$r->post('/create/log')->over( authenticated => 1 )->to( 'ChangeLog#createlog',   namespace => $namespace );
	$r->get('/newlogcount')->over( authenticated => 1 )->to( 'ChangeLog#newlogcount', namespace => $namespace );

	# -- Configuredrac - Configure Dell DRAC settings (RAID, BIOS, etc)
	$r->post('/configuredrac')->over( authenticated => 1 )->to( 'Dell#configuredrac', namespace => $namespace );

	# -- Configfiles
	$r->route('/genfiles/:mode/:id/#filename')->via('GET')->over( authenticated => 1 )->to( 'ConfigFiles#genfiles', namespace => $namespace );
	$r->route('/genfiles/:mode/byprofile/:profile/CRConfig.xml')->via('GET')->over( authenticated => 1 )
		->to( 'ConfigFiles#genfiles_crconfig_profile', namespace => $namespace );
	$r->route('/genfiles/:mode/bycdnname/:cdnname/CRConfig.xml')->via('GET')->over( authenticated => 1 )
		->to( 'ConfigFiles#genfiles_crconfig_cdnname', namespace => $namespace );
	$r->route('/snapshot_crconfig')->via( 'GET', 'POST' )->over( authenticated => 1 )->to( 'ConfigFiles#snapshot_crconfig', namespace => $namespace );
	$r->post('/upload_ccr_compare')->over( authenticated => 1 )->to( 'ConfigFiles#diff_ccr_xml_file', namespace => $namespace );

	# -- Asn
	$r->get('/asns')->over( authenticated => 1 )->to( 'Asn#index', namespace => $namespace );
	$r->get('/asns/add')->over( authenticated => 1 )->to( 'Asn#add', namespace => $namespace );
	$r->post('/asns/create')->over( authenticated => 1 )->to( 'Asn#create', namespace => $namespace );
	$r->get('/asns/:id/delete')->over( authenticated => 1 )->to( 'Asn#delete', namespace => $namespace );
	$r->post('/asns/:id/update')->over( authenticated => 1 )->to( 'Asn#update', namespace => $namespace );
	$r->route('/asns/:id/:mode')->via('GET')->over( authenticated => 1 )->to( 'Asn#view', namespace => $namespace );

	# -- CDNs
	$r->get('/cdns')->over( authenticated => 1 )->to( 'Cdn#index', namespace => $namespace );
	$r->get('/cdn/add')->over( authenticated => 1 )->to( 'Cdn#add', namespace => $namespace );
	$r->post('/cdn/create')->over( authenticated => 1 )->to( 'Cdn#create', namespace => $namespace );
	$r->get('/cdn/:id/delete')->over( authenticated => 1 )->to( 'Cdn#delete', namespace => $namespace );

	# mode is either 'edit' or 'view'.
	$r->route('/cdn/:mode/:id')->via('GET')->over( authenticated => 1 )->to( 'Cdn#view', namespace => $namespace );
	$r->post('/cdn/:id/update')->over( authenticated => 1 )->to( 'Cdn#update', namespace => $namespace );

	$r->get('/cdns/:cdn_name/dnsseckeys/add')->over( authenticated => 1 )->to( 'DnssecKeys#add', namespace => $namespace );
	$r->get('/cdns/:cdn_name/dnsseckeys/addksk')->over( authenticated => 1 )->to( 'DnssecKeys#addksk', namespace => $namespace );
	$r->post('/cdns/dnsseckeys/create')->over( authenticated => 1 )->to( 'DnssecKeys#create', namespace => $namespace );
	$r->post('/cdns/dnsseckeys/genksk')->over( authenticated => 1 )->to( 'DnssecKeys#genksk', namespace => $namespace );
	$r->get('/cdns/dnsseckeys')->to( 'DnssecKeys#index', namespace => $namespace );
	$r->get('/cdns/:cdn_name/dnsseckeys/manage')->over( authenticated => 1 )->to( 'DnssecKeys#manage', namespace => $namespace );
	$r->post('/cdns/dnsseckeys/activate')->over( authenticated => 1 )->to( 'DnssecKeys#activate', namespace => $namespace );

	# -- Dell - print boxes
	$r->get('/dells')->over( authenticated => 1 )->to( 'Dell#dells', namespace => $namespace );

	# -- Division
	$r->get('/divisions')->over( authenticated => 1 )->to( 'Division#index', namespace => $namespace );
	$r->get('/division/add')->over( authenticated => 1 )->to( 'Division#add', namespace => $namespace );
	$r->post('/division/create')->over( authenticated => 1 )->to( 'Division#create', namespace => $namespace );
	$r->get('/division/:id/edit')->over( authenticated => 1 )->to( 'Division#edit', namespace => $namespace );
	$r->post('/division/:id/update')->over( authenticated => 1 )->to( 'Division#update', namespace => $namespace );
	$r->get('/division/:id/delete')->over( authenticated => 1 )->to( 'Division#delete', namespace => $namespace );

	# -- DeliverysSrvice
	$r->get('/ds/add')->over( authenticated => 1 )->to( 'DeliveryService#add',  namespace => $namespace );
	$r->get('/ds/:id')->over( authenticated => 1 )->to( 'DeliveryService#edit', namespace => $namespace );
	$r->post('/ds/create')->over( authenticated => 1 )->to( 'DeliveryService#create', namespace => $namespace );
	$r->get('/ds/:id/delete')->over( authenticated => 1 )->to( 'DeliveryService#delete', namespace => $namespace );
	$r->post('/ds/:id/update')->over( authenticated => 1 )->to( 'DeliveryService#update', namespace => $namespace );

	# -- Keys - SSL Key management
	$r->get('/ds/:id/sslkeys/add')->to( 'SslKeys#add', namespace => $namespace );
	$r->post('/ds/sslkeys/create')->over( authenticated => 1 )->to( 'SslKeys#create', namespace => $namespace );

	# -- Keys - SSL Key management
	$r->get('/ds/:id/urlsigkeys/add')->to( 'UrlSigKeys#add', namespace => $namespace );

	# JvD: ded route?? # $r->get('/ds_by_id/:id')->over( authenticated => 1 )->to('DeliveryService#ds_by_id', namespace => $namespace );
	$r->get('/healthdatadeliveryservice')->to( 'DeliveryService#readdeliveryservice', namespace => $namespace );
	$r->get('/delivery_services')->over( authenticated => 1 )->to( 'DeliveryService#index', namespace => $namespace );

	# -- DeliveryServiceserver
	$r->post('/dss/:id/update')->over( authenticated => 1 )->to( 'DeliveryServiceServer#assign_servers', namespace => $namespace )
		;    # update and create are the same... ?
	$r->post('/update/cpdss/:to_server')->over( authenticated => 1 )->to( 'DeliveryServiceServer#clone_server', namespace => $namespace );
	$r->route('/dss/:id/edit')->via('GET')->over( authenticated => 1 )->to( 'DeliveryServiceServer#edit', namespace => $namespace );
	$r->route('/cpdssiframe/:mode/:id')->via('GET')->over( authenticated => 1 )->to( 'DeliveryServiceServer#cpdss_iframe', namespace => $namespace );
	$r->post('/create/dsserver')->over( authenticated => 1 )->to( 'DeliveryServiceServer#create', namespace => $namespace );

	# -- DeliveryServiceTmuser
	$r->post('/dstmuser')->over( authenticated => 1 )->to( 'DeliveryServiceTmUser#create', namespace => $namespace );
	$r->get('/dstmuser/:ds/:tm_user_id/delete')->over( authenticated => 1 )->to( 'DeliveryServiceTmUser#delete', namespace => $namespace );

	# -- Federation
	$r->get('/federation')->over( authenticated => 1 )->to( 'Federation#index', namespace => $namespace );
	$r->get('/federation/:federation_id/delete')->name("federation_delete")->over( authenticated => 1 )->to( 'Federation#delete', namespace => $namespace );
	$r->get('/federation/:federation_id/edit')->name("federation_edit")->over( authenticated => 1 )->to( 'Federation#edit', namespace => $namespace );
	$r->get('/federation/add')->name('federation_add')->over( authenticated => 1 )->to( 'Federation#add', namespace => $namespace );
	$r->post('/federation')->name('federation_create')->to( 'Federation#create', namespace => $namespace );
	$r->post('/federation/:federation_id')->name('federation_update')->to( 'Federation#update', namespace => $namespace );
	$r->get( "/federation/resolvers" => [ format => [qw(json)] ] )->to( 'Federation#resolvers', namespace => $namespace );
	$r->get( "/federation/users"     => [ format => [qw(json)] ] )->to( 'Federation#users',     namespace => $namespace );

	# -- Gendbdump - Get DB dump
	$r->get('/dbdump')->over( authenticated => 1 )->to( 'GenDbDump#dbdump', namespace => $namespace );

	# -- Geniso - From the Tools tab:
	$r->route('/geniso')->via('GET')->over( authenticated => 1 )->to( 'GenIso#geniso', namespace => $namespace );
	$r->route('/iso_download')->via('GET')->over( authenticated => 1 )->to( 'GenIso#iso_download', namespace => $namespace );

	# -- Hardware
	$r->get('/hardware')->over( authenticated => 1 )->to( 'Hardware#hardware', namespace => $namespace );
	$r->get('/hardware/:filter/:byvalue')->over( authenticated => 1 )->to( 'Hardware#hardware', namespace => $namespace );

	# -- Health - Parameters for rascal
	$r->get('/health')->to( 'Health#healthprofile', namespace => $namespace );
	$r->get('/healthfull')->to( 'Health#healthfull', namespace => $namespace );
	$r->get('/health/:cdnname')->to( 'Health#rascal_config', namespace => $namespace );

	# -- Job - These are for internal/agent job operations
	$r->post('/job/external/new')->to( 'Job#newjob', namespace => $namespace );
	$r->get('/job/external/view/:id')->to( 'Job#read_job_by_id', namespace => $namespace );
	$r->post('/job/external/cancel/:id')->to( 'Job#canceljob', namespace => $namespace );
	$r->get('/job/external/result/view/:id')->to( 'Job#readresult', namespace => $namespace );
	$r->get('/job/external/status/view/all')->to( 'Job#readstatus', namespace => $namespace );
	$r->get('/job/agent/viewpendingjobs/:id')->over( authenticated => 1 )->to( 'Job#viewagentjob', namespace => $namespace );
	$r->post('/job/agent/new')->over( authenticated => 1 )->to( 'Job#newagent', namespace => $namespace );
	$r->post('/job/agent/result/new')->over( authenticated => 1 )->to( 'Job#newresult', namespace => $namespace );
	$r->get('/job/agent/statusupdate/:id')->over( authenticated => 1 )->to( 'Job#jobstatusupdate', namespace => $namespace );
	$r->get('/job/agent/view/all')->over( authenticated => 1 )->to( 'Job#readagent', namespace => $namespace );
	$r->get('/job/view/all')->over( authenticated => 1 )->to( 'Job#listjob', namespace => $namespace );
	$r->get('/job/agent/new')->over( authenticated => 1 )->to( 'Job#addagent', namespace => $namespace );
	$r->get('/job/new')->over( authenticated => 1 )->to( 'Job#addjob', namespace => $namespace );
	$r->get('/jobs')->over( authenticated => 1 )->to( 'Job#jobs', namespace => $namespace );

	$r->get('/custom_charts')->over( authenticated => 1 )->to( 'CustomCharts#custom', namespace => $namespace );
	$r->get('/custom_charts_single')->over( authenticated => 1 )->to( 'CustomCharts#custom_single_chart', namespace => $namespace );
	$r->get('/custom_charts_single/cache/#cdn/#cdn_location/:cache/:stat')->over( authenticated => 1 )
		->to( 'CustomCharts#custom_single_chart', namespace => $namespace );
	$r->get('/custom_charts_single/ds/#cdn/#cdn_location/:ds/:stat')->over( authenticated => 1 )
		->to( 'CustomCharts#custom_single_chart', namespace => $namespace );
	$r->get('/uploadservercsv')->over( authenticated => 1 )->to( 'UploadServerCsv#uploadservercsv', namespace => $namespace );
	$r->get('/generic_uploader')->over( authenticated => 1 )->to( 'GenericUploader#generic', namespace => $namespace );
	$r->post('/upload_handler')->over( authenticated => 1 )->to( 'UploadHandler#upload', namespace => $namespace );
	$r->post('/uploadhandlercsv')->over( authenticated => 1 )->to( 'UploadHandlerCsv#upload', namespace => $namespace );

	# -- Cachegroupparameter
	$r->post('/cachegroupparameter/create')->over( authenticated => 1 )->to( 'CachegroupParameter#create', namespace => $namespace );
	$r->get('/cachegroupparameter/#cachegroup/#parameter/delete')->over( authenticated => 1 )->to( 'CachegroupParameter#delete', namespace => $namespace );

	# -- Options
	$r->options('/')->to( 'Cdn#options', namespace => $namespace );
	$r->options('/*')->to( 'Cdn#options', namespace => $namespace );

	# -- Ort
	$r->route('/ort/:hostname/ort1')->via('GET')->over( authenticated => 1 )->to( 'Ort#ort1', namespace => $namespace );
	$r->route('/ort/:hostname/packages')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_package_versions', namespace => $namespace );
	$r->route('/ort/:hostname/package/:package')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_package_version', namespace => $namespace );
	$r->route('/ort/:hostname/chkconfig')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_chkconfig', namespace => $namespace );
	$r->route('/ort/:hostname/chkconfig/:package')->via('GET')->over( authenticated => 1 )->to( 'Ort#get_package_chkconfig', namespace => $namespace );

	# -- Parameter
	$r->post('/parameter/create')->over( authenticated => 1 )->to( 'Parameter#create', namespace => $namespace );
	$r->get('/parameter/:id/delete')->over( authenticated => 1 )->to( 'Parameter#delete', namespace => $namespace );
	$r->post('/parameter/:id/update')->over( authenticated => 1 )->to( 'Parameter#update', namespace => $namespace );
	$r->get('/parameters')->over( authenticated => 1 )->to( 'Parameter#index', namespace => $namespace );
	$r->get('/parameters/:filter/:byvalue')->over( authenticated => 1 )->to( 'Parameter#index', namespace => $namespace );
	$r->get('/parameter/add')->over( authenticated => 1 )->to( 'Parameter#add', namespace => $namespace );
	$r->route('/parameter/:id')->via('GET')->over( authenticated => 1 )->to( 'Parameter#view', namespace => $namespace );

	# -- PhysLocation
	$r->get('/phys_locations')->over( authenticated => 1 )->to( 'PhysLocation#index', namespace => $namespace );
	$r->post('/phys_location/create')->over( authenticated => 1 )->to( 'PhysLocation#create', namespace => $namespace );
	$r->get('/phys_location/add')->over( authenticated => 1 )->to( 'PhysLocation#add', namespace => $namespace );

	# mode is either 'edit' or 'view'.
	$r->route('/phys_location/:id/edit')->via('GET')->over( authenticated => 1 )->to( 'PhysLocation#edit', namespace => $namespace );
	$r->get('/phys_location/:id/delete')->over( authenticated => 1 )->to( 'PhysLocation#delete', namespace => $namespace );
	$r->post('/phys_location/:id/update')->over( authenticated => 1 )->to( 'PhysLocation#update', namespace => $namespace );

	# -- Profile
	$r->get('/profile/add')->over( authenticated => 1 )->to( 'Profile#add', namespace => $namespace );
	$r->get('/profile/edit/:id')->over( authenticated => 1 )->to( 'Profile#edit', namespace => $namespace );
	$r->route('/profile/:id/view')->via('GET')->over( authenticated => 1 )->to( 'Profile#view', namespace => $namespace );
	$r->route('/cmpprofile/:profile1/:profile2')->via('GET')->over( authenticated => 1 )->to( 'Profile#compareprofile', namespace => $namespace );
	$r->route('/cmpprofile/aadata/:profile1/:profile2')->via('GET')->over( authenticated => 1 )->to( 'Profile#acompareprofile', namespace => $namespace );
	$r->post('/profile/create')->over( authenticated => 1 )->to( 'Profile#create', namespace => $namespace );
	$r->get('/profile/import')->over( authenticated => 1 )->to( 'Profile#import', namespace => $namespace );
	$r->post('/profile/doImport')->over( authenticated => 1 )->to( 'Profile#doImport', namespace => $namespace );
	$r->get('/profile/:id/delete')->over( authenticated => 1 )->to( 'Profile#delete', namespace => $namespace );
	$r->post('/profile/:id/update')->over( authenticated => 1 )->to( 'Profile#update', namespace => $namespace );

	# select available Profile, DS or Server
	$r->get('/availableprofile/:paramid')->over( authenticated => 1 )->to( 'Profile#availableprofile', namespace => $namespace );
	$r->route('/profile/:id/export')->via('GET')->over( authenticated => 1 )->to( 'Profile#export', namespace => $namespace );
	$r->get('/profiles')->over( authenticated => 1 )->to( 'Profile#index', namespace => $namespace );

	# -- Profileparameter
	$r->post('/profileparameter/create')->over( authenticated => 1 )->to( 'ProfileParameter#create', namespace => $namespace );
	$r->get('/profileparameter/:profile/:parameter/delete')->over( authenticated => 1 )->to( 'ProfileParameter#delete', namespace => $namespace );

	# -- Rascalstatus
	$r->get('/edge_health')->over( authenticated => 1 )->to( 'RascalStatus#health', namespace => $namespace );
	$r->get('/rascalstatus')->over( authenticated => 1 )->to( 'RascalStatus#health', namespace => $namespace );

	# -- Region
	$r->get('/regions')->over( authenticated => 1 )->to( 'Region#index', namespace => $namespace );
	$r->get('/region/add')->over( authenticated => 1 )->to( 'Region#add', namespace => $namespace );
	$r->post('/region/create')->over( authenticated => 1 )->to( 'Region#create', namespace => $namespace );
	$r->get('/region/:id/edit')->over( authenticated => 1 )->to( 'Region#edit', namespace => $namespace );
	$r->post('/region/:id/update')->over( authenticated => 1 )->to( 'Region#update', namespace => $namespace );
	$r->get('/region/:id/delete')->over( authenticated => 1 )->to( 'Region#delete', namespace => $namespace );

	# -- Server
	$r->post('/server/:name/status/:state')->over( authenticated => 1 )->to( 'Server#rest_update_server_status', namespace => $namespace );
	$r->get('/server/:name/status')->over( authenticated => 1 )->to( 'Server#get_server_status', namespace => $namespace );
	$r->get('/server/:key/key')->over( authenticated => 1 )->to( 'Server#get_redis_key', namespace => $namespace );
	$r->get('/servers')->over( authenticated => 1 )->to( 'Server#index', namespace => $namespace );
	$r->get('/server/add')->over( authenticated => 1 )->to( 'Server#add', namespace => $namespace );
	$r->post('/server/:id/update')->over( authenticated => 1 )->to( 'Server#update', namespace => $namespace );
	$r->get('/server/:id/delete')->over( authenticated => 1 )->to( 'Server#delete', namespace => $namespace );
	$r->route('/server/:id/:mode')->via('GET')->over( authenticated => 1 )->to( 'Server#view', namespace => $namespace );
	$r->post('/server/create')->over( authenticated => 1 )->to( 'Server#create', namespace => $namespace );
	$r->post('/server/updatestatus')->over( authenticated => 1 )->to( 'Server#updatestatus', namespace => $namespace );

	# -- Serverstatus
	$r->get('/server_check')->to( 'server_check#server_check', namespace => $namespace );

	# -- Staticdnsentry
	$r->route('/staticdnsentry/:id/edit')->via('GET')->over( authenticated => 1 )->to( 'StaticDnsEntry#edit', namespace => $namespace );
	$r->post('/staticdnsentry/:dsid/update')->over( authenticated => 1 )->to( 'StaticDnsEntry#update_assignments', namespace => $namespace );
	$r->get('/staticdnsentry/:id/delete')->over( authenticated => 1 )->to( 'StaticDnsEntry#delete', namespace => $namespace );

	# -- Status
	$r->post('/status/create')->over( authenticated => 1 )->to( 'Status#create', namespace => $namespace );
	$r->get('/status/delete/:id')->over( authenticated => 1 )->to( 'Status#delete', namespace => $namespace );
	$r->post('/status/update/:id')->over( authenticated => 1 )->to( 'Status#update', namespace => $namespace );

	# -- Tools
	$r->get('/tools')->over( authenticated => 1 )->to( 'Tools#tools', namespace => $namespace );
	$r->get('/tools/db_dump')->over( authenticated => 1 )->to( 'Tools#db_dump', namespace => $namespace );
	$r->get('/tools/queue_updates')->over( authenticated => 1 )->to( 'Tools#queue_updates', namespace => $namespace );
	$r->get('/tools/snapshot_crconfig')->over( authenticated => 1 )->to( 'Tools#snapshot_crconfig', namespace => $namespace );
	$r->get('/tools/diff_crconfig/:cdn_name')->over( authenticated => 1 )->to( 'Tools#diff_crconfig_iframe', namespace => $namespace );
	$r->get('/tools/write_crconfig/:cdn_name')->over( authenticated => 1 )->to( 'Tools#write_crconfig', namespace => $namespace );
	$r->get('/tools/invalidate_content/')->over( authenticated => 1 )->to( 'Tools#invalidate_content', namespace => $namespace );

	# -- Topology - CCR Config, rewrote in json
	$r->route('/genfiles/:mode/bycdnname/:cdnname/CRConfig.json')->via('GET')->over( authenticated => 1 )
		->to( 'Topology#ccr_config', namespace => $namespace );

	$r->get('/types')->over( authenticated => 1 )->to( 'Types#index', namespace => $namespace );
	$r->route('/types/add')->via('GET')->over( authenticated => 1 )->to( 'Types#add', namespace => $namespace );
	$r->route('/types/create')->via('POST')->over( authenticated => 1 )->to( 'Types#create', namespace => $namespace );
	$r->route('/types/:id/update')->over( authenticated => 1 )->to( 'Types#update', namespace => $namespace );
	$r->route('/types/:id/delete')->over( authenticated => 1 )->to( 'Types#delete', namespace => $namespace );
	$r->route('/types/:id/:mode')->via('GET')->over( authenticated => 1 )->to( 'Types#view', namespace => $namespace );

	# -- Update bit - Process updates - legacy stuff.
	$r->get('/update/:host_name')->over( authenticated => 1 )->to( 'Server#readupdate', namespace => $namespace );
	$r->post('/update/:host_name')->over( authenticated => 1 )->to( 'Server#postupdate', namespace => $namespace );
	$r->post('/postupdatequeue/:id')->over( authenticated => 1 )->to( 'Server#postupdatequeue', namespace => $namespace );
	$r->post('/postupdatequeue/:cdn/#cachegroup')->over( authenticated => 1 )->to( 'Server#postupdatequeue', namespace => $namespace );

	# -- User
	$r->post('/user/register/send')->over( authenticated => 1 )->name('user_register_send')->to( 'User#send_registration', namespace => $namespace );
	$r->get('/users')->name("user_index")->over( authenticated => 1 )->to( 'User#index', namespace => $namespace );
	$r->get('/user/:id/edit')->name("user_edit")->over( authenticated => 1 )->to( 'User#edit', namespace => $namespace );
	$r->get('/user/add')->name('user_add')->over( authenticated => 1 )->to( 'User#add', namespace => $namespace );
	$r->get('/user/register')->name('user_register')->to( 'User#register', namespace => $namespace );
	$r->post('/user/:id/reset_password')->name('user_reset_password')->to( 'User#reset_password', namespace => $namespace );
	$r->post('/user')->name('user_create')->to( 'User#create', namespace => $namespace );
	$r->post('/user/:id')->name('user_update')->to( 'User#update', namespace => $namespace );

	# -- Utils
	$r->get('/utils/close_fancybox')->over( authenticated => 1 )->to( 'Utils#close_fancybox', namespace => $namespace );

	# -- Visualstatus
	$r->get('/visualstatus/:matchstring')->over( authenticated => 1 )->to( 'VisualStatus#graphs', namespace => $namespace );
	$r->get('/visualstatus_redis/:matchstring')->over( authenticated => 1 )->to( 'VisualStatus#graphs_redis', namespace => $namespace );
	$r->get('/redis/#match/:start/:end/:interval')->over( authenticated => 1 )->to( 'Redis#stats', namespace => 'UI' );
	$r->get('/dailysummary')->over( authenticated => 1 )->to( 'VisualStatus#daily_summary', namespace => $namespace );

	# deprecated - see: /api/$version/servers.json and /api/1.1/servers/hostname/:host_name/details.json
	# duplicate route
	$r->get('/healthdataserver')->to( 'Server#index_response', namespace => $namespace );

	# deprecated - see: /api/$version/traffic_monitor/stats.json
	# $r->get('/rascalstatus/getstats')->over( authenticated => 1 )->to( 'RascalStatus#get_host_stats', namespace => $namespace );

	# deprecated - see: /api/$version/redis/info/#shortname
	$r->get('/redis/info/#shortname')->over( authenticated => 1 )->to( 'Redis#info', namespace => $namespace );

	# deprecated - see: /api/$version/redis/match/#match/start_date/:start
	$r->get('/redis/#match/:start_date/:end_date/:interval')->over( authenticated => 1 )->to( 'Redis#stats', namespace => $namespace );

	# select * from table where id=ID;
	$r->get('/server_by_id/:id')->over( authenticated => 1 )->to( 'Server#server_by_id', namespace => $namespace );

}

sub api_routes {
	my $self      = shift;
	my $r         = shift;
	my $version   = shift;
	my $namespace = shift;

	# -- API DOCS
	$r->get( "/api/$version/docs" => [ format => [qw(json)] ] )->to( 'ApiDocs#index', namespace => $namespace );

	# -- CACHE GROUPS - #NEW
	# NOTE: any 'trimmed' urls will potentially go away with keys= support
	# -- orderby=key&key=name (where key is the database column)
	# -- query parameter options ?orderby=key&keys=name (where key is the database column)
	$r->get( "/api/$version/cachegroups" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cachegroup#index', namespace => $namespace );
	$r->get( "/api/$version/cachegroups/trimmed" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cachegroup#index_trimmed', namespace => $namespace );

	$r->get( "/api/$version/cachegroup/:parameter_id/parameter" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cachegroup#by_parameter_id', namespace => $namespace );
	$r->get( "/api/$version/cachegroupparameters" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'CachegroupParameter#index', namespace => $namespace );
	$r->get( "/api/$version/cachegroups/:parameter_id/parameter/available" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cachegroup#available_for_parameter', namespace => $namespace );

	# -- Federation
	$r->get( "/internal/api/$version/federations" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Federation#index', namespace => $namespace );
	$r->get( "/api/$version/federations" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Federation#external_index', namespace => $namespace );
	$r->post("/api/$version/federations")->over( authenticated => 1 )->to( 'Federation#add', namespace => $namespace );
	$r->delete("/api/$version/federations")->over( authenticated => 1 )->to( 'Federation#delete', namespace => $namespace );
	$r->put("/api/$version/federations")->over( authenticated => 1 )->to( 'Federation#update', namespace => $namespace );

	# -- CDN -- #NEW
	$r->get( "/api/$version/cdns"            => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#index', namespace => $namespace );
	$r->get( "/api/$version/cdns/name/:name" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#name',  namespace => $namespace );

	# -- CHANGE LOG - #NEW
	$r->get( "/api/$version/logs"            => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ChangeLog#index', namespace => $namespace );
	$r->get( "/api/$version/logs/:days/days" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ChangeLog#index', namespace => $namespace );
	$r->get( "/api/$version/logs/newcount"   => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'ChangeLog#newlogcount', namespace => $namespace );

	# -- CRANS - #NEW
	$r->get( "/api/$version/asns" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Asn#index', namespace => $namespace );

	# -- HWINFO - #NEW
	# Supports: ?orderby=key
	$r->get("/api/$version/hwinfo")->over( authenticated => 1 )->to( 'HwInfo#index', namespace => $namespace );

	# -- KEYS
	#ping riak server
	$r->get("/api/$version/keys/ping")->over( authenticated => 1 )->to( 'Keys#ping_riak', namespace => $namespace );

	$r->get("/api/$version/riak/ping")->over( authenticated => 1 )->to( 'Riak#ping', namespace => $namespace );

	$r->get( "/api/$version/riak/bucket/#bucket/key/#key/values" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Riak#get', namespace => $namespace );

	# -- DELIVERY SERVICE
	# USED TO BE - GET /api/$version/services
	$r->get( "/api/$version/deliveryservices" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#delivery_services', namespace => $namespace );

	# USED TO BE - GET /api/$version/services/:id
	$r->get( "/api/$version/deliveryservices/:id" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#delivery_services', namespace => $namespace );

	# -- DELIVERY SERVICE: Health
	# USED TO BE - GET /api/$version/services/:id/health
	$r->get( "/api/$version/deliveryservices/:id/health" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#health', namespace => $namespace );

	# -- DELIVERY SERVICE: Capacity
	# USED TO BE - GET /api/$version/services/:id/capacity
	$r->get( "/api/$version/deliveryservices/:id/capacity" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#capacity', namespace => $namespace );

	# -- DELIVERY SERVICE: Routing
	# USED TO BE - GET /api/$version/services/:id/routing
	$r->get( "/api/$version/deliveryservices/:id/routing" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#routing', namespace => $namespace );

	# -- DELIVERY SERVICE: State
	# USED TO BE - GET /api/$version/services/:id/state
	$r->get( "/api/$version/deliveryservices/:id/state" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryService#state', namespace => $namespace );

	## -- DELIVERY SERVICE: SSL Keys
	## Support for SSL private keys, certs, and csrs
	#gets the latest key by default unless a version query param is provided with ?version=x
	$r->get( "/api/$version/deliveryservices/xmlId/:xmlid/sslkeys" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'SslKeys#view_by_xml_id', namespace => 'API::DeliveryService' );

	#"pristine hostname"
	$r->get( "/api/$version/deliveryservices/hostname/#hostname/sslkeys" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'SslKeys#view_by_hostname', namespace => 'API::DeliveryService' );

	#generate new
	$r->post("/api/$version/deliveryservices/sslkeys/generate")->over( authenticated => 1 )->to( 'SslKeys#generate', namespace => 'API::DeliveryService' );

	#add existing
	$r->post("/api/$version/deliveryservices/sslkeys/add")->over( authenticated => 1 )->to( 'SslKeys#add', namespace => 'API::DeliveryService' );

	#deletes the latest key by default unless a version query param is provided with ?version=x
	$r->get( "/api/$version/deliveryservices/xmlId/:xmlid/sslkeys/delete" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'SslKeys#delete', namespace => 'API::DeliveryService' );

	# -- KEYS Url Sig
	$r->post("/api/$version/deliveryservices/xmlId/:xmlId/urlkeys/generate")->over( authenticated => 1 )
		->to( 'KeysUrlSig#generate', namespace => 'API::DeliveryService' );
	$r->get( "/api/$version/deliveryservices/xmlId/:xmlId/urlkeys" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'KeysUrlSig#view_by_xmlid', namespace => 'API::DeliveryService' );

	#       ->over( authenticated => 1 )->to( 'DeliveryService#get_summary', namespace => $namespace );
	# -- DELIVERY SERVICE SERVER - #NEW
	# Supports ?orderby=key
	$r->get("/api/$version/deliveryserviceserver")->over( authenticated => 1 )->to( 'DeliveryServiceServer#index', namespace => $namespace );

	# -- EXTENSIONS
	$r->get( "/api/$version/to_extensions" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ToExtension#index', namespace => $namespace );
	$r->post("/api/$version/to_extensions")->over( authenticated => 1 )->to( 'ToExtension#update', namespace => $namespace );
	$r->post("/api/$version/to_extensions/:id/delete")->over( authenticated => 1 )->to( 'ToExtension#delete', namespace => $namespace );

	# -- PARAMETER #NEW
	# Supports ?orderby=key
	$r->get( "/api/$version/parameters" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Parameter#index', namespace => $namespace );
	$r->get( "/api/$version/parameters/profile/:name" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Parameter#profile', namespace => $namespace );

	# -- PHYS_LOCATION #NEW
	# Supports ?orderby=key
	$r->get("/api/$version/phys_locations")->over( authenticated => 1 )->to( 'PhysLocation#index', namespace => $namespace );
	$r->get("/api/$version/phys_locations/trimmed")->over( authenticated => 1 )->to( 'PhysLocation#index_trimmed', namespace => $namespace );

	# -- PROFILES - #NEW
	# Supports ?orderby=key
	$r->get( "/api/$version/profiles" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Profile#index', namespace => $namespace );

	$r->get( "/api/$version/profiles/trimmed" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Profile#index_trimmed', namespace => $namespace );

	# -- PROFILE PARAMETERS - #NEW
	# Supports ?orderby=key
	$r->get( "/api/$version/profileparameters" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'ProfileParameter#index', namespace => $namespace );

	# -- REGION #NEW
	# Supports ?orderby=key
	$r->get("/api/$version/regions")->over( authenticated => 1 )->to( 'Region#index', namespace => $namespace );

	# -- ROLES #NEW
	# Supports ?orderby=key
	$r->get("/api/$version/roles")->over( authenticated => 1 )->to( 'Role#index', namespace => $namespace );

	# -- SERVER #NEW
	$r->get( "/api/$version/servers"         => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Server#index',   namespace => $namespace );
	$r->get( "/api/$version/servers/summary" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Server#summary', namespace => $namespace );
	$r->get( "/api/$version/servers/hostname/:name/details" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Server#details', namespace => $namespace );
	$r->get( "/api/$version/servers/checks" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'ServerCheck#read', namespace => $namespace );
	$r->get( "/api/$version/servercheck/aadata" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'ServerCheck#aadata', namespace => $namespace );
	$r->post("/api/$version/servercheck")->over( authenticated => 1 )->to( 'ServerCheck#update', namespace => $namespace );

	# -- STATUS #NEW
	# Supports ?orderby=key
	$r->get("/api/$version/statuses")->over( authenticated => 1 )->to( 'Status#index', namespace => $namespace );

	# -- STATIC DNS ENTRIES #NEW
	$r->get("/api/$version/staticdnsentries")->over( authenticated => 1 )->to( 'StaticDnsEntry#index', namespace => $namespace );

	# -- SYSTEM
	$r->get( "/api/$version/system/info" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'System#get_info', namespace => $namespace );

	# TM Status #NEW #in use # JvD
	$r->get( "/api/$version/traffic_monitor/stats" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'TrafficMonitor#get_host_stats', namespace => $namespace );

	# -- RIAK #NEW
	$r->get("/api/$version/riak/stats")->over( authenticated => 1 )->to( 'Riak#stats', namespace => $namespace );

	# -- TYPE #NEW
	# Supports ?orderby=key
	$r->get("/api/$version/types")->over( authenticated => 1 )->to( 'Types#index', namespace => $namespace );
	$r->get("/api/$version/types/trimmed")->over( authenticated => 1 )->to( 'Types#index_trimmed', namespace => $namespace );

	# -- CDN
	# USED TO BE - Nothing, this is new
	$r->get( "/api/$version/cdns/:name/health" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#health', namespace => $namespace );

	# USED TO BE - GET /api/$version/health.json
	$r->get( "/api/$version/cdns/health" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#health', namespace => $namespace );

	# USED TO BE - GET /api/$version/capacity.json
	$r->get( "/api/$version/cdns/capacity" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#capacity', namespace => $namespace );

	# USED TO BE - GET /api/$version/routing.json
	$r->get( "/api/$version/cdns/routing" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#routing', namespace => $namespace );

	#WARNING: this is an intentionally "unauthenticated" route for the Portal Home Page.
	# USED TO BE - GET /api/$version/metrics/g/:metric/:start/:end/s.json
	$r->get( "/api/$version/cdns/metric_types/:metric_type/start_date/:start_date/end_date/:end_date" => [ format => [qw(json)] ] )
		->to( 'Cdn#metrics', namespace => $namespace );

	## -- CDNs: DNSSEC Keys
	## Support for DNSSEC zone signing, key signing, and private keys
	#gets the latest key by default unless a version query param is provided with ?version=x
	$r->get( "/api/$version/cdns/name/:name/dnsseckeys" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cdn#dnssec_keys', namespace => $namespace );

	#generate new
	$r->post("/api/$version/cdns/dnsseckeys/generate")->over( authenticated => 1 )->to( 'Cdn#dnssec_keys_generate', namespace => $namespace );

	#delete
	$r->get( "/api/$version/cdns/name/:name/dnsseckeys/delete" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'Cdn#delete_dnssec_keys', namespace => $namespace );

	#checks expiration of keys and re-generates if necessary.  Used by Cron.
	$r->get( "/internal/api/$version/cdns/dnsseckeys/refresh" => [ format => [qw(json)] ] )->to( 'Cdn#dnssec_keys_refresh', namespace => $namespace );

	# -- CDN: Topology
	# USED TO BE - GET /api/$version/configs/cdns
	$r->get( "/api/$version/cdns/configs" => [ format => [qw(json)] ] )->via('GET')->over( authenticated => 1 )
		->to( 'Cdn#get_cdns', namespace => $namespace );

	# USED TO BE - GET /api/$version/configs/routing/:cdn_name
	$r->get( "/api/$version/cdns/:name/configs/routing" => [ format => [qw(json)] ] )->via('GET')->over( authenticated => 1 )
		->to( 'Cdn#configs_routing', namespace => $namespace );

	# -- CDN: domains #NEW
	$r->get( "/api/$version/cdns/domains" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Cdn#domains', namespace => $namespace );

	# -- USERS
	$r->get( "/api/$version/users" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'User#index', namespace => $namespace );
	$r->post("/api/$version/user/login")->to( 'User#login', namespace => $namespace );
	$r->get("/api/$version/user/:id/deliveryservices/available")->over( authenticated => 1 )
		->to( 'User#get_available_deliveryservices', namespace => $namespace );
	$r->post("/api/$version/user/login/token")->to( 'User#token_login', namespace => $namespace );
	$r->post("/api/$version/user/logout")->over( authenticated => 1 )->to( 'Cdn#tool_logout', namespace => $namespace );

	# TO BE REFACTORED TO /api/$version/deliveryservices/:id/jobs/keyword/PURGE
	# USED TO BE - GET /api/$version/user/jobs/purge.json

	# USED TO BE - POST /api/$version/user/password/reset
	$r->post("/api/$version/user/reset_password")->to( 'User#reset_password', namespace => $namespace );

	# USED TO BE - GET /api/$version/user/profile.json
	$r->get( "/api/$version/user/current" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'User#current', namespace => $namespace );

	# USED TO BE - POST /api/$version/user/job/purge
	$r->get( "/api/$version/user/current/jobs" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'Job#index', namespace => $namespace );
	$r->post("/api/$version/user/current/jobs")->over( authenticated => 1 )->to( 'Job#create', namespace => $namespace );

	# USED TO BE - POST /api/$version/user/profile.json
	$r->post("/api/$version/user/current/update")->over( authenticated => 1 )->to( 'User#update_current', namespace => $namespace );

	$r->get( "/api/$version/cdns/:name/configs/monitoring" => [ format => [qw(json)] ] )->via('GET')->over( authenticated => 1 )
		->to( 'Cdn#configs_monitoring', namespace => $namespace );

	$r->get( "/api/$version/stats_summary" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'StatsSummary#index', namespace => $namespace );
	$r->post("/api/$version/stats_summary/create")->over( authenticated => 1 )->to( 'StatsSummary#create', namespace => $namespace );

	# -- Ping - health check for CodeBig
	$r->get(
		"/api/$version/ping" => sub {
			my $self = shift;
			$self->render(
				json => {
					ping => "pong"
				}
			);
		}
	);

}

sub api_1_0_routes {
	my $self      = shift;
	my $r         = shift;
	my $namespace = shift;

	# ------------------------------------------------------------------------
	# API Routes
	# ------------------------------------------------------------------------
	# -- Parameter 1.0 API
	# deprecated - see: /api/$version/crans
	$r->get('/datacrans')->over( authenticated => 1 )->to( 'Asn#index', namespace => $namespace );
	$r->get('/datacrans/orderby/:orderby')->over( authenticated => 1 )->to( 'Asn#index', namespace => $namespace );

	# deprecated - see: /api/$version/locations
	$r->get('/datalocation')->over( authenticated => 1 )->to( 'Cachegroup#read', namespace => $namespace );

	# deprecated - see: /api/$version/locations
	$r->get('/datalocation/orderby/:orderby')->over( authenticated => 1 )->to( 'Cachegroup#read', namespace => $namespace );
	$r->get('/datalocationtrimmed')->over( authenticated => 1 )->to( 'Cachegroup#readlocationtrimmed', namespace => $namespace );

	# deprecated - see: /api/$version/locationparameters
	$r->get('/datalocationparameter')->over( authenticated => 1 )->to( 'CachegroupParameter#index', namespace => $namespace );

	# deprecated - see: /api/$version/logs
	$r->get('/datalog')->over( authenticated => 1 )->to( 'ChangeLog#readlog', namespace => $namespace );
	$r->get('/datalog/:days')->over( authenticated => 1 )->to( 'ChangeLog#readlog', namespace => $namespace );

	# deprecated - see: /api/$version/parameters
	$r->get('/dataparameter')->over( authenticated => 1 )->to( 'Parameter#readparameter', namespace => $namespace );
	$r->get('/dataparameter/#profile_name')->over( authenticated => 1 )->to( 'Parameter#readparameter_for_profile', namespace => $namespace );
	$r->get('/dataparameter/orderby/:orderby')->over( authenticated => 1 )->to( 'Parameter#readparameter', namespace => $namespace );

	# deprecated - see: /api/$version/profiles
	$r->get('/dataprofile')->over( authenticated => 1 )->to( 'Profile#readprofile', namespace => $namespace );
	$r->get('/dataprofile/orderby/:orderby')->over( authenticated => 1 )->to( 'Profile#readprofile', namespace => $namespace );
	$r->get('/dataprofiletrimmed')->over( authenticated => 1 )->to( 'Profile#readprofiletrimmed', namespace => $namespace );

	# deprecated - see: /api/$version/hwinfo
	$r->get('/datahwinfo')->over( authenticated => 1 )->to( 'HwInfo#readhwinfo', namespace => $namespace );
	$r->get('/datahwinfo/orderby/:orderby')->over( authenticated => 1 )->to( 'HwInfo#readhwinfo', namespace => $namespace );

	# deprecated - see: /api/$version/profileparameters
	$r->get('/dataprofileparameter')->over( authenticated => 1 )->to( 'ProfileParameter#read', namespace => $namespace );
	$r->get('/dataprofileparameter/orderby/:orderby')->over( authenticated => 1 )->to( 'ProfileParameter#read', namespace => $namespace );

	# deprecated - see: /api/$version/deliveryserviceserver
	$r->get('/datalinks')->over( authenticated => 1 )->to( 'DataAll#data_links', namespace => $namespace );
	$r->get('/datalinks/orderby/:orderby')->over( authenticated => 1 )->to( 'DataAll#data_links', namespace => $namespace );

	# deprecated - see: /api/$version/deliveryserviceserver
	$r->get('/datadeliveryserviceserver')->over( authenticated => 1 )->to( 'DeliveryServiceServer#read', namespace => $namespace );

	# deprecated - see: /api/$version/cdn/domains
	$r->get('/datadomains')->over( authenticated => 1 )->to( 'DataAll#data_domains', namespace => $namespace );

	# deprecated - see: /api/$version/user/:id/deliveryservices/available.json
	$r->get('/availableds/:id')->over( authenticated => 1 )->to( 'DataAll#availableds', namespace => $namespace );

	# deprecated - see: /api/$version/deliveryservices.json
	#$r->get('/datadeliveryservice')->over( authenticated => 1 )->to('DeliveryService#read', namespace => $namespace );
	$r->get('/datadeliveryservice')->to( 'DeliveryService#read', namespace => $namespace );
	$r->get('/datadeliveryservice/orderby/:orderby')->over( authenticated => 1 )->to( 'DeliveryService#read', namespace => $namespace );

	# deprecated - see: /api/$version/deliveryservices.json
	$r->get('/datastatus')->over( authenticated => 1 )->to( 'Status#index', namespace => $namespace );
	$r->get('/datastatus/orderby/:orderby')->over( authenticated => 1 )->to( 'Status#index', namespace => $namespace );

	# deprecated - see: /api/$version/users.json
	$r->get('/datauser')->over( authenticated => 1 )->to( 'User#read', namespace => $namespace );
	$r->get('/datauser/orderby/:orderby')->over( authenticated => 1 )->to( 'User#read', namespace => $namespace );

	# deprecated - see: /api/$version/phys_locations.json
	$r->get('/dataphys_location')->over( authenticated => 1 )->to( 'PhysLocation#readphys_location', namespace => $namespace );
	$r->get('/dataphys_locationtrimmed')->over( authenticated => 1 )->to( 'PhysLocation#readphys_locationtrimmed', namespace => $namespace );

	# deprecated - see: /api/$version/regions.json
	$r->get('/dataregion')->over( authenticated => 1 )->to( 'PhysLocation#readregion', namespace => $namespace );

	# deprecated - see: /api/$version/roles.json
	$r->get('/datarole')->over( authenticated => 1 )->to( 'Role#read', namespace => $namespace );
	$r->get('/datarole/orderby/:orderby')->over( authenticated => 1 )->to( 'Role#read', namespace => $namespace );

	# deprecated - see: /api/$version/servers.json and /api/1.1/servers/hostname/:host_name/details.json
	# WARNING: unauthenticated
	#TODO JvD over auth after we have rascal pointed over!!
	$r->get('/dataserver')->to( 'Server#index_response', namespace => $namespace );
	$r->get('/dataserver/orderby/:orderby')->to( 'Server#index_response', namespace => $namespace );
	$r->get('/dataserverdetail/select/:select')->over( authenticated => 1 )->to( 'Server#serverdetail', namespace => $namespace )
		;    # legacy route - rm me later

	# deprecated - see: /api/$version//api/1.1/staticdnsentries.json
	$r->get('/datastaticdnsentry')->over( authenticated => 1 )->to( 'StaticDnsEntry#read', namespace => $namespace );

	# -- Type
	# deprecated - see: /api/$version/types.json
	$r->get('/datatype')->over( authenticated => 1 )->to( 'Types#readtype', namespace => $namespace );
	$r->get('/datatypetrimmed')->over( authenticated => 1 )->to( 'Types#readtypetrimmed', namespace => $namespace );
	$r->get('/datatype/orderby/:orderby')->over( authenticated => 1 )->to( 'Types#readtype', namespace => $namespace );

}

sub traffic_stats_routes {
	my $self      = shift;
	my $r         = shift;
	my $version   = shift;
	my $namespace = "Extensions::TrafficStats::API";

	$r->get( "/api/$version/cdns/usage/overview" => [ format => [qw(json)] ] )->to( 'CdnStats#get_usage_overview', namespace => $namespace );
	$r->get( "/api/$version/deliveryservice_stats" => [ format => [qw(json)] ] )->over( authenticated => 1 )
		->to( 'DeliveryServiceStats#index', namespace => $namespace );
	$r->get( "/api/$version/cache_stats" => [ format => [qw(json)] ] )->over( authenticated => 1 )->to( 'CacheStats#index', namespace => $namespace );
	$r->get( "internal/api/$version/current_bandwidth" => [ format => [qw(json)] ] )->to( 'CacheStats#current_bandwidth', namespace => $namespace );
	$r->get( "internal/api/$version/current_connections" => [ format => [qw(json)] ] )->to( 'CacheStats#current_connections', namespace => $namespace );
	$r->get( "internal/api/$version/current_capacity" => [ format => [qw(json)] ] )->to( 'CacheStats#current_capacity', namespace => $namespace );
}

sub catch_all {
	my $self      = shift;
	my $r         = shift;
	my $namespace = shift;

	# -- CATCH ALL
	$r->get('/api/(*everything)')->to( 'Cdn#catch_all', namespace => $namespace );
	$r->post('/api/(*everything)')->to( 'Cdn#catch_all', namespace => $namespace );
	$r->put('/api/(*everything)')->to( 'Cdn#catch_all', namespace => $namespace );
	$r->delete('/api/(*everything)')->to( 'Cdn#catch_all', namespace => $namespace );

	$r->get(
		'/(*everything)' => sub {
			my $self = shift;

			if ( defined( $self->current_user() ) ) {
				$self->render( template => "not_found", status => 404 );
			}
			else {
				$self->flash( login_msg => "Unauthorized . Please log in ." );
				$self->render( controller => 'cdn', action => 'loginpage', layout => undef, status => 401 );
			}
		}
	);
}

1;
