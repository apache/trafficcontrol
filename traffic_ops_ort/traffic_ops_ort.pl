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

use strict;
use warnings;
use feature qw(switch);
no if $] >= 5.018, warnings => qw( experimental::smartmatch );
use JSON;
use File::Basename;
use File::Path;
use Fcntl qw(:flock);
use MIME::Base64;
use Getopt::Long;

$| = 1;
my $date           = `/bin/date`;
chomp($date);
print "$date\n";

# supported redhat/centos releases
my %supported_el_release = ( "EL6" => 1, "EL7" => 1, "EL8" => 1);

my $dispersion = 300;
my $retries = 5;
my $wait_for_parents = 1;
my $login_dispersion = 0;
my $reval_wait_time = 60;
my $reval_in_use = 0;
my $rev_proxy_disable = 0;
my $skip_os_check = 0;
my $override_hostname_short = '';
my $to_timeout_ms = 30000;
my $syncds_updates_ipallow = 0;
my $traffic_ops_insecure = 0;

GetOptions( "dispersion=i"       => \$dispersion, # dispersion (in seconds)
            "retries=i"          => \$retries,
            "wait_for_parents=i" => \$wait_for_parents,
            "login_dispersion=i" => \$login_dispersion,
            "rev_proxy_disable=i" => \$rev_proxy_disable,
            "skip_os_check=i" => \$skip_os_check,
            "override_hostname_short=s" => \$override_hostname_short,
            "to_timeout_ms=i" => \$to_timeout_ms,
            "syncds_updates_ipallow=i" => \$syncds_updates_ipallow,
            "traffic_ops_insecure=i" => \$traffic_ops_insecure,
          );

if ( $#ARGV < 1 ) {
	&usage();
}

my $log_level = 0;
$ARGV[1] = uc( $ARGV[1] );
given ( $ARGV[1] ) {
	when ("ALL")   { $log_level = 255; }
	when ("TRACE") { $log_level = 127; }
	when ("DEBUG") { $log_level = 63; }
	when ("INFO")  { $log_level = 31; }
	when ("WARN")  { $log_level = 15; }
	when ("ERROR") { $log_level = 7; }
	when ("FATAL") { $log_level = 3; }
	when ("NONE")  { $log_level = 1; }
	default        { &usage(); }
}

my $TO_URL = undef;
my $TO_USER = undef;
my $TO_PASS = undef;

if ( defined( $ARGV[2] ) ) {
	if ( $ARGV[2] !~ /^https*:\/\/.*$/ ) {
		&usage();
	}
	else {
		$TO_URL = $ARGV[2];
		$TO_URL =~ s/\/*$//g;
	}
}
else {
	&usage();
}

if ( defined( $ARGV[3] ) ) {
	if ( $ARGV[3] !~ m/^.*\:.*$/ ) {
		&usage();
	}
	else {
		( $TO_USER, $TO_PASS ) = split( /:/, $ARGV[3] );
	}
}
else {
	&usage();
}

#### Script mode constants ####
my $INTERACTIVE = 0;
my $REPORT      = 1;
my $BADASS      = 2;
my $SYNCDS      = 3;
my $REVALIDATE  = 4;
#### Logging constants for bit shifting ####
my $ALL   = 7;
my $TRACE = 6;
my $DEBUG = 5;
my $INFO  = 4;
my $WARN  = 3;
my $ERROR = 2;
my $FATAL = 1;
my $NONE  = 0;

my $RELEASE = &os_version();
( $log_level >> $DEBUG ) && print "DEBUG OS release is $RELEASE.\n";

my $script_mode = &check_script_mode();
&check_run_user();
&check_only_copy_running();
&check_log_level();

#### Constants to track update status ####
my $UPDATE_TROPS_NOTNEEDED  = 0;
my $UPDATE_TROPS_NEEDED     = 1;
my $UPDATE_TROPS_SUCCESSFUL = 2;
my $UPDATE_TROPS_FAILED     = 3;

#### Other constants #####
my $START_FAILED        = 0;
my $START_SUCCESSFUL    = 1;
my $ALREADY_RUNNING     = 2;
my $START_NOT_ATTEMPTED = 3;
my $CLEAR               = 0;
my $PLUGIN_NO           = 0;
my $PLUGIN_YES          = 1;
#### Constants for config file changes ####
my $CFG_FILE_UNCHANGED         = 0;
my $CFG_FILE_NOT_PROCESSED     = 1;
my $CFG_FILE_CHANGED           = 2;
my $CFG_FILE_PREREQ_FAILED     = 3;
my $CFG_FILE_ALREADY_PROCESSED = 4;

my $unixtime       = time();
my $hostname_short = `/bin/hostname -s`;
if ($override_hostname_short ne '') {
	$hostname_short = $override_hostname_short;
}
chomp($hostname_short);

my $atstccfg_timeout_arg='';
if (length $to_timeout_ms > 0) {
		$atstccfg_timeout_arg = "--traffic-ops-timeout-milliseconds=$to_timeout_ms";
}

my $atstccfg_cmd = '/opt/ort/atstccfg';

my $atstccfg_arg_disable_proxy = '';
if ($rev_proxy_disable != 0) {
	$atstccfg_arg_disable_proxy = '--traffic-ops-disable-proxy';
}

my $atstccfg_insecure_arg = "";
if ($traffic_ops_insecure == 1) {
	$atstccfg_insecure_arg = "--traffic-ops-insecure";
}

my $TMP_BASE  = "/tmp/ort";

my $LOG_BASE  = "/var/log/ort"; # TODO add inferring ORT install location, and allowing / vs /opt install
my $atstccfg_log_path = $LOG_BASE . '/atstccfg.log';

# add any special yum options for your environment here; this variable is used with all yum commands
my $YUM_OPTS = "";
( $log_level >> $DEBUG ) && print "DEBUG YUM_OPTS: $YUM_OPTS.\n";

my $TS_HOME      = "/opt/trafficserver";
my $TRAFFIC_CTL = $TS_HOME . "/bin/traffic_ctl";

my $out          = `/usr/bin/yum $YUM_OPTS clean expire-cache 2>&1`;
my $return       = &check_output($out);
my @config_files = ();

#### Process reboot tracker
my $traffic_ctl_needed           = 0;
my $sysctl_p_needed              = 0;
my $trafficserver_restart_needed = 0;

#### Process installed tracker
my $installed_new_ssl_keys    = 0;
my %install_tracker;

my $cfg_file_tracker = undef;

my $ats_config_dir = get_ats_config_dir();

####-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-####
#### Start main flow
####-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-####


#### Delete /tmp dirs older than one week
if ( $script_mode == $BADASS || $script_mode == $INTERACTIVE || $script_mode == $SYNCDS || $script_mode == $REVALIDATE ) {
	&smart_mkdir($TMP_BASE);
	&clean_tmp_dirs();
}

my $header_comment = &get_header_comment();

#### If this is a syncds run, check to see if we can bail.
my $syncds_update = 0;

if ( $script_mode == $REVALIDATE ) {
	( $syncds_update ) = &check_revalidate_state();
	if ( $syncds_update < 1 ) {
		exit 1;
	}
}
else {
	( $syncds_update ) = &check_syncds_state();
}

$cfg_file_tracker = &get_cfg_file_list( $hostname_short, $script_mode );

if ( $script_mode == $REVALIDATE ) {
	( $log_level >> $INFO ) && print "\nINFO: ======== Revalidating, no package processing needed ========\n";
}
else {
	( $log_level >> $INFO ) && print "\nINFO: ======== Start processing packages ========\n";
	&process_packages( $hostname_short );
	# get the ats user's UID after package installation in case this is the initial badass
	( $log_level >> $INFO ) && print "\nINFO: ======== Start second package processing run ========\n";
	&process_chkconfig( $hostname_short );
}

&process_config_files();

foreach my $file ( keys ( %{$cfg_file_tracker} ) ) {
	if ( exists($cfg_file_tracker->{$file}->{'remap_plugin_config_file'}) && $cfg_file_tracker->{$file}->{'remap_plugin_config_file'} ) {
		if ( exists($cfg_file_tracker->{$file}->{'change_applied'}) && $cfg_file_tracker->{$file}->{'change_applied'} ) {
			( $log_level >> $DEBUG ) && print "\nDEBUG $file is a remap plugin config file, and was changed. remap.config needs touched.  ========\n";
			&touch_file('remap.config');
			last;
		}
	}
}

if ( ($installed_new_ssl_keys) && !$cfg_file_tracker->{'ssl_multicert.config'}->{'change_applied'} ) {
	my $return = &touch_file('ssl_multicert.config');
	if ($return) {
		if ( $syncds_update == $UPDATE_TROPS_NEEDED ) {
			$syncds_update = $UPDATE_TROPS_SUCCESSFUL;
		}
		$traffic_ctl_needed++;
	}
}

&start_restart_services();

if ( $sysctl_p_needed && $script_mode != $SYNCDS ) {
	&run_sysctl_p();
}

if ( $script_mode != $REPORT ) {
	&update_trops();
}

####-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-####
#### End main flow
####-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-####

####-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-####
#### Subroutines
####-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-####

# Returns the ATS config directory, if it can find it.
# Tries rpm (yum) first, then falls back to find.
# If it fails to find it, logs an error and returns the empty string.
sub get_ats_config_dir {
	my $dir = `rpm -ql trafficserver | grep -E 'etc/trafficserver\$' | tail -1`;
	$dir =~ s/^\s+|\s+$//g; # trim leading and trailing whitespace
	if ( $dir eq "" ) {
		$dir = `find / -type d -path '*/etc/trafficserver' | tail -1`;
		$dir =~ s/^\s+|\s+$//g; # trim leading and trailing whitespace
	}
	if ( ! length $dir ) {
		# if it became undefined somehow, make sure we're returning ""
		$dir = "";
	}
	if ( $dir eq "" ) {
		( $log_level >> $ERROR ) && print "ERROR Failed to find config directory, using empty string!\n";
	}
	return $dir;
}

sub revalidate_while_sleeping {
	$syncds_update = &check_revalidate_state(1);
	if ( $syncds_update > 0 ) {
		$script_mode = $REVALIDATE;
		$cfg_file_tracker = &get_cfg_file_list( $hostname_short, $script_mode );

		&process_config_files();

		&start_restart_services();

		&update_trops();

		$traffic_ctl_needed = 0;
	}
}

sub os_version {
	my $release = "UNKNOWN";
	if (`uname -r` =~ m/.+(el\d)((?:\.\w+)|(?:_\w+))*\.x86_64/)  {
		$release = uc $1;
	}
	if (!exists $supported_el_release{$release} && !$skip_os_check) {
		die("skip_os_check: $skip_os_check dispersion: $dispersion unsupported el_version: $release");
	}
	return $release;
}

sub usage {
	print "====-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-====\n";
	print "Usage: ./traffic_ops_ort.pl <Mode> <Log_Level> <Traffic_Ops_URL> <Traffic_Ops_Login> [optional flags]\n";
	print "\t<Mode> = interactive - asks questions during config process.\n";
	print "\t<Mode> = report - prints config differences and exits.\n";
	print "\t<Mode> = badass - attempts to fix all config differences that it can.\n";
	print "\t<Mode> = syncds - syncs delivery services with what is configured in Traffic Ops.\n";
	print "\t<Mode> = revalidate - checks for updated revalidations in Traffic Ops and applies them.  Requires Traffic Ops 2.1.\n";
	print "\n";
	print "\t<Log_Level> => ALL, TRACE, DEBUG, INFO, WARN, ERROR, FATAL, NONE\n";
	print "\n";
	print "\t<Traffic_Ops_URL> = URL to Traffic Ops host. Example: https://trafficops.company.net\n";
	print "\n";
	print "\t<Traffic_Ops_Login> => Example: 'username:password' \n";
	print "\n\t[optional flags]:\n";
	print "\t   dispersion=<time>              => wait a random number between 0 and <time> before starting. Default = 300.\n";
	print "\t   login_dispersion=<time>        => wait a random number between 0 and <time> before login. Default = 0.\n";
	print "\t   retries=<number>               => retry connection to Traffic Ops URL <number> times. Default = 3.\n";
	print "\t   wait_for_parents=<0|1>         => do not update if parent_pending = 1 in the update json. Default = 1, wait for parents.\n";
	print "\t   rev_proxy_disable=<0|1>        => bypass the reverse proxy even if one has been configured Default = 0.\n";
	print "\t   skip_os_check=<0|1>            => bypass the check for a supported CentOS version. Default = 0.\n";
	print "\t   override_hostname_short=<text> => override the short hostname of the OS for config generation. Default = ''.\n";
	print "\t   to_timeout_ms=<time>           => the Traffic Ops request timeout in milliseconds. Default = 30000 (30 seconds).\n";
	print "\t   syncds_updates_ipallow=<0|1>   => Update ip_allow.config in syncds mode, which may trigger an ATS bug blocking random addresses on load! Default = 0, only update on badass and restart.\n";
	print "\t   traffic_ops_insecure=<0|1>     => Turns off certificate checking when connecting to Traffic Ops.\n";
	print "====-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-====\n";
	exit 1;
}

