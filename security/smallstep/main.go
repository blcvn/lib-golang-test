package main

import (
	"log"
	"os"

	smallstep "github.com/blcvn/lib-golang-test/security/smallstep/lib"

	"github.com/smallstep/cli/utils"
	"go.step.sm/cli-utils/errs"
)

func main() {
	// 1. create gstate certificate if not existed yet
	gstateName := "gstate01"
	_, err := os.ReadFile(gstateName + ".crt")
	if err != nil {
		log.Printf("error while read cert file %s => create a new one", gstateName+".crt")

		if err := smallstep.GenerateCertificate(gstateName); err != nil {
			log.Panicf("error while GenerateCertificate: %s", err.Error())
		}
	}

	// 2. register rootCA for chaincode (init)
	crtFile := "./root_ca.crt"
	crtBytes, err := utils.ReadFile(crtFile)
	if err != nil {
		log.Panicf("error while read Priv file: %s", errs.FileError(err, crtFile).Error())
	}
	chaincodeService.registerRootCA(crtBytes)

	// 3. Gstate signs on UTXO using its priv key
	if err := InitGstate(gstateName+".key", gstateName+".crt"); err != nil {
		log.Panicf("error while InitGstate: %s", err.Error())
	}

	newUTXO := chaincodeService.getMintUTXO() // fake UTXO
	signature, err := gstate.Sign(&newUTXO)
	if err != nil {
		log.Panicf("error while gstate.Sign: %s", err.Error())
	}
	newUTXO.Signature = signature
	newUTXO.SignerCert = gstate.signer.GetCertBytes()

	// 4. Chaincode verify signature
	if err := chaincodeService.verify(&newUTXO, gstateName); err != nil {
		log.Panicf("error while verify: %s", err.Error())
	}

	log.Println("perfect")
}
