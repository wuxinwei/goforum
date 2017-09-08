package services

import (
	"testing"

	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/google/go-cmp/cmp"
	"github.com/wuxinwei/goforum/models"
)

func TestCacheService(t *testing.T) {
	c := NewCacheService()
	user := models.User{
		Username: "wuxinwei",
		Name:     "markwoo",
		Birthday: time.Now(),
		Country:  "china",
		Province: "Hunan",
		City:     "Loudi",
	}
	if err := c.Set(user.Username, &user, Test); err != nil {
		t.Errorf("Set Action, Want: nil, Got: %s", err)
	}
	var newUser models.User
	if err := c.Get(user.Username, &newUser, Test); err != nil {
		t.Errorf("Get Action, Want: nil, Got: %s", err)
	}
	if !cmp.Equal(&user, &newUser) {
		t.Errorf("No equal, user: %#v, newUser: %#v", user, newUser)
	}
	if err := c.Del(user.Username, Test); err != nil {
		t.Errorf("Del Action, Want: nil, Got: %s", err)
	} else if err := c.Get(user.Username, &newUser, Test); err != nil {
		if err != redis.ErrNil {
			t.Errorf("Get Action, Want: nil, Got: %s", err)
		}
	} else {
		t.Errorf("Delete failed")
	}
}
