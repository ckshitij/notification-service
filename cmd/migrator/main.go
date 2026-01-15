package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ckshitij/notify-srv/internal/config"
	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/mysql"
)

func main() {
	migrationsDir := "./migrations"

	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

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

	var migrations []string

	err = filepath.Walk(migrationsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".sql") {
			return nil
		}

		switch *direction {
		case "up":
			if strings.HasSuffix(info.Name(), ".up.sql") {
				migrations = append(migrations, path)
			}
		case "down":
			if strings.HasSuffix(info.Name(), ".down.sql") {
				migrations = append(migrations, path)
			}
		default:
			return fmt.Errorf("invalid direction: %s", *direction)
		}

		return nil
	})

	if err != nil {
		log.Fatal(ctx, "failed to read migrations", logger.Error(err))
	}

	// Ensure deterministic order
	sort.Strings(migrations)

	if len(migrations) == 0 {
		log.Info(ctx, "no migrations to run")
		return
	}

	log.Info(ctx, "starting migrations", logger.String("direction", *direction))

	for _, migration := range migrations {
		log.Info(ctx, "executing migration", logger.String("file", migration))

		sqlBytes, err := os.ReadFile(migration)
		if err != nil {
			log.Fatal(ctx, "failed to read migration file", logger.Error(err))
		}

		if _, err := database.Conn().ExecContext(ctx, string(sqlBytes)); err != nil {
			log.Fatal(ctx, "migration failed", logger.Error(err), logger.String("file", migration))
		}
	}

	log.Info(ctx, "migrations completed successfully")
}
