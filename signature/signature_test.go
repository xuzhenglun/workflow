package signature

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func Test_SignAndVerify(t *testing.T) {
	signer := NewSigner("../keys.json")

	reader := bufio.NewReader(rand.Reader)
	var buff [102400]byte

	for i := 0; i < len(buff); i++ {
		buff[i], _ = reader.ReadByte()
	}
	sig := signer.Crypto.Sign(buff[:]).ToBytes()
	if i := signer.Crypto.Verify(buff[:], sig); i == nil {
		t.Log("Sign Pass")
	}
}

func Test_SignAndVerifyWithBase64(t *testing.T) {
	signer := NewSigner("../keys.json")

	reader := bufio.NewReader(rand.Reader)
	var buff [102400]byte

	for i := 0; i < len(buff); i++ {
		buff[i], _ = reader.ReadByte()
	}

	sig := signer.Crypto.Sign(buff[:]).ToBase64()
	sigbyte, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		t.Error("Fail")
	}
	if i := signer.Crypto.Verify(buff[:], sigbyte); i == nil {
		t.Log("Sign Pass")
	}
}
