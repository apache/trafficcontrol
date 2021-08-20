#!/usr/bin/env perl 
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
use English;
use Getopt::Long;
use FileHandle;
use Cwd;
use Data::Dumper;
use File::Find::Rule;
use File::Path qw(make_path remove_tree);

my $usage = "\n"
	. "Usage:  $PROGRAM_NAME --gpg-key=[your-signed-key-id] --release-no=[release-to-create]\t\n\n"
	. "Example:  $PROGRAM_NAME --gpg-key=75AFDE1 --release-no=RELEASE-1.1.0 \n\n"
	. "Purpose:  This script automates the release process for the Traffic Control cdn.\n"
	. "\nFlags:   \n\n"
	. "--gpg-key          - Your gpg-key id. ie: 774ACED1\n"
	. "--release-no       - The release-no name you want to cut. ie: 1.1.0\n"
	. "--git-short-hash   - (optional) The git hash that will be used to reference the release. ie: da4aab57d \n"
	. "--git-remote-url   - (optional) Overrides the git repo URL where the release will be pulled and sent (mostly for testing). ie: git\@github.com:yourrepo/trafficcontrol.git \n"
	. "--dry-run          - (optional) Simulation mode which will NOT apply any changes. \n"
	. "--debug            - (optional) Show debug output\n"
	. "\nArguments:   \n\n"
	. "cut        - Cut the release branch, tag the release then make the branch, tag public.\n"
	. "cleanup    - Reverses the release steps in case you messed up.\n"
	. "pushdoc    - Upload documentation to the public website.\n";

my $git_remote_name = 'official';

my $git_remote_url = 'git@github.com:apache/trafficcontrol.git';

my $gpg_key;
my $release_no;

# Example: 1.1.0
my $version;

# Example: 1.2.0
my $next_version;

# Example: 1.1.x
my $new_branch;

# Example: RC0
my $build_no;

# Example: RC1
my $next_build_no;

# Example: 774ACED1
my $git_short_hash;

# Keeps track of the branch to determine RC flow
my $branch_exists;

my $rc;
my $dry_run = 0;
my $debug   = 0;
my $working_dir;

GetOptions(
	"gpg-key=s"        => \$gpg_key,
	"release-no=s"     => \$release_no,
	"git-short-hash=s" => \$git_short_hash,
	"git-remote-url=s" => \$git_remote_url,
	"dry-run!"         => \$dry_run,
	"debug!"           => \$debug
);

#TODO: drichardson - Preflight check for commands 'git', 's3cmd' , '
#                  - Add validation logic here for required flags
#                  - Upload Release (s3cmd)

STDERR->autoflush(1);
my $argument = shift(@ARGV);

if ( defined($argument) ) {

	if ( $argument eq 'cut' ) {
		fetch_branch();
		my $prompt = "Continue with creating the RELEASE?";

		if ( prompt_yn($prompt) ) {

			# Only tag the release
			print "branch_exists #-> (" . $branch_exists . ")\n";
			if ($branch_exists) {
				add_official_remote();
				publish_version_file( $new_branch, $version );
				tag_and_push();
			}
			else {
				add_official_remote();
				cut_new_release();
				tag_and_push();
			}

		}
		else {
			exit(0);
		}
	}
	elsif ( $argument eq 'cleanup' ) {
		my $prompt = "\n\nAre you sure you want to cleanup the RELEASE? (" . $release_no . ")";
		if ( prompt_yn($prompt) ) {
			fetch_branch();
			cleanup_release();
		}
	}
	else {
		print $usage;
	}
}
else {
	print $usage;
}

exit(0);

