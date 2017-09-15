package services

import (
	"testing"

	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fpay/gopress"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/wuxinwei/goforum/conf"
	"github.com/wuxinwei/goforum/models"
)

var tl *gopress.Logger

func init() {
	tl = gopress.NewLogger()
	tl.SetLevel(log.DEBUG)
	tl.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}

func TestConnect(t *testing.T) {
	s := NewElasticService(tl)
	if s == nil {
		t.Logf("Want non-nil, Got nil")
	}
}

func TestMapping(t *testing.T) {
	s := NewElasticService(tl)
	if s == nil {
		t.Logf("Want non-nil, Got nil")
	}
	if err := s.Mapping("test", conf.PostMapping); err != nil {
		t.Errorf("Want nil, Got %s", err)
	}
}

func TestIndexing(t *testing.T) {
	TestMapping(t)
	s := NewElasticService(tl)
	if s == nil {
		t.Errorf("Want non-nil, Got nil")
	}

	sample := models.Post{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    1,
		Title:     gofakeit.HackerNoun(),
		Author:    gofakeit.Name(),
		Content:   gofakeit.HackerPhrase(),
	}
	if err := s.Indexing(&sample, "test", "post"); err != nil {
		t.Errorf("want nil, Got %s", err)
	}
}

func TestSearchPost(t *testing.T) {
	s := NewElasticService(tl)
	if s == nil {
		t.Errorf("Want non-nil, Got nil")
	}
	posts, err := s.SearchPost("content", "input asdfwq", true)
	if err != nil {
		t.Errorf("Want non-nil, Got %s", err)
	}
	t.Logf("posts: %#v", posts)
}
