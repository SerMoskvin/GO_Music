package api

// Mapper базовый интерфейс для всех мапперов
type Mapper[CreateDTO any, UpdateDTO any, ResponseDTO any, Domain any] interface {
	ToDomain(dto *CreateDTO) *Domain
	UpdateDomain(domain *Domain, dto *UpdateDTO)
	ToResponse(domain *Domain) *ResponseDTO
	ToResponseList(domains []*Domain) []*ResponseDTO
}

// CRUDMapper интерфейс для CRUD операций
type CRUDMapper[CreateDTO any, UpdateDTO any, ResponseDTO any, Domain any] interface {
	Mapper[CreateDTO, UpdateDTO, ResponseDTO, Domain]
}
