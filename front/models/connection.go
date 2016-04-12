package model

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type NeedArgs struct {
	Args []string
}

type ActivityInfo struct {
	Code int
	Msg  map[string]string
}

type User struct {
	Name             string
	SignTime         time.Time
	ExpireTime       time.Time
	SignTimeString   string
	ExpireTimeString string
}

type Activities []string

func GetArgs(url, auth string) *NeedArgs {
	body := GetResponse(url, auth)

	var args NeedArgs

	if err := json.Unmarshal(body, &args); err != nil {
		log.Println(err)
		return nil
	}
	return &args
}

func GetInfo(url, auth string) *ActivityInfo {
	body := GetResponse(url, auth)

	var info ActivityInfo
	log.Println("BODY:", string(body))

	if err := json.Unmarshal(body, &info); err != nil {
		log.Println(err)
		return nil
	}
	return &info
}

func GetAllActivites(url, auth string) *Activities {
	body := GetResponse(url, auth)

	var activities Activities
	log.Println("BODY:", string(body))

	if err := json.Unmarshal(body, &activities); err != nil {
		log.Println(err)
		return nil
	}
	for k, v := range activities {
		log.Println(k, v)
	}

	sort.Sort(SortableString(activities))
	return &activities
}

func GetResponse(url, auth string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	req.Header.Add("auth", auth)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return body
}

func GetUserInfo(auth string) *User {
	str, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		log.Println(err)
		return nil
	}

	var user User
	if err = json.Unmarshal(str, &user); err != nil {
		log.Println(err)
		return nil
	}

	user.ExpireTimeString = user.ExpireTime.Format(time.RFC822)
	user.SignTimeString = user.SignTime.Format(time.RFC822)

	return &user
}

func Post(url string, r *http.Request) ([]byte, error) {
	//r.ParseForm()
	req, err := http.NewRequest("POST", url, strings.NewReader(r.Form.Encode()))
	if err != nil {
		return nil, err
	}

	auth, err := r.Cookie("auth")
	if err != nil {
		return nil, err
	}
	req.Header.Add("auth", auth.Value)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}
