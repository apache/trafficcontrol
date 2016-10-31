package Utils::Properties;
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

use strict;
use warnings;

use Config::Properties;

my $propsRef = ();
my $file;
my $properties;

sub new {
	my $class = shift;
	$file = shift;
	my $fh;

	unless (length $file) {
		die "Usage:  Utils::Properties->loadProperties (\$propertiesFile)";
	}

	$properties = Config::Properties->new ();

	#
	# If the properties file does not exist or is empty, open the file 
	# read / write to insure that it is created so that saveProperty()
	# may be used.  Then make sure the hash pointed to by propsRef is
	# intialized so that it is blessed.
	# 
	if (! -f $file || -z $file) {
		open ($fh, "+>", $file) or die "$file: $!";
		
		$propsRef->{"file"} = $file;
		
		close ($fh);
	}
	#
	# The properties file exists and is not empty.  Load the
	# name value pairs it contains.
	else {
		open ($fh, "<", $file) or die "$file: $!";

		$properties->load ($fh);

		my %props = $properties->properties;

		foreach my $k (keys %props) {
			$propsRef->{$k} = $props{$k};
		}

		close ($fh);
	}
	
	bless $propsRef, $class;

	return $propsRef;
}

sub getProperty {
	my $class = shift;
	my $name = shift;

	unless (length $name) {
		die "Usage:  Utils::Properties->getProperty (\$name)";
	}

	if (defined  $propsRef->{$name}) {
		return $propsRef->{$name}; 
	}
	else {
		return "";
	}
}

sub saveProperty {
	my $class = shift;
	my $name = shift;
	my $value = shift;
	my $header = shift;

	unless (length $name && length $value) {
		die "Usage:  Utils::Properties->saveProperty (\$name,\$value,\$header)";
	}

	open (my $fh, ">", $file) or die "$file: $!";

	$properties->setProperty ($name, $value);

	if (length $header) {
		$properties->store ($fh, $header);
	}
	else {
		$properties->store ($fh);
	}

	close ($fh);
}

1;
