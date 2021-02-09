package db_commons

import (
	"errors"
	"github.com/jinzhu/gorm"
	"kutty-kumar/db_commons/model"
)

type GORMRepository struct {
	db      *gorm.DB
	factory db_commons.DomainFactory
}

func NewGORMRepository(db *gorm.DB, factory db_commons.DomainFactory) *GORMRepository {
	return &GORMRepository{
		db:      db,
		factory: factory,
	}
}

func (r *GORMRepository) GetById(id uint64, creator db_commons.EntityCreator) (error, db_commons.Base) {
	entity := creator()
	if err := r.db.Where("id = (?)", id).Find(entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) GetByExternalId(externalId string, creator db_commons.EntityCreator) (error, db_commons.Base) {
	entity := creator()
	if err := r.db.Where("external_id = (?)", externalId).Find(entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) MultiGetByExternalId(externalIds [] string, creator func() []db_commons.Base) (error, []db_commons.Base) {
	entities := creator()
	if err := r.db.Where("external_id IN (?)", externalIds).Find(entities).Error; err != nil {
		return err, nil
	}
	return nil, entities
}

func (r *GORMRepository) Create(base db_commons.Base) (error, db_commons.Base) {
	if err := r.db.Create(base).Error; err != nil {
		return err, nil
	}
	return nil, base
}

func (r *GORMRepository) Update(externalId string, updatedBase db_commons.Base, creator db_commons.EntityCreator) (error, db_commons.Base) {
	err, entity := r.GetByExternalId(externalId, creator)
	if err != nil {
		return err, nil
	}
	entity.Merge(updatedBase)
	if err := r.db.Table(string(entity.GetName())).Model(entity).Update(entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) Search(params map[string]string, creator db_commons.EntityCreator) (error, []db_commons.Base) {
	return errors.New("not implemented"), nil
}
