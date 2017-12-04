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

use strict;
use warnings;
use feature qw(switch);
use JSON;
use WWW::Curl::Easy;

my $global;

&init();
&get_traffic_ops_cookie(); 
&get_traffic_ops_ort();
&get_traffic_ops_monitor_cfg();
&get_default_monitor_cfg();
&diff_default_traffic_ops_cfg();
&read_disk_monitor_cfg();

my $update_needed = &check_update_needed();
if ($update_needed) {
	&write_monitor_cfg();
}
else {
	print "DEBUG: Config on disk does not need an update, exiting.\n";
	exit 0;
}

sub init {
	&init_hostname();
	&init_time();
	&init_cli_params();
}

sub init_hostname {
	my $fqdn = `/bin/hostname -f`; chomp($fqdn);
	my ($hostname, undef) = split(/\./, $fqdn, 2);
	if (!defined($hostname)) {
		die("FATAL: Unable to determine host name; please ensure this machine properly configured with FQDN.");
	}
	$global->{'host_name'} = $hostname;
	print "DEBUG: Found hostname: " . $hostname . "\n";
}

sub init_time {
	my $unixtime = time();
	my $date = `/bin/date`; chomp($date);
	$global->{'date'} = $date;
}

sub get_traffic_ops_ort {
	my $ort_url = $global->{'traffic_ops_host'} . "/ort/" . $global->{'host_name'} . "/ort1";
	my $result = &curl_me($ort_url);
	my $ort_ref = decode_json($result);
	if ( defined ($ort_ref->{'profile'}->{'name'}) ) {
		print "DEBUG: Found profile from traffic_ops: " . $ort_ref->{'profile'}->{'name'} . "\n";
		$global->{'traffic_ops_data'}->{'profile'} = $ort_ref->{'profile'}->{'name'};
	}
	else {
		print "ERROR: No profile found in traffic_ops!\n";
	}
	if ( defined ($ort_ref->{'other'}->{'CDN_name'}) ) {
		print "DEBUG: Found CDN name from traffic_ops: " . $ort_ref->{'other'}->{'CDN_name'} . "\n";
		$global->{'traffic_ops_data'}->{'CDN_name'} = $ort_ref->{'other'}->{'CDN_name'};
	}
	else {
		print "ERROR: No CDN name found in traffic_ops, bailing!\n";
		exit 1;
	}
	if ( defined ($ort_ref->{'config_files'}->{'rascal-config.txt'}->{'location'}) ) {
		print "DEBUG: Found location for rascal-config.txt from traffic_ops: " . $ort_ref->{'config_files'}->{'rascal-config.txt'}->{'location'} . "\n";
		$global->{'location'}->{'traffic_monitor_config'} = $ort_ref->{'config_files'}->{'rascal-config.txt'}->{'location'};
	}
	else {
		print "ERROR: No location for rascal-config.txt found in traffic_ops, bailing!\n";
		exit 1;
	}
}

sub get_traffic_ops_monitor_cfg {
	my $health_url = $global->{'traffic_ops_host'} . "/health/" . $global->{'traffic_ops_data'}->{'CDN_name'};
	my $result = &curl_me($health_url);
	my $health_ref = decode_json($result);
	$global->{'traffic_ops_data'}->{'traffic_monitor_config'} = $health_ref->{'traffic_monitor_config'};
	if ( !exists($global->{'traffic_ops_data'}->{'traffic_monitor_config'}) ) {
		print "FATAL: Monitor config not found! Bailing!\n";
		exit 1;
	}
	$global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{'cdnName'} = $global->{'traffic_ops_data'}->{'CDN_name'};
	my $tm_host = $global->{'traffic_ops_host'};
	$tm_host =~ s/(https?\:\/\/)(.*)/$2/;
	$global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{'tm.hostname'} = $tm_host;

	my ($tm_username, $tm_password) = split(/:/, $global->{'traffic_ops_login'}, 2);
	$global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{'tm.auth.username'} = $tm_username;
	$global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{'tm.auth.password'} = $tm_password;

	if (exists $global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{'CDN_name'}) {
		delete $global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{'CDN_name'};
	}
}

sub validate_traffic_ops_monitor_cfg {
	my @missing_params;
	foreach my $param ( keys %{$global->{'disk'}->{'traffic_monitor_config'}} ) {
		if (!exists($global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{$param}) ) {
			push (@missing_params, $param);	
		}
	}
	if (scalar(@missing_params) ) {
		$" = ',';
		print "FATAL: These params are missing from the traffic_ops config: @missing_params \n";
		exit 2;
	}
}

