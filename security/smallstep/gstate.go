package main

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	smallstep "github.com/blcvn/lib-golang-test/security/smallstep/lib"

	"github.com/smallstep/cli/utils"
	"go.step.sm/cli-utils/errs"
)

type MyGStateService struct {
	signer *smallstep.Signer
}

var (
	gstate *MyGStateService
)

func InitGstate(privPath, crtFile string) error {
	gstate = &MyGStateService{}

	// private key
	privBytes, err := os.ReadFile(privPath)
	if err != nil {
		log.Printf("error while read Priv file: %s", err.Error())
		return err
	}
	block, _ := pem.Decode(privBytes)
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	// certificate
	crtBytes, err := utils.ReadFile(crtFile)
	if err != nil {
		log.Printf("error while read Crt file: %s", errs.FileError(err, crtFile).Error())
		return err
	}

	gstate.signer = smallstep.InitSigner(privateKey, crtBytes)
	return nil
}

func (g *MyGStateService) Sign(utxo *UTXO) ([]byte, error) {
	return g.signer.Sign(utxo.Hash)
}
