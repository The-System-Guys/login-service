package main

import (
	"time"

	routers "github.com/The-System-Guys/login-service/api/router"
	"github.com/The-System-Guys/login-service/components"
	"github.com/The-System-Guys/login-service/middleware"
	"github.com/The-System-Guys/login-service/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config, err := util.LoadConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Could not load environment variables")
	}

	db, err := gorm.Open(postgres.Open(config.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to Postgres:")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		AllowMethods:     []string{"GET", "PUSH", "POST"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.Recover())
	appCtx := components.NewAppContext(db, config.SecretKey)
	routers.AuthenticationRoute(router, appCtx)
	router.Run(":8080")
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}