sub write_monitor_cfg {
	my $monitor_config->{'traffic_monitor_config'} = $global->{'traffic_ops_data'}->{'traffic_monitor_config'};
	my $monitor_config_json = JSON->new->utf8->indent->encode($monitor_config);
	if ($global->{'write_mode'} eq 'prompt' ) {
		print "DEBUG: Proposed traffic_monitor_config: \n$monitor_config_json\n";
		my $select = &get_answer();
		if ($select	eq 'Y') {
			&write_monitor_cfg_to_disk($monitor_config_json);
		}
		else {
			print "You elected not to write config to disk, exiting.\n";
			exit 0;
		}
	}
	elsif ($global->{'write_mode'} eq 'auto' ) {
		&write_monitor_cfg_to_disk($monitor_config_json);
	}
}

sub curl_me {
    my $url = shift;
	my $curl = WWW::Curl::Easy->new;
	my $response_body;
	open(my $fileb, ">", \$response_body);
	$curl->setopt(CURLOPT_VERBOSE, 0);
	if ($url =~ m/https/) {
		$curl->setopt(CURLOPT_SSL_VERIFYHOST, 0);
		$curl->setopt(CURLOPT_SSL_VERIFYPEER, 0);
		$curl->setopt(CURLOPT_USERPWD, $global->{'traffic_ops_login'});
	}
    $curl->setopt(CURLOPT_FOLLOWLOCATION, 1);
    $curl->setopt(CURLOPT_CONNECTTIMEOUT, 5);
    $curl->setopt(CURLOPT_TIMEOUT, 15);
    $curl->setopt(CURLOPT_HEADER,0);
    $curl->setopt(CURLOPT_COOKIE, $global->{'traffic_ops_cookie'});
    $curl->setopt(CURLOPT_URL, $url);
	$curl->setopt(CURLOPT_WRITEDATA, $fileb); 
    #$curl->setopt(CURLOPT_HTTPHEADER, @( 'Connection: Keep-Alive', 'Keep-Alive: 300'));
    my $retcode = $curl->perform;
    my $response_code = $curl->getinfo(CURLINFO_HTTP_CODE);
	if ($response_code != 200) {
		print "FATAL: Got HTTP $response_code response for $url! Cannot continue, bailing!\n";
		exit 1;
	}
    if ($response_body =~ m/html/ || $response_body !~ m/\{/ || $response_body !~ m/\}/ || $response_body !~ m/\:/ ) {
        print "FATAL: $url did not return valid JSON!\n";
        exit 1;
    }
    my $size = length($response_body);
    if ($size == 0) {
        print "FATAL: URL: $url returned empty!! Bailing!\n";
        exit 1;
    }
    return $response_body;

}

sub get_traffic_ops_cookie {

    my ( $u, $p ) = split( /:/, $global->{'traffic_ops_login'});
    my $url = $global->{'traffic_ops_host'} . "/login"; 
    my $curl = WWW::Curl::Easy->new;
    my $response_header;
    open(my $fileb, ">", \$response_header);
    $curl->setopt(CURLOPT_VERBOSE, 0);
    if ($url =~ m/https/) {
            $curl->setopt(CURLOPT_SSL_VERIFYHOST, 0);
            $curl->setopt(CURLOPT_SSL_VERIFYPEER, 0);
    }
    $curl->setopt(CURLOPT_POST, 1);
    $curl->setopt(CURLOPT_POSTFIELDS, "u=$u&p=$p");
    $curl->setopt(CURLOPT_FOLLOWLOCATION, 0);
    $curl->setopt(CURLOPT_CONNECTTIMEOUT, 5);
    $curl->setopt(CURLOPT_TIMEOUT, 15);
    $curl->setopt(CURLOPT_HEADER,0);
    $curl->setopt(CURLOPT_URL, $url);
    $curl->setopt(CURLOPT_WRITEHEADER, $fileb);
    my $retcode = $curl->perform;
    my $response_code = $curl->getinfo(CURLINFO_HTTP_CODE);
    if ($response_code != 302) {
        print "FATAL: Got HTTP $response_code response for $url! Cannot continue, bailing!\n";
       exit 1;
    }
    my $size = length($response_header);
    if ($size == 0) {
        print "FATAL: URL: $url returned empty!! Bailing!\n";
        exit 1;
    }
    (my @lines) = split (/\n/, $response_header);
    (my @cookies) = grep /Set-Cookie/, @lines; 
    foreach my $cookie (@cookies) {
	if ($cookie =~ m/Set-Cookie/ && !defined($global->{'traffic_ops_cookie'}) ) {
	    (my $dum, $global->{'traffic_ops_cookie'}) = split(/ /, $cookie);
	    $global->{'traffic_ops_cookie'} =~ s/\;//g;
	    last;
	}
    }
    if ( !defined($global->{'traffic_ops_cookie'}) ) {
        print "FATAL: Didn't get cookie from traffic_ops! Bailing!\n";
	exit 1;
    }
}

sub get_default_monitor_cfg {
    my $default_cfg = `/opt/traffic_monitor/bin/config-doc.sh`;
    my $cfg_json = decode_json($default_cfg);
    $global->{'default'}->{'traffic_monitor_config'} = $cfg_json;
}

