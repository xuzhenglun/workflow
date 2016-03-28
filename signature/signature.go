package signature

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/xuzhenglun/workflow/algo"
)

type Signature struct {
	Key    interface{}
	Method string
	Crypto algo.Algorithm
}

func NewSigner(keys string) *Signature {
	var s Signature

	f, err := ioutil.ReadFile(keys)
	if err != nil {
		log.Panic(err)
	}

	if err := json.Unmarshal(f, &s); err != nil {
		log.Panic(err)
	}
	if s.Crypto, err = algo.Init[s.Method](s.Key); err != nil {
		log.Panic(err)
	}
	return &s
}
