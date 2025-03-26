package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/EtincelleCoworking/availability-monitor/internal/admin"
	"github.com/EtincelleCoworking/availability-monitor/internal/monitor"
	"github.com/EtincelleCoworking/availability-monitor/internal/server"
	"github.com/EtincelleCoworking/availability-monitor/internal/web"
	"github.com/EtincelleCoworking/availability-monitor/pkg/db"
)

type Configuration struct {
	DbUrl              url.URL       `env:"DB_URL,required"`
	AdminUsername      string        `env:"ADMIN_USERNAME,required"`
	AdminPasswordHash  string        `env:"ADMIN_PASSWORD_HASH,required"`
	Port               int           `env:"PORT" envDefault:"7000"`
	MaybeFreeTimeout   time.Duration `env:"MAYBE_FREE_TIMEOUT" envDefault:"3m"`
	FreeTimeout        time.Duration `env:"FREE_TIMEOUT" envDefault:"5m"`
	UnreachableTimeout time.Duration `env:"UNREACHABLE_TIMEOUT" envDefault:"15m"`
}

var webApi = flag.Bool("web", false, "enables web API")
var adminApi = flag.Bool("admin", false, "enables admin API")
var monitorSvc = flag.Bool("monitor", false, "enables monitor service")
var ui = flag.Bool("dashboard", false, "enables UI")

//go:embed static/*
var staticFiles embed.FS

func main() {
	flag.Parse()
	logger := log.Default()

	config := Configuration{}
	err := env.Parse(&config)
	if err != nil {
		logger.Fatal(err)
	}

	if !*adminApi && !*webApi && !*monitorSvc && !*ui {
		*adminApi = true
		*webApi = true
		*monitorSvc = true
		*ui = true
	}

	connection, err := sqlx.Connect("postgres", config.DbUrl.String())
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
		credentials := admin.NewBasicCredentials(config.AdminUsername, config.AdminPasswordHash, logger)
		admApi := admin.NewAdminApi(sensorManager, credentials, logger)

		admApi.MapHandlers(httpServer.Echo, "admin")
		logger.Print("Admin API enabled")
	}
	if *monitorSvc {
		monitorService := monitor.NewMonitorService(
			sensorManager,
			logger,
			config.UnreachableTimeout,
			config.MaybeFreeTimeout,
			config.FreeTimeout,
		)
		monitorService.MapHandlers(httpServer.Echo)
		monitorService.Run(ctx)
		logger.Print("Monitor Service enabled")
	}
	if *ui {
		static := middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "static",
			Index:      "home.html",
			Filesystem: http.FS(staticFiles),
		})
		httpServer.Use(static)

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