sub diff_default_traffic_ops_cfg {
    foreach my $param ( sort keys %{$global->{'traffic_ops_data'}->{'traffic_monitor_config'}} ) {
		if (!exists $global->{'default'}->{'traffic_monitor_config'}->{$param}) {
			print "WARN: Param found in traffic_ops, but not used in Monitor: '$param'\n";
			delete ($global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{$param});
		}
	}
    foreach my $param ( sort keys %{$global->{'default'}->{'traffic_monitor_config'}} ) {
		my $data = $global->{'default'}->{'traffic_monitor_config'}->{$param};
        if (!exists ($global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{$param})) {
            if (exists $data->{'defaultValue'}) {
                printf ("WARN: Param not in traffic_ops: %-40s description: %-120s Using default value of: %-40s \n", $param, $data->{'description'}, $data->{'defaultValue'});
                $global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{$param} = $data->{'defaultValue'};
            }
            else {
                print "FATAL: $param has no default value, and is not in config from traffic_opsonkeys.\n";
            }
        }
    }
}

sub check_update_needed {
	my $update_needed = 0;
	foreach my $param ( sort keys %{$global->{'traffic_ops_data'}->{'traffic_monitor_config'}} ) {
		if (!exists ( $global->{'disk'}->{'traffic_monitor_config'}->{$param} ) ) {
			print "DEBUG: $param needed in config, but does not exist in config on disk.\n";
			$update_needed++;
			next;
		}
		else {
			if ( $global->{'disk'}->{'traffic_monitor_config'}->{$param} ne $global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{$param} ) {
				print "DEBUG: $param value on disk (" . $global->{'disk'}->{'traffic_monitor_config'}->{$param} . ") does not match value needed in config (" . $global->{'traffic_ops_data'}->{'traffic_monitor_config'}->{$param} . ").\n";
				$update_needed++; 
			}
		}
	}
	return $update_needed;
}

sub read_disk_monitor_cfg {
	my $disk_fname = $global->{'location'}->{'traffic_monitor_config'} . "/traffic_monitor_config.js";

	if (! -f $disk_fname) {
		print "WARN: $disk_fname does not exist\n";
		$global->{'disk'}->{'traffic_monitor_config'} = undef;
		return();
	}

	open my $fh, '<', $disk_fname || die("FATAL: Can't open $disk_fname: $!");
	local $/ = undef;
	my $disk_cfg = <$fh>;
	close ($fh);
	my $disk_cfg_json = decode_json($disk_cfg);	
	$global->{'disk'}->{'traffic_monitor_config'} = $disk_cfg_json->{'traffic_monitor_config'};
}

sub usage {
	print "====-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-====\n";
	print "Usage: perl traffic_monitor_config.pl <Traffic Operations Host> <Traffic Operations Login> <Write Mode>\n";
	print "\t<Traffic Operations Host> => format like:\n";
	print "\t\thttps://tm-host.company.net\n";
	print "\n";
	print "\t<Traffic Operations Login> => format like:\n";
	print "\t\tadmin:password\n";
	print "\n";
	print "\t<Write Mode> => choose:\n";
	print "\t\t[auto] Automatic  -- Write new config changes to disk, if needed.\n";
	print "\t\t[prompt] Prompt   -- Prompt before writing config changes to disk, if needed.\n";
	print "====-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-====\n";
	exit 1;
}

sub init_cli_params {
	if ( $#ARGV != 2 ) {
		&usage();
	}
	if ( $ARGV[0] !~ m/^https?\:\/\// ) {
		&usage();
	}
	if ( $ARGV[1] !~ m/\:/ ) {
		&usage();
	}
	if ( $ARGV[2] ne 'auto' && $ARGV[2] ne 'prompt' ) {
		&usage();
	}
	else {
		$global->{'traffic_ops_host'} = $ARGV[0];
		print "DEBUG: traffic_ops selected: " . $global->{'traffic_ops_host'} . "\n";
		$global->{'traffic_ops_login'} = $ARGV[1];
		print "DEBUG: traffic_ops login: " . $global->{'traffic_ops_login'} . "\n";
		$global->{'write_mode'} = $ARGV[2];
		print "DEBUG: Config write mode: " . $global->{'write_mode'} . "\n";
	}
}

sub write_monitor_cfg_to_disk {
	my $monitor_config_json = shift;
	open my $fh, '>', $global->{'location'}->{'traffic_monitor_config'} . "/traffic_monitor_config.js" || die "Can't open " . $global->{'location'}->{'traffic_monitor_config'} . "\n";
	print "DEBUG: Writing " . $global->{'location'}->{'traffic_monitor_config'} . "/traffic_monitor_config.js\n";
	print $fh $monitor_config_json;
	close $fh; 
}

sub get_answer {
	my $select;
	while (!defined($select) || ($select ne 'Y' && $select ne 'n')) {
		print "----------------------------------------------\n";
		print "----OK to write this config to disk? (Y/n) [n]";
		$select = <STDIN>;
		chomp($select);
		print "----------------------------------------------\n";
	}
	return $select;
}
