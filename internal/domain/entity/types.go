package entity

import "time"

// ObjectType представляет тип объекта недвижимости
type ObjectType string

const (
	ObjectTypeResidentialHouse ObjectType = "residential_house" // Жилой дом
	ObjectTypeApartment        ObjectType = "apartment"         // Квартира
	ObjectTypeRoom             ObjectType = "room"              // Комната
	ObjectTypeNonResidential   ObjectType = "non_residential"   // Нежилое помещение
)

// PersonType представляет тип лица (физическое/юридическое)
type PersonType string

const (
	PersonTypeIndividual PersonType = "individual" // Физическое лицо
	PersonTypeLegal      PersonType = "legal"      // Юридическое лицо
)

// AuditEntry представляет запись в истории изменений
type AuditEntry struct {
	Timestamp   time.Time
	Action      string
	User        string
	Description string
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
