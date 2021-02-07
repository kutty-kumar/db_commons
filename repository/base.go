package db_commons

import "github.com/kutty-kumar/db_commons/model"

type BaseRepository interface {
	GetById(id uint64, creator db_commons.EntityCreator) (error,db_commons.Base)
	GetByExternalId(externalId string, creator db_commons.EntityCreator) (error, db_commons.Base)
	MultiGetByExternalId(externalIds[] string, creator func() []db_commons.Base) (error, []db_commons.Base)
	Create(base db_commons.Base) (error, db_commons.Base)
	Update(externalId string, updatedBase db_commons.Base, creator db_commons.EntityCreator) (error, db_commons.Base)
	Search(params map[string]string, creator db_commons.EntityCreator) (error, []db_commons.Base)
}
