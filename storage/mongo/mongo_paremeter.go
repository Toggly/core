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

type mongoParameterStorage struct {
	log     zerolog.Logger
	owner   string
	project string
	env     string
	group   string
	ctx     context.Context
	db      *mongo.Database
}

func (s *mongoParameterStorage) collection() *mongo.Collection {
	return s.db.Collection("parameter")
}

func (s *mongoParameterStorage) filter(code string) bson.M {
	f := bson.M{
		"owner":       s.owner,
		"project":     s.project,
		"environment": s.env,
		"group":       s.group,
	}
	if code != "" {
		f["code"] = code
	}
	return f
}

func (s *mongoParameterStorage) list(codes ...string) (list []*domain.Parameter, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	filter := s.filter("")
	if len(codes) > 0 {
		in := bson.A{}
		for _, c := range codes {
			in = append(in, c)
		}
		filter["code"] = bson.M{"$in": in}
	}
	cur, err := s.collection().Find(ctxT, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("DB error")
		return nil, err
	}
	defer cur.Close(ctxT)
	list = []*domain.Parameter{}
	for cur.Next(s.ctx) {
		item := &domain.Parameter{}
		err := cur.Decode(item)
		if err != nil {
			s.log.Error().Err(err).Msg("Can't decode parameter object")
			return nil, err
		}
		list = append(list, item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *mongoParameterStorage) List() (list []*domain.Parameter, err error) {
	return s.list()
}

func (s *mongoParameterStorage) Get(code string) (parameter *domain.Parameter, err error) {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	err = s.collection().FindOne(ctxT, s.filter(code)).Decode(&parameter)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, storage.ErrNotFound
		default:
			return nil, err
		}
	}
	return parameter, nil
}

func (s *mongoParameterStorage) GetBatch(codes ...string) ([]*domain.Parameter, error) {
	return s.list(codes...)
}

func (s *mongoParameterStorage) Delete(code string) error {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	if _, err := s.collection().DeleteOne(ctxT, s.filter(code)); err != nil {
		s.log.Error().Err(err).Msg("Can't delete parameter")
		return err
	}
	return nil
}

func (s *mongoParameterStorage) ensureIndex(ctxT context.Context) error {
	idx := mongo.IndexModel{
		Keys: []bsonx.Elem{
			bsonx.Elem{Key: "owner", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "project", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "environment", Value: bsonx.Int32(1)},
			bsonx.Elem{Key: "group", Value: bsonx.Int32(1)},
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

func (s *mongoParameterStorage) checkRelations(parameter *domain.Parameter) error {
	if s.owner != parameter.Owner {
		s.log.Error().Msgf("Wrong owner. Expected: %s, got: %s", s.owner, parameter.Owner)
		return storage.ErrEntityRelationsBroken
	}
	if s.project != parameter.Project {
		s.log.Error().Msgf("Wrong project. Expected: %s, got: %s", s.project, parameter.Project)
		return storage.ErrEntityRelationsBroken
	}
	if s.env != parameter.Environment {
		s.log.Error().Msgf("Wrong environment. Expected: %s, got: %s", s.env, parameter.Environment)
		return storage.ErrEntityRelationsBroken
	}
	if s.group != parameter.Group {
		s.log.Error().Msgf("Wrong group. Expected: %s, got: %s", s.group, parameter.Group)
		return storage.ErrEntityRelationsBroken
	}
	return nil
}

func (s *mongoParameterStorage) Save(parameter *domain.Parameter) error {
	if err := s.checkRelations(parameter); err != nil {
		return err
	}
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	if err := s.ensureIndex(ctxT); err != nil {
		return nil
	}
	if _, err := s.collection().InsertOne(ctxT, parameter); err != nil {
		s.log.Error().Err(err).Msg("Can't insert group")
		if e, ok := err.(mongo.WriteErrors); ok && len(e) > 0 {
			switch e[0].Code {
			case 11000:
				return &storage.ErrUniqueIndex{Type: "parameter", Key: parameter.Key()}
			}
		}
		return err
	}
	return nil
}

func (s *mongoParameterStorage) Update(parameter *domain.Parameter) error {
	if err := s.checkRelations(parameter); err != nil {
		return err
	}
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	if err := s.ensureIndex(ctxT); err != nil {
		return nil
	}
	if err := s.collection().FindOneAndReplace(ctxT, s.filter(parameter.Code), parameter).Decode(&domain.Group{}); err != nil {
		if err == mongo.ErrNoDocuments {
			return storage.ErrNotFound
		}
		s.log.Error().Err(err).Msg("Can't update parameter")
		return err
	}
	return nil
}
