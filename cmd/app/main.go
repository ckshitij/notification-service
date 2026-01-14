package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ckshitij/notify-srv/internal/config"
	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/repository"
	"github.com/ckshitij/notify-srv/internal/repository/mysql"
	"github.com/ckshitij/notify-srv/internal/server"
	"github.com/ckshitij/notify-srv/internal/template"
)

func main() {

	cfg, err := config.Load("./config/config.yml")
	if err != nil {
		panic(err)
	}

	log, err := logger.NewZapLogger(cfg.App.Env, cfg.App.LogLevel)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	database, err := mysql.New(cfg.MySQL)
	if err != nil {
		log.Fatal(ctx, "failed to connect to mysql", logger.Error(err))
	}
	defer database.Close()

	moduleRoutes := map[string]http.Handler{
		"/v1/templates": template.NewTemplateRoutes(repository.NewTemplateRepository(database.Conn())),
	}
	router := server.NewRouter(log, database, moduleRoutes)

	addr := ":" + fmt.Sprint(cfg.App.Port)
	srv := server.New(addr, log, router)

	go func() {
		if err := srv.Start(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatal(ctx, "server failed to start", logger.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error(ctx, "server shutdown failed", logger.Error(err))
	}

	log.Info(ctx, "server exited cleanly")
}
