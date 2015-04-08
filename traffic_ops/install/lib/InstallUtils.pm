#
#
#
#

package InstallUtils;

use Term::ReadPassword;

my $self = {};

sub new {
	my ($class) = @_;

	return (bless ($self, $class));
}

sub execCommand {
	my ($class, $command, @args) = @_;
	my $pid = fork ();
	my $result = 0;

	if ($pid == 0) {
		exec ($command, @args);
		exit 0;
	}
	else {
		wait;
		$result = $?;
		if ($result != 0) {
			print "ERROR executing: $commands,  args: " . join (' ', @args) . "\n";
		}
	}
	return $result;
}

sub promptUser {
    my ($class, $promptString, $defaultValue, $noEcho) = @_;

    if ($defaultValue) {
        print $promptString, " [", $defaultValue, "]:  ";
    }
    else {
        print $promptString, ":  ";
    }

    if (defined $noEcho && $noEcho)  {
        my $response = read_password('');
        if ((!defined $response || $response eq '') && (defined $defaultValue && $defaultValue ne '')) {
            $response = $defaultValue;
        }
        return $response
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
