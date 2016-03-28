package user

import (
	"encoding/base64"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/xuzhenglun/workflow/signature"
)

type Groups []*string

type User struct {
	Name     string
	Group    Groups
	Signture string
	Time     time.Time
}

type Auth struct{ signature.Signature }

func (this Groups) Len() int {
	return len(this)
}

func (this Groups) Less(i, j int) bool {
	return (*this[i]) < (*this[j])
}

func (this Groups) Swap(i, j int) {
	temp := this[i]
	this[i] = this[j]
	this[j] = temp
}

func (this User) sum() []byte {
	sum := this.Name
	for _, v := range this.Group {
		sum = sum + *v
	}
	sum = sum + fmt.Sprintf("%v", this.Time)
	return []byte(sum)
}

func (this Auth) IsValid(user User) bool {
	sort.Sort(user.Group)
	sum := user.sum()
	sigbyte, err := base64.StdEncoding.DecodeString(user.Signture)
	if err != nil {
		log.Println(err)
		return false
	}
	if err := this.Crypto.Verify(sum, sigbyte); err == nil {
		return true
	} else {
		log.Println(err)
		return false
	}
}

func (this Auth) IsBelong(user User, group string) bool {
	if !this.IsValid(user) {
		return false
	}
	for _, v := range user.Group {
		if *v == group {
			return true
		}
	}
	return false
}
