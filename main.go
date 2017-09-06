package main

import (
	"log"

	"github.com/fpay/gopress"
	log2 "github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	_ "github.com/wuxinwei/goforum/conf"
	"github.com/wuxinwei/goforum/controllers"
	"github.com/wuxinwei/goforum/services"
)

func main() {
	// create server
	s := gopress.NewServer(gopress.ServerOptions{
		Port: 3000,
	})
	s.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	s.Logger.SetLevel(log2.WARN)

	// init and register services
	s.RegisterServices(
		services.NewDbService(),
	)

	// register middleware
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global", nil),
	)

	// init and register controllers
	s.RegisterControllers(
		controllers.NewUserController(),
	)

	if err := s.Start(); err != nil {
		log.Fatalf("Server start failed: %s", err)
	}
}
