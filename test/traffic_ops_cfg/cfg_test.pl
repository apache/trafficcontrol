#!/usr/bin/perl
#
# Some misc subs to help with testing.
#
use warnings;
use strict;
use JSON;
use Data::Dumper;
use Term::ReadKey;
use LWP::UserAgent;
use File::Find;
use File::Spec;

#use Data::Compare;
use Test::Deep;
use Test::More;
use List::Compare;

my $config;
my $ref_db_file  = "/tmp/to_ref.mysql";
my $tmp_dir_base = "/tmp/files";
my $tmp_dir = $tmp_dir_base . "/ref";
my $CURL_OPTS;
my $cookie;

&configure( $ARGV[1] );

if ( $ARGV[0] eq "getref" ) {
	&get_ref();
}
elsif ( $ARGV[0] eq "getnew" ) {
	&get_new();
}
elsif ( $ARGV[0] eq "compare" ) {
	&do_the_compare();
}
else {
	print "Help\n";
}

if (0) {
	&build_rpms();
	&build_runner_containers();
	&load_mysql_database();
	&pg_migrate();

	# get a cookie from the reference system; cookie and CURL_OPTs are globals
#	if ( !defined( $config->{ref_to_passwd} ) ) {
#		$config->{ref_to_passwd} = &get_to_passwd( $config->{ref_to_user} );
#	}
#	my $to_login = $config->{ref_to_user} . ":" . $config->{ref_to_passwd};
#	$cookie = &get_cookie( $config->{ref_to_url}, $to_login );
#	$CURL_OPTS = "-H 'Cookie: $cookie' -w %{response_code} -k -L -s -S --connect-timeout 5 --retry 5 --retry-delay 5 --basic";
#
#	&get_files( $config->{ref_to_url} );
#	&get_crconfigs( $config->{ref_to_url} );
#
## start riak; it needs to be available at 127.0.0.1:8088
## update server set status=1 where type=41;
## insert into server (host_name, domain_name, type, ip_address, ip_gateway, ip_netmask, profile, tcp_port,interface_name, phys_location, cachegroup,status, cdn_id) values ('riak-local', 'cdnlab.comcast.net', 41, '127.0.0.1', '127.0.0.2', '255.255.255.0', 51, 8088,'eth0', 1, 1, 2, 2);
## /etc/hosts: 127.0.0.1   localhost riak-local.cdnlab.comcast.net
#
#	# start morbo like export `MOJO_MODE=test 	./bin/start.pl`
#	# get a cookie from the system we're testing; cookie and CURL_OPTs are globals
#
#	if ( !defined( $config->{to_passwd} ) ) {
#		$config->{to_passwd} = &get_to_passwd( $config->{to_user} );
#	}
#	my $to_login = $config->{to_user} . ":" . $config->{to_passwd};
#	$cookie = &get_cookie( $config->{ref_to_url}, $to_login );
#	$CURL_OPTS = "-H 'Cookie: $cookie' -w %{response_code} -k -L -s -S --connect-timeout 5 --retry 5 --retry-delay 5 --basic";
#
#	$tmp_dir = $tmp_dir_base . "/new";
#	&get_files( $config->{to_url} );
#	&get_crconfigs( $config->{to_url} );
#
#	#
#	&compare_all_files();
#
}

sub get_ref {

	# get a cookie from the reference system; cookie and CURL_OPTs are globals
	if ( !defined( $config->{ref_to_passwd} ) ) {
		$config->{ref_to_passwd} = &get_to_passwd( $config->{ref_to_user} );
	}
	my $to_login = $config->{ref_to_user} . ":" . $config->{ref_to_passwd};
	$cookie = &get_cookie( $config->{ref_to_url}, $to_login );
	$CURL_OPTS = "-H 'Cookie: $cookie' -w %{response_code} -k -L -s -S --connect-timeout 5 --retry 5 --retry-delay 5 --basic";

	$tmp_dir = $tmp_dir_base . "/ref";

	&get_files( $config->{ref_to_url} );
	&get_crconfigs( $config->{ref_to_url} );
}
#
sub get_new {
	if ( !defined( $config->{to_passwd} ) ) {
		$config->{to_passwd} = &get_to_passwd( $config->{to_user} );
	}
	my $to_login = $config->{to_user} . ":" . $config->{to_passwd};
	$cookie = &get_cookie( $config->{ref_to_url}, $to_login );
	$CURL_OPTS = "-H 'Cookie: $cookie' -w %{response_code} -k -L -s -S --connect-timeout 5 --retry 5 --retry-delay 5 --basic";

	$tmp_dir = $tmp_dir_base . "/new";

	&get_files( $config->{to_url} );     # old style
#    &get_files_new( $config->{to_url} ); # derek style
	&get_crconfigs( $config->{to_url} );
}

