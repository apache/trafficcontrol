use utf8;
package Schema::Result::DeliveryservicePathPrefix;

# Created by DBIx::Class::Schema::Loader
# DO NOT MODIFY THE FIRST PART OF THIS FILE

=head1 NAME

Schema::Result::DeliveryservicePathPrefix

=cut

use strict;
use warnings;

use base 'DBIx::Class::Core';

=head1 TABLE: C<deliveryservice_path_prefix>

=cut

__PACKAGE__->table("deliveryservice_path_prefix");

=head1 ACCESSORS

=head2 deliveryservice

  data_type: 'integer'
  is_foreign_key: 1
  is_nullable: 0

=head2 path_prefix

  data_type: 'varchar'
  is_nullable: 0
  size: 255

=head2 last_updated

  data_type: 'timestamp'
  datetime_undef_if_invalid: 1
  default_value: current_timestamp
  is_nullable: 1

=cut

__PACKAGE__->add_columns(
  "deliveryservice",
  { data_type => "integer", is_foreign_key => 1, is_nullable => 0 },
  "path_prefix",
  { data_type => "varchar", is_nullable => 0, size => 255 },
  "last_updated",
  {
    data_type => "timestamp",
    datetime_undef_if_invalid => 1,
    default_value => \"current_timestamp",
    is_nullable => 1,
  },
);

=head1 PRIMARY KEY

=over 4

=item * L</deliveryservice>

=item * L</path_prefix>

=back

=cut

__PACKAGE__->set_primary_key("deliveryservice", "path_prefix");

=head1 RELATIONS

=head2 deliveryservice

Type: belongs_to

Related object: L<Schema::Result::Deliveryservice>

=cut

__PACKAGE__->belongs_to(
  "deliveryservice",
  "Schema::Result::Deliveryservice",
  { id => "deliveryservice" },
  { is_deferrable => 1, on_delete => "CASCADE", on_update => "CASCADE" },
);


# Created by DBIx::Class::Schema::Loader v0.07043 @ 2016-11-11 14:33:49
# DO NOT MODIFY THIS OR ANYTHING ABOVE! md5sum:JTCGdZIfVw1GqE2JGW94YA


# You can replace this text with custom code or comments, and it will be preserved on regeneration
1;
