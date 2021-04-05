package db_commons

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type GORMRepository struct {
	db               *gorm.DB
	creator          EntityCreator
	externalIdSetter ExternalIdSetter
}

func (r *GORMRepository) GetDb() *gorm.DB {
	return r.db
}

func NewGORMRepository(db *gorm.DB, entityCreator EntityCreator, setter ExternalIdSetter) *GORMRepository {
	return &GORMRepository{
		db:               db,
		creator:          entityCreator,
		externalIdSetter: setter,
	}
}

func (r *GORMRepository) GetById(id uint64) (error, Base) {
	entity := r.creator()
	if err := r.db.Where("id = ?", id).First(&entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) GetByExternalId(externalId string) (error, Base) {
	entity := r.creator()
	if err := r.db.Table(string(entity.GetName())).Where("external_id = ?", externalId).First(entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) populateRows(rows *sql.Rows) (error, []Base) {
	var models []Base
	for rows.Next() {
		entity := r.creator()
		entity, err := entity.FromSqlRow(rows)
		if err != nil {
			return err, nil
		}
		models = append(models, entity)
	}
	return nil, models
}

func (r *GORMRepository) MultiGetByExternalId(externalIds [] string) (error, []Base) {
	entity := r.creator()
	rows, err := r.db.Table(string(entity.GetName())).Where("external_id IN (?)", externalIds).Rows()
	if err != nil {
		return err, nil
	}
	return r.populateRows(rows)
}

func (r *GORMRepository) generateExternalId(base Base) (error, string) {
	if base.GetExternalId() == "" {
		uid := uuid.NewV4()
		return nil, uid.String()
	}
	return nil, base.GetExternalId()
}

func (r *GORMRepository) Create(base Base) (error, Base) {
	err, externalId := r.generateExternalId(base)
	if err != nil {
		return err, nil
	}
	r.externalIdSetter(externalId, base)
	if err := r.db.Create(base).Error; err != nil {
		return err, nil
	}
	return nil, base
}

func (r *GORMRepository) Update(externalId string, updatedBase Base) (error, Base) {
	err, entity := r.GetByExternalId(externalId)
	if err != nil {
		return err, nil
	}
	entity.Merge(updatedBase)
	if err := r.db.Table(string(entity.GetName())).Model(entity).Update(entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) Search(params map[string]string) (error, []Base) {
	return errors.New("not implemented"), nil
}
