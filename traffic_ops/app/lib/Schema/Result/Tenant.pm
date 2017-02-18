use utf8;
package Schema::Result::Tenant;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::Tenant

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<tenant>

=cut

__PACKAGE__->table("tenant");

=head1 ACCESSORS

=head2 id

  data_type: 'bigint'
  is_nullable: 0

=head2 name

  data_type: 'text'
  is_nullable: 0

=head2 last_updated

  data_type: 'timestamp with time zone'
  default_value: current_timestamp
  is_nullable: 1
  original: {default_value => \"now()"}

=cut

__PACKAGE__->add_columns(
  "id",
  { data_type => "bigint", is_nullable => 0 },
  "name",
  { data_type => "text", is_nullable => 0 },
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

=head1 UNIQUE CONSTRAINTS

=head2 C<tenant_name_key>

=over 4

=item * L</name>

=back

=cut

__PACKAGE__->add_unique_constraint("tenant_name_key", ["name"]);

=head1 RELATIONS

=head2 cdns

Type: has_many

Related object: L<Schema::Result::Cdn>

=cut

__PACKAGE__->has_many(
  "cdns",
  "Schema::Result::Cdn",
  { "foreign.tenant_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 deliveryservices

Type: has_many

Related object: L<Schema::Result::Deliveryservice>

=cut

__PACKAGE__->has_many(
  "deliveryservices",
  "Schema::Result::Deliveryservice",
  { "foreign.tenant_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);

=head2 tm_users

Type: has_many

Related object: L<Schema::Result::TmUser>

=cut

__PACKAGE__->has_many(
  "tm_users",
  "Schema::Result::TmUser",
  { "foreign.tenant_id" => "self.id" },
  { cascade_copy => 0, cascade_delete => 0 },
);


# Created by DBIx::Class::Schema::Loader v0.07046 @ 2017-02-18 09:32:59
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:GwD6l7apYu+ouFMo7Jv/Tg


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
