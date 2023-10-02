rm gen/*.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 3650 -nodes -keyout gen/ca-key.pem -out gen/ca-cert.pem -subj "/C=VN/ST=CA/L=VCB/O=Vnpay/OU=Blockchain/CN=*.vnpay.vn/emailAddress=admin@vnpay.vn"

echo "CA's self-signed certificate"
openssl x509 -in gen/ca-cert.pem -noout -text

#########  Generate Key and Certificate for each nodes
declare -a arr=("gate1" "gate2" "gen1" "gen2" "gen3")

## now loop through the above array
for name in "${arr[@]}"
do
    echo "Generating certificate and key for $name"

    # 1. Generate  private key and certificate signing request (CSR)
    openssl req -newkey rsa:4096 -nodes -keyout gen/$name-key.pem -out gen/$name-req.pem -subj "/C=VN/ST=Server/L=VCB/O=Vnpay/OU=Blockchain/CN=*.vnpay.vn/emailAddress=admin@vnpay.vn"

    # 2. Use CA's private key to sign web server's CSR and get back the signed certificate
    # openssl x509 -req -in gen/$name-req.pem -days 3650 -CA gen/ca-cert.pem -CAkey gen/ca-key.pem -CAcreateserial -out gen/$name-cert.pem -extfile ./server-ext.cnf
    openssl x509 -req -extfile ./vnpay-ext.cnf -in gen/$name-req.pem -days 3650 -CA gen/ca-cert.pem -CAkey gen/ca-key.pem -CAcreateserial -out gen/$name-cert.pem  

    # 3. Verify certficate 
    echo "Verifying signed certificate of $name "
    openssl x509 -in gen/$name-cert.pem -noout -text
    openssl verify -CAfile gen/ca-cert.pem gen/$name-cert.pem
done 