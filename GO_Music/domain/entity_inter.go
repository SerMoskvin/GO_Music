package domain

import "time"

type Entity[ID comparable] interface {
	GetID() ID
	SetID(id ID)
	Validate() error
}

// ParseDMY преобразует "DD.MM.YYYY" в time.Time
func ParseDMY(dateStr string) (time.Time, error) {
	return time.Parse("02.01.2006", dateStr)
}

// ToDMY преобразует time.Time в "DD.MM.YYYY"
func ToDMY(t time.Time) string {
	return t.Format("02.01.2006")
}
