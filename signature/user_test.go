package signature

import (
	"testing"
	"time"

	"github.com/xuzhenglun/workflow/database"
)

func Test_NewLicse(T *testing.T) {
	group := []string{"shit", "holy", "start"}

	db := database.NewMongoDB("")
	sign := NewSigner("../keys.json", db)

	user := sign.NewUser("Reficul", group, time.Now().Add(24*30*time.Hour))

	l := sign.NewLicese(user)

	T.Log(l)

	if !sign.Verify([]byte(l), "shit") {
		T.Error("Can't Verify")
	}
}
