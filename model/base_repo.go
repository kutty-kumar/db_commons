package db_commons

import (
	"github.com/jinzhu/gorm"
)

type BaseRepository interface {
	GetById(id uint64) (error, Base)
	GetByExternalId(externalId string) (error, Base)
	MultiGetByExternalId(externalIds [] string) (error, []Base)
	Create(base Base) (error, Base)
	Update(externalId string, updatedBase Base) (error, Base)
	Search(params map[string]string) (error, []Base)
	GetDb() *gorm.DB
}

type BaseDao struct {
	BaseRepository
}

func NewBaseGORMDao(db *gorm.DB, creator EntityCreator, externalIdSetter ExternalIdSetter) BaseDao {
	return BaseDao{
		NewGORMRepository(db, creator, externalIdSetter),
	}
}
