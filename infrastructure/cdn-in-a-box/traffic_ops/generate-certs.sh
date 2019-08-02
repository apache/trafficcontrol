#!/bin/bash
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

export X509_CA_DEFAULT_NAME="ca"
export X509_CA_DEFAULT_COUNTRY="ZZ"
export X509_CA_DEFAULT_STATE="SomeState"
export X509_CA_DEFAULT_CITY="SomeCity"
export X509_CA_DEFAULT_COMPANY="SomeCompany"
export X509_CA_DEFAULT_ORG="SomeOrganization"
export X509_CA_DEFAULT_ORGUNIT="SomeOrgUnit"
export X509_CA_DEFAULT_EMAIL="no-reply@some.host.test"
export X509_CA_DEFAULT_DIGEST="sha256"
export X509_CA_DEFAULT_DURATION_DAYS="3650"
export X509_CA_DEFAULT_KEYTYPE="rsa"
export X509_CA_DEFAULT_KEYSIZE=4096
export X509_CA_DEFAULT_UMASK=0002
export X509_CA_DEFAULT_DIR="$PWD/ca-default"

x509v3_init()
{
  if [[ $X509_CA_INITIALIZED = true ]] ; then
    echo "ERROR: Already initialized."
    return 2
  fi

  export X509_CA_DIR="${X509_CA_DIR:-$X509_CA_DEFAULT_DIR}"
  export X509_CA_NAME="${X509_CA_NAME:-$X509_CA_DEFAULT_NAME}"
  export X509_CA_COUNTRY="${X509_CA_COUNTRY:-$X509_CA_DEFAULT_COUNTRY}"
  export X509_CA_STATE="${X509_CA_STATE:-$X509_CA_DEFAULT_STATE}"
  export X509_CA_CITY="${X509_CA_CITY:-$X509_CA_DEFAULT_CITY}"
  export X509_CA_COMPANY="${X509_CA_COMPANY:-$X509_CA_DEFAULT_COMPANY}"
  export X509_CA_ORG="${X509_CA_ORG:-$X509_CA_DEFAULT_ORG}"
  export X509_CA_ORGUNIT="${X509_CA_ORGUNIT:-$X509_CA_DEFAULT_ORGUNIT}"
  export X509_CA_EMAIL="${X509_CA_EMAIL:-$X509_CA_DEFAULT_EMAIL}"
  export X509_CA_DIGEST="${X509_CA_DIGEST:-$X509_CA_DEFAULT_DIGEST}"
  export X509_CA_DURATION_DAYS="${X509_CA_DURATION_DAYS:-$X509_CA_DEFAULT_DURATION_DAYS}"
  export X509_CA_KEYTYPE="${X509_CA_KEYTYPE:-$X509_CA_DEFAULT_KEYTYPE}"
  export X509_CA_KEYSIZE="${X509_CA_KEYSIZE:-$X509_CA_DEFAULT_KEYSIZE}"
  export X509_CA_UMASK="${X509_CA_UMASK:-$X509_CA_DEFAULT_UMASK}"
  export X509_CA_INITIALIZED=true

  # Runtime
  export X509_CA_ROOT_CERT_FILE="$X509_CA_DIR/${X509_CA_NAME}-root.crt"
  export X509_CA_ROOT_KEY_FILE="$X509_CA_DIR/${X509_CA_NAME}-root.key"
  export X509_CA_ROOT_CONFIG_FILE="$X509_CA_DIR/${X509_CA_NAME}-root.config"
  export X509_CA_ROOT_DB_FILE="$X509_CA_DIR/${X509_CA_NAME}-root.db"
  export X509_CA_ROOT_REVOKE_FILE="$X509_CA_DIR/${X509_CA_NAME}-root.crl"
  export X509_CA_ROOT_SERIAL_FILE="$X509_CA_DIR/${X509_CA_NAME}-root.serial"

  export X509_CA_INTR_REQUEST_FILE="$X509_CA_DIR/${X509_CA_NAME}-intr.csr"
  export X509_CA_INTR_CERT_FILE="$X509_CA_DIR/${X509_CA_NAME}-intr.crt"
  export X509_CA_INTR_KEY_FILE="$X509_CA_DIR/${X509_CA_NAME}-intr.key"
  export X509_CA_INTR_CONFIG_FILE="$X509_CA_DIR/${X509_CA_NAME}-intr.config"
  export X509_CA_INTR_DB_FILE="$X509_CA_DIR/${X509_CA_NAME}-intr.db"
  export X509_CA_INTR_REVOKE_FILE="$X509_CA_DIR/${X509_CA_NAME}-intr.crl"
  export X509_CA_INTR_SERIAL_FILE="$X509_CA_DIR/${X509_CA_NAME}-intr.serial"

  export X509_CA_CERT_FULL_CHAIN_FILE="$X509_CA_DIR/${X509_CA_NAME}-fullchain.crt"
  export X509_CA_ENV_FILE="$X509_CA_DIR/environment"

  # If no X509_CA directory exists, create it
  if [ -d "$X509_CA_DIR" ] ; then
    echo "ERROR: Previous X509v3 CA Directory Exists."
    return 3
  fi

  # Set the Umask
  umask $X509_CA_UMASK

  # Create CA Certificate
  mkdir -p "$X509_CA_DIR"
  x509v3_create_root_ca
  x509v3_create_intermediate_ca

  # Create the CA environment file
  x509v3_dump_env

  return $?
}

