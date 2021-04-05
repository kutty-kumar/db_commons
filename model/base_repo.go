package db_commons

import (
	"github.com/jinzhu/gorm"
)

type BaseRepository interface {
	GetById(id uint64, creator EntityCreator) (error, Base)
	GetByExternalId(externalId string, creator EntityCreator) (error, Base)
	MultiGetByExternalId(externalIds [] string, creator EntityCreator) (error, []Base)
	Create(base Base) (error, Base)
	Update(externalId string, updatedBase Base, creator EntityCreator) (error, Base)
	Search(params map[string]string, creator EntityCreator) (error, []Base)
	GetDb() *gorm.DB
	GetMapping(domainName DomainName) EntityCreator
	GetExternalIdSetter() ExternalIdSetter
}

type BaseDao struct {
	BaseRepository
}

func NewBaseDao(db *gorm.DB, factory DomainFactory,
	externalIdSetter func(externalId string, base Base) Base) BaseDao {
	persistence := NewGORMRepository(db, factory, externalIdSetter)
	return BaseDao{
		persistence,
	}
}
