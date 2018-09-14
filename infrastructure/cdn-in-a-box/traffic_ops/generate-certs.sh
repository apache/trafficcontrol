#!/usr/bin/env bash 
################################################################################
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
# 
#   http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
################################################################################
#
# Name: generate_certs.sh
#
# Purpose: Generate a full self-signed CA and associated x509 certificates
#           and keys needed for a full CDN-In-A-Box 
# 
# Expected ENVIRONMENT variables for this script to run properly
# 
# CERT_COUNTRY        - x509 subject line 
# CERT_COMPANY        - x509 subject line
# CERT_STATE          - x509 subject line
# CERT_CITY           - x509 subject line
# CERT_ORG            - x509 subject line 
# CERT_KEYSIZE        - RSA key size (1024, 2048, 4096) 
# CERT_DURATION       - Certificate Duration in days 
# CERT_CA_NAME        - Certificate Authority Short Name
# CERT_DIR            - Directory name under the enroller dir that contains the
#                       generated SSL cert/key pairs such as "ssl"
#
# Additional variables created by this script will be in:
#    $ENROLLER_DIR + / + $ENROLLER_SSL_DIR + / + $CERT_ENV_FILE
#
# CERT_DOMAIN         - Copy of $DOMAIN
# CERT_CA_CERT_FILE   - Absolute path to CA Certificate
# CERT_CA_KEY_FILE    - Absolute path to CA Key
# CERT_ENV_FILE       - Source of additional variables for other containers.
# CERT_DONE_FILE      - Path to poll to know when SSL generation is complete.
# CERT_STATIC_HOSTS   - Space delimited infra host list
# CERT_WILDCARD_HOSTS - Space delimited wildcard host list
# CERT_{NAME}_CERT    - Component certificate (ex. NAME could be TRAFFICOPS)
# CERT_{NAME}_KEY     - Component certificate
#
################################################################################
# ToDo:
################################################################################
# 1) Ability to run the script with cmdline arguments
# 2) Ability to create more certs without re-creating another CA
################################################################################

# Extra that need to be constructed
CERT_DOMAIN="$DOMAIN"
CERT_LOGS_DIR="$CERT_DIR/logs"

CERT_CA_CERT_FILE="$CERT_DIR/$CERT_CA_NAME.$DOMAIN.crt"
CERT_CA_KEY_FILE="$CERT_DIR/$CERT_CA_NAME.$DOMAIN.key"

# Hosts that will have SSL certs generated
CERT_STATIC_HOSTS="$DB_HOST $EDGE_HOST $MID_HOST $ORIGIN_HOST $ENROLLER_HOST $TM_HOST $TO_HOST $TO_PERL_HOST $TP_HOST $TM_HOST $TV_HOST $TR_HOST $TS_HOST"
CERT_WILDCARD_HOSTS="$CDN_DOMAINS"

#
# Funciton create_ssl_cir
# Params: none
#
create_output_dir()
{
  rm -rf "$CERT_DIR"
  mkdir -p "$CERT_DIR" "$CERT_LOGS_DIR"

  ret=$?
  
  if ((ret!=0)) ; then
    echo "ERROR: Can't Create SSL Certificate directory [$CERT_DIR] (retcode=$ret)"
    exit $ret
  fi

  return $ret
}

#
# Function: create_key
# Params: no arguments
#
create_key()
{
  local host="$1"
  local filename="$host.$CERT_DOMAIN.key"
  local keypath="$CERT_DIR/$filename"

  echo "Generating RSA Key"
  echo " - Size=$CERT_KEYSIZE"
  echo " - Path=$keypath"

  openssl genrsa -out "$keypath" "$CERT_KEYSIZE" >"$CERT_LOGS_DIR/create_key_$host.log" 2>&1

  ret=$?

  if ((ret!=0)) ; then
    echo " ! ERROR: Can't create RSA key [$keypath] with size [$CERT_KEYSIZE] (retcode=$ret)"     
    exit $ret
  fi

  echo " + Complete" 

  return $ret
}

#
# Function: create_ca
# Params: no arguments
#
create_ca() 
{

  create_key "$CERT_CA_NAME" 

  openssl req -x509 -new -nodes -key "$CERT_CA_KEY_FILE" \
     -sha256 -days $CERT_DURATION -out "$CERT_CA_CERT_FILE" \
     -subj "/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_ORG ($CERT_DOMAIN)" >"$CERT_LOGS_DIR/create_ca_${CERT_CA_NAME}.log" 2>&1

  ret=$?

  echo "Creating x509 Certificate Authority"
  echo " - Name=$CERT_CA_NAME"
  echo " - Domain=$CERT_DOMAIN"
  echo " - CertPath=$CERT_CA_CERT_FILE"
  echo " - KeyPath=$CERT_CA_KEY_FILE"
  if ((ret!=0));then
    echo " ! ERROR: Can't create CA certificate at [$CERT_DIR] (retcode=$ret)"
    exit $ret
  fi

  echo " + Complete."

  >| "$CERT_ENV_FILE"

  return $ret
}