x509v3_create_root_ca()
{
  # Use today's epoch date for the first serial number
  echo "$(date +%s)" > "$X509_CA_ROOT_SERIAL_FILE"

  # Create the CA index file
  touch "$X509_CA_ROOT_DB_FILE"

  local cn="$X509_CA_ORG Root CA"
  local ca_name="$X509_CA_NAME-root"

  local ca_config=""`
  `"[ca]\n"`
  `"default_ca = $ca_name\n\n"`
  `"[$ca_name]\n"`
  `"new_certs_dir = $X509_CA_DIR\n"`
  `"certificate = $X509_CA_ROOT_CERT_FILE\n"`
  `"private_key = $X509_CA_ROOT_KEY_FILE\n"`
  `"serial = $X509_CA_ROOT_SERIAL_FILE\n"`
  `"database = $X509_CA_ROOT_DB_FILE\n"`
  `"default_md = $X509_CA_DIGEST\n"`
  `"default_days = $X509_CA_DURATION_DAYS\n"`
  `"prompt = no\n"`
  `"preserve = no\n\n"`
  `"[policy_match]\n"`
  `"countryName = match\n"`
  `"stateOrProvinceName = match\n"`
  `"organizationName = match\n"`
  `"organizationalUnitName = optional\n"`
  `"commonName = supplied\n"`
  `"emailAddress = optional\n\n"`
  `"[policy_anything]\n"`
  `"countryName = optional\n"`
  `"stateOrProvinceName = optional\n"`
  `"localityName = optional\n"`
  `"organizationName = optional\n"`
  `"organizationalUnitName = optional\n"`
  `"commonName = supplied\n"`
  `"emailAddress = optional\n\n"`
  `"[req]\n"`
  `"default_bits = $X509_CA_KEYSIZE\n"`
  `"default_md = $X509_CA_DIGEST\n"`
  `"default_days = $X509_CA_DURATION_DAYS\n"`
  `"distinguished_name = req_dn\n"`
  `"string_mask = nombstr\n"`
  `"x509_extensions = v3_ca\n\n"`
  `"[req_dn]\n"`
  `"countryName = Country Name (2 letter code)\n"`
  `"countryName_default = $X509_CA_COUNTRY\n"`
  `"countryName_min = 2\n"`
  `"countryName_max = 2\n"`
  `"stateOrProvinceName = State or Province Name (full name)\n"`
  `"stateOrProvinceName_default = $X509_CA_STATE\n"`
  `"localityName = Locality Name (eg, city)\n"`
  `"localityName_default = $X509_CA_CITY\n"`
  `"0.organizationName = Organization Name (eg, company)\n"`
  `"0.organizationName_default = $X509_CA_ORG\n"`
  `"organizationalUnitName = Organizational Unit Name (eg, section)\n"`
  `"organizationalUnitName_default = $X509_CA_ORGUNIT\n"`
  `"commonName = Common Name (eg, YOUR name)\n"`
  `"commonName_max = 64\n"`
  `"emailAddress = Email Address\n"`
  `"emailAddress_max = 64\n"`
  `"emailAddress_default = $X509_CA_EMAIL\n\n"`
  `"[v3_ca]\n"`
  `"subjectKeyIdentifier = hash\n"`
  `"authorityKeyIdentifier = keyid:always,issuer\n"`
  `"basicConstraints = critical, CA:true\n"`
  `"keyUsage = critical, digitalSignature, keyCertSign\n\n"`
  `"[v3_intermediate_ca]\n"`
  `"subjectKeyIdentifier = hash\n"`
  `"authorityKeyIdentifier = keyid:always,issuer\n"`
  `"basicConstraints = critical, CA:true, pathlen:0\n"`
  `"keyUsage = critical, digitalSignature, keyCertSign\n"

  echo "Writing Root CA Config File"
  echo -e "$ca_config" > "$X509_CA_ROOT_CONFIG_FILE"

  echo "Creating Root CA certificate for [$ca_name]."
  # Create new CA Certificate and Key
  openssl req -x509 -nodes -extensions v3_ca \
    -days "$((X509_CA_DURATION_DAYS*2))" \
    -config "$X509_CA_ROOT_CONFIG_FILE" \
    -newkey "$X509_CA_KEYTYPE:$X509_CA_KEYSIZE" \
    -keyout "$X509_CA_ROOT_KEY_FILE" \
    -out "$X509_CA_ROOT_CERT_FILE" \
    -subj "/C=$X509_CA_COUNTRY/ST=$X509_CA_STATE/L=$X509_CA_CITY/O=$X509_CA_ORG/OU=$X509_CA_ORG/CN=$cn/emailAddress=$X509_CA_EMAIL/" \
    > "$X509_CA_DIR/x509_create_root_ca.log" 2>&1

  retcode=$?

  echo "CreateX509 Root CA RetCode=$retcode"

  return $retcode
}