sub process_cfg_file {
	my $cfg_file = shift;

	if ( $cfg_file eq "" ) {
		( $log_level >> $ERROR ) && print "ERROR Config file name is empty!\n";
		$cfg_file_tracker->{$cfg_file}->{'audit_failed'}++;
		return $CFG_FILE_NOT_PROCESSED;
	}

	my $result = $cfg_file_tracker->{$cfg_file}->{'body'};

	my $return_code = 0;
	return $CFG_FILE_ALREADY_PROCESSED
		if ( defined( $cfg_file_tracker->{$cfg_file}->{'audit_complete'} ) && $cfg_file_tracker->{$cfg_file}->{'audit_complete'} > 0 );

	( $log_level >> $INFO ) && print "\nINFO: ======== Start processing config file: $cfg_file ========\n";

	my $cfg_dir = dirname($cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'});
	if ( length($cfg_dir) == 0 ) {
			( $log_level >> $ERROR ) && print "ERROR No location information for $cfg_file.\n";
		return $CFG_FILE_NOT_PROCESSED;
	}

	&smart_mkdir($cfg_dir);

	# Process ##OVERRIDE## remap rules (from anymap rawtext)
	if ( $cfg_file eq "remap.config" ) {
		$result = &adv_preprocessing_remap( \$result );
	}

	my @db_file_lines = @{ &scrape_unencode_text($result) };
	@db_file_lines = @{ &scrape_canned_comments(\@db_file_lines) };
	my $file = $cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'};
	return $CFG_FILE_PREREQ_FAILED if ( !&prereqs_ok( $cfg_file, \@db_file_lines ) );

	my @disk_file_lines;
	if ( -e $file ) {
		return $CFG_FILE_NOT_PROCESSED if ( !&can_read_write_file($cfg_file) );
		@disk_file_lines = @{ &open_file_get_contents($file) };
		@disk_file_lines = @{ &scrape_canned_comments(\@disk_file_lines) };
	}

	# First, check if the file to be generated would be identical including order
	my $change_needed = ( join( '\0', @disk_file_lines ) ne join( '\0', @db_file_lines ) );

	# if different, look deeper to see if we care about the diffs
	if ( $change_needed ) {
		# diff_file_lines has all the debug we want
		my @return = &diff_file_lines( $cfg_file, \@db_file_lines, \@disk_file_lines );
		my $order_dependent = ( $cfg_file eq 'logs_xml.config' || $cfg_file =~ m/\.cer$/ || $cfg_file =~ m/hdr\_rw\_(.*)\.config$/ );

		# if the files aren't order dependent then relax the criteria
		if ( ! $order_dependent ) {
			my @db_lines_missing   = @{ shift(@return) };
			my @disk_lines_missing = @{ shift(@return) };

			if ( scalar(@disk_lines_missing) == 0 && scalar(@db_lines_missing) == 0 ) {
				# all lines accounted for
				$change_needed = undef;
			}
		}
	}

	if ($change_needed && $cfg_file eq "ip_allow.config" && $syncds_updates_ipallow != 1) {
		if ($script_mode == $BADASS) {
			$trafficserver_restart_needed++;
		} else {
			( $log_level >> $ERROR ) && print "ERROR Not in badass mode, but ip_allow.config changed! Changing that file will cause ATS to break the next time it Reloads! Ignoring file!! This will cause this server to reject any new servers! ORT must be run in badass mode to get the ip_allow.config change and permit the necessary client!\n";
			$change_needed = undef;
		}
	}

	if ( $change_needed ) {
		$cfg_file_tracker->{$cfg_file}{'change_needed'}++;
		( $log_level >> $ERROR ) && print "ERROR $file needs updated.\n";
		&backup_file( $cfg_file, \$result );
	}
	else {
		( $log_level >> $INFO ) && print "INFO: All lines match TrOps for config file: $cfg_file.\n";
		$cfg_file_tracker->{$cfg_file}{'change_needed'} = 0;
		( $log_level >> $TRACE ) && print "TRACE Setting change not needed for $cfg_file.\n";
		$return_code = $CFG_FILE_UNCHANGED;
	}

	if ( $cfg_file eq "50-ats.rules" ) {
		&adv_processing_udev( \@db_file_lines );
	}

	( $log_level >> $INFO )
		&& print "INFO: ======== End processing config file: $cfg_file for service: " . $cfg_file_tracker->{$cfg_file}->{'service'} . " ========\n";
	$cfg_file_tracker->{$cfg_file}->{'audit_complete'}++;

	return $return_code;
}

sub systemd_service_set {
	my $systemd_service = shift;
	my $systemd_service_enable = shift;

	my $command = "/bin/systemctl $systemd_service_enable $systemd_service";
	`$command 2>/dev/null`;
	if ($? == 0) {
		return 1;
	} else {
		return 0;
	}
}

sub systemd_service_chk {
	my $service = shift;

	my $status = "disabled";
	open(FH, "/bin/systemctl list-unit-files ${service}.service|") or die ("/bin/systemctl: $!");
	while(<FH>) {
		chomp($_);
		if ($_ =~ m/$service\.service\s(\w+)/) {
			$status = $1;
		}
	}
	close(FH);
	return $status;
}

sub systemd_service_status {
	my $pkg_name = shift;
	my $running_string;
	my $running = 0;
	my $pid;
	my $prog;

	open(FH, "/bin/systemctl status $pkg_name|") or die ("/bin/systemctl $!");
	while(<FH>) {
		chomp ($_);
		if ($_ =~ m/\s+Active:\s+active\s\(running\)/) {
			$running = 1;
		}
		if ($_ =~ m/\s+Main\sPID:\s(\d+)\s+\((\w+)\)/) {
			$pid = $1;
			$prog = $2
		}
	}
	close(FH);
	if ($running) {
		$running_string = "$prog (pid $pid) is running...";
	} else {
		$running_string = "$pkg_name is stopped";
	}

	return $running_string;
}

sub start_service {
	my $pkg_name = shift;

	( $log_level >> $DEBUG ) && print "DEBUG start_service called for $pkg_name.\n";

	my $pkg_running;
	if (($RELEASE eq "EL7") || ($RELEASE eq "EL8")) {
		$pkg_running = &systemd_service_status($pkg_name);
	} else {
		$pkg_running  = `/sbin/service $pkg_name status`;
	}
	my $running_string = "";
	if ( $pkg_name eq "trafficserver" ) {
		$running_string = "traffic_manager|traffic_cop";
	}
	else {
		$running_string = $pkg_name;
	}
	if ( $running_string ne "" ) {
		if ( $pkg_running !~ m/$running_string \(pid\s+(\d+)\) is running.../ ) {
			if ( $script_mode == $REPORT || $script_mode == $SYNCDS ) {
				( $log_level >> $ERROR ) && print "ERROR $pkg_name is not running.\n";
				$pkg_running = $START_NOT_ATTEMPTED;
			}
			elsif ( $script_mode == $BADASS ) {
				( $log_level >> $ERROR ) && print "ERROR $pkg_name needs started. Trying to do that now.\n";
				my $pkg_start_output = `/sbin/service $pkg_name start`;
				my $pkg_started = 0;
				if (($RELEASE eq "EL7") || ($RELEASE eq "EL8")) {
					my $_st = &systemd_service_status($pkg_name);
					if ($_st =~ m/\(pid\s+(\d+)\) is running.../) {
						$pkg_started++;
					}
				} else {
					( my @output_lines ) = split( /\n/, $pkg_start_output );
					foreach my $ol (@output_lines) {
						if ( $ol =~ m/\[.*\]/ && $ol =~ m/OK/ ) {
							$pkg_started++;
						}
					}
				}
				if ($pkg_started) {
					( $log_level >> $ERROR ) && print "ERROR $pkg_name started successfully.\n";
					$pkg_running = $START_SUCCESSFUL;
				}
				else {
					$pkg_start_output =~ s/\n/\t/g;
					$pkg_start_output =~ s/\r/\t/g;
					( $log_level >> $ERROR ) && print "ERROR $pkg_name failed to start, error is: $pkg_start_output.\n";
					$pkg_running = $START_FAILED;
				}
			}
			elsif ( $script_mode == $INTERACTIVE ) {
				my $select = 'Y';
				( $log_level >> $ERROR ) && print "ERROR $pkg_name is not running. Should I start it now? (Y/n) [n]";
				$select = <STDIN>;
				chomp($select);
				if ( $select =~ m/Y/ ) {
					( $log_level >> $ERROR ) && print "ERROR $pkg_name needs started. Trying to do that now.\n";
					my $pkg_start_output = `/sbin/service $pkg_name start`;
					my $pkg_started = 0;
					if (($RELEASE eq "EL7") || ($RELEASE eq "EL8")) {
						my $_st = &systemd_service_status($pkg_name);
						if ($_st =~ m/\(pid\s+(\d+)\) is running.../) {
							$pkg_started++;
						}
					} else {
						( my @output_lines ) = split( /\n/, $pkg_start_output );
						foreach my $ol (@output_lines) {
							if ( $ol =~ m/\[.*\]/ && $ol =~ m/OK/ ) {
								$pkg_started++;
							}
						}
					}
					if ($pkg_started) {
						( $log_level >> $DEBUG ) && print "DEBUG $pkg_name started successfully.\n";
						$pkg_running = $START_SUCCESSFUL;
					}
					else {
						$pkg_start_output =~ s/\n/\t/g;
						( $log_level >> $ERROR ) && print "ERROR $pkg_name failed to start, error is: $pkg_start_output.\n";
						$pkg_running = $START_FAILED;
					}
				}
			}
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR $pkg_name is running.\n";
			$pkg_running = $ALREADY_RUNNING;
		}
	}
	else {
		( $log_level >> $FATAL ) && print "FATAL Unrecognized service: $pkg_name. Not starting $pkg_name.\n";
		$pkg_running = $START_NOT_ATTEMPTED;
	}
	return $pkg_running;
}

sub restart_service {
	my $pkg_name = $_[0];

	my $pkg_running;
	if (($RELEASE eq "EL7") || ($RELEASE eq "EL8")) {
		$pkg_running = &systemd_service_status($pkg_name);
	} else {
		$pkg_running  = `/sbin/service $pkg_name status`;
	}
	my $running_string = "";
	if ( $pkg_name eq "trafficserver" ) {
		$running_string = "traffic_manager|traffic_cop";
	}
	if ( $running_string ne "" ) {
		if ( $pkg_running =~ m/$running_string \(pid  (\d+)\) is running.../ ) {
			if ( $script_mode == $REPORT ) {
				( $log_level >> $ERROR ) && print "ERROR $pkg_name needs to be restarted. Please run 'service $pkg_name restart' to fix.\n";
			}
			if ( $script_mode == $BADASS ) {
				( $log_level >> $ERROR ) && print "ERROR Trying to restart $pkg_name.\n";
				my $pkg_start_output = `/sbin/service $pkg_name restart`;
				( my @output_lines ) = split( /\n/, $pkg_start_output );
				my $pkg_started = 0;
				foreach my $ol (@output_lines) {
					if ( $ol =~ m/\[.*\]/ && $ol =~ m/OK/ ) {
						$pkg_started++;
					}
				}
				if ($pkg_started) {
					( $log_level >> $ERROR ) && print "ERROR $pkg_name restarted successfully.\n";
					$pkg_running = $START_SUCCESSFUL;
				}
				else {
					$pkg_start_output =~ s/\n/\t/g;
					( $log_level >> $ERROR ) && print "ERROR $pkg_name failed to restart, error is: $pkg_start_output.\n";
					$pkg_running = $START_FAILED;
				}
			}
			if ( $script_mode == $INTERACTIVE ) {
				my $select = 'Y';
				( $log_level >> $ERROR ) && print "ERROR $pkg_name needs to be restarted. Should I restart it now? (Y/n) [n]";
				$select = <STDIN>;
				chomp($select);
				if ( $select =~ m/Y/ ) {
					( $log_level >> $DEBUG ) && print "DEBUG Trying to restart $pkg_name.\n";
					my $pkg_start_output = `/sbin/service $pkg_name restart`;
					( my @output_lines ) = split( /\n/, $pkg_start_output );
					my $pkg_started = 0;
					foreach my $ol (@output_lines) {
						if ( $ol =~ m/\[.*\]/ && $ol =~ m/OK/ ) {
							$pkg_started++;
						}
					}
					if ($pkg_started) {
						( $log_level >> $DEBUG ) && print "DEBUG $pkg_name restarted successfully.\n";
						$pkg_running = $START_SUCCESSFUL;
					}
					else {
						$pkg_start_output =~ s/\n/\t/g;
						( $log_level >> $ERROR ) && print "ERROR $pkg_name failed to restart, error is: $pkg_start_output.\n";
						$pkg_running = $START_FAILED;
					}
				}
			}
		}
		else {
			( $log_level >> $DEBUG ) && print "DEBUG $pkg_name is not running! This shouldn't happnen, $pkg_name must have died recently!\n";
			$pkg_running = $START_FAILED;
		}
	}
	else {
		( $log_level >> $FATAL ) && print "FATAL Unrecognized service: $pkg_name. Not restarting $pkg_name.\n";
		$pkg_running = $START_NOT_ATTEMPTED;
	}
	return $pkg_running;
}

sub smart_mkdir {
	my $dir = shift;

	if ( !-d ($dir) ) {
		if ( $script_mode == $BADASS || $script_mode == $INTERACTIVE || $script_mode == $SYNCDS || $script_mode == $REVALIDATE ) {
			( $log_level >> $TRACE ) && print "TRACE Directory to create if needed: $dir\n";
			system("/bin/mkdir -p $dir");
			if ( $dir =~ m/config_trops/ ) {
				( $log_level >> $DEBUG )
					&& print "DEBUG Temp directory created: $dir. Config files from Traffic Ops will be placed here for future processing.\n";
			}
			elsif ( $dir =~ m/config_bkp/ ) {
				( $log_level >> $DEBUG ) && print "DEBUG Backup directory created: $dir. Config files will be backed up here.\n";
			}
			else {
				( $log_level >> $DEBUG ) && print "DEBUG Directory created: $dir.\n";
			}
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR Directory: $dir doesn't exist, and was not created.\n";
		}
	}
	else {
		( $log_level >> $TRACE ) && print "TRACE Directory: $dir exists.\n";
	}
}

sub clean_tmp_dirs {
	my $old_time = $unixtime - 604800;
	( $log_level >> $ERROR ) && print "ERROR Deleting directories older than $old_time\n";
	opendir( DIR, $TMP_BASE ) || err ("Could not open $TMP_BASE: $!\n");
	my @dirs = grep( /\d{10}/, readdir(DIR) );
	closedir(DIR);
	foreach my $dir (@dirs) {
		if ( $dir <= $old_time ) {
			( $log_level >> $ERROR ) && print "ERROR Deleting directory $TMP_BASE/$dir\n";
			system("rm -rf $TMP_BASE/$dir");
		}
	}
}

sub update_trops {
	my $update_result = 0;
	if ( $syncds_update == $UPDATE_TROPS_NOTNEEDED ) {
		( $log_level >> $DEBUG ) && print "DEBUG Traffic Ops does not require an update at this time.\n";
		return 0;
	}
	elsif ( $syncds_update == $UPDATE_TROPS_FAILED ) {
		( $log_level >> $ERROR )
			&& print "ERROR Traffic Ops requires an update, but applying the update locally failed. Traffic Ops is not being updated!\n";
		return 1;
	}
	elsif ( $syncds_update == $UPDATE_TROPS_SUCCESSFUL ) {
		( $log_level >> $ERROR ) && print "ERROR Traffic Ops required an update, and it was applied successfully. Clearing update state in Traffic Ops.\n";
		$update_result++;
	}
	elsif ( $syncds_update == $UPDATE_TROPS_NEEDED ) {
		( $log_level >> $ERROR )
			&& print
			"ERROR Traffic Ops is signaling that an update is ready to be applied, but none was found! Clearing update state in Traffic Ops anyway.\n";
		$update_result++;
	}
	if ($update_result) {
		#need to know if reval_pending is supported
		my $upd_json = get_update_status();

		my $upd_pending = ( defined( $upd_json->{'upd_pending'} ) ) ? $upd_json->{'upd_pending'} : undef;
		my $reval_pending = ( defined( $upd_json->{'reval_pending'} ) ) ? $upd_json->{'reval_pending'} : undef;

		if ( $script_mode == $INTERACTIVE ) {
			( $log_level >> $ERROR ) && print "ERROR Traffic Ops needs updated. Should I do that now? [Y/n] (n): ";
			my $select = 'n';
			$select = <STDIN>;
			chomp($select);
			if ( $select =~ m/Y/ ) {
				&send_update_to_trops($CLEAR);
			}
			else {
				( $log_level >> $ERROR )
					&& print "ERROR Traffic Ops needs updated. You elected not to do that now; you should probably do that manually.\n";
			}
		}
		elsif ( $script_mode == $BADASS || $script_mode == $SYNCDS ) {
			if ( defined $reval_pending ) {
				&send_update_to_trops($CLEAR, $reval_pending );
			}
			else {
				&send_update_to_trops($CLEAR, $CLEAR);
			}
		}
		elsif ( $script_mode = $REVALIDATE ) {
			&send_update_to_trops($upd_pending, $CLEAR);
		}
	}
}

sub send_update_to_trops {
	my $status = shift;
	my $reval_status = shift;

	my $upd_str='false';
	my $reval_str='false';
	if ( $status != 0 ) {
		$upd_str='true';
	}
	if ( $reval_status != 0 ) {
		$reval_str='true';
	}

	my $response = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --set-queue-status=$upd_str --set-reval-status=$reval_str 2>>$atstccfg_log_path`;
	my $atstccfg_exit_code = $?;
	if ($atstccfg_exit_code != 0) {
		( $log_level >> $ERROR ) && printf("ERROR sending update status with atstccfg (via Traffic Ops). See $atstccfg_log_path.\n");
	}

	( $log_level >> $DEBUG ) && print "DEBUG Response from Traffic Ops is: " . $response . ".\n";
}

sub get_update_status {
	my $upd_ref = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --get-data=update-status 2>>$atstccfg_log_path`;
	my $atstccfg_exit_code = $?;
	if ($atstccfg_exit_code != 0) {
		( $log_level >> $ERROR ) && printf("ERROR getting update status from atstccfg (via Traffic Ops). See $atstccfg_log_path.\n");
		exit 1;
	}

	my $upd_json = decode_json($upd_ref);

	##Some versions of Traffic Ops had the 1.3 API but did not have the use_reval_pending field.  If this field is not present, exit.
	if ( !defined( $upd_json->{'use_reval_pending'} ) ) {
		my $info_ref = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --get-data=system-info 2>>$atstccfg_log_path`;
		my $atstccfg_exit_code = $?;
		if ($atstccfg_exit_code != 0) {
			( $log_level >> $ERROR ) && printf("ERROR Unable to get status of use_reval_pending parameter.  Stopping.\n");
			exit 1;
		}
		my $info_json = decode_json($info_ref);
		if (defined( $info_json->{'use_reval_pending'} ) ) {
			$reval_in_use = $info_json->{'use_reval_pending'};
		}
	}
	else {
		$reval_in_use = $upd_json->{'use_reval_pending'};
	}
	return $upd_json;
}

sub check_revalidate_state {
	my $sleep_override = shift;

	my $syncds_update = 0;

	( $log_level >> $DEBUG ) && print "DEBUG Checking revalidate state.\n";
	if ( $script_mode == $REVALIDATE || $sleep_override == 1 ) {
		## The herd is about to get /update/<hostname>

		my $upd_json = get_update_status();

		if ( $reval_in_use == 0 ) {
			( $log_level >> $ERROR ) && print "ERROR Update URL: Instant invalidate is not enabled.  Separated revalidation requires upgrading to Traffic Ops version 2.2 and enabling this feature.\n";
			return($UPDATE_TROPS_NOTNEEDED);
		}
		my $reval_pending = $upd_json->{'reval_pending'};
		if ( $reval_pending == 1 ) {
			( $log_level >> $ERROR ) && print "ERROR Traffic Ops is signaling that a revalidation is waiting to be applied.\n";
			$syncds_update = $UPDATE_TROPS_NEEDED;

			my $parent_reval_pending = $upd_json->{'parent_reval_pending'};
			if ( !defined($parent_reval_pending) ) {
				( $log_level >> $ERROR ) && print "ERROR Update URL: did not have an parent_reval_pending key.  Separated revalidation requires upgrading to Traffic Ops version 2.2.  Unable to continue!\n";
				return($UPDATE_TROPS_NOTNEEDED);
			}
			if ( $parent_reval_pending == 1 ) {
				( $log_level >> $ERROR ) && print "ERROR Traffic Ops is signaling that my parents need to revalidate.\n";
				return($UPDATE_TROPS_NOTNEEDED);
			}
		}
		elsif ( $script_mode == $REVALIDATE && $reval_pending != 1 ) {
			( $log_level >> $ERROR ) && print "ERROR In revalidate mode, but no update needs to be applied. I'm outta here.\n";
			return($UPDATE_TROPS_NOTNEEDED);
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR Traffic Ops is signaling that no revalidations are waiting to be applied.\n";
		}

		my $stj = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --get-data=statuses 2>>$atstccfg_log_path`;
		my $atstccfg_exit_code = $?;
		if ( $atstccfg_exit_code != 0 ) {
			( $log_level >> $ERROR ) && print "Statuses URL: returned $stj! Skipping creation of status file.\n";
		}

		my $statuses = decode_json($stj);
		my $my_status = ( defined( $upd_json->{'status'} ) ) ? $upd_json->{'status'} : undef;

		if ( defined($my_status) ) {
			( $log_level >> $DEBUG ) && print "DEBUG Found $my_status status from Traffic Ops.\n";
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR Returning; did not find status from Traffic Ops!\n";
			return ($syncds_update);
		}

		my $status_dir  = dirname($0) . "/status";
		my $status_file = $status_dir . "/" . $my_status;

		if ( !-f $status_file ) {
			( $log_level >> $ERROR ) && print "ERROR status file $status_file does not exist.\n";
		}

		for my $status ( @{$statuses} ) {
			next if ( $status->{name} eq $my_status );
			my $other_status = $status_dir . "/" . $status->{name};

			if ( -f $other_status && $status->{name} ne $my_status ) {
				( $log_level >> $ERROR ) && print "ERROR Other status file $other_status exists.\n";
				if ( $script_mode != $REPORT ) {
					( $log_level >> $DEBUG ) && print "DEBUG Removing $other_status\n";
					unlink($other_status);
				}
			}
		}
	}

	return ($syncds_update);
}


sub check_syncds_state {

	my $syncds_update = 0;
	my $random_duration = int( rand($dispersion) );

	( $log_level >> $DEBUG ) && print "DEBUG Checking syncds state.\n";
	if ( $script_mode == $SYNCDS || $script_mode == $BADASS || $script_mode == $REPORT ) {
		## The herd is about to get /update/<hostname>
		## need to check if revalidation is being used first.

		my $upd_json = get_update_status();

		my $upd_pending = ( defined( $upd_json->{'upd_pending'} ) ) ? $upd_json->{'upd_pending'} : undef;
		if ( !defined($upd_pending) ) {
			( $log_level >> $ERROR ) && print "ERROR Update URL: did not have an upd_pending key.\n";
			if ( $script_mode != $SYNCDS ) {
				return $syncds_update;
			}
			else {
				( $log_level >> $ERROR ) && print "ERROR Invalid JSON for update_status. Exiting, not sure what else to do.\n";
				exit 1;
			}
		}

		if ( $upd_pending == 1 ) {
			( $dispersion > 0 ) && &sleep_timer($random_duration);
			( $log_level >> $ERROR ) && print "ERROR Traffic Ops is signaling that an update is waiting to be applied.\n";
			$syncds_update = $UPDATE_TROPS_NEEDED;

			my $parent_pending = ( defined( $upd_json->{'parent_pending'} ) ) ? $upd_json->{'parent_pending'} : undef;
			my $parent_reval_pending = ( defined( $upd_json->{'parent_reval_pending'} ) ) ? $upd_json->{'parent_reval_pending'} : undef;
			if ( !defined($parent_pending) ) {
				( $log_level >> $ERROR ) && print "ERROR Update URL: did not have an parent_pending key.\n";
				if ( $script_mode != $SYNCDS ) {
					return $syncds_update;
				}
				else {
					( $log_level >> $ERROR ) && print "ERROR Invalid JSON for update_status. Exiting, not sure what else to do.\n";
					exit 1;
				}
			}
			if ( $parent_pending == 1 && $wait_for_parents == 1 && $reval_in_use == 0) {
				( $log_level >> $ERROR ) && print "ERROR Traffic Ops is signaling that my parents need an update.\n";
				if ( $script_mode == $SYNCDS ) {
					if ( $dispersion > 0 ) {
						( $log_level >> $WARN ) && print "WARN In syncds mode, sleeping for " . $dispersion . "s to see if the update my parents need is cleared.\n";
						( $dispersion > 0 ) && &sleep_timer($dispersion);
					}
					$upd_json = get_update_status();

					$parent_pending = ( defined( $upd_json->{'parent_pending'} ) ) ? $upd_json->{'parent_pending'} : undef;
					if ( !defined($parent_pending) ) {
						( $log_level >> $ERROR ) && print "ERROR Invalid JSON for update_status. Exiting, not sure what else to do.\n";
					}
					if ( $parent_pending == 1 || $parent_reval_pending == 1 ) {
						( $log_level >> $ERROR ) && print "ERROR My parents still need an update, bailing.\n";
						exit 1;

					}
					else {
						( $log_level >> $DEBUG ) && print "DEBUG The update on my parents cleared; continuing.\n";
					}
				}
			}
			else {
				( $log_level >> $DEBUG ) && print "DEBUG Traffic Ops is signaling that my parents do not need an update, or wait_for_parents == 0.\n";
			}
		}
		elsif ( $script_mode == $SYNCDS && $upd_pending != 1 ) {
			( $log_level >> $ERROR ) && print "ERROR In syncds mode, but no syncds update needs to be applied. Running revalidation before exiting.\n";
			&revalidate_while_sleeping();
			exit 0;
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR Traffic Ops is signaling that no update is waiting to be applied.\n";
		}

		my $stj = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --get-data=statuses 2>>$atstccfg_log_path`;
		my $atstccfg_exit_code = $?;
		if ( $atstccfg_exit_code != 0 ) {
			( $log_level >> $ERROR ) && print "Statuses URL: returned $stj! Skipping creation of status file.\n";
		}

		my $statuses = decode_json($stj);
		my $my_status = ( defined( $upd_json->{'status'} ) ) ? $upd_json->{'status'} : undef;

		if ( defined($my_status) ) {
			( $log_level >> $DEBUG ) && print "DEBUG Found $my_status status from Traffic Ops.\n";
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR Returning; did not find status from Traffic Ops!\n";
			return ($syncds_update);
		}

		my $status_dir  = dirname($0) . "/status";
		my $status_file = $status_dir . "/" . $my_status;

		if ( !-f $status_file ) {
			( $log_level >> $ERROR ) && print "ERROR status file $status_file does not exist.\n";
		}

		for my $status ( @{$statuses} ) {
			next if ( $status->{name} eq $my_status );
			my $other_status = $status_dir . "/" . $status->{name};

			if ( -f $other_status && $status->{name} ne $my_status ) {
				( $log_level >> $ERROR ) && print "ERROR Other status file $other_status exists.\n";
				if ( $script_mode != $REPORT ) {
					( $log_level >> $DEBUG ) && print "DEBUG Removing $other_status\n";
					unlink($other_status);
				}
			}
		}

		if ( $script_mode != $REPORT ) {
			if ( !-d $status_dir ) {
				mkpath($status_dir);
			}

			if ( !-f $status_file ) {
				my $r = open( FH, "> $status_file" );

				if ( !$r ) {
					( $log_level >> $ERROR ) && print "ERROR Unable to touch $status_file\n";
				}
				else {
					close(FH);
				}
			}
		}
	}
	return ($syncds_update);
}

sub sleep_rand {
	my $duration = int( rand(shift) );

	( $log_level >> $WARN ) && print "WARN Sleeping for $duration seconds: ";

	for ( my $i = $duration; $i > 0; $i-- ) {
		( $log_level >> $WARN ) && print ".";
		sleep 1;
	}
	( $log_level >> $WARN ) && print "\n";
}

sub sleep_timer {
	my $duration = shift;
	my $reval_clock = $reval_wait_time;

	my $proper_script_mode = $script_mode;

	if ( $reval_in_use == 1 && $proper_script_mode != $BADASS ) {
		( $log_level >> $WARN ) && print "WARN Performing a revalidation check before sleeping... \n";
		&revalidate_while_sleeping();
		( $log_level >> $WARN ) && print "WARN Revalidation check complete.\n";
	}
	if ( $duration < $reval_clock || $reval_in_use == 0 || $script_mode == $BADASS ) {
		( $log_level >> $WARN ) && print "WARN Sleeping for $duration seconds: ";
	}
	else {
		( $log_level >> $WARN ) && print "WARN $reval_clock seconds until next revalidation check.\n";
		( $log_level >> $WARN ) && print "WARN $duration seconds remaining in dispersion sleep period\n";
		( $log_level >> $WARN ) && print "WARN Sleeping for $reval_clock seconds: ";
	}
	for ( my $i = $duration; $i > 0; $i-- ) {
		( $log_level >> $WARN ) && print ".";
		sleep 1;
		$reval_clock--;
		if ($reval_clock < 1 && $script_mode != $BADASS && $reval_in_use == 1 ) {
			( $log_level >> $WARN ) && print "\n";
			( $log_level >> $WARN ) && print "WARN Interrupting dispersion sleep period for revalidation check. \n";
			&revalidate_while_sleeping();
			$reval_clock = $reval_wait_time;
			if ($reval_clock < $i ) {
				( $log_level >> $WARN ) && print "WARN Revalidation check complete. $reval_clock seconds until next revalidation check.\n";
				( $log_level >> $WARN ) && print "WARN $i seconds remaining in dispersion sleep period\n";
				( $log_level >> $WARN ) && print "WARN Sleeping for $reval_clock seconds: ";
			}
			else {
				( $log_level >> $WARN ) && print "WARN Revalidation check complete. $i seconds remaining in dispersion sleep period.\n";
				( $log_level >> $WARN ) && print "WARN Sleeping for $i seconds: ";
			}
		}
	}

	$script_mode = $proper_script_mode;

	( $log_level >> $WARN ) && print "\n";
}

sub process_config_files{
	( $log_level >> $INFO ) && print "\nINFO: ======== Start processing config files ========\n";

	foreach my $file ( keys %{$cfg_file_tracker} ) {
		( $log_level >> $DEBUG ) && print "DEBUG Starting processing of config file: $file\n";
		my $cfg_dir = dirname($cfg_file_tracker->{$file}->{'headers'}->{'path'});
		my $return = undef;
		if (
			$script_mode == $SYNCDS
			&& (   $file eq "records.config"
				|| $file eq "remap.config"
				|| $file eq "parent.config"
				|| $file eq "cache.config"
				|| $file eq "hosting.config"
				|| $file =~ m/url\_sig\_(.*)\.config$/
				|| $file =~ m/uri\_signing\_(.*)\.config$/
				|| $file =~ m/hdr\_rw\_(.*)\.config$/
				|| $file eq "regex_revalidate.config"
				|| $file eq "astats.config"
				|| $file =~ m/cacheurl\_(.*)\.config$/
				|| $file =~ m/regex\_remap\_(.*)\.config$/
				|| $file =~ m/\.cer$/
				|| $file =~ m/\.key$/
				|| $file eq "logs_xml.config"
				|| $file eq "logging.config"
				|| $file eq "logging.yaml"
				|| $file eq "ssl_multicert.config" )
			)
		{
			if ( package_installed("trafficserver") ) {
				( $log_level >> $DEBUG ) && print "DEBUG In syncds mode, I'm about to process config file: $file\n";
				$cfg_file_tracker->{$file}->{'service'} = "trafficserver";
				$return = &process_cfg_file($file);
			}
			else {
				( $log_level >> $FATAL ) && print "FATAL In syncds mode, but trafficserver isn't installed. Bailing.\n";
				exit 1;
			}
		}
		elsif ($script_mode == $SYNCDS
			&& $file =~ m/\_facts/
			&& $cfg_dir =~ m/\/opt\/ort/ )
		{
			( $log_level >> $DEBUG ) && print "DEBUG In syncds mode, I'm about to process config file: $file\n";
			$cfg_file_tracker->{$file}->{'service'} = "puppet";
			$return = &process_cfg_file($file);
		}
		elsif ( $script_mode == $SYNCDS && $cfg_dir =~ m/cron/ ) {
			( $log_level >> $DEBUG ) && print "DEBUG In syncds mode, I'm about to process config file: $file\n";
			$cfg_file_tracker->{$file}->{'service'} = "system";
			$return = &process_cfg_file($file);
		}
		elsif ( $script_mode == $SYNCDS ) {
			( $log_level >> $DEBUG ) && print "DEBUG In syncds mode, I'm about to process config file: $file\n";
			$cfg_file_tracker->{$file}->{'service'} = "trafficserver";
			$return = &process_cfg_file($file);
		}
		elsif ( $script_mode != $SYNCDS ) {
			if (
				package_installed("trafficserver")
				&& ( $cfg_dir =~ m/trafficserver/ || $cfg_dir =~ m/udev/ )
				)
			{
				$cfg_file_tracker->{$file}->{'service'} = "trafficserver";
				$return = &process_cfg_file($file);
			}
			elsif ( $file eq "sysctl.conf" || $file eq "50-ats.rules" || $file =~ m/cron/ ) {
				$cfg_file_tracker->{$file}->{'service'} = "system";
				$return = &process_cfg_file($file);
			}
			elsif ( $file =~ m/\_facts/ ) {
				$cfg_file_tracker->{$file}->{'service'} = "puppet";
				$return = &process_cfg_file($file);
			}
			else {
				( $log_level >> $WARN ) && print "WARN $file is being processed with an unknown service\n";
				$cfg_file_tracker->{$file}->{'service'} = "unknown";
				$return = &process_cfg_file($file);
			}
		}
		if ( defined($return) && $return == $CFG_FILE_PREREQ_FAILED ) {
			$syncds_update = $UPDATE_TROPS_FAILED;
		}
	}

	foreach my $file ( keys %{$cfg_file_tracker} ) {
		if (   $cfg_file_tracker->{$file}->{'change_needed'}
			&& !$cfg_file_tracker->{$file}->{'change_applied'}
			&& $cfg_file_tracker->{$file}->{'audit_complete'}
			&& !$cfg_file_tracker->{$file}->{'prereq_failed'}
			&& !$cfg_file_tracker->{$file}->{'audit_failed'} )
		{
			if ( $file eq "plugin.config" && $cfg_file_tracker->{'remap.config'}->{'prereq_failed'} ) {
				( $log_level >> $ERROR )
					&& print "ERROR plugin.config changed. However, prereqs failed for remap.config so I am skipping updates for plugin.config.\n";
				next;
			}
			elsif ( $file eq "remap.config" && $cfg_file_tracker->{'plugin.config'}->{'prereq_failed'} ) {
				( $log_level >> $ERROR )
					&& print "ERROR remap.config changed. However, prereqs failed for plugin.config so I am skipping updates for remap.config.\n";
				next;
			}
			else {
				( $log_level >> $DEBUG ) && print "DEBUG Prereqs passed for replacing $file on disk with that in Traffic Ops.\n";
				&replace_cfg_file($file);
			}
		}
	}
	( $log_level >> $INFO ) && print "\nINFO: ======== End processing config files ========\n\n";
}

sub touch_file {
	my $return = 0;
	my $file   = shift;
	$file = $cfg_file_tracker->{$file}->{'headers'}->{'path'};
	( $log_level >> $DEBUG ) && print "DEBUG About to touch $file.\n";

	if ( $script_mode == $INTERACTIVE ) {
		( $log_level >> $ERROR ) && print "ERROR $file needs touched. Should I do that now? [Y/n] (n): ";
		my $select = 'n';
		$select = <STDIN>;
		chomp($select);
		if ( $select =~ m/Y/ ) {
			$return = &touch_this_file($file);
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR $file was not touched.\n";
		}
	}
	elsif ( $script_mode == $BADASS || $script_mode == $SYNCDS ) {
		( $log_level >> $ERROR ) && print "ERROR $file needs touched. Doing that now.\n";
		$return = &touch_this_file($file);
	}
	return $return;
}

sub touch_this_file {
	my $file    = shift;
	my $result  = `/bin/touch $file 2>&1`;
	my $success = 0;
	chomp($result);
	if ( $result =~ m/cannot touch/ || $result =~ m/Permission denied/ || $result =~ m/No such file or directory/ ) {
		( $log_level >> $ERROR ) && print "ERROR $file was not touched successfully. Error: $result.\n";
		$success = 0;
	}
	else {
		( $log_level >> $DEBUG ) && print "DEBUG $file was touched successfully.\n";
		$success++;
	}
	return $success;
}

sub run_traffic_ctl {
	my $output = `$TRAFFIC_CTL config reload 2>&1`;
	if ( $output !~ m/error/ ) {
		( $log_level >> $DEBUG ) && print "DEBUG traffic_ctl run successful.\n";
		if ( $syncds_update == $UPDATE_TROPS_NEEDED ) {
			$syncds_update = $UPDATE_TROPS_SUCCESSFUL;
		}
	}
	else {
		if ( $syncds_update == $UPDATE_TROPS_NEEDED ) {
			( $log_level >> $ERROR ) && print "ERROR traffic_ctl run failed. Updating Traffic Ops anyway.\n";
			$syncds_update = $UPDATE_TROPS_SUCCESSFUL;
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR traffic_ctl run failed.\n";
		}
	}
}

sub check_plugins {
	my $cfg_file       = shift;
	my $file_lines_ref = shift;
	my @file_lines     = @{$file_lines_ref};
	my $return_code    = 0;

	( $log_level >> $DEBUG ) && print "DEBUG Checking plugins for $cfg_file\n";

	if ( $cfg_file eq "plugin.config" ) {
		( $log_level >> $DEBUG ) && print "DEBUG Entering advanced processing for plugin.config.\n";
		foreach my $linep (@file_lines) {
			if ( $linep =~ m/^\#/ ) { next; }
			( my $plugin_name ) = split( /\s+/, $linep );
			$plugin_name =~ s/\s+//g;
			( $log_level >> $DEBUG ) && print "DEBUG Found plugin $plugin_name in $cfg_file.\n";

			my $return_code = &check_this_plugin($plugin_name);

			if ( $return_code == $PLUGIN_YES ) {
				( $log_level >> $DEBUG ) && print "DEBUG Package for plugin: $plugin_name is installed.\n";
			}
			elsif ( $return_code == $PLUGIN_NO ) {
				( $log_level >> $ERROR ) && print "ERROR Package for plugin: $plugin_name is not installed!\n";
				$cfg_file_tracker->{$cfg_file}->{'prereq_failed'}++;
			}
		}
	}
	elsif ( $cfg_file eq "remap.config" ) {
		( $log_level >> $DEBUG ) && print "DEBUG Entering advanced processing for remap.config\n";
		foreach my $liner (@file_lines) {
			if ( $liner =~ m/^\#/ ) { next; }
			( my @parts ) = split( /\@plugin\=/, $liner );
			foreach my $i ( 1..$#parts ) {
				( my $plugin_name, my $plugin_config_file ) = split( /\@pparam\=/, $parts[$i] );
				if (defined( $plugin_config_file ) ) {
					# Subblock for lasting out of.
					{
						($plugin_config_file) = split( /\s+/, $plugin_config_file);

						# Skip parameters that start with '-' or 'proxy.config.', since those are probabably parameters, not config files.
						last if $plugin_config_file =~ m/^-/; # Exit subblock.
						last if $plugin_config_file =~ m/^proxy.config./;

						( my @parts ) = split( /\//, $plugin_config_file );
						$plugin_config_file = $parts[$#parts];
						$plugin_config_file =~ s/\s+//g;
						if ( !exists($cfg_file_tracker->{$plugin_config_file}->{'remap_plugin_config_file'} ) && $plugin_config_file !~ /.lua$/ ) {
							$cfg_file_tracker->{$plugin_config_file}->{'remap_plugin_config_file'} = 1;
						}
					}
				}
				else {
					($plugin_name) = split(/\s/, $plugin_name);
				}
				$plugin_name =~ s/\s//g;
				( $log_level >> $DEBUG ) && print "DEBUG Found plugin $plugin_name in $cfg_file.\n";

				$return_code = &check_this_plugin($plugin_name);

				if ( $return_code == $PLUGIN_YES ) {
					( $log_level >> $DEBUG ) && print "DEBUG Package for plugin: $plugin_name is installed.\n";
				}
				elsif ( $return_code == $PLUGIN_NO ) {
					( $log_level >> $ERROR ) && print "ERROR Package for plugin: $plugin_name is not installed\n";
					$cfg_file_tracker->{$cfg_file}->{'prereq_failed'}++;
				}
			}
		}
	}
	( $log_level >> $TRACE ) && print "TRACE Returning $return_code for checking plugins for $cfg_file.\n";

	return $return_code;
}

my %checked_plugins = ();
sub check_this_plugin {
	my $plugin      = shift;

	if ( exists( $checked_plugins{$plugin} ) ) {
		return ($checked_plugins{$plugin});
	}

	my $full_plugin = $TS_HOME . "/libexec/trafficserver/" . $plugin;
	( $log_level >> $DEBUG ) && print "DEBUG Checking package dependency for plugin: $plugin.\n";

	my $provided = package_provides($full_plugin);

	if ($provided) {
		if ( package_was_installed($provided) ) {
			$trafficserver_restart_needed++;
		}

		$checked_plugins{$plugin} = $PLUGIN_YES;
		return ($PLUGIN_YES);
	}
	else {
		$checked_plugins{$plugin} = $PLUGIN_NO;
		return ($PLUGIN_NO);
	}
}

sub replace_cfg_file {
	my $cfg_file    = shift;
	my $return_code = 0;
	my $select      = 2;

	if ( $script_mode == $INTERACTIVE ) {
		( $log_level >> $ERROR )
			&& print
			"ERROR $cfg_file on disk needs updated with one from Traffic Ops. [1] override files on disk with data in Traffic Ops, [2] ignore and continue. (2): ";
		my $input = <STDIN>;
		chomp($input);
		if ( $input =~ m/\d/ ) {
			$select = $input;
		}
	}
	if ( $select == 1 || $script_mode == $BADASS || $script_mode == $SYNCDS || $script_mode == $REVALIDATE ) {
		( $log_level >> $ERROR )
			&& print "ERROR Copying "
			. $cfg_file_tracker->{$cfg_file}->{'backup_from_trops'} . " to "
			. $cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'}
			. "\n";
		system("/bin/cp $cfg_file_tracker->{$cfg_file}->{'backup_from_trops'} $cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'}");
		if ( $cfg_file =~ /cron/ ) {
			chown 0, 0, "$cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'}";
		}
		else {
			my $ats_uid  = getpwnam("ats");
			chown $ats_uid, $ats_uid, "$cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'}";
		}
		$cfg_file_tracker->{$cfg_file}->{'change_applied'}++;
		( $log_level >> $TRACE ) && print "TRACE Setting change applied for $cfg_file.\n";
		$return_code = $CFG_FILE_CHANGED;
		&process_reload_restarts($cfg_file);
	}
	elsif ( $select == 2 && $script_mode != $REPORT ) {
		( $log_level >> $ERROR ) && print "ERROR You elected not to replace $cfg_file with version from Traffic Ops.\n";
		$cfg_file_tracker->{$cfg_file}->{'change_applied'} = 0;
		$return_code = $CFG_FILE_UNCHANGED;
	}
	else {
		$cfg_file_tracker->{$cfg_file}->{'change_applied'} = 0;
		$return_code = $CFG_FILE_UNCHANGED;
	}
	return $return_code;
}

sub process_reload_restarts {

	my $cfg_file = shift;
	my $cfg_dir = dirname($cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'});
	( $log_level >> $DEBUG ) && print "DEBUG Applying config for: $cfg_file.\n";

	if ( $cfg_file =~ m/url\_sig\_(.*)\.config/ ) {
		( $log_level >> $DEBUG ) && print "DEBUG New keys were installed in: $cfg_file, touch remap.config, and traffic_ctl config reload needed.\n";
		$traffic_ctl_needed++;
	}
	elsif ( $cfg_file =~ m/uri\_signing\_(.*)\.config/ ) {
		( $log_level >> $DEBUG ) && print "DEBUG New keys were installed in: $cfg_file, touch remap.config, and traffic_ctl config reload needed.\n";
		$traffic_ctl_needed++;
	}
	elsif ( $cfg_file =~ m/hdr\_rw\_(.*)\.config/ ) {
		( $log_level >> $DEBUG ) && print "DEBUG New/changed header rewrite rule, installed in: $cfg_file. Later I will attempt to touch remap.config.\n";
		$traffic_ctl_needed++;
	}
	elsif ( $cfg_file eq "plugin.config" || $cfg_file eq "50-ats.rules" ) {
		( $log_level >> $DEBUG ) && print "DEBUG $cfg_file changed, trafficserver restart needed.\n";
		$trafficserver_restart_needed++;
	}
	elsif ( $cfg_dir =~ m/ssl/ && ( $cfg_file =~ m/\.cer$/ || $cfg_file =~ m/\.key$/ ) ) {
		( $log_level >> $DEBUG ) && print "DEBUG SSL key/cert $cfg_file changed, touch ssl_multicert.config, and traffic_ctl config reload needed.\n";
		$installed_new_ssl_keys++;
		$traffic_ctl_needed++;
	}
	elsif ( $cfg_dir =~ m/trafficserver/ ) {
		( $log_level >> $DEBUG ) && print "DEBUG $cfg_file changed, traffic_ctl config reload needed.\n";
		$traffic_ctl_needed++;
	}
	elsif ( $cfg_file eq "sysctl.conf" ) {
		( $log_level >> $DEBUG ) && print "DEBUG $cfg_file changed, 'sysctl -p' needed.\n";
		$sysctl_p_needed++;
	}
	elsif ( $cfg_file =~ m/\_facts/ ) {
		( $log_level >> $DEBUG ) && print "DEBUG Puppet facts file $cfg_file changed.\n";
		$UPDATE_TROPS_SUCCESSFUL = 1;
	}
	elsif ( $cfg_file =~ m/cron/ ) {
		( $log_level >> $DEBUG ) && print "DEBUG Cron file $cfg_file changed.\n";
		$UPDATE_TROPS_SUCCESSFUL = 1;
	}
}

sub check_output {
	my $out = shift;
	if ( defined($out) ) {
		$out =~ s/(\n+|\t+|\r+|\s+)/ /g;
		if ( $out =~ m/error/i ) {
			( $log_level >> $ERROR ) && print "ERROR $out\n";
			return 1;
		}
		else {
			return 0;
		}
	}
	else {
		return 1;
	}
}

sub check_script_mode {
	#### No default script mode
	my $script_mode = undef;
	if ( $ARGV[0] eq "interactive" ) {
		( $log_level >> $DEBUG ) && print "DEBUG Script running in interactive mode.\n";
		$script_mode = $INTERACTIVE;
	}
	elsif ( $ARGV[0] eq "report" ) {
		( $log_level >> $DEBUG ) && print "DEBUG Script running in report mode.\n";
		$script_mode = $REPORT;
	}
	elsif ( $ARGV[0] eq "badass" ) {
		( $log_level >> $DEBUG ) && print "DEBUG Script running in badass mode.\n";
		$script_mode = $BADASS;
	}
	elsif ( $ARGV[0] eq "syncds" ) {
		( $log_level >> $DEBUG ) && print "DEBUG Script running in syncds mode.\n";
		$script_mode = $SYNCDS;
	}
	elsif ( $ARGV[0] eq "revalidate" ) {
		( $log_level >> $DEBUG ) && print "DEBUG Script running in revalidate mode.\n";
		$script_mode = $REVALIDATE;
	}
	else {
		( $log_level >> $FATAL ) && print "FATAL You did not specify a valid mode. Exiting.\n";
		&usage();
		exit 1;
	}
	return $script_mode;

}

sub check_run_user {
	my $run_user = `/usr/bin/id`;
	chomp($run_user);
	if (   ( $run_user !~ m/uid\=0\(root\)/ && $run_user !~ m/gid\=0\(root\)/ && $run_user !~ m/groups\=0\(root\)/ )
		&& ( $script_mode == $INTERACTIVE || $script_mode == $BADASS || $script_mode == $SYNCDS ) )
	{
		( $log_level >> $FATAL ) && print "FATAL For interactive, badass, or syncds mode, you must run script as root user. Exiting.\n";
		exit 1;
	}
	else {
		( $log_level >> $TRACE ) && print "TRACE run user is $run_user.\n";
	}
}

sub check_log_level {
	if ( ( $script_mode == $INTERACTIVE ) && !( $log_level >> $ERROR ) ) {
		print "FATAL Sorry, for interactive mode, the log level must be at least ERROR, exiting.\n";
		exit 1;
	}
}

# get_cfg_file_list gets the config files from atstccfg via Traffic Ops.
# See parse_multipart_config_files for the return type.
sub get_cfg_file_list {
	my $host_name = shift;
	my $script_mode = shift;

	my $atstccfg_reval_arg = '';
	if ( $script_mode == $REVALIDATE ) {
		$atstccfg_reval_arg = '--revalidate-only';
	}

	my $result = `$atstccfg_cmd --dir='$ats_config_dir' $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$host_name' $atstccfg_reval_arg --log-location-error=stderr --log-location-warning=stderr --log-location-info=null 2>>$atstccfg_log_path`;
	my $atstccfg_exit_code = $?;
	if ($atstccfg_exit_code != 0) {
		( $log_level >> $ERROR ) && printf("ERROR getting config files from atstccfg via Traffic Ops. See $atstccfg_log_path for details\n");
		exit 1;
	}

	return &parse_multipart_config_files($result);
}

# parse_multipart_config_files parses the multipart/mixed message returned by atstccfg, and returns a map of file names to content.
# Returns a hash of the file name to a sub-hash. Each sub-hash contains the values 'headers' and 'body'.
#  The 'headers' is a hash of lowercased header names to header values.
#  The 'body' is the string body
sub parse_multipart_config_files {
	my $multipart_txt = shift;

	# Note this doesn't enforce RFC compliance for the boundary. We assume any char is valid, where in reality only certain characters are allowed in the unquoted value. See RFC2616s3.6
	my $boundary = '';
	if ($multipart_txt =~ m/boundary="/) {
		($boundary) = $multipart_txt =~ m/boundary="((?:[^"\\]|\\.)*)"/g; # this regex gets the quoted boundary="foo" value
	} else {
		($boundary) = $multipart_txt =~ m/boundary=([^ \r\n\t]+)/g; # this regex gets the unquoted boundary=foo value
	}

	if ( length($boundary) == 0 ) {
		( $log_level >> $FATAL ) && print "FATAL Error getting package list from Traffic Ops! Could not find boundary for multipart message\n";
		exit 1;
	}

	my $last_boundary = "--" . $boundary . "--"; # multipart ends in --boundary--
	$multipart_txt =~ s/$last_boundary//;       # remove it

	my @files = split("--" . $boundary . "\r\n", $multipart_txt);

	my %all_files;
	for my $i (1 .. $#files) { # start at 1, because 0 is the MIME-Version and Content-Type: multipart/mixed
		my $file = $files[$i];
		my @headers_body = split("\r\n\r\n", $file);

		if ( @headers_body < 2 ) {
			print "FATAL Error getting package list from Traffic Ops! malformed headers on file $i\n";
			exit 1;
		}

		my @headers_arr = split("\r\n", $headers_body[0]);
		my %headers;
		for my $i (0 .. $#headers_arr) {
				my $header = $headers_arr[$i];
				my @header_name_val = split(": ", $header, 2);
				my $header_name = lc($header_name_val[0]);
				my $header_val = $header_name_val[1];
				$headers{$header_name} = $header_val;
		}

		my $file_name = basename($headers{'path'});
		if ( length($file_name) == 0 ) {
			print "FATAL Error getting package list from Traffic Ops! Headers on file $i missing Path!\n";
			exit 1;
		}

		my %file_obj;
		$file_obj{'headers'} = \%headers;
		my $body = $headers_body[1];
		$body =~ s/\s+$//; # trim trailing whitespace
		$body = $body . "\n"; # add a single trailing newline, POSIX files require trailing newlines
		$file_obj{'body'} = $body;
		$all_files{$file_name} = \%file_obj;
	}

	return \%all_files;
}

