#!/bin/bash

# Create directory for certificates if it doesn't exist
mkdir -p certs
cd certs

# Generate private key
openssl genrsa -out key.pem 2048

# Generate CSR (Certificate Signing Request)
# Note: Format adjusted for Windows compatibility
openssl req -new -key key.pem -out csr.pem \
    -subj "//CN=localhost\O=Test Proxy\C=US"

# Generate self-signed certificate
openssl x509 -req -days 365 \
    -in csr.pem \
    -signkey key.pem \
    -out cert.pem

# Clean up CSR
rm -f csr.pem

echo "Generated certificates in certs directory:"
ls -l cert.pem key.pem
echo ""
echo "Press any key to continue..."
read -n 1
