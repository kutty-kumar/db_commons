package db_commons

type BaseRepository interface {
	GetById(id uint64, creator EntityCreator) (error, Base)
	GetByExternalId(externalId string, creator EntityCreator) (error, Base)
	MultiGetByExternalId(externalIds[] string, creator func() []Base) (error, []Base)
	Create(base Base, externalIdSetter ExternalIdSetter) (error, Base)
	Update(externalId string, updatedBase Base, creator EntityCreator) (error, Base)
	Search(params map[string]string, creator EntityCreator) (error, []Base)
}
