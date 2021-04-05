package db_commons

type BaseSvc struct {
	persistence      BaseRepository
	domainName       DomainName
	externalIdSetter ExternalIdSetter
	entityCreator    EntityCreator
}

func (b *BaseSvc) Init(repo BaseRepository, domainName DomainName, setter ExternalIdSetter, creator EntityCreator) {
	b.persistence = repo
	b.domainName = domainName
	b.externalIdSetter = setter
	b.entityCreator = creator
}

func (b *BaseSvc) FindById(id uint64) (error, Base) {
	return b.persistence.GetById(id, b.entityCreator)
}

func (b *BaseSvc) FindByExternalId(id string) (error, Base) {
	return b.persistence.GetByExternalId(id, b.entityCreator)
}

func (b *BaseSvc) MultiGetByExternalId(ids []string) (error, []Base) {
	return b.persistence.MultiGetByExternalId(ids, b.entityCreator)
}

func (b *BaseSvc) Create(base Base) (error, Base) {
	return b.persistence.Create(base, b.externalIdSetter)
}

func (b *BaseSvc) Update(id string, base Base) (error, Base) {
	return b.persistence.Update(id, base, b.entityCreator)
}

func (b *BaseSvc) GetPersistence() BaseRepository {
	return b.persistence
}