sub get_header_comment {
	my $toolname;

	my $result = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --get-data=system-info 2>>$atstccfg_log_path`;
	my $atstccfg_exit_code = $?;
	if ($atstccfg_exit_code != 0) {
			( $log_level >> $ERROR ) && printf("ERROR Unable to get system info. Stopping.\n");
			exit 1;
	}
	my $result_ref = decode_json($result);
	if ( defined( $result_ref->{'tm.toolname'} ) ) {
		$toolname = $result_ref->{'tm.toolname'};
		( $log_level >> $INFO ) && printf("INFO Found tm.toolname: $toolname\n");
	}
	else {
		print "ERROR Did not find tm.toolaname!\n";
		$toolname = "";
	}
	return $toolname;
}

sub __package_action {
	my $action        = shift;
	my @argument_list = @_;

	my $arguments   = join( " ", @argument_list );
	my $yum_command = "/usr/bin/yum $YUM_OPTS $action $arguments 2>&1";
	my $out         = `$yum_command`;

	# yum exits 0 if successful
	if ( $? != 0 ) {
		( $log_level >> $ERROR ) && print "ERROR Execution of $yum_command failed!\n";
		( $log_level >> $ERROR ) && print "ERROR Output: $out\n";

		return (0);
	}
	else {
		( $log_level >> $TRACE ) && print "TRACE Successfully executed $yum_command\n";

		#($log_level >> $DEBUG) && print "DEBUG Output: $out\n";

		return (1);
	}
}

sub get_full_package_name {
	my $package = shift;
	my $version = shift;
	return ( $package . "-" . $version );
}

sub package_provides {
	my $filename = shift || die("Please supply the full path to the file to verify");

	my $out = `/bin/rpm -qf $filename 2>&1`;

	if ( defined($out) ) {
		chomp($out);
	}

	if ( $? == 0 ) {

		# return package name that provides $filename
		return ($out);
	}
	else {
		return (0);
	}
}

sub package_requires {
	my $package_name = shift;
	my @package_list = ();

	my $out = `/bin/rpm -q --whatrequires $package_name 2>&1`;

	if ( defined($out) ) {
		chomp($out);
	}

	if ( $? == 0 ) {
		@package_list = split( /\n/, $out );
	}

	return (@package_list);
}

sub package_was_installed {
	my $package_name = shift;

	if ( exists( $install_tracker{$package_name} ) ) {
		( $log_level >> $TRACE ) && print "TRACE $package_name was installed during this run, returning true\n";
		return (1);
	}
	else {
		( $log_level >> $TRACE ) && print "TRACE $package_name was not installed during this run, returning false\n";
		return (0);
	}
}

my %packages_installed = ();
sub package_installed {
	my $package_name    = shift;
	my $package_version = shift;
	my @package_list    = ();

	if ( defined($package_version) ) {
		$package_name = $package_name . "-" . $package_version;
	}

	if ( exists( $packages_installed{$package_name} ) ) {
		return ($packages_installed{$package_name});
	}

	my $out = `/bin/rpm -q $package_name 2>&1`;

	# rpm returns 0 if installed, 1 if not installed
	if ( $? == 0 ) {

		# installed
		# remove the newlines (hence not using an array for $out)
		@package_list = split( /\n/, $out );
	}

	$packages_installed{$package_name} = @package_list;
	return (@package_list);
}

sub packages_available {
	my @package_list    = @_;
	my $package_missing = 0;

	for my $package (@package_list) {
		my $result = __package_action( "info", $package );

		if ($result) {
			( $log_level >> $TRACE ) && print "TRACE $package is available\n";
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR $package is not available in the yum repo(s)!\n";
			$package_missing = 1;
		}
	}

	if ($package_missing) {
		return (0);
	}
	else {
		return (1);
	}
}

sub install_packages {
	my @package_list = @_;

	if ( __package_action( "install", "-y", @package_list ) ) {
		for my $pkg (@package_list) {
			$install_tracker{$pkg} = 1;
		}

		return (1);
	}
	else {
		return (0);
	}
}

sub remove_packages {
	my @package_list = @_;

	return ( __package_action( "remove", "-y", @package_list ) );
}

sub process_packages {
	my $host_name = shift;

	my $proceed = 0;

	my $result = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --get-data=packages 2>>$atstccfg_log_path`;
	my $atstccfg_exit_code = $?;
	if ($atstccfg_exit_code != 0) {
		( $log_level >> $FATAL ) && print "FATAL Error getting package list from Traffic Ops!\n";
			exit 1;
	}

		my %package_map;
		my @package_list = @{ decode_json($result) };

		# iterate through to build the uninstall list
		for my $package (@package_list) {
			my $full_package = get_full_package_name( $package->{"name"}, $package->{"version"} );

			# check to see if any package is installed that has this package's basename (no version)
			for my $installed_package ( package_installed( $package->{name} ) ) {
				if ( exists( $package_map{"uninstall"}{$full_package} ) ) {
					( $log_level >> $INFO ) && print "INFO $full_package: Already marked for removal.\n";
					next;
				}
				elsif ( $installed_package eq $full_package ) {

					# skip this package if it's the correct version
					( $log_level >> $INFO ) && print "INFO $full_package: Currently installed and not marked for removal.\n";
					next;
				}

				if ( $script_mode == $REPORT ) {
					( $log_level >> $FATAL ) && print "ERROR $installed_package: Currently installed and needs to be removed.\n";
				}
				else {
					( $log_level >> $TRACE ) && print "TRACE $installed_package: Currently installed, marked for removal.\n";
				}

				$package_map{"uninstall"}{$installed_package} = 1;

				# add any dependent packages to the list of packages to uninstall
				for my $dependent_package ( package_requires( $package->{name} ) ) {
					if ( $script_mode == $REPORT ) {
						( $log_level >> $FATAL )
							&& print "ERROR $dependent_package: Currently installed and depends on " . $package->{name} . "and needs to be removed.\n";
					}
					else {
						( $log_level >> $TRACE )
							&& print "TRACE $dependent_package: Currently installed and depends on " . $package->{name} . ", marked for removal.\n";
					}

					$package_map{"uninstall"}{$dependent_package} = 1;
				}
			}
		}

		# iterate through to build the install list
		for my $package (@package_list) {
			my $full_package = get_full_package_name( $package->{"name"}, $package->{"version"} );
			if ( !package_installed( $package->{name}, $package->{version} ) ) {
				if ( $script_mode == $REPORT ) {
					( $log_level >> $FATAL ) && print "ERROR $full_package: Needs to be installed.\n";
				}
				else {
					( $log_level >> $TRACE ) && print "TRACE $full_package: Needs to be installed.\n";
				}

				$package_map{"install"}{$full_package} = 1;
			}
			elsif ( exists( $package_map{"uninstall"}{$full_package} ) ) {
				if ( $script_mode == $REPORT ) {
					( $log_level >> $FATAL ) && print "ERROR $full_package: Marked for removal and needs to be installed.\n";
				}
				else {
					( $log_level >> $TRACE ) && print "TRACE $full_package: Marked for removal and needs to be installed.\n";
				}

				$package_map{"install"}{$full_package} = 1;
			}
			else {
				# if the correct version is already installed not marked for removal we don't want to do anything..
				if ( $script_mode == $REPORT ) {
					( $log_level >> $INFO ) && print "INFO $full_package: Currently installed and not marked for removal.\n";
				}
				else {
					( $log_level >> $TRACE ) && print "TRACE $full_package: Currently installed and not marked for removal.\n";
				}
			}
		}

		my @install_packages   = keys( %{ $package_map{"install"} } );
		my @uninstall_packages = keys( %{ $package_map{"uninstall"} } );

		if ( scalar(@install_packages) > 0 || scalar(@uninstall_packages) > 0 ) {

			if ( packages_available(@install_packages) ) {
				my $uninstalled = ( scalar(@uninstall_packages) > 0 ) ? 0 : 1;
				( $log_level >> $TRACE ) && print "TRACE All packages available.. proceeding..\n";

				if ( $script_mode == $BADASS ) {
					$proceed = 1;
				}
				elsif ( $script_mode == $INTERACTIVE && scalar(@uninstall_packages) > 0 ) {
					( $log_level >> $INFO )
						&& print "INFO The following packages must be uninstalled before proceeding:\n  - " . join( "\n  - ", @uninstall_packages ) . "\n";
					if ( get_answer("Should I uninstall them now?") && get_answer("Are you sure you want to proceed with the uninstallation?") ) {
						$proceed = 1;
					}
					else {
						$proceed = 0;
					}
				}

				if ( $proceed && scalar(@uninstall_packages) > 0 ) {
					if ( remove_packages(@uninstall_packages) ) {
						( $log_level >> $INFO )
							&& print "INFO Successfully uninstalled the following packages:\n  - " . join( "\n  - ", @uninstall_packages ) . "\n";
						$uninstalled = 1;
					}
					else {
						( $log_level >> $ERROR )
							&& print "ERROR Unable to uninstall the following packages:\n  - " . join( "\n  - ", @uninstall_packages ) . "\n";
						$proceed = 0;
					}
				}

				if ( $uninstalled && $script_mode == $INTERACTIVE && scalar(@install_packages) > 0 ) {
					( $log_level >> $INFO ) && print "INFO The following packages must be installed:\n  - " . join( "\n  - ", @install_packages ) . "\n";
					if ( get_answer("Should I install them now?") && get_answer("Are you sure you want to proceed with the installation?") ) {
						$proceed = 1;
					}
					else {
						$proceed = 0;
					}
				}

				if ( $uninstalled && $proceed && scalar(@install_packages) > 0 ) {
					if ( install_packages(@install_packages) ) {
						( $log_level >> $INFO )
							&& print "INFO Successfully installed the following packages:\n  - " . join( "\n  - ", @install_packages ) . "\n";
						$syncds_update = $UPDATE_TROPS_SUCCESSFUL;
					}
					else {
						( $log_level >> $ERROR )
							&& print "ERROR Unable to install the following packages:\n  - " . join( "\n  - ", @install_packages ) . "\n";
					}
				}
				elsif ( scalar(@install_packages) == 0 ) {
					( $log_level >> $INFO ) && print "INFO All of the required packages are installed.\n";
				}
			}
			else {
				( $log_level >> $ERROR ) && print "ERROR Not all of the required packages are available in the configured yum repo(s)!\n";
			}
		}
		else {
			if ( $script_mode == $REPORT ) {
				( $log_level >> $INFO ) && print "INFO All required packages are installed.\n";
			}
			else {
				( $log_level >> $TRACE ) && print "TRACE All required packages are installed.\n";
			}
		}
}

