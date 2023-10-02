package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const PeerGRPCTimeout = 30 * time.Second

type FabricConfigSigner struct {
	// PeerName     string `json:"peerName"`
	MSPID        string `json:"mspID"`
	IdentityPath string `json:"identityPath"`
	KeyPath      string `json:"keyPath"`
}

type FabricConfigPeer struct {
	PeerName        string `json:"peerName"`
	PeerAddress     string `json:"peerAddress"`
	TLSRootCertFile string `json:"tlsRootCertFile"`
	TLSKeyFile      string `json:"tlsKeyFile"`
	TLSCertFile     string `json:"tlsCertFile"`
}

func LoadConfig() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	certBytes, err := ioutil.ReadFile(filepath.Join(cwd, "testdata", "cert.pem"))
	if err != nil {
		return err
	}
	keyBytes, err := ioutil.ReadFile(filepath.Join(cwd, "testdata", "key.pem"))
	if err != nil {
		return err
	}

	fmt.Println("Cert: ", certBytes)
	fmt.Println("key: ", keyBytes)

	return nil
}
