package cluster

import "bytes"

// CertificateComparator returns whether some relation holds for two given certificates
type CertificateComparator func([]byte, []byte) bool

type MemberMapping struct {
	id2stub       map[uint64]*Stub
	SamePublicKey CertificateComparator
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

// LookupByClientCert retrieves a Stub with the given client certificate
func (mp MemberMapping) LookupByClientCert(cert []byte) *Stub {
	for _, stub := range mp.id2stub {
		if mp.SamePublicKey(stub.ClientTLSCert, cert) {
			return stub
		}
	}
	return nil
}

// LookupByIdentity retrieves a Stub by Identity
func (mp MemberMapping) LookupByIdentity(identity []byte) *Stub {
	for _, stub := range mp.id2stub {
		if bytes.Equal(identity, stub.Identity) {
			return stub
		}
	}
	return nil
}

// ServerCertificates returns a set of the server certificates
// represented as strings
func (mp MemberMapping) ServerCertificates() StringSet {
	res := make(StringSet)
	for _, member := range mp.id2stub {
		res[string(member.ServerTLSCert)] = struct{}{}
	}
	return res
}
