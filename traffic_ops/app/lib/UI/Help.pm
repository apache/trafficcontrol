package UI::Help;
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

use UI::Utils;
use Mojo::Base 'Mojolicious::Controller';

sub about {
	my $self = shift;

	my $tool_name = $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.toolname', config_file => 'global' ] } )
		->get_column('value')->single();
	my $tool_instance = $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.instance_name', config_file => 'global' ] } )
		->get_column('value')->single();
	my $tool_version = &tm_version();
	my $tool_info_url = $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.infourl', config_file => 'global' ] } )
		->get_column('value')->single();
	my $tool_logo_url = $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.logourl', config_file => 'global' ] } )
		->get_column('value')->single();
	my $git_rev = `git rev-list HEAD -1`;

	$self->stash(
		tool_name          => $tool_name,
		tool_version       => $tool_version,
		tool_instance	   => $tool_instance,
		tool_info_url      => $tool_info_url,
		tool_logo_url      => $tool_logo_url,
		tool_mojo_version  => $Mojolicious::VERSION,
		tool_git_rev       => $git_rev,
		# tool_db_type       => (split(/:/, $Schema::dsn))[1], 
		tool_db_type       => $Schema::dsn, 
	);

	&navbarpage($self);
}
sub releasenotes {
	my $self = shift;

	my $tool_name = $self->db->resultset('Parameter')->search( { -and => [ name => 'tm.toolname', config_file => 'global' ] } )
		->get_column('value')->single();
	my $tool_version = &tm_version();
	# open( my $fh, '<', 'doc/releasenotes.txt' );

	my $file = '../doc/releasenotes.txt';
	open (FH, "< $file") or die "Can't open $file for read: $!";
	my @lines;
	while (<FH>) {
    	push (@lines, $_);
	}
	my $rn_text = join(" ", @lines);
	close FH or die "Cannot close $file: $!";

 
	$self->stash(
		tool_name          => $tool_name,
		tool_version       => $tool_version,
		rn_text 		   => $rn_text
	);

	&navbarpage($self);
}
1;