sub fetch_branch {

	clone_repo_to_tmp();

	parse_variables();
	chdir $working_dir;
	( $rc, $branch_exists ) = check_branch_exists();

	# if not passed in as option, then determine git hash
	if ( !defined($git_short_hash) ) {
		( $rc, $git_short_hash ) = get_git_short_hash();
	}

	my $new_branch_info = <<"INFO";
\nNEW Release Summary
Git Repo       : $git_remote_url
Version        : $version
Branch         : $new_branch
Tag            : $release_no
Git Short Hash : $git_short_hash
INFO

	my $release_candidate_info = <<"INFO";
\nRelease CANDIDATE Summary
Git Repo       : $git_remote_url
Version        : $version
Branch         : $new_branch
Next Tag       : $release_no
Git Short Hash : $git_short_hash
INFO

	if ( $release_no !~ /RC/ ) {
		print $new_branch_info;
	}
	elsif ($branch_exists) {
		print $release_candidate_info;
	}
	else {
		print $new_branch_info;
	}

}

sub get_git_short_hash {
	my $cmd = "git log --pretty=format:'%h' -n 1";
	my ( $rc, $git_hash ) = run_and_capture_command( $cmd, "force" );
	if ( $rc > 0 ) {
		print " Failed to retrieve git hash : " . $cmd . " \n ";

		#exit(1);
	}
	return $rc, $git_hash;
}

sub check_branch_exists {
	my $cmd = sprintf( "git checkout %s", $new_branch );
	my $output;
	my $git_branch;
	( $rc, $output ) = run_and_capture_command( $cmd, "force" );

	$cmd = sprintf( "git rev-parse --verify %s", $new_branch );
	( $rc, $git_branch ) = run_and_capture_command( $cmd, "force" );
	if ( $rc > 0 ) {
		$branch_exists = 0;
	}
	else {
		$branch_exists = 1;
	}
	return $rc, $branch_exists;
}

sub clone_repo_to_tmp {
	my $tmp_dir = "/tmp";
	my $tc_dir  = "trafficcontrol";
	$working_dir = sprintf( "%s/%s", $tmp_dir, $tc_dir );
	remove_tree($working_dir);
	chdir $tmp_dir;
	print "Cloning output to: " . $working_dir . "\n";
	my $cmd = "git clone " . $git_remote_url;
	chdir $working_dir;

	my $rc = run_command( $cmd, "force" );
	if ( $rc > 0 ) {
		print " Failed to clone repo : " . $cmd . " \n ";

		#exit(1);
	}
}

sub parse_variables {
	my $major;
	my $minor;
	my $patch;
	if ( $release_no =~ /RC/ ) {
		( $major, $minor, $patch, $build_no ) = ( $release_no =~ /RELEASE-(\d).(\d).(\d)-(.*)/ );
		my ( $rc_build_no, $x ) = ( $build_no =~ /RC(\d+)/ );
		my $next_build_no_version = $rc_build_no + 1;
		$next_build_no = sprintf( "RC%d", $next_build_no_version );
	}
	else {
		( $major, $minor, $patch ) = ( $release_no =~ /RELEASE-(\d).(\d).(\d)/ );
	}

	$version = sprintf( "%s.%s.%s", $major, $minor, $patch );
	my $next_minor = $minor + 1;
	$next_version = sprintf( "%s.%s.%s", $major, $next_minor, $patch );
	$new_branch = sprintf( "%s.%s.x", $major, $minor );

}

sub add_official_remote {
	my $cmd = "git remote add official " . $git_remote_url;
	my $rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Added new origin: " . $git_remote_name . " " . $git_remote_url . "\n\n";
	}
}

sub cut_new_release {

	print "Creating new branch\n";
	my $cmd = "git checkout -b " . $new_branch;
	my $rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to checkout new branch" . $cmd . "\n";
	}

	publish_version_file( "master", $next_version );

}

sub publish_version_file {

	my $branch  = shift;
	my $version = shift;
	my $cmd     = "git checkout " . $branch;
	my $rc      = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to checkout new branch" . $cmd . "\n";
	}
	update_version_file($version);
	$cmd = "git commit -m 'RELEASE: Syncing VERSION file' VERSION";
	$rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to run:" . $cmd . "\n";
	}

	print "Updating 'VERSION' file\n";
	$cmd = "git push official " . $branch;
	$rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to push official to master" . $cmd . "\n";
	}

}

