package signature

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/xuzhenglun/workflow/algo"
)

type Database interface {
	IsInBlackList(string, string) bool
	Revocate(string, string) error
}

type Signature struct {
	Key    interface{}
	Method string
	Crypto algo.Algorithm
	DB     Database
}

func NewSigner(keys string, db Database) *Signature {
	var s Signature

	s.DB = db

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
