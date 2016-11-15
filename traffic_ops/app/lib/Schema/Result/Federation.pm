use utf8;
package Schema::Result::Federation;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Federation

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<federation>

=cut

__PACKAGE__->table("federation");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_auto_increment: 1
  is_nullable: 0
  sequence: 'federation_id_seq'

=head2 cname

  data_type: 'varchar'
  is_nullable: 0
  size: 1024

=head2 description

  data_type: 'varchar'
  is_nullable: 1
  size: 1024

=head2 ttl

  data_type: 'integer'
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 1
  original: {default_value => \"now()"}

=cut

__PACKAGE__->add_columns(
  "id",
  {
    data_type         => "bigint",
    is_auto_increment => 1,
    is_nullable       => 0,
    sequence          => "federation_id_seq",
  },
  "cname",
  { data_type => "varchar", is_nullable => 0, size => 1024 },
  "description",
  { data_type => "varchar", is_nullable => 1, size => 1024 },
  "ttl",
  { data_type => "integer", is_nullable => 0 },
  "last_updated",
  {
    data_type     => "timestamp with time zone",
    default_value => \"current_timestamp",
    is_nullable   => 1,
    original      => { default_value => \"now()" },
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</id>

=back

=cut

__PACKAGE__->set_primary_key("id");

=head1 RELATIONS

=head2 federation_deliveryservices

Type: has_many

Related object: L<Schema::Result::FederationDeliveryservice>

=cut

__PACKAGE__->has_many(
  "federation_deliveryservices",
  "Schema::Result::FederationDeliveryservice",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_federation_resolvers

Type: has_many

Related object: L<Schema::Result::FederationFederationResolver>

=cut

__PACKAGE__->has_many(
  "federation_federation_resolvers",
  "Schema::Result::FederationFederationResolver",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 federation_tmusers

Type: has_many

Related object: L<Schema::Result::FederationTmuser>

=cut

__PACKAGE__->has_many(
  "federation_tmusers",
  "Schema::Result::FederationTmuser",
  { "foreign.federation" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07045 @ 2016-11-15 08:31:12
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:tiblLFygkfyk5wvu7f2yQw


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
