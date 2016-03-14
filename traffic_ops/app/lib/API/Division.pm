package API::Division;

use UI::Utils;
use UI::Division;
use Mojo::Base 'Mojolicious::Controller';
use Data::Dumper;
use JSON;
use MojoPlugins::Response;


sub create{
    my $self = shift;
    my $params = $self->req->json;
    if (!defined($params)) {
        return $self->alert("parameters must be in JSON format,  please check!");
    }

    if ( !&is_oper($self) ) {
          return $self->alert( { Error => " - You must be an ADMIN or OPER to perform this operation!" } );
    }

    my $name = $params->{name};
    if (!defined($name)) {
        return $self->alert("division 'name' is not given.");
    }

    #Check for duplicate division name
    my $existing_division = $self->db->resultset('Division')->search( { name => $name } )->get_column('name')->single();
    if ($existing_division) {
        return $self->alert("A division with name \"$name\" already exists." );
    }

    my $insert = $self->db->resultset('Division')->create( { name => $name } );
    $insert->insert();

    my $response;
    my $rs = $self->db->resultset('Division')->find( { id => $insert->id } );
    if (defined($rs)) {
        $response->{id}     = $rs->id;
        $response->{name}   = $rs->name;
        return $self->success($response);
    }
    return $self->alert("create division failed.");
}

1;