sub set_chkconfig {
	my $service   = shift;
	my $run_level = shift;
	my $setting   = shift;

	if ( !defined($service) || !defined($run_level) || !defined($setting) ) {
		die("Please supply a service, run level (0-6) and setting, in that order");
	}
	elsif ( $run_level !~ m/^[0-6]$/ ) {
		die("Please supply a numeric run level (0-6)");
	}

	my $command = "/sbin/chkconfig --level $run_level $service $setting";
	my $output  = `$command 2>&1`;

	chomp($output);

	( $log_level >> $TRACE ) && print "TRACE $command returned $?, output: $output\n";

	if ( $? == 0 ) {
		return (1);
	}
	else {
		return (0);
	}
}

sub chkconfig_matches {
	my $service          = shift || die("Please supply a service");
	my $service_settings = shift || die("Please supply a chkconfig string to verify");

	( $log_level >> $TRACE ) && print "TRACE Checking whether ${service}'s chkconfig output matches $service_settings.\n";

	# systemd check.
	# This will work for now as  it trys to map from chkconfig run level settings to systemd enabled/disabled state.
	# I think that a new generic endpoint should be added to traffic opts for chkconfig and systemd state settings and that functions
	# here in the ort script should abstract the checking of chkconfig/systemd states with traffic ops.
	if (($RELEASE eq "EL7") || ($RELEASE eq "EL8")) {
		my $service_state = systemd_service_chk($service);
		if ($service_state eq "enabled") {
			if ($service_settings =~ m/on/) {
				( $log_level >> $INFO ) && print "INFO chkconfig output for $service matches $service_settings.\n";
				return 1;
			} else {
				( $log_level >> $ERROR ) && print "ERROR chkconfig output for $service does not match what we expect...\n";
				return 0;
			}
		} else {
			if ($service_settings =~ m/on/) {
				( $log_level >> $ERROR ) && print "ERROR chkconfig output for $service does not match what we expect...\n";
				return 0;
			} else {
				( $log_level >> $INFO ) && print "INFO chkconfig output for $service matches $service_settings.\n";
				return 1;
			}
		}
	} else {
		my $command = "/sbin/chkconfig --list $service";
		my $output  = `$command 2>&1`;
		chomp($output);

		if ( $? == 0 ) {
			if ( $output =~ m/^$service\s+$service_settings$/ ) {
				( $log_level >> $INFO ) && print "INFO chkconfig output for $service matches $service_settings.\n";
				return (1);
			}
			else {
				( $log_level >> $ERROR ) && print "ERROR chkconfig output for $service does not match what we expect...\n";
				( $log_level >> $TRACE ) && print "TRACE $output != $service_settings.\n";
				return (0);
			}
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR $command returned non-zero ($?), output: $output.\n";

			return (0);
		}
	}
}

