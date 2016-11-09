#
#
#
#

package InstallUtils;

use Term::ReadPassword;
use base qw{ Exporter };
our @EXPORT_OK = qw{ execCommand randomWord promptUser promptRequired promptPassword promptPasswordVerify trim readJson writeJson writePerl};
our %EXPORT_TAGS = ( all => \@EXPORT_OK );

sub execCommand {
	my ( $cmd, @args ) = @_;
	system( $cmd, @args );
	my $result = $? >> 8;
	return $result;
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

sub promptUser {
	my ( $promptString, $defaultValue, $noEcho ) = @_;

	if ($defaultValue) {
		print $promptString, " [", $defaultValue, "]:  ";
	}
	else {
		print $promptString, ":  ";
	}

	if ( defined $noEcho && $noEcho ) {
		my $response = read_password('');
		if ( ( !defined $response || $response eq '' ) && ( defined $defaultValue && $defaultValue ne '' ) ) {
			$response = $defaultValue;
		}
		return $response;
	}
	else {
		$| = 1;
		$_ = <STDIN>;
		chomp;

		if ("$defaultValue") {
			return $_ ? $_ : $defaultValue;
		}
		else {
			return $_;
		}
		return $_;
	}
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
	open( my $fh, '<', $file ) or return;
	local $/;    # slurp mode
	my $text = <$fh>;
	undef $fh;
	return JSON->new->utf8->decode($text);
}

sub writeJson {
	my $file = shift;
	open( my $fh, '>', $file ) or die("open(): $!");
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
 	my $dumper = Data::Dumper->new([ $data ]);
 
 	# print without var names and with simple indentation
 	print $fh $dumper->Terse(1)->Dump();
 	close $fh;
}

1;
