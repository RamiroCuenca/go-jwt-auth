package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"sync"

	"github.com/golang-jwt/jwt/v4"
)

var (
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	once      sync.Once
)

// We left them private (Only can be used by this package)
// so that the other packages can not manipulate them

// Singleton that assign a value to this vars only once (EXPORTED)
func LoadCertificates(privateFile, publicFile string) error {
	var err error

	once.Do(func() {
		err = loadCertificates(privateFile, publicFile)
	})

	return err
}

// Loads the certicates and send them to be parsed
func loadCertificates(privateFile, publicFile string) error {
	// It's cool to use ioutil.ReadFile beacause it reads the file and
	// close it after it got the value.
	// If we use os.Open, we should close manually the document.
	privateBytes, err := ioutil.ReadFile(privateFile)
	if err != nil {
		return err
	}

	publicBytes, err := ioutil.ReadFile(publicFile)
	if err != nil {
		return err
	}

	return parseRSA(privateBytes, publicBytes)
}

// Parces the certicates
func parseRSA(privateBytes, publicBytes []byte) error {
	var err error
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		return err
	}

	return nil
}