sub process_chkconfig {
	my $host_name = shift;

	my $proceed = 0;

	my $result = `$atstccfg_cmd $atstccfg_insecure_arg $atstccfg_timeout_arg $atstccfg_arg_disable_proxy --traffic-ops-user='$TO_USER' --traffic-ops-password='$TO_PASS' --traffic-ops-url='$TO_URL' --cache-host-name='$hostname_short' --log-location-error=stderr --log-location-warning=stderr --log-location-info=null --get-data=chkconfig 2>>$atstccfg_log_path`;

	my $atstccfg_exit_code = $?;
	if ($atstccfg_exit_code != 0) {
		( $log_level >> $FATAL ) && print "FATAL Error getting package list from Traffic Ops!\n";
			exit 1;
	}

	if ( defined($result) && $result ne "" && $result !~ m/^\d{3}$/ ) {
		my @chkconfig_list = @{ decode_json($result) };

		for my $chkconfig (@chkconfig_list) {
			if ( package_installed( $chkconfig->{"name"} ) ) {
				if ( !chkconfig_matches( $chkconfig->{"name"}, $chkconfig->{"value"} ) ) {
					if ( $script_mode == $BADASS || $script_mode == $INTERACTIVE ) {
						my $fixit = 0;

						if ( $script_mode == $INTERACTIVE ) {
							if ( get_answer("Are you sure you would like to correct chkconfig for $chkconfig->{name}?") ) {
								$fixit = 1;
							}
						}
						else {
							$fixit = 1;
						}

						if ($fixit) {
							#use systemd commands by mapping chkconfig runlrvrld to either enable or disable.
							if (($RELEASE eq "EL7") || ($RELEASE eq "EL8")) {
								my $systemd_service_enable = "disable";
								if ($chkconfig->{"value"} =~ m/on/) {
									$systemd_service_enable = "enable";
								}
								if (&systemd_service_set($chkconfig->{"name"}, $systemd_service_enable)) {
									( $log_level >> $ERROR ) && print "ERROR $chkconfig->{name}: has been set to $systemd_service_enable\n";
								} else {
									( $log_level >> $ERROR ) && print "ERROR failed to set the systemd service for $chkconfig->{name} to $systemd_service_enable\n";
								}
							} else {
								my (@levels) = split( /\s+/, $chkconfig->{"value"} );

								if ( scalar(@levels) == 7 ) {
									( $log_level >> $TRACE ) && print "TRACE $chkconfig->{name}: Split chkconfig into " . join( ", ", @levels ) . "\n";

									for my $level (@levels) {
										my ( $run_level, $setting ) = split( /:/, $level );

										if ( defined($run_level) && defined($setting) ) {
											( $log_level >> $TRACE ) && print "TRACE $chkconfig->{name}: Setting run level $run_level to $setting\n";

											if ( !set_chkconfig( $chkconfig->{"name"}, $run_level, $setting ) ) {
												( $log_level >> $ERROR ) && print "ERROR $chkconfig->{name}: Unable to set run level $run_level to $setting!\n";
											}
										}
										else {
											( $log_level >> $ERROR ) && print "ERROR $chkconfig->{name}: $level is not what we expected!\n";
										}
									}

									if ( chkconfig_matches( $chkconfig->{"name"}, $chkconfig->{"value"} ) ) {
										( $log_level >> $INFO ) && print "INFO Successfully set chkconfig for $chkconfig->{name}.\n";
									}
									else {
										( $log_level >> $ERROR ) && print "FATAL Unable to set chkconfig values for $chkconfig->{name}!\n";
									}
								}
								else {
									( $log_level >> $ERROR ) && print "ERROR $chkconfig->{name}: $chkconfig->{value} is not what we expected!\n";
								}
							}
						}
					}
					elsif ( $script_mode == $REPORT ) {
						( $log_level >> $INFO ) && print "INFO chkconfig for $chkconfig->{name} DOES NOT MATCH $chkconfig->{value}.\n";
					}
				}
				else {
					if ( $script_mode == $REPORT ) {
						( $log_level >> $INFO ) && print "INFO chkconfig for $chkconfig->{name} matches $chkconfig->{value}.\n";
					}
					else {
						( $log_level >> $TRACE ) && print "TRACE chkconfig for $chkconfig->{name} matches $chkconfig->{value}.\n";
					}
				}
			}
			else {
				( $log_level >> $ERROR ) && print "ERROR $chkconfig->{name} is not installed!\n";
			}
		}
	}
	else {
		( $log_level >> $ERROR ) && print "ERROR No chkconfig parameters returned.\n";
	}
}

