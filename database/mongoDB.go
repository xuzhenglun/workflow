package database

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuzhenglun/workflow/core"
)

type MongoDB struct {
	Db       *mgo.Collection
	BL       *mgo.Collection
	mux      sync.Mutex
	Activity *core.Activity
}

func NewMongoDB(URI string) MongoDB {
	var db MongoDB
	session, err := mgo.Dial(URI)
	if err != nil {
		log.Panic(err)
	}
	db.Db = session.DB("workflow").C("process")
	db.BL = session.DB("workflow").C("blacklist")
	return db
}

func (this MongoDB) AddRow(args ...map[string]string) error {
	toAdd := bson.M{}
	toAdd["_id"] = bson.NewObjectId()
	for i, _ := range args {
		for k, v := range args[i] {
			toAdd[k] = v
		}
	}
	if _, ok := toAdd["pass"]; !ok {
		toAdd["pass"] = "0"
	}
	return this.Db.Insert(toAdd)
}

func (this MongoDB) ModifyRow(args ...map[string]string) error {
	id := bson.ObjectIdHex(args[1]["id"])
	qurry := this.Db.FindId(id)
	var saved bson.M
	if err := qurry.One(&saved); err != nil {
		log.Println(err)
		return err
	}
	for i, _ := range args {
		for k, v := range args[i] {
			saved[k] = v
		}
	}
	delete(saved, "id")
	if saved["pass"] == "true" {
		saved["pass"] = "1"
	}
	if saved["pass"] == "false" {
		saved["pass"] = "0"
	}

	return this.Db.UpdateId(id, saved)
}

func (this MongoDB) FindRow(id string, needArgs ...string) (string, error) {
	var ret []string
	qurry := this.Db.FindId(bson.ObjectIdHex(id))
	saved := bson.M{}

	if err := qurry.One(&saved); err != nil {
		log.Println(err)
		return "", err
	}
	j := simplejson.New()
	for _, v := range needArgs {
		j.Set(v, saved[v])
	}

	if r, err := j.Encode(); err != nil {
		log.Println(err)
		return "", err
	} else {
		ret = append(ret, string(r))
	}

	if l := len(ret); l < 1 {
		return "[]", nil
	} else if l == 1 {
		return ret[0], nil
	} else if ret, err := json.Marshal(ret); err != nil {
		log.Println(err)
		return "", err
	} else {
		return string(ret), nil
	}
}

func (this MongoDB) DeleteRow(id ...string) error {
	for _, i := range id {
		Id := bson.ObjectIdHex(i)
		if err := this.Db.RemoveId(Id); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (this MongoDB) GetJustDone(id string) (string, string, error) {
	qurry := this.Db.FindId(bson.ObjectIdHex(id))
	saved := bson.M{}

	if err := qurry.One(&saved); err != nil {
		log.Println(err)
		return "", "", DBError{When: time.Now(), What: "Can't find target"}
	} else {
		justdone, ok1 := saved["JustDone"]
		pass, ok2 := saved["pass"]
		if ok1 && ok2 {
			return justdone.(string), pass.(string), nil
		} else {
			return justdone.(string), pass.(string), DBError{When: time.Now(), What: "Empty response"}
		}
	}
}

func (this MongoDB) GetList(action string, pass string) (string, error) {
	req := bson.M{"JustDone": action, "pass": pass}
	qurry := this.Db.Find(req)

	list := []struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	}{}
	if err := qurry.All(&list); err != nil {
		log.Println(err)
		return "", err
	} else {
		var l []bson.ObjectId
		for _, v := range list {
			l = append(l, v.Id)
		}
		if ret, err := json.Marshal(l); err != nil {
			log.Println(err)
			return "", err
		} else {
			return string(ret), nil
		}
	}
}

func (this MongoDB) IsInBlackList(name, sig string) bool {
	n, err := this.BL.Find(bson.M{name: sig}).Count()
	log.Println("--->", n, err)
	if err != nil {
		log.Println(err)
		return false
	}
	if n < 1 {
		return false
	}
	return true
}

func (this MongoDB) Revocate(name, sig string) error {
	return this.BL.Insert(bson.M{name: sig})
}
