package Extensions::Helper::TestHelper;
#
# Copyright 2011-2014, Comcast Corporation. This software and its contents are
# Comcast confidential and proprietary. It cannot be used, disclosed, or
# distributed without Comcast's prior written permission. Modification of this
# software is only allowed at the direction of Comcast Corporation. All allowed
# modifications must be provided to Comcast Corporation.
#
use Data::Dumper;

sub new {
	my $self  = {};
	my $class = shift;

	return ( bless( $self, $class ) );
}

sub prepend_extensions {
	print( "PERL5LIB: " . Dumper(@INC) . "\n" );
	my $to_ext_lib = $ENV{"TO_EXTENSIONS_LIB"};
	if ( defined($to_ext_lib) ) {
		if ( -e $to_ext_lib ) {
			unshift( @INC, $to_ext_lib );
			print "Found TO_EXTENSIONS_LIB prepending to library path: $to_ext_lib\n";
		}
		else {
			print "\nWARNING TO_EXTENSIONS_LIB environment variable is defined as $to_ext_lib but does not exist.\n\n";
		}
	}
}

1;