sub get_answer {
	my $question = shift || die("Please supply a question");

	my $answer = "";

	while ( $answer !~ /^(y|n)$/i ) {
		( $log_level >> $INFO ) && print "INFO $question [Y/n]: ";
		$answer = <STDIN>;
		chomp($answer);
	}

	if ( $answer =~ /^y$/i ) {
		return (1);
	}
	else {
		return (0);
	}
}

sub start_restart_services {
	my $ats_running = 0;
	#### Start ATS
	if ( package_installed("trafficserver") ) {
		( $log_level >> $DEBUG ) && print "DEBUG trafficserver is installed.\n";
		$ats_running = &start_service("trafficserver");
		if ( $ats_running == $START_SUCCESSFUL ) {
			$traffic_ctl_needed = 0;
			( $log_level >> $DEBUG ) && print "DEBUG trafficserver was just started, no need to run $TRAFFIC_CTL config reload.\n";
		}
		elsif ( $ats_running == $START_FAILED ) {
			$traffic_ctl_needed = 0;
			( $log_level >> $DEBUG ) && print "DEBUG trafficserver failed to start, running $TRAFFIC_CTL config reload will also fail.\n";
		}
		elsif ( $ats_running == $START_NOT_ATTEMPTED ) {
			( $log_level >> $DEBUG ) && print "DEBUG trafficserver was not attempted to be started.\n";
		}
	}

	#### Advanced ATS processing
	if ( $ats_running == $ALREADY_RUNNING && $traffic_ctl_needed && !$trafficserver_restart_needed ) {
		if ( $script_mode == $REPORT ) {
			( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed. '$TRAFFIC_CTL config reload' needs to be run.\n";
		}
		elsif ( $script_mode == $BADASS || $script_mode == $SYNCDS || $script_mode == $REVALIDATE ) {
			( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed. Running '$TRAFFIC_CTL config reload' now.\n";
			&run_traffic_ctl();
		}
		elsif ( $script_mode == $INTERACTIVE ) {
			my $select = 'n';
			( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed. '$TRAFFIC_CTL config reload' needs to be run. Should I do that now? (Y/[n]):";
			$select = <STDIN>;
			chomp($select);
			if ( $select =~ m/Y/ ) {
				&run_traffic_ctl();
				( $log_level >> $DEBUG ) && print "DEBUG traffic_ctl run successful.\n";
				if ( $syncds_update == $UPDATE_TROPS_NEEDED ) {
					$syncds_update = $UPDATE_TROPS_SUCCESSFUL;
				}
			}
			else {
				( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed. '$TRAFFIC_CTL config reload' was not run.\n";
				if ( $syncds_update == $UPDATE_TROPS_NEEDED ) {
					( $log_level >> $ERROR ) && print "ERROR $TRAFFIC_CTL config reload was not run, so Traffic Ops was not updated!\n";
					$syncds_update = $UPDATE_TROPS_FAILED;
				}
			}
		}
	}
	elsif ( $traffic_ctl_needed && ( $ats_running == $START_FAILED || $ats_running == $START_NOT_ATTEMPTED ) ) {
		( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed. The new config will be picked up the next time ATS is started.\n";
		if ( $syncds_update == $UPDATE_TROPS_NEEDED ) {
			( $log_level >> $ERROR ) && print "ERROR $TRAFFIC_CTL config reload was not run, but Traffic Ops is being updated anyway.\n";
			$syncds_update = $UPDATE_TROPS_SUCCESSFUL;
		}
	}
	elsif ( $ats_running && $trafficserver_restart_needed ) {
		if ( $script_mode == $REPORT ) {
			( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed, trafficserver needs to be restarted (service trafficserver restart).\n";
		}
		elsif ( $script_mode == $INTERACTIVE ) {
			my $select = 'n';
			( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed, trafficserver needs to be restarted. Should I do that now? (Y/[n]):";
			$select = <STDIN>;
			chomp($select);
			if ( $select =~ m/Y/ ) {
				my $result = &restart_service("trafficserver");
			}
			else {
				( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed, but trafficserver was not restarted.\n";
			}
		}
		elsif ( $script_mode == $BADASS ) {
			( $log_level >> $ERROR ) && print "ERROR ATS configuration has changed, trafficserver needs to be restarted.\n";
			my $result = &restart_service("trafficserver");
		}
	}
	#### End processing ATS

	#### Start teakd
	if ( package_installed("teakd") ) {
		( $log_level >> $DEBUG ) && print "DEBUG teakd is installed.\n";
		&start_service("teakd");

		# Do something here in the future.
	}
}

sub run_sysctl_p {
	if ( $script_mode == $INTERACTIVE ) {
		my $select = 'n';
		( $log_level >> $ERROR ) && print "ERROR sysctl configuration has changed. 'sysctl -p' needs to be run. Should I do that now? (Y/[n]):";
		$select = <STDIN>;
		chomp($select);
		if ( $select =~ m/Y/ ) {
			my $out    = `sysctl -p 2>&1`;
			my $return = &check_output($out);
			if ( !$return ) {
				( $log_level >> $DEBUG ) && print "DEBUG sysctl -p run successful.\n";
			}
			else {
				( $log_level >> $ERROR ) && print "ERROR sysctl -p failed.\n";
			}
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR sysctl configuration has changed. 'sysctl -p' was not run.\n";
		}
	}
	elsif ( $script_mode == $BADASS ) {
		my $out    = `sysctl -p 2>&1`;
		my $return = &check_output($out);
		if ( !$return ) {
			( $log_level >> $DEBUG ) && print "DEBUG sysctl -p run successful.\n";
		}
		else {
			( $log_level >> $ERROR ) && print "ERROR sysctl -p failed.\n";
		}
	}
}

sub scrape_unencode_text {
	my $text = shift;

	( my @file_lines ) = split( /\n/, $text );
	my @lines;

	foreach my $line (@file_lines) {
		( $log_level >> $TRACE ) && print "TRACE Line from cfg file in TrOps:\t$line\n";
		$line =~ s/\s+/ /g;
		$line =~ s/(^\s+|\s+$)//g;
		$line =~ s/amp\;//g;
		$line =~ s/\&gt\;/\>/g;
		$line =~ s/\&lt\;/\</g;
		chomp($line);
		next if ( $line =~ m/^$/ );

		push( @lines, $line );
	}

	return \@lines;
}

sub scrape_canned_comments {
	my $linesin = $_[0];

	my @linesout;

	foreach my $line (@$linesin) {
		if ( $line =~ m/^\#/ ) {
			if ( $line =~ m/DO NOT EDIT - Generated for / || $line =~ m/$header_comment/ || $line =~ m/TRAFFIC OPS NOTE\:/ || $line =~ m/^##OVERRID.*##/ ) {
				next;
			}
		}
		push( @linesout, $line );
	}

	return \@linesout;
}

sub can_read_write_file {
	my $filename = shift;
	my $file     = $cfg_file_tracker->{$filename}->{'headers'}->{'path'};

	my $username = $ENV{LOGNAME} || $ENV{USER} || getpwuid($<);
	( $log_level >> $TRACE ) && print "TRACE User to validate $file against: $username\n";

	if ( !-R $file ) {
		( $log_level >> $ERROR ) && print "ERROR $file is not readable by $username!\n";
		$cfg_file_tracker->{$filename}->{'audit_failed'}++;
		return 0;
	}

	if ( !-W $file && $script_mode != $REPORT ) {
		( $log_level >> $ERROR ) && print "ERROR $file is not writable by $username!\n";
		$cfg_file_tracker->{$filename}->{'audit_failed'}++;
		return 0;
	}

	( $log_level >> $TRACE ) && print "TRACE RW perms okay for $filename!\n";
	return 1;
}

sub open_file_get_contents {
	my $file = shift;
	my @disk_file_lines;

	( $log_level >> $DEBUG ) && print "DEBUG Opening file from disk:\t$file.\n";
	open my $fh, '<', $file || ( ( $log_level >> $ERROR ) && print "ERROR Can't open $file: $!\n" );

	while (<$fh>) {
		my $line = $_;
		$line =~ s/\s+/ /g;
		$line =~ s/(^\s+|\s+$)//g;
		chomp($line);
		( $log_level >> $TRACE ) && print "TRACE Line from cfg file on disk:\t$line.\n";

		# strip empty lines
		if ( $line =~ m/^$/ ) {
			next;
		}
		push( @disk_file_lines, $line );
	}
	close $fh;

	return \@disk_file_lines;
}

sub prereqs_ok {
	my $filename       = shift;
	my $file_lines_ref = shift;

	( $log_level >> $DEBUG ) && print "DEBUG Starting to check prereqs for:\t$filename.\n";

	if ( $filename eq "plugin.config" || $filename eq "remap.config" ) {
		&check_plugins( $filename, $file_lines_ref );
		if ( $cfg_file_tracker->{$filename}->{'prereq_failed'} ) {
			( $log_level >> $ERROR ) && print "ERROR Prereqs failed for $filename!\n";
			return 0;
		}
	}
	return 1;
}

sub diff_file_lines {
	my $cfg_file        = shift;
	my @db_file_lines   = @{ $_[0] };
	my @disk_file_lines = @{ $_[1] };

	my %db_file_lines   = map { $_ => 1 } @db_file_lines;
	my %disk_file_lines = map { $_ => 1 } @disk_file_lines;

	my @db_lines_missing;
	my @disk_lines_missing;

	my $file = $cfg_file_tracker->{$cfg_file}->{'headers'}->{'path'};

	foreach my $line (@db_file_lines) {
		( $log_level >> $TRACE ) && print "TRACE Line from TrOps: $line!\n";
		if ( !exists $disk_file_lines{$line} ) {
			#### Float compare
			if ( $line =~ m/FLOAT/ ) {
				( my $disk_dum, my $disk_name, my $disk_type, my $disk_val ) = split( /\s/, $line );
				foreach my $l ( keys %db_file_lines ) {
					( my $db_dum, my $db_name, my $db_type, my $db_val ) = split( /\s/, $l );
					if ( $db_name eq $disk_name && $db_type eq $disk_type ) {
						if ( abs( $disk_val - $db_val ) > 0.00001 ) {
							push( @disk_lines_missing, $line );
						}
					}
				}
			}
			elsif ( ( $line =~ m/DO NOT EDIT - Generated for / && $line =~ m/$header_comment/ ) || $line =~ m/TRAFFIC OPS NOTE\:/ ) {
				my $found_it = 0;
				foreach my $line_disk (@disk_file_lines) {
					if ( ( $line =~ m/DO NOT EDIT - Generated for / && $line =~ m/$header_comment/ ) || $line =~ m/TRAFFIC OPS NOTE\:/ ) {
						$found_it++;
					}
				}
				if ( !$found_it ) {
					push( @disk_lines_missing, $line );
				}
			}
			else {
				push( @disk_lines_missing, $line );
			}
		}
	}
	foreach my $line (@disk_file_lines) {
		( $log_level >> $TRACE ) && print "TRACE Line from disk : $line!\n";
		if ( !exists $db_file_lines{$line} ) {
			#### Float compare
			if ( $line =~ m/FLOAT/ ) {
				( my $db_dum, my $db_name, my $db_type, my $db_val ) = split( /\s/, $line );
				foreach my $l (@disk_file_lines) {
					( my $disk_dum, my $disk_name, my $disk_type, my $disk_val ) = split( /\s/, $l );
					if ( $db_name eq $disk_name && $db_type eq $disk_type ) {
						if ( abs( $disk_val - $db_val ) > 0.00001 ) {
							push( @db_lines_missing, $line );
						}
					}
				}
			}
			elsif ( ( $line =~ m/DO NOT EDIT - Generated for / && $line =~ m/$header_comment/ ) || $line =~ m/TRAFFIC OPS NOTE\:/ ) {
				next;
			}
			else {
				push( @db_lines_missing, $line );
			}
		}
	}

	if ( scalar(@db_lines_missing) || scalar(@disk_lines_missing) ) {
		( $log_level >> $ERROR ) && print "ERROR Lines for $file from Traffic Ops do not match file on disk.\n";
	}
	if ( scalar(@db_lines_missing) ) {
		my $line_count = scalar(@db_lines_missing);
		( $log_level >> $DEBUG ) && print "DEBUG $line_count lines are missing from file that is in Traffic Ops.\n";
		foreach my $line (@db_lines_missing) {
			( $log_level >> $ERROR ) && print "ERROR Config file $cfg_file line only on disk :\t$line\n";
		}
	}

	if ( scalar(@disk_lines_missing) ) {
		my $line_count = scalar(@disk_lines_missing);
		( $log_level >> $DEBUG ) && print "DEBUG $line_count lines are missing from file that is on disk.\n";
		foreach my $line (@disk_lines_missing) {
			( $log_level >> $ERROR ) && print "ERROR Config file $cfg_file line only in TrOps:\t$line\n";
		}
	}

	return ( \@db_lines_missing, \@disk_lines_missing );
}

sub backup_file {
	my $filename   = shift;
	my $result_ref = shift;

	my $result   = ${$result_ref};
	my $file     = $cfg_file_tracker->{$filename}->{'headers'}->{'path'};

	if ( $script_mode != $REPORT ) {
		my $ats_uid  = getpwnam("ats");
		my $bkp_dir;
		my $bkp_file;
		if ( -e $file ) {
			( $log_level >> $ERROR ) && print "ERROR Creating backup of file on disk for $filename.\n";
			$bkp_dir  = $TMP_BASE . "/" . $unixtime . "/" . $cfg_file_tracker->{$filename}->{'service'} . "/config_bkp/";
			$bkp_file = $bkp_dir . $filename;
			&smart_mkdir($bkp_dir);
			( $log_level >> $DEBUG ) && print "DEBUG Backup file: $bkp_file.\n";
			$cfg_file_tracker->{$filename}->{'backup_from_disk'} = $bkp_file;
			system("/bin/cp $file $bkp_file");
		}
		else {
			( $log_level >> $DEBUG ) && print "DEBUG Config file: $file doesn't exist. No need to back up.\n";
		}
		( $log_level >> $ERROR ) && print "ERROR Creating backup of file in TrOps for $filename.\n";
		$bkp_dir  = $TMP_BASE . "/" . $unixtime . "/" . $cfg_file_tracker->{$filename}->{'service'} . "/config_trops/";
		$bkp_file = $bkp_dir . $filename;
		&smart_mkdir($bkp_dir);
		( $log_level >> $DEBUG ) && print "DEBUG Backup file: $bkp_file.\n";
		$cfg_file_tracker->{$filename}->{'backup_from_trops'} = $bkp_file;
		open my $fh, '>', $bkp_file || die "Can't open $bkp_file for writing!\n";
		print $fh $result;
		chmod oct(644), $fh;
		chown $ats_uid, $ats_uid, $fh;
		close $fh;
	}
	return 0;
}

sub adv_preprocessing_remap {
	my $buffer = ${ $_[0] };

	( my @file_lines ) = split( /\n/, $buffer );

	if ( 1 < $#file_lines ) { #header line is always present, so look for 2 lines or more
		( $log_level >> $DEBUG ) && print "DEBUG Entering advanced pre-processing for remap.config.\n";

		# key on the FROM remap
		my %override_hash=();

		# sweep for override lines
		foreach my $line (@file_lines) {
			if ( $line =~ m/^##OVERRIDE##\s*/ ) {
				my $newline = $line =~ s/^##OVERRIDE##\s*//r;
				( my @fields ) = split( /\s+/, $newline, 4 );
				# ensure <rule> <from> <to>
				if ( 2 < $#fields ) {
					$override_hash{$fields[1]}=1;
				}
			}
		}

		# handle override lines, etc
		my $overrides = keys %override_hash;
		if (0 < $overrides) {
			( $log_level >> $INFO ) && print "INFO Entering ##OVERRIDE## processing for remap.config.\n";

			my @out_lines;
			my $lineno = 0;
			my $cnt = 0;

			foreach my $line (@file_lines) {
				if ( $line =~ m/DO NOT EDIT - Generated for/ ) {
					push( @out_lines, $line );
				}
				elsif ( $line =~ m/^##OVERRIDE##\s*/ ) {
					$line =~ s/^##OVERRIDE##\s*//;
					( $log_level >> $DEBUG ) && print "DEBUG new override line $lineno: \'$line\'\n";
					push( @out_lines, "##OVERRIDE##" );
					push( @out_lines, $line );
				}
				elsif ( $line !~ m/^\s*#/ ) {
					( my @fields ) = split( /\s+/, $line, 4 );
					# ensure <rule> <from> <to>
					if ( 2 < $#fields && exists($override_hash{$fields[1]}) ) {
						( $log_level >> $DEBUG ) && print "DEBUG old overriden line $lineno: \'$line\'\n";
						$cnt++;
						$line = "##OVERRIDDEN## " . $line;
					}

					push( @out_lines, $line );
				}
				else {
					push( @out_lines, $line );
				}

				$lineno++;
			}

			( $log_level >> $WARN ) && print "WARN Overrode $cnt old remap rule(s) with $overrides new remap rule(s).\n";
			$buffer = join( "\n", @out_lines ) . "\n";
		}
	}

	return $buffer;
}

sub adv_processing_udev {
	my @db_file_lines = @{ $_[0] };

	( $log_level >> $DEBUG ) && print "DEBUG Entering advanced processing for 50-ats.rules.\n";
	foreach my $line50 (@db_file_lines) {
		if ( $line50 =~ m/KERNEL/ && $line50 =~ m/OWNER/ ) {
			( my $dev, my $should_own ) = split( /,/, $line50 );
			$dev =~ s/KERNEL\s*\=\=\s*//g;
			$dev =~ s/\"//g;
			$should_own =~ s/ OWNER\s*:?\=\s*//g;
			$should_own =~ s/\"//g;

			my $dev_path = "/dev/$dev";
			my $dc       = undef;

			next if ( $should_own eq "root" );

			my $ats_uid = `/usr/bin/id $should_own 2>&1`;

			if ( $ats_uid =~ m/No such user/ ) {
				( $log_level >> $ERROR ) && print "ERROR User: $should_own does not exist! Skipping future checks for $dev_path\n";
				next;
			}

			chomp($ats_uid);
			$ats_uid =~ s/\((.*)$//g;
			$ats_uid =~ s/uid\=//g;

			if ( -e $dev_path ) {
				( $log_level >> $TRACE ) && print "TRACE Found device in 50-ats.rules: $dev_path.\n";
				( $dc, $dc, $dc, $dc, my $uid, $dc, $dc, $dc, $dc, $dc, $dc, $dc, $dc ) = stat($dev_path);
				if ( $uid != $ats_uid ) {
					( $log_level >> $ERROR ) && print "ERROR Device $dev_path is owned by $uid, not $should_own ($ats_uid)\n";
				}
				( my @df_lines ) = split( /\n/, `/bin/df` );
				foreach my $l (@df_lines) {
					if ( $l =~ m/$dev_path\d/ ) {
						( $log_level >> $FATAL ) && print "FATAL Device /dev/$dev has an active partition and a file system!!\n";
					}
				}
			}
			else {
				open( DEV, "ls /dev/* |" ) or ( $log_level >> $FATAL ) && print "FATAL Couldn't get /dev/ listing: $!\n";
				while ( my $dnode = <DEV> ) {
					next unless ( $dnode =~ m!$dev_path! );

					chomp $dnode;
					next if ( $dnode =~ m!/dev/sda[0-9]*! );

					( $log_level >> $TRACE ) && print "TRACE Found device in 50-ats.rules: $dnode.\n";
					( $dc, $dc, $dc, $dc, my $uid, $dc, $dc, $dc, $dc, $dc, $dc, $dc, $dc ) = stat($dnode);
					if ( $uid != $ats_uid ) {
						( $log_level >> $ERROR ) && print "ERROR Device $dnode is owned by $uid, not $should_own ($ats_uid)\n";
					}
					( my @df_lines ) = split( /\n/, `/bin/df` );
					foreach my $l (@df_lines) {
						if ( $l =~ m/$dnode/ ) {
							( $log_level >> $FATAL ) && print "FATAL Device /dev/$dev has an active partition and a file system!!\n";
						}
					}
				}
				close(DEV);
			}
		}
	}
	return 0;
}

{
	my $fh;

	sub check_only_copy_running {
		return if $fh;
		open $fh, '<', $0 or die $!;

		unless ( flock( $fh, LOCK_EX | LOCK_NB ) ) {
			( $log_level >> $FATAL ) && print "FATAL $0 is already running. Exiting.\n";
			exit 1;
		}
	}
}
