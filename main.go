package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"sts/db"
	"sts/handlers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

var DBSession *mgo.Session

func init() {
	decimal.DivisionPrecision = 2
}

func main() {
	loadConfig()
	DBSession = db.InitDB(viper.Sub("db"))
	defer DBSession.Close()

	app := echo.New()
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(setUpContext)

	app.GET("/take", appHandler(handlers.Take))
	app.GET("/fund", appHandler(handlers.Fund))
	app.GET("/announceTournament", appHandler(handlers.AnnounceTournament))
	app.GET("/joinTournament", appHandler(handlers.JoinTournament))
	app.POST("/resultTournament", appHandler(handlers.ResultTournament))
	app.GET("/balance", appHandler(handlers.Balance))
	app.GET("/reset", appHandler(handlers.Reset))

	appConfig := viper.Sub("app")
	host := appConfig.GetString("addrs") + ":" + appConfig.GetString("port")

	go func() {
		if err := app.Start(host); err != nil {
			app.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}

}

func loadConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic("Can't read config: " + err.Error())
	}
}

func appHandler(callback handlers.AppHandler) echo.HandlerFunc {
	return func(context echo.Context) error {
		appContext := context.(handlers.AppContext)
		appContext.Session = DBSession.Copy()
		defer func() { appContext.Session.Close() }()

		return callback(appContext)
	}
}

func setUpContext(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		appContext := handlers.AppContext{
			Context: context,
			Session: nil,
		}

		return handler(appContext)
	}
}
