chmod +x generate-cert.sh
./generate-cert.sh

curl -v -k https://localhost:8443/get -H "Host: httpbin.org"

echo ""
echo "Press any key to continue..."
read -n 1
