package Test::TrafficServerConfig;

use strict;
use warnings;
use Carp qw/cluck/;

BEGIN {
	use Exporter;
	our @EXPORT_OK = qw{ loadConfigFile loadConfig };
}

my $keyval_re = qr/
	(\w+)=               # key=
	(
		"[^"]*"  |       # quoted string
		[^"\s]*          # unquoted value (no spaces)
	)
	(?:\s+|$)            # white space or end-of-line
	/x;

sub parseConfigLine {
	my $line = shift;
	my %h;

	$line =~ s/^\s+//;
	$line =~ s/\s+$//;

	while ( $line =~ /${keyval_re}\s*/g ) {
		my ( $k, $v ) = ( $1, $2 );

		# remove surrounding quotes if there
		$v =~ s/^"(.*)"$/$1/;
		if ( $k =~ /parent/ ) {
			$h{$k} = [ split /;/, $v ];
		}
		else {
			$h{$k} = $v;
		}
	}
	return \%h;
}

sub parseConfig {
	my $lines = shift;
	my $uq    = qr/(?<!\\)"/;
	my @config;
	for my $line ( split /\n/, $lines ) {
		next if $line =~ /^\s*#/;
		push @config, parseConfigLine($line);
	}
	return \@config;
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

	return parseConfig($txt);
}

1;
