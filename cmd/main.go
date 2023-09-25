package main

import (
	"context"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/docs"
	"github.com/semicolon-indonesia/wealthy-backend/infrastructures/databases"
	"github.com/semicolon-indonesia/wealthy-backend/infrastructures/instrumentations"
	"github.com/semicolon-indonesia/wealthy-backend/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// ----------------------------------------------------------------------------------------------------------------
	// READING ENV FILE
	// ----------------------------------------------------------------------------------------------------------------
	if err := godotenv.Load(constants.ENVPATH); err != nil {
		logrus.Error(err.Error())
		panic(err.Error())
	}

	appMode := os.Getenv("APP_MODE")

	if appMode == constants.PROD {
		gin.SetMode(gin.ReleaseMode)
	}
}

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	var (
		portBuilder   strings.Builder
		traceProvider *trace.TracerProvider
	)

	logrus.Info("starting application")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ----------------------------------------------------------------------------------------------------------------
	// DATABASE SETUP
	// ----------------------------------------------------------------------------------------------------------------
	dbConnection, err := databases.NewDBConnection()
	if err != nil {
		logrus.Error("can not connected database. reason : ", err.Error())
		return
	}

	// ----------------------------------------------------------------------------------------------------------------
	// ROUTER SETUP
	// ----------------------------------------------------------------------------------------------------------------
	route := gin.Default()
	route.Use(databases.DBInContext(dbConnection))
	route.NoRoute(routers.NoRoute)
	routers.RouterConfig(route)
	routers.API(&route.RouterGroup, dbConnection)

	// ----------------------------------------------------------------------------------------------------------------
	// INSTRUMENTATION SETUP
	// ----------------------------------------------------------------------------------------------------------------
	traceProvider = instrumentations.OpenTelemetryExporter(os.Getenv("APP_NAME"), os.Getenv("APP_MODE"))
	otel.SetTracerProvider(traceProvider)

	// ----------------------------------------------------------------------------------------------------------------
	// SWAGGER SETUP
	// ----------------------------------------------------------------------------------------------------------------
	docs.SwaggerInfo.Title = "Example API"
	docs.SwaggerInfo.Description = "This is a sample documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// ----------------------------------------------------------------------------------------------------------------
	// SERVER RUNNING
	// ----------------------------------------------------------------------------------------------------------------

	portBuilder.WriteString(":")
	portBuilder.WriteString(os.Getenv("APP_PORT"))

	srv := &http.Server{
		Addr:         portBuilder.String(),
		Handler:      route,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Error("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	logrus.Warn("shutting down gracefully, press Ctrl+C again to force 🔴")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown 🔴: ", err)
	}

	logrus.Warn("application closed 🔴")

}
