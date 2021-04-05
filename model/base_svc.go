package db_commons

type BaseSvc struct {
	persistence BaseRepository
}

func (b *BaseSvc) Init(repo BaseRepository) {
	b.persistence = repo
}

func (b *BaseSvc) FindById(id uint64) (error, Base) {
	return b.persistence.GetById(id)
}

func (b *BaseSvc) FindByExternalId(id string) (error, Base) {
	return b.persistence.GetByExternalId(id)
}

func (b *BaseSvc) MultiGetByExternalId(ids []string) (error, []Base) {
	return b.persistence.MultiGetByExternalId(ids)
}

func (b *BaseSvc) Create(base Base) (error, Base) {
	return b.persistence.Create(base)
}

func (b *BaseSvc) Update(id string, base Base) (error, Base) {
	return b.persistence.Update(id, base)
}

func (b *BaseSvc) GetPersistence() BaseRepository {
	return b.persistence
}
