node1:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go run main.go  --port 9001 --id 1 --cluster 127.0.0.1:8001,127.0.0.1:8002,127.0.0.1:8003  --raftport 8001

node2:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go run main.go  --port 9002 --id 2 --cluster 127.0.0.1:8001,127.0.0.1:8002,127.0.0.1:8003 --raftport 8002

node3:
	GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go run main.go  --port 9003 --id 3 --cluster 127.0.0.1:8001,127.0.0.1:8002,127.0.0.1:8003 --raftport 8003