sub do_the_compare {
	&compare_all_files();
}
#
#&compare_files( "/tmp/files/ref/odol-atsec-uta-14/parent.config", "/tmp/files/new/odol-atsec-uta-14/parent.config" );

done_testing();
exit(0);

##################################################################################################
##################################################################################################

# compare all files in $tmp_dir_base . "/ref" to $tmp_dir_base . "/new"
sub compare_all_files {
	find( \&compare, $tmp_dir_base . "/ref" );
}

# th real work for compare_all_files
sub compare {
	my $f1 = File::Spec->rel2abs($_);
	if ( -f $f1 ) {
		( my $f2 = $f1 ) =~ s/ref/new/;
		if ( $f1 !~ /.json$/ ) {
			&compare_files( $f1, $f2 );
		}
		else {
			&compare_files( $f1, $f2, 1 );
		}
	}
}

# read a parent.config line into an object
sub get_parent_config_item {
	my $line = shift;
	my $config;
	my @parts = split( /\s+/, $line );
	foreach my $p (@parts) {
		my ( $key, $val ) = split( /=/, $p );
		$config->{$key} = $val;
	}
	return $config;
}

# parent.config is a bit different, in that the order of parents is irrelevant for urlhas and consistent_hash
sub compare_parent_dot_configs {
	my $f  = shift;
	my $d1 = shift;
	my $d2 = shift;

	if ( $d1 eq $d2 ) {
		return 0;
	}
	my @lines1 = split( /\n/, $d1 );
	my @lines2 = split( /\n/, $d2 );

	my $full_config1;
	foreach my $line (@lines1) {
		my $config_line = &get_parent_config_item($line);
		$full_config1->{ $config_line->{dest_domain} } = $config_line;
	}
	my $full_config2;
	foreach my $line (@lines2) {
		my $config_line = &get_parent_config_item($line);
		$full_config2->{ $config_line->{dest_domain} } = $config_line;
	}

	foreach my $domain ( keys %{$full_config1} ) {
		my $config1 = $full_config1->{$domain};
		my $config2 = $full_config2->{$domain};
		if ( defined( $config1->{round_robin} ) && $config1->{round_robin} =~ /hash/ ) {
			my $pstring = $config1->{parent};
			$pstring =~ s/\"//g;
			foreach my $parent ( split( /;/, $pstring ) ) {
				$config1->{parents_hash}->{$parent} = 1;
			}
			my $pstring = $config2->{parent};
			$pstring =~ s/"//g;
			foreach my $parent ( split( /;/, $pstring ) ) {
				$config2->{parents_hash}->{$parent} = 1;
			}
			$config1->{parent} = undef;
			$config2->{parent} = undef;
		}
		my $ok = cmp_deeply( $config1, $config2, "parent.config deep compare for $f:$domain" );

		if ( !$ok ) {
			print Dumper($config1);
			print Dumper($config2);
		}
	}
}

