#!/usr/bin/perl

#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

package InstallUtils;

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


use Term::ReadKey;
use JSON;
use IO::Pipe;
use base qw{ Exporter };
our @EXPORT_OK = qw{ execCommand randomWord promptUser promptRequired promptPassword promptPasswordVerify trim readJson writeJson writePerl errorOut logger rotateLog};
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

my $logFile;
my $debug;

sub initLogger {
    $debug   = shift;
    $logFile = shift;
}

sub execCommand {
        my ( $cmd, @args ) = @_;
        system( $cmd, @args );
        my $result = $? >> 8;
        return $result;
}

# log the error and then kill the process
sub errorOut {
    logger( @_, "error" );
    die;
}

# moves a log to file to a backup file with the same name appended with .bkp
# This function is intended to keep log file sizes low and is called from postinstall
sub rotateLog {
    my $logFileName = shift;

    if ( !-f $logFileName ) {
        logger( "Log file '$logFileName' does not exist - not rotating log", "warn" );
        return;
    }

    execCommand( '/bin/mv', '-f', $logFileName, $logFileName . '.bkp' );
    logger( "Rotated log $logFileName", "info" );
}

# outputs logging messages to terminal and log file
sub logger {
    my $output = shift;
    my $type   = shift;

    # optional custom log file to use instead of main log file used by postinstall
    # cpan uses a custom log file because of its size
    my $customLogFile = shift;

    my $message = $output;
    if ( index( $message, "\n" ) == -1 ) {
        $message = $message . "\n";
    }

    # if in debug mode or message is more critical than info print to console
    if ( $debug || ( defined $type && $type ne "" && $type ne "info" ) ) {
        print($message);
    }

    # output to log file
    my $fh;
    my $result = 0;
    if ( defined $customLogFile && $customLogFile ne "" ) {
        open $fh, '>>', $customLogFile or die("Couldn't open log file '$customLogFile'");
        $result = 1;
    }
    else {
        if ($logFile) {
            open( $fh, '>>', $logFile ) or die("Couldn't open log file '$logFile'");
            $result = 1;
        }
    }

    if ($result) {
        print $fh localtime . ": " . uc($type) . ' ' . $message;
        close $fh;
    }
}

sub randomWord {
    my $length = shift || 12;
    my $secret = '';
    while ( length($secret) < $length ) {
        my $c = chr( rand(0x7F) );
        if ( $c =~ /\w/ ) {
            $secret .= $c;
        }
    }
    return $secret;
}

# Any checks to user input configurations can be implemented in this functions.
sub sanitize{
    my ($promptString, $userInput) = @_;
    if (index($promptString, "Human-readable CDN Name") != -1){
        if ($userInput =~ m/[^-a-zA-Z0-9_.]/){
            print "Invalid characters in user input. Try again.\n\n";
            return false;
        }
    }
    return true;
}

sub promptUser {
    my ( $promptString, $defaultValue, $noEcho ) = @_;

    if ( defined $noEcho && $noEcho ) {
        # Set echo mode to off via ReadMode 2
        ReadMode 2;
    }

    # Check user input for invalid characters
    my $sanitized = false;
    while ($sanitized eq false){
        if ($defaultValue) {
            print $promptString, " [", $defaultValue, "]: ";
        }
        else {
            print $promptString, ": ";
        }
        $| = 1;
        $_ = <STDIN>;
        chomp;
        $sanitized = sanitize($promptString, $_)
    }

    if ( defined $noEcho && $noEcho ) {
        # Set echo mode to on via ReadMode 1
        ReadMode 1;
        # Print extra line because echo mode was off during the STDIN
        print "\n";
    }

    if ("$defaultValue") {
        return $_ ? $_ : $defaultValue;
    }

    # If we have gotten here, return what is left
    return $_;
}

sub promptRequired {
    my $val = '';
    while ( length($val) == 0 ) {
        $val = promptUser(@_);
    }
    return $val;
}

sub promptPassword {
    my $prompt = shift;
    my $pw = promptRequired( $prompt, '', 1 );
    return $pw;
}

sub promptPasswordVerify {
    my $prompt = shift;
    my $pw     = shift;

    while (1) {
        $pw = promptPassword($prompt);
        my $verify = promptPassword("Re-Enter $prompt");
        last if $pw eq $verify;
        print "\nError: passwords do not match, try again.\n\n";
    }
    return $pw;
}

sub trim {
    my $str = shift;

    $str =~ s/^\s+//;
    $str =~ s/^\s+$//;

    return $str;
}

sub readJson {
    my $file = shift;
    open( my $fh, '<', $file ) or die("open(): $!");
    local $/;    # slurp mode
    my $text = <$fh>;
    undef $fh;
    return JSON->new->utf8->decode($text);
}

sub writeJson {
    my $file = shift;
    open( my $fh, '>', $file ) or die("open(): $!");
    logger( "Writing json to $file", "info" );
    foreach my $data (@_) {
        my $json_text = JSON->new->utf8->pretty->encode($data);
        print $fh $json_text, "\n";
    }
    close $fh;
}

sub writePerl {
    my $file = shift;
    my $data = shift;

    open( my $fh, '>', $file ) or die("open(): $!");
    my $dumper = Data::Dumper->new( [$data] );

    # print without var names and with simple indentation
    print $fh $dumper->Quotekeys(0)->Terse(1)->Dump();
    close $fh;
}

1;
