package entity

// UtilityConnection представляет подключение к инженерным сетям
type UtilityConnection struct {
	// Площадь с централизованным подключением (кв.м)
	Centralized float64 `json:"centralized"`

	// Площадь с автономным подключением (кв.м)
	Autonomous float64 `json:"autonomous"`
}

// Utilities содержит информацию о благоустройстве объекта
type Utilities struct {
	// Водоснабжение
	Water UtilityConnection `json:"water"`

	// Канализация
	Sewerage UtilityConnection `json:"sewerage"`

	// Отопление
	Heating UtilityConnection `json:"heating"`

	// Горячее водоснабжение
	HotWater UtilityConnection `json:"hot_water"`

	// Газоснабжение
	Gas UtilityConnection `json:"gas"`

	// Электроснабжение
	Electricity UtilityConnection `json:"electricity"`

	// Другие элементы благоустройства
	Other string `json:"other,omitempty"`
}

// IsValid проверяет корректность данных о благоустройстве
func (u *Utilities) IsValid() error {
	// Проверка что площади неотрицательные
	if u.Water.Centralized < 0 || u.Water.Autonomous < 0 {
		return ValidationError{Field: "water", Message: "площадь не может быть отрицательной"}
	}

	if u.Sewerage.Centralized < 0 || u.Sewerage.Autonomous < 0 {
		return ValidationError{Field: "sewerage", Message: "площадь не может быть отрицательной"}
	}

	if u.Heating.Centralized < 0 || u.Heating.Autonomous < 0 {
		return ValidationError{Field: "heating", Message: "площадь не может быть отрицательной"}
	}

	if u.HotWater.Centralized < 0 || u.HotWater.Autonomous < 0 {
		return ValidationError{Field: "hot_water", Message: "площадь не может быть отрицательной"}
	}

	if u.Gas.Centralized < 0 || u.Gas.Autonomous < 0 {
		return ValidationError{Field: "gas", Message: "площадь не может быть отрицательной"}
	}

	if u.Electricity.Centralized < 0 || u.Electricity.Autonomous < 0 {
		return ValidationError{Field: "electricity", Message: "площадь не может быть отрицательной"}
	}

	return nil
}
