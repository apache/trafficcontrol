package UI::Render;
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

require Exporter;
our @ISA = qw(Exporter);

use constant READ  => 10;
use constant OPER  => 20;
use constant ADMIN => 30;

our %EXPORT_TAGS = ( 'all' => [ qw(build_json) ] );
our @EXPORT_OK = ( @{ $EXPORT_TAGS{all} } );
our @EXPORT    = ( @{ $EXPORT_TAGS{all} } );

sub build_json{
   my $self = shift;
   my $rs_data = shift;
   my $default_columns = shift;
   my $columns;

   if ( defined $self->param('columns') ) {
     $columns = $self->param('columns');
   }
   else {
     $columns = $default_columns;
   }

   my (@columns) = split(/,/, $columns);
   my %columns;
   foreach my $col (@columns) {
     $columns{$col} = defined;
   }

   my @data;
   my @cols = grep { exists $columns{$_} } $rs_data->result_source->columns;     

   while ( my $row = $rs_data->next ) {
     my %parameter;
     foreach my $col ( @cols ) {
        $parameter{$col}=$row->$col;
     }
     push (@data, \%parameter);
    }
    return \@data;
}

1;
