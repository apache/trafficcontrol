package UI::UploadHandler;
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

use Data::Dumper;
use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';
use Mojo::Base 'Mojolicious::Plugin';
use UI::Server;
use JSON;
use Mojo::JSON;
use Mojo::Upload;
use Mojo::Log;

sub upload {
  my $self = shift;
  my $serverPath = '/tmp/';
  my $url = $self->req->url->to_abs;
  my $userinfo = $self->req->url->to_abs->userinfo;
  my $host = $self->req->url->to_abs->host;
  #return $self->render(json => Dumper($self->req));
  my $upload = $self->param('file-0');
  return $self->render_exception('[.csv files only] Upload size exceeded the 5MByte(5242880 bytes) limit!')
    if $upload->size > 5242880;
  my $upload2 = $upload->move_to($serverPath . $upload->filename);
  return $self->render(json => "{\"success\":true,\"serverpath\":\"" . $serverPath . "\",\"filename\":\"" . $upload->filename . "\",\"size\":" . $upload->size . "}");
}

1;

