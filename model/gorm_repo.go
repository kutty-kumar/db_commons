package db_commons

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type GORMRepository struct {
	db      *gorm.DB
	factory DomainFactory
}

func NewGORMRepository(db *gorm.DB, factory DomainFactory) *GORMRepository {
	return &GORMRepository{
		db:      db,
		factory: factory,
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
	if err := r.db.Where("external_id = ?", externalId).First(&entity).Error; err != nil {
		return err, nil
	}
	return nil, entity
}

func (r *GORMRepository) MultiGetByExternalId(externalIds [] string, creator func() []Base) (error, []Base) {
	entities := creator()
	if err := r.db.Where("external_id IN ?", externalIds).Find(&entities).Error; err != nil {
		return err, nil
	}
	return nil, entities
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
