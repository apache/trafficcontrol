package Utils::JsonConfig;
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
#
#

# Utility for loading .conf file into a resulting hash
use strict;
use warnings;
use JSON;
use Data::Dumper;

my $json_ref = ();
my $file;

sub new {
	my $class = shift;
	$file = shift;
	my $fh;

	unless ( length $file ) {
		die "Usage:  Utils::JsonConfig->new (\$file)";
	}

	if ( !-f $file || -z $file ) {
		die "$file: $!";
	}
	open( $fh, "<", $file ) or die "$file: $!";

	my $json_fh = <$fh>;
	$json_ref = decode_json($json_fh);
	$json_ref->{"file"} = $file;
	close($fh);

	bless $json_ref, $class;

	return $json_ref;
}

sub get {
	my $key = shift;

	unless ( length $key ) {
		die "Usage:  Utils::JsonConfig->get (\$key)";
	}

	if ( defined $json_ref->{$key} ) {
		return $json_ref->{$key};
	}
	else {
		return "";
	}
}

sub load_conf {
	my $self           = shift;
	my $mode           = shift;
	my $conf_file_name = shift;

	local $/;    #Enable 'slurp' mode
	my $conf_path = "conf/" . $mode . "/" . $conf_file_name;
	return Utils::JsonConfig->new($conf_path);
}

1;
