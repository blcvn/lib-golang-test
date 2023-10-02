node1:
	go run main.go --id=1 --port=9121 --cluster="1;http://127.0.0.1:9021,2;http://127.0.0.1:9022,3;http://127.0.0.1:9023" 
 
node2:
	go run main.go  --id=2  --port=9122  --cluster="1;http://127.0.0.1:9021,2;http://127.0.0.1:9022,3;http://127.0.0.1:9023"

node3:
	go run main.go  --id=3 --learn=true  --port=9123   --cluster="1;http://127.0.0.1:9021,2;http://127.0.0.1:9022,3;http://127.0.0.1:9023"
 
