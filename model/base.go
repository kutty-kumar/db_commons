package db_commons

import "time"

const (
	active   = 0
	inactive = 1
)

type Status int

var statusMapping map[string]int
var statusReverseMapping map[int]string

func init(){
	statusMapping = make(map[string]int)
	statusReverseMapping = make(map[int]string)
	statusMapping["active"] = active
	statusMapping["inactive"] = inactive
	statusReverseMapping[active] = "active"
	statusReverseMapping[inactive] = "inactive"
}

func GetStatusInt(status string) int {
	return statusMapping[status]
}

func GetStatusStr(status int) string {
	return statusReverseMapping[status]
}

type EntityCreator func() Base

type DomainName string

type DomainFactory struct {
	entityMappings map[DomainName]EntityCreator
}

func (d *DomainFactory) RegisterMapping(domainName DomainName, creator EntityCreator){
	d.entityMappings[domainName] = creator
}

func (d *DomainFactory) GetMapping(domainName DomainName) EntityCreator {
	return d.entityMappings[domainName]
}

func NewDomainFactory() *DomainFactory {
	return  &DomainFactory{entityMappings: make(map[DomainName]EntityCreator)}
}

type Base interface {
	GetExternalId() string
	GetName() DomainName
	GetId() uint64
	GetStatus() Status
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() time.Time
	ToDto() interface{}
	FillProperties(dto interface{}) Base
	Merge(other Base) Base
}

type Attribute interface {
	GetKey() string
	GetValue() string
}

type AttributeWithLanguage interface {
	Attribute
	GetLanguage() string
}
