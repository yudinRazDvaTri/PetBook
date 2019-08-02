package init

import (
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"os"
	"test/pkg/utils"
)

const (
	privKeyPath = "./keys/app.rsa"
	pubKeyPath  = "./keys/app.rsa.pub"
)

func init() {

	file, err := os.OpenFile("petbook.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	utils.Logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		utils.Error(err, "Error occurred while trying to read private key from file.\n")
	}

	utils.SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		utils.Error(err, "Error occurred while trying to parse private key from file.\n")
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		utils.Error(err, "Error occurred while trying to read public key from file.\n")
	}

	utils.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		utils.Error(err, "Error occurred while trying to read private key from file.\n")
	}
}
