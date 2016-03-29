package signature

import (
	"testing"
	"time"
)

func Test_NewLicse(T *testing.T) {
	group := []string{"shit", "holy", "start"}

	sign := NewSigner("../keys.json")

	user := sign.NewUser("Reficul", group, time.Now().Add(24*30*time.Hour))

	l := sign.NewLicese(user)

	T.Log(l)

	if !sign.Verify([]byte(l), "shit") {
		T.Error("Can't Verify")
	}
}