x509v3_create_intermediate_ca()
{
  # Use today's epoch date for the first serial number
  echo "$(date +%s)" > "$X509_CA_INTR_SERIAL_FILE"

  # Create the CA index file
  touch "$X509_CA_INTR_DB_FILE"

  local cn="$X509_CA_ORG Intermediate CA"
  local ca_name="$X509_CA_NAME-intr"

  local ca_config=""`
  `"[ca]\n"`
  `"default_ca = $ca_name\n\n"`
  `"[$ca_name]\n"`
  `"new_certs_dir = $X509_CA_DIR\n"`
  `"certificate = $X509_CA_INTR_CERT_FILE\n"`
  `"private_key = $X509_CA_INTR_KEY_FILE\n"`
  `"serial = $X509_CA_INTR_SERIAL_FILE\n"`
  `"database = $X509_CA_INTR_DB_FILE\n"`
  `"default_md = $X509_CA_DIGEST\n"`
  `"default_days = $X509_CA_DURATION_DAYS\n"`
  `"prompt = no\n"`
  `"preserve = no\n\n"`
  `"[policy_match]\n"`
  `"countryName = match\n"`
  `"stateOrProvinceName = match\n"`
  `"organizationName = match\n"`
  `"organizationalUnitName = optional\n"`
  `"commonName = supplied\n"`
  `"emailAddress = optional\n\n"`
  `"[policy_anything]\n"`
  `"countryName = optional\n"`
  `"stateOrProvinceName = optional\n"`
  `"localityName = optional\n"`
  `"organizationName = optional\n"`
  `"organizationalUnitName = optional\n"`
  `"commonName = supplied\n"`
  `"emailAddress = optional\n\n"`
  `"[req]\n"`
  `"default_bits = $X509_CA_KEYSIZE\n"`
  `"default_md = $X509_CA_DIGEST\n"`
  `"default_days = $X509_CA_DURATION_DAYS\n"`
  `"distinguished_name = req_dn\n"`
  `"string_mask = nombstr\n"`
  `"x509_extensions = v3_ca\n\n"`
  `"[req_dn]\n"`
  `"countryName = Country Name (2 letter code)\n"`
  `"countryName_default = $X509_CA_COUNTRY\n"`
  `"countryName_min = 2\n"`
  `"countryName_max = 2\n"`
  `"stateOrProvinceName = State or Province Name (full name)\n"`
  `"stateOrProvinceName_default = $X509_CA_STATE\n"`
  `"localityName = Locality Name (eg, city)\n"`
  `"localityName_default = $X509_CA_CITY\n"`
  `"0.organizationName = Organization Name (eg, company)\n"`
  `"0.organizationName_default = $X509_CA_ORG\n"`
  `"organizationalUnitName = Organizational Unit Name (eg, section)\n"`
  `"organizationalUnitName_default = $X509_CA_ORGUNIT\n"`
  `"commonName = Common Name (eg, YOUR name)\n"`
  `"commonName_max = 64\n"`
  `"emailAddress = Email Address\n"`
  `"emailAddress_max = 64\n"`
  `"emailAddress_default = $X509_CA_EMAIL\n\n"`
  `"[v3_ca]\n"`
  `"subjectKeyIdentifier = hash\n"`
  `"authorityKeyIdentifier = keyid:always,issuer\n"`
  `"basicConstraints = critical, CA:true\n"`
  `"keyUsage = critical, digitalSignature, keyCertSign\n\n"`
  `"[v3_intermediate_ca]\n"`
  `"subjectKeyIdentifier = hash\n"`
  `"authorityKeyIdentifier = keyid:always,issuer\n"`
  `"basicConstraints = critical, CA:true, pathlen:0\n"`
  `"keyUsage = critical, digitalSignature, keyCertSign\n"

  echo "Writing Intemediate CA openssl CA Config File"
  echo -e "$ca_config" > "$X509_CA_INTR_CONFIG_FILE"

  echo "Creating Intermediate CA certificate request for [$ca_name]."
  # Create new CA Certificate and Key
  openssl req -new -nodes \
    -days "$((X509_CA_DURATION_DAYS))" \
    -config "$X509_CA_INTR_CONFIG_FILE" \
    -newkey "$X509_CA_KEYTYPE:$X509_CA_KEYSIZE" \
    -keyout "$X509_CA_INTR_KEY_FILE" \
    -subj "/C=$X509_CA_COUNTRY/ST=$X509_CA_STATE/L=$X509_CA_CITY/O=$X509_CA_ORG/OU=$X509_CA_ORG/CN=$cn/emailAddress=$X509_CA_EMAIL/" \
    -out "$X509_CA_INTR_REQUEST_FILE"
    > "$X509_CA_DIR/x509_create_intermediate_csr.log" 2>&1
  retcode=$?
  echo "CreateX509 Intermediate CA RetCode=$retcode"

  echo "Signing x509v3 Intermediate CA with Root CA Certificate..."
  # Sign with the CA
  openssl ca -config "$X509_CA_ROOT_CONFIG_FILE" \
    -batch \
    -extensions v3_intermediate_ca \
    -policy policy_anything \
    -out "$X509_CA_INTR_CERT_FILE" \
    -infiles "$X509_CA_INTR_REQUEST_FILE" \
    > "$X509_CA_DIR/x509_sign_intermediate_ca.log" 2>&1

  retcode=$?
  echo "Sign X509 Req RetCode=$retcode"

  # Build X509 CA full chain 
  cat $X509_CA_INTR_CERT_FILE $X509_CA_ROOT_CERT_FILE > $X509_CA_CERT_FULL_CHAIN_FILE
  
  return $retcode
}

