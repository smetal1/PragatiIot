# Navigate to mosquitto/certs directory
cd mosquitto/certs

# Create a Certificate Authority (CA)
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 -out ca.crt -subj "/CN=CA"

# Create Server Certificates
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=mosquitto"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256
