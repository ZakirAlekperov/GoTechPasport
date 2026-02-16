package entity

// Address представляет адрес объекта недвижимости
// согласно правилам Государственного адресного реестра (ГАР)
type Address struct {
	// Субъект Российской Федерации
	Subject string `json:"subject"`

	// Административный район (округ)
	District string `json:"district,omitempty"`

	// Город (населенный пункт)
	City string `json:"city,omitempty"`

	// Район города
	CityDistrict string `json:"city_district,omitempty"`

	// Улица (переулок)
	Street string `json:"street,omitempty"`

	// Номер дома
	House string `json:"house"`

	// Строение (корпус)
	Building string `json:"building,omitempty"`

	// Квартира
	Apartment string `json:"apartment,omitempty"`

	// Комната
	Room string `json:"room,omitempty"`

	// Почтовый индекс
	PostalCode string `json:"postal_code,omitempty"`
}

// IsValid проверяет корректность заполнения адреса
func (a *Address) IsValid() error {
	if a.Subject == "" {
		return ValidationError{Field: "subject", Message: "субъект РФ обязателен"}
	}

	if a.House == "" {
		return ValidationError{Field: "house", Message: "номер дома обязателен"}
	}

	return nil
}

// FullAddress возвращает полный адрес в виде строки
func (a *Address) FullAddress() string {
	parts := []string{}

	if a.PostalCode != "" {
		parts = append(parts, a.PostalCode)
	}
	if a.Subject != "" {
		parts = append(parts, a.Subject)
	}
	if a.District != "" {
		parts = append(parts, a.District)
	}
	if a.City != "" {
		parts = append(parts, a.City)
	}
	if a.CityDistrict != "" {
		parts = append(parts, a.CityDistrict)
	}
	if a.Street != "" {
		parts = append(parts, a.Street)
	}
	if a.House != "" {
		parts = append(parts, "д. "+a.House)
	}
	if a.Building != "" {
		parts = append(parts, "корп. "+a.Building)
	}
	if a.Apartment != "" {
		parts = append(parts, "кв. "+a.Apartment)
	}
	if a.Room != "" {
		parts = append(parts, "ком. "+a.Room)
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ", "
		}
		result += part
	}

	return result
}
