package Extensions::CustomDatasourceList;
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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

#
# To add a data source extension:
#
## Start Extensions List .pm Anchor ## DO NOT REMOVE OR CHANGE THIS LINE
use Extensions::DATASOURCE_STUB;
## End Extensions List .pm Anchor ## DO NOT REMOVE OR CHANGE THIS LINE

sub new {
	my $self  = {};
	my $class = shift;
	my $args  = shift;

	return ( bless( $self, $class ) );
}

# Note: Should we create a dispatch table here?
sub hash_ref {
	my $ext_hash_ref = {
## Start Extensions List hash Anchor ## DO NOT REMOVE OR CHANGE THIS LINE
		'stats_long_term' => 'Extensions::DATASOURCE_STUB::stats_long_term',
		'get_config'      => 'Extensions::DATASOURCE_STUB::get_config',
## End Extensions List hash Anchor ## DO NOT REMOVE OR CHANGE THIS LINE
	};

	return $ext_hash_ref;
}
1;
