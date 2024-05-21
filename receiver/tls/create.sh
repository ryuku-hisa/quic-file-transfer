#!/bin/bash

# 証明書の情報を設定
COUNTRY="JP"
STATE="Tokyo"
LOCALITY="Tokyo"
ORGANIZATION="Ina-lab"
COMMON_NAME="localhost"

# OpenSSL設定ファイルを生成
cat <<EOL > openssl.cnf
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext
x509_extensions    = v3_req # The extentions to add to the self signed cert

[ req_distinguished_name ]
countryName                 = Country Name (2 letter code)
countryName_default         = $COUNTRY
stateOrProvinceName         = State or Province Name (full name)
stateOrProvinceName_default = $STATE
localityName                = Locality Name (eg, city)
localityName_default        = $LOCALITY
organizationName            = Organization Name (eg, company)
organizationName_default    = $ORGANIZATION
commonName                  = Common Name (eg, fully qualified host name)
commonName_default          = $COMMON_NAME

[ req_ext ]
subjectAltName = @alt_names

[ v3_req ]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[ alt_names ]
DNS.1   = $COMMON_NAME
EOL

# 秘密鍵を生成
openssl genrsa -out server.key 2048

# 証明書署名要求(CSR)を生成
openssl req -new -key server.key -out server.csr -config openssl.cnf -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORGANIZATION/CN=$COMMON_NAME"

# 自己署名証明書を生成
openssl x509 -req -in server.csr -signkey server.key -out server.crt -extensions v3_req -extfile openssl.cnf

# 不要なファイルを削除
rm server.csr openssl.cnf

echo "Successfully created certificate: server.crt, server.key"
