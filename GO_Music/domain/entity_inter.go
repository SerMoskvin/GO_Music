package domain

import "time"

type Entity[ID comparable] interface {
	GetID() ID
	SetID(id ID)
	Validate() error
}

// ParseDMY преобразует "DD.MM.YYYY" в time.Time
func ParseDMY(dateStr string) time.Time {
	t, _ := time.Parse("02.01.2006", dateStr)
	return t
}

// ToDMY преобразует time.Time в "DD.MM.YYYY"
func ToDMY(t time.Time) string {
	return t.Format("02.01.2006")
}

func ParseTimeHM(s string) time.Time {
	const layout = "15:04"
	t, _ := time.Parse(layout, s)
	return t
}
