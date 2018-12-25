// Toggly Server
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Toggly/core/api/engine"
	"github.com/Toggly/core/rest"
	"github.com/Toggly/core/storage/mongo"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const logo = `
_______  _____   ______  ______       __   __      _______ _______  ______ _    _ _______  ______
   |    |     | |  ____ |  ____ |       \_/        |______ |______ |_____/  \  /  |______ |_____/
   |    |_____| |_____| |_____| |_____   |         ______| |______ |    \_   \/   |______ |    \_                                                                                                
`

var version = "development"

type options struct {
	Version       bool   `short:"v" long:"version" description:"Show version"`
	Port          int    `short:"p" long:"port" default:"8080" env:"TOGGLY_SRV_PORT" description:"Port"`
	NoLogo        bool   `long:"no-logo" description:"Do not display logo"`
	BasePath      string `long:"base-path" default:"/api" env:"TOGGLY_SRV_BASE_PATH" description:"Rest API base path"`
	StoreMongoURL string `long:"store-mongo-url" default:"mongodb://localhost:27017" env:"TOGGLY_SRV_STORE_MONGO_URL" description:"Mongo connection url"`
	StoreMongoDB  string `long:"store-mongo-db" default:"toggly" env:"TOGGLY_SRV_STORE_MONGO_DB" description:"Mongo database name"`
	CacheType     string `long:"cache-type" env:"TOGGLY_SRV_CACHE_TYPE" choice:"memory" choice:"redis" default:"memory" description:"Cache type"`
	CacheRedisURL string `long:"cache-redis-url" env:"TOGGLY_SRV_CACHE_REDIS_URL" description:"Redis connection url"`
	Debug         bool   `long:"debug" env:"TOGGLY_SRV_DEBUG" description:"Debug mode"`
}

func main() {

	var opts options

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	if !opts.NoLogo {
		fmt.Print(logo)
		fmt.Printf("\n  ver. %s\n\n", version)
	}

	logLevel := zerolog.InfoLevel

	if opts.Debug {
		logLevel = zerolog.DebugLevel
		log.Logger = log.With().Caller().Logger()
	}
	zerolog.SetGlobalLevel(logLevel)

	ctx, cancel := context.WithCancel(context.Background())

	log := log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warn().Msg("Interrupt signal")
		cancel()
	}()

	dataStorage, err := mongo.NewMongoDataStorage(ctx, opts.StoreMongoURL, opts.StoreMongoDB, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't create mongo client")
	}

	err = dataStorage.Connect()
	if err != nil {
		log.Fatal().Err(err).Msg("Can't open storage connection")
	}

	server := &rest.Server{
		Version:  version,
		API:      &engine.APIEngine{Storage: dataStorage, Log: log},
		Log:      log,
		LogLevel: logLevel,
	}

	log.Info().Msg("API server started")
	server.Run(ctx, opts.Port, opts.BasePath)
	log.Warn().Msg("Application terminated")
}