# 
# Function: create_cert
# Description: Creates key, certificate request, and certificate.  
#              Request is signed by the previously generated CA
# Params:
#   1) short hostname
#   2) wildcard 
#
create_cert()
{
  local host="$1"
  local do_wc="$2"
  local cn="$host.$CERT_DOMAIN"
  local request="$CERT_DIR/$host.$CERT_DOMAIN.csr"
  local cert="$CERT_DIR/$host.$CERT_DOMAIN.crt"
  local key="$CERT_DIR/$host.$CERT_DOMAIN.key"
  local uc_cert="$(echo $host | tr '[a-z]' '[A-Z]' | tr '-' '_' )_CERT=\"$cert\""
  local uc_key="$(echo $host | tr '[a-z]' '[A-Z]' | tr '-' '_' )_KEY=\"$key\""

  if [ "$do_wc" = "wildcard" ] ; then 
     cn="*.$host.$CERT_DOMAIN"
  fi

  create_key "$host"

  echo "Creating x509 Certificate Signing Request"
  echo " - CN=$cn"
  echo " - Path=$request"
	  
  openssl req -new -key $key -out $request \
     -subj "/CN=$cn/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_ORG/OU=$host" 
  ret=$? > "$CERT_LOGS_DIR/cert_req_${cn}.log" 2>&1

  if ((ret!=0)) ; then
    echo " ! ERROR: Can't create x509 certificate request for [$host] in [$CERT_LOGS_DIR] (retcode=$ret)" 
    exit $ret
  fi 

  echo " + Complete."

  echo "Signing x509 Certificate Request"
  echo " - CN=$cn"
  echo " - CAfile=$CERT_CA_CERT_FILE"

  openssl x509 -req -in "$request" \
    -CA "$CERT_CA_CERT_FILE" \
    -CAkey "$CERT_CA_KEY_FILE" -CAcreateserial \
    -out "$cert" \
    -days "$CERT_DURATION" -sha256 > "$CERT_LOGS_DIR/sign_cert_${cn}.log" 2>&1

  ret=$?

  if ((ret!=0)) ; then
    echo " ! ERROR: Can't sign x509 certificate [$host] with CA in [$CERT_DIR] (retcode=$ret)" 
    exit $ret
  fi 

  echo " + Complete."

  # Save cert info to CA environment file
  echo "CERT_${uc_cert}" >> $CERT_ENV_FILE
  echo "CERT_${uc_key}" >> $CERT_ENV_FILE

  return $ret
}

verify_cert()
{
  local host="$1"
  local fqdn="$host.$CERT_DOMAIN"
  local cn="$host.$CERT_DOMAIN"
  local cert="$CERT_DIR/$host.$CERT_DOMAIN.crt"

  openssl verify -CAfile "$CERT_CA_CERT_FILE" "$cert" 
  ret=$?

  if ((ret!=0)) ; then
    echo "ERROR: Can't verify certificate [$fqdn] against CA [$CERT_CA_CERT_FILE] (retcode=$ret)"
    return $ret
  fi

  return $ret
}

print_cert_banner()
{
  echo "--------------------------------------------------------------------------------"
  echo "Creating Shared CA for CDN-In-A-Box"
  echo "--------------------------------------------------------------------------------"
  echo "Static Certificates:"
  for chost in $CERT_STATIC_HOSTS
  do
    echo "  + CN=$chost.$CERT_DOMAIN" 
  done
  echo
  echo "Wildcard Certificates:"
  for whost in $CERT_WILDCARD_HOSTS
  do
    echo "  + CN=*.$whost.$CERT_DOMAIN" 
  done
  echo "--------------------------------------------------------------------------------"
}

save_cert_environment()
{
  local tmp_file="$(mktemp)"
  set | grep '^CERT' >> "$tmp_file"
  cat "$CERT_ENV_FILE" >> "$tmp_file"
  cat "$tmp_file" | sort > "$CERT_ENV_FILE" 
  cat "$tmp_file" | sed 's/^/export /' > "$CERT_ENV_FILE"
  rm -f "$tmp_file"
}

print_cert_environment()
{
  source "$CERT_ENV_FILE"
  echo "--------------------------------------------------------------------------------"
  echo "Final Shared Certificate Environment"
  echo "--------------------------------------------------------------------------------"
  set | grep ^CERT | sort
  echo "--------------------------------------------------------------------------------"
}

# Umask is wide open so all containers can read the certs regardless of uid/gid
umask 0000

print_cert_banner
create_output_dir
create_ca

# Create static certs for one or more infra hosts
for chost in $CERT_STATIC_HOSTS
do
  create_cert "$chost"
done

# Create wildcard certs for one or more CDN sub-domains
for whost in $CERT_WILDCARD_HOSTS
do
  create_cert "$whost" "wildcard"
done


# Validate all generated certs against self-signed CA
echo "--------------------------------------------------------------------------------"
echo "Cryptographically validating certs against CA [$CERT_CA_NAME] ..."
echo "--------------------------------------------------------------------------------"
for host in $CERT_STATIC_HOSTS $CDNS
do
  echo -en "Validating "
  verify_cert "$host"
done

save_cert_environment
print_cert_environment

touch "$CERT_DONE_FILE"

echo "--------------------------------------------------------------------------------"
echo "Shared CA and all SSL certificates have been successfully created."
echo "--------------------------------------------------------------------------------"
