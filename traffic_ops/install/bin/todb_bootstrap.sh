#!/usr/bin/env bash

# To bypass the password prompts for automation, please set TODB_USERNAME_PASSWORD=<yourpassword> before you invoke

# Example:
#
#    $ TODB_USERNAME_PASSWORD=<yourpassword> ./todb_bootstrap.sh
#
TODB_USERNAME=traffic_ops
TODB_NAME=traffic_ops

if [[ -z $TODB_USERNAME ]]; then
    echo "Using environment database user: $TODB_USERNAME"
fi

if [[ -z $TODB_USERNAME_PASSWORD ]]; then
   while true; do
    read -s -p "Please ENTER the new password for database user '$TODB_USERNAME': " password
    echo
    read -s -p "Please CONFIRM enter the new password for database user '$TODB_USERNAME' again: " password_confirm
    echo
    [ "$password" = "$password_confirm" ] && break
    echo "Passwords do not match, please try again"
   done
   TODB_USERNAME_PASSWORD=$password
else
    echo "Using environment database password"
fi
echo "Setting up database role: $TODB_USERNAME"
psql -U postgres -h localhost -c "CREATE USER $TODB_USERNAME WITH ENCRYPTED PASSWORD '$TODB_USERNAME_PASSWORD';"
createdb $TODB_NAME --owner $TODB_USERNAME -U postgres -h localhost

echo "Successfully set up database '$TODB_NAME' with role '$TODB_USERNAME'"
