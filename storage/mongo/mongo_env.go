package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Toggly/core/domain"
	"github.com/Toggly/core/storage"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/rs/zerolog"
)

type mongoEnvironmentStorage struct {
	log     zerolog.Logger
	owner   string
	project string
	ctx     context.Context
	db      *mongo.Database
}

func (s *mongoEnvironmentStorage) collection() *mongo.Collection {
	return s.db.Collection("env")
}

func (s *mongoEnvironmentStorage) List() ([]*domain.Environment, error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	cur, err := s.collection().Find(ctxT, bson.M{"owner": s.owner, "project": s.project})
	if err != nil {
		s.log.Error().Err(err).Msg("DB error")
		return nil, err
	}
	defer cur.Close(ctxT)
	list := make([]*domain.Environment, 0)
	for cur.Next(s.ctx) {
		var item domain.Environment
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		list = append(list, &item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *mongoEnvironmentStorage) Get(code string) (env *domain.Environment, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	err = s.collection().FindOne(ctxT, bson.M{"owner": s.owner, "project": s.project, "code": code}).Decode(&env)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, storage.ErrNotFound
		default:
			return nil, err
		}
	}
	return env, nil
}

func (s *mongoEnvironmentStorage) Delete(code string) error {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	res, err := s.collection().DeleteOne(ctxT, bson.M{"owner": s.owner, "project": s.project, "code": code})
	s.log.Debug().Int64("count", res.DeletedCount).Msg("Environment deleted")
	return err
}

func (s *mongoEnvironmentStorage) ensureIndex(ctxT context.Context) error {
	idx := mongo.IndexModel{
		Keys: []bsonx.Elem{
			bsonx.Elem{Key: "owner", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "project", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "code", Value: bsonx.Int32(1)},
		},
		Options: []bsonx.Elem{bsonx.Elem{Key: "unique", Value: bsonx.Boolean(true)}},
	}
	name, err := s.collection().Indexes().CreateOne(ctxT, idx)
	if err != nil {
		s.log.Error().Err(err).Msg("Can't create index")
		return err
	}
	s.log.Debug().Str("name", name).Msg("Index created")
	return nil
}

func (s *mongoEnvironmentStorage) checkRelations(env *domain.Environment) error {
	if s.owner != env.Owner {
		s.log.Error().Msgf("Wrong owner. Expected: %s, got: %s", s.owner, env.Owner)
		return storage.ErrEntityRelationsBroken
	}
	if s.project != env.Project {
		s.log.Error().Msgf("Wrong project. Expected: %s, got: %s", s.project, env.Project)
		return storage.ErrEntityRelationsBroken
	}
	return nil
}

func (s *mongoEnvironmentStorage) Save(env *domain.Environment) error {
	if err := s.checkRelations(env); err != nil {
		return err
	}
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	if err := s.ensureIndex(ctxT); err != nil {
		return nil
	}

	res, err := s.collection().InsertOne(ctxT, env)
	// TODO: check unique index error
	s.log.Debug().Str("id", fmt.Sprintf("%v", res.InsertedID)).Msg("Environment inserted")
	return err
}

func (s *mongoEnvironmentStorage) Update(env *domain.Environment) error {
	if err := s.checkRelations(env); err != nil {
		return err
	}
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	if err := s.ensureIndex(ctxT); err != nil {
		return nil
	}

	res := s.collection().FindOneAndReplace(ctxT, bson.M{"owner": s.owner, "project": s.project, "code": env.Code}, env)
	// TODO: check not found error
	return res.Err()
}
