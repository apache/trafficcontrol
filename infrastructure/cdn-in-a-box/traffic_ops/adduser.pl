#!/usr/bin/env perl

use strict;
use Crypt::ScryptKDF qw{ scrypt_hash };

if ($#ARGV < 2) {
    die "Usage: $ARGV[0] <username> <password> <role>\n";
}

my $username = shift // 'admin';
my $password = shift or die "Password is required\n";
my $role = shift // 'admin';

# Skip the insert if the admin 'username' is already there.
my $hashed_passwd = hash_pass( $password );
print <<"ADMIN";
insert into tm_user (username, role, local_passwd, confirm_local_passwd)
    values  ('$username',
            (select id from role where name = '$role'),
            '$hashed_passwd',
            '$hashed_passwd' )
    ON CONFLICT (username) DO NOTHING;
ADMIN

sub hash_pass {
    my $pass = shift;
    return scrypt_hash($pass, \64, 16384, 8, 1, 64);
}
