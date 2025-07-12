package domain

type Entity[ID comparable] interface {
	GetID() ID
	SetID(id ID)
	Validate() error
}
