package main

import (
	"log"

	"github.com/fpay/gopress"
	"github.com/labstack/echo-contrib/session"
	elog "github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/wuxinwei/goforum/conf"
	_ "github.com/wuxinwei/goforum/conf"
	"github.com/wuxinwei/goforum/controllers"
	"github.com/wuxinwei/goforum/services"
	validator "gopkg.in/asaskevich/govalidator.v4"
)

func main() {
	validator.SetFieldsRequiredByDefault(true)
	// create server
	s := gopress.NewServer(gopress.ServerOptions{
		Port: 3000,
	})
	s.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	s.Logger.SetLevel(elog.DEBUG)

	// init and register services
	s.RegisterServices(
		services.NewDbService(),
		services.NewCacheService(),
	)

	store, err := s.App().Services.Get("cache").(*services.CacheService).SessionStore(conf.SessionSecret)
	if err != nil {
		s.Logger.Errorf("get session store failed: %s", err)
	}
	// register middleware
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global", nil),
		session.Middleware(store),
	)

	// init and register controllers
	s.RegisterControllers(
		controllers.NewUserController(),
		controllers.NewPostsController(),
	)

	if err := s.Start(); err != nil {
		log.Fatalf("Server start failed: %s", err)
	}
}
