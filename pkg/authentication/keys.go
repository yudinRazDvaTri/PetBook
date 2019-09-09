package authentication

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/dpgolang/PetBook/pkg/logger"
	"github.com/gorilla/securecookie"
)

type rsaKeys struct {
	SignKey   *rsa.PrivateKey
	VerifyKey *rsa.PublicKey
}

var RsaKeys = rsaKeys{}
var hashKey = securecookie.GenerateRandomKey(32)
var SCookie = securecookie.New(hashKey, nil)

func init() {
	reader := rand.Reader
	bitSize := 2048
	var err error

	RsaKeys.SignKey, err = rsa.GenerateKey(reader, bitSize)
	if err != nil {
		logger.FatalError(err, "Error occurred while trying to generate rsa rsaKeys.\n")
	}
	RsaKeys.VerifyKey = &RsaKeys.SignKey.PublicKey
}
