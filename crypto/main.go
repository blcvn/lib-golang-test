package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/blcvn/lib-golang-test/crypto/schnorr"

	"github.com/blcvn/lib-golang-test/crypto/btcec"
)

func main() {
	AggregateSignatures()
}

func Sign() {
	var message [32]byte

	privateKey, _ := new(big.Int).SetString("B7E151628AED2A6ABF7158809CF4F3C762E7160F38B4DA56A784D9045190CFEF", 16)
	msg, _ := hex.DecodeString("243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89")
	copy(message[:], msg)

	signature, err := schnorr.Sign(privateKey, message)
	if err != nil {
		fmt.Printf("The signing is failed: %v\n", err)
	}
	fmt.Printf("The signature is: %x\n", signature)
}
func Verify() {
	var (
		publicKey [33]byte
		message   [32]byte
		signature [64]byte
	)

	pk, _ := hex.DecodeString("02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659")
	copy(publicKey[:], pk)
	msg, _ := hex.DecodeString("243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89")
	copy(message[:], msg)
	sig, _ := hex.DecodeString("2A298DACAE57395A15D0795DDBFD1DCB564DA82B0F269BC70A74F8220429BA1D1E51A22CCEC35599B8F266912281F8365FFC2D035A230434A1A64DC59F7013FD")
	copy(signature[:], sig)

	if result, err := schnorr.Verify(publicKey, message, signature); err != nil {
		fmt.Printf("The signature verification failed: %v\n", err)
	} else if result {
		fmt.Println("The signature is valid.")
	}
}
func BatchVerify() {
	var (
		publicKey  [33]byte
		message    [32]byte
		signature  [64]byte
		publicKeys [][33]byte
		messages   [][32]byte
		signatures [][64]byte
	)

	pk, _ := hex.DecodeString("02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659")
	copy(publicKey[:], pk)
	publicKeys = append(publicKeys, publicKey)
	pk, _ = hex.DecodeString("03FAC2114C2FBB091527EB7C64ECB11F8021CB45E8E7809D3C0938E4B8C0E5F84B")
	copy(publicKey[:], pk)
	publicKeys = append(publicKeys, publicKey)
	pk, _ = hex.DecodeString("026D7F1D87AB3BBC8BC01F95D9AECE1E659D6E33C880F8EFA65FACF83E698BBBF7")
	copy(publicKey[:], pk)
	publicKeys = append(publicKeys, publicKey)
	msg, _ := hex.DecodeString("243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89")
	copy(message[:], msg)
	messages = append(messages, message)
	msg, _ = hex.DecodeString("5E2D58D8B3BCDF1ABADEC7829054F90DDA9805AAB56C77333024B9D0A508B75C")
	copy(message[:], msg)
	messages = append(messages, message)
	msg, _ = hex.DecodeString("B2F0CD8ECB23C1710903F872C31B0FD37E15224AF457722A87C5E0C7F50FFFB3")
	copy(message[:], msg)
	messages = append(messages, message)
	sig, _ := hex.DecodeString("2A298DACAE57395A15D0795DDBFD1DCB564DA82B0F269BC70A74F8220429BA1D1E51A22CCEC35599B8F266912281F8365FFC2D035A230434A1A64DC59F7013FD")
	copy(signature[:], sig)
	signatures = append(signatures, signature)
	sig, _ = hex.DecodeString("00DA9B08172A9B6F0466A2DEFD817F2D7AB437E0D253CB5395A963866B3574BE00880371D01766935B92D2AB4CD5C8A2A5837EC57FED7660773A05F0DE142380")
	copy(signature[:], sig)
	signatures = append(signatures, signature)
	sig, _ = hex.DecodeString("68CA1CC46F291A385E7C255562068357F964532300BEADFFB72DD93668C0C1CAC8D26132EB3200B86D66DE9C661A464C6B2293BB9A9F5B966E53CA736C7E504F")
	copy(signature[:], sig)
	signatures = append(signatures, signature)

	if result, err := schnorr.BatchVerify(publicKeys, messages, signatures); err != nil {
		fmt.Printf("The signature verification failed: %v\n", err)
	} else if result {
		fmt.Println("The signature is valid.")
	}
}

func AggregateSignatures() {
	var (
		publicKey [33]byte
		message   [32]byte
	)

	privateKey1, _ := new(big.Int).SetString("B7E151628AED2A6ABF7158809CF4F3C762E7160F38B4DA56A784D9045190CFEF", 16)
	privateKey2, _ := new(big.Int).SetString("C90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B14E5C7", 16)
	msg, _ := hex.DecodeString("243F6A8885A308D313198A2E03707344A4093822299F31D0082EFA98EC4E6C89")
	copy(message[:], msg)

	privateKeys := []*big.Int{privateKey1, privateKey2}
	signature, _ := schnorr.AggregateSignatures(privateKeys, message)

	Curve := btcec.S256()

	// verifying an aggregated signature
	pk, _ := hex.DecodeString("02DFF1D77F2A671C5F36183726DB2341BE58FEAE1DA2DECED843240F7B502BA659")
	copy(publicKey[:], pk)
	P1x, P1y := schnorr.Unmarshal(Curve, publicKey[:])

	pk, _ = hex.DecodeString("03FAC2114C2FBB091527EB7C64ECB11F8021CB45E8E7809D3C0938E4B8C0E5F84B")
	copy(publicKey[:], pk)
	P2x, P2y := schnorr.Unmarshal(Curve, publicKey[:])
	Px, Py := Curve.Add(P1x, P1y, P2x, P2y)

	copy(publicKey[:], schnorr.Marshal(Curve, Px, Py))

	if result, err := schnorr.Verify(publicKey, message, signature); err != nil {
		fmt.Printf("The signature verification failed: %v\n", err)
	} else if result {
		fmt.Println("The signature is valid.")
	}
}
