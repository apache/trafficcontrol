#!/usr/bin/perl
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
use warnings;
use strict;
use JSON;
use Data::Dumper;
use Term::ReadKey;
use LWP::UserAgent;
use File::Find;
use File::Spec;
use Time::HiRes qw(gettimeofday tv_interval);

#use Data::Compare;
use Test::Deep;
use Test::More;
use List::Compare;

my $config;
my $tmp_dir_base = "/tmp/files";
my $tmp_dir      = $tmp_dir_base . "/ref";
my $CURL_OPTS;
my $cookie;

&configure( $ARGV[1] );

my $perform_snapshot = $config->{perform_snapshot};

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

	&get_files( $config->{to_url} );    # old style

	#    &get_files_new( $config->{to_url} ); # derek style
	&get_crconfigs( $config->{to_url} );
}

sub do_the_compare {
	&compare_all_files();
}

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
		if ($perform_snapshot) {
			print "Generating CRConfig for " . $cdn->{name};
			my $start = [gettimeofday];
			&curl_me( $to_url . "/tools/write_crconfig/" . $cdn->{name} );
			my $load_time = tv_interval($start);
			print " time: " . $load_time . "\n";
		}
		print "Getting CRConfig for " . $cdn->{name};
			my $start = [gettimeofday];
		my $fcontents = &curl_me( $to_url . '/CRConfig-Snapshots/' . $cdn->{name} . '/CRConfig.json' );
		open( my $fh, '>', $dir . '/CRConfig.json' );
			my $load_time = tv_interval($start);
			print " time: " . $load_time . "\n";
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

			print "Getting " . $sample_server . " " . $filename;
			my $start     = [gettimeofday];
			my $fcontents = &curl_me( $to_url . '/genfiles/view/' . $profile_sample{$sample_server} . "/" . $filename );
			my $load_time = tv_interval($start);
			print " load_time: " . $load_time . "\n";
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

sub get_to_passwd {
	my $user = shift;

	print "Type Reference Traffic Ops passwd for " . $user . ":";
	ReadMode('noecho');    # don't echo
	chomp( my $passwd = <STDIN> );
	ReadMode(0);           # back to normal
	print "\n";
	return $passwd;
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
	$config = $json->decode($json_text);
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

