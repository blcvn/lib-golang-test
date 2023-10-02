package transport

import (
	"google.golang.org/grpc"
)

// ConnByCertMap maps certificates represented as strings
// to gRPC connections
type ConnByCertMap map[string]*grpc.ClientConn

// Lable used for TLS Export Keying Material call
const KeyingMaterialLabel = "orderer v3 authentication label"

// Lookup looks up a certificate and returns the connection that was mapped
// to the certificate, and whether it was found or not
func (cbc ConnByCertMap) Lookup(cert []byte) (*grpc.ClientConn, bool) {
	conn, ok := cbc[string(cert)]
	return conn, ok
}

// Put associates the given connection to the certificate
func (cbc ConnByCertMap) Put(cert []byte, conn *grpc.ClientConn) {
	cbc[string(cert)] = conn
}

// Remove removes the connection that is associated to the given certificate
func (cbc ConnByCertMap) Remove(cert []byte) {
	delete(cbc, string(cert))
}

// Size returns the size of the connections by certificate mapping
func (cbc ConnByCertMap) Size() int {
	return len(cbc)
}

// MemberMapping defines NetworkMembers by their ID
// and enables to lookup stubs by a certificate
type MemberMapping struct {
	id2stub map[uint64]*Stub
}

func NewMemberMapping() *MemberMapping {
	return &MemberMapping{
		id2stub: map[uint64]*Stub{},
	}
}

// Foreach applies the given function on all stubs in the mapping
func (mp *MemberMapping) Foreach(f func(id uint64, stub *Stub)) {
	for id, stub := range mp.id2stub {
		f(id, stub)
	}
}

// Put inserts the given stub to the MemberMapping
func (mp *MemberMapping) Put(stub *Stub) {
	mp.id2stub[stub.ID] = stub
}

// Remove removes the stub with the given ID from the MemberMapping
func (mp *MemberMapping) Remove(ID uint64) {
	delete(mp.id2stub, ID)
}

// ByID retrieves the Stub with the given ID from the MemberMapping
func (mp MemberMapping) ByID(ID uint64) *Stub {
	return mp.id2stub[ID]
}

// // LookupByIdentity retrieves a Stub by Identity
// func (mp MemberMapping) LookupByIdentity(identity []byte) *Stub {
// 	for _, stub := range mp.id2stub {
// 		if bytes.Equal(identity, stub.Identity) {
// 			return stub
// 		}
// 	}
// 	return nil
// }
