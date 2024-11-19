package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/EtincelleCoworking/availability-monitor/internal/admin"
	"github.com/EtincelleCoworking/availability-monitor/internal/monitor"
	"github.com/EtincelleCoworking/availability-monitor/internal/server"
	"github.com/EtincelleCoworking/availability-monitor/internal/web"
	"github.com/EtincelleCoworking/availability-monitor/pkg/db"
)

type Configuration struct {
	DbUrl             string `required:"true" split_words:"true"`
	AdminUsername     string `required:"true" split_words:"true"`
	AdminPasswordHash string `required:"true" split_words:"true"`
	Port              int    `default:"7000" split_words:"true"`
	// duration, in seconds, after which the sensor is unreachable
	UnreachableTimeout int `default:"900" split_words:"true"`
	// duration, in seconds, after which the sensor is maybe free
	MaybeFreeTimeout int `default:"180" split_words:"true"`
	// duration, in seconds, after which the sensor is free
	FreeTimeout int `default:"300" split_words:"true"`
}

var webApi = flag.Bool("web", false, "enables web API")
var adminApi = flag.Bool("admin", false, "enables admin API")
var monitorSvc = flag.Bool("monitor", false, "enables monitor service")
var ui = flag.Bool("dashboard", false, "enables UI")

func main() {
	flag.Parse()
	logger := log.Default()

	config := Configuration{}
	err := envconfig.Process("", &config)
	if err != nil {
		logger.Fatal(err)
	}

	if !*adminApi && !*webApi && !*monitorSvc && !*ui {
		*adminApi = true
		*webApi = true
		*monitorSvc = true
		*ui = true
	}

	connection, err := sqlx.Connect("postgres", config.DbUrl)
	if err != nil {
		logger.Fatal(err)
	}
	if err := prepareDb(connection); err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	httpServer := server.NewServer()

	sensorManager := db.NewDbSensorManager(connection, logger)

	if *webApi {
		reporter := db.NewDbReporter(connection, logger)
		webApi := web.NewWebApi(reporter, logger)
		webApi.MapHandlers(httpServer.Echo, "api")
		logger.Print("Web API enabled")
	}
	if *adminApi {
		credentials := admin.NewBasicCredentials(config.AdminUsername, config.AdminPasswordHash)
		admApi := admin.NewAdminApi(sensorManager, credentials, logger)

		admApi.MapHandlers(httpServer.Echo, "admin")
		logger.Print("Admin API enabled")
	}
	if *monitorSvc {
		monitorService := monitor.NewMonitorService(
			sensorManager,
			logger,
			time.Duration(config.UnreachableTimeout)*time.Second,
			time.Duration(config.MaybeFreeTimeout)*time.Second,
			time.Duration(config.FreeTimeout)*time.Second,
		)
		monitorService.MapHandlers(httpServer.Echo)
		monitorService.Run(ctx)
		logger.Print("Monitor Service enabled")
	}
	if *ui {
		httpServer.Static("/", "./pkg/views/static")
		httpServer.GET("/", func(c echo.Context) error {
			return c.Render(http.StatusOK, "home.html", nil)
		})
		logger.Print("UI enabled")
	}

	go httpServer.Start(fmt.Sprintf(":%d", config.Port))

	// listen for shutdown signal
	shutdownCtx, clearShutdownCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer clearShutdownCtx()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-shutdownCtx.Done()
	logger.Println("shutting down")
	cancel()

	// timeout if server shutdown takes too long
	shudownTimeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelTimeout()
	if err := httpServer.Shutdown(shudownTimeoutCtx); err != nil {
		logger.Fatal(err)
	}

}

func prepareDb(db *sqlx.DB) error {
	var schema = `
CREATE TABLE IF NOT EXISTS sensors (
	mac text NOT NULL UNIQUE,
	location text NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT(NOW()	),
    last_seen_at timestamp NOT NULL,
	last_motion_at timestamp NULL
);
`
	_, err := db.Exec(schema)
	return err
}
