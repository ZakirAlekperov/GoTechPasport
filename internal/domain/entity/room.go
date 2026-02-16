package entity

// Room представляет помещение в составе объекта (для экспликации)
type Room struct {
	// Литера объекта
	Litera string `json:"litera"`

	// Этаж
	Floor string `json:"floor"`

	// Номер помещения
	RoomNumber string `json:"room_number"`

	// Назначение помещения (жилая комната, кухня, ванная и т.п.)
	Purpose string `json:"purpose"`

	// Площадь помещения (кв.м)
	Area float64 `json:"area"`

	// Жилая площадь (кв.м)
	LivingArea float64 `json:"living_area"`

	// Вспомогательная площадь (кв.м)
	AuxiliaryArea float64 `json:"auxiliary_area"`

	// Высота помещения (м)
	Height float64 `json:"height"`

	// Самовольно переустроенная площадь (кв.м)
	UnauthorizedArea float64 `json:"unauthorized_area"`

	// Примечание
	Note string `json:"note,omitempty"`
}

// IsValid проверяет корректность данных помещения
func (r *Room) IsValid() error {
	if r.Litera == "" {
		return ValidationError{Field: "litera", Message: "литера обязательна"}
	}

	if r.Floor == "" {
		return ValidationError{Field: "floor", Message: "этаж обязателен"}
	}

	if r.RoomNumber == "" {
		return ValidationError{Field: "room_number", Message: "номер помещения обязателен"}
	}

	if r.Purpose == "" {
		return ValidationError{Field: "purpose", Message: "назначение помещения обязательно"}
	}

	if r.Area <= 0 {
		return ValidationError{Field: "area", Message: "площадь должна быть больше 0"}
	}

	if r.LivingArea < 0 {
		return ValidationError{Field: "living_area", Message: "жилая площадь не может быть отрицательной"}
	}

	if r.AuxiliaryArea < 0 {
		return ValidationError{Field: "auxiliary_area", Message: "вспомогательная площадь не может быть отрицательной"}
	}

	if r.Height < 0 {
		return ValidationError{Field: "height", Message: "высота не может быть отрицательной"}
	}

	return nil
}
