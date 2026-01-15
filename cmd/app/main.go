package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ckshitij/notify-srv/internal/config"
	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/pkg/notification"
	notfysql "github.com/ckshitij/notify-srv/internal/pkg/notification/mysql"
	tmplsql "github.com/ckshitij/notify-srv/internal/pkg/template/mysql"

	"github.com/ckshitij/notify-srv/internal/mysql"
	"github.com/ckshitij/notify-srv/internal/pkg/renderer"
	"github.com/ckshitij/notify-srv/internal/pkg/senders/email"
	"github.com/ckshitij/notify-srv/internal/pkg/senders/inapp"
	"github.com/ckshitij/notify-srv/internal/pkg/senders/slack"
	"github.com/ckshitij/notify-srv/internal/pkg/template"
	"github.com/ckshitij/notify-srv/internal/server"
	"github.com/ckshitij/notify-srv/internal/shared"
)

func processModules(ctx context.Context, database *mysql.DB, cfg *config.Config, log logger.Logger) map[string]http.Handler {
	senders := map[shared.Channel]notification.Sender{
		shared.ChannelEmail: email.New(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.From,
			cfg.SMTP.User,
			cfg.SMTP.Pass,
		),
		shared.ChannelSlack: slack.New(cfg.Slack.WebhookURL),
		shared.ChannelInApp: inapp.New(database.Conn()),
	}

	renderer := renderer.NewGoTemplateRenderer()
	templateRepo := tmplsql.NewTemplateRepository(database, log)
	templateService := template.NewTemplateService(templateRepo, renderer)

	notificationRepo := notfysql.NewNotificationRepository(database, log)
	notificationSrv := notification.NewNotificationService(notificationRepo, renderer, senders, templateRepo, log)
	scheduler := notification.NewSchedular(notificationSrv, notificationRepo, log, 5*time.Second, 50)

	go scheduler.Run(ctx)

	return map[string]http.Handler{
		"/v1/templates":     template.NewTemplateRoutes(templateService),
		"/v1/notifications": notification.NewNotificationRoutes(notificationSrv),
	}
}

func main() {

	configPath := flag.String("config", "./config/config.yml", "pass the config file path")
	swaggerFilePath := flag.String("swagger", "./api/openapi.yaml", "path to openapi spec file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
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

	schedularCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handlers := processModules(schedularCtx, database, cfg, log)
	router := server.NewRouter(log, database, *swaggerFilePath, handlers)

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
