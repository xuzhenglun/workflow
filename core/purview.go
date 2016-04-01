package core

import (
	"errors"
	"strings"

	"github.com/yuin/gopher-lua"
)

type Purview interface {
	Verify(req []byte, target string) bool
	GetGroupsOfUser(auth string) *[]string
}

func (this *VMs) SetAuther(p Purview) {
	this.Auth = p
}

func (this *VMs) GetPurviewActivities(auth string) ([]string, error) {
	g := this.Auth.GetGroupsOfUser(auth)
	if g == nil {
		return nil, errors.New("Have no any purview")
	}

	group := make(map[string]bool)
	for _, v := range *g {
		group[v] = true
	}

	activities := []string{}
	l := lua.NewState()
	defer l.Close()
	l.DoString(this.Scripts)
	for k, _ := range this.Activities {
		activity := FindActivityByName(l, k)
		if activity.Groups == "" && !strings.Contains(activity.Name, "_reader") {
			activities = append(activities, activity.Name)
		} else {
			if v, ok := group[activity.Groups]; v && ok {
				activities = append(activities, k)
			}
		}
	}
	return activities, nil
}
