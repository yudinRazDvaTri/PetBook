package authentication

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/dpgolang/PetBook/pkg/logger"
)

type keys struct {
	SignKey   *rsa.PrivateKey
	VerifyKey *rsa.PublicKey
}

var Keys = keys{}

func init() {
	reader := rand.Reader
	bitSize := 2048
	var err error

	Keys.SignKey, err = rsa.GenerateKey(reader, bitSize)
	if err != nil {
		logger.FatalError(err, "Error occurred while trying to generate rsa keys.\n")
	}
	Keys.VerifyKey = &Keys.SignKey.PublicKey

	// Saving keys to directory. Currently it is not needed.
	/*
	   func savePEMKey(fileName string, key *rsa.PrivateKey) {
	   	outFile, err := os.Create(fileName)
	   	if err != nil {
	   		logger.FatalError(err)
	   	}
	   	defer outFile.Close()

	   	var privateKey = &pem.Block{
	   		Type:  "PRIVATE KEY",
	   		Bytes: x509.MarshalPKCS1PrivateKey(key),
	   	}

	   	err = pem.Encode(outFile, privateKey)
	   	if err != nil {
	   		logger.FatalError(err)
	   	}
	   }

	   func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	   	asn1Bytes, err := asn1.Marshal(pubkey)
	   	if err != nil {
	   		logger.FatalError(err)
	   	}

	   	var pemkey = &pem.Block{
	   		Type:  "PUBLIC KEY",
	   		Bytes: asn1Bytes,
	   	}

	   	pemfile, err := os.Create(fileName)
	   	if err != nil {
	   		logger.FatalError(err)
	   	}
	   	defer pemfile.Close()

	   	err = pem.Encode(pemfile, pemkey)
	   	if err != nil {
	   		logger.FatalError(err)
	   	}
	*/
}
