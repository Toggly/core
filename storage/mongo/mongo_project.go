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

type mongoProjectStorage struct {
	log   zerolog.Logger
	owner string
	ctx   context.Context
	db    *mongo.Database
}

func (s *mongoProjectStorage) collection() *mongo.Collection {
	return s.db.Collection("project")
}

func (s *mongoProjectStorage) List() ([]*domain.Project, error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	cur, err := s.collection().Find(ctxT, bson.M{"owner": s.owner})
	if err != nil {
		s.log.Error().Err(err).Msg("DB error")
		return nil, err
	}
	defer cur.Close(ctxT)
	list := make([]*domain.Project, 0)
	for cur.Next(s.ctx) {
		var item domain.Project
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

func (s *mongoProjectStorage) Get(code string) (project *domain.Project, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	err = s.collection().FindOne(ctxT, bson.M{"owner": s.owner, "code": code}).Decode(&project)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, storage.ErrNotFound
		default:
			return nil, err
		}
	}
	return project, nil
}

func (s *mongoProjectStorage) Delete(code string) error {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	res, err := s.collection().DeleteOne(ctxT, bson.M{"owner": s.owner, "code": code})
	s.log.Debug().Int64("count", res.DeletedCount).Msg("Project deleted")
	return err
}

func (s *mongoProjectStorage) ensureIndex(ctxT context.Context) error {
	idx := mongo.IndexModel{
		Keys: []bsonx.Elem{
			bsonx.Elem{Key: "owner", Value: bsonx.Int32(1)},
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

func (s *mongoProjectStorage) checkRelations(project *domain.Project) error {
	if s.owner != project.Owner {
		s.log.Error().Msgf("Wrong owner. Expected: %s, got: %s", s.owner, project.Owner)
		return storage.ErrEntityRelationsBroken
	}
	return nil
}

func (s *mongoProjectStorage) Save(project *domain.Project) error {
	if err := s.checkRelations(project); err != nil {
		return err
	}

	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	if err := s.ensureIndex(ctxT); err != nil {
		return err
	}

	res, err := s.collection().InsertOne(ctxT, project)
	if err != nil {
		s.log.Error().Err(err).Msg("Can't insert project")
		if e, ok := err.(mongo.WriteErrors); ok && len(e) > 0 {
			switch e[0].Code {
			case 11000:
				return &storage.ErrUniqueIndex{Type: "project", Key: project.Key()}
			}
		}
		return err
	}

	s.log.Debug().Str("id", fmt.Sprintf("%v %v", res, err)).Msg("Project inserted")
	return nil
}

func (s *mongoProjectStorage) Update(project *domain.Project) error {
	if err := s.checkRelations(project); err != nil {
		return err
	}

	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	if err := s.ensureIndex(ctxT); err != nil {
		return nil
	}

	res := s.collection().FindOneAndReplace(ctxT, bson.M{"owner": s.owner, "code": project.Code}, project)
	var proj domain.Project
	err := res.Decode(&proj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return storage.ErrNotFound
		}
		s.log.Error().Err(err).Msg("Can't decode project")
		return err
	}
	return res.Err()
}
