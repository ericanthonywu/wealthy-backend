package main

import (
	"context"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/infrastructures/databases"
	"github.com/wealthy-app/wealthy-backend/routers"
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
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

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
		portBuilder strings.Builder
	)

	logrus.Info("starting application")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbConnection, err := databases.NewDBConnection()
	if err != nil {
		logrus.Error("can not connected database. reason : ", err.Error())
		return
	}

	route := gin.Default()
	route.Static("/images", "./images")
	route.Use(databases.DBInContext(dbConnection))
	route.NoRoute(routers.NoRoute)
	routers.RouterConfig(route)
	routers.API(&route.RouterGroup, dbConnection)

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
			logrus.Error(err.Error())
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
