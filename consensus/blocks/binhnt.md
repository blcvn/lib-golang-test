# Folder Structure 
- App:  new code to support start server 
- client: Client call direct grpc to test 
- comm: start connection server 
- consensus:  (Get all from fabric )
    + common: 
    + etcdraft: support raft 
    + protoutil: support tool 

- broadcast: start service of orderer 

- grpcServers: 
    - Broadcast Server: support for client send broadcast and delivery message 
    - Cluster Server 
    - Cluster Node Server 