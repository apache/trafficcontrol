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

TO_DIR="/opt/traffic_ops/app"; export TO_DIR
TO_EXT_PRIVATE_LIB="/opt/traffic_ops_extensions/private/lib"; export TO_PRIVATE_LIB
PERL5LIB=$TO_EXT_PRIVATE_LIB:$TO_DIR/lib:$TO_DIR/local/lib/perl5:$PERL5LIB; export PERL5LIB

# Setup GOPATH 
GOPATH="/opt/traffic_ops/go"; export GOPATH
GOBIN=$GOPATH/bin

# Setup PATH
PATH=$PATH:$GOBIN:/usr/local/go/bin