x509v3_gen_alt_names()
{
  local names="$1"

  [ -z "$name" ] && return 1

  local alt_names_text=""`
  `"[alt_names]\n"

  local count=1
  for n in $names
  do
    alt_names_text="${alt_names_text}DNS.$count = $n"
    ((count+=1))
  done

  echo "$alt_names_text"

  return 0
}

x509v3_create_cert()
{
  local name="$1"
  local cn="$2"
  local alt_names="$3"

  echo "name=[$name]"
  echo "cn=[$cn]"
  echo "alt_names=[$alt_names]"

  # TODO: Error Checking

  local config_file="$X509_CA_DIR/$cn.config"
  local exten_file="$X509_CA_DIR/$cn.exten"
  local cert_file="$X509_CA_DIR/$cn.crt"
  local key_file="$X509_CA_DIR/$cn.key"
  local request_file="$X509_CA_DIR/$cn.csr"
  local env_name=$(echo -e "$name" | tr '[a-z]' '[A-Z]' | sed 's/\./_/g')

  # CN is always included in SAN list as it is required by all modern web browsers.
  cn="*.${cn}"
  alt_names=$(x509v3_gen_alt_names "$cn $san_list")

  local request_config=""`
  `"[req]\n"`
  `"encrypt_key = no\n"`
  `"prompt = no\n"`
  `"utf8 = yes\n"`
  `"default_md = $X509_CA_DIGEST\n"`
  `"default_bits = $X509_CA_KEYSIZE\n"`
  `"distinguished_name = dn\n"`
  `"[dn]\n"`
  `"C = $X509_CA_COUNTRY\n"`
  `"ST = $X509_CA_STATE\n"`
  `"L = $X509_CA_CITY\n"`
  `"O  = $X509_CA_ORG\n"`
  `"CN = $cn\n\n"

  local exten_config=""`
  `"[v3_ext]\n"`
  `"basicConstraints = CA:FALSE\n"`
  `"nsCertType = server\n"`
  `"subjectKeyIdentifier = hash\n"`
  `"authorityKeyIdentifier = keyid,issuer:always\n"`
  `"keyUsage = critical, digitalSignature, keyEncipherment\n"`
  `"extendedKeyUsage = serverAuth\n"`
  `"subjectAltName=@alt_names\n\n"

  echo "Creating x509v3 request for cn=$cn type $type..."

  # Create OpenSSL config file this request
  echo -e "${request_config}${alt_names}" > "$config_file"
  echo -e "${exten_config}${alt_names}" > "$exten_file"

  # Create the x509 request config file
  openssl req -new -nodes \
    -config "$config_file" \
    -newkey "$X509_CA_KEYTYPE:$X509_CA_KEYSIZE" \
    -keyout "$key_file" \
    -out "$request_file" \
    > "$X509_CA_DIR/x509_create_request_${name}.log" 2>&1

  retcode=$?
  echo "Create X509 Req RetCode=$retcode"

  echo "Signing x509v3 request for cn=$cn type $type..."
  # Sign with the CA
  openssl ca -batch \
    -policy policy_anything \
    -config "$X509_CA_INTR_CONFIG_FILE" \
    -out "$cert_file" \
    -extensions v3_ext \
    -extfile "$exten_file" \
    -infiles "$request_file" \
    > "$X509_CA_DIR/x509_create_signrequest_${name}.log" 2>&1

  retcode=$?
  echo "Sign X509 Req RetCode=$retcode"

  echo "X509_${env_name}_CONFIG_FILE=\"$config_file\"" >> "$X509_CA_ENV_FILE"
  echo "X509_${env_name}_EXTEN_FILE=\"$exten_file\"" >> "$X509_CA_ENV_FILE"
  echo "X509_${env_name}_CERT_FILE=\"$cert_file\"" >> "$X509_CA_ENV_FILE"
  echo "X509_${env_name}_KEY_FILE=\"$key_file\"" >> "$X509_CA_ENV_FILE"
  echo "X509_${env_name}_REQUEST_FILE=\"$request_file\"" >> "$X509_CA_ENV_FILE"
}


x509v3_dump_env()
{
  touch "$X509_CA_ENV_FILE"
  tmp_file="$(mktemp)"
  cat "$X509_CA_ENV_FILE" > "$tmp_file"
  env | grep -E '^X509_' >> "$tmp_file"
  set | grep -E '^X509_' >> "$tmp_file"
  sort "$tmp_file" | uniq | sed 's/^/export /' > "$X509_CA_ENV_FILE"
  sync ; sleep 1
  rm -f "$tmp_file"
}
