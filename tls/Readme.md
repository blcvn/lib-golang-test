# 1. Start server 
go run main.go  serve  --port 5001 --name=client.vnpay.vn

# 2. Start client 

GODEBUG="x509ignoreCN=0" go run main.go  echo localhost --name=client.vnpay.vn --port=5001 