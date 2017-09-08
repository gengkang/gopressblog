package main

import (
	"blog/config"
	"blog/controllers"
	"blog/middlewares"
	"blog/services"

	"github.com/fpay/gopress"
)

const (
	// ConfigFile config file path
	ConfigFile = "config/config.yaml"
	// TimeFormat time format str
)

func main() {
	// create server
	s := gopress.NewServer(gopress.ServerOptions{
		Port: 3000,
	})

	// opt
	opts := &config.Options{}
	opts.Database = &services.DBOptions{}
	opts.ScoreRule = &services.ScoreRule{}
	config.GetConfig(ConfigFile, opts)

	// services register
	dbs := services.NewDBService(opts.Database.DBType, opts.Database)
	vs := services.NewValidatorService()
	score := services.NewScoreService(opts.ScoreRule)
	s.RegisterServices(dbs, vs, score)

	// register middlewares
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global", gopress.NewLogger()),
	)

	// RouteGroups route groups
	needLoginMiddlewares := []gopress.MiddlewareFunc{middlewares.NewAuthMiddleware()}

	authGroup := s.App().Group("/blog", needLoginMiddlewares...)
	//init and register controllers
	s.RegisterControllers(
		controllers.NewIndexController(),
		controllers.NewUserController(),
		controllers.NewPostController(authGroup),
		controllers.NewCommentController(authGroup),
		controllers.NewAccountController(authGroup),
	)

	// static path
	s.App().Static("/assets", "assets")
	//
	s.Start()
}
