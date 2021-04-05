package db_commons

import (
	"github.com/jinzhu/gorm"
)

type BaseRepository interface {
	GetById(id uint64, creator EntityCreator) (error, Base)
	GetByExternalId(externalId string, creator EntityCreator) (error, Base)
	MultiGetByExternalId(externalIds [] string, creator EntityCreator) (error, []Base)
	Create(base Base, setter ExternalIdSetter) (error, Base)
	Update(externalId string, updatedBase Base, creator EntityCreator) (error, Base)
	Search(params map[string]string, creator EntityCreator) (error, []Base)
	GetDb() *gorm.DB
}

type BaseDao struct {
	BaseRepository
}

func NewBaseGORMDao(db *gorm.DB) BaseDao {
	return BaseDao{
		NewGORMRepository(db),
	}
}