sub compare_files {
	my $f1   = shift;
	my $f2   = shift;
	my $json = shift || 0;

	open my $fh, '<', $f1 or print "$f1 is missing\n";
	my ( $d1, $d2 );
	while (<$fh>) {
		next if (/^#/);
		next if ( $f1 =~ /_xml.config$/ && /^\s*<!--.*-->\s*$/ );
		$d1 .= $_;
	}
	close($fh);
	open $fh, '<', $f2 or print "$f2 is missing\n";
	while (<$fh>) {
		next if (/^#/);
		next if ( $f2 =~ /_xml.config$/ && /^\s*<!--.*-->\s*$/ );
		$d2 .= $_;
	}
	close($fh);
	if ( !defined($d1) || !defined($d2) ) {
		return;
	}

	if ( $f1 =~ /parent.config$/ ) {
		&compare_parent_dot_configs( $f1, $d1, $d2 );
	}
	elsif ( $f1 =~ /CRConfig.json$/ || $f1 =~ /ort1$/ ) {
		my $h1 = JSON->new->allow_nonref->utf8->decode($d1);
		my $h2 = JSON->new->allow_nonref->utf8->decode($d2);
		if ( defined( $h1->{stats} ) ) {
			$h1->{stats}->{tm_user}    = $h2->{stats}->{tm_user};
			$h1->{stats}->{date}       = $h2->{stats}->{date};
			$h1->{stats}->{tm_version} = $h2->{stats}->{tm_version};
			$h1->{stats}->{tm_path}    = $h2->{stats}->{tm_path};
			$h1->{stats}->{tm_host}    = $h2->{stats}->{tm_host};
		}
		my $ok = cmp_deeply( $h1, $h2, "compare $f1" );
	}
	else {
		my $ok = cmp_deeply( $d1, $d2, "compare $f1" );
	}
}

sub copy_riak_config {

}

sub get_crconfigs {
	my $to_url = shift;

	my $to_cdn_url = $to_url . '/api/1.2/cdns.json';
	my $result     = &curl_me($to_cdn_url);
	my $cdn_json   = decode_json($result);

	my %profile_sample;
	foreach my $cdn ( @{ $cdn_json->{response} } ) {
		next unless $cdn->{name} ne "ALL";
		my $dir = $tmp_dir . '/cdn-' . $cdn->{name};
		system( 'mkdir -p ' . $dir );
		print "Generating CRConfig for " . $cdn->{name} . "\n";
		&curl_me($to_url . "/tools/write_crconfig/" . $cdn->{name});
		print "Getting CRConfig for " . $cdn->{name} . "\n";
		my $fcontents = &curl_me( $to_url . '/CRConfig-Snapshots/' . $cdn->{name} . '/CRConfig.json' );
		open( my $fh, '>', $dir . '/CRConfig.json' );
		print $fh $fcontents;
		close $fh;
	}
}

sub get_files {
	my $to_url = shift;

	my $to_server_url = $to_url . '/api/1.2/servers.json';
	my $result        = &curl_me($to_server_url);
	my $server_json   = decode_json($result);

	my %profile_sample;
	foreach my $server ( @{ $server_json->{response} } ) {
		$profile_sample{ $server->{profile} } = $server->{hostName};
	}

	foreach my $sample_server ( keys %profile_sample ) {
		my $dir = $tmp_dir . '/' . $profile_sample{$sample_server};
		system( 'mkdir -p ' . $dir );
		my $to_ort1_url = $to_url . '/ort/' . $profile_sample{$sample_server} . '/ort1';
		my $result      = &curl_me($to_ort1_url);
		open( my $fh, '>', $dir . '/ort1' );
		print $fh $result;
		close $fh;
		my $file_list_json = decode_json($result);

		foreach my $filename ( keys %{ $file_list_json->{config_files} } ) {

			next unless  $filename eq "parent.config";
			print "Getting " . $sample_server . " " . $filename . "\n";
			my $fcontents = &curl_me( $to_url . '/genfiles/view/' . $profile_sample{$sample_server} . "/" . $filename );
			open( my $fh, '>', $dir . '/' . $filename );
			print $fh $fcontents;
			close $fh;
		}
	}
}

sub get_files_new {
	my $to_url = shift;

	my $to_server_url = $to_url . '/api/1.2/servers.json';
	my $result        = &curl_me($to_server_url);
	my $server_json   = decode_json($result);

	my %profile_sample;
	foreach my $server ( @{ $server_json->{response} } ) {
		$profile_sample{ $server->{profile} } = $server->{hostName};
	}

	foreach my $sample_server ( keys %profile_sample ) {
		my $dir = $tmp_dir . '/' . $profile_sample{$sample_server};
		system( 'mkdir -p ' . $dir );
		my $to_ort1_url = $to_url . '/api/1.2/server/' . $profile_sample{$sample_server} . '/configfiles/ats.json';
		my $new_mode    = 1;
		if ( $sample_server !~ /^EDGE/ && $sample_server !~ /^MID/ && $sample_server !~ /TEAK/ ) {
			$to_ort1_url = $to_url . '/ort/' . $profile_sample{$sample_server} . '/ort1';
			$new_mode    = 0;
		}
		my $result = &curl_me($to_ort1_url);
		open( my $fh, '>', $dir . '/ats.json' );
		print $fh $result;
		close $fh;
		my $file_list_json = decode_json($result);

		foreach my $filename ( keys %{ $file_list_json->{config_files} } ) {
			my $url;
			if ( defined( $file_list_json->{config_files}->{$filename}->{scope} ) ) {
				$url = $to_url . $file_list_json->{config_files}->{$filename}->{API_URI};

				#				my $scope   = $file_list_json->{config_files}->{$filename}->{scope};
				#				my $cdn     = $file_list_json->{other}->{CDN_name};
				#				my $profile = $file_list_json->{profile}->{name};
				#				if ( $scope eq "cdn" ) {
				#					$url = $to_url . '/api/1.2/' . $scope . "/" . $cdn . "/configfiles/ats/" . $filename;
				#				}
				#				elsif ( $scope eq "profile" ) {
				#					$url = $to_url . '/api/1.2/' . $scope . "/" . $profile . "/configfiles/ats/" . $filename;
				#				}
				#				elsif ( $scope eq "server" ) {
				#					$url = $to_url . '/api/1.2/' . $scope . "/" . $profile_sample{$sample_server} . "/configfiles/ats/" . $filename;
				#				}
			}
			else {
				$url = $to_url . '/genfiles/view/' . $profile_sample{$sample_server} . "/" . $filename;
			}
			print "Getting " . $sample_server . " " . $filename . " (url " . $url . ")\n";
			my $fcontents = &curl_me($url);
			open( my $fh, '>', $dir . '/' . $filename );
			print $fh $fcontents;
			close $fh;
		}
	}
}

sub load_mysql_database {
	my $cmd  = "mysql ";
	my $args = "-h " . $config->{mysql_db_host} . " -u " . $config->{mysql_dbadmin_user} . " -p" . $config->{mysql_dbadmin_passwd};
	$cmd .= $args;
	my $bash_cmd = "echo drop database " . $config->{mysql_db_name} . " | " . $cmd;
	print $bash_cmd . "\n";
	system($bash_cmd );
	$bash_cmd = "echo create database " . $config->{mysql_db_name} . " | " . $cmd;
	print $bash_cmd . "\n";
	system($bash_cmd );
	$bash_cmd = $cmd . " " . $config->{mysql_db_name} . " < " . $ref_db_file;
	print $bash_cmd . "\n";
	system($bash_cmd );
}

# perform the migration from mysql -> postgres
sub pg_migrate {

	my $drop_cmd = "echo drop database " . $config->{pg_db_name} . " | psql postgres";
	print $drop_cmd . "\n";
	system($drop_cmd);

	my $cr_cmd = "echo create database " . $config->{pg_db_name} . " | psql postgres";
	print $cr_cmd . "\n";
	system($cr_cmd);

	my $cmd = "pgloader  --cast 'type tinyint to smallint drop typemod' --cast 'type varchar to text drop typemod'";
	$cmd .= " --cast 'type double to numeric drop typemod'";
	my $args =
		  " mysql://"
		. $config->{mysql_dbadmin_user} . ":"
		. $config->{mysql_dbadmin_passwd} . "@"
		. $config->{mysql_db_host} . ":"
		. $config->{mysql_db_port};
	$args .= "/" . $config->{mysql_db_name};
	$args .= " postgresql://" . "/" . $config->{pg_db_name};    # TODO add username / passwd if you have it.
	$cmd  .= $args;
	print $cmd . "\n";
	system($cmd);
	chdir( $config->{working_dir} . "/traffic_ops/app" ) || die "can't chdir to " . $config->{working_dir} . "/traffic_ops/app";
	$cmd = "psql " . $config->{pg_db_name} . "< db/convert_bools.sql";
	system($cmd);
	$cmd = "goose -env=" . $config->{goose_env} . " up";
	system($cmd);
}

sub get_to_passwd {
	my $user = shift;

	print "Type Reference Traffic Ops passwd for " . $user . ":";
	ReadMode('noecho');    # don't echo
	chomp( my $passwd = <STDIN> );
	ReadMode(0);           # back to normal
	print "\n";
	return $passwd;
}

# TODO JvD: finish
sub get_reference_database {

	#"https://tm.comcast.net/dbdump?filename=to-backup-ipcdn-tools-03.cdnlab.comcast.net-20170114222140.mysql

}

# TODO - need to get postgres container
sub start_runner_containers {

	#docker run --name my-traffic-vault --hostname my-traffic-vault --net cdnet --env ADMIN_PASS=riakadminsecret --env USER_PASS=marginallylesssecret
	# --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000
	# --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=cdnet --detach traffic_vault:1.6.0
	# at some point you'll have to have done `docker network create cdnet`
	my $tv_args = " --name traffic-vault --hostname traffic-vault --net cdnet --env ADMIN_PASS=riakadminsecret";
	$tv_args .= " --env USER_PASS=marginallylesssecret --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver";
	$tv_args .= " --env CERT_COMPANY=NotComcast --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot";
	$tv_args .= " --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=cdnet --detach traffic_vault:" . $config->{git_branch};

	my $to_args =
		" --name traffic-ops --hostname my-traffic-ops --net cdnet --publish 443:443 --env MYSQL_IP=my-traffic-ops-mysql --env MYSQL_PORT=3306 --env MYSQL_ROOT_PASS=secretrootpass --env MYSQL_TRAFFIC_OPS_PASS=supersecretpassword --env ADMIN_USER=superroot --env ADMIN_PASS=supersecreterpassward --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env TRAFFIC_VAULT_PASS=marginallylesssecret --env DOMAIN=cdnet --detach traffic_ops:1.5.1"

}

sub build_rpms {
	my $dir = $config->{working_dir};
	chdir($dir) || die( "Can't chdir to " . $dir );

	$ENV{'BRANCH'}  = $config->{git_branch};
	$ENV{'GITREPO'} = $config->{git_repo};

	foreach my $builder (qw/traffic_monitor_build traffic_ops_build traffic_portal_build traffic_router_build traffic_stats_build/) {
		system( "docker-compose -f infrastructure/docker/build/docker-compose.yml up " . $builder );
	}
}

sub build_runner_containers {
	foreach my $runner (qw/traffic_monitor traffic_ops traffic_router traffic_stats/) {
		my $dir = $config->{working_dir} . "/infrastructure/docker/" . $runner;
		chdir($dir) || die( "Can't chdir to " . $dir );
		my $rpm_filename = $runner . "*" . "el7.x86_64.rpm";
		my $rpm          = $config->{working_dir} . "/infrastructure/docker/build/artifacts/" . $rpm_filename;
		my $cp_cmd       = "cp " . $rpm . " .";
		print $cp_cmd . " \n";
		system($cp_cmd);
		my $branch     = $config->{git_branch};
		my $args       = " --rm --build-arg RPM=" . $rpm_filename . " --tag " . $runner . ":" . $branch . " .";
		my $dbuild_cmd = "docker build " . $args;
		print $dbuild_cmd. "\n";
		system($dbuild_cmd );
	}
}

# read the config json.
sub configure {
	my $filename = shift;

	my $json_text = do {
		open( my $json_fh, "<:encoding(UTF-8)", $filename )
			or die("Can't open \$filename\": $!\n");
		local $/;
		<$json_fh>;
	};

	my $json = JSON->new;
	#print Dumper($config);
	$config = $json->decode($json_text);
	#print Dumper($config);

}

## rest is from other scripts, should probably be replaced by something better.
sub curl_me {
	my $url           = shift;
	my $retry_counter = 5;
	my $result        = `/usr/bin/curl $CURL_OPTS $url 2>&1`;

	while ( $result =~ m/^curl\: \(\d+\)/ && $retry_counter > 0 ) {
		$result =~ s/(\r|\c|\f|\t|\n)/ /g;
		print "WARN Error receiving $url: $result\n";
		$retry_counter--;
		sleep 5;
		$result = `/usr/bin/curl $CURL_OPTS $url 2>&1`;
	}
	if ( $result =~ m/^curl\: \(\d+\)/ && $retry_counter == 0 ) {
		print "FATAL $url returned in error five times!\n";
		exit 1;
	}
	else {
		#print "INFO Success receiving $url.\n";
	}

	my (@chars) = split( //, $result );
	my $response_code = pop(@chars) . pop(@chars) . pop(@chars);
	$response_code = reverse($response_code);

	#print "DEBUG Received $response_code for $url.\n";
	if ( $response_code >= 400 ) {
		print "ERROR Received error code $response_code for $url!\n";
		return $response_code;
	}
	for ( 0 .. 2 ) { chop($result) }

	if ( $url =~ m/\.json$/ ) {
		eval {
			decode_json($result);
			1;
		} or do {
			my $error = $@;
			print "FATAL $url did not return valid JSON: $result | error: $error\n";
			exit 1;
			}
	}
	my $size = length($result);
	if ( $size == 0 ) {
		print "FATAL URL: $url returned empty!! Bailing!\n";
		exit 1;
	}
	return $result;
}

sub get_cookie {
	my $tm_host  = shift;
	my $tm_login = shift;
	my ( $u, $p ) = split( /:/, $tm_login );

	my $cmd = "curl -vLks -X POST -d 'u=" . $u . "' -d 'p=" . $p . "' " . $tm_host . "/login -o /dev/null 2>&1 | grep Set-Cookie | awk '{print \$3}'";

	#print utput_log_fh "DEBUG Getting cookie with $cmd.\n";
	my $cookie = `$cmd`;
	chomp $cookie;
	$cookie =~ s/;$//;
	if ( $cookie =~ m/mojolicious/ ) {

		#print "DEBUG Cookie is $cookie.\n";
		return $cookie;
	}
	else {
		print "ERROR Cookie not found from Traffic Ops!\n";
		return 0;
	}
}

