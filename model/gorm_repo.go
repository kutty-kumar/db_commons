package db_commons

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type GORMRepository struct {
	db               *gorm.DB
	factory          DomainFactory
	externalIdSetter func(externalId string, base Base) Base
}

func (r *GORMRepository) GetDb() *gorm.DB {
	return r.db
}

func NewGORMRepository(db *gorm.DB, factory DomainFactory, extIdSetter func(externalId string, base Base) Base) *GORMRepository {
	return &GORMRepository{
		db:               db,
		factory:          factory,
		externalIdSetter: extIdSetter,
	}
}

func (r *GORMRepository) GetById(id uint64, creator EntityCreator) (error, Base) {
	entity := creator()
	if err := r.db.Where("id = ?", id).First(&entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) GetByExternalId(externalId string, creator EntityCreator) (error, Base) {
	entity := creator()
	if err := r.db.Table(string(entity.GetName())).Where("external_id = ?", externalId).First(entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) populateRows(creator EntityCreator, rows *sql.Rows) (error, []Base) {
	var models []Base
	for rows.Next() {
		entity := creator()
		entity, err := entity.FromSqlRow(rows)
		if err != nil {
			return err, nil
		}
		models = append(models, entity)
	}
	return nil, models
}

func (r *GORMRepository) MultiGetByExternalId(externalIds [] string, creator EntityCreator) (error, []Base) {
	entity := creator()
	rows, err := r.db.Table(string(entity.GetName())).Where("external_id IN (?)", externalIds).Rows()
	if err != nil {
		return err, nil
	}
	return r.populateRows(creator, rows)
}

func (r *GORMRepository) generateExternalId(base Base) (error, string) {
	if base.GetExternalId() == "" {
		uid := uuid.NewV4()
		return nil, uid.String()
	}
	return nil, base.GetExternalId()
}

func (r *GORMRepository) Create(base Base, externalIdSetter ExternalIdSetter) (error, Base) {
	err, externalId := r.generateExternalId(base)
	if err != nil {
		return err, nil
	}
	externalIdSetter(externalId, base)
	if err := r.db.Create(base).Error; err != nil {
		return err, nil
	}
	return nil, base
}

func (r *GORMRepository) Update(externalId string, updatedBase Base, creator EntityCreator) (error, Base) {
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

func (r *GORMRepository) Search(params map[string]string, creator EntityCreator) (error, []Base) {
	return errors.New("not implemented"), nil
}
