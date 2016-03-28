package algo

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

type RSA struct {
	prikey *rsa.PrivateKey
	pubkey *rsa.PublicKey
}

type Result struct {
	R []byte
	E error
}

func RsaInit(key interface{}) (Algorithm, error) {
	var defaultrsa RSA
	keys, ok := key.(map[string]interface{})
	if !ok {
		log.Panic("Invalid RSA JSON")
	}
	prikey, ok1 := keys["Prikey"].(string)
	pubkey, ok2 := keys["Pubkey"].(string)
	if !ok1 || !ok2 {
		log.Panic("Invalid RSA Key")
	}

	defaultrsa.prikey = ParseKeys(prikey).(*rsa.PrivateKey)
	defaultrsa.pubkey = ParseKeys(pubkey).(*rsa.PublicKey)

	return &defaultrsa, nil
}

func (this RSA) Encrypt(origindata []byte) Result {
	return Result{}
}

func (this RSA) Decrypt([]byte) Result {
	return Result{}
}

func (this RSA) Sign(data []byte) Result {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	r, e := rsa.SignPKCS1v15(rand.Reader, this.prikey, crypto.SHA256, d)
	return Result{R: r, E: e}
}

func (this RSA) Verify(message, sig []byte) error {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(this.pubkey, crypto.SHA256, d, sig)
}

func ParseKeys(key string) interface{} {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		log.Panic(errors.New("Invalid RSA in Json, maybe wrong farmat?"))
	}
	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			log.Panic(err)
		}
		return key
	case "PRIVIATE KEY":
		fallthrough
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Panic(err)
		}
		return key

	default:
		log.Panic(fmt.Errorf("ssh: unsupported key type %q", block.Type))
	}
	return nil
}

func (this Result) ToBase64() string {
	return base64.StdEncoding.EncodeToString(this.R)
}

func (this Result) ToBytes() []byte {
	return this.R
}