sub tag_and_push {
	print "Signing new tag based upon your gpg key\n";
	my $comment = "Release " . $version;
	my $cmd = sprintf( "git tag -s -u %s -m '%s' %s", $gpg_key, $comment, $release_no );
	$rc = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to tag and push" . $cmd . "\n";
	}

	print "Making new release tag and branch publicly available.\n";

	#$cmd = "git push official " . $new_branch;
	$cmd = "git push --follow-tags official " . $new_branch;
	$rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to tag release" . $cmd . "\n";
	}
}

sub cleanup_release {

	if ($debug) {
		print "gpg_key #-> (" . $gpg_key . ")\n";
		print "release_no #-> (" . $release_no . ")\n";
		print "dry_run #-> (" . $dry_run . ")\n";
	}

	my $cmd = "git remote add official " . $git_remote_url;
	my $rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Added new origin: " . $git_remote_name . " " . $git_remote_url . "\n\n";
	}
	else {
		print "Found Official : " . $git_remote_name . " " . $git_remote_url . "\n\n";
	}

	update_version_file($version);
	$cmd = "git commit -m 'RELEASE: Decrementing VERSION file' VERSION";
	$rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to run:" . $cmd . "\n";
	}

	print "Updating 'VERSION' file\n";
	$cmd = "git push official master";
	$rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to run:" . $cmd . "\n";
	}

	print "Creating new branch\n";
	$cmd = "git push origin --delete " . $new_branch;
	$rc  = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to run:" . $cmd . "\n";
	}

	my $comment = "Release " . $version;
	$cmd = sprintf( "git tag -d %s", $release_no );
	$rc = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to run:" . $cmd . "\n";
	}

	$cmd = sprintf( "git push origin :refs/tags/%s", $release_no );
	$rc = run_command($cmd);
	if ( $rc > 0 ) {
		print "Failed to run:" . $cmd . "\n";
	}

}

sub update_version_file {

	my $version_no = shift;

	my $version_file_name = "VERSION";
	open my $fh, '<', $version_file_name or die "error opening $version_file_name $!";
	my $data = do { local $/; <$fh> };

	if ($dry_run) {
		print "Would have updated VERSION file to: " . $version_no . "\n";
	}
	else {
		print "PRIOR version: " . $data . "\n";
		print "Version: " . $version_no . "\n";
		open( $fh, '>', $version_file_name ) or die "Could not open file '$version_file_name' $!";
		print $fh $version_no . "\n";
		close $fh;
	}
}

sub deploy_documentation {

}

sub prompt {
	my ($query) = @_;    # take a prompt string as argument
	local $| = 1;        # activate autoflush to immediately show the prompt
	print $query;
	chomp( my $answer = <STDIN> );
	return $answer;
}

sub prompt_yn {
	my ($query) = @_;
	my $answer = prompt("$query (Y/N): ");
	return lc($answer) eq 'y';
}

sub run_and_capture_command {
	my ( $cmd, $force ) = @_;
	if ( $dry_run && ( !defined($force) ) ) {
		print "Simulating cmd:> " . $cmd . "\n\n";
		return 0;
	}
	else {
		if ($debug) {
			print "Capturing COMMAND> " . $cmd . "\n\n";
		}
		my $cmd_output = `$cmd </dev/null`;

		#my $cmd_output = `$cmd >/dev/null 2>&1`;
		return $?, $cmd_output;
	}
}

sub run_command {
	my ( $cmd, $force ) = @_;
	if ( $dry_run && ( !defined($force) ) ) {
		print "Simulating cmd:> " . $cmd . "\n\n";
		return 0;
	}
	else {
		if ($debug) {
			print "Executing COMMAND> " . $cmd . "\n\n";
		}
		system($cmd);

		return $?;
	}
}
