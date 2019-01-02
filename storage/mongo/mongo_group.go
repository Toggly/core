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

type mongoGroupStorage struct {
	log     zerolog.Logger
	owner   string
	env     string
	project string
	ctx     context.Context
	db      *mongo.Database
}

func (s *mongoGroupStorage) collection() *mongo.Collection {
	return s.db.Collection("grp")
}

func (s *mongoGroupStorage) filter(code string) bson.M {
	f := bson.M{
		"owner":       s.owner,
		"project":     s.project,
		"environment": s.env,
	}
	if code != "" {
		f["code"] = code
	}
	return f
}

func (s *mongoGroupStorage) List() (list []*domain.Group, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	cur, err := s.collection().Find(ctxT, s.filter(""))
	if err != nil {
		s.log.Error().Err(err).Msg("DB error")
		return nil, err
	}
	defer cur.Close(ctxT)
	list = []*domain.Group{}
	for cur.Next(s.ctx) {
		item := &domain.Group{}
		err := cur.Decode(item)
		if err != nil {
			s.log.Error().Err(err).Msg("Can't decode group object")
			return nil, err
		}
		list = append(list, item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *mongoGroupStorage) Get(code string) (grp *domain.Group, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	err = s.collection().FindOne(ctxT, s.filter(code)).Decode(&grp)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, storage.ErrNotFound
		default:
			return nil, err
		}
	}
	return grp, nil
}

func (s *mongoGroupStorage) Delete(code string) error {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	if _, err := s.collection().DeleteOne(ctxT, s.filter(code)); err != nil {
		s.log.Error().Err(err).Msg("Can't delete group")
		return err
	}
	return nil
}

func (s *mongoGroupStorage) ensureIndex(ctxT context.Context) error {
	idx := mongo.IndexModel{
		Keys: []bsonx.Elem{
			bsonx.Elem{Key: "owner", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "project", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "environment", Value: bsonx.Int32(1)},
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

func (s *mongoGroupStorage) checkRelations(grp *domain.Group) error {
	if s.owner != grp.Owner {
		s.log.Error().Msgf("Wrong owner. Expected: %s, got: %s", s.owner, grp.Owner)
		return storage.ErrEntityRelationsBroken
	}
	if s.project != grp.Project {
		s.log.Error().Msgf("Wrong project. Expected: %s, got: %s", s.project, grp.Project)
		return storage.ErrEntityRelationsBroken
	}
	if s.env != grp.Environment {
		s.log.Error().Msgf("Wrong environment. Expected: %s, got: %s", s.env, grp.Environment)
		return storage.ErrEntityRelationsBroken
	}
	return nil
}

func (s *mongoGroupStorage) Save(grp *domain.Group) error {
	if err := s.checkRelations(grp); err != nil {
		return err
	}
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	if err := s.ensureIndex(ctxT); err != nil {
		return nil
	}
	if _, err := s.collection().InsertOne(ctxT, grp); err != nil {
		s.log.Error().Err(err).Msg("Can't insert group")
		if e, ok := err.(mongo.WriteErrors); ok && len(e) > 0 {
			switch e[0].Code {
			case 11000:
				return &storage.ErrUniqueIndex{Type: "group", Key: grp.Key()}
			}
		}
		return err
	}
	return nil
}

func (s *mongoGroupStorage) Update(grp *domain.Group) error {
	if err := s.checkRelations(grp); err != nil {
		return err
	}
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	if err := s.ensureIndex(ctxT); err != nil {
		return nil
	}
	if err := s.collection().FindOneAndReplace(ctxT, s.filter(grp.Code), grp).Decode(&domain.Group{}); err != nil {
		if err == mongo.ErrNoDocuments {
			return storage.ErrNotFound
		}
		s.log.Error().Err(err).Msg("Can't update group")
		return err
	}
	return nil
}
