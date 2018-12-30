package mongo_test

import (
	"context"
	"os"

	driver "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/Toggly/core/storage"
	"github.com/Toggly/core/storage/mongo"
)

var logger = log.Output(zerolog.ConsoleWriter{
	Out:     os.Stdout,
	NoColor: true,
}).Level(zerolog.DebugLevel)

func getDB() storage.DataStorage {
	ctx := context.Background()
	dataStorage, err := mongo.NewDataStorage(ctx, "mongodb://localhost:27017", "toggly_storage_test", logger)
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
	err = client.Database("toggly_storage_test").Drop(ctx)
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
