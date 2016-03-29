package signature

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"sort"
	"time"
)

const SUPERUSER = "root"

type Groups []*string

type User struct {
	Name       string
	Group      Groups
	Signture   string
	SignTime   time.Time
	ExpireTime time.Time
}

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

func (this Groups) String() string {
	ret := ""
	for _, v := range this {
		ret = ret + *v + ", "
	}
	return ret
}

func (this User) sum() []byte {
	sum := this.Name + this.SignTime.String() + this.ExpireTime.String()
	for _, v := range this.Group {
		sum = sum + *v
	}
	return []byte(sum)
}

func (this Signature) IsCrtValid(user User) bool {
	sort.Sort(user.Group)
	sum := user.sum()
	sigbyte, err := base64.StdEncoding.DecodeString(user.Signture)
	if err != nil {
		log.Println(err)
		return false
	}
	if err := this.Crypto.Verify(sum, sigbyte); err == nil {
		log.Println("Package Signature:: IsCrtValid: PASS")
		return true
	} else {
		log.Println("Package Signature:: IsCrtValid:", err)
		return false
	}
}

func (this Signature) IsBelong(user User, group string) bool {
	if !this.IsCrtValid(user) {
		return false
	}
	for _, v := range user.Group {
		if *v == group || *v == SUPERUSER {
			return true
		}
	}
	log.Println("Package Signature:: the user is not at target group")
	log.Println("User:", user.Group, "target:", group)
	return false
}

func (this Signature) Verify(requestHead []byte, targetGroup string) bool {
	//未分组的表示该Activity无需权限，任何人都可访问
	if targetGroup == "" {
		return true
	}

	req, err := base64.StdEncoding.DecodeString(string(requestHead))
	if err != nil {
		log.Println("Package Signature::Verify:", err)
		return false
	}

	var user User
	if err := json.Unmarshal(req, &user); err != nil {
		log.Println("Package Signature::Verify:", err)
		return false
	}

	return this.IsBelong(user, targetGroup)
}

func (this Signature) NewUser(name string, groups []string, expire time.Time) User {
	user := User{Name: name, ExpireTime: expire}
	user.Group = make([]*string, len(groups))
	for i, _ := range groups {
		user.Group[i] = &(groups[i])
	}
	sort.Sort(user.Group)
	return user
}

func (this Signature) NewLicese(user User) string {
	user.SignTime = time.Now()
	user.Signture = this.Crypto.Sign(user.sum()).ToBase64()

	j, err := json.Marshal(user)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString([]byte(j))
}
