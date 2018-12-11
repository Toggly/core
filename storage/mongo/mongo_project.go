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
	"github.com/rs/zerolog"
)

type mongoOwnerStorage struct {
	log   zerolog.Logger
	owner string
	ctx   context.Context
	db    *mongo.Database
}

func (s *mongoOwnerStorage) Projects() storage.ProjectStorage {
	return &mongoProjectStorage{
		log:   s.log,
		owner: s.owner,
		ctx:   s.ctx,
		db:    s.db,
	}
}

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
	cur, err := s.collection().Find(ctxT, nil)
	if err != nil {
		s.log.Error().Err(err).Msg("DB error")
		return nil, err
	}
	defer cur.Close(ctxT)
	for cur.Next(s.ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return nil, nil
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
	return nil
}

func (s *mongoProjectStorage) Save(project *domain.Project) error {
	ctxT, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()
	res, err := s.collection().InsertOne(ctxT, project)
	s.log.Debug().Str("id", fmt.Sprintf("%v", res.InsertedID)).Msg("Project inserted")
	return err
}

func (s *mongoProjectStorage) Update(project *domain.Project) error {
	return nil
}
