# Start 
- node1: 
make node1

- node2: 
make node2

- node3: 
make node3


# Query data 
curl -L http://127.0.0.1:9121/test -XPUT -d test2

curl -L http://127.0.0.1:9121/test
curl -L http://127.0.0.1:9122/test
curl -L http://127.0.0.1:9123/test

curl -L http://127.0.0.1:9124/test
# Note
- App: Application using raft (example)
- API:  from etdc/api
- client: 
- lib 
- server
- v3: etcd/raft/v3 