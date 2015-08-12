#
#
#
#

package InstallUtils;

use Term::ReadPassword;

sub execCommand {
	my ( $cmd, @args ) = @_;
	system( $cmd, @args );
	my $result = $? >> 8;
	return $result;
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

sub trim {
	my $str = shift;

	$str =~ s/^\s+//;
	$str =~ s/^\s+$//;

	return $str;
}

1;
