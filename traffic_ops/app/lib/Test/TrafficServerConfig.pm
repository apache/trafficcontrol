package Test::TrafficServerConfig;

use strict;
use warnings;

BEGIN {
	use Exporter;
	our @EXPORT_OK = qw{ loadConfigFile loadConfig };
}

sub loadConfigFile {
	my $cf = shift;
	if ( !-f $cf ) {
		return {};
	}

	open my $cfh, '<', $cf or return {};

	local $/;    # slurp mode
	my $txt = <$cfh>;
	close $cfh;

	return loadConfig($txt);
}

sub loadConfig {
	my $lines = shift;
	my $uq    = qr/(?<!\\)"/;
	my @config;
	for my $line ( split /\n/, $lines ) {

		$line =~ s{\s*#.*}{};
		next unless $line =~ /\S/;

		# create hash to represent this one line and populate
		my %h;
		my $cur = '';

		# split on unescaped quotes -- avoids need for quote-pairing
		my @b = split /$uq/, $line;

		while ( scalar @b > 0 ) {
			$cur .= shift @b;
			while ( $cur =~ /^[,\s]*(\w+)=+(.*)/ ) {
				my ( $k, $remainder ) = ( $1, $2 );
				if ( length $remainder == 0 ) {

					# incomplete -- add next chunk
					last;
				}
				if ( $remainder =~ /^$uq(.*?)$uq(.*)/ or $remainder =~ /(\S+)(.*)/ ) {
					$h{$k} = $1;
					$cur = $2;
				}
				else {
					die "Malformed? $k=$remainder";
				}
			}
		}
		if ( length $cur != 0 ) {
			die "$cur left over.  Malformed line?: $line";
		}

		push @config, \%h;
	}
	return \@config;
}

1;
