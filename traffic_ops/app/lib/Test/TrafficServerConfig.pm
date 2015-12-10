package Test::TrafficServerConfig;

use strict;
use warnings;
use Carp qw/cluck/;

BEGIN {
	use Exporter;
	our @EXPORT_OK = qw{ loadConfigFile loadConfig };
}

sub parseConfigLine {
	my $line = shift;
	my $uq   = qr/(?<!\\)"/;
	my %h;

	$line =~ s/^\s+//;
	$line =~ s/\s+$//;

	my $qs;    # in quoted string
	while ( $line =~ s/^\s*(.+?)($uq|\s)// ) {
		my ( $pre, $found ) = ( $1, $2 );
		if ( !defined $qs ) {

			# not in quoted string
			if ( $found !~ qr/$uq/ ) {
				my ( $k, $v ) = split /=/, $pre;
				$h{$k} = $v;
			}
			else {
				# start of quoted string skip spaces until next quote;
				$qs = 1;
			}
		}
		else {
			if ( $found eq '"' ) {

				# end of quote -- go back to looking for space
				undef $qs;
			}

			# else keep looking for end quote
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
