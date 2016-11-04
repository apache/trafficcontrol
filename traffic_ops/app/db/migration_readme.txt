#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

Create a goose migration:
Information on goose can be found here:  https://bitbucket.org/liamstask/goose

1. Download Goose ->  go get bitbucket.org/liamstask/goose/cmd/goose
2. In the db directory create a dbconf.yml file. 
  -> A sample should be there
  -> NOTE:  dbconf.yml CANNOT contain tabs!
3. from the /opt/tm directory create your migration
	$ goose -env=mysql create foober sql
	goose: created /opt/tm/db/migrations/20141006092229_foober.sql
4. Add your database migration under the goose up section.  Add any rollback under the goose down section.
5. run 'goose -env=mysql up' to perform the migration
