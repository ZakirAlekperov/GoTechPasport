package entity

// GeneralInfo содержит общие сведения об объекте недвижимости
type GeneralInfo struct {
	// Назначение объекта
	Purpose string `json:"purpose"`

	// Фактическое использование
	ActualUsage string `json:"actual_usage"`

	// Год постройки
	ConstructionYear int `json:"construction_year"`

	// Общая площадь (кв.м)
	TotalArea float64 `json:"total_area"`

	// Жилая площадь (кв.м)
	LivingArea float64 `json:"living_area"`

	// Количество этажей надземной части
	FloorsAboveGround int `json:"floors_above_ground"`

	// Количество этажей подземной части
	FloorsUnderground int `json:"floors_underground"`

	// Примечание
	Note string `json:"note,omitempty"`
}

// IsValid проверяет корректность общих сведений
func (g *GeneralInfo) IsValid() error {
	if g.Purpose == "" {
		return ValidationError{Field: "purpose", Message: "назначение объекта обязательно"}
	}

	if g.ConstructionYear < 1800 || g.ConstructionYear > 2100 {
		return ValidationError{Field: "construction_year", Message: "некорректный год постройки"}
	}

	if g.TotalArea <= 0 {
		return ValidationError{Field: "total_area", Message: "общая площадь должна быть больше 0"}
	}

	if g.LivingArea < 0 {
		return ValidationError{Field: "living_area", Message: "жилая площадь не может быть отрицательной"}
	}

	if g.LivingArea > g.TotalArea {
		return ValidationError{Field: "living_area", Message: "жилая площадь не может превышать общую"}
	}

	if g.FloorsAboveGround < 0 {
		return ValidationError{Field: "floors_above_ground", Message: "количество этажей не может быть отрицательным"}
	}

	if g.FloorsUnderground < 0 {
		return ValidationError{Field: "floors_underground", Message: "количество этажей не может быть отрицательным"}
	}

	return nil
}
