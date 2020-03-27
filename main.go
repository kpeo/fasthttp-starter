package main

import (
	"fasthttp-starter/internal/setup"
	"github.com/fasthttp/router"
	"github.com/urfave/cli/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	stdlog "log"
	"os"
)

var (
	BuildTime string // Unixtime of app's build
	Commit    string // Recent Git commit hash
)

var cfg = setup.Settings{
	Name:      "My project",
	Version:   "0.0.1",
	File:      "config.yml",
	BuildTime: BuildTime,
	Commit:    Commit,
}

func main() {

	app := &cli.App{
		Name:      cfg.Name,
		Usage:     "My service",
		Version:   cfg.Version,
		Copyright: "(c) 2020 My Company",
		Authors: []*cli.Author{
			{
				Name:  "Developer",
				Email: "developer@developers.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",
				Value:       cfg.File,
				DefaultText: cfg.File,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		stdlog.Fatal("Error running app: ", zap.Error(err))
	}

	// Init Settings
	settings, _ := setup.NewSettings(&cfg)
	cfg.Config = settings

	// Init Logger
	log, _ := setup.NewLogger(setup.NewLoggerConfig(cfg.Config), &cfg)
	cfg.Logger = log

	// Init Database
	db, _ := setup.NewDatabase(settings)
	cfg.Db = db
	defer db.Close()

	// Init Router
	r := router.New()
	cfg.Router = r

	r.GET("/", indexHandler)

	serverPort := cfg.Config.GetString("server.port")
	if serverPort == "" {
		serverPort = "8080"
	}
	err = fasthttp.ListenAndServe(serverPort, r.Handler)
	if err != nil {
		stdlog.Fatal("Error while running server: ", zap.Error(err))
	}
}

func indexHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome to " + cfg.Name + " v" + cfg.Version + "(" + cfg.BuildTime + "," + cfg.Commit + ")")
}
