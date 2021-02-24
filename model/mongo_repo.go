package db_commons

import (
	"context"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	db      *mongo.Database
	factory DomainFactory
	ctx     context.Context
}

func NewMongoRepository(db *mongo.Database, factory DomainFactory) *MongoRepository {
	return &MongoRepository{
		db:      db,
		factory: factory,
		ctx:     context.TODO(),
	}
}

func (r *MongoRepository) GetById(id uint64, creator EntityCreator) (error, Base) {
	entity := creator()
	filter := primitive.E{Key: "_id", Value: id}
	if err := r.findOne(&entity, filter); err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *MongoRepository) GetByExternalId(externalId string, creator EntityCreator) (error, Base) {
	entity := creator()
	filter := primitive.E{Key: "externalId", Value: externalId}
	if err := r.findOne(&entity, filter); err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *MongoRepository) MultiGetByExternalId(externalIds []string, creator func() []Base) (error, []Base) {
	entities := creator()
	fmt.Printf("entities <%v>\n", entities)
	filter := bson.M{}
	filter["externalId"] = bson.M{
		"$in": externalIds,
	}
	if err := r.find(&entities, filter); err != nil {
		return err, nil
	}
	return nil, entities
}

func (r *MongoRepository) Create(base Base, externalIdSetter ExternalIdSetter) (error, Base) {
	err, externalId := r.generateExternalId(base)
	if err != nil {
		return err, nil
	}
	externalIdSetter(externalId, base)
	_, err = r.db.Collection(string(base.GetName())).InsertOne(r.ctx, &base)
	if err != nil {
		return err, nil
	}
	return nil, base
}

func (r *MongoRepository) Update(externalId string, updatedBase Base, creator EntityCreator) (error, Base) {
	err, entity := r.GetByExternalId(externalId, creator)
	if err != nil {
		return err, nil
	}
	filter := primitive.E{Key: "externalId", Value: externalId}
	entity.Merge(updatedBase)
	_, err = r.db.Collection(string(entity.GetName())).UpdateOne(r.ctx, filter, &entity)
	return nil, entity
}

func (r *MongoRepository) Search(params map[string]string, creator EntityCreator) (error, []Base) {
	return errors.New("not implemented"), nil
}

func (r *MongoRepository) findOne(entity *Base, filter interface{}) error {
	queryResult := r.db.Collection(string((*entity).GetName())).FindOne(r.ctx, filter)
	if queryResult.Err() != nil {
		return queryResult.Err()
	}
	return queryResult.Decode(&entity)
}

func (r *MongoRepository) find(entity *[]Base, filter interface{}) error {
	queryResult, err := r.db.Collection(string((*entity)[0].GetName())).Find(r.ctx, filter)
	if err != nil {
		return err
	}
	i := 0
	for queryResult.Next(r.ctx) {
		err := queryResult.Decode(&(*entity)[i])
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func (r *MongoRepository) generateExternalId(base Base) (error, string) {
	if base.GetExternalId() == "" {
		uid := uuid.NewV4()
		return nil, uid.String()
	}
	return nil, base.GetExternalId()
}
