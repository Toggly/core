package engine_test

import (
	"context"
	"os"

	"github.com/Toggly/core/storage"
	"github.com/Toggly/core/storage/mongo"
	driver "github.com/mongodb/mongo-go-driver/mongo"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger = log.Output(zerolog.ConsoleWriter{
	Out:     os.Stdout,
	NoColor: true,
}).Level(zerolog.DebugLevel)

func getDB() storage.DataStorage {
	ctx := context.Background()
	dataStorage, err := mongo.NewMongoDataStorage(ctx, "mongodb://localhost:27017", "toggly_api_test", logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Can't create storage")
	}
	err = dataStorage.Connect()
	if err != nil {
		logger.Fatal().Err(err).Msg("Can't connect")
	}
	return dataStorage
}

func dropDB() {
	client, err := driver.NewClient("mongodb://localhost:27017")
	if err != nil {
		logger.Fatal().Err(err).Msg("Can't connect to mongo")
	}
	ctx := context.Background()
	client.Connect(ctx)
	err = client.Database("toggly_api_test").Drop(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Can't drop db")
	}
}

func beforeTest() {
	dropDB()
}

func afterTest() {
	dropDB()
}
