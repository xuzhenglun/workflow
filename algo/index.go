package algo

type Algorithm interface {
	Encrypt([]byte) Result
	Decrypt([]byte) Result
	Sign([]byte) Result
	Verify([]byte, []byte) error
}

var method Algorithm
var Init = map[string]func(interface{}) (Algorithm, error){
	"RSA": RsaInit,
}
