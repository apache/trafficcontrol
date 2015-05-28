package Extensions::ConfigList;
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
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';

#
# To add a config file extension:
#
# 1) add a parameter like the below, and associate it in the profile of the server:
# +-----+----------+--------------------+---------------+---------------------+
# | id  | name     | config_file        | value         | last_updated        |
# +-----+----------+--------------------+---------------+---------------------+
# | 875 | location | to_ext_ttt.config  | /opt/ttt/etc  | 2015-02-01 12:31:55 |
# +-----+----------+--------------------+---------------+---------------------+
#
# This will create the file ttt.config (note the to_ext_ prefix in the parameter)
#
# 2) Create a .pm file that has the perl sub to generate the config file. This sub will have
# access to all the $self->db style things. Return the text for the config file in this sub
#
# 3) Add the use line for your .pm in the list below
#
# 4) add the filename (without the to_ext_ prefix) the $ext_hash_ref below

## Start Extensions List .pm Anchor ## DO NOT REMOVE OR CHANGE THIS LINE
## End Extensions List .pm Anchor ## DO NOT REMOVE OR CHANGE THIS LINE

sub hash_ref {
	my $ext_hash_ref = {
## Start Extensions List hash Anchor ## DO NOT REMOVE OR CHANGE THIS LINE
## End Extensions List hash Anchor ## DO NOT REMOVE OR CHANGE THIS LINE
	};

	return $ext_hash_ref;
}

1;
