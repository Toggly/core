package mongo

import (
	"context"
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

func (s *mongoProjectStorage) filter(code string) bson.M {
	f := bson.M{
		"owner": s.owner,
	}
	if code != "" {
		f["code"] = code
	}
	return f
}

func (s *mongoProjectStorage) List() (list []*domain.Project, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	cur, err := s.collection().Find(ctxT, s.filter(""))
	if err != nil {
		s.log.Error().Err(err).Msg("DB error")
		return nil, err
	}
	defer cur.Close(ctxT)
	list = []*domain.Project{}
	for cur.Next(s.ctx) {
		item := &domain.Project{}
		err := cur.Decode(item)
		if err != nil {
			s.log.Error().Err(err).Msg("Can't decode project object")
			return nil, err
		}
		list = append(list, item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *mongoProjectStorage) Get(code string) (project *domain.Project, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	err = s.collection().FindOne(ctxT, s.filter(code)).Decode(&project)
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
	if _, err := s.collection().DeleteOne(ctxT, s.filter(code)); err != nil {
		s.log.Error().Err(err).Msg("Can't delete project")
		return err
	}
	return nil
}

func (s *mongoProjectStorage) ensureIndex(ctxT context.Context) error {
	idx := mongo.IndexModel{
		Keys: []bsonx.Elem{
			bsonx.Elem{Key: "owner", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "code", Value: bsonx.Int32(1)},
		},
		Options: []bsonx.Elem{bsonx.Elem{Key: "unique", Value: bsonx.Boolean(true)}},
	}
	if _, err := s.collection().Indexes().CreateOne(ctxT, idx); err != nil {
		s.log.Error().Err(err).Msg("Can't create index")
		return err
	}
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
	if _, err := s.collection().InsertOne(ctxT, project); err != nil {
		s.log.Error().Err(err).Msg("Can't insert project")
		if e, ok := err.(mongo.WriteErrors); ok && len(e) > 0 {
			switch e[0].Code {
			case 11000:
				return &storage.ErrUniqueIndex{Type: "project", Key: project.Key()}
			}
		}
		return err
	}
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
	if err := s.collection().FindOneAndReplace(ctxT, s.filter(project.Code), project).Decode(&domain.Project{}); err != nil {
		if err == mongo.ErrNoDocuments {
			return storage.ErrNotFound
		}
		s.log.Error().Err(err).Msg("Can't update project")
		return err
	}
	return nil
}
