package mongo

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/Toggly/core/storage"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// NewMongoDataStorage returns mongo storage implementation
func NewMongoDataStorage(ctx context.Context, url, dbName string, log zerolog.Logger) (storage.DataStorage, error) {
	client, err := mongo.NewClient(url)
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	go func() {
		<-ctx.Done()
		err := client.Disconnect(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Can't disconnect Mongo storage")
		}
		log.Info().Msg("Mongo storage disconnected")
	}()
	return &mongoStorage{
		ctx:    ctx,
		client: client,
		db:     db,
		log:    log,
	}, nil
}

type mongoStorage struct {
	ctx    context.Context
	client *mongo.Client
	log    zerolog.Logger
	db     *mongo.Database
}

func (s *mongoStorage) Connect() error {
	return s.client.Connect(s.ctx)
}

func (s *mongoStorage) ForOwner(owner string) storage.OwnerStorage {
	return &mongoOwnerStorage{
		log:   s.log,
		owner: owner,
		ctx:   s.ctx,
		db:    s.db,
	}
}
